package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/crypto/ed25519"
	"blockchain/smcsdk/sdk/crypto/sha3"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
	"encoding/hex"
	"fmt"
	"strconv"
)

//excellencies This is struct of contract
//@:contract:excellencies
//@:version:2.0
//@:organization:orgCZkw5xz9DYa3h5pJ2CzZSuGHRCj2ot5xq
//@:author:ef94556a937618c72ffaf173b1533c533d77aa3ea2a63f053bb904feefe5a92f
type Excellencies struct {
	sdk sdk.ISmartContract

	//@:public:store:cache
	secretSigner types.PubKey // Check to sign the public key

	//@:public:store
	betInfo map[string]map[string]*BetInfo // key1:string = commitStr , key2:string = index

	////@:public:store:cache
	//lockedInBets map[string]bn.Number // Lock amount (unit cong) key: currency name

	//@:public:store
	mapSetting *MapSetting

	//@:public:store:cache
	recFeeInfo []RecFeeInfo

	//@:public:store
	roundInfo map[string]RoundInfo //key1轮标识

	//@:public:store:cache
	poolAmount map[types.Address]map[string]bn.Number

	//@:public:store
	playerIndex map[string]map[string]*PlayerIndexes // key1:string = commitStr , key2:string = index

	//@:public:store:cache
	grandPrizer map[types.Address]map[string]GrandPrizer //key1 =saler key2 = tokenName
}

//@:public:receipt
type receipt interface {
	emitSetSecretSigner(newSecretSigner types.PubKey)
	emitSetSettings(Settings map[string]Setting, BetExpirationBlocks, PoolFeeRatio, CarveUpPoolRatio, GrandPrizeRatio int64)
	emitPlaceBet(tokenName string, amount bn.Number, betData []BetData, possibleWinAmount bn.Number, commitLastBlock, betCount int64, commit, signData []byte, refAddress types.Address)
	emitSetRecFeeInfo(info []RecFeeInfo)
	emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number)
	emitSettleBet(reveal, commit []byte, banker GamerInfo, players map[string]GamerInfo, startIndex, endIndex int64, amountOfWin, amountOfUnLock map[string]bn.Number, finished bool, poolAmount map[string]map[string]bn.Number)
	emitWithdrawWin(commit []byte, amountOfWin, amountOfUnLock map[string]bn.Number, poolAmount map[string]map[string]bn.Number)
	emitRefundBet(commit []byte, refundCount int64, unlockAmount map[string]bn.Number, finished bool)
}

//InitChain Constructor of this excellencies
//@:constructor
func (sg *Excellencies) InitChain() {

	// init data
	sg.mapSetting = new(MapSetting)
	sg.mapSetting.Settings = make(map[string]Setting, 0)
	sg.mapSetting.Settings[sg.sdk.Helper().GenesisHelper().Token().Name()] = Setting{
		bn.N(1950000000000),
		bn.N(20000000000),
		bn.N(100000000),
		50,
		bn.N(300000),
		100,
	}
	sg.mapSetting.BetExpirationBlocks = 250
	sg.mapSetting.GrandPrizeRatio = 500
	sg.mapSetting.CarveUpPoolRatio = 50
	sg.mapSetting.PoolFeeRatio = 20

	sg._setMapSetting(sg.mapSetting)
	//sg._setLockedInBets(sg.sdk.Helper().GenesisHelper().Token().Name(), bn.N(0))

}

//@:import:directsale
type directsale interface {
	Income(salerAddr types.Address, lockAmount bn.Number, note string)
	Pay(salerAddr types.Address, tokenName string, incomeAmount, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string)
	Refund(salerAddr types.Address, tokenName string, unlockAmount bn.Number, amountList []bn.Number, toAddrList []types.Address, note string)
}

//SetSecretSigner - Set up the public key
//@:public:method:gas[5000]
func (sg *Excellencies) SetSecretSigner(newSecretSigner types.PubKey) {

	sdk.RequireOwner()
	sdk.Require(len(newSecretSigner) == 32,
		types.ErrInvalidParameter, "length of newSecretSigner must be 32 bytes")

	//Save to database
	sg._setSecretSigner(newSecretSigner)

	// fire event
	sg.emitSetSecretSigner(newSecretSigner)
}

//SetSettings - Change game settings
//@:public:method:gas[5000]
func (sg *Excellencies) SetSettings(newSettinsStr string) {
	sdk.RequireOwner()

	newSettings := sg.checkSettings(newSettinsStr)

	sg._setMapSetting(newSettings)
	// fire event
	sg.emitSetSettings(newSettings.Settings, newSettings.BetExpirationBlocks, newSettings.PoolFeeRatio, newSettings.CarveUpPoolRatio, newSettings.GrandPrizeRatio)
}

