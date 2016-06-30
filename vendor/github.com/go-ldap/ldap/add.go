//
// https://tools.ietf.org/html/rfc4511
//
// AddRequest ::= [APPLICATION 8] SEQUENCE {
//      entry           LDAPDN,
//      attributes      AttributeList }
//
// AttributeList ::= SEQUENCE OF attribute Attribute

package ldap

import (
	"errors"
	"log"

	"gopkg.in/asn1-ber.v1"
)

type Attribute struct {
	Type string
	Vals []string
}

func (a *Attribute) encode() *ber.Packet {
	seq := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Attribute")
	seq.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, a.Type, "Type"))
	set := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "AttributeValue")
	for _, value := range a.Vals {
		set.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, value, "Vals"))
	}
	seq.AppendChild(set)
	return seq
}

type AddRequest struct {
	DN         string
	Attributes []Attribute
}

func (a AddRequest) encode() *ber.Packet {
	request := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ApplicationAddRequest, nil, "Add Request")
	request.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, a.DN, "DN"))
	attributes := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Attributes")
	for _, attribute := range a.Attributes {
		attributes.AppendChild(attribute.encode())
	}
	request.AppendChild(attributes)
	return request
}

func (a *AddRequest) Attribute(attrType string, attrVals []string) {
	a.Attributes = append(a.Attributes, Attribute{Type: attrType, Vals: attrVals})
}

func NewAddRequest(dn string) *AddRequest {
	return &AddRequest{
		DN: dn,
	}

}

func (l *Conn) Add(addRequest *AddRequest) error {
	packet := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Request")
	packet.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, l.nextMessageID(), "MessageID"))
	packet.AppendChild(addRequest.encode())

	l.Debug.PrintPacket(packet)

	msgCtx, err := l.sendMessage(packet)
	if err != nil {
		return err
	}
	defer l.finishMessage(msgCtx)

	l.Debug.Printf("%d: waiting for response", msgCtx.id)
	packetResponse, ok := <-msgCtx.responses
	if !ok {
		return NewError(ErrorNetwork, errors.New("ldap: response channel closed"))
	}
	packet, err = packetResponse.ReadPacket()
	l.Debug.Printf("%d: got response %p", msgCtx.id, packet)
	if err != nil {
		return err
	}

	if l.Debug {
		if err := addLDAPDescriptions(packet); err != nil {
			return err
		}
		ber.PrintPacket(packet)
	}

	if packet.Children[1].Tag == ApplicationAddResponse {
		resultCode, resultDescription := getLDAPResultCode(packet)
		if resultCode != 0 {
			return NewError(resultCode, errors.New(resultDescription))
		}
	} else {
		log.Printf("Unexpected Response: %d", packet.Children[1].Tag)
	}

	l.Debug.Printf("%d: returning", msgCtx.id)
	return nil
}
