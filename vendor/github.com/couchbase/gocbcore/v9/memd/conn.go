package memd

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

// Conn represents a memcached protocol connection.
type Conn struct {
	stream io.ReadWriter

	headerBuf       []byte
	enabledFeatures map[HelloFeature]bool
}

// NewConn creates a new connection object which can be used to perform
// reading and writing of packets.
func NewConn(stream io.ReadWriter) *Conn {
	return &Conn{
		stream:          stream,
		headerBuf:       make([]byte, 24),
		enabledFeatures: make(map[HelloFeature]bool),
	}
}

// EnableFeature enables a particular feature on this connection.
func (c *Conn) EnableFeature(feature HelloFeature) {
	c.enabledFeatures[feature] = true
}

// IsFeatureEnabled indicates whether a particular feature is enabled
// on this particular connection.  Note that this is directly based on
// calls to EnableFeature and is not controlled by the library.
func (c *Conn) IsFeatureEnabled(feature HelloFeature) bool {
	if enabled, ok := c.enabledFeatures[feature]; ok {
		return enabled
	}
	return false
}

// WritePacket writes a packet to the network.
func (c *Conn) WritePacket(pkt *Packet) error {
	encodedKey := pkt.Key
	extras := pkt.Extras
	if c.IsFeatureEnabled(FeatureCollections) {
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

	// We automatically upgrade a packet from normal Req or Res magic into
	// the frame variant depending on the usage of them.
	pktMagic := pkt.Magic
	if framesLen > 0 {
		if pktMagic == CmdMagicReq {
			if !c.IsFeatureEnabled(FeatureAltRequests) {
				return errors.New("cannot use frames in req packets without enabling the feature")
			}

			pktMagic = cmdMagicReqExt
		} else if pktMagic == CmdMagicRes {
			pktMagic = cmdMagicResExt
		} else {
			return errors.New("cannot use frames with an unsupported magic")
		}
	}

	// Go appears to do some clever things in regards to writing data
	//   to the kernel for network dispatch.  Having a write buffer
	//   per-server that is re-used actually hinders performance...
	// For now, we will simply create a new buffer and let it be GC'd.
	buffer := make([]byte, 24+keyLen+extLen+valLen+framesLen)

	buffer[0] = uint8(pktMagic)
	buffer[1] = uint8(pkt.Command)

	// This is safe to do without checking the magic as we check the magic
	// above before incrementing the framesLen variable
	if framesLen > 0 {
		buffer[2] = uint8(framesLen)
		buffer[3] = uint8(keyLen)
	} else {
		binary.BigEndian.PutUint16(buffer[2:], uint16(keyLen))
	}
	buffer[4] = byte(extLen)
	buffer[5] = pkt.Datatype

	if pkt.Magic == CmdMagicReq {
		if pkt.Status != 0 {
			return errors.New("cannot specify status in a request packet")
		}

		binary.BigEndian.PutUint16(buffer[6:], pkt.Vbucket)
	} else if pkt.Magic == CmdMagicRes {
		if pkt.Vbucket != 0 {
			return errors.New("cannot specify vbucket in a response packet")
		}

		binary.BigEndian.PutUint16(buffer[6:], uint16(pkt.Status))
	} else {
		return errors.New("cannot encode status/vbucket for unknown packet magic")
	}

	binary.BigEndian.PutUint32(buffer[8:], uint32(len(buffer)-24))
	binary.BigEndian.PutUint32(buffer[12:], pkt.Opaque)
	binary.BigEndian.PutUint64(buffer[16:], pkt.Cas)

	bodyPos := 24

	// Generate the framing extra data

	// makeFrameHeader will take a FrameType and len and then encode it into a 4:4 bit
	// frame header.  Note that this does not account for sizing overruns as this is meant
	// to be done by the specific commands.
	makeFrameHeader := func(ftype frameType, len uint8) uint8 {
		ftypeNum := uint8(ftype)
		return (ftypeNum << 4) | (len << 0)
	}

	if pkt.BarrierFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use barrier frame in non-request packets")
		}

		buffer[bodyPos] = makeFrameHeader(frameTypeReqBarrier, 0)
		bodyPos++
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
			buffer[bodyPos+0] = makeFrameHeader(frameTypeReqSyncDurability, 1)
			buffer[bodyPos+1] = uint8(pkt.DurabilityLevelFrame.DurabilityLevel)
			bodyPos += 2
		} else {
			durabilityTimeoutMillis := pkt.DurabilityTimeoutFrame.DurabilityTimeout / time.Millisecond
			if durabilityTimeoutMillis > 65535 {
				durabilityTimeoutMillis = 65535
			}

			buffer[bodyPos+0] = makeFrameHeader(frameTypeReqSyncDurability, 3)
			buffer[bodyPos+1] = uint8(pkt.DurabilityLevelFrame.DurabilityLevel)
			binary.BigEndian.PutUint16(buffer[bodyPos+2:], uint16(durabilityTimeoutMillis))
			bodyPos += 4
		}
	}
	if pkt.StreamIDFrame != nil {
		if pkt.Magic != CmdMagicReq {
			return errors.New("cannot use stream id frame in non-request packets")
		}

		buffer[bodyPos+0] = makeFrameHeader(frameTypeReqStreamID, 2)
		binary.BigEndian.PutUint16(buffer[bodyPos+1:], pkt.StreamIDFrame.StreamID)
		bodyPos += 3
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
			buffer[bodyPos+0] = makeFrameHeader(frameTypeReqOpenTracing, uint8(traceCtxLen))
			copy(buffer[bodyPos+1:], pkt.OpenTracingFrame.TraceContext)
			bodyPos += 1 + traceCtxLen
		} else {
			buffer[bodyPos+0] = makeFrameHeader(frameTypeReqOpenTracing, 15)
			buffer[bodyPos+1] = uint8(traceCtxLen - 15)
			copy(buffer[bodyPos+2:], pkt.OpenTracingFrame.TraceContext)
			bodyPos += 2 + traceCtxLen
		}
	}

	if pkt.ServerDurationFrame != nil {
		if pkt.Magic != CmdMagicRes {
			return errors.New("cannot use server duration frame in non-response packets")
		}
		if !c.IsFeatureEnabled(FeatureDurations) {
			return errors.New("cannot use server duration frames without enabling the feature")
		}

		serverDurationEnc := EncodeSrvDura16(pkt.ServerDurationFrame.ServerDuration)

		buffer[bodyPos+0] = makeFrameHeader(frameTypeResSrvDuration, 2)
		binary.BigEndian.PutUint16(buffer[bodyPos+1:], serverDurationEnc)
		bodyPos += 3
	}

	if len(pkt.UnsupportedFrames) > 0 {
		return errors.New("cannot send packets with unsupported frames")
	}

	// Copy the extras into the body of the packet
	copy(buffer[bodyPos:], extras)
	bodyPos += len(extras)

	// Copy the encoded key into the body of the packet
	copy(buffer[bodyPos:], encodedKey)
	bodyPos += len(encodedKey)

	// Copy the value into the body of the packet
	copy(buffer[bodyPos:], pkt.Value)

	bytesWritten, err := c.stream.Write(buffer)
	if err != nil {
		return err
	}
	if bytesWritten != len(buffer) {
		return io.ErrShortWrite
	}

	return nil
}

