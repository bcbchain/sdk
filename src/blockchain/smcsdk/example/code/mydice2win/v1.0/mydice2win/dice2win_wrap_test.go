package mydice2win

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl/object"
	"blockchain/smcsdk/sdkimpl/sdkhelper"
	"blockchain/smcsdk/utest"
	"fmt"
)

var (
	contractName       = "mydice2win" //contract name
	contractMethods    = []string{"SetSecretSigner(types.PubKey)", "SetSettings(string)", "SetRecvFeeInfos(string)", "WithdrawFunds(string,types.Address,bn.Number)", "PlaceBet(bn.Number,int64,int64,[]byte,[]byte,types.Address)", "SettleBet([]byte)", "RefundBet([]byte)"}
	contractInterfaces = []string{"PlaceBet(bn.Number,int64,int64,[]byte,[]byte,types.Address)"}
	orgID              = "orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer"
)

//TestObject This is a struct for test
type TestObject struct {
	obj *Dice2Win
}

//FuncRecover recover panic by Assert
func FuncRecover(err *types.Error) {
	if rerr := recover(); rerr != nil {
		if _, ok := rerr.(types.Error); ok {
			err.ErrorCode = rerr.(types.Error).ErrorCode
			err.ErrorDesc = rerr.(types.Error).ErrorDesc
			fmt.Println(err)
		} else {
			panic(rerr)
		}
	}
}

//NewTestObject This is a function
func NewTestObject(sender sdk.IAccount) *TestObject {
	return &TestObject{&Dice2Win{sdk: utest.UTP.ISmartContract}}
}

//transfer This is a method of TestObject
func (t *TestObject) transfer(balance bn.Number) *TestObject {
	contract := t.obj.sdk.Message().Contract()
	utest.Transfer(t.obj.sdk.Message().Sender(), t.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), balance)
	t.obj.sdk = sdkhelper.OriginNewMessage(t.obj.sdk, contract, t.obj.sdk.Message().MethodID(), t.obj.sdk.Message().(*object.Message).OutputReceipts())
	return t
}

//setSender This is a method of TestObject
func (t *TestObject) setSender(sender sdk.IAccount) *TestObject {
	t.obj.sdk = utest.SetSender(sender.Address())
	return t
}

//run This is a method of TestObject
func (t *TestObject) run() *TestObject {
	t.obj.sdk = utest.ResetMsg()
	return t
}

//InitChain This is a method of TestObject
func (t *TestObject) InitChain() {
	utest.NextBlock(1)
	t.obj.InitChain()
	utest.Commit()
	return
}

//SetSecretSigner This is a method of TestObject
func (t *TestObject) SetSecretSigner(newSecretSigner types.PubKey) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetSecretSigner(newSecretSigner)
	utest.Commit()
	return
}

//SetSettings This is a method of TestObject
func (t *TestObject) SetSettings(newSettingsStr string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetSettings(newSettingsStr)
	utest.Commit()
	return
}

//SetRecvFeeInfos This is a method of TestObject
func (t *TestObject) SetRecvFeeInfos(recvFeeInfosStr string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetRecvFeeInfos(recvFeeInfosStr)
	utest.Commit()
	return
}

//WithdrawFunds This is a method of TestObject
func (t *TestObject) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.WithdrawFunds(tokenName, beneficiary, withdrawAmount)
	utest.Commit()
	return
}

//PlaceBet This is a method of TestObject
func (t *TestObject) PlaceBet(betMask bn.Number, modulo, commitLastBlock int64, commit, signData []byte, refAddress types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.PlaceBet(betMask, modulo, commitLastBlock, commit, signData, refAddress)
	utest.Commit()
	return
}

//SettleBet This is a method of TestObject
func (t *TestObject) SettleBet(reveal []byte) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SettleBet(reveal)
	utest.Commit()
	return
}

//RefundBet This is a method of TestObject
func (t *TestObject) RefundBet(commit []byte) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.RefundBet(commit)
	utest.Commit()
	return
}
