package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Livez(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
