package radix

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

// ActionProperties describes various properties of an Action. It should be
// expected that more fields will be added to this struct as time goes forward,
// though the zero values of those new fields will have sane default behaviors.
type ActionProperties struct {
	// Keys describes which redis keys an Action will act on. An empty/nil slice
	// maybe used if no keys are being acted on. The slice may contain duplicate
	// values.
	Keys []string

	// CanRetry indicates, in the event of a cluster node returning a MOVED or
	// ASK error, the Action can be retried on a different node.
	CanRetry bool

	// CanPipeline indicates that an Action can be pipelined alongside other
	// Actions for which this property is also true.
	CanPipeline bool

	// CanShareConn indicates that an Action can be Perform'd on a Conn
	// concurrently with other Actions for which this property is also true.
	CanShareConn bool
}

// Action performs a task using a Conn.
type Action interface {
	// Properties returns an ActionProperties value for the Action. Multiple
	// calls to Properties should always yield the same ActionProperties value.
	Properties() ActionProperties

	// Perform actually performs the Action using an existing Conn.
	Perform(ctx context.Context, c Conn) error
}

var noKeyCmds = map[string]bool{
	"SENTINEL": true,

	"CLUSTER":   true,
	"READONLY":  true,
	"READWRITE": true,
	"ASKING":    true,

	"AUTH":   true,
	"ECHO":   true,
	"PING":   true,
	"QUIT":   true,
	"SELECT": true,
	"SWAPDB": true,

	"KEYS":      true,
	"MIGRATE":   true,
	"OBJECT":    true,
	"RANDOMKEY": true,
	"WAIT":      true,
	"SCAN":      true,

	"EVAL":    true,
	"EVALSHA": true,
	"SCRIPT":  true,

	"BGREWRITEAOF": true,
	"BGSAVE":       true,
	"CLIENT":       true,
	"COMMAND":      true,
	"CONFIG":       true,
	"DBSIZE":       true,
	"DEBUG":        true,
	"FLUSHALL":     true,
	"FLUSHDB":      true,
	"INFO":         true,
	"LASTSAVE":     true,
	"MONITOR":      true,
	"ROLE":         true,
	"SAVE":         true,
	"SHUTDOWN":     true,
	"SLAVEOF":      true,
	"SLOWLOG":      true,
	"SYNC":         true,
	"TIME":         true,

	"DISCARD": true,
	"EXEC":    true,
	"MULTI":   true,
	"UNWATCH": true,
	"WATCH":   true,
}

var blockingCmds = map[string]bool{
	"WAIT": true,

	// taken from https://github.com/joomcode/redispipe#limitations
	"BLPOP":      true,
	"BRPOP":      true,
	"BRPOPLPUSH": true,

	"BZPOPMIN": true,
	"BZPOPMAX": true,

	"XREAD":      true,
	"XREADGROUP": true,

	"SAVE": true,
}

func marshalBlobString(prevErr error, w io.Writer, str string, o *resp.Opts) error {
	if prevErr != nil {
		return prevErr
	}
	return resp3.BlobString{S: str}.MarshalRESP(w, o)
}

func marshalBlobStringBytes(prevErr error, w io.Writer, b []byte, o *resp.Opts) error {
	if prevErr != nil {
		return prevErr
	}
	return resp3.BlobStringBytes{B: b}.MarshalRESP(w, o)
}

////////////////////////////////////////////////////////////////////////////////

type cmdAction struct {
	rcv        interface{}
	cmd        string
	args       []string
	properties ActionProperties

	flattenErr error
}

// BREAM: Benchmarks Rule Everything Around Me.
var cmdActionPool sync.Pool

func getCmdAction() *cmdAction {
	if ci := cmdActionPool.Get(); ci != nil {
		return ci.(*cmdAction)
	}
	return new(cmdAction)
}

func (c *cmdAction) Properties() ActionProperties {
	return c.properties
}

