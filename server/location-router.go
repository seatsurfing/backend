package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type LocationRouter struct {
}

type CreateLocationRequest struct {
	Name                  string `json:"name" validate:"required"`
	Description           string `json:"description"`
	MaxConcurrentBookings uint   `json:"maxConcurrentBookings"`
	Timezone              string `json:"timezone"`
}

type GetLocationResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	MapWidth       uint   `json:"mapWidth"`
	MapHeight      uint   `json:"mapHeight"`
	MapMimeType    string `json:"mapMimeType"`
	CreateLocationRequest
}

type GetMapResponse struct {
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

func (router *LocationRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/loadsampledata", router.loadSampleData).Methods("POST")
	s.HandleFunc("/{id}/map", router.getMap).Methods("GET")
	s.HandleFunc("/{id}/map", router.setMap).Methods("POST")
	s.HandleFunc("/{id}", router.getOne).Methods("GET")
	s.HandleFunc("/{id}", router.update).Methods("PUT")
	s.HandleFunc("/{id}", router.delete).Methods("DELETE")
	s.HandleFunc("/", router.create).Methods("POST")
	s.HandleFunc("/", router.getAll).Methods("GET")
}

func (router *LocationRouter) getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetLocationRepository().GetOne(vars["id"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAccessOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	res := router.copyToRestModel(e)
	SendJSON(w, res)
}

func (router *LocationRouter) getAll(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	list, err := GetLocationRepository().GetAll(user.OrganizationID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetLocationResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *LocationRouter) update(w http.ResponseWriter, r *http.Request) {
	var m CreateLocationRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	e, err := GetLocationRepository().GetOne(vars["id"])
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanSpaceAdminOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	if m.Timezone != "" {
		if !isValidTimeZone(m.Timezone) {
			SendBadRequest(w)
			return
		}
	}
	eNew := router.copyFromRestModel(&m)
	eNew.ID = e.ID
	eNew.OrganizationID = e.OrganizationID
	if err := GetLocationRepository().Update(eNew); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *LocationRouter) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetLocationRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	user := GetRequestUser(r)
	if !CanSpaceAdminOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	if err := GetLocationRepository().Delete(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *LocationRouter) create(w http.ResponseWriter, r *http.Request) {
	var m CreateLocationRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	e := router.copyFromRestModel(&m)
	e.OrganizationID = user.OrganizationID
	if !CanSpaceAdminOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	if m.Timezone != "" {
		if !isValidTimeZone(m.Timezone) {
			SendBadRequest(w)
			return
		}
	}
	if err := GetLocationRepository().Create(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendCreated(w, e.ID)
}

func (router *LocationRouter) getMap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetLocationRepository().GetOne(vars["id"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAccessOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	locationMap, err := GetLocationRepository().GetMap(e)
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	res := &GetMapResponse{
		Width:    locationMap.Width,
		Height:   locationMap.Height,
		MimeType: locationMap.MimeType,
		Data:     base64.StdEncoding.EncodeToString(locationMap.Data),
	}
	SendJSON(w, res)
}

func (router *LocationRouter) setMap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetLocationRepository().GetOne(vars["id"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	user := GetRequestUser(r)
	if !CanSpaceAdminOrg(user, e.OrganizationID) {
		SendForbidden(w)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		SendBadRequest(w)
		return
	}
	image, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		SendBadRequest(w)
		return
	}
	locationMap := &LocationMap{
		Width:    uint(image.Width),
		Height:   uint(image.Height),
		MimeType: format,
		Data:     data,
	}
	if err := GetLocationRepository().SetMap(e, locationMap); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *LocationRouter) loadSampleData(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	org, err := GetOrganizationRepository().GetOne(user.OrganizationID)
	if err != nil {
		SendInternalServerError(w)
		return
	}
	GetOrganizationRepository().createSampleData(org)
}

func (router *LocationRouter) copyFromRestModel(m *CreateLocationRequest) *Location {
	e := &Location{}
	e.Name = m.Name
	e.Description = m.Description
	e.MaxConcurrentBookings = m.MaxConcurrentBookings
	e.Timezone = m.Timezone
	return e
}

func (router *LocationRouter) copyToRestModel(e *Location) *GetLocationResponse {
	m := &GetLocationResponse{}
	m.ID = e.ID
	m.OrganizationID = e.OrganizationID
	m.Name = e.Name
	m.MapMimeType = e.MapMimeType
	m.MapWidth = e.MapWidth
	m.MapHeight = e.MapHeight
	m.Description = e.Description
	m.MaxConcurrentBookings = e.MaxConcurrentBookings
	m.Timezone = e.Timezone
	return m
}
