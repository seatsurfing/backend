package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
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
	adminUser := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(adminUser.ID)
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
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00Z\", \"leave\": \"2030-09-01T17:00:00Z\"}"
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
	checkTestString(t, "2030-09-01T08:30:00+02:00", resBody.Enter.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, "2030-09-01T17:00:00+02:00", resBody.Leave.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, spaceID, resBody.Space.ID)
	checkTestString(t, "H234", resBody.Space.Name)
	checkTestString(t, locationID, resBody.Space.Location.ID)
	checkTestString(t, "Location 1", resBody.Space.Location.Name)

	// 3. Update by admin
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:45:00Z\", \"leave\": \"2030-09-01T18:15:00Z\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "2030-09-01T08:45:00+02:00", resBody2.Enter.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, "2030-09-01T18:15:00+02:00", resBody2.Leave.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, spaceID, resBody2.Space.ID)
	checkTestString(t, "H234", resBody2.Space.Name)
	checkTestString(t, locationID, resBody2.Space.Location.ID)
	checkTestString(t, "Location 1", resBody2.Space.Location.Name)

	// 3. Update by Non-admin
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T09:00:00Z\", \"leave\": \"2030-09-01T18:15:00Z\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, user.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestBookingsCreateNonExistingUser(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")
	GetSettingsRepository().Set(org.ID, SettingAllowBookingsNonExistingUsers.Name, "1")

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
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00Z\", \"leave\": \"2030-09-01T17:00:00Z\", \"userEmail\": \"new-user@test.com\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Check user
	newUser, _ := GetUserRepository().GetByEmail("new-user@test.com")
	checkTestBool(t, true, newUser != nil)

	// Check booking
	booking, _ := GetBookingRepository().GetOne(id)
	checkTestBool(t, true, booking != nil)
}

func TestBookingsCreateNonExistingUserNoAdmin(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")
	GetSettingsRepository().Set(org.ID, SettingAllowBookingsNonExistingUsers.Name, "1")

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

	// Create
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00Z\", \"leave\": \"2030-09-01T17:00:00Z\", \"userEmail\": \"new-user@test.com\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBookingsCreateNonExistingUserNotEnabled(t *testing.T) {
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
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00Z\", \"leave\": \"2030-09-01T17:00:00Z\", \"userEmail\": \"new-user@test.com\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestBookingsCreateNonExistingUserForeignDomain(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")
	GetSettingsRepository().Set(org.ID, SettingAllowBookingsNonExistingUsers.Name, "1")

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
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00Z\", \"leave\": \"2030-09-01T17:00:00Z\", \"userEmail\": \"new-user@test2.com\"}"
	req = newHTTPRequest("POST", "/booking/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
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
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00Z\", \"leave\": \"2030-09-01T17:00:00Z\"}"
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
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-06-01T08:30:00Z\", \"leave\": \"2030-06-01T17:00:00Z\"}"
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
	checkTestString(t, "2030-06-01T08:30:00+02:00", resBody[0].Enter.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, "2030-09-01T08:30:00+02:00", resBody[1].Enter.Format(JsDateTimeFormatWithTimezone))
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
	adminUser := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(adminUser.ID)
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

	// Update foreign booking as admin
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:30:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

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

func TestBookingsConflictTooClose(t *testing.T) {
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

	// Create booking for tomorrow
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05-07:00")
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\":" + "\"" + tomorrow + "\"" + ", \"leave\":" + "\"" + tomorrow + "\"" + "}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Create another booking booking (not for tomorrow but for the next week)
	another_date := time.Now().Add(168 * time.Hour).Format("2006-01-02T15:04:05-07:00")
	payload = "{\"spaceId\": \"" + spaceID + "\", \"enter\":" + "\"" + another_date + "\"" + ", \"leave\":" + "\"" + another_date + "\"" + "}"
	req = newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id2 := res.Header().Get("X-Object-Id")

	// Delete with Error for tomorrow booking
	req = newHTTPRequest("DELETE", "/booking/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// Delete without Error for next week booking
	req = newHTTPRequest("DELETE", "/booking/"+id2, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
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

func TestBookingsDeleteSpaceAdmin(t *testing.T) {
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

	// Switch to space admin user
	user3 := createTestUserInOrgWithName(org, uuid.New().String()+"@test.com", UserRoleSpaceAdmin)
	loginResponse3 := loginTestUser(user3.ID)

	// Delete
	req = newHTTPRequest("DELETE", "/booking/"+id, loginResponse3.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
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
	user := createTestUserInOrg(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * -2).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, false, res)
}

func TestBookingsValidBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * 8).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "12")
	user := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * 14).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, false, res)

	res = router.isValidBookingDuration(m, org.ID, adminUser)
	checkTestBool(t, true, res)

	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "0")
	res = router.isValidBookingDuration(m, org.ID, adminUser)
	checkTestBool(t, false, res)

}

