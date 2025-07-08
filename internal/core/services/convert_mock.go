package services

import (
	"encoding/xml"
	"fmt"
	"time"
)

type mockConversor struct{}

func NewMockConversor() Conversor {
	return &mockConversor{}
}

func (m *mockConversor) Convert(req ConvertRequest) (*ConvertResult, error) {
	// Generar un XML simulado
	type DummyXML struct {
		XMLName     xml.Name `xml:"Document"`
		Type        string   `xml:"Type"`
		GeneratedAt string   `xml:"GeneratedAt"`
		EmisorRUC   string   `xml:"Emisor"`
		MontoTotal  float64  `xml:"Total"`
	}

	doc := DummyXML{
		Type:        req.DocumentType,
		GeneratedAt: time.Now().Format(time.RFC3339),
		EmisorRUC:   fmt.Sprintf("%v", req.Data["emisor"]),
		MontoTotal:  parseMonto(req.Data["total"]),
	}

	xmlBytes, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	return &ConvertResult{
		DocumentID: fmt.Sprintf("%s-%d", req.DocumentType, time.Now().Unix()),
		XML:        string(xmlBytes),
		Hash:       "mocked-hash-123",
	}, nil
}

// helper básico para monto
func parseMonto(val interface{}) float64 {
	if f, ok := val.(float64); ok {
		return f
	}
	return 0.0
}