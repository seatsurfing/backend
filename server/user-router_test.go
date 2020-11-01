package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestUserCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	username := uuid.New().String() + "@test.com"
	payload := "{\"email\": \"" + username + "\", \"password\": \"12345678\"}"
	req := newHTTPRequest("POST", "/user/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	userID := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/user/"+userID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetUserResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, username, resBody.Email)
	checkTestString(t, org.ID, resBody.OrganizationID)
	checkTestString(t, "", resBody.AuthProviderID)
	checkTestBool(t, true, resBody.RequirePassword)

	// 3. Update
	username = uuid.New().String() + "@test.com"
	payload = "{\"email\": \"" + username + "\", \"password\": \"\"}"
	req = newHTTPRequest("PUT", "/user/"+userID, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/user/"+userID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetUserResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, username, resBody2.Email)
	checkTestString(t, org.ID, resBody2.OrganizationID)
	checkTestString(t, "", resBody2.AuthProviderID)
	checkTestBool(t, true, resBody2.RequirePassword)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/user/"+userID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/user/"+userID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestUserForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	username := uuid.New().String() + "@test.com"
	payload := "{\"email\": \"" + username + "\", \"password\": \"12345678\"}"
	req := newHTTPRequest("POST", "/user/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// 2. Read
	req = newHTTPRequest("GET", "/user/"+user.ID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// 3. Update
	req = newHTTPRequest("PUT", "/user/"+user.ID, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/user/"+user.ID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestUserSetPassword(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"password": "12345678"}`
	req := newHTTPRequest("PUT", "/user/"+user.ID+"/password", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	user2, err := GetUserRepository().GetOne(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	checkTestBool(t, true, GetUserRepository().CheckPassword(string(user2.HashedPassword), "12345678"))
}

func TestUserSubscriptionExceeded(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	GetSettingsRepository().Set(org.ID, SettingSubscriptionMaxUsers.Name, "1")

	username := uuid.New().String() + "@test.com"
	payload := "{\"email\": \"" + username + "\", \"password\": \"12345678\"}"
	req := newHTTPRequest("POST", "/user/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusPaymentRequired, res.Code)
}

func TestUserGetCount(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	req := newHTTPRequest("GET", "/user/count", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetUserCountResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, 1, resBody.Count)
}

// TODO test domain in org!
