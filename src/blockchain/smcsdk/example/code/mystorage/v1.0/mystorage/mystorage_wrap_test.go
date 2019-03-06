package mystorage

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
	contractName       = "mystorage" //contract name
	contractMethods    = []string{"Set(uint64)", "Get()TYPE"}
	contractInterfaces = []string{}
	orgID              = "orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer"
)

//TestObject This is a struct for test
type TestObject struct {
	obj *MyStorage
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
	return &TestObject{&MyStorage{sdk: utest.UTP.ISmartContract}}
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

//Set This is a method of TestObject
func (t *TestObject) Set(data uint64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Set(data)
	utest.Commit()
	return
}

//Get This is a method of TestObject
func (t *TestObject) Get() (result0 uint64, err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	result0 = t.obj.Get()
	utest.Commit()
	return
}
