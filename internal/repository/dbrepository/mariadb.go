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
	VALUES (?, ?, ?, ?, ?, ?, ?);`

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

//SearchAvailabilityByDatesAndRoomID return true if availability exists, otherwise false
func (m *mariaDBRepository) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	// transaction given 3 seconds to complete, after which connection will be released
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `SELECT COUNT(id) FROM room_restrictions WHERE roomID = ? AND ? < end_date AND ? > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *mariaDBRepository) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// transaction given 3 seconds to complete, after which connection will be released
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms_available []models.Room

	query := `SELECT r.id, r.room_name FROM rooms r WHERE r.id NOT IN 
  (SELECT rr.room_id FROM room_restrictions rr WHERE ? < rr.end_date AND ? > rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms_available, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms_available, err
		}
		rooms_available = append(rooms_available, room)
	}

	if err = rows.Err(); err != nil {
		return rooms_available, err
	}

	return rooms_available, nil

}

// GetRoomByID gets room details, especially room name, by id
func (m *mariaDBRepository) GetRoomByID(id int) (models.Room, error) {
	// transaction given 3 seconds to complete, after which connection will be released
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `SELECT * FROM rooms WHERE id = ?;`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}
