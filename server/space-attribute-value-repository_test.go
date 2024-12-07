package main

import (
	"log"
	"testing"
)

func TestSpaceAttributeValueRepositoryCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")

	sa1 := &SpaceAttribute{
		OrganizationID:     org.ID,
		Label:              "Test 123",
		Type:               SettingTypeBool,
		SpaceApplicable:    true,
		LocationApplicable: true,
	}
	err := GetSpaceAttributeRepository().Create(sa1)
	checkTestBool(t, true, err == nil)

	sa2 := &SpaceAttribute{
		OrganizationID:     org.ID,
		Label:              "Test 456",
		Type:               SettingTypeString,
		SpaceApplicable:    false,
		LocationApplicable: false,
	}
	err = GetSpaceAttributeRepository().Create(sa2)
	checkTestBool(t, true, err == nil)

	l1 := &Location{
		OrganizationID: org.ID,
		Name:           "L1",
	}
	err = GetLocationRepository().Create(l1)
	checkTestBool(t, true, err == nil)

	s1 := &Space{
		LocationID: l1.ID,
		Name:       "S1",
	}
	err = GetSpaceRepository().Create(s1)
	checkTestBool(t, true, err == nil)

	// Set
	err = GetSpaceAttributeValueRepository().Set(sa1.ID, l1.ID, SpaceAttributeValueEntityTypeLocation, "val1.1")
	checkTestBool(t, true, err == nil)

	err = GetSpaceAttributeValueRepository().Set(sa2.ID, s1.ID, SpaceAttributeValueEntityTypeSpace, "val2.1")
	checkTestBool(t, true, err == nil)

	// Update
	err = GetSpaceAttributeValueRepository().Set(sa1.ID, l1.ID, SpaceAttributeValueEntityTypeLocation, "val1.2")
	checkTestBool(t, true, err == nil)

	err = GetSpaceAttributeValueRepository().Set(sa2.ID, s1.ID, SpaceAttributeValueEntityTypeSpace, "val2.2")
	checkTestBool(t, true, err == nil)

	list1, err := GetSpaceAttributeValueRepository().GetAllForEntity(l1.ID, SpaceAttributeValueEntityTypeLocation)
	log.Println(err)
	checkTestBool(t, true, err == nil)
	checkTestInt(t, 1, len(list1))
	checkTestString(t, "val1.2", list1[0].Value)

	list2, err := GetSpaceAttributeValueRepository().GetAllForEntity(s1.ID, SpaceAttributeValueEntityTypeSpace)
	checkTestBool(t, true, err == nil)
	checkTestInt(t, 1, len(list2))
	checkTestString(t, "val2.2", list2[0].Value)
}
