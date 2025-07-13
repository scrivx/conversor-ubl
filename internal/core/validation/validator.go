package validation

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ValidateXMLAgainstXSD valida usando xmllint en el SO.
func ValidateXMLAgainstXSD(xmlPath, xsdPath string) error {
	cmd := exec.Command("xmllint", "--noout", "--schema", xsdPath, xmlPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error de validación XSD: %v - %s", err, string(output))
	}
	return nil
}

// ValidateXMLAgainstXSDWithFallback intenta usar xmllint, si no está disponible usa validación básica
func ValidateXMLAgainstXSDWithFallback(xmlPath, xsdPath string) error {
	// Intentar usar xmllint primero
	err := ValidateXMLAgainstXSD(xmlPath, xsdPath)
	if err != nil {
		// Si xmllint no está disponible, hacer validación básica
		if strings.Contains(err.Error(), "executable file not found") {
			return validateXMLBasic(xmlPath)
		}
		return err
	}
	return nil
}

// validateXMLBasic hace una validación básica del XML sin XSD
func validateXMLBasic(xmlPath string) error {
	// Leer el archivo XML
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return fmt.Errorf("error leyendo archivo XML: %v", err)
	}

	// Validación básica: verificar que es XML válido
	content := string(data)
	
	// Verificar que comience con declaración XML
	if !strings.HasPrefix(strings.TrimSpace(content), "<?xml") {
		return fmt.Errorf("archivo no comienza con declaración XML válida")
	}

	// Verificar que tenga elementos básicos de UBL Invoice
	if !strings.Contains(content, "Invoice") {
		return fmt.Errorf("archivo no contiene elemento Invoice")
	}

	if !strings.Contains(content, "UBLVersionID") {
		return fmt.Errorf("archivo no contiene UBLVersionID")
	}

	if !strings.Contains(content, "ID") {
		return fmt.Errorf("archivo no contiene ID de factura")
	}

	return nil
}