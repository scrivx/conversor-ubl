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
	var req services.ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DocumentType != "Factura" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Tipo de documento no soportado"})
		return
	}

	result, err := services.ConvertInvoice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/xml", []byte(result.XML))
}