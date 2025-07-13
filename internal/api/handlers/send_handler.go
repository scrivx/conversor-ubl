package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/core/services"
)

type SendRequest struct {
	InvoiceID string `json:"invoice_id"`
	XML       string `json:"xml"`
}

// POST /send
func SendHandler(c *gin.Context) {
	var req SendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.PrepareAndValidate(req.XML, req.InvoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
