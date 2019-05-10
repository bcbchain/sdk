package everycolor

import (
	"blockchain/algorithm"
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/utest"
	"common/jsoniter"
	"common/wal"
	"encoding/hex"
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/go-crypto"
	"gopkg.in/check.v1"
	"math"
	"testing"
)

const (
	keystore  = ".keystore"
	ownerName = "local_owner"
	password  = "Wty@12345678"
)

var (
	cdc = amino.NewCodec()
)

func init() {
	crypto.RegisterAmino(cdc)
	crypto.SetChainId("local")
	wal.NewAccount(keystore, ownerName, password)
}

//Test This is a function
func Test(t *testing.T) { check.TestingT(t) }

//MySuite This is a struct
type MySuite struct{}

var _ = check.Suite(&MySuite{})

//TestEverycolor_SetSecretSigner This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_SetSecretSigner(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	acct, _ := wal.LoadAccount(keystore, ownerName, password)
	pbk := acct.PubKey().(crypto.PubKeyEd25519)
	pubKey := pbk[:]

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

//TestEverycolor_SetSettings This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_SetSettings(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	genesisTokenName := test.obj.sdk.Helper().GenesisHelper().Token().Name()
	accounts := utest.NewAccounts(genesisTokenName, bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	mapSettings := MapSetting{Settings: map[string]Setting{}}
	set := make(map[string]Setting, 0)
	settings := set[genesisTokenName]
	settings.FeeRatio = 50
	settings.FeeMiniNum = bn.N(300000)
	settings.SendToCltRatio = 100
	settings.MaxProfit = bn.N(2E12)
	settings.MaxLimit = bn.N(2E10)
	settings.MinLimit = bn.N(1E8)
	mapSettings.BetExpirationBlocks = 250
	mapSettings.Settings[genesisTokenName] = settings
	resBytes1, _ := jsoniter.Marshal(mapSettings)

	settings.MaxProfit = bn.N(2E12)
	settings.MaxLimit = bn.N(2E9)
	settings.MinLimit = bn.N(1E10)
	mapSettings.Settings[genesisTokenName] = settings
	resBytes2, _ := jsoniter.Marshal(mapSettings)
	// 代币名称为空
	settings = mapSettings.Settings[""]
	settings.FeeRatio = 50
	settings.FeeMiniNum = bn.N(300000)
	settings.SendToCltRatio = 100
	settings.MaxProfit = bn.N(2E12)
	settings.MaxLimit = bn.N(2E10)
	settings.MinLimit = bn.N(1E8)
	mapSettings.BetExpirationBlocks = 250
	mapSettings.Settings[""] = settings
	resBytes3, _ := jsoniter.Marshal(mapSettings)

	settings = mapSettings.Settings[genesisTokenName]
	delete(mapSettings.Settings, "")
	settings.FeeRatio = 50
	settings.FeeMiniNum = bn.N(300000)
	settings.SendToCltRatio = 100
	settings.MaxProfit = bn.N(2E12)
	settings.MaxLimit = bn.N(0)
	mapSettings.BetExpirationBlocks = 250
	mapSettings.Settings[genesisTokenName] = settings
	resBytes4, _ := jsoniter.Marshal(mapSettings)

	settings.MaxLimit = bn.N(2E10)
	settings.MinLimit = bn.N(-1)
	mapSettings.Settings[genesisTokenName] = settings
	resBytes5, _ := jsoniter.Marshal(mapSettings)

	settings.MaxProfit = bn.N(math.MinInt64)
	settings.MaxLimit = bn.N(2E18)
	settings.MinLimit = bn.N(2E8)
	mapSettings.Settings[genesisTokenName] = settings
	resBytes6, _ := jsoniter.Marshal(mapSettings)

	settings.MaxProfit = bn.N(2E12)
	settings.MaxLimit = bn.N(2E11)
	settings.MinLimit = bn.N(2E10)
	settings.FeeMiniNum = bn.N(-1)
	mapSettings.Settings[genesisTokenName] = settings
	resBytes7, _ := jsoniter.Marshal(mapSettings)

	settings.FeeMiniNum = bn.N(300000)
	settings.FeeRatio = -1
	mapSettings.Settings[genesisTokenName] = settings
	resBytes8, _ := jsoniter.Marshal(mapSettings)

	settings.FeeRatio = 1001
	resBytes9, _ := jsoniter.Marshal(mapSettings)

	settings.FeeRatio = 50
	settings.SendToCltRatio = -1
	mapSettings.Settings[genesisTokenName] = settings
	resBytes10, _ := jsoniter.Marshal(mapSettings)

	settings.SendToCltRatio = 1001
	mapSettings.Settings[genesisTokenName] = settings
	resBytes11, _ := jsoniter.Marshal(mapSettings)

	settings.SendToCltRatio = 100
	mapSettings.BetExpirationBlocks = -1
	mapSettings.Settings[genesisTokenName] = settings
	resBytes12, _ := jsoniter.Marshal(mapSettings)

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
	for i, item := range tests {
		fmt.Println("12323123123", i)
		utest.AssertError(test.run().setSender(item.account).SetSettings(string(item.settings)), item.code)

	}
}

