package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type AirportID int

const (
	AirportIDUnknown AirportID = iota
	AirportIDTaipei
	AirportIDTokyo
	AirportIDShanghai
	AirportIDHongKong
	AirportIDSingapore
)

type FlightStatus int

const (
	FlightStatusUnknown FlightStatus = iota
	FlightStatusEnabled
	FlightStatusDisabled
)

type Flight struct {
	ID                 string
	Number             string
	DepartureAirportID AirportID
	ArrivalAirportID   AirportID
	DepartureTime      time.Time
	DurationSeconds    int64
	Status             FlightStatus
	TotalSeats         int
	OverbookedSeats    int
	AvailableSeats     int
	Price              decimal.Decimal
}

type FlightPlan struct {
	Flights []Flight
	Price   decimal.Decimal
}
