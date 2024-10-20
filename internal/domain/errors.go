package domain

import "errors"

var (
	ErrFlightNotFound         = errors.New("flight not found")
	ErrUnavailableFlightSeats = errors.New("unavailable flight seats")
)
