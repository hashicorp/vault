package gore

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

// Sentinel is a special Redis process that monitors other Redis instances,
// does fail-over, notifies client status of all monitored instances.
type Sentinel struct {
	servers   []string
	conn      *Conn
	subConn   *Conn // A dedicated connection for pubsub
	subs      *Subscriptions
	mutex     *sync.Mutex
	state     int
	instances map[string]*instance
}

// NewSentinel returns new Sentinel
func NewSentinel() *Sentinel {
	return &Sentinel{
		mutex:     &sync.Mutex{},
		state:     connStateNotConnected,
		instances: make(map[string]*instance),
	}
}

// AddServer adds new sentinel servers. Only one sentinel server is active
// at any time. If this server fails, gore will connect to other sentinel
// servers immediately.
//
// AddServer can be called at anytime, to add new server on the fly.
// In production environment, you should always have at least 3 sentinel
// servers up and running.
func (s *Sentinel) AddServer(addresses ...string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.servers = append(s.servers, addresses...)
}

// Dial connects to one sentinel server in the list. If it fails to connect,
// it moves to the next on the list. If all servers cannot be connected,
// Init return error.
func (s *Sentinel) Dial() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.state != connStateNotConnected {
		return nil
	}
	return s.connect()
}

// Close gracefully closes the sentinel and all monitored connections
func (s *Sentinel) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.close()
}

// GetPool returns a pool of connection from a pool name.
// If the pool has not been retrieved before, gore will attempt to
// fetch the address from the sentinel server, and initialize connections
// with this address. The application should never call this function repeatedly
// to get the same pool, because internal locking can cause performance to drop.
// An error can be returned if the pool name is not monitored by the sentinel,
// or the redis server is currently dead, or the redis server cannot be connected
// (for example: firewall issues).
func (s *Sentinel) GetPool(name string) (*Pool, error) {
	return s.getPool(name, "")
}

// GetPoolWithPassword returns a pool of connection to a password-protected instance
func (s *Sentinel) GetPoolWithPassword(name string, password string) (*Pool, error) {
	return s.getPool(name, password)
}

func (s *Sentinel) getPool(name string, password string) (*Pool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ins, ok := s.instances[name]; ok {
		return ins.pool, nil
	}
	rep, err := NewCommand("SENTINEL", "master", name).Run(s.conn)
	if err != nil {
		return nil, err
	}
	if rep.IsError() {
		return nil, ErrNil
	}
	master, err := rep.Map()
	if err != nil {
		return nil, err
	}
	flags := strings.Split(master["flags"], ",")
	for _, flag := range flags {
		if flag == "s_down" || flag == "o_down" {
			return nil, ErrNotConnected
		}
	}
	ins := &instance{
		name:    name,
		address: master["ip"] + ":" + master["port"],
		state:   connStateConnected,
		pool:    &Pool{sentinel: true, Password: password},
	}
	err = ins.pool.Dial(ins.address)
	if err != nil {
		return nil, err
	}
	s.instances[name] = ins
	return ins.pool, nil
}

// GetCluster returns a cluster monitored by the sentinel.
// The name of the cluster will determine name of Redis instances.
// For example, if the cluster name is "mycluster", the instances' name
// maybe "mycluster1", "mycluster2", ...
func (s *Sentinel) GetCluster(name string) (c *Cluster, err error) {
	return s.getCluster(name, "")
}

// GetClusterWithPassword returns a password-protected cluster monitored by the sentinel.
func (s *Sentinel) GetClusterWithPassword(name string, password string) (c *Cluster, err error) {
	return s.getCluster(name, password)
}

func (s *Sentinel) getCluster(name string, password string) (c *Cluster, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rep, err := NewCommand("SENTINEL", "masters").Run(s.conn)
	if err != nil {
		return nil, err
	}
	replies, err := rep.Array()
	if err != nil {
		return nil, err
	}
	if len(replies) == 0 {
		return nil, ErrNoShard
	}
	instances := make(map[string]*instance)
	defer func() {
		if err != nil {
			for _, ins := range instances {
				ins.pool.Close()
			}
		}
	}()
	for _, r := range replies {
		master, err := r.Map()
		if err != nil {
			return nil, err
		}
		suffix := strings.TrimPrefix(master["name"], name)
		if !suffixRegex.MatchString(suffix) {
			continue
		}
		ins := &instance{
			name:    master["name"],
			address: master["ip"] + ":" + master["port"],
			state:   connStateConnected,
			pool:    &Pool{sentinel: true, Password: password},
		}
		err = ins.pool.Dial(ins.address)
		if err != nil {
			return nil, err
		}
		instances[ins.name] = ins
	}
	c = NewCluster()
	c.sentinel = true
	for _, ins := range instances {
		s.instances[ins.name] = ins
		c.addresses = append(c.addresses, &addressWithPassword{ins.address, password})
		c.shards = append(c.shards, ins.pool)
	}
	return c, nil
}

