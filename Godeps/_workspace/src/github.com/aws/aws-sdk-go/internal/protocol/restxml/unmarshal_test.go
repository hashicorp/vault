package restxml_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/internal/protocol/restxml"
	"github.com/aws/aws-sdk-go/internal/protocol/xml/xmlutil"
	"github.com/aws/aws-sdk-go/internal/signer/v4"
	"github.com/aws/aws-sdk-go/internal/util"
	"github.com/stretchr/testify/assert"
)

var _ bytes.Buffer // always import bytes
var _ http.Request
var _ json.Marshaler
var _ time.Time
var _ xmlutil.XMLNode
var _ xml.Attr
var _ = ioutil.Discard
var _ = util.Trim("")
var _ = url.Values{}
var _ = io.EOF

type OutputService1ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService1ProtocolTest client.
func NewOutputService1ProtocolTest(config *aws.Config) *OutputService1ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice1protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService1ProtocolTest{service}
}

// newRequest creates a new request for a OutputService1ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService1TestCaseOperation1 = "OperationName"

// OutputService1TestCaseOperation1Request generates a request for the OutputService1TestCaseOperation1 operation.
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1Request(input *OutputService1TestShapeOutputService1TestCaseOperation1Input) (req *aws.Request, output *OutputService1TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService1TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService1TestShapeOutputService1TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService1TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1(input *OutputService1TestShapeOutputService1TestCaseOperation1Input) (*OutputService1TestShapeOutputShape, error) {
	req, out := c.OutputService1TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

const opOutputService1TestCaseOperation2 = "OperationName"

// OutputService1TestCaseOperation2Request generates a request for the OutputService1TestCaseOperation2 operation.
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation2Request(input *OutputService1TestShapeOutputService1TestCaseOperation2Input) (req *aws.Request, output *OutputService1TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService1TestCaseOperation2,
	}

	if input == nil {
		input = &OutputService1TestShapeOutputService1TestCaseOperation2Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService1TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation2(input *OutputService1TestShapeOutputService1TestCaseOperation2Input) (*OutputService1TestShapeOutputShape, error) {
	req, out := c.OutputService1TestCaseOperation2Request(input)
	err := req.Send()
	return out, err
}

type OutputService1TestShapeOutputService1TestCaseOperation1Input struct {
	metadataOutputService1TestShapeOutputService1TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService1TestShapeOutputService1TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService1TestShapeOutputService1TestCaseOperation2Input struct {
	metadataOutputService1TestShapeOutputService1TestCaseOperation2Input `json:"-" xml:"-"`
}

type metadataOutputService1TestShapeOutputService1TestCaseOperation2Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService1TestShapeOutputShape struct {
	Char *string `type:"character"`

	Double *float64 `type:"double"`

	FalseBool *bool `type:"boolean"`

	Float *float64 `type:"float"`

	ImaHeader *string `location:"header" type:"string"`

	ImaHeaderLocation *string `location:"header" locationName:"X-Foo" type:"string"`

	Long *int64 `type:"long"`

	Num *int64 `locationName:"FooNum" type:"integer"`

	Str *string `type:"string"`

	Timestamp *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	TrueBool *bool `type:"boolean"`

	metadataOutputService1TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService1TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService2ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService2ProtocolTest client.
func NewOutputService2ProtocolTest(config *aws.Config) *OutputService2ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice2protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService2ProtocolTest{service}
}

// newRequest creates a new request for a OutputService2ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService2TestCaseOperation1 = "OperationName"

// OutputService2TestCaseOperation1Request generates a request for the OutputService2TestCaseOperation1 operation.
func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1Request(input *OutputService2TestShapeOutputService2TestCaseOperation1Input) (req *aws.Request, output *OutputService2TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService2TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService2TestShapeOutputService2TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService2TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1(input *OutputService2TestShapeOutputService2TestCaseOperation1Input) (*OutputService2TestShapeOutputShape, error) {
	req, out := c.OutputService2TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService2TestShapeOutputService2TestCaseOperation1Input struct {
	metadataOutputService2TestShapeOutputService2TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService2TestShapeOutputService2TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService2TestShapeOutputShape struct {
	Blob []byte `type:"blob"`

	metadataOutputService2TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService2TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService3ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService3ProtocolTest client.
func NewOutputService3ProtocolTest(config *aws.Config) *OutputService3ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice3protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService3ProtocolTest{service}
}

// newRequest creates a new request for a OutputService3ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService3TestCaseOperation1 = "OperationName"

// OutputService3TestCaseOperation1Request generates a request for the OutputService3TestCaseOperation1 operation.
func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1Request(input *OutputService3TestShapeOutputService3TestCaseOperation1Input) (req *aws.Request, output *OutputService3TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService3TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService3TestShapeOutputService3TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService3TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1(input *OutputService3TestShapeOutputService3TestCaseOperation1Input) (*OutputService3TestShapeOutputShape, error) {
	req, out := c.OutputService3TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService3TestShapeOutputService3TestCaseOperation1Input struct {
	metadataOutputService3TestShapeOutputService3TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService3TestShapeOutputService3TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService3TestShapeOutputShape struct {
	ListMember []*string `type:"list"`

	metadataOutputService3TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService3TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService4ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService4ProtocolTest client.
func NewOutputService4ProtocolTest(config *aws.Config) *OutputService4ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice4protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService4ProtocolTest{service}
}

// newRequest creates a new request for a OutputService4ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService4TestCaseOperation1 = "OperationName"

// OutputService4TestCaseOperation1Request generates a request for the OutputService4TestCaseOperation1 operation.
func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1Request(input *OutputService4TestShapeOutputService4TestCaseOperation1Input) (req *aws.Request, output *OutputService4TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService4TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService4TestShapeOutputService4TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService4TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1(input *OutputService4TestShapeOutputService4TestCaseOperation1Input) (*OutputService4TestShapeOutputShape, error) {
	req, out := c.OutputService4TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService4TestShapeOutputService4TestCaseOperation1Input struct {
	metadataOutputService4TestShapeOutputService4TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService4TestShapeOutputService4TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService4TestShapeOutputShape struct {
	ListMember []*string `locationNameList:"item" type:"list"`

	metadataOutputService4TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService4TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService5ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService5ProtocolTest client.
func NewOutputService5ProtocolTest(config *aws.Config) *OutputService5ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice5protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService5ProtocolTest{service}
}

// newRequest creates a new request for a OutputService5ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService5TestCaseOperation1 = "OperationName"

// OutputService5TestCaseOperation1Request generates a request for the OutputService5TestCaseOperation1 operation.
func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1Request(input *OutputService5TestShapeOutputService5TestCaseOperation1Input) (req *aws.Request, output *OutputService5TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService5TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService5TestShapeOutputService5TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService5TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1(input *OutputService5TestShapeOutputService5TestCaseOperation1Input) (*OutputService5TestShapeOutputShape, error) {
	req, out := c.OutputService5TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService5TestShapeOutputService5TestCaseOperation1Input struct {
	metadataOutputService5TestShapeOutputService5TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService5TestShapeOutputService5TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService5TestShapeOutputShape struct {
	ListMember []*string `type:"list" flattened:"true"`

	metadataOutputService5TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService5TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService6ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService6ProtocolTest client.
func NewOutputService6ProtocolTest(config *aws.Config) *OutputService6ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice6protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService6ProtocolTest{service}
}

// newRequest creates a new request for a OutputService6ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService6TestCaseOperation1 = "OperationName"

// OutputService6TestCaseOperation1Request generates a request for the OutputService6TestCaseOperation1 operation.
func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1Request(input *OutputService6TestShapeOutputService6TestCaseOperation1Input) (req *aws.Request, output *OutputService6TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService6TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService6TestShapeOutputService6TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService6TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1(input *OutputService6TestShapeOutputService6TestCaseOperation1Input) (*OutputService6TestShapeOutputShape, error) {
	req, out := c.OutputService6TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService6TestShapeOutputService6TestCaseOperation1Input struct {
	metadataOutputService6TestShapeOutputService6TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService6TestShapeOutputService6TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService6TestShapeOutputShape struct {
	Map map[string]*OutputService6TestShapeSingleStructure `type:"map"`

	metadataOutputService6TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService6TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService6TestShapeSingleStructure struct {
	Foo *string `locationName:"foo" type:"string"`

	metadataOutputService6TestShapeSingleStructure `json:"-" xml:"-"`
}

type metadataOutputService6TestShapeSingleStructure struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService7ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService7ProtocolTest client.
func NewOutputService7ProtocolTest(config *aws.Config) *OutputService7ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice7protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService7ProtocolTest{service}
}

// newRequest creates a new request for a OutputService7ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService7TestCaseOperation1 = "OperationName"

// OutputService7TestCaseOperation1Request generates a request for the OutputService7TestCaseOperation1 operation.
func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1Request(input *OutputService7TestShapeOutputService7TestCaseOperation1Input) (req *aws.Request, output *OutputService7TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService7TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService7TestShapeOutputService7TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService7TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1(input *OutputService7TestShapeOutputService7TestCaseOperation1Input) (*OutputService7TestShapeOutputShape, error) {
	req, out := c.OutputService7TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService7TestShapeOutputService7TestCaseOperation1Input struct {
	metadataOutputService7TestShapeOutputService7TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService7TestShapeOutputService7TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService7TestShapeOutputShape struct {
	Map map[string]*string `type:"map" flattened:"true"`

	metadataOutputService7TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService7TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService8ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService8ProtocolTest client.
func NewOutputService8ProtocolTest(config *aws.Config) *OutputService8ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice8protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService8ProtocolTest{service}
}

// newRequest creates a new request for a OutputService8ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService8TestCaseOperation1 = "OperationName"

// OutputService8TestCaseOperation1Request generates a request for the OutputService8TestCaseOperation1 operation.
func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1Request(input *OutputService8TestShapeOutputService8TestCaseOperation1Input) (req *aws.Request, output *OutputService8TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService8TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService8TestShapeOutputService8TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService8TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1(input *OutputService8TestShapeOutputService8TestCaseOperation1Input) (*OutputService8TestShapeOutputShape, error) {
	req, out := c.OutputService8TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService8TestShapeOutputService8TestCaseOperation1Input struct {
	metadataOutputService8TestShapeOutputService8TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService8TestShapeOutputService8TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService8TestShapeOutputShape struct {
	Map map[string]*string `locationNameKey:"foo" locationNameValue:"bar" type:"map"`

	metadataOutputService8TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService8TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService9ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService9ProtocolTest client.
func NewOutputService9ProtocolTest(config *aws.Config) *OutputService9ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice9protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService9ProtocolTest{service}
}

// newRequest creates a new request for a OutputService9ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService9TestCaseOperation1 = "OperationName"

// OutputService9TestCaseOperation1Request generates a request for the OutputService9TestCaseOperation1 operation.
func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1Request(input *OutputService9TestShapeOutputService9TestCaseOperation1Input) (req *aws.Request, output *OutputService9TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService9TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService9TestShapeOutputService9TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService9TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1(input *OutputService9TestShapeOutputService9TestCaseOperation1Input) (*OutputService9TestShapeOutputShape, error) {
	req, out := c.OutputService9TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService9TestShapeOutputService9TestCaseOperation1Input struct {
	metadataOutputService9TestShapeOutputService9TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService9TestShapeOutputService9TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService9TestShapeOutputShape struct {
	Data *OutputService9TestShapeSingleStructure `type:"structure"`

	Header *string `location:"header" locationName:"X-Foo" type:"string"`

	metadataOutputService9TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService9TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure" payload:"Data"`
}

type OutputService9TestShapeSingleStructure struct {
	Foo *string `type:"string"`

	metadataOutputService9TestShapeSingleStructure `json:"-" xml:"-"`
}

type metadataOutputService9TestShapeSingleStructure struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService10ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService10ProtocolTest client.
func NewOutputService10ProtocolTest(config *aws.Config) *OutputService10ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice10protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService10ProtocolTest{service}
}

// newRequest creates a new request for a OutputService10ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService10ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService10TestCaseOperation1 = "OperationName"

// OutputService10TestCaseOperation1Request generates a request for the OutputService10TestCaseOperation1 operation.
func (c *OutputService10ProtocolTest) OutputService10TestCaseOperation1Request(input *OutputService10TestShapeOutputService10TestCaseOperation1Input) (req *aws.Request, output *OutputService10TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService10TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService10TestShapeOutputService10TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService10TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService10ProtocolTest) OutputService10TestCaseOperation1(input *OutputService10TestShapeOutputService10TestCaseOperation1Input) (*OutputService10TestShapeOutputShape, error) {
	req, out := c.OutputService10TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService10TestShapeOutputService10TestCaseOperation1Input struct {
	metadataOutputService10TestShapeOutputService10TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService10TestShapeOutputService10TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService10TestShapeOutputShape struct {
	Stream []byte `type:"blob"`

	metadataOutputService10TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService10TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure" payload:"Stream"`
}

type OutputService11ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService11ProtocolTest client.
func NewOutputService11ProtocolTest(config *aws.Config) *OutputService11ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice11protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService11ProtocolTest{service}
}

// newRequest creates a new request for a OutputService11ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService11TestCaseOperation1 = "OperationName"

// OutputService11TestCaseOperation1Request generates a request for the OutputService11TestCaseOperation1 operation.
func (c *OutputService11ProtocolTest) OutputService11TestCaseOperation1Request(input *OutputService11TestShapeOutputService11TestCaseOperation1Input) (req *aws.Request, output *OutputService11TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService11TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService11TestShapeOutputService11TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService11TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService11ProtocolTest) OutputService11TestCaseOperation1(input *OutputService11TestShapeOutputService11TestCaseOperation1Input) (*OutputService11TestShapeOutputShape, error) {
	req, out := c.OutputService11TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService11TestShapeOutputService11TestCaseOperation1Input struct {
	metadataOutputService11TestShapeOutputService11TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService11TestShapeOutputService11TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService11TestShapeOutputShape struct {
	Char *string `location:"header" locationName:"x-char" type:"character"`

	Double *float64 `location:"header" locationName:"x-double" type:"double"`

	FalseBool *bool `location:"header" locationName:"x-false-bool" type:"boolean"`

	Float *float64 `location:"header" locationName:"x-float" type:"float"`

	Integer *int64 `location:"header" locationName:"x-int" type:"integer"`

	Long *int64 `location:"header" locationName:"x-long" type:"long"`

	Str *string `location:"header" locationName:"x-str" type:"string"`

	Timestamp *time.Time `location:"header" locationName:"x-timestamp" type:"timestamp" timestampFormat:"iso8601"`

	TrueBool *bool `location:"header" locationName:"x-true-bool" type:"boolean"`

	metadataOutputService11TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService11TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService12ProtocolTest struct {
	*aws.Service
}

// New returns a new OutputService12ProtocolTest client.
func NewOutputService12ProtocolTest(config *aws.Config) *OutputService12ProtocolTest {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "outputservice12protocoltest",
		APIVersion:  "",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	return &OutputService12ProtocolTest{service}
}

// newRequest creates a new request for a OutputService12ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService12ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	return req
}

const opOutputService12TestCaseOperation1 = "OperationName"

// OutputService12TestCaseOperation1Request generates a request for the OutputService12TestCaseOperation1 operation.
func (c *OutputService12ProtocolTest) OutputService12TestCaseOperation1Request(input *OutputService12TestShapeOutputService12TestCaseOperation1Input) (req *aws.Request, output *OutputService12TestShapeOutputShape) {
	op := &aws.Operation{
		Name: opOutputService12TestCaseOperation1,
	}

	if input == nil {
		input = &OutputService12TestShapeOutputService12TestCaseOperation1Input{}
	}

	req = c.newRequest(op, input, output)
	output = &OutputService12TestShapeOutputShape{}
	req.Data = output
	return
}

func (c *OutputService12ProtocolTest) OutputService12TestCaseOperation1(input *OutputService12TestShapeOutputService12TestCaseOperation1Input) (*OutputService12TestShapeOutputShape, error) {
	req, out := c.OutputService12TestCaseOperation1Request(input)
	err := req.Send()
	return out, err
}

type OutputService12TestShapeOutputService12TestCaseOperation1Input struct {
	metadataOutputService12TestShapeOutputService12TestCaseOperation1Input `json:"-" xml:"-"`
}

type metadataOutputService12TestShapeOutputService12TestCaseOperation1Input struct {
	SDKShapeTraits bool `type:"structure"`
}

type OutputService12TestShapeOutputShape struct {
	String *string `type:"string"`

	metadataOutputService12TestShapeOutputShape `json:"-" xml:"-"`
}

type metadataOutputService12TestShapeOutputShape struct {
	SDKShapeTraits bool `type:"structure" payload:"String"`
}

//
// Tests begin here
//

func TestOutputService1ProtocolTestScalarMembersCase1(t *testing.T) {
	svc := NewOutputService1ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Str>myname</Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><Timestamp>2015-01-25T08:00:00Z</Timestamp></OperationNameResponse>"))
	req, out := svc.OutputService1TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("ImaHeader", "test")
	req.HTTPResponse.Header.Set("X-Foo", "abc")

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "a", *out.Char)
	assert.Equal(t, 1.3, *out.Double)
	assert.Equal(t, false, *out.FalseBool)
	assert.Equal(t, 1.2, *out.Float)
	assert.Equal(t, "test", *out.ImaHeader)
	assert.Equal(t, "abc", *out.ImaHeaderLocation)
	assert.Equal(t, int64(200), *out.Long)
	assert.Equal(t, int64(123), *out.Num)
	assert.Equal(t, "myname", *out.Str)
	assert.Equal(t, time.Unix(1.4221728e+09, 0).UTC().String(), out.Timestamp.String())
	assert.Equal(t, true, *out.TrueBool)

}

func TestOutputService1ProtocolTestScalarMembersCase2(t *testing.T) {
	svc := NewOutputService1ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Str></Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><Timestamp>2015-01-25T08:00:00Z</Timestamp></OperationNameResponse>"))
	req, out := svc.OutputService1TestCaseOperation2Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("ImaHeader", "test")
	req.HTTPResponse.Header.Set("X-Foo", "abc")

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "a", *out.Char)
	assert.Equal(t, 1.3, *out.Double)
	assert.Equal(t, false, *out.FalseBool)
	assert.Equal(t, 1.2, *out.Float)
	assert.Equal(t, "test", *out.ImaHeader)
	assert.Equal(t, "abc", *out.ImaHeaderLocation)
	assert.Equal(t, int64(200), *out.Long)
	assert.Equal(t, int64(123), *out.Num)
	assert.Equal(t, "", *out.Str)
	assert.Equal(t, time.Unix(1.4221728e+09, 0).UTC().String(), out.Timestamp.String())
	assert.Equal(t, true, *out.TrueBool)

}

func TestOutputService2ProtocolTestBlobCase1(t *testing.T) {
	svc := NewOutputService2ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><Blob>dmFsdWU=</Blob></OperationNameResult>"))
	req, out := svc.OutputService2TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "value", string(out.Blob))

}

func TestOutputService3ProtocolTestListsCase1(t *testing.T) {
	svc := NewOutputService3ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><ListMember><member>abc</member><member>123</member></ListMember></OperationNameResult>"))
	req, out := svc.OutputService3TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "abc", *out.ListMember[0])
	assert.Equal(t, "123", *out.ListMember[1])

}

func TestOutputService4ProtocolTestListWithCustomMemberNameCase1(t *testing.T) {
	svc := NewOutputService4ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><ListMember><item>abc</item><item>123</item></ListMember></OperationNameResult>"))
	req, out := svc.OutputService4TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "abc", *out.ListMember[0])
	assert.Equal(t, "123", *out.ListMember[1])

}

func TestOutputService5ProtocolTestFlattenedListCase1(t *testing.T) {
	svc := NewOutputService5ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><ListMember>abc</ListMember><ListMember>123</ListMember></OperationNameResult>"))
	req, out := svc.OutputService5TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "abc", *out.ListMember[0])
	assert.Equal(t, "123", *out.ListMember[1])

}

func TestOutputService6ProtocolTestNormalMapCase1(t *testing.T) {
	svc := NewOutputService6ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><Map><entry><key>qux</key><value><foo>bar</foo></value></entry><entry><key>baz</key><value><foo>bam</foo></value></entry></Map></OperationNameResult>"))
	req, out := svc.OutputService6TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "bam", *out.Map["baz"].Foo)
	assert.Equal(t, "bar", *out.Map["qux"].Foo)

}

func TestOutputService7ProtocolTestFlattenedMapCase1(t *testing.T) {
	svc := NewOutputService7ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><Map><key>qux</key><value>bar</value></Map><Map><key>baz</key><value>bam</value></Map></OperationNameResult>"))
	req, out := svc.OutputService7TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "bam", *out.Map["baz"])
	assert.Equal(t, "bar", *out.Map["qux"])

}

func TestOutputService8ProtocolTestNamedMapCase1(t *testing.T) {
	svc := NewOutputService8ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResult><Map><entry><foo>qux</foo><bar>bar</bar></entry><entry><foo>baz</foo><bar>bam</bar></entry></Map></OperationNameResult>"))
	req, out := svc.OutputService8TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "bam", *out.Map["baz"])
	assert.Equal(t, "bar", *out.Map["qux"])

}

func TestOutputService9ProtocolTestXMLPayloadCase1(t *testing.T) {
	svc := NewOutputService9ProtocolTest(nil)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Foo>abc</Foo></OperationNameResponse>"))
	req, out := svc.OutputService9TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("X-Foo", "baz")

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "abc", *out.Data.Foo)
	assert.Equal(t, "baz", *out.Header)

}

func TestOutputService10ProtocolTestStreamingPayloadCase1(t *testing.T) {
	svc := NewOutputService10ProtocolTest(nil)

	buf := bytes.NewReader([]byte("abc"))
	req, out := svc.OutputService10TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "abc", string(out.Stream))

}

func TestOutputService11ProtocolTestScalarMembersInHeadersCase1(t *testing.T) {
	svc := NewOutputService11ProtocolTest(nil)

	buf := bytes.NewReader([]byte(""))
	req, out := svc.OutputService11TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("x-char", "a")
	req.HTTPResponse.Header.Set("x-double", "1.5")
	req.HTTPResponse.Header.Set("x-false-bool", "false")
	req.HTTPResponse.Header.Set("x-float", "1.5")
	req.HTTPResponse.Header.Set("x-int", "1")
	req.HTTPResponse.Header.Set("x-long", "100")
	req.HTTPResponse.Header.Set("x-str", "string")
	req.HTTPResponse.Header.Set("x-timestamp", "Sun, 25 Jan 2015 08:00:00 GMT")
	req.HTTPResponse.Header.Set("x-true-bool", "true")

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "a", *out.Char)
	assert.Equal(t, 1.5, *out.Double)
	assert.Equal(t, false, *out.FalseBool)
	assert.Equal(t, 1.5, *out.Float)
	assert.Equal(t, int64(1), *out.Integer)
	assert.Equal(t, int64(100), *out.Long)
	assert.Equal(t, "string", *out.Str)
	assert.Equal(t, time.Unix(1.4221728e+09, 0).UTC().String(), out.Timestamp.String())
	assert.Equal(t, true, *out.TrueBool)

}

func TestOutputService12ProtocolTestStringCase1(t *testing.T) {
	svc := NewOutputService12ProtocolTest(nil)

	buf := bytes.NewReader([]byte("operation result string"))
	req, out := svc.OutputService12TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req)
	restxml.Unmarshal(req)
	assert.NoError(t, req.Error)

	// assert response
	assert.NotNil(t, out) // ensure out variable is used
	assert.Equal(t, "operation result string", *out.String)

}
