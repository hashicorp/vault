package memd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"
)

// bufPool - Thread safe pool containing packet write buffers i.e. they should be used to write a single packet to the
// TCP socket.
var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0))
	},
}

// aquireWriteBuf - Returns a pointer to a write buffer which is ready to be used, ensure the buffer is released using
// the 'releaseWriteBuf' function.
func aquireWriteBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

// releaseWriteBuf - Reset the buffer so that it's clean for the next user (note that this retains the underlying
// storage for future writes) and then return it to the pool.
func releaseWriteBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

// Conn represents a memcached protocol connection.
type Conn struct {
	reader io.Reader
	writer io.Writer

	headerBuf [24]byte

	collectionsEnabled bool
	enabledFeatures    map[HelloFeature]bool
}

// NewConn creates a new connection object which can be used to perform
// reading and writing of packets.
func NewConn(stream io.ReadWriter) *Conn {
	return &Conn{
		// Use the provided writer directly without buffering writes; if we did so, we may end up writing packets into
		// the buffer which remain in userspace indefinitely.
		writer: stream,

		// Wrap the provided reader with a buffer that's 20MB in size; this is the same size used for connections inside
		// KV engine.
		reader:          bufio.NewReaderSize(stream, 20*1024*1024),
		enabledFeatures: make(map[HelloFeature]bool),
	}
}

// EnableFeature enables a particular feature on this connection.
func (c *Conn) EnableFeature(feature HelloFeature) {
	c.enabledFeatures[feature] = true

	if feature == FeatureCollections {
		c.collectionsEnabled = true
	}
}

// IsFeatureEnabled indicates whether a particular feature is enabled
// on this particular connection.  Note that this is directly based on
// calls to EnableFeature and is not controlled by the library.
func (c *Conn) IsFeatureEnabled(feature HelloFeature) bool {
	enabled, ok := c.enabledFeatures[feature]
	return ok && enabled
}

