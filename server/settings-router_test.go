package main

import (
	"bytes"
	"encoding/json"
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

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	payload = `[]`
	req = newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
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
	checkTestInt(t, 2, len(resBody))
	checkTestString(t, SettingAllowAnyUser.Name, resBody[0].Name)
	checkTestString(t, SettingMaxBookingsPerUser.Name, resBody[1].Name)
	checkTestString(t, "1", resBody[0].Value)
	checkTestString(t, "5", resBody[1].Value)

	payload = `[{"name": "allow_any_user", "value": "0"}, {"name": "max_bookings_per_user", "value": "3"}]`
	req = newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestInt(t, 2, len(resBody2))
	checkTestString(t, SettingAllowAnyUser.Name, resBody2[0].Name)
	checkTestString(t, SettingMaxBookingsPerUser.Name, resBody2[1].Name)
	checkTestString(t, "0", resBody2[0].Value)
	checkTestString(t, "3", resBody2[1].Value)
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
