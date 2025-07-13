package ziputil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateZIP(xmlPath, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	w, err := zipWriter.Create(filepath.Base(xmlPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, xmlFile)
	return err
}
