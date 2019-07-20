package mystorage

import (
	"blockchain/smcsdk/common/gls"
	types2 "blockchain/smcsdk/sdk/types"
	"blockchain/types"
	"math"
	"testing"

	ut "blockchain/smcsdk/utest"
	"gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type MySuite struct{}

var _ = check.Suite(&MySuite{})

func (mysuit *MySuite) TestMyStorage_Set(c *check.C) {
	ut.Init(orgID)
	contractOwner := ut.DeployContract(c, contractName, orgID, contractMethods, nil)
	gls.Mgr.SetValues(gls.Values{gls.SDKKey: ut.UTP.ISmartContract}, func() {
		test := NewTestObject(contractOwner)

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return t.Set(0)
		})
		ut.AssertSDB("/storedData", 0)
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == 0)

			return err
		})

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return t.Set(2000348989)
		})
		ut.AssertSDB("/storedData", 2000348989)
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == 2000348989)

			return err
		})

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return t.Set(math.MaxUint64)
		})
		var maxUint64 uint64 = math.MaxUint64
		ut.AssertSDB("/storedData", &maxUint64)
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == math.MaxUint64)

			return err
		})
	})
}

func (mysuit *MySuite) TestMyStorage_Get(c *check.C) {
	ut.Init(orgID)
	contractOwner := ut.DeployContract(c, contractName, orgID, contractMethods, nil)
	gls.Mgr.SetValues(gls.Values{gls.SDKKey: ut.UTP.ISmartContract}, func() {
		test := NewTestObject(contractOwner)

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return t.Set(0)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == 0)

			return err
		})

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(1)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == 1)

			return err
		})

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(2)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == 2)

			return err
		})

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(3)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(4)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(6)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(9)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(10)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(11)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == 11)

			return err
		})

		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			return test.Set(math.MaxUint64)
		})
		test.run(types.CodeOK, func(t *TestObject) types2.Error {
			retV, err := test.Get()
			ut.Assert(retV == math.MaxUint64)

			return err
		})
	})
}
