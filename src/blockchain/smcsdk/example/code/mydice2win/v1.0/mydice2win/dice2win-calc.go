package mydice2win

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

// Get the expected win amount after fee is subtracted.
func (dw *Dice2Win) getDiceWinAmount(amount, modulo, rollUnder bn.Number) (winAmount, fee bn.Number) {

	sdk.Require(rollUnder.CmpI(0) > 0 && rollUnder.Cmp(modulo) <= 0,
		types.ErrInvalidParameter, "Win probability out of range")

	settings := dw._settings()
	winAmount = bn.N(0)
	fee = bn.N(0)

	// 手续费
	fee = amount.MulI(settings.FeeRatio).DivI(perMille)
	if fee.CmpI(settings.FeeMinimum) < 0 {
		fee = bn.N(settings.FeeMinimum)
	}

	sdk.Require(fee.Cmp(amount) <= 0,
		types.ErrInvalidParameter, "Bet doesn't even cover fee")

	//计算最后能赢多少钱
	winAmount = amount.Mul(modulo).Div(rollUnder).Sub(fee)

	return
}
