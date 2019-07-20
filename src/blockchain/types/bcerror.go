package types

//BcError structure of bcerror
type BcError struct {
	ErrorCode uint32 // Error code
	ErrorDesc string // Error description
}

// Error() gets error description with error code
func (bcerror *BcError) Error() string {
	if bcerror.ErrorDesc != "" {
		return bcerror.ErrorDesc
	}

	for _, error := range bcErrors {
		if error.ErrorCode == bcerror.ErrorCode {
			return error.ErrorDesc
		}
	}
	return ""
}

//CodeOK means success
const (
	CodeOK = 200 + iota
)

// ErrMarshal For smart contracts custom errors
const (
	ErrMarshal = 500 + iota
	ErrCallRPC
	ErrDealFailed
	ErrAccountLocked
)

//ErrCheckTx beginning error code of checkTx
const (
	ErrCheckTx = 600 + iota
)

//ErrDeliverTx beginning error code of deliverTx
const (
	ErrDeliverTx = 700 + iota
)

const (
	ErrCodeNoAuthorization = 1000 + iota
)

// For lowlevel (stateDB, go libs, 3rd party) errors
// only set error code and uses original error message
const (
	ErrCodeLowLevelError = 5000 + iota
)

var bcErrors = []BcError{
	{CodeOK, ""},

	{ErrMarshal, "Json marshal error"},
	{ErrCallRPC, "Call rpc error"},
	{ErrDealFailed, "The deal failed"},
	{ErrAccountLocked, "Account is locked"},

	{ErrCheckTx, "CheckTx failed"},

	//ErrCodeNoAuthorization
	{ErrCodeNoAuthorization, "No authorization"},

	{ErrDeliverTx, "DeliverTx failed"},
}
