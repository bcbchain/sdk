package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
	"strconv"
)

const (
	PERMILLI       = 1000 // permillage
	MAXSETTLECOUNT = 100  // max settle count per time
	MAXREFUNDCOUNT = 100  // max refund count per time
)

func (sg *Excellencies) checkPossibleWinAmount(tokenName string, amount, possibleWinAmount, feeNum bn.Number, roundInfo *RoundInfo) {

	sdk.Require(
		possibleWinAmount.Sub(feeNum).Cmp(roundInfo.Setting.Settings[tokenName].MaxProfit) <= 0,
		types.ErrInvalidParameter,
		"PossibleWinAmount should be smaller than maxProfit",
	)

	return
}

func (sg *Excellencies) checkBetInfo(tokenName string, amount bn.Number, betDataSlice []BetData, roundInfo *RoundInfo) (modes []BaseMode) {

	// tokenName 必须在settings里设置
	setting, ok := roundInfo.Setting.Settings[tokenName]
	sdk.Require(
		ok == true,
		types.ErrInvalidParameter,
		"TokenName is not support in contract ",
	)

	//长度检查
	sdk.Require(
		len(betDataSlice) != 0,
		types.ErrInvalidParameter,
		"betData slice cannot empty ",
	)

	//检查下注总金额是否匹配
	totalAmount := bn.N(0)
	forx.Range(betDataSlice, func(_, data BetData) bool {
		totalAmount = totalAmount.Add(data.BetAmount)

		//amount 精度应该等于PERMILLI cong
		sdk.Require(
			data.BetAmount.ModI(PERMILLI).CmpI(0) == 0,
			types.ErrInvalidParameter,
			"Amount per bet mode accuracy should be equal"+strconv.Itoa(PERMILLI),
		)
		return forx.Continue
	})

	sdk.Require(
		amount.Cmp(totalAmount) == 0,
		types.ErrInvalidParameter,
		"Amount param is invalid",
	)

	sdk.Require(
		amount.Cmp(setting.MinLimit) >= 0,
		types.ErrInvalidParameter,
		tokenName+":"+"Amount should be bigger than MinLimit",
	)

	sdk.Require(
		amount.Cmp(setting.MaxLimit) <= 0,
		types.ErrInvalidParameter,
		tokenName+":"+"Amount should be smaller than MaxLimit",
	)

	modes = CreateBetMode(betDataSlice)

	return
}

func (sg *Excellencies) checkSettings(newSettings string) (resultSettings *MapSetting) {

	resultSettings = &MapSetting{}
	jsonErr := jsoniter.Unmarshal([]byte(newSettings), resultSettings)

	sdk.RequireNotError(jsonErr, types.ErrInvalidParameter)
	//check settings empty
	sdk.Require(
		len(resultSettings.Settings) != 0,
		types.ErrInvalidParameter,
		"Settings can not empty ")

	//check BetExpirationBlocks
	sdk.Require(
		resultSettings.BetExpirationBlocks > 0 && resultSettings.BetExpirationBlocks <= (1<<31-1),
		types.ErrInvalidParameter,
		"BetExpirationBlocks must be bigger than zero and smaller than "+strconv.Itoa(1<<31-1),
	)

	//check Settings
	forx.Range(resultSettings.Settings, func(tokenName string, value Setting) bool {

		//代币名称不为空
		sdk.Require(
			tokenName != "",
			types.ErrInvalidParameter,
			"TokenName should not be empty",
		)

		//代币必须为链上支持的币
		token := sg.sdk.Helper().TokenHelper().TokenOfName(tokenName)

		sdk.Require(
			token != nil,
			types.ErrInvalidParameter,
			"TokenName error",
		)

		sdk.Require(
			value.MaxProfit.CmpI(0) >= 0,
			types.ErrInvalidParameter,
			tokenName+":"+"MaxProfit can not be negative",
		)

		sdk.Require(
			value.MaxLimit.CmpI(0) > 0,
			types.ErrInvalidParameter,
			tokenName+":"+"MaxLimit must be bigger than zero",
		)

		sdk.Require(
			value.MinLimit.CmpI(0) > 0 && value.MaxLimit.Cmp(value.MinLimit) > 0,
			types.ErrInvalidParameter,
			tokenName+":"+"MinLimit must be bigger than zero and small than MaxLimit",
		)

		sdk.Require(
			value.SendToCltRatio >= 0 && value.SendToCltRatio < PERMILLI,
			types.ErrInvalidParameter,
			tokenName+":"+"SendToCltRatio must be bigger than zero and smaller than "+strconv.Itoa(PERMILLI),
		)

		sdk.Require(
			value.FeeRatio > 0 && value.FeeRatio < PERMILLI,
			types.ErrInvalidParameter,
			tokenName+":"+"FeeRatio must be bigger than zero and  smaller than "+strconv.Itoa(PERMILLI),
		)

		sdk.Require(
			value.FeeMiniNum.CmpI(0) > 0 && value.FeeMiniNum.Cmp(value.MinLimit.MulI(value.FeeRatio).DivI(PERMILLI)) <= 0,
			types.ErrInvalidParameter,
			tokenName+":"+"FeeMiniNum must bigger than zero and smaller than resultSettings.MinLimit * resultSettings.FeeRatio/"+strconv.Itoa(PERMILLI),
		)

		return true
	})

	return
}

