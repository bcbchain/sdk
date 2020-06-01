//unittestplatform
//tx.go 实现与交易数据相关的功能，包括：
// 创建交易数据对象，设置交易数据的发送人，初始化交易数据消息，获取当前交易的发送人、获取当前合约等。

package utest

import (
	"github.com/bcbchain/sdk/common/gls"
	"github.com/bcbchain/sdk/sdk"
	"github.com/bcbchain/sdk/sdk/bn"
	"github.com/bcbchain/sdk/sdk/std"
	"github.com/bcbchain/sdk/sdk/types"
	"github.com/bcbchain/sdk/sdkimpl"
	"github.com/bcbchain/sdk/sdkimpl/object"
)

//SetSender set sender to
func SetSender(_sender types.Address) sdk.ISmartContract {

	// 更新dg.SmartContract的tx
	bcbgls.Mgr.SetValues(bcbgls.Values{bcbgls.SDKKey: UTP.ISmartContract}, func() {
		_tx := object.NewTx(UTP.ISmartContract, "", 500000, 0, []byte(""), _sender)
		api := UTP.ISmartContract.(*sdkimpl.SmartContract)
		api.SetTx(_tx)

		senderAcct := object.NewAccount(api, _sender)
		_msg := object.NewMessage(UTP.ISmartContract,
			UTP.ISmartContract.Message().Contract(),
			UTP.ISmartContract.Message().MethodID(),
			UTP.ISmartContract.Message().Items(),
			senderAcct.Address(),
			senderAcct.Address(),
			UTP.ISmartContract.Message().Origins(),
			UTP.ISmartContract.Message().InputReceipts(),
		)
		api.SetMessage(_msg)
	})

	return UTP.ISmartContract
}

//ResetMsg reset message
func ResetMsg() sdk.ISmartContract {

	bcbgls.Mgr.SetValues(bcbgls.Values{bcbgls.SDKKey: UTP.ISmartContract}, func() {
		_message := object.NewMessage(UTP.ISmartContract,
			UTP.Message().Contract(),
			UTP.Message().MethodID(),
			UTP.Message().Items(),
			UTP.Message().Sender().Address(),
			UTP.Message().Payer().Address(),
			UTP.Message().Origins(),
			UTP.Message().InputReceipts())
		api := UTP.ISmartContract.(*sdkimpl.SmartContract)
		api.SetMessage(_message)
	})

	return UTP.ISmartContract
}

//GetSender get sender object
func GetSender() sdk.IAccount { return UTP.Tx().Signer() }

//GetContract get contract object
func GetContract() sdk.IContract { return UTP.Message().Contract() }

func (ut *UtPlatform) setTxSender(_sender types.Address) sdk.ISmartContract {

	bcbgls.Mgr.SetValues(bcbgls.Values{bcbgls.SDKKey: ut.ISmartContract}, func() {
		tx := ut.Tx()
		acct := object.NewAccount(ut.ISmartContract, _sender)
		o := object.NewTx(ut.ISmartContract, tx.Note(), tx.GasLimit(), tx.GasLeft(), tx.TxHash(), acct.Address())

		api := ut.ISmartContract.(*sdkimpl.SmartContract)
		api.SetTx(o)
	})

	return ut.ISmartContract
}

//SetReceipt sets receipt
func (ut *UtPlatform) SetReceipt(token string, from, to types.Address, value bn.Number) (err types.Error) {
	tr := std.Transfer{
		Token: token,
		From:  from,
		To:    to,
		Value: value,
	}
	rh := ut.Helper().ReceiptHelper()
	rh.Emit(&tr)

	return
}
