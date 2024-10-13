package certgen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func generateSelfSignedCert(certFile, keyFile, commonName string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1-year certificate

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})

	log.Println("Self-signed certificate and key generated successfully:", commonName)
	return nil
}

func SetupCertificates() {
	certDir := "cert/certs"
	if err := os.MkdirAll(certDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create cert/certs directory: %v", err)
	}

	serverCertFile := filepath.Join(certDir, ".server-cert.pem")
	serverKeyFile := filepath.Join(certDir, ".server-key.pem")
	clientCertFile := filepath.Join(certDir, ".client-cert.pem")
	clientKeyFile := filepath.Join(certDir, ".client-key.pem")

	if err := generateSelfSignedCert(serverCertFile, serverKeyFile, "My Server"); err != nil {
		log.Fatalf("Failed to generate server certificate: %v", err)
	}

	if err := generateSelfSignedCert(clientCertFile, clientKeyFile, "My Client"); err != nil {
		log.Fatalf("Failed to generate client certificate: %v", err)
	}
}
