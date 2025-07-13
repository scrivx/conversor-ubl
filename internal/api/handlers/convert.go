package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/core/services"
)

// ConvertHandler expone POST /convert
func ConvertHandler(c *gin.Context) {
	var req services.ConvertRequest

	// Parsear JSON de entrada
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar tipo de documento
	if strings.ToLower(strings.TrimSpace(req.DocumentType)) != "invoice" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DocumentType debe ser 'invoice'"})
		return
	}
	// Ejecutar conversión
	result, err := services.ConvertInvoice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con XML generado
	c.JSON(http.StatusOK, gin.H{
		"xml": result.XML,
	})
}