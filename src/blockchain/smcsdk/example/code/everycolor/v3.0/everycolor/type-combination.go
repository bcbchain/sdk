package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strings"
)

//组合模式
type Combination struct {
	Base BaseMode
}

func (com *Combination) ToLockAmount() bn.Number {
	odds := int64(0)
	counts := strings.Count(com.Base.BetMsg, "*")
	effStr := com.Base.GetBetType()
	switch counts {
	case 0:
		//三字组合
		if strings.Count(effStr, string(effStr[0])) == 3 {
			// 三数相同
			odds = com.Base.OddsMap["NNN"]
		} else if strings.Count(effStr, string(effStr[0])) == 2 ||
			strings.Count(effStr, string(effStr[1])) == 2 {
			odds = com.Base.OddsMap["NNP"]
			//二数相同
		} else {
			//三数不同
			odds = com.Base.OddsMap["NMP"]
		}
	case 1:
		//二字组合
		if effStr[0] == effStr[1] {
			//二数相同
			odds = com.Base.OddsMap["NN*"]
		} else {
			//二数不同
			odds = com.Base.OddsMap["NM*"]
		}
	case 2:
		//一字组合
		odds = com.Base.OddsMap["N**"]
	case 4:
		odds = com.Base.OddsMap["N****"]
	}

	return com.Base.BetAmount.MulI(odds).DivI(PERCENT)
}

func (com *Combination) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {
	odds := int64(0)
	effStr := com.Base.GetEffectiveNum(lotteryNum)

	counts := strings.Count(com.Base.BetMsg, "*")
	betStr := com.Base.GetBetType()
	switch counts {
	case 0:
		//三字组合
		if strings.Count(betStr, string(betStr[0])) == 3 {
			// 三数相同
			if betStr == effStr {
				odds = com.Base.OddsMap["NNN"]
			}
		} else if strings.Count(betStr, string(betStr[0])) == 2 ||
			strings.Count(betStr, string(betStr[1])) == 2 {
			//二数相同
			if isWin(effStr, betStr) {
				odds = com.Base.OddsMap["NNP"]
			}

		} else {
			//三数不同
			if isWin(effStr, betStr) {
				odds = com.Base.OddsMap["NMP"]
			}
		}
	case 1:
		//二字组合
		if betStr[0] == betStr[1] {
			//二数相同
			if strings.Count(effStr, string(betStr[0])) == 2 {
				odds = com.Base.OddsMap["NN*"]
			}

		} else {
			//二数不同
			if strings.Contains(effStr, string(betStr[0])) &&
				strings.Contains(effStr, string(betStr[1])) {
				odds = com.Base.OddsMap["NM*"]
			}

		}
	case 2:
		//一字组合
		if strings.Contains(effStr, string(betStr[0])) {
			odds = com.Base.OddsMap["N**"]
		}
	case 4:
		//全五
		if strings.Contains(effStr, string(com.Base.BetMsg[0])) {
			odds = com.Base.OddsMap["N****"]
		}

	}

	winAmount, feeNum := getEachWinAndFee(com.Base.BetAmount, odds, feeRatio)

	return winAmount, feeNum
}

func (com *Combination) Init(data BetData) {

	com.checkBetData(data.BetMsg)

	com.Base.BetMsg = data.BetMsg
	com.Base.BetAmount = data.BetAmount
	com.Base.OddsMap = map[string]int64{
		"N****": 24308, "N**": 36477,
		"NM*": 181055, "NN*": 348714,
		"NMP": 1625500, "NNP": 3250500,
		"NNN": 9750500,
	}

	return
}

func (com *Combination) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	sdk.Require(
		strings.Count(betMsg, "X") == 2 || strings.Count(betMsg, "*") == 4,
		types.ErrInvalidParameter,
		"betStr format is error",
	)

	if strings.Count(betMsg, "*") == 4 {
		effBet := strings.Replace(betMsg, "*", "", -1)
		sdk.Require(
			effBet >= "0" && effBet <= "9",
			types.ErrInvalidParameter,
			"betMsg type is error",
		)
	} else {

		sdk.Require(
			isLocation(betMsg),
			types.ErrInvalidParameter,
			"betMsg type is error ",
		)

		effBet := strings.Replace(betMsg, "X", "", -1)
		ok := true
		switch strings.Count(effBet, "*") {
		case 0:
			if (effBet[0] > '9' || effBet[0] < '0') ||
				(effBet[1] > '9' || effBet[1] < '0') ||
				(effBet[2] > '9' || effBet[2] < '0') {
				ok = false
			}
		case 1:
			if (effBet[0] > '9' || effBet[0] < '0') ||
				(effBet[1] > '9' || effBet[1] < '0') {
				ok = false
			}
		case 2:
			if effBet[0] > '9' || effBet[0] < '0' {
				ok = false
			}
		}
		if !ok {
			sdk.Require(
				ok,
				types.ErrInvalidParameter,
				"betMsg type is error ",
			)
			return
		}

	}

	return
}

func isWin(effStr, betStr string) bool {
	sortEff := sortStr(effStr)
	sortBet := sortStr(betStr)
	//二数相同
	if sortBet[0] == sortEff[0] && sortBet[1] == sortEff[1] && sortBet[2] == sortEff[2] {
		return true
	}
	return false
}
