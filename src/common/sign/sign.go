package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// ParsePublicKey parses a PEM encoded private key.
func ParsePublicKey(pemBytes []byte) (Checker, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("load PublicKey: no key found")
	}

	var rawKey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsaKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawKey = rsaKey
	default:
		return nil, fmt.Errorf("load PublicKey: unsupported key type %q", block.Type)
	}

	return newCheckerFromKey(rawKey)
}

// ParsePrivateKey parses a PEM encoded private key.
func ParsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("privKey decode: no key found")
	}

	var rawKey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawKey = rsaKey
	default:
		return nil, fmt.Errorf("privKey parser: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawKey)
}

// Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keyType to the data.
	Sign(data []byte) ([]byte, error)
}

// Checker is can create signatures that verify against a public key.
type Checker interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keyType to the data.
	CheckSign(data []byte, sig []byte) error
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var signer Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		signer = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("signer constructor: unsupported key type %T", k)
	}
	return signer, nil
}

func newCheckerFromKey(k interface{}) (Checker, error) {
	var checker Checker
	switch t := k.(type) {
	case *rsa.PublicKey:
		checker = &rsaPublicKey{t}
	default:
		return nil, fmt.Errorf("checker constructor: unsupported key type %T", k)
	}
	return checker, nil
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

// CheckSign verifies the message using a rsa-sha256 signature
func (r *rsaPublicKey) CheckSign(message []byte, sig []byte) error {
	h := sha256.New()
	if _, err := h.Write(message); err != nil {
		return err
	}
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
}
