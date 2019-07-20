package myballot

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/forx"
	"blockchain/smcsdk/sdk/types"
)

//Ballot a demo smart contract for voting with delegation.
//@:contract:myballot
//@:version:1.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:b37e7627431feb18123b81bcf1f41ffd37efdb90513d48ff2c7f8a0c27a9d06c
type Ballot struct {
	sdk sdk.ISmartContract

	//chairperson this declares a state variable that stores a chairperson's
	//            address for the contract
	//@:public:store:cache
	chairperson string

	//voters this declares a state variable that stores a 'Voter' struct for
	//       each possible address
	//@:public:store
	voters map[types.Address]Voter

	//proposals a dynamically-sized array of 'Proposal' structs
	//@:public:store:cache
	proposals []Proposal
}

//InitChain init when deployed on the blockChain first time
//@:constructor
func (ballot *Ballot) InitChain() {
}

//Init create a new (voting) vote for each proposal in 'proposal Names'
//@:public:method:gas[500]
func (ballot *Ballot) Init(proposalNames []string) {
	sender := ballot.sdk.Message().Sender().Address()

	// Only cntract's owner can perform init
	sdk.RequireOwner()

	proposals := ballot._proposals()
	sdk.Require(len(proposals) <= 0,
		types.ErrUserDefined, "Already inited")

	chairperson := sender
	ballot._setChairperson(chairperson)

	voter := ballot._voters(chairperson)
	voter.weight = 1
	ballot._setVoters(chairperson, voter)

	// For each of the provided proposal names,
	// create a new 'Proposal' object and add it to the end of the array
	forx.Range(proposalNames, func(i int, pName string) {
		proposals = append(proposals,
			Proposal{
				name:      pName,
				voteCount: 0,
			})
	})
	ballot._setProposals(proposals)
}

//GiveRightToVote give `voter` the right to vote on this ballot.
//                may only be called by 'chairperson'.
//@:public:method:gas[500]
func (ballot *Ballot) GiveRightToVote(voterAddr types.Address) {
	sender := ballot.sdk.Message().Sender().Address()
	chairperson := ballot._chairperson()
	sdk.Require(sender == chairperson,
		types.ErrNoAuthorization, "Only chairperson can give right to vote.")

	voter := ballot._voters(voterAddr)
	sdk.Require(voter.voted == false,
		types.ErrUserDefined, "The voter already voted.")
	sdk.Require(voter.weight == 0,
		types.ErrUserDefined, "The voter's weight must be zero.")

	voter.weight = 1
	ballot._setVoters(voterAddr, voter)
}

//Delegate Delegate your vote to the voter 'to'
//@:public:method:gas[1500]
func (ballot *Ballot) Delegate(to types.Address) {
	sender := ballot.sdk.Message().Sender().Address()
	sendVoter := ballot._voters(sender)

	sdk.Require(sendVoter.voted == false,
		types.ErrUserDefined, "You already voted.")
	sdk.Require(to != sender,
		types.ErrUserDefined, "Self-delegation is disallowed.")

	// Forward the delegation as long as 'to' also delegated.
	// In general, such loops are very dangerous, because if they run too
	// long, they might need more gas than is available in a block.
	// In this case, the delegation will not be executed, but in other
	// situations, such loops might cause a contract to get "stuck" completely.
	toVoter := ballot._voters(to)
	forx.Range(func() bool {
		return toVoter.delegate != ""
	}, func(i int) {
		to = toVoter.delegate
		toVoter = ballot._voters(to)

		// We found a loop in the delegation, not allowed.
		sdk.Require(to != sender,
			types.ErrUserDefined, "Found loop in delegation.")
	})

	sendVoter.voted = true
	sendVoter.delegate = to
	delegate := toVoter
	if delegate.voted {
		// If the delegate already voted,
		// directly add to the number of votes
		proposals := ballot._proposals()
		proposals[int(delegate.vote)].voteCount += sendVoter.weight
		ballot._setProposals(proposals)
	} else {
		// If the delegate did not vote yet,
		// add to her weight.
		delegate.weight += sendVoter.weight
		ballot._setVoters(to, delegate)
	}
	return
}

//Vote give your vote (including votes delegated to you)
//     to proposal `proposals[proposal].name`.
//@:public:method:gas[500]
func (ballot *Ballot) Vote(proposal uint) {
	sender := ballot.sdk.Message().Sender().Address()
	sendVoter := ballot._voters(sender)

	sdk.Require(sendVoter.voted == false,
		types.ErrUserDefined, "You already voted.")

	proposals := ballot._proposals()
	sdk.Require(proposal < uint(len(proposals)),
		types.ErrUserDefined, "Proposal is out of index.")

	sendVoter.voted = true
	sendVoter.vote = proposal
	proposals[int(proposal)].voteCount += sendVoter.weight
	ballot._setProposals(proposals)
}

//WinningProposal computes the winning proposal taking all
//                previous votes into account.
//@:public:method:gas[500]
func (ballot *Ballot) WinningProposal() (winningProposal uint) {
	var winningVoteCount uint

	proposals := ballot._proposals()
	forx.Range(proposals, func(i int, proposal Proposal) {
		if proposal.voteCount > winningVoteCount {
			winningVoteCount = proposal.voteCount
			winningProposal = uint(i)
		}
	})
	return
}

//WinnerName calls winningProposal() function to get the index
//           of the winner contained in the proposals array and then
//           returns the name of the winner
//@:public:method:gas[500]
func (ballot *Ballot) WinnerName() (winnerName string) {
	proposals := ballot._proposals()
	if len(proposals) > 0 {
		winnerName = proposals[ballot.WinningProposal()].name
	}
	return
}
