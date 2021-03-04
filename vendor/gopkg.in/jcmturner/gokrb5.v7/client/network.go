package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"gopkg.in/jcmturner/gokrb5.v7/iana/errorcode"
	"gopkg.in/jcmturner/gokrb5.v7/messages"
)

// SendToKDC performs network actions to send data to the KDC.
func (cl *Client) sendToKDC(b []byte, realm string) ([]byte, error) {
	var rb []byte
	if cl.Config.LibDefaults.UDPPreferenceLimit == 1 {
		//1 means we should always use TCP
		rb, errtcp := cl.sendKDCTCP(realm, b)
		if errtcp != nil {
			if e, ok := errtcp.(messages.KRBError); ok {
				return rb, e
			}
			return rb, fmt.Errorf("communication error with KDC via TCP: %v", errtcp)
		}
		return rb, nil
	}
	if len(b) <= cl.Config.LibDefaults.UDPPreferenceLimit {
		//Try UDP first, TCP second
		rb, errudp := cl.sendKDCUDP(realm, b)
		if errudp != nil {
			if e, ok := errudp.(messages.KRBError); ok && e.ErrorCode != errorcode.KRB_ERR_RESPONSE_TOO_BIG {
				// Got a KRBError from KDC
				// If this is not a KRB_ERR_RESPONSE_TOO_BIG we will return immediately otherwise will try TCP.
				return rb, e
			}
			// Try TCP
			r, errtcp := cl.sendKDCTCP(realm, b)
			if errtcp != nil {
				if e, ok := errtcp.(messages.KRBError); ok {
					// Got a KRBError
					return r, e
				}
				return r, fmt.Errorf("failed to communicate with KDC. Attempts made with UDP (%v) and then TCP (%v)", errudp, errtcp)
			}
			rb = r
		}
		return rb, nil
	}
	//Try TCP first, UDP second
	rb, errtcp := cl.sendKDCTCP(realm, b)
	if errtcp != nil {
		if e, ok := errtcp.(messages.KRBError); ok {
			// Got a KRBError from KDC so returning and not trying UDP.
			return rb, e
		}
		rb, errudp := cl.sendKDCUDP(realm, b)
		if errudp != nil {
			if e, ok := errudp.(messages.KRBError); ok {
				// Got a KRBError
				return rb, e
			}
			return rb, fmt.Errorf("failed to communicate with KDC. Attempts made with TCP (%v) and then UDP (%v)", errtcp, errudp)
		}
	}
	return rb, nil
}

// dialKDCTCP establishes a UDP connection to a KDC.
func dialKDCUDP(count int, kdcs map[int]string) (*net.UDPConn, error) {
	i := 1
	for i <= count {
		udpAddr, err := net.ResolveUDPAddr("udp", kdcs[i])
		if err != nil {
			return nil, fmt.Errorf("error resolving KDC address: %v", err)
		}

		conn, err := net.DialTimeout("udp", udpAddr.String(), 5*time.Second)
		if err == nil {
			if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				return nil, err
			}
			// conn is guaranteed to be a UDPConn
			return conn.(*net.UDPConn), nil
		}
		i++
	}
	return nil, errors.New("error in getting a UDP connection to any of the KDCs")
}

// dialKDCTCP establishes a TCP connection to a KDC.
func dialKDCTCP(count int, kdcs map[int]string) (*net.TCPConn, error) {
	i := 1
	for i <= count {
		tcpAddr, err := net.ResolveTCPAddr("tcp", kdcs[i])
		if err != nil {
			return nil, fmt.Errorf("error resolving KDC address: %v", err)
		}

		conn, err := net.DialTimeout("tcp", tcpAddr.String(), 5*time.Second)
		if err == nil {
			if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				return nil, err
			}
			// conn is guaranteed to be a TCPConn
			return conn.(*net.TCPConn), nil
		}
		i++
	}
	return nil, errors.New("error in getting a TCP connection to any of the KDCs")
}

