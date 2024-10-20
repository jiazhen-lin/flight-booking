package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jiazhen-lin/flight-booking/internal/service"
)

type flightHandler struct {
	service service.FlightService
}

func NewFlightHandler(service service.FlightService) *flightHandler {
	return &flightHandler{service: service}
}

type searchReq struct {
	DepartureAirportID int    `json:"departure_airport_id"`
	ArrivalAirportID   int    `json:"arrival_airport_id"`
	DepartureDate      string `json:"departure_date"`
	Cursor             string `json:"cursor"`
	Limit              int    `json:"limit"`
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// todo: implement service.Search

	c.JSON(http.StatusOK, searchResp{})
}

type bookReq struct {
	FlightIDs []string `json:"flight_ids"`
}

func (h *flightHandler) Book(c *gin.Context) {
	var req bookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// todo: implement service.Book

	c.Status(http.StatusNoContent)
}
