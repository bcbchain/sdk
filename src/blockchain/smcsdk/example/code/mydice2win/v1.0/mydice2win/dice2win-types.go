package mydice2win

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

const (
	maxModulo     = 100  // is a number of equiprobable outcomes in a game
	maxMaskModulo = 40   // max mask modulo
	perMille      = 1000 // base of permille
)

// Bet - bet information
type Bet struct {
	TokenName        string        `json:"tokenName"`        // Token name
	Amount           bn.Number     `json:"amount"`           // Wager amount in cong.
	Modulo           int64         `json:"modulo"`           // Modulo of a game.
	RollUnder        bn.Number     `json:"rollUnder"`        // Number of winning outcomes, used to compute winning payment (* modulo/rollUnder), and used instead of mask for games with modulo > MAX_MASK_MODULO.
	Mask             bn.Number     `json:"mask"`             // Bit mask representing winning bet outcomes (see MAX_MASK_MODULO comment).
	PlaceBlockNumber int64         `json:"placeBlockNumber"` // Block number of placeBet tx.
	Gambler          types.Address `json:"gambler"`          // Address of a gambler, used to pay out winning bets.
}

func (b *Bet) init() {
	b.TokenName = ""
	b.Amount = bn.N(0)
	b.Gambler = ""
	b.Mask = bn.N(0)
	b.Modulo = 0
	b.PlaceBlockNumber = 0
	b.RollUnder = bn.N(0)
}

// Settings - contract settings
type Settings struct {
	TokenNames          map[string]struct{} `json:"tokenNames"`          //支持的代币名称列表
	MinBet              int64               `json:"minBet"`              //最小的投注金额（单位cong）
	MaxBet              int64               `json:"maxBet"`              //最大的投注金额（单位cong）
	MaxProfit           int64               `json:"maxProfit"`           //最大收益<MaxAmount（单位cong）
	FeeRatio            int64               `json:"feeRatio"`            //总的手续费
	FeeMinimum          int64               `json:"feeMinimum"`          //最小的手续费总额（单位cong）
	SendToCltRatio      int64               `json:"sendToCltRatio"`      //手续费中转给clt的比例
	BetExpirationBlocks int64               `json:"betExpirationBlocks"` //开奖的区块间隔数（开奖的区块要小于，下注区块加上区块间隔）
}

// RecvFeeInfo - recv fee info
type RecvFeeInfo struct {
	Ratio   int64         `json:"ratio"`   //减掉clt后，手续费分配比例
	Address types.Address `json:"address"` //接收手续费的地址列表
}
