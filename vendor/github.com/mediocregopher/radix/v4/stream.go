package radix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v4/internal/bytesutil"
	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

// StreamEntryID represents an ID used in a Redis stream with the format <time>-<seq>.
type StreamEntryID struct {
	// Time is the first part of the ID, which is based on the time of the server that Redis runs on.
	Time uint64

	// Seq is the sequence number of the ID for entries with the same Time value.
	Seq uint64
}

// Before returns true if s comes before o in a stream (is less than o).
func (s StreamEntryID) Before(o StreamEntryID) bool {
	if s.Time != o.Time {
		return s.Time < o.Time
	}

	return s.Seq < o.Seq
}

// Prev returns the previous stream entry ID or s if there is no prior id (s is 0-0).
func (s StreamEntryID) Prev() StreamEntryID {
	if s.Seq > 0 {
		s.Seq--
		return s
	}

	if s.Time > 0 {
		s.Time--
		s.Seq = math.MaxUint64
		return s
	}

	return s
}

// Next returns the next stream entry ID or s if there is no higher id (s is 18446744073709551615-18446744073709551615).
func (s StreamEntryID) Next() StreamEntryID {
	if s.Seq < math.MaxUint64 {
		s.Seq++
		return s
	}

	if s.Time < math.MaxUint64 {
		s.Time++
		s.Seq = 0
		return s
	}

	return s
}

var _ resp.Marshaler = (*StreamEntryID)(nil)
var _ resp.Unmarshaler = (*StreamEntryID)(nil)

var maxUint64Len = len(strconv.FormatUint(math.MaxUint64, 10))

func (s *StreamEntryID) bytes() []byte {
	b := make([]byte, 0, maxUint64Len*2+1)
	b = strconv.AppendUint(b, s.Time, 10)
	b = append(b, '-')
	b = strconv.AppendUint(b, s.Seq, 10)
	return b
}

// MarshalRESP implements the resp.Marshaler interface.
func (s *StreamEntryID) MarshalRESP(w io.Writer, o *resp.Opts) error {
	return resp3.BlobStringBytes{B: s.bytes()}.MarshalRESP(w, o)
}

var errInvalidStreamID = errors.New("invalid stream entry id")

// UnmarshalRESP implements the resp.Unmarshaler interface.
func (s *StreamEntryID) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	buf := o.GetBytes()
	defer o.PutBytes(buf)

	bsb := resp3.BlobStringBytes{B: (*buf)[:0]}
	if err := bsb.UnmarshalRESP(br, o); err != nil {
		return err
	}

	split := bytes.IndexByte(bsb.B, '-')
	if split == -1 {
		return errInvalidStreamID
	}

	time, err := bytesutil.ParseUint(bsb.B[:split])
	if err != nil {
		return errInvalidStreamID
	}

	seq, err := bytesutil.ParseUint(bsb.B[split+1:])
	if err != nil {
		return errInvalidStreamID
	}

	s.Time, s.Seq = time, seq
	return nil
}

var _ fmt.Stringer = (*StreamEntryID)(nil)

// String returns the ID in the format <time>-<seq> (the same format used by
// Redis).
//
// String implements the fmt.Stringer interface.
func (s StreamEntryID) String() string {
	return string(s.bytes())
}

// StreamEntry is an entry in a stream as returned by XRANGE, XREAD and
// XREADGROUP.
type StreamEntry struct {
	// ID is the ID of the entry in a stream.
	ID StreamEntryID

	// Fields contains the fields and values for the stream entry.
	Fields [][2]string
}

var _ resp.Unmarshaler = (*StreamEntry)(nil)

var errInvalidStreamEntry = errors.New("invalid stream entry")

// UnmarshalRESP implements the resp.Unmarshaler interface.
func (s *StreamEntry) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	var ah resp3.ArrayHeader
	if err := ah.UnmarshalRESP(br, o); err != nil {
		return err
	} else if ah.NumElems != 2 {
		return errInvalidStreamEntry
	} else if err := s.ID.UnmarshalRESP(br, o); err != nil {
		return err
	}

	if err := ah.UnmarshalRESP(br, o); err != nil {
		return err
	} else if ah.NumElems == 0 {
		// if NumElems is zero that means the Fields are actually nil, since
		// it's not possible to submit a stream entry with zero fields.
		s.Fields = s.Fields[:0]
		return nil
	} else if ah.NumElems%2 != 0 {
		return errInvalidStreamEntry
	} else if s.Fields == nil {
		s.Fields = make([][2]string, 0, ah.NumElems)
	}

	var bs resp3.BlobString
	for i := 0; i < ah.NumElems; i += 2 {
		if err := bs.UnmarshalRESP(br, o); err != nil {
			return err
		}
		key := bs.S
		if err := bs.UnmarshalRESP(br, o); err != nil {
			return err
		}
		s.Fields = append(s.Fields, [2]string{key, bs.S})
	}
	return nil
}