func (sg *Excellencies) checkRecFeeInfo(infos []RecFeeInfo) {
	sdk.Require(len(infos) > 0,
		types.ErrInvalidParameter, "The length of RecvFeeInfos must be larger than zero")

	allRatio := int64(0)
	forx.Range(infos, func(_, info RecFeeInfo) {
		sdk.Require(info.RecFeeRatio > 0,
			types.ErrInvalidParameter, "ratio must be larger than zero")
		sdk.RequireAddress(info.RecFeeAddr)
		sdk.Require(info.RecFeeAddr != sg.sdk.Message().Contract().Account(),
			types.ErrInvalidParameter, "address cannot be contract account address")

		allRatio += info.RecFeeRatio
	})

	//The allocation ratio set must add up to 1000
	sdk.Require(allRatio <= 1000,
		types.ErrInvalidParameter, "The sum of ratio must be less or equal 1000")
}

func (sg *Excellencies) checkRoundInfoPlaceBet(roundInfo *RoundInfo) {

	sdk.Require(
		len(roundInfo.Commit) != 0,
		types.ErrInvalidParameter,
		"The round is not exist",
	)

	//要结算的轮的下注高度要小于结算高度
	sdk.Require(
		roundInfo.FirstBlockHeight < sg.sdk.Block().Height(),
		types.ErrInvalidParameter,
		"SettleBet block can not be in the same block as placeBet or before.",
	)

	if roundInfo.SettledCount == 0 {
		//首次结算，轮信息不能过期
		sdk.Require(
			roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks > sg.sdk.Block().Height(),
			types.ErrInvalidParameter,
			"This round is time out",
		)
	} else {
		sdk.Require(
			roundInfo.BetCount != roundInfo.SettledCount,
			types.ErrInvalidParameter,
			"This round is complete",
		)

	}

	return
}

//Transfer to fee's receiving address
func (sg *Excellencies) calcCltAndRecAmount(feeNum bn.Number, setting Setting, clt types.Address) (addressList []types.Address, amountList []bn.Number) {

	amountList = make([]bn.Number, 0)
	addressList = make([]types.Address, 0)

	sentToCltFee := feeNum.MulI(setting.SendToCltRatio).DivI(PERMILLI)

	//转给clt
	if sentToCltFee.CmpI(0) > 0 {
		amountList = append(amountList, sentToCltFee)
		addressList = append(addressList, clt)
	}

	recAmount := feeNum.Sub(sentToCltFee)
	if recAmount.CmpI(0) > 0 {
		recfees := sg._recFeeInfo()
		forx.Range(recfees, func(_, rec RecFeeInfo) bool {
			if recAmount.MulI(rec.RecFeeRatio).DivI(PERMILLI).CmpI(0) > 0 {
				amountList = append(amountList, recAmount.MulI(rec.RecFeeRatio).DivI(PERMILLI))
				addressList = append(addressList, rec.RecFeeAddr)
			}

			return forx.Continue
		})
	}

	return

}

func (sg *Excellencies) checkCommit(commitStr string) {

	tmp := sg._chkRoundInfo(commitStr)
	// 新的commit 允许下注
	if tmp == false {
		return

	} else {
		//当前轮已开奖或已过期，不允许下注
		roundInfo := sg._roundInfo(commitStr)
		sdk.Require(
			(len(roundInfo.Players) != 0 || len(roundInfo.Banker.Cards) != 0) ||
				(sg.sdk.Block().Height() <= roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks),
			types.ErrInvalidParameter,
			"Current round is lottery or timeout",
		)
	}

	return
}

func (sg *Excellencies) checkRoundInfoRefund(roundInfo *RoundInfo) {

	sdk.Require(len(roundInfo.Commit) != 0,
		types.ErrInvalidParameter,
		"The round is not exist",
	)
	sdk.Require(len(roundInfo.Players) == 0 && len(roundInfo.Banker.Cards) == 0,
		types.ErrInvalidParameter,
		"The round was run a lottery",
	)
	sdk.Require(roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks < sg.sdk.Block().Height(),
		types.ErrInvalidParameter,
		"The round never timeout",
	)
	sdk.Require(roundInfo.RefundCount < roundInfo.BetCount,
		types.ErrInvalidParameter,
		"The round is refunded",
	)

	return
}
