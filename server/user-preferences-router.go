package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type UserPreferencesRouter struct {
}

func (router *UserPreferencesRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/{name}", router.getPreference).Methods("GET")
	s.HandleFunc("/{name}", router.setPreference).Methods("PUT")
	s.HandleFunc("/", router.getAll).Methods("GET")
	s.HandleFunc("/", router.setAll).Methods("PUT")
}

func (router *UserPreferencesRouter) getPreference(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	vars := mux.Vars(r)
	if !router.isValidPreferenceName(vars["name"]) {
		SendNotFound(w)
		return
	}
	value, err := GetUserPreferencesRepository().Get(user.ID, vars["name"])
	if err != nil {
		log.Println(err)
		SendNotFound(w)
		return
	}
	SendJSON(w, value)
}

func (router *UserPreferencesRouter) setPreference(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	var value SetSettingsRequest
	if UnmarshalValidateBody(r, &value) != nil {
		SendBadRequest(w)
		return
	}
	vars := mux.Vars(r)
	if !router.isValidPreferenceName(vars["name"]) {
		SendNotFound(w)
		return
	}
	if !router.isValidPreferncesType(vars["name"], value.Value) {
		SendBadRequest(w)
		return
	}
	if !router.isValidPreferenceValue(vars["name"], value.Value) {
		SendBadRequest(w)
		return
	}
	err := router.doSetOne(user.ID, vars["name"], value.Value)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	SendUpdated(w)
}

func (router *UserPreferencesRouter) getAll(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	list, err := GetUserPreferencesRepository().GetAll(user.ID)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	res := []*GetSettingsResponse{}
	for _, e := range list {
		m := router.copyToRestModel(e)
		res = append(res, m)
	}
	SendJSON(w, res)
}

func (router *UserPreferencesRouter) setAll(w http.ResponseWriter, r *http.Request) {
	user := GetRequestUser(r)
	var list []GetSettingsResponse
	if err := UnmarshalBody(r, &list); err != nil {
		log.Println(err)
		SendBadRequest(w)
		return
	}
	for _, e := range list {
		if !router.isValidPreferenceName(e.Name) {
			SendNotFound(w)
			return
		}
		if !router.isValidPreferncesType(e.Name, e.Value) {
			SendBadRequest(w)
			return
		}
		if !router.isValidPreferenceValue(e.Name, e.Value) {
			SendBadRequest(w)
			return
		}
		err := router.doSetOne(user.ID, e.Name, e.Value)
		if err != nil {
			log.Println(err)
			SendInternalServerError(w)
			return
		}
	}
	SendUpdated(w)
}

func (router *UserPreferencesRouter) doSetOne(userID, name, value string) error {
	err := GetUserPreferencesRepository().Set(userID, name, value)
	return err
}

func (router *UserPreferencesRouter) isValidPreferenceName(name string) bool {
	if name == PreferenceEnterTime.Name ||
		name == PreferenceWorkdayStart.Name ||
		name == PreferenceWorkdayEnd.Name ||
		name == PreferenceWorkdays.Name ||
		name == PreferenceLocation.Name {
		return true
	}
	return false
}

func (router *UserPreferencesRouter) getPreferenceType(name string) SettingType {
	if name == PreferenceEnterTime.Name {
		return PreferenceEnterTime.Type
	}
	if name == PreferenceWorkdayStart.Name {
		return PreferenceWorkdayStart.Type
	}
	if name == PreferenceWorkdayEnd.Name {
		return PreferenceWorkdayEnd.Type
	}
	if name == PreferenceWorkdays.Name {
		return PreferenceWorkdays.Type
	}
	if name == PreferenceLocation.Name {
		return PreferenceLocation.Type
	}
	return 0
}

func (router *UserPreferencesRouter) isValidPreferncesType(name string, value string) bool {
	settingType := router.getPreferenceType(name)
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
	if settingType == SettingTypeIntArray {
		tokens := strings.Split(value, ",")
		ok := true
		for _, token := range tokens {
			if _, err := strconv.Atoi(token); err != nil {
				ok = false
			}
		}
		return ok
	}
	return false
}

func (router *UserPreferencesRouter) isValidPreferenceValue(name string, value string) bool {
	if name == PreferenceEnterTime.Name {
		i, _ := strconv.Atoi(value)
		if !(i == PreferenceEnterTimeNow || i == PreferenceEnterTimeNextDay || i == PreferenceEnterTimeNextWorkday) {
			return false
		}
	}
	if name == PreferenceWorkdayStart.Name {
		i, _ := strconv.Atoi(value)
		if i < 0 || i > 24 {
			return false
		}
	}
	if name == PreferenceWorkdayEnd.Name {
		i, _ := strconv.Atoi(value)
		if i < 0 || i > 24 {
			return false
		}
	}
	if name == PreferenceWorkdays.Name {
		tokens := strings.Split(value, ",")
		ok := true
		for _, token := range tokens {
			if workday, err := strconv.Atoi(token); err != nil || workday < 0 || workday > 6 {
				ok = false
			}
		}
		return ok
	}
	return true
}

func (router *UserPreferencesRouter) copyToRestModel(e *UserPreference) *GetSettingsResponse {
	m := &GetSettingsResponse{}
	m.Name = e.Name
	m.Value = e.Value
	return m
}
