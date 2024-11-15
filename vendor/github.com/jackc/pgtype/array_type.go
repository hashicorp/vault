package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/jackc/pgio"
)

// ArrayType represents an array type. While it implements Value, this is only in service of its type conversion duties
// when registered as a data type in a ConnType. It should not be used directly as a Value. ArrayType is a convenience
// type for types that do not have a concrete array type.
type ArrayType struct {
	elements   []ValueTranscoder
	dimensions []ArrayDimension

	typeName   string
	newElement func() ValueTranscoder

	elementOID uint32
	status     Status
}

func NewArrayType(typeName string, elementOID uint32, newElement func() ValueTranscoder) *ArrayType {
	return &ArrayType{typeName: typeName, elementOID: elementOID, newElement: newElement}
}

func (at *ArrayType) NewTypeValue() Value {
	return &ArrayType{
		elements:   at.elements,
		dimensions: at.dimensions,
		status:     at.status,

		typeName:   at.typeName,
		elementOID: at.elementOID,
		newElement: at.newElement,
	}
}

func (at *ArrayType) TypeName() string {
	return at.typeName
}

func (dst *ArrayType) setNil() {
	dst.elements = nil
	dst.dimensions = nil
	dst.status = Null
}

func (dst *ArrayType) Set(src interface{}) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		dst.setNil()
		return nil
	}

	sliceVal := reflect.ValueOf(src)
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("cannot set non-slice")
	}

	if sliceVal.IsNil() {
		dst.setNil()
		return nil
	}

	dst.elements = make([]ValueTranscoder, sliceVal.Len())
	for i := range dst.elements {
		v := dst.newElement()
		err := v.Set(sliceVal.Index(i).Interface())
		if err != nil {
			return err
		}

		dst.elements[i] = v
	}
	dst.dimensions = []ArrayDimension{{Length: int32(len(dst.elements)), LowerBound: 1}}
	dst.status = Present

	return nil
}

func (dst ArrayType) Get() interface{} {
	switch dst.status {
	case Present:
		elementValues := make([]interface{}, len(dst.elements))
		for i := range dst.elements {
			elementValues[i] = dst.elements[i].Get()
		}
		return elementValues
	case Null:
		return nil
	default:
		return dst.status
	}
}

func (src *ArrayType) AssignTo(dst interface{}) error {
	ptrSlice := reflect.ValueOf(dst)
	if ptrSlice.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot assign to non-pointer")
	}

	sliceVal := ptrSlice.Elem()
	sliceType := sliceVal.Type()

	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("cannot assign to pointer to non-slice")
	}

	switch src.status {
	case Present:
		slice := reflect.MakeSlice(sliceType, len(src.elements), len(src.elements))
		elemType := sliceType.Elem()

		for i := range src.elements {
			ptrElem := reflect.New(elemType)
			err := src.elements[i].AssignTo(ptrElem.Interface())
			if err != nil {
				return err
			}

			slice.Index(i).Set(ptrElem.Elem())
		}

		sliceVal.Set(slice)
		return nil
	case Null:
		sliceVal.Set(reflect.Zero(sliceType))
		return nil
	}

	return fmt.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *ArrayType) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		dst.setNil()
		return nil
	}

	uta, err := ParseUntypedTextArray(string(src))
	if err != nil {
		return err
	}

	var elements []ValueTranscoder

	if len(uta.Elements) > 0 {
		elements = make([]ValueTranscoder, len(uta.Elements))

		for i, s := range uta.Elements {
			elem := dst.newElement()
			var elemSrc []byte
			if s != "NULL" {
				elemSrc = []byte(s)
			}
			err = elem.DecodeText(ci, elemSrc)
			if err != nil {
				return err
			}

			elements[i] = elem
		}
	}

	dst.elements = elements
	dst.dimensions = uta.Dimensions
	dst.status = Present

	return nil
}

func (dst *ArrayType) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		dst.setNil()
		return nil
	}

	var arrayHeader ArrayHeader
	rp, err := arrayHeader.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	var elements []ValueTranscoder

	if len(arrayHeader.Dimensions) == 0 {
		dst.elements = elements
		dst.dimensions = arrayHeader.Dimensions
		dst.status = Present
		return nil
	}

	elementCount := arrayHeader.Dimensions[0].Length
	for _, d := range arrayHeader.Dimensions[1:] {
		elementCount *= d.Length
	}

	elements = make([]ValueTranscoder, elementCount)

	for i := range elements {
		elem := dst.newElement()
		elemLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4
		var elemSrc []byte
		if elemLen >= 0 {
			elemSrc = src[rp : rp+elemLen]
			rp += elemLen
		}
		err = elem.DecodeBinary(ci, elemSrc)
		if err != nil {
			return err
		}

		elements[i] = elem
	}

	dst.elements = elements
	dst.dimensions = arrayHeader.Dimensions
	dst.status = Present

	return nil
}

func (src ArrayType) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	if len(src.dimensions) == 0 {
		return append(buf, '{', '}'), nil
	}

	buf = EncodeTextArrayDimensions(buf, src.dimensions)

	// dimElemCounts is the multiples of elements that each array lies on. For
	// example, a single dimension array of length 4 would have a dimElemCounts of
	// [4]. A multi-dimensional array of lengths [3,5,2] would have a
	// dimElemCounts of [30,10,2]. This is used to simplify when to render a '{'
	// or '}'.
	dimElemCounts := make([]int, len(src.dimensions))
	dimElemCounts[len(src.dimensions)-1] = int(src.dimensions[len(src.dimensions)-1].Length)
	for i := len(src.dimensions) - 2; i > -1; i-- {
		dimElemCounts[i] = int(src.dimensions[i].Length) * dimElemCounts[i+1]
	}

	inElemBuf := make([]byte, 0, 32)
	for i, elem := range src.elements {
		if i > 0 {
			buf = append(buf, ',')
		}

		for _, dec := range dimElemCounts {
			if i%dec == 0 {
				buf = append(buf, '{')
			}
		}

		elemBuf, err := elem.EncodeText(ci, inElemBuf)
		if err != nil {
			return nil, err
		}
		if elemBuf == nil {
			buf = append(buf, `NULL`...)
		} else {
			buf = append(buf, QuoteArrayElementIfNeeded(string(elemBuf))...)
		}

		for _, dec := range dimElemCounts {
			if (i+1)%dec == 0 {
				buf = append(buf, '}')
			}
		}
	}

	return buf, nil
}

func (src ArrayType) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	arrayHeader := ArrayHeader{
		Dimensions: src.dimensions,
		ElementOID: int32(src.elementOID),
	}

	for i := range src.elements {
		if src.elements[i].Get() == nil {
			arrayHeader.ContainsNull = true
			break
		}
	}

	buf = arrayHeader.EncodeBinary(ci, buf)

	for i := range src.elements {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		elemBuf, err := src.elements[i].EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if elemBuf != nil {
			buf = elemBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *ArrayType) Scan(src interface{}) error {
	if src == nil {
		return dst.DecodeText(nil, nil)
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src ArrayType) Value() (driver.Value, error) {
	buf, err := src.EncodeText(nil, nil)
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}

	return string(buf), nil
}
