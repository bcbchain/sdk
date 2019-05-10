package directsale

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
	"strconv"
)

//DirectSale This is struct of contract
//@:contract:directsale
//@:version:1.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:8867990fbcf775ae279637ef8bdb8f300f61edfc85be70df739f5c3460788822
type DirectSale struct {
	sdk sdk.ISmartContract

	//This is setting which is to store in db
	//@:public:store
	settings Settings

	//@:public:store
	salers map[types.Address]Saler

	//@:public:store
	managers []types.Address

	//@:public:store
	salerToApps map[types.Address]map[string]SalerApp //key1 = 会员地址， key2 = 合约名称

	//@:public:store
	global map[string]Global //key = tokenName
}

//@:public:receipt
type receipt interface {
	emitRegister(refAddr types.Address, newSaler Saler)
	emitSetManager(managersAddrs []types.Address)
	emitSetSettings(setting Settings)
	emitConfigApplication(salerAddr types.Address, appInfoList map[string]map[string]AppStatByToken)
	emitConfigInvestmentRatio(saleAddrs types.Address, ratio int64)
	emitSettle(settleInfo *Settle)
	emitIncome(contractName, tokenName string, amount bn.Number, salerAddr types.Address, Note string)
	emitPay(contractNme, tokenName string, unlockAmount bn.Number, amount []bn.Number, toAddrList []types.Address, salerAddr types.Address, note string)
	emitRefund(contractNme, tokenName string, unlockAmount bn.Number, amount []bn.Number, toAddrList []types.Address, salerAddr types.Address, note string)
}

//InitChain Constructor of this DirectSale
//@:constructor
func (ds *DirectSale) InitChain() {
	//init settings data
	settings := Settings{}

	settings.Foster = 50
	settings.ShareReward = 300
	settings.PointReward = bn.N(3E11)
	settings.SalerExpense = bn.N(1E13)
	//单位是 (cong)
	settings.StarSettings = []starReward{
		{bn.N(1E15), 800},                // (0-100]万 	     1星
		{bn.N(2E15), 810},                // (100-200]万    2星
		{bn.N(5E15), 820},                // (200-500]万    3星
		{bn.N(1E16), 830},                // (500-1000]万   4星
		{bn.N(2E16), 840},                // (1000-2000]万  5星
		{bn.N(5E16), 850},                // (2000-5000]万  6星
		{bn.N(1E17), 860},                // (5000-10000]万 7星
		{bn.N(9223372036854775807), 870}, // (10000-int64最大值]万 8星 大约92亿
	}

	ds._setSettings(settings)

	managers := make([]types.Address, 0)
	managers = append(managers, ds.sdk.Message().Contract().Owner())
	ds._setManagers(managers)

}

//Set manager
//@:public:method:gas[5000]
func (ds *DirectSale) SetManager(managerAddrs []types.Address) {
	//only contract owner just can do it
	sdk.RequireOwner()

	sdk.Require(
		len(managerAddrs) != 0,
		types.ErrInvalidParameter,
		"managerAddrs not exist")

	//地址有效性和要求检查
	//for _, managerAddr := range managerAddrs {
	//	sdk.RequireAddress(managerAddr)
	//}
	forx.Range(managerAddrs, func(i int, managerAddr types.Address) bool {
		sdk.RequireAddress(managerAddr)
		return true
	})

	ds._setManagers(managerAddrs)

	//fire event
	ds.emitSetManager(managerAddrs)
	return
}

//Set setting
//@:public:method:gas[5000]
func (ds *DirectSale) SetSettings(setting string) {
	//检查是否设置管理者
	ds.checkSender()

	sdk.Require(
		len(setting) != 0,
		types.ErrInvalidParameter,
		"setting not exist")

	//检查settings
	resultSettings := ds.checkSettings(setting)

	ds._setSettings(resultSettings)

	ds.emitSetSettings(resultSettings)
	return
}

