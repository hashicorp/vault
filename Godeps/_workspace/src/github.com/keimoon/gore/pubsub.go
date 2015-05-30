package gore

import (
	"sync"
	"time"
)

// Message is a nofitication from a subscribed channel or pchannel
type Message struct {
	// message or pmessage
	Type string
	// The channel/pchannel the client subscribed to. For example: "test", "te*"
	Channel string
	// The channel that publisher published to. For example "test", "text"
	OriginalChannel string
	// The payload
	Message []byte
}

// Publish a message to a channel over conn
func Publish(conn *Conn, channel string, message interface{}) error {
	_, err := NewCommand("PUBLISH", channel, message).Run(conn)
	return err
}

// Subscriptions keeps all SUBSCRIBE and PSUBSCRIBE channels, and handles
// all errors and re-subcription process when the connection is down then reconnected.
// A dedicated connection should be used for Subscriptions
type Subscriptions struct {
	channels       map[string]struct{}
	pchannels      map[string]struct{}
	conn           *Conn
	closed         bool
	messageChannel chan *Message
	lock           sync.Mutex
	ready          bool
	readyChannel   chan bool
	// Sentinel will set this to true to handle read error from sentinel server.
	throwError bool
}

// NewSubscriptions returns new Subscriptions
func NewSubscriptions(conn *Conn) *Subscriptions {
	s := &Subscriptions{
		channels:       make(map[string]struct{}),
		pchannels:      make(map[string]struct{}),
		conn:           conn,
		messageChannel: make(chan *Message, 100),
		ready:          true,
		readyChannel:   make(chan bool, 1),
		throwError:     false,
	}
	go s.receive()
	return s
}

// Subscribe subscribes to a list of channels
func (s *Subscriptions) Subscribe(channel ...string) error {
	return s.do("SUBSCRIBE", channel...)
}

// PSubscribe subscribes to a list of channels with given pattern
func (s *Subscriptions) PSubscribe(channel ...string) error {
	return s.do("PSUBSCRIBE", channel...)
}

// Unsubscribe unsubscribes to a list of channels
func (s *Subscriptions) Unsubscribe(channel ...string) error {
	return s.do("UNSUBSCRIBE", channel...)
}

// PUnsubscribe unsubscribes to a list of channels with given pattern
func (s *Subscriptions) PUnsubscribe(channel ...string) error {
	return s.do("PUNSUBSCRIBE", channel...)
}

// Close terminates the subscriptions.
// The connection is NOT closed. You should close it if you do not want to
// use anymore.
func (s *Subscriptions) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	close(s.readyChannel)
	close(s.messageChannel)
}

// Message returns a channel for receiving message event.
// A nil message indicates the channel is closed.
// The channel should be used from a separated goroutine.
// For example:
//
//    for message := range subs.Message() {
//       if message == nil {
//           break
//       }
//       ...
//    }
func (s *Subscriptions) Message() chan *Message {
	return s.messageChannel
}

func (s *Subscriptions) receive() {
	for {
		if s.closed {
			break
		}
		for !s.ready {
			<-s.readyChannel
			if s.closed {
				break
			}
		}
		for {
			rep, err := readReply(s.conn)
			if err != nil {
				if s.throwError {
					s.messageChannel <- nil
					return
				}
				s.ready = false
				go s.resubscribe()
				break
			}
			if !rep.IsArray() {
				continue
			}
			replies, _ := rep.Array()
			if len(replies) < 3 || !replies[0].IsString() || !replies[1].IsString() || !replies[2].IsString() {
				continue
			}
			switch string(replies[0].stringValue) {
			case "message":
				s.messageChannel <- &Message{
					Type:            "message",
					Channel:         string(replies[1].stringValue),
					OriginalChannel: string(replies[1].stringValue),
					Message:         replies[2].stringValue,
				}
			case "pmessage":
				s.messageChannel <- &Message{
					Type:            "pmessage",
					Channel:         string(replies[1].stringValue),
					OriginalChannel: string(replies[2].stringValue),
					Message:         replies[3].stringValue,
				}
			}
		}

	}
}

func (s *Subscriptions) do(command string, channel ...string) error {
	if len(channel) == 0 {
		return nil
	}
	if s.conn.state != connStateConnected {
		return ErrNotConnected
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	err := NewCommand(command, s.makeArgs(channel...)...).Send(s.conn)
	if err == nil {
		switch {
		case command[0] == 'S':
			for _, ch := range channel {
				s.channels[ch] = struct{}{}
			}
		case command[0] == 'U':
			for _, ch := range channel {
				delete(s.channels, ch)
			}
		case command[1] == 'S':
			for _, ch := range channel {
				s.pchannels[ch] = struct{}{}
			}
		default:
			for _, ch := range channel {
				delete(s.pchannels, ch)
			}
		}
	}
	return err
}

func (s *Subscriptions) makeArgs(channel ...string) []interface{} {
	args := make([]interface{}, len(channel))
	for i := range channel {
		args[i] = channel[i]
	}
	return args
}

func (s *Subscriptions) resubscribe() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for {
		if s.conn.state != connStateConnected {
			time.Sleep(2 * time.Second)
			continue
		}
		channels := []string{}
		for ch := range s.channels {
			channels = append(channels, ch)
		}
		pchannels := []string{}
		for ch := range s.pchannels {
			pchannels = append(pchannels, ch)
		}
		err := NewCommand("SUBSCRIBE", s.makeArgs(channels...)...).Send(s.conn)
		if err != nil {
			continue
		}
		err = NewCommand("PSUBSCRIBE", s.makeArgs(pchannels...)...).Send(s.conn)
		if err != nil {
			continue
		}
		break
	}
	s.ready = true
	s.readyChannel <- true
}
