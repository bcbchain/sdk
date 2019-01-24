package myballot

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/types"
)

//Ballot a demo contract
//@:contract:myballot
//@:version:1.0
//@:organization:orgBtjfCSPCAJ84uQWcpNr74NLMWYm5SXzer
//@:author:b37e7627431feb18123b81bcf1f41ffd37efdb90513d48ff2c7f8a0c27a9d06c
type Ballot struct {
	sdk sdk.ISmartContract

	//chairperson 这声明了一个状态变量，为这个合约存储一个主席地址
	//@:public:store:cache
	chairperson string

	//voters 这声明了一个状态变量，为每个可能的地址存储一个 `Voter`
	//@:public:store
	voters map[types.Address]Voter

	//proposals 一个 `Proposal` 结构类型的动态数组
	//@:public:store:cache
	proposals []Proposal
}

//InitChain init when deployed on the blockChain first time
//@:constructor
func (ballot *Ballot) InitChain() {
}

//Init 为 `proposalNames` 中的每个提案，创建一个新的（投票）表决
//@:public:method:gas[500]
func (ballot *Ballot) Init(proposalNames []string) (error types.Error) {
	sender := ballot.sdk.Message().Sender().Address()

	// Only cntract's owner can perform init
	sdk.RequireOwner(ballot.sdk)

	proposals := ballot._proposals()
	sdk.Require(len(proposals) <= 0,
		types.ErrUserDefined, "Already inited")

	chairperson := sender
	ballot._setChairperson(chairperson)

	voter := ballot._voters(chairperson)
	voter.weight = 1
	ballot._setVoters(chairperson, voter)

	//对于提供的每个提案名称，
	//创建一个新的 Proposal 对象并把它添加到数组的末尾。
	for i := 0; i < len(proposalNames); i++ {
		proposals = append(proposals,
			Proposal{
				name:      proposalNames[i],
				voteCount: 0,
			})
	}
	ballot._setProposals(proposals)

	error.ErrorCode = types.CodeOK
	return
}

//GiveRightToVote 授权 `voterAddr` 对这个（投票）表决进行投票
//               只有 `chairperson` 可以调用该函数。
//@:public:method:gas[500]
func (ballot *Ballot) GiveRightToVote(voterAddr types.Address) (error types.Error) {
	// 若 `sdk.Require` 的第一个参数的计算结果为 `false`，
	// 则终止执行，撤销所有对状态的改动。
	// 你也可以在 require 的第三个参数中提供一个对错误情况的详细解释。
	sender := ballot.sdk.Message().Sender().Address()
	chairperson := ballot._chairperson()
	sdk.Require(sender != chairperson,
		types.ErrNoAuthorization, "Only chairperson can give right to vote.")

	voter := ballot._voters(voterAddr)
	sdk.Require(voter.voted,
		types.ErrUserDefined, "The voter already voted.")

	sdk.Require(voter.weight != 0,
		types.ErrUserDefined, "The voter's weight must be zero.")

	voter.weight = 1
	ballot._setVoters(voterAddr, voter)
	return
}

//Delegate 把你的投票委托到投票者 `to`。
//@:public:method:gas[1500]
func (ballot *Ballot) Delegate(to types.Address) (error types.Error) {
	sender := ballot.sdk.Message().Sender().Address()
	sendVoter := ballot._voters(sender)
	sdk.Require(sendVoter.voted,
		types.ErrUserDefined, "You already voted.")

	sdk.Require(to == sender,
		types.ErrUserDefined, "Self-delegation is disallowed.")

	// 委托是可以传递的，只要被委托者 `to` 也设置了委托。
	// 一般来说，这种循环委托是危险的。因为，如果传递的链条太长，
	// 则可能需消耗的gas要多于区块中剩余的（大于区块设置的gasLimit），
	// 这种情况下，委托不会被执行。
	// 而在另一些情况下，如果形成闭环，则会让合约完全卡住。
	toVoter := ballot._voters(to)
	for toVoter.delegate != "" {
		to = toVoter.delegate
		toVoter = ballot._voters(to)

		// 不允许闭环委托
		sdk.Require(to == sender,
			types.ErrUserDefined, "Found loop in delegation.")
	}

	sendVoter.voted = true
	sendVoter.delegate = to
	delegate := toVoter
	if delegate.voted {
		// 若被委托者已经投过票了，直接增加得票数
		proposals := ballot._proposals()
		proposals[int(delegate.vote)].voteCount += sendVoter.weight
		ballot._setProposals(proposals)
	} else {
		// 若被委托者还没投票，增加委托者的权重
		delegate.weight += sendVoter.weight
		ballot._setVoters(to, delegate)
	}
	return
}

//Vote 把你的票(包括委托给你的票)，
//     投给提案 `proposals[proposal].name`.
//@:public:method:gas[500]
func (ballot *Ballot) Vote(proposal uint) (error types.Error) {
	sender := ballot.sdk.Message().Sender().Address()
	sendVoter := ballot._voters(sender)
	sdk.Require(sendVoter.voted,
		types.ErrUserDefined, "You already voted.")

	proposals := ballot._proposals()
	sdk.Require(proposal >= uint(len(proposals)),
		types.ErrUserDefined, "Proposal is out of index.")

	sendVoter.voted = true
	sendVoter.vote = proposal
	proposals[int(proposal)].voteCount += sendVoter.weight
	ballot._setProposals(proposals)
	return
}

//WinningProposal 结合之前所有的投票，计算出最终胜出的提案
//@:public:method:gas[500]
func (ballot *Ballot) WinningProposal() (winningProposal uint) {
	var winningVoteCount uint
	proposals := ballot._proposals()
	for p := 0; p < len(proposals); p++ {
		if proposals[p].voteCount > winningVoteCount {
			winningVoteCount = proposals[p].voteCount
			winningProposal = uint(p)
		}
	}
	return
}

//WinnerName 调用 WinningProposal() 函数以获取提案数组中获胜者
//           的索引，并以此返回获胜者的名称
//@:public:method:gas[500]
func (ballot *Ballot) WinnerName() (winnerName string) {
	proposals := ballot._proposals()
	if len(proposals) > 0 {
		winnerName = proposals[ballot.WinningProposal()].name
	}
	return
}
