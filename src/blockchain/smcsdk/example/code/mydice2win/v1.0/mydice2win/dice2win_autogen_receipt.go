package mydice2win

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

var _ receipt = (*Dice2Win)(nil)

//emitSetSecretSigner This is a method of Dice2Win
func (dw *Dice2Win) emitSetSecretSigner(newSecretSigner types.PubKey) {
	type setSecretSigner struct {
		NewSecretSigner types.PubKey `json:"newSecretSigner"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		setSecretSigner{
			NewSecretSigner: newSecretSigner,
		},
	)
}

//emitSetSettings This is a method of Dice2Win
func (dw *Dice2Win) emitSetSettings(tokenNames []string, minBet, maxBet, maxProfit, feeRatio, feeMinimum, sendToCltRatio, betExpirationBlocks int64) {
	type setSettings struct {
		TokenNames          []string `json:"tokenNames"`
		MinBet              int64    `json:"minBet"`
		MaxBet              int64    `json:"maxBet"`
		MaxProfit           int64    `json:"maxProfit"`
		FeeRatio            int64    `json:"feeRatio"`
		FeeMinimum          int64    `json:"feeMinimum"`
		SendToCltRatio      int64    `json:"sendToCltRatio"`
		BetExpirationBlocks int64    `json:"betExpirationBlocks"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		setSettings{
			TokenNames:          tokenNames,
			MinBet:              minBet,
			MaxBet:              maxBet,
			MaxProfit:           maxProfit,
			FeeRatio:            feeRatio,
			FeeMinimum:          feeMinimum,
			SendToCltRatio:      sendToCltRatio,
			BetExpirationBlocks: betExpirationBlocks,
		},
	)
}

//emitSetRecvFeeInfos This is a method of Dice2Win
func (dw *Dice2Win) emitSetRecvFeeInfos(infos []RecvFeeInfo) {
	type setRecvFeeInfos struct {
		Infos []RecvFeeInfo `json:"infos"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		setRecvFeeInfos{
			Infos: infos,
		},
	)
}

//emitWithdrawFunds This is a method of Dice2Win
func (dw *Dice2Win) emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {
	type withdrawFunds struct {
		TokenName      string        `json:"tokenName"`
		Beneficiary    types.Address `json:"beneficiary"`
		WithdrawAmount bn.Number     `json:"withdrawAmount"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		withdrawFunds{
			TokenName:      tokenName,
			Beneficiary:    beneficiary,
			WithdrawAmount: withdrawAmount,
		},
	)
}

//emitPlaceBet This is a method of Dice2Win
func (dw *Dice2Win) emitPlaceBet(tokenName string, gambler types.Address, amount, betMask, possibleWinAmount bn.Number, modulo, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {
	type placeBet struct {
		TokenName         string        `json:"tokenName"`
		Gambler           types.Address `json:"gambler"`
		Amount            bn.Number     `json:"amount"`
		BetMask           bn.Number     `json:"betMask"`
		PossibleWinAmount bn.Number     `json:"possibleWinAmount"`
		Modulo            int64         `json:"modulo"`
		CommitLastBlock   int64         `json:"commitLastBlock"`
		Commit            []byte        `json:"commit"`
		SignData          []byte        `json:"signData"`
		RefAddress        types.Address `json:"refAddress"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		placeBet{
			TokenName:         tokenName,
			Gambler:           gambler,
			Amount:            amount,
			BetMask:           betMask,
			PossibleWinAmount: possibleWinAmount,
			Modulo:            modulo,
			CommitLastBlock:   commitLastBlock,
			Commit:            commit,
			SignData:          signData,
			RefAddress:        refAddress,
		},
	)
}

//emitSettleBet This is a method of Dice2Win
func (dw *Dice2Win) emitSettleBet(tokenName string, reveal []byte, gambler types.Address, diceWinAmount bn.Number, winNumber int64) {
	type settleBet struct {
		TokenName     string        `json:"tokenName"`
		Reveal        []byte        `json:"reveal"`
		Gambler       types.Address `json:"gambler"`
		DiceWinAmount bn.Number     `json:"diceWinAmount"`
		WinNumber     int64         `json:"winNumber"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		settleBet{
			TokenName:     tokenName,
			Reveal:        reveal,
			Gambler:       gambler,
			DiceWinAmount: diceWinAmount,
			WinNumber:     winNumber,
		},
	)
}

//emitRefundBet This is a method of Dice2Win
func (dw *Dice2Win) emitRefundBet(commit []byte, tokenName string, gambler types.Address, amount bn.Number) {
	type refundBet struct {
		Commit    []byte        `json:"commit"`
		TokenName string        `json:"tokenName"`
		Gambler   types.Address `json:"gambler"`
		Amount    bn.Number     `json:"amount"`
	}

	dw.sdk.Helper().ReceiptHelper().Emit(
		refundBet{
			Commit:    commit,
			TokenName: tokenName,
			Gambler:   gambler,
			Amount:    amount,
		},
	)
}
