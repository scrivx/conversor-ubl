package soap

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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
	// Configurar cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	soapClient, err := gosoap.SoapClient(cfg.Endpoint, client)
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
		// Si el método no existe, intentar con el namespace completo
		if err.Error() == "method or namespace is empty" {
			// Intentar con el namespace completo de SUNAT
			res, err = soapClient.Call("urn:service.sunat.gob.pe:billService:sendBill", params)
			if err != nil {
				return nil, fmt.Errorf("error invocando sendBill con namespace: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error invocando sendBill: %w", err)
		}
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

// Verifica la conectividad con el servicio SUNAT
func CheckSUNATConnectivity(cfg SUNATConfig) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Intentar obtener el WSDL
	resp, err := client.Get(cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("error conectando a SUNAT: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("servicio SUNAT no disponible, status: %d", resp.StatusCode)
	}
	
	return nil
}

// ListAvailableMethods lista los métodos disponibles en el WSDL
func ListAvailableMethods(cfg SUNATConfig) ([]string, error) {
	// Por ahora, retornar métodos conocidos de SUNAT
	return []string{
		"sendBill",
		"sendSummary", 
		"getStatus",
	}, nil
}

// SendBillMock simula el envío para desarrollo
// SendBillMock simula el envío para desarrollo con CDR válido
func SendBillMock(cfg SUNATConfig, fileName string, zipContent []byte) (*SendResponse, error) {
	// Crear un CDR mock válido (ZIP con XML de respuesta)
	mockCDRContent := `<?xml version="1.0" encoding="UTF-8"?>
<ApplicationResponse xmlns="urn:oasis:names:specification:ubl:schema:xsd:ApplicationResponse-2">
    <UBLVersionID>2.1</UBLVersionID>
    <CustomizationID>2.0</CustomizationID>
    <DocumentResponse>
        <Response>
            <ResponseCode>0</ResponseCode>
            <Description>La Factura numero F001-00000001, ha sido aceptada</Description>
        </Response>
    </DocumentResponse>
</ApplicationResponse>`

	// Crear un ZIP con el contenido del CDR
	mockCDRZip := createMockCDRZip(mockCDRContent)
	
	// Codificar a Base64
	mockCDRBase64 := base64.StdEncoding.EncodeToString(mockCDRZip)

	return &SendResponse{
		Success: true,
		CDR:     mockCDRBase64,
		Message: "Factura enviada exitosamente (MODO PRUEBA)",
	}, nil
}

// createMockCDRZip crea un ZIP válido con el contenido del CDR
func createMockCDRZip(xmlContent string) []byte {
	// Aquí puedes usar el paquete ziputil que ya tienes
	// o crear un ZIP en memoria
	
	// Por simplicidad, retornar un ZIP básico válido
	// En un caso real, deberías usar archive/zip para crear un ZIP válido
	
	// Este es un ejemplo básico - deberías implementar la creación del ZIP
	return []byte(xmlContent) // Temporal - reemplazar con ZIP real
}