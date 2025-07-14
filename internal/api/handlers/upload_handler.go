package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/pkg/soap"
)

type UploadRequest struct {
	InvoiceID string `json:"invoice_id"`
	FileName  string `json:"file_name"` // Nuevo campo con el nombre del archivo generado
}

// POST /upload
func UploadHandler(c *gin.Context) {
	var req UploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar el nuevo nombre del archivo si está disponible, sino usar el invoiceID
	fileName := req.FileName
	if fileName == "" {
		fileName = req.InvoiceID
	}

	zipPath := filepath.Join("temp", fileName+".zip")

	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ZIP no encontrado"})
		return
	}

	cfg := soap.SUNATConfig{
		RUC:      "20123456789",
		Usuario:  "MODDATOS",
		Clave:    "moddatos",
		Endpoint: "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService?wsdl",
	}

	// Para desarrollo, usar modo de prueba por defecto
	// En producción, descomentar las líneas siguientes para usar el servicio real
	/*
	if err := soap.CheckSUNATConnectivity(cfg); err != nil {
		// Si el servicio real no está disponible, usar modo de prueba
		fmt.Printf("Servicio SUNAT no disponible, usando modo de prueba: %v\n", err)
		res, err := soap.SendBillMock(cfg, fileName+".zip", zipData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res, err := soap.SendBill(cfg, fileName+".zip", zipData)
	*/



	// Usar modo de prueba para desarrollo
	res, err := soap.SendBillMock(cfg, fileName+".zip", zipData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Si viene un CDR, guárdalo correctamente
	if res.CDR != "" {
		if err := saveCDR(fileName, res.CDR); err != nil {
			fmt.Printf("Error guardando CDR: %v\n", err)
			// No fallar la respuesta por esto, solo log el error
		}
	}

	c.JSON(http.StatusOK, res)
}

// saveCDR guarda el CDR decodificado correctamente
func saveCDR(invoiceID, cdrBase64 string) error {
	// Validar que el CDR no esté vacío
	if cdrBase64 == "" {
		return fmt.Errorf("CDR vacío")
	}

	// Decodificar el Base64
	decoded, err := base64.StdEncoding.DecodeString(cdrBase64)
	if err != nil {
		return fmt.Errorf("error decodificando CDR: %w", err)
	}

	// Crear el directorio temp si no existe
	if err := os.MkdirAll("temp", 0755); err != nil {
		return fmt.Errorf("error creando directorio temp: %w", err)
	}

	// Crear el archivo CDR
	cdrPath := filepath.Join("temp", invoiceID+"_CDR.zip")
	if err := os.WriteFile(cdrPath, decoded, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo CDR: %w", err)
	}

	fmt.Printf("CDR guardado exitosamente en: %s\n", cdrPath)
	return nil
}