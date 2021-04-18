package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type JWTResponse struct {
	JWT string `json:"jwt"`
}

type Claims struct {
	Email    string `json:"email"`
	UserID   string `json:"userID"`
	OrgAdmin bool   `json:"admin"`
	jwt.StandardClaims
}

type AuthPreflightRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type InitPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CompletePasswordResetRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

type AuthPreflightResponse struct {
	Organization    *GetOrganizationResponse         `json:"organization"`
	AuthProviders   []*GetAuthProviderPublicResponse `json:"authProviders"`
	RequirePassword bool                             `json:"requirePassword"`
}

type AuthPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthRouter struct {
}

func (router *AuthRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/verify/{id}", router.verify).Methods("GET")
	s.HandleFunc("/{id}/login/{type}", router.login).Methods("GET")
	s.HandleFunc("/{id}/callback", router.callback).Methods("GET")
	s.HandleFunc("/preflight", router.preflight).Methods("POST")
	s.HandleFunc("/login", router.loginPassword).Methods("POST")
	s.HandleFunc("/initpwreset", router.initPasswordReset).Methods("POST")
	s.HandleFunc("/pwreset/{id}", router.completePasswordReset).Methods("POST")
}

func (router *AuthRouter) initPasswordReset(w http.ResponseWriter, r *http.Request) {
	var m InitPasswordResetRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	user, err := GetUserRepository().GetByEmail(m.Email)
	if user == nil || err != nil {
		SendNotFound(w)
		return
	}
	if user.HashedPassword == "" {
		SendNotFound(w)
		return
	}
	org, err := GetOrganizationRepository().GetOne(user.OrganizationID)
	if org == nil || err != nil {
		SendNotFound(w)
		return
	}
	authState := &AuthState{
		AuthProviderID: GetSettingsRepository().getNullUUID(),
		Expiry:         time.Now().Add(time.Hour * 1),
		AuthStateType:  AuthResetPasswordRequest,
		Payload:        user.ID,
	}
	GetAuthStateRepository().Create(authState)
	router.SendPasswordResetEmail(user, authState.ID, org)
	SendUpdated(w)
}

func (router *AuthRouter) completePasswordReset(w http.ResponseWriter, r *http.Request) {
	var m CompletePasswordResetRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	authState, err := GetAuthStateRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	if authState.AuthStateType != AuthResetPasswordRequest {
		SendNotFound(w)
		return
	}
	user, err := GetUserRepository().GetOne(authState.Payload)
	if user == nil || err != nil {
		SendNotFound(w)
		return
	}
	if user.HashedPassword == "" {
		SendNotFound(w)
		return
	}
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword(m.Password))
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *AuthRouter) preflight(w http.ResponseWriter, r *http.Request) {
	var m AuthPreflightRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	res := router.getPreflightResponse(&m)
	if res == nil {
		SendNotFound(w)
		return
	}
	user, err := GetUserRepository().GetByEmail(m.Email)
	if err != nil {
		log.Println(err)
		SendJSON(w, res)
		return
	}
	res.RequirePassword = (user.HashedPassword != "")
	SendJSON(w, res)
}

func (router *AuthRouter) loginPassword(w http.ResponseWriter, r *http.Request) {
	var m AuthPasswordRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	user, err := GetUserRepository().GetByEmail(m.Email)
	if err != nil {
		SendNotFound(w)
		return
	}
	if user.HashedPassword == "" {
		SendNotFound(w)
		return
	}
	if !GetUserRepository().CheckPassword(string(user.HashedPassword), m.Password) {
		SendNotFound(w)
		return
	}
	claims := router.createClaims(user)
	jwt := router.createJWT(claims)
	res := &JWTResponse{
		JWT: jwt,
	}
	SendJSON(w, res)
}

func (router *AuthRouter) handleAtlassianVerify(authState *AuthState, w http.ResponseWriter) {
	user, err := GetUserRepository().GetByAtlassianID(authState.Payload)
	if err != nil {
		SendNotFound(w)
		return
	}
	GetAuthStateRepository().Delete(authState)
	claims := router.createClaims(user)
	jwt := router.createJWT(claims)
	res := &JWTResponse{
		JWT: jwt,
	}
	SendJSON(w, res)
}

