package mycoin

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

//Mycoin a demo contract for digital coin
//@:contract:mycoin
//@:version:1.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:b37e7627431feb18123b81bcf1f41ffd37efdb90513d48ff2c7f8a0c27a9d06c
type Mycoin struct {
	sdk sdk.ISmartContract

	//@:public:store
	totalSupply bn.Number

	//@:public:store
	balanceOf map[types.Address]bn.Number
}

const oneToken int64 = 1000000000

//InitChain init when deployed on the blockchain first time
//@:constructor
func (mc *Mycoin) InitChain() {
	owner := mc.sdk.Message().Contract().Owner()
	totalSupply := bn.N1(1000000, oneToken)
	mc._setTotalSupply(totalSupply)
	mc._setBalanceOf(owner, totalSupply)
}

//@:public:receipt
type receipt interface {
	emitTransferMyCoin(token, from, to types.Address, value bn.Number)
}

//Transfer transfer coins from sender to another
//@:public:method:gas[500]
//@:public:interface:gas[450]
func (mc *Mycoin) Transfer(to types.Address, value bn.Number) {
	sdk.Require(value.IsPositive(),
		types.ErrInvalidParameter, "value must be positive")

	totalSupply := mc._totalSupply()
	sdk.Require(totalSupply.IsPositive(),
		types.ErrUserDefined, "Must init first")

	sender := mc.sdk.Message().Sender().Address()
	balanceOfSender := mc._balanceOf(sender).Sub(value)
	sdk.Require(balanceOfSender.IsGEI(0),
		types.ErrInsufficientBalance, "")

	receiver := to
	balanceOfReceiver := mc._balanceOf(receiver).Add(value)

	mc._setBalanceOf(sender, balanceOfSender)
	mc._setBalanceOf(receiver, balanceOfReceiver)

	mc.emitTransferMyCoin(
		mc.sdk.Message().Contract().Address(),
		sender,
		receiver,
		value)
}
