package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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

var (
	ErrAlreadyExists = errors.New("resource already exists")
)

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
	err := router.doSetOne(user.OrganizationID, vars["name"], value.Value)
	if err != nil {
		log.Println(err)
		if errors.Is(err, ErrAlreadyExists) {
			SendAleadyExists(w)
		} else {
			SendInternalServerError(w)
		}
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
		err := router.doSetOne(user.OrganizationID, e.Name, e.Value)
		if err != nil {
			log.Println(err)
			if errors.Is(err, ErrAlreadyExists) {
				SendAleadyExists(w)
			} else {
				SendInternalServerError(w)
			}
			return
		}
	}
	SendUpdated(w)
}

func (router *SettingsRouter) doSetOne(organizationID, name, value string) error {
	if name == SettingConfluenceClientID.Name {
		currentClientID, err := GetSettingsRepository().Get(organizationID, name)
		if err != nil {
			return err
		}
		value = strings.TrimSpace(value)
		if currentClientID == value {
			// Nothing to changge
			return nil
		}
		if value != "" {
			// Check if any other org has this Client ID
			orgIDs, err := GetSettingsRepository().GetOrganizationIDsByValue(name, value)
			if err != nil {
				return err
			}
			if len(orgIDs) > 0 {
				return ErrAlreadyExists
			}
		}
		if currentClientID != "" {
			if value == "" {
				// If Client ID is removed: Delete all Confluence users
				users, err := GetUserRepository().GetUsersWithAtlassianID(organizationID)
				if err != nil {
					return err
				}
				for _, user := range users {
					if user.Email == string(user.AtlassianID) {
						GetUserRepository().Delete(user)
					} else {
						user.AtlassianID = ""
						GetUserRepository().Update(user)
					}
				}
			} else {
				// Else, change User IDs
				if err := GetUserRepository().UpdateAtlassianClientID(organizationID, currentClientID, value); err != nil {
					return err
				}
			}
		}
	}
	err := GetSettingsRepository().Set(organizationID, name, value)
	return err
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
		name == SettingConfluenceServerSharedSecret.Name ||
		name == SettingConfluenceClientID.Name ||
		name == SettingConfluenceAnonymous.Name ||
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
	if name == SettingConfluenceClientID.Name {
		return SettingConfluenceClientID.Type
	}
	if name == SettingConfluenceServerSharedSecret.Name {
		return SettingConfluenceServerSharedSecret.Type
	}
	if name == SettingConfluenceAnonymous.Name {
		return SettingConfluenceAnonymous.Type
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
	if settingType == SettingTypeString {
		return true
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
