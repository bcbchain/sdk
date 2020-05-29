package helper

import (
	"github.com/bcbchain/sdk/sdk"
)

// NewBlockChainHelper factory method fro create IBlockChainHelper object
func NewBlockChainHelper(smc sdk.ISmartContract) sdk.IBlockChainHelper {
	o := BlockChainHelper{}
	o.SetSMC(smc)
	return &o
}
