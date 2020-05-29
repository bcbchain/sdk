package llstate

import (
	"github.com/bcbchain/sdk/sdk"
	"github.com/bcbchain/sdk/sdkimpl"
)

// NewLowLevelSDB factory method to create LowLevelSDB
func NewLowLevelSDB(smc sdk.ISmartContract, transID, txID int64) sdkimpl.ILowLevelSDB {
	o := LowLevelSDB{cache: make(map[string][]byte)}
	o.Init(transID, txID)
	o.SetSMC(smc)

	return &o
}
