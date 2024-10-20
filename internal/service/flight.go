package service

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/jiazhen-lin/flight-booking/internal/domain"
	"github.com/sirupsen/logrus"
)

type SearchFilter struct {
	DepartureAirportID domain.AirportID
	ArrivalAirportID   domain.AirportID
	DepartureDate      time.Time
	Cursor             string
	Limit              int
}

type BookParams struct {
	FlightID string
	Seats    int
}

type FlightService interface {
	// Search returns a list of flight plans based on the given parameters.
	Search(ctx context.Context, filter SearchFilter) (plans []domain.FlightPlan, nextCursor string, err error)
	// Book books specific flights for the given parameters.
	Book(ctx context.Context, params BookParams) error
}

type flightService struct {
	repo domain.FlightRepository
}

var _ FlightService = (*flightService)(nil)

func NewFlightService(repo domain.FlightRepository) *flightService {
	return &flightService{repo: repo}
}

func (s *flightService) Search(ctx context.Context, filter SearchFilter) ([]domain.FlightPlan, string, error) {
	departureTimeFrom := filter.DepartureDate
	departureTimeTo := departureTimeFrom.AddDate(0, 0, 1)

	// parse cursor to implement pagination
	if filter.Cursor != "" {
		cursor, err := parseCursor(filter.Cursor)
		if err != nil {
			// ignore cursor error
			logrus.Errorf("parse cursor: %+v, error: %+v", filter.Cursor, err)
		} else {
			if cursor.lastDepartureTime.After(departureTimeFrom) {
				departureTimeFrom = cursor.lastDepartureTime.Add(1 * time.Second)
			}
		}
	}

	limit := filter.Limit
	flights, err := s.repo.List(ctx, domain.ListFilter{
		DepartureAirportID: filter.DepartureAirportID,
		ArrivalAirportID:   filter.ArrivalAirportID,
		DepartureTimeFrom:  departureTimeFrom,
		DepartureTimeTo:    departureTimeTo,
		Limit:              limit + 1,
	})
	if err != nil {
		return nil, "", err
	}
	if len(flights) == 0 {
		return nil, "", nil
	}

	var nextCursor string
	if len(flights) > limit {
		cursor := cursor{
			lastDepartureTime: flights[limit-1].DepartureTime,
		}
		nextCursor = cursor.String()
		flights = flights[:limit]
	}

	plans := make([]domain.FlightPlan, len(flights))
	for i, flight := range flights {
		plans[i] = domain.FlightPlan{
			Flights: []domain.Flight{flight},
			Price:   flight.Price,
		}
	}

	return plans, nextCursor, nil
}

type cursor struct {
	lastDepartureTime time.Time
}

func (c cursor) String() string {
	params := url.Values{}
	params.Add("departure_timestamp", strconv.FormatInt(c.lastDepartureTime.Unix(), 10))
	return params.Encode()
}

func parseCursor(c string) (*cursor, error) {
	values, err := url.ParseQuery(c)
	if err != nil {
		return nil, err
	}
	departureTimestamp, err := strconv.ParseInt(values.Get("departure_timestamp"), 10, 64)
	if err != nil {
		return nil, err
	}
	return &cursor{
		lastDepartureTime: time.Unix(departureTimestamp, 0),
	}, nil
}

func (s *flightService) Book(ctx context.Context, params BookParams) error {
	if err := s.repo.Book(ctx, params.FlightID, params.Seats); err != nil {
		return err
	}
	return nil
}