//TestEverycolor_SetRecvFeeInfo This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_SetRecvFeeInfo(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}
	ratio1 := 500
	address1 := "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE"
	ratio2 := 501
	ratio3 := 450
	address3 := "lo9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE"

	address4 := test.obj.sdk.Helper().BlockChainHelper().CalcAccountFromName(contractName, orgID)
	ratio4 := -1

	recvFeeInfo := new(RecvFeeInfo)
	recvFeeInfo.RecvFeeAddr = make([]string, 0)
	recvFeeInfo.RecvFeeRatio = make([]int64, 0)
	recvFeeInfo.RecvFeeAddr = []string{address1}
	recvFeeInfo.RecvFeeRatio = []int64{int64(ratio1)}
	resBytes1, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo.RecvFeeAddr = make([]string, 0)
	recvFeeInfo.RecvFeeRatio = make([]int64, 0)
	resBytes2, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo.RecvFeeAddr = []string{address1, address1}
	recvFeeInfo.RecvFeeRatio = []int64{int64(ratio1), int64(ratio2)}
	resBytes3, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo.RecvFeeAddr = []string{address1, address3}
	recvFeeInfo.RecvFeeRatio = []int64{int64(ratio1), int64(ratio2)}
	resBytes4, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo.RecvFeeAddr = []string{address1, address4}
	recvFeeInfo.RecvFeeRatio = []int64{int64(ratio1), int64(ratio3)}
	resBytes5, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo.RecvFeeAddr = []string{address1, address3}
	recvFeeInfo.RecvFeeRatio = []int64{int64(ratio1), int64(ratio4)}
	resBytes6, _ := jsoniter.Marshal(recvFeeInfo)

	recvFeeInfo.RecvFeeAddr = []string{address1}
	recvFeeInfo.RecvFeeRatio = []int64{int64(ratio1)}
	resBytes1, _ = jsoniter.Marshal(recvFeeInfo)

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
		utest.AssertError(test.run().setSender(item.account).SetRecvFeeInfo(string(item.infos)), item.code)
	}
}

//TestEverycolor_WithdrawFunds This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_WithdrawFunds(c *check.C) {
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

