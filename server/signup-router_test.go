package main

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestSignup(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foo@bar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "12345678", 
		"country": "DE",
		"language": "de",
		"acceptTerms": true
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hallo Foo Bar,"))
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "To: Foo Bar <foo@bar.com>"))

	// Extract Confirm ID from email
	rx := regexp.MustCompile(`/confirm/(.*)?\n`)
	confirmTokens := rx.FindStringSubmatch(SendMailMockContent)
	checkTestInt(t, 2, len(confirmTokens))
	confirmID := confirmTokens[1]

	// Confirm signup (Double Opt In)
	req = newHTTPRequest("POST", "/signup/confirm/"+confirmID, "", nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hallo Foo Bar,"))
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "To: foo@bar.com"))
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "admin@testorg.on.seatsurfing.local"))

	// Check if login works
	payload = `{"email": "admin@testorg.on.seatsurfing.local", "password": "12345678"}`
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	// Verify signup confirm is not possible anymore
	req = newHTTPRequest("POST", "/signup/confirm/"+confirmID, "", nil)
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestSignupLanguageEN(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foo@bar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "12345678", 
		"country": "DE",
		"language": "en",
		"acceptTerms": true
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hello Foo Bar,"))
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "To: Foo Bar <foo@bar.com>"))
}

func TestSignupCountryES(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foo@bar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "12345678", 
		"country": "ES",
		"language": "en",
		"acceptTerms": true
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hello Foo Bar,"))
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "To: Foo Bar <foo@bar.com>"))
}

func TestSignupNotAcceptTerms(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foo@bar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "12345678", 
		"country": "DE",
		"language": "de",
		"acceptTerms": false
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSignupInvalidEmail(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foobar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "12345678", 
		"country": "DE",
		"language": "de",
		"acceptTerms": true
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSignupShortPassword(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foo@bar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "123456", 
		"country": "DE",
		"language": "de",
		"acceptTerms": true
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSignupDomainConflict(t *testing.T) {
	clearTestDB()

	createTestOrg("testorg.on.seatsurfing.local")

	// Perform Signup
	payload := `{
		"firstname": "",
		"lastname": "",
		"email": "foo@bar.com", 
		"organization": "Test Org", 
		"domain": "testorg", 
		"contactFirstname": "Foo", 
		"contactLastname": "Bar", 
		"password": "12345678", 
		"country": "DE",
		"language": "de",
		"acceptTerms": true
		}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestSignupEmailConflictSignup(t *testing.T) {
	clearTestDB()

	// Perform Signup
	payload := `{
			"firstname": "",
			"lastname": "",
			"email": "foo@bar.com",
			"organization": "Test Org",
			"domain": "testorg",
			"contactFirstname": "Foo",
			"contactLastname": "Bar",
			"password": "12345678",
			"country": "DE",
			"language": "de",
			"acceptTerms": true
			}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	payload = `{
			"firstname": "",
			"lastname": "",
			"email": "foo@bar.com",
			"organization": "Test Org",
			"domain": "testorg2",
			"contactFirstname": "Foo",
			"contactLastname": "Bar",
			"password": "12345678",
			"country": "DE",
			"language": "de",
			"acceptTerms": true
			}`
	req = newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestSignupEmailConflictExistingOrg(t *testing.T) {
	clearTestDB()

	org := createTestOrg("testorg.on.seatsurfing.app")
	org.ContactEmail = "foo@bar.com"
	GetOrganizationRepository().Update(org)

	payload := `{
			"firstname": "",
			"lastname": "",
			"email": "foo@bar.com",
			"organization": "Test Org",
			"domain": "testorg2",
			"contactFirstname": "Foo",
			"contactLastname": "Bar",
			"password": "12345678",
			"country": "DE",
			"language": "de",
			"acceptTerms": true
			}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestSignupNonEuropeanCountry(t *testing.T) {
	clearTestDB()

	payload := `{
			"firstname": "",
			"lastname": "",
			"email": "foo@bar.com",
			"organization": "Test Org",
			"domain": "testorg2",
			"contactFirstname": "Foo",
			"contactLastname": "Bar",
			"password": "12345678",
			"country": "US",
			"language": "de",
			"acceptTerms": true
			}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestSignupInvalidLanguage(t *testing.T) {
	clearTestDB()

	payload := `{
			"firstname": "",
			"lastname": "",
			"email": "foo@bar.com",
			"organization": "Test Org",
			"domain": "testorg2",
			"contactFirstname": "Foo",
			"contactLastname": "Bar",
			"password": "12345678",
			"country": "DE",
			"language": "tr",
			"acceptTerms": true
			}`
	req := newHTTPRequest("POST", "/signup/", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}
