package services

import (
	"bytes"
	"encoding/xml"
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
	XMLName        xml.Name      `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`
	Extensions     UBLExtensions `xml:"ext:UBLExtensions"`
	UBLInvoiceBody ubl.Invoice   `xml:""`
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
	// 1️⃣ Construir la estructura UBL base
	invoice := ubl.Invoice{
		UBLVersionID:         "2.1",
		CustomizationID:      "2.0",
		ID:                   req.Data["id"].(string),
		IssueDate:            time.Now().Format("2006-01-02"),
		DocumentCurrencyCode: "PEN",
		AccountingSupplierParty: ubl.SupplierParty{
			CustomerAssignedAccountID: req.Data["emisor_ruc"].(string),
			Party: ubl.Party{
				PartyName: []ubl.PartyName{{Name: req.Data["emisor_nombre"].(string)}},
				PartyLegalEntity: []ubl.PartyLegalEntity{{
					RegistrationName: req.Data["emisor_razon"].(string),
					CompanyID:        req.Data["emisor_ruc"].(string),
				}},
			},
		},
		AccountingCustomerParty: ubl.CustomerParty{
			CustomerAssignedAccountID: req.Data["cliente_ruc"].(string),
			Party: ubl.Party{
				PartyLegalEntity: []ubl.PartyLegalEntity{{
					RegistrationName: req.Data["cliente_razon"].(string),
					CompanyID:        req.Data["cliente_ruc"].(string),
				}},
			},
		},
		LegalMonetaryTotal: ubl.LegalMonetaryTotal{
			PayableAmount: ubl.MonetaryAmount{
				Value:      req.Data["total"].(float64),
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
				Value:      req.Data["total"].(float64),
				CurrencyID: "PEN",
			},
			Item: ubl.Item{
				Name: req.Data["item_nombre"].(string),
			},
			Price: ubl.Price{
				PriceAmount: ubl.MonetaryAmount{
					Value:      req.Data["total"].(float64),
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

	signedXML, err := signature.SignXML(buf.String(), cert, key)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Insertar firma en UBLExtensions
	wrapped := UBLInvoiceWithExtensions{
		Extensions: UBLExtensions{
			Extension: []UBLExtension{
				{
					ExtensionContent: ExtensionContent{
						XML: signedXML,
					},
				},
			},
		},
		UBLInvoiceBody: invoice,
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
