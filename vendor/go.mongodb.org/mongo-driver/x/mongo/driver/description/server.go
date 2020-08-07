// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/tag"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/address"
)

// UnsetRTT is the unset value for a round trip time.
const UnsetRTT = -1 * time.Millisecond

// SelectedServer represents a selected server that is a member of a topology.
type SelectedServer struct {
	Server
	Kind TopologyKind
}

// Server represents a description of a server. This is created from an isMaster
// command.
type Server struct {
	Addr address.Address

	AverageRTT            time.Duration
	AverageRTTSet         bool
	Compression           []string // compression methods returned by server
	CanonicalAddr         address.Address
	ElectionID            primitive.ObjectID
	HeartbeatInterval     time.Duration
	LastError             error
	LastUpdateTime        time.Time
	LastWriteTime         time.Time
	MaxBatchCount         uint32
	MaxDocumentSize       uint32
	MaxMessageSize        uint32
	Members               []address.Address
	ReadOnly              bool
	SessionTimeoutMinutes uint32
	SetName               string
	SetVersion            uint32
	Tags                  tag.Set
	Kind                  ServerKind
	WireVersion           *VersionRange

	SaslSupportedMechs []string // user-specific from server handshake
}

// NewServer creates a new server description from the given parameters.
func NewServer(addr address.Address, response bsoncore.Document) Server {
	desc := Server{Addr: addr, CanonicalAddr: addr, LastUpdateTime: time.Now().UTC()}
	elements, err := response.Elements()
	if err != nil {
		desc.LastError = err
		return desc
	}
	var ok bool
	var isReplicaSet, isMaster, hidden, secondary, arbiterOnly bool
	var msg string
	var version VersionRange
	var hosts, passives, arbiters []string
	for _, element := range elements {
		switch element.Key() {
		case "arbiters":
			var err error
			arbiters, err = decodeStringSlice(element, "arbiters")
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "arbiterOnly":
			arbiterOnly, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'arbiterOnly' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "compression":
			var err error
			desc.Compression, err = decodeStringSlice(element, "compression")
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "electionId":
			desc.ElectionID, ok = element.Value().ObjectIDOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'electionId' to be a objectID but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "hidden":
			hidden, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'hidden' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "hosts":
			var err error
			hosts, err = decodeStringSlice(element, "hosts")
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "ismaster":
			isMaster, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'isMaster' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "isreplicaset":
			isReplicaSet, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'isreplicaset' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "lastWrite":
			lastWrite, ok := element.Value().DocumentOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'lastWrite' to be a document but it's a BSON %s", element.Value().Type)
				return desc
			}
			dateTime, err := lastWrite.LookupErr("lastWriteDate")
			if err == nil {
				dt, ok := dateTime.DateTimeOK()
				if !ok {
					desc.LastError = fmt.Errorf("expected 'lastWriteDate' to be a datetime but it's a BSON %s", dateTime.Type)
					return desc
				}
				desc.LastWriteTime = time.Unix(dt/1000, dt%1000*1000000).UTC()
			}
		case "logicalSessionTimeoutMinutes":
			i64, ok := element.Value().AsInt64OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'logicalSessionTimeoutMinutes' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.SessionTimeoutMinutes = uint32(i64)
		case "maxBsonObjectSize":
			i64, ok := element.Value().AsInt64OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'maxBsonObjectSize' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.MaxDocumentSize = uint32(i64)
		case "maxMessageSizeBytes":
			i64, ok := element.Value().AsInt64OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'maxMessageSizeBytes' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.MaxMessageSize = uint32(i64)
		case "maxWriteBatchSize":
			i64, ok := element.Value().AsInt64OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'maxWriteBatchSize' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.MaxBatchCount = uint32(i64)
		case "me":
			me, ok := element.Value().StringValueOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'me' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.CanonicalAddr = address.Address(me).Canonicalize()
		case "maxWireVersion":
			version.Max, ok = element.Value().AsInt32OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'maxWireVersion' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "minWireVersion":
			version.Min, ok = element.Value().AsInt32OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'minWireVersion' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "msg":
			msg, ok = element.Value().StringValueOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'msg' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "ok":
			okay, ok := element.Value().AsInt32OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'ok' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
			if okay != 1 {
				desc.LastError = errors.New("not ok")
				return desc
			}
		case "passives":
			var err error
			passives, err = decodeStringSlice(element, "passives")
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "readOnly":
			desc.ReadOnly, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'readOnly' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "saslSupportedMechs":
			var err error
			desc.SaslSupportedMechs, err = decodeStringSlice(element, "saslSupportedMechs")
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "secondary":
			secondary, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'secondary' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "setName":
			desc.SetName, ok = element.Value().StringValueOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'setName' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "setVersion":
			i64, ok := element.Value().AsInt64OK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'setVersion' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.SetVersion = uint32(i64)
		case "tags":
			m, err := decodeStringMap(element, "tags")
			if err != nil {
				desc.LastError = err
				return desc
			}
			desc.Tags = tag.NewTagSetFromMap(m)
		}
	}

	for _, host := range hosts {
		desc.Members = append(desc.Members, address.Address(host).Canonicalize())
	}

	for _, passive := range passives {
		desc.Members = append(desc.Members, address.Address(passive).Canonicalize())
	}

	for _, arbiter := range arbiters {
		desc.Members = append(desc.Members, address.Address(arbiter).Canonicalize())
	}

	desc.Kind = Standalone

	if isReplicaSet {
		desc.Kind = RSGhost
	} else if desc.SetName != "" {
		if isMaster {
			desc.Kind = RSPrimary
		} else if hidden {
			desc.Kind = RSMember
		} else if secondary {
			desc.Kind = RSSecondary
		} else if arbiterOnly {
			desc.Kind = RSArbiter
		} else {
			desc.Kind = RSMember
		}
	} else if msg == "isdbgrid" {
		desc.Kind = Mongos
	}

	desc.WireVersion = &version

	return desc
}

