package gocbcoreps

type ConnState uint8

const (
	ConnStateOffline ConnState = iota
	ConnStateDegraded
	ConnStateOnline
)