// SetRecFeeInfo - Set ratio of fee and receiver's account address
//@:public:method:gas[5000]
func (sg *Excellencies) SetRecFeeInfo(recFeeInfoStr string) {

	sdk.RequireOwner()

	info := make([]RecFeeInfo, 0)
	err := jsoniter.Unmarshal([]byte(recFeeInfoStr), &info)
	sdk.RequireNotError(err, types.ErrInvalidParameter)
	//Check that the parameters are valid
	sg.checkRecFeeInfo(info)

	sg._setRecFeeInfo(info)
	// fire event
	sg.emitSetRecFeeInfo(info)
}

//PlaceBet - place bet
//@:public:method:gas[500]
func (sg *Excellencies) PlaceBet(betJson string, commitLastBlock int64, commit, signData []byte, refAddress types.Address, saler types.Address) {
	//contract owner cannot do it
	sendrAddr := sg.sdk.Message().Sender().Address()
	sdk.Require(
		sendrAddr != sg.sdk.Message().Contract().Owner(),
		types.ErrNoAuthorization,
		"Contract owner cannot do PlaceBet()",
	)

	// check secretSinger
	secret := sg._secretSigner()
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
		sg.sdk.Block().Height() <= commitLastBlock,
		types.ErrInvalidParameter,
		"Commit has expired",
	)

	// Check that commit is valid - it has not expired and its signature is valid.
	commitStr := hex.EncodeToString(commit)
	sg.checkCommit(commitStr)

	// check amount info
	amount := bn.N(0)
	tokenName := ""

	// get transfer receipt and save value
	mapsettings := sg._mapSetting()
	transferReceipts := sg.sdk.Message().GetTransferToMe()
	forx.Range(transferReceipts, func(_, receipt std.Transfer) bool {

		token := sg.sdk.Helper().TokenHelper().TokenOfAddress(receipt.Token)
		if _, ok := mapsettings.Settings[token.Name()]; ok {
			tokenName = token.Name()
			amount = receipt.Value
			return forx.Break
		}
		return forx.Continue
	})

	sdk.Require(
		tokenName != "",
		types.ErrInvalidParameter,
		"Token name is invalid",
	)

	sg.PayGrandPrize(saler)

	roundInfo := RoundInfo{}
	if !sg._chkRoundInfo(commitStr) {
		roundInfo = RoundInfo{
			commit, // start new round
			make(map[types.Address]map[string]map[string]bn.Number, 0),
			make(map[types.Address]map[string]bn.Number, 0),
			GamerInfo{},
			make(map[string]GamerInfo, 0),
			0,
			0,
			0,
			sg.sdk.Block().Height(),
			sg._mapSetting(),
		}
		sg._setRoundInfo(commitStr, roundInfo)
	}
	roundInfo = sg._roundInfo(commitStr)

	poolAmount := bn.N(0)
	if sg._chkPoolAmount(saler, tokenName) {
		poolAmount = sg._poolAmount(saler, tokenName)
	}
	_, ok := roundInfo.SettledPoolAmount[saler]

	if !ok {
		roundInfo.SettledPoolAmount[saler] = make(map[string]bn.Number)
	}
	if !ok {
		roundInfo.SettledPoolAmount[saler][tokenName] = bn.N(0)
	}
	roundInfo.SettledPoolAmount[saler][tokenName] = poolAmount

	var betDataSlice []BetData
	jsonErr := jsoniter.Unmarshal([]byte(betJson), &betDataSlice)
	sdk.RequireNotError(jsonErr, types.ErrInvalidParameter)

	// Check the betData
	modes := sg.checkBetInfo(tokenName, amount, betDataSlice, &roundInfo)

	// calc locked amount
	totalLockedAmount := bn.N(0)
	feeNum := bn.N(0)
	setting := roundInfo.Setting.Settings[tokenName]
	forx.Range(modes, func(_, item BaseMode) bool {

		lockedAmount := item.ToLockAmount()
		totalLockedAmount = totalLockedAmount.Add(lockedAmount)
		item.SetRoundInfoTotalBuy(saler, tokenName, &roundInfo)
		return true
	})

	//check Possible Amount
	feeNum = amount.MulI(setting.FeeRatio).DivI(PERMILLI)
	sg.checkPossibleWinAmount(tokenName, amount, totalLockedAmount, feeNum, &roundInfo)

	// 跨合约调用
	sg.directsale().Income(saler, totalLockedAmount, "")

	bet := BetInfo{
		sg.sdk.Block().Time(),
		sendrAddr,
		saler,
		tokenName,
		amount,
		betDataSlice,
		bn.N(0),
		make([]bn.Number, 0),
		false,
	}

	ok = sg._chkPlayerIndex(commitStr, sendrAddr)
	index := &PlayerIndexes{}
	if ok == false {
		index = &PlayerIndexes{make([]int64, 0)}
	} else {
		index = sg._playerIndex(commitStr, sendrAddr)
	}
	index.BetIndexes = append(index.BetIndexes, roundInfo.BetCount)
	sg._setPlayerIndex(commitStr, sendrAddr, index)
	sg._setBetInfo(commitStr, fmt.Sprintf("%d", roundInfo.BetCount), &bet)

	//fire event
	sg.emitPlaceBet(
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

	roundInfo.BetCount += 1
	sg._setRoundInfo(commitStr, roundInfo)
}

