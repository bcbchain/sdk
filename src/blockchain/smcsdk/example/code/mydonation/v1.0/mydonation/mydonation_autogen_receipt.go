package mydonation

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

var _ receipt = (*Mydonation)(nil)

//emitAddDonee This is a method of Mydonation
func (m *Mydonation) emitAddDonee(donee types.Address) {
	type addDonee struct {
		Donee types.Address `json:"donee"`
	}

	m.sdk.Helper().ReceiptHelper().Emit(
		addDonee{
			Donee: donee,
		},
	)
}

//emitDelDonee This is a method of Mydonation
func (m *Mydonation) emitDelDonee(donee types.Address) {
	type delDonee struct {
		Donee types.Address `json:"donee"`
	}

	m.sdk.Helper().ReceiptHelper().Emit(
		delDonee{
			Donee: donee,
		},
	)
}

//emitDonate This is a method of Mydonation
func (m *Mydonation) emitDonate(from, donee types.Address, value, balance bn.Number) {
	type donate struct {
		From    types.Address `json:"from"`
		Donee   types.Address `json:"donee"`
		Value   bn.Number     `json:"value"`
		Balance bn.Number     `json:"balance"`
	}

	m.sdk.Helper().ReceiptHelper().Emit(
		donate{
			From:    from,
			Donee:   donee,
			Value:   value,
			Balance: balance,
		},
	)
}

//emitTransferDonation This is a method of Mydonation
func (m *Mydonation) emitTransferDonation(donee types.Address, value, balance bn.Number) {
	type transferDonation struct {
		Donee   types.Address `json:"donee"`
		Value   bn.Number     `json:"value"`
		Balance bn.Number     `json:"balance"`
	}

	m.sdk.Helper().ReceiptHelper().Emit(
		transferDonation{
			Donee:   donee,
			Value:   value,
			Balance: balance,
		},
	)
}
