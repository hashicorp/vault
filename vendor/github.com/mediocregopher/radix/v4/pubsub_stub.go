package radix

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mediocregopher/radix/v4/internal/proc"
	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

var errPubSubMode = resp3.SimpleError{
	S: "ERR only (P)SUBSCRIBE / (P)UNSUBSCRIBE / PING / QUIT allowed in this context",
}

type multiMarshal []interface{}

func (mm multiMarshal) MarshalRESP(w io.Writer, o *resp.Opts) error {
	for _, m := range mm {
		if err := resp3.Marshal(w, m, o); err != nil {
			return err
		}
	}
	return nil
}

type pubSubStub struct {
	proc *proc.Proc
	Conn
	fn   func(context.Context, []string) interface{}
	inCh <-chan PubSubMessage

	pubsubMode      bool
	subbed, psubbed map[string]bool

	// this is only used for tests
	mDoneCh chan struct{}
}

// NewPubSubConnStub returns a stubbed Conn, much like NewStubConn does, which
// pretends it is a Conn to a real redis instance, but is instead using the
// given callback to service requests. It is primarily useful for writing tests.
//
// NewPubSubConnStub differs from NewStubConn in that EncodeDecode calls for the
// (P)SUBSCRIBE, (P)UNSUBSCRIBE, and PING commands will be intercepted and
// handled as per redis' expected pubsub functionality. A PubSubMessage may be
// written to the returned channel at any time, and if the returned Conn has had
// (P)SUBSCRIBE called matching that PubSubMessage then the PubSubMessage will
// be written to the Conn's internal buffer.
//
// This is intended to be used for mocking services which can perform both
// normal redis commands and pubsub (e.g. a real redis instance, redis
// sentinel). The returned Conn can be passed into NewPubSubConn.
//
// remoteNetwork and remoteAddr can be empty, but if given will be used as the
// return from the Addr method.
//
func NewPubSubConnStub(remoteNetwork, remoteAddr string, fn func(context.Context, []string) interface{}) (Conn, chan<- PubSubMessage) {
	if fn == nil {
		fn = func(_ context.Context, args []string) interface{} {
			return fmt.Errorf("command %#v not supported by stub", args)
		}
	}

	ch := make(chan PubSubMessage)
	s := &pubSubStub{
		proc:    proc.New(),
		fn:      fn,
		inCh:    ch,
		subbed:  map[string]bool{},
		psubbed: map[string]bool{},
		mDoneCh: make(chan struct{}, 1),
	}
	s.Conn = NewStubConn(remoteNetwork, remoteAddr, s.innerFn)
	s.proc.Run(s.spin)
	return s, ch
}

func (s *pubSubStub) innerFn(ctx context.Context, ss []string) interface{} {
	var res interface{}
	err := s.proc.WithLock(func() error {
		writeRes := func(mm multiMarshal, cmd, subj string) multiMarshal {
			c := len(s.subbed) + len(s.psubbed)
			s.pubsubMode = c > 0
			return append(mm, []interface{}{cmd, subj, c})
		}

		switch strings.ToUpper(ss[0]) {
		case "PING":
			if !s.pubsubMode {
				res = resp3.SimpleString{S: "PONG"}
			} else {
				res = []string{"PONG", ""}
			}
		case "SUBSCRIBE":
			var mm multiMarshal
			for _, channel := range ss[1:] {
				s.subbed[channel] = true
				mm = writeRes(mm, "subscribe", channel)
			}
			res = mm
		case "UNSUBSCRIBE":
			var mm multiMarshal
			for _, channel := range ss[1:] {
				delete(s.subbed, channel)
				mm = writeRes(mm, "unsubscribe", channel)
			}
			res = mm
		case "PSUBSCRIBE":
			var mm multiMarshal
			for _, pattern := range ss[1:] {
				s.psubbed[pattern] = true
				mm = writeRes(mm, "psubscribe", pattern)
			}
			res = mm
		case "PUNSUBSCRIBE":
			var mm multiMarshal
			for _, pattern := range ss[1:] {
				delete(s.psubbed, pattern)
				mm = writeRes(mm, "punsubscribe", pattern)
			}
			res = mm
		case "MESSAGE":
			m := PubSubMessage{
				Type:    "message",
				Channel: ss[1],
				Message: []byte(ss[2]),
			}

			var mm multiMarshal
			if s.subbed[m.Channel] {
				mm = append(mm, m)
			}
			res = mm
		case "PMESSAGE":
			m := PubSubMessage{
				Type:    "pmessage",
				Pattern: ss[1],
				Channel: ss[2],
				Message: []byte(ss[3]),
			}

			var mm multiMarshal
			if s.psubbed[m.Pattern] {
				mm = append(mm, m)
			}
			res = mm
		default:
			if s.pubsubMode {
				return errPubSubMode
			}
			res = s.fn(ctx, ss)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return res
}

func (s *pubSubStub) Close() error {
	return s.proc.Close(func() error {
		return s.Conn.Close()
	})
}

func (s *pubSubStub) spin(ctx context.Context) {
	for {
		select {
		case m, ok := <-s.inCh:
			if !ok {
				panic("PubSubStub message channel was closed")
			}
			if m.Type == "" {
				if m.Pattern == "" {
					m.Type = "message"
				} else {
					m.Type = "pmessage"
				}
			}
			if err := s.Conn.EncodeDecode(ctx, m, nil); err != nil {
				panic(fmt.Sprintf("error encoding message in PubSubStub: %s", err))
			}
			select {
			case s.mDoneCh <- struct{}{}:
			default:
			}
		case <-ctx.Done():
			return
		}
	}
}