var suffixRegex = regexp.MustCompile("^\\d+$")

func (s *Sentinel) connect() (err error) {
	for i, server := range s.servers {
		s.conn, err = DialTimeout(server, time.Duration(Config.ConnectTimeout)*time.Second)
		if err != nil {
			continue
		}
		s.subConn, err = DialTimeout(server, time.Duration(Config.ConnectTimeout)*time.Second)
		if err != nil {
			s.conn.Close()
			continue
		}
		s.state = connStateConnected
		s.subs = NewSubscriptions(s.subConn)
		s.subs.throwError = true
		err = s.subs.Subscribe("+sdown", "-sdown", "+odown", "-odown", "+switch-master")
		if err != nil {
			s.close()
			continue
		}
		s.servers = append(s.servers[0:i], s.servers[i+1:]...)
		s.servers = append(s.servers, server)
		go s.monitor()
		return nil
	}
	return ErrNotConnected
}

func (s *Sentinel) close() {
	s.state = connStateNotConnected
	s.subs.Close()
	s.subConn.Close()
	s.conn.Close()
	for _, ins := range s.instances {
		ins.pool.Close()
	}
	s.instances = make(map[string]*instance)
}

func (s *Sentinel) fail() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.state == connStateConnected {
		s.state = connStateNotConnected
		s.reconnect()
	}
}

func (s *Sentinel) reconnect() {
	s.subs.Close()
	s.subConn.Close()
	s.conn.Close()
	sleepTime := Config.ReconnectTime
	for {
		err := s.connect()
		if err == nil {
			break
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
		if sleepTime < 30 {
			sleepTime += 2
		}
	}
}

func (s *Sentinel) monitor() {
	for message := range s.subs.Message() {
		if message == nil {
			s.fail()
			return
		}
		s.mutex.Lock()
		ins := s.getInstanceFromMessage(message)
		if ins == nil {
			s.mutex.Unlock()
			continue
		}
		if message.Channel == "+sdown" || message.Channel == "+odown" {
			ins.down(message)
		} else if message.Channel == "-sdown" || message.Channel == "-odown" {
			ins.up(message)
		} else if message.Channel == "+switch-master" {
			ins.switchMaster(s)
		}
		s.mutex.Unlock()
	}
}

func (s *Sentinel) getInstanceFromMessage(message *Message) *instance {
	if message.Channel == "+sdown" || message.Channel == "+odown" ||
		message.Channel == "-sdown" || message.Channel == "-odown" {
		pieces := strings.Split(string(message.Message), " ")
		if len(pieces) < 2 || pieces[0] != "master" {
			return nil
		}
		return s.instances[pieces[1]]
	} else if message.Channel == "+switch-master" {
		pieces := strings.Split(string(message.Message), " ")
		if len(pieces) < 1 {
			return nil
		}
		return s.instances[pieces[0]]
	}
	return nil
}

func (s *Sentinel) getInstanceAddress(name string) string {
	for {
		rep, err := NewCommand("SENTINEL", "get-master-addr-by-name", name).Run(s.conn)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if !rep.IsArray() {
			return ""
		}
		result := []string{}
		rep.Slice(&result)
		if len(result) != 2 {
			return ""
		}
		return result[0] + ":" + result[1]
	}
}

type instance struct {
	name    string
	address string
	sdown   bool
	odown   bool
	pool    *Pool
	state   int
}

func (ins *instance) down(message *Message) {
	if ins.state != connStateConnected {
		return
	}
	ins.state = connStateNotConnected
	if message.Channel == "+sdown" {
		ins.sdown = true
	} else if message.Channel == "+odown" {
		ins.odown = true
	}
	ins.pool.sentinelGonnaLetYouDown()
}

func (ins *instance) up(message *Message) {
	if ins.state == connStateConnected {
		return
	}
	if message.Channel == "-sdown" {
		ins.sdown = false
	} else if message.Channel == "-odown" {
		ins.odown = false
	}
	if !ins.sdown && !ins.odown {
		ins.pool.sentinelGonnaGiveYouUp()
		ins.state = connStateConnected
	}
}

func (ins *instance) switchMaster(s *Sentinel) {
	address := s.getInstanceAddress(ins.name)
	if address == "" {
		// WTF
		return
	}
	if ins.state == connStateConnected {
		ins.pool.sentinelGonnaLetYouDown()
	}
	ins.pool.address = address
	ins.pool.sentinelGonnaGiveYouUp()
	ins.state = connStateConnected
}
