package signature

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

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