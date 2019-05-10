package directsale

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

var _ receipt = (*DirectSale)(nil)

//emitRegister This is a method of DirectSale
func (ds *DirectSale) emitRegister(refAddr types.Address, newSaler Saler) {
	type register struct {
		RefAddr  types.Address `json:"refAddr"`
		NewSaler Saler         `json:"newSaler"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		register{
			RefAddr:  refAddr,
			NewSaler: newSaler,
		},
	)
}

//emitSetManager This is a method of DirectSale
func (ds *DirectSale) emitSetManager(managersAddrs []types.Address) {
	type setManager struct {
		ManagersAddrs []types.Address `json:"managersAddrs"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		setManager{
			ManagersAddrs: managersAddrs,
		},
	)
}

//emitSetSettings This is a method of DirectSale
func (ds *DirectSale) emitSetSettings(setting Settings) {
	type setSettings struct {
		Setting Settings `json:"setting"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		setSettings{
			Setting: setting,
		},
	)
}

//emitConfigApplication This is a method of DirectSale
func (ds *DirectSale) emitConfigApplication(salerAddr types.Address, appInfoList map[string]map[string]AppStatByToken) {
	type configApplication struct {
		SalerAddr   types.Address                        `json:"salerAddr"`
		AppInfoList map[string]map[string]AppStatByToken `json:"appInfoList"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		configApplication{
			SalerAddr:   salerAddr,
			AppInfoList: appInfoList,
		},
	)
}

//emitConfigInvestmentRatio This is a method of DirectSale
func (ds *DirectSale) emitConfigInvestmentRatio(saleAddrs types.Address, ratio int64) {
	type configInvestmentRatio struct {
		SaleAddrs types.Address `json:"saleAddrs"`
		Ratio     int64         `json:"ratio"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		configInvestmentRatio{
			SaleAddrs: saleAddrs,
			Ratio:     ratio,
		},
	)
}

//emitSettle This is a method of DirectSale
func (ds *DirectSale) emitSettle(settleInfo *Settle) {
	type settle struct {
		SettleInfo *Settle `json:"settleInfo"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		settle{
			SettleInfo: settleInfo,
		},
	)
}

//emitIncome This is a method of DirectSale
func (ds *DirectSale) emitIncome(contractName, tokenName string, amount bn.Number, salerAddr types.Address, Note string) {
	type income struct {
		ContractName string        `json:"contractName"`
		TokenName    string        `json:"tokenName"`
		Amount       bn.Number     `json:"amount"`
		SalerAddr    types.Address `json:"salerAddr"`
		Note         string        `json:"Note"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		income{
			ContractName: contractName,
			TokenName:    tokenName,
			Amount:       amount,
			SalerAddr:    salerAddr,
			Note:         Note,
		},
	)
}

//emitPay This is a method of DirectSale
func (ds *DirectSale) emitPay(contractNme, tokenName string, unlockAmount bn.Number, amount []bn.Number, toAddrList []types.Address, salerAddr types.Address, note string) {
	type pay struct {
		ContractNme  string          `json:"contractNme"`
		TokenName    string          `json:"tokenName"`
		UnlockAmount bn.Number       `json:"unlockAmount"`
		Amount       []bn.Number     `json:"amount"`
		ToAddrList   []types.Address `json:"toAddrList"`
		SalerAddr    types.Address   `json:"salerAddr"`
		Note         string          `json:"note"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		pay{
			ContractNme:  contractNme,
			TokenName:    tokenName,
			UnlockAmount: unlockAmount,
			Amount:       amount,
			ToAddrList:   toAddrList,
			SalerAddr:    salerAddr,
			Note:         note,
		},
	)
}

//emitRefund This is a method of DirectSale
func (ds *DirectSale) emitRefund(contractNme, tokenName string, unlockAmount bn.Number, amount []bn.Number, toAddrList []types.Address, salerAddr types.Address, note string) {
	type refund struct {
		ContractNme  string          `json:"contractNme"`
		TokenName    string          `json:"tokenName"`
		UnlockAmount bn.Number       `json:"unlockAmount"`
		Amount       []bn.Number     `json:"amount"`
		ToAddrList   []types.Address `json:"toAddrList"`
		SalerAddr    types.Address   `json:"salerAddr"`
		Note         string          `json:"note"`
	}

	ds.sdk.Helper().ReceiptHelper().Emit(
		refund{
			ContractNme:  contractNme,
			TokenName:    tokenName,
			UnlockAmount: unlockAmount,
			Amount:       amount,
			ToAddrList:   toAddrList,
			SalerAddr:    salerAddr,
			Note:         note,
		},
	)
}