func (router *AuthRouter) verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authState, err := GetAuthStateRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	if authState.AuthStateType == AuthAtlassian {
		router.handleAtlassianVerify(authState, w)
		return
	}
	if authState.AuthStateType != AuthResponseCache {
		SendNotFound(w)
		return
	}
	provider, err := GetAuthProviderRepository().GetOne(authState.AuthProviderID)
	if err != nil {
		SendNotFound(w)
		return
	}
	user, err := GetUserRepository().GetByEmail(authState.Payload)
	// TODO Change email to auth server ID???
	if err != nil {
		org, err := GetOrganizationRepository().GetOne(provider.OrganizationID)
		if err != nil {
			SendInternalServerError(w)
			return
		}
		if !GetUserRepository().canCreateUser(org) {
			SendPaymentRequired(w)
			return
		}
		user := &User{
			Email:          authState.Payload,
			OrganizationID: org.ID,
			OrgAdmin:       false,
			SuperAdmin:     false,
		}
		GetUserRepository().Create(user)
	}
	if user.OrganizationID != provider.OrganizationID {
		SendBadRequest(w)
		return
	}
	GetAuthStateRepository().Delete(authState)
	claims := router.createClaims(user)
	jwt := router.createJWT(claims)
	res := &JWTResponse{
		JWT: jwt,
	}
	SendJSON(w, res)
}

func (router *AuthRouter) login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loginType := vars["type"]
	if loginType != "web" && loginType != "app" && loginType != "ui" {
		SendBadRequest(w)
		return
	}
	provider, err := GetAuthProviderRepository().GetOne(vars["id"])
	if err != nil {
		SendTemporaryRedirect(w, router.getRedirectFailedUrl(loginType))
		return
	}
	config := router.getConfig(provider)
	authState := &AuthState{
		AuthProviderID: provider.ID,
		Expiry:         time.Now().Add(time.Minute * 5),
		AuthStateType:  AuthRequestState,
		Payload:        loginType,
	}
	if err := GetAuthStateRepository().Create(authState); err != nil {
		SendTemporaryRedirect(w, router.getRedirectFailedUrl(loginType))
		return
	}
	url := config.AuthCodeURL(authState.ID)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (router *AuthRouter) callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider, err := GetAuthProviderRepository().GetOne(vars["id"])
	if err != nil {
		SendTemporaryRedirect(w, router.getRedirectFailedUrl("ui"))
		return
	}
	claims, loginType, err := router.getUserInfo(provider, r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		SendTemporaryRedirect(w, router.getRedirectFailedUrl(loginType))
		return
	}
	if !router.isValidEmailForOrg(provider, claims.Email) {
		SendTemporaryRedirect(w, router.getRedirectFailedUrl(loginType))
		return
	}
	allowAnyUser, _ := GetSettingsRepository().Get(provider.OrganizationID, SettingAllowAnyUser.Name)
	if allowAnyUser != "1" {
		_, err := GetUserRepository().GetByEmail(claims.Email)
		if err != nil {
			SendTemporaryRedirect(w, router.getRedirectFailedUrl(loginType))
			return
		}
	}
	authState := &AuthState{
		AuthProviderID: provider.ID,
		Expiry:         time.Now().Add(time.Minute * 5),
		AuthStateType:  AuthResponseCache,
		Payload:        claims.Email,
	}
	if err := GetAuthStateRepository().Create(authState); err != nil {
		log.Println(err)
		SendTemporaryRedirect(w, router.getRedirectFailedUrl(loginType))
		return
	}
	SendTemporaryRedirect(w, router.getRedirectSuccessUrl(loginType, authState))
}

func (router *AuthRouter) getRedirectSuccessUrl(loginType string, authState *AuthState) string {
	if loginType == "app" {
		return GetConfig().AppURL + "login/success/" + authState.ID
	} else if loginType == "ui" {
		return GetConfig().FrontendURL + "ui/login/success/" + authState.ID
	} else {
		return GetConfig().FrontendURL + "admin/login/success/" + authState.ID
	}
}