// SettleBet - The lottery and settlement
//@:public:method:gas[500]
func (sg *Excellencies) SettleBet(reveal []byte, settleCount int64) {
	sdk.Require(len(reveal) > 0,
		types.ErrInvalidParameter, "Commit should be not exist")

	sdk.RequireOwner()
	hexCommit := hex.EncodeToString(sha3.Sum256(reveal))
	if !sg._chkRoundInfo(hexCommit) {
		sdk.Require(false, types.ErrInvalidParameter, "Commit should be not exist")
	}

	roundInfo := sg._roundInfo(hexCommit)
	sg.checkRoundInfoPlaceBet(&roundInfo)

	if roundInfo.SettledCount == 0 {
		random := sha3.Sum256(reveal, sg.sdk.Block().BlockHash(), sg.sdk.Block().RandomNumber())
		sg.Lottery(random, &roundInfo) // 计算开牌结果的种类
	}

	//Initial index
	startIndex := roundInfo.SettledCount
	lastIndex := startIndex + settleCount - 1

	if lastIndex > roundInfo.BetCount {
		lastIndex = roundInfo.BetCount
	}

	//The actual winning money
	totalPossibleWinAmount := make(map[string]bn.Number, 0)
	totalWinAmount := make(map[string]bn.Number, 0)

	poolTokens := make(map[string]map[string]bn.Number)
	grandSaler := make(map[string]int64)

	clt := sg.sdk.Helper().ContractHelper().ContractOfAddress("clt").Account()
	account := sg.sdk.Message().Contract().Account()

	forx.Range(startIndex, lastIndex+1, func(index int) bool {

		betInfo := sg._betInfo(hexCommit, fmt.Sprintf("%d", index))
		if betInfo.Settled == true {
			return forx.Continue
		}

		//The currency name of the bet such as BCB
		tokenName := betInfo.TokenName

		possibleWinAmount, sgWinAmount, feeNum, poolAmount := sg.settleBet(betInfo, int64(index), hexCommit, &roundInfo)

		salerAddr := betInfo.Saler
		//save this round saler address
		_, ok := grandSaler[betInfo.Saler]
		if !ok {
			grandSaler[salerAddr] = 0
		}

		addressList, amountList := sg.calcCltAndRecAmount(
			feeNum,
			roundInfo.Setting.Settings[tokenName],
			clt,
		)
		//CLT & RecFee reward
		if sgWinAmount.CmpI(0) > 0 {
			addressList = append(addressList, betInfo.Gambler)
			amountList = append(amountList, sgWinAmount)
		}

		// pool amount to contract account
		addressList = append(addressList, account)
		amountList = append(amountList, poolAmount)

		//transfer by directsale contract
		sg.directsale().Pay(
			salerAddr,
			betInfo.TokenName,
			betInfo.Amount,
			possibleWinAmount,
			[]bn.Number{sgWinAmount},
			[]types.Address{betInfo.Gambler},
			"",
		)

		// total lockeAmount & feeNum & ecWinAmount
		_, ok = totalPossibleWinAmount[betInfo.TokenName]
		if ok == false {
			totalPossibleWinAmount[betInfo.TokenName] = bn.N(0)
			totalWinAmount[betInfo.TokenName] = bn.N(0)
		}
		totalPossibleWinAmount[betInfo.TokenName] = totalPossibleWinAmount[betInfo.TokenName].Add(possibleWinAmount)
		totalWinAmount[betInfo.TokenName] = totalWinAmount[betInfo.TokenName].Add(sgWinAmount)
		//增加奖池金额
		pool := bn.N(0)
		if sg._chkPoolAmount(salerAddr, tokenName) {
			pool = sg._poolAmount(salerAddr, tokenName)
		}

		amount := pool.Add(poolAmount)
		sg._setPoolAmount(salerAddr, tokenName, amount)

		_, ok = poolTokens[salerAddr]
		if !ok {
			poolTokens[salerAddr] = make(map[string]bn.Number)
		}

		_, ok = poolTokens[salerAddr][tokenName]
		if !ok {
			poolTokens[salerAddr][tokenName] = bn.N(0)
		}
		poolTokens[salerAddr][tokenName] = poolTokens[salerAddr][tokenName].Add(amount)
		return forx.Continue
	})

	roundInfo.SettledCount = lastIndex
	sg._setRoundInfo(hexCommit, roundInfo)

	//Transfer to other handling address
	forx.Range(grandSaler, func(address types.Address, _ int) bool {
		sg.PayGrandPrize(address)
		return true
	})

	//Send the receipt
	sg.emitSettleBet(
		reveal,
		roundInfo.Commit,
		roundInfo.Banker,
		roundInfo.Players,
		startIndex, lastIndex,
		totalWinAmount,
		totalPossibleWinAmount,
		roundInfo.BetCount == roundInfo.SettledCount,
		poolTokens,
	)

}