func (c *cmdAction) MarshalRESP(w io.Writer, o *resp.Opts) error {
	if c.flattenErr != nil {
		return c.flattenErr
	}

	err := resp3.ArrayHeader{NumElems: len(c.args) + 1}.MarshalRESP(w, o)
	err = marshalBlobString(err, w, c.cmd, o)
	for i := range c.args {
		err = marshalBlobString(err, w, c.args[i], o)
	}
	return err
}

func (c *cmdAction) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	if err := resp3.Unmarshal(br, c.rcv, o); err != nil {
		return err
	}
	cmdActionPool.Put(c)
	return nil
}

func (c *cmdAction) Perform(ctx context.Context, conn Conn) error {
	return conn.EncodeDecode(ctx, c, c)
}

func (c *cmdAction) String() string {
	ss := make([]string, 0, len(c.args)+1)
	ss = append(ss, strings.ToUpper(c.cmd))
	ss = append(ss, c.args...)
	for i := range ss {
		ss[i] = strconv.QuoteToASCII(ss[i])
	}
	return "[" + strings.Join(ss, " ") + "]"
}

// CmdConfig is used to create redis command Actions with particular settings. All
// fields are optional, all methods are thread-safe.
//
// The global Cmd and FlatCmd functions are shortcuts for using an empty
// CmdConfig{}.
//
// This can be useful for working with custom commands provided by Redis
// modules for which the built in logic may return sub-optimal properties or in
// case some commands are known to be slow in some cases and therefore should
// not use connection sharing.
//
// All methods on CmdConfig are safe for concurrent use from different
// goroutines.
type CmdConfig struct {
	// ActionProperties is an optional callback that will be called when
	// creating a new Action using CmdConfig.Cmd or CmdConfig.FlatCmd and is
	// used to set the ActionProperties for the new Action.
	//
	// If ActionProperties is nil, DefaultActionProperties will be used.
	ActionProperties func(cmd string, args ...string) ActionProperties
}

// DefaultActionProperties returns an ActionProperties instance for the given
// Redis command and it's argument.
//
// The returned ActionProperties should work well with all standard Redis
// commands, including allowing the command to be used in pipelines and for
// connection sharing if the command is non-blocking, but may not return
// correct results for custom commands provided by Redis Modules or unreleased
// commands.
func DefaultActionProperties(cmd string, args ...string) ActionProperties {
	isBlocking := blockingCmds[strings.ToUpper(cmd)]
	properties := ActionProperties{
		CanRetry:     true,
		CanPipeline:  !isBlocking,
		CanShareConn: !isBlocking,
	}
	cmd = strings.ToUpper(cmd)
	switch {
	case noKeyCmds[cmd] || len(args) == 0:
	case cmd == "BITOP" && len(args) > 1:
		properties.Keys = args[1:]
	case cmd == "MEMORY" && len(args) > 1 && strings.ToUpper(args[0]) == "USAGE":
		properties.Keys = args[1:2]
	case cmd == "MSET":
		properties.Keys = keysFromKeyValuePairs(args)
	case cmd == "XINFO":
		if len(args) >= 2 {
			properties.Keys = args[1:2]
		}
	case cmd == "XGROUP" && len(args) > 1:
		properties.Keys = args[1:2]
	case cmd == "XREAD" || cmd == "XREADGROUP":
		properties.Keys = findStreamsKeys(args)
	default:
		properties.Keys = args[:1]
	}
	return properties
}

func findStreamsKeys(args []string) []string {
	for i, arg := range args {
		if strings.ToUpper(arg) != "STREAMS" {
			continue
		}

		// after STREAMS only stream keys and IDs can be given and since there must be the same number of keys and ids
		// we can just take half of remaining arguments as keys. If the number of IDs does not match the number of
		// keys the command will fail later when send to Redis so no need for us to handle that case.
		ids := len(args[i+1:]) / 2

		return args[i+1 : len(args)-ids]
	}

	return nil
}

