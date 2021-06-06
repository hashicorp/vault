package operation

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"

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
	speculativeAuth    bsoncore.Document
	topologyVersion    *description.TopologyVersion
	maxAwaitTimeMS     *int64

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

// SpeculativeAuthenticate sets the document to be used for speculative authentication.
func (im *IsMaster) SpeculativeAuthenticate(doc bsoncore.Document) *IsMaster {
	im.speculativeAuth = doc
	return im
}

// TopologyVersion sets the TopologyVersion to be used for heartbeats.
func (im *IsMaster) TopologyVersion(tv *description.TopologyVersion) *IsMaster {
	im.topologyVersion = tv
	return im
}

// MaxAwaitTimeMS sets the maximum time for the sever to wait for topology changes during a heartbeat.
func (im *IsMaster) MaxAwaitTimeMS(awaitTime int64) *IsMaster {
	im.maxAwaitTimeMS = &awaitTime
	return im
}

// Result returns the result of executing this operation.
func (im *IsMaster) Result(addr address.Address) description.Server {
	return description.NewServer(addr, im.res)
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
	if im.speculativeAuth != nil {
		dst = bsoncore.AppendDocumentElement(dst, "speculativeAuthenticate", im.speculativeAuth)
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
	dst = bsoncore.AppendInt32Element(dst, "isMaster", 1)

	if tv := im.topologyVersion; tv != nil {
		var tvIdx int32

		tvIdx, dst = bsoncore.AppendDocumentElementStart(dst, "topologyVersion")
		dst = bsoncore.AppendObjectIDElement(dst, "processId", tv.ProcessID)
		dst = bsoncore.AppendInt64Element(dst, "counter", tv.Counter)
		dst, _ = bsoncore.AppendDocumentEnd(dst, tvIdx)
	}
	if im.maxAwaitTimeMS != nil {
		dst = bsoncore.AppendInt64Element(dst, "maxAwaitTimeMS", *im.maxAwaitTimeMS)
	}

	return dst, nil
}

// Execute runs this operation.
func (im *IsMaster) Execute(ctx context.Context) error {
	if im.d == nil {
		return errors.New("an IsMaster must have a Deployment set before Execute can be called")
	}

	return im.createOperation().Execute(ctx, nil)
}

// StreamResponse gets the next streaming isMaster response from the server.
func (im *IsMaster) StreamResponse(ctx context.Context, conn driver.StreamerConnection) error {
	return im.createOperation().ExecuteExhaust(ctx, conn, nil)
}

func (im *IsMaster) createOperation() driver.Operation {
	return driver.Operation{
		Clock:      im.clock,
		CommandFn:  im.command,
		Database:   "admin",
		Deployment: im.d,
		ProcessResponseFn: func(response bsoncore.Document, _ driver.Server, _ description.Server, _ int) error {
			im.res = response
			return nil
		},
	}
}

// GetDescription retrieves the server description for the given connection. This function implements the Handshaker
// interface.
func (im *IsMaster) GetDescription(ctx context.Context, _ address.Address, c driver.Connection) (description.Server, error) {
	err := driver.Operation{
		Clock:      im.clock,
		CommandFn:  im.handshakeCommand,
		Deployment: driver.SingleConnectionDeployment{c},
		Database:   "admin",
		ProcessResponseFn: func(response bsoncore.Document, _ driver.Server, _ description.Server, _ int) error {
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
