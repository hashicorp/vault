package operation

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/tag"
	"go.mongodb.org/mongo-driver/version"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// IsMaster is used to run the isMaster handshake operation.
type IsMaster struct {
	appname            string
	compressors        []string
	saslSupportedMechs string
	d                  driver.Deployment
	clock              *session.ClusterClock

	res bsoncore.Document
}

var _ driver.Handshaker = (*IsMaster)(nil)

// NewIsMaster constructs an IsMaster.
func NewIsMaster() *IsMaster { return &IsMaster{} }

// AppName sets the application name in the client metadata sent in this operation.
func (im *IsMaster) AppName(appname string) *IsMaster {
	im.appname = appname
	return im
}

// ClusterClock sets the cluster clock for this operation.
func (im *IsMaster) ClusterClock(clock *session.ClusterClock) *IsMaster {
	if im == nil {
		im = new(IsMaster)
	}

	im.clock = clock
	return im
}

// Compressors sets the compressors that can be used.
func (im *IsMaster) Compressors(compressors []string) *IsMaster {
	im.compressors = compressors
	return im
}

// SASLSupportedMechs retrieves the supported SASL mechanism for the given user when this operation
// is run.
func (im *IsMaster) SASLSupportedMechs(username string) *IsMaster {
	im.saslSupportedMechs = username
	return im
}

// Deployment sets the Deployment for this operation.
func (im *IsMaster) Deployment(d driver.Deployment) *IsMaster {
	im.d = d
	return im
}

// Result returns the result of executing this operation.
func (im *IsMaster) Result(addr address.Address) description.Server {
	desc := description.Server{Addr: addr, CanonicalAddr: addr, LastUpdateTime: time.Now().UTC()}
	elements, err := im.res.Elements()
	if err != nil {
		desc.LastError = err
		return desc
	}
	var ok bool
	var isReplicaSet, isMaster, hidden, secondary, arbiterOnly bool
	var msg string
	var version description.VersionRange
	var hosts, passives, arbiters []string
	for _, element := range elements {
		switch element.Key() {
		case "arbiters":
			var err error
			arbiters, err = im.decodeStringSlice(element, "arbiters")
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
			desc.Compression, err = im.decodeStringSlice(element, "compression")
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
			hosts, err = im.decodeStringSlice(element, "hosts")
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
			passives, err = im.decodeStringSlice(element, "passives")
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
			desc.SaslSupportedMechs, err = im.decodeStringSlice(element, "saslSupportedMechs")
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
			m, err := im.decodeStringMap(element, "tags")
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

	desc.Kind = description.Standalone

	if isReplicaSet {
		desc.Kind = description.RSGhost
	} else if desc.SetName != "" {
		if isMaster {
			desc.Kind = description.RSPrimary
		} else if hidden {
			desc.Kind = description.RSMember
		} else if secondary {
			desc.Kind = description.RSSecondary
		} else if arbiterOnly {
			desc.Kind = description.RSArbiter
		} else {
			desc.Kind = description.RSMember
		}
	} else if msg == "isdbgrid" {
		desc.Kind = description.Mongos
	}

	desc.WireVersion = &version

	return desc
}

func (im *IsMaster) decodeStringSlice(element bsoncore.Element, name string) ([]string, error) {
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

func (im *IsMaster) decodeStringMap(element bsoncore.Element, name string) (map[string]string, error) {
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

// handshakeCommand appends all necessary command fields as well as client metadata, SASL supported mechs, and compression.
func (im *IsMaster) handshakeCommand(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst, err := im.command(dst, desc)
	if err != nil {
		return dst, err
	}

	if im.saslSupportedMechs != "" {
		dst = bsoncore.AppendStringElement(dst, "saslSupportedMechs", im.saslSupportedMechs)
	}
	var idx int32
	idx, dst = bsoncore.AppendArrayElementStart(dst, "compression")
	for i, compressor := range im.compressors {
		dst = bsoncore.AppendStringElement(dst, strconv.Itoa(i), compressor)
	}
	dst, _ = bsoncore.AppendArrayEnd(dst, idx)

	// append client metadata
	idx, dst = bsoncore.AppendDocumentElementStart(dst, "client")

	didx, dst := bsoncore.AppendDocumentElementStart(dst, "driver")
	dst = bsoncore.AppendStringElement(dst, "name", "mongo-go-driver")
	dst = bsoncore.AppendStringElement(dst, "version", version.Driver)
	dst, _ = bsoncore.AppendDocumentEnd(dst, didx)

	didx, dst = bsoncore.AppendDocumentElementStart(dst, "os")
	dst = bsoncore.AppendStringElement(dst, "type", runtime.GOOS)
	dst = bsoncore.AppendStringElement(dst, "architecture", runtime.GOARCH)
	dst, _ = bsoncore.AppendDocumentEnd(dst, didx)

	dst = bsoncore.AppendStringElement(dst, "platform", runtime.Version())
	if im.appname != "" {
		didx, dst = bsoncore.AppendDocumentElementStart(dst, "application")
		dst = bsoncore.AppendStringElement(dst, "name", im.appname)
		dst, _ = bsoncore.AppendDocumentEnd(dst, didx)
	}
	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)

	return dst, nil
}

// command appends all necessary command fields.
func (im *IsMaster) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	return bsoncore.AppendInt32Element(dst, "isMaster", 1), nil
}

// Execute runs this operation.
func (im *IsMaster) Execute(ctx context.Context) error {
	if im.d == nil {
		return errors.New("an IsMaster must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		Clock:      im.clock,
		CommandFn:  im.command,
		Database:   "admin",
		Deployment: im.d,
		ProcessResponseFn: func(response bsoncore.Document, _ driver.Server, _ description.Server) error {
			im.res = response
			return nil
		},
	}.Execute(ctx, nil)
}

// GetDescription retrieves the server description for the given connection. This function implements the Handshaker
// interface.
func (im *IsMaster) GetDescription(ctx context.Context, _ address.Address, c driver.Connection) (description.Server, error) {
	err := driver.Operation{
		Clock:      im.clock,
		CommandFn:  im.handshakeCommand,
		Deployment: driver.SingleConnectionDeployment{c},
		Database:   "admin",
		ProcessResponseFn: func(response bsoncore.Document, _ driver.Server, _ description.Server) error {
			im.res = response
			return nil
		},
	}.Execute(ctx, nil)
	if err != nil {
		return description.Server{}, err
	}
	return im.Result(c.Address()), nil
}

// FinishHandshake implements the Handshaker interface. This is a no-op function because a non-authenticated connection
// does not do anything besides the initial isMaster for a handshake.
func (im *IsMaster) FinishHandshake(context.Context, driver.Connection) error {
	return nil
}