//Save data and issue rewards
//@:public:method:gas[500]
func (ds *DirectSale) Register(refAddr types.Address) {
	//检查是否满足注册条件注册
	ds.checkTransferAccept()
	newSaler := Saler{RefAddr: refAddr}

	newSalerAddr := ds.sdk.Message().Sender().Address()
	ds.sdk.Message().Contract().Address()
	//没有推荐人
	if refAddr == "" {
		ds._setSalers(newSalerAddr, newSaler)
		//fire event
		ds.emitRegister(refAddr, newSaler)
		return
	}

	//判断直推人是否存在
	sdk.Require(
		ds._chkSalers(refAddr),
		types.ErrInvalidParameter,
		"refAddr does not exist")
	ds.saveSaler(newSalerAddr, newSaler, refAddr)

	//保存直推人信息
	refSaler := ds._salers(refAddr)
	refSaler.RefCounts += 1
	ds._setSalers(refAddr, refSaler)

	//发放直推奖
	setting := ds._settings()
	recommendReward := setting.SalerExpense.MulI(int64(setting.ShareReward)).DivI(int64(PROPORTION))
	ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account()).TransferByName(TOKENNAME, refAddr, recommendReward)

	//发放见点奖
	ds.pointReward(newSalerAddr, setting.PointReward)
	// fire event
	ds.emitRegister(refAddr, newSaler)
	return
}

//Config about apps
//@:public:method:gas[5000]
func (ds *DirectSale) ConfigApplication(salerAddr types.Address, contractNames []string) {
	//检查权限
	ds.checkSender()

	//check salerAddr
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrInvalidParameter,
		"saler is nothingness in contract",
	)

	newContracts := ds.updateContractNames(salerAddr, contractNames)

	salerApps := make(map[string]map[string]AppStatByToken, 0) //key1 == contract // key2 == tokenName
	//for _, contract := range newContracts {
	//	if ds._chkSalerToApps(salerAddr, contract) {
	//		apps := ds._salerToApps(salerAddr, contract)
	//		salerApps[contract] = apps.AppInfoMap
	//	} else {
	//		salerApps[contract] = map[string]AppStatByToken{}
	//	}
	//}
	forx.Range(newContracts, func(i int, contract string) bool {
		if ds._chkSalerToApps(salerAddr, contract) {
			apps := ds._salerToApps(salerAddr, contract)
			salerApps[contract] = apps.AppInfoMap
		} else {
			salerApps[contract] = map[string]AppStatByToken{}
		}
		return true
	})
	//fire event
	ds.emitConfigApplication(salerAddr, salerApps)
}

//Config about Investment
//@:public:method:gas[5000]
func (ds *DirectSale) ConfigInvestmentRatio(salerAddr types.Address, ratio int64) {
	//检查权限
	ds.checkSender()

	//check saler
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrInvalidParameter,
		"saler is nothingness in contract",
	)

	//check ratio
	sdk.Require(
		ratio > 0 && ratio <= 1000,
		types.ErrInvalidParameter,
		"InvestmentRatio must be bigger than zero smaller or equal"+strconv.Itoa(PROPORTION),
	)

	saler := ds._salers(salerAddr)
	sdk.Require(
		saler.InvestRate != ratio,
		types.ErrInvalidParameter,
		"ratio do not equal saler.InvestRate ",
	)

	//get salerInfo
	isInvestZero := true
	forx.Range(saler.Accounts, func(tokenName string, value SalerAccount) bool {
		sdk.Require(
			value.FlowAmount.CmpI(0) == 0,
			types.ErrInvalidParameter,
			"saler's Accounts FlowAmount must be zero",
		)
		if value.InvestAmount.CmpI(0) != 0 {
			isInvestZero = false
		}

		return true
	})

	//if all tokenName->InvestAmount == 0
	if isInvestZero {
		saler.InvestRate = ratio
		ds._setSalers(salerAddr, saler)
		//fire event
		ds.emitConfigInvestmentRatio(salerAddr, ratio)
		return
	}

	//set ratio
	ds.configRatio(salerAddr, saler, int64(ratio))

	//fire event
	ds.emitConfigInvestmentRatio(salerAddr, ratio)

	return
}

