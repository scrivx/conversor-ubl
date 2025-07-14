package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/scrivx/conversor-ubl/internal/pkg/signature"
	"github.com/scrivx/conversor-ubl/internal/pkg/ubl"
)

type ConvertRequest struct {
	DocumentType string                 `json:"document_type"`
	Data         map[string]interface{} `json:"data"`
}

type ConvertResult struct {
	XML string `json:"xml"`
}

type UBLInvoiceWithExtensions struct {
	XMLName    xml.Name      `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`
	Extensions UBLExtensions `xml:"ext:UBLExtensions"`
	ubl.Invoice `xml:",inline"`
}


type UBLExtensions struct {
	Extension []UBLExtension `xml:"ext:UBLExtension"`
}

type UBLExtension struct {
	ExtensionContent ExtensionContent `xml:"ext:ExtensionContent"`
}

type ExtensionContent struct {
	XML string `xml:",innerxml"`
}

func ConvertInvoice(req ConvertRequest) (*ConvertResult, error) {
	// Validar campos requeridos
	if req.Data["id"] == nil || req.Data["emisor_ruc"] == nil || req.Data["emisor_nombre"] == nil ||
		req.Data["emisor_razon"] == nil || req.Data["cliente_ruc"] == nil || req.Data["cliente_razon"] == nil ||
		req.Data["item_nombre"] == nil || req.Data["total"] == nil {
		return nil, fmt.Errorf("todos los campos son requeridos")
	}

	// Validar tipos de datos
	id, ok := req.Data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'id' debe ser string")
	}
	emisorRuc, ok := req.Data["emisor_ruc"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'emisor_ruc' debe ser string")
	}
	emisorNombre, ok := req.Data["emisor_nombre"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'emisor_nombre' debe ser string")
	}
	emisorRazon, ok := req.Data["emisor_razon"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'emisor_razon' debe ser string")
	}
	clienteRuc, ok := req.Data["cliente_ruc"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'cliente_ruc' debe ser string")
	}
	clienteRazon, ok := req.Data["cliente_razon"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'cliente_razon' debe ser string")
	}
	itemNombre, ok := req.Data["item_nombre"].(string)
	if !ok {
		return nil, fmt.Errorf("campo 'item_nombre' debe ser string")
	}
	total, ok := req.Data["total"].(float64)
	if !ok {
		return nil, fmt.Errorf("campo 'total' debe ser número")
	}

	// 1️⃣ Construir la estructura UBL base
	invoice := ubl.Invoice{
		UBLVersionID:         "2.1",
		CustomizationID:      "2.0",
		ID:                   id,
		IssueDate:            time.Now().Format("2006-01-02"),
		DocumentCurrencyCode: "PEN",
		AccountingSupplierParty: ubl.SupplierParty{
			CustomerAssignedAccountID: emisorRuc,
			Party: ubl.Party{
				PartyName: []ubl.PartyName{{Name: emisorNombre}},
				PartyLegalEntity: []ubl.PartyLegalEntity{{
					RegistrationName: emisorRazon,
					CompanyID:        emisorRuc,
				}},
			},
		},
		AccountingCustomerParty: ubl.CustomerParty{
			CustomerAssignedAccountID: clienteRuc,
			Party: ubl.Party{
				PartyLegalEntity: []ubl.PartyLegalEntity{{
					RegistrationName: clienteRazon,
					CompanyID:        clienteRuc,
				}},
			},
		},
		LegalMonetaryTotal: ubl.LegalMonetaryTotal{
			PayableAmount: ubl.MonetaryAmount{
				Value:      total,
				CurrencyID: "PEN",
			},
		},
		InvoiceLines: []ubl.InvoiceLine{{
			ID: "1",
			InvoicedQuantity: ubl.Quantity{
				Value:    1,
				UnitCode: "NIU",
			},
			LineExtensionAmount: ubl.MonetaryAmount{
				Value:      total,
				CurrencyID: "PEN",
			},
			Item: ubl.Item{
				Name: itemNombre,
			},
			Price: ubl.Price{
				PriceAmount: ubl.MonetaryAmount{
					Value:      total,
					CurrencyID: "PEN",
				},
			},
		}},
	}

	// 2️⃣ Serializar sin firma
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	xmlEncoder := xml.NewEncoder(&buf)
	xmlEncoder.Indent("", "  ")
	if err := xmlEncoder.Encode(invoice); err != nil {
		return nil, err
	}

	// 3️⃣ Cargar el .pfx y firmar con xmlsig
	cert, key, err := signature.LoadKeyPairFromPFX("certificados/C23022479065.pfx", "CRIV123")
	if err != nil {
		return nil, err
	}

	// Usar la nueva función que genera la firma como elemento estructurado
	signatureElement, err := signature.SignXMLAsElement(buf.String(), cert, key)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Insertar firma en UBLExtensions
	wrapped := UBLInvoiceWithExtensions{
		Extensions: UBLExtensions{
			Extension: []UBLExtension{
				{
					ExtensionContent: ExtensionContent{
						XML: signatureElement,
					},
				},
			},
		},
		Invoice: invoice,
	}

	// 5️⃣ Serializar final
	var out bytes.Buffer
	out.WriteString(xml.Header)
	encoder := xml.NewEncoder(&out)
	encoder.Indent("", "  ")
	if err := encoder.Encode(wrapped); err != nil {
		return nil, err
	}

	return &ConvertResult{XML: out.String()}, nil
}
