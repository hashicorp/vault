package plugin

import (
	"net/rpc"

	log_std "log"

	log "github.com/hashicorp/go-hclog"
)

type LoggerClient struct {
	client *rpc.Client
}

func (l *LoggerClient) Trace(msg string, args ...interface{}) {
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	l.client.Call("Plugin.Trace", cArgs, &struct{}{})
}

func (l *LoggerClient) Debug(msg string, args ...interface{}) {
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	l.client.Call("Plugin.Debug", cArgs, &struct{}{})
}

func (l *LoggerClient) Info(msg string, args ...interface{}) {
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	l.client.Call("Plugin.Info", cArgs, &struct{}{})
}
func (l *LoggerClient) Warn(msg string, args ...interface{}) {
	var reply LoggerReply
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	err := l.client.Call("Plugin.Warn", cArgs, &reply)
	_ = err

	//if err != nil {
	//	return err
	//}
	//if reply.Error != nil {
	//	return reply.Error
	//}
	//
	//return nil
}
func (l *LoggerClient) Error(msg string, args ...interface{}) {
	var reply LoggerReply
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	err := l.client.Call("Plugin.Error", cArgs, &reply)
	_ = err
	//if err != nil {
	//	return err
	//}
	//if reply.Error != nil {
	//	return reply.Error
	//}
	//
	//return nil
}

func (l *LoggerClient) Fatal(msg string, args ...interface{}) {
	// NOOP since it's not actually used within vault
	return
}

func (l *LoggerClient) Log(level int, msg string, args []interface{}) {
	cArgs := &LoggerArgs{
		Level: level,
		Msg:   msg,
		Args:  args,
	}
	l.client.Call("Plugin.Log", cArgs, &struct{}{})
}

func (l *LoggerClient) SetLevel(level int) {
	l.client.Call("Plugin.SetLevel", level, &struct{}{})
}

func (l *LoggerClient) IsTrace() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsTrace", new(interface{}), &reply)
	return reply.IsTrue
}
func (l *LoggerClient) IsDebug() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsDebug", new(interface{}), &reply)
	return reply.IsTrue
}

func (l *LoggerClient) IsInfo() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsInfo", new(interface{}), &reply)
	return reply.IsTrue
}

func (l *LoggerClient) IsWarn() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsWarn", new(interface{}), &reply)
	return reply.IsTrue
}

func (l *LoggerClient) IsError() bool {
	//var reply LoggerReply
	//l.client.Call("Plugin.IsError", new(interface{}), &reply)
	//return reply.IsTrue
	return false
}

// TODO: fix
func (l *LoggerClient) Named(s string) log.Logger {
	//var reply LoggerReply
	//l.client.Call("Plugin.IsError", new(interface{}), &reply)
	//return reply.IsTrue
	return log.Default()
}

// TODO: fix
func (l *LoggerClient) ResetNamed(s string) log.Logger {
	//var reply LoggerReply
	//l.client.Call("Plugin.IsError", new(interface{}), &reply)
	//return reply.IsTrue
	return log.Default()
}

// TODO: fix
func (l *LoggerClient) With(args ...interface{}) log.Logger {
	//var reply LoggerReply
	//l.client.Call("Plugin.IsError", new(interface{}), &reply)
	//return reply.IsTrue
	return log.Default()
}

// TODO: fix
func (l *LoggerClient) StandardLogger(opts *log.StandardLoggerOptions) *log_std.Logger {
	//var reply LoggerReply
	//l.client.Call("Plugin.IsError", new(interface{}), &reply)
	//return reply.IsTrue
	return &log_std.Logger{}
}

type LoggerServer struct {
	logger log.Logger
}

func (l *LoggerServer) Trace(args *LoggerArgs, _ *struct{}) error {
	l.logger.Trace(args.Msg, args.Args...)
	return nil
}

func (l *LoggerServer) Debug(args *LoggerArgs, _ *struct{}) error {
	l.logger.Debug(args.Msg, args.Args...)
	return nil
}

func (l *LoggerServer) Info(args *LoggerArgs, _ *struct{}) error {
	l.logger.Info(args.Msg, args.Args...)
	return nil
}

func (l *LoggerServer) Warn(args *LoggerArgs, reply *LoggerReply) error {
	//err := l.logger.Warn(args.Msg, args.Args...)
	l.logger.Warn(args.Msg, args.Args...)
	//if err != nil {
	//	*reply = LoggerReply{
	//		Error: wrapError(err),
	//	}
	//	return nil
	//}
	return nil
}

func (l *LoggerServer) Error(args *LoggerArgs, reply *LoggerReply) error {
	l.logger.Error(args.Msg, args.Args...)
	//if err != nil {
	//	*reply = LoggerReply{
	//		Error: wrapError(err),
	//	}
	//	return nil
	//}
	return nil
}

func (l *LoggerServer) Log(args *LoggerArgs, _ *struct{}) error {
	// TODO: FIX
	//l.logger.Log(args.Level, args.Msg, args.Args)
	return nil
}

func (l *LoggerServer) SetLevel(args int, _ *struct{}) error {
	// TODO: FIX
	//l.logger.SetLevel(args)
	return nil
}

func (l *LoggerServer) IsTrace(args interface{}, reply *LoggerReply) error {
	result := l.logger.IsTrace()
	*reply = LoggerReply{
		IsTrue: result,
	}
	return nil
}

func (l *LoggerServer) IsDebug(args interface{}, reply *LoggerReply) error {
	result := l.logger.IsDebug()
	*reply = LoggerReply{
		IsTrue: result,
	}
	return nil
}

func (l *LoggerServer) IsInfo(args interface{}, reply *LoggerReply) error {
	result := l.logger.IsInfo()
	*reply = LoggerReply{
		IsTrue: result,
	}
	return nil
}

func (l *LoggerServer) IsWarn(args interface{}, reply *LoggerReply) error {
	result := l.logger.IsWarn()
	*reply = LoggerReply{
		IsTrue: result,
	}
	return nil
}

type LoggerArgs struct {
	Level int
	Msg   string
	Args  []interface{}
}

// LoggerReply contains the RPC reply. Not all fields may be used
// for a particular RPC call.
type LoggerReply struct {
	IsTrue bool
	Error  error
}
