package helper

import (
	"github.com/bcbchain/sdk/sdk"
	"github.com/bcbchain/sdk/sdk/types"
	"github.com/bcbchain/sdk/sdkimpl"
	"github.com/bcbchain/sdk/sdkimpl/object"
)

// AccountHelper the account helper information
type AccountHelper struct {
	smc sdk.ISmartContract //指向智能合约API对象指针
}

var _ sdk.IAccountHelper = (*AccountHelper)(nil)
var _ sdkimpl.IAcquireSMC = (*AccountHelper)(nil)

// SMC get smc object
func (ah *AccountHelper) SMC() sdk.ISmartContract { return ah.smc }

// SetSMC set smc object
func (ah *AccountHelper) SetSMC(smc sdk.ISmartContract) { ah.smc = smc }

// AccountOf get account with address
func (ah *AccountHelper) AccountOf(addr types.Address) sdk.IAccount {
	return object.NewAccount(ah.smc, addr)
}

// AccountOfPubKey create account with pubKey
func (ah *AccountHelper) AccountOfPubKey(pubKey types.PubKey) sdk.IAccount {
	return object.NewAccountWithPubKey(ah.smc, pubKey)
}
