package utest

import (
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
	"time"

	"github.com/tendermint/go-crypto"
)

const genesisStr = `{
  "app_hash": "",
  "app_state": {
    "token": {
      "address": "localAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1TA",
      "owner": "localBsHvWxKkScTSpkF5gPFhrWegN2yosrZV9",
      "version": "",
      "name": "LOC",
      "symbol": "LOC",
      "totalSupply": {"v": 2000000000000000000},
      "addSupplyEnabled": false,
      "burnEnabled": false,
      "gasPrice": 2500
    },
    "contracts": [
      {
        "name": "system",
        "address": "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
        "owner": "localBsHvWxKkScTSpkF5gPFhrWegN2yosrZV9",
        "version": "",
        "codeHash": "BE6DC61A323C4A027337F181CEBB0BE105C5530E974ED7573FA86DAC94FDC347",
        "effectHeight": 1,
        "loseHeight": 0,
        "token": "",
        "orgID": "orgAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1T",
        "methods": [
          {
            "methodId": "1b0b7ce1",
            "prototype": "NewValidator(string,smc.PubKey,smc.Address,uint64)smc.Error",
            "gas": 50000
          },
          {
            "methodId": "bcc890f0",
            "prototype": "SetPower(smc.PubKey,uint64)smc.Error",
            "gas": 20000
          },
          {
            "methodId": "757a5721",
            "prototype": "SetRewardAddr(smc.PubKey,smc.Address)smc.Error",
            "gas": 20000
          },
          {
            "methodId": "52d1c625",
            "prototype": "ForbidInternalContract(smc.Address,uint64)smc.Error",
            "gas": 50000
          },
          {
            "methodId": "48c26f01",
            "prototype": "DeployInternalContract(string,string,[]string,[]uint64,smc.Hash,uint64)(smc.Address,smc.Error)",
            "gas": 50000
          },
          {
            "methodId": "92622878",
            "prototype": "SetRewardStrategy(string,uint64)smc.Error",
            "gas": 50000
          }
        ]
      },
      {
        "name": "token-basic",
        "address": "localAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1TA",
        "owner": "localBsHvWxKkScTSpkF5gPFhrWegN2yosrZV9",
        "version": "",
        "codeHash": "BF463D47763A9430DFD5AE2FBCF02BEABE44DD2C40BA1FBECEC9CF5A4EFA3FE5",
        "effectHeight": 1,
        "loseHeight": 0,
        "token": "localAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1TA",
        "orgID": "orgAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1T",
        "methods": [
          {
            "methodId": "af0228bc",
            "prototype": "Transfer(smc.Address,big.Int)smc.Error",
            "gas": 500
          },
          {
            "methodId": "cc50053d",
            "prototype": "SetGasPrice(uint64)smc.Error",
            "gas": 2000
          },
          {
            "methodId": "adcf5578",
            "prototype": "SetGasBasePrice(uint64)smc.Error",
            "gas": 2000
          }
        ]
      },
      {
        "name": "token-issue",
        "address": "local5sAkG9xgzjuuhQn3FVf8K3FCrk5RAtab6",
        "owner": "localBsHvWxKkScTSpkF5gPFhrWegN2yosrZV9",
        "version": "1.0",
        "codeHash": "22D14A57C3D0F3388E93501E3D4636E38D32FA4EBBA054BA7674CAC9B6EB0ED6",
        "effectHeight": 1,
        "loseHeight": 0,
        "token": "",
        "orgID": "orgAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1T",
        "methods": [
          {
            "methodId": "825b55bb",
            "prototype": "NewToken(string,string,big.Int,bool,bool,uint64)(smc.Address,smc.Error)",
            "gas": 20000
          },
          {
            "methodId": "af0228bc",
            "prototype": "Transfer(smc.Address,big.Int)smc.Error",
            "gas": 600
          },
          {
            "methodId": "b8f4b0d8",
            "prototype": "BatchTransfer([]smc.Address,big.Int)smc.Error",
            "gas": 6000
          },
          {
            "methodId": "4a439c44",
            "prototype": "AddSupply(big.Int)smc.Error",
            "gas": 2400
          },
          {
            "methodId": "5a0e0fa3",
            "prototype": "Burn(big.Int)smc.Error",
            "gas": 2400
          },
          {
            "methodId": "a6be2c35",
            "prototype": "SetOwner(smc.Address)smc.Error",
            "gas": 2400
          },
          {
            "methodId": "cc50053d",
            "prototype": "SetGasPrice(uint64)smc.Error",
            "gas": 2400
          }
        ]
      },
      {
        "name": "myplayerbook",
        "address": "localWkNWzXyqMmumfxfXva2QV1qKa3aroVyu",
        "account": "local2uGLmMnsHauRUjyjQKGdXchUxpRMM8oeD",
        "owner": "localBsHvWxKkScTSpkF5gPFhrWegN2yosrZV9",
        "version": "2.0",
        "codeHash": "43A15EC506F3864126E78FD3E1A265D9EAF5D436E776C9B0200D77E57B76B7ED",
        "effectHeight": 1,
        "loseHeight": 0,
        "token": "",
        "orgID": "orgAJrbk6Wdf7TCbunrXXS5kKvbWVszhC1T",
        "methods": [
          {
            "methodId": "e463fdb2",
            "prototype": "RegisterName(string)(types.Error)",
            "gas": 500
          }
        ]
      }
    ]
  },
  "chain_id": "local",
  "genesis_time": "2018-07-14T14:43:53.1778867+08:00",
  "validators": [
    {
      "name": "local",
      "reward_addr": "local8LtT8AonWgJ8nMCEdAR5UGrbRfUmuoeiz",
      "power": 10
    }
  ],
  "mode": 1
}
`

// GenesisValidator validator's info
type GenesisValidator struct {
	RewardAddr string        `json:"reward_addr"`
	PubKey     crypto.PubKey `json:"pub_key,omitempty"` // No Key In Genesis File,so omit empty
	Power      int64         `json:"power"`
	Name       string        `json:"name"`
}

// TokenContract token contract
type TokenContract struct {
	Address      types.Address `json:"address"`      // 合约地址
	EffectHeight int64         `json:"effectHeight"` // 合约生效的区块高度
}

// AppState app state
type AppState struct {
	GnsToken     std.Token      `json:"token"`
	GnsContracts []std.Contract `json:"contracts"`
}

type genesis struct {
	GenesisTime  time.Time          `json:"genesis_time"`
	ChainID      string             `json:"chain_id"`
	Validators   []GenesisValidator `json:"validators"`
	AppHash      types.HexBytes     `json:"app_hash"`
	AppStateJSON AppState           `json:"app_state"`
	Mode         int                `json:"mode,omitempty"`
}
