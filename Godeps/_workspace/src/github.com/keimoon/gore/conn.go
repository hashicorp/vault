package gore

import (
	"bufio"
	"net"
	"sync"
	"time"
)

const (
	connStateNotConnected = iota
	connStateConnected
	connStateReconnecting
)

// Conn holds a persistent connection to a redis server
type Conn struct {
	address        string
	tcpConn        net.Conn
	state          int
	mutex          sync.Mutex
	rb             *bufio.Reader
	wb             *bufio.Writer
	sentinel       bool
	RequestTimeout time.Duration
	isClosed       bool
	password       string
}

// Dial opens a TCP connection with a redis server.
func Dial(address string) (*Conn, error) {
	conn := &Conn{
		RequestTimeout: time.Duration(Config.RequestTimeout) * time.Second,
	}
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.connect(address, 0)
	return conn, err
}

// DialTimeout opens a TCP connection with a redis server with a connection timeout
func DialTimeout(address string, timeout time.Duration) (*Conn, error) {
	conn := &Conn{
		RequestTimeout: time.Duration(Config.RequestTimeout) * time.Second,
	}
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.connect(address, timeout)
	return conn, err
}

// Auth makes authentication with redis server
func (c *Conn) Auth(password string) error {
	c.password = password
	if c.password == "" {
		return nil
	}
	rep, err := NewCommand("AUTH", password).Run(c)
	if err != nil {
		return err
	}
	if !rep.IsOk() {
		return ErrAuth
	}
	return nil
}

// Close closes the connection
func (c *Conn) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isClosed = true
	if c.state == connStateNotConnected || c.tcpConn == nil {
		return nil
	}
	c.state = connStateNotConnected
	return c.tcpConn.Close()
}

// IsConnected returns true if connection is okay
func (c *Conn) IsConnected() bool {
	return c.state == connStateConnected
}

// GetAddress returns connection address
func (c *Conn) GetAddress() string {
	return c.address
}

// Lock locks the whole connection
func (c *Conn) Lock() {
	c.mutex.Lock()
}

// Unlock unlocks the whole connection
func (c *Conn) Unlock() {
	c.mutex.Unlock()
}

func (c *Conn) connect(address string, timeout time.Duration) error {
	if c.state == connStateConnected {
		return nil
	}
	var err error
	c.address = address
	if timeout == 0 {
		c.tcpConn, err = net.Dial("tcp", address)
	} else {
		c.tcpConn, err = net.DialTimeout("tcp", address, timeout)
	}
	if err == nil {
		c.state = connStateConnected
		c.rb = bufio.NewReader(c.tcpConn)
		c.wb = bufio.NewWriter(c.tcpConn)
	}
	return err
}

func (c *Conn) fail() {
	if !c.sentinel {
		c.mutex.Lock()
		if c.state == connStateReconnecting {
			c.mutex.Unlock()
			return
		}
		c.tcpConn.Close()
		c.state = connStateReconnecting
		c.mutex.Unlock()
		go c.reconnect()
	}
}

func (c *Conn) reconnect() {
	sleepTime := Config.ReconnectTime
	for {
		c.mutex.Lock()
		if c.isClosed {
			c.mutex.Unlock()
			break
		}
		if err := c.connect(c.address, 0); err == nil {
			c.mutex.Unlock()
			break
		}
		c.mutex.Unlock()
		time.Sleep(time.Duration(sleepTime) * time.Second)
		if sleepTime < 30 {
			sleepTime += 2
		}
	}
	if c.password != "" {
		c.Auth(c.password)
	}
}
