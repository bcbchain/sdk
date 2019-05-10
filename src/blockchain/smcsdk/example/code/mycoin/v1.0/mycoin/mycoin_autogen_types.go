package mycoin

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

// TransferParam structure of parameters of Transfer() of v2.0
type TransferParam struct {
	To    types.Address
	Value bn.Number
}
