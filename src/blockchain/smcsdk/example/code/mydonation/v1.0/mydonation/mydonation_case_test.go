package mydonation

import (
	"blockchain/smcsdk/utest"
	"fmt"
	"gopkg.in/check.v1"
	"testing"
)

//Test This is a function
func Test(t *testing.T) { check.TestingT(t) }

//MySuite This is a struct
type MySuite struct{}

var _ = check.Suite(&MySuite{})

//TestMydonation_AddDonee This is a method of MySuite
func (mysuit *MySuite) TestMydonation_AddDonee(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()

	fmt.Println("=== Run UnitTestcase: AddDonee(donee types.Address)")
	//TODO
	fmt.Println("--- Pass UnitTestcase: AddDonee(donee types.Address)")
}

//TestMydonation_DelDonee This is a method of MySuite
func (mysuit *MySuite) TestMydonation_DelDonee(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()

	fmt.Println("=== Run UnitTestcase: DelDonee(donee types.Address)")
	//TODO
	fmt.Println("--- Pass UnitTestcase: DelDonee(donee types.Address)")
}

//TestMydonation_Donate This is a method of MySuite
func (mysuit *MySuite) TestMydonation_Donate(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()

	fmt.Println("=== Run UnitTestcase: Donate(donee types.Address)")
	//TODO
	fmt.Println("--- Pass UnitTestcase: Donate(donee types.Address)")
}

//TestMydonation_Transfer This is a method of MySuite
func (mysuit *MySuite) TestMydonation_Transfer(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()

	fmt.Println("=== Run UnitTestcase: Transfer(donee types.Address, value bn.Number)")
	//TODO
	fmt.Println("--- Pass UnitTestcase: Transfer(donee types.Address, value bn.Number)")
}
