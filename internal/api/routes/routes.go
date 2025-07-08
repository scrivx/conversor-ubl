package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/api/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Conversor UBL API v1"})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/convert", handlers.ConvertHandler)
}