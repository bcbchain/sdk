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

//Donate This is a method of MySuite
func (mysuit *MySuite) TestDemo_Donate(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)

	gls.Mgr.SetValues(gls.Values{gls.SDKKey: utest.UTP.ISmartContract}, func() {
		test := NewTestObject(contractOwner)
		test.setSender(contractOwner).InitChain()
		mysuit.test_Donate(contractOwner, test)
	})
}

func (mysuit *MySuite) test_Donate(owner sdk.IAccount, test *TestObject) {
	fmt.Println("=== Run  UnitTest case: Donate(donee types.Address)")

	fmt.Println("=== prepare for test")
	halfCoin := bn.N(500000000)
	oneCoin := bn.N(1000000000)
	oneHalfCoin := bn.N(1500000000)
	twoCoin := bn.N(2000000000)
	utest.Transfer(nil, owner.Address(), "TSC", twoCoin)
	utest.Transfer(nil, owner.Address(), "BTC", twoCoin)
	a1 := utest.NewAccount("TSC", oneCoin)
	a2 := utest.NewAccount("TSC", oneCoin)
	a3 := utest.NewAccount("TSC", oneCoin)
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).AddDonee(a1.Address())
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
		{owner, a2.Address(), errDoneeNotExist},
		{owner, a3.Address(), errDoneeNotExist},
	}
	//run
	for _, c := range cases {
		test.run(c.codeExpect, func(t *TestObject) types.Error {
			return t.setSender(c.sender).Donate(c.addr)
		})
	}
	fmt.Println("=== pass")

	fmt.Println("=== test for receipt of transfer")
	//run
	test.run(types.ErrInvalidParameter, func(t *TestObject) types.Error {
		return t.setSender(owner).Donate(a1.Address())
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(owner)
		utest.Assert(t.transfer("TSC", halfCoin) != nil)
		return t.Donate(a1.Address())
	})
	test.run(types.ErrInvalidParameter, func(t *TestObject) types.Error {
		t.setSender(owner)
		utest.Assert(test.transfer("TSC", halfCoin) != nil)
		utest.Assert(test.transfer("TSC", halfCoin) != nil)
		return t.Donate(a1.Address())
	})
	test.run(types.ErrInvalidParameter, func(t *TestObject) types.Error {
		t.setSender(owner)
		utest.Assert(t.transfer("BTC", halfCoin) != nil)
		return t.Donate(a1.Address())
	})
	fmt.Println("=== pass")

	fmt.Println("=== test for business logic")
	//prepare
	utest.AssertSDB(keyOfDonation(a1.Address()), &halfCoin)
	//run
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(a2)
		utest.Assert(t.transfer(halfCoin) != nil)
		return t.Donate(a1.Address())
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(a3)
		utest.Assert(t.transfer(halfCoin) != nil)
		return t.Donate(a1.Address())
	})
	//check
	utest.AssertSDB(keyOfDonation(a1.Address()), &oneHalfCoin)
	fmt.Println("=== pass")
}
