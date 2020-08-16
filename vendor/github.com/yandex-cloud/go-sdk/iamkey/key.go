// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package iamkey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	yaml2json "github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	yaml "gopkg.in/yaml.v2"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"github.com/yandex-cloud/go-sdk/pkg/sdkerrors"
)

var (
	_ json.Marshaler   = &Key{}
	_ json.Unmarshaler = &Key{}
	_ yaml.Marshaler   = &Key{}
	_ yaml.Unmarshaler = &Key{}
)

// New creates new Key from IAM Key Service Create response.
func New(created *iam.CreateKeyResponse) *Key {
	if created == nil {
		panic("nil key")
	}
	public := created.GetKey()
	key := &Key{
		Id:           public.GetId(),
		Subject:      nil,
		CreatedAt:    public.GetCreatedAt(),
		Description:  public.GetDescription(),
		KeyAlgorithm: public.GetKeyAlgorithm(),
		PublicKey:    public.GetPublicKey(),
		PrivateKey:   created.GetPrivateKey(),
	}
	switch subj := public.GetSubject().(type) {
	case *iam.Key_ServiceAccountId:
		key.Subject = &Key_ServiceAccountId{
			ServiceAccountId: subj.ServiceAccountId,
		}
	case *iam.Key_UserAccountId:
		key.Subject = &Key_UserAccountId{
			UserAccountId: subj.UserAccountId,
		}
	case nil:
		// Do nothing.
	default:
		panic(fmt.Sprintf("unexpected key subject: %#v", subj))
	}
	return key
}

// UnmarshalJSON unmarshals IAM Key JSON data.
// Both snake_case (gRPC API) and camelCase (REST API) fields are accepted.
func (m *Key) UnmarshalJSON(data []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(data), m)
}

func (m *Key) MarshalJSON() ([]byte, error) {
	marshaller := &jsonpb.Marshaler{OrigName: true}
	buf := &bytes.Buffer{}
	err := marshaller.Marshal(buf, m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalYAML unmarshals IAM Key YAML data.
// Both snake_case (gRPC API) and camelCase (REST API) fields are accepted.
func (m *Key) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var obj yaml.MapSlice
	err := unmarshal(&obj)
	if err != nil {
		return err
	}
	yamlData, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	jsonData, err := yaml2json.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}
	return m.UnmarshalJSON(jsonData)
}

func (m *Key) MarshalYAML() (interface{}, error) {
	jsonData, err := m.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var obj yaml.MapSlice
	err = yaml.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// ReadFromJSONFile reads IAM Key from JSON bytes.
func ReadFromJSONBytes(keyBytes []byte) (*Key, error) {
	key := &Key{}
	err := json.Unmarshal(keyBytes, key)
	if err != nil {
		return nil, sdkerrors.WithMessage(err, "key unmarshal fail")
	}
	return key, nil
}

// ReadFromJSONFile reads IAM Key from JSON file.
func ReadFromJSONFile(path string) (*Key, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, sdkerrors.WithMessagef(err, "key file '%s' read fail", path)
	}
	return ReadFromJSONBytes(data)
}

// WriteToJSONFile writes key to file in JSON format.
// File permissions will be 0600, because private key part is sensitive data.
func WriteToJSONFile(path string, key *Key) error {
	data, err := json.MarshalIndent(key, "", "   ")
	if err != nil {
		return sdkerrors.WithMessage(err, "key marshal fail")
	}
	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		return sdkerrors.WithMessagef(err, "file '%s' write fail", path)
	}
	return nil
}