// WritePacket writes a packet to the network.
func (c *Conn) WritePacket(pkt *Packet) error {
	encodedKey := pkt.Key
	extras := pkt.Extras
	if c.collectionsEnabled {
		if pkt.Command == CmdObserve {
			// While it's possible that the Observe operation is in fact supported with collections
			// enabled, we don't currently implement that operation for simplicity, as the key is
			// actually hidden away in the value data instead of the usual key data.
			return errors.New("the observe operation is not supported with collections enabled")
		}

		if IsCommandCollectionEncoded(pkt.Command) {
			collEncodedKey := make([]byte, 0, len(encodedKey)+5)
			collEncodedKey = AppendULEB128_32(collEncodedKey, pkt.CollectionID)
			collEncodedKey = append(collEncodedKey, encodedKey...)
			encodedKey = collEncodedKey
		} else if pkt.Command == CmdGetRandom {
			// GetRandom expects the cid to be in the extras
			// GetRandom MUST not have any extras if not using collections so we're ok to just set it.
			// It also doesn't expect the collection ID to be leb encoded.
			extras = make([]byte, 4)
			binary.BigEndian.PutUint32(extras, pkt.CollectionID)
		} else {
			if pkt.CollectionID > 0 {
				return errors.New("cannot encode collection id with a non-collection command")
			}
		}
	} else {
		if pkt.CollectionID > 0 {
			return errors.New("cannot encode collection id without the feature enabled")
		}
	}

	extLen := len(extras)
	keyLen := len(encodedKey)
	valLen := len(pkt.Value)

	framesLen := 0
	if pkt.BarrierFrame != nil {
		framesLen++
	}
	if pkt.DurabilityLevelFrame != nil {
		if pkt.DurabilityTimeoutFrame == nil {
			framesLen += 2
		} else {
			framesLen += 4
		}
	}
	if pkt.StreamIDFrame != nil {
		framesLen += 3
	}
	if pkt.OpenTracingFrame != nil {
		traceCtxLen := len(pkt.OpenTracingFrame.TraceContext)
		if traceCtxLen < 15 {
			framesLen += 1 + traceCtxLen
		} else {
			framesLen += 2 + traceCtxLen
		}
	}
	if pkt.ServerDurationFrame != nil {
		framesLen += 3
	}
	if pkt.UserImpersonationFrame != nil {
		userLen := len(pkt.UserImpersonationFrame.User)
		if userLen < 15 {
			framesLen += 1 + userLen
		} else {
			framesLen += 2 + userLen
		}
	}
	if pkt.PreserveExpiryFrame != nil {
		framesLen += 1
	}

	// We automatically upgrade a packet from normal Req or Res magic into
	// the frame variant depending on the usage of them.
	pktMagic := pkt.Magic
	if framesLen > 0 {
		switch pktMagic {
		case CmdMagicReq:
			if !c.IsFeatureEnabled(FeatureAltRequests) {
				return errors.New("cannot use frames in req packets without enabling the feature")
			}

			pktMagic = cmdMagicReqExt
		case CmdMagicRes:
			pktMagic = cmdMagicResExt
		default:
			return errors.New("cannot use frames with an unsupported magic")
		}
	}

	buffer := aquireWriteBuf()
	defer releaseWriteBuf(buffer)

	buffer.WriteByte(byte(pktMagic))
	buffer.WriteByte(byte(pkt.Command))

	// This is safe to do without checking the magic as we check the magic
	// above before incrementing the framesLen variable
	if framesLen > 0 {
		buffer.WriteByte(byte(framesLen))
		buffer.WriteByte(byte(keyLen))
	} else {
		writeUint16(buffer, uint16(keyLen))
	}

	buffer.WriteByte(byte(extLen))
	buffer.WriteByte(pkt.Datatype)

	switch pkt.Magic {
	case CmdMagicReq:
		if pkt.Status != 0 {
			return errors.New("cannot specify status in a request packet")
		}

		writeUint16(buffer, pkt.Vbucket)
	case CmdMagicRes:
		if pkt.Vbucket != 0 {
			return errors.New("cannot specify vbucket in a response packet")
		}

		writeUint16(buffer, uint16(pkt.Status))
	default:
		return errors.New("cannot encode status/vbucket for unknown packet magic")
	}

	writeUint32(buffer, uint32(keyLen+extLen+valLen+framesLen))
	writeUint32(buffer, pkt.Opaque)
	writeUint64(buffer, pkt.Cas)

	// Generate the framing extra data

	if pkt.BarrierFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use barrier frame in non-request packets")
		}

		writeFrameHeader(buffer, frameTypeReqBarrier, 0)
	}

	if pkt.DurabilityLevelFrame != nil || pkt.DurabilityTimeoutFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use durability level frame in non-request packets")
		}

		if !c.IsFeatureEnabled(FeatureSyncReplication) {
			return errors.New("cannot use sync replication frames without enabling the feature")
		}

		if pkt.DurabilityLevelFrame == nil && pkt.DurabilityTimeoutFrame != nil {
			return errors.New("cannot encode durability timeout frame without durability level frame")
		}

		if pkt.DurabilityTimeoutFrame == nil {
			writeFrameHeader(buffer, frameTypeReqSyncDurability, 1)
			buffer.WriteByte(byte(pkt.DurabilityLevelFrame.DurabilityLevel))
		} else {
			durabilityTimeoutMillis := pkt.DurabilityTimeoutFrame.DurabilityTimeout / time.Millisecond
			if durabilityTimeoutMillis > 65535 {
				durabilityTimeoutMillis = 65535
			}

			writeFrameHeader(buffer, frameTypeReqSyncDurability, 3)
			buffer.WriteByte(byte(pkt.DurabilityLevelFrame.DurabilityLevel))
			writeUint16(buffer, uint16(durabilityTimeoutMillis))
		}
	}

	if pkt.StreamIDFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use stream id frame in non-request packets")
		}

		writeFrameHeader(buffer, frameTypeReqStreamID, 2)
		writeUint16(buffer, pkt.StreamIDFrame.StreamID)
	}

	if pkt.OpenTracingFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use open tracing frame in non-request packets")
		}

		if !c.IsFeatureEnabled(FeatureOpenTracing) {
			return errors.New("cannot use open tracing frames without enabling the feature")
		}

		traceCtxLen := len(pkt.OpenTracingFrame.TraceContext)
		if traceCtxLen < 15 {
			writeFrameHeader(buffer, frameTypeReqOpenTracing, uint8(traceCtxLen))
			buffer.Write(pkt.OpenTracingFrame.TraceContext)
		} else {
			writeFrameHeader(buffer, frameTypeReqOpenTracing, 15)
			buffer.WriteByte(uint8(traceCtxLen - 15))
			buffer.Write(pkt.OpenTracingFrame.TraceContext)
		}
	}

	if pkt.ServerDurationFrame != nil {
		if pkt.Magic != CmdMagicRes {
			return errors.New("cannot use server duration frame in non-response packets")
		}

		if !c.IsFeatureEnabled(FeatureDurations) {
			return errors.New("cannot use server duration frames without enabling the feature")
		}

		writeFrameHeader(buffer, frameTypeResSrvDuration, 2)
		writeUint16(buffer, EncodeSrvDura16(pkt.ServerDurationFrame.ServerDuration))
	}

	if pkt.UserImpersonationFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use user impersonation frame in non-request packets")
		}

		userCtxLen := len(pkt.UserImpersonationFrame.User)
		if userCtxLen < 15 {
			writeFrameHeader(buffer, frameTypeReqUserImpersonation, uint8(userCtxLen))
			buffer.Write(pkt.UserImpersonationFrame.User)
		} else {
			writeFrameHeader(buffer, frameTypeReqUserImpersonation, 15)
			buffer.WriteByte(byte(userCtxLen - 15))
			buffer.Write(pkt.UserImpersonationFrame.User)
		}
	}

	if pkt.PreserveExpiryFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use preserve expiry frame in non-request packets")
		}

		if !c.IsFeatureEnabled(FeaturePreserveExpiry) {
			return errors.New("cannot use preserve expiry frames without enabling the feature")
		}

		writeFrameHeader(buffer, frameTypeReqPreserveExpiry, 0)
	}

	if len(pkt.UnsupportedFrames) > 0 {
		return errors.New("cannot send packets with unsupported frames")
	}

	// Copy the extras into the body of the packet
	buffer.Write(extras)

	// Copy the encoded key into the body of the packet
	buffer.Write(encodedKey)

	// Copy the value into the body of the packet
	buffer.Write(pkt.Value)

	n, err := c.writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	if n != buffer.Len() {
		return io.ErrShortWrite
	}

	return nil
}

