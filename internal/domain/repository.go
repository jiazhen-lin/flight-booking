package domain

import (
	"context"
	"time"
)

type ListFilter struct {
	DepartureTimeFrom time.Time
	DepartureTimeTo   time.Time
	AvailableSeats    *bool
	CursorID          *int64
	Limit             int
}

type FlightRepository interface {
	List(ctx context.Context, filter ListFilter) ([]Flight, error)
	Book(ctx context.Context, seats []FlightSeats) error
}
