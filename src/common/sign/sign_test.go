package sign

import (
	"encoding/hex"
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestSign(t *testing.T) {
	signer, err := loadPrivateKey("private.pem")
	if err != nil {
		fmt.Printf("signer is damaged: %v", err)
	}

	toSign := "date: Thu, 05 Jan 2012 21:31:40 GMT"
	//srcByte, err := ioutil.ReadFile("contract.tar.gz")
	//if err != nil {
	//	fmt.Println("source file read error:", err)
	//	return
	//}

	signed, err := signer.Sign([]byte(toSign))
	//signed, err := signer.Sign(srcByte)
	if err != nil {
		fmt.Printf("could not sign request: %v", err)
	}
	sig := hex.EncodeToString(signed)
	fmt.Printf("Signature: %v\n", sig)

	sigByte, err := hex.DecodeString(sig)
	if err != nil {
		fmt.Println("decode sig cause error:", err)
	}

	parser, err := loadPublicKey("public.pem")
	if err != nil {
		fmt.Printf("could not load publicKey: %v", err)
	}

	err = parser.CheckSign([]byte(toSign), sigByte)
	// err = parser.CheckSign(srcByte, sigByte)

	assert.Equal(t, err, nil)

}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPublicKey(path string) (Checker, error) {
	return ParsePublicKey([]byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`))
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(path string) (Signer, error) {
	return ParsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`))
}