// ReadPacket reads a packet from the network.
func (c *Conn) ReadPacket() (*Packet, int, error) {
	pkt := AcquirePacket()

	// Read the entire 24-byte header first
	_, err := io.ReadFull(c.reader, c.headerBuf[:])
	if err != nil {
		return nil, 0, err
	}

	// Grab the length of the full body
	bodyLen := binary.BigEndian.Uint32(c.headerBuf[8:])

	// Read the remaining bytes of the body
	bodyBuf := make([]byte, bodyLen)
	_, err = io.ReadFull(c.reader, bodyBuf)
	if err != nil {
		return nil, 0, err
	}

	pktMagic := CmdMagic(c.headerBuf[0])
	switch pktMagic {
	case CmdMagicReq, cmdMagicReqExt:
		pkt.Magic = CmdMagicReq
		pkt.Vbucket = binary.BigEndian.Uint16(c.headerBuf[6:])
	case CmdMagicRes, cmdMagicResExt:
		pkt.Magic = CmdMagicRes
		pkt.Status = StatusCode(binary.BigEndian.Uint16(c.headerBuf[6:]))
	default:
		return nil, 0, errors.New("cannot decode status/vbucket for unknown packet magic")
	}

	pkt.Command = CmdCode(c.headerBuf[1])
	pkt.Datatype = c.headerBuf[5]
	pkt.Opaque = binary.BigEndian.Uint32(c.headerBuf[12:])
	pkt.Cas = binary.BigEndian.Uint64(c.headerBuf[16:])

	var (
		extLen    = int(c.headerBuf[4])
		keyLen    = int(binary.BigEndian.Uint16(c.headerBuf[2:]))
		framesLen int
	)

	if pktMagic == cmdMagicReqExt || pktMagic == cmdMagicResExt {
		framesLen = int(c.headerBuf[2])
		keyLen = int(c.headerBuf[3])
	}

	if framesLen > 0 {
		var (
			framesBuf = bodyBuf[:framesLen]
			framePos  int
		)

		for framePos < framesLen {
			frameHeader := framesBuf[framePos]
			framePos++

			frType := frameType((frameHeader & 0xF0) >> 4)
			if frType == 15 {
				frType = 15 + frameType(framesBuf[framePos])
				framePos++
			}

			frameLen := int((frameHeader & 0x0F) >> 0)
			if frameLen == 15 {
				frameLen = 15 + int(framesBuf[framePos])
				framePos++
			}

			frameBody := framesBuf[framePos : framePos+frameLen]
			framePos += frameLen

			switch pktMagic {
			case cmdMagicReqExt:
				if frType == frameTypeReqBarrier && frameLen == 0 {
					pkt.BarrierFrame = &BarrierFrame{}
				} else if frType == frameTypeReqSyncDurability && (frameLen == 1 || frameLen == 3) {
					pkt.DurabilityLevelFrame = &DurabilityLevelFrame{
						DurabilityLevel: DurabilityLevel(frameBody[0]),
					}
					if frameLen == 3 {
						durabilityTimeoutMillis := binary.BigEndian.Uint16(frameBody[1:])
						pkt.DurabilityTimeoutFrame = &DurabilityTimeoutFrame{
							DurabilityTimeout: time.Duration(durabilityTimeoutMillis) * time.Millisecond,
						}
					} else {
						// We follow the semantic that duplicate frames overwrite previous ones,
						// since the timeout frame is 'virtual' to us, we need to clear it in case
						// this is a duplicate frame.
						pkt.DurabilityTimeoutFrame = nil
					}
				} else if frType == frameTypeReqStreamID && frameLen == 2 {
					pkt.StreamIDFrame = &StreamIDFrame{
						StreamID: binary.BigEndian.Uint16(frameBody),
					}
				} else if frType == frameTypeReqOpenTracing {
					pkt.OpenTracingFrame = &OpenTracingFrame{
						TraceContext: frameBody,
					}
				} else {
					// If we don't understand this frame type, we record it as an
					// UnsupportedFrame (as opposed to dropping it blindly)
					pkt.UnsupportedFrames = append(pkt.UnsupportedFrames, UnsupportedFrame{
						Type: frType,
						Data: frameBody,
					})
				}
			case cmdMagicResExt:
				if frType == frameTypeResSrvDuration && frameLen == 2 {
					serverDurationEnc := binary.BigEndian.Uint16(frameBody)
					pkt.ServerDurationFrame = &ServerDurationFrame{
						ServerDuration: DecodeSrvDura16(serverDurationEnc),
					}
				} else {
					// If we don't understand this frame type, we record it as an
					// UnsupportedFrame (as opposed to dropping it blindly)
					pkt.UnsupportedFrames = append(pkt.UnsupportedFrames, UnsupportedFrame{
						Type: frType,
						Data: frameBody,
					})
				}
			default:
				return nil, 0, errors.New("got unexpected magic when decoding frames")
			}
		}
	}

	pkt.Extras = bodyBuf[framesLen : framesLen+extLen]
	pkt.Key = bodyBuf[framesLen+extLen : framesLen+extLen+keyLen]
	pkt.Value = bodyBuf[framesLen+extLen+keyLen:]

	if c.collectionsEnabled {
		if pkt.Command == CmdObserve {
			// While it's possible that the Observe operation is in fact supported with collections
			// enabled, we don't currently implement that operation for simplicity, as the key is
			// actually hidden away in the value data instead of the usual key data.
			return nil, 0, errors.New("the observe operation is not supported with collections enabled")
		}

		if keyLen > 0 && IsCommandCollectionEncoded(pkt.Command) {
			collectionID, idLen, err := DecodeULEB128_32(pkt.Key)
			if err != nil {
				return nil, 0, err
			}

			pkt.Key = pkt.Key[idLen:]
			pkt.CollectionID = collectionID
		}
	}

	return pkt, 24 + int(bodyLen), nil
}

