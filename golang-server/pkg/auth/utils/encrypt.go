package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func ParsePublicKey(publicKeyPem string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

// Function to encrypt data using RSA-OAEP with SHA-256
func EncryptWithPublicKey(message string, publicKey *rsa.PublicKey) (string, error) {
	label := []byte("") // optional label
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, []byte(message), label)
	if err != nil {
		return "", err
	}

	// Return the ciphertext as a base64-encoded string for easy transmission
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Encrypt(message string, publicKeyPem string) (string, error) {
	publicKey, err := ParsePublicKey(publicKeyPem)
	if err != nil {
		return "", err
	}
	return EncryptWithPublicKey(message, publicKey)
}
