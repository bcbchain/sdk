package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strconv"
	"strings"
)

//总和大小单双模式
type TotalPoint struct {
	Base BaseMode
}

func (tp *TotalPoint) ToLockAmount() bn.Number {
	var odds int64 = 0
	effBet := strings.Replace(tp.Base.BetMsg, "*", "", -1)
	odds = tp.Base.OddsMap[effBet]
	return tp.Base.BetAmount.MulI(odds).DivI(PERCENT)
}

func (tp *TotalPoint) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {

	var winNum int
	for _, item := range tp.Base.GetEffectiveNum(lotteryNum) {
		num, _ := strconv.Atoi(string(item))
		winNum += num
	}
	var odds int64 = 0
	switch string(tp.Base.BetMsg[0]) {
	case "O":
		if winNum%2 == 1 {
			odds = tp.Base.OddsMap["O"]
		}
	case "E":
		if winNum%2 == 0 {
			odds = tp.Base.OddsMap["E"]
		}
	case "S":
		if winNum < 23 {
			odds = tp.Base.OddsMap["S"]
		}
	case "B":
		if winNum >= 23 {
			odds = tp.Base.OddsMap["B"]
		}
	}

	winAmount, feeNum := getEachWinAndFee(tp.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (tp *TotalPoint) Init(data BetData) {

	tp.checkBetData(data.BetMsg)

	tp.Base.BetMsg = data.BetMsg
	tp.Base.BetAmount = data.BetAmount
	tp.Base.OddsMap = map[string]int64{
		"O": 20000, "E": 20000,
		"S": 20000, "B": 20000,
	}

	return
}

func (tp *TotalPoint) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	sdk.Require(
		strings.Count(betMsg, "****") == 1,
		types.ErrInvalidParameter,
		"betMsg format is error",
	)

	tmp := betMsg[0]
	sdk.Require(
		tmp == 'O' || tmp == 'E' ||
			tmp == 'S' || tmp == 'B',
		types.ErrInvalidParameter,
		"betMsg is error",
	)

	return
}
