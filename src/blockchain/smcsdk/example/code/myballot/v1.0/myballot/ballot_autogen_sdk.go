package myballot

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk: This is a method of Ballot
func (b *Ballot) SetSdk(sdk sdk.ISmartContract) {
	b.sdk = sdk
}

//GetSdk: This is a method of Ballot
func (b *Ballot) GetSdk() sdk.ISmartContract {
	return b.sdk
}
