// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package operation

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// Copy from bb.yandex-team.ru/cloud/cloud-go/pkg/protoutil/any.go
func UnmarshalAny(msg *any.Any) (proto.Message, error) {
	if msg == nil {
		return nil, nil
	}
	box := &ptypes.DynamicAny{}
	err := ptypes.UnmarshalAny(msg, box)
	if err != nil {
		return nil, err
	}
	return box.Message, nil
}
