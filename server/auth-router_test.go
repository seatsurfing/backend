package main

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestPasswordReset(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	user.HashedPassword = NullString(GetUserRepository().GetHashedPassword("12345678"))
	GetUserRepository().Update(user)

	// Init password reset
	payload := "{ \"email\": \"" + user.Email + "\" }"
	req := newHTTPRequest("POST", "/auth/initpwreset", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hallo "+user.Email+","))

	// Extract Confirm ID from email
	rx := regexp.MustCompile(`\\?id=(.*)?\n`)
	confirmTokens := rx.FindStringSubmatch(SendMailMockContent)
	checkTestInt(t, 2, len(confirmTokens))
	confirmID := confirmTokens[1]

	// Complete password reset
	payload = `{
			"password": "abcd1234"
		}`
	req = newHTTPRequest("POST", "/auth/pwreset/"+confirmID, "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	// Test login with old password
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"12345678\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)

	// Test login with new password
	payload = "{ \"email\": \"" + user.Email + "\", \"password\": \"abcd1234\" }"
	req = newHTTPRequest("POST", "/auth/login", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
}

func TestPasswordResetAdminHandling(t *testing.T) {
	clearTestDB()

	org := createTestOrg("test" + GetConfig().SignupDomain)
	org.ContactEmail = "test-mailer@seatsurfing.de"
	GetOrganizationRepository().Update(org)
	user := &User{
		Email:          GetConfig().SignupAdmin + "@test" + GetConfig().SignupDomain,
		OrganizationID: org.ID,
		OrgAdmin:       true,
		SuperAdmin:     false,
		HashedPassword: NullString(GetUserRepository().GetHashedPassword("12345678")),
	}
	if err := GetUserRepository().Create(user); err != nil {
		panic(err)
	}

	// Init password reset
	payload := "{ \"email\": \"" + user.Email + "\" }"
	req := newHTTPRequest("POST", "/auth/initpwreset", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "Hallo "+user.Email+","))
	checkTestBool(t, true, strings.Contains(SendMailMockContent, "To: "+org.ContactEmail+" <"+org.ContactEmail+">"))
}
