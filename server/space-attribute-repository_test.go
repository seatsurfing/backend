package main

import "testing"

func TestSpaceAttributeRepositoryCRUD(t *testing.T) {
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

	sa11, err := GetSpaceAttributeRepository().GetOne(sa1.ID)
	checkTestBool(t, true, err == nil)
	checkTestString(t, sa1.ID, sa11.ID)
	checkTestString(t, sa1.Label, sa11.Label)
	checkTestInt(t, int(sa1.Type), int(sa11.Type))
	checkTestBool(t, sa1.LocationApplicable, sa11.LocationApplicable)
	checkTestBool(t, sa1.SpaceApplicable, sa11.SpaceApplicable)

	sa21, err := GetSpaceAttributeRepository().GetOne(sa2.ID)
	checkTestBool(t, true, err == nil)
	checkTestString(t, sa2.ID, sa21.ID)
	checkTestString(t, sa2.Label, sa21.Label)
	checkTestInt(t, int(sa2.Type), int(sa21.Type))
	checkTestBool(t, sa2.LocationApplicable, sa21.LocationApplicable)
	checkTestBool(t, sa2.SpaceApplicable, sa21.SpaceApplicable)

	list, err := GetSpaceAttributeRepository().GetAll(org.ID)
	checkTestBool(t, true, err == nil)
	checkTestInt(t, 2, len(list))
	checkTestString(t, sa1.Label, list[0].Label)
	checkTestString(t, sa2.Label, list[1].Label)

	GetSpaceAttributeRepository().Delete(sa1)

	list, err = GetSpaceAttributeRepository().GetAll(org.ID)
	checkTestBool(t, true, err == nil)
	checkTestInt(t, 1, len(list))
	checkTestString(t, sa2.Label, list[0].Label)
}
