package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/ochom/gutils/logs"
)

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func HashText(certText, text string) string {
	certContent := []byte(certText)

	// Decode PEM block
	block, _ := pem.Decode(certContent)
	if block == nil {
		logs.Error("decoding certificate pem")
		return ""
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		logs.Error("parsing certificate: %v", err)
		return ""
	}

	pub := cert.PublicKey.(*rsa.PublicKey)
	msg := []byte(text)

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		logs.Error("encrypting message: %v", err)
		return ""
	}

	return Encode(cipher)
}
