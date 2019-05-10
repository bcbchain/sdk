package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strings"
)

//定位模式
type Location struct {
	Base BaseMode
}

func (loc *Location) ToLockAmount() bn.Number {

	odds := int64(0)
	if len(loc.Base.GetBetType()) == 2 {
		odds = loc.Base.OddsMap["NM"]
	} else {
		odds = loc.Base.OddsMap["NMP"]
	}

	return loc.Base.BetAmount.MulI(odds).DivI(PERCENT)
}

func (loc *Location) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {
	odds := int64(0)
	effStr := loc.Base.GetEffectiveNum(lotteryNum)
	if effStr == loc.Base.GetBetType() {
		if len(loc.Base.GetBetType()) == 2 {
			odds = loc.Base.OddsMap["NM"]
		} else {
			odds = loc.Base.OddsMap["NMP"]
		}
	}
	winAmount, feeNum := getEachWinAndFee(loc.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (loc *Location) Init(data BetData) {

	loc.checkBetData(data.BetMsg)

	loc.Base.BetMsg = data.BetMsg
	loc.Base.BetAmount = data.BetAmount
	loc.Base.OddsMap = map[string]int64{
		"NM": 975500, "NMP": 9750500,
	}

	return
}

func (loc *Location) checkBetData(betMsg string) {
	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	switch strings.Count(betMsg, "X") {
	case 2:
		sdk.Require(
			isLocation(betMsg),
			types.ErrInvalidParameter,
			"betMsg bet type is error ",
		)
	case 3:
		break
	default:
		sdk.Require(
			false,
			types.ErrInvalidParameter,
			"betStr format is error ",
		)
	}

	effBet := strings.Replace(betMsg, "X", "", -1)

	for _, item := range effBet {
		sdk.Require(
			string(item) >= "0" && string(item) <= "9",
			types.ErrInvalidParameter,
			"betMsg type is error ",
		)
	}

	return
}