// StreamEntries is a stream name and set of entries as returned by XREAD and
// XREADGROUP. The results from a call to XREAD(GROUP) can be unmarshaled into a
// []StreamEntries.
type StreamEntries struct {
	Stream  string
	Entries []StreamEntry
}

// UnmarshalRESP implements the resp.Unmarshaler interface.
func (s *StreamEntries) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	// For RESP2 we get an array of 2 elements, for RESP 3 we are already inside an map so there is no array header.
	if ok, _ := resp3.NextMessageIs(br, resp3.ArrayHeaderPrefix); ok {
		var ah resp3.ArrayHeader
		if err := ah.UnmarshalRESP(br, o); err != nil {
			return err
		} else if ah.NumElems != 2 {
			return errors.New("invalid xread[group] response")
		}
	}

	var stream resp3.BlobString
	if err := stream.UnmarshalRESP(br, o); err != nil {
		return err
	}
	s.Stream = stream.S

	var ah resp3.ArrayHeader
	if err := ah.UnmarshalRESP(br, o); err != nil {
		return err
	}

	s.Entries = make([]StreamEntry, ah.NumElems)
	for i := range s.Entries {
		if err := s.Entries[i].UnmarshalRESP(br, o); err != nil {
			return err
		}
	}
	return nil
}

// streamEntriesMap implements parsing of StreamEntries from XREAD[GROUP] for
// both RESP2 and RESP3 which use different ways to represent the stream names.
type streamEntriesMap []StreamEntries

func (s *streamEntriesMap) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	if err := resp3.DiscardAttribute(br, o); err != nil {
		return err
	}

	if ok, _ := resp3.NextMessageIs(br, resp3.MapHeaderPrefix); ok {
		return s.unmarshalRESP3(br, o)
	}

	return s.unmarshalRESP2(br, o)
}

func (s *streamEntriesMap) unmarshalRESP2(br resp.BufferedReader, o *resp.Opts) error {
	return resp3.Unmarshal(br, (*[]StreamEntries)(s), o)
}

func (s *streamEntriesMap) unmarshalRESP3(br resp.BufferedReader, o *resp.Opts) error {
	var mh resp3.MapHeader
	if err := mh.UnmarshalRESP(br, o); err != nil {
		return err
	}

	// NOTE: This does not handle streamed map responses, but current Redis
	// versions don't use streamed maps for XREAD[GROUP] responses so unless
	// this changes, we panic for now.
	if mh.StreamedMapHeader {
		panic("streamed map response from XREAD[GROUP] not supported")
	}

	ss := *s
	if cap(ss) >= mh.NumPairs {
		ss = ss[:mh.NumPairs]
	} else {
		ss = make([]StreamEntries, mh.NumPairs)
	}

	for i := range ss {
		if err := ss[i].UnmarshalRESP(br, o); err != nil {
			return err
		}
	}

	*s = ss
	return nil
}

// ErrNoStreamEntries is returned by StreamReader's Next method to indicate that
// there were no stream entries left to be read.
var ErrNoStreamEntries = errors.New("no stream entries")

// StreamReader allows reading StreamEntrys sequentially from one or more
// streams.
type StreamReader interface {
	// Next returns a new entry for any of the configured streams. If no new
	// entries are available then Next uses the context's deadline to determine
	// how long to block for (via the BLOCK argument to XREAD(GROUP)). If the
	// context has no deadline then Next will block indefinitely.
	//
	// Next returns ErrNoStreamEntries if there were no entries to be returned.
	// In general Next should be called again after receiving this error.
	//
	// The StreamReader should not be used again if an error which is not
	// ErrNoStreamEntries is returned.
	Next(context.Context) (stream string, entry StreamEntry, err error)
}

// StreamReaderConfig is used to create StreamReader instances with particular
// settings. All fields are optional, all methods are thread-safe.
type StreamReaderConfig struct {
	// Group is an optional consumer group name.
	//
	// If Group is not empty reads will use XREADGROUP with the Group as the
	// group name and Consumer as the consumer name. XREAD will be used
	// otherwise.
	Group string

	// Consumer is an optional consumer name for use with Group.
	Consumer string

	// NoAck enables passing the NOACK flag to XREADGROUP.
	NoAck bool

	// NoBlock disables blocking when no new data is available.
	NoBlock bool

	// Count can be used to limit the number of entries retrieved by each
	// internal redis call to XREAD(GROUP). Can be set to -1 to indicate no
	// limit.
	//
	// Defaults to 20.
	Count int
}

func (cfg StreamReaderConfig) withDefaults() StreamReaderConfig {
	if cfg.Count == -1 {
		cfg.Count = 0
	} else if cfg.Count == 0 {
		cfg.Count = 10
	}
	return cfg
}

