package ubl

import (
	"encoding/xml"
)

// Invoice represents the root UBL Invoice element
type Invoice struct {
	XMLName                xml.Name                `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`
	UBLVersionID           string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 UBLVersionID,omitempty"`
	CustomizationID        string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CustomizationID,omitempty"`
	ProfileID              string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ProfileID,omitempty"`
	ID                     string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	CopyIndicator          bool                    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CopyIndicator,omitempty"`
	UUID                   string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 UUID,omitempty"`
	IssueDate              string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IssueDate"`
	IssueTime              string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IssueTime,omitempty"`
	DueDate                string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 DueDate,omitempty"`
	InvoiceTypeCode        string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 InvoiceTypeCode,omitempty"`
	Note                   []string                `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Note,omitempty"`
	TaxPointDate           string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxPointDate,omitempty"`
	DocumentCurrencyCode   string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 DocumentCurrencyCode,omitempty"`
	TaxCurrencyCode        string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxCurrencyCode,omitempty"`
	PricingCurrencyCode    string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PricingCurrencyCode,omitempty"`
	PaymentCurrencyCode    string                  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PaymentCurrencyCode,omitempty"`
	LineCountNumeric       int                     `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 LineCountNumeric,omitempty"`
	AccountingSupplierParty SupplierParty          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 AccountingSupplierParty"`
	AccountingCustomerParty CustomerParty          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 AccountingCustomerParty"`
	PaymentMeans           []PaymentMeans          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PaymentMeans,omitempty"`
	PaymentTerms           []PaymentTerms          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PaymentTerms,omitempty"`
	TaxTotal               []TaxTotal              `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxTotal,omitempty"`
	LegalMonetaryTotal     LegalMonetaryTotal      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 LegalMonetaryTotal"`
	InvoiceLines           []InvoiceLine           `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 InvoiceLine"`
	AllowanceCharge []AllowanceCharge `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 AllowanceCharge,omitempty"`

}


// AllowanceCharge represents a discount or surcharge
type AllowanceCharge struct {
	ChargeIndicator           bool           `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ChargeIndicator"`
	AllowanceChargeReasonCode string         `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 AllowanceChargeReasonCode,omitempty"`
	Amount                    MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Amount"`
}


// SupplierParty represents the supplier party information
type SupplierParty struct {
	CustomerAssignedAccountID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CustomerAssignedAccountID,omitempty"`
	SupplierAssignedAccountID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 SupplierAssignedAccountID,omitempty"`
	Party                     Party  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Party"`
}

// CustomerParty represents the customer party information
type CustomerParty struct {
	CustomerAssignedAccountID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CustomerAssignedAccountID,omitempty"`
	SupplierAssignedAccountID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 SupplierAssignedAccountID,omitempty"`
	Party                     Party  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Party"`
}

// Party represents party information
type Party struct {
	EndpointID       string            `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 EndpointID,omitempty"`
	PartyIdentification []PartyIdentification `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyIdentification,omitempty"`
	PartyName        []PartyName       `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyName,omitempty"`
	PostalAddress    *PostalAddress    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PostalAddress,omitempty"`
	PartyTaxScheme   []PartyTaxScheme  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyTaxScheme,omitempty"`
	PartyLegalEntity []PartyLegalEntity `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyLegalEntity,omitempty"`
	Contact          *Contact          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Contact,omitempty"`
}

// PartyIdentification represents party identification
type PartyIdentification struct {
	ID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
}

// PartyName represents party name
type PartyName struct {
	Name string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
}

// PostalAddress represents postal address
type PostalAddress struct {
	StreetName           string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 StreetName,omitempty"`
	AdditionalStreetName string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 AdditionalStreetName,omitempty"`
	CityName             string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CityName,omitempty"`
	PostalZone           string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PostalZone,omitempty"`
	CountrySubentity     string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CountrySubentity,omitempty"`
	Country              *Country `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Country,omitempty"`
}

// Country represents country information
type Country struct {
	IdentificationCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IdentificationCode,omitempty"`
	Name               string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
}

// PartyTaxScheme represents party tax scheme
type PartyTaxScheme struct {
	RegistrationName string    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 RegistrationName,omitempty"`
	CompanyID        string    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CompanyID,omitempty"`
	TaxScheme        TaxScheme `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxScheme"`
}

// PartyLegalEntity represents party legal entity
type PartyLegalEntity struct {
	RegistrationName string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 RegistrationName,omitempty"`
	CompanyID        string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CompanyID,omitempty"`
}

// Contact represents contact information
type Contact struct {
	Name             string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
	Telephone        string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Telephone,omitempty"`
	ElectronicMail   string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ElectronicMail,omitempty"`
}

// PaymentMeans represents payment means
type PaymentMeans struct {
	ID                    string                 `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID,omitempty"`
	PaymentMeansCode      string                 `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PaymentMeansCode"`
	PaymentDueDate        string                 `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PaymentDueDate,omitempty"`
	PaymentChannelCode    string                 `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PaymentChannelCode,omitempty"`
	InstructionID         string                 `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 InstructionID,omitempty"`
	InstructionNote       []string               `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 InstructionNote,omitempty"`
	PayeeFinancialAccount *FinancialAccount      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PayeeFinancialAccount,omitempty"`
}

