package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl/object"
	"blockchain/smcsdk/sdkimpl/sdkhelper"
	"blockchain/smcsdk/utest"
	"fmt"
)

var (
	contractName       = "excellencies" //contract name
	contractMethods    = []string{"SetSecretSigner(types.PubKey)", "SetSettings(string)", "SetRecFeeInfo(string)", "PlaceBet(string,int64,[]byte,[]byte,types.Address,types.Address)", "SettleBet([]byte,int64)", "WithdrawWin([]byte)", "RefundBets([]byte,int64)"}
	contractInterfaces = []string{}
	orgID              = "orgCZkw5xz9DYa3h5pJ2CzZSuGHRCj2ot5xq"
)

//TestObject This is a struct for test
type TestObject struct {
	obj *Excellencies
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
	return &TestObject{&Excellencies{sdk: utest.UTP.ISmartContract}}
}

//transfer This is a method of TestObject
func (t *TestObject) transfer(args ...interface{}) *TestObject {
	contract := t.obj.sdk.Message().Contract()
	utest.Transfer(t.obj.sdk.Message().Sender(), contract.Account(), args...)
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
func (t *TestObject) SetSettings(newSettinsStr string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetSettings(newSettinsStr)
	utest.Commit()
	return
}

//SetRecFeeInfo This is a method of TestObject
func (t *TestObject) SetRecFeeInfo(recFeeInfoStr string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetRecFeeInfo(recFeeInfoStr)
	utest.Commit()
	return
}

//PlaceBet This is a method of TestObject
func (t *TestObject) PlaceBet(betJson string, commitLastBlock int64, commit, signData []byte, refAddress types.Address, saler types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.PlaceBet(betJson, commitLastBlock, commit, signData, refAddress, saler)
	utest.Commit()
	return
}

//SettleBet This is a method of TestObject
func (t *TestObject) SettleBet(reveal []byte, settleCount int64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SettleBet(reveal, settleCount)
	utest.Commit()
	return
}

//WithdrawWin This is a method of TestObject
func (t *TestObject) WithdrawWin(commit []byte) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.WithdrawWin(commit)
	utest.Commit()
	return
}

//RefundBets This is a method of TestObject
func (t *TestObject) RefundBets(commit []byte, refundCount int64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.RefundBets(commit, refundCount)
	utest.Commit()
	return
}
