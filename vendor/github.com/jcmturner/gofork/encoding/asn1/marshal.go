// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asn1

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"time"
	"unicode/utf8"
)

// A forkableWriter is an in-memory buffer that can be
// 'forked' to create new forkableWriters that bracket the
// original. After
//    pre, post := w.fork()
// the overall sequence of bytes represented is logically w+pre+post.
type forkableWriter struct {
	*bytes.Buffer
	pre, post *forkableWriter
}

func newForkableWriter() *forkableWriter {
	return &forkableWriter{new(bytes.Buffer), nil, nil}
}

func (f *forkableWriter) fork() (pre, post *forkableWriter) {
	if f.pre != nil || f.post != nil {
		panic("have already forked")
	}
	f.pre = newForkableWriter()
	f.post = newForkableWriter()
	return f.pre, f.post
}

func (f *forkableWriter) Len() (l int) {
	l += f.Buffer.Len()
	if f.pre != nil {
		l += f.pre.Len()
	}
	if f.post != nil {
		l += f.post.Len()
	}
	return
}

func (f *forkableWriter) writeTo(out io.Writer) (n int, err error) {
	n, err = out.Write(f.Bytes())
	if err != nil {
		return
	}

	var nn int

	if f.pre != nil {
		nn, err = f.pre.writeTo(out)
		n += nn
		if err != nil {
			return
		}
	}

	if f.post != nil {
		nn, err = f.post.writeTo(out)
		n += nn
	}
	return
}

func marshalBase128Int(out *forkableWriter, n int64) (err error) {
	if n == 0 {
		err = out.WriteByte(0)
		return
	}

	l := 0
	for i := n; i > 0; i >>= 7 {
		l++
	}

	for i := l - 1; i >= 0; i-- {
		o := byte(n >> uint(i*7))
		o &= 0x7f
		if i != 0 {
			o |= 0x80
		}
		err = out.WriteByte(o)
		if err != nil {
			return
		}
	}

	return nil
}

func marshalInt64(out *forkableWriter, i int64) (err error) {
	n := int64Length(i)

	for ; n > 0; n-- {
		err = out.WriteByte(byte(i >> uint((n-1)*8)))
		if err != nil {
			return
		}
	}

	return nil
}

func int64Length(i int64) (numBytes int) {
	numBytes = 1

	for i > 127 {
		numBytes++
		i >>= 8
	}

	for i < -128 {
		numBytes++
		i >>= 8
	}

	return
}

func marshalBigInt(out *forkableWriter, n *big.Int) (err error) {
	if n.Sign() < 0 {
		// A negative number has to be converted to two's-complement
		// form. So we'll subtract 1 and invert. If the
		// most-significant-bit isn't set then we'll need to pad the
		// beginning with 0xff in order to keep the number negative.
		nMinus1 := new(big.Int).Neg(n)
		nMinus1.Sub(nMinus1, bigOne)
		bytes := nMinus1.Bytes()
		for i := range bytes {
			bytes[i] ^= 0xff
		}
		if len(bytes) == 0 || bytes[0]&0x80 == 0 {
			err = out.WriteByte(0xff)
			if err != nil {
				return
			}
		}
		_, err = out.Write(bytes)
	} else if n.Sign() == 0 {
		// Zero is written as a single 0 zero rather than no bytes.
		err = out.WriteByte(0x00)
	} else {
		bytes := n.Bytes()
		if len(bytes) > 0 && bytes[0]&0x80 != 0 {
			// We'll have to pad this with 0x00 in order to stop it
			// looking like a negative number.
			err = out.WriteByte(0)
			if err != nil {
				return
			}
		}
		_, err = out.Write(bytes)
	}
	return
}

func marshalLength(out *forkableWriter, i int) (err error) {
	n := lengthLength(i)

	for ; n > 0; n-- {
		err = out.WriteByte(byte(i >> uint((n-1)*8)))
		if err != nil {
			return
		}
	}

	return nil
}

func lengthLength(i int) (numBytes int) {
	numBytes = 1
	for i > 255 {
		numBytes++
		i >>= 8
	}
	return
}

