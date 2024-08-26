package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func EcdsaTokenGenerate() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed for generate key:%v", err)
	}

	// Encode private key to PKCS#8 ASN.1 DER format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to encode the private key: %v", err)
	}

	privateKeyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Write the PEM block to a file
	file, err := os.Create("ecdsa_private_key.pem")
	if err != nil {
		log.Fatalf("Failed to generate the file: %v", err)
	}
	defer file.Close()

	err = pem.Encode(file, privateKeyPem)
	if err != nil {
		log.Fatalf("Failed to write PEM block to file: %v", err)
	}
	fmt.Println("ECDSA private key generated and saved to ecdsa_private_key.pem")
}