//@:public:method:gas[500]
func (ds *DirectSale) Invest() {
	//check sender
	salerAddr := ds.sdk.Message().Sender().Address()
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrNoAuthorization,
		"saler is nothingness in contract",
	)

	// get transfer info
	correct := false
	var transferReceipt std.Transfer
	//for _, investTransfer := range ds.sdk.Message().GetTransferToMe() {
	//
	//	if ds._chkSalers(investTransfer.From) &&
	//		investTransfer.To == ds.sdk.Message().Contract().Account() {
	//		transferReceipt.Token = investTransfer.Token
	//		transferReceipt.Value = investTransfer.Value
	//		correct = true
	//		break
	//	}
	//}
	forx.Range(ds.sdk.Message().GetTransferToMe(), func(i int, investTransfer *std.Transfer) bool {
		if ds._chkSalers(investTransfer.From) &&
			investTransfer.To == ds.sdk.Message().Contract().Account() {
			transferReceipt.Token = investTransfer.Token
			transferReceipt.Value = investTransfer.Value
			correct = true
			return forx.Break
		}
		return true
	})
	sdk.Require(
		correct,
		types.ErrNoAuthorization,
		"TnvestTransfer information error",
	)

	saler := ds._salers(transferReceipt.From)
	sdk.Require(
		saler.InvestRate == 0,
		types.ErrInvalidParameter,
		"Saler InvestRate can not be zero",
	)

	tokenName := ds.sdk.Helper().TokenHelper().TokenOfAddress(transferReceipt.Token).Name()

	// check tokenName
	account := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	sdk.Require(
		account.BalanceOfName(tokenName).CmpI(0) == 0,
		types.ErrInsufficientBalance,
		"Contract is not support this token",
	)

	global := Global{}
	if !ds._chkGlobal(tokenName) {
		global = Global{bn.N(0), bn.N(0)}
	} else {
		global = ds._global(tokenName)
	}

	// check contract balance
	totalAdd := transferReceipt.Value.MulI(PROPORTION).DivI(int64(saler.InvestRate))
	contractInvest := totalAdd.Sub(transferReceipt.Value)
	sdk.Require(
		account.BalanceOfName(tokenName).Sub(global.PoolBalance).Cmp(contractInvest) >= 0,
		types.ErrInsufficientBalance,
		"Insufficient balance",
	)

	//update global
	global.PoolBalance = global.PoolBalance.Add(totalAdd)
	ds._setGlobal(tokenName, global)

	//update saler account
	temp, ok := saler.Accounts[tokenName]
	if !ok {
		temp.InvestAmount = bn.N(0)
		temp.IncomeBalance = bn.N(0)
		temp.FlowAmount = bn.N(0)
		temp.FlowAmount = bn.N(0)
		temp.TotalBalance = bn.N(0)
	}
	temp.InvestAmount = temp.InvestAmount.Add(totalAdd)
	temp.TotalBalance = temp.TotalBalance.Add(totalAdd)
	saler.Accounts[tokenName] = temp

	ds._setSalers(salerAddr, saler)

	return
}

//@:public:method:gas[500]
func (ds *DirectSale) Settle(salerAddr string, exchangeRates map[string]int64) {
	//check sender
	ds.checkSender()

	//check saler
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrInvalidParameter,
		"saler is nothingness in contract",
	)

	saler := ds._salers(salerAddr)

	totalFlowAmount := bn.N(0)
	forx.Range(saler.Accounts, func(tokenName string, value SalerAccount) bool {
		val, ok := exchangeRates[tokenName]
		sdk.Require(
			ok,
			types.ErrInvalidParameter,
			"The saler could not support this tokenName:"+tokenName,
		)
		totalFlowAmount = totalFlowAmount.Add(value.FlowAmount.Mul(bn.N(val).DivI(EXCHANGERATIO)))

		return true
	})

	starLevel, forsterFoeward, salerForward := ds.settleFlow(salerAddr, saler, totalFlowAmount)
	settle := Settle{
		salerAddr,
		exchangeRates,
		starLevel,
		forsterFoeward,
		salerForward,
		totalFlowAmount,
	}
	//fire event
	ds.emitSettle(&settle)
}

//@:public:method:gas[5000]
func (ds *DirectSale) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {
	//check sender
	sdk.RequireOwner()

	//check address
	sdk.RequireAddress(beneficiary)

	//check tokenName
	// check tokenName
	account := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	sdk.Require(
		account.BalanceOfName(tokenName).CmpI(0) == 0,
		types.ErrInsufficientBalance,
		"Contract is not support this token",
	)

	if !ds._chkGlobal(tokenName) {
		sdk.Require(withdrawAmount.CmpI(0) > 0 &&
			withdrawAmount.Cmp(account.BalanceOfName(tokenName)) <= 0,
			types.ErrInvalidParameter,
			"WithdrawAmount is error",
		)

		account.TransferByName(tokenName, beneficiary, withdrawAmount)
		return
	}

	global := ds._global(tokenName)

	//check withdrawAmount
	sdk.Require(
		withdrawAmount.CmpI(0) > 0 &&
			withdrawAmount.Cmp(account.BalanceOfName(tokenName).Sub(global.PoolBalance)) <= 0,
		types.ErrInvalidParameter,
		"WithdrawAmount is error",
	)

	//transfer
	account.TransferByName(tokenName, ds.sdk.Message().Contract().Owner(), withdrawAmount)

	return
}