func marshalTagAndLength(out *forkableWriter, t tagAndLength) (err error) {
	b := uint8(t.class) << 6
	if t.isCompound {
		b |= 0x20
	}
	if t.tag >= 31 {
		b |= 0x1f
		err = out.WriteByte(b)
		if err != nil {
			return
		}
		err = marshalBase128Int(out, int64(t.tag))
		if err != nil {
			return
		}
	} else {
		b |= uint8(t.tag)
		err = out.WriteByte(b)
		if err != nil {
			return
		}
	}

	if t.length >= 128 {
		l := lengthLength(t.length)
		err = out.WriteByte(0x80 | byte(l))
		if err != nil {
			return
		}
		err = marshalLength(out, t.length)
		if err != nil {
			return
		}
	} else {
		err = out.WriteByte(byte(t.length))
		if err != nil {
			return
		}
	}

	return nil
}

func marshalBitString(out *forkableWriter, b BitString) (err error) {
	paddingBits := byte((8 - b.BitLength%8) % 8)
	err = out.WriteByte(paddingBits)
	if err != nil {
		return
	}
	_, err = out.Write(b.Bytes)
	return
}

func marshalObjectIdentifier(out *forkableWriter, oid []int) (err error) {
	if len(oid) < 2 || oid[0] > 2 || (oid[0] < 2 && oid[1] >= 40) {
		return StructuralError{"invalid object identifier"}
	}

	err = marshalBase128Int(out, int64(oid[0]*40+oid[1]))
	if err != nil {
		return
	}
	for i := 2; i < len(oid); i++ {
		err = marshalBase128Int(out, int64(oid[i]))
		if err != nil {
			return
		}
	}

	return
}

func marshalPrintableString(out *forkableWriter, s string) (err error) {
	b := []byte(s)
	for _, c := range b {
		if !isPrintable(c) {
			return StructuralError{"PrintableString contains invalid character"}
		}
	}

	_, err = out.Write(b)
	return
}

func marshalIA5String(out *forkableWriter, s string) (err error) {
	b := []byte(s)
	for _, c := range b {
		if c > 127 {
			return StructuralError{"IA5String contains invalid character"}
		}
	}

	_, err = out.Write(b)
	return
}

func marshalUTF8String(out *forkableWriter, s string) (err error) {
	_, err = out.Write([]byte(s))
	return
}

func marshalTwoDigits(out *forkableWriter, v int) (err error) {
	err = out.WriteByte(byte('0' + (v/10)%10))
	if err != nil {
		return
	}
	return out.WriteByte(byte('0' + v%10))
}

func marshalFourDigits(out *forkableWriter, v int) (err error) {
	var bytes [4]byte
	for i := range bytes {
		bytes[3-i] = '0' + byte(v%10)
		v /= 10
	}
	_, err = out.Write(bytes[:])
	return
}

func outsideUTCRange(t time.Time) bool {
	year := t.Year()
	return year < 1950 || year >= 2050
}

func marshalUTCTime(out *forkableWriter, t time.Time) (err error) {
	year := t.Year()

	switch {
	case 1950 <= year && year < 2000:
		err = marshalTwoDigits(out, year-1900)
	case 2000 <= year && year < 2050:
		err = marshalTwoDigits(out, year-2000)
	default:
		return StructuralError{"cannot represent time as UTCTime"}
	}
	if err != nil {
		return
	}

	return marshalTimeCommon(out, t)
}

func marshalGeneralizedTime(out *forkableWriter, t time.Time) (err error) {
	year := t.Year()
	if year < 0 || year > 9999 {
		return StructuralError{"cannot represent time as GeneralizedTime"}
	}
	if err = marshalFourDigits(out, year); err != nil {
		return
	}

	return marshalTimeCommon(out, t)
}

