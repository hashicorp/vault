package gocbcore

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

const defaultReaderBufSize = 20 * 1024 * 1024

type memdConn interface {
	LocalAddr() string
	RemoteAddr() string
	WritePacket(*memd.Packet) error
	ReadPacket() (*memd.Packet, int, error)
	Close() error
	Release()

	EnableFeature(feature memd.HelloFeature)
	IsFeatureEnabled(feature memd.HelloFeature) bool
}

type wrappedReadWriteCloser struct {
	*bufio.Reader
	io.Writer
	io.Closer
}

// readerBufPools - Map of buffer size to thread safe pool containing packet reader buffers.
var readerBufPools = map[int]*sync.Pool{}
var readerBufPoolsLock sync.Mutex

// acquireReadBuf - Returns a pointer to a read buffer which is ready to be used, ensure the buffer is released using
// the 'releaseWriteBuf' function.
func acquireReadBuf(stream io.Reader, bufSize int) *bufio.Reader {
	readerBufPoolsLock.Lock()
	bufPool, ok := readerBufPools[bufSize]
	if !ok {
		bufPool = &sync.Pool{}
		readerBufPools[bufSize] = bufPool
	}
	readerBufPoolsLock.Unlock()

	iReader := bufPool.Get()
	var reader *bufio.Reader
	if iReader == nil {
		reader = bufio.NewReaderSize(stream, bufSize)
	} else {
		var ok bool
		reader, ok = iReader.(*bufio.Reader)
		if ok {
			reader.Reset(stream)
		} else {
			reader = bufio.NewReaderSize(stream, bufSize)
		}
	}

	return reader
}

// releaseReadBuf - Reset the buffer so that it's clean for the next user (note that this retains the underlying
// storage for future reads) and then return it to the pool.
func releaseReadBuf(buf *bufio.Reader, bufSize int) {
	buf.Reset(nil)
	readerBufPoolsLock.Lock()
	bufPool, ok := readerBufPools[bufSize]
	if !ok {
		readerBufPoolsLock.Unlock()
		logWarnf("Attempted to release a read buffer for a buffer size without a registered pool")
		return
	}
	bufPool.Put(buf)
	readerBufPoolsLock.Unlock()
}

type memdConnWrap struct {
	localAddr  string
	remoteAddr string
	conn       *memd.Conn
	baseConn   *wrappedReadWriteCloser
	bufSize    int
}

func (s *memdConnWrap) LocalAddr() string {
	return s.localAddr
}

func (s *memdConnWrap) RemoteAddr() string {
	return s.remoteAddr
}

func (s *memdConnWrap) WritePacket(pkt *memd.Packet) error {
	return s.conn.WritePacket(pkt)
}

func (s *memdConnWrap) ReadPacket() (*memd.Packet, int, error) {
	return s.conn.ReadPacket()
}

func (s *memdConnWrap) EnableFeature(feature memd.HelloFeature) {
	s.conn.EnableFeature(feature)
}

func (s *memdConnWrap) IsFeatureEnabled(feature memd.HelloFeature) bool {
	return s.conn.IsFeatureEnabled(feature)
}

func (s *memdConnWrap) Close() error {
	return s.baseConn.Close()
}

// Release is not thread safe and should not be called whilst there are pending calls, such as ReadPacket.
func (s *memdConnWrap) Release() {
	if s.baseConn == nil {
		logWarnf("Release called on already released connection")
		return
	}
	releaseReadBuf(s.baseConn.Reader, s.bufSize)
	s.baseConn = nil
}

func dialMemdConn(ctx context.Context, address string, tlsConfig *tls.Config, deadline time.Time, bufSize uint) (memdConn, error) {
	d := net.Dialer{
		Deadline: deadline,
	}

	dialID := formatCbUID(randomCbUID())
	logDebugf("Dialling new client connection for %s, dial id = %s", address, dialID)

	baseConn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		logDebugf("Failed to dial client connection for %s, dial id = %s", address, dialID)
		return nil, err
	}

	logDebugf("Dialled new client connection for %s, dial id = %s", address, dialID)

	tcpConn, isTCPConn := baseConn.(*net.TCPConn)
	if !isTCPConn || tcpConn == nil {
		return nil, errCliInternalError
	}

	err = tcpConn.SetNoDelay(false)
	if err != nil {
		logWarnf("Failed to disable TCP nodelay (%s)", err)
	}

	var conn io.ReadWriteCloser = tcpConn
	if tlsConfig != nil {
		tlsConn := tls.Client(tcpConn, tlsConfig)
		err = tlsConn.Handshake()
		if err != nil {
			return nil, err
		}

		conn = tlsConn
	}

	if bufSize == 0 {
		bufSize = defaultReaderBufSize
	}

	c := &wrappedReadWriteCloser{
		Reader: acquireReadBuf(conn, int(bufSize)),
		Writer: conn,
		Closer: conn,
	}

	return &memdConnWrap{
		conn:       memd.NewConn(c),
		baseConn:   c,
		localAddr:  baseConn.LocalAddr().String(),
		remoteAddr: address,
		bufSize:    int(bufSize),
	}, nil
}
