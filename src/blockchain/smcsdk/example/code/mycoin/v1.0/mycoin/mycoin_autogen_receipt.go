package mycoin

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

var _ receipt = (*Mycoin)(nil)

//emitTransferMyCoin This is a method of Mycoin
func (m *Mycoin) emitTransferMyCoin(token, from, to types.Address, value bn.Number) {
	type transferMyCoin struct {
		Token types.Address `json:"token"`
		From  types.Address `json:"from"`
		To    types.Address `json:"to"`
		Value bn.Number     `json:"value"`
	}

	m.sdk.Helper().ReceiptHelper().Emit(
		transferMyCoin{
			Token: token,
			From:  from,
			To:    to,
			Value: value,
		},
	)
}