//WithdrawWin - Player settlement
//@:public:method:gas[500]
func (sg *Excellencies) WithdrawWin(commit []byte) {
	hexCommit := hex.EncodeToString(commit)
	temp := sg._roundInfo(hexCommit)
	fmt.Println(temp)
	if !sg._chkRoundInfo(hexCommit) {
		sdk.Require(false, types.ErrInvalidParameter, "Commit should be not exist")
	}

	roundInfo := sg._roundInfo(hexCommit)
	//The bet height of the round to be settled should be less than the settlement height
	sdk.Require(roundInfo.FirstBlockHeight < sg.sdk.Block().Height(),
		types.ErrInvalidParameter, "SettleBet block can not be in the same block as placeBet, or before.")

	sdk.Require(roundInfo.BetCount > roundInfo.SettledCount,
		types.ErrInvalidParameter, "This round is complete")

	// get plyrRnd
	betIndexs := sg._playerIndex(hexCommit, sg.sdk.Message().Sender().Address())

	//The actual winning money
	totalPossibleWinAmount := make(map[string]bn.Number, 0)
	totalWinAmount := make(map[string]bn.Number, 0)

	poolTokens := make(map[string]map[string]bn.Number)
	grandSaler := make(map[string]int64)

	clt := sg.sdk.Helper().ContractHelper().ContractOfAddress("clt").Account()

	forx.Range(betIndexs.BetIndexes, func(_, index int64) bool {

		betInfo := sg._betInfo(hexCommit, fmt.Sprintf("%d", index))
		if betInfo.Settled == true {
			return forx.Continue
		}

		//The currency name of the bet such as BCB
		tokenName := betInfo.TokenName

		possibleWinAmount, sgWinAmount, feeNum, poolAmount := sg.settleBet(betInfo, index, hexCommit, &roundInfo)

		salerAddr := betInfo.Saler
		//save this round saler address
		_, ok := grandSaler[betInfo.Saler]
		if !ok {
			grandSaler[salerAddr] = 0
		}

		addressList, amountList := sg.calcCltAndRecAmount(
			feeNum,
			roundInfo.Setting.Settings[tokenName],
			clt,
		)
		if sgWinAmount.CmpI(0) > 0 {
			addressList = append(addressList, betInfo.Gambler)
			amountList = append(amountList, sgWinAmount)
		}

		//transfer by directsale contract
		sg.directsale().Pay(
			salerAddr,
			betInfo.TokenName,
			betInfo.Amount,
			possibleWinAmount,
			[]bn.Number{sgWinAmount},
			[]types.Address{betInfo.Gambler},
			"",
		)

		// total lockeAmount & feeNum & ecWinAmount
		_, ok = totalPossibleWinAmount[betInfo.TokenName]
		if ok == false {
			totalPossibleWinAmount[betInfo.TokenName] = bn.N(0)
			totalWinAmount[betInfo.TokenName] = bn.N(0)
		}
		totalPossibleWinAmount[betInfo.TokenName] = totalPossibleWinAmount[betInfo.TokenName].Add(possibleWinAmount)
		totalWinAmount[betInfo.TokenName] = totalWinAmount[betInfo.TokenName].Add(sgWinAmount)

		//增加奖池金额
		pool := bn.N(0)
		if sg._chkPoolAmount(salerAddr, tokenName) {
			pool = sg._poolAmount(salerAddr, tokenName)
		}

		amount := pool.Add(poolAmount)
		sg._setPoolAmount(salerAddr, tokenName, amount)

		_, ok = poolTokens[salerAddr]
		if !ok {
			poolTokens[salerAddr] = make(map[string]bn.Number)
		}

		_, ok = poolTokens[salerAddr][tokenName]
		if !ok {
			poolTokens[salerAddr][tokenName] = bn.N(0)
		}
		poolTokens[salerAddr][tokenName] = poolTokens[salerAddr][tokenName].Add(amount)
		return forx.Continue
	})

	//Transfer to other handling address

	forx.Range(grandSaler, func(address types.Address, _ int) bool {
		sg.PayGrandPrize(address)
		return forx.Continue
	})

	//fire event
	sg.emitWithdrawWin(commit, totalWinAmount, totalPossibleWinAmount, poolTokens)
}

