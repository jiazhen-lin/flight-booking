package service

import (
	"context"
	"time"

	"github.com/jiazhen-lin/flight-booking/internal/domain"
)

type SearchFilter struct {
	DepartureAirportID *domain.AirportID
	ArrivalAirportID   *domain.AirportID
	DepartureDate      time.Time
	Cursor             string
	Limit              int
}

type BookParams struct {
	FlightIDs []string
}

type FlightService interface {
	// Search returns a list of flight plans based on the given parameters.
	Search(ctx context.Context, filter SearchFilter) (plans *domain.FlightPlan, nextCursor string, err error)
	// Book books specific flights for the given parameters.
	Book(ctx context.Context, params BookParams) error
}

type flightService struct{}

func NewFlightService() FlightService {
	return &flightService{}
}

func (s *flightService) Search(ctx context.Context, filter SearchFilter) (*domain.FlightPlan, string, error) {
	return nil, "", nil
}

func (s *flightService) Book(ctx context.Context, params BookParams) error {
	return nil
}