//@:public:method:gas[500]
func (ds *DirectSale) RefundInvestment(tokenName string, refundAmount bn.Number) {
	//check sender
	salerAddr := ds.sdk.Message().Sender().Address()
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrNoAuthorization,
		"Sender must be contract saler",
	)

	saler := ds._salers(salerAddr)

	//check tokenName
	salerAccount, ok := saler.Accounts[tokenName]
	sdk.Require(
		ok,
		types.ErrInvalidParameter,
		"This saler is not support the tokenName",
	)

	tempAmount := salerAccount.TotalBalance.Sub(salerAccount.LockedAmount)
	//check refundAmount
	sdk.Require(
		refundAmount.CmpI(0) > 0 &&
			refundAmount.Cmp(tempAmount.MulI(int64(saler.InvestRate)).DivI(PROPORTION)) <= 0,
		types.ErrInvalidParameter,
		"RefundAmount should be bigger than zero and "+
			"smaller than salerAccount.TotalBalance.Sub(salerAccount.LockedAmount).MulI(int64(saler.InvestRate)).DivI("+
			strconv.Itoa(PROPORTION)+")",
	)

	total := refundAmount.DivI(PROPORTION).MulI(int64(saler.InvestRate))

	//update global
	global := ds._global(tokenName)
	global.PoolBalance = global.PoolBalance.Sub(total)
	ds._setGlobal(tokenName, global)

	//update saler account
	salerAcc := SalerAccount{}
	temp := saler.Accounts[tokenName]
	salerAcc.InvestAmount = temp.InvestAmount.Sub(total)
	salerAcc.IncomeBalance = temp.IncomeBalance
	salerAcc.LockedAmount = temp.LockedAmount
	salerAcc.FlowAmount = temp.FlowAmount
	salerAcc.TotalBalance = temp.TotalBalance.Sub(total)
	saler.Accounts[tokenName] = salerAcc

	ds._setSalers(salerAddr, saler)

	account := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	account.TransferByName(tokenName, salerAddr, refundAmount)

	return

}

//--------------------------interface------------------------------//

//@:public:interface:gas[100]
func (ds *DirectSale) Income(salerAddr types.Address, lockAmount bn.Number, note string) {

	//check saler
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrInvalidParameter,
		"saler is nothingness in contract",
	)

	// get transfer info
	correct := false
	var transferReceipt std.Transfer
	//for _, investTransfer := range ds.sdk.Message().GetTransferToMe() {
	//	if investTransfer.To == ds.sdk.Message().Contract().Account() &&
	//		investTransfer.From == salerAddr {
	//		transferReceipt.Token = investTransfer.Token
	//		transferReceipt.Value = investTransfer.Value
	//		correct = true
	//		break
	//	}
	//}
	forx.Range(ds.sdk.Message().GetTransferToMe(), func(i int, investTransfer *std.Transfer) bool {
		if investTransfer.To == ds.sdk.Message().Contract().Account() &&
			investTransfer.From == salerAddr {
			transferReceipt.Token = investTransfer.Token
			transferReceipt.Value = investTransfer.Value
			correct = true
			return forx.Break
		}
		return true
	})

	sdk.Require(
		correct,
		types.ErrNoAuthorization,
		"TnvestTransfer information error",
	)

	//check salerApps
	correct = false
	saler := ds._salers(salerAddr)
	originsLength := len(ds.sdk.Message().Origins())
	contractName := ds.sdk.Helper().ContractHelper().ContractOfAddress(ds.sdk.Message().Origins()[originsLength-1]).Name()
	//for _, app := range saler.ContractNames {
	//	if app == contractName {
	//		correct = true
	//		break
	//	}
	//}
	forx.Range(saler.ContractNames, func(i int, app string) bool {
		if app == contractName {
			correct = true
			return forx.Break
		}
		return true
	})
	sdk.Require(
		correct,
		types.ErrInvalidParameter,
		"This saler can not support"+contractName,
	)

	tokenName := ds.sdk.Helper().TokenHelper().TokenOfAddress(transferReceipt.Token).Name()
	amount := transferReceipt.Value

	sdk.Require(
		amount.Cmp(lockAmount) <= 0,
		types.ErrInvalidParameter,
		"LockAmount should be bigger than transfer amount",
	)

	//update Saler，Apps & global
	ds.incomeUpdate(salerAddr, saler, contractName, tokenName, lockAmount, amount)

	//fire event
	ds.emitIncome(contractName, tokenName, amount, salerAddr, note)

	return
}

