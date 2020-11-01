package main

import "testing"

func TestSpacesCount(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")

	l1 := &Location{
		OrganizationID: org.ID,
		Name:           "L1",
	}
	GetLocationRepository().Create(l1)

	s1 := &Space{
		LocationID: l1.ID,
		Name:       "S1",
	}
	GetSpaceRepository().Create(s1)
	s2 := &Space{
		LocationID: l1.ID,
		Name:       "S2",
	}
	GetSpaceRepository().Create(s2)

	res, err := GetSpaceRepository().GetCount(org.ID)
	if err != nil {
		t.Fatal(err)
	}
	checkTestInt(t, 2, res)
}
