package helper

import "github.com/bcbchain/sdk/sdk"

// NewIBCStubHelper factory method to create IBCStubHelper
func NewIBCStubHelper(smc sdk.ISmartContract) sdk.IIBCStubHelper {
	o := IBCStubHelper{}
	o.SetSMC(smc)
	return &o
}
