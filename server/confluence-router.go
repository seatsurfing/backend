package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ConfluenceServerClaims struct {
	UserName string `json:"user"`
	UserKey  string `json:"key"`
	jwt.StandardClaims
}

type ConfluenceRouter struct {
}

func (router *ConfluenceRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/{orgID}/{jwt}", router.serverLogin).Methods("GET")
}

func (router *ConfluenceRouter) serverLogin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org, err := GetOrganizationRepository().GetOne(vars["orgID"])
	if err != nil || org == nil {
		SendNotFound(w)
		return
	}
	sharedSecret, err := GetSettingsRepository().Get(org.ID, SettingConfluenceServerSharedSecret.Name)
	if err != nil || sharedSecret == "" {
		SendBadRequest(w)
		return
	}
	claims := &ConfluenceServerClaims{}
	token, err := jwt.ParseWithClaims(vars["jwt"], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sharedSecret), nil
	})
	if err != nil {
		log.Println("JWT header verification failed: parsing JWT failed with: " + err.Error())
		SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/failed")
		return
	}
	if !token.Valid {
		log.Println("JWT header verification failed: invalid JWT")
		SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/failed")
		return
	}
	allowAnonymous, _ := GetSettingsRepository().GetBool(org.ID, SettingConfluenceAnonymous.Name)
	userID := router.getUserEmailServer(org, claims, allowAnonymous)
	if userID == "" {
		SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/confluence/anonymous")
		return
	}
	_, err = GetUserRepository().GetByAtlassianID(userID)
	if err != nil {
		if !GetUserRepository().canCreateUser(org) {
			SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/failed")
			return
		}
		user := &User{
			Email:          userID,
			AtlassianID:    NullString(userID),
			OrganizationID: org.ID,
			OrgAdmin:       false,
			SuperAdmin:     false,
		}
		GetUserRepository().Create(user)
	}
	authState := &AuthState{
		AuthProviderID: GetSettingsRepository().getNullUUID(),
		Expiry:         time.Now().Add(time.Minute * 5),
		AuthStateType:  AuthAtlassian,
		Payload:        userID,
	}
	if err := GetAuthStateRepository().Create(authState); err != nil {
		SendInternalServerError(w)
		return
	}
	SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/success/"+authState.ID)
}

func (router *ConfluenceRouter) getUserEmailServer(org *Organization, claims *ConfluenceServerClaims, allowAnonymous bool) string {
	userAccountID := "confluence-" + claims.UserName
	if claims.UserName == "" {
		if !allowAnonymous {
			return ""
		}
		userAccountID = "confluence-anonymous-" + uuid.New().String()
	}
	domains, err := GetOrganizationRepository().GetDomains(org)
	if err != nil {
		return ""
	}
	domain := ""
	for _, curDomain := range domains {
		if curDomain.Active {
			domain = curDomain.DomainName
		}
	}
	return userAccountID + "@" + domain
}
