package mydice2win

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_setSecretSigner This is a method of Dice2Win
func (dw *Dice2Win) _setSecretSigner(v types.PubKey) {
	dw.sdk.Helper().StateHelper().McSet("/secretSigner", &v)
}

//_secretSigner This is a method of Dice2Win
func (dw *Dice2Win) _secretSigner() types.PubKey {

	return *dw.sdk.Helper().StateHelper().McGetEx("/secretSigner", new(types.PubKey)).(*types.PubKey)
}

//_clrSecretSigner This is a method of Dice2Win
func (dw *Dice2Win) _clrSecretSigner() {
	dw.sdk.Helper().StateHelper().McClear("/secretSigner")
}

//_chkSecretSigner This is a method of Dice2Win
func (dw *Dice2Win) _chkSecretSigner() bool {
	return dw.sdk.Helper().StateHelper().Check("/secretSigner")
}

//_McChkSecretSigner This is a method of Dice2Win
func (dw *Dice2Win) _McChkSecretSigner() bool {
	return dw.sdk.Helper().StateHelper().McCheck("/secretSigner")
}

//_setBet This is a method of Dice2Win
func (dw *Dice2Win) _setBet(k string, v *Bet) {
	dw.sdk.Helper().StateHelper().Set(fmt.Sprintf("/bet/%v", k), v)
}

//_bet This is a method of Dice2Win
func (dw *Dice2Win) _bet(k string) *Bet {

	return dw.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/bet/%v", k), new(Bet)).(*Bet)
}

//_chkBet This is a method of Dice2Win
func (dw *Dice2Win) _chkBet(k string) bool {
	return dw.sdk.Helper().StateHelper().Check(fmt.Sprintf("/bet/%v", k))
}

//_setLockedAmount This is a method of Dice2Win
func (dw *Dice2Win) _setLockedAmount(k string, v bn.Number) {
	dw.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/lockedAmount/%v", k), &v)
}

//_lockedAmount This is a method of Dice2Win
func (dw *Dice2Win) _lockedAmount(k string) bn.Number {
	temp := bn.N(0)
	return *dw.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/lockedAmount/%v", k), &temp).(*bn.Number)
}

//_clrLockedAmount This is a method of Dice2Win
func (dw *Dice2Win) _clrLockedAmount(k string) {
	dw.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/lockedAmount/%v", k))
}

//_chkLockedAmount This is a method of Dice2Win
func (dw *Dice2Win) _chkLockedAmount(k string) bool {
	return dw.sdk.Helper().StateHelper().Check(fmt.Sprintf("/lockedAmount/%v", k))
}

//_McChkLockedAmount This is a method of Dice2Win
func (dw *Dice2Win) _McChkLockedAmount(k string) bool {
	return dw.sdk.Helper().StateHelper().McCheck(fmt.Sprintf("/lockedAmount/%v", k))
}

//_setSettings This is a method of Dice2Win
func (dw *Dice2Win) _setSettings(v *Settings) {
	dw.sdk.Helper().StateHelper().McSet("/settings", v)
}

//_settings This is a method of Dice2Win
func (dw *Dice2Win) _settings() *Settings {

	return dw.sdk.Helper().StateHelper().McGetEx("/settings", new(Settings)).(*Settings)
}

//_clrSettings This is a method of Dice2Win
func (dw *Dice2Win) _clrSettings() {
	dw.sdk.Helper().StateHelper().McClear("/settings")
}

//_chkSettings This is a method of Dice2Win
func (dw *Dice2Win) _chkSettings() bool {
	return dw.sdk.Helper().StateHelper().Check("/settings")
}

//_McChkSettings This is a method of Dice2Win
func (dw *Dice2Win) _McChkSettings() bool {
	return dw.sdk.Helper().StateHelper().McCheck("/settings")
}

//_setRecvFeeInfos This is a method of Dice2Win
func (dw *Dice2Win) _setRecvFeeInfos(v []RecvFeeInfo) {
	dw.sdk.Helper().StateHelper().McSet("/recvFeeInfos", &v)
}

//_recvFeeInfos This is a method of Dice2Win
func (dw *Dice2Win) _recvFeeInfos() []RecvFeeInfo {

	return *dw.sdk.Helper().StateHelper().McGetEx("/recvFeeInfos", new([]RecvFeeInfo)).(*[]RecvFeeInfo)
}

//_clrRecvFeeInfos This is a method of Dice2Win
func (dw *Dice2Win) _clrRecvFeeInfos() {
	dw.sdk.Helper().StateHelper().McClear("/recvFeeInfos")
}

//_chkRecvFeeInfos This is a method of Dice2Win
func (dw *Dice2Win) _chkRecvFeeInfos() bool {
	return dw.sdk.Helper().StateHelper().Check("/recvFeeInfos")
}

//_McChkRecvFeeInfos This is a method of Dice2Win
func (dw *Dice2Win) _McChkRecvFeeInfos() bool {
	return dw.sdk.Helper().StateHelper().McCheck("/recvFeeInfos")
}
