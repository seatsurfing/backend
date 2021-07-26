package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestSettingsForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceServerSharedSecret.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	payload = `[]`
	req = newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestSettingsReadPublic(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	allowedSettings := []string{
		SettingMaxBookingsPerUser.Name,
		SettingMaxDaysInAdvance.Name,
		SettingMaxBookingDurationHours.Name,
		SettingDailyBasisBooking.Name,
		SettingShowNames.Name,
	}
	forbiddenSettings := []string{
		SettingDatabaseVersion.Name,
		SettingAllowAnyUser.Name,
		SettingConfluenceServerSharedSecret.Name,
		SettingConfluenceClientID.Name,
		SettingConfluenceAnonymous.Name,
		SettingActiveSubscription.Name,
		SettingSubscriptionMaxUsers.Name,
	}

	for _, name := range allowedSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusOK, res.Code)
	}

	for _, name := range forbiddenSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusForbidden, res.Code)
	}

	req := newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, len(allowedSettings), len(resBody))
	found := 0
	for _, name := range allowedSettings {
		for _, cur := range resBody {
			if name == cur.Name {
				found++
			}
		}
	}
	checkTestInt(t, len(allowedSettings), found)
}

func TestSettingsReadAdmin(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	allowedSettings := []string{
		SettingMaxBookingsPerUser.Name,
		SettingMaxDaysInAdvance.Name,
		SettingMaxBookingDurationHours.Name,
		SettingDailyBasisBooking.Name,
		SettingShowNames.Name,
		SettingAllowAnyUser.Name,
		SettingConfluenceServerSharedSecret.Name,
		SettingConfluenceClientID.Name,
		SettingConfluenceAnonymous.Name,
		SettingActiveSubscription.Name,
		SettingSubscriptionMaxUsers.Name,
	}
	forbiddenSettings := []string{
		SettingDatabaseVersion.Name,
	}

	for _, name := range allowedSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusOK, res.Code)
	}

	for _, name := range forbiddenSettings {
		req := newHTTPRequest("GET", "/setting/"+name, loginResponse.UserID, nil)
		res := executeTestRequest(req)
		checkTestResponseCode(t, http.StatusForbidden, res.Code)
	}

	req := newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, len(allowedSettings), len(resBody))
	found := 0
	for _, name := range allowedSettings {
		for _, cur := range resBody {
			if name == cur.Name {
				found++
			}
		}
	}
	checkTestInt(t, len(allowedSettings), found)
}

func TestSettingsCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "1", resBody)

	payload = `{"value": "0"}`
	req = newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 string
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "0", resBody2)
}

func TestSettingsCRUDMany(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)
	GetDatabase().DB().Exec("TRUNCATE settings")

	payload := `[{"name": "allow_any_user", "value": "1"}, {"name": "max_bookings_per_user", "value": "5"}]`
	req := newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestInt(t, 2, len(resBody))
	checkTestString(t, SettingAllowAnyUser.Name, resBody[0].Name)
	checkTestString(t, SettingMaxBookingsPerUser.Name, resBody[1].Name)
	checkTestString(t, "1", resBody[0].Value)
	checkTestString(t, "5", resBody[1].Value)

	payload = `[{"name": "allow_any_user", "value": "0"}, {"name": "max_bookings_per_user", "value": "3"}]`
	req = newHTTPRequest("PUT", "/setting/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []GetSettingsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestInt(t, 2, len(resBody2))
	checkTestString(t, SettingAllowAnyUser.Name, resBody2[0].Name)
	checkTestString(t, SettingMaxBookingsPerUser.Name, resBody2[1].Name)
	checkTestString(t, "0", resBody2[0].Value)
	checkTestString(t, "3", resBody2[1].Value)
}

