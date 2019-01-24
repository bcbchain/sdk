package mystorage

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of MyStorage
func (ms *MyStorage) SetSdk(sdk sdk.ISmartContract) {
	ms.sdk = sdk
}

//GetSdk This is a method of MyStorage
func (ms *MyStorage) GetSdk() sdk.ISmartContract {
	return ms.sdk
}
