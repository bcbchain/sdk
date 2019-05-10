package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strconv"
	"strings"
)

//单球模式
type Ball struct {
	Base BaseMode
}

func (ball *Ball) ToLockAmount() bn.Number {
	betType := ball.Base.GetBetType()
	if betType >= "0" && betType <= "9" {
		betType = "N"
	}

	// 扣除手续费之前可能赢得钱
	return ball.Base.BetAmount.MulI(ball.Base.OddsMap[betType]).DivI(PERCENT)
}

func (ball *Ball) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {

	winNum, _ := strconv.Atoi(ball.Base.GetEffectiveNum(lotteryNum))
	betType := ball.Base.GetBetType()
	var odds int64 = 0
	switch betType {
	case "O":
		if winNum%2 == 1 {
			odds = ball.Base.OddsMap["O"]
		}
	case "E":
		if winNum%2 != 1 {
			odds = ball.Base.OddsMap["E"]
		}
	case "S":
		if winNum < 5 {
			odds = ball.Base.OddsMap["S"]
		}
	case "B":
		if winNum >= 5 {
			odds = ball.Base.OddsMap["B"]
		}
	default: //处理押注数字的情况
		if betType == strconv.Itoa(winNum) {
			odds = ball.Base.OddsMap["N"]
		}
	}

	winAmount, feeNum := getEachWinAndFee(ball.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (ball *Ball) Init(data BetData) {

	ball.checkBetData(data.BetMsg)

	ball.Base.BetMsg = data.BetMsg
	ball.Base.BetAmount = data.BetAmount
	ball.Base.OddsMap = map[string]int64{
		"O": 20000, "E": 20000,
		"S": 20000, "B": 20000,
		"N": 98000,
	}

	return
}

func (ball *Ball) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)
	sdk.Require(
		strings.Count(betMsg, "X") == 4,
		types.ErrInvalidParameter,
		"betMsg format is error",
	)

	effBet := strings.Replace(betMsg, "X", "", -1)
	sdk.Require(
		effBet == "O" || effBet == "E" || effBet == "S" || effBet == "B" || (effBet >= "0" && effBet <= "9"),
		types.ErrInvalidParameter,
		"betMsg type is error",
	)

	return
}
