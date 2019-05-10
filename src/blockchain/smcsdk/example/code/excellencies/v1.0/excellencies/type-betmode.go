package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/types"
)

const (
	PERCENT = 10000 // 赔率的万分比
)

type BaseMode struct {
	BetMsg    string
	TokenName string
	BetAmount bn.Number
	Odds      int64
}

func (base *BaseMode) Init(data BetData) {
	base.BetMsg = data.PlayerId
	base.BetAmount = data.BetAmount
	base.Odds = 20000

	return
}

func (base *BaseMode) ToLockAmount() (possibleWin bn.Number) {
	// 通用的算可能赢得钱，
	possibleWin = bn.N(0)
	possibleWin = possibleWin.Add(base.BetAmount.MulI(base.Odds).DivI(PERCENT))

	return
}

func (base *BaseMode) SetRoundInfoTotalBuy(saler types.Address, tokenName string, roundInfo *RoundInfo) {
	// 通用的算可能赢得钱，
	_, ok := roundInfo.TotalBuy[saler]
	if !ok {
		roundInfo.TotalBuy[saler] = make(map[string]map[string]bn.Number, 0)
	}

	temp, ok := roundInfo.TotalBuy[saler][base.BetMsg]
	if !ok {
		temp = make(map[string]bn.Number, 0)
	}
	buyAmount, ok := temp[tokenName]
	if !ok {
		buyAmount = bn.N(0)
	}
	roundInfo.TotalBuy[tokenName][base.BetMsg][tokenName] = buyAmount.Add(base.BetAmount)

	return
}

func (base *BaseMode) GetWinAmount(roundInfo *RoundInfo) {
	// 通用的算可能赢得钱，
	//possibleWin = bn.N(0)
	//possibleWin = possibleWin.Add(base.BetAmount.MulI(base.Odds).DivI(PERCENT))

	return
}

func (base *BaseMode) WinAmount(salerAddress types.Address, roundInfo *RoundInfo, feeRatio int64) (eachWin, eachFee, poolWin bn.Number) {

	poolWin = bn.N(0)
	eachFee = base.BetAmount.MulI(feeRatio).DivI(PERMILLI)
	if roundInfo.Players[base.BetMsg].IsWin {
		eachWin = base.BetAmount.MulI(base.Odds).DivI(PERCENT).Sub(eachFee)
		if roundInfo.Players[base.BetMsg].Points > 9 && roundInfo.SettledPoolAmount[salerAddress][base.TokenName].CmpI(0) > 0 {
			poolWin := roundInfo.SettledPoolAmount[salerAddress][base.TokenName].Div(roundInfo.TotalBuy[salerAddress][base.BetMsg][base.TokenName]).Mul(base.BetAmount)
			eachWin = eachWin.Add(poolWin)
		}

	} else {
		eachWin = bn.N(0)
	}
	return
}

func CreateBetMode(betData []BetData) (modes []BaseMode) {

	modes = make([]BaseMode, 0, 0)

	forx.Range(betData, func(_, item BetData) bool {
		switch item.PlayerId {
		case "A", "B", "C", "D":
			var bs BaseMode
			bs.Init(item)
			modes = append(modes, bs)
		default:
			sdk.Require(
				false,
				types.ErrInvalidParameter,
				"The bet mode is error",
			)
		}
		return forx.Continue
	})

	return
}

// 得到每一笔的手续费和赢得的钱
func getEachWinAndFee(betAmount bn.Number, odds, feeRatio int64) (ecWinAmount bn.Number, feeNum bn.Number) {
	//手续费的处理
	feeNum = bn.N(0)
	ecWinAmount = betAmount.MulI(odds).DivI(PERCENT)

	feeNum = betAmount.MulI(feeRatio).DivI(PERMILLI)

	if ecWinAmount.CmpI(0) > 0 {
		ecWinAmount = ecWinAmount.Sub(feeNum)
	}

	return
}
