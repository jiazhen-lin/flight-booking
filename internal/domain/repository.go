package domain

import (
	"context"
	"time"
)

type ListFilter struct {
	DepartureAirportID AirportID
	ArrivalAirportID   AirportID
	DepartureTimeFrom  time.Time
	DepartureTimeTo    time.Time
	Limit              int
}

type FlightRepository interface {
	List(ctx context.Context, filter ListFilter) ([]Flight, error)
	Book(ctx context.Context, flightID string, seats int) error
}