// ReadPacket reads a packet from the network.
func (c *Conn) ReadPacket() (*Packet, int, error) {
	var pkt Packet

	// We use a single byte blob to read all headers to avoid allocating a bunch
	// of identical buffers when we only need one
	headerBuf := c.headerBuf

	// Read the entire 24-byte header first
	_, err := io.ReadFull(c.stream, headerBuf)
	if err != nil {
		return nil, 0, err
	}

	// Grab the length of the full body
	bodyLen := binary.BigEndian.Uint32(headerBuf[8:])

	// Read the remaining bytes of the body
	bodyBuf := make([]byte, bodyLen)
	_, err = io.ReadFull(c.stream, bodyBuf)
	if err != nil {
		return nil, 0, err
	}

	pktMagic := CmdMagic(headerBuf[0])
	if pktMagic == cmdMagicReqExt {
		pkt.Magic = CmdMagicReq
	} else if pktMagic == cmdMagicResExt {
		pkt.Magic = CmdMagicRes
	} else {
		pkt.Magic = pktMagic
	}

	pkt.Command = CmdCode(headerBuf[1])
	pkt.Datatype = headerBuf[5]
	pkt.Opaque = binary.BigEndian.Uint32(headerBuf[12:])
	pkt.Cas = binary.BigEndian.Uint64(headerBuf[16:])

	if pktMagic == CmdMagicReq || pktMagic == cmdMagicReqExt {
		pkt.Vbucket = binary.BigEndian.Uint16(headerBuf[6:])
	} else if pktMagic == CmdMagicRes || pktMagic == cmdMagicResExt {
		pkt.Status = StatusCode(binary.BigEndian.Uint16(headerBuf[6:]))
	} else {
		return nil, 0, errors.New("cannot decode status/vbucket for unknown packet magic")
	}

	extLen := int(headerBuf[4])
	keyLen := 0
	framesLen := 0
	if pktMagic == cmdMagicReqExt || pktMagic == cmdMagicResExt {
		framesLen = int(headerBuf[2])
		keyLen = int(headerBuf[3])
	} else {
		keyLen = int(binary.BigEndian.Uint16(headerBuf[2:]))
	}

	bodyPos := 0

	if framesLen > 0 {
		framesBuf := bodyBuf[bodyPos : bodyPos+framesLen]
		framePos := 0
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

			if pktMagic == cmdMagicReqExt {
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
			} else if pktMagic == cmdMagicResExt {
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
			} else {
				return nil, 0, errors.New("got unexpected magic when decoding frames")
			}
		}

		bodyPos += framesLen
	}

	pkt.Extras = bodyBuf[bodyPos : bodyPos+extLen]
	bodyPos += extLen

	keyVal := bodyBuf[bodyPos : bodyPos+keyLen]
	bodyPos += keyLen
	if c.IsFeatureEnabled(FeatureCollections) {
		if pkt.Command == CmdObserve {
			// While it's possible that the Observe operation is in fact supported with collections
			// enabled, we don't currently implement that operation for simplicity, as the key is
			// actually hidden away in the value data instead of the usual key data.
			return nil, 0, errors.New("the observe operation is not supported with collections enabled")
		}

		if IsCommandCollectionEncoded(pkt.Command) && keyLen > 0 {
			collectionID, idLen, err := DecodeULEB128_32(keyVal)
			if err != nil {
				return nil, 0, err
			}

			keyVal = keyVal[idLen:]
			pkt.CollectionID = collectionID
		}
	}
	pkt.Key = keyVal

	pkt.Value = bodyBuf[bodyPos:]

	return &pkt, 24 + int(bodyLen), nil
}
