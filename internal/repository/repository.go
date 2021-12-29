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

	GetAllRooms() ([]models.Room, error)

	GetAllReservations() ([]models.Reservation, error)
	GetNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(admin models.Reservation) error
	DeleteReservation(id int) error
	UpdateReservationProcessed(id int, processed uint8) error
	GetRoomRestrictionsByDate(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
	InsertRoomBlock(roomID int, startDate, endDate time.Time) error
	DeleteRoomBlock(id int) error
}
