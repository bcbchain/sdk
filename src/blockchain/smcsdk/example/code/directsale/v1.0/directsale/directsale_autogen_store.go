package directsale

import (
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_setSettings This is a method of DirectSale
func (ds *DirectSale) _setSettings(v Settings) {
	ds.sdk.Helper().StateHelper().Set("/settings", &v)
}

//_settings This is a method of DirectSale
func (ds *DirectSale) _settings() Settings {

	return *ds.sdk.Helper().StateHelper().GetEx("/settings", new(Settings)).(*Settings)
}

//_chkSettings This is a method of DirectSale
func (ds *DirectSale) _chkSettings() bool {
	return ds.sdk.Helper().StateHelper().Check("/settings")
}

//_setSalers This is a method of DirectSale
func (ds *DirectSale) _setSalers(k types.Address, v Saler) {
	ds.sdk.Helper().StateHelper().Set(fmt.Sprintf("/salers/%v", k), &v)
}

//_salers This is a method of DirectSale
func (ds *DirectSale) _salers(k types.Address) Saler {

	return *ds.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/salers/%v", k), new(Saler)).(*Saler)
}

//_chkSalers This is a method of DirectSale
func (ds *DirectSale) _chkSalers(k types.Address) bool {
	return ds.sdk.Helper().StateHelper().Check(fmt.Sprintf("/salers/%v", k))
}

//_setManagers This is a method of DirectSale
func (ds *DirectSale) _setManagers(v []types.Address) {
	ds.sdk.Helper().StateHelper().Set("/managers", &v)
}

//_managers This is a method of DirectSale
func (ds *DirectSale) _managers() []types.Address {

	return *ds.sdk.Helper().StateHelper().GetEx("/managers", new([]types.Address)).(*[]types.Address)
}

//_chkManagers This is a method of DirectSale
func (ds *DirectSale) _chkManagers() bool {
	return ds.sdk.Helper().StateHelper().Check("/managers")
}

//_setSalerToApps This is a method of DirectSale
func (ds *DirectSale) _setSalerToApps(k1 types.Address, k2 string, v SalerApp) {
	ds.sdk.Helper().StateHelper().Set(fmt.Sprintf("/salerToApps/%v/%v", k1, k2), &v)
}

//_salerToApps This is a method of DirectSale
func (ds *DirectSale) _salerToApps(k1 types.Address, k2 string) SalerApp {

	return *ds.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/salerToApps/%v/%v", k1, k2), new(SalerApp)).(*SalerApp)
}

//_chkSalerToApps This is a method of DirectSale
func (ds *DirectSale) _chkSalerToApps(k1 types.Address, k2 string) bool {
	return ds.sdk.Helper().StateHelper().Check(fmt.Sprintf("/salerToApps/%v/%v", k1, k2))
}

//_setGlobal This is a method of DirectSale
func (ds *DirectSale) _setGlobal(k string, v Global) {
	ds.sdk.Helper().StateHelper().Set(fmt.Sprintf("/global/%v", k), &v)
}

//_global This is a method of DirectSale
func (ds *DirectSale) _global(k string) Global {

	return *ds.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/global/%v", k), new(Global)).(*Global)
}

//_chkGlobal This is a method of DirectSale
func (ds *DirectSale) _chkGlobal(k string) bool {
	return ds.sdk.Helper().StateHelper().Check(fmt.Sprintf("/global/%v", k))
}
