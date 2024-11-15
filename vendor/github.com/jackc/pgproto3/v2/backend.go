package pgproto3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Backend acts as a server for the PostgreSQL wire protocol version 3.
type Backend struct {
	cr ChunkReader
	w  io.Writer

	// Frontend message flyweights
	bind           Bind
	cancelRequest  CancelRequest
	_close         Close
	copyFail       CopyFail
	copyData       CopyData
	copyDone       CopyDone
	describe       Describe
	execute        Execute
	flush          Flush
	functionCall   FunctionCall
	gssEncRequest  GSSEncRequest
	parse          Parse
	query          Query
	sslRequest     SSLRequest
	startupMessage StartupMessage
	sync           Sync
	terminate      Terminate

	bodyLen    int
	msgType    byte
	partialMsg bool
	authType   uint32
}

const (
	minStartupPacketLen = 4     // minStartupPacketLen is a single 32-bit int version or code.
	maxStartupPacketLen = 10000 // maxStartupPacketLen is MAX_STARTUP_PACKET_LENGTH from PG source.
)

// NewBackend creates a new Backend.
func NewBackend(cr ChunkReader, w io.Writer) *Backend {
	return &Backend{cr: cr, w: w}
}

// Send sends a message to the frontend.
func (b *Backend) Send(msg BackendMessage) error {
	buf, err := msg.Encode(nil)
	if err != nil {
		return err
	}

	_, err = b.w.Write(buf)
	return err
}

// ReceiveStartupMessage receives the initial connection message. This method is used of the normal Receive method
// because the initial connection message is "special" and does not include the message type as the first byte. This
// will return either a StartupMessage, SSLRequest, GSSEncRequest, or CancelRequest.
func (b *Backend) ReceiveStartupMessage() (FrontendMessage, error) {
	buf, err := b.cr.Next(4)
	if err != nil {
		return nil, err
	}
	msgSize := int(binary.BigEndian.Uint32(buf) - 4)

	if msgSize < minStartupPacketLen || msgSize > maxStartupPacketLen {
		return nil, fmt.Errorf("invalid length of startup packet: %d", msgSize)
	}

	buf, err = b.cr.Next(msgSize)
	if err != nil {
		return nil, translateEOFtoErrUnexpectedEOF(err)
	}

	code := binary.BigEndian.Uint32(buf)

	switch code {
	case ProtocolVersionNumber:
		err = b.startupMessage.Decode(buf)
		if err != nil {
			return nil, err
		}
		return &b.startupMessage, nil
	case sslRequestNumber:
		err = b.sslRequest.Decode(buf)
		if err != nil {
			return nil, err
		}
		return &b.sslRequest, nil
	case cancelRequestCode:
		err = b.cancelRequest.Decode(buf)
		if err != nil {
			return nil, err
		}
		return &b.cancelRequest, nil
	case gssEncReqNumber:
		err = b.gssEncRequest.Decode(buf)
		if err != nil {
			return nil, err
		}
		return &b.gssEncRequest, nil
	default:
		return nil, fmt.Errorf("unknown startup message code: %d", code)
	}
}

// Receive receives a message from the frontend. The returned message is only valid until the next call to Receive.
func (b *Backend) Receive() (FrontendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, translateEOFtoErrUnexpectedEOF(err)
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
		if b.bodyLen < 0 {
			return nil, errors.New("invalid message with negative body length received")
		}
	}

	var msg FrontendMessage
	switch b.msgType {
	case 'B':
		msg = &b.bind
	case 'C':
		msg = &b._close
	case 'D':
		msg = &b.describe
	case 'E':
		msg = &b.execute
	case 'F':
		msg = &b.functionCall
	case 'f':
		msg = &b.copyFail
	case 'd':
		msg = &b.copyData
	case 'c':
		msg = &b.copyDone
	case 'H':
		msg = &b.flush
	case 'P':
		msg = &b.parse
	case 'p':
		switch b.authType {
		case AuthTypeSASL:
			msg = &SASLInitialResponse{}
		case AuthTypeSASLContinue:
			msg = &SASLResponse{}
		case AuthTypeSASLFinal:
			msg = &SASLResponse{}
		case AuthTypeGSS, AuthTypeGSSCont:
			msg = &GSSResponse{}
		case AuthTypeCleartextPassword, AuthTypeMD5Password:
			fallthrough
		default:
			// to maintain backwards compatability
			msg = &PasswordMessage{}
		}
	case 'Q':
		msg = &b.query
	case 'S':
		msg = &b.sync
	case 'X':
		msg = &b.terminate
	default:
		return nil, fmt.Errorf("unknown message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, translateEOFtoErrUnexpectedEOF(err)
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)
	return msg, err
}

// SetAuthType sets the authentication type in the backend.
// Since multiple message types can start with 'p', SetAuthType allows
// contextual identification of FrontendMessages. For example, in the
// PG message flow documentation for PasswordMessage:
//
//			Byte1('p')
//
//	     Identifies the message as a password response. Note that this is also used for
//			GSSAPI, SSPI and SASL response messages. The exact message type can be deduced from
//			the context.
//
// Since the Frontend does not know about the state of a backend, it is important
// to call SetAuthType() after an authentication request is received by the Frontend.
func (b *Backend) SetAuthType(authType uint32) error {
	switch authType {
	case AuthTypeOk,
		AuthTypeCleartextPassword,
		AuthTypeMD5Password,
		AuthTypeSCMCreds,
		AuthTypeGSS,
		AuthTypeGSSCont,
		AuthTypeSSPI,
		AuthTypeSASL,
		AuthTypeSASLContinue,
		AuthTypeSASLFinal:
		b.authType = authType
	default:
		return fmt.Errorf("authType not recognized: %d", authType)
	}

	return nil
}
