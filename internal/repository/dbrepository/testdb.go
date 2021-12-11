package dbrepository

import (
	"errors"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

func (m *testDBRepository) AllAdministrators() bool {
	return true
}

// InsertReservation inserts a new reservation record into database
func (m *testDBRepository) InsertReservation(rsvn models.Reservation) (int64, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction in database
func (m *testDBRepository) InsertRoomRestriction(rest models.RoomRestriction) error {
	return nil
}

//SearchAvailabilityByDatesAndRoomID return true if availability exists, otherwise false
func (m *testDBRepository) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepository) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms_available []models.Room
	return rooms_available, nil
}

// GetRoomByID gets room details, especially room name, by id
func (m *testDBRepository) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("#0003: attempting to return room number greater than number of rooms available")
	}
	return room, nil
}
