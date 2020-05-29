package helper

import (
	"github.com/AeReach/sdk/sdk"
)

// NewStateHelper factory method for create IStateHelper
func NewStateHelper(smc sdk.ISmartContract) sdk.IStateHelper {
	o := StateHelper{}
	o.SetSMC(smc)
	return &o
}
