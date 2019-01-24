package mycoin

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/utest"
	"gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type MySuite struct{}

var _ = check.Suite(&MySuite{})

func (mysuit *MySuite) TestMycoin_Transfer(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, nil)
	test := NewTestObject(contractOwner)

	acct := utest.NewAccount(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1000000000))
	if acct == nil {
		panic("初始化账户失败")
	}

	test.run().InitChain()

	test.run().setSender(contractOwner).Transfer(acct.Address(), bn.N(1000000000))
}
