package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SpaceRouter struct {
}

type CreateSpaceRequest struct {
	Name     string `json:"name" validate:"required"`
	X        uint   `json:"x"`
	Y        uint   `json:"y"`
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`
	Rotation uint   `json:"rotation"`
}

type GetSpaceResponse struct {
	ID         string              `json:"id"`
	Available  bool                `json:"available"`
	LocationID string              `json:"locationId"`
	Location   GetLocationResponse `json:"location"`
	CreateSpaceRequest
}

type GetSpaceAvailabilityRequest struct {
	Enter time.Time `json:"enter" validate:"required"`
	Leave time.Time `json:"leave" validate:"required"`
}

func (router *SpaceRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/availability", router.getAvailability).Methods("POST")
	s.HandleFunc("/{id}", router.getOne).Methods("GET")
	s.HandleFunc("/{id}", router.update).Methods("PUT")
	s.HandleFunc("/{id}", router.delete).Methods("DELETE")
	s.HandleFunc("/", router.create).Methods("POST")
	s.HandleFunc("/", router.getAll).Methods("GET")
}

func (router *SpaceRouter) getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetSpaceRepository().GetOne(vars["id"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	location, err := GetLocationRepository().GetOne(e.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAccessOrg(user, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	res := router.copyToRestModel(e)
	SendJSON(w, res)
}

func (router *SpaceRouter) getAvailability(w http.ResponseWriter, r *http.Request) {
	var m GetSpaceAvailabilityRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	location, err := GetLocationRepository().GetOne(vars["locationId"])
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAccessOrg(user, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	list, err := GetSpaceRepository().GetAllInTime(location.ID, m.Enter, m.Leave)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetSpaceResponse{}
	for _, e := range list {
		m := &GetSpaceResponse{}
		m.ID = e.ID
		m.LocationID = e.LocationID
		m.Name = e.Name
		m.X = e.X
		m.Y = e.Y
		m.Width = e.Width
		m.Height = e.Height
		m.Rotation = e.Rotation
		m.Available = e.Available
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *SpaceRouter) getAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location, err := GetLocationRepository().GetOne(vars["locationId"])
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAccessOrg(user, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	list, err := GetSpaceRepository().GetAll(location.ID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetSpaceResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *SpaceRouter) update(w http.ResponseWriter, r *http.Request) {
	var m CreateSpaceRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	e := router.copyFromRestModel(&m)
	e.ID = vars["id"]
	e.LocationID = vars["locationId"]
	location, err := GetLocationRepository().GetOne(e.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAdminOrg(user, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if err := GetSpaceRepository().Update(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *SpaceRouter) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetSpaceRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	location, err := GetLocationRepository().GetOne(e.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAdminOrg(user, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if err := GetSpaceRepository().Delete(e); err != nil {
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *SpaceRouter) create(w http.ResponseWriter, r *http.Request) {
	var m CreateSpaceRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	e := router.copyFromRestModel(&m)
	e.LocationID = vars["locationId"]
	location, err := GetLocationRepository().GetOne(e.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	user := GetRequestUser(r)
	if !CanAdminOrg(user, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if err := GetSpaceRepository().Create(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendCreated(w, e.ID)
}

func (router *SpaceRouter) copyFromRestModel(m *CreateSpaceRequest) *Space {
	e := &Space{}
	e.Name = m.Name
	e.X = m.X
	e.Y = m.Y
	e.Width = m.Width
	e.Height = m.Height
	e.Rotation = m.Rotation
	return e
}

func (router *SpaceRouter) copyToRestModel(e *Space) *GetSpaceResponse {
	m := &GetSpaceResponse{}
	m.ID = e.ID
	m.LocationID = e.LocationID
	m.Name = e.Name
	m.X = e.X
	m.Y = e.Y
	m.Width = e.Width
	m.Height = e.Height
	m.Rotation = e.Rotation
	return m
}
