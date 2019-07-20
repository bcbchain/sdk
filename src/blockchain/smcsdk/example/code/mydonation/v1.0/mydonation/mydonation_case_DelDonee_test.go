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

//DelDonee This is a method of MySuite
func (mysuit *MySuite) TestDemo_DelDonee(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)

	gls.Mgr.SetValues(gls.Values{gls.SDKKey: utest.UTP.ISmartContract}, func() {
		test := NewTestObject(contractOwner)
		test.setSender(contractOwner).InitChain()
		mysuit.test_DelDonee(contractOwner, test)
	})
}

func (mysuit *MySuite) test_DelDonee(owner sdk.IAccount, test *TestObject) {
	fmt.Println("=== Run  UnitTest case: DelDonee(donee types.Address)")

	fmt.Println("=== prepare for test")
	zero := bn.N(0)
	oneCoin := bn.N(1000000000)
	halfCoin := bn.N(500000000)
	a1 := utest.NewAccount("TSC", oneCoin)
	a2 := utest.NewAccount("TSC", oneCoin)
	a3 := utest.NewAccount("TSC", oneCoin)
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).AddDonee(a1.Address())
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).AddDonee(a2.Address())
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(a1)
		utest.Assert(t.transfer(halfCoin) != nil)
		return t.Donate(a2.Address())
	})
	fmt.Println("=== pass")

	fmt.Println("=== test for authorization")
	test.run(types.ErrNoAuthorization, func(t *TestObject) types.Error {
		return t.setSender(a2).DelDonee(a1.Address())
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
		{owner, owner.Address(), errDoneeNotExist},
		{owner, utest.GetContract().Address(), errDoneeNotExist},
		{owner, utest.GetContract().Account().Address(), errDoneeNotExist},
		{owner, a3.Address(), errDoneeNotExist},
	}
	//run
	for _, c := range cases {
		test.run(c.codeExpect, func(t *TestObject) types.Error {
			return t.setSender(c.sender).DelDonee(c.addr)
		})
	}
	fmt.Println("=== pass")

	fmt.Println("=== test for business logic 1")
	//prepare
	utest.AssertSDB(keyOfDonation(a1.Address()), &zero)
	//run
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).DelDonee(a1.Address())
	})
	//check
	utest.AssertSDB(keyOfDonation(a1.Address()), nil)
	//run
	test.run(errDoneeNotExist, func(t *TestObject) types.Error {
		t.setSender(owner)
		return t.DelDonee(a1.Address())
	})
	//check
	utest.AssertSDB(keyOfDonation(a1.Address()), nil)
	fmt.Println("=== pass")

	fmt.Println("=== test for business logic 2")
	//prepare
	utest.AssertSDB(keyOfDonation(a2.Address()), &halfCoin)
	//run
	test.run(errDonationExist, func(t *TestObject) types.Error {
		return t.setSender(owner).DelDonee(a2.Address())
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(owner)
		return t.Transfer(a2.Address(), halfCoin)
	})

	utest.AssertSDB(keyOfDonation(a2.Address()), &zero)
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(owner)
		return t.DelDonee(a2.Address())
	})
	//check
	utest.AssertSDB(keyOfDonation(a2.Address()), nil)
	//run
	test.run(errDoneeNotExist, func(t *TestObject) types.Error {
		t.setSender(owner)
		return t.DelDonee(a2.Address())
	})
	//check
	utest.AssertBalance(a2, "TSC", oneCoin.Add(halfCoin))
	fmt.Println("=== pass")
}
