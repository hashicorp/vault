package zk

import "fmt"

// TODO: (v2) enum type for CreateMode API.
const (
	FlagPersistent                  = 0
	FlagEphemeral                   = 1
	FlagSequence                    = 2
	FlagEphemeralSequential         = 3
	FlagContainer                   = 4
	FlagTTL                         = 5
	FlagPersistentSequentialWithTTL = 6
)

type createMode struct {
	flag         int32
	isEphemeral  bool
	isSequential bool
	isContainer  bool
	isTTL        bool
}

// parsing a flag integer into the CreateMode needed to call the correct
// Create RPC to Zookeeper.
//
// NOTE: This parse method is designed to be able to copy and paste the same
// CreateMode ENUM constructors from Java:
// https://github.com/apache/zookeeper/blob/master/zookeeper-server/src/main/java/org/apache/zookeeper/CreateMode.java
func parseCreateMode(flag int32) (createMode, error) {
	switch flag {
	case FlagPersistent:
		return createMode{0, false, false, false, false}, nil
	case FlagEphemeral:
		return createMode{1, true, false, false, false}, nil
	case FlagSequence:
		return createMode{2, false, true, false, false}, nil
	case FlagEphemeralSequential:
		return createMode{3, true, true, false, false}, nil
	case FlagContainer:
		return createMode{4, false, false, true, false}, nil
	case FlagTTL:
		return createMode{5, false, false, false, true}, nil
	case FlagPersistentSequentialWithTTL:
		return createMode{6, false, true, false, true}, nil
	default:
		return createMode{}, fmt.Errorf("invalid flag value: [%v]", flag)
	}
}
