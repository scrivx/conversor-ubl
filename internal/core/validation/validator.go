package validation

import (
	"os/exec"
)

// ValidateXMLAgainstXSD valida usando xmllint en el SO.
func ValidateXMLAgainstXSD(xmlPath, xsdPath string) error {
	cmd := exec.Command("xmllint", "--noout", "--schema", xsdPath, xmlPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if len(output) > 0 {
		return nil
	}
	return nil
}