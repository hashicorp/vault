package memd

import "time"

// BarrierFrame is used to signal to the server that this command should be
// barriered and must not be executed concurrently with other commands.
type BarrierFrame struct {
	// Barrier frames have no additional configuration, but their existence
	// triggers the barriering behaviour.
}

// DurabilityLevelFrame allows you to specify a durability level for an
// operation through the frame extras.
type DurabilityLevelFrame struct {
	DurabilityLevel DurabilityLevel
}

// DurabilityTimeoutFrame allows you to specify a specific timeout for
// durability operations to timeout.  Note that this frame is actually
// an extension of DurabilityLevelFrame and requires that frame to also
// be used in order to function.
type DurabilityTimeoutFrame struct {
	DurabilityTimeout time.Duration
}

// StreamIDFrame provides information about which stream this particular
// operation is related to (used for DCP streams).
type StreamIDFrame struct {
	StreamID uint16
}

// OpenTracingFrame allows open tracing context information to be included
// along with a command which is being performed.
type OpenTracingFrame struct {
	TraceContext []byte
}

// ServerDurationFrame allows the server to return information about the
// period of time an operation took to complete.
type ServerDurationFrame struct {
	ServerDuration time.Duration
}

// UnsupportedFrame is used to include an unsupported frame type in the
// packet data to enable further processing if needed.
type UnsupportedFrame struct {
	Type frameType
	Data []byte
}

// Packet represents a single request or response packet being exchanged
// between two clients.
type Packet struct {
	Magic        CmdMagic
	Command      CmdCode
	Datatype     uint8
	Status       StatusCode
	Vbucket      uint16
	Opaque       uint32
	Cas          uint64
	CollectionID uint32
	Key          []byte
	Extras       []byte
	Value        []byte

	BarrierFrame           *BarrierFrame
	DurabilityLevelFrame   *DurabilityLevelFrame
	DurabilityTimeoutFrame *DurabilityTimeoutFrame
	StreamIDFrame          *StreamIDFrame
	OpenTracingFrame       *OpenTracingFrame
	ServerDurationFrame    *ServerDurationFrame
	UnsupportedFrames      []UnsupportedFrame
}