func TestBookingsDailyBasisBookingValid(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	user := createTestUserInOrg(org)
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &BookingRequest{
		Enter: time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave: time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsDailyBasisBookingSameDayValid(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	user := createTestUserInOrg(org)
	tm := time.Now().UTC()

	m := &BookingRequest{
		Enter: time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave: time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsDailyBasisBookingInvalidEnter(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	user := createTestUserInOrg(org)
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &BookingRequest{
		Enter: time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 1, 0, 0, tm.Location()),
		Leave: time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, false, res)
}

func TestBookingsDailyBasisBookingInvalidLeave(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	user := createTestUserInOrg(org)
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &BookingRequest{
		Enter: time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave: time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 50, 59, 0, tm.Location()),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, false, res)
}

func TestBookingsDailyBasisBookingRoundBookingDurationUp(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "12")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	user := createTestUserInOrg(org)
	tm := time.Now().Add(time.Hour * 24).UTC()

	m := &BookingRequest{
		Enter: time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()),
		Leave: time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsValidBorderBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "3")
	user := createTestUserInOrg(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * 4).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidBorderBookingDuration(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "3")
	user := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * 4).Add(time.Minute * 1).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingDuration(m, org.ID, user)
	checkTestBool(t, false, res)

	res = router.isValidBookingDuration(m, org.ID, adminUser)
	checkTestBool(t, true, res)

	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "0")
	res = router.isValidBookingDuration(m, org.ID, adminUser)
	checkTestBool(t, false, res)
}

func TestBookingsPastEnterDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")
	user := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * -25).UTC(),
		Leave: time.Now().Add(time.Hour * 1).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, false, res)

	// also admins cannot book in past
	res = router.isValidBookingAdvance(m, org.ID, adminUser)
	checkTestBool(t, false, res)
}

func TestBookingsEarlyMorningEnterDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")
	user := createTestUserInOrg(org)

	now := time.Now().UTC()
	m := &BookingRequest{
		Enter: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		Leave: time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location()),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsValidFutureAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")
	user := createTestUserInOrg(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 2 * 24).UTC(),
		Leave: time.Now().Add(time.Hour * 2 * 24).Add(time.Hour * 5).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsValidBorderAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")
	user := createTestUserInOrg(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 5 * 24).Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * 5 * 24).Add(time.Hour * 5).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, true, res)
}

func TestBookingsInvalidBorderAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")
	user := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 6 * 24).Add(time.Hour * 1).UTC(),
		Leave: time.Now().Add(time.Hour * 6 * 24).Add(time.Hour * 5).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, false, res)

	res = router.isValidBookingAdvance(m, org.ID, adminUser)
	checkTestBool(t, true, res)

	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "0")
	res = router.isValidBookingAdvance(m, org.ID, adminUser)
	checkTestBool(t, false, res)

}

