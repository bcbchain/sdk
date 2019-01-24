package types

//BcError structure of bcerror
type BcError struct {
	ErrorCode uint32 // Error code
	ErrorDesc string // Error description
}

//todo 移到bclib，改名BcError
// Error() gets error description with error code
func (bcerror *BcError) Error() string {
	if bcerror.ErrorDesc != "" {
		return bcerror.ErrorDesc
	}

	// TODO: it would be better to use binary search
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
	ErrOutOfRange
	ErrNeedPositiveNumber
)

//ErrCheckTx beginning error code of checkTx
const (
	ErrCheckTx = 600 + iota
	ErrCheckInsufficientBalance
)

//ErrDeliverTx beginning error code of deliverTx
const (
	ErrDeliverTx = 700 + iota
	ErrDeliverInsufficientBalance
)

var bcErrors = []BcError{
	{CodeOK, ""},

	{ErrMarshal, "Json marshal error"},
	{ErrCallRPC, "Call rpc error"},
	{ErrOutOfRange, "Out of range"},
	{ErrNeedPositiveNumber, "Must positive number"},

	{ErrCheckTx, "CheckTx failed"},
	{ErrCheckInsufficientBalance, "Insufficient balance"},

	{ErrDeliverTx, "DeliverTx failed"},
	{ErrDeliverInsufficientBalance, "Insufficient balance"},
}
