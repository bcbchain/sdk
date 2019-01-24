package utest

import (
	"blockchain/smcsdk/sdk/crypto/sha3"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl"
	"blockchain/smcsdk/sdkimpl/object"
	"fmt"
	"math/big"
	"time"
)

var (
	//BlockHeight block height
	BlockHeight int64
	//LastNumTxs number of txs
	LastNumTxs int32

	tx      = []types.HexBytes{types.HexBytes("YmNiPHR4Pi52MS4yVVpHSlJRZDYxYTdMa296MzJuRWJQQU5RQVJhWG9jZWQ1dDR5THRvZkVtaENWMkc3eWVrNVgzdU1UUzdlazU5ODZNTExRZ3ZRRllEb3ZQZjYyb3RjUjNLQ3p0VU5tcE1xU1l1SE5EREhHeGVCdUZZb0xyRk5LU3cxdEFZb2t1NGt4RlouPDE+LllUZ2lBMWdkREdpMkw4aG44enhlNWp5d2Y0bTFvZ3o4OXFRd0R1Y0ZCcERhQjZSOFkzOW5MUERyZ0FOZUxYVzNmZGdNd2o4WWFBUkRmTlZLTWlyTGRKQ1FLMkZham1ScFl4OHRZdEx0ZWlKdUhyUlgzakc3ZU1Nc01FVVdXeEdDaUxCNG5DMXkxVUs1NTk4ODg4WUEy")}
	block   = []byte(`{"chainID":"local","blockHash":"663666426932665034586B70685A5A323757424F786845473079493D","height":1,"time":1542436677,"numTxs":1,"dataHash":"4369496472305A546757635671336C3131326D31773539684557593D","proposerAddress":"localCUh7Zsb7PBgLwHJVok2QaMhbW64HNK4FU","rewardAddress":"localCUh7Zsb7PBgLwHJVok2QaMhbW64HNK4FU","randomNumber":"596D4E694E3074346546704E64314E3057556448626B74704D33684F546D3035625735744E33426B563352474F576445","version":"1.0","lastBlockHash":"4369496472305A546757635671336C3131326D31773539684557593D","lastCommitHash":"4D31724C4573745A68314B30624A644D6B4A4B53625774324450673D","lastAppHash":"36324446453643413937353539313437324435323143413943303845413243453242444646333546363231363230393132393937324636414133334633314234","lastFee":1500000}`)
	stateDB map[string][]byte
)

func initStateDB() {
	BlockHeight = 0
	LastNumTxs = 0

	stateDB = make(map[string][]byte)
	stateDB[std.KeyOfAppState()] = block
}

func setToDB(key string, value []byte) {
	stateDB[key] = value
}

func getBlock(height int64) []byte {
	return stateDB[std.KeyOfAppState()]
}

func build(transID, txID int64, meta std.ContractMeta) (result std.BuildResult) {
	return
}

func sdbGet(transID, txID int64, key string) []byte {
	result := std.GetResult{Msg: "ok"}
	if data, ok := stateDB[key]; ok == true {
		result.Code = types.CodeOK
		result.Data = data
	} else {
		result.Code = types.ErrInvalidParameter
		result.Msg = "invalid key"
	}

	resBytes, _ := jsoniter.Marshal(result)

	return resBytes
}

func sdbSet(transID, txID int64, values map[string][]byte) {
	for k, v := range values {
		stateDB[k] = v
	}
}

// SdbGet 供合约运行服务的测试程序使用
func SdbGet(transID, txID int64, key string) []byte {
	return stateDB[key]
}

//SdbSet set sdb
func SdbSet(transID, txID int64, values map[string][]byte) {
	for k, v := range values {
		stateDB[k] = v
	}
}

//GetBlock get block data
func GetBlock(numTxs int32) []byte {
	BlockHeight++

	block := std.Block{
		ChainID:         utChainID,
		BlockHash:       sha3.Sum256(big.NewInt(BlockHeight).Bytes()),
		Height:          BlockHeight,
		Time:            time.Now().Unix(),
		NumTxs:          numTxs,
		DataHash:        sha3.Sum256(big.NewInt(BlockHeight + 150000000000000).Bytes()),
		ProposerAddress: CalcAccountFromPubKey([]byte("pp123456789012345678901234567890")),
		RewardAddress:   CalcAccountFromPubKey([]byte("rw123456789012345678901234567890")),
		RandomNumber:    sha3.Sum256(big.NewInt(BlockHeight + 983377333372898).Bytes()),
		Version:         "",
		LastBlockHash:   sha3.Sum256(big.NewInt(BlockHeight - 1).Bytes()),
		LastCommitHash:  sha3.Sum256(big.NewInt(BlockHeight + 100000000000000).Bytes()),
		LastAppHash:     sha3.Sum256(big.NewInt(BlockHeight + 200000000000000).Bytes()),
		LastFee:         int64(LastNumTxs) * 500 * 2500}

	resBytes, err := jsoniter.Marshal(block)
	if err != nil {
		panic(err.Error())
	}
	LastNumTxs = numTxs

	return resBytes
}

//NextBlock generate next block data
func NextBlock(_numTxs int32) []byte {

	BlockHeight++
	block := object.NewBlock(UTP.ISmartContract.(*sdkimpl.SmartContract),
		utChainID,
		"",
		sha3.Sum256(big.NewInt(BlockHeight).Bytes()),
		sha3.Sum256(big.NewInt(BlockHeight+150000000000000).Bytes()),
		BlockHeight,
		time.Now().Unix(),
		_numTxs,
		CalcAccountFromPubKey([]byte("pp123456789012345678901234567890")),
		CalcAccountFromPubKey([]byte("rw123456789012345678901234567890")),
		sha3.Sum256(big.NewInt(BlockHeight+983377333372898).Bytes()),
		sha3.Sum256(big.NewInt(BlockHeight-1).Bytes()),
		sha3.Sum256(big.NewInt(BlockHeight+100000000000000).Bytes()),
		sha3.Sum256(big.NewInt(BlockHeight+200000000000000).Bytes()),
		int64(LastNumTxs)*500*2500)

	resBytes, err := jsoniter.Marshal(block)
	if err != nil {
		panic(err.Error())
	}
	LastNumTxs = _numTxs

	smc := UTP.ISmartContract.(*sdkimpl.SmartContract)
	smc.SetBlock(block)

	key := fmt.Sprintf("/block/%d", block.Height())
	b := make(map[string][]byte)

	b[key], err = jsoniter.Marshal(block)
	if err != nil {
		panic(err.Error())
	}
	sdbSet(0, 0, b)
	return resBytes
}

func data(key string, resBytes []byte) []byte {
	var getResult std.GetResult
	err := jsoniter.Unmarshal(resBytes, &getResult)
	if err != nil {
		sdkimpl.Logger.Fatalf("Cannot unmarshal get result struct, key=%s, error=%v\nbytes=%v", key, err, resBytes)
		sdkimpl.Logger.Flush()
		panic(err)
	} else if getResult.Code != types.CodeOK {
		sdkimpl.Logger.Debugf("Cannot find key=%s in stateDB, error=%s", getResult.Msg)
		return nil
	}

	return getResult.Data
}
