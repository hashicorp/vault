package sarama

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

// Broker represents a single Kafka broker connection. All operations on this object are entirely concurrency-safe.
type Broker struct {
	conf *Config
	rack *string

	id            int32
	addr          string
	correlationID int32
	conn          net.Conn
	connErr       error
	lock          sync.Mutex
	opened        int32
	responses     chan responsePromise
	done          chan bool

	registeredMetrics []string

	incomingByteRate       metrics.Meter
	requestRate            metrics.Meter
	requestSize            metrics.Histogram
	requestLatency         metrics.Histogram
	outgoingByteRate       metrics.Meter
	responseRate           metrics.Meter
	responseSize           metrics.Histogram
	requestsInFlight       metrics.Counter
	brokerIncomingByteRate metrics.Meter
	brokerRequestRate      metrics.Meter
	brokerRequestSize      metrics.Histogram
	brokerRequestLatency   metrics.Histogram
	brokerOutgoingByteRate metrics.Meter
	brokerResponseRate     metrics.Meter
	brokerResponseSize     metrics.Histogram
	brokerRequestsInFlight metrics.Counter

	kerberosAuthenticator GSSAPIKerberosAuth
}

// SASLMechanism specifies the SASL mechanism the client uses to authenticate with the broker
type SASLMechanism string

const (
	// SASLTypeOAuth represents the SASL/OAUTHBEARER mechanism (Kafka 2.0.0+)
	SASLTypeOAuth = "OAUTHBEARER"
	// SASLTypePlaintext represents the SASL/PLAIN mechanism
	SASLTypePlaintext = "PLAIN"
	// SASLTypeSCRAMSHA256 represents the SCRAM-SHA-256 mechanism.
	SASLTypeSCRAMSHA256 = "SCRAM-SHA-256"
	// SASLTypeSCRAMSHA512 represents the SCRAM-SHA-512 mechanism.
	SASLTypeSCRAMSHA512 = "SCRAM-SHA-512"
	SASLTypeGSSAPI      = "GSSAPI"
	// SASLHandshakeV0 is v0 of the Kafka SASL handshake protocol. Client and
	// server negotiate SASL auth using opaque packets.
	SASLHandshakeV0 = int16(0)
	// SASLHandshakeV1 is v1 of the Kafka SASL handshake protocol. Client and
	// server negotiate SASL by wrapping tokens with Kafka protocol headers.
	SASLHandshakeV1 = int16(1)
	// SASLExtKeyAuth is the reserved extension key name sent as part of the
	// SASL/OAUTHBEARER initial client response
	SASLExtKeyAuth = "auth"
)

// AccessToken contains an access token used to authenticate a
// SASL/OAUTHBEARER client along with associated metadata.
type AccessToken struct {
	// Token is the access token payload.
	Token string
	// Extensions is a optional map of arbitrary key-value pairs that can be
	// sent with the SASL/OAUTHBEARER initial client response. These values are
	// ignored by the SASL server if they are unexpected. This feature is only
	// supported by Kafka >= 2.1.0.
	Extensions map[string]string
}

// AccessTokenProvider is the interface that encapsulates how implementors
// can generate access tokens for Kafka broker authentication.
type AccessTokenProvider interface {
	// Token returns an access token. The implementation should ensure token
	// reuse so that multiple calls at connect time do not create multiple
	// tokens. The implementation should also periodically refresh the token in
	// order to guarantee that each call returns an unexpired token.  This
	// method should not block indefinitely--a timeout error should be returned
	// after a short period of inactivity so that the broker connection logic
	// can log debugging information and retry.
	Token() (*AccessToken, error)
}

// SCRAMClient is a an interface to a SCRAM
// client implementation.
type SCRAMClient interface {
	// Begin prepares the client for the SCRAM exchange
	// with the server with a user name and a password
	Begin(userName, password, authzID string) error
	// Step steps client through the SCRAM exchange. It is
	// called repeatedly until it errors or `Done` returns true.
	Step(challenge string) (response string, err error)
	// Done should return true when the SCRAM conversation
	// is over.
	Done() bool
}

type responsePromise struct {
	requestTime   time.Time
	correlationID int32
	headerVersion int16
	packets       chan []byte
	errors        chan error
}

// NewBroker creates and returns a Broker targeting the given host:port address.
// This does not attempt to actually connect, you have to call Open() for that.
func NewBroker(addr string) *Broker {
	return &Broker{id: -1, addr: addr}
}