// writeUint16 - Similar to 'bytes.BigEndian.PutUint16' accept we write directly into the provided buffer.
func writeUint16(buffer *bytes.Buffer, n uint16) {
	buffer.WriteByte(byte(n >> 8))
	buffer.WriteByte(byte(n))
}

// writeUint32 - Similar to 'bytes.BigEndian.PutUint32' accept we write directly into the provided buffer.
func writeUint32(buffer *bytes.Buffer, n uint32) {
	buffer.WriteByte(byte(n >> 24))
	buffer.WriteByte(byte(n >> 16))
	buffer.WriteByte(byte(n >> 8))
	buffer.WriteByte(byte(n))
}

// writeUint64 - Similar to 'bytes.BigEndian.PutUint64' accept we write directly into the provided buffer.
func writeUint64(buffer *bytes.Buffer, n uint64) {
	buffer.WriteByte(byte(n >> 56))
	buffer.WriteByte(byte(n >> 48))
	buffer.WriteByte(byte(n >> 40))
	buffer.WriteByte(byte(n >> 32))
	buffer.WriteByte(byte(n >> 24))
	buffer.WriteByte(byte(n >> 16))
	buffer.WriteByte(byte(n >> 8))
	buffer.WriteByte(byte(n))
}

// writeFrameHeader - Write a single byte containing information about the following frame directly into the provided
// buffer.
func writeFrameHeader(buffer *bytes.Buffer, frameType frameType, frameLen uint8) {
	buffer.WriteByte(uint8(frameType)<<4 | frameLen)
}
