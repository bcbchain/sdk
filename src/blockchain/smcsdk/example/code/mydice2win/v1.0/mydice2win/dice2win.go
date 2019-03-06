package mydice2win

import (
	"encoding/hex"

	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/crypto/ed25519"
	"blockchain/smcsdk/sdk/crypto/sha3"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
)

//Dice2Win a demo contract
//@:contract:mydice2win
//@:version:1.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:2b6cab2f53a83f2d08807010533bc53785edfda0aed55028336914ccebadbc94
type Dice2Win struct {
	sdk sdk.ISmartContract

	//@:public:store:cache
	secretSigner types.PubKey

	//@:public:store
	bet map[string]*Bet // key = token name

	//@:public:store:cache
	lockedAmount map[string]bn.Number // key = token name

	//@:public:store:cache
	settings *Settings

	//@:public:store:cache
	recvFeeInfos []RecvFeeInfo
}

//@:public:receipt
type receipt interface {
	emitSetSecretSigner(newSecretSigner types.PubKey)
	emitSetSettings(tokenNames []string, minBet, maxBet, maxProfit, feeRatio, feeMinimum, sendToCltRatio, betExpirationBlocks int64)
	emitSetRecvFeeInfos(infos []RecvFeeInfo)
	emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number)
	emitPlaceBet(tokenName string, gambler types.Address, amount, betMask, possibleWinAmount bn.Number, modulo, commitLastBlock int64, commit, signData []byte, refAddress types.Address)
	emitSettleBet(tokenName string, reveal []byte, gambler types.Address, diceWinAmount bn.Number, winNumber int64)
	emitRefundBet(commit []byte, tokenName string, gambler types.Address, amount bn.Number)
}

// InitChain - construct function
//@:constructor
func (dw *Dice2Win) InitChain() {
	// init data
	settings := Settings{}
	settings.TokenNames = []string{dw.sdk.Helper().GenesisHelper().Token().Name()}
	settings.MaxProfit = 2E12
	settings.MaxBet = 2E10
	settings.MinBet = 1E8
	settings.SendToCltRatio = 100
	settings.FeeRatio = 50
	settings.FeeMinimum = 300000
	settings.BetExpirationBlocks = 250

	dw._setSettings(&settings)
	dw._setLockedAmount(dw.sdk.Helper().GenesisHelper().Token().Name(), bn.N(0))
}

// SetSecretSigner - Set the secret signer
//@:public:method:gas[500]
func (dw *Dice2Win) SetSecretSigner(newSecretSigner types.PubKey) {

	sdk.RequireOwner(dw.sdk)
	sdk.Require(len(newSecretSigner) == 32,
		types.ErrInvalidParameter, "length of newSecretSigner must be 32 bytes")

	// save secret signer
	dw._setSecretSigner(newSecretSigner)

	// fire event
	dw.emitSetSecretSigner(newSecretSigner)
}

// SetSettings - Change game settings
//@:public:method:gas[500]
func (dw *Dice2Win) SetSettings(newSettingsStr string) {

	sdk.RequireOwner(dw.sdk)

	//只有在全部结算完成，退款完成后，才能设置settings
	settings := dw._settings()
	for _, tokenName := range settings.TokenNames {
		lockedAmount := dw._lockedAmount(tokenName)
		sdk.Require(lockedAmount.CmpI(0) == 0,
			types.ErrUserDefined, "only lockedAmount is zero that can do SetSettings()")
	}

	newSettings := new(Settings)
	err := jsoniter.Unmarshal([]byte(newSettingsStr), newSettings)
	sdk.RequireNotError(err, types.ErrInvalidParameter)

	dw.checkSettings(newSettings)
	dw._setSettings(newSettings)

	// fire event
	dw.emitSetSettings(
		newSettings.TokenNames,
		newSettings.MinBet,
		newSettings.MaxBet,
		newSettings.MaxProfit,
		newSettings.FeeRatio,
		newSettings.FeeMinimum,
		newSettings.SendToCltRatio,
		newSettings.BetExpirationBlocks,
	)
}

// SetRecvFeeInfos - Set ratio of fee and receiver's account address
//@:public:method:gas[500]
func (dw *Dice2Win) SetRecvFeeInfos(recvFeeInfosStr string) {

	sdk.RequireOwner(dw.sdk)

	infos := make([]RecvFeeInfo, 0)
	err := jsoniter.Unmarshal([]byte(recvFeeInfosStr), &infos)
	sdk.RequireNotError(err, types.ErrInvalidParameter)

	dw.checkRecvFeeInfos(infos)
	dw._setRecvFeeInfos(infos)

	// fire event
	dw.emitSetRecvFeeInfos(infos)
}

