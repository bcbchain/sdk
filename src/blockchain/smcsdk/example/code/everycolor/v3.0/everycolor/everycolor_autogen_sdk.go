package everycolor

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk - set sdk
func (e *Everycolor) SetSdk(sdk sdk.ISmartContract) {
	e.sdk = sdk
}

//GetSdk - get sdk
func (e *Everycolor) GetSdk() sdk.ISmartContract {
	return e.sdk
}
