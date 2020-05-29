package object

import (
	"github.com/bcbchain/sdk/sdk"
	"github.com/bcbchain/sdk/sdk/types"
)

// NewTx factory method for create tx with all tx's property
func NewTx(smc sdk.ISmartContract, note string, gasLimit int64, gasLeft int64, txHash []byte, sender types.Address) sdk.ITx {
	signer := NewAccount(smc, sender)

	o := &Tx{
		note:     note,
		gasLimit: gasLimit,
		gasLeft:  gasLeft,
		txHash:   txHash,
		signer:   signer,
	}
	o.SetSMC(smc)

	return o
}
