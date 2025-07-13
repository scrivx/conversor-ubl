package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func HashXML(xmlPath string) (string, error) {
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return "", err
	}

	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:]), nil
}
