package soap

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/tiaguinho/gosoap"
)

// Configuración SUNAT
type SUNATConfig struct {
	RUC      string
	Usuario  string
	Clave    string
	Endpoint string
}

// Respuesta de envío
type SendResponse struct {
	Success bool   `json:"success"`
	Ticket  string `json:"ticket,omitempty"`
	CDR     string `json:"cdr,omitempty"`
	Message string `json:"message"`
}

// Envía el ZIP con sendBill (Factura, Boleta)
func SendBill(cfg SUNATConfig, fileName string, zipContent []byte) (*SendResponse, error) {
	soapClient, err := gosoap.SoapClient(cfg.Endpoint, &http.Client{})
	if err != nil {
		return nil, fmt.Errorf("error creando cliente SOAP: %w", err)
	}

	// Encode ZIP a Base64
	zipB64 := base64.StdEncoding.EncodeToString(zipContent)

	// Params de sendBill
	params := gosoap.Params{
		"fileName":    fileName,
		"contentFile": zipB64,
	}

	// Ejecutar sendBill
	res, err := soapClient.Call("sendBill", params)
	if err != nil {
		return nil, fmt.Errorf("error invocando sendBill: %w", err)
	}

	// SUNAT devuelve el CDR como Base64 en el body de la respuesta
	cdr := string(res.Body)

	return &SendResponse{
		Success: true,
		CDR:     cdr,
		Message: "Factura enviada y aceptada",
	}, nil
}

// Envía el ZIP con sendSummary (Resumen de boletas)
func SendSummary(cfg SUNATConfig, fileName string, zipContent []byte) (*SendResponse, error) {
	soapClient, err := gosoap.SoapClient(cfg.Endpoint, &http.Client{})
	if err != nil {
		return nil, fmt.Errorf("error creando cliente SOAP: %w", err)
	}

	zipB64 := base64.StdEncoding.EncodeToString(zipContent)

	params := gosoap.Params{
		"fileName":    fileName,
		"contentFile": zipB64,
	}

	res, err := soapClient.Call("sendSummary", params)
	if err != nil {
		return nil, fmt.Errorf("error invocando sendSummary: %w", err)
	}

	// El ticket viene en el body de la respuesta
	ticket := string(res.Body)

	return &SendResponse{
		Success: true,
		Ticket:  ticket,
		Message: "Resumen enviado, ticket generado",
	}, nil
}

// Consulta Ticket (getStatus)
func GetStatus(cfg SUNATConfig, ticket string) (*SendResponse, error) {
	soapClient, err := gosoap.SoapClient(cfg.Endpoint, &http.Client{})
	if err != nil {
		return nil, fmt.Errorf("error creando cliente SOAP: %w", err)
	}

	params := gosoap.Params{
		"ticket": ticket,
	}

	res, err := soapClient.Call("getStatus", params)
	if err != nil {
		return nil, fmt.Errorf("error invocando getStatus: %w", err)
	}

	// El CDR viene en el body de la respuesta
	cdr := string(res.Body)

	return &SendResponse{
		Success: true,
		CDR:     cdr,
		Message: "CDR recuperado con éxito",
	}, nil
}

func DecodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}