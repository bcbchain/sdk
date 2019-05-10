package excellencies

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_secretSigner This is a method of Excellencies
func (e *Excellencies) _secretSigner() types.PubKey {

	return *e.sdk.Helper().StateHelper().McGetEx("/secretSigner", new(types.PubKey)).(*types.PubKey)
}

//_setSecretSigner This is a method of Excellencies
func (e *Excellencies) _setSecretSigner(v types.PubKey) {
	e.sdk.Helper().StateHelper().McSet("/secretSigner", &v)
}

//_chkSecretSigner This is a method of Excellencies
func (e *Excellencies) _chkSecretSigner() bool {
	return e.sdk.Helper().StateHelper().Check("/secretSigner")
}

//_clrSecretSigner This is a method of Excellencies
func (e *Excellencies) _clrSecretSigner() {
	e.sdk.Helper().StateHelper().McClear("/secretSigner")
}

//_delSecretSigner This is a method of Excellencies
func (e *Excellencies) _delSecretSigner() {
	e.sdk.Helper().StateHelper().Delete("/secretSigner")
}

//_betInfo This is a method of Excellencies
func (e *Excellencies) _betInfo(k1 string, k2 string) *BetInfo {

	return e.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/betInfo/%v/%v", k1, k2), new(BetInfo)).(*BetInfo)
}

//_setBetInfo This is a method of Excellencies
func (e *Excellencies) _setBetInfo(k1 string, k2 string, v *BetInfo) {
	e.sdk.Helper().StateHelper().Set(fmt.Sprintf("/betInfo/%v/%v", k1, k2), v)
}

//_chkBetInfo This is a method of Excellencies
func (e *Excellencies) _chkBetInfo(k1 string, k2 string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/betInfo/%v/%v", k1, k2))
}

//_delBetInfo This is a method of Excellencies
func (e *Excellencies) _delBetInfo(k1 string, k2 string) {
	e.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/betInfo/%v/%v", k1, k2))
}

//_mapSetting This is a method of Excellencies
func (e *Excellencies) _mapSetting() *MapSetting {

	return e.sdk.Helper().StateHelper().GetEx("/mapSetting", new(MapSetting)).(*MapSetting)
}

//_setMapSetting This is a method of Excellencies
func (e *Excellencies) _setMapSetting(v *MapSetting) {
	e.sdk.Helper().StateHelper().Set("/mapSetting", v)
}

//_chkMapSetting This is a method of Excellencies
func (e *Excellencies) _chkMapSetting() bool {
	return e.sdk.Helper().StateHelper().Check("/mapSetting")
}

//_delMapSetting This is a method of Excellencies
func (e *Excellencies) _delMapSetting() {
	e.sdk.Helper().StateHelper().Delete("/mapSetting")
}

//_recFeeInfo This is a method of Excellencies
func (e *Excellencies) _recFeeInfo() []RecFeeInfo {

	return *e.sdk.Helper().StateHelper().McGetEx("/recFeeInfo", new([]RecFeeInfo)).(*[]RecFeeInfo)
}

//_setRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _setRecFeeInfo(v []RecFeeInfo) {
	e.sdk.Helper().StateHelper().McSet("/recFeeInfo", &v)
}

//_chkRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _chkRecFeeInfo() bool {
	return e.sdk.Helper().StateHelper().Check("/recFeeInfo")
}

//_clrRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _clrRecFeeInfo() {
	e.sdk.Helper().StateHelper().McClear("/recFeeInfo")
}

//_delRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _delRecFeeInfo() {
	e.sdk.Helper().StateHelper().Delete("/recFeeInfo")
}

//_roundInfo This is a method of Excellencies
func (e *Excellencies) _roundInfo(k string) RoundInfo {

	return *e.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/roundInfo/%v", k), new(RoundInfo)).(*RoundInfo)
}

//_setRoundInfo This is a method of Excellencies
func (e *Excellencies) _setRoundInfo(k string, v RoundInfo) {
	e.sdk.Helper().StateHelper().Set(fmt.Sprintf("/roundInfo/%v", k), &v)
}

//_chkRoundInfo This is a method of Excellencies
func (e *Excellencies) _chkRoundInfo(k string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/roundInfo/%v", k))
}

//_delRoundInfo This is a method of Excellencies
func (e *Excellencies) _delRoundInfo(k string) {
	e.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/roundInfo/%v", k))
}

//_poolAmount This is a method of Excellencies
func (e *Excellencies) _poolAmount(k1 types.Address, k2 string) bn.Number {
	temp := bn.N(0)
	return *e.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/poolAmount/%v/%v", k1, k2), &temp).(*bn.Number)
}

//_setPoolAmount This is a method of Excellencies
func (e *Excellencies) _setPoolAmount(k1 types.Address, k2 string, v bn.Number) {
	e.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/poolAmount/%v/%v", k1, k2), &v)
}

//_chkPoolAmount This is a method of Excellencies
func (e *Excellencies) _chkPoolAmount(k1 types.Address, k2 string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/poolAmount/%v/%v", k1, k2))
}

//_clrPoolAmount This is a method of Excellencies
func (e *Excellencies) _clrPoolAmount(k1 types.Address, k2 string) {
	e.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/poolAmount/%v/%v", k1, k2))
}

//_delPoolAmount This is a method of Excellencies
func (e *Excellencies) _delPoolAmount(k1 types.Address, k2 string) {
	e.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/poolAmount/%v/%v", k1, k2))
}

//_playerIndex This is a method of Excellencies
func (e *Excellencies) _playerIndex(k1 string, k2 string) *PlayerIndexes {

	return e.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/playerIndex/%v/%v", k1, k2), new(PlayerIndexes)).(*PlayerIndexes)
}

//_setPlayerIndex This is a method of Excellencies
func (e *Excellencies) _setPlayerIndex(k1 string, k2 string, v *PlayerIndexes) {
	e.sdk.Helper().StateHelper().Set(fmt.Sprintf("/playerIndex/%v/%v", k1, k2), v)
}

//_chkPlayerIndex This is a method of Excellencies
func (e *Excellencies) _chkPlayerIndex(k1 string, k2 string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/playerIndex/%v/%v", k1, k2))
}

//_delPlayerIndex This is a method of Excellencies
func (e *Excellencies) _delPlayerIndex(k1 string, k2 string) {
	e.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/playerIndex/%v/%v", k1, k2))
}

//_grandPrizer This is a method of Excellencies
func (e *Excellencies) _grandPrizer(k1 types.Address, k2 string) GrandPrizer {

	return *e.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/grandPrizer/%v/%v", k1, k2), new(GrandPrizer)).(*GrandPrizer)
}

//_setGrandPrizer This is a method of Excellencies
func (e *Excellencies) _setGrandPrizer(k1 types.Address, k2 string, v GrandPrizer) {
	e.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/grandPrizer/%v/%v", k1, k2), &v)
}

//_chkGrandPrizer This is a method of Excellencies
func (e *Excellencies) _chkGrandPrizer(k1 types.Address, k2 string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/grandPrizer/%v/%v", k1, k2))
}

//_clrGrandPrizer This is a method of Excellencies
func (e *Excellencies) _clrGrandPrizer(k1 types.Address, k2 string) {
	e.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/grandPrizer/%v/%v", k1, k2))
}

//_delGrandPrizer This is a method of Excellencies
func (e *Excellencies) _delGrandPrizer(k1 types.Address, k2 string) {
	e.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/grandPrizer/%v/%v", k1, k2))
}
