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

//Transfer This is a method of MySuite
func (mysuit *MySuite) TestDemo_Transfer(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)

	gls.Mgr.SetValues(gls.Values{gls.SDKKey: utest.UTP.ISmartContract}, func() {
		test := NewTestObject(contractOwner)
		test.setSender(contractOwner).InitChain()
		mysuit.test_Transfer(contractOwner, test)
	})
}

func (mysuit *MySuite) test_Transfer(owner sdk.IAccount, test *TestObject) {
	fmt.Println("=== Run  UnitTest case: Transfer(donee types.Address, value bn.Number)")

	fmt.Println("=== prepare for test")
	zero := bn.N(0)
	oneCoin := bn.N(1000000000)
	halfCoin := bn.N(500000000)
	utest.Transfer(nil, owner.Address(), bn.N(2).Mul(oneCoin))
	a1 := utest.NewAccount("TSC", oneCoin)
	a2 := utest.NewAccount("TSC", oneCoin)
	a3 := utest.NewAccount("TSC", oneCoin)
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).AddDonee(a1.Address())
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(owner)
		utest.Assert(test.transfer(oneCoin) != nil)
		return t.Donate(a1.Address())
	})
	utest.AssertOK(test.AddDonee(a2.Address()))
	fmt.Println("=== pass")

	fmt.Println("=== test for authorization")
	//run
	test.run(types.ErrNoAuthorization, func(t *TestObject) types.Error {
		t.setSender(a2)
		utest.Assert(test.transfer(oneCoin) != nil)
		return t.Transfer(a1.Address(), halfCoin)
	})
	fmt.Println("=== pass")

	fmt.Println("=== test for parameters")
	//prepare
	var cases = []struct {
		sender     sdk.IAccount
		addr       types.Address
		amount     bn.Number
		codeExpect uint32
	}{
		{owner, "", halfCoin, types.ErrInvalidAddress},
		{owner, "local", halfCoin, types.ErrInvalidAddress},
		{owner, "localhshskhjkshfsswtsysyst6t76ddsg7s7w", halfCoin, types.ErrInvalidAddress},
		{owner, owner.Address(), halfCoin, errDoneeNotExist},
		{owner, utest.GetContract().Address(), halfCoin, errDoneeNotExist},
		{owner, utest.GetContract().Account().Address(), halfCoin, errDoneeNotExist},
		{owner, a3.Address(), halfCoin, errDoneeNotExist},
		{owner, a1.Address(), bn.N(-1), types.ErrInvalidParameter},
		{owner, a1.Address(), bn.N(0), types.ErrInvalidParameter},
		{owner, a1.Address(), bn.N(1), types.CodeOK},
		{owner, a1.Address(), oneCoin, errDonationNotEnough},
	}
	//run
	for _, c := range cases {
		test.run(c.codeExpect, func(t *TestObject) types.Error {
			return t.setSender(c.sender).Transfer(c.addr, c.amount)
		})
	}
	fmt.Println("=== pass")

	fmt.Println("=== test for business logic")
	//prepare
	x := oneCoin.SubI(1)
	utest.AssertSDB(keyOfDonation(a1.Address()), &x)
	//run
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		t.setSender(owner)
		utest.Assert(test.transfer(bn.N(1)) != nil)
		return t.Donate(a1.Address())
	})
	//check
	utest.AssertSDB(keyOfDonation(a1.Address()), &oneCoin)
	//run
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).Transfer(a1.Address(), halfCoin)
	})
	test.run(types.CodeOK, func(t *TestObject) types.Error {
		return t.setSender(owner).Transfer(a1.Address(), halfCoin)
	})
	//check
	utest.AssertSDB(keyOfDonation(a1.Address()), &zero)
	utest.AssertBalance(a1, "TSC", bn.N(2).Mul(oneCoin).AddI(1))
	fmt.Println("=== pass")
}
