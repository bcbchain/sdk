package myballot

import (
	"blockchain/smcsdk/sdk/types"
)

// Voter this declares a new complex type which will
//       be used for variables later.
//       it will represent a single voter.
type Voter struct {
	weight   uint // weight is accumulated by delegation
	voted    bool // if true, that person already voted
	delegate types.Address // person delegated to
	vote     uint // index of the voted proposal
}

//Proposal this is a type for a single proposal.
type Proposal struct {
	name      string // short name (up to 32 bytes)
	voteCount uint // number of accumulated votes
}
