package wal

import (
	"blockchain/algorithm"
	"bytes"
	"common/fs"
	"common/sig"
	"common/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	cmn "github.com/tendermint/tmlibs/common"
	"golang.org/x/crypto/sha3"
)

var cdc = amino.NewCodec()

func init() {
	crypto.RegisterAmino(cdc)
}

// ----- account struct -----
type Account struct {
	Name         string         `json:"name"`
	PrivateKey   crypto.PrivKey `json:"privateKey"`
	Hash         []byte         `json:"hash"`
	keyStoreFile string
}

func NewAccount(keyStoreDir, name, password string) (acct *Account, err error) {
	privateKey := crypto.GenPrivKeyEd25519()
	return ImportAccount(keyStoreDir, name, password, privateKey)
}

func ImportAccount(keyStoreDir, name, password string, privKey crypto.PrivKey) (acct *Account, err error) {
	privateKey := privKey.(crypto.PrivKeyEd25519)
	keyStoreFile := filepath.Join(keyStoreDir, name+".wal")

	acct = &Account{
		Name:         name,
		PrivateKey:   privateKey,
		keyStoreFile: keyStoreFile,
	}

	sha256 := sha3.New256()
	sha256.Write([]byte(name))
	sha256.Write(privateKey[:])
	acct.Hash = sha256.Sum(nil)

	if err = acct.save(password, true); err != nil {
		acct = nil
	}
	return
}

func LoadAccount(keyStoreDir, name, password string) (acct *Account, err error) {
	acct = &Account{}
	keyStoreFile := filepath.Join(keyStoreDir, name+".wal")

	fmt.Println(keyStoreFile)
	walBytes, err := ioutil.ReadFile(keyStoreFile)
	if err != nil {
		return nil, errors.New("account does not exist")
	}

	passwordBytes := make([]byte, 0)
	if password == "" {
		passwordBytes, err = utils.CheckPassword("Enter password (" + name + "): ")
		if err != nil {
			return nil, err
		}
	} else {
		passwordBytes = []byte(password)
	}

	jsonBytes, err := algorithm.DecryptWithPassword(walBytes, passwordBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("the password is wrong, err info : %s", err)
	}
	err = cdc.UnmarshalJSON(jsonBytes, acct)
	if err != nil {
		return nil, err
	}

	privkey := acct.PrivateKey.(crypto.PrivKeyEd25519)
	sha256 := sha3.New256()
	sha256.Write([]byte(acct.Name))
	sha256.Write(privkey[:])
	hash := sha256.Sum(nil)
	if bytes.Equal(hash, acct.Hash) == false {
		return nil, fmt.Errorf("verify hash of wallet failed")
	}

	acct.keyStoreFile = keyStoreFile
	return acct, nil
}

func (acct *Account) Save(password string) (err error) {
	return acct.save(password, false)
}

func (acct *Account) save(password string, notAllowExist bool) (err error) {

	privkey := acct.PrivateKey.(crypto.PrivKeyEd25519)
	sha256 := sha3.New256()
	sha256.Write([]byte(acct.Name))
	sha256.Write(privkey[:])
	hash := sha256.Sum(nil)
	if bytes.Equal(hash, acct.Hash) == false {
		return fmt.Errorf("verify hash of wallet failed")
	}

	if acct.keyStoreFile == "" {
		return errors.New("no key store file specified in account object")
	}
	if ok, _ := fs.PathExists(acct.keyStoreFile); ok && notAllowExist {
		return errors.New("key store file is already exist")
	}

	keyStoreDir := filepath.Dir(acct.keyStoreFile)
	if ok, _ := fs.PathExists(keyStoreDir); !ok {
		if ok, err = fs.MakeDir(keyStoreDir); err != nil {
			return err
		}
	}

	passwordBytes := []byte(password)
	if password == "" {
		passwordBytes, err = utils.GetAndCheckPassword(
			"Enter  password ("+acct.Name+"): ",
			"Repeat password ("+acct.Name+"): ")
		if err != nil {
			return err
		}
	} else {
		passwordBytes = []byte(password)
	}

	jsonBytes, err := cdc.MarshalJSON(acct)
	if err != nil {
		return err
	}
	walBytes := algorithm.EncryptWithPassword(jsonBytes, passwordBytes, nil)
	err = cmn.WriteFileAtomic(acct.keyStoreFile, walBytes, 0600)
	if err != nil {
		return err
	}
	return
}

func (acct *Account) PubKey() (pubKey crypto.PubKey) {
	return acct.PrivateKey.PubKey()
}

func (acct *Account) Address(chainId string) (address string) {
	return acct.PrivateKey.PubKey().AddressByChainID(chainId)
}

func (acct *Account) Sign(data []byte) (sigInfo *sig.Ed25519Sig, err error) {
	return sig.Sign(acct.PrivateKey, data)
}

func (acct *Account) Sign2File(data []byte, sigFile string) (err error) {
	return sig.Sign2File(acct.PrivateKey, data, sigFile)
}

func (acct *Account) SignBinFile(binFile, sigFile string) (err error) {
	return sig.SignBinFile(acct.PrivateKey, binFile, sigFile)
}

func (acct *Account) SignTextFile(textFile, sigFile string) (err error) {
	return sig.SignTextFile(acct.PrivateKey, textFile, sigFile)
}
