package myballot

import (
	"fmt"
)

//_setChairperson: This is a method of Ballot
func (b *Ballot) _setChairperson(v string) {
	b.sdk.Helper().StateHelper().McSet("/chairperson", &v)
}

//_chairperson: This is a method of Ballot
func (b *Ballot) _chairperson() string {

	return *b.sdk.Helper().StateHelper().McGetEx("/chairperson", new(string)).(*string)
}

//_clrChairperson: This is a method of Ballot
func (b *Ballot) _clrChairperson() {
	b.sdk.Helper().StateHelper().McClear("/chairperson")
}

//_chkChairperson: This is a method of Ballot
func (b *Ballot) _chkChairperson() bool {
	return b.sdk.Helper().StateHelper().Check("/chairperson")
}

//_McChkChairperson: This is a method of Ballot
func (b *Ballot) _McChkChairperson() bool {
	return b.sdk.Helper().StateHelper().McCheck("/chairperson")
}

//_setVoters: This is a method of Ballot
func (b *Ballot) _setVoters(k string, v Voter) {
	b.sdk.Helper().StateHelper().Set(fmt.Sprintf("/voters/%v", k), &v)
}

//_voters: This is a method of Ballot
func (b *Ballot) _voters(k string) Voter {

	return *b.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/voters/%v", k), new(Voter)).(*Voter)
}

//_chkVoters: This is a method of Ballot
func (b *Ballot) _chkVoters(k string) bool {
	return b.sdk.Helper().StateHelper().Check(fmt.Sprintf("/voters/%v", k))
}

//_setProposals: This is a method of Ballot
func (b *Ballot) _setProposals(v []Proposal) {
	b.sdk.Helper().StateHelper().McSet("/proposals", &v)
}

//_proposals: This is a method of Ballot
func (b *Ballot) _proposals() []Proposal {

	return *b.sdk.Helper().StateHelper().McGetEx("/proposals", new([]Proposal)).(*[]Proposal)
}

//_clrProposals: This is a method of Ballot
func (b *Ballot) _clrProposals() {
	b.sdk.Helper().StateHelper().McClear("/proposals")
}

//_chkProposals: This is a method of Ballot
func (b *Ballot) _chkProposals() bool {
	return b.sdk.Helper().StateHelper().Check("/proposals")
}

//_McChkProposals: This is a method of Ballot
func (b *Ballot) _McChkProposals() bool {
	return b.sdk.Helper().StateHelper().McCheck("/proposals")
}
