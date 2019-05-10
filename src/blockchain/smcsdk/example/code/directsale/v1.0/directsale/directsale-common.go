package directsale

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
)

func (ds *DirectSale) checkSender() {
	//检查权限
	managerList := ds._managers()
	var isManager bool
	//for _, manager := range managerList {
	//	if ds.sdk.Message().Sender().Address() == manager {
	//		return
	//	}
	//}
	forx.Range(managerList, func(i int, manager types.Address) bool {
		if ds.sdk.Message().Sender().Address() == manager {
			isManager = true
			return forx.Break
		}
		return true
	})

	sdk.Require(
		isManager,
		types.ErrNoAuthorization,
		"if contract have Manager,contract sender must be manager",
	)
}
func (ds *DirectSale) checkSettings(_settings string) (resultSettings Settings) {

	resultSettings = Settings{}
	jsonErr := jsoniter.Unmarshal([]byte(_settings), &resultSettings)
	sdk.RequireNotError(jsonErr, types.ErrInvalidParameter)

	//check settings
	//for index, starSetting := range resultSettings.StarSettings {
	//	if index == len(resultSettings.StarSettings)-1 {
	//		break
	//	}
	//	sdk.Require(
	//		starSetting.MaxIncome.CmpI(0) >= 0 && starSetting.MaxIncome.Cmp(resultSettings.StarSettings[index+1].MaxIncome) <= 0 && starSetting.DividendRatio < resultSettings.StarSettings[index+1].DividendRatio,
	//		types.ErrInvalidParameter,
	//		"StarSetting parameter error",
	//	)
	//}
	forx.Range(resultSettings.StarSettings, func(index int, starSetting starReward) bool {
		if index == len(resultSettings.StarSettings)-1 {
			return forx.Break
		}
		sdk.Require(
			starSetting.MaxIncome.CmpI(0) >= 0 && starSetting.MaxIncome.Cmp(resultSettings.StarSettings[index+1].MaxIncome) <= 0 && starSetting.DividendRatio < resultSettings.StarSettings[index+1].DividendRatio,
			types.ErrInvalidParameter,
			"StarSetting parameter error",
		)
		return true
	})
	sdk.Require(
		resultSettings.Foster >= 0 && resultSettings.Foster <= 1000,
		types.ErrInvalidParameter,
		"Foster parameter error")
	sdk.Require(
		resultSettings.ShareReward >= 0 && resultSettings.ShareReward <= 1000,
		types.ErrInvalidParameter,
		"Foster parameter error")
	sdk.Require(
		resultSettings.SalerExpense.Cmp(resultSettings.PointReward) >= 0,
		types.ErrInvalidParameter,
		"SalerExpense must be more than PointReward")
	return
}

//check contract names
func checkContractNames(names []string) map[string]string {
	contractMap := make(map[string]string, 0)
	//for _, name := range names {
	//	// check Name != ""
	//	sdk.Require(
	//		name != "",
	//		types.ErrInvalidParameter,
	//		"Contract Name could not be empty",
	//	)
	//	contractMap[name] = ""
	//}
	forx.Range(names, func(i int, name string) bool {
		// check Name != ""
		sdk.Require(
			name != "",
			types.ErrInvalidParameter,
			"Contract Name could not be empty",
		)
		contractMap[name] = ""
		return true
	})
	return contractMap
}

func (ds *DirectSale) checkArguments(salerAddr types.Address, tokenName string, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address) (Saler, SalerAccount, bn.Number) {

	sdk.RequireAddress(salerAddr)

	//check saler
	sdk.Require(
		ds._chkSalers(salerAddr),
		types.ErrInvalidParameter,
		"saler is nothingness in contract",
	)

	saler := ds._salers(salerAddr)
	salerAcc, ok := saler.Accounts[tokenName]

	//check tokenName
	sdk.Require(
		ok,
		types.ErrInvalidParameter,
		"This saler can not support this"+tokenName,
	)

	//check unlockAmount
	sdk.Require(
		unlockAmount.CmpI(0) > 0 && unlockAmount.Cmp(salerAcc.LockedAmount) <= 0,
		types.ErrInvalidParameter,
		"UnlockAmount should be bigger than zero and smaller than lockedAmount",
	)

	//check amountList
	totalAmount := bn.N(0)
	//for _, amount := range amountList {
	//	sdk.Require(
	//		amount.CmpI(0) > 0,
	//		types.ErrInvalidParameter,
	//		"Each amount in amountList should be bigger than zero",
	//	)
	//	totalAmount = totalAmount.Add(amount)
	//
	//}
	forx.Range(amountList, func(i int, amount bn.Number) bool {
		sdk.Require(
			amount.CmpI(0) > 0,
			types.ErrInvalidParameter,
			"Each amount in amountList should be bigger than zero",
		)
		totalAmount = totalAmount.Add(amount)
		return true
	})
	sdk.Require(
		totalAmount.Cmp(unlockAmount) <= 0,
		types.ErrInvalidParameter,
		"AmountList sum should be bigger than 0 and smaller than totalBalance",
	)
	sdk.Require(
		totalAmount.CmpI(0) > 0 && totalAmount.Cmp(salerAcc.TotalBalance) <= 0,
		types.ErrInvalidParameter,
		"AmountList sum should be bigger than 0 and smaller than totalBalance",
	)

	//check toAddrList
	//for _, toAddr := range toAddrList {
	//	sdk.Require(
	//		toAddr != "",
	//		types.ErrInvalidParameter,
	//		"Each toAddr in toAddrList can not be empty",
	//	)
	//}
	forx.Range(toAddrList, func(i int, toAddr types.Address) bool {
		sdk.Require(
			toAddr != "",
			types.ErrInvalidParameter,
			"Each toAddr in toAddrList can not be empty",
		)
		return true
	})
	sdk.Require(
		len(amountList) == len(toAddrList),
		types.ErrInvalidParameter,
		"AmountList length should be equal toAddrList length",
	)

	return saler, salerAcc, totalAmount
}
