package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
	"encoding/binary"

	"strconv"
)

func (ec *Everycolor) TokenAddress(tokenName string) types.Address {
	tokenAddress := ec.sdk.Helper().TokenHelper().TokenOfName(tokenName)
	return tokenAddress.Address()
}

//转账给fee的接收地址
func (e *Everycolor) transferToRecvFeeAddr(tokenName string, recvFee bn.Number) {

	if recvFee.CmpI(0) <= 0 {
		return
	}

	rfi := e._recvFeeInfos()
	account := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
	totalTransferFee := bn.N(0)

	for k, feeRatio := range rfi.RecvFeeRatio {
		if rfi.RecvFeeAddr[k] != e.sdk.Message().Contract().Address() {
			account.TransferByName(
				tokenName,
				rfi.RecvFeeAddr[k],
				recvFee.MulI(feeRatio).DivI(PERMILLI),
			)
			totalTransferFee.Add(recvFee.MulI(feeRatio).DivI(PERMILLI))
		}

	}

	return
}

func (e *Everycolor) checkSetRecvFeeInfo(_recFeeInfoStr string) (recFeeInfo *RecvFeeInfo) {
	recFeeInfo = &RecvFeeInfo{}
	jsonErr := jsoniter.Unmarshal([]byte(_recFeeInfoStr), recFeeInfo)
	sdk.RequireNotError(jsonErr, types.ErrInvalidParameter)

	e.checkRecvFeeInfo(recFeeInfo.RecvFeeRatio, recFeeInfo.RecvFeeAddr)

	return
}

func (e *Everycolor) checkRecvFeeInfo(recvFeeRatio []int64, recvFeeAddr []types.Address) {

	sdk.Require(
		len(recvFeeRatio) > 0,
		types.ErrInvalidParameter,
		"The length of RecFeeRatio must be larger than zero",
	)

	sdk.Require(
		len(recvFeeAddr) > 0,
		types.ErrInvalidParameter,
		"The length of RecFeeAddr must be larger than zero",
	)

	sdk.Require(
		len(recvFeeRatio) == len(recvFeeAddr),
		types.ErrInvalidParameter,
		"The RecFeeAddr's length must be equal with RecFeeRatio's length",
	)

	allRatio := int64(0)
	for k, ratio := range recvFeeRatio {
		sdk.Require(
			ratio > 0,
			types.ErrInvalidParameter,
			"RecRatio must be larger than zero",
		)

		allRatio += ratio
		//检查地址是否合法
		sdk.RequireAddress(recvFeeAddr[k])

		sdk.Require(
			recvFeeAddr[k] != e.sdk.Message().Contract().Account(),
			types.ErrInvalidParameter,
			"The contract account address cannot be set as the transfer fee address")

	}

	//设置的分配比例加起来必须等于PERMILLI
	sdk.Require(
		allRatio <= PERMILLI,
		types.ErrInvalidParameter,
		"The all ratio must add up to smaller than "+strconv.Itoa(PERMILLI))

	return
}

func (e *Everycolor) checkSettings(_settings string) (resultSettings *MapSetting) {

	resultSettings = &MapSetting{}
	jsonErr := jsoniter.Unmarshal([]byte(_settings), resultSettings)

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
		token := e.sdk.Helper().TokenHelper().TokenOfName(tokenName)

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

func (e *Everycolor) checkBetInfo(tokenName string, amount bn.Number, betDataSlice []BetData, roundInfo *RoundInfo) (modes []IMode) {

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
	for _, k := range betDataSlice {
		totalAmount = totalAmount.Add(k.BetAmount)

		//amount 精度应该等于PERMILLI cong
		sdk.Require(
			k.BetAmount.ModI(PERMILLI).CmpI(0) == 0,
			types.ErrInvalidParameter,
			"Amount per bet mode accuracy should be equal"+strconv.Itoa(PERMILLI),
		)

	}

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
func (e *Everycolor) checkPossibleWinAmount(tokenName string, amount, possibleWinAmount, feeNum bn.Number, roundInfo *RoundInfo) {

	// 先查询LokedInBets中有没有此项代币
	ok := e._chkLockedAmount(tokenName)
	if ok == false {
		// 没有就存到map中
		e._setLockedAmount(tokenName, bn.N(0))
	}

	// Check whether contract account has enough funds to process this bet.
	balance := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
	locked := e._lockedAmount(tokenName)
	sdk.Require(
		locked.Add(possibleWinAmount).Cmp(balance.BalanceOfName(tokenName).Add(amount)) < 0,
		types.ErrInvalidParameter,
		"Cannot afford to lose this bet",
	)

	sdk.Require(
		possibleWinAmount.Sub(feeNum).Cmp(roundInfo.Setting.Settings[tokenName].MaxProfit) <= 0,
		types.ErrInvalidParameter,
		"PossibleWinAmount should be smaller than maxProfit",
	)

	return
}

func (e *Everycolor) checkCommit(commitStr string) {

	tmp := e._chkRoundInfo(commitStr)
	// 新的commit 允许下注
	if tmp == false {
		return

	} else {
		//当前轮已开奖或已过期，不允许下注
		roundInfo := e._roundInfo(commitStr)
		sdk.Require(
			len(roundInfo.WinNumber) != 0 ||
				(e.sdk.Block().Height() <= roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks),
			types.ErrInvalidParameter,
			"Current round is lottery or timeout",
		)
	}

	return
}

func (e *Everycolor) checkRoundInfoPlaceBet(roundInfo *RoundInfo) {

	sdk.Require(
		len(roundInfo.Commit) != 0,
		types.ErrInvalidParameter,
		"The round is not exist",
	)

	//要结算的轮的下注高度要小于结算高度
	sdk.Require(
		roundInfo.FirstBlockHeight < e.sdk.Block().Height(),
		types.ErrInvalidParameter,
		"SettleBet block can not be in the same block as placeBet or before.",
	)

	if len(roundInfo.WinNumber) == 0 {
		//首次结算，轮信息不能过期
		sdk.Require(
			roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks > e.sdk.Block().Height(),
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

func (e *Everycolor) checkRoundInfoRefund(roundInfo *RoundInfo) {

	sdk.Require(len(roundInfo.Commit) != 0,
		types.ErrInvalidParameter,
		"The round is not exist",
	)
	sdk.Require(len(roundInfo.WinNumber) == 0,
		types.ErrInvalidParameter,
		"The round was run a lottery",
	)
	sdk.Require(roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks < e.sdk.Block().Height(),
		types.ErrInvalidParameter,
		"The round never timeout",
	)
	sdk.Require(roundInfo.RefundCount < roundInfo.BetCount,
		types.ErrInvalidParameter,
		"The round is refunded",
	)

	return
}

func IntToByte(val int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val))
	return buf
}
