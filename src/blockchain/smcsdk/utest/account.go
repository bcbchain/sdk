// unittestplatform
// account.go 实现测试账户管理功能，包括创建账户地址等

package utest

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl"
	"blockchain/smcsdk/sdkimpl/object"
	"crypto/rand"

	"github.com/tendermint/go-crypto"
)

//FuncRecover recover panic by Assert
func FuncRecover(err *types.Error) {
	err.ErrorCode = types.CodeOK
	if rerr := recover(); rerr != nil {
		if _, ok := rerr.(types.Error); ok {
			err.ErrorCode = rerr.(types.Error).ErrorCode
			err.ErrorDesc = rerr.(types.Error).ErrorDesc
		} else {
			panic(rerr)
		}
	}
}

func newRandPubKey() []byte {
	tmp := make([]byte, 32)
	_, err := rand.Read(tmp)
	if err != nil {
		panic(err.Error())
	}
	return tmp
}

// CalcAddressFromPubKey calc address from public key
func CalcAddressFromPubKey(_pubKey []byte) types.Address {
	pk := crypto.PubKeyEd25519FromBytes(_pubKey)
	return pk.Address()
}

// NewAccount generate a new account object with a given token and balance
func NewAccount(tokenName string, balance bn.Number) sdk.IAccount {
	addr := CalcAddressFromPubKey(newRandPubKey())
	UTP.accountList = append(UTP.accountList, addr)

	if balance.IsGreaterThanI(0) {
		Transfer(nil, tokenName, addr, balance)
		UTP.ISmartContract.(*sdkimpl.SmartContract).Commit()
	}

	return object.NewAccount(UTP.ISmartContract, addr)
}

// NewAccounts generate some new ccount objects with a given token and balance
func NewAccounts(tokenName string, balance bn.Number, count int) []sdk.IAccount {
	accounts := make([]sdk.IAccount, 0)
	for i := 0; i < count; i++ {
		pubKey := newRandPubKey()
		addr := CalcAddressFromPubKey(pubKey)
		UTP.accountList = append(UTP.accountList, addr)

		if balance.IsGreaterThanI(0) {
			Transfer(nil, tokenName, addr, balance)
			UTP.ISmartContract.(*sdkimpl.SmartContract).Commit()
		}

		accounts = append(accounts, object.NewAccount(UTP.ISmartContract, addr))
	}

	return accounts
}

//GetAccount get an account
func (ut *UtPlatform) GetAccount(index int) types.Address {
	if index >= len(ut.accountList) {
		return ""
	}

	return ut.accountList[index]
}

//Transfer transfer token to account
func Transfer(sender sdk.IAccount, tokenName, addr string, value bn.Number) (err types.Error) {
	defer FuncRecover(&err)

	if value.CmpI(0) > 0 {
		contract := UTP.Message().Contract()

		var ic sdk.IToken
		if tokenName == "" {
			//转本合约代币
			ic = UTP.Helper().TokenHelper().Token()
		} else {
			//代币，只能调用自己合约的代币
			ic = UTP.Helper().TokenHelper().TokenOfName(tokenName)
			if ic != nil {
				tempContract := UTP.Helper().ContractHelper().ContractOfToken(ic.Address())
				UTP.Message().(*object.Message).SetContract(tempContract)
			}
		}

		if ic == nil {
			err.ErrorCode = types.ErrUserDefined // 使用sdk中未使用的错误定义，避免干扰测试结果
			err.ErrorDesc = "Invalid token name=" + tokenName
			return
		}

		if sender == nil {
			sender = object.NewAccount(UTP.ISmartContract, ic.Owner())
		}
		sender.TransferByToken(ic.Address(), addr, value)
		UTP.Message().(*object.Message).SetContract(contract)
	}

	return
}
