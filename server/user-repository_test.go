package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestUsersCRUD(t *testing.T) {
	clearTestDB()

	// Create
	user := &User{
		Email:          uuid.New().String() + "@test.com",
		OrganizationID: "73980078-f4d7-40ff-9211-a7bcbf8d1981",
	}
	GetUserRepository().Create(user)
	checkStringNotEmpty(t, user.ID)

	// Read
	user2, err := GetUserRepository().GetOne(user.ID)
	if err != nil {
		t.Fatalf("Expected non-nil user")
	}
	checkTestString(t, user.ID, user2.ID)
	checkTestString(t, "73980078-f4d7-40ff-9211-a7bcbf8d1981", user.OrganizationID)

	// Update
	user2 = &User{
		ID:             user.ID,
		OrganizationID: "61bf23af-0310-4d2b-b401-21c31d60c2c4",
	}
	GetUserRepository().Update(user2)

	// Read
	user3, err := GetUserRepository().GetOne(user.ID)
	if err != nil {
		t.Fatalf("Expected non-nil user")
	}
	checkTestString(t, user.ID, user3.ID)
	checkTestString(t, "61bf23af-0310-4d2b-b401-21c31d60c2c4", user3.OrganizationID)

	// Delete
	GetUserRepository().Delete(user)
	_, err = GetUserRepository().GetOne(user.ID)
	if err == nil {
		t.Fatalf("Expected nil user")
	}
}

func TestUsersCount(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	createTestUserInOrg(org)
	createTestUserInOrg(org)

	res, err := GetUserRepository().GetCount(org.ID)
	if err != nil {
		t.Fatal(err)
	}
	checkTestInt(t, 2, res)
}
