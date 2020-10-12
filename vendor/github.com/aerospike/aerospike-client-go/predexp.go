// Copyright 2017-2019 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"fmt"
	"math"
	"strconv"
)

const (
	_AS_PREDEXP_UNKNOWN_BIN uint16 = math.MaxUint16

	_AS_PREDEXP_AND uint16 = 1
	_AS_PREDEXP_OR  uint16 = 2
	_AS_PREDEXP_NOT uint16 = 3

	_AS_PREDEXP_INTEGER_VALUE uint16 = 10
	_AS_PREDEXP_STRING_VALUE  uint16 = 11
	_AS_PREDEXP_GEOJSON_VALUE uint16 = 12

	_AS_PREDEXP_INTEGER_BIN uint16 = 100
	_AS_PREDEXP_STRING_BIN  uint16 = 101
	_AS_PREDEXP_GEOJSON_BIN uint16 = 102
	_AS_PREDEXP_LIST_BIN    uint16 = 103
	_AS_PREDEXP_MAP_BIN     uint16 = 104

	_AS_PREDEXP_INTEGER_VAR uint16 = 120
	_AS_PREDEXP_STRING_VAR  uint16 = 121
	_AS_PREDEXP_GEOJSON_VAR uint16 = 122

	_AS_PREDEXP_REC_DEVICE_SIZE   uint16 = 150
	_AS_PREDEXP_REC_LAST_UPDATE   uint16 = 151
	_AS_PREDEXP_REC_VOID_TIME     uint16 = 152
	_AS_PREDEXP_REC_DIGEST_MODULO uint16 = 153

	_AS_PREDEXP_INTEGER_EQUAL     uint16 = 200
	_AS_PREDEXP_INTEGER_UNEQUAL   uint16 = 201
	_AS_PREDEXP_INTEGER_GREATER   uint16 = 202
	_AS_PREDEXP_INTEGER_GREATEREQ uint16 = 203
	_AS_PREDEXP_INTEGER_LESS      uint16 = 204
	_AS_PREDEXP_INTEGER_LESSEQ    uint16 = 205

	_AS_PREDEXP_STRING_EQUAL   uint16 = 210
	_AS_PREDEXP_STRING_UNEQUAL uint16 = 211
	_AS_PREDEXP_STRING_REGEX   uint16 = 212

	_AS_PREDEXP_GEOJSON_WITHIN   uint16 = 220
	_AS_PREDEXP_GEOJSON_CONTAINS uint16 = 221

	_AS_PREDEXP_LIST_ITERATE_OR    uint16 = 250
	_AS_PREDEXP_MAPKEY_ITERATE_OR  uint16 = 251
	_AS_PREDEXP_MAPVAL_ITERATE_OR  uint16 = 252
	_AS_PREDEXP_LIST_ITERATE_AND   uint16 = 253
	_AS_PREDEXP_MAPKEY_ITERATE_AND uint16 = 254
	_AS_PREDEXP_MAPVAL_ITERATE_AND uint16 = 255
)

// ----------------

type PredExp interface {
	String() string
	marshaledSize() int
	marshal(*baseCommand) error
}

type predExpBase struct {
}

func (e *predExpBase) marshaledSize() int {
	return 2 + 4 // sizeof(TAG) + sizeof(LEN)
}

func (e *predExpBase) marshalTL(cmd *baseCommand, tag uint16, len uint32) {
	cmd.WriteUint16(tag)
	cmd.WriteUint32(len)
}

// ---------------- predExpAnd

type predExpAnd struct {
	predExpBase
	nexpr uint16 // number of child expressions
}

// String implements the Stringer interface
func (e *predExpAnd) String() string {
	return "AND"
}

// NewPredExpAnd creates an AND predicate. Argument describes the number of expressions.
func NewPredExpAnd(nexpr uint16) *predExpAnd {
	return &predExpAnd{nexpr: nexpr}
}

func (e *predExpAnd) marshaledSize() int {
	return e.predExpBase.marshaledSize() + 2
}

