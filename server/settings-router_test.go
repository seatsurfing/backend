package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

func TestSettingsForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceServerSharedSecret.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	payload = `[]`
	req = newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestSettingsReadPublic(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	allowedSettings := []string{
		SettingMaxBookingsPerUser.Name,
		SettingMaxConcurrentBookingsPerUser.Name,
		SettingMaxDaysInAdvance.Name,
		SettingMaxBookingDurationHours.Name,
		SettingDailyBasisBooking.Name,
		SettingNoAdminRestrictions.Name,
		SettingShowNames.Name,
		SettingMinBookingDurationHours.Name,
		SettingAllowBookingsNonExistingUsers.Name,
		SettingDefaultTimezone.Name,
		SysSettingVersion,
	}
	forbiddenSettings := []string{
		SettingDatabaseVersion.Name,
		SettingAllowAnyUser.Name,
		SettingConfluenceServerSharedSecret.Name,
		SettingConfluenceAnonymous.Name,
		SettingActiveSubscription.Name,
		SettingSubscriptionMaxUsers.Name,
	}

	for _, name := range allowedSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusOK, res.Code)
	}

	for _, name := range forbiddenSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusForbidden, res.Code)
	}

	req := newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, len(allowedSettings), len(resBody))
	found := 0
	for _, name := range allowedSettings {
		for _, cur := range resBody {
			if name == cur.Name {
				found++
			}
		}
	}
	checkTestInt(t, len(allowedSettings), found)
}

func TestSettingsReadAdmin(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	allowedSettings := []string{
		SettingMaxBookingsPerUser.Name,
		SettingMaxConcurrentBookingsPerUser.Name,
		SettingMaxDaysInAdvance.Name,
		SettingMaxBookingDurationHours.Name,
		SettingDailyBasisBooking.Name,
		SettingMinBookingDurationHours.Name,
		SettingNoAdminRestrictions.Name,
		SettingShowNames.Name,
		SettingAllowBookingsNonExistingUsers.Name,
		SettingAllowAnyUser.Name,
		SettingConfluenceServerSharedSecret.Name,
		SettingConfluenceAnonymous.Name,
		SettingActiveSubscription.Name,
		SettingSubscriptionMaxUsers.Name,
		SettingDefaultTimezone.Name,
		SysSettingOrgSignupDelete,
		SysSettingVersion,
	}
	forbiddenSettings := []string{
		SettingDatabaseVersion.Name,
	}

	for _, name := range allowedSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusOK, res.Code)
	}

	for _, name := range forbiddenSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusForbidden, res.Code)
	}

	req := newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, len(allowedSettings), len(resBody))
	found := 0
	for _, name := range allowedSettings {
		for _, cur := range resBody {
			if name == cur.Name {
				found++
			}
		}
	}
	checkTestInt(t, len(allowedSettings), found)
}

func TestSettingsCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "1", resBody)

	payload = `{"value": "0"}`
	req = newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 string
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "0", resBody2)
}

func TestSettingsCRUDMany(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetDatabase().DB().Exec("TRUNCATE settings")

	payload := `[{"name": "allow_any_user", "value": "1"}, {"name": "max_bookings_per_user", "value": "5"}]`
	req := newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, 4, len(resBody))
	checkTestString(t, SettingAllowAnyUser.Name, resBody[0].Name)
	checkTestString(t, SettingMaxBookingsPerUser.Name, resBody[1].Name)
	checkTestString(t, SysSettingOrgSignupDelete, resBody[2].Name)
	checkTestString(t, SysSettingVersion, resBody[3].Name)
	checkTestString(t, "1", resBody[0].Value)
	checkTestString(t, "5", resBody[1].Value)
	checkTestString(t, GetProductVersion(), resBody[3].Value)

	payload = `[{"name": "allow_any_user", "value": "0"}, {"name": "max_bookings_per_user", "value": "3"}]`
	req = newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestInt(t, 4, len(resBody2))
	checkTestString(t, SettingAllowAnyUser.Name, resBody2[0].Name)
	checkTestString(t, SettingMaxBookingsPerUser.Name, resBody2[1].Name)
	checkTestString(t, SysSettingOrgSignupDelete, resBody2[2].Name)
	checkTestString(t, SysSettingVersion, resBody2[3].Name)
	checkTestString(t, "0", resBody2[0].Value)
	checkTestString(t, "3", resBody2[1].Value)

}

func TestSettingsMinHoursBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetDatabase().DB().Exec("TRUNCATE settings")

	payload := `[{"name": "min_booking_duration_hours", "value": "2"}]`
	req := newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody3 []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody3)
	log.Println(resBody3)
	checkTestInt(t, 3, len(resBody3))
	checkTestString(t, SettingMinBookingDurationHours.Name, resBody3[0].Name)
	checkTestString(t, SysSettingOrgSignupDelete, resBody3[1].Name)
	checkTestString(t, SysSettingVersion, resBody3[2].Name)
	checkTestString(t, "2", resBody3[0].Value)
}

func TestSettingsInvalidName(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/setting/test123", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestSettingsInvalidBool(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "2"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSettingsInvalidInt(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "test"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingMaxBookingsPerUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSettingsInvalidTimezone(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "Europe/Hamburg"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingDefaultTimezone.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)

	payload = `{"value": "Europe/Berlin"}`
	req = newHTTPRequest("PUT", "/setting/"+SettingDefaultTimezone.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
}
