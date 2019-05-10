package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
)

type MapSetting struct {
	Settings            map[string]Setting `json:"settings"`            //设置信息
	BetExpirationBlocks int64              `json:"betExpirationBlocks"` //开奖的区块间隔数（开奖的区块要小于，下注区块加上区块间隔）
	PoolFeeRatio        int64              `json:"poolFeeRatio"`        //从赢家扣除的比例，会积累到奖池
	CarveUpPoolRatio    int64              `json:"carveUpPoolRatio"`    //出现三公时，瓜分奖池的比例
	GrandPrizeRatio     int64              `json:"grandPrizeRatio"`     //24小时倒计时，赢得大奖的奖励比例
}
type Setting struct {
	MaxProfit      bn.Number `json:"maxProfit"`      //最大收益<MaxAmount（单位cong）
	MaxLimit       bn.Number `json:"maxLimit"`       //最大下注金额（单位cong）
	MinLimit       bn.Number `json:"minLimit"`       //最小的投注金额（单位cong）
	FeeRatio       int64     `json:"feeRatio"`       //总的手续费
	FeeMiniNum     bn.Number `json:"feeMiniNum"`     //最小的手续费总额（单位cong）
	SendToCltRatio int64     `json:"sendToCltRatio"` //手续费中转给clt的比例
}

type RecFeeInfo struct {
	RecFeeRatio int64         `json:"recFeeRatio"` // Commission allocation ratio
	RecFeeAddr  types.Address `json:"recFeeAddr"`  // List of addresses to receive commissions
}

type RoundInfo struct {
	Commit            []byte                                            `json:"commit"`            // round random number hash
	TotalBuy          map[types.Address]map[string]map[string]bn.Number `json:"totalBuy"`          // round total buy,key1=saler address ,key2 = playerID,key3=tokenName
	SettledPoolAmount map[types.Address]map[string]bn.Number            `json:"settledPoolAmount"` // round settledPoolAmount key1=saler, key2=tokenName
	Banker            GamerInfo                                         `json:"banker"`            // banker info
	Players           map[string]GamerInfo                              `json:"players"`           // Game Result Player Type
	BetCount          int64                                             `json:"betCount"`          // round people's number
	SettledCount      int64                                             `json:"settledCount"`      // round settle count
	RefundCount       int64                                             `json:"refundCount"`       // round refund count
	FirstBlockHeight  int64                                             `json:"firstBlockHeight"`  // round init block
	Setting           *MapSetting                                       `json:"mapSettings"`       // round setting
}

type BetData struct {
	PlayerId  string    `json:"playerId"`  // A, B, C, D
	BetAmount bn.Number `json:"betAmount"` // Betting amount
}

type BetInfo struct {
	BetTime       int64         `json:"betTime"`       // player’s bet time
	Gambler       types.Address `json:"gambler"`       // Player betting address
	Saler         types.Address `json:"saler"`         // Operator address
	TokenName     string        `json:"tokenName"`     // Players bet on currency names
	Amount        bn.Number     `json:"amount"`        // Players bet the total amount
	BetData       []BetData     `json:"betData"`       // Player betting details
	WinAmount     bn.Number     `json:"winAmount"`     // Players this bet the largest bonus
	EachWinAmount []bn.Number   `json:"eachWinAmount"` // each bet win amount
	Settled       bool          `json:"settled"`       //  the current bet has been settled
}

func NewBetInfo(tokenName string, gambler types.Address) *BetInfo {
	return &BetInfo{
		TokenName:     tokenName,
		Gambler:       gambler,
		Amount:        bn.N(0),
		BetData:       make([]BetData, 0),
		WinAmount:     bn.N(0),
		EachWinAmount: make([]bn.Number, 0),
		Settled:       false,
	}
}

func (bi *BetInfo) UpdateBetInfo(bet BetData) {
	data := bi.BetData
	bi.BetData = append(data, bet)
	number := bet.BetAmount
	bi.Amount = bi.Amount.Add(number)
}

type GrandPrizer struct {
	BetTime   int64         `json:"betTime"`   //Clearing time
	Gambler   types.Address `json:"gambler"`   //Players address
	WinAmount bn.Number     `json:"winAmount"` //The GrandPrize  amount
}

func BuildBetData(jsonData []byte) []BetData {
	data := make([]BetData, 0)
	unmarshal := jsoniter.Unmarshal(jsonData, &data)
	sdk.RequireNotError(unmarshal, types.ErrInvalidParameter)
	return data
}

func CheckList(t string, s []string) (flag bool) {
	for _, v := range s {
		if v == t {
			flag = true
			return
		}
	}
	return
}

type PlayerIndexes struct {
	BetIndexes []int64 `json:"betIndexes"` //player bet by index
}
