package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jiazhen-lin/flight-booking/internal/app/router"
	"github.com/jiazhen-lin/flight-booking/internal/service"
)

type Application struct {
	FlightService service.FlightService
}

func NewApplication(
	flightService service.FlightService,
) *Application {
	return &Application{FlightService: flightService}
}

func NewHTTPServer(addr string, app *Application) (*http.Server, error) {
	r := gin.Default()
	r.ContextWithFallback = true

	v1 := r.Group("/v1")

	// register routers
	v1.GET("/livez", router.Livez)

	flightHandler := router.NewFlightHandler(app.FlightService)
	flightGroup := v1.Group("/flights")
	flightGroup.GET("/search", flightHandler.Search)

	bookingRateLimiter, err := service.NewMemoryTokenBucketRateLimiter(service.TokenBucketConfig{
		Key:      "/v1/flights/book",
		Duration: 100 * time.Millisecond,
		Burst:    10,
	})
	if err != nil {
		return nil, err
	}
	flightGroup.POST("/book",
		router.Limiter(bookingRateLimiter),
		flightHandler.Book,
	)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}, nil
}
