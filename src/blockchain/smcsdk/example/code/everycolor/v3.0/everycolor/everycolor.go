package everycolor

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/crypto/ed25519"
	"blockchain/smcsdk/sdk/crypto/sha3"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/types"
	"fmt"
	"strconv"

	"blockchain/smcsdk/sdk/jsoniter"
	"encoding/hex"
)

//Everycolor This is struct of contract
//@:contract:everycolor
//@:version:3.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:bbe07ac1e7f26bc65918ddac5e49aad3fe813a576e4407f201d2005ca0bb7c36
type Everycolor struct {
	sdk sdk.ISmartContract

	//@:public:store:cache
	secretSigner types.PubKey

	//@:public:store
	betInfo map[string]map[string]*BetInfo // key1:string = commitStr , key2:string = index

	//@:public:store
	lockedAmount map[string]bn.Number // key = token name

	//@:public:store
	mapSetting *MapSetting

	//@:public:store:cache
	recvFeeInfos RecvFeeInfo

	//@:public:store
	roundInfo map[string]*RoundInfo

	//@:public:store
	playerIndex map[string]map[string]*PlayerIndexes
}

//@:public:receipt
type receipt interface {
	emitSetSecretSigner(newSecretSigner types.PubKey)
	emitSetSettings(settings map[string]Setting, betExpirationBlocks int64)
	emitSetRecvFeeInfos(infos RecvFeeInfo)
	emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number)
	emitPlaceBet(tokenName string, amount bn.Number, betData []BetData, possibleWinAmount bn.Number, commitLastBlock, betCount int64, commit, signData []byte, refAddress types.Address)
	emitSettleBet(reveal, commit []byte, winNumber string, startIndex, endIndex int64, amountOfWin, amountOfUnLock map[string]bn.Number, finished bool)
	emitWithdrawWin(commit []byte, amountOfWin, unLockAmount map[string]bn.Number)
	emitRefundBet(commit []byte, refundCount int64, unlockAmount map[string]bn.Number, finished bool)
}

//InitChain Constructor of this Everycolor
//@:constructor
func (e *Everycolor) InitChain() {
	// init data
	e.mapSetting = new(MapSetting)
	e.mapSetting.Settings = make(map[string]Setting, 0)
	e.mapSetting.Settings[e.sdk.Helper().GenesisHelper().Token().Name()] = Setting{
		bn.N(1950000000000),
		bn.N(20000000000),
		bn.N(100000000),
		50,
		bn.N(300000),
		100,
	}
	e.mapSetting.BetExpirationBlocks = 250

	e._setMapSetting(e.mapSetting)
	e._setLockedAmount(e.sdk.Helper().GenesisHelper().Token().Name(), bn.N(0))

}

// See comment for "secretSigner" variable.
//@:public:method:gas[500]
func (e *Everycolor) SetSecretSigner(newSecretSigner types.PubKey) {

	//only contract owner just can do it
	sdk.RequireOwner()

	//check SecretSigner length
	sdk.Require(
		len(newSecretSigner) == 32,
		types.ErrInvalidParameter,
		"length of newSecretSingner must be 32 bytes",
	)

	// set store
	e._setSecretSigner(newSecretSigner)
	// fire event
	e.emitSetSecretSigner(newSecretSigner)
	return
}

// Change max bet reward. Setting this to zero effectively disables betting.
//@:public:method:gas[500]
func (e *Everycolor) SetSettings(settings string) {

	//only contract owner just can do it
	sdk.RequireOwner()
	resultSettings := e.checkSettings(settings)
	// set
	e._setMapSetting(resultSettings)
	// fire event
	e.emitSetSettings(resultSettings.Settings, resultSettings.BetExpirationBlocks)
	return
}

//@:public:method:gas[500]
func (e *Everycolor) SetRecvFeeInfo(recvFeeInfoStr string) {

	// only contract owner just can do it
	sdk.RequireOwner()
	recFeeInfo := e.checkSetRecvFeeInfo(recvFeeInfoStr)

	rfi := e._recvFeeInfos()
	rfi.RecvFeeRatio = recFeeInfo.RecvFeeRatio
	rfi.RecvFeeAddr = recFeeInfo.RecvFeeAddr
	e._setRecvFeeInfos(rfi)
	//fire event
	e.emitSetRecvFeeInfos(*recFeeInfo)
	return
}

