package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestBookingsEmptyResult(t *testing.T) {
	clearTestDB()
	loginResponse := createLoginTestUser()

	req := newHTTPRequest("GET", "/booking/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestBookingsCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "2030-09-01T08:30:00", resBody.Enter.Format(JsDateTimeFormat))
	checkTestString(t, "2030-09-01T17:00:00", resBody.Leave.Format(JsDateTimeFormat))
	checkTestString(t, spaceID, resBody.Space.ID)
	checkTestString(t, "H234", resBody.Space.Name)
	checkTestString(t, locationID, resBody.Space.Location.ID)
	checkTestString(t, "Location 1", resBody.Space.Location.Name)

	// 3. Update
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:45:00+02:00\", \"leave\": \"2030-09-01T18:15:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "2030-09-01T08:45:00", resBody2.Enter.Format(JsDateTimeFormat))
	checkTestString(t, "2030-09-01T18:15:00", resBody2.Leave.Format(JsDateTimeFormat))
	checkTestString(t, spaceID, resBody2.Space.ID)
	checkTestString(t, "H234", resBody2.Space.Name)
	checkTestString(t, locationID, resBody2.Space.Location.ID)
	checkTestString(t, "Location 1", resBody2.Space.Location.Name)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestBookingsList(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// Create #1
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	enter, _ := ParseJSDate("2019-09-01T08:30:00+02:00")
	leave, _ := ParseJSDate("2019-09-01T07:00:00+02:00")
	b2 := &Booking{
		SpaceID: spaceID,
		UserID:  loginResponse.UserID,
		Enter:   enter,
		Leave:   leave,
	}
	GetBookingRepository().Create(b2)

	// Create #3
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-06-01T08:30:00+02:00\", \"leave\": \"2030-06-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	req = newHTTPRequest("GET", "/booking/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 2 {
		t.Fatalf("Expected array with 2 elements")
	}
	checkTestString(t, "2030-06-01T08:30:00", resBody[0].Enter.Format(JsDateTimeFormat))
	checkTestString(t, "2030-09-01T08:30:00", resBody[1].Enter.Format(JsDateTimeFormat))
}

func TestBookingsGetForeign(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Switch to non-admin user 2
	user3 := createTestUserInOrg(org)
	loginResponse3 := loginTestUser(user3.ID)

	// 2. Read
	req, _ = http.NewRequest("GET", "/booking/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+getTestJWT(loginResponse3.UserID))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBookingsUpdateForeign(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// Create booking
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Create location #2
	payload = `{"name": "Location 2"}`
	req = newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID2 := res.Header().Get("X-Object-Id")

	// Create space #2
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID2+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID2 := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user3 := createTestUserInOrg(org)
	loginResponse3 := loginTestUser(user3.ID)

	// Update
	payload = "{\"spaceId\": \"" + spaceID2 + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req, _ = http.NewRequest("PUT", "/booking/"+id, bytes.NewBufferString(payload))
	req.Header.Set("Authorization", "Bearer "+getTestJWT(loginResponse3.UserID))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBookingsCreateForeign(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test2.com")
	user2 := createTestUserOrgAdminDomain(org, "test2.com")
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch user
	loginResponse3 := createLoginTestUserParams()

	// Create booking
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse3.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBookingsDeleteForeign(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	// Create booking
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user3 := createTestUserInOrg(org)
	loginResponse3 := loginTestUser(user3.ID)

	// Delete
	req = newHTTPRequest("DELETE", "/booking/"+id, loginResponse3.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBookingsConflictEnd(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Create #1
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T15:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestBookingsConflictStart(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Create #1
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T07:00:00+02:00\", \"leave\": \"2030-09-01T09:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestBookingsConflictInner(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Create #1
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T09:00:00+02:00\", \"leave\": \"2030-09-01T16:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestBookingsConflictOuter(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Create #1
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T07:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestBookingsConflictUpdateSelf(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Create
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Update
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T07:00:00+02:00\", \"leave\": \"2030-09-01T09:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
}

func TestBookingsConflictUpdateOther(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Create #1
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create #2
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T17:30:00+02:00\", \"leave\": \"2030-09-01T22:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Update #2
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T07:00:00+02:00\", \"leave\": \"2030-09-01T09:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestBookingsNegativeBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * -2).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsValidBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * 8).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "12")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * 14).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsDailyBasisBookingValid(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &CreateBookingRequest{
		Enter:   time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave:   time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsDailyBasisBookingSameDayValid(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	tm := time.Now().UTC()

	m := &CreateBookingRequest{
		Enter:   time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave:   time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsDailyBasisBookingInvalidEnter(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &CreateBookingRequest{
		Enter:   time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 1, 0, 0, tm.Location()),
		Leave:   time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsDailyBasisBookingInvalidLeave(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &CreateBookingRequest{
		Enter:   time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave:   time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 50, 59, 0, tm.Location()),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsDailyBasisBookingRoundBookingDurationUp(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "12")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &CreateBookingRequest{
		Enter:   time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave:   time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsValidBorderBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "3")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * 4).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidBorderBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "3")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * 4).Add(time.Minute * 1).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsPastEnterDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * -5).UTC(),
		Leave:   time.Now().Add(time.Hour * -2).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsValidFutureAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 2 * 24).UTC(),
		Leave:   time.Now().Add(time.Hour * 2 * 24).Add(time.Hour * 5).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsValidBorderAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 5 * 24).Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * 5 * 24).Add(time.Hour * 5).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidBorderAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 6 * 24).Add(time.Hour * 1).UTC(),
		Leave:   time.Now().Add(time.Hour * 6 * 24).Add(time.Hour * 5).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsInvalidFutureAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")

	m := &CreateBookingRequest{
		Enter:   time.Now().Add(time.Hour * 7 * 24).UTC(),
		Leave:   time.Now().Add(time.Hour * 7 * 24).Add(time.Hour * 5).UTC(),
		SpaceID: "",
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID)
	checkTestBool(t, false, res)
}

func TestBookingsValidMaxUpcomingBookings(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user := createTestUserInOrg(org)

	router := &BookingRouter{}
	res := router.isValidMaxUpcomingBookings(org.ID, user.ID)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidMaxUpcomingBookings(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user := createTestUserInOrg(org)

	l := &Location{
		Name:           "Test",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l)
	s := &Space{
		Name:       "Test",
		LocationID: l.ID,
	}
	GetSpaceRepository().Create(s)
	b := &Booking{
		Enter:   time.Now().Add(time.Hour * 6 * 24).UTC(),
		Leave:   time.Now().Add(time.Hour * 6 * 24).Add(time.Hour * 5).UTC(),
		SpaceID: s.ID,
		UserID:  user.ID,
	}
	GetBookingRepository().Create(b)

	router := &BookingRouter{}
	res := router.isValidMaxUpcomingBookings(org.ID, user.ID)
	checkTestBool(t, false, res)
}
