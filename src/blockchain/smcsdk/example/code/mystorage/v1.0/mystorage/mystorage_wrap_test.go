package mystorage

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl"
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

//setSender This is a method of TestObject
func (t *TestObject) setSender(sender sdk.IAccount) *TestObject {
	t.obj.sdk = utest.SetSender(sender.Address())
	return t
}

// run This is a method of TestObject
func (t *TestObject) run(errCode uint32, f func(t *TestObject) types.Error) {
	utest.SetFlag(true)
	msg := t.obj.sdk.Message()
	smc := t.obj.sdk
	// new message, empty input
	sdkhelper.OriginNewMessage(smc, smc.Message().Contract(), smc.Message().MethodID(), nil)

	t.resetMsg(t.obj.sdk.Message().Origins(), nil)

	err := f(t)

	utest.AssertError(err, errCode)

	if err.ErrorCode == types.CodeOK {
		utest.Commit()
	} else {
		utest.Rollback()
	}
	utest.SetFlag(false)

	t.obj.sdk.(*sdkimpl.SmartContract).SetMessage(msg)
}

// addOrigins This is a method of TestObject
func (t *TestObject) addOrigins(newOrigins []string) {
	smc := t.obj.sdk
	oldO := smc.Message().Origins()
	oldO = append(oldO, newOrigins...)

	t.resetMsg(oldO, smc.Message().InputReceipts())
}

// emitReceipt This is a method of TestObject
func (t *TestObject) emitReceipt(receipt interface{}) {
	t.obj.sdk.Helper().ReceiptHelper().Emit(receipt)

	t.resetMsg(t.obj.sdk.Message().Origins(), t.obj.sdk.Message().(*object.Message).OutputReceipts())
}

func (t *TestObject) resetMsg(origins []types.Address, receipts []types.KVPair) {
	smc := t.obj.sdk

	inR := smc.Message().InputReceipts()
	if receipts != nil {
		inR = append(inR, receipts...)
	}

	smc.(*sdkimpl.SmartContract).SetMessage(object.NewMessage(smc,
		smc.Message().Contract(),
		smc.Message().MethodID(),
		smc.Message().Items(),
		smc.Message().Sender().Address(),
		smc.Message().Payer().Address(),
		origins,
		inR))
}

//transfer This is a method of TestObject
func (t *TestObject) transfer(args ...interface{}) *TestObject {
	utest.Assert(utest.GetFlag())

	contract := t.obj.sdk.Message().Contract()
	t.obj.sdk.Message().Sender().(*object.Account).SetSMC(t.obj.sdk)
	if utest.Transfer(t.obj.sdk.Message().Sender(), contract.Account().Address(), args...).ErrorCode != types.CodeOK {
		return nil
	}
	t.resetMsg(t.obj.sdk.Message().Origins(), t.obj.sdk.Message().(*object.Message).OutputReceipts())
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
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Set(data)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//Get This is a method of TestObject
func (t *TestObject) Get() (result0 uint64, err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	result0 = t.obj.Get()
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}
