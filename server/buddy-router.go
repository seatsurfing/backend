package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type BuddyRouter struct {
}

type BuddyRequest struct {
	BuddyID           string `json:"buddyId" validate:"required"`
	BuddyEmail        string `json:"buddyEmail"`
	BuddyFirstBooking string `json:"buddyFirstBooking"`
}

type CreateBuddyRequest struct {
	BuddyRequest
}

type GetBuddyResponse struct {
	ID string `json:"id"`
	CreateBuddyRequest
}

func (router *BuddyRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/{id}", router.delete).Methods("DELETE")
	s.HandleFunc("/", router.create).Methods("POST")
	s.HandleFunc("/", router.getAll).Methods("GET")
}

func (router *BuddyRouter) getAll(w http.ResponseWriter, r *http.Request) {
	list, err := GetBuddyRepository().GetAllByOwner(GetRequestUserID(r))
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetBuddyResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *BuddyRouter) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetBuddyRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	if e.OwnerID != GetRequestUserID(r) {
		SendForbidden(w)
	}
	if err := GetBuddyRepository().Delete(e); err != nil {
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *BuddyRouter) create(w http.ResponseWriter, r *http.Request) {
	var m CreateBuddyRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}

	buddyUser, err := GetUserRepository().GetOne(m.BuddyID)
	if err != nil {
		SendBadRequest(w)
		return
	}

	e := &Buddy{}
	e.BuddyID = buddyUser.ID
	e.OwnerID = GetRequestUserID(r)
	if err := GetBuddyRepository().Create(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendCreated(w, e.ID)
}

func (router *BuddyRouter) copyToRestModel(e *BuddyDetails) *GetBuddyResponse {
	m := &GetBuddyResponse{}
	m.ID = e.ID
	m.BuddyID = e.BuddyID
	m.BuddyEmail = e.BuddyEmail
	// Assuming GetOne returns a pointer to BookingDetails
	bookingDetails, _ := GetBookingRepository().GetFirstUpcomingBookingByUserID(e.BuddyID)
	// Use * to dereference the pointer
	actualBookingDetails := *bookingDetails

	// Assuming bookingDetails.Enter is of type time.Time
	// Use .Format to convert it to a string
	enterString := actualBookingDetails.Enter.Format("02-01-2006")
	space := actualBookingDetails.Space.Location.Name
	m.BuddyFirstBooking = enterString + " at the " + space
	return m
}