// WithdrawFunds - Funds withdrawal
//@:public:method:gas[500]
func (dw *Dice2Win) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {

	sdk.RequireOwner(dw.sdk)
	sdk.Require(withdrawAmount.CmpI(0) > 0,
		types.ErrInvalidParameter, "withdrawAmount must be larger than zero")

	account := dw.sdk.Helper().AccountHelper().AccountOf(dw.sdk.Message().Contract().Account())
	lockedAmount := dw._lockedAmount(tokenName)
	unlockedAmount := account.BalanceOfName(tokenName).Sub(lockedAmount)
	sdk.Require(unlockedAmount.Cmp(withdrawAmount) >= 0,
		types.ErrInvalidParameter, "Not enough funds")

	// transfer to beneficiary
	account.TransferByName(tokenName, beneficiary, withdrawAmount)

	// fire event
	dw.emitWithdrawFunds(tokenName, beneficiary, withdrawAmount)
}

// PlaceBet - Issued by the gambler to place a bet
//@:public:method:gas[500]
//@:public:interface:gas[400]
func (dw *Dice2Win) PlaceBet(betMask bn.Number, modulo, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {

	//contract owner cannot do it
	sdk.Require(dw.sdk.Message().Sender().Address() != dw.sdk.Message().Contract().Owner(),
		types.ErrNoAuthorization, "contract owner cannot do PlaceBet")

	// Check that commit is valid - it has not expired and its signature is valid and must be new
	data := append(bn.N(commitLastBlock).Bytes(), commit...)
	hexCommit := hex.EncodeToString(commit)
	sdk.Require(ed25519.VerifySign(dw._secretSigner(), data, signData),
		types.ErrInvalidParameter, "ECDSA signature is not valid")
	sdk.Require(dw.sdk.Block().Height() <= commitLastBlock,
		types.ErrInvalidParameter, "Commit has expired")
	sdk.Require(dw._chkBet(hexCommit) == false,
		types.ErrInvalidParameter, "Commit should be new")

	// Validate input data ranges.
	sdk.Require(modulo > 1 && modulo <= maxModulo,
		types.ErrInvalidParameter, "Modulo should be within range")
	sdk.Require(betMask.CmpI(0) > 0 && betMask.Cmp(dw.maxBetMask()) < 0,
		types.ErrInvalidParameter, "Mask should be within range")

	amount := bn.N(0)
	tokenName := ""

	// get transfer receipt and save value
	settings := dw._settings()
	for _, tkName := range settings.TokenNames {
		transferReceipt := dw.sdk.Message().GetTransferToMe(tkName)
		if transferReceipt != nil {
			tokenName = tkName
			amount = transferReceipt.Value
			break
		}
	}
	sdk.Require(tokenName != "" && amount.CmpI(0) > 0,
		types.ErrUserDefined, "Must transfer tokens to me before place a bet")
	sdk.Require(amount.CmpI(settings.MinBet) >= 0 && amount.CmpI(settings.MaxBet) <= 0,
		types.ErrInvalidParameter, "Amount should be within range")

	var rollUnder, mask bn.Number

	if modulo <= maxMaskModulo {
		//当modulo<40的时候，这个变量保存了mask里面有多少个1,比如mask==6==0000110/mask==3==0000011（rollUnder=2）  mask==7==000111（rollUnder=3）
		rollUnder = betMask.Mul(dw.popCntMult()).And(dw.popCntMask()).Mod(dw.popCntModulo())
		mask = betMask
	} else {
		sdk.Require(betMask.CmpI(0) > 0 && betMask.CmpI(modulo) <= 0,
			types.ErrInvalidParameter, "High modulo range, betMask larger than modulo")

		rollUnder = betMask
	}

	// Calc possible win amount and fee
	possibleWinAmount, feeAmount := dw.getDiceWinAmount(amount, bn.N(modulo), rollUnder)

	// Enforce max profit limit
	sdk.Require(possibleWinAmount.CmpI(settings.MaxProfit) <= 0,
		types.ErrInvalidParameter, "MaxProfit limit violation")

	// Check whether contract account has enough funds to process this bet.
	contractAcct := dw.sdk.Helper().AccountHelper().AccountOf(dw.sdk.Message().Contract().Account())
	totalLockedAmount := dw._lockedAmount(tokenName)
	totalLockedAmount = totalLockedAmount.Add(possibleWinAmount).Add(feeAmount)
	totalUnlockedAmount := contractAcct.BalanceOfName(tokenName).Sub(totalLockedAmount)
	sdk.Require(totalUnlockedAmount.Cmp(possibleWinAmount) >= 0,
		types.ErrInvalidParameter, "Cannot afford to lose this bet")
	dw._setLockedAmount(tokenName, totalLockedAmount)

	bet := &Bet{}
	bet.TokenName = tokenName
	bet.Amount = amount
	bet.Modulo = modulo
	bet.Mask = mask
	bet.RollUnder = rollUnder
	bet.PlaceBlockNumber = dw.sdk.Block().Height()
	bet.Gambler = dw.sdk.Message().Sender().Address()
	dw._setBet(hex.EncodeToString(commit), bet)

	//fire event
	dw.emitPlaceBet(tokenName, bet.Gambler, amount, betMask, possibleWinAmount, modulo, commitLastBlock, commit, signData, refAddress)
}

// SettleBet - Settle the bet and transfer winMoney、cltFee and other fee to destination
//@:public:method:gas[500]
func (dw *Dice2Win) SettleBet(reveal []byte) {

	// "commit" for bet settlement can only be obtained by hashing a "reveal".
	hexCommit := hex.EncodeToString(sha3.Sum256(reveal))
	sdk.Require(dw._chkBet(hexCommit) == true,
		types.ErrInvalidParameter, "Commit should be exist")

	settings := dw._settings()
	bet := dw._bet(hexCommit)
	sdk.Require(bet.Amount.CmpI(0) > 0,
		types.ErrUserDefined, "Bet has already settled.")
	sdk.Require(dw.sdk.Block().Height() > bet.PlaceBlockNumber,
		types.ErrUserDefined, "SettleBet block can not be in the same block as placeBet, or before.")
	sdk.Require(dw.sdk.Block().Height() <= (bet.PlaceBlockNumber+settings.BetExpirationBlocks),
		types.ErrUserDefined, "The lottery time is out of date")

	//get random and to mod modulo
	entropy := bn.NBytes(sha3.Sum256(reveal, dw.sdk.Block().BlockHash(), dw.sdk.Block().RandomNumber()))
	dice := entropy.Mod(bn.N(bet.Modulo)) //用随机数对modulo取膜
	diceWinAmount, feeAmount := dw.getDiceWinAmount(bet.Amount, bn.N(bet.Modulo), bet.RollUnder)
	diceWin := bn.N(0)

	// Determine dice outcome
	if bet.Modulo <= maxMaskModulo {
		// For small modulo games, check the outcome against a bit mask.
		if bn.N(2).Exp(dice).And(bet.Mask).CmpI(0) != 0 {
			//比如dice是0,你选的mask==000011==3,那么2的0次方等于1==000001， 与上mask==1 ！= 0 ，所以中奖
			diceWin = diceWinAmount
		}
	} else {
		// For larger modulos, check inclusion into half-open interval.
		if dice.Cmp(bet.RollUnder) < 0 {
			diceWin = diceWinAmount
		}
	}

	// Unlock the bet amount, regardless of the outcome.
	lockedAmount := dw._lockedAmount(bet.TokenName)
	lockedAmount = lockedAmount.Sub(diceWinAmount).Sub(feeAmount)
	dw._setLockedAmount(bet.TokenName, lockedAmount)

	// Send the win funds to gambler.
	contractAcct := dw.sdk.Helper().AccountHelper().AccountOf(dw.sdk.Message().Contract().Account())
	if diceWin.CmpI(0) > 0 {
		contractAcct.TransferByName(bet.TokenName, bet.Gambler, diceWin)
	}

	// Move bet into 'processed' state already.
	bet.Amount = bn.N(0)
	dw._setBet(hexCommit, bet)

	// Send fee to clt
	sentToCltFee := feeAmount.MulI(settings.SendToCltRatio).DivI(perMille)
	if settings.SendToCltRatio > 0 {
		cltAcct := dw.sdk.Helper().BlockChainHelper().CalcAccountFromName("clt", "")
		contractAcct.TransferByName(bet.TokenName, cltAcct, sentToCltFee)
	}

	// Send fee to recvFeeInfos
	dw.transferToRecvFeeAddr(bet.TokenName, feeAmount.Sub(sentToCltFee))

	//fire event
	dw.emitSettleBet(bet.TokenName, reveal, bet.Gambler, diceWin, dice.V.Int64())
}

// RefundBet - Refund bet when it's out of time
//@:public:method:gas[500]
func (dw *Dice2Win) RefundBet(commit []byte) {

	hexCommit := hex.EncodeToString(commit)
	sdk.Require(dw._chkBet(hexCommit) == true,
		types.ErrInvalidParameter, "Commit should be exist")

	bet := dw._bet(hexCommit)
	betAmount := bet.Amount
	sdk.Require(bet.Amount.CmpI(0) > 0,
		types.ErrInvalidParameter, "Bet should not be settled")

	// Check that bet has already expired.
	settings := dw._settings()
	sdk.Require(dw.sdk.Block().Height() > bet.PlaceBlockNumber+settings.BetExpirationBlocks,
		types.ErrUserDefined, "The refund time has not been arrived")

	possibleWinAmount, feeAmount := dw.getDiceWinAmount(bet.Amount, bn.N(bet.Modulo), bet.RollUnder)
	lockedAmount := dw._lockedAmount(bet.TokenName)
	lockedAmount = lockedAmount.Sub(possibleWinAmount).Sub(feeAmount)
	dw._setLockedAmount(bet.TokenName, lockedAmount)

	//Send the funds to gambler.
	contractAcct := dw.sdk.Helper().AccountHelper().AccountOf(dw.sdk.Message().Contract().Account())
	contractAcct.TransferByName(bet.TokenName, bet.Gambler, bet.Amount)

	// Move bet into 'processed' state, release funds.
	bet.Amount = bn.N(0)
	dw._setBet(hexCommit, bet)

	// fire event
	dw.emitRefundBet(commit, bet.TokenName, bet.Gambler, betAmount)
}
