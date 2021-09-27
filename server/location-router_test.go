package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLocationsEmptyResult(t *testing.T) {
	clearTestDB()
	loginResponse := createLoginTestUser()

	req := newHTTPRequest("GET", "/location/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 0 {
		t.Fatalf("Expected empty array")
	}
}

func TestLocationsForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// Create
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	org2 := createTestOrg("test2.com")
	user2 := createTestUserOrgAdmin(org2)
	loginResponse2 := loginTestUser(user2.ID)

	// Get from other org
	req = newHTTPRequest("GET", "/location/"+id, loginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// Update location from other org
	payload = `{"name": "Location 2"}`
	req = newHTTPRequest("PUT", "/location/"+id, loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	// Delete location from other org
	req = newHTTPRequest("DELETE", "/location/"+id, loginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestLocationsCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/location/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetLocationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "Location 1", resBody.Name)
	checkTestString(t, "", resBody.Description)
	checkTestInt(t, 0, int(resBody.MaxConcurrentBookings))

	// 3. Update
	payload = `{"name": "Location 2", "description": "Test 123", "maxConcurrentBookings": 20}`
	req = newHTTPRequest("PUT", "/location/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/location/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetLocationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "Location 2", resBody2.Name)
	checkTestString(t, "Test 123", resBody2.Description)
	checkTestInt(t, 20, int(resBody2.MaxConcurrentBookings))

	// 4. Delete
	req = newHTTPRequest("DELETE", "/location/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/location/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestLocationsList(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"name": "Location 1"}`
	req := newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	payload = `{"name": "Location 2"}`
	req = newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	payload = `{"name": "Location 0"}`
	req = newHTTPRequest("POST", "/location/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	req = newHTTPRequest("GET", "/location/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetLocationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements")
	}
	checkTestString(t, "Location 0", resBody[0].Name)
	checkTestString(t, "Location 1", resBody[1].Name)
	checkTestString(t, "Location 2", resBody[2].Name)
}

func TestLocationsUpload(t *testing.T) {
	resp, err := http.Get("https://upload.wikimedia.org/wikipedia/commons/7/70/Claybury_Asylum%2C_first_floor_plan._Wellcome_L0023316.jpg")
	if err != nil {
		t.Fatal("Could not load example image")
	}
	checkTestResponseCode(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Could not read body from example image")
	}

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

	// Upload
	req = newHTTPRequest("POST", "/location/"+id+"/map", loginResponse.UserID, bytes.NewBuffer(data))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Get metadata
	req = newHTTPRequest("GET", "/location/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetLocationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "jpeg", resBody.MapMimeType)
	checkTestUint(t, 4895, resBody.MapWidth)
	checkTestUint(t, 3504, resBody.MapHeight)

	// Retrieve
	req = newHTTPRequest("GET", "/location/"+id+"/map", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBodyMap *GetMapResponse
	json.Unmarshal(res.Body.Bytes(), &resBodyMap)
	data2, err := base64.StdEncoding.DecodeString(resBodyMap.Data)
	if err != nil {
		t.Fatal(err)
	}
	checkTestUint(t, uint(len(data)), uint(len(data2)))
	checkTestUint(t, 0, uint(bytes.Compare(data, data2)))
}
