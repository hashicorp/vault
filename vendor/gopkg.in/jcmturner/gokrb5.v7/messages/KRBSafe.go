package messages

import (
	"fmt"
	"time"

	"github.com/jcmturner/gofork/encoding/asn1"
	"gopkg.in/jcmturner/gokrb5.v7/iana/asnAppTag"
	"gopkg.in/jcmturner/gokrb5.v7/iana/msgtype"
	"gopkg.in/jcmturner/gokrb5.v7/krberror"
	"gopkg.in/jcmturner/gokrb5.v7/types"
)

/*
KRB-SAFE        ::= [APPLICATION 20] SEQUENCE {
	pvno            [0] INTEGER (5),
	msg-type        [1] INTEGER (20),
	safe-body       [2] KRB-SAFE-BODY,
	cksum           [3] Checksum
}

KRB-SAFE-BODY   ::= SEQUENCE {
	user-data       [0] OCTET STRING,
	timestamp       [1] KerberosTime OPTIONAL,
	usec            [2] Microseconds OPTIONAL,
	seq-number      [3] UInt32 OPTIONAL,
	s-address       [4] HostAddress,
	r-address       [5] HostAddress OPTIONAL
}
*/

// KRBSafe implements RFC 4120 KRB_SAFE: https://tools.ietf.org/html/rfc4120#section-5.6.1.
type KRBSafe struct {
	PVNO     int            `asn1:"explicit,tag:0"`
	MsgType  int            `asn1:"explicit,tag:1"`
	SafeBody KRBSafeBody    `asn1:"explicit,tag:2"`
	Cksum    types.Checksum `asn1:"explicit,tag:3"`
}

// KRBSafeBody implements the KRB_SAFE_BODY of KRB_SAFE.
type KRBSafeBody struct {
	UserData       []byte            `asn1:"explicit,tag:0"`
	Timestamp      time.Time         `asn1:"generalized,optional,explicit,tag:1"`
	Usec           int               `asn1:"optional,explicit,tag:2"`
	SequenceNumber int64             `asn1:"optional,explicit,tag:3"`
	SAddress       types.HostAddress `asn1:"explicit,tag:4"`
	RAddress       types.HostAddress `asn1:"optional,explicit,tag:5"`
}

// Unmarshal bytes b into the KRBSafe struct.
func (s *KRBSafe) Unmarshal(b []byte) error {
	_, err := asn1.UnmarshalWithParams(b, s, fmt.Sprintf("application,explicit,tag:%v", asnAppTag.KRBSafe))
	if err != nil {
		return processUnmarshalReplyError(b, err)
	}
	expectedMsgType := msgtype.KRB_SAFE
	if s.MsgType != expectedMsgType {
		return krberror.NewErrorf(krberror.KRBMsgError, "message ID does not indicate a KRB_SAFE. Expected: %v; Actual: %v", expectedMsgType, s.MsgType)
	}
	return nil
}