// Open tries to connect to the Broker if it is not already connected or connecting, but does not block
// waiting for the connection to complete. This means that any subsequent operations on the broker will
// block waiting for the connection to succeed or fail. To get the effect of a fully synchronous Open call,
// follow it by a call to Connected(). The only errors Open will return directly are ConfigurationError or
// AlreadyConnected. If conf is nil, the result of NewConfig() is used.
func (b *Broker) Open(conf *Config) error {
	if !atomic.CompareAndSwapInt32(&b.opened, 0, 1) {
		return ErrAlreadyConnected
	}

	if conf == nil {
		conf = NewConfig()
	}

	err := conf.Validate()
	if err != nil {
		return err
	}

	b.lock.Lock()

	go withRecover(func() {
		defer b.lock.Unlock()

		dialer := conf.getDialer()
		b.conn, b.connErr = dialer.Dial("tcp", b.addr)
		if b.connErr != nil {
			Logger.Printf("Failed to connect to broker %s: %s\n", b.addr, b.connErr)
			b.conn = nil
			atomic.StoreInt32(&b.opened, 0)
			return
		}
		if conf.Net.TLS.Enable {
			b.conn = tls.Client(b.conn, validServerNameTLS(b.addr, conf.Net.TLS.Config))
		}

		b.conn = newBufConn(b.conn)
		b.conf = conf

		// Create or reuse the global metrics shared between brokers
		b.incomingByteRate = metrics.GetOrRegisterMeter("incoming-byte-rate", conf.MetricRegistry)
		b.requestRate = metrics.GetOrRegisterMeter("request-rate", conf.MetricRegistry)
		b.requestSize = getOrRegisterHistogram("request-size", conf.MetricRegistry)
		b.requestLatency = getOrRegisterHistogram("request-latency-in-ms", conf.MetricRegistry)
		b.outgoingByteRate = metrics.GetOrRegisterMeter("outgoing-byte-rate", conf.MetricRegistry)
		b.responseRate = metrics.GetOrRegisterMeter("response-rate", conf.MetricRegistry)
		b.responseSize = getOrRegisterHistogram("response-size", conf.MetricRegistry)
		b.requestsInFlight = metrics.GetOrRegisterCounter("requests-in-flight", conf.MetricRegistry)
		// Do not gather metrics for seeded broker (only used during bootstrap) because they share
		// the same id (-1) and are already exposed through the global metrics above
		if b.id >= 0 {
			b.registerMetrics()
		}

		if conf.Net.SASL.Enable {
			b.connErr = b.authenticateViaSASL()

			if b.connErr != nil {
				err = b.conn.Close()
				if err == nil {
					Logger.Printf("Closed connection to broker %s\n", b.addr)
				} else {
					Logger.Printf("Error while closing connection to broker %s: %s\n", b.addr, err)
				}
				b.conn = nil
				atomic.StoreInt32(&b.opened, 0)
				return
			}
		}

		b.done = make(chan bool)
		b.responses = make(chan responsePromise, b.conf.Net.MaxOpenRequests-1)

		if b.id >= 0 {
			Logger.Printf("Connected to broker at %s (registered as #%d)\n", b.addr, b.id)
		} else {
			Logger.Printf("Connected to broker at %s (unregistered)\n", b.addr)
		}
		go withRecover(b.responseReceiver)
	})

	return nil
}

// Connected returns true if the broker is connected and false otherwise. If the broker is not
// connected but it had tried to connect, the error from that connection attempt is also returned.
func (b *Broker) Connected() (bool, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.conn != nil, b.connErr
}

//Close closes the broker resources
func (b *Broker) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.conn == nil {
		return ErrNotConnected
	}

	close(b.responses)
	<-b.done

	err := b.conn.Close()

	b.conn = nil
	b.connErr = nil
	b.done = nil
	b.responses = nil

	b.unregisterMetrics()

	if err == nil {
		Logger.Printf("Closed connection to broker %s\n", b.addr)
	} else {
		Logger.Printf("Error while closing connection to broker %s: %s\n", b.addr, err)
	}

	atomic.StoreInt32(&b.opened, 0)

	return err
}

// ID returns the broker ID retrieved from Kafka's metadata, or -1 if that is not known.
func (b *Broker) ID() int32 {
	return b.id
}

// Addr returns the broker address as either retrieved from Kafka's metadata or passed to NewBroker.
func (b *Broker) Addr() string {
	return b.addr
}

// Rack returns the broker's rack as retrieved from Kafka's metadata or the
// empty string if it is not known.  The returned value corresponds to the
// broker's broker.rack configuration setting.  Requires protocol version to be
// at least v0.10.0.0.
func (b *Broker) Rack() string {
	if b.rack == nil {
		return ""
	}
	return *b.rack
}

