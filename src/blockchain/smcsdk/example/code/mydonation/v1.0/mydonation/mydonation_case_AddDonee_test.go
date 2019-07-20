package mydonation

import (
	"blockchain/smcsdk/common/gls"
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/utest"
	"fmt"

	"gopkg.in/check.v1"
)

func keyOfDonation(addr types.Address) string {
	return fmt.Sprintf("/donations/%v", addr)
}

//AddDonee This is a method of MySuite
func (mysuit *MySuite) TestDemo_AddDonee(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)

	gls.Mgr.SetValues(gls.Values{gls.SDKKey: utest.UTP.ISmartContract}, func() {
		test := NewTestObject(contractOwner)
		test.setSender(contractOwner).InitChain()
		mysuit.test_AddDonee(contractOwner, test)
	})
}

func (mysuit *MySuite) test_AddDonee(owner sdk.IAccount, test *TestObject) {
	fmt.Println("=== Run  UnitTest case: AddDonee(donee types.Address)")

	fmt.Println("=== prepare for test")
	zero := bn.N(0)
	oneCoin := bn.N(1000000000)
	a1 := utest.NewAccount("TSC", oneCoin)
	a2 := utest.NewAccount("TSC", oneCoin)
	fmt.Println("=== pass")

	fmt.Println("=== test for authorization")
	test.run(types.ErrNoAuthorization, func(t *TestObject) types.Error {
		return t.setSender(a1).AddDonee(a2.Address())
	})
	fmt.Println("=== pass")

	fmt.Println("=== test for parameters")
	//prepare
	var cases = []struct {
		sender     sdk.IAccount
		addr       types.Address
		codeExpect uint32
	}{
		{owner, "", types.ErrInvalidAddress},
		{owner, "local", types.ErrInvalidAddress},
		{owner, "localhshskhjkshfsswtsysyst6t76ddsg7s7w", types.ErrInvalidAddress},
		{owner, owner.Address(), errDoneeCannotBeOwner},
		{owner, utest.GetContract().Address(), errDoneeCannotBeSmc},
		{owner, utest.GetContract().Account().Address(), errDoneeCannotBeSmc},
	}
	//run
	for _, c := range cases {
		test.run(c.codeExpect, func(t *TestObject) types.Error {
			return t.setSender(c.sender).AddDonee(c.addr)
		})
	}
	fmt.Println("=== pass")

	fmt.Println("=== test for business logic")
	//prepare
	utest.AssertSDB(keyOfDonation(a1.Address()), nil)
	//run
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).AddDonee(a1.Address())
	})
	test.run(errDoneeAlreadyExist, func(t *TestObject) types.Error {
		return t.setSender(owner).AddDonee(a1.Address())
	})
	//check
	utest.AssertSDB(keyOfDonation(a1.Address()), &zero)
	utest.AssertSDB(keyOfDonation(a2.Address()), nil)
	fmt.Println("=== pass")
}
