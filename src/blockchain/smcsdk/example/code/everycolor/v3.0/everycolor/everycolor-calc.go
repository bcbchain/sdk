package everycolor

import (
	"blockchain/smcsdk/sdk/bn"
	"fmt"
	"strconv"
)

func getLotteryNum(random []byte) (winNumber string) {

	entropy := bn.NBS(random)
	winNumber = fmt.Sprintf("%05d", entropy.ModI(FIVEBALLMODEL).V.Uint64())

	return
}
func (e *Everycolor) SettleBet(bet *BetInfo, index int64, commitStr string, roundInfo *RoundInfo) (possibleWinAmount, ecWinAmount, feeNum bn.Number) {

	possibleWinAmount = bn.N(0)
	ecWinAmount = bn.N(0)
	feeNum = bn.N(0)

	modes := CreateBetMode(bet.BetData)

	//Calculate the amount of money locked  for each bet
	for _, item := range modes {

		setting := roundInfo.Setting.Settings
		possibleWin := item.ToLockAmount()
		possibleWinAmount = possibleWinAmount.Add(possibleWin)
		eachWin, eachFee := item.WinAmount(roundInfo.WinNumber, setting[bet.TokenName].FeeRatio)

		bet.EachWinAmount = append(bet.EachWinAmount, eachWin)
		// 计算总的费用
		ecWinAmount = ecWinAmount.Add(eachWin)
		feeNum = feeNum.Add(eachFee)

	}
	//Calculate the amount of Win money  for each bet
	if ecWinAmount.CmpI(0) > 0 {
		account := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
		//Send the funds to gambler.
		account.TransferByName(bet.TokenName, bet.Gambler, ecWinAmount)

	}

	//save BetInfo
	bet.WinAmount = ecWinAmount
	bet.Settled = true
	e._setBetInfo(commitStr, strconv.Itoa(int(index)), bet)

	return
}
func (e *Everycolor) RatioTransfer(tokenName string, feeNum bn.Number, setting Setting) {

	account := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())

	sentToCltFee := feeNum.MulI(setting.SendToCltRatio).DivI(PERMILLI)
	//转给clt
	if sentToCltFee.CmpI(0) > 0 {
		account.TransferByName(
			tokenName,
			e.sdk.Helper().BlockChainHelper().CalcAccountFromName("clt", ""),
			sentToCltFee,
		)

	}

	//转给除clt外的其他接收地址
	e.transferToRecvFeeAddr(tokenName, feeNum.Sub(sentToCltFee))

}

func (e *Everycolor) Refund(bet *BetInfo, index int64, commitStr string) (totalPossibleWinAmount bn.Number) {
	account := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())

	account.TransferByName(bet.TokenName, bet.Gambler, bet.Amount)

	bet.Settled = true
	e._setBetInfo(commitStr, strconv.Itoa(int(index)), bet)
	modes := CreateBetMode(bet.BetData)

	totalPossibleWinAmount = bn.N(0)

	for _, item := range modes {
		possibleWinAmount := item.ToLockAmount()

		totalPossibleWinAmount = totalPossibleWinAmount.Add(possibleWinAmount)
	}

	return
}
