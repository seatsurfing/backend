package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestOrganizationsEmptyResult(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	req := newHTTPRequest("GET", "/organization/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []string
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 1 {
		t.Fatalf("Expected array with one element (auto-created)")
	}
}

func TestOrganizationsForbidden(t *testing.T) {
	clearTestDB()
	loginResponse := createLoginTestUser()
	org := createTestOrg("testing.com")

	req := newHTTPRequest("GET", "/organization/", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("POST", "/organization/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("DELETE", "/organization/"+org.ID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("PUT", "/organization/"+org.ID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)

	req = newHTTPRequest("GET", "/organization/"+org.ID, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestOrganizationsCRUD(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// 1. Create
	payload := `{
		"name": "Some Company Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// 2. Read
	req = newHTTPRequest("GET", "/organization/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetOrganizationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "Some Company Ltd.", resBody.Name)
	checkTestString(t, "Foo", resBody.Firstname)
	checkTestString(t, "Bar", resBody.Lastname)
	checkTestString(t, "foo@seatsurfing.de", resBody.Email)
	checkTestString(t, "DE", resBody.Country)
	checkTestString(t, "de", resBody.Language)

	// 3. Update
	payload = `{
		"name": "Some Company 2 Ltd.",
		"firstname": "Foo 2",
		"lastname": "Bar 2",
		"email": "foo2@seatsurfing.de",
		"country": "AT",
		"language": "us"
	}`
	req = newHTTPRequest("PUT", "/organization/"+id, loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/organization/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetOrganizationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "Some Company 2 Ltd.", resBody2.Name)
	checkTestString(t, "Foo 2", resBody2.Firstname)
	checkTestString(t, "Bar 2", resBody2.Lastname)
	checkTestString(t, "foo2@seatsurfing.de", resBody2.Email)
	checkTestString(t, "AT", resBody2.Country)
	checkTestString(t, "us", resBody2.Language)

	// 4. Delete
	req = newHTTPRequest("DELETE", "/organization/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Read
	req = newHTTPRequest("GET", "/organization/"+id, loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestOrganizationsGetByDomain(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// Create organization
	payload := `{
		"name": "Some Company Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Add domain 1
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/test1.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Add domain 2
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/test2.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Get by domain 1 (created by super admin, so it's verified from the start)
	req = newHTTPRequest("GET", "/organization/domain/test1.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	// Verify both domains
	org, _ := GetOrganizationRepository().GetOne(id)
	GetOrganizationRepository().ActivateDomain(org, "test1.com")
	GetOrganizationRepository().ActivateDomain(org, "test2.com")

	// Get by domain 1
	req = newHTTPRequest("GET", "/organization/domain/test1.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetOrganizationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	checkTestString(t, "Some Company Ltd.", resBody.Name)

	// Get by domain 2
	req = newHTTPRequest("GET", "/organization/domain/test2.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 *GetOrganizationResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	checkTestString(t, "Some Company Ltd.", resBody.Name)

	// Get by unknown domain
	req = newHTTPRequest("GET", "/organization/domain/test3.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestOrganizationsDomainsCRUD(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// Create organization
	payload := `{
		"name": "Some Company Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	// Add domain 1
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/test1.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Add domain 2
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/test2.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Add domain 3
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/abc.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Get domain list
	req = newHTTPRequest("GET", "/organization/"+id+"/domain/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetDomainResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 3 {
		t.Fatalf("Expected array with 3 elements, got %d", len(resBody))
	}
	checkTestString(t, "abc.com", resBody[0].DomainName)
	checkTestString(t, "test1.com", resBody[1].DomainName)
	checkTestString(t, "test2.com", resBody[2].DomainName)
	checkTestBool(t, true, resBody[0].Active)
	checkTestBool(t, true, resBody[1].Active)
	checkTestBool(t, true, resBody[2].Active)

	// Remove 2
	req = newHTTPRequest("DELETE", "/organization/"+id+"/domain/test2.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Get domain list
	req = newHTTPRequest("GET", "/organization/"+id+"/domain/", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody2 []*GetDomainResponse
	json.Unmarshal(res.Body.Bytes(), &resBody2)
	if len(resBody2) != 2 {
		t.Fatalf("Expected array with 2 elements")
	}
	checkTestString(t, "abc.com", resBody[0].DomainName)
	checkTestString(t, "test1.com", resBody[1].DomainName)
	checkTestBool(t, true, resBody[0].Active)
	checkTestBool(t, true, resBody[1].Active)
}

func TestOrganizationsDomainsPreventAdminDomainDelete(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	req := newHTTPRequest("DELETE", "/organization/"+org.ID+"/domain/test.com", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestOrganizationsVerifyDNS(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// Create organization
	payload := `{
		"name": "Some Company Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id := res.Header().Get("X-Object-Id")

	org, _ := GetOrganizationRepository().GetOne(id)
	adminUser := createTestUserOrgAdmin(org)
	adminLoginResponse := loginTestUser(adminUser.ID)

	// Add domain
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/testcase.seatsurfing.de", adminLoginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Fake verify token
	GetDatabase().DB().Exec("UPDATE organizations_domains "+
		"SET verify_token = '65e51a4b-339f-4b24-b376-f9d866057b38' "+
		"WHERE domain = LOWER($1) AND organization_id = $2",
		"testcase.seatsurfing.de", id)

	// Verify domain
	req = newHTTPRequest("POST", "/organization/"+id+"/domain/testcase.seatsurfing.de/verify", adminLoginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Get domain list
	req = newHTTPRequest("GET", "/organization/"+id+"/domain/", adminLoginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody []*GetDomainResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)
	if len(resBody) != 1 {
		t.Fatalf("Expected array with 1 elements, got %d", len(resBody))
	}
	checkTestString(t, "testcase.seatsurfing.de", resBody[0].DomainName)
	checkTestBool(t, true, resBody[0].Active)
}

func TestOrganizationsAddDomainConflict(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// Create organization 1
	payload := `{
		"name": "Some Company Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id1 := res.Header().Get("X-Object-Id")

	// Create organization 2
	payload = `{
		"name": "Some Company 2 Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req = newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id2 := res.Header().Get("X-Object-Id")

	// Add domain to org 1 and activate it
	req = newHTTPRequest("POST", "/organization/"+id1+"/domain/test1.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	org1, _ := GetOrganizationRepository().GetOne(id1)
	GetOrganizationRepository().ActivateDomain(org1, "test1.com")

	// Try to add same domain to org 2
	req = newHTTPRequest("POST", "/organization/"+id2+"/domain/test1.com", loginResponse.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestOrganizationsAddDomainNoConflictBecauseInactive(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// Create organization 1
	payload := `{
		"name": "Some Company 1 Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id1 := res.Header().Get("X-Object-Id")

	// Create organization 2
	payload = `{
		"name": "Some Company 2 Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req = newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id2 := res.Header().Get("X-Object-Id")

	org1, _ := GetOrganizationRepository().GetOne(id1)
	adminUser1 := createTestUserOrgAdmin(org1)
	adminLoginResponse1 := loginTestUser(adminUser1.ID)

	org2, _ := GetOrganizationRepository().GetOne(id2)
	adminUser2 := createTestUserOrgAdmin(org2)
	adminLoginResponse2 := loginTestUser(adminUser2.ID)

	// Add domain to org 1
	req = newHTTPRequest("POST", "/organization/"+id1+"/domain/test1.com", adminLoginResponse1.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Add same domain to org 2
	req = newHTTPRequest("POST", "/organization/"+id2+"/domain/test1.com", adminLoginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
}

func TestOrganizationsAddDomainActivateConflicting(t *testing.T) {
	clearTestDB()
	user := createTestUserSuperAdmin()
	loginResponse := loginTestUser(user.ID)

	// Create organization 1
	payload := `{
		"name": "Some Company 1 Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req := newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id1 := res.Header().Get("X-Object-Id")

	// Create organization 2
	payload = `{
		"name": "Some Company 2 Ltd.",
		"firstname": "Foo",
		"lastname": "Bar",
		"email": "foo@seatsurfing.de",
		"country": "DE",
		"language": "de"
	}`
	req = newHTTPRequest("POST", "/organization/", loginResponse.UserID, bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	id2 := res.Header().Get("X-Object-Id")

	org1, _ := GetOrganizationRepository().GetOne(id1)
	adminUser1 := createTestUserOrgAdmin(org1)
	adminLoginResponse1 := loginTestUser(adminUser1.ID)

	org2, _ := GetOrganizationRepository().GetOne(id2)
	adminUser2 := createTestUserOrgAdmin(org2)
	adminLoginResponse2 := loginTestUser(adminUser2.ID)

	// Add domain to org 1
	req = newHTTPRequest("POST", "/organization/"+id1+"/domain/testcase.seatsurfing.de", adminLoginResponse1.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Add same domain to org 2
	req = newHTTPRequest("POST", "/organization/"+id2+"/domain/testcase.seatsurfing.de", adminLoginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	// Fake verify tokens
	_, err := GetDatabase().DB().Exec("UPDATE organizations_domains "+
		"SET verify_token = '65e51a4b-339f-4b24-b376-f9d866057b38' "+
		"WHERE domain = LOWER($1) AND organization_id IN ($2, $3)",
		"testcase.seatsurfing.de", id1, id2)
	if err != nil {
		t.Fatal(err)
	}

	// Activate domain in org 1
	req = newHTTPRequest("POST", "/organization/"+id1+"/domain/testcase.seatsurfing.de/verify", adminLoginResponse1.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Try to activate same domain in org 2
	req = newHTTPRequest("POST", "/organization/"+id2+"/domain/testcase.seatsurfing.de/verify", adminLoginResponse2.UserID, nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestOrganizationsDelete(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	req := newHTTPRequest("DELETE", "/organization/"+org.ID, loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Verify
	users, _ := GetUserRepository().GetAll(org.ID, 100, 0)
	checkTestInt(t, 0, len(users))

}
