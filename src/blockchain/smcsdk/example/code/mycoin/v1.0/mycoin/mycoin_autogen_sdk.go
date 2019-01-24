package mycoin

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of Mycoin
func (m *Mycoin) SetSdk(sdk sdk.ISmartContract) {
	m.sdk = sdk
}

//GetSdk This is a method of Mycoin
func (m *Mycoin) GetSdk() sdk.ISmartContract {
	return m.sdk
}
