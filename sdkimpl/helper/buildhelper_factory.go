package helper

import "github.com/bcbchain/sdk/sdk"

// NewBuildHelper factory method to create IBuildHelper object
func NewBuildHelper(smc sdk.ISmartContract) sdk.IBuildHelper {
	o := BuildHelper{}
	o.SetSMC(smc)
	return &o
}
