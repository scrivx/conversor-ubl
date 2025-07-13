package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scrivx/conversor-ubl/internal/core/validation"
	"github.com/scrivx/conversor-ubl/internal/pkg/hashutil"
	"github.com/scrivx/conversor-ubl/internal/pkg/ziputil"
)

type SendResult struct {
	XMLFile  string `json:"xml_file"`
	ZipFile  string `json:"zip_file"`
	Hash     string `json:"hash"`
	Message  string `json:"message"`
}

func PrepareAndValidate(xmlContent, invoiceID string) (*SendResult, error) {
	// Paths
	baseName := fmt.Sprintf("%s.xml", invoiceID)
	xmlPath := filepath.Join("temp", baseName)
	zipPath := filepath.Join("temp", fmt.Sprintf("%s.zip", invoiceID))
	xsdPath := "schemas/xsd/maindoc/UBL-Invoice-2.1.xsd"

	// Guardar XML
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0644); err != nil {
		return nil, fmt.Errorf("error guardando XML: %v", err)
	}

	// Validar contra XSD
	if err := validation.ValidateXMLAgainstXSD(xmlPath, xsdPath); err != nil {
		return nil, fmt.Errorf("XML inválido: %v", err)
	}

	// Crear ZIP
	if err := ziputil.CreateZIP(xmlPath, zipPath); err != nil {
		return nil, fmt.Errorf("error creando ZIP: %v", err)
	}

	// Generar hash
	hash, err := hashutil.HashXML(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("error generando hash: %v", err)
	}

	return &SendResult{
		XMLFile: baseName,
		ZipFile: filepath.Base(zipPath),
		Hash:    hash,
		Message: "XML validado, zipeado y listo para envío a SUNAT.",
	}, nil
}