// StreamConfig is used to configure the reading behavior of individual streams
// being read by a StreamReader. Exactly one field should be filled in.
type StreamConfig struct {

	// After indicates that only entries newer than the given ID will be
	// returned. If Group is set on the outer StreamReaderConfig then only
	// pending entries newer than the given ID will be returned.
	//
	// The zero StreamEntryID value is a valid value here.
	After StreamEntryID

	// Latest indicates that only entries added after the first call to Next
	// should be returned. If Group is set on the outer StreamReaderConfig then
	// only entries which haven't been delivered to other consumers will be
	// returned.
	Latest bool

	// PendingThenLatest can only be used if Group is set on the outer
	// StreamReaderConfig. The reader will first return entries which are marked
	// as pending for the consumer. Once all pending entries are consumed then
	// the reader will switch to returning entries which haven't been delivered
	// to other consumers.
	PendingThenLatest bool
}

// streamReader implements the StreamReader interface.
type streamReader struct {
	c   Client
	cfg StreamReaderConfig

	streams    []string
	streamCfgs map[string]StreamConfig
	ids        map[string]string

	cmd       string   // command. either XREAD or XREADGROUP
	fixedArgs []string // fixed arguments that always come directly after the command

	unread streamEntriesMap
	err    error
}

// New returns a new StreamReader for the given Client. The StreamReader will
// read from the streams given as the keys of the map.
func (cfg StreamReaderConfig) New(c Client, streamCfgs map[string]StreamConfig) StreamReader {
	sr := &streamReader{
		c:          c,
		cfg:        cfg.withDefaults(),
		streamCfgs: streamCfgs,

		// pre-allocated up to the maximumim potential arguments.
		// (GROUP + group + consumer) + (BLOCK block) + (COUNT count) + NOACK +
		// (STREAMS + streams... + ids...)
		fixedArgs: make([]string, 0, 3+2+2+1+1+len(streamCfgs)*2),
	}

	if sr.cfg.Group != "" {
		sr.cmd = "XREADGROUP"
		sr.fixedArgs = append(sr.fixedArgs, "GROUP", sr.cfg.Group, sr.cfg.Consumer)
	} else {
		sr.cmd = "XREAD"
	}

	if sr.cfg.Count > 0 {
		sr.fixedArgs = append(sr.fixedArgs, "COUNT", strconv.Itoa(sr.cfg.Count))
	}
	if sr.cfg.Group != "" && sr.cfg.NoAck {
		sr.fixedArgs = append(sr.fixedArgs, "NOACK")
	}

	sr.streams = make([]string, 0, len(streamCfgs))
	sr.ids = make(map[string]string, len(streamCfgs))
	for stream, streamCfg := range streamCfgs {
		sr.streams = append(sr.streams, stream)
		if streamCfg.Latest {
			if sr.cfg.Group == "" {
				sr.ids[stream] = "$"
			} else {
				sr.ids[stream] = ">"
			}
		} else {
			sr.ids[stream] = streamCfg.After.String()
		}
	}

	return sr
}

func (sr *streamReader) backfill(ctx context.Context) error {
	args := sr.fixedArgs

	if !sr.cfg.NoBlock {
		now := time.Now()
		if deadline, ok := ctx.Deadline(); ok {
			if d := deadline.Sub(now); d > 200*time.Millisecond {
				// to give us some wiggle room we only block for half the context
				// timeout.
				d /= 2
				args = append(args, "BLOCK", strconv.Itoa(int(d/time.Millisecond)))
			}
		}
	}

	args = append(args, "STREAMS")
	args = append(args, sr.streams...)
	for _, s := range sr.streams {
		args = append(args, sr.ids[s])
	}

	if err := sr.c.Do(ctx, Cmd(&sr.unread, sr.cmd, args...)); err != nil {
		return fmt.Errorf("calling %s: %w", sr.cmd, err)
	}

	// run through returned entries and update ids for the next call, as needed
	for _, sre := range sr.unread {
		if len(sre.Entries) == 0 {
			streamCfg := sr.streamCfgs[sre.Stream]
			if sr.cfg.Group != "" && streamCfg.PendingThenLatest {
				sr.ids[sre.Stream] = ">"
			}
		} else if sr.cfg.Group == "" || sr.ids[sre.Stream] != ">" {
			sr.ids[sre.Stream] = sre.Entries[len(sre.Entries)-1].ID.String()
		}
	}

	return nil
}

func (sr *streamReader) Next(ctx context.Context) (stream string, entry StreamEntry, err error) {
	if sr.err != nil {
		return "", StreamEntry{}, sr.err
	}

	var backfillCalled bool // we only call backfill once per Next
	for {
		if len(sr.unread) == 0 {
			if backfillCalled {
				break
			} else if sr.err = sr.backfill(ctx); sr.err != nil {
				return "", StreamEntry{}, sr.err
			}
			backfillCalled = true
		}

		for len(sr.unread) > 0 {
			i := len(sr.unread) - 1
			if len(sr.unread[i].Entries) == 0 {
				sr.unread = sr.unread[:i]
				continue
			}

			entry := sr.unread[i].Entries[0]
			sr.unread[i].Entries = sr.unread[i].Entries[1:]
			return sr.unread[i].Stream, entry, nil
		}
	}

	return "", StreamEntry{}, ErrNoStreamEntries
}