func (router *AuthRouter) getRedirectFailedUrl(loginType string) string {
	if loginType == "app" {
		return GetConfig().AppURL + "login/failed"
	} else if loginType == "ui" {
		return GetConfig().FrontendURL + "ui/login/failed"
	} else {
		return GetConfig().FrontendURL + "admin/login/failed"
	}
}

func (router *AuthRouter) isValidEmailForOrg(provider *AuthProvider, email string) bool {
	org, err := GetOrganizationRepository().GetOne(provider.OrganizationID)
	if err != nil {
		return false
	}
	return GetOrganizationRepository().isValidEmailForOrg(email, org)
}

func (router *AuthRouter) getUserInfo(provider *AuthProvider, state string, code string) (*Claims, string, error) {
	// Verify state string
	authState, err := GetAuthStateRepository().GetOne(state)
	if err != nil {
		return nil, "", fmt.Errorf("state not found for id %s", state)
	}
	if authState.AuthProviderID != provider.ID {
		return nil, "", fmt.Errorf("auth providers don't match")
	}
	defer GetAuthStateRepository().Delete(authState)
	// Exchange authorization code for an access token
	config := router.getConfig(provider)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}
	// Get user info from resource server
	client := &http.Client{}
	req, err := http.NewRequest("GET", provider.UserInfoURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed creating http request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	response, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}
	// Extract email address from JSON response
	var result map[string]interface{}
	json.Unmarshal([]byte(contents), &result)
	claims := &Claims{
		Email: result[provider.UserInfoEmailField].(string),
	}
	return claims, authState.Payload, nil
}

func (router *AuthRouter) SendPasswordResetEmail(user *User, ID string, org *Organization) error {
	email := user.Email
	if strings.Contains(email, GetConfig().SignupAdmin+"@") && strings.Contains(email, GetConfig().SignupDomain) {
		email = org.ContactEmail
	}
	vars := map[string]string{
		"recipientName":  user.Email,
		"recipientEmail": email,
		"confirmID":      ID,
	}
	return sendEmail(user.Email, "info@seatsurfing.de", EmailTemplateResetpassword, org.Language, vars)
}

func (router *AuthRouter) getConfig(provider *AuthProvider) *oauth2.Config {
	config := &oauth2.Config{
		RedirectURL:  GetConfig().PublicURL + "auth/" + provider.ID + "/callback",
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Scopes:       strings.Split(provider.Scopes, ","),
		Endpoint: oauth2.Endpoint{
			AuthURL:   provider.AuthURL,
			TokenURL:  provider.TokenURL,
			AuthStyle: oauth2.AuthStyle(provider.AuthStyle),
		},
	}
	return config
}

func (router *AuthRouter) createClaims(user *User) *Claims {
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		OrgAdmin: user.OrgAdmin,
	}
	return claims
}

func (router *AuthRouter) createJWT(claims *Claims) string {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(60 * 24 * 14 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwtString, err := accessToken.SignedString([]byte(GetConfig().JwtSigningKey))
	if err != nil {
		return ""
	}
	return jwtString
}

func (router *AuthRouter) getOrgForEmail(email string) *Organization {
	mailParts := strings.Split(email, "@")
	if len(mailParts) != 2 {
		return nil
	}
	domain := strings.ToLower(mailParts[1])
	org, err := GetOrganizationRepository().GetOneByDomain(domain)
	if err != nil {
		log.Println(err)
		return nil
	}
	return org
}

func (router *AuthRouter) getPreflightResponse(req *AuthPreflightRequest) *AuthPreflightResponse {
	org := router.getOrgForEmail(req.Email)
	if org == nil {
		return nil
	}
	list, err := GetAuthProviderRepository().GetAll(org.ID)
	if err != nil {
		return nil
	}
	res := &AuthPreflightResponse{
		Organization: &GetOrganizationResponse{
			ID: org.ID,
			CreateOrganizationRequest: CreateOrganizationRequest{
				Name: org.Name,
			},
		},
		RequirePassword: false,
		AuthProviders:   []*GetAuthProviderPublicResponse{},
	}
	for _, e := range list {
		m := &GetAuthProviderPublicResponse{}
		m.ID = e.ID
		m.Name = e.Name
		res.AuthProviders = append(res.AuthProviders, m)
	}
	return res
}
