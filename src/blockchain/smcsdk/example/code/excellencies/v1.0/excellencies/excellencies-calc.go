package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/types"
	"strconv"
)

func (e *Excellencies) MaybeWinAmountAndFeeByList(amount bn.Number, betInfo BetInfo) (totalMaybeWinAmount, fee bn.Number) {
	fee = bn.N(0)
	totalMaybeWinAmount = bn.N(0)
	settings := e._mapSetting()
	setting := settings.Settings[betInfo.TokenName]
	// fee
	fee = amount.MulI(setting.FeeRatio).DivI(PERMILLE)
	if fee.Cmp(setting.FeeMiniNum) < 0 {
		fee = setting.FeeMiniNum
	}

	sdk.Require(fee.Cmp(amount) <= 0,
		types.ErrInvalidParameter, "Bet doesn't even cover fee")

	forx.Range(betInfo.BetData, func(_, bet BetData) bool {

		totalMaybeWinAmount = totalMaybeWinAmount.Add(bet.BetAmount.MulI(ODDS))
		return forx.Continue
	})

	return
}

func (e *Excellencies) Lottery(random []byte, roundInfo *RoundInfo) {
	banker := GamerInfo{}
	players := make([]GamerInfo, 4)

	card := Dealer{}
	card.Init(random)

	// get cards
	forx.Range(0, 3, func(index int) bool {
		forx.Range(0, len(players), func(i int) bool {
			players[i].AddCard(card.GetCard())

			return forx.Continue
		})

		banker.AddCard(card.GetCard())
		return forx.Continue
	})

	forx.Range(0, len(players), func(i int) bool {
		banker.JudgeWin(&players[i])

		return forx.Continue
	})

	mapPlayers := map[string]GamerInfo{"A": players[0], "B": players[1], "C": players[2], "D": players[3]}
	roundInfo.Banker = banker
	roundInfo.Players = mapPlayers

	return
}

func (sg *Excellencies) settleBet(bet *BetInfo, index int64, commitStr string, roundInfo *RoundInfo) (possibleWinAmount, sgWinAmount, feeNum, addPool bn.Number) {

	possibleWinAmount = bn.N(0)
	sgWinAmount = bn.N(0)
	feeNum = bn.N(0)
	sgPoolWin := bn.N(0)
	modes := CreateBetMode(bet.BetData)

	//Calculate the amount of money locked  for each bet
	setting := roundInfo.Setting.Settings
	forx.Range(modes, func(_, item BaseMode) bool {
		possibleWin := item.ToLockAmount()
		possibleWinAmount = possibleWinAmount.Add(possibleWin)
		eachWin, eachFee, poolWin := item.WinAmount(bet.Saler, roundInfo, setting[bet.TokenName].FeeRatio)
		sgPoolWin = sgPoolWin.Add(poolWin)

		bet.EachWinAmount = append(bet.EachWinAmount, eachWin)
		// 计算总的费用
		sgWinAmount = sgWinAmount.Add(eachWin)

		feeNum = feeNum.Add(eachFee)
		return forx.Continue
	})

	//Calculate the amount of Win money  for each bet
	//too pool
	addPool = sgWinAmount.MulI(roundInfo.Setting.PoolFeeRatio).DivI(PERMILLI)
	sgWinAmount = sgWinAmount.Sub(addPool)

	//save BetInfo
	bet.WinAmount = sgWinAmount
	bet.Settled = true
	sg._setBetInfo(commitStr, strconv.Itoa(int(index)), bet)

	sg.SetGrandPrizer(bet) //检查更新奖励最大者

	return
}

func (sg *Excellencies) Refund(bet *BetInfo, index int64, commitStr string) (totalPossibleWinAmount bn.Number) {

	bet.Settled = true
	sg._setBetInfo(commitStr, strconv.Itoa(int(index)), bet)
	modes := CreateBetMode(bet.BetData)

	totalPossibleWinAmount = bn.N(0)

	forx.Range(modes, func(_, item BaseMode) bool {
		possibleWinAmount := item.ToLockAmount()
		totalPossibleWinAmount = totalPossibleWinAmount.Add(possibleWinAmount)
		return forx.Continue
	})

	return
}

func (sg *Excellencies) PayGrandPrize(address types.Address) {
	mapsettings := sg._mapSetting()
	nowTime := sg.sdk.Block().Time()
	forx.Range(mapsettings.Settings, func(tokenName string, _ Setting) bool {
		if sg._chkGrandPrizer(address, tokenName) {
			prizer := sg._grandPrizer(address, tokenName)
			if nowTime-prizer.BetTime > 24*3600 {
				//发奖
				account := sg.sdk.Helper().AccountHelper().AccountOf(sg.sdk.Message().Contract().Account())
				poolAmount := sg._poolAmount(address, tokenName)
				rewardAmount := poolAmount.MulI(mapsettings.GrandPrizeRatio).DivI(PERMILLI)
				account.TransferByName(tokenName, prizer.Gambler, rewardAmount)

				sg._setPoolAmount(address, tokenName, poolAmount.Sub(rewardAmount))

				prizer := GrandPrizer{0, "", bn.N(0)}
				sg._setGrandPrizer(address, tokenName, prizer)
			}
		}
		return forx.Continue

	})

}

func (sg *Excellencies) SetGrandPrizer(bet *BetInfo) {

	if sg._chkGrandPrizer(bet.Saler, bet.TokenName) {
		prizer := sg._grandPrizer(bet.Saler, bet.TokenName)
		if prizer.WinAmount.Cmp(bet.WinAmount) >= 0 {
			return
		}
	}

	prizer := GrandPrizer{bet.BetTime, bet.Gambler, bet.WinAmount}
	sg._setGrandPrizer(bet.Saler, bet.TokenName, prizer)
}