// SetAverageRTT sets the average round trip time for this server description.
func (s Server) SetAverageRTT(rtt time.Duration) Server {
	s.AverageRTT = rtt
	if rtt == UnsetRTT {
		s.AverageRTTSet = false
	} else {
		s.AverageRTTSet = true
	}

	return s
}

// DataBearing returns true if the server is a data bearing server.
func (s Server) DataBearing() bool {
	return s.Kind == RSPrimary ||
		s.Kind == RSSecondary ||
		s.Kind == Mongos ||
		s.Kind == Standalone
}

// SelectServer selects this server if it is in the list of given candidates.
func (s Server) SelectServer(_ Topology, candidates []Server) ([]Server, error) {
	for _, candidate := range candidates {
		if candidate.Addr == s.Addr {
			return []Server{candidate}, nil
		}
	}
	return nil, nil
}

func decodeStringSlice(element bsoncore.Element, name string) ([]string, error) {
	arr, ok := element.Value().ArrayOK()
	if !ok {
		return nil, fmt.Errorf("expected '%s' to be an array but it's a BSON %s", name, element.Value().Type)
	}
	vals, err := arr.Values()
	if err != nil {
		return nil, err
	}
	var strs []string
	for _, val := range vals {
		str, ok := val.StringValueOK()
		if !ok {
			return nil, fmt.Errorf("expected '%s' to be an array of strings, but found a BSON %s", name, val.Type)
		}
		strs = append(strs, str)
	}
	return strs, nil
}

func decodeStringMap(element bsoncore.Element, name string) (map[string]string, error) {
	doc, ok := element.Value().DocumentOK()
	if !ok {
		return nil, fmt.Errorf("expected '%s' to be a document but it's a BSON %s", name, element.Value().Type)
	}
	elements, err := doc.Elements()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, element := range elements {
		key := element.Key()
		value, ok := element.Value().StringValueOK()
		if !ok {
			return nil, fmt.Errorf("expected '%s' to be a document of strings, but found a BSON %s", name, element.Value().Type)
		}
		m[key] = value
	}
	return m, nil
}
