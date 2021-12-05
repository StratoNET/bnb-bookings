package repository

import "github.com/StratoNET/bnb-bookings/internal/models"

type DatabaseRepository interface {
	AllAdministrators() bool

	InsertReservation(rsvn models.Reservation) (int64, error)
	InsertRoomRestriction(rest models.RoomRestriction) error
}