func keysFromKeyValuePairs(pairs []string) []string {
	keys := make([]string, len(pairs)/2)
	for i := range keys {
		keys[i] = pairs[i*2]
	}
	return keys
}

func (cfg CmdConfig) actionProperties(cmd string, args ...string) ActionProperties {
	pf := cfg.ActionProperties
	if pf == nil {
		pf = DefaultActionProperties
	}
	return pf(cmd, args...)
}

// Cmd works like the global Cmd function but can be additionally configured using
// fields on CmdConfig. See the global Cmd's documentation for further details.
func (cfg CmdConfig) Cmd(rcv interface{}, cmd string, args ...string) Action {
	c := getCmdAction()
	*c = cmdAction{
		rcv:  rcv,
		cmd:  cmd,
		args: args,
	}
	c.properties = cfg.actionProperties(c.cmd, c.args...)
	return c
}

// FlatCmd works like the global FlatCmd function but can be additionally configured using
// fields on CmdConfig. See the global FlatCmd's documentation for further details.
func (cfg CmdConfig) FlatCmd(rcv interface{}, cmd string, args ...interface{}) Action {
	c := getCmdAction()
	*c = cmdAction{
		rcv: rcv,
		cmd: cmd,
	}
	c.args, c.flattenErr = resp3.Flatten(args, nil)
	c.properties = cfg.actionProperties(c.cmd, c.args...)
	return c
}

// Cmd is used to perform a redis command and retrieve a result. It should not
// be passed into Do more than once.
//
// If the receiver value of Cmd is nil then the result is discarded.
//
// If the receiver value of Cmd is a primitive, a slice/map, or a struct then a
// pointer must be passed in. It may also be an io.Writer, an
// encoding.Text/BinaryUnmarshaler, or a resp.Unmarshaler.
//
// The Action returned by Cmd also implements resp.Marshaler.
//
// See CmdConfig's documentation if more configurability is required, e.g. if
// using commands provided by Redis Modules.
func Cmd(rcv interface{}, cmd string, args ...string) Action {
	return (CmdConfig{}).Cmd(rcv, cmd, args...)
}

// FlatCmd is like Cmd, but the arguments can be of almost any type, and FlatCmd
// will automatically flatten them into a single array of strings. Like Cmd, a
// FlatCmd should not be passed into Do more than once.
//
// FlatCmd supports using a resp.LenReader (an io.Reader with a Len() method) as
// an argument. *bytes.Buffer is an example of a LenReader, and the resp package
// has a NewLenReader function which can wrap an existing io.Reader.
//
// FlatCmd supports encoding.Text/BinaryMarshalers, big.Float, and big.Int.
//
// The receiver to FlatCmd follows the same rules as for Cmd.
//
// The Action returned by FlatCmd implements resp.Marshaler.
//
// See CmdConfig's documentation if more configurability is required, e.g. if
// using commands provided by Redis Modules.
func FlatCmd(rcv interface{}, cmd string, args ...interface{}) Action {
	return (CmdConfig{}).FlatCmd(rcv, cmd, args...)
}

////////////////////////////////////////////////////////////////////////////////

// Maybe is a type which wraps a receiver being unmarshaled into. When
// unmarshaling takes place Maybe will also populate its other fields
// accordingly.
type Maybe struct {
	// Rcv is the receiver which will be unmarshaled into.
	Rcv interface{}

	// Null will be true if a null RESP value is unmarshaled.
	Null bool

	// Empty will be true if an empty aggregated RESP type (array, set, map,
	// push, or attribute) is unmarshaled.
	Empty bool
}

// UnmarshalRESP implements the method for the resp.Unmarshaler interface.
func (mb *Maybe) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	var rm resp3.RawMessage
	if err := rm.UnmarshalRESP(br, o); err != nil {
		return err
	}

	mb.Null = rm.IsNull()
	mb.Empty = rm.IsEmpty()
	return rm.UnmarshalInto(mb.Rcv, o)
}