//@:public:method:gas[500]
func (sg *Excellencies) RefundBets(commit []byte, refundCount int64) {

	// Check that bet is in 'active' state.
	commitStr := hex.EncodeToString(commit)
	roundInfo := sg._roundInfo(commitStr)

	// Check that bet is in 'active' state.
	sg.checkRoundInfoRefund(&roundInfo)

	// get plyrRnd
	betIndexs := sg._playerIndex(commitStr, sg.sdk.Message().Sender().Address())
	refunded := int64(0)
	//temp calc unlock amount
	tmp := make(map[string]bn.Number, 0)
	grandSaler := make(map[string]int64)

	// First, the refund request originator all the refund
	forx.Range(betIndexs.BetIndexes, func(_, betIndex int) bool {
		bet := sg._betInfo(commitStr, strconv.Itoa(int(betIndex)))
		if bet.Settled == false {
			possibleWinAmount := sg.Refund(bet, int64(betIndex), commitStr)

			salerAddress := bet.Saler
			sg.directsale().Refund(
				salerAddress,
				bet.TokenName,
				possibleWinAmount,
				[]bn.Number{bet.Amount},
				[]types.Address{bet.Gambler},
				"",
			)
			// tmp for unlock amount
			_, ok := tmp[bet.TokenName]
			if ok == false {
				tmp[bet.TokenName] = bn.N(0)
			}
			tmp[bet.TokenName] = tmp[bet.TokenName].Add(possibleWinAmount)

			refunded += 1
			//save this round saler address
			_, ok = grandSaler[salerAddress]
			if !ok {
				grandSaler[salerAddress] = 0
			}
		}
		return forx.Continue
	})

	if refundCount > MAXREFUNDCOUNT || refundCount <= 0 {
		refundCount = MAXREFUNDCOUNT
	}
	// setup data
	strartIndex := roundInfo.SettledCount
	lastIndex := strartIndex + refundCount - 1
	if lastIndex > roundInfo.BetCount-1 {
		lastIndex = roundInfo.BetCount - 1
	}

	forx.Range(strartIndex, lastIndex, func(index int64) bool {
		bet := sg._betInfo(commitStr, strconv.Itoa(int(index)))
		if bet.Settled == false {
			possibleWinAmount := sg.Refund(bet, index, commitStr)

			salerAddress := bet.Saler
			//refund by directsale
			sg.directsale().Refund(
				bet.Saler,
				bet.TokenName,
				possibleWinAmount,
				[]bn.Number{bet.Amount},
				[]types.Address{bet.Gambler},
				"",
			)

			// tmp for unlock amount
			_, ok := tmp[bet.TokenName]
			if ok == false {
				tmp[bet.TokenName] = bn.N(0)
			}
			tmp[bet.TokenName] = tmp[bet.TokenName].Add(possibleWinAmount)

			refunded += 1

			_, ok = grandSaler[salerAddress]
			if !ok {
				grandSaler[salerAddress] = 0
			}
		}
		return forx.Continue
	})

	roundInfo.RefundCount = roundInfo.RefundCount + refunded
	roundInfo.SettledCount = lastIndex + 1
	sg._setRoundInfo(commitStr, roundInfo)

	//Transfer to other handling address
	forx.Range(grandSaler, func(address types.Address, _ int) bool {
		sg.PayGrandPrize(address)
		return forx.Continue
	})

	//fire event
	sg.emitRefundBet(commit, refunded, tmp, roundInfo.RefundCount == roundInfo.BetCount)
	return
}
