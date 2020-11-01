package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuthProvidersEmptyResult(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	req := newHTTPRequest("GET", "/auth-provider/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestAuthProvidersForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	userAdmin := createTestUserOrgAdmin(org)
	loginResponseAdmin := loginTestUser(userAdmin.ID)
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"name": "Test", "providerType": 1, "clientId": "test1", "clientSecret": "test2", "authUrl": "http://test.com/1", "tokenUrl": "http://test.com/2", "authStyle": 0, "scopes": "http://test.com/3", "userInfoUrl": "http://test.com/userinfo", "userInfoEmailField": "email"}`
	req := newHTTPRequest("POST", "/auth-provider/", loginResponseAdmin.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	req = newHTTPRequest("GET", "/auth-provider/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("POST", "/auth-provider/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("DELETE", "/auth-provider/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("PUT", "/auth-provider/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/auth-provider/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestAuthProvidersCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	userAdmin := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(userAdmin.ID)

	// 1. Create
	payload := `{"name": "Test", "providerType": 1, "clientId": "test1", "clientSecret": "test2", "authUrl": "http://test.com/1", "tokenUrl": "http://test.com/2", "authStyle": 0, "scopes": "http://test.com/3", "userInfoUrl": "http://test.com/userinfo", "userInfoEmailField": "email"}`
	req := newHTTPRequest("POST", "/auth-provider/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/auth-provider/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetAuthProviderResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "Test", resBody.Name)
	checkTestString(t, "test1", resBody.ClientID)
	checkTestString(t, "test2", resBody.ClientSecret)
	checkTestString(t, "http://test.com/1", resBody.AuthURL)
	checkTestString(t, "http://test.com/2", resBody.TokenURL)
	checkTestInt(t, 0, resBody.AuthStyle)
	checkTestString(t, "http://test.com/3", resBody.Scopes)
	checkTestString(t, org.ID, resBody.OrganizationID)
	checkTestString(t, "http://test.com/userinfo", resBody.UserInfoURL)
	checkTestString(t, "email", resBody.UserInfoEmailField)
	checkTestInt(t, int(OAuth2), resBody.ProviderType)

	// 3. Update
	payload = `{"name": "Test_2", "providerType": 1, "clientId": "test1_2", "clientSecret": "test2_2", "authUrl": "http://test.com/1_2", "tokenUrl": "http://test.com/2_2", "authStyle": 1, "scopes": "http://test.com/3_2", "userInfoUrl": "http://test.com/userinfo_2", "userInfoEmailField": "email_2"}`
	req = newHTTPRequest("PUT", "/auth-provider/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/auth-provider/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetAuthProviderResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "Test_2", resBody2.Name)
	checkTestString(t, "test1_2", resBody2.ClientID)
	checkTestString(t, "test2_2", resBody2.ClientSecret)
	checkTestString(t, "http://test.com/1_2", resBody2.AuthURL)
	checkTestString(t, "http://test.com/2_2", resBody2.TokenURL)
	checkTestInt(t, 1, resBody2.AuthStyle)
	checkTestString(t, "http://test.com/3_2", resBody2.Scopes)
	checkTestString(t, org.ID, resBody2.OrganizationID)
	checkTestString(t, "http://test.com/userinfo_2", resBody2.UserInfoURL)
	checkTestString(t, "email_2", resBody2.UserInfoEmailField)
	checkTestInt(t, int(OAuth2), resBody2.ProviderType)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/auth-provider/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/auth-provider/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestAuthProvidersGetPublicForOrg(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	userAdmin := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(userAdmin.ID)

	// Create 1
	payload := `{"name": "Test", "providerType": 1, "clientId": "test1", "clientSecret": "test2", "authUrl": "http://test.com/1", "tokenUrl": "http://test.com/2", "authStyle": 0, "scopes": "http://test.com/3", "userInfoUrl": "http://test.com/userinfo", "userInfoEmailField": "email"}`
	req := newHTTPRequest("POST", "/auth-provider/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id1 := res.Header().Get("X-Object-Id")

	// Create 2
	payload = `{"name": "Test2", "providerType": 2, "clientId": "test2", "clientSecret": "test3", "authUrl": "http://test.com/7", "tokenUrl": "http://test.com/8", "authStyle": 0, "scopes": "http://test.com/9", "userInfoUrl": "http://test.com/userinfo", "userInfoEmailField": "email"}`
	req = newHTTPRequest("POST", "/auth-provider/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id2 := res.Header().Get("X-Object-Id")

	// Get Public List
	req = newHTTPRequest("GET", "/auth-provider/org/"+org.ID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetAuthProviderPublicResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 2 {
		t.Fatalf("Expected array with 2 elements")
	}
	checkTestString(t, id1, resBody[0].ID)
	checkTestString(t, "Test", resBody[0].Name)
	checkTestString(t, id2, resBody[1].ID)
	checkTestString(t, "Test2", resBody[1].Name)
}
