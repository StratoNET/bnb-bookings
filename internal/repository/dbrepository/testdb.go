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

// GetAdministratorByID does exactly that
func (m *testDBRepository) GetAdministratorByID(id int) (models.Administrator, error) {
	var admin models.Administrator
	return admin, nil
}

// UpdateAdministrator updates an administrator record in the database
func (m *testDBRepository) UpdateAdministrator(admin models.Administrator) error {
	return nil
}

// AuthenticateAdministrator does exactly that
func (m *testDBRepository) AuthenticateAdministrator(email, password string) (int, string, error) {
	return 1, "", nil
}

// GetAllRooms returns all rooms as a slice of models.Room
func (m *testDBRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetAllReservations returns all reservations as a slice of models.Reservation
func (m *testDBRepository) GetAllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// GetNewReservations returns only new reservations as a slice of models.Reservation
func (m *testDBRepository) GetNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// GetReservationByID returns only one reservation as a models.Reservation
func (m *testDBRepository) GetReservationByID(id int) (models.Reservation, error) {
	var r models.Reservation
	return r, nil
}

// UpdateReservation updates a reservation record in the database
func (m *testDBRepository) UpdateReservation(admin models.Reservation) error {
	return nil
}

// DeleteReservation deletes a reservation from the database by id
func (m *testDBRepository) DeleteReservation(id int) error {
	return nil
}

// UpdateReservationProcessed updates processed level of a reservation by id
func (m *testDBRepository) UpdateReservationProcessed(id int, processed uint8) error {
	return nil
}

// GetRoomRestrictionsByDate returns all rooms restrictions by room id, for a date range, as a slice of models.RoomRestriction
func (m *testDBRepository) GetRoomRestrictionsByDate(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

// InsertRoomBlock inserts an owner block restriction for a given room
func (m *testDBRepository) InsertRoomBlock(roomID int, startDate, endDate time.Time) error {
	return nil
}

// DeleteRoomBlock deletes an owner block restriction for a room by id
func (m *testDBRepository) DeleteRoomBlock(id int) error {
	return nil
}
