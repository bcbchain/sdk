package mydonation

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_donations This is a method of Mydonation
func (m *Mydonation) _donations(k types.Address) bn.Number {
	temp := bn.N(0)
	return *m.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/donations/%v", k), &temp).(*bn.Number)
}

//_chkDonations This is a method of Mydonation
func (m *Mydonation) _chkDonations(k types.Address) bool {
	return m.sdk.Helper().StateHelper().Check(fmt.Sprintf("/donations/%v", k))
}

//_setDonations This is a method of Mydonation
func (m *Mydonation) _setDonations(k types.Address, v bn.Number) {
	m.sdk.Helper().StateHelper().Set(fmt.Sprintf("/donations/%v", k), &v)
}

//_delDonations This is a method of Mydonation
func (m *Mydonation) _delDonations(k types.Address) {
	m.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/donations/%v", k))
}
