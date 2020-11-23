package gocbcore

import (
	"crypto/tls"
	"io"
	"net"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

type memdConn interface {
	LocalAddr() string
	RemoteAddr() string
	WritePacket(*memd.Packet) error
	ReadPacket() (*memd.Packet, int, error)
	Close() error

	EnableFeature(feature memd.HelloFeature)
	IsFeatureEnabled(feature memd.HelloFeature) bool
}

type memdConnWrap struct {
	localAddr  string
	remoteAddr string
	conn       *memd.Conn
	baseConn   io.Closer
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

func dialMemdConn(address string, tlsConfig *tls.Config, deadline time.Time) (memdConn, error) {
	d := net.Dialer{
		Deadline: deadline,
	}

	baseConn, err := d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	tcpConn, isTCPConn := baseConn.(*net.TCPConn)
	if !isTCPConn || tcpConn == nil {
		return nil, errCliInternalError
	}

	err = tcpConn.SetNoDelay(false)
	if err != nil {
		logWarnf("Failed to disable TCP nodelay (%s)", err)
	}

	var conn io.ReadWriteCloser
	if tlsConfig == nil {
		conn = tcpConn
	} else {
		tlsConn := tls.Client(tcpConn, tlsConfig)
		err = tlsConn.Handshake()
		if err != nil {
			return nil, err
		}

		conn = tlsConn
	}

	return &memdConnWrap{
		conn:       memd.NewConn(conn),
		baseConn:   conn,
		localAddr:  baseConn.LocalAddr().String(),
		remoteAddr: address,
	}, nil
}
