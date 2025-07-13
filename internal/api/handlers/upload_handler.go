package handlers

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/pkg/soap"
)

type UploadRequest struct {
	InvoiceID string `json:"invoice_id"`
}

// POST /upload
func UploadHandler(c *gin.Context) {
	var req UploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zipPath := filepath.Join("temp", req.InvoiceID+".zip")

	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ZIP no encontrado"})
		return
	}

	cfg := soap.SUNATConfig{
		RUC:      "20123456789",                  // <- Tu RUC
		Usuario:  "MODDATOS",                     // Usuario Secundario DEMO
		Clave:    "moddatos",                     // Clave DEMO
		Endpoint: "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService",
	}

	res, err := soap.SendBill(cfg, req.InvoiceID+".zip", zipData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Si viene un CDR, guárdalo
	if res.CDR != "" {
		cdrPath := filepath.Join("temp", req.InvoiceID+"_CDR.zip")
		cdrData, _ := os.Create(cdrPath)
		defer cdrData.Close()

		decoded, _ := base64.StdEncoding.DecodeString(res.CDR)
		cdrData.Write(decoded)
	}

	c.JSON(http.StatusOK, res)
}
