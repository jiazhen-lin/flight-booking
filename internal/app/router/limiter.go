package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jiazhen-lin/flight-booking/internal/service"
)

func Limiter(limiter service.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := limiter.Allow(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
