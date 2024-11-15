/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"encoding/base64"
	"math/big"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	proto3 "google.golang.org/protobuf/types/known/structpb"
)

// Helpers to generate protobuf values and Cloud Spanner types.

func stringProto(s string) *proto3.Value {
	return &proto3.Value{Kind: stringKind(s)}
}

func stringKind(s string) *proto3.Value_StringValue {
	return &proto3.Value_StringValue{StringValue: s}
}

func stringType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_STRING}
}

func boolProto(b bool) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_BoolValue{BoolValue: b}}
}

func boolType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_BOOL}
}

func intProto(n int64) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: strconv.FormatInt(n, 10)}}
}

func intType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_INT64}
}

func float32Proto(n float32) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_NumberValue{NumberValue: float64(n)}}
}

func float32Type() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_FLOAT32}
}

func floatProto(n float64) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_NumberValue{NumberValue: n}}
}

func floatType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_FLOAT64}
}

func numericProto(n *big.Rat) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: NumericString(n)}}
}

func numericType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_NUMERIC}
}

func pgNumericType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_NUMERIC, TypeAnnotation: sppb.TypeAnnotationCode_PG_NUMERIC}
}

func pgOidType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_INT64, TypeAnnotation: sppb.TypeAnnotationCode_PG_OID}
}

func jsonType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_JSON}
}

func pgJsonbType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_JSON, TypeAnnotation: sppb.TypeAnnotationCode_PG_JSONB}
}

func bytesProto(b []byte) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: base64.StdEncoding.EncodeToString(b)}}
}

func bytesType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_BYTES}
}

func timeProto(t time.Time) *proto3.Value {
	return stringProto(t.UTC().Format(time.RFC3339Nano))
}

func timeType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_TIMESTAMP}
}

func dateProto(d civil.Date) *proto3.Value {
	return stringProto(d.String())
}

func dateType() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_DATE}
}

func listProto(p ...*proto3.Value) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_ListValue{ListValue: &proto3.ListValue{Values: p}}}
}

func listValueProto(p ...*proto3.Value) *proto3.ListValue {
	return &proto3.ListValue{Values: p}
}

func listType(t *sppb.Type) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_ARRAY, ArrayElementType: t}
}

func mkField(n string, t *sppb.Type) *sppb.StructType_Field {
	return &sppb.StructType_Field{Name: n, Type: t}
}

func structType(fields ...*sppb.StructType_Field) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_STRUCT, StructType: &sppb.StructType{Fields: fields}}
}

func nullProto() *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_NullValue{NullValue: proto3.NullValue_NULL_VALUE}}
}

func protoMessageType(fqn string) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_PROTO, ProtoTypeFqn: fqn}
}

func protoEnumType(fqn string) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_ENUM, ProtoTypeFqn: fqn}
}

func protoMessageProto(m proto.Message) *proto3.Value {
	var b, _ = proto.Marshal(m)
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: base64.StdEncoding.EncodeToString(b)}}
}

func protoEnumProto(e protoreflect.Enum) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: strconv.FormatInt(int64(e.Number()), 10)}}
}
