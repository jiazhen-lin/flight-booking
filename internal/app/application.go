package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jiazhen-lin/flight-booking/internal/app/router"
	"github.com/jiazhen-lin/flight-booking/internal/service"
)

type Application struct {
	FlightService service.FlightService
}

func NewApplication(flightService service.FlightService) *Application {
	return &Application{FlightService: flightService}
}

func NewHTTPServer(addr string, app *Application) *http.Server {
	r := gin.Default()
	r.ContextWithFallback = true

	v1 := r.Group("/v1")

	// register routers
	v1.GET("/livez", router.Livez)

	flightHandler := router.NewFlightHandler(app.FlightService)
	flightGroup := v1.Group("/flights")
	flightGroup.GET("/search", flightHandler.Search)
	flightGroup.POST("/book", flightHandler.Book)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