////////////////////////////////////////////////////////////////////////////////

// Tuple is a helper type which can be used when unmarshaling a RESP array.
// Each element of Tuple should be a pointer receiver which the corresponding
// element of the RESP array will be unmarshaled into, or nil to skip that
// element. The length of Tuple must match the length of the RESP array being
// unmarshaled.
//
// Tuple is useful when unmarshaling the results from commands like EXEC and
// EVAL.
type Tuple []interface{}

// UnmarshalRESP implements the method for the resp.Unmarshaler interface.
func (t Tuple) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	var ah resp3.ArrayHeader
	if err := ah.UnmarshalRESP(br, o); err != nil {
		return err
	} else if ah.NumElems != len(t) {
		for i := 0; i < ah.NumElems; i++ {
			if err := resp3.Unmarshal(br, nil, o); err != nil {
				return err
			}
		}
		return resp.ErrConnUsable{
			Err: fmt.Errorf("expected array of size %d but got array of size %d", len(t), ah.NumElems),
		}
	}

	var retErr error
	for i := 0; i < ah.NumElems; i++ {
		if err := resp3.Unmarshal(br, t[i], o); err != nil {
			// if the message was discarded then we can just continue, this
			// method will return the first error it sees
			if !errors.As(err, new(resp.ErrConnUsable)) {
				return err
			} else if retErr == nil {
				retErr = err
			}
		}
	}
	return retErr
}

////////////////////////////////////////////////////////////////////////////////

// EvalScript contains the body of a script to be used with redis' EVAL
// functionality. Call Cmd on a EvalScript to actually create an Action which
// can be run.
type EvalScript struct {
	script, sum string
}

// NewEvalScript initializes a EvalScript instance with the given script.
func NewEvalScript(script string) EvalScript {
	sumRaw := sha1.Sum([]byte(script))
	sum := hex.EncodeToString(sumRaw[:])
	return EvalScript{
		script: script,
		sum:    sum,
	}
}

var (
	evalsha = []byte("EVALSHA")
	eval    = []byte("EVAL")
)

type evalAction struct {
	EvalScript
	keys, args []string
	rcv        interface{}

	flattenErr error
	eval       bool
}

// Cmd is like the top-level Cmd but it uses the the EvalScript to perform an
// EVALSHA command (and will automatically fallback to EVAL as necessary).
func (es EvalScript) Cmd(rcv interface{}, keys []string, args ...string) Action {
	return &evalAction{
		EvalScript: es,
		keys:       keys,
		args:       args,
		rcv:        rcv,
	}
}

// FlatCmd is like the top level FlatCmd except it uses the EvalScript to
// perform an EVALSHA command (and will automatically fallback to EVAL as
// necessary).
func (es EvalScript) FlatCmd(rcv interface{}, keys []string, args ...interface{}) Action {
	ec := &evalAction{
		EvalScript: es,
		keys:       keys,
		rcv:        rcv,
	}
	ec.args, ec.flattenErr = resp3.Flatten(args, nil)
	return ec
}

func (ec *evalAction) Properties() ActionProperties {
	return ActionProperties{
		Keys:     ec.keys,
		CanRetry: true,

		// EvalScript doesn't work within a Pipeline because the initial EVALSHA
		// might return a NOSCRIPT, in which case a second EVAL must be
		// performed. If done in a Pipeline there would be no opportunity to
		// perform the second EVAL.
		CanPipeline: false,

		CanShareConn: true,
	}
}

