package gocb

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// gocbZapCore is a wrapper around our own Logger type which satisfies the zapcore.Core interface.
// This allows us to forward logging created by zap into any Logger specified by the user.
type gocbZapCore struct {
	enc zapcore.Encoder
}

func (g *gocbZapCore) clone() *gocbZapCore {
	return &gocbZapCore{
		enc: g.enc.Clone(),
	}
}

func newZapLogger() *zap.Logger {
	// This is pretty barebones as we just want to receive messages and we'll deal with them from there.
	return zap.New(&gocbZapCore{
		enc: zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:       "msg",
			SkipLineEnding:   true,
			ConsoleSeparator: " ",
		}),
	})
}

func (g *gocbZapCore) Enabled(level zapcore.Level) bool {
	// We don't know the log level so just pass-through messages allowing the higher level logger to filter.
	return true
}

func (g *gocbZapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := g.clone()

	for i := range fields {
		fields[i].AddTo(clone.enc)
	}

	return clone
}

func (g *gocbZapCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return checked.AddCore(entry, g)
}

func (g *gocbZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Using wrapEntry allows us to defer the call to EncodeEntry until it's
	// actually required - which means that we do not have to encode entries which
	// are at log levels which won't actually get logged.

	// offset of 2 lifts the line number out of this function, to  the actual origin.
	logExf(g.logLevel(entry.Level), 2, "%s", g.wrapEntry(entry, fields))

	return nil
}

func (g *gocbZapCore) Sync() error {
	return nil
}

func (g *gocbZapCore) logLevel(level zapcore.Level) LogLevel {
	switch level {
	case zapcore.FatalLevel:
		return LogError
	case zapcore.PanicLevel:
		return LogError
	case zapcore.DPanicLevel:
		return LogError
	case zapcore.ErrorLevel:
		return LogError
	case zapcore.WarnLevel:
		return LogWarn
	case zapcore.InfoLevel:
		return LogInfo
	default:
		return LogDebug
	}
}

func (g *gocbZapCore) wrapEntry(entry zapcore.Entry, fields []zapcore.Field) *zapLazyEntry {
	return &zapLazyEntry{
		wrapped: entry,
		enc:     g.enc,
		fields:  fields,
	}
}

type zapLazyEntry struct {
	wrapped zapcore.Entry
	enc     zapcore.Encoder
	fields  []zapcore.Field
}

func (z *zapLazyEntry) String() string {
	buf, err := z.enc.EncodeEntry(z.wrapped, z.fields)
	if err != nil {
		return "failed to encode log entry: " + err.Error()
	}

	str := buf.String()
	buf.Free()

	return str
}