func TestSettingsInvalidName(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "1"}`
	req := newHTTPRequest("PUT", "/setting/test123", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestSettingsInvalidBool(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "2"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingAllowAnyUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSettingsInvalidInt(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	payload := `{"value": "test"}`
	req := newHTTPRequest("PUT", "/setting/"+SettingMaxBookingsPerUser.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSettingsInitiallySetAtlassianClientID(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	clientID := uuid.New().String()

	payload := "{\"value\": \"" + clientID + "\"}"
	req := newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, clientID, resBody)
}

func TestSettingsSetSameAtlassianClientID(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	clientID := uuid.New().String()

	payload := "{\"value\": \"" + clientID + "\"}"
	req := newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, clientID, resBody)
}

func TestSettingsSetNonConflictingAtlassianClientID(t *testing.T) {
	clearTestDB()
	org1 := createTestOrg("test1.com")
	user1 := createTestUserOrgAdmin(org1)
	loginResponse1 := loginTestUser(user1.ID)
	org2 := createTestOrg("test2.com")
	user2 := createTestUserOrgAdmin(org2)
	loginResponse2 := loginTestUser(user2.ID)

	clientID1 := uuid.New().String()
	clientID2 := uuid.New().String()

	payload := "{\"value\": \"" + clientID1 + "\"}"
	req := newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse1.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	payload = "{\"value\": \"" + clientID2 + "\"}"
	req = newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	var resBody string

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceClientID.Name, loginResponse1.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, clientID1, resBody)

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceClientID.Name, loginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, clientID2, resBody)
}

func TestSettingsSetConflictingAtlassianClientID(t *testing.T) {
	clearTestDB()
	org1 := createTestOrg("test1.com")
	user1 := createTestUserOrgAdmin(org1)
	loginResponse1 := loginTestUser(user1.ID)
	org2 := createTestOrg("test2.com")
	user2 := createTestUserOrgAdmin(org2)
	loginResponse2 := loginTestUser(user2.ID)

	clientID := uuid.New().String()

	payload := "{\"value\": \"" + clientID + "\"}"
	req := newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse1.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req = newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse2.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)

	var resBody string

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceClientID.Name, loginResponse1.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, clientID, resBody)

	req = newHTTPRequest("GET", "/setting/"+SettingConfluenceClientID.Name, loginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "", resBody)
}

func TestSettingsRemoveAtlassianClientID(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// Set Client ID
	clientID := uuid.New().String()
	payload := "{\"value\": \"" + clientID + "\"}"
	req := newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Create users
	user2 := createTestUserInOrgDomain(org, "test.com")
	user3 := createTestUserInOrgDomain(org, "test.com")
	user3.AtlassianID = NullString(uuid.New().String() + "@" + clientID)
	GetUserRepository().Update(user3)
	user4 := createTestUserInOrgDomain(org, clientID)
	user4.AtlassianID = NullString(user4.Email)
	GetUserRepository().Update(user4)

	// Unset Client ID
	payload = "{\"value\": \"\"}"
	req = newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Check surviving users
	users, _ := GetUserRepository().GetAll(org.ID, 10000, 0)
	checkTestInt(t, 3, len(users))
	for _, u := range users {
		checkTestString(t, "", string(u.AtlassianID))
		if !(u.Email == user.Email || u.Email == user2.Email || u.Email == user3.Email) {
			t.Errorf("got unknown email %s", u.Email)
		}
	}
}

func TestSettingsUpdateAtlassianClientID(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	// Set Client ID
	oldClientID := uuid.New().String()
	payload := "{\"value\": \"" + oldClientID + "\"}"
	req := newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Create users
	user2 := createTestUserInOrgDomain(org, "test.com")
	user3 := createTestUserInOrgDomain(org, "test.com")
	user3.AtlassianID = NullString(uuid.New().String() + "@" + oldClientID)
	GetUserRepository().Update(user3)
	user4 := createTestUserInOrgDomain(org, oldClientID)
	user4.AtlassianID = NullString(user4.Email)
	GetUserRepository().Update(user4)

	// Change Client ID
	newClientID := uuid.New().String()
	payload = "{\"value\": \"" + newClientID + "\"}"
	req = newHTTPRequest("PUT", "/setting/"+SettingConfluenceClientID.Name, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Check users
	users, _ := GetUserRepository().GetAll(org.ID, 10000, 0)
	checkTestInt(t, 4, len(users))
	checkedCases := 0
	for _, u := range users {
		if u.Email == user.Email {
			checkTestString(t, "", string(u.AtlassianID))
			checkedCases++
		} else if u.Email == user2.Email {
			checkTestString(t, "", string(u.AtlassianID))
			checkedCases++
		} else if u.Email == user3.Email {
			checkTestString(t, strings.ReplaceAll(string(user3.AtlassianID), "@"+oldClientID, "@"+newClientID), string(u.AtlassianID))
			checkedCases++
		} else if u.Email == string(u.AtlassianID) {
			checkTestString(t, strings.ReplaceAll(string(user4.AtlassianID), "@"+oldClientID, "@"+newClientID), string(u.AtlassianID))
			checkedCases++
		}
	}
	checkTestInt(t, 4, checkedCases)
}
