package excellencies

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

//InterfacedirectsaleStub This is a interface stub of directsale
type InterfacedirectsaleStub struct {
}

//directsaleStub This is method of Excellencies
func (e *Excellencies) directsale() *InterfacedirectsaleStub {
	return &InterfacedirectsaleStub{}
}

//Income This is a method of InterfacedirectsaleStub
func (is *InterfacedirectsaleStub) Income(salerAddr types.Address, lockAmount bn.Number, note string) {
	return
}

//Pay This is a method of InterfacedirectsaleStub
func (is *InterfacedirectsaleStub) Pay(salerAddr types.Address, tokenName string, incomeAmount, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string) {
	return
}

//Refund This is a method of InterfacedirectsaleStub
func (is *InterfacedirectsaleStub) Refund(salerAddr types.Address, tokenName string, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string) {
	return
}
