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
		Endpoint: "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService?wsdl",
	}

	// Para desarrollo, usar modo de prueba por defecto
	// En producción, descomentar las líneas siguientes para usar el servicio real
	/*
	if err := soap.CheckSUNATConnectivity(cfg); err != nil {
		// Si el servicio real no está disponible, usar modo de prueba
		fmt.Printf("Servicio SUNAT no disponible, usando modo de prueba: %v\n", err)
		res, err := soap.SendBillMock(cfg, req.InvoiceID+".zip", zipData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res, err := soap.SendBill(cfg, req.InvoiceID+".zip", zipData)
	*/
	
	// Usar modo de prueba para desarrollo
	res, err := soap.SendBillMock(cfg, req.InvoiceID+".zip", zipData)
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
