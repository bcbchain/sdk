package directsale

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
)

func (ds *DirectSale) saveSaler(newSalerAddr types.Address, newSaler Saler, refAddr types.Address) {
	var startSalerAddrList []types.Address

	//先存入直推人地址
	startSalerAddrList = append(startSalerAddrList, refAddr)

	//层次内先左后右依次向下层遍历
	//for _, startSalerAddr := range startSalerAddrList {
	//	startSaler := ds._salers(startSalerAddr)
	//	if len(startSaler.Sons) < 2 {
	//		startSaler.Sons = append(startSaler.Sons, newSalerAddr)
	//		newSaler.Parent = startSalerAddr
	//		break
	//	} else {
	//		startSalerAddrList = append(startSalerAddrList, startSaler.Sons...)
	//	}
	//}
	forx.Range(startSalerAddrList, func(i int, startSalerAddr types.Address) bool {
		startSaler := ds._salers(startSalerAddr)
		if len(startSaler.Sons) < 2 {
			startSaler.Sons = append(startSaler.Sons, newSalerAddr)
			newSaler.Parent = startSalerAddr
			return forx.Break
		} else {
			startSalerAddrList = append(startSalerAddrList, startSaler.Sons...)
		}
		return true
	})
	ds._setSalers(newSalerAddr, newSaler)
	//fire event
	ds.emitRegister(refAddr, newSaler)
}

func (ds *DirectSale) pointReward(newSalerAddr types.Address, pointRewardAmount bn.Number) {
	reward := map[int]int64{
		1:  1,
		2:  1,
		3:  1,
		4:  1,
		5:  2,
		6:  3,
		7:  4,
		8:  5,
		9:  6,
		10: 7,
		11: 8,
	}
	startSalerAddr := ds._salers(newSalerAddr).Parent

	//向上遍历并发放见点奖
	//for index := 1; index <= MAXLEVEL && startSalerAddr != ""; index++ {
	//	parentSaler := ds._salers(startSalerAddr)
	//	if parentSaler.RefCounts >= reward[index] {
	//		ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account()).TransferByName(TOKENNAME, startSalerAddr, pointRewardAmount)
	//	}
	//	startSalerAddr = parentSaler.Parent
	//}
	forx.Range(func(index int) bool {
		if index <= MAXLEVEL && startSalerAddr != "" {
			parentSaler := ds._salers(startSalerAddr)
			if parentSaler.RefCounts >= reward[index] {
				ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account()).TransferByName(TOKENNAME, startSalerAddr, pointRewardAmount)
			}
			startSalerAddr = parentSaler.Parent
			index++
			return true
		}
		return forx.Break
	})

}

func (ds *DirectSale) checkTransferAccept() {
	//取得交易信息
	transferReceipt := ds.sdk.Message().GetTransferToMe()
	var isApplicant bool
	//for _, registerTransfer := range transferReceipt {
	//	token := ds.sdk.Helper().TokenHelper().TokenOfAddress(registerTransfer.Token)
	//
	//	//检查代币名称，检查转账金额，转账人是申请人
	//	if token.Name() == TOKENNAME &&
	//		registerTransfer.Value.Cmp(ds.settings.SalerExpense) == 0 &&
	//		registerTransfer.From == ds.sdk.Message().Sender().Address() {
	//		return
	//	}
	//}
	forx.Range(transferReceipt, func(i int, registerTransfer *std.Transfer) bool {
		token := ds.sdk.Helper().TokenHelper().TokenOfAddress(registerTransfer.Token)
		//检查代币名称，检查转账金额，转账人是申请人
		if token.Name() == TOKENNAME &&
			registerTransfer.Value.Cmp(ds.settings.SalerExpense) == 0 &&
			registerTransfer.From == ds.sdk.Message().Sender().Address() {
			isApplicant = true
			return forx.Break
		}
		return true
	})
	//sdk.Require(
	//	false,
	//	types.ErrNoAuthorization,
	//	"Registration information error",
	//)
	sdk.Require(
		isApplicant,
		types.ErrNoAuthorization,
		"Registration information error",
	)
}

