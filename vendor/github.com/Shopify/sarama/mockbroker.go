package sarama

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	expectationTimeout = 500 * time.Millisecond
)

type GSSApiHandlerFunc func([]byte) []byte

type requestHandlerFunc func(req *request) (res encoderWithHeader)

// RequestNotifierFunc is invoked when a mock broker processes a request successfully
// and will provides the number of bytes read and written.
type RequestNotifierFunc func(bytesRead, bytesWritten int)

// MockBroker is a mock Kafka broker that is used in unit tests. It is exposed
// to facilitate testing of higher level or specialized consumers and producers
// built on top of Sarama. Note that it does not 'mimic' the Kafka API protocol,
// but rather provides a facility to do that. It takes care of the TCP
// transport, request unmarshaling, response marshaling, and makes it the test
// writer responsibility to program correct according to the Kafka API protocol
// MockBroker behaviour.
//
// MockBroker is implemented as a TCP server listening on a kernel-selected
// localhost port that can accept many connections. It reads Kafka requests
// from that connection and returns responses programmed by the SetHandlerByMap
// function. If a MockBroker receives a request that it has no programmed
// response for, then it returns nothing and the request times out.
//
// A set of MockRequest builders to define mappings used by MockBroker is
// provided by Sarama. But users can develop MockRequests of their own and use
// them along with or instead of the standard ones.
//
// When running tests with MockBroker it is strongly recommended to specify
// a timeout to `go test` so that if the broker hangs waiting for a response,
// the test panics.
//
// It is not necessary to prefix message length or correlation ID to your
// response bytes, the server does that automatically as a convenience.
type MockBroker struct {
	brokerID      int32
	port          int32
	closing       chan none
	stopper       chan none
	expectations  chan encoderWithHeader
	listener      net.Listener
	t             TestReporter
	latency       time.Duration
	handler       requestHandlerFunc
	notifier      RequestNotifierFunc
	history       []RequestResponse
	lock          sync.Mutex
	gssApiHandler GSSApiHandlerFunc
}

// RequestResponse represents a Request/Response pair processed by MockBroker.
type RequestResponse struct {
	Request  protocolBody
	Response encoder
}

// SetLatency makes broker pause for the specified period every time before
// replying.
func (b *MockBroker) SetLatency(latency time.Duration) {
	b.latency = latency
}

// SetHandlerByMap defines mapping of Request types to MockResponses. When a
// request is received by the broker, it looks up the request type in the map
// and uses the found MockResponse instance to generate an appropriate reply.
// If the request type is not found in the map then nothing is sent.
func (b *MockBroker) SetHandlerByMap(handlerMap map[string]MockResponse) {
	b.setHandler(func(req *request) (res encoderWithHeader) {
		reqTypeName := reflect.TypeOf(req.body).Elem().Name()
		mockResponse := handlerMap[reqTypeName]
		if mockResponse == nil {
			return nil
		}
		return mockResponse.For(req.body)
	})
}

// SetNotifier set a function that will get invoked whenever a request has been
// processed successfully and will provide the number of bytes read and written
func (b *MockBroker) SetNotifier(notifier RequestNotifierFunc) {
	b.lock.Lock()
	b.notifier = notifier
	b.lock.Unlock()
}

// BrokerID returns broker ID assigned to the broker.
func (b *MockBroker) BrokerID() int32 {
	return b.brokerID
}

// History returns a slice of RequestResponse pairs in the order they were
// processed by the broker. Note that in case of multiple connections to the
// broker the order expected by a test can be different from the order recorded
// in the history, unless some synchronization is implemented in the test.
func (b *MockBroker) History() []RequestResponse {
	b.lock.Lock()
	history := make([]RequestResponse, len(b.history))
	copy(history, b.history)
	b.lock.Unlock()
	return history
}

// Port returns the TCP port number the broker is listening for requests on.
func (b *MockBroker) Port() int32 {
	return b.port
}

// Addr returns the broker connection string in the form "<address>:<port>".
func (b *MockBroker) Addr() string {
	return b.listener.Addr().String()
}

