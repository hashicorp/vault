package types

// Reference: https://www.ietf.org/rfc/rfc4120.txt
// Section: 5.2.8

import (
	"github.com/jcmturner/gofork/encoding/asn1"
)

/*
KerberosFlags

For several message types, a specific constrained bit string type,
KerberosFlags, is used.

KerberosFlags   ::= BIT STRING (SIZE (32..MAX))
-- minimum number of bits shall be sent,
-- but no fewer than 32

Compatibility note: The following paragraphs describe a change from
the RFC 1510 description of bit strings that would result in
incompatility in the case of an implementation that strictly
conformed to ASN.1 DER and RFC 1510.

ASN.1 bit strings have multiple uses.  The simplest use of a bit
string is to contain a vector of bits, with no particular meaning
attached to individual bits.  This vector of bits is not necessarily
a multiple of eight bits long.  The use in Kerberos of a bit string
as a compact boolean vector wherein each element has a distinct
meaning poses some problems.  The natural notation for a compact
boolean vector is the ASN.1 "NamedBit" notation, and the DER require
that encodings of a bit string using "NamedBit" notation exclude any
trailing zero bits.  This truncation is easy to neglect, especially
given C language implementations that naturally choose to store
boolean vectors as 32-bit integers.

For example, if the notation for KDCOptions were to include the
"NamedBit" notation, as in RFC 1510, and a KDCOptions value to be
encoded had only the "forwardable" (bit number one) bit set, the DER
encoding MUST include only two bits: the first reserved bit
("reserved", bit number zero, value zero) and the one-valued bit (bit
number one) for "forwardable".

Most existing implementations of Kerberos unconditionally send 32
bits on the wire when encoding bit strings used as boolean vectors.
This behavior violates the ASN.1 syntax used for flag values in RFC
1510, but it occurs on such a widely installed base that the protocol
description is being modified to accommodate it.

Consequently, this document removes the "NamedBit" notations for
individual bits, relegating them to comments.  The size constraint on
the KerberosFlags type requires that at least 32 bits be encoded at
all times, though a lenient implementation MAY choose to accept fewer
than 32 bits and to treat the missing bits as set to zero.

Currently, no uses of KerberosFlags specify more than 32 bits' worth
of flags, although future revisions of this document may do so.  When
more than 32 bits are to be transmitted in a KerberosFlags value,
future revisions to this document will likely specify that the
smallest number of bits needed to encode the highest-numbered one-
valued bit should be sent.  This is somewhat similar to the DER
encoding of a bit string that is declared with the "NamedBit"
notation.
*/

// NewKrbFlags returns an ASN1 BitString struct of the right size for KrbFlags.
func NewKrbFlags() asn1.BitString {
	f := asn1.BitString{}
	f.Bytes = make([]byte, 4)
	f.BitLength = len(f.Bytes) * 8
	return f
}

// SetFlags sets the flags of an ASN1 BitString.
func SetFlags(f *asn1.BitString, j []int) {
	for _, i := range j {
		SetFlag(f, i)
	}
}

// SetFlag sets a flag in an ASN1 BitString.
func SetFlag(f *asn1.BitString, i int) {
	for l := len(f.Bytes); l < 4; l++ {
		(*f).Bytes = append((*f).Bytes, byte(0))
		(*f).BitLength = len((*f).Bytes) * 8
	}
	//Which byte?
	b := i / 8
	//Which bit in byte
	p := uint(7 - (i - 8*b))
	(*f).Bytes[b] = (*f).Bytes[b] | (1 << p)
}

// UnsetFlags unsets flags in an ASN1 BitString.
func UnsetFlags(f *asn1.BitString, j []int) {
	for _, i := range j {
		UnsetFlag(f, i)
	}
}

// UnsetFlag unsets a flag in an ASN1 BitString.
func UnsetFlag(f *asn1.BitString, i int) {
	for l := len(f.Bytes); l < 4; l++ {
		(*f).Bytes = append((*f).Bytes, byte(0))
		(*f).BitLength = len((*f).Bytes) * 8
	}
	//Which byte?
	b := i / 8
	//Which bit in byte
	p := uint(7 - (i - 8*b))
	(*f).Bytes[b] = (*f).Bytes[b] &^ (1 << p)
}

// IsFlagSet tests if a flag is set in the ASN1 BitString.
func IsFlagSet(f *asn1.BitString, i int) bool {
	//Which byte?
	b := i / 8
	//Which bit in byte
	p := uint(7 - (i - 8*b))
	if (*f).Bytes[b]&(1<<p) != 0 {
		return true
	}
	return false
}
