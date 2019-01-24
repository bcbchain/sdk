package sdk

import (
	"blockchain/smcsdk/sdk/types"
)

// Require method for require result
func Require(expr bool, errCode uint32, errInfo string) {
	if expr == false {
		err := types.Error{ErrorCode: errCode, ErrorDesc: errInfo}
		panic(err)
	}
}

// RequireNotError method for require result
func RequireNotError(err error, errCode uint32) {
	if err != nil {
		err := types.Error{ErrorCode: errCode, ErrorDesc: err.Error()}
		panic(err)
	}
}

// RequireOwner  method for require owner
func RequireOwner(sdk ISmartContract) {
	if sdk.Message().Sender().Address() != sdk.Message().Contract().Owner() {
		err := types.Error{ErrorCode: types.ErrNoAuthorization, ErrorDesc: "only contract owner just can do it"}
		panic(err)
	}
}

// RequireAddress  method for require address
func RequireAddress(sdk ISmartContract, addr types.Address) {
	if err := sdk.Helper().BlockChainHelper().CheckAddress(addr); err != nil {
		err2 := types.Error{ErrorCode: types.ErrInvalidAddress, ErrorDesc: err.Error()}
		panic(err2)
	}
}
