package main

import "testing"

func TestAuthAttemptRepositoryBanSimple(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrgWithName(org, "u1@test.com", UserRoleUser)

	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 1
	if err := GetAuthAttemptRepository().RecordLoginAttempt(user, false); err != nil {
		t.Error(err)
	}
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 2
	if err := GetAuthAttemptRepository().RecordLoginAttempt(user, false); err != nil {
		t.Error(err)
	}
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 3
	if err := GetAuthAttemptRepository().RecordLoginAttempt(user, false); err != nil {
		t.Error(err)
	}
	checkTestBool(t, true, authAttemptRepositoryIsUserDisabled(t, user.ID))
}

func TestAuthAttemptRepositoryBanWithSuccess(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user := createTestUserInOrgWithName(org, "u1@test.com", UserRoleUser)

	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 1
	GetAuthAttemptRepository().RecordLoginAttempt(user, false)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 2
	GetAuthAttemptRepository().RecordLoginAttempt(user, false)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Successful Login
	GetAuthAttemptRepository().RecordLoginAttempt(user, true)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 1
	GetAuthAttemptRepository().RecordLoginAttempt(user, false)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 2
	GetAuthAttemptRepository().RecordLoginAttempt(user, false)
	checkTestBool(t, false, authAttemptRepositoryIsUserDisabled(t, user.ID))

	// Attempt 3
	GetAuthAttemptRepository().RecordLoginAttempt(user, false)
	checkTestBool(t, true, authAttemptRepositoryIsUserDisabled(t, user.ID))
}

func authAttemptRepositoryIsUserDisabled(t *testing.T, userID string) bool {
	user, err := GetUserRepository().GetOne(userID)
	if err != nil {
		t.Error(err)
	}
	return user.Disabled
}
