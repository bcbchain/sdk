package sdk

import (
	"blockchain/smcsdk/common/gls"
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

// RequireOwner method for require owner
func RequireOwner() {
	var sdk ISmartContract
	if iSDK, ok := gls.Mgr.GetValue(gls.SDKKey); !ok {
		err := types.Error{ErrorCode: types.ErrStubDefined, ErrorDesc: "gls cannot get sdk"}
		panic(err)
	} else {
		sdk = iSDK.(ISmartContract)
	}

	if sdk.Message().Sender().Address() != sdk.Message().Contract().Owner().Address() {
		err := types.Error{ErrorCode: types.ErrNoAuthorization, ErrorDesc: "only contract owner just can do it"}
		panic(err)
	}
}

// RequireAddress method for require address
func RequireAddress(addr types.Address) {
	var sdk ISmartContract
	if iSDK, ok := gls.Mgr.GetValue(gls.SDKKey); !ok {
		err := types.Error{ErrorCode: types.ErrStubDefined, ErrorDesc: "gls cannot get sdk"}
		panic(err)
	} else {
		sdk = iSDK.(ISmartContract)
	}

	if err := sdk.Helper().BlockChainHelper().CheckAddress(addr); err != nil {
		err2 := types.Error{ErrorCode: types.ErrInvalidAddress, ErrorDesc: err.Error()}
		panic(err2)
	}
}