func marshalTimeCommon(out *forkableWriter, t time.Time) (err error) {
	_, month, day := t.Date()

	err = marshalTwoDigits(out, int(month))
	if err != nil {
		return
	}

	err = marshalTwoDigits(out, day)
	if err != nil {
		return
	}

	hour, min, sec := t.Clock()

	err = marshalTwoDigits(out, hour)
	if err != nil {
		return
	}

	err = marshalTwoDigits(out, min)
	if err != nil {
		return
	}

	err = marshalTwoDigits(out, sec)
	if err != nil {
		return
	}

	_, offset := t.Zone()

	switch {
	case offset/60 == 0:
		err = out.WriteByte('Z')
		return
	case offset > 0:
		err = out.WriteByte('+')
	case offset < 0:
		err = out.WriteByte('-')
	}

	if err != nil {
		return
	}

	offsetMinutes := offset / 60
	if offsetMinutes < 0 {
		offsetMinutes = -offsetMinutes
	}

	err = marshalTwoDigits(out, offsetMinutes/60)
	if err != nil {
		return
	}

	err = marshalTwoDigits(out, offsetMinutes%60)
	return
}

func stripTagAndLength(in []byte) []byte {
	_, offset, err := parseTagAndLength(in, 0)
	if err != nil {
		return in
	}
	return in[offset:]
}

func marshalBody(out *forkableWriter, value reflect.Value, params fieldParameters) (err error) {
	switch value.Type() {
	case flagType:
		return nil
	case timeType:
		t := value.Interface().(time.Time)
		if params.timeType == TagGeneralizedTime || outsideUTCRange(t) {
			return marshalGeneralizedTime(out, t)
		} else {
			return marshalUTCTime(out, t)
		}
	case bitStringType:
		return marshalBitString(out, value.Interface().(BitString))
	case objectIdentifierType:
		return marshalObjectIdentifier(out, value.Interface().(ObjectIdentifier))
	case bigIntType:
		return marshalBigInt(out, value.Interface().(*big.Int))
	}

	switch v := value; v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return out.WriteByte(255)
		} else {
			return out.WriteByte(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return marshalInt64(out, v.Int())
	case reflect.Struct:
		t := v.Type()

		startingField := 0

		// If the first element of the structure is a non-empty
		// RawContents, then we don't bother serializing the rest.
		if t.NumField() > 0 && t.Field(0).Type == rawContentsType {
			s := v.Field(0)
			if s.Len() > 0 {
				bytes := make([]byte, s.Len())
				for i := 0; i < s.Len(); i++ {
					bytes[i] = uint8(s.Index(i).Uint())
				}
				/* The RawContents will contain the tag and
				 * length fields but we'll also be writing
				 * those ourselves, so we strip them out of
				 * bytes */
				_, err = out.Write(stripTagAndLength(bytes))
				return
			} else {
				startingField = 1
			}
		}

		for i := startingField; i < t.NumField(); i++ {
			var pre *forkableWriter
			pre, out = out.fork()
			err = marshalField(pre, v.Field(i), parseFieldParameters(t.Field(i).Tag.Get("asn1")))
			if err != nil {
				return
			}
		}
		return
	case reflect.Slice:
		sliceType := v.Type()
		if sliceType.Elem().Kind() == reflect.Uint8 {
			bytes := make([]byte, v.Len())
			for i := 0; i < v.Len(); i++ {
				bytes[i] = uint8(v.Index(i).Uint())
			}
			_, err = out.Write(bytes)
			return
		}

		// jtasn1 Pass on the tags to the members but need to unset explicit switch and implicit value
		//var fp fieldParameters
		params.explicit = false
		params.tag = nil
		for i := 0; i < v.Len(); i++ {
			var pre *forkableWriter
			pre, out = out.fork()
			err = marshalField(pre, v.Index(i), params)
			if err != nil {
				return
			}
		}
		return
	case reflect.String:
		switch params.stringType {
		case TagIA5String:
			return marshalIA5String(out, v.String())
		case TagPrintableString:
			return marshalPrintableString(out, v.String())
		default:
			return marshalUTF8String(out, v.String())
		}
	}

	return StructuralError{"unknown Go type"}
}

