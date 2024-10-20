package domain

import "time"

type BookingStatus int

const (
	BookingStatusEnabled BookingStatus = iota
	BookingStatusDisabled
)

type Booking struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	FlightID  string
	Seats     int
	Status    BookingStatus
}
