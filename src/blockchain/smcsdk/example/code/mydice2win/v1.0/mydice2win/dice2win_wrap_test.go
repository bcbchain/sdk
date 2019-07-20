package mydice2win

import (
	"blockchain/smcsdk/common/gls"
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl"
	"blockchain/smcsdk/sdkimpl/object"
	"blockchain/smcsdk/sdkimpl/sdkhelper"
	"blockchain/smcsdk/utest"
	"bytes"
	"common/jsoniter"
	"fmt"
	"reflect"
	"strings"
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

type setSecretSigner struct {
	NewSecretSigner types.PubKey `json:"newSecretSigner"`
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
	return &TestObject{obj: &Dice2Win{sdk: utest.UTP.ISmartContract}}
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

	var err types.Error
	gls.Mgr.SetValues(gls.Values{gls.SDKKey: smc}, func() {
		err = f(t)
	})

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

//SetSecretSigner This is a method of TestObject
func (t *TestObject) SetSecretSigner(newSecretSigner types.PubKey) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.SetSecretSigner(newSecretSigner)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//SetSettings This is a method of TestObject
func (t *TestObject) SetSettings(newSettingsStr string) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.SetSettings(newSettingsStr)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//SetRecvFeeInfos This is a method of TestObject
func (t *TestObject) SetRecvFeeInfos(recvFeeInfosStr string) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.SetRecvFeeInfos(recvFeeInfosStr)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//WithdrawFunds This is a method of TestObject
func (t *TestObject) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.WithdrawFunds(tokenName, beneficiary, withdrawAmount)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//PlaceBet This is a method of TestObject
func (t *TestObject) PlaceBet(betMask bn.Number, modulo, commitLastBlock int64, commit, signData []byte, refAddress types.Address) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.PlaceBet(betMask, modulo, commitLastBlock, commit, signData, refAddress)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//SettleBet This is a method of TestObject
func (t *TestObject) SettleBet(reveal []byte) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.SettleBet(reveal)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
}

//RefundBet This is a method of TestObject
func (t *TestObject) RefundBet(commit []byte) (err types.Error) {
	utest.Assert(utest.GetFlag())

	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.UTP.ISmartContract = t.obj.sdk
	utest.NextBlock(1)
	t.obj.RefundBet(commit)
	t.resetMsg(t.obj.sdk.Message().Origins(),
		t.obj.sdk.Message().(*object.Message).OutputReceipts())

	return
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
