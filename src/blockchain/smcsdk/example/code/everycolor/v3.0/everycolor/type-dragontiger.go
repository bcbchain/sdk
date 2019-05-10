package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strings"
)

//龙虎和模式
type DragonTiger struct {
	Base BaseMode
}

func (dt *DragonTiger) ToLockAmount() bn.Number {

	return dt.Base.PossibleWinAmount()
}

func (dt *DragonTiger) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {
	effStr := dt.Base.GetEffectiveNum(lotteryNum)
	odds := int64(0)
	winAmount := bn.N(0)
	feeNum := bn.N(0)

	betType := dt.Base.GetBetType()
	if dt.checkEffStr(effStr) == betType {
		odds = dt.Base.OddsMap[betType]
	}

	winAmount, feeNum = getEachWinAndFee(dt.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (dt *DragonTiger) Init(data BetData) {

	dt.checkBetData(data.BetMsg)

	dt.Base.BetMsg = data.BetMsg
	dt.Base.BetAmount = data.BetAmount
	dt.Base.OddsMap = map[string]int64{
		"D*": 22166, "*T": 22166,
		"==": 98000,
	}

	return
}

func (dt *DragonTiger) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	if strings.Contains(betMsg, "DXXX*") ||
		strings.Contains(betMsg, "*XXXT") ||
		strings.Contains(betMsg, "=XXX=") {
		return
	} else {
		sdk.Require(
			false,
			types.ErrInvalidParameter,
			"betMsg type is error",
		)
	}
	return
}

func (dt *DragonTiger) checkEffStr(effStr string) string {
	if effStr[0] > effStr[1] {
		return "D*"
	} else if effStr[0] < effStr[1] {
		return "*T"
	} else {
		return "=="
	}

}
