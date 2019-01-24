package myballot

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/utest"
)

var (
	contractName    = "myballot" //contract name
	contractMethods = []string{"Init([]string)types.Error", "GiveRightToVote(types.Address)types.Error", "Delegate(types.Address)types.Error", "Vote(uint)types.Error", "WinningProposal()uint", "WinnerName()string"}
	orgId           = "orgNUjCm1i8RcoW2kVTbDw4vKW6jzfMxewJH"
)

//TestObject: This is a struct for test
type TestObject struct {
	obj *Ballot
}

//NewTestObject: This is a function
func NewTestObject(sender sdk.IAccount) *TestObject {
	return &TestObject{&Ballot{sdk: utest.UTP.ISmartContract}}
}

//transfer: This is a method of TestObject
func (t *TestObject) transfer(balance bn.Number) *TestObject {
	t.obj.sdk.Message().Sender().Transfer(t.obj.sdk.Message().Contract().Account(), balance)
	return t
}

//setSender: This is a method of TestObject
func (t *TestObject) setSender(sender sdk.IAccount) *TestObject {
	t.obj.sdk = utest.SetSender(sender.Address())
	return t
}

//run: This is a method of TestObject
func (t *TestObject) run() *TestObject {
	t.obj.sdk = utest.ResetMsg()
	return t
}

//Init: This is a method of TestObject
func (t *TestObject) Init(proposalNames []string) (error types.Error) {
	utest.NextBlock(1)
	result0 := t.obj.Init(proposalNames)
	if result0.ErrorCode != types.CodeOK {
		//one day rollback
	}
	utest.Commit()
	return result0
}

//GiveRightToVote: This is a method of TestObject
func (t *TestObject) GiveRightToVote(voterAddr types.Address) (error types.Error) {
	utest.NextBlock(1)
	result0 := t.obj.GiveRightToVote(voterAddr)
	if result0.ErrorCode != types.CodeOK {
		//one day rollback
	}
	utest.Commit()
	return result0
}

//Delegate: This is a method of TestObject
func (t *TestObject) Delegate(to types.Address) (error types.Error) {
	utest.NextBlock(1)
	result0 := t.obj.Delegate(to)
	if result0.ErrorCode != types.CodeOK {
		//one day rollback
	}
	utest.Commit()
	return result0
}

//Vote: This is a method of TestObject
func (t *TestObject) Vote(proposal uint) (error types.Error) {
	utest.NextBlock(1)
	result0 := t.obj.Vote(proposal)
	if result0.ErrorCode != types.CodeOK {
		//one day rollback
	}
	utest.Commit()
	return result0
}

//WinningProposal: This is a method of TestObject
func (t *TestObject) WinningProposal() (winningProposal uint) {
	utest.NextBlock(1)
	result0 := t.obj.WinningProposal()

	utest.Commit()
	return result0
}

//WinnerName: This is a method of TestObject
func (t *TestObject) WinnerName() (winnerName string) {
	utest.NextBlock(1)
	result0 := t.obj.WinnerName()

	utest.Commit()
	return result0
}