func TestBookingsInvalidFutureAdvanceDate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5")
	user := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	m := &BookingRequest{
		Enter: time.Now().Add(time.Hour * 7 * 24).UTC(),
		Leave: time.Now().Add(time.Hour * 7 * 24).Add(time.Hour * 5).UTC(),
	}

	router := &BookingRouter{}
	res := router.isValidBookingAdvance(m, org.ID, user)
	checkTestBool(t, false, res)

	res = router.isValidBookingAdvance(m, org.ID, adminUser)
	checkTestBool(t, true, res)

	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "0")
	res = router.isValidBookingAdvance(m, org.ID, adminUser)
	checkTestBool(t, false, res)
}

func TestBookingsValidMaxUpcomingBookings(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user := createTestUserInOrg(org)

	router := &BookingRouter{}
	res := router.isValidMaxUpcomingBookings(org.ID, user)
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
	res := router.isValidMaxUpcomingBookings(org.ID, user)
	checkTestBool(t, false, res)
}

func TestBookingsMaxConcurrentOK(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	user3 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 3
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T13:00:00+02:00\", \"leave\": \"2030-09-01T18:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsSwitchToWinterTime(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "1000")
	user1 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	now := time.Now().UTC()
	enter := time.Date(2025, 10, 26, 0, 0, 0, 0, now.Location())
	leave := time.Date(2025, 10, 26, 23, 59, 59, 0, now.Location())
	if now.Compare(enter) > 0 {
		// Skip test
		t.Log("Skipping test TestBookingsSwitchToWinterTime")
		return
	}

	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"" + enter.Format(JsDateTimeFormatWithTimezone) + "\", \"leave\": \"" + leave.Format(JsDateTimeFormatWithTimezone) + "\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsSwitchToSummerTime(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "1000")
	user1 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	now := time.Now().UTC()
	enter := time.Date(2026, 3, 29, 0, 0, 0, 0, now.Location())
	leave := time.Date(2026, 3, 29, 23, 59, 59, 0, now.Location())
	if now.Compare(enter) > 0 {
		// Skip test
		t.Log("Skipping test TestBookingsSwitchToSummerTime")
		return
	}

	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"" + enter.Format(JsDateTimeFormatWithTimezone) + "\", \"leave\": \"" + leave.Format(JsDateTimeFormatWithTimezone) + "\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsSameDay(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingDailyBasisBooking.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxBookingDurationHours.Name, "24")
	user1 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	now := time.Now().UTC()
	enter := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	leave := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"" + enter.Format(JsDateTimeFormatWithTimezone) + "\", \"leave\": \"" + leave.Format(JsDateTimeFormatWithTimezone) + "\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsMaxConcurrentLimitExceeded(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	user3 := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 3
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
	checkTestString(t, strconv.Itoa(ResponseCodeBookingLocationMaxConcurrent), res.Header().Get("X-Error-Code"))

	// Create booking 3 as admin -> should also NOT possible
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
	checkTestString(t, strconv.Itoa(ResponseCodeBookingLocationMaxConcurrent), res.Header().Get("X-Error-Code"))
}

func TestBookingsMaxConcurrentLimitOKOnUpdate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Modify booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T10:00:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
}

func TestBookingsMaxConcurrentLimitExceededOnUpdate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	user3 := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 3
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T13:00:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Modify booking 3
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T13:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
	checkTestString(t, strconv.Itoa(ResponseCodeBookingLocationMaxConcurrent), res.Header().Get("X-Error-Code"))

	// Modify booking 3 as admin (cannot override concurrency)
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T13:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
	checkTestString(t, strconv.Itoa(ResponseCodeBookingLocationMaxConcurrent), res.Header().Get("X-Error-Code"))
}

func TestBookingsMaxConcurrentLimitExceededHeadRequest(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	user3 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)

	//	 |------------------------| #1 - OK
	//	|------------|              #2 - OK
	//	         |-----------|      #3 - NOK

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 3
	payload = "{\"locationId\": \"" + l.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/precheck/", user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
	checkTestString(t, strconv.Itoa(ResponseCodeBookingLocationMaxConcurrent), res.Header().Get("X-Error-Code"))
}