func (ds *DirectSale) updateContractNames(salerAddr types.Address, contractNames []string) []string {

	//check contractNames
	contractMap := checkContractNames(contractNames)

	// first time to set contract names for saler
	saler := ds._salers(salerAddr)
	if len(saler.ContractNames) == 0 {
		saler.ContractNames = contractNames
		ds._setSalers(salerAddr, saler)
		return saler.ContractNames
	}

	// find out will del contract Names
	tempDelSlice := make([]string, 0)
	//for _, salerContract := range saler.ContractNames {
	//	_, ok := contractMap[salerContract]
	//	if ok == false {
	//		tempDelSlice = append(tempDelSlice, salerContract)
	//	}
	//}
	forx.Range(saler.ContractNames, func(i int, salerContract string) bool {
		_, ok := contractMap[salerContract]
		if ok == false {
			tempDelSlice = append(tempDelSlice, salerContract)
		}
		return true
	})
	// judge del contarct
	appToken := AppStatByToken{}
	//for _, temp := range tempDelSlice {
	//	if ds._chkSalerToApps(salerAddr, temp) {
	//		salerApp := ds._salerToApps(salerAddr, temp)
	//		isDelContract := true
	//		forx.Range(salerApp, func(tokenName string, value AppStatByToken) bool {
	//			if value.LockedAmount.CmpI(0) != 0 {
	//				contractMap[temp] = ""
	//				isDelContract = false
	//				return forx.Break
	//			}
	//			appToken.LockedAmount = bn.N(0)
	//			appToken.FlowAmount = bn.N(0)
	//			appToken.PayCount = 0
	//			salerApp.AppInfoMap[tokenName] = appToken
	//
	//			return true
	//		})
	//
	//		if isDelContract == true {
	//			ds._setSalerToApps(salerAddr, temp, salerApp)
	//		}
	//
	//	}
	//}
	forx.Range(tempDelSlice, func(i int, temp string) bool {
		if ds._chkSalerToApps(salerAddr, temp) {
			salerApp := ds._salerToApps(salerAddr, temp)
			isDelContract := true
			forx.Range(salerApp, func(tokenName string, value AppStatByToken) bool {
				if value.LockedAmount.CmpI(0) != 0 {
					contractMap[temp] = ""
					isDelContract = false
					return forx.Break
				}
				appToken.LockedAmount = bn.N(0)
				appToken.FlowAmount = bn.N(0)
				appToken.PayCount = 0
				salerApp.AppInfoMap[tokenName] = appToken

				return true
			})

			if isDelContract == true {
				ds._setSalerToApps(salerAddr, temp, salerApp)
			}

		}
		return true
	})
	// add contarct
	forx.Range(contractMap, func(name, value string) bool {
		saler.ContractNames = append(saler.ContractNames, name)

		return true
	})

	ds._setSalers(salerAddr, saler)
	return saler.ContractNames
}

func (ds *DirectSale) configRatio(salerAddr types.Address, saler Saler, ratio int64) {
	salerAccount := SalerAccount{}
	account := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	forx.Range(saler.Accounts, func(tokenName string, value SalerAccount) bool {

		ownerInvest := value.TotalBalance.Sub(value.IncomeBalance).MulI(PROPORTION - saler.InvestRate).DivI(PROPORTION)
		newOwnerInvest := ownerInvest.MulI(PROPORTION).DivI(PROPORTION - ratio).Sub(value.TotalBalance.Sub(value.IncomeBalance).Sub(ownerInvest))
		global := ds._global(tokenName)
		balance := account.BalanceOfName(tokenName)
		increase := newOwnerInvest.Sub(ownerInvest)

		// 合约owner 出资金额增加
		if increase.CmpI(0) > 0 {
			if balance.Sub(global.PoolBalance).Sub(increase).CmpI(0) < 0 {
				// 合约没有足够得余额去支付，减少会员的出资额
				salerInvest := value.TotalBalance.Sub(value.IncomeBalance).MulI(saler.InvestRate).DivI(PROPORTION)
				newSalerInvest := salerInvest.MulI(PROPORTION).DivI(ratio).Sub(salerInvest)
				increase = newSalerInvest.Sub(newSalerInvest)

				//判断减少后的金额与锁定金额相比
				sdk.Require(
					value.InvestAmount.Sub(increase).Cmp(value.LockedAmount) >= 0,
					types.ErrInvalidParameter,
					"New ratio lead to new Invest Amount has not enough to pay ",
				)
				//Send the funds to saler.
				account.TransferByName(tokenName, salerAddr, bn.NB(increase.V.Abs(increase.V)))
			}
		} else {
			//判断减少后的金额与锁定金额相比
			sdk.Require(
				value.InvestAmount.Add(increase).Cmp(value.LockedAmount) >= 0,
				types.ErrInvalidParameter,
				"New ratio lead to new Invest Amount has not enough to pay ",
			)
		}

		// salerAccount
		salerAccount.InvestAmount = value.InvestAmount.Add(increase)
		salerAccount.TotalBalance = value.TotalBalance.Add(increase)
		salerAccount.LockedAmount = value.LockedAmount
		salerAccount.FlowAmount = value.FlowAmount
		salerAccount.IncomeBalance = value.IncomeBalance
		saler.Accounts[tokenName] = salerAccount

		// global
		global.PoolBalance = global.PoolBalance.Add(increase)
		ds._setGlobal(tokenName, global)

		return true
	})

	ds._setSalers(salerAddr, saler)
}

