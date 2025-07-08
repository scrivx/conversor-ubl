package services

type ConvertResult struct {
	DocumentID string `json:"document_id"`
	XML        string `json:"xml"`  // base64 o XML plano (por ahora)
	Hash       string `json:"hash"` // opcional
}

type ConvertRequest struct {
	DocumentType string                 `json:"document_type"`
	Data         map[string]interface{} `json:"data"`
}

type Conversor interface {
	Convert(req ConvertRequest) (*ConvertResult, error)
}