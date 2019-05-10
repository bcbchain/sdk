package mydonation

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of Mydonation
func (m *Mydonation) SetSdk(sdk sdk.ISmartContract) {
	m.sdk = sdk
}

//GetSdk This is a method of Mydonation
func (m *Mydonation) GetSdk() sdk.ISmartContract {
	return m.sdk
}
