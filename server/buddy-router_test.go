package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestBuddiesEmptyResult(t *testing.T) {
	clearTestDB()
	loginResponse := createLoginTestUser()

	req := newHTTPRequest("GET", "/buddy/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestBuddiesCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create buddy users
	buddyUser1 := createTestUserInOrg(org)

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	payload := "{\"buddyId\": \"" + buddyUser1.ID + "\"}"
	req := newHTTPRequest("POST", "/buddy/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read all buddies and ensure buddy was created correctly
	req = newHTTPRequest("GET", "/buddy/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetBuddyResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 1 {
		t.Fatalf("Expected array with 1 element")
	}
	checkTestString(t, buddyUser1.ID, resBody[0].BuddyID)
	checkTestString(t, id, resBody[0].ID)

	// 3. Delete
	req = newHTTPRequest("DELETE", "/buddy/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// 4. Read all buddies and ensure buddy was removed correctly
	req = newHTTPRequest("GET", "/buddy/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []*GetBuddyResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	if len(resBody2) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestDeleteBuddyOfAnotherUser(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create buddy users
	buddyUser1 := createTestUserInOrg(org)

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// Create
	payload := "{\"buddyId\": \"" + buddyUser1.ID + "\"}"
	req := newHTTPRequest("POST", "/buddy/", user2.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Delete
	req = newHTTPRequest("DELETE", "/buddy/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBuddiesCreateWithMissingUser(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")
	GetSettingsRepository().Set(org.ID, SettingAllowBookingsNonExistingUsers.Name, "1")

	// Create
	payload := "{\"buddyId\": \"" + uuid.New().String() + "\"}"
	req := newHTTPRequest("POST", "/buddy/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestBuddiesList(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create buddy users
	buddyUser1 := createTestUserInOrg(org)
	buddyUser2 := createTestUserInOrg(org)
	buddyUser3 := createTestUserInOrg(org)

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// Create #1
	payload := "{\"buddyId\": \"" + buddyUser1.ID + "\"}"
	req := newHTTPRequest("POST", "/buddy/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	payload = "{\"buddyId\": \"" + buddyUser2.ID + "\"}"
	req = newHTTPRequest("POST", "/buddy/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #3(for a different user)
	payload = "{\"buddyId\": \"" + buddyUser3.ID + "\"}"
	req = newHTTPRequest("POST", "/buddy/", buddyUser2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Read all buddies for user 1
	req = newHTTPRequest("GET", "/buddy/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetBuddyResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 2 {
		t.Fatalf("Expected array with 2 elements")
	}
	acceptedBuddyIDs := []string{buddyUser1.ID, buddyUser2.ID}
	if !contains(acceptedBuddyIDs, resBody[0].BuddyID) {
		t.Fatalf("Expected %s to one of %#v", resBody[0].BuddyID, acceptedBuddyIDs)
	}
	if !contains(acceptedBuddyIDs, resBody[1].BuddyID) {
		t.Fatalf("Expected %s to one of %#v", resBody[1].BuddyID, acceptedBuddyIDs)
	}

	// Read all buddies for user 2
	req = newHTTPRequest("GET", "/buddy/", buddyUser2.ID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []*GetBuddyResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	if len(resBody2) != 1 {
		t.Fatalf("Expected array with 1 elements")
	}
	checkTestString(t, buddyUser3.ID, resBody2[0].BuddyID)
}
