package everycolor

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

const (
	PERMILLI       = 1000   // permillage
	MAXSETTLECOUNT = 100    // max settle count per time
	MAXREFUNDCOUNT = 100    // max refund count per time
	FIVEBALLMODEL  = 100000 // get five winnumber
)

type RoundInfo struct {
	Commit           []byte               `json:"commit"`           // round random number hash
	TotalBuy         map[string]bn.Number `json:"totalBuy"`         // round total BCB for buy
	WinNumber        string               `json:"winNumber"`        // lottery number
	BetCount         int64                `json:"betCount"`         // round people's number
	SettledCount     int64                `json:"settledCount"`     // round settle count
	RefundCount      int64                `json:"refundCount"`      // round refund count
	FirstBlockHeight int64                `json:"firstBlockHeight"` // round init block
	Setting          *MapSetting          `json:"mapSettings"`      // round setting
}

type BetInfo struct {
	Gambler       types.Address `json:"gambler"`       // gambler address
	TokenName     string        `json:"tokenName"`     // bet tokenName
	Amount        bn.Number     `json:"amount"`        // bet amount
	BetData       []BetData     `json:"betData"`       // bet data
	WinAmount     bn.Number     `json:"winAmount"`     // win amount
	EachWinAmount []bn.Number   `json:"eachWinAmount"` // each bet win amount
	Settled       bool          `json:"settled"`       // bet settled
}

type MapSetting struct {
	Settings            map[string]Setting `json:"settings"`            //设置信息
	BetExpirationBlocks int64              `json:"betExpirationBlocks"` //开奖的区块间隔数（开奖的区块要小于，下注区块加上区块间隔）
}
type Setting struct {
	MaxProfit      bn.Number `json:"maxProfit"`      //最大收益<MaxAmount（单位cong）
	MaxLimit       bn.Number `json:"maxLimit"`       //最大下注金额（单位cong）
	MinLimit       bn.Number `json:"minLimit"`       //最小的投注金额（单位cong）
	FeeRatio       int64     `json:"feeRatio"`       //总的手续费
	FeeMiniNum     bn.Number `json:"feeMiniNum"`     //最小的手续费总额（单位cong）
	SendToCltRatio int64     `json:"sendToCltRatio"` //手续费中转给clt的比例
}

type RecvFeeInfo struct {
	RecvFeeRatio []int64         `json:"recvFeeRatio"` //减掉clt后，手续费分配比例
	RecvFeeAddr  []types.Address `json:"recvFeeAddr"`  //接收手续费的地址列表
}
type PlayerIndexes struct {
	BetIndexes []int64 `json:"betIndexes"` //player bet by index
}
