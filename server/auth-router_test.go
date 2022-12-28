package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuthPasswordLogin(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword("12345678"))
	GetUserRepository().Update(user)

	// Log in
	payload := "{ \"email\": \"" + user.Email + "\", \"password\": \"12345678\" }"
	req := newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *JWTResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestBool(t, true, len(resBody.AccessToken) > 32)
	checkTestBool(t, true, len(resBody.RefreshToken) == 36)

	// Test access token
	req = newHTTPRequestWithAccessToken("GET", "/user/me", resBody.AccessToken, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetUserResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, user.Email, resBody2.Email)
}

func TestAuthPasswordLoginBan(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword("12345678"))
	GetUserRepository().Update(user)

	// Attempt 1
	payload := "{ \"email\": \"" + user.Email + "\", \"password\": \"12345670\" }"
	req := newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 2
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"12345679\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 3
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"12345671\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
	checkTestBool(t, true, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Would be successful, but fails cause banned
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"12345678\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
	checkTestBool(t, true, authAttemptRepositoryIsUserDisabled(t, user.ID))
}

func TestAuthRefresh(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword("12345678"))
	GetUserRepository().Update(user)

	// Log in
	payload := "{ \"email\": \"" + user.Email + "\", \"password\": \"12345678\" }"
	req := newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *JWTResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestBool(t, true, len(resBody.AccessToken) > 32)
	checkTestBool(t, true, len(resBody.RefreshToken) == 36)

	// Sleep to ensure new access token
	time.Sleep(time.Second * 2)

	// Refresh access token
	payload = "{ \"refreshToken\": \"" + resBody.RefreshToken + "\" }"
	req = newHTTPRequest("POST", "/auth/refresh", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody3 *JWTResponse
	json.Unmarshal(res.Body.Bytes(), &resBody3)
	checkTestBool(t, true, len(resBody3.AccessToken) > 32)
	checkTestBool(t, true, len(resBody3.RefreshToken) == 36)
	checkTestBool(t, false, resBody3.AccessToken == resBody.AccessToken)
	checkTestBool(t, false, resBody3.RefreshToken == resBody.RefreshToken)

	// Test refreshed access token
	req = newHTTPRequestWithAccessToken("GET", "/user/me", resBody3.AccessToken, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetUserResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, user.Email, resBody2.Email)
}

func TestAuthRefreshNonExistent(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword("12345678"))
	GetUserRepository().Update(user)

	// Refresh access token
	payload := "{ \"refreshToken\": \"" + uuid.New().String() + "\" }"
	req := newHTTPRequest("POST", "/auth/refresh", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestAuthPasswordReset(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword("12345678"))
	GetUserRepository().Update(user)

	// Init password reset
	payload := "{ \"email\": \"" + user.Email + "\" }"
	req := newHTTPRequest("POST", "/auth/initpwreset", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hallo "+user.Email+","))

	// Extract Confirm ID from email
	rx := regexp.MustCompile(`/resetpw/(.*)?\n`)
	confirmTokens := rx.FindStringSubmatch(SendMailMockContent)
	checkTestInt(t, 2, len(confirmTokens))
	confirmID := confirmTokens[1]

	// Complete password reset
	payload = `{
			"password": "abcd1234"
		}`
	req = newHTTPRequest("POST", "/auth/pwreset/"+confirmID, "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Test login with old password
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"12345678\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)

	// Test login with new password
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"abcd1234\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
}

func TestAuthSingleOrg(t *testing.T) {
	clearTestDB()
	createTestOrg("test.com")

	req := newHTTPRequestWithAccessToken("GET", "/auth/singleorg", "", nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	var resBody *AuthPreflightResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestBool(t, false, resBody.RequirePassword)
	checkTestBool(t, false, resBody.Organization == nil)
	checkTestString(t, "Test Org", resBody.Organization.Name)
}

func TestAuthSingleOrgWithMultipleOrgs(t *testing.T) {
	clearTestDB()
	createTestOrg("test1.com")
	createTestOrg("test2.com")

	req := newHTTPRequestWithAccessToken("GET", "/auth/singleorg", "", nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}
