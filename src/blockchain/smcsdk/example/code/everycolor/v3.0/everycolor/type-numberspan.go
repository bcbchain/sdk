package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"sort"
	"strconv"
	"strings"
)

//跨度模式
type NumberSpan struct {
	Base BaseMode
}

func (ns *NumberSpan) ToLockAmount() bn.Number {

	return ns.Base.PossibleWinAmount()
}

func (ns *NumberSpan) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {
	effStr := ns.Base.GetEffectiveNum(lotteryNum)
	odds := int64(0)
	var NumSlice []int
	for _, item := range effStr {
		NumSlice = append(NumSlice, int(item))
	}
	sort.Ints(NumSlice)
	//最后一个减去第一个
	numSpan, _ := strconv.Atoi(string(ns.Base.GetBetType()[0]))
	if NumSlice[2]-NumSlice[0] == numSpan {
		odds = ns.Base.OddsMap[ns.Base.GetBetType()]
	}
	winAmount, feeNum := getEachWinAndFee(ns.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (ns *NumberSpan) Init(data BetData) {
	//check
	ns.checkBetData(data.BetMsg)

	ns.Base.BetMsg = data.BetMsg
	ns.Base.BetAmount = data.BetAmount
	ns.Base.OddsMap = map[string]int64{
		"0**": 975500, "1**": 181055,
		"2**": 102062, "3**": 77880,
		"4**": 68208, "5**": 65500,
		"6**": 68208, "7**": 77880,
		"8**": 102062, "9**": 181055,
	}

	return
}

func (ns *NumberSpan) checkBetData(betMsg string) {
	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	sdk.Require(
		strings.Count(betMsg, "X") == 2 && strings.Count(string(betMsg), "**") == 1,
		types.ErrInvalidParameter,
		"betMsg format is error ",
	)

	index := strings.Index(betMsg, "**")

	if index > 0 {
		sdk.Require(
			string(betMsg[index-1]) >= "0" && string(betMsg[index-1]) <= "9",
			types.ErrInvalidParameter,
			"betMsg value is error",
		)

	} else {
		sdk.Require(
			false,
			types.ErrInvalidParameter,
			"betMsg type is error",
		)

	}

	return
}
