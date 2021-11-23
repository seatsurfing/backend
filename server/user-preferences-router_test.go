package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestUserPreferencesCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/preference/"+PreferenceEnterTime.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/preference/"+PreferenceEnterTime.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "1", resBody)

	payload = `{"value": "2"}`
	req = newHTTPRequest("PUT", "/preference/"+PreferenceEnterTime.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/preference/"+PreferenceEnterTime.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 string
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "2", resBody2)
}

func TestUserPreferencesCRUDMany(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)
	GetDatabase().DB().Exec("TRUNCATE users_preferences")

	payload := `[{"name": "enter_time", "value": "1"}, {"name": "workday_start", "value": "5"}]`
	req := newHTTPRequest("PUT", "/preference/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/preference/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, 2, len(resBody))
	checkTestString(t, PreferenceEnterTime.Name, resBody[0].Name)
	checkTestString(t, PreferenceWorkdayStart.Name, resBody[1].Name)
	checkTestString(t, "1", resBody[0].Value)
	checkTestString(t, "5", resBody[1].Value)

	payload = `[{"name": "enter_time", "value": "2"}, {"name": "workday_start", "value": "3"}]`
	req = newHTTPRequest("PUT", "/preference/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/preference/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestInt(t, 2, len(resBody2))
	checkTestString(t, PreferenceEnterTime.Name, resBody2[0].Name)
	checkTestString(t, PreferenceWorkdayStart.Name, resBody2[1].Name)
	checkTestString(t, "2", resBody2[0].Value)
	checkTestString(t, "3", resBody2[1].Value)
}
