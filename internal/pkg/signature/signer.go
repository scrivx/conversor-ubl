package signature

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	dsig "github.com/russellhaering/goxmldsig"
	"software.sslmate.com/src/go-pkcs12"
)

func LoadKeyPairFromPFX(pfxPath string, password string) (*x509.Certificate, crypto.PrivateKey, error) {
	pfxData, err := os.ReadFile(pfxPath)
	if err != nil {
		return nil, nil, err
	}

	privateKey, certificate, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, nil, err
	}

	return certificate, privateKey, nil
}

func SignXML(xmlInput string, cert *x509.Certificate, key crypto.PrivateKey) (string, error) {
	ctx := dsig.NewDefaultSigningContext(dsig.NewSigner(privateRSAKey(key)))
	ctx.Hash = crypto.SHA1

	ctx.SetSignatureMethod(dsig.RSASHA1SignatureMethod)
	ctx.Canonicalizer = dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	ctx.SetCertificateChain([][]byte{cert.Raw})

	signedXML, err := ctx.SignEnveloped([]byte(xmlInput))
	if err != nil {
		return "", err
	}

	return string(signedXML), nil
}

func privateRSAKey(key crypto.PrivateKey) *rsa.PrivateKey {
	if rsaKey, ok := key.(*rsa.PrivateKey); ok {
		return rsaKey
	}
	panic(errors.New("clave privada no es RSA"))
}
