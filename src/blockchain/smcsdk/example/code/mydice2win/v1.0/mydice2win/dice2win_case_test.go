package mydice2win

import (
	"blockchain/algorithm"
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/utest"
	"common/keys"
	"common/kms"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"testing"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/go-crypto"

	cmn "github.com/tendermint/tmlibs/common"
	"gopkg.in/check.v1"
)

const (
	ownerName = "local_owner"
	password  = "12345678"
)

var (
	cdc = amino.NewCodec()
)

func init() {
	crypto.RegisterAmino(cdc)
	crypto.SetChainId("local")
	kms.InitKMS("./.keystore", "local_mode", "", "", "0x1003")
	kms.GenPrivKey(ownerName, []byte(password))
}

// Hook up goCheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type MySuite struct{}

var _ = check.Suite(&MySuite{})

//TestDice2Win_SetSecretSigner is a method of MySuite
func (mysuit *MySuite) TestDice2Win_SetSecretSigner(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	pubKey, _ := kms.GetPubKey(ownerName, []byte(password))

	account := utest.NewAccount(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1000000000))

	var tests = []struct {
		account sdk.IAccount
		pubKey  []byte
		desc    string
		code    uint32
	}{
		{contractOwner, pubKey, "--正常流程--", types.CodeOK},
		{contractOwner, []byte("0xff"), "--异常流程--公钥长度不正确--", types.ErrInvalidParameter},
		{account, pubKey, "--异常流程--非owner调用--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).SetSecretSigner(item.pubKey), item.code)
	}
}

//TestDice2Win_SetSettings is a method of MySuite
func (mysuit *MySuite) TestDice2Win_SetSettings(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	settings := Settings{}
	settings.TokenNames = []string{test.obj.sdk.Helper().GenesisHelper().Token().Name()}
	settings.MaxBet = 2E10
	settings.MinBet = 1E8
	settings.MaxProfit = 2E12
	settings.FeeMinimum = 300000
	settings.FeeRatio = 50
	settings.SendToCltRatio = 100
	settings.BetExpirationBlocks = 250
	resBytes1, _ := jsoniter.Marshal(settings)

	settings.MaxBet = 2E9
	settings.MinBet = 2E10
	resBytes2, _ := jsoniter.Marshal(settings)

	settings.MaxBet = 2E10
	settings.MinBet = 2E8
	settings.TokenNames = []string{}
	resBytes3, _ := jsoniter.Marshal(settings)

	settings.TokenNames = []string{test.obj.sdk.Helper().GenesisHelper().Token().Name()}
	settings.MaxBet = 0
	resBytes4, _ := jsoniter.Marshal(settings)

	settings.MaxBet = 2E10
	settings.MinBet = -1
	resBytes5, _ := jsoniter.Marshal(settings)

	settings.MinBet = 2E8
	settings.MaxProfit = math.MinInt64
	resBytes6, _ := jsoniter.Marshal(settings)

	settings.MaxProfit = 2E12
	settings.FeeMinimum = -1
	resBytes7, _ := jsoniter.Marshal(settings)

	settings.FeeMinimum = 300000
	settings.FeeRatio = -1
	resBytes8, _ := jsoniter.Marshal(settings)

	settings.FeeRatio = 1001
	resBytes9, _ := jsoniter.Marshal(settings)

	settings.FeeRatio = 50
	settings.SendToCltRatio = -1
	resBytes10, _ := jsoniter.Marshal(settings)

	settings.SendToCltRatio = 1001
	resBytes11, _ := jsoniter.Marshal(settings)

	settings.SendToCltRatio = 100
	settings.BetExpirationBlocks = -1
	resBytes12, _ := jsoniter.Marshal(settings)

	var tests = []struct {
		account  sdk.IAccount
		settings []byte
		desc     string
		code     uint32
	}{
		{contractOwner, resBytes1, "--正常流程--", types.CodeOK},
		{contractOwner, resBytes2, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes3, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes4, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes5, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes6, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes7, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes8, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes9, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes10, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes11, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes12, "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], resBytes1, "--异常流程--", types.ErrNoAuthorization},
	}

	test.run().setSender(contractOwner).InitChain()
	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).SetSettings(string(item.settings)), item.code)
	}
}