// sendKDCUDP sends bytes to the KDC via UDP.
func (cl *Client) sendKDCUDP(realm string, b []byte) ([]byte, error) {
	var r []byte
	count, kdcs, err := cl.Config.GetKDCs(realm, false)
	if err != nil {
		return r, err
	}
	conn, err := dialKDCUDP(count, kdcs)
	if err != nil {
		return r, err
	}
	r, err = cl.sendUDP(conn, b)
	if err != nil {
		return r, err
	}
	return checkForKRBError(r)
}

// sendKDCTCP sends bytes to the KDC via TCP.
func (cl *Client) sendKDCTCP(realm string, b []byte) ([]byte, error) {
	var r []byte
	count, kdcs, err := cl.Config.GetKDCs(realm, true)
	if err != nil {
		return r, err
	}
	conn, err := dialKDCTCP(count, kdcs)
	if err != nil {
		return r, err
	}
	rb, err := cl.sendTCP(conn, b)
	if err != nil {
		return r, err
	}
	return checkForKRBError(rb)
}

// sendUDP sends bytes to connection over UDP.
func (cl *Client) sendUDP(conn *net.UDPConn, b []byte) ([]byte, error) {
	var r []byte
	defer conn.Close()
	_, err := conn.Write(b)
	if err != nil {
		return r, fmt.Errorf("error sending to (%s): %v", conn.RemoteAddr().String(), err)
	}
	udpbuf := make([]byte, 4096)
	n, _, err := conn.ReadFrom(udpbuf)
	r = udpbuf[:n]
	if err != nil {
		return r, fmt.Errorf("sending over UDP failed to %s: %v", conn.RemoteAddr().String(), err)
	}
	if len(r) < 1 {
		return r, fmt.Errorf("no response data from %s", conn.RemoteAddr().String())
	}
	return r, nil
}

// sendTCP sends bytes to connection over TCP.
func (cl *Client) sendTCP(conn *net.TCPConn, b []byte) ([]byte, error) {
	defer conn.Close()
	var r []byte
	/*
		RFC https://tools.ietf.org/html/rfc4120#section-7.2.2
		Each request (KRB_KDC_REQ) and response (KRB_KDC_REP or KRB_ERROR)
		sent over the TCP stream is preceded by the length of the request as
		4 octets in network byte order.  The high bit of the length is
		reserved for future expansion and MUST currently be set to zero.  If
		a KDC that does not understand how to interpret a set high bit of the
		length encoding receives a request with the high order bit of the
		length set, it MUST return a KRB-ERROR message with the error
		KRB_ERR_FIELD_TOOLONG and MUST close the TCP stream.
		NB: network byte order == big endian
	*/
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, uint32(len(b)))
	if err != nil {
		return r, err
	}
	b = append(buf.Bytes(), b...)

	_, err = conn.Write(b)
	if err != nil {
		return r, fmt.Errorf("error sending to KDC (%s): %v", conn.RemoteAddr().String(), err)
	}

	sh := make([]byte, 4, 4)
	_, err = conn.Read(sh)
	if err != nil {
		return r, fmt.Errorf("error reading response size header: %v", err)
	}
	s := binary.BigEndian.Uint32(sh)

	rb := make([]byte, s, s)
	_, err = io.ReadFull(conn, rb)
	if err != nil {
		return r, fmt.Errorf("error reading response: %v", err)
	}
	if len(rb) < 1 {
		return r, fmt.Errorf("no response data from KDC %s", conn.RemoteAddr().String())
	}
	return rb, nil
}

// checkForKRBError checks if the response bytes from the KDC are a KRBError.
func checkForKRBError(b []byte) ([]byte, error) {
	var KRBErr messages.KRBError
	if err := KRBErr.Unmarshal(b); err == nil {
		return b, KRBErr
	}
	return b, nil
}
