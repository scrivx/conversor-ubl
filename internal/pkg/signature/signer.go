package signature

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/amdonov/xmlsig"
	"software.sslmate.com/src/go-pkcs12"
)

func LoadKeyPairFromPFX(pfxPath, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	pfxData, err := os.ReadFile(pfxPath)
	if err != nil {
		return nil, nil, err
	}

	privateKey, certificate, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, nil, err
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("clave privada no es RSA")
	}

	return certificate, rsaKey, nil
}

func SignXML(xmlInput string, cert *x509.Certificate, key *rsa.PrivateKey) (string, error) {
	if cert == nil || key == nil {
		return "", errors.New("certificado o clave privada no pueden ser nil")
	}

	// Convertir certificado y clave a formato tls.Certificate
	tlsCert, err := createTLSCertificate(cert, key)
	if err != nil {
		return "", err
	}

	// Crear el firmador con xmlsig
	signer, err := xmlsig.NewSigner(tlsCert)
	if err != nil {
		return "", err
	}

	// Firmar el XML directamente
	signedBytes, err := signer.Sign([]byte(xmlInput))
	if err != nil {
		return "", err
	}

	return string(signedBytes), nil
}

// Función auxiliar para convertir x509.Certificate y rsa.PrivateKey a tls.Certificate
func createTLSCertificate(cert *x509.Certificate, key *rsa.PrivateKey) (tls.Certificate, error) {
	// Convertir certificado a PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	// Convertir clave privada a PEM
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Crear tls.Certificate desde PEM
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tlsCert, nil
}

// Función alternativa con opciones personalizadas
func SignXMLWithOptions(xmlInput string, cert *x509.Certificate, key *rsa.PrivateKey) (string, error) {
	if cert == nil || key == nil {
		return "", errors.New("certificado o clave privada no pueden ser nil")
	}

	// Convertir a tls.Certificate
	tlsCert, err := createTLSCertificate(cert, key)
	if err != nil {
		return "", err
	}

	// Crear firmador con opciones personalizadas
	options := xmlsig.SignerOptions{
		SignatureAlgorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		DigestAlgorithm:    "http://www.w3.org/2001/04/xmlenc#sha256",
	}

	signer, err := xmlsig.NewSignerWithOptions(tlsCert, options)
	if err != nil {
		return "", err
	}

	// Firmar el XML
	signedBytes, err := signer.Sign([]byte(xmlInput))
	if err != nil {
		return "", err
	}

	return string(signedBytes), nil
}

// SignatureElement representa el elemento ds:Signature
type SignatureElement struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	SignedInfo     SignedInfo `xml:"SignedInfo"`
	SignatureValue SignatureValue `xml:"SignatureValue"`
	KeyInfo        KeyInfo `xml:"KeyInfo"`
}

type SignedInfo struct {
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}

type CanonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type Reference struct {
	URI          string      `xml:"URI,attr"`
	Transforms   Transforms  `xml:"Transforms"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string      `xml:"DigestValue"`
}

type Transforms struct {
	Transform []Transform `xml:"Transform"`
}

type Transform struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type DigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureValue struct {
	Value string `xml:",chardata"`
}

type KeyInfo struct {
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

// SignXMLAsElement firma el XML y devuelve solo el elemento ds:Signature
func SignXMLAsElement(xmlInput string, cert *x509.Certificate, key *rsa.PrivateKey) (string, error) {
	if cert == nil || key == nil {
		return "", errors.New("certificado o clave privada no pueden ser nil")
	}

	// Convertir certificado y clave a formato tls.Certificate
	tlsCert, err := createTLSCertificate(cert, key)
	if err != nil {
		return "", err
	}

	// Crear el firmador con xmlsig
	signer, err := xmlsig.NewSigner(tlsCert)
	if err != nil {
		return "", err
	}

	// Firmar el XML
	signedBytes, err := signer.Sign([]byte(xmlInput))
	if err != nil {
		return "", err
	}

	// La librería xmlsig devuelve solo los bytes de la firma, no el elemento XML
	// Necesitamos generar manualmente el elemento ds:Signature
	signatureElement, err := generateSignatureElement([]byte(signedBytes), cert)
	if err != nil {
		return "", err
	}

	return signatureElement, nil
}

// generateSignatureElement genera el elemento ds:Signature manualmente
func generateSignatureElement(signatureBytes []byte, cert *x509.Certificate) (string, error) {
	// Convertir certificado a base64
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	
	// Convertir firma a base64
	signatureB64 := string(signatureBytes)
	
	// Generar el elemento ds:Signature
	signatureXML := fmt.Sprintf(`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <ds:SignedInfo>
    <ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
    <ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/>
    <ds:Reference URI="">
      <ds:Transforms>
        <ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/>
      </ds:Transforms>
      <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
      <ds:DigestValue>%s</ds:DigestValue>
    </ds:Reference>
  </ds:SignedInfo>
  <ds:SignatureValue>%s</ds:SignatureValue>
  <ds:KeyInfo>
    <ds:X509Data>
      <ds:X509Certificate>%s</ds:X509Certificate>
    </ds:X509Data>
  </ds:KeyInfo>
</ds:Signature>`, signatureB64, signatureB64, strings.TrimSpace(string(certPEM)))
	
	return signatureXML, nil
}