//@:public:interface:gas[100]
func (ds *DirectSale) Pay(salerAddr types.Address, tokenName string, incomeAmount, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string) {

	saler, salerAcc, totalAmount := ds.checkArguments(salerAddr, tokenName, unlockAmount, amountList, toAddrList)
	//check income amount
	sdk.Require(
		incomeAmount.CmpI(0) > 0 && incomeAmount.Cmp(salerAcc.TotalBalance) <= 0,
		types.ErrInvalidParameter,
		"IncomeAmount should be bigger than zero and smaller than totalBalance",
	)

	//transfer to each addr
	account := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	//for index, amount := range amountList {
	//	account.TransferByName(tokenName, toAddrList[index], amount)
	//}
	forx.Range(amountList, func(index int, amount bn.Number) bool {
		account.TransferByName(tokenName, toAddrList[index], amount)
		return true
	})
	//update saler
	salerAcc.LockedAmount = salerAcc.LockedAmount.Sub(unlockAmount)
	salerAcc.TotalBalance = salerAcc.TotalBalance.Sub(totalAmount)
	salerAcc.FlowAmount = salerAcc.FlowAmount.Add(incomeAmount)
	salerAcc.IncomeBalance = salerAcc.IncomeBalance.Sub(incomeAmount)
	saler.Accounts[tokenName] = salerAcc
	ds._setSalers(salerAddr, saler)

	//global
	global := ds._global(tokenName)
	global.TotalLocked = global.TotalLocked.Sub(unlockAmount)
	global.PoolBalance = global.TotalLocked.Sub(totalAmount)
	ds._setGlobal(tokenName, global)

	// apps
	originsLength := len(ds.sdk.Message().Origins())
	contractName := ds.sdk.Helper().ContractHelper().ContractOfAddress(ds.sdk.Message().Origins()[originsLength-1]).Name()
	salerApp := ds._salerToApps(salerAddr, contractName)
	salerAppToken := salerApp.AppInfoMap[tokenName]
	salerAppToken.LockedAmount = salerAppToken.LockedAmount.Sub(unlockAmount)
	salerAppToken.PayCount += 1
	salerAppToken.FlowAmount = salerAppToken.FlowAmount.Add(incomeAmount)
	salerApp.AppInfoMap[tokenName] = salerAppToken
	ds._setSalerToApps(salerAddr, contractName, salerApp)

	//fire event
	ds.emitPay(contractName, tokenName, unlockAmount, amountList, toAddrList, salerAddr, note)
	return

}

//@:public:interface:gas[100]
func (ds *DirectSale) Refund(salerAddr types.Address, tokenName string, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string) {

	saler, salerAcc, totalAmount := ds.checkArguments(salerAddr, tokenName, unlockAmount, amountList, toAddrList)

	//transfer to each addr
	account := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	//for index, amount := range amountList {
	//	account.TransferByName(tokenName, toAddrList[index], amount)
	//}
	forx.Range(amountList, func(index int, amount bn.Number) bool {
		account.TransferByName(tokenName, toAddrList[index], amount)
		return true
	})
	//update saler
	salerAcc.LockedAmount = salerAcc.LockedAmount.Sub(unlockAmount)
	salerAcc.TotalBalance = salerAcc.TotalBalance.Sub(totalAmount)
	salerAcc.IncomeBalance = salerAcc.IncomeBalance.Sub(totalAmount)
	saler.Accounts[tokenName] = salerAcc
	ds._setSalers(salerAddr, saler)

	//global
	global := ds._global(tokenName)
	global.TotalLocked = global.TotalLocked.Sub(unlockAmount)
	global.PoolBalance = global.TotalLocked.Sub(totalAmount)
	ds._setGlobal(tokenName, global)

	// apps
	originsLength := len(ds.sdk.Message().Origins())
	contractName := ds.sdk.Helper().ContractHelper().ContractOfAddress(ds.sdk.Message().Origins()[originsLength-1]).Name()
	salerApp := ds._salerToApps(salerAddr, contractName)
	salerAppToken := salerApp.AppInfoMap[tokenName]
	salerAppToken.LockedAmount = salerAppToken.LockedAmount.Sub(unlockAmount)
	salerApp.AppInfoMap[tokenName] = salerAppToken
	ds._setSalerToApps(salerAddr, contractName, salerApp)

	//fire event
	ds.emitRefund(salerAddr, tokenName, unlockAmount, amountList, toAddrList, salerAddr, note)
}
