package main

import (
	"log"
	"testing"
	"time"
)

func TestBookingRepositoryPresenceReport(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	user1 := createTestUserInOrgWithName(org, "u1@test.com", UserRoleUser)
	user2 := createTestUserInOrgWithName(org, "u2@test.com", UserRoleUser)
	user3 := createTestUserInOrgWithName(org, "u3@test.com", UserRoleUser)

	// Prepare
	l := &Location{
		Name:           "Test",
		OrganizationID: org.ID,
	}
	GetLocationRepository().Create(l)
	s1 := &Space{Name: "Test 1", LocationID: l.ID}
	GetSpaceRepository().Create(s1)

	tomorrow := time.Now().Add(24 * time.Hour)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 8, 0, 0, 0, tomorrow.Location())

	// Create booking
	b1_1 := &Booking{
		UserID:  user1.ID,
		SpaceID: s1.ID,
		Enter:   tomorrow.Add(0 * time.Hour),
		Leave:   tomorrow.Add(8 * time.Hour),
	}
	GetBookingRepository().Create(b1_1)
	b1_2 := &Booking{
		UserID:  user1.ID,
		SpaceID: s1.ID,
		Enter:   tomorrow.Add((24 + 0) * time.Hour),
		Leave:   tomorrow.Add((24 + 8) * time.Hour),
	}
	GetBookingRepository().Create(b1_2)
	b2_1 := &Booking{
		UserID:  user2.ID,
		SpaceID: s1.ID,
		Enter:   tomorrow.Add((24*2 + 0) * time.Hour),
		Leave:   tomorrow.Add((24*2 + 8) * time.Hour),
	}
	GetBookingRepository().Create(b2_1)

	end := tomorrow.Add(24 * 7 * time.Hour)
	res, err := GetBookingRepository().GetPresenceReport(org.ID, nil, tomorrow, end, 99999, 0)

	checkTestBool(t, true, err == nil)
	for _, item := range res {
		log.Println(item)
	}
	checkTestInt(t, 3, len(res))
	const DateFormat string = "2006-01-02"

	checkTestString(t, user1.Email, res[0].User.Email)
	checkTestInt(t, 1, res[0].Presence[tomorrow.Add(24*0*time.Hour).Format(DateFormat)])
	checkTestInt(t, 1, res[0].Presence[tomorrow.Add(24*1*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[0].Presence[tomorrow.Add(24*2*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[0].Presence[tomorrow.Add(24*3*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[0].Presence[tomorrow.Add(24*4*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[0].Presence[tomorrow.Add(24*5*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[0].Presence[tomorrow.Add(24*6*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[0].Presence[tomorrow.Add(24*7*time.Hour).Format(DateFormat)])

	checkTestString(t, user2.Email, res[1].User.Email)
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*0*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*1*time.Hour).Format(DateFormat)])
	checkTestInt(t, 1, res[1].Presence[tomorrow.Add(24*2*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*3*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*4*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*5*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*6*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[1].Presence[tomorrow.Add(24*7*time.Hour).Format(DateFormat)])

	checkTestString(t, user3.Email, res[2].User.Email)
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*0*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*1*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*2*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*3*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*4*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*5*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*6*time.Hour).Format(DateFormat)])
	checkTestInt(t, 0, res[2].Presence[tomorrow.Add(24*7*time.Hour).Format(DateFormat)])
}
