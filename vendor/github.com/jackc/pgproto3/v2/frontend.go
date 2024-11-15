package pgproto3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Frontend acts as a client for the PostgreSQL wire protocol version 3.
type Frontend struct {
	cr ChunkReader
	w  io.Writer

	// Backend message flyweights
	authenticationOk                AuthenticationOk
	authenticationCleartextPassword AuthenticationCleartextPassword
	authenticationMD5Password       AuthenticationMD5Password
	authenticationGSS               AuthenticationGSS
	authenticationGSSContinue       AuthenticationGSSContinue
	authenticationSASL              AuthenticationSASL
	authenticationSASLContinue      AuthenticationSASLContinue
	authenticationSASLFinal         AuthenticationSASLFinal
	backendKeyData                  BackendKeyData
	bindComplete                    BindComplete
	closeComplete                   CloseComplete
	commandComplete                 CommandComplete
	copyBothResponse                CopyBothResponse
	copyData                        CopyData
	copyInResponse                  CopyInResponse
	copyOutResponse                 CopyOutResponse
	copyDone                        CopyDone
	dataRow                         DataRow
	emptyQueryResponse              EmptyQueryResponse
	errorResponse                   ErrorResponse
	functionCallResponse            FunctionCallResponse
	noData                          NoData
	noticeResponse                  NoticeResponse
	notificationResponse            NotificationResponse
	parameterDescription            ParameterDescription
	parameterStatus                 ParameterStatus
	parseComplete                   ParseComplete
	readyForQuery                   ReadyForQuery
	rowDescription                  RowDescription
	portalSuspended                 PortalSuspended

	bodyLen    int
	msgType    byte
	partialMsg bool
	authType   uint32
}

// NewFrontend creates a new Frontend.
func NewFrontend(cr ChunkReader, w io.Writer) *Frontend {
	return &Frontend{cr: cr, w: w}
}

// Send sends a message to the backend.
func (f *Frontend) Send(msg FrontendMessage) error {
	buf, err := msg.Encode(nil)
	if err != nil {
		return err
	}
	_, err = f.w.Write(buf)
	return err
}

func translateEOFtoErrUnexpectedEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

// Receive receives a message from the backend. The returned message is only valid until the next call to Receive.
func (f *Frontend) Receive() (BackendMessage, error) {
	if !f.partialMsg {
		header, err := f.cr.Next(5)
		if err != nil {
			return nil, translateEOFtoErrUnexpectedEOF(err)
		}

		f.msgType = header[0]
		f.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		f.partialMsg = true
		if f.bodyLen < 0 {
			return nil, errors.New("invalid message with negative body length received")
		}
	}

	msgBody, err := f.cr.Next(f.bodyLen)
	if err != nil {
		return nil, translateEOFtoErrUnexpectedEOF(err)
	}

	f.partialMsg = false

	var msg BackendMessage
	switch f.msgType {
	case '1':
		msg = &f.parseComplete
	case '2':
		msg = &f.bindComplete
	case '3':
		msg = &f.closeComplete
	case 'A':
		msg = &f.notificationResponse
	case 'c':
		msg = &f.copyDone
	case 'C':
		msg = &f.commandComplete
	case 'd':
		msg = &f.copyData
	case 'D':
		msg = &f.dataRow
	case 'E':
		msg = &f.errorResponse
	case 'G':
		msg = &f.copyInResponse
	case 'H':
		msg = &f.copyOutResponse
	case 'I':
		msg = &f.emptyQueryResponse
	case 'K':
		msg = &f.backendKeyData
	case 'n':
		msg = &f.noData
	case 'N':
		msg = &f.noticeResponse
	case 'R':
		var err error
		msg, err = f.findAuthenticationMessageType(msgBody)
		if err != nil {
			return nil, err
		}
	case 's':
		msg = &f.portalSuspended
	case 'S':
		msg = &f.parameterStatus
	case 't':
		msg = &f.parameterDescription
	case 'T':
		msg = &f.rowDescription
	case 'V':
		msg = &f.functionCallResponse
	case 'W':
		msg = &f.copyBothResponse
	case 'Z':
		msg = &f.readyForQuery
	default:
		return nil, fmt.Errorf("unknown message type: %c", f.msgType)
	}

	err = msg.Decode(msgBody)
	return msg, err
}

// Authentication message type constants.
// See src/include/libpq/pqcomm.h for all
// constants.
const (
	AuthTypeOk                = 0
	AuthTypeCleartextPassword = 3
	AuthTypeMD5Password       = 5
	AuthTypeSCMCreds          = 6
	AuthTypeGSS               = 7
	AuthTypeGSSCont           = 8
	AuthTypeSSPI              = 9
	AuthTypeSASL              = 10
	AuthTypeSASLContinue      = 11
	AuthTypeSASLFinal         = 12
)

func (f *Frontend) findAuthenticationMessageType(src []byte) (BackendMessage, error) {
	if len(src) < 4 {
		return nil, errors.New("authentication message too short")
	}
	f.authType = binary.BigEndian.Uint32(src[:4])

	switch f.authType {
	case AuthTypeOk:
		return &f.authenticationOk, nil
	case AuthTypeCleartextPassword:
		return &f.authenticationCleartextPassword, nil
	case AuthTypeMD5Password:
		return &f.authenticationMD5Password, nil
	case AuthTypeSCMCreds:
		return nil, errors.New("AuthTypeSCMCreds is unimplemented")
	case AuthTypeGSS:
		return &f.authenticationGSS, nil
	case AuthTypeGSSCont:
		return &f.authenticationGSSContinue, nil
	case AuthTypeSSPI:
		return nil, errors.New("AuthTypeSSPI is unimplemented")
	case AuthTypeSASL:
		return &f.authenticationSASL, nil
	case AuthTypeSASLContinue:
		return &f.authenticationSASLContinue, nil
	case AuthTypeSASLFinal:
		return &f.authenticationSASLFinal, nil
	default:
		return nil, fmt.Errorf("unknown authentication type: %d", f.authType)
	}
}

// GetAuthType returns the authType used in the current state of the frontend.
// See SetAuthType for more information.
func (f *Frontend) GetAuthType() uint32 {
	return f.authType
}
