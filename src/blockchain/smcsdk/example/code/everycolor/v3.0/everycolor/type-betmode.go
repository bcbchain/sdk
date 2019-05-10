package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"strings"
)

const (
	PERCENT = 10000 // 赔率的万分比
)

type BetData struct {
	BetMode   string    `json:"betMode"`   // bet mode
	BetMsg    string    `json:"betMsg"`    // bet small type
	BetAmount bn.Number `json:"betAmount"` // each bet amount
}

type IMode interface {
	ToLockAmount() bn.Number
	WinAmount(string, int64) (bn.Number, bn.Number)
	Init(BetData)
}

type BaseMode struct {
	BetMsg    string
	BetAmount bn.Number
	OddsMap   map[string]int64
}

func (base *BaseMode) PossibleWinAmount() (possibleWin bn.Number) {
	// 通用的算可能赢得钱，
	possibleWin = bn.N(0)
	betType := base.GetBetType()
	possibleWin = possibleWin.Add(base.BetAmount.MulI(base.OddsMap[betType]).DivI(PERCENT))

	return
}

func (base *BaseMode) GetBetType() string {
	return strings.Replace(base.BetMsg, "X", "", -1)
}

func (base *BaseMode) GetEffectiveNum(lottery string) string {

	effectiveNum := ""
	for i, item := range base.BetMsg {
		if item != 'X' {
			effectiveNum += string(lottery[i])
		}
	}

	return effectiveNum
}

func CreateBetMode(betData []BetData) (modes []IMode) {

	modes = make([]IMode, 0, 0)

	for _, item := range betData {
		switch item.BetMode {
		case "A":
			ball := new(Ball)
			ball.Init(item)
			modes = append(modes, ball)
		case "B":
			gf := new(GoldenFlowers)
			gf.Init(item)
			modes = append(modes, gf)
		case "C":
			com := new(Combination)
			com.Init(item)
			modes = append(modes, com)
		case "D":
			ori := new(Location)
			ori.Init(item)
			modes = append(modes, ori)
		case "E":
			gs := new(GroupSelection)
			gs.Init(item)
			modes = append(modes, gs)
		case "F":
			ns := new(NumberSpan)
			ns.Init(item)
			modes = append(modes, ns)
		case "G":
			ts := new(TwoSide)
			ts.Init(item)
			modes = append(modes, ts)
		case "H":
			dt := new(DragonTiger)
			dt.Init(item)
			modes = append(modes, dt)
		case "I":
			tp := new(TotalPoint)
			tp.Init(item)
			modes = append(modes, tp)
		default:
			sdk.Require(
				false,
				types.ErrInvalidParameter,
				"The bet mode is error",
			)
		}

	}

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

//计算手续费
func calcFeeAmount(betAmount bn.Number, feeRatio int64) bn.Number {

	feeAmount := betAmount.MulI(feeRatio).DivI(PERMILLI)
	return feeAmount
}