func (e *predExpAnd) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_AND, 2)
	cmd.WriteUint16(e.nexpr)
	return nil
}

// ---------------- predExpOr

type predExpOr struct {
	predExpBase
	nexpr uint16 // number of child expressions
}

// String implements the Stringer interface
func (e *predExpOr) String() string {
	return "OR"
}

// NewPredExpOr creates an OR predicate. Argument describes the number of expressions.
func NewPredExpOr(nexpr uint16) *predExpOr {
	return &predExpOr{nexpr: nexpr}
}

func (e *predExpOr) marshaledSize() int {
	return e.predExpBase.marshaledSize() + 2
}

func (e *predExpOr) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_OR, 2)
	cmd.WriteUint16(e.nexpr)
	return nil
}

// ---------------- predExpNot

type predExpNot struct {
	predExpBase
}

// String implements the Stringer interface
func (e *predExpNot) String() string {
	return "NOT"
}

// NewPredExpNot creates a NOT predicate
func NewPredExpNot() *predExpNot {
	return &predExpNot{}
}

func (e *predExpNot) marshaledSize() int {
	return e.predExpBase.marshaledSize()
}

func (e *predExpNot) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_NOT, 0)
	return nil
}

// ---------------- predExpIntegerValue

type predExpIntegerValue struct {
	predExpBase
	val int64
}

// String implements the Stringer interface
func (e *predExpIntegerValue) String() string {
	return strconv.FormatInt(e.val, 10)
}

// NewPredExpIntegerValue embeds an int64 value in a predicate expression.
func NewPredExpIntegerValue(val int64) *predExpIntegerValue {
	return &predExpIntegerValue{val: val}
}

func (e *predExpIntegerValue) marshaledSize() int {
	return e.predExpBase.marshaledSize() + 8
}

func (e *predExpIntegerValue) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_INTEGER_VALUE, 8)
	cmd.WriteInt64(e.val)
	return nil
}

// ---------------- predExpStringValue

type predExpStringValue struct {
	predExpBase
	val string
}

// String implements the Stringer interface
func (e *predExpStringValue) String() string {
	return "'" + e.val + "'"
}

// NewPredExpStringValue embeds a string value in a predicate expression.
func NewPredExpStringValue(val string) *predExpStringValue {
	return &predExpStringValue{val: val}
}

func (e *predExpStringValue) marshaledSize() int {
	return e.predExpBase.marshaledSize() + len(e.val)
}

func (e *predExpStringValue) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_STRING_VALUE, uint32(len(e.val)))
	cmd.WriteString(e.val)
	return nil
}

// ---------------- predExpGeoJSONValue

type predExpGeoJSONValue struct {
	predExpBase
	val string
}

// String implements the Stringer interface
func (e *predExpGeoJSONValue) String() string {
	return e.val
}

// NewPredExpGeoJSONValue embeds a GeoJSON value in a predicate expression.
func NewPredExpGeoJSONValue(val string) *predExpGeoJSONValue {
	return &predExpGeoJSONValue{val: val}
}

func (e *predExpGeoJSONValue) marshaledSize() int {
	return e.predExpBase.marshaledSize() +
		1 + // flags
		2 + // ncells
		len(e.val) // strlen value
}

func (e *predExpGeoJSONValue) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_GEOJSON_VALUE, uint32(1+2+len(e.val)))
	cmd.WriteByte(uint8(0))
	cmd.WriteUint16(0)
	cmd.WriteString(e.val)
	return nil
}

// ---------------- predExp???Bin

type predExpBin struct {
	predExpBase
	name string
	tag  uint16 // not marshaled
}

// String implements the Stringer interface
func (e *predExpBin) String() string {
	// FIXME - This is not currently distinguished from a var.
	return e.name
}

// NewPredExpUnknownBin creates a Bin predicate expression which its type is not known.
func NewPredExpUnknownBin(name string) *predExpBin {
	return &predExpBin{name: name, tag: _AS_PREDEXP_UNKNOWN_BIN}
}