func TestBookingsMaxConcurrentLimitComplex(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	user3 := createTestUserInOrg(org)
	user4 := createTestUserInOrg(org)
	user5 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 2,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)
	s4 := &Space{Name: "Test 4", LocationID: l.ID}
	GetSpaceRepository().Create(s4)
	s5 := &Space{Name: "Test 5", LocationID: l.ID}
	GetSpaceRepository().Create(s5)

	//	|------------|                   #1 - OK  (07:30 - 12:00)
	//	                      |-----|    #2 - OK  (16:00 - 19:00)
	//	 |------------------------|      #3 - OK  (08:30 - 17:00)
	//	              |------|           #4 - OK  (12:30 - 15:30)
	//	                   |-----|       #5 - NOK (15:00 - 16:45)

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 3
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T08:30:00+02:00\", \"leave\": \"2030-09-01T17:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 4
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T12:30:00+02:00\", \"leave\": \"2030-09-01T15:30:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user4.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 5
	payload = "{\"spaceId\": \"" + s5.ID + "\", \"enter\": \"2030-09-01T15:00:00+02:00\", \"leave\": \"2030-09-01T16:45:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user5.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestBookingsMaxConcurrentLimitOKComplex(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	user3 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 1,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)

	//	|------------|                   #1 - OK  (07:30 - 12:00)
	//	                      |-----|    #2 - OK  (16:00 - 19:00)
	//	             |--------|          #3 - OK  (08:30 - 17:00)

	// Create booking 1
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 2
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Create booking 3
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T12:00:00+02:00\", \"leave\": \"2030-09-01T16:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user3.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsConvertTimestampDefaultSetting(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingDefaultTimezone.Name, "US/Central")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	user1 := createTestUserInOrg(org)

	l := &Location{
		Name:           "Test",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	// Create booking
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T09:30:00Z\", \"leave\": \"2030-09-01T12:00:00Z\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Read booking
	req = newHTTPRequest("GET", "/booking/"+id, user1.ID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "2030-09-01T09:30:00-05:00", resBody.Enter.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, "2030-09-01T12:00:00-05:00", resBody.Leave.Format(JsDateTimeFormatWithTimezone))
}

