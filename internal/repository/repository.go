package repository

import (
	"time"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

type DatabaseRepository interface {
	AllAdministrators() bool

	InsertReservation(rsvn models.Reservation) (int64, error)
	InsertRoomRestriction(rest models.RoomRestriction) error
	SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
}
