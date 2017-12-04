package radius

import (
	"strconv"
)

// Code defines the RADIUS packet type.
type Code int

// Standard RADIUS packet codes.
const (
	CodeAccessRequest      Code = 1
	CodeAccessAccept       Code = 2
	CodeAccessReject       Code = 3
	CodeAccountingRequest  Code = 4
	CodeAccountingResponse Code = 5
	CodeAccessChallenge    Code = 11
	CodeStatusServer       Code = 12
	CodeStatusClient       Code = 13
	CodeDisconnectRequest  Code = 40
	CodeDisconnectACK      Code = 41
	CodeDisconnectNAK      Code = 42
	CodeCoARequest         Code = 43
	CodeCoAACK             Code = 44
	CodeCoANAK             Code = 45
	CodeReserved           Code = 255
)

// String returns a string representation of the code.
func (c Code) String() string {
	switch c {
	case CodeAccessRequest:
		return `Access-Request`
	case CodeAccessAccept:
		return `Access-Accept`
	case CodeAccessReject:
		return `Access-Reject`
	case CodeAccountingRequest:
		return `Accounting-Request`
	case CodeAccountingResponse:
		return `Accounting-Response`
	case CodeAccessChallenge:
		return `Access-Challenge`
	case CodeStatusServer:
		return `Status-Server`
	case CodeStatusClient:
		return `Status-Client`
	case CodeDisconnectRequest:
		return `Disconnect-Request`
	case CodeDisconnectACK:
		return `Disconnect-ACK`
	case CodeDisconnectNAK:
		return `Disconnect-NAK`
	case CodeCoARequest:
		return `CoA-Request`
	case CodeCoAACK:
		return `CoA-ACK`
	case CodeCoANAK:
		return `CoA-NAK`
	case CodeReserved:
		return `Reserved`
	}
	return "Code(" + strconv.Itoa(int(c)) + ")"
}
