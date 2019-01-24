package mycoin

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_setTotalSupply This is a method of Mycoin
func (m *Mycoin) _setTotalSupply(v bn.Number) {
	m.sdk.Helper().StateHelper().Set("/totalSupply", &v)
}

//_totalSupply This is a method of Mycoin
func (m *Mycoin) _totalSupply() bn.Number {
	temp := bn.N(0)
	return *m.sdk.Helper().StateHelper().GetEx("/totalSupply", &temp).(*bn.Number)
}

//_chkTotalSupply This is a method of Mycoin
func (m *Mycoin) _chkTotalSupply() bool {
	return m.sdk.Helper().StateHelper().Check("/totalSupply")
}

//_setBalanceOf This is a method of Mycoin
func (m *Mycoin) _setBalanceOf(k types.Address, v bn.Number) {
	m.sdk.Helper().StateHelper().Set(fmt.Sprintf("/balanceOf/%v", k), &v)
}

//_balanceOf This is a method of Mycoin
func (m *Mycoin) _balanceOf(k types.Address) bn.Number {
	temp := bn.N(0)
	return *m.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/balanceOf/%v", k), &temp).(*bn.Number)
}

//_chkBalanceOf This is a method of Mycoin
func (m *Mycoin) _chkBalanceOf(k types.Address) bool {
	return m.sdk.Helper().StateHelper().Check(fmt.Sprintf("/balanceOf/%v", k))
}