// NewPredExpUnknownBin creates a Bin predicate expression which its type is integer.
func NewPredExpIntegerBin(name string) *predExpBin {
	return &predExpBin{name: name, tag: _AS_PREDEXP_INTEGER_BIN}
}

// NewPredExpUnknownBin creates a Bin predicate expression which its type is String.
func NewPredExpStringBin(name string) *predExpBin {
	return &predExpBin{name: name, tag: _AS_PREDEXP_STRING_BIN}
}

// NewPredExpUnknownBin creates a Bin predicate expression which its type is GeoJSON.
func NewPredExpGeoJSONBin(name string) *predExpBin {
	return &predExpBin{name: name, tag: _AS_PREDEXP_GEOJSON_BIN}
}

// NewPredExpUnknownBin creates a Bin predicate expression which its type is List.
func NewPredExpListBin(name string) *predExpBin {
	return &predExpBin{name: name, tag: _AS_PREDEXP_LIST_BIN}
}

// NewPredExpUnknownBin creates a Bin predicate expression which its type is Map.
func NewPredExpMapBin(name string) *predExpBin {
	return &predExpBin{name: name, tag: _AS_PREDEXP_MAP_BIN}
}

func (e *predExpBin) marshaledSize() int {
	return e.predExpBase.marshaledSize() + len(e.name)
}

func (e *predExpBin) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, e.tag, uint32(len(e.name)))
	cmd.WriteString(e.name)
	return nil
}

// ---------------- predExp???Var

type predExpVar struct {
	predExpBase
	name string
	tag  uint16 // not marshaled
}

// String implements the Stringer interface
func (e *predExpVar) String() string {
	// FIXME - This is not currently distinguished from a bin.
	return e.name
}

// NewPredExpIntegerVar creates 64 bit integer variable used in list/map iterations.
func NewPredExpIntegerVar(name string) *predExpVar {
	return &predExpVar{name: name, tag: _AS_PREDEXP_INTEGER_VAR}
}

// NewPredExpStringVar creates string variable used in list/map iterations.
func NewPredExpStringVar(name string) *predExpVar {
	return &predExpVar{name: name, tag: _AS_PREDEXP_STRING_VAR}
}

// NewPredExpGeoJSONVar creates GeoJSON variable used in list/map iterations.
func NewPredExpGeoJSONVar(name string) *predExpVar {
	return &predExpVar{name: name, tag: _AS_PREDEXP_GEOJSON_VAR}
}

func (e *predExpVar) marshaledSize() int {
	return e.predExpBase.marshaledSize() + len(e.name)
}

func (e *predExpVar) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, e.tag, uint32(len(e.name)))
	cmd.WriteString(e.name)
	return nil
}

// ---------------- predExpMD (RecDeviceSize, RecLastUpdate, RecVoidTime)

type predExpMD struct {
	predExpBase
	tag uint16 // not marshaled
}

// String implements the Stringer interface
func (e *predExpMD) String() string {
	switch e.tag {
	case _AS_PREDEXP_REC_DEVICE_SIZE:
		return "rec.DeviceSize"
	case _AS_PREDEXP_REC_LAST_UPDATE:
		return "rec.LastUpdate"
	case _AS_PREDEXP_REC_VOID_TIME:
		return "rec.Expiration"
	case _AS_PREDEXP_REC_DIGEST_MODULO:
		return "rec.DigestModulo"
	default:
		panic("Invalid Metadata tag.")
	}
}

func (e *predExpMD) marshaledSize() int {
	return e.predExpBase.marshaledSize()
}

func (e *predExpMD) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, e.tag, 0)
	return nil
}

// NewPredExpRecDeviceSize creates record size on disk predicate
func NewPredExpRecDeviceSize() *predExpMD {
	return &predExpMD{tag: _AS_PREDEXP_REC_DEVICE_SIZE}
}

// NewPredExpRecLastUpdate creates record last update predicate
func NewPredExpRecLastUpdate() *predExpMD {
	return &predExpMD{tag: _AS_PREDEXP_REC_LAST_UPDATE}
}

