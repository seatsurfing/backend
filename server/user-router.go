package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserRouter struct {
}

type CreateUserRequest struct {
	Email          string `json:"email" validate:"required"`
	OrgAdmin       bool   `json:"admin"`
	SuperAdmin     bool   `json:"superAdmin"`
	AuthProviderID string `json:"authProviderId"`
	Password       string `json:"password"`
}

type GetUserResponse struct {
	ID              string                  `json:"id"`
	OrganizationID  string                  `json:"organizationId"`
	Organization    GetOrganizationResponse `json:"organization"`
	RequirePassword bool                    `json:"requirePassword"`
	CreateUserRequest
}

type GetUserCountResponse struct {
	Count int `json:"count"`
}

type SetPasswordRequest struct {
	Password string `json:"password"`
}

func (router *UserRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/count", router.getCount).Methods("GET")
	s.HandleFunc("/me", router.getSelf).Methods("GET")
	s.HandleFunc("/{id}", router.getOne).Methods("GET")
	s.HandleFunc("/{id}/password", router.setPassword).Methods("PUT")
	s.HandleFunc("/{id}", router.update).Methods("PUT")
	s.HandleFunc("/{id}", router.delete).Methods("DELETE")
	s.HandleFunc("/", router.create).Methods("POST")
	s.HandleFunc("/", router.getAll).Methods("GET")
}

func (router *UserRouter) getCount(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	num, _ := GetUserRepository().GetCount(user.OrganizationID)
	m := &GetUserCountResponse{
		Count: num,
	}
	SendJSON(w, m)
}

func (router *UserRouter) setPassword(w http.ResponseWriter, r *http.Request) {
	var m SetPasswordRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	e, err := GetUserRepository().GetOne(vars["id"])
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAdminOrg(user, e.OrganizationID) && (user.ID != e.ID) {
		SendForbidden(w)
		return
	}
	e.HashedPassword = NullString(GetUserRepository().GetHashedPassword(m.Password))
	if err := GetUserRepository().Update(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *UserRouter) getSelf(w http.ResponseWriter, r *http.Request) {
	e := GetRequestUser(r)
	if e == nil {
		SendNotFound(w)
		return
	}
	org, err := GetOrganizationRepository().GetOne(e.OrganizationID)
	if err != nil {
		SendInternalServerError(w)
		return
	}
	res := router.copyToRestModel(e, false)
	res.Organization = GetOrganizationResponse{
		ID: org.ID,
		CreateOrganizationRequest: CreateOrganizationRequest{
			Name: org.Name,
		},
	}
	SendJSON(w, res)
}

func (router *UserRouter) getOne(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	vars := mux.Vars(r)
	e, err := GetUserRepository().GetOne(vars["id"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	if e.OrganizationID != user.OrganizationID {
		SendForbidden(w)
		return
	}
	res := router.copyToRestModel(e, true)
	SendJSON(w, res)
}

func (router *UserRouter) getAll(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	list, err := GetUserRepository().GetAll(user.OrganizationID, 1000, 0)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetUserResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e, true)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *UserRouter) update(w http.ResponseWriter, r *http.Request) {
	var m CreateUserRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	e, err := GetUserRepository().GetOne(vars["id"])
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAdminOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	eNew := router.copyFromRestModel(&m)
	eNew.ID = e.ID
	eNew.OrganizationID = e.OrganizationID
	eNew.SuperAdmin = e.SuperAdmin
	eNew.HashedPassword = e.HashedPassword
	org, err := GetOrganizationRepository().GetOne(e.OrganizationID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	if !GetOrganizationRepository().isValidEmailForOrg(user.Email, org) {
		SendBadRequest(w)
		return
	}
	if err := GetUserRepository().Update(eNew); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *UserRouter) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetUserRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAdminOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	if err := GetUserRepository().Delete(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *UserRouter) create(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	var m CreateUserRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	e := router.copyFromRestModel(&m)
	e.OrganizationID = user.OrganizationID
	org, err := GetOrganizationRepository().GetOne(e.OrganizationID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	if !router.canCreateUser(org) {
		SendPaymentRequired(w)
		return
	}
	if !GetOrganizationRepository().isValidEmailForOrg(e.Email, org) {
		SendBadRequest(w)
		return
	}
	if err := GetUserRepository().Create(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendCreated(w, e.ID)
}

func (router *UserRouter) canCreateUser(org *Organization) bool {
	maxUsers, _ := GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	curUsers, _ := GetUserRepository().GetCount(org.ID)
	return curUsers < maxUsers
}

func (router *UserRouter) copyFromRestModel(m *CreateUserRequest) *User {
	e := &User{}
	e.Email = m.Email
	e.OrgAdmin = m.OrgAdmin
	e.SuperAdmin = false
	if m.Password != "" {
		e.HashedPassword = NullString(GetUserRepository().GetHashedPassword(m.Password))
		e.AuthProviderID = NullString("")
	} else {
		e.AuthProviderID = NullString(m.AuthProviderID)
	}
	return e
}

func (router *UserRouter) copyToRestModel(e *User, admin bool) *GetUserResponse {
	m := &GetUserResponse{}
	m.ID = e.ID
	m.OrganizationID = e.OrganizationID
	m.Email = e.Email
	m.OrgAdmin = e.OrgAdmin
	m.SuperAdmin = e.SuperAdmin
	m.RequirePassword = (e.HashedPassword != "")
	if admin {
		m.AuthProviderID = string(e.AuthProviderID)
	}
	return m
}
