package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSpacesSameOrgForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user2 := createTestUserOrgAdmin(org)
	loginResponse2 := loginTestUser(user2.ID)

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Create space
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+id+"/space/", loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	spaceID := res.Header().Get("X-Object-Id")

	// Switch to non-admin user
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	payload = `{"name": "Location 1"}`
	req = newHTTPRequest("POST", "/location/"+id+"/space/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	payload = `{"name": "Location 1"}`
	req = newHTTPRequest("PUT", "/location/"+id+"/space/"+spaceID, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("DELETE", "/location/"+id+"/space/"+spaceID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestSpacesEmptyResult(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Get spaces
	req = newHTTPRequest("GET", "/location/"+id+"/space/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestSpacesCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// 1. Create
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/location/"+locationID+"/space/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "H234", resBody.Name)
	checkTestUint(t, 50, resBody.X)
	checkTestUint(t, 100, resBody.Y)
	checkTestUint(t, 200, resBody.Width)
	checkTestUint(t, 300, resBody.Height)
	checkTestUint(t, 90, resBody.Rotation)

	// 3. Update
	payload = `{"name": "H235", "x": 51, "y": 101, "width": 201, "height": 301, "rotation": 91}`
	req = newHTTPRequest("PUT", "/location/"+locationID+"/space/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/location/"+locationID+"/space/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "H235", resBody2.Name)
	checkTestUint(t, 51, resBody2.X)
	checkTestUint(t, 101, resBody2.Y)
	checkTestUint(t, 201, resBody2.Width)
	checkTestUint(t, 301, resBody2.Height)
	checkTestUint(t, 91, resBody2.Rotation)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/location/"+locationID+"/space/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/location/"+locationID+"/space/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestSpacesList(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	locationID, _, _, _ := createTestSpaces(t, loginResponse)

	req := newHTTPRequest("GET", "/location/"+locationID+"/space/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "H234", resBody[0].Name)
	checkTestString(t, "H235", resBody[1].Name)
	checkTestString(t, "H236", resBody[2].Name)
}

func TestSpacesAvailabilityOuter(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	locationID, spaceID, _, _ := createTestSpaces(t, loginResponse)

	// Create booking
	payload := "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T06:00:00+02:00\", \"leave\": \"2030-09-01T18:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Check
	payload = `{"enter": "2030-09-01T08:30:00+02:00", "leave": "2030-09-01T17:00:00+02:00"}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/availability", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "H234", resBody[0].Name)
	checkTestString(t, "H235", resBody[1].Name)
	checkTestString(t, "H236", resBody[2].Name)
	checkTestBool(t, false, resBody[0].Available)
	checkTestBool(t, true, resBody[1].Available)
	checkTestBool(t, true, resBody[2].Available)
}

func TestSpacesAvailabilityInner(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	locationID, spaceID, _, _ := createTestSpaces(t, loginResponse)

	// Create booking
	payload := "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T09:00:00+02:00\", \"leave\": \"2030-09-01T11:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Check
	payload = `{"enter": "2020-09-01T08:30:00+02:00", "leave": "2030-09-01T17:00:00+02:00"}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/availability", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "H234", resBody[0].Name)
	checkTestString(t, "H235", resBody[1].Name)
	checkTestString(t, "H236", resBody[2].Name)
	checkTestBool(t, false, resBody[0].Available)
	checkTestBool(t, true, resBody[1].Available)
	checkTestBool(t, true, resBody[2].Available)
}

func TestSpacesAvailabilityStart(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	locationID, spaceID, _, _ := createTestSpaces(t, loginResponse)

	// Create booking
	payload := "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T07:00:00+02:00\", \"leave\": \"2030-09-01T09:00:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Check
	payload = `{"enter": "2030-09-01T08:30:00+02:00", "leave": "2030-09-01T17:00:00+02:00"}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/availability", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "H234", resBody[0].Name)
	checkTestString(t, "H235", resBody[1].Name)
	checkTestString(t, "H236", resBody[2].Name)
	checkTestBool(t, false, resBody[0].Available)
	checkTestBool(t, true, resBody[1].Available)
	checkTestBool(t, true, resBody[2].Available)
}

func TestSpacesAvailabilityEnd(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetSettingsRepository().Set(org.ID, SettingMaxDaysInAdvance.Name, "5000")

	locationID, spaceID, _, _ := createTestSpaces(t, loginResponse)

	// Create booking
	payload := "{\"spaceId\": \"" + spaceID + "\", \"enter\": \"2030-09-01T16:30:00+02:00\", \"leave\": \"2030-09-01T17:30:00+02:00\"}"
	req := newHTTPRequest("POST", "/booking/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Check
	payload = `{"enter": "2030-09-01T08:30:00+02:00", "leave": "2030-09-01T17:00:00+02:00"}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/availability", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "H234", resBody[0].Name)
	checkTestString(t, "H235", resBody[1].Name)
	checkTestString(t, "H236", resBody[2].Name)
	checkTestBool(t, false, resBody[0].Available)
	checkTestBool(t, true, resBody[1].Available)
	checkTestBool(t, true, resBody[2].Available)
}

func TestSpacesAvailabilityNoBookings(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	locationID, _, _, _ := createTestSpaces(t, loginResponse)

	payload := `{"enter": "2020-09-01T08:30:00+02:00", "leave": "2020-09-01T17:00:00+02:00"}`
	req := newHTTPRequest("POST", "/location/"+locationID+"/space/availability", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetSpaceResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "H234", resBody[0].Name)
	checkTestString(t, "H235", resBody[1].Name)
	checkTestString(t, "H236", resBody[2].Name)
	checkTestBool(t, true, resBody[0].Available)
	checkTestBool(t, true, resBody[1].Available)
	checkTestBool(t, true, resBody[2].Available)
}

func createTestSpaces(t *testing.T, loginResponse *LoginResponse) (lID, s1ID, s2ID, s3ID string) {
	// Create location
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	locationID := res.Header().Get("X-Object-Id")

	// Create #1
	payload = `{"name": "H234", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	space1ID := res.Header().Get("X-Object-Id")

	// Create #2
	payload = `{"name": "H236", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	space2ID := res.Header().Get("X-Object-Id")

	// Create #3
	payload = `{"name": "H235", "x": 50, "y": 100, "width": 200, "height": 300, "rotation": 90}`
	req = newHTTPRequest("POST", "/location/"+locationID+"/space/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	space3ID := res.Header().Get("X-Object-Id")

	return locationID, space1ID, space2ID, space3ID
}
