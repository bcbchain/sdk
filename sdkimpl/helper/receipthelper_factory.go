package helper

import (
	"github.com/bcbchain/sdk/sdk"
)

// NewReceiptHelper factory method to create IReceiptHelper
func NewReceiptHelper(smc sdk.ISmartContract) sdk.IReceiptHelper {
	o := ReceiptHelper{}
	o.SetSMC(smc)
	return &o
}