func (ds *DirectSale) settleFlow(salerAddr types.Address, saler Saler, totalFlowAmount bn.Number) (starLevel int, toSalerMap, toRefMap map[string]bn.Number) {
	//get contract account
	accountCon := ds.sdk.Helper().AccountHelper().AccountOf(ds.sdk.Message().Contract().Account())
	settings := ds._settings()
	//get start
	salerAccount := SalerAccount{}
	starLevel = 0
	toSalerMap = map[string]bn.Number{}
	toRefMap = map[string]bn.Number{}
	//for index, value := range settings.StarSettings {
	//	isbiggest := false
	//	if index == len(settings.StarSettings)-1 && totalFlowAmount.Cmp(value.MaxIncome) > 0 {
	//		isbiggest = true
	//	}
	//	if totalFlowAmount.Cmp(value.MaxIncome) <= 0 || isbiggest {
	//
	//		//according to the setting rewards to transfer
	//		starLevel = index + 1
	//		forx.Range(saler.Accounts, func(tokenName string, salerAcc SalerAccount) bool {
	//
	//			//calc this saler earnings
	//			earnings := salerAcc.TotalBalance.Sub(salerAcc.IncomeBalance).Sub(salerAcc.InvestAmount)
	//			tempMod := bn.N(0)
	//			if earnings.CmpI(0) > 0 {
	//				sdk.Require(
	//					salerAcc.TotalBalance.Sub(earnings).Sub(salerAcc.IncomeBalance).Cmp(salerAcc.LockedAmount) >= 0,
	//					types.ErrInsufficientBalance,
	//					"lockedAmount should be smaller than new SalerAccount.TotalBalance",
	//				)
	//				//给saler 发会员奖
	//				toSaler := earnings.MulI(saler.InvestRate).DivI(PROPORTION).MulI(settings.StarSettings[index].DividendRatio).DivI(PROPORTION)
	//				if toSaler.CmpI(PROPORTION) >= 0 {
	//					toSaler = toSaler.Sub(toSaler.ModI(PROPORTION))
	//					accountCon.TransferByName(tokenName, salerAddr, toSaler)
	//					toSalerMap[tokenName] = toSaler
	//				}
	//				tempMod = tempMod.Add(toSaler.ModI(PROPORTION))
	//
	//				//给上级发抚育奖
	//				if saler.RefAddr != "" {
	//					toRef := earnings.MulI(saler.InvestRate).DivI(PROPORTION).MulI(settings.Foster).DivI(PROPORTION)
	//					if toRef.CmpI(PROPORTION) >= 0 {
	//						toRef = toRef.Sub(toRef.ModI(PROPORTION))
	//						accountCon.TransferByName(tokenName, salerAddr, toRef)
	//						toRefMap[tokenName] = toRef
	//					}
	//					tempMod = tempMod.Add(toRef.ModI(PROPORTION))
	//				}
	//
	//				// 结算完当月流水清零 update SalerAccount
	//				salerAccount.FlowAmount = bn.N(0)
	//				salerAccount.InvestAmount = salerAcc.InvestAmount
	//				salerAccount.IncomeBalance = salerAcc.IncomeBalance
	//				salerAccount.TotalBalance = salerAcc.TotalBalance.Sub(earnings).Add(tempMod)
	//
	//				//save in saler
	//				saler.Accounts[tokenName] = salerAccount
	//
	//				//update new Global
	//				global := ds._global(tokenName)
	//				global.PoolBalance = global.PoolBalance.Sub(earnings)
	//				ds._setGlobal(tokenName, global)
	//			}
	//
	//			return true
	//		})
	//		break
	//	}
	//}
	forx.Range(settings.StarSettings, func(index int, value starReward) bool {
		isbiggest := false
		if index == len(settings.StarSettings)-1 && totalFlowAmount.Cmp(value.MaxIncome) > 0 {
			isbiggest = true
		}
		if totalFlowAmount.Cmp(value.MaxIncome) <= 0 || isbiggest {

			//according to the setting rewards to transfer
			starLevel = index + 1
			forx.Range(saler.Accounts, func(tokenName string, salerAcc SalerAccount) bool {

				//calc this saler earnings
				earnings := salerAcc.TotalBalance.Sub(salerAcc.IncomeBalance).Sub(salerAcc.InvestAmount)
				tempMod := bn.N(0)
				if earnings.CmpI(0) > 0 {
					sdk.Require(
						salerAcc.TotalBalance.Sub(earnings).Sub(salerAcc.IncomeBalance).Cmp(salerAcc.LockedAmount) >= 0,
						types.ErrInsufficientBalance,
						"lockedAmount should be smaller than new SalerAccount.TotalBalance",
					)
					//给saler 发会员奖
					toSaler := earnings.MulI(saler.InvestRate).DivI(PROPORTION).MulI(settings.StarSettings[index].DividendRatio).DivI(PROPORTION)
					if toSaler.CmpI(PROPORTION) >= 0 {
						toSaler = toSaler.Sub(toSaler.ModI(PROPORTION))
						accountCon.TransferByName(tokenName, salerAddr, toSaler)
						toSalerMap[tokenName] = toSaler
					}
					tempMod = tempMod.Add(toSaler.ModI(PROPORTION))

					//给上级发抚育奖
					if saler.RefAddr != "" {
						toRef := earnings.MulI(saler.InvestRate).DivI(PROPORTION).MulI(settings.Foster).DivI(PROPORTION)
						if toRef.CmpI(PROPORTION) >= 0 {
							toRef = toRef.Sub(toRef.ModI(PROPORTION))
							accountCon.TransferByName(tokenName, salerAddr, toRef)
							toRefMap[tokenName] = toRef
						}
						tempMod = tempMod.Add(toRef.ModI(PROPORTION))
					}

					// 结算完当月流水清零 update SalerAccount
					salerAccount.FlowAmount = bn.N(0)
					salerAccount.InvestAmount = salerAcc.InvestAmount
					salerAccount.IncomeBalance = salerAcc.IncomeBalance
					salerAccount.TotalBalance = salerAcc.TotalBalance.Sub(earnings).Add(tempMod)

					//save in saler
					saler.Accounts[tokenName] = salerAccount

					//update new Global
					global := ds._global(tokenName)
					global.PoolBalance = global.PoolBalance.Sub(earnings)
					ds._setGlobal(tokenName, global)
				}

				return true
			})
			return forx.Break
		}
		return true
	})
	ds._setSalers(salerAddr, saler)
	return
}