//TestEverycolor_PlaceBet This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_PlaceBet(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	contract := utest.UTP.Message().Contract()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 5)
	if accounts == nil {
		panic("初始化newOwner失败")
	}
	genesisTokenName := test.obj.sdk.Helper().GenesisHelper().Token().Name()
	commitLastBlock, pubKey, _, commit, signData := PlaceBetHelper(100)
	//utest.AssertError(err, types.CodeOK)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)
	//var i int64 = 0
	//for {
	//	if i > 2 {
	//		return
	//	}
	//	betData := []BetData{{i + 1, bn.N(1000000000)}}
	//	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	//													                                 // PlaceBet(betInfoJson string, commitLastBlock int64,betIndex string, commit, signData []byte, refAddress types.Address)
	//	utest.AssertError(test.run().setSender(accounts[i]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock,"hello", commit, signData[:], ""), types.CodeOK)
	//	i++
	//}
	//var i int64 = 0
	//for {
	//	if i > 2 {
	//		return
	//	}
	betData := []BetData{{"A", "XXXX5", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"A", "XXXX5", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	// PlaceBet(betInfoJson string, commitLastBlock int64,betIndex string, commit, signData []byte, refAddress types.Address)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	//	i++
	//}
}

//TestEverycolor_SettleBets This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_SettleBets(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)
	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), test.obj.sdk.Message().Contract().Account(), bn.N(1E15))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E14), 5)
	if accounts == nil {
		panic("初始化newOwner失败")
	}
	genesisTokenName := test.obj.sdk.Helper().GenesisHelper().Token().Name()
	commitLastBlock, pubKey, reveal, commit, signData := PlaceBetHelper(100)
	//utest.AssertError(err, types.CodeOK)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)
	//var i int64 = 0
	//for {
	//	if i > 2 {
	//		break
	//	}
	//	betData := []BetData{{i + 1, bn.N(1000000000)}}
	//	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	//	//utest.AssertError(test.run().setSender(accounts[i]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	//	utest.AssertError(test.run().setSender(accounts[i]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock,"hello", commit, signData[:], ""), types.CodeOK)
	//	i++
	//
	//}
	betData := []BetData{{"A", "XXXX5", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", "AAAXX", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	betData2 := []BetData{{"H", "DXXX*", bn.N(1000000000)}}
	betDataJsonBytes2, _ := jsoniter.Marshal(betData2)
	betData3 := []BetData{{"F", "5**XX", bn.N(1000000000)}}
	betDataJsonBytes3, _ := jsoniter.Marshal(betData3)
	betData4 := []BetData{{"D", "123XX", bn.N(1000000000)}}
	betDataJsonBytes4, _ := jsoniter.Marshal(betData4)
	// PlaceBet(betInfoJson string, commitLastBlock int64,betIndex string, commit, signData []byte, refAddress types.Address)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBets(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes4), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	//	i++
	utest.AssertError(test.run().setSender(contractOwner).SettleBets(reveal, 2), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBets(reveal, 4), types.CodeOK)

}

//TestEverycolor_WithdrawWin This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_WithdrawWin(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)
	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), test.obj.sdk.Message().Contract().Account(), bn.N(1E14))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 5)
	if accounts == nil {
		panic("初始化newOwner失败")
	}
	//commitLastBlock, pubKey, reveal, commit, signData, _ := PlaceBetHelper(100)
	commitLastBlock, pubKey, reveal, commit, signData := PlaceBetHelper(100)
	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)
	betData := []BetData{{"A", "XXXX5", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", "AAAXX", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	betData2 := []BetData{{"H", "DXXX*", bn.N(1000000000)}}
	betDataJsonBytes2, _ := jsoniter.Marshal(betData2)
	betData3 := []BetData{{"F", "5**XX", bn.N(1000000000)}}
	betDataJsonBytes3, _ := jsoniter.Marshal(betData3)
	betData4 := []BetData{{"D", "123XX", bn.N(1000000000)}}
	betDataJsonBytes4, _ := jsoniter.Marshal(betData4)

	genesisTokenName := test.obj.sdk.Helper().GenesisHelper().Token().Name()
	// PlaceBet(betInfoJson string, commitLastBlock int64,betIndex string, commit, signData []byte, refAddress types.Address)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes4), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	//utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	//utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 3), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBets(reveal, 1), types.CodeOK)
	//utest.AssertError(test.run().setSender(contractOwner).WithdrawWin(commit), types.ErrInvalidParameter)
	utest.AssertError(test.run().setSender(accounts[3]).WithdrawWin(commit), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBets(reveal, -1), types.CodeOK)
}

//TestEverycolor_RefundBets This is a method of MySuite
func (mysuit *MySuite) TestEverycolor_RefundBets(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)
	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), test.obj.sdk.Message().Contract().Account(), bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 2)
	if accounts == nil {
		panic("初始化newOwner失败")
	}
	genesisTokenName := test.obj.sdk.Helper().GenesisHelper().Token().Name()
	commitLastBlock, pubKey, _, commit, signData := PlaceBetHelper(100)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	betData := []BetData{{"A", "XXXX5", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(genesisTokenName, bn.N(1000000000), string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	// set bet time out
	count := 0
	for {
		utest.NextBlock(1)
		count++
		if count > 250 {
			break
		}
	}
	utest.AssertError(test.run().setSender(contractOwner).RefundBets(commit, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).RefundBets(commit, 1), types.CodeOK)
}

//hempHeight 想对于下注高度和生效高度之间的差值
//acct 合约的owner
func PlaceBetHelper(tempHeight int64) (commitLastBlock int64, pubKey [32]byte, reveal, commit []byte, signData [64]byte) {
	acct, _ := wal.LoadAccount(".keystore", "local_owner", password)

	localBlockHeight := utest.UTP.ISmartContract.Block().Height()

	pubKey = acct.PubKey().(crypto.PubKeyEd25519)

	commitLastBlock = localBlockHeight + tempHeight
	decode := crypto.CRandBytes(32)
	revealStr := hex.EncodeToString(algorithm.SHA3256(decode))
	reveal, _ = hex.DecodeString(revealStr)

	commit = algorithm.SHA3256(reveal)

	signByte := append(bn.N(commitLastBlock).Bytes(), commit...)
	signData = acct.PrivateKey.Sign(signByte).(crypto.SignatureEd25519)

	return
}
