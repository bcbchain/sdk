package helper

import (
	"github.com/AeReach/sdk/sdk"
)

// NewAccountHelper factory method for AccountHelper
func NewAccountHelper(smc sdk.ISmartContract) sdk.IAccountHelper {
	o := AccountHelper{}
	o.SetSMC(smc)
	return &o
}
