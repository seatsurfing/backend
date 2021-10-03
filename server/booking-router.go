package main

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type BookingRouter struct {
}

type BookingRequest struct {
	Enter time.Time `json:"enter" validate:"required"`
	Leave time.Time `json:"leave" validate:"required"`
}

type CreateBookingRequest struct {
	SpaceID string `json:"spaceId" validate:"required"`
	BookingRequest
}

type PreCreateBookingRequest struct {
	LocationID string `json:"locationID" validate:"required"`
	BookingRequest
}

type GetBookingResponse struct {
	ID        string           `json:"id"`
	UserID    string           `json:"userId"`
	UserEmail string           `json:"userEmail"`
	Space     GetSpaceResponse `json:"space"`
	CreateBookingRequest
}

type GetBookingFilterRequest struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
}

func (router *BookingRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/filter/", router.getFiltered).Methods("POST")
	s.HandleFunc("/precheck/", router.preBookingCreateCheck).Methods("POST")
	s.HandleFunc("/{id}", router.getOne).Methods("GET")
	s.HandleFunc("/{id}", router.update).Methods("PUT")
	s.HandleFunc("/{id}", router.delete).Methods("DELETE")
	s.HandleFunc("/", router.create).Methods("POST")
	s.HandleFunc("/", router.getAll).Methods("GET")
}