// Funds withdrawal to cover costs of everycolor.win operation.
//@:public:method:gas[500]
func (e *Everycolor) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {
	// only contract owner just can do it
	sdk.RequireOwner()

	sdk.Require(
		tokenName != "",
		types.ErrInvalidParameter, "tokenName cannot be empty",
	)

	sdk.Require(
		withdrawAmount.CmpI(0) > 0,
		types.ErrInvalidParameter, "WithdrawAmount must be larger than zero",
	)

	accounts := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
	balance := accounts.BalanceOfName(tokenName)
	sdk.Require(
		withdrawAmount.Cmp(balance) <= 0,
		types.ErrInvalidParameter,
		"WithdrawAmount cannot be larger than balance",
	)

	sdk.Require(
		beneficiary != "",
		types.ErrInvalidParameter,
		"beneficiary cannot be empty",
	)

	lockAmount := e._lockedAmount(tokenName)
	sdk.Require(
		lockAmount.Add(withdrawAmount).Cmp(balance) <= 0,
		types.ErrInvalidParameter,
		"Not enough funds",
	)

	// transfer to beneficiary
	accounts.TransferByName(tokenName, beneficiary, withdrawAmount)
	//fire event
	e.emitWithdrawFunds(tokenName, beneficiary, withdrawAmount)

	return

}