func TestBookingsConvertTimestamp(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	user1 := createTestUserInOrg(org)

	l := &Location{
		Name:           "Test",
		OrganizationID: org.ID,
		Timezone:       "US/Central",
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	// Create booking
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T09:30:00Z\", \"leave\": \"2030-09-01T12:00:00Z\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Read booking
	req = newHTTPRequest("GET", "/booking/"+id, user1.ID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "2030-09-01T09:30:00-05:00", resBody.Enter.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, "2030-09-01T12:00:00-05:00", resBody.Leave.Format(JsDateTimeFormatWithTimezone))

	// Update Booking
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T08:45:00Z\", \"leave\": \"2030-09-01T15:15:00Z\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read booking
	req = newHTTPRequest("GET", "/booking/"+id, user1.ID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetBookingResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "2030-09-01T08:45:00-05:00", resBody2.Enter.Format(JsDateTimeFormatWithTimezone))
	checkTestString(t, "2030-09-01T15:15:00-05:00", resBody2.Leave.Format(JsDateTimeFormatWithTimezone))
}

func TestBookingsPresenceReport(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user1 := createTestUserInOrgWithName(org, "u1@test.com", UserRoleUser)
	user2 := createTestUserInOrgWithName(org, "u2@test.com", UserRoleUser)
	user3 := createTestUserInOrgWithName(org, "u3@test.com", UserRoleSpaceAdmin)

	// Prepare
	l := &Location{
		Name:           "Test",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	tomorrow := time.Now().Add(24 * time.Hour)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 8, 0, 0, 0, tomorrow.Location())

	// Create booking
	b1_1 := &Booking{
		UserID:  user1.ID,
		SpaceID: s1.ID,
		Enter:   tomorrow.Add(0 * time.Hour),
		Leave:   tomorrow.Add(8 * time.Hour),
	}
	GetBookingRepository().Create(b1_1)
	b1_2 := &Booking{
		UserID:  user1.ID,
		SpaceID: s1.ID,
		Enter:   tomorrow.Add((24 + 0) * time.Hour),
		Leave:   tomorrow.Add((24 + 8) * time.Hour),
	}
	GetBookingRepository().Create(b1_2)
	b2_1 := &Booking{
		UserID:  user2.ID,
		SpaceID: s1.ID,
		Enter:   tomorrow.Add((24*2 + 0) * time.Hour),
		Leave:   tomorrow.Add((24*2 + 8) * time.Hour),
	}
	GetBookingRepository().Create(b2_1)

	end := tomorrow.Add(24 * 7 * time.Hour)
	end = time.Date(end.Year(), end.Month(), end.Day(), 8, 0, 0, 0, end.Location())
	payload := "{\"start\": \"" + tomorrow.Format(JsDateTimeFormatWithTimezone) + "\", \"end\": \"" + end.Format(JsDateTimeFormatWithTimezone) + "\"}"
	req := newHTTPRequest("POST", "/booking/report/presence/", user3.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetPresenceReportResult
	json.Unmarshal(res.Body.Bytes(), &resBody)

	checkTestInt(t, 3, len(resBody.Users))
	checkTestInt(t, 8, len(resBody.Dates))

	checkTestString(t, user1.ID, resBody.Users[0].UserID)
	checkTestString(t, user1.Email, resBody.Users[0].Email)
	checkTestString(t, user2.ID, resBody.Users[1].UserID)
	checkTestString(t, user2.Email, resBody.Users[1].Email)
	checkTestString(t, user3.ID, resBody.Users[2].UserID)
	checkTestString(t, user3.Email, resBody.Users[2].Email)

	const DateFormat string = "2006-01-02"
	checkTestString(t, tomorrow.Add(24*0*time.Hour).Format(DateFormat), resBody.Dates[0])
	checkTestString(t, tomorrow.Add(24*1*time.Hour).Format(DateFormat), resBody.Dates[1])
	checkTestString(t, tomorrow.Add(24*2*time.Hour).Format(DateFormat), resBody.Dates[2])
	checkTestString(t, tomorrow.Add(24*3*time.Hour).Format(DateFormat), resBody.Dates[3])
	checkTestString(t, tomorrow.Add(24*4*time.Hour).Format(DateFormat), resBody.Dates[4])
	checkTestString(t, tomorrow.Add(24*5*time.Hour).Format(DateFormat), resBody.Dates[5])
	checkTestString(t, tomorrow.Add(24*6*time.Hour).Format(DateFormat), resBody.Dates[6])
	checkTestString(t, tomorrow.Add(24*7*time.Hour).Format(DateFormat), resBody.Dates[7])

	checkTestInt(t, 1, resBody.Presences[0][0])
	checkTestInt(t, 1, resBody.Presences[0][1])
	checkTestInt(t, 0, resBody.Presences[0][2])
	checkTestInt(t, 0, resBody.Presences[0][3])
	checkTestInt(t, 0, resBody.Presences[0][4])
	checkTestInt(t, 0, resBody.Presences[0][5])
	checkTestInt(t, 0, resBody.Presences[0][6])
	checkTestInt(t, 0, resBody.Presences[0][7])

	checkTestInt(t, 0, resBody.Presences[1][0])
	checkTestInt(t, 0, resBody.Presences[1][1])
	checkTestInt(t, 1, resBody.Presences[1][2])
	checkTestInt(t, 0, resBody.Presences[1][3])
	checkTestInt(t, 0, resBody.Presences[1][4])
	checkTestInt(t, 0, resBody.Presences[1][5])
	checkTestInt(t, 0, resBody.Presences[1][6])
	checkTestInt(t, 0, resBody.Presences[1][7])

	checkTestInt(t, 0, resBody.Presences[2][0])
	checkTestInt(t, 0, resBody.Presences[2][1])
	checkTestInt(t, 0, resBody.Presences[2][2])
	checkTestInt(t, 0, resBody.Presences[2][3])
	checkTestInt(t, 0, resBody.Presences[2][4])
	checkTestInt(t, 0, resBody.Presences[2][5])
	checkTestInt(t, 0, resBody.Presences[2][6])
	checkTestInt(t, 0, resBody.Presences[2][7])
}

func TestBookingsUserConcurrentOk(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "50")
	GetSettingsRepository().Set(org.ID, SettingMaxConcurrentBookingsPerUser.Name, "1")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 10,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)
	s4 := &Space{Name: "Test 4", LocationID: l.ID}
	GetSpaceRepository().Create(s4)
	s5 := &Space{Name: "Test 5", LocationID: l.ID}
	GetSpaceRepository().Create(s5)

	// all with overlap

	// user one books
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user two books
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another away from first
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T20:00:00+02:00\", \"leave\": \"2030-09-01T20:25:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another with start as another ends, this should be ok
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T12:00:00+02:00\", \"leave\": \"2030-09-01T15:25:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another bordering with start time, this should be ok
	payload = "{\"spaceId\": \"" + s5.ID + "\", \"enter\": \"2030-09-01T05:00:00+02:00\", \"leave\": \"2030-09-01T07:30:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsUserConcurrentExceedLimit(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "50")
	GetSettingsRepository().Set(org.ID, SettingMaxConcurrentBookingsPerUser.Name, "2")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 10,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)
	s4 := &Space{Name: "Test 4", LocationID: l.ID}
	GetSpaceRepository().Create(s4)
	s5 := &Space{Name: "Test 5", LocationID: l.ID}
	GetSpaceRepository().Create(s5)

	// all with overlap

	// user one books
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user two books
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another with overlap
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:30:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// border start
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T06:00:00+02:00\", \"leave\": \"2030-09-01T11:40:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)

	// border end
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T11:50:00+02:00\", \"leave\": \"2030-09-01T14:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)

	// surround
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T06:00:00+02:00\", \"leave\": \"2030-09-01T13:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)

	// within
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T09:00:00+02:00\", \"leave\": \"2030-09-01T11:31:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestBookingsUserConcurrentNoLimit(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "50")
	GetSettingsRepository().Set(org.ID, SettingMaxConcurrentBookingsPerUser.Name, "0")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 10,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)
	s4 := &Space{Name: "Test 4", LocationID: l.ID}
	GetSpaceRepository().Create(s4)
	s5 := &Space{Name: "Test 5", LocationID: l.ID}
	GetSpaceRepository().Create(s5)

	// all with overlap

	// user one books
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user two books
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books some more, plenty more, no errors
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	payload = "{\"spaceId\": \"" + s5.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestBookingsUserConcurrentLimitOkOnUpdate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "50")
	GetSettingsRepository().Set(org.ID, SettingMaxConcurrentBookingsPerUser.Name, "2")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 10,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)
	s4 := &Space{Name: "Test 4", LocationID: l.ID}
	GetSpaceRepository().Create(s4)
	s5 := &Space{Name: "Test 5", LocationID: l.ID}
	GetSpaceRepository().Create(s5)

	// all with overlap

	// user one books
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user two books
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another with overlap
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:30:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another with different overlap
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T08:00:00+02:00\", \"leave\": \"2030-09-01T11:25:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// user moves last booking, still within the concurrency rules
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T07:00:00+02:00\", \"leave\": \"2030-09-01T11:15:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
}
func TestBookingsUserConcurrentLimitExceededOnUpdate(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	GetSettingsRepository().Set(org.ID, SettingMaxBookingsPerUser.Name, "50")
	GetSettingsRepository().Set(org.ID, SettingMaxConcurrentBookingsPerUser.Name, "2")
	user1 := createTestUserInOrg(org)
	user2 := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 10,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)
	s2 := &Space{Name: "Test 2", LocationID: l.ID}
	GetSpaceRepository().Create(s2)
	s3 := &Space{Name: "Test 3", LocationID: l.ID}
	GetSpaceRepository().Create(s3)
	s4 := &Space{Name: "Test 4", LocationID: l.ID}
	GetSpaceRepository().Create(s4)
	s5 := &Space{Name: "Test 5", LocationID: l.ID}
	GetSpaceRepository().Create(s5)

	// all with overlap

	// user one books
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user two books
	payload = "{\"spaceId\": \"" + s2.ID + "\", \"enter\": \"2030-09-01T16:00:00+02:00\", \"leave\": \"2030-09-01T19:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user2.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another with overlap
	payload = "{\"spaceId\": \"" + s3.ID + "\", \"enter\": \"2030-09-01T11:30:00+02:00\", \"leave\": \"2030-09-01T15:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// user one books another with different overlap
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T08:00:00+02:00\", \"leave\": \"2030-09-01T11:25:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// user moves last booking, now overlaps with 2 previous
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T13:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, user1.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
	checkTestString(t, strconv.Itoa(ResponseCodeBookingMaxConcurrentForUser), res.Header().Get("X-Error-Code"))

	// admin move last booking, now overlaps with 2 previous, but should be ok
	payload = "{\"spaceId\": \"" + s4.ID + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T13:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+id, adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
}

