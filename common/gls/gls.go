package bcbgls

import "github.com/jtolds/gls"

var (
	Mgr    = gls.NewContextManager()
	SDKKey = gls.GenSym()
)

type Values = gls.Values
