package tx2

import (
	"blockchain/smcsdk/sdk/rlp"
	"blockchain/types"
	"bytes"
	"common/kms"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
)

// TxParse 解析一笔交易（包含签名验证）的接口函数，将结果填入Transaction数据结构，其中Data字段为RLP编码的合约调用参数
func TxParse(txString string) (tx types.Transaction, pubKey crypto.PubKeyEd25519, err error) {
	MAC := chainID + "<tx>"
	Version := "v2"
	SignerNumber := "<1>"
	strs := strings.Split(txString, ".")

	if strs[0] != MAC || strs[1] != Version || strs[3] != SignerNumber {
		err = errors.New("tx data error")
		return
	}

	txData := base58.Decode(strs[2])
	sigBytes := base58.Decode(strs[4])

	reader := bytes.NewReader(sigBytes)
	var siginfo kms.Ed25519Sig
	err = rlp.Decode(reader, &siginfo)
	if err != nil {
		return
	}

	if !siginfo.PubKey.VerifyBytes(txData, siginfo.SigValue) {
		err = errors.New("verify sig fail")
		return
	}
	pubKey = siginfo.PubKey

	//RLP解码Transaction结构
	reader = bytes.NewReader(txData)
	err = rlp.Decode(reader, &tx)
	return
}
