package launcher

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

// generateEd25519KeyPair creates an Ed25519 key pair for testing.
func generateEd25519KeyPair() (pub ed25519.PublicKey, priv ed25519.PrivateKey, err error) {
	return ed25519.GenerateKey(rand.Reader)
}

// marshalPrivateKey encodes an Ed25519 private key in OpenSSH-compatible PEM format.
func marshalPrivateKey(priv ed25519.PrivateKey) []byte {
	bytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		panic(err)
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytes,
	}
	return pem.EncodeToMemory(block)
}
