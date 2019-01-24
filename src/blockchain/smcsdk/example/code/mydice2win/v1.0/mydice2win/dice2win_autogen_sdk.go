package mydice2win

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of Dice2Win
func (dw *Dice2Win) SetSdk(sdk sdk.ISmartContract) {
	dw.sdk = sdk
}

//GetSdk This is a method of Dice2Win
func (dw *Dice2Win) GetSdk() sdk.ISmartContract {
	return dw.sdk
}
