package directsale

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
	contractName       = "directsale" //contract name
	contractMethods    = []string{"SetManager([]types.Address)", "SetSettings(string)", "Register(types.Address)", "ConfigApplication(types.Address,[]string)", "ConfigInvestmentRatio(types.Address,int64)", "Invest()", "Settle(string,map[string]int64)", "WithdrawFunds(string,types.Address,bn.Number)", "RefundInvestment(string,bn.Number)"}
	contractInterfaces = []string{"Income(types.Address,bn.Number,string)", "Pay(types.Address,string,bn.Number,bn.Number,[]bn.Number,[]types.Address,string)", "Refund(types.Address,string,bn.Number,[]bn.Number,[]types.Address,string)"}
	orgID              = "orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer"
)

//TestObject This is a struct for test
type TestObject struct {
	obj *DirectSale
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
	return &TestObject{&DirectSale{sdk: utest.UTP.ISmartContract}}
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

//SetManager This is a method of TestObject
func (t *TestObject) SetManager(managerAddrs []types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetManager(managerAddrs)
	utest.Commit()
	return
}

//SetSettings This is a method of TestObject
func (t *TestObject) SetSettings(setting string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetSettings(setting)
	utest.Commit()
	return
}

//Register This is a method of TestObject
func (t *TestObject) Register(refAddr types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Register(refAddr)
	utest.Commit()
	return
}

//ConfigApplication This is a method of TestObject
func (t *TestObject) ConfigApplication(salerAddr types.Address, contractNames []string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.ConfigApplication(salerAddr, contractNames)
	utest.Commit()
	return
}

//ConfigInvestmentRatio This is a method of TestObject
func (t *TestObject) ConfigInvestmentRatio(salerAddr types.Address, ratio int64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.ConfigInvestmentRatio(salerAddr, ratio)
	utest.Commit()
	return
}

//Invest This is a method of TestObject
func (t *TestObject) Invest() (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Invest()
	utest.Commit()
	return
}

//Settle This is a method of TestObject
func (t *TestObject) Settle(salerAddr string, exchangeRates map[string]int64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Settle(salerAddr, exchangeRates)
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

//RefundInvestment This is a method of TestObject
func (t *TestObject) RefundInvestment(tokenName string, refundAmount bn.Number) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.RefundInvestment(tokenName, refundAmount)
	utest.Commit()
	return
}

//Income This is a method of TestObject
func (t *TestObject) Income(salerAddr types.Address, lockAmount bn.Number, note string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Income(salerAddr, lockAmount, note)
	utest.Commit()
	return
}

//Pay This is a method of TestObject
func (t *TestObject) Pay(salerAddr types.Address, tokenName string, incomeAmount, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Pay(salerAddr, tokenName, incomeAmount, unlockAmount, amountList, toAddrList, note)
	utest.Commit()
	return
}

//Refund This is a method of TestObject
func (t *TestObject) Refund(salerAddr types.Address, tokenName string, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.Refund(salerAddr, tokenName, unlockAmount, amountList, toAddrList, note)
	utest.Commit()
	return
}
