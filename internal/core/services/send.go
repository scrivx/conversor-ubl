package services

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrivx/conversor-ubl/internal/core/validation"
	"github.com/scrivx/conversor-ubl/internal/pkg/hashutil"
	"github.com/scrivx/conversor-ubl/internal/pkg/ziputil"
)

type SendResult struct {
	XMLFile  string `json:"xml_file"`
	ZipFile  string `json:"zip_file"`
	Hash     string `json:"hash"`
	Message  string `json:"message"`
	FileName string `json:"file_name"` // Nuevo campo con el nombre del archivo generado
}

// InvoiceInfo contiene la información extraída del XML
type InvoiceInfo struct {
	RUC    string
	TIPO   string
	SERIE  string
	NUMERO string
}

// extractInvoiceInfo extrae RUC, TIPO, SERIE y NUMERO del XML
func extractInvoiceInfo(xmlContent string) (*InvoiceInfo, error) {
	// Parsear el XML para extraer la información
	type Invoice struct {
		ID string `xml:"ID"`
		AccountingSupplierParty struct {
			CustomerAssignedAccountID string `xml:"CustomerAssignedAccountID"`
		} `xml:"AccountingSupplierParty"`
	}

	var invoice Invoice
	if err := xml.Unmarshal([]byte(xmlContent), &invoice); err != nil {
		return nil, fmt.Errorf("error parseando XML: %v", err)
	}

	// Extraer RUC del emisor
	ruc := invoice.AccountingSupplierParty.CustomerAssignedAccountID
	if ruc == "" {
		return nil, fmt.Errorf("no se encontró RUC del emisor")
	}

	// Extraer TIPO, SERIE y NUMERO del ID
	// El ID tiene formato: F001-123 (TIPO+SERIE-NUMERO)
	id := invoice.ID
	if id == "" {
		return nil, fmt.Errorf("no se encontró ID de la factura")
	}

	// Separar por el guión
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("formato de ID inválido: %s", id)
	}

	tipoSerie := parts[0]
	numero := parts[1]

	// Extraer TIPO y SERIE del primer componente
	if len(tipoSerie) < 2 {
		return nil, fmt.Errorf("formato de TIPO-SERIE inválido: %s", tipoSerie)
	}

	tipo := tipoSerie[0:1]  // Primera letra (F, B, etc.)
	serie := tipoSerie[1:]  // Resto (001, 002, etc.)

	return &InvoiceInfo{
		RUC:    ruc,
		TIPO:   tipo,
		SERIE:  serie,
		NUMERO: numero,
	}, nil
}

func PrepareAndValidate(xmlContent, invoiceID string) (*SendResult, error) {
	// Extraer información del XML para generar el nombre correcto
	info, err := extractInvoiceInfo(xmlContent)
	if err != nil {
		return nil, fmt.Errorf("error extrayendo información del XML: %v", err)
	}

	// Generar nombre del archivo con formato: RUC-TIPO-SERIE-NUMERO
	fileName := fmt.Sprintf("%s-%s-%s-%s", info.RUC, info.TIPO, info.SERIE, info.NUMERO)
	
	// Paths
	baseName := fmt.Sprintf("%s.xml", fileName)
	xmlPath := filepath.Join("temp", baseName)
	zipPath := filepath.Join("temp", fmt.Sprintf("%s.zip", fileName))
	xsdPath := "schemas/xsd/maindoc/UBL-Invoice-2.1.xsd"

	// Guardar XML
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0644); err != nil {
		return nil, fmt.Errorf("error guardando XML: %v", err)
	}

	// Validar contra XSD con fallback
	if err := validation.ValidateXMLAgainstXSDWithFallback(xmlPath, xsdPath); err != nil {
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
		FileName: fileName,
	}, nil
}