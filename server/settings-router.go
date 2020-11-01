package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SettingsRouter struct {
}

type SetSettingsRequest struct {
	Value string `json:"value"`
}

type GetSettingsResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (router *SettingsRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/{name}", router.getSetting).Methods("GET")
	s.HandleFunc("/{name}", router.setSetting).Methods("PUT")
	s.HandleFunc("/", router.getAll).Methods("GET")
	s.HandleFunc("/", router.setAll).Methods("PUT")
}

func (router *SettingsRouter) getSetting(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	vars := mux.Vars(r)
	if !router.isValidSettingNameReadAdmin(vars["name"]) {
		SendNotFound(w)
		return
	}
	value, err := GetSettingsRepository().Get(user.OrganizationID, vars["name"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	SendJSON(w, value)
}

func (router *SettingsRouter) setSetting(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	var value SetSettingsRequest
	if UnmarshalValidateBody(r, &value) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	if !router.isValidSettingNameWrite(vars["name"]) {
		SendNotFound(w)
		return
	}
	if !router.isValidSettingType(vars["name"], value.Value) {
		SendBadRequest(w)
		return
	}
	err := GetSettingsRepository().Set(user.OrganizationID, vars["name"], value.Value)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *SettingsRouter) getAll(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAccessOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	orgAdmin := CanAdminOrg(user, user.OrganizationID)
	list, err := GetSettingsRepository().GetAll(user.OrganizationID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetSettingsResponse{}
	for _, e := range list {
		if (orgAdmin && router.isValidSettingNameReadAdmin(e.Name)) || (router.isValidSettingNameReadPublic(e.Name)) {
			m := router.copyToRestModel(e)
			res = append(res, m)
		}
	}
	SendJSON(w, res)
}

func (router *SettingsRouter) setAll(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	if !CanAdminOrg(user, user.OrganizationID) {
		SendForbidden(w)
		return
	}
	var list []GetSettingsResponse
	if err := UnmarshalBody(r, &list); err != nil {
		log.Println(err)
		SendBadRequest(w)
		return
	}
	for _, e := range list {
		if !router.isValidSettingNameWrite(e.Name) {
			SendNotFound(w)
			return
		}
		if !router.isValidSettingType(e.Name, e.Value) {
			SendBadRequest(w)
			return
		}
		err := GetSettingsRepository().Set(user.OrganizationID, e.Name, e.Value)
		if err != nil {
			log.Println(err)
			SendInternalServerError(w)
			return
		}
	}
	SendUpdated(w)
}

func (router *SettingsRouter) copyToRestModel(e *OrgSetting) *GetSettingsResponse {
	m := &GetSettingsResponse{}
	m.Name = e.Name
	m.Value = e.Value
	return m
}

func (router *SettingsRouter) isValidSettingNameReadPublic(name string) bool {
	if router.isValidSettingNameWrite(name) ||
		name == SettingSubscriptionMaxUsers.Name ||
		name == SettingActiveSubscription.Name {
		return true
	}
	return false
}

func (router *SettingsRouter) isValidSettingNameReadAdmin(name string) bool {
	if router.isValidSettingNameReadPublic(name) ||
		name == SettingFastSpringAccountID.Name ||
		name == SettingFastSpringSubscriptionID.Name {
		return true
	}
	return false
}

func (router *SettingsRouter) isValidSettingNameWrite(name string) bool {
	if name == SettingAllowAnyUser.Name ||
		name == SettingMaxBookingsPerUser.Name ||
		name == SettingMaxDaysInAdvance.Name ||
		name == SettingMaxBookingDurationHours.Name {
		return true
	}
	return false
}

func (router *SettingsRouter) getSettingType(name string) SettingType {
	if name == SettingAllowAnyUser.Name {
		return SettingAllowAnyUser.Type
	}
	if name == SettingMaxBookingsPerUser.Name {
		return SettingMaxBookingsPerUser.Type
	}
	if name == SettingMaxDaysInAdvance.Name {
		return SettingMaxDaysInAdvance.Type
	}
	if name == SettingMaxBookingDurationHours.Name {
		return SettingMaxBookingDurationHours.Type
	}
	return 0
}

func (router *SettingsRouter) isValidSettingType(name string, value string) bool {
	settingType := router.getSettingType(name)
	if settingType == 0 {
		return false
	}
	if settingType == SettingTypeBool && (value == "1" || value == "0") {
		return true
	}
	if settingType == SettingTypeInt {
		if _, err := strconv.Atoi(value); err == nil {
			return true
		}
	}
	return false
}
