package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestSearchForbidden(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrg(org)
	loginResponse := loginTestUser(user.ID)

	req := newHTTPRequest("GET", "/search/test", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusForbidden, res.Code)
}

func TestSearchUsers(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	u1 := &User{
		Email:          "this.is.max@test.com",
		OrganizationID: org.ID,
	}
	GetUserRepository().Create(u1)
	u2 := &User{
		Email:          "max.it.is@test.com",
		OrganizationID: org.ID,
	}
	GetUserRepository().Create(u2)
	u3 := &User{
		Email:          "other.name@test.com",
		OrganizationID: org.ID,
	}
	GetUserRepository().Create(u3)

	req := newHTTPRequest("GET", "/search/max", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetSearchResultsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)

	checkTestInt(t, 2, len(resBody.Users))
	checkTestString(t, u2.Email, resBody.Users[0].Email)
	checkTestString(t, u1.Email, resBody.Users[1].Email)
}

func TestSearchLocations(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	l1 := &Location{
		Name:           "Frankfurt 1",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l1)
	l2 := &Location{
		Name:           "Frankfurt 2",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l2)
	l3 := &Location{
		Name:           "Berlin 1",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l3)

	req := newHTTPRequest("GET", "/search/frank", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetSearchResultsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)

	checkTestInt(t, 2, len(resBody.Locations))
	checkTestString(t, l1.Name, resBody.Locations[0].Name)
	checkTestString(t, l2.Name, resBody.Locations[1].Name)
}

func TestSearchSpaces(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserOrgAdmin(org)
	loginResponse := loginTestUser(user.ID)

	l1 := &Location{
		Name:           "Frankfurt 1",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l1)

	s1 := &Space{
		Name:       "H123",
		LocationID: l1.ID,
	}
	GetSpaceRepository().Create(s1)
	s2 := &Space{
		Name:       "H234",
		LocationID: l1.ID,
	}
	GetSpaceRepository().Create(s2)
	s3 := &Space{
		Name:       "G123",
		LocationID: l1.ID,
	}
	GetSpaceRepository().Create(s3)

	req := newHTTPRequest("GET", "/search/123", loginResponse.UserID, nil)
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	var resBody *GetSearchResultsResponse
	json.Unmarshal(res.Body.Bytes(), &resBody)

	checkTestInt(t, 2, len(resBody.Spaces))
	checkTestString(t, s3.Name, resBody.Spaces[0].Name)
	checkTestString(t, s1.Name, resBody.Spaces[1].Name)
}
