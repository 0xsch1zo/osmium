package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func GenerateKey() (*rsa.PrivateKey, error) {
	const keySize = 2048
	return rsa.GenerateKey(rand.Reader, keySize)
}

func PrivRsaToPem(privateKey *rsa.PrivateKey) string {
	privatePem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return string(privatePem)
}

func PemToPrivRsa(privateKeyPem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		return nil, errors.New("Failed to decode PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func PubRsaToPem(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	publicPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicPem), nil
}

func PemToPubRsa(publicKeyPem string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil {
		return nil, errors.New("Failed to decode PEM block containing the key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	t, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Key is not of RSA format")
	}

	return t, nil
}
