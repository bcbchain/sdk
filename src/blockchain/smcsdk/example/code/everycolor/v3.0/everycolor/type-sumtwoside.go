package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strconv"
	"strings"
)

//和数双面模式
type TwoSide struct {
	Base BaseMode
}

func (ts *TwoSide) ToLockAmount() bn.Number {

	return ts.Base.PossibleWinAmount()
}

func (ts *TwoSide) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {

	var winNum int
	for _, item := range ts.Base.GetEffectiveNum(lotteryNum) {
		num, _ := strconv.Atoi(string(item))
		winNum += num
	}
	var odds int64 = 0
	switch ts.Base.GetBetType() {
	case "OO":
		if winNum%2 == 1 {
			odds = ts.Base.OddsMap["OO"]
		}
	case "EE":
		if winNum%2 == 0 {
			odds = ts.Base.OddsMap["EE"]
		}

	}
	winAmount, feeNum := getEachWinAndFee(ts.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (ts *TwoSide) Init(data BetData) {

	ts.checkBetData(data.BetMsg)

	ts.Base.BetMsg = data.BetMsg
	ts.Base.BetAmount = data.BetAmount
	ts.Base.OddsMap = map[string]int64{
		"OO": 20000, "EE": 20000,
	}

	return
}

func (ts *TwoSide) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	sdk.Require(
		strings.Count(betMsg, "X") == 3,
		types.ErrInvalidParameter,
		"betMsg format is error",
	)

	sdk.Require(
		strings.Count(betMsg, "O") == 2 || strings.Count(betMsg, "E") == 2,
		types.ErrInvalidParameter,
		"betMsg type is error",
	)

	return
}
