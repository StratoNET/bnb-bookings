package models

import "time"

//Administrator is the administrator model
type Administrator struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RestrictionCategory is the restriction category model
type RestrictionCategory struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	RoomID    int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomResriction is the room restriction model (NB: LastInsertId() requires ReservationID as type int64)
type RoomRestriction struct {
	ID                  int
	RoomID              int
	ReservationID       int64
	RestrictionID       int
	StartDate           time.Time
	EndDate             time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Room                Room
	Reservation         Reservation
	RestrictionCategory RestrictionCategory
}

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
