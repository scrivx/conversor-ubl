package services

import (
	"encoding/xml"
	"time"

	"github.com/scrivx/conversor-ubl/internal/pkg/signature"
	"github.com/scrivx/conversor-ubl/internal/pkg/ubl"
)

//Variable global
const dummySignature = `<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <ds:SignedInfo>
    <ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
    <ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
    <ds:Reference URI="">
      <ds:Transforms>
        <ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/>
      </ds:Transforms>
      <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
      <ds:DigestValue>dummyvalue==</ds:DigestValue>
    </ds:Reference>
  </ds:SignedInfo>
  <ds:SignatureValue>ABC123==</ds:SignatureValue>
</ds:Signature>`




type ConvertResult struct {
	// DocumentID string `json:"document_id"`
	XML string `json:"xml"` // base64 o XML plano (por ahora)
	// Hash       string `json:"hash"` // opcional
}

type ConvertRequest struct {
	DocumentType string                 `json:"document_type"`
	Data         map[string]interface{} `json:"data"`
}

type Conversor interface {
	Convert(req ConvertRequest) (*ConvertResult, error)
}

type UBLInvoiceWithExtensions struct {
	XMLName        xml.Name     `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`
	Extensions     UBLExtensions `xml:"ext:UBLExtensions"`
	UBLInvoiceBody ubl.Invoice   `xml:",inline"`
}

type UBLExtensions struct {
	Extension []UBLExtension `xml:"ext:UBLExtension"`
}

type UBLExtension struct {
	ExtensionContent ExtensionContent `xml:"ext:ExtensionContent"`
}

type ExtensionContent struct {
	XML string `xml:",innerxml"` // aquí va la firma como string XML crudo
}


func ConvertInvoice(req ConvertRequest) (*ConvertResult, error) {
	invoice := ubl.Invoice{
		UBLVersionID:         "2.1",
		CustomizationID:      "2.0",
		ID:                   req.Data["id"].(string),
		IssueDate:            time.Now().Format("2006-01-02T15:04:05"),
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

	// 1. Serializar XML sin firma
	xmlUnsigned, err := xml.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return nil, err
	}

	// 2. Cargar PFX
	cert, key, err := signature.LoadKeyPairFromPFX("certificados/C23022479065.pfx", "Ch14pp32023")
	if err != nil {
		return nil, err
	}

	// 3. Firmar
	signedXML, err := signature.SignXML(string(xmlUnsigned), cert, key)
	if err != nil {
		return nil, err
	}

	// 4. Insertar firma en UBLExtensions
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

	// 5. Serializar completo
	finalXML, err := xml.MarshalIndent(wrapped, "", "  ")
	if err != nil {
		return nil, err
	}

	return &ConvertResult{XML: xml.Header + string(finalXML)}, nil
}

