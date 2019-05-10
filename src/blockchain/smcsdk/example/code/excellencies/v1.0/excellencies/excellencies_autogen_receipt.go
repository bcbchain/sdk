package excellencies

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

var _ receipt = (*Excellencies)(nil)

//emitSetSecretSigner This is a method of Excellencies
func (e *Excellencies) emitSetSecretSigner(newSecretSigner types.PubKey) {
	type setSecretSigner struct {
		NewSecretSigner types.PubKey `json:"newSecretSigner"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setSecretSigner{
			NewSecretSigner: newSecretSigner,
		},
	)
}

//emitSetSettings This is a method of Excellencies
func (e *Excellencies) emitSetSettings(Settings map[string]Setting, BetExpirationBlocks, PoolFeeRatio, CarveUpPoolRatio, GrandPrizeRatio int64) {
	type setSettings struct {
		Settings            map[string]Setting `json:"Settings"`
		BetExpirationBlocks int64              `json:"BetExpirationBlocks"`
		PoolFeeRatio        int64              `json:"PoolFeeRatio"`
		CarveUpPoolRatio    int64              `json:"CarveUpPoolRatio"`
		GrandPrizeRatio     int64              `json:"GrandPrizeRatio"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setSettings{
			Settings:            Settings,
			BetExpirationBlocks: BetExpirationBlocks,
			PoolFeeRatio:        PoolFeeRatio,
			CarveUpPoolRatio:    CarveUpPoolRatio,
			GrandPrizeRatio:     GrandPrizeRatio,
		},
	)
}

//emitPlaceBet This is a method of Excellencies
func (e *Excellencies) emitPlaceBet(tokenName string, amount bn.Number, betData []BetData, possibleWinAmount bn.Number, commitLastBlock, betCount int64, commit, signData []byte, refAddress types.Address) {
	type placeBet struct {
		TokenName         string        `json:"tokenName"`
		Amount            bn.Number     `json:"amount"`
		BetData           []BetData     `json:"betData"`
		PossibleWinAmount bn.Number     `json:"possibleWinAmount"`
		CommitLastBlock   int64         `json:"commitLastBlock"`
		BetCount          int64         `json:"betCount"`
		Commit            []byte        `json:"commit"`
		SignData          []byte        `json:"signData"`
		RefAddress        types.Address `json:"refAddress"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		placeBet{
			TokenName:         tokenName,
			Amount:            amount,
			BetData:           betData,
			PossibleWinAmount: possibleWinAmount,
			CommitLastBlock:   commitLastBlock,
			BetCount:          betCount,
			Commit:            commit,
			SignData:          signData,
			RefAddress:        refAddress,
		},
	)
}

//emitSetRecFeeInfo This is a method of Excellencies
func (e *Excellencies) emitSetRecFeeInfo(info []RecFeeInfo) {
	type setRecFeeInfo struct {
		Info []RecFeeInfo `json:"info"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setRecFeeInfo{
			Info: info,
		},
	)
}

//emitWithdrawFunds This is a method of Excellencies
func (e *Excellencies) emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {
	type withdrawFunds struct {
		TokenName      string        `json:"tokenName"`
		Beneficiary    types.Address `json:"beneficiary"`
		WithdrawAmount bn.Number     `json:"withdrawAmount"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		withdrawFunds{
			TokenName:      tokenName,
			Beneficiary:    beneficiary,
			WithdrawAmount: withdrawAmount,
		},
	)
}

//emitSettleBet This is a method of Excellencies
func (e *Excellencies) emitSettleBet(reveal, commit []byte, banker GamerInfo, players map[string]GamerInfo, startIndex, endIndex int64, amountOfWin, amountOfUnLock map[string]bn.Number, finished bool, poolAmount map[string]map[string]bn.Number) {
	type settleBet struct {
		Reveal         []byte                          `json:"reveal"`
		Commit         []byte                          `json:"commit"`
		Banker         GamerInfo                       `json:"banker"`
		Players        map[string]GamerInfo            `json:"players"`
		StartIndex     int64                           `json:"startIndex"`
		EndIndex       int64                           `json:"endIndex"`
		AmountOfWin    map[string]bn.Number            `json:"amountOfWin"`
		AmountOfUnLock map[string]bn.Number            `json:"amountOfUnLock"`
		Finished       bool                            `json:"finished"`
		PoolAmount     map[string]map[string]bn.Number `json:"poolAmount"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		settleBet{
			Reveal:         reveal,
			Commit:         commit,
			Banker:         banker,
			Players:        players,
			StartIndex:     startIndex,
			EndIndex:       endIndex,
			AmountOfWin:    amountOfWin,
			AmountOfUnLock: amountOfUnLock,
			Finished:       finished,
			PoolAmount:     poolAmount,
		},
	)
}

//emitWithdrawWin This is a method of Excellencies
func (e *Excellencies) emitWithdrawWin(commit []byte, amountOfWin, amountOfUnLock map[string]bn.Number, poolAmount map[string]map[string]bn.Number) {
	type withdrawWin struct {
		Commit         []byte                          `json:"commit"`
		AmountOfWin    map[string]bn.Number            `json:"amountOfWin"`
		AmountOfUnLock map[string]bn.Number            `json:"amountOfUnLock"`
		PoolAmount     map[string]map[string]bn.Number `json:"poolAmount"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		withdrawWin{
			Commit:         commit,
			AmountOfWin:    amountOfWin,
			AmountOfUnLock: amountOfUnLock,
			PoolAmount:     poolAmount,
		},
	)
}

//emitRefundBet This is a method of Excellencies
func (e *Excellencies) emitRefundBet(commit []byte, refundCount int64, unlockAmount map[string]bn.Number, finished bool) {
	type refundBet struct {
		Commit       []byte               `json:"commit"`
		RefundCount  int64                `json:"refundCount"`
		UnlockAmount map[string]bn.Number `json:"unlockAmount"`
		Finished     bool                 `json:"finished"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		refundBet{
			Commit:       commit,
			RefundCount:  refundCount,
			UnlockAmount: unlockAmount,
			Finished:     finished,
		},
	)
}