func TestBookingsNonExistingUsers(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	GetSettingsRepository().Set(org.ID, SettingAllowBookingsNonExistingUsers.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingNoAdminRestrictions.Name, "1")
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, strconv.Itoa(365*10))
	user := createTestUserInOrg(org)
	adminUser := createTestUserOrgAdmin(org)

	l := &Location{
		Name:                  "Test",
		MaxConcurrentBookings: 10,
		OrganizationID:        org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	// admin books
	payload := "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"noobie@test.com\", \"enter\": \"2030-09-01T07:30:00+02:00\", \"leave\": \"2030-09-01T12:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", adminUser.ID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	bookingID := res.Header().Get("X-Object-Id")

	// modify to another user
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"noobie2@test.com\", \"enter\": \"2030-09-02T11:00:00+02:00\", \"leave\": \"2030-09-02T13:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+bookingID, adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// user books
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"noobie3@test.com\", \"enter\": \"2030-09-03T07:30:00+02:00\", \"leave\": \"2030-09-03T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// user tries to overtake foreign booking
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"" + user.Email + "\", \"enter\": \"2030-09-01T11:00:00+02:00\", \"leave\": \"2030-09-01T13:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+bookingID, user.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// user books normal
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"enter\": \"2030-09-04T07:30:00+02:00\", \"leave\": \"2030-09-04T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", user.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	bookingID = res.Header().Get("X-Object-Id")

	// user tries to change to new user
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"noobie4@test.com\", \"enter\": \"2030-09-04T07:30:00+02:00\", \"leave\": \"2030-09-04T12:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+bookingID, user.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// user tries to change to existing user
	payload = "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"noobie@test.com\", \"enter\": \"2030-09-04T07:30:00+02:00\", \"leave\": \"2030-09-04T12:00:00+02:00\"}"
	req = newHTTPRequest("PUT", "/booking/"+bookingID, user.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// disallow feature
	GetSettingsRepository().Set(org.ID, SettingAllowBookingsNonExistingUsers.Name, "0")

	payload = "{\"spaceId\": \"" + s1.ID + "\", \"userEmail\": \"noobie5@test.com\", \"enter\": \"2030-09-05T07:30:00+02:00\", \"leave\": \"2030-09-05T12:00:00+02:00\"}"
	req = newHTTPRequest("POST", "/booking/", adminUser.ID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

}
