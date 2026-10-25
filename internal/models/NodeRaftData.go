package models

type NodeRaftData struct {
	term        uint64
	index       uint64
	commitIndex uint64
}

func (this *NodeRaftData) GetCommitIndex() uint64 {
	return this.commitIndex
}
func (this *NodeRaftData) IncrementCommitIndex() {
	this.commitIndex = this.commitIndex + 1
}

func (this *NodeRaftData) GetTerm() uint64 {
	return this.term
}

func (this *NodeRaftData) GetIndex() uint64 {
	return this.index
}

func (this *NodeRaftData) IncrementTerm() {
	this.index = this.term + 1
}

func (this *NodeRaftData) IncrementIndex() {
	this.index = this.index + 1
}

func (this *NodeRaftData) Init() {
	this.index = 0
	this.term = 0
}
