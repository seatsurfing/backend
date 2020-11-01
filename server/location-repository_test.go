package main

import "testing"

func TestLocationsCount(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")

	l1 := &Location{
		OrganizationID: org.ID,
		Name:           "L1",
	}
	GetLocationRepository().Create(l1)

	res, err := GetLocationRepository().GetCount(org.ID)
	if err != nil {
		t.Fatal(err)
	}
	checkTestInt(t, 1, res)
}
