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
	Flights         []flightResp `json:"flights"`
	Price           string       `json:"price"`
	DurationSeconds int64        `json:"duration_seconds"`
}

type flightResp struct {
	ID                 int64  `json:"id"`
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	departureDate, err := time.Parse(time.DateOnly, req.DepartureDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 10 // by default
	if req.Limit > 0 {
		limit = req.Limit
	}

	paths, nextCursor, err := h.service.Search(c, service.SearchFilter{
		DepartureAirportID: domain.AirportID(req.DepartureAirportID),
		ArrivalAirportID:   domain.AirportID(req.ArrivalAirportID),
		DepartureDate:      departureDate,
		Cursor:             req.Cursor,
		Limit:              limit,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := searchResp{
		Flights:    make([]flightDetailResp, len(paths)),
		NextCursor: nextCursor,
	}
	for i, path := range paths {
		detail := flightDetailResp{
			Flights:         make([]flightResp, len(path.Flights)),
			Price:           path.Price.String(),
			DurationSeconds: path.DurationSeconds,
		}
		for j, flight := range path.Flights {
			detail.Flights[j] = flightResp{
				ID:                 flight.ID,
				Number:             flight.Number,
				DepartureAirportID: int(flight.DepartureAirportID),
				ArrivalAirportID:   int(flight.ArrivalAirportID),
				DepartureTimestamp: flight.DepartureTime.Unix(),
				DurationSeconds:    flight.DurationSeconds,
				TotalSeats:         flight.TotalSeats,
				AvailableSeats:     flight.AvailableSeats,
				Price:              flight.Price.String(),
			}
		}
		resp.Flights[i] = detail
	}

	c.JSON(http.StatusOK, resp)
}

type bookReq struct {
	Flights []flightSeats `json:"flights"`
}

type flightSeats struct {
	FlightID int64 `json:"flight_id"`
	Seats    int   `json:"seats"`
}

func (h *flightHandler) Book(c *gin.Context) {
	var req bookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := service.BookParams{
		Seats: make([]domain.FlightSeats, len(req.Flights)),
	}
	for i, f := range req.Flights {
		params.Seats[i] = domain.FlightSeats{
			FlightID: f.FlightID,
			Seats:    f.Seats,
		}
	}

	if err := h.service.Book(c, params); err != nil {
		if errors.Is(err, domain.ErrFlightNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, domain.ErrUnavailableFlightSeats) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
