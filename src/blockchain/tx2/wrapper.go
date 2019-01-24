package tx2

import (
	"common/kms"
	"encoding/hex"
	"strings"

	"blockchain/smcsdk/sdk/rlp"
	"blockchain/types"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data/base58"
	"github.com/tendermint/tmlibs/common"
)

var (
	chainID string
)

// Init - chainID
func Init(_chainID string) {
	chainID = _chainID
}

// WrapInvokeParams - wrap contract parameters
func WrapInvokeParams(params ...interface{}) []common.HexBytes {
	paramsRlp := make([]common.HexBytes, 0)
	for _, param := range params {
		var paramRlp []byte
		var err error

		paramRlp, err = rlp.EncodeToBytes(param)
		if err != nil {
			panic(err)
		}
		paramsRlp = append(paramsRlp, paramRlp)
	}
	return paramsRlp
}

// WrapPayload - wrap contracts to payload byte
func WrapPayload(nonce uint64, gasLimit int64, note string, messages ...types.Message) []byte {

	type transaction struct {
		Nonce    uint64
		GasLimit int64
		Note     string
		Messages []types.Message
	}
	tx := transaction{
		Nonce:    nonce,
		GasLimit: gasLimit,
		Note:     note,
		Messages: messages,
	}
	txRlp, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return txRlp
}

// WrapTx - sign the payload to string
// privateKey的格式:
// name:password
// enprivatekey:password
// 0x十六进制表示的私钥数据
//nolint unhandled
func WrapTx(payload []byte, privateKey string) string {
	var sigInfo kms.Ed25519Sig
	var isHexPrivKey = strings.HasPrefix(privateKey, "0x")
	var segPrivKey = strings.Split(privateKey, ":")

	if isHexPrivKey && len(segPrivKey) == 1 {
		var privKey crypto.PrivKey
		var pubKey crypto.PubKey

		hexData := privateKey[2:]
		privKeyBytes, err := hex.DecodeString(hexData)
		if err != nil {
			panic(err.Error())
		}
		privKey = crypto.PrivKeyEd25519FromBytes(privKeyBytes)
		pubKey = privKey.PubKey()

		sigInfo = kms.Ed25519Sig{
			SigType:  "ed25519",
			PubKey:   pubKey.(crypto.PubKeyEd25519),
			SigValue: privKey.Sign(payload).(crypto.SignatureEd25519),
		}
	} else if len(segPrivKey) == 2 {
		si, err := kms.SignData(segPrivKey[0], segPrivKey[1], payload)
		if err != nil {
			panic(err.Error())
		}
		sigInfo = *si
	} else {
		panic("Invalid private key format")
	}

	size, r, err := rlp.EncodeToReader(sigInfo)
	if err != nil {
		panic(err.Error())
	}
	sig := make([]byte, size)
	r.Read(sig)

	payloadString := base58.Encode(payload)
	sigString := base58.Encode(sig)

	MAC := string(chainID) + "<tx>"
	Version := "v2"
	SignerNumber := "<1>"

	return MAC + "." + Version + "." + payloadString + "." + SignerNumber + "." + sigString
}
