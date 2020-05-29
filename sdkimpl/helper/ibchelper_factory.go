package helper

import "github.com/bcbchain/sdk/sdk"

// NewReceiptHelper factory method to create IReceiptHelper
func NewIBCHelper(smc sdk.ISmartContract) sdk.IIBCHelper {
	o := IBCHelper{}
	o.SetSMC(smc)
	return &o
}