//GetMetadata send a metadata request and returns a metadata response or error
func (b *Broker) GetMetadata(request *MetadataRequest) (*MetadataResponse, error) {
	response := new(MetadataResponse)

	err := b.sendAndReceive(request, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//GetConsumerMetadata send a consumer metadata request and returns a consumer metadata response or error
func (b *Broker) GetConsumerMetadata(request *ConsumerMetadataRequest) (*ConsumerMetadataResponse, error) {
	response := new(ConsumerMetadataResponse)

	err := b.sendAndReceive(request, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//FindCoordinator sends a find coordinate request and returns a response or error
func (b *Broker) FindCoordinator(request *FindCoordinatorRequest) (*FindCoordinatorResponse, error) {
	response := new(FindCoordinatorResponse)

	err := b.sendAndReceive(request, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//GetAvailableOffsets return an offset response or error
func (b *Broker) GetAvailableOffsets(request *OffsetRequest) (*OffsetResponse, error) {
	response := new(OffsetResponse)

	err := b.sendAndReceive(request, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//Produce returns a produce response or error
func (b *Broker) Produce(request *ProduceRequest) (*ProduceResponse, error) {
	var (
		response *ProduceResponse
		err      error
	)

	if request.RequiredAcks == NoResponse {
		err = b.sendAndReceive(request, nil)
	} else {
		response = new(ProduceResponse)
		err = b.sendAndReceive(request, response)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

//Fetch returns a FetchResponse or error
func (b *Broker) Fetch(request *FetchRequest) (*FetchResponse, error) {
	response := new(FetchResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//CommitOffset return an Offset commit response or error
func (b *Broker) CommitOffset(request *OffsetCommitRequest) (*OffsetCommitResponse, error) {
	response := new(OffsetCommitResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//FetchOffset returns an offset fetch response or error
func (b *Broker) FetchOffset(request *OffsetFetchRequest) (*OffsetFetchResponse, error) {
	response := new(OffsetFetchResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//JoinGroup returns a join group response or error
func (b *Broker) JoinGroup(request *JoinGroupRequest) (*JoinGroupResponse, error) {
	response := new(JoinGroupResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//SyncGroup returns a sync group response or error
func (b *Broker) SyncGroup(request *SyncGroupRequest) (*SyncGroupResponse, error) {
	response := new(SyncGroupResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//LeaveGroup return a leave group response or error
func (b *Broker) LeaveGroup(request *LeaveGroupRequest) (*LeaveGroupResponse, error) {
	response := new(LeaveGroupResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//Heartbeat returns a heartbeat response or error
func (b *Broker) Heartbeat(request *HeartbeatRequest) (*HeartbeatResponse, error) {
	response := new(HeartbeatResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//ListGroups return a list group response or error
func (b *Broker) ListGroups(request *ListGroupsRequest) (*ListGroupsResponse, error) {
	response := new(ListGroupsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DescribeGroups return describe group response or error
func (b *Broker) DescribeGroups(request *DescribeGroupsRequest) (*DescribeGroupsResponse, error) {
	response := new(DescribeGroupsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//ApiVersions return api version response or error
func (b *Broker) ApiVersions(request *ApiVersionsRequest) (*ApiVersionsResponse, error) {
	response := new(ApiVersionsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//CreateTopics send a create topic request and returns create topic response
func (b *Broker) CreateTopics(request *CreateTopicsRequest) (*CreateTopicsResponse, error) {
	response := new(CreateTopicsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DeleteTopics sends a delete topic request and returns delete topic response
func (b *Broker) DeleteTopics(request *DeleteTopicsRequest) (*DeleteTopicsResponse, error) {
	response := new(DeleteTopicsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//CreatePartitions sends a create partition request and returns create
//partitions response or error
func (b *Broker) CreatePartitions(request *CreatePartitionsRequest) (*CreatePartitionsResponse, error) {
	response := new(CreatePartitionsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//AlterPartitionReassignments sends a alter partition reassignments request and
//returns alter partition reassignments response
func (b *Broker) AlterPartitionReassignments(request *AlterPartitionReassignmentsRequest) (*AlterPartitionReassignmentsResponse, error) {
	response := new(AlterPartitionReassignmentsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//ListPartitionReassignments sends a list partition reassignments request and
//returns list partition reassignments response
func (b *Broker) ListPartitionReassignments(request *ListPartitionReassignmentsRequest) (*ListPartitionReassignmentsResponse, error) {
	response := new(ListPartitionReassignmentsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DeleteRecords send a request to delete records and return delete record
//response or error
func (b *Broker) DeleteRecords(request *DeleteRecordsRequest) (*DeleteRecordsResponse, error) {
	response := new(DeleteRecordsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DescribeAcls sends a describe acl request and returns a response or error
func (b *Broker) DescribeAcls(request *DescribeAclsRequest) (*DescribeAclsResponse, error) {
	response := new(DescribeAclsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//CreateAcls sends a create acl request and returns a response or error
func (b *Broker) CreateAcls(request *CreateAclsRequest) (*CreateAclsResponse, error) {
	response := new(CreateAclsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DeleteAcls sends a delete acl request and returns a response or error
func (b *Broker) DeleteAcls(request *DeleteAclsRequest) (*DeleteAclsResponse, error) {
	response := new(DeleteAclsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//InitProducerID sends an init producer request and returns a response or error
func (b *Broker) InitProducerID(request *InitProducerIDRequest) (*InitProducerIDResponse, error) {
	response := new(InitProducerIDResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//AddPartitionsToTxn send a request to add partition to txn and returns
//a response or error
func (b *Broker) AddPartitionsToTxn(request *AddPartitionsToTxnRequest) (*AddPartitionsToTxnResponse, error) {
	response := new(AddPartitionsToTxnResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//AddOffsetsToTxn sends a request to add offsets to txn and returns a response
//or error
func (b *Broker) AddOffsetsToTxn(request *AddOffsetsToTxnRequest) (*AddOffsetsToTxnResponse, error) {
	response := new(AddOffsetsToTxnResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//EndTxn sends a request to end txn and returns a response or error
func (b *Broker) EndTxn(request *EndTxnRequest) (*EndTxnResponse, error) {
	response := new(EndTxnResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//TxnOffsetCommit sends a request to commit transaction offsets and returns
//a response or error
func (b *Broker) TxnOffsetCommit(request *TxnOffsetCommitRequest) (*TxnOffsetCommitResponse, error) {
	response := new(TxnOffsetCommitResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DescribeConfigs sends a request to describe config and returns a response or
//error
func (b *Broker) DescribeConfigs(request *DescribeConfigsRequest) (*DescribeConfigsResponse, error) {
	response := new(DescribeConfigsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//AlterConfigs sends a request to alter config and return a response or error
func (b *Broker) AlterConfigs(request *AlterConfigsRequest) (*AlterConfigsResponse, error) {
	response := new(AlterConfigsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DeleteGroups sends a request to delete groups and returns a response or error
func (b *Broker) DeleteGroups(request *DeleteGroupsRequest) (*DeleteGroupsResponse, error) {
	response := new(DeleteGroupsResponse)

	if err := b.sendAndReceive(request, response); err != nil {
		return nil, err
	}

	return response, nil
}

//DescribeLogDirs sends a request to get the broker's log dir paths and sizes
func (b *Broker) DescribeLogDirs(request *DescribeLogDirsRequest) (*DescribeLogDirsResponse, error) {
	response := new(DescribeLogDirsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// readFull ensures the conn ReadDeadline has been setup before making a
// call to io.ReadFull
func (b *Broker) readFull(buf []byte) (n int, err error) {
	if err := b.conn.SetReadDeadline(time.Now().Add(b.conf.Net.ReadTimeout)); err != nil {
		return 0, err
	}

	return io.ReadFull(b.conn, buf)
}

// write  ensures the conn WriteDeadline has been setup before making a
// call to conn.Write
func (b *Broker) write(buf []byte) (n int, err error) {
	if err := b.conn.SetWriteDeadline(time.Now().Add(b.conf.Net.WriteTimeout)); err != nil {
		return 0, err
	}

	return b.conn.Write(buf)
}

func (b *Broker) send(rb protocolBody, promiseResponse bool, responseHeaderVersion int16) (*responsePromise, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.conn == nil {
		if b.connErr != nil {
			return nil, b.connErr
		}
		return nil, ErrNotConnected
	}

	if !b.conf.Version.IsAtLeast(rb.requiredVersion()) {
		return nil, ErrUnsupportedVersion
	}

	req := &request{correlationID: b.correlationID, clientID: b.conf.ClientID, body: rb}
	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil {
		return nil, err
	}

	requestTime := time.Now()
	// Will be decremented in responseReceiver (except error or request with NoResponse)
	b.addRequestInFlightMetrics(1)
	bytes, err := b.write(buf)
	b.updateOutgoingCommunicationMetrics(bytes)
	if err != nil {
		b.addRequestInFlightMetrics(-1)
		return nil, err
	}
	b.correlationID++

	if !promiseResponse {
		// Record request latency without the response
		b.updateRequestLatencyAndInFlightMetrics(time.Since(requestTime))
		return nil, nil
	}

	promise := responsePromise{requestTime, req.correlationID, responseHeaderVersion, make(chan []byte), make(chan error)}
	b.responses <- promise

	return &promise, nil
}

func (b *Broker) sendAndReceive(req protocolBody, res protocolBody) error {
	responseHeaderVersion := int16(-1)
	if res != nil {
		responseHeaderVersion = res.headerVersion()
	}

	promise, err := b.send(req, res != nil, responseHeaderVersion)
	if err != nil {
		return err
	}

	if promise == nil {
		return nil
	}

	select {
	case buf := <-promise.packets:
		return versionedDecode(buf, res, req.version())
	case err = <-promise.errors:
		return err
	}
}

func (b *Broker) decode(pd packetDecoder, version int16) (err error) {
	b.id, err = pd.getInt32()
	if err != nil {
		return err
	}

	host, err := pd.getString()
	if err != nil {
		return err
	}

	port, err := pd.getInt32()
	if err != nil {
		return err
	}

	if version >= 1 {
		b.rack, err = pd.getNullableString()
		if err != nil {
			return err
		}
	}

	b.addr = net.JoinHostPort(host, fmt.Sprint(port))
	if _, _, err := net.SplitHostPort(b.addr); err != nil {
		return err
	}

	return nil
}

func (b *Broker) encode(pe packetEncoder, version int16) (err error) {
	host, portstr, err := net.SplitHostPort(b.addr)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return err
	}

	pe.putInt32(b.id)

	err = pe.putString(host)
	if err != nil {
		return err
	}

	pe.putInt32(int32(port))

	if version >= 1 {
		err = pe.putNullableString(b.rack)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Broker) responseReceiver() {
	var dead error

	for response := range b.responses {
		if dead != nil {
			// This was previously incremented in send() and
			// we are not calling updateIncomingCommunicationMetrics()
			b.addRequestInFlightMetrics(-1)
			response.errors <- dead
			continue
		}

		var headerLength = getHeaderLength(response.headerVersion)
		header := make([]byte, headerLength)

		bytesReadHeader, err := b.readFull(header)
		requestLatency := time.Since(response.requestTime)
		if err != nil {
			b.updateIncomingCommunicationMetrics(bytesReadHeader, requestLatency)
			dead = err
			response.errors <- err
			continue
		}

		decodedHeader := responseHeader{}
		err = versionedDecode(header, &decodedHeader, response.headerVersion)
		if err != nil {
			b.updateIncomingCommunicationMetrics(bytesReadHeader, requestLatency)
			dead = err
			response.errors <- err
			continue
		}
		if decodedHeader.correlationID != response.correlationID {
			b.updateIncomingCommunicationMetrics(bytesReadHeader, requestLatency)
			// TODO if decoded ID < cur ID, discard until we catch up
			// TODO if decoded ID > cur ID, save it so when cur ID catches up we have a response
			dead = PacketDecodingError{fmt.Sprintf("correlation ID didn't match, wanted %d, got %d", response.correlationID, decodedHeader.correlationID)}
			response.errors <- dead
			continue
		}

		buf := make([]byte, decodedHeader.length-int32(headerLength)+4)
		bytesReadBody, err := b.readFull(buf)
		b.updateIncomingCommunicationMetrics(bytesReadHeader+bytesReadBody, requestLatency)
		if err != nil {
			dead = err
			response.errors <- err
			continue
		}

		response.packets <- buf
	}
	close(b.done)
}

func getHeaderLength(headerVersion int16) int8 {
	if headerVersion < 1 {
		return 8
	} else {
		// header contains additional tagged field length (0), we don't support actual tags yet.
		return 9
	}
}

func (b *Broker) authenticateViaSASL() error {
	switch b.conf.Net.SASL.Mechanism {
	case SASLTypeOAuth:
		return b.sendAndReceiveSASLOAuth(b.conf.Net.SASL.TokenProvider)
	case SASLTypeSCRAMSHA256, SASLTypeSCRAMSHA512:
		return b.sendAndReceiveSASLSCRAMv1()
	case SASLTypeGSSAPI:
		return b.sendAndReceiveKerberos()
	default:
		return b.sendAndReceiveSASLPlainAuth()
	}
}

func (b *Broker) sendAndReceiveKerberos() error {
	b.kerberosAuthenticator.Config = &b.conf.Net.SASL.GSSAPI
	if b.kerberosAuthenticator.NewKerberosClientFunc == nil {
		b.kerberosAuthenticator.NewKerberosClientFunc = NewKerberosClient
	}
	return b.kerberosAuthenticator.Authorize(b)
}

func (b *Broker) sendAndReceiveSASLHandshake(saslType SASLMechanism, version int16) error {
	rb := &SaslHandshakeRequest{Mechanism: string(saslType), Version: version}

	req := &request{correlationID: b.correlationID, clientID: b.conf.ClientID, body: rb}
	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil {
		return err
	}

	requestTime := time.Now()
	// Will be decremented in updateIncomingCommunicationMetrics (except error)
	b.addRequestInFlightMetrics(1)
	bytes, err := b.write(buf)
	b.updateOutgoingCommunicationMetrics(bytes)
	if err != nil {
		b.addRequestInFlightMetrics(-1)
		Logger.Printf("Failed to send SASL handshake %s: %s\n", b.addr, err.Error())
		return err
	}
	b.correlationID++

	header := make([]byte, 8) // response header
	_, err = b.readFull(header)
	if err != nil {
		b.addRequestInFlightMetrics(-1)
		Logger.Printf("Failed to read SASL handshake header : %s\n", err.Error())
		return err
	}

	length := binary.BigEndian.Uint32(header[:4])
	payload := make([]byte, length-4)
	n, err := b.readFull(payload)
	if err != nil {
		b.addRequestInFlightMetrics(-1)
		Logger.Printf("Failed to read SASL handshake payload : %s\n", err.Error())
		return err
	}

	b.updateIncomingCommunicationMetrics(n+8, time.Since(requestTime))
	res := &SaslHandshakeResponse{}

	err = versionedDecode(payload, res, 0)
	if err != nil {
		Logger.Printf("Failed to parse SASL handshake : %s\n", err.Error())
		return err
	}

	if res.Err != ErrNoError {
		Logger.Printf("Invalid SASL Mechanism : %s\n", res.Err.Error())
		return res.Err
	}

	Logger.Print("Successful SASL handshake. Available mechanisms: ", res.EnabledMechanisms)
	return nil
}

// Kafka 0.10.x supported SASL PLAIN/Kerberos via KAFKA-3149 (KIP-43).
// Kafka 1.x.x onward added a SaslAuthenticate request/response message which
// wraps the SASL flow in the Kafka protocol, which allows for returning
// meaningful errors on authentication failure.
//
// In SASL Plain, Kafka expects the auth header to be in the following format
// Message format (from https://tools.ietf.org/html/rfc4616):
//
//   message   = [authzid] UTF8NUL authcid UTF8NUL passwd
//   authcid   = 1*SAFE ; MUST accept up to 255 octets
//   authzid   = 1*SAFE ; MUST accept up to 255 octets
//   passwd    = 1*SAFE ; MUST accept up to 255 octets
//   UTF8NUL   = %x00 ; UTF-8 encoded NUL character
//
//   SAFE      = UTF1 / UTF2 / UTF3 / UTF4
//                  ;; any UTF-8 encoded Unicode character except NUL
//
// With SASL v0 handshake and auth then:
// When credentials are valid, Kafka returns a 4 byte array of null characters.
// When credentials are invalid, Kafka closes the connection.
//
// With SASL v1 handshake and auth then:
// When credentials are invalid, Kafka replies with a SaslAuthenticate response
// containing an error code and message detailing the authentication failure.
func (b *Broker) sendAndReceiveSASLPlainAuth() error {
	// default to V0 to allow for backward compatibility when SASL is enabled
	// but not the handshake
	if b.conf.Net.SASL.Handshake {
		handshakeErr := b.sendAndReceiveSASLHandshake(SASLTypePlaintext, b.conf.Net.SASL.Version)
		if handshakeErr != nil {
			Logger.Printf("Error while performing SASL handshake %s\n", b.addr)
			return handshakeErr
		}
	}

	if b.conf.Net.SASL.Version == SASLHandshakeV1 {
		return b.sendAndReceiveV1SASLPlainAuth()
	}
	return b.sendAndReceiveV0SASLPlainAuth()
}

// sendAndReceiveV0SASLPlainAuth flows the v0 sasl auth NOT wrapped in the kafka protocol
func (b *Broker) sendAndReceiveV0SASLPlainAuth() error {
	length := len(b.conf.Net.SASL.AuthIdentity) + 1 + len(b.conf.Net.SASL.User) + 1 + len(b.conf.Net.SASL.Password)
	authBytes := make([]byte, length+4) //4 byte length header + auth data
	binary.BigEndian.PutUint32(authBytes, uint32(length))
	copy(authBytes[4:], []byte(b.conf.Net.SASL.AuthIdentity+"\x00"+b.conf.Net.SASL.User+"\x00"+b.conf.Net.SASL.Password))

	requestTime := time.Now()
	// Will be decremented in updateIncomingCommunicationMetrics (except error)
	b.addRequestInFlightMetrics(1)
	bytesWritten, err := b.write(authBytes)
	b.updateOutgoingCommunicationMetrics(bytesWritten)
	if err != nil {
		b.addRequestInFlightMetrics(-1)
		Logger.Printf("Failed to write SASL auth header to broker %s: %s\n", b.addr, err.Error())
		return err
	}

	header := make([]byte, 4)
	n, err := b.readFull(header)
	b.updateIncomingCommunicationMetrics(n, time.Since(requestTime))
	// If the credentials are valid, we would get a 4 byte response filled with null characters.
	// Otherwise, the broker closes the connection and we get an EOF
	if err != nil {
		Logger.Printf("Failed to read response while authenticating with SASL to broker %s: %s\n", b.addr, err.Error())
		return err
	}

	Logger.Printf("SASL authentication successful with broker %s:%v - %v\n", b.addr, n, header)
	return nil
}

// sendAndReceiveV1SASLPlainAuth flows the v1 sasl authentication using the kafka protocol
func (b *Broker) sendAndReceiveV1SASLPlainAuth() error {
	correlationID := b.correlationID

	requestTime := time.Now()

	// Will be decremented in updateIncomingCommunicationMetrics (except error)
	b.addRequestInFlightMetrics(1)
	bytesWritten, err := b.sendSASLPlainAuthClientResponse(correlationID)
	b.updateOutgoingCommunicationMetrics(bytesWritten)

	if err != nil {
		b.addRequestInFlightMetrics(-1)
		Logger.Printf("Failed to write SASL auth header to broker %s: %s\n", b.addr, err.Error())
		return err
	}

	b.correlationID++

	bytesRead, err := b.receiveSASLServerResponse(&SaslAuthenticateResponse{}, correlationID)
	b.updateIncomingCommunicationMetrics(bytesRead, time.Since(requestTime))

	// With v1 sasl we get an error message set in the response we can return
	if err != nil {
		Logger.Printf("Error returned from broker during SASL flow %s: %s\n", b.addr, err.Error())
		return err
	}

	return nil
}

// sendAndReceiveSASLOAuth performs the authentication flow as described by KIP-255
// https://cwiki.apache.org/confluence/pages/viewpage.action?pageId=75968876
func (b *Broker) sendAndReceiveSASLOAuth(provider AccessTokenProvider) error {
	if err := b.sendAndReceiveSASLHandshake(SASLTypeOAuth, SASLHandshakeV1); err != nil {
		return err
	}

	token, err := provider.Token()
	if err != nil {
		return err
	}

	message, err := buildClientFirstMessage(token)
	if err != nil {
		return err
	}

	challenged, err := b.sendClientMessage(message)
	if err != nil {
		return err
	}

	if challenged {
		// Abort the token exchange. The broker returns the failure code.
		_, err = b.sendClientMessage([]byte(`\x01`))
	}

	return err
}

// sendClientMessage sends a SASL/OAUTHBEARER client message and returns true
// if the broker responds with a challenge, in which case the token is
// rejected.
func (b *Broker) sendClientMessage(message []byte) (bool, error) {
	requestTime := time.Now()
	// Will be decremented in updateIncomingCommunicationMetrics (except error)
	b.addRequestInFlightMetrics(1)
	correlationID := b.correlationID

	bytesWritten, err := b.sendSASLOAuthBearerClientMessage(message, correlationID)
	b.updateOutgoingCommunicationMetrics(bytesWritten)
	if err != nil {
		b.addRequestInFlightMetrics(-1)
		return false, err
	}

	b.correlationID++

	res := &SaslAuthenticateResponse{}
	bytesRead, err := b.receiveSASLServerResponse(res, correlationID)

	requestLatency := time.Since(requestTime)
	b.updateIncomingCommunicationMetrics(bytesRead, requestLatency)

	isChallenge := len(res.SaslAuthBytes) > 0

	if isChallenge && err != nil {
		Logger.Printf("Broker rejected authentication token: %s", res.SaslAuthBytes)
	}

	return isChallenge, err
}

func (b *Broker) sendAndReceiveSASLSCRAMv1() error {
	if err := b.sendAndReceiveSASLHandshake(b.conf.Net.SASL.Mechanism, SASLHandshakeV1); err != nil {
		return err
	}

	scramClient := b.conf.Net.SASL.SCRAMClientGeneratorFunc()
	if err := scramClient.Begin(b.conf.Net.SASL.User, b.conf.Net.SASL.Password, b.conf.Net.SASL.SCRAMAuthzID); err != nil {
		return fmt.Errorf("failed to start SCRAM exchange with the server: %s", err.Error())
	}

	msg, err := scramClient.Step("")
	if err != nil {
		return fmt.Errorf("failed to advance the SCRAM exchange: %s", err.Error())
	}

	for !scramClient.Done() {
		requestTime := time.Now()
		// Will be decremented in updateIncomingCommunicationMetrics (except error)
		b.addRequestInFlightMetrics(1)
		correlationID := b.correlationID
		bytesWritten, err := b.sendSaslAuthenticateRequest(correlationID, []byte(msg))
		b.updateOutgoingCommunicationMetrics(bytesWritten)
		if err != nil {
			b.addRequestInFlightMetrics(-1)
			Logger.Printf("Failed to write SASL auth header to broker %s: %s\n", b.addr, err.Error())
			return err
		}

		b.correlationID++
		challenge, err := b.receiveSaslAuthenticateResponse(correlationID)
		if err != nil {
			b.addRequestInFlightMetrics(-1)
			Logger.Printf("Failed to read response while authenticating with SASL to broker %s: %s\n", b.addr, err.Error())
			return err
		}

		b.updateIncomingCommunicationMetrics(len(challenge), time.Since(requestTime))
		msg, err = scramClient.Step(string(challenge))
		if err != nil {
			Logger.Println("SASL authentication failed", err)
			return err
		}
	}

	Logger.Println("SASL authentication succeeded")
	return nil
}

func (b *Broker) sendSaslAuthenticateRequest(correlationID int32, msg []byte) (int, error) {
	rb := &SaslAuthenticateRequest{msg}
	req := &request{correlationID: correlationID, clientID: b.conf.ClientID, body: rb}
	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil {
		return 0, err
	}

	return b.write(buf)
}

func (b *Broker) receiveSaslAuthenticateResponse(correlationID int32) ([]byte, error) {
	buf := make([]byte, responseLengthSize+correlationIDSize)
	_, err := b.readFull(buf)
	if err != nil {
		return nil, err
	}

	header := responseHeader{}
	err = versionedDecode(buf, &header, 0)
	if err != nil {
		return nil, err
	}

	if header.correlationID != correlationID {
		return nil, fmt.Errorf("correlation ID didn't match, wanted %d, got %d", b.correlationID, header.correlationID)
	}

	buf = make([]byte, header.length-correlationIDSize)
	_, err = b.readFull(buf)
	if err != nil {
		return nil, err
	}

	res := &SaslAuthenticateResponse{}
	if err := versionedDecode(buf, res, 0); err != nil {
		return nil, err
	}
	if res.Err != ErrNoError {
		return nil, res.Err
	}
	return res.SaslAuthBytes, nil
}

// Build SASL/OAUTHBEARER initial client response as described by RFC-7628
// https://tools.ietf.org/html/rfc7628
func buildClientFirstMessage(token *AccessToken) ([]byte, error) {
	var ext string

	if token.Extensions != nil && len(token.Extensions) > 0 {
		if _, ok := token.Extensions[SASLExtKeyAuth]; ok {
			return []byte{}, fmt.Errorf("the extension `%s` is invalid", SASLExtKeyAuth)
		}
		ext = "\x01" + mapToString(token.Extensions, "=", "\x01")
	}

	resp := []byte(fmt.Sprintf("n,,\x01auth=Bearer %s%s\x01\x01", token.Token, ext))

	return resp, nil
}

// mapToString returns a list of key-value pairs ordered by key.
// keyValSep separates the key from the value. elemSep separates each pair.
func mapToString(extensions map[string]string, keyValSep string, elemSep string) string {
	buf := make([]string, 0, len(extensions))

	for k, v := range extensions {
		buf = append(buf, k+keyValSep+v)
	}

	sort.Strings(buf)

	return strings.Join(buf, elemSep)
}

func (b *Broker) sendSASLPlainAuthClientResponse(correlationID int32) (int, error) {
	authBytes := []byte(b.conf.Net.SASL.AuthIdentity + "\x00" + b.conf.Net.SASL.User + "\x00" + b.conf.Net.SASL.Password)
	rb := &SaslAuthenticateRequest{authBytes}
	req := &request{correlationID: correlationID, clientID: b.conf.ClientID, body: rb}
	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil {
		return 0, err
	}

	return b.write(buf)
}

func (b *Broker) sendSASLOAuthBearerClientMessage(initialResp []byte, correlationID int32) (int, error) {
	rb := &SaslAuthenticateRequest{initialResp}

	req := &request{correlationID: correlationID, clientID: b.conf.ClientID, body: rb}

	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil {
		return 0, err
	}

	return b.write(buf)
}

func (b *Broker) receiveSASLServerResponse(res *SaslAuthenticateResponse, correlationID int32) (int, error) {
	buf := make([]byte, responseLengthSize+correlationIDSize)
	bytesRead, err := b.readFull(buf)
	if err != nil {
		return bytesRead, err
	}

	header := responseHeader{}
	err = versionedDecode(buf, &header, 0)
	if err != nil {
		return bytesRead, err
	}

	if header.correlationID != correlationID {
		return bytesRead, fmt.Errorf("correlation ID didn't match, wanted %d, got %d", b.correlationID, header.correlationID)
	}

	buf = make([]byte, header.length-correlationIDSize)
	c, err := b.readFull(buf)
	bytesRead += c
	if err != nil {
		return bytesRead, err
	}

	if err := versionedDecode(buf, res, 0); err != nil {
		return bytesRead, err
	}

	if res.Err != ErrNoError {
		return bytesRead, res.Err
	}

	return bytesRead, nil
}

func (b *Broker) updateIncomingCommunicationMetrics(bytes int, requestLatency time.Duration) {
	b.updateRequestLatencyAndInFlightMetrics(requestLatency)
	b.responseRate.Mark(1)

	if b.brokerResponseRate != nil {
		b.brokerResponseRate.Mark(1)
	}

	responseSize := int64(bytes)
	b.incomingByteRate.Mark(responseSize)
	if b.brokerIncomingByteRate != nil {
		b.brokerIncomingByteRate.Mark(responseSize)
	}

	b.responseSize.Update(responseSize)
	if b.brokerResponseSize != nil {
		b.brokerResponseSize.Update(responseSize)
	}
}

func (b *Broker) updateRequestLatencyAndInFlightMetrics(requestLatency time.Duration) {
	requestLatencyInMs := int64(requestLatency / time.Millisecond)
	b.requestLatency.Update(requestLatencyInMs)

	if b.brokerRequestLatency != nil {
		b.brokerRequestLatency.Update(requestLatencyInMs)
	}

	b.addRequestInFlightMetrics(-1)
}

func (b *Broker) addRequestInFlightMetrics(i int64) {
	b.requestsInFlight.Inc(i)
	if b.brokerRequestsInFlight != nil {
		b.brokerRequestsInFlight.Inc(i)
	}
}

func (b *Broker) updateOutgoingCommunicationMetrics(bytes int) {
	b.requestRate.Mark(1)
	if b.brokerRequestRate != nil {
		b.brokerRequestRate.Mark(1)
	}

	requestSize := int64(bytes)
	b.outgoingByteRate.Mark(requestSize)
	if b.brokerOutgoingByteRate != nil {
		b.brokerOutgoingByteRate.Mark(requestSize)
	}

	b.requestSize.Update(requestSize)
	if b.brokerRequestSize != nil {
		b.brokerRequestSize.Update(requestSize)
	}
}

func (b *Broker) registerMetrics() {
	b.brokerIncomingByteRate = b.registerMeter("incoming-byte-rate")
	b.brokerRequestRate = b.registerMeter("request-rate")
	b.brokerRequestSize = b.registerHistogram("request-size")
	b.brokerRequestLatency = b.registerHistogram("request-latency-in-ms")
	b.brokerOutgoingByteRate = b.registerMeter("outgoing-byte-rate")
	b.brokerResponseRate = b.registerMeter("response-rate")
	b.brokerResponseSize = b.registerHistogram("response-size")
	b.brokerRequestsInFlight = b.registerCounter("requests-in-flight")
}

func (b *Broker) unregisterMetrics() {
	for _, name := range b.registeredMetrics {
		b.conf.MetricRegistry.Unregister(name)
	}
	b.registeredMetrics = nil
}

func (b *Broker) registerMeter(name string) metrics.Meter {
	nameForBroker := getMetricNameForBroker(name, b)
	b.registeredMetrics = append(b.registeredMetrics, nameForBroker)
	return metrics.GetOrRegisterMeter(nameForBroker, b.conf.MetricRegistry)
}

func (b *Broker) registerHistogram(name string) metrics.Histogram {
	nameForBroker := getMetricNameForBroker(name, b)
	b.registeredMetrics = append(b.registeredMetrics, nameForBroker)
	return getOrRegisterHistogram(nameForBroker, b.conf.MetricRegistry)
}

func (b *Broker) registerCounter(name string) metrics.Counter {
	nameForBroker := getMetricNameForBroker(name, b)
	b.registeredMetrics = append(b.registeredMetrics, nameForBroker)
	return metrics.GetOrRegisterCounter(nameForBroker, b.conf.MetricRegistry)
}

func validServerNameTLS(addr string, cfg *tls.Config) *tls.Config {
	if cfg == nil {
		cfg = &tls.Config{}
	}
	if cfg.ServerName != "" {
		return cfg
	}

	c := cfg.Clone()
	sn, _, err := net.SplitHostPort(addr)
	if err != nil {
		Logger.Println(fmt.Errorf("failed to get ServerName from addr %w", err))
	}
	c.ServerName = sn
	return c
}
