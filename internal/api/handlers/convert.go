package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/core/services"
)

var conversor services.Conversor = services.NewMockConversor()

type ConvertRequest struct {
	DocumentType string                 `json:"document_type" binding:"required"`
	Data         map[string]interface{} `json:"data" binding:"required"`
}

func ConvertHandler(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := conversor.Convert(services.ConvertRequest(req))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}