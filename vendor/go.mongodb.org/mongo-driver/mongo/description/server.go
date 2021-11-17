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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/tag"
)

// SelectedServer augments the Server type by also including the TopologyKind of the topology that includes the server.
// This type should be used to track the state of a server that was selected to perform an operation.
type SelectedServer struct {
	Server
	Kind TopologyKind
}

// Server contains information about a node in a cluster. This is created from isMaster command responses. If the value
// of the Kind field is LoadBalancer, only the Addr and Kind fields will be set. All other fields will be set to the
// zero value of the field's type.
type Server struct {
	Addr address.Address

	Arbiters              []string
	AverageRTT            time.Duration
	AverageRTTSet         bool
	Compression           []string // compression methods returned by server
	CanonicalAddr         address.Address
	ElectionID            primitive.ObjectID
	HeartbeatInterval     time.Duration
	Hosts                 []string
	LastError             error
	LastUpdateTime        time.Time
	LastWriteTime         time.Time
	MaxBatchCount         uint32
	MaxDocumentSize       uint32
	MaxMessageSize        uint32
	Members               []address.Address
	Passives              []string
	Primary               address.Address
	ReadOnly              bool
	ServiceID             *primitive.ObjectID // Only set for servers that are deployed behind a load balancer.
	SessionTimeoutMinutes uint32
	SetName               string
	SetVersion            uint32
	Tags                  tag.Set
	TopologyVersion       *TopologyVersion
	Kind                  ServerKind
	WireVersion           *VersionRange
}

// NewServer creates a new server description from the given isMaster command response.
func NewServer(addr address.Address, response bson.Raw) Server {
	desc := Server{Addr: addr, CanonicalAddr: addr, LastUpdateTime: time.Now().UTC()}
	elements, err := response.Elements()
	if err != nil {
		desc.LastError = err
		return desc
	}
	var ok bool
	var isReplicaSet, isWritablePrimary, hidden, secondary, arbiterOnly bool
	var msg string
	var version VersionRange
	for _, element := range elements {
		switch element.Key() {
		case "arbiters":
			var err error
			desc.Arbiters, err = internal.StringSliceFromRawElement(element)
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
			desc.Compression, err = internal.StringSliceFromRawElement(element)
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
			desc.Hosts, err = internal.StringSliceFromRawElement(element)
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "isWritablePrimary":
			isWritablePrimary, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'isWritablePrimary' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "ismaster":
			isWritablePrimary, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'ismaster' to be a boolean but it's a BSON %s", element.Value().Type)
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
			desc.Passives, err = internal.StringSliceFromRawElement(element)
			if err != nil {
				desc.LastError = err
				return desc
			}
		case "primary":
			primary, ok := element.Value().StringValueOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'primary' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			}
			desc.Primary = address.Address(primary)
		case "readOnly":
			desc.ReadOnly, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'readOnly' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "secondary":
			secondary, ok = element.Value().BooleanOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'secondary' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			}
		case "serviceId":
			oid, ok := element.Value().ObjectIDOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'serviceId' to be an ObjectId but it's a BSON %s", element.Value().Type)
			}
			desc.ServiceID = &oid
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
		case "topologyVersion":
			doc, ok := element.Value().DocumentOK()
			if !ok {
				desc.LastError = fmt.Errorf("expected 'topologyVersion' to be a document but it's a BSON %s", element.Value().Type)
				return desc
			}

			desc.TopologyVersion, err = NewTopologyVersion(doc)
			if err != nil {
				desc.LastError = err
				return desc
			}

			if internal.SetMockServiceID {
				desc.ServiceID = &desc.TopologyVersion.ProcessID
			}
		}
	}

	for _, host := range desc.Hosts {
		desc.Members = append(desc.Members, address.Address(host).Canonicalize())
	}

	for _, passive := range desc.Passives {
		desc.Members = append(desc.Members, address.Address(passive).Canonicalize())
	}

	for _, arbiter := range desc.Arbiters {
		desc.Members = append(desc.Members, address.Address(arbiter).Canonicalize())
	}

	desc.Kind = Standalone

	if isReplicaSet {
		desc.Kind = RSGhost
	} else if desc.SetName != "" {
		if isWritablePrimary {
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

// NewDefaultServer creates a new unknown server description with the given address.
func NewDefaultServer(addr address.Address) Server {
	return NewServerFromError(addr, nil, nil)
}

// NewServerFromError creates a new unknown server description with the given parameters.
func NewServerFromError(addr address.Address, err error, tv *TopologyVersion) Server {
	return Server{
		Addr:            addr,
		LastError:       err,
		Kind:            Unknown,
		TopologyVersion: tv,
	}
}

// SetAverageRTT sets the average round trip time for this server description.
func (s Server) SetAverageRTT(rtt time.Duration) Server {
	s.AverageRTT = rtt
	s.AverageRTTSet = true
	return s
}

// DataBearing returns true if the server is a data bearing server.
func (s Server) DataBearing() bool {
	return s.Kind == RSPrimary ||
		s.Kind == RSSecondary ||
		s.Kind == Mongos ||
		s.Kind == Standalone
}

// LoadBalanced returns true if the server is a load balancer or is behind a load balancer.
func (s Server) LoadBalanced() bool {
	return s.Kind == LoadBalancer || s.ServiceID != nil
}

// String implements the Stringer interface
func (s Server) String() string {
	str := fmt.Sprintf("Addr: %s, Type: %s",
		s.Addr, s.Kind)
	if len(s.Tags) != 0 {
		str += fmt.Sprintf(", Tag sets: %s", s.Tags)
	}

	if s.AverageRTTSet {
		str += fmt.Sprintf(", Average RTT: %d", s.AverageRTT)
	}

	if s.LastError != nil {
		str += fmt.Sprintf(", Last error: %s", s.LastError)
	}
	return str
}

func decodeStringMap(element bson.RawElement, name string) (map[string]string, error) {
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

// Equal compares two server descriptions and returns true if they are equal
func (s Server) Equal(other Server) bool {
	if s.CanonicalAddr.String() != other.CanonicalAddr.String() {
		return false
	}

	if !sliceStringEqual(s.Arbiters, other.Arbiters) {
		return false
	}

	if !sliceStringEqual(s.Hosts, other.Hosts) {
		return false
	}

	if !sliceStringEqual(s.Passives, other.Passives) {
		return false
	}

	if s.Primary != other.Primary {
		return false
	}

	if s.SetName != other.SetName {
		return false
	}

	if s.Kind != other.Kind {
		return false
	}

	if s.LastError != nil || other.LastError != nil {
		if s.LastError == nil || other.LastError == nil {
			return false
		}
		if s.LastError.Error() != other.LastError.Error() {
			return false
		}
	}

	if !s.WireVersion.Equals(other.WireVersion) {
		return false
	}

	if len(s.Tags) != len(other.Tags) || !s.Tags.ContainsAll(other.Tags) {
		return false
	}

	if s.SetVersion != other.SetVersion {
		return false
	}

	if s.ElectionID != other.ElectionID {
		return false
	}

	if s.SessionTimeoutMinutes != other.SessionTimeoutMinutes {
		return false
	}

	// If TopologyVersion is nil for both servers, CompareToIncoming will return -1 because it assumes that the
	// incoming response is newer. We want the descriptions to be considered equal in this case, though, so an
	// explicit check is required.
	if s.TopologyVersion == nil && other.TopologyVersion == nil {
		return true
	}
	return s.TopologyVersion.CompareToIncoming(other.TopologyVersion) == 0
}

func sliceStringEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
