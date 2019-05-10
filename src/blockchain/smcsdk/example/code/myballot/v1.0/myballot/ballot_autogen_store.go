package myballot

import (
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_chairperson This is a method of Ballot
func (b *Ballot) _chairperson() string {

	return *b.sdk.Helper().StateHelper().McGetEx("/chairperson", new(string)).(*string)
}

//_setChairperson This is a method of Ballot
func (b *Ballot) _setChairperson(v string) {
	b.sdk.Helper().StateHelper().McSet("/chairperson", &v)
}

//_chkChairperson This is a method of Ballot
func (b *Ballot) _chkChairperson() bool {
	return b.sdk.Helper().StateHelper().Check("/chairperson")
}

//_clrChairperson This is a method of Ballot
func (b *Ballot) _clrChairperson() {
	b.sdk.Helper().StateHelper().McClear("/chairperson")
}

//_delChairperson This is a method of Ballot
func (b *Ballot) _delChairperson() {
	b.sdk.Helper().StateHelper().Delete("/chairperson")
}

//_voters This is a method of Ballot
func (b *Ballot) _voters(k types.Address) Voter {

	return *b.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/voters/%v", k), new(Voter)).(*Voter)
}

//_setVoters This is a method of Ballot
func (b *Ballot) _setVoters(k types.Address, v Voter) {
	b.sdk.Helper().StateHelper().Set(fmt.Sprintf("/voters/%v", k), &v)
}

//_chkVoters This is a method of Ballot
func (b *Ballot) _chkVoters(k types.Address) bool {
	return b.sdk.Helper().StateHelper().Check(fmt.Sprintf("/voters/%v", k))
}

//_delVoters This is a method of Ballot
func (b *Ballot) _delVoters(k types.Address) {
	b.sdk.Helper().StateHelper().Delete(fmt.Sprintf("/voters/%v", k))
}

//_proposals This is a method of Ballot
func (b *Ballot) _proposals() []Proposal {

	return *b.sdk.Helper().StateHelper().McGetEx("/proposals", new([]Proposal)).(*[]Proposal)
}

//_setProposals This is a method of Ballot
func (b *Ballot) _setProposals(v []Proposal) {
	b.sdk.Helper().StateHelper().McSet("/proposals", &v)
}

//_chkProposals This is a method of Ballot
func (b *Ballot) _chkProposals() bool {
	return b.sdk.Helper().StateHelper().Check("/proposals")
}

//_clrProposals This is a method of Ballot
func (b *Ballot) _clrProposals() {
	b.sdk.Helper().StateHelper().McClear("/proposals")
}

//_delProposals This is a method of Ballot
func (b *Ballot) _delProposals() {
	b.sdk.Helper().StateHelper().Delete("/proposals")
}
