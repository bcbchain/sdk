package mycoin

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl/object"
	"blockchain/smcsdk/sdkimpl/sdkhelper"
	"blockchain/smcsdk/utest"
)

var (
	contractName    = "mycoin" //contract name
	contractMethods = []string{"Transfer(types.Address,bn.Number)"}
	orgID           = "orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer"
)

//TestObject This is a struct for test
type TestObject struct {
	obj *Mycoin
}

//FuncRecover recover panic by Assert
func FuncRecover(err *types.Error) {
	if rerr := recover(); rerr != nil {
		if _, ok := rerr.(types.Error); ok {
			err.ErrorCode = rerr.(types.Error).ErrorCode
			err.ErrorDesc = rerr.(types.Error).ErrorDesc
		} else {
			panic(rerr)
		}
	}
}

//NewTestObject This is a function
func NewTestObject(sender sdk.IAccount) *TestObject {
	return &TestObject{&Mycoin{sdk: utest.UTP.ISmartContract}}
}

//transfer This is a method of TestObject
func (t *TestObject) transfer(balance bn.Number) *TestObject {
	t.obj.sdk.Message().Sender().TransferByName(t.obj.sdk.Helper().GenesisHelper().Token().Name(), t.obj.sdk.Message().Contract().Account(), balance)
	t.obj.sdk = sdkhelper.OriginNewMessage(t.obj.sdk, t.obj.sdk.Message().Contract(), t.obj.sdk.Message().MethodID(), t.obj.sdk.Message().(*object.Message).OutputReceipts())
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

//Transfer This is a method of TestObject
func (t *TestObject) Transfer(to types.Address, value bn.Number) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Transfer(to, value)
	utest.Commit()
	return
}
