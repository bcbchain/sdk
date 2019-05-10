package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"sort"
	"strings"
)

//金花模式
type GoldenFlowers struct {
	Base BaseMode
}

func (gf *GoldenFlowers) ToLockAmount() bn.Number {

	return gf.Base.PossibleWinAmount()
}

func (gf *GoldenFlowers) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {
	effStr := gf.Base.GetEffectiveNum(lotteryNum)
	odds := int64(0)
	betType := gf.Base.GetBetType()

	if gf.checkEffStr(effStr) == betType {
		odds = gf.Base.OddsMap[betType]
	}
	winAmount, feeNum := getEachWinAndFee(gf.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (gf *GoldenFlowers) Init(data BetData) {

	gf.checkBetData(data.BetMsg)

	gf.Base.BetMsg = data.BetMsg
	gf.Base.BetAmount = data.BetAmount
	gf.Base.OddsMap = map[string]int64{
		"AAA": 975500, "ABC": 163000,
		"AA*": 36611, "ABD": 27583,
		"ACE": 33000,
	}

	return
}

func (gf *GoldenFlowers) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	sdk.Require(
		strings.Count(betMsg, "X") == 2,
		types.ErrInvalidParameter,
		"betStr format is error ",
	)

	if strings.Contains(betMsg, "AAA") ||
		strings.Contains(betMsg, "ABC") ||
		strings.Contains(betMsg, "AA*") ||
		strings.Contains(betMsg, "ABD") ||
		strings.Contains(betMsg, "ACE") {
		return
	} else {
		sdk.Require(
			false,
			types.ErrInvalidParameter,
			"betStr type is error ",
		)
	}
	return
}

func sortStr(effStr string) []int {
	var numSlice []int
	numSlice = append(numSlice, int(effStr[0]))
	numSlice = append(numSlice, int(effStr[1]))
	numSlice = append(numSlice, int(effStr[2]))
	sort.Ints(numSlice)
	return numSlice
}

func isABC(numSlice []int) bool {

	if (numSlice[0]+1 == numSlice[1] && numSlice[1]+1 == numSlice[2]) ||
		(numSlice[0] == '0' && numSlice[2] == '9' && (numSlice[1] == '8' || numSlice[1] == '1')) {
		return true
	}
	return false
}

func isAA(effStr string) bool {

	if strings.Count(effStr, string(effStr[0])) == 2 ||
		strings.Count(effStr, string(effStr[1])) == 2 {
		return true
	}
	return false
}

func isABD(numSlice []int) bool {
	// 0,*,9
	if numSlice[0] == '0' &&
		numSlice[2] == '9' &&
		numSlice[1] != '8' &&
		numSlice[1] != '1' &&
		numSlice[1] != '0' &&
		numSlice[1] != '9' {
		return true
	}
	if (numSlice[1]-numSlice[0] == 1 && numSlice[2]-numSlice[1] >= 2) ||
		(numSlice[2]-numSlice[1] == 1 && numSlice[1]-numSlice[0] >= 2) {
		return true
	}
	return false
}

func (gf *GoldenFlowers) checkEffStr(effStr string) string {

	numSlice := sortStr(effStr)

	if strings.Count(effStr, string(effStr[0])) == 3 {
		return "AAA"
	} else if isABC(numSlice) {
		return "ABC"

	} else if isAA(effStr) {
		return "AA*"

	} else if isABD(numSlice) {
		return "ABD"
	} else {
		return "ACE"
	}

}
