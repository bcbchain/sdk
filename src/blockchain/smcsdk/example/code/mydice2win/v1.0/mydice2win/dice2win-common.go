package mydice2win

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

func (dw *Dice2Win) maxBetMask() bn.Number {
	return bn.N(2).Exp(bn.N(maxMaskModulo))
}

func (dw *Dice2Win) popCntMult() bn.Number {
	return bn.NewNumberStringBase("0000000000002000000000100000000008000000000400000000020000000001", 16)
}

func (dw *Dice2Win) popCntMask() bn.Number {
	return bn.NewNumberStringBase("0001041041041041041041041041041041041041041041041041041041041041", 16)
}

func (dw *Dice2Win) popCntModulo() bn.Number {
	return bn.NewNumberStringBase("3F", 16)
}

//转账给fee的接收地址
func (dw *Dice2Win) transferToRecvFeeAddr(tokenName types.Address, recvFee bn.Number) {
	if recvFee.CmpI(0) <= 0 {
		return
	}

	infos := dw._recvFeeInfos()
	account := dw.sdk.Message().Contract().Account()
	forx.Range(infos, func(i int, info RecvFeeInfo) bool {
		account.TransferByName(tokenName, info.Address, recvFee.MulI(info.Ratio).DivI(perMille))
		return true
	})
}

func (dw *Dice2Win) checkRecvFeeInfos(infos []RecvFeeInfo) {
	sdk.Require(len(infos) > 0,
		types.ErrInvalidParameter, "The length of RecvFeeInfos must be larger than zero")

	allRatio := int64(0)
	forx.Range(infos, func(i int, info RecvFeeInfo) bool {
		sdk.Require(info.Ratio > 0,
			types.ErrInvalidParameter, "ratio must be larger than zero")
		sdk.RequireAddress(info.Address)
		sdk.Require(info.Address != dw.sdk.Message().Contract().Account().Address(),
			types.ErrInvalidParameter, "address cannot be contract account address")

		allRatio += info.Ratio

		return true
	})

	//设置的分配比例加起来必须等于1000
	sdk.Require(allRatio <= 1000, types.ErrInvalidParameter,
		"The sum of ratio must be less or equal 1000")
}

func (dw *Dice2Win) checkSettings(newSettings *Settings) {

	sdk.Require(len(newSettings.TokenNames) > 0,
		types.ErrInvalidParameter, "tokenNames cannot be empty")

	forx.Range(newSettings.TokenNames, func(tokenName string, v struct{}) bool {
		token := dw.sdk.Helper().TokenHelper().TokenOfName(tokenName)
		sdk.Require(token != nil,
			types.ErrInvalidParameter, fmt.Sprintf("tokenName=%s is not exist", tokenName))

		return true
	})

	sdk.Require(newSettings.MaxBet > 0,
		types.ErrInvalidParameter, "MaxBet must be bigger than zero")

	sdk.Require(newSettings.MaxProfit >= 0,
		types.ErrInvalidParameter, "MaxProfit can not be negative")

	sdk.Require(newSettings.MinBet > 0 && newSettings.MinBet < newSettings.MaxBet,
		types.ErrInvalidParameter, "MinBet must be bigger than zero and smaller than MaxBet")

	sdk.Require(newSettings.SendToCltRatio >= 0 && newSettings.SendToCltRatio < perMille,
		types.ErrInvalidParameter,
		fmt.Sprintf("SendToCltRatio must be bigger than zero and smaller than %d", perMille))

	sdk.Require(newSettings.FeeRatio > 0 && newSettings.FeeRatio < perMille,
		types.ErrInvalidParameter,
		fmt.Sprintf("FeeRatio must be bigger than zero and  smaller than %d", perMille))

	sdk.Require(newSettings.FeeMinimum > 0,
		types.ErrInvalidParameter, "FeeMinimum must be bigger than zero")

	sdk.Require(newSettings.BetExpirationBlocks > 0,
		types.ErrInvalidParameter, "BetExpirationBlocks must be bigger than zero")
}
