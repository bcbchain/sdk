package std

import "github.com/AeReach/sdk/sdk/types"

type Message struct {
	SmcAddress types.Address `json:"smcAddress"`
	Method     string        `json:"method"`
	To         string        `json:"to"`
	Value      string        `json:"value"`
}

type TxResult struct {
	TxHash      string
	Code        uint32
	Log         string
	BlockHeight int64
	From        string
	Note        string
	Message     []Message
}