// NewPredExpRecVoidTime creates record expiration time predicate expressed in nanoseconds since 1970-01-01 epoch as 64 bit integer.
func NewPredExpRecVoidTime() *predExpMD {
	return &predExpMD{tag: _AS_PREDEXP_REC_VOID_TIME}
}

// ---------------- predExpMDDigestModulo

type predExpMDDigestModulo struct {
	predExpBase
	mod int32
}

// String implements the Stringer interface
func (e *predExpMDDigestModulo) String() string {
	return "rec.DigestModulo"
}

func (e *predExpMDDigestModulo) marshaledSize() int {
	return e.predExpBase.marshaledSize() + 4
}

func (e *predExpMDDigestModulo) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_REC_DIGEST_MODULO, 4)
	cmd.WriteInt32(e.mod)
	return nil
}

// NewPredExpRecDigestModulo creates a digest modulo record metadata value predicate expression.
// The digest modulo expression assumes the value of 4 bytes of the
// record's key digest modulo as its argument.
// This predicate is available in Aerospike server versions 3.12.1+
//
// For example, the following sequence of predicate expressions
// selects records that have digest(key) % 3 == 1):
// stmt.SetPredExp(
// 		NewPredExpRecDigestModulo(3),
// 		NewPredExpIntegerValue(1),
// 		NewPredExpIntegerEqual(),
// )
func NewPredExpRecDigestModulo(mod int32) *predExpMDDigestModulo {
	return &predExpMDDigestModulo{mod: mod}
}

// ---------------- predExpCompare

type predExpCompare struct {
	predExpBase
	tag uint16 // not marshaled
}

// String implements the Stringer interface
func (e *predExpCompare) String() string {
	switch e.tag {
	case _AS_PREDEXP_INTEGER_EQUAL, _AS_PREDEXP_STRING_EQUAL:
		return "="
	case _AS_PREDEXP_INTEGER_UNEQUAL, _AS_PREDEXP_STRING_UNEQUAL:
		return "!="
	case _AS_PREDEXP_INTEGER_GREATER:
		return ">"
	case _AS_PREDEXP_INTEGER_GREATEREQ:
		return ">="
	case _AS_PREDEXP_INTEGER_LESS:
		return "<"
	case _AS_PREDEXP_INTEGER_LESSEQ:
		return "<="
	case _AS_PREDEXP_STRING_REGEX:
		return "~="
	case _AS_PREDEXP_GEOJSON_CONTAINS:
		return "CONTAINS"
	case _AS_PREDEXP_GEOJSON_WITHIN:
		return "WITHIN"
	default:
		panic(fmt.Sprintf("unexpected predicate tag: %d", e.tag))
	}
}

func (e *predExpCompare) marshaledSize() int {
	return e.predExpBase.marshaledSize()
}

func (e *predExpCompare) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, e.tag, 0)
	return nil
}

// NewPredExpIntegerEqual creates Equal predicate for integer values
func NewPredExpIntegerEqual() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_INTEGER_EQUAL}
}

// NewPredExpIntegerUnequal creates NotEqual predicate for integer values
func NewPredExpIntegerUnequal() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_INTEGER_UNEQUAL}
}

// NewPredExpIntegerGreater creates Greater Than predicate for integer values
func NewPredExpIntegerGreater() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_INTEGER_GREATER}
}

// NewPredExpIntegerGreaterEq creates Greater Than Or Equal predicate for integer values
func NewPredExpIntegerGreaterEq() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_INTEGER_GREATEREQ}
}

// NewPredExpIntegerLess creates Less Than predicate for integer values
func NewPredExpIntegerLess() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_INTEGER_LESS}
}

// NewPredExpIntegerLessEq creates Less Than Or Equal predicate for integer values
func NewPredExpIntegerLessEq() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_INTEGER_LESSEQ}
}

// NewPredExpStringEqual creates Equal predicate for string values
func NewPredExpStringEqual() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_STRING_EQUAL}
}

// NewPredExpStringUnequal creates Not Equal predicate for string values
func NewPredExpStringUnequal() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_STRING_UNEQUAL}
}