// Close terminates the broker blocking until it stops internal goroutines and
// releases all resources.
func (b *MockBroker) Close() {
	close(b.expectations)
	if len(b.expectations) > 0 {
		buf := bytes.NewBufferString(fmt.Sprintf("mockbroker/%d: not all expectations were satisfied! Still waiting on:\n", b.BrokerID()))
		for e := range b.expectations {
			_, _ = buf.WriteString(spew.Sdump(e))
		}
		b.t.Error(buf.String())
	}
	close(b.closing)
	<-b.stopper
}

// setHandler sets the specified function as the request handler. Whenever
// a mock broker reads a request from the wire it passes the request to the
// function and sends back whatever the handler function returns.
func (b *MockBroker) setHandler(handler requestHandlerFunc) {
	b.lock.Lock()
	b.handler = handler
	b.lock.Unlock()
}

func (b *MockBroker) serverLoop() {
	defer close(b.stopper)
	var err error
	var conn net.Conn

	go func() {
		<-b.closing
		err := b.listener.Close()
		if err != nil {
			b.t.Error(err)
		}
	}()

	wg := &sync.WaitGroup{}
	i := 0
	for conn, err = b.listener.Accept(); err == nil; conn, err = b.listener.Accept() {
		wg.Add(1)
		go b.handleRequests(conn, i, wg)
		i++
	}
	wg.Wait()
	Logger.Printf("*** mockbroker/%d: listener closed, err=%v", b.BrokerID(), err)
}

func (b *MockBroker) SetGSSAPIHandler(handler GSSApiHandlerFunc) {
	b.gssApiHandler = handler
}

func (b *MockBroker) readToBytes(r io.Reader) ([]byte, error) {
	var (
		bytesRead   int
		lengthBytes = make([]byte, 4)
	)

	if _, err := io.ReadFull(r, lengthBytes); err != nil {
		return nil, err
	}

	bytesRead += len(lengthBytes)
	length := int32(binary.BigEndian.Uint32(lengthBytes))

	if length <= 4 || length > MaxRequestSize {
		return nil, PacketDecodingError{fmt.Sprintf("message of length %d too large or too small", length)}
	}

	encodedReq := make([]byte, length)
	if _, err := io.ReadFull(r, encodedReq); err != nil {
		return nil, err
	}

	bytesRead += len(encodedReq)

	fullBytes := append(lengthBytes, encodedReq...)

	return fullBytes, nil
}

func (b *MockBroker) isGSSAPI(buffer []byte) bool {
	return buffer[4] == 0x60 || bytes.Equal(buffer[4:6], []byte{0x05, 0x04})
}

func (b *MockBroker) handleRequests(conn io.ReadWriteCloser, idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		_ = conn.Close()
	}()
	Logger.Printf("*** mockbroker/%d/%d: connection opened", b.BrokerID(), idx)
	var err error

	abort := make(chan none)
	defer close(abort)
	go func() {
		select {
		case <-b.closing:
			_ = conn.Close()
		case <-abort:
		}
	}()

	var bytesWritten int
	var bytesRead int
	for {
		buffer, err := b.readToBytes(conn)
		if err != nil {
			Logger.Printf("*** mockbroker/%d/%d: invalid request: err=%+v, %+v", b.brokerID, idx, err, spew.Sdump(buffer))
			b.serverError(err)
			break
		}

		bytesWritten = 0
		if !b.isGSSAPI(buffer) {
			req, br, err := decodeRequest(bytes.NewReader(buffer))
			bytesRead = br
			if err != nil {
				Logger.Printf("*** mockbroker/%d/%d: invalid request: err=%+v, %+v", b.brokerID, idx, err, spew.Sdump(req))
				b.serverError(err)
				break
			}

			if b.latency > 0 {
				time.Sleep(b.latency)
			}

			b.lock.Lock()
			res := b.handler(req)
			b.history = append(b.history, RequestResponse{req.body, res})
			b.lock.Unlock()

			if res == nil {
				Logger.Printf("*** mockbroker/%d/%d: ignored %v", b.brokerID, idx, spew.Sdump(req))
				continue
			}
			Logger.Printf("*** mockbroker/%d/%d: served %v -> %v", b.brokerID, idx, req, res)

			encodedRes, err := encode(res, nil)
			if err != nil {
				b.serverError(err)
				break
			}
			if len(encodedRes) == 0 {
				b.lock.Lock()
				if b.notifier != nil {
					b.notifier(bytesRead, 0)
				}
				b.lock.Unlock()
				continue
			}

			resHeader := b.encodeHeader(res.headerVersion(), req.correlationID, uint32(len(encodedRes)))
			if _, err = conn.Write(resHeader); err != nil {
				b.serverError(err)
				break
			}
			if _, err = conn.Write(encodedRes); err != nil {
				b.serverError(err)
				break
			}
			bytesWritten = len(resHeader) + len(encodedRes)
		} else {
			// GSSAPI is not part of kafka protocol, but is supported for authentication proposes.
			// Don't support history for this kind of request as is only used for test GSSAPI authentication mechanism
			b.lock.Lock()
			res := b.gssApiHandler(buffer)
			b.lock.Unlock()
			if res == nil {
				Logger.Printf("*** mockbroker/%d/%d: ignored %v", b.brokerID, idx, spew.Sdump(buffer))
				continue
			}
			if _, err = conn.Write(res); err != nil {
				b.serverError(err)
				break
			}
			bytesWritten = len(res)
		}

		b.lock.Lock()
		if b.notifier != nil {
			b.notifier(bytesRead, bytesWritten)
		}
		b.lock.Unlock()
	}
	Logger.Printf("*** mockbroker/%d/%d: connection closed, err=%v", b.BrokerID(), idx, err)
}

