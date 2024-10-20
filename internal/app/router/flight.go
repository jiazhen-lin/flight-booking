package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jiazhen-lin/flight-booking/internal/domain"
	"github.com/jiazhen-lin/flight-booking/internal/service"
)

type flightHandler struct {
	service service.FlightService
}

func NewFlightHandler(service service.FlightService) *flightHandler {
	return &flightHandler{service: service}
}

type searchReq struct {
	DepartureAirportID int `form:"departure_airport_id"`
	ArrivalAirportID   int `form:"arrival_airport_id"`
	// format: 2024-10-21
	DepartureDate string `form:"departure_date"`
	Cursor        string `form:"cursor"`
	Limit         int    `form:"limit"`
}

type searchResp struct {
	Flights    []flightDetailResp `json:"flights"`
	NextCursor string             `json:"next_cursor"`
}

type flightDetailResp struct {
	Flights []flightResp `json:"flights"`
	Price   string       `json:"price"`
}

type flightResp struct {
	ID                 string `json:"id"`
	Number             string `json:"number"`
	DepartureAirportID int    `json:"departure_airport_id"`
	ArrivalAirportID   int    `json:"arrival_airport_id"`
	DepartureTimestamp int64  `json:"departure_timestamp"`
	DurationSeconds    int64  `json:"duration_seconds"`
	TotalSeats         int    `json:"total_seats"`
	AvailableSeats     int    `json:"available_seats"`
	Price              string `json:"price"`
}

func (h *flightHandler) Search(c *gin.Context) {
	var req searchReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	departureDate, err := time.Parse(time.DateOnly, req.DepartureDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plans, nextCursor, err := h.service.Search(c, service.SearchFilter{
		DepartureAirportID: domain.AirportID(req.DepartureAirportID),
		ArrivalAirportID:   domain.AirportID(req.ArrivalAirportID),
		DepartureDate:      departureDate,
		Cursor:             req.Cursor,
		Limit:              req.Limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := searchResp{
		Flights:    make([]flightDetailResp, len(plans)),
		NextCursor: nextCursor,
	}
	for i, plan := range plans {
		flight := flightResp{
			ID:                 plan.Flights[0].ID,
			Number:             plan.Flights[0].Number,
			DepartureAirportID: int(plan.Flights[0].DepartureAirportID),
			ArrivalAirportID:   int(plan.Flights[0].ArrivalAirportID),
			DepartureTimestamp: plan.Flights[0].DepartureTime.Unix(),
			DurationSeconds:    plan.Flights[0].DurationSeconds,
			TotalSeats:         plan.Flights[0].TotalSeats,
			AvailableSeats:     plan.Flights[0].AvailableSeats,
			Price:              plan.Flights[0].Price.String(),
		}
		resp.Flights[i] = flightDetailResp{
			Flights: []flightResp{flight},
			Price:   plan.Price.String(),
		}
	}

	c.JSON(http.StatusOK, resp)
}

type bookReq struct {
	Flights []bookFlight `json:"flights"`
}

type bookFlight struct {
	FlightID string `json:"flight_id"`
	Seats    int    `json:"seats"`
}

func (h *flightHandler) Book(c *gin.Context) {
	var req bookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Book(c, service.BookParams{
		FlightID: req.Flights[0].FlightID,
		Seats:    req.Flights[0].Seats,
	}); err != nil {
		if errors.Is(err, domain.ErrFlightNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, domain.ErrUnavailableFlightSeats) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
