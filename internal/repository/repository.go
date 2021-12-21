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
	GetRoomByID(id int) (models.Room, error)

	GetAdministratorByID(id int) (models.Administrator, error)
	UpdateAdministrator(admin models.Administrator) error
	AuthenticateAdministrator(email, password string) (int, string, error)

	GetAllReservations() ([]models.Reservation, error)
}
