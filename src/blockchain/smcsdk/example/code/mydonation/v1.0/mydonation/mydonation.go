package mydonation

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/std"
	"blockchain/smcsdk/sdk/types"
)

//Mydonation This is struct of contract
//@:contract:mydonation
//@:version:1.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:b37e7627431feb18123b81bcf1f41ffd37efdb90513d48ff2c7f8a0c27a9d06c
type Mydonation struct {
	sdk sdk.ISmartContract

	//Total donations received by donees
	//@:public:store
	donations map[types.Address]bn.Number // key=address of donee
}

const (
	errDoneeCannotBeOwner = 55000 + iota
	errDoneeCannotBeSmc
	errDoneeAlreadyExist
	errDoneeNotExist
	errDonationExist
	errDonationNotEnough
)

//@:public:receipt
type receipt interface {
	emitAddDonee(donee types.Address)
	emitDelDonee(donee types.Address)
	emitDonate(from, donee types.Address, value, balance bn.Number)
	emitTransferDonation(donee types.Address, value, balance bn.Number)
}

//InitChain Constructor of this Mydonation
//@:constructor
func (d *Mydonation) InitChain() {
}

//AddDonee Add a new donee
//@:public:method:gas[500]
func (d *Mydonation) AddDonee(donee types.Address) {
	sdk.RequireOwner()
	sdk.RequireAddress(donee)
	sdk.Require(donee != d.sdk.Message().Sender().Address(),
		errDoneeCannotBeOwner, "Donee can not be owner")
	sdk.Require(donee != d.sdk.Message().Contract().Address(),
		errDoneeCannotBeSmc, "Donee can not be this smart contract")
	sdk.Require(donee != d.sdk.Message().Contract().Account(),
		errDoneeCannotBeSmc, "Donee can not be account of this smart contract")
	sdk.Require(!d._chkDonations(donee),
		errDoneeAlreadyExist, "Donee already exists")

	d._setDonations(donee, bn.N(0))

	//emit receipt
	d.emitAddDonee(donee)
}

//Donate delete a donee
//@:public:method:gas[500]
func (d *Mydonation) DelDonee(donee types.Address) {
	sdk.RequireOwner()
	sdk.RequireAddress(donee)
	sdk.Require(d._chkDonations(donee),
		errDoneeNotExist, "Donee does not exist")
	sdk.Require(d._donations(donee).IsEqualI(0),
		errDonationExist, "Donation exists")

	d._delDonations(donee)

	//emit receipt
	d.emitDelDonee(donee)
}

//Donate Charitable donors donate money to smart contract
//@:public:method:gas[500]
func (d *Mydonation) Donate(donee types.Address) {
	sdk.RequireAddress(donee)
	sdk.Require(d._chkDonations(donee),
		errDoneeNotExist, "Donee does not exist")

	var valTome *std.Transfer
	token := d.sdk.Helper().GenesisHelper().Token()
	forx.Range(d.sdk.Message().GetTransferToMe(), func(i int, receipt *std.Transfer) bool {
		sdk.Require(receipt.Token == token.Address(),
			types.ErrInvalidParameter, "Accept donations in genesis token only")
		sdk.Require(valTome == nil,
			types.ErrInvalidParameter, "Accept only one donation at a time")
		valTome = receipt

		return true
	})
	sdk.Require(valTome != nil,
		types.ErrInvalidParameter, "Please transfer token to me first")

	balance := d._donations(donee).Add(valTome.Value)
	d._setDonations(donee, balance)

	//emit receipt
	d.emitDonate(
		d.sdk.Message().Sender().Address(),
		donee,
		valTome.Value,
		balance,
	)
}

//Withdraw To transfer donations to donee
//@:public:method:gas[500]
func (d *Mydonation) Transfer(donee types.Address, value bn.Number) {
	sdk.RequireOwner()
	sdk.RequireAddress(donee)
	sdk.Require(value.IsGreaterThanI(0),
		types.ErrInvalidParameter, "Parameter \"value\" must be greater than 0")
	sdk.Require(d._chkDonations(donee),
		errDoneeNotExist, "Donee does not exist")
	sdk.Require(d._donations(donee).IsGE(value),
		errDonationNotEnough, "Donation is not enough")

	account := d.sdk.Helper().AccountHelper().AccountOf(d.sdk.Message().Contract().Account())
	token := d.sdk.Helper().GenesisHelper().Token()
	account.TransferByToken(token.Address(), donee, value)
	balance := d._donations(donee).Sub(value)
	d._setDonations(donee, balance)

	//emit receipt
	d.emitTransferDonation(
		donee,
		value,
		balance,
	)
}
