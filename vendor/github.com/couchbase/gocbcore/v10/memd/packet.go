package memd

import (
	"sync"
	"time"
)

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

// UserImpersonationFrame is used to indicate a user to impersonate.
// Internal: This should never be used and is not supported.
type UserImpersonationFrame struct {
	User []byte
}

// PreserveExpiryFrame is used to indicate that the server should preserve the
// expiry time for existing document.
type PreserveExpiryFrame struct {
	// Preserve Expiry frames have no extra configuration, but their existence
	// triggers the preserve expiry behaviour.
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
	UserImpersonationFrame *UserImpersonationFrame
	PreserveExpiryFrame    *PreserveExpiryFrame
	UnsupportedFrames      []UnsupportedFrame
}

// packetPool - Thread safe pool containing memcached packet structures. Used by the memcached connection when reading
// packets from the TCP socket.
var packetPool = sync.Pool{
	New: func() interface{} {
		return &Packet{}
	},
}

// AcquirePacket - Retrieve a packet from the internal pool. Note that the packet should be returned to the pool to
// avoid unnecessary allocations.
func AcquirePacket() *Packet {
	return packetPool.Get().(*Packet)
}

// ReleasePacket - Return a packet to the internal pool. Note that the packet will be reset, removing any active
// pointers to existing data structures.
func ReleasePacket(packet *Packet) {
	*packet = Packet{}
	packetPool.Put(packet)
}
