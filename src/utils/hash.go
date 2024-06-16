package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"os"

	"github.com/ochom/gutils/logs"
)

func HashText(certPath, text string) string {
	certFile, err := os.OpenFile(certPath, os.O_RDONLY, 0)
	if err != nil {
		logs.Error("reading certificate file: %v", err)
		return ""
	}

	defer certFile.Close()

	certContent, err := io.ReadAll(certFile)
	if err != nil {
		logs.Error("reading certificate content: %v", err)
		return ""
	}

	// Decode PEM block
	block, _ := pem.Decode(certContent)
	if block == nil {
		logs.Error("parsing certificate: %v", err)
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

	return base64.StdEncoding.EncodeToString(cipher)
}