// PaymentTerms represents payment terms
type PaymentTerms struct {
	ID           string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID,omitempty"`
	PaymentMeansID []string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PaymentMeansID,omitempty"`
	PrepaidPaymentReferenceID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PrepaidPaymentReferenceID,omitempty"`
	Note         []string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Note,omitempty"`
}

// FinancialAccount represents financial account
type FinancialAccount struct {
	ID                string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID,omitempty"`
	Name              string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
	AccountTypeCode   string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 AccountTypeCode,omitempty"`
	FinancialInstitutionBranch *FinancialInstitutionBranch `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 FinancialInstitutionBranch,omitempty"`
}

// FinancialInstitutionBranch represents financial institution branch
type FinancialInstitutionBranch struct {
	ID   string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID,omitempty"`
	Name string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
}

// TaxTotal represents tax total
type TaxTotal struct {
	TaxAmount    MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxAmount"`
	TaxSubtotal  []TaxSubtotal  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxSubtotal,omitempty"`
}

// TaxSubtotal represents tax subtotal
type TaxSubtotal struct {
	TaxableAmount MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxableAmount"`
	TaxAmount     MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxAmount"`
	TaxCategory   TaxCategory    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxCategory"`
}

// TaxCategory represents tax category
type TaxCategory struct {
	ID               string    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID,omitempty"`
	Name             string    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
	Percent          float64   `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Percent,omitempty"`
	TaxExemptionReasonCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxExemptionReasonCode,omitempty"`
	TaxExemptionReason []string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxExemptionReason,omitempty"`
	TaxScheme        TaxScheme `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxScheme"`
}

// TaxScheme represents tax scheme
type TaxScheme struct {
	ID   string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	Name string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
}

// LegalMonetaryTotal represents legal monetary total
type LegalMonetaryTotal struct {
	LineExtensionAmount   MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 LineExtensionAmount,omitempty"`
	TaxExclusiveAmount    MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxExclusiveAmount,omitempty"`
	TaxInclusiveAmount    MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxInclusiveAmount,omitempty"`
	AllowanceTotalAmount  MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 AllowanceTotalAmount,omitempty"`
	ChargeTotalAmount     MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ChargeTotalAmount,omitempty"`
	PrepaidAmount         MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PrepaidAmount,omitempty"`
	PayableRoundingAmount MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PayableRoundingAmount,omitempty"`
	PayableAmount         MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PayableAmount"`
}

// InvoiceLine represents invoice line
type InvoiceLine struct {
	ID                  string              `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	UUID                string              `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 UUID,omitempty"`
	Note                []string            `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Note,omitempty"`
	InvoicedQuantity    Quantity            `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 InvoicedQuantity"`
	LineExtensionAmount MonetaryAmount      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 LineExtensionAmount"`
	AccountingCost      string              `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 AccountingCost,omitempty"`
	OrderLineReference  []OrderLineReference `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 OrderLineReference,omitempty"`
	TaxTotal            []TaxTotal          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxTotal,omitempty"`
	Item                Item                `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Item"`
	Price               Price               `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Price"`
}

// OrderLineReference represents order line reference
type OrderLineReference struct {
	LineID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 LineID"`
}

// Item represents item
type Item struct {
	Description                []string               `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Description,omitempty"`
	Name                       string                 `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name,omitempty"`
	BuyersItemIdentification   *ItemIdentification    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 BuyersItemIdentification,omitempty"`
	SellersItemIdentification  *ItemIdentification    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 SellersItemIdentification,omitempty"`
	StandardItemIdentification *ItemIdentification    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 StandardItemIdentification,omitempty"`
	OriginCountry              *Country               `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 OriginCountry,omitempty"`
	CommodityClassification    []CommodityClassification `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 CommodityClassification,omitempty"`
	ClassifiedTaxCategory      []TaxCategory          `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ClassifiedTaxCategory,omitempty"`
	AdditionalItemProperty     []ItemProperty         `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 AdditionalItemProperty,omitempty"`
}

// ItemIdentification represents item identification
type ItemIdentification struct {
	ID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
}

// CommodityClassification represents commodity classification
type CommodityClassification struct {
	ItemClassificationCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ItemClassificationCode,omitempty"`
}

// ItemProperty represents item property
type ItemProperty struct {
	Name  string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
	Value string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Value,omitempty"`
}

// Price represents price
type Price struct {
	PriceAmount         MonetaryAmount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PriceAmount"`
	BaseQuantity        Quantity       `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 BaseQuantity,omitempty"`
	PriceChangeReason   []string       `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PriceChangeReason,omitempty"`
	PriceTypeCode       string         `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PriceTypeCode,omitempty"`
	PriceType           string         `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PriceType,omitempty"`
}

// MonetaryAmount represents a monetary amount with currency
type MonetaryAmount struct {
	Value      float64 `xml:",chardata"`
	CurrencyID string  `xml:"currencyID,attr,omitempty"`
}

// Quantity represents a quantity with unit of measure
type Quantity struct {
	Value    float64 `xml:",chardata"`
	UnitCode string  `xml:"unitCode,attr,omitempty"`
}