func (b *MockBroker) encodeHeader(headerVersion int16, correlationId int32, payloadLength uint32) []byte {
	headerLength := uint32(8)

	if headerVersion >= 1 {
		headerLength = 9
	}

	resHeader := make([]byte, headerLength)
	binary.BigEndian.PutUint32(resHeader, payloadLength+headerLength-4)
	binary.BigEndian.PutUint32(resHeader[4:], uint32(correlationId))

	if headerVersion >= 1 {
		binary.PutUvarint(resHeader[8:], 0)
	}

	return resHeader
}

func (b *MockBroker) defaultRequestHandler(req *request) (res encoderWithHeader) {
	select {
	case res, ok := <-b.expectations:
		if !ok {
			return nil
		}
		return res
	case <-time.After(expectationTimeout):
		return nil
	}
}

func (b *MockBroker) serverError(err error) {
	isConnectionClosedError := false
	if _, ok := err.(*net.OpError); ok {
		isConnectionClosedError = true
	} else if err == io.EOF {
		isConnectionClosedError = true
	} else if err.Error() == "use of closed network connection" {
		isConnectionClosedError = true
	}

	if isConnectionClosedError {
		return
	}

	b.t.Errorf(err.Error())
}

// NewMockBroker launches a fake Kafka broker. It takes a TestReporter as provided by the
// test framework and a channel of responses to use.  If an error occurs it is
// simply logged to the TestReporter and the broker exits.
func NewMockBroker(t TestReporter, brokerID int32) *MockBroker {
	return NewMockBrokerAddr(t, brokerID, "localhost:0")
}

// NewMockBrokerAddr behaves like newMockBroker but listens on the address you give
// it rather than just some ephemeral port.
func NewMockBrokerAddr(t TestReporter, brokerID int32, addr string) *MockBroker {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return NewMockBrokerListener(t, brokerID, listener)
}

// NewMockBrokerListener behaves like newMockBrokerAddr but accepts connections on the listener specified.
func NewMockBrokerListener(t TestReporter, brokerID int32, listener net.Listener) *MockBroker {
	var err error

	broker := &MockBroker{
		closing:      make(chan none),
		stopper:      make(chan none),
		t:            t,
		brokerID:     brokerID,
		expectations: make(chan encoderWithHeader, 512),
		listener:     listener,
	}
	broker.handler = broker.defaultRequestHandler

	Logger.Printf("*** mockbroker/%d listening on %s\n", brokerID, broker.listener.Addr().String())
	_, portStr, err := net.SplitHostPort(broker.listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	broker.port = int32(tmp)

	go broker.serverLoop()

	return broker
}

func (b *MockBroker) Returns(e encoderWithHeader) {
	b.expectations <- e
}