//TestDice2Win_SetRecvFeeInfo is a method of MySuite
func (mysuit *MySuite) TestDice2Win_SetRecvFeeInfos(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	recvFeeInfo := make([]RecvFeeInfo, 0)
	resBytes2, _ := jsoniter.Marshal(recvFeeInfo)
	item := RecvFeeInfo{
		Ratio:   500,
		Address: "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recvFeeInfo = append(recvFeeInfo, item)
	resBytes1, _ := jsoniter.Marshal(recvFeeInfo)

	item1 := RecvFeeInfo{
		Ratio:   501,
		Address: "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recvFeeInfo = append(recvFeeInfo, item1)
	resBytes3, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo = append(recvFeeInfo[:1], recvFeeInfo[2:]...)
	item2 := RecvFeeInfo{
		Ratio:   450,
		Address: "lo9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recvFeeInfo = append(recvFeeInfo, item2)
	resBytes4, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo = append(recvFeeInfo[:1], recvFeeInfo[2:]...)
	item3 := RecvFeeInfo{
		Ratio:   500,
		Address: test.obj.sdk.Helper().BlockChainHelper().CalcAccountFromName(contractName),
	}
	recvFeeInfo = append(recvFeeInfo, item3)
	resBytes5, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo = append(recvFeeInfo[:1], recvFeeInfo[2:]...)
	item4 := RecvFeeInfo{
		Ratio:   -1,
		Address: "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recvFeeInfo = append(recvFeeInfo, item4)
	resBytes6, _ := jsoniter.Marshal(recvFeeInfo)

	var tests = []struct {
		account sdk.IAccount
		infos   []byte
		desc    string
		code    uint32
	}{
		{contractOwner, resBytes1, "--正常流程--", types.CodeOK},
		{contractOwner, resBytes2, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes3, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes4, "--异常流程--", types.ErrInvalidAddress},
		{contractOwner, resBytes5, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes6, "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], resBytes1, "--异常流程--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).SetRecvFeeInfos(string(item.infos)), item.code)
	}
}

//TestDice2Win_WithdrawFunds is a method of MySuite
func (mysuit *MySuite) TestDice2Win_WithdrawFunds(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisToken := test.obj.sdk.Helper().GenesisHelper().Token()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	contractAccount := utest.UTP.Helper().ContractHelper().ContractOfName(contractName).Account()

	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)
	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contractAccount, bn.N(1E11))

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	test.run().setSender(contractOwner).InitChain()

	var tests = []struct {
		account        sdk.IAccount
		tokenName      string
		beneficiary    types.Address
		withdrawAmount bn.Number
		desc           string
		code           uint32
	}{
		{contractOwner, genesisToken.Name(), contractOwner.Address(), bn.N(1E10), "--正常流程--", types.CodeOK},
		{contractOwner, genesisToken.Name(), accounts[0].Address(), bn.N(1E10), "--正常流程--", types.CodeOK},
		{contractOwner, genesisToken.Name(), contractOwner.Address(), bn.N(1E15), "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, genesisToken.Name(), contractOwner.Address(), bn.N(-1), "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, genesisToken.Name(), contractAccount, bn.N(1E10), "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, "xt", contractOwner.Address(), bn.N(1E10), "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], genesisToken.Name(), contractOwner.Address(), bn.N(1E10), "--异常流程--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).WithdrawFunds(item.tokenName, item.beneficiary, item.withdrawAmount), item.code)
	}
}

//TestDice2Win_PlaceBet is a method of MySuite
func (mysuit *MySuite) TestDice2Win_PlaceBet(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	contract := utest.UTP.Message().Contract()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), bn.N(1E11))

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock1, pubKey, _, commit1, signData1 := PlaceBetHelper(100)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)
	hexStr := hex.EncodeToString(pubKey[:])
	fmt.Println(hexStr)

	var tests = []struct {
		account         sdk.IAccount
		amount          bn.Number
		betMask         bn.Number
		modulo          int64
		commitLastBlock int64
		commit          []byte
		signData        []byte
		desc            string
		code            uint32
	}{
		{accounts[0], bn.N(1000000000), bn.N(1), 2, commitLastBlock1, commit1, signData1[:], "--正常流程--", types.CodeOK},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).transfer(item.amount).PlaceBet(item.betMask, item.modulo, item.commitLastBlock, item.commit, item.signData, ""), item.code)
	}
}

//TestDice2Win_SettleBet is a method of MySuite
func (mysuit *MySuite) TestDice2Win_SettleBet(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), test.obj.sdk.Message().Contract().Account(), bn.N(1E11))

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock1, pubKey, reveal, commit1, signData1 := PlaceBetHelper(100)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(bn.N(1), 2, commitLastBlock1, commit1, signData1[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal), types.CodeOK)
}

//TestDice2Win_RefundBet is a method of MySuite
func (mysuit *MySuite) TestDice2Win_RefundBet(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), test.obj.sdk.Message().Contract().Account(), bn.N(1E11))

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock1, pubKey, _, commit1, signData1 := PlaceBetHelper(100)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(bn.N(1), 2, commitLastBlock1, commit1, signData1[:], ""), types.CodeOK)
	// set bet time out
	count := 0
	for {
		utest.NextBlock(1)
		count++
		if count > 250 {
			break
		}
	}
	utest.AssertError(test.run().setSender(contractOwner).RefundBet(commit1), types.CodeOK)
}

//hempHeight 想对于下注高度和生效高度之间的差值
//acct 合约的owner
func PlaceBetHelper(tempHeight int64) (commitLastBlock int64, pubKey [32]byte, reveal, commit []byte, signData [64]byte) {
	acct := Load("./.keystore/local_owner.wal", []byte(password), nil)

	localBlockHeight := utest.UTP.ISmartContract.Block().Height()

	pubKey = acct.PubKey.(crypto.PubKeyEd25519)

	commitLastBlock = localBlockHeight + tempHeight
	decode := crypto.CRandBytes(32)
	revealStr := hex.EncodeToString(algorithm.SHA3256(decode))
	reveal, _ = hex.DecodeString(revealStr)

	commit = algorithm.SHA3256(reveal)

	signByte := append(bn.N(commitLastBlock).Bytes(), commit...)
	signData = acct.PrivKey.Sign(signByte).(crypto.SignatureEd25519)

	return
}

func Load(keystorePath string, password, fingerprint []byte) (acct *keys.Account) {
	if keystorePath == "" {
		cmn.PanicSanity("Cannot loads account because keystorePath not set")
	}

	walBytes, mErr := ioutil.ReadFile(keystorePath)
	sdk.RequireNotError(mErr, types.ErrInvalidParameter)

	jsonBytes, mErr := algorithm.DecryptWithPassword(walBytes, password, fingerprint)
	sdk.RequireNotError(mErr, types.ErrInvalidParameter)

	acct = new(keys.Account)
	mErr = cdc.UnmarshalJSON(jsonBytes, acct)
	sdk.RequireNotError(mErr, types.ErrInvalidParameter)

	acct.KeystorePath = keystorePath
	return
}