// NewPredExpGeoJSONWithin creates Within Region predicate for GeoJSON values
func NewPredExpGeoJSONWithin() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_GEOJSON_WITHIN}
}

// NewPredExpGeoJSONContains creates Region Contains predicate for GeoJSON values
func NewPredExpGeoJSONContains() *predExpCompare {
	return &predExpCompare{tag: _AS_PREDEXP_GEOJSON_CONTAINS}
}

// ---------------- predExpStringRegex

type predExpStringRegex struct {
	predExpBase
	cflags uint32 // cflags
}

// String implements the Stringer interface
func (e *predExpStringRegex) String() string {
	return "regex:"
}

// NewPredExpStringRegex creates a Regex predicate
func NewPredExpStringRegex(cflags uint32) *predExpStringRegex {
	return &predExpStringRegex{cflags: cflags}
}

func (e *predExpStringRegex) marshaledSize() int {
	return e.predExpBase.marshaledSize() + 4
}

func (e *predExpStringRegex) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, _AS_PREDEXP_STRING_REGEX, 4)
	cmd.WriteUint32(e.cflags)
	return nil
}

// ---------------- predExp???Iterate???

type predExpIter struct {
	predExpBase
	name string
	tag  uint16 // not marshaled
}

// String implements the Stringer interface
func (e *predExpIter) String() string {
	switch e.tag {
	case _AS_PREDEXP_LIST_ITERATE_OR:
		return "list_iterate_or using \"" + e.name + "\":"
	case _AS_PREDEXP_MAPKEY_ITERATE_OR:
		return "mapkey_iterate_or using \"" + e.name + "\":"
	case _AS_PREDEXP_MAPVAL_ITERATE_OR:
		return "mapval_iterate_or using \"" + e.name + "\":"
	case _AS_PREDEXP_LIST_ITERATE_AND:
		return "list_iterate_and using \"" + e.name + "\":"
	case _AS_PREDEXP_MAPKEY_ITERATE_AND:
		return "mapkey_iterate_and using \"" + e.name + "\":"
	case _AS_PREDEXP_MAPVAL_ITERATE_AND:
		return "mapval_iterate_and using \"" + e.name + "\":"
	default:
		panic("Invalid Metadata tag.")
	}
}

// NewPredExpListIterateOr creates an Or iterator predicate for list items
func NewPredExpListIterateOr(name string) *predExpIter {
	return &predExpIter{name: name, tag: _AS_PREDEXP_LIST_ITERATE_OR}
}

// NewPredExpMapKeyIterateOr creates an Or iterator predicate on map keys
func NewPredExpMapKeyIterateOr(name string) *predExpIter {
	return &predExpIter{name: name, tag: _AS_PREDEXP_MAPKEY_ITERATE_OR}
}

// NewPredExpMapValIterateOr creates an Or iterator predicate on map values
func NewPredExpMapValIterateOr(name string) *predExpIter {
	return &predExpIter{name: name, tag: _AS_PREDEXP_MAPVAL_ITERATE_OR}
}

// NewPredExpListIterateAnd creates an And iterator predicate for list items
func NewPredExpListIterateAnd(name string) *predExpIter {
	return &predExpIter{name: name, tag: _AS_PREDEXP_LIST_ITERATE_AND}
}

// NewPredExpMapKeyIterateAnd creates an And iterator predicate on map keys
func NewPredExpMapKeyIterateAnd(name string) *predExpIter {
	return &predExpIter{name: name, tag: _AS_PREDEXP_MAPKEY_ITERATE_AND}
}

// NewPredExpMapKeyIterateAnd creates an And iterator predicate on map values
func NewPredExpMapValIterateAnd(name string) *predExpIter {
	return &predExpIter{name: name, tag: _AS_PREDEXP_MAPVAL_ITERATE_AND}
}

func (e *predExpIter) marshaledSize() int {
	return e.predExpBase.marshaledSize() + len(e.name)
}

func (e *predExpIter) marshal(cmd *baseCommand) error {
	e.marshalTL(cmd, e.tag, uint32(len(e.name)))
	cmd.WriteString(e.name)
	return nil
}