func (router *BookingRouter) getFiltered(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	var m GetBookingFilterRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	list, err := GetBookingRepository().GetAllByOrg(user.OrganizationID, m.Start, m.End)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetBookingResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *BookingRouter) getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetBookingRepository().GetOne(vars["id"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	if e.UserID != GetRequestUserID(r) {
		SendForbidden(w)
		return
	}
	res := router.copyToRestModel(e)
	SendJSON(w, res)
}

func (router *BookingRouter) getAll(w http.ResponseWriter, r *http.Request) {
	list, err := GetBookingRepository().GetAllByUser(GetRequestUserID(r), time.Now().UTC())
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetBookingResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *BookingRouter) update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetBookingRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	var m CreateBookingRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	space, err := GetSpaceRepository().GetOne(m.SpaceID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	location, err := GetLocationRepository().GetOne(space.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	if !CanAccessOrg(GetRequestUser(r), location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if e.UserID != GetRequestUserID(r) {
		SendForbidden(w)
		return
	}
	if valid, code := router.isValidBookingRequest(&m.BookingRequest, GetRequestUserID(r), location.OrganizationID, true); !valid {
		SendBadRequestCode(w, code)
		return
	}
	if !router.isValidConcurrent(&m.BookingRequest, location, e.ID) {
		SendBadRequestCode(w, ResponseCodeBookingLocationMaxConcurrent)
		return
	}
	conflicts, err := GetBookingRepository().GetConflicts(m.SpaceID, m.Enter, m.Leave, e.ID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	if len(conflicts) > 0 {
		SendAleadyExists(w)
		return
	}
	eNew := router.copyFromRestModel(&m)
	eNew.ID = vars["id"]
	eNew.UserID = GetRequestUserID(r)
	if err := GetBookingRepository().Update(eNew); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *BookingRouter) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e, err := GetBookingRepository().GetOne(vars["id"])
	if err != nil {
		SendNotFound(w)
		return
	}
	space, err := GetSpaceRepository().GetOne(e.SpaceID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	location, err := GetLocationRepository().GetOne(space.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	if !CanAccessOrg(GetRequestUser(r), location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if e.UserID != GetRequestUserID(r) {
		SendForbidden(w)
		return
	}
	if err := GetBookingRepository().Delete(e); err != nil {
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *BookingRouter) checkBookingCreateUpdate(m *BookingRequest, location *Location, requestUser *User, bookingID string) (bool, int) {
	isUpdate := bookingID != ""
	if valid, code := router.isValidBookingRequest(m, requestUser.ID, location.OrganizationID, isUpdate); !valid {
		return false, code
	}
	if !router.isValidConcurrent(m, location, bookingID) {
		return false, ResponseCodeBookingLocationMaxConcurrent
	}
	return true, 0
}

func (router *BookingRouter) preBookingCreateCheck(w http.ResponseWriter, r *http.Request) {
	var m PreCreateBookingRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	location, err := GetLocationRepository().GetOne(m.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	requestUser := GetRequestUser(r)
	if !CanAccessOrg(requestUser, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if valid, code := router.checkBookingCreateUpdate(&m.BookingRequest, location, requestUser, ""); !valid {
		SendBadRequestCode(w, code)
		return
	}
	SendUpdated(w)
}

func (router *BookingRouter) create(w http.ResponseWriter, r *http.Request) {
	var m CreateBookingRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	space, err := GetSpaceRepository().GetOne(m.SpaceID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	location, err := GetLocationRepository().GetOne(space.LocationID)
	if err != nil {
		SendBadRequest(w)
		return
	}
	requestUser := GetRequestUser(r)
	if !CanAccessOrg(requestUser, location.OrganizationID) {
		SendForbidden(w)
		return
	}
	if valid, code := router.checkBookingCreateUpdate(&m.BookingRequest, location, requestUser, ""); !valid {
		SendBadRequestCode(w, code)
		return
	}
	conflicts, err := GetBookingRepository().GetConflicts(m.SpaceID, m.Enter, m.Leave, "")
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	if len(conflicts) > 0 {
		SendAleadyExists(w)
		return
	}
	e := router.copyFromRestModel(&m)
	e.UserID = GetRequestUserID(r)
	if err := GetBookingRepository().Create(e); err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendCreated(w, e.ID)
}

func (router *BookingRouter) isValidBookingDuration(m *BookingRequest, orgID string) bool {
	dailyBasisBooking, _ := GetSettingsRepository().GetBool(orgID, SettingDailyBasisBooking.Name)
	maxDurationHours, _ := GetSettingsRepository().GetInt(orgID, SettingMaxBookingDurationHours.Name)
	if dailyBasisBooking && (maxDurationHours%24 != 0) {
		maxDurationHours += (24 - (maxDurationHours % 24))
	}
	duration := math.Floor(m.Leave.Sub(m.Enter).Minutes()) / 60
	if duration < 0 || duration > float64(maxDurationHours) {
		return false
	}
	durationNotRounded := int(math.Round(m.Leave.Sub(m.Enter).Minutes()) / 60)
	if dailyBasisBooking && (durationNotRounded%24 != 0) {
		return false
	}
	return true
}

func (router *BookingRouter) isValidBookingAdvance(m *BookingRequest, orgID string) bool {
	maxAdvanceDays, _ := GetSettingsRepository().GetInt(orgID, SettingMaxDaysInAdvance.Name)
	now := time.Now().UTC()
	dailyBasisBooking, _ := GetSettingsRepository().GetBool(orgID, SettingDailyBasisBooking.Name)
	if dailyBasisBooking {
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		now = now.Add(-12 * time.Hour)
	}
	advanceDays := math.Floor(m.Enter.Sub(now).Hours() / 24)
	if advanceDays < 0 || advanceDays > float64(maxAdvanceDays) {
		return false
	}
	return true
}

func (router *BookingRouter) isValidMaxUpcomingBookings(orgID string, userID string) bool {
	maxUpcoming, _ := GetSettingsRepository().GetInt(orgID, SettingMaxBookingsPerUser.Name)
	curUpcoming, _ := GetBookingRepository().GetAllByUser(userID, time.Now().UTC())
	return len(curUpcoming) < maxUpcoming
}

func (router *BookingRouter) isValidBookingRequest(m *BookingRequest, userID string, orgID string, isUpdate bool) (bool, int) {
	if !router.isValidBookingDuration(m, orgID) {
		return false, ResponseCodeBookingInvalidBookingDuration
	}
	if !router.isValidBookingAdvance(m, orgID) {
		return false, ResponseCodeBookingTooManyDaysInAdvance
	}
	if !isUpdate {
		if !router.isValidMaxUpcomingBookings(orgID, userID) {
			return false, ResponseCodeBookingTooManyUpcomingBookings
		}
	}
	return true, 0
}

func (router *BookingRouter) isValidConcurrent(m *BookingRequest, location *Location, bookingID string) bool {
	if location.MaxConcurrentBookings == 0 {
		return true
	}
	bookings, err := GetBookingRepository().GetConcurrent(location.ID, m.Enter, m.Leave, bookingID)
	if err != nil {
		log.Println(err)
		return false
	}
	if bookings >= int(location.MaxConcurrentBookings) {
		return false
	}
	return true
}

func (router *BookingRouter) copyFromRestModel(m *CreateBookingRequest) *Booking {
	e := &Booking{}
	e.SpaceID = m.SpaceID
	e.Enter = m.Enter
	e.Leave = m.Leave
	return e
}

func (router *BookingRouter) copyToRestModel(e *BookingDetails) *GetBookingResponse {
	m := &GetBookingResponse{}
	m.ID = e.ID
	m.UserID = e.UserID
	m.UserEmail = e.UserEmail
	m.SpaceID = e.SpaceID
	m.Enter = e.Enter
	m.Leave = e.Leave
	m.Space.ID = e.Space.ID
	m.Space.LocationID = e.Space.LocationID
	m.Space.Name = e.Space.Name
	m.Space.Location.ID = e.Space.Location.ID
	m.Space.Location.Name = e.Space.Location.Name
	return m
}
