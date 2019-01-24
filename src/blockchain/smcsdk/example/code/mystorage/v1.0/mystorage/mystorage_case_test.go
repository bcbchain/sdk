package mystorage

import (
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
	test := NewTestObject(contractOwner)

	test.run().Set(0)
	ut.AssertSDB("/mystorage/storedData", 0)
	retV, err := test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == 0)

	test.run().Set(2000348989)
	ut.AssertSDB("/mystorage/storedData", 2000348989)
	retV, err = test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == 2000348989)

	test.run().Set(math.MaxUint64)
	var maxUint64 uint64 = math.MaxUint64
	ut.AssertSDB("/mystorage/storedData", &maxUint64)
	retV, err = test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == math.MaxUint64)
}

func (mysuit *MySuite) TestMyStorage_Get(c *check.C) {
	ut.Init(orgID)
	contractOwner := ut.DeployContract(c, contractName, orgID, contractMethods, nil)
	test := NewTestObject(contractOwner)

	test.run().Set(0)
	retV, err := test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == 0)

	test.run().Set(1)
	retV, err = test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == 1)

	test.run().Set(2)
	retV, err = test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == 2)

	test.run().Set(3)
	test.run().Set(4)
	test.run().Set(6)
	test.run().Set(9)
	test.run().Set(10)
	test.run().Set(11)
	retV, err = test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == 11)

	test.run().Set(math.MaxUint64)
	retV, err = test.run().Get()
	ut.AssertError(err, types.CodeOK)
	ut.Assert(retV == math.MaxUint64)
}
