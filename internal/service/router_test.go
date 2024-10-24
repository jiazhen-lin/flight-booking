package service

import (
	"testing"
	"time"

	"github.com/jiazhen-lin/flight-booking/internal/domain"
	"github.com/stretchr/testify/require"
)

var (
	flights = []domain.Flight{
		{
			// direct flight
			ID:                 1,
			DepartureAirportID: 1,
			ArrivalAirportID:   2,
			DepartureTime:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			DurationSeconds:    3600 * 10, // 10 hours
		},
		{
			ID:                 2,
			DepartureAirportID: 2,
			ArrivalAirportID:   1,
			DepartureTime:      time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			DurationSeconds:    3600 * 11,
		},
		{
			ID:                 3,
			DepartureAirportID: 1,
			ArrivalAirportID:   3,
			DepartureTime:      time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC),
			DurationSeconds:    3600 * 8,
		},
		{
			ID:                 4,
			DepartureAirportID: 3,
			ArrivalAirportID:   2,
			DepartureTime:      time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			DurationSeconds:    3600 * 4,
		},
		{
			// direct flight
			ID:                 5,
			DepartureAirportID: 1,
			ArrivalAirportID:   2,
			DepartureTime:      time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC),
			DurationSeconds:    3600 * 10, // 11 hours
		},
		{
			// direct flight but not in the same day
			ID:                 6,
			DepartureAirportID: 1,
			ArrivalAirportID:   2,
			DepartureTime:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			DurationSeconds:    3600 * 11, // 11 hours
		},
	}
)

func TestFindPaths(t *testing.T) {
	paths, err := findShortestNPaths(
		flights,
		1, 2,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		10)
	require.NoError(t, err)
	require.Equal(t, 3, len(paths))
	require.Equal(t, flights[0].ID, paths[0].flights[0].ID)
	require.Equal(t, flights[4].ID, paths[1].flights[0].ID)
	require.Equal(t, flights[2].ID, paths[2].flights[0].ID)
	require.Equal(t, flights[3].ID, paths[2].flights[1].ID)
}
