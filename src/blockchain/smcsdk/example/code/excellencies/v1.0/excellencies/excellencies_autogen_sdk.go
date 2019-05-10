package excellencies

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of Excellencies
func (e *Excellencies) SetSdk(sdk sdk.ISmartContract) {
	e.sdk = sdk
}

//GetSdk This is a method of Excellencies
func (e *Excellencies) GetSdk() sdk.ISmartContract {
	return e.sdk
}