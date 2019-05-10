package directsale

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of DirectSale
func (ds *DirectSale) SetSdk(sdk sdk.ISmartContract) {
	ds.sdk = sdk
}

//GetSdk This is a method of DirectSale
func (ds *DirectSale) GetSdk() sdk.ISmartContract {
	return ds.sdk
}