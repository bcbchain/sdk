package mydonation

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl"
	"blockchain/smcsdk/sdkimpl/llstate"
	"blockchain/smcsdk/sdkimpl/object"
	"blockchain/smcsdk/sdkimpl/sdkhelper"
	"blockchain/smcsdk/utest"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/check.v1"
)

var (
	contractName       = "mydonation" //contract name
	contractMethods    = []string{"AddDonee(types.Address)", "DelDonee(types.Address)", "Donate(types.Address)", "Transfer(types.Address,bn.Number)"}
	contractInterfaces = []string{}
	orgID              = "orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer"
)

//Test This is a function
func Test(t *testing.T) { check.TestingT(t) }

//MySuite This is a struct
type MySuite struct{}

var _ = check.Suite(&MySuite{})

//TestObject This is a struct for test
type TestObject struct {
	obj        *Mydonation
	origins    []types.Address
	isSetBlock bool
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
	return &TestObject{obj: &Mydonation{sdk: utest.UTP.ISmartContract}}
}

//transfer This is a method of TestObject
func (t *TestObject) transfer(args ...interface{}) *TestObject {
	utest.Assert(utest.GetFlag())

	contract := t.obj.sdk.Message().Contract()
	t.obj.sdk.Message().Sender().(*object.Account).SetSMC(t.obj.sdk)
	if utest.TransferEx(t.obj.sdk.Message().Sender(), contract.Account().Address(), args...).ErrorCode != types.CodeOK {
		return nil
	}
	t.resetMsg(t.obj.sdk.Message().Origins(), t.obj.sdk.Message().(*object.Message).OutputReceipts())
	return t
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
	newll := llstate.NewLowLevelSDB(t.obj.sdk, 0, 0)
	t.obj.sdk.(*sdkimpl.SmartContract).SetLlState(newll)
}

// runf This is a method of TestObject
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
}

func (t *TestObject) assertReceipt(index int, value interface{}) {
	outReceipts := t.obj.sdk.Message().(*object.Message).InputReceipts()

	utest.Assert(index < len(outReceipts) && index >= 0)

	receipt := outReceipts[index]

	name := receiptName(value)
	utest.Assert(strings.HasSuffix(string(receipt.Key), name))

	var r std.Receipt
	err := jsoniter.Unmarshal(receipt.Value, &r)
	utest.Assert(err == nil)

	res, err := jsoniter.Marshal(value)
	utest.Assert(err == nil)

	utest.Assert(bytes.Equal(res, r.Bytes))
}

func (t *TestObject) assertReceiptNil() {
	utest.Assert(len(t.obj.sdk.Message().InputReceipts()) == 0)
}

func receiptName(receipt interface{}) string {
	typeOfInterface := reflect.TypeOf(receipt).String()

	if strings.HasPrefix(typeOfInterface, "std.") {
		prefixLen := len("std.")
		return "std::" + strings.ToLower(typeOfInterface[prefixLen:prefixLen+1]) + typeOfInterface[prefixLen+1:]
	}

	return typeOfInterface
}

//Set blockInfo
func (t *TestObject) SetNextBlock(block std.Block) {
	utest.NextBlockEx(1,
		block.Height,
		block.Time,
		block.LastFee,
		block.BlockHash,
		block.DataHash,
		block.LastBlockHash,
		block.LastCommitHash,
		block.LastAppHash,
		block.ProposerAddress,
		block.RewardAddress,
		block.RandomNumber,
	)

	t.isSetBlock = true
	return
}

//InitChain This is a method of TestObject
func (t *TestObject) InitChain() {
	utest.NextBlock(1)
	t.obj.InitChain()
	utest.Commit()
	return
}

//AddDonee This is a method of TestObject
func (t *TestObject) AddDonee(donee types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	if !t.isSetBlock {
		utest.NextBlock(1)
	}
	t.obj.AddDonee(donee)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())
	t.isSetBlock = false
	return
}

//DelDonee This is a method of TestObject
func (t *TestObject) DelDonee(donee types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	if !t.isSetBlock {
		utest.NextBlock(1)
	}
	t.obj.DelDonee(donee)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())
	t.isSetBlock = false
	return
}

//Donate This is a method of TestObject
func (t *TestObject) Donate(donee types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	if !t.isSetBlock {
		utest.NextBlock(1)
	}
	t.obj.Donate(donee)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())
	t.isSetBlock = false
	return
}

//Transfer This is a method of TestObject
func (t *TestObject) Transfer(donee types.Address, value bn.Number) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	if !t.isSetBlock {
		utest.NextBlock(1)
	}
	t.obj.Transfer(donee, value)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())
	t.isSetBlock = false
	return
}