func (ds *DirectSale) incomeUpdate(salerAddr types.Address, saler Saler, contractName, tokenName string, lockAmount, amount bn.Number) {

	salerAcc, ok := saler.Accounts[tokenName]
	sdk.Require(
		ok,
		types.ErrNoAuthorization,
		"This saler is not support "+tokenName,
	)

	//record amount & check lockedAmount
	sdk.Require(
		lockAmount.CmpI(0) > 0 &&
			lockAmount.Cmp(salerAcc.TotalBalance.Add(amount)) <= 0,
		types.ErrInvalidParameter,
		"LockAmount is error ",
	)
	//update saler
	salerAcc.LockedAmount = salerAcc.LockedAmount.Add(lockAmount)
	salerAcc.IncomeBalance = salerAcc.IncomeBalance.Add(amount)
	salerAcc.TotalBalance = salerAcc.TotalBalance.Add(amount)
	saler.Accounts[tokenName] = salerAcc

	ds._setSalers(salerAddr, saler)

	//update salerApp
	salerApp := SalerApp{}
	if !ds._chkSalerToApps(salerAddr, contractName) {
		salerApp = SalerApp{make(map[string]AppStatByToken, 0)}
	} else {
		salerApp = ds._salerToApps(salerAddr, contractName)
	}

	salerAppToken, ok := salerApp.AppInfoMap[tokenName]
	if !ok {
		salerAppToken.FlowAmount = bn.N(0)
		salerAppToken.LockedAmount = bn.N(0)
		salerAppToken.PayCount = 0
	}
	salerAppToken.LockedAmount = salerAppToken.LockedAmount.Add(lockAmount)
	salerAppToken.FlowAmount = bn.N(0)
	salerApp.AppInfoMap[tokenName] = salerAppToken
	ds._setSalerToApps(salerAddr, contractName, salerApp)

	//update global
	global := Global{}
	if !ds._chkGlobal(tokenName) {
		global.PoolBalance = bn.N(0)
		global.TotalLocked = bn.N(0)
	}
	global.TotalLocked = global.TotalLocked.Add(lockAmount)
	global.PoolBalance = global.PoolBalance.Add(amount)
	ds._setGlobal(tokenName, global)
}
