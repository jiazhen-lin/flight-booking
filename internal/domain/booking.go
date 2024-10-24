package domain

import "time"

type BookingStatus int

const (
	BookingStatusEnabled BookingStatus = iota
	BookingStatusDisabled
)

type Booking struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    BookingStatus
}

type BookingSeats struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	BookingID int64
	FlightID  int64
	Seats     int
}

type FlightSeats struct {
	FlightID int64
	Seats    int
}