// Bet placing transaction - issued by the player.
//@:public:method:gas[500]
func (e *Everycolor) PlaceBet(tokenName string, amount bn.Number, betData string, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {
	//contract owner cannot do it
	sendrAddr := e.sdk.Message().Sender().Address()
	sdk.Require(
		sendrAddr != e.sdk.Message().Contract().Owner(),
		types.ErrNoAuthorization,
		"Contract owner cannot do PlaceBet()",
	)

	// check secretSinger
	secret := e._secretSigner()
	sdk.Require(
		len(secret) != 0,
		types.ErrInvalidParameter,
		"Must set secretSigner first",
	)

	// check signData
	data := append(bn.N(commitLastBlock).Bytes(), commit...)
	sdk.Require(
		ed25519.VerifySign(secret, data, signData),
		types.ErrInvalidParameter,
		"ECDSA signature is not valid",
	)

	sdk.Require(
		e.sdk.Block().Height() <= commitLastBlock,
		types.ErrInvalidParameter,
		"Commit has expired",
	)

	// Check that commit is valid - it has not expired and its signature is valid.
	commitStr := hex.EncodeToString(commit)
	e.checkCommit(commitStr)

	var roundInfo *RoundInfo
	if !e._chkRoundInfo(commitStr) {
		roundInfo = &RoundInfo{
			commit, // start new round
			make(map[string]bn.Number, 0),
			"",
			0,
			0,
			0,
			e.sdk.Block().Height(),
			e._mapSetting(),
		}
		e._setRoundInfo(commitStr, roundInfo)
	}
	roundInfo = e._roundInfo(commitStr)
	var betDataSlice []BetData
	jsonErr := jsoniter.Unmarshal([]byte(betData), &betDataSlice)
	sdk.RequireNotError(jsonErr, types.ErrInvalidParameter)

	// Check the betData
	modes := e.checkBetInfo(tokenName, amount, betDataSlice, roundInfo)

	// Winning amount and jackpot increase.
	totalLockedAmount := bn.N(0)
	feeNum := bn.N(0)
	setting := roundInfo.Setting.Settings[tokenName]
	for _, item := range modes {
		lockedAmount := item.ToLockAmount()
		totalLockedAmount = totalLockedAmount.Add(lockedAmount)

	}
	//check Possible Amount
	feeNum = amount.MulI(setting.FeeRatio).DivI(PERMILLI)
	e.checkPossibleWinAmount(tokenName, amount, totalLockedAmount, feeNum, roundInfo)

	ok := e._chkLockedAmount(tokenName)
	if ok == false {
		// 没有就存到map中
		e._setLockedAmount(tokenName, bn.N(0))
	}

	lockedAmount := e._lockedAmount(tokenName)
	e._setLockedAmount(tokenName, lockedAmount.Add(totalLockedAmount))

	//transfer to contract account
	// get transfer receipt and save value
	transferReceipt := e.sdk.Message().GetTransferToMe()
	sdk.Require(transferReceipt != nil,
		types.ErrInvalidParameter,
		"Player was not transfer to contract",
	)

	bet := BetInfo{
		sendrAddr,
		tokenName,
		amount,
		betDataSlice,
		bn.N(0),
		make([]bn.Number, 0),
		false,
	}

	ok = e._chkPlayerIndex(commitStr, sendrAddr)
	index := &PlayerIndexes{}
	if ok == false {
		index = &PlayerIndexes{make([]int64, 0)}
	} else {
		index = e._playerIndex(commitStr, sendrAddr)

	}
	index.BetIndexes = append(index.BetIndexes, roundInfo.BetCount)
	e._setPlayerIndex(commitStr, sendrAddr, index)
	e._setBetInfo(commitStr, fmt.Sprintf("%d", roundInfo.BetCount), &bet)

	//fire event
	e.emitPlaceBet(
		tokenName,
		amount,
		betDataSlice,
		totalLockedAmount,
		commitLastBlock,
		roundInfo.BetCount,
		commit,
		signData,
		refAddress,
	)

	_, ok = roundInfo.TotalBuy[tokenName]
	if ok == false {
		roundInfo.TotalBuy[tokenName] = bn.N(0)
	}
	roundInfo.TotalBuy[tokenName] = roundInfo.TotalBuy[tokenName].Add(amount)
	roundInfo.BetCount += 1
	e._setRoundInfo(commitStr, roundInfo)
}

//@:public:method:gas[500]
func (e *Everycolor) SettleBets(reveal []byte, settleCount int64) {

	//check length
	sdk.Require((len(reveal)) > 0,
		types.ErrInvalidParameter,
		"Commit should be not exist",
	)
	// "commit" for bet settlement can only be obtained by hashing a "reveal".
	commit := sha3.Sum256(reveal)
	commitStr := hex.EncodeToString(sha3.Sum256(reveal))

	roundInfo := e._roundInfo(commitStr)
	e.checkRoundInfoPlaceBet(roundInfo)
	// Lottery
	if len(roundInfo.WinNumber) == 0 {
		random := sha3.Sum256(reveal, e.sdk.Block().BlockHash(), e.sdk.Block().RandomNumber())
		winNumber := getLotteryNum(random) // 计算开牌结果的种类
		roundInfo.WinNumber = winNumber
	}

	if settleCount > MAXSETTLECOUNT || settleCount <= 0 {
		settleCount = MAXSETTLECOUNT
	}
	startIndex := roundInfo.SettledCount
	lastIndex := startIndex + settleCount - 1
	if lastIndex > roundInfo.BetCount-1 {
		lastIndex = roundInfo.BetCount - 1
	}

	totalPossibleWinAmount := make(map[string]bn.Number, 0)
	totalWinAmount := make(map[string]bn.Number, 0)
	totalFeeNum := make(map[string]bn.Number, 0)

	settings := roundInfo.Setting.Settings
	for index := startIndex; index <= lastIndex; index++ {

		bet := e._betInfo(commitStr, strconv.Itoa(int(index)))
		if bet.Settled {
			continue
		}
		possibleWinAmount, ecWinAmount, feeNum := e.SettleBet(bet, index, commitStr, roundInfo)

		// total lockeAmount & feeNum & ecWinAmount
		_, ok := totalPossibleWinAmount[bet.TokenName]
		if ok == false {
			totalPossibleWinAmount[bet.TokenName] = bn.N(0)
			totalWinAmount[bet.TokenName] = bn.N(0)
			totalFeeNum[bet.TokenName] = bn.N(0)
		}
		totalPossibleWinAmount[bet.TokenName] = totalPossibleWinAmount[bet.TokenName].Add(possibleWinAmount)
		totalWinAmount[bet.TokenName] = totalWinAmount[bet.TokenName].Add(ecWinAmount)
		totalFeeNum[bet.TokenName] = totalFeeNum[bet.TokenName].Add(feeNum)

	}

	forx.Range(totalPossibleWinAmount, func(tokenName string, value bn.Number) bool {
		//根据k的值减锁定金额
		e._setLockedAmount(commitStr, e._lockedAmount(tokenName).Sub(value))
		e.RatioTransfer(tokenName, totalFeeNum[tokenName], settings[tokenName])

		return true
	})

	roundInfo.SettledCount = lastIndex + 1
	e._setRoundInfo(commitStr, roundInfo)

	//fire event
	e.emitSettleBet(reveal,
		commit,
		roundInfo.WinNumber,
		startIndex,
		lastIndex,
		totalWinAmount,
		totalPossibleWinAmount,
		roundInfo.BetCount == roundInfo.SettledCount,
	)
	return
}

//@:public:method:gas[500]
func (e *Everycolor) WithdrawWin(commit []byte) {

	// setup data
	commitStr := hex.EncodeToString(commit)
	roundInfo := e._roundInfo(commitStr)

	sdk.Require(len(roundInfo.WinNumber) != 0,
		types.ErrInvalidParameter,
		"The round never run a lottery",
	)

	sdk.Require(roundInfo.SettledCount < roundInfo.BetCount,
		types.ErrInvalidParameter,
		"The round is complete",
	)

	// get plyrRnd
	betIndexs := e._playerIndex(commitStr, e.sdk.Message().Sender().Address())

	totalPossibleWinAmount := make(map[string]bn.Number, 0)
	totalWinAmount := make(map[string]bn.Number, 0)
	totalFeeNum := make(map[string]bn.Number, 0)

	for _, betIndex := range betIndexs.BetIndexes {
		bet := e._betInfo(commitStr, strconv.Itoa(int(betIndex)))
		if bet.Settled {
			continue
		}
		possibleWinAmount, ecWinAmount, feeNum := e.SettleBet(bet, betIndex, commitStr, roundInfo)

		// total lockeAmount & feeNum & ecWinAmount
		_, ok := totalPossibleWinAmount[bet.TokenName]
		if ok == false {
			totalPossibleWinAmount[bet.TokenName] = bn.N(0)
			totalWinAmount[bet.TokenName] = bn.N(0)
			totalFeeNum[bet.TokenName] = bn.N(0)
		}
		totalPossibleWinAmount[bet.TokenName] = totalPossibleWinAmount[bet.TokenName].Add(possibleWinAmount)
		totalWinAmount[bet.TokenName] = totalWinAmount[bet.TokenName].Add(ecWinAmount)
		totalFeeNum[bet.TokenName] = totalFeeNum[bet.TokenName].Add(feeNum)
	}

	//遍历字符串数组
	forx.Range(totalPossibleWinAmount, func(tokenName string, value bn.Number) bool {
		//根据k的值减锁定金额
		e._setLockedAmount(commitStr, e._lockedAmount(tokenName).Sub(value))

		e.RatioTransfer(tokenName, totalFeeNum[tokenName], roundInfo.Setting.Settings[tokenName])

		return true
	})

	// fire event
	e.emitWithdrawWin(commit, totalWinAmount, totalPossibleWinAmount)
	return
}

//@:public:method:gas[500]
func (e *Everycolor) RefundBets(commit []byte, refundCount int64) {

	// Check that bet is in 'active' state.
	commitStr := hex.EncodeToString(commit)
	roundInfo := e._roundInfo(commitStr)

	// Check that bet is in 'active' state.
	e.checkRoundInfoRefund(roundInfo)

	// get plyrRnd
	betIndexs := e._playerIndex(commitStr, e.sdk.Message().Sender().Address())
	refunded := int64(0)
	//temp calc unlock amount
	tmp := make(map[string]bn.Number, 0)

	// First, the refund request originator all the refund
	for _, betIndex := range betIndexs.BetIndexes {
		bet := e._betInfo(commitStr, strconv.Itoa(int(betIndex)))
		if bet.Settled == false {
			possibleWinAmount := e.Refund(bet, betIndex, commitStr)

			// tmp for unlock amount
			_, ok := tmp[bet.TokenName]
			if ok == false {
				tmp[bet.TokenName] = bn.N(0)
			}
			tmp[bet.TokenName] = tmp[bet.TokenName].Add(possibleWinAmount)

			//set new locked amount
			e._setLockedAmount(commitStr, e._lockedAmount(bet.TokenName).Sub(possibleWinAmount))
			refunded += 1
		}
	}

	if refundCount > MAXREFUNDCOUNT || refundCount <= 0 {
		refundCount = MAXREFUNDCOUNT
	}
	// setup data
	index := roundInfo.SettledCount
	lastIndex := index + refundCount - 1
	if lastIndex > roundInfo.BetCount-1 {
		lastIndex = roundInfo.BetCount - 1
	}

	for index <= lastIndex {
		bet := e._betInfo(commitStr, strconv.Itoa(int(index)))
		if bet.Settled == false {
			possibleWinAmount := e.Refund(bet, index, commitStr)

			// tmp for unlock amount
			_, ok := tmp[bet.TokenName]
			if ok == false {
				tmp[bet.TokenName] = bn.N(0)
			}
			tmp[bet.TokenName] = tmp[bet.TokenName].Add(possibleWinAmount)

			//set new locked amount
			e._setLockedAmount(commitStr, e._lockedAmount(bet.TokenName).Sub(possibleWinAmount))
			refunded += 1
		}
		index += 1
	}

	roundInfo.RefundCount = roundInfo.RefundCount + refunded
	roundInfo.SettledCount = lastIndex + 1
	e._setRoundInfo(commitStr, roundInfo)

	//fire event
	e.emitRefundBet(commit, refunded, tmp, roundInfo.RefundCount == roundInfo.BetCount)
	return
}
