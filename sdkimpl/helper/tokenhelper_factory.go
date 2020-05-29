package helper

import (
	"github.com/AeReach/sdk/sdk"
)

// NewTokenHelper factory method to create TokenHelper
func NewTokenHelper(smc sdk.ISmartContract) sdk.ITokenHelper {
	o := TokenHelper{}
	o.SetSMC(smc)
	return &o
}
