package memd

import (
	"bytes"
	"fmt"
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

// ReadUnitsFrame allows the server to return information about the
// number of read units used by a command.
type ReadUnitsFrame struct {
	ReadUnits uint16
}

// WriteUnitsFrame allows the server to return information about the
// number of write units used by a command.
type WriteUnitsFrame struct {
	WriteUnits uint16
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
	ReadUnitsFrame         *ReadUnitsFrame
	WriteUnitsFrame        *WriteUnitsFrame
	UnsupportedFrames      []UnsupportedFrame
}

func (pak *Packet) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(
		&buffer,
		"memd.Packet{Magic:%#02x(%s), Command:%#02x(%s), Datatype:%#02x, Status:%#04x(%s), Vbucket:%d(%#04x), Opaque:%#08x, "+
			"Cas: %#08x, CollectionID:%d(%#08x), Barrier:%t\nKey:\n%sValue:\n%sExtras:\n%s",
		uint8(pak.Magic),
		pak.Magic,
		pak.Command,
		pak.Command.Name(),
		pak.Datatype,
		uint16(pak.Status),
		pak.Status,
		pak.Vbucket,
		pak.Vbucket,
		pak.Opaque,
		pak.Cas,
		pak.CollectionID,
		pak.CollectionID,
		pak.BarrierFrame != nil,
		bytesToHexAsciiString(pak.Key),
		bytesToHexAsciiString(pak.Value),
		bytesToHexAsciiString(pak.Extras),
	)

	if pak.DurabilityLevelFrame != nil {
		fmt.Fprintf(&buffer, "\nDurability Level: %#02x", pak.DurabilityLevelFrame.DurabilityLevel)

		if pak.DurabilityTimeoutFrame != nil {
			fmt.Fprintf(&buffer, "\nDurability Level Timeout: %s", pak.DurabilityTimeoutFrame.DurabilityTimeout)
		}
	}

	if pak.StreamIDFrame != nil {
		fmt.Fprintf(&buffer, "\nStreamID: %#02x", pak.StreamIDFrame.StreamID)
	}

	if pak.OpenTracingFrame != nil {
		fmt.Fprintf(&buffer, "\nTrace Context:\n%s", bytesToHexAsciiString(pak.OpenTracingFrame.TraceContext))
	}

	if pak.ServerDurationFrame != nil {
		fmt.Fprintf(&buffer, "\nServer Duration: %s", pak.ServerDurationFrame.ServerDuration)
	}

	if pak.UserImpersonationFrame != nil {
		fmt.Fprintf(&buffer, "\nUser: %s", string(pak.UserImpersonationFrame.User))
	}

	if pak.PreserveExpiryFrame != nil {
		fmt.Fprintf(&buffer, "\nPreserve Expiry: true")
	}

	if len(pak.UnsupportedFrames) > 0 {
		fmt.Fprintf(&buffer, "\nUnsupported frames:")
		for _, frame := range pak.UnsupportedFrames {
			fmt.Fprintf(&buffer, "\nFrame type: %02x, data: %s", frame.Type, bytesToHexAsciiString(frame.Data))
		}
	}

	fmt.Fprintf(&buffer, "}")

	return buffer.String()
}

func bytesToHexAsciiString(bytes []byte) string {
	out := ""
	var ascii [16]byte
	n := (len(bytes) + 15) &^ 15
	for i := 0; i < n; i++ {
		// include the line numbering at beginning of every line
		if i%16 == 0 {
			out += fmt.Sprintf("%4d", i)
		}

		// extra space between blocks of 8 bytes
		if i%8 == 0 {
			out += " "
		}

		// if we have bytes left, print the hex
		if i < len(bytes) {
			out += fmt.Sprintf(" %02X", bytes[i])
		} else {
			out += "   "
		}

		// build the ascii
		if i >= len(bytes) {
			ascii[i%16] = ' '
		} else if bytes[i] < 32 || bytes[i] > 126 {
			ascii[i%16] = '.'
		} else {
			ascii[i%16] = bytes[i]
		}

		// at the end of the line, print the newline.
		if i%16 == 15 {
			out += fmt.Sprintf("  %s\n", string(ascii[:]))
		}
	}
	return out
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
