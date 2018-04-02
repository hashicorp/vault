package plugin

import hclog "github.com/hashicorp/go-hclog"

type LoggerServer struct {
	logger hclog.Logger
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
	l.logger.Warn(args.Msg, args.Args...)
	return nil
}

func (l *LoggerServer) Error(args *LoggerArgs, reply *LoggerReply) error {
	l.logger.Error(args.Msg, args.Args...)
	return nil
}

func (l *LoggerServer) Log(args *LoggerArgs, _ *struct{}) error {

	switch translateLevel(args.Level) {

	case hclog.Trace:
		l.logger.Trace(args.Msg, args.Args...)

	case hclog.Debug:
		l.logger.Debug(args.Msg, args.Args...)

	case hclog.Info:
		l.logger.Info(args.Msg, args.Args...)

	case hclog.Warn:
		l.logger.Warn(args.Msg, args.Args...)

	case hclog.Error:
		l.logger.Error(args.Msg, args.Args...)

	case hclog.NoLevel:
	}
	return nil
}

func (l *LoggerServer) SetLevel(args int, _ *struct{}) error {
	level := translateLevel(args)
	l.logger = hclog.New(&hclog.LoggerOptions{Level: level})
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

func translateLevel(logxiLevel int) hclog.Level {

	switch logxiLevel {

	case 1000, 10:
		// logxi.LevelAll, logxi.LevelTrace:
		return hclog.Trace

	case 7:
		// logxi.LevelDebug:
		return hclog.Debug

	case 6, 5:
		// logxi.LevelInfo, logxi.LevelNotice:
		return hclog.Info

	case 4:
		// logxi.LevelWarn:
		return hclog.Warn

	case 3, 2, 1, -1:
		// logxi.LevelError, logxi.LevelFatal, logxi.LevelAlert, logxi.LevelEmergency:
		return hclog.Error
	}
	return hclog.NoLevel
}
