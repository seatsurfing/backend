package main

import (
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
)

type SpaceRepository struct {
}

type Space struct {
	ID         string
	LocationID string
	Name       string
	X          uint
	Y          uint
	Width      uint
	Height     uint
	Rotation   uint
}

type SpaceAvailabilityBookingEntry struct {
	BookingID string
	UserID    string
	UserEmail string
	Enter     time.Time
	Leave     time.Time
}

type SpaceAvailability struct {
	Space
	Available bool
	Bookings  []*SpaceAvailabilityBookingEntry
}

type SpaceDetails struct {
	Location Location
	Space
}

var spaceRepository *SpaceRepository
var spaceRepositoryOnce sync.Once

func GetSpaceRepository() *SpaceRepository {
	spaceRepositoryOnce.Do(func() {
		spaceRepository = &SpaceRepository{}
		_, err := GetDatabase().DB().Exec("CREATE TABLE IF NOT EXISTS spaces (" +
			"id uuid DEFAULT uuid_generate_v4(), " +
			"location_id uuid NOT NULL, " +
			"name VARCHAR NOT NULL, " +
			"x INTEGER, " +
			"y INTEGER, " +
			"width INTEGER, " +
			"height INTEGER, " +
			"rotation INTEGER, " +
			"PRIMARY KEY (id))")
		if err != nil {
			panic(err)
		}
		_, err = GetDatabase().DB().Exec("CREATE INDEX IF NOT EXISTS idx_spaces_location_id ON spaces(location_id)")
		if err != nil {
			panic(err)
		}
	})
	return spaceRepository
}

func (r *SpaceRepository) RunSchemaUpgrade(curVersion, targetVersion int) {
	// No updates yet
}

func (r *SpaceRepository) Create(e *Space) error {
	var id string
	err := GetDatabase().DB().QueryRow("INSERT INTO spaces "+
		"(name, location_id, x, y, width, height, rotation) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) "+
		"RETURNING id",
		e.Name, e.LocationID, e.X, e.Y, e.Width, e.Height, e.Rotation).Scan(&id)
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (r *SpaceRepository) GetOne(id string) (*Space, error) {
	e := &Space{}
	err := GetDatabase().DB().QueryRow("SELECT id, location_id, name, x, y, width, height, rotation "+
		"FROM spaces "+
		"WHERE id = $1",
		id).Scan(&e.ID, &e.LocationID, &e.Name, &e.X, &e.Y, &e.Width, &e.Height, &e.Rotation)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *SpaceRepository) GetAllInTime(locationID string, enter, leave time.Time) ([]*SpaceAvailability, error) {
	var result []*SpaceAvailability
	subQueryWhere := "bookings.space_id = spaces.id AND (" +
		"($1 BETWEEN bookings.enter_time AND bookings.leave_time) OR " +
		"($2 BETWEEN bookings.enter_time AND bookings.leave_time) OR " +
		"(bookings.enter_time BETWEEN $1 AND $2) OR " +
		"(bookings.leave_time BETWEEN $1 AND $2)" +
		")"
	rows, err := GetDatabase().DB().Query("SELECT id, location_id, name, x, y, width, height, rotation, "+
		"NOT EXISTS(SELECT id FROM bookings WHERE "+subQueryWhere+"), "+
		"ARRAY(SELECT CONCAT(users.id, '@@@', users.email, '@@@', bookings.enter_time, '@@@', bookings.leave_time, '@@@', bookings.id) FROM bookings INNER JOIN users ON users.id = bookings.user_id WHERE "+subQueryWhere+" ORDER BY bookings.enter_time ASC) "+
		"FROM spaces "+
		"WHERE location_id = $3 "+
		"ORDER BY name", enter, leave, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &SpaceAvailability{}
		var bookingUserNames []string
		err = rows.Scan(&e.ID, &e.LocationID, &e.Name, &e.X, &e.Y, &e.Width, &e.Height, &e.Rotation, &e.Available, pq.Array(&bookingUserNames))
		for _, bookingUserName := range bookingUserNames {
			tokens := strings.Split(bookingUserName, "@@@")
			timeFormat := "2006-01-02 15:04:05"
			enter, _ := time.Parse(timeFormat, tokens[2])
			leave, _ := time.Parse(timeFormat, tokens[3])
			entry := &SpaceAvailabilityBookingEntry{
				BookingID: tokens[4],
				UserID:    tokens[0],
				UserEmail: tokens[1],
				Enter:     enter,
				Leave:     leave,
			}
			e.Bookings = append(e.Bookings, entry)
		}
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *SpaceRepository) GetByKeyword(organizationID string, keyword string) ([]*Space, error) {
	var result []*Space
	rows, err := GetDatabase().DB().Query("SELECT spaces.id, spaces.location_id, spaces.name, spaces.x, spaces.y, spaces.width, spaces.height, spaces.rotation "+
		"FROM spaces "+
		"INNER JOIN locations ON locations.id = spaces.location_id "+
		"WHERE locations.organization_id = $1 AND LOWER(spaces.name) LIKE '%' || $2 || '%'"+
		"ORDER BY spaces.name", organizationID, strings.ToLower(keyword))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &Space{}
		err = rows.Scan(&e.ID, &e.LocationID, &e.Name, &e.X, &e.Y, &e.Width, &e.Height, &e.Rotation)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *SpaceRepository) GetAll(locationID string) ([]*Space, error) {
	var result []*Space
	rows, err := GetDatabase().DB().Query("SELECT id, location_id, name, x, y, width, height, rotation "+
		"FROM spaces "+
		"WHERE location_id = $1 "+
		"ORDER BY name", locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &Space{}
		err = rows.Scan(&e.ID, &e.LocationID, &e.Name, &e.X, &e.Y, &e.Width, &e.Height, &e.Rotation)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}
func (r *SpaceRepository) Update(e *Space) error {
	_, err := GetDatabase().DB().Exec("UPDATE spaces SET "+
		"location_id = $1, "+
		"name = $2, "+
		"x = $3, "+
		"y = $4, "+
		"width = $5, "+
		"height = $6, "+
		"rotation = $7 "+
		"WHERE id = $8",
		e.LocationID, e.Name, e.X, e.Y, e.Width, e.Height, e.Rotation, e.ID)
	return err
}

func (r *SpaceRepository) Delete(e *Space) error {
	// if _, err := GetDatabase().DB().Exec("DELETE FROM bookings WHERE bookings.space_id = $1", e.ID); err != nil {
	// 	return err
	// }
	_, err := GetDatabase().DB().Exec("DELETE FROM spaces WHERE id = $1", e.ID)
	return err
}

func (r *SpaceRepository) GetCount(organizationID string) (int, error) {
	var res int
	err := GetDatabase().DB().QueryRow("SELECT COUNT(spaces.id) "+
		"FROM spaces "+
		"INNER JOIN locations ON locations.id = spaces.location_id "+
		"WHERE locations.organization_id = $1",
		organizationID).Scan(&res)
	return res, err
}
