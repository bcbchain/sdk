package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/types"
	"sort"
	"strings"
)

//组选模式
type GroupSelection struct {
	Base BaseMode
}

func (gs *GroupSelection) ToLockAmount() bn.Number {
	betType := gs.Base.GetBetType()
	var odds int64
	if strings.Count(betType, "*") == 1 {
		odds = gs.Base.OddsMap["NM*"]
	} else {
		odds = gs.Base.OddsMap["NMP"]
	}

	return gs.Base.BetAmount.MulI(odds).DivI(PERCENT)
}

func (gs *GroupSelection) WinAmount(lotteryNum string, feeRatio int64) (bn.Number, bn.Number) {

	var winNumSlice, betNumSlice []string
	winNum := gs.Base.GetEffectiveNum(lotteryNum)
	for _, item := range winNum {
		winNumSlice = append(winNumSlice, string(item))
	}

	betType := gs.Base.GetBetType()
	for _, item := range betType {
		betNumSlice = append(betNumSlice, string(item))
	}
	var odds int64 = 0
	switch strings.Count(betType, "*") {
	case 1:
		betNum := strings.Replace(betType, "*", "", -1)
		if winNumSlice[0] != winNumSlice[2] {
			if (winNumSlice[0] == winNumSlice[1] && (winNumSlice[0] == string(betNum[0]) || winNumSlice[0] == string(betNum[1]))) ||
				(winNumSlice[1] == winNumSlice[2] && (winNumSlice[1] == string(betNum[0]) || winNumSlice[1] == string(betNum[1]))) {
				if strings.Contains(winNum, string(betNum[0])) && strings.Contains(winNum, string(betNum[1])) {
					odds = gs.Base.OddsMap["NM*"]
				}

			}
		}

	case 0:
		sort.Strings(winNumSlice)
		sort.Strings(betNumSlice)

		if betNumSlice[0] == winNumSlice[0] && betNumSlice[1] == winNumSlice[1] && betNumSlice[2] == winNumSlice[2] {
			odds = gs.Base.OddsMap["NMP"]
		}
	}
	winAmount, feeNum := getEachWinAndFee(gs.Base.BetAmount, odds, feeRatio)
	return winAmount, feeNum
}

func (gs *GroupSelection) Init(data BetData) {

	gs.checkBetData(data.BetMsg)

	gs.Base.BetMsg = data.BetMsg
	gs.Base.BetAmount = data.BetAmount
	gs.Base.OddsMap = map[string]int64{
		"NM*": 1625500, "NMP": 1625500,
	}

	return
}

func (gs *GroupSelection) checkBetData(betMsg string) {

	sdk.Require(
		len(betMsg) == 5,
		types.ErrInvalidParameter,
		"betMsg length is error",
	)

	sdk.Require(
		strings.Count(betMsg, "X") == 2,
		types.ErrInvalidParameter,
		"betMsg format \"X\" is error",
	)

	sdk.Require(
		strings.Count(betMsg, "*") == 1 || strings.Count(string(betMsg), "*") == 0,
		types.ErrInvalidParameter,
		"betMsg format \"*\" is error",
	)

	sdk.Require(
		isLocation(betMsg),
		types.ErrInvalidParameter,
		"betMsg type is error ",
	)

	betEff := strings.Replace(betMsg, "*", "", -1)
	betEff = strings.Replace(betEff, "X", "", -1)
	sdk.Require(
		len(removeRepeatStr(betEff)) == len(betEff),
		types.ErrInvalidParameter,
		"betMsg type is error ",
	)

	for _, item := range betMsg {
		if item != 'X' && item != '*' {
			sdk.Require(
				item >= '0' && item <= '9',
				types.ErrInvalidParameter,
				"betMsg type is not in range",
			)
		}
	}

	return
}

func removeRepeatStr(oldStr string) (newStr string) {

	strSlice := strings.Split(oldStr, "")

	tempMap := map[string]int64{}
	// 存放不重复主键
	for _, e := range strSlice {
		tempMap[e] = 0
	}
	//拼接成新的字符串
	newStr = ""
	forx.Range(tempMap, func(str string, value int64) bool {
		newStr = newStr + str

		return true
	})

	return newStr
}

func isLocation(betMsg string) bool {
	first := strings.Index(betMsg, "X")
	last := strings.LastIndex(betMsg, "X")

	if (first == 0 && last == 4) ||
		(first == 0 && last == 1) ||
		(first == 3 && last == 4) {
		return true
	}
	return false
}
