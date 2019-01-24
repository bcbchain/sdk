package myballot

import (
	"blockchain/smcsdk/utest"
	"gopkg.in/check.v1"
	"testing"
)

//Test: This is a function
func Test(t *testing.T) { check.TestingT(t) }

//MySuite: This is a struct
type MySuite struct{}

var _ = check.Suite(&MySuite{})

//TestBallot_Init: This is a method of MySuite
func (mysuit *MySuite) TestBallot_Init(c *check.C) {
	utest.Init(orgId)
	contractOwner := utest.DeployContract(c, contractName, orgId, contractMethods, nil)
	test := NewTestObject(contractOwner)

	//TODO
	test.obj.Init([]string{})
}

//TestBallot_GiveRightToVote: This is a method of MySuite
func (mysuit *MySuite) TestBallot_GiveRightToVote(c *check.C) {
	utest.Init(orgId)
	contractOwner := utest.DeployContract(c, contractName, orgId, contractMethods, nil)
	test := NewTestObject(contractOwner)

	//TODO
	test.obj.GiveRightToVote("")
}

//TestBallot_Delegate: This is a method of MySuite
func (mysuit *MySuite) TestBallot_Delegate(c *check.C) {
	utest.Init(orgId)
	contractOwner := utest.DeployContract(c, contractName, orgId, contractMethods, nil)
	test := NewTestObject(contractOwner)

	//TODO
	test.obj.Delegate("")
}

//TestBallot_Vote: This is a method of MySuite
func (mysuit *MySuite) TestBallot_Vote(c *check.C) {
	utest.Init(orgId)
	contractOwner := utest.DeployContract(c, contractName, orgId, contractMethods, nil)
	test := NewTestObject(contractOwner)

	//TODO
	test.obj.Vote(1)
}

//TestBallot_WinningProposal: This is a method of MySuite
func (mysuit *MySuite) TestBallot_WinningProposal(c *check.C) {
	utest.Init(orgId)
	contractOwner := utest.DeployContract(c, contractName, orgId, contractMethods, nil)
	test := NewTestObject(contractOwner)

	//TODO
	test.obj.WinningProposal()
}

//TestBallot_WinnerName: This is a method of MySuite
func (mysuit *MySuite) TestBallot_WinnerName(c *check.C) {
	utest.Init(orgId)
	contractOwner := utest.DeployContract(c, contractName, orgId, contractMethods, nil)
	test := NewTestObject(contractOwner)

	//TODO
	test.obj.WinnerName()
}
