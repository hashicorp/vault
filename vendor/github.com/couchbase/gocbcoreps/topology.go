package gocbcoreps

type Node struct {
	NodeID      string
	ServerGroup string
}

type DataNode struct {
	Node *Node

	LocalVbuckets []uint32
	GroupVbuckets []uint32
}

type VbucketRouting struct {
	Nodes       []*DataNode
	NumVbuckets uint
}

type Topology struct {
	Revision []uint64

	Nodes          []*Node
	VbucketRouting *VbucketRouting
}