func marshalField(out *forkableWriter, v reflect.Value, params fieldParameters) (err error) {
	if !v.IsValid() {
		return fmt.Errorf("asn1: cannot marshal nil value")
	}
	// If the field is an interface{} then recurse into it.
	if v.Kind() == reflect.Interface && v.Type().NumMethod() == 0 {
		return marshalField(out, v.Elem(), params)
	}

	if v.Kind() == reflect.Slice && v.Len() == 0 && params.omitEmpty {
		return
	}

	if params.optional && params.defaultValue != nil && canHaveDefaultValue(v.Kind()) {
		defaultValue := reflect.New(v.Type()).Elem()
		defaultValue.SetInt(*params.defaultValue)

		if reflect.DeepEqual(v.Interface(), defaultValue.Interface()) {
			return
		}
	}

	// If no default value is given then the zero value for the type is
	// assumed to be the default value. This isn't obviously the correct
	// behaviour, but it's what Go has traditionally done.
	if params.optional && params.defaultValue == nil {
		if reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
			return
		}
	}

	if v.Type() == rawValueType {
		rv := v.Interface().(RawValue)
		if len(rv.FullBytes) != 0 {
			_, err = out.Write(rv.FullBytes)
		} else {
			err = marshalTagAndLength(out, tagAndLength{rv.Class, rv.Tag, len(rv.Bytes), rv.IsCompound})
			if err != nil {
				return
			}
			_, err = out.Write(rv.Bytes)
		}
		return
	}

	tag, isCompound, ok := getUniversalType(v.Type())
	if !ok {
		err = StructuralError{fmt.Sprintf("unknown Go type: %v", v.Type())}
		return
	}
	class := ClassUniversal

	if params.timeType != 0 && tag != TagUTCTime {
		return StructuralError{"explicit time type given to non-time member"}
	}

	// jtasn1 updated to allow slices of strings
	if params.stringType != 0 && !(tag == TagPrintableString || (v.Kind() == reflect.Slice && tag == 16 && v.Type().Elem().Kind() == reflect.String)) {
		return StructuralError{"explicit string type given to non-string member"}
	}

	switch tag {
	case TagPrintableString:
		if params.stringType == 0 {
			// This is a string without an explicit string type. We'll use
			// a PrintableString if the character set in the string is
			// sufficiently limited, otherwise we'll use a UTF8String.
			for _, r := range v.String() {
				if r >= utf8.RuneSelf || !isPrintable(byte(r)) {
					if !utf8.ValidString(v.String()) {
						return errors.New("asn1: string not valid UTF-8")
					}
					tag = TagUTF8String
					break
				}
			}
		} else {
			tag = params.stringType
		}
	case TagUTCTime:
		if params.timeType == TagGeneralizedTime || outsideUTCRange(v.Interface().(time.Time)) {
			tag = TagGeneralizedTime
		}
	}

	if params.set {
		if tag != TagSequence {
			return StructuralError{"non sequence tagged as set"}
		}
		tag = TagSet
	}

	tags, body := out.fork()

	err = marshalBody(body, v, params)
	if err != nil {
		return
	}

	bodyLen := body.Len()

	var explicitTag *forkableWriter
	if params.explicit {
		explicitTag, tags = tags.fork()
	}

	if !params.explicit && params.tag != nil {
		// implicit tag.
		tag = *params.tag
		class = ClassContextSpecific
	}

	err = marshalTagAndLength(tags, tagAndLength{class, tag, bodyLen, isCompound})
	if err != nil {
		return
	}

	if params.explicit {
		err = marshalTagAndLength(explicitTag, tagAndLength{
			class:      ClassContextSpecific,
			tag:        *params.tag,
			length:     bodyLen + tags.Len(),
			isCompound: true,
		})
	}

	return err
}

// Marshal returns the ASN.1 encoding of val.
//
// In addition to the struct tags recognised by Unmarshal, the following can be
// used:
//
//	ia5:		causes strings to be marshaled as ASN.1, IA5 strings
//	omitempty:	causes empty slices to be skipped
//	printable:	causes strings to be marshaled as ASN.1, PrintableString strings.
//	utf8:		causes strings to be marshaled as ASN.1, UTF8 strings
func Marshal(val interface{}) ([]byte, error) {
	var out bytes.Buffer
	v := reflect.ValueOf(val)
	f := newForkableWriter()
	err := marshalField(f, v, fieldParameters{})
	if err != nil {
		return nil, err
	}
	_, err = f.writeTo(&out)
	return out.Bytes(), err
}
