package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSpaceAttributesEmptyResult(t *testing.T) {
	clearTestDB()
	loginResponse := createLoginTestUser()

	req := newHTTPRequest("GET", "/space-attribute/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestSpaceAttributesCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	payload := `{"label": "Test 123", "type": 3, "spaceApplicable": true, "locationApplicable": false}`
	req := newHTTPRequest("POST", "/space-attribute/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/space-attribute/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetSpaceAttributeResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "Test 123", resBody.Label)
	checkTestBool(t, true, resBody.SpaceApplicable)
	checkTestBool(t, false, resBody.LocationApplicable)

	// 3. Update
	payload = `{"label": "Test 456", "type": 2, "spaceApplicable": false, "locationApplicable": true}`
	req = newHTTPRequest("PUT", "/space-attribute/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/space-attribute/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetSpaceAttributeResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "Test 456", resBody2.Label)
	checkTestBool(t, false, resBody2.SpaceApplicable)
	checkTestBool(t, true, resBody2.LocationApplicable)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/space-attribute/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/space-attribute/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}
