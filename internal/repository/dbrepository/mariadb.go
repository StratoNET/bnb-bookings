package dbrepository

import (
	"context"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

func (m *mariaDBRepository) AllAdministrators() bool {
	return true
}

// InsertReservation inserts a new reservation record into database
func (m *mariaDBRepository) InsertReservation(rsvn models.Reservation) (int64, error) {
	// transaction given 3 seconds to complete, after which connection will be released
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO reservations (room_id, first_name, last_name, email, phone, start_date, end_date, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	res, err := m.DB.ExecContext(ctx, stmt,
		rsvn.RoomID,
		rsvn.FirstName,
		rsvn.LastName,
		rsvn.Email,
		rsvn.Phone,
		rsvn.StartDate,
		rsvn.EndDate,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// InsertRoomRestriction inserts a room restriction in database
func (m *mariaDBRepository) InsertRoomRestriction(rest models.RoomRestriction) error {
	// transaction given 3 seconds to complete, after which connection will be released
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (room_id, reservation_id, restriction_id, start_date, end_date, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt,
		rest.RoomID,
		rest.ReservationID,
		rest.RestrictionID,
		rest.StartDate,
		rest.EndDate,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
