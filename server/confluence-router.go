package main

import (
	"log"
	"net/http"
	"time"

	gonnect "github.com/craftamap/atlas-gonnect"
	"github.com/craftamap/atlas-gonnect/middleware"
	"github.com/gorilla/mux"
)

type ConfluenceRouter struct {
	Addon *gonnect.Addon
}

func (router *ConfluenceRouter) setupRoutes(s *mux.Router) {
	s.Handle("/macro", middleware.NewAuthenticationMiddleware(router.Addon, false)(
		http.HandlerFunc(router.macro),
	))
}

func (router *ConfluenceRouter) macro(w http.ResponseWriter, r *http.Request) {
	clientID := r.Context().Value("clientKey").(string)
	orgIDs, err := GetSettingsRepository().GetOrganizationIDsByValue(SettingConfluenceClientID.Name, clientID)
	if (err != nil) || (len(orgIDs) > 1) {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	if len(orgIDs) == 0 {
		SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/confluence/"+clientID)
		return
	}
	// TODO Auth user
	userID := router.getUserEmail(r)
	user, err := GetUserRepository().GetByAtlassianID(userID)
	if err != nil {
		user = &User{
			Email:          userID,
			AtlassianID:    NullString(userID),
			OrganizationID: orgIDs[0],
			OrgAdmin:       false,
			SuperAdmin:     false,
		}
		GetUserRepository().Create(user)
	}
	authState := &AuthState{
		AuthProviderID: "00000000-0000-0000-0000-000000000000",
		Expiry:         time.Now().Add(time.Minute * 5),
		AuthStateType:  AuthAtlassian,
		Payload:        userID,
	}
	if err := GetAuthStateRepository().Create(authState); err != nil {
		SendInternalServerError(w)
		return
	}
	SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/login/success/"+authState.ID)
	/*
		log.Println("Confluence plugin request for user " + r.Context().Value("userAccountId").(string))
		log.Println("Client key " + r.Context().Value("clientKey").(string))
		SendTemporaryRedirect(w, GetConfig().FrontendURL+"ui/")
	*/
}

func (router *ConfluenceRouter) getUserEmail(r *http.Request) string {
	userAccountID := r.Context().Value("userAccountId").(string)
	clientKey := r.Context().Value("clientKey").(string)
	return userAccountID + "@" + clientKey
	/*
		httpClient, err := hostrequest.FromRequest(router.Addon, r)
		if err != nil {
			log.Println("1")
			return "", err
		}
		log.Println(httpClient)
		//request, err := http.NewRequest("GET", "https://team-1609087764900.atlassian.net/wiki/rest/api/user/current", http.NoBody)
		request, err := http.NewRequest("GET", "https://team-1609087764900.atlassian.net/wiki/rest/api/user/email?accountId="+r.Context().Value("userAccountId").(string), http.NoBody)
		if err != nil {
			log.Println("2")
			return "", err
		}
		request, err = httpClient.AsAddon(request)
		//request, err = httpClient.AsUser(request, r.Context().Value("userAccountId").(string))
		if err != nil {
			log.Println("2.5")
			return "", err
		}
		log.Println(request)
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Println("3")
			return "", err
		}
		log.Println(response)
		return "", nil
	*/
}