func (ec *evalAction) MarshalRESP(w io.Writer, o *resp.Opts) error {
	if ec.flattenErr != nil {
		return ec.flattenErr
	}

	// EVAL(SHA) script/sum numkeys keys... args...
	ah := resp3.ArrayHeader{NumElems: 3 + len(ec.keys) + len(ec.args)}
	err := ah.MarshalRESP(w, o)

	if ec.eval {
		err = marshalBlobStringBytes(err, w, eval, o)
		err = marshalBlobString(err, w, ec.script, o)
	} else {
		err = marshalBlobStringBytes(err, w, evalsha, o)
		err = marshalBlobString(err, w, ec.sum, o)
	}

	err = marshalBlobString(err, w, strconv.Itoa(len(ec.keys)), o)
	for i := range ec.keys {
		err = marshalBlobString(err, w, ec.keys[i], o)
	}
	for i := range ec.args {
		err = marshalBlobString(err, w, ec.args[i], o)
	}
	return err
}

func (ec *evalAction) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	return resp3.Unmarshal(br, ec.rcv, o)
}

func (ec *evalAction) Perform(ctx context.Context, conn Conn) error {
	run := func(eval bool) error {
		ec.eval = eval
		return conn.EncodeDecode(ctx, ec, ec)
	}

	err := run(false)
	if rErr := (resp3.SimpleError{}); errors.As(err, &rErr) && strings.HasPrefix(rErr.Error(), "NOSCRIPT") {
		err = run(true)
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////

type pipelineMarshalerUnmarshaler struct {
	marshal, unmarshalInto interface{}
	err                    error
}

// pipeline contains all the fields of Pipeline as well as some methods we'd
// rather not expose to users.
type pipeline struct {
	// preallocated buffers of slices to avoid allocating for small pipelines
	actionsBuf        [8]Action
	mmBuf             [8]pipelineMarshalerUnmarshaler
	propertiesKeysBuf [8]string

	actions    []Action
	mm         []pipelineMarshalerUnmarshaler
	properties ActionProperties

	// Conn is only set during the Perform method. It's primary purpose is to
	// provide for the methods which aren't EncodeDecode during the inner
	// Actions' Perform calls, in the unlikely event that they are needed.
	Conn
}

var _ Conn = new(pipeline)

func (p *pipeline) init() {
	p.actions = p.actionsBuf[:0]
	p.mm = p.mmBuf[:0]
	p.properties = ActionProperties{
		Keys:         p.propertiesKeysBuf[:0],
		CanPipeline:  true, // obviously
		CanShareConn: true,
	}
}

func (p *pipeline) reset() {
	p.actions = p.actions[:0]
	p.mm = p.mm[:0]
	p.properties.Keys = p.properties.Keys[:0]
	p.properties.CanShareConn = true
}

func (p *pipeline) append(a Action) {
	props := a.Properties()
	if !props.CanPipeline {
		panic(fmt.Sprintf("can't pipeline Action of type %T: %+v", a, a))
	}
	p.properties.Keys = append(p.properties.Keys, props.Keys...)
	p.properties.CanShareConn = p.properties.CanShareConn && props.CanShareConn
	p.actions = append(p.actions, a)
}

func (p *pipeline) EncodeDecode(_ context.Context, m, u interface{}) error {
	p.mm = append(p.mm, pipelineMarshalerUnmarshaler{
		marshal: m, unmarshalInto: u,
	})
	return nil
}

func (p *pipeline) MarshalRESP(w io.Writer, o *resp.Opts) error {
	for i := range p.mm {
		if p.mm[i].marshal == nil {
			// skip
		} else if err := resp3.Marshal(w, p.mm[i].marshal, o); err == nil {
			// ok
		} else if errors.As(err, new(resp.ErrConnUsable)) {
			// if the connection is still usable then mark this
			// pipelineMarshalerUnmarshaler as having had an error but continue
			// on to the rest...
			p.mm[i].err = err
		} else {
			return err
		}
	}
	return nil
}

func (p *pipeline) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	for i := range p.mm {
		if p.mm[i].unmarshalInto == nil || p.mm[i].err != nil {
			// skip
		} else if err := resp3.Unmarshal(br, p.mm[i].unmarshalInto, o); err == nil {
			// ok
		} else if errors.As(err, new(resp.ErrConnUsable)) {
			// if the connection is still usable then mark this
			// pipelineMarshalerUnmarshaler as having had an error but continue
			// on to the rest...
			p.mm[i].err = err
		} else {
			return err
		}
	}
	return nil
}

func (p *pipeline) Perform(ctx context.Context, c Conn) error {
	p.Conn = c
	defer func() { p.Conn = nil }()

	for _, action := range p.actions {
		// any errors that happen within Perform will not be IO errors, because
		// pipeline is suppressing all potential IO errors
		if err := action.Perform(ctx, p); err != nil {
			return resp.ErrConnUsable{Err: err}
		}
	}

	// pipeline's Marshal/UnmarshalRESP methods only return conn-unusable
	// errors. There's not much to be done for those except return the error for
	// the entire pipeline.
	if err := c.EncodeDecode(ctx, p, p); err != nil {
		return err
	}

	// look through the errors encountered on individual
	// pipelineMarshalerUnmarshalers, they will be conn-usable. Return the first
	// one in a nicely formatted way so the caller knows that something went
	// wrong but the connection is still usable.
	for _, m := range p.mm {
		if m.err != nil {
			return fmt.Errorf("command %+v in pipeline returned error: %w", m.marshal, m.err)
		}
	}

	return nil
}

// Pipeline is an Action which combines multiple commands into a single network
// round-trip. Pipeline accumulates commands via its Append method. When
// Pipeline is performed (i.e. passed into a Client's Do method) it will first
// write all commands as a single write operation and then read all command
// responses with a single read operation.
//
// Pipeline may be Reset in order to re-use an instance for multiple sets of
// commands. A Pipeline may _not_ be performed multiple times without being
// Reset in between.
//
// NOTE that, while a Pipeline performs all commands on a single Conn, it
// shouldn't be used by itself for MULTI/EXEC transactions, because if there's
// an error it won't discard the incomplete transaction. Use WithConn or
// EvalScript for transactional functionality instead.
type Pipeline struct {
	pipeline pipeline
}

// NewPipeline returns a Pipeline instance to which Actions can be Appended.
func NewPipeline() *Pipeline {
	p := &Pipeline{}
	p.pipeline.init()
	return p
}

// Reset discards all Actions and resets all internal state. A Pipeline with
// Reset called on it is equivalent to one returned by NewPipeline.
func (p *Pipeline) Reset() {
	p.pipeline.reset()
}

// Append adds the Action to the end of the list of Actions to pipeline
// together. This will panic if given an Action without the CanPipeline property
// set to true.
func (p *Pipeline) Append(a Action) {
	p.pipeline.append(a)
}

// Properties implements the method for the Action interface.
func (p *Pipeline) Properties() ActionProperties {
	return p.pipeline.properties
}

// Perform implements the method for the Action interface.
func (p *Pipeline) Perform(ctx context.Context, c Conn) error {
	return p.pipeline.Perform(ctx, c)
}

////////////////////////////////////////////////////////////////////////////////

type withConn struct {
	key [1]string // use array to avoid allocation in Properties
	fn  func(context.Context, Conn) error
}

// WithConn is used to perform a set of independent Actions on the same Conn.
//
// key should be a key which one or more of the inner Actions is going to act
// on, or "" if no keys are being acted on or the keys aren't yet known. key is
// generally only necessary when using Cluster.
//
// The callback function is what should actually carry out the inner actions,
// and the error it returns will be passed back up immediately.
//
// NOTE that WithConn only ensures all inner Actions are performed on the same
// Conn, it doesn't make them transactional. Use MULTI/WATCH/EXEC within a
// WithConn for transactions, or use EvalScript.
func WithConn(key string, fn func(context.Context, Conn) error) Action {
	return &withConn{[1]string{key}, fn}
}

func (wc *withConn) Properties() ActionProperties {
	return ActionProperties{
		Keys: wc.key[:],
	}
}

func (wc *withConn) Perform(ctx context.Context, c Conn) error {
	return wc.fn(ctx, c)
}
