package dbrepository

import (
	"errors"
	"log"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

func (m *testDBRepository) AllAdministrators() bool {
	return true
}

// InsertReservation inserts a new reservation record into database
func (m *testDBRepository) InsertReservation(rsvn models.Reservation) (int64, error) {
	if rsvn.RoomID == 99 {
		return 0, errors.New("insert room reservation failed")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction in database
func (m *testDBRepository) InsertRoomRestriction(rest models.RoomRestriction) error {
	if rest.RoomID == 9999 {
		return errors.New("insert room restriction failed")
	}
	return nil
}

//SearchAvailabilityByDatesAndRoomID return true if availability exists, otherwise false
func (m *testDBRepository) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	// test to fail query
	layout := "02/01/2006"
	testDate, err := time.Parse(layout, "01/01/2099")
	if err != nil {
		log.Println(err)
	}

	if start == testDate {
		return false, errors.New("SearchAvailabilityByDatesAndRoomID query failed")
	}

	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepository) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms_available []models.Room
	// if the start date is after 31/12/2098 return empty slice, indicating no rooms are available
	layout := "02/01/2006"
	date := "31/12/2098"
	dt, err := time.Parse(layout, date)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "01/01/2100")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		// no availability case, empty slice, specific fail date
		return rooms_available, errors.New("tested NO AVAILABILITY, SPECIFIC FAIL DATE case")
	}

	if start.After(dt) {
		return rooms_available, nil
	}

	// otherwise put entry into slice, indicating that some room is available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms_available = append(rooms_available, room)
	return rooms_available, nil
}

// GetRoomByID gets room details, especially room name, by id
func (m *testDBRepository) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("attempting to return room number greater than number of rooms available")
	}
	return room, nil
}
