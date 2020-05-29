package helper

import (
	"github.com/AeReach/sdk/sdk"
)

// NewGenesisHelper factory method to create IGenesisHelper
func NewGenesisHelper(smc sdk.ISmartContract) sdk.IGenesisHelper {
	o := GenesisHelper{}
	o.SetSMC(smc)
	return &o
}
