package grpcgcp

import (
	"strings"

	"google.golang.org/grpc/grpclog"
)

const (
	FINE   = 90
	FINEST = 99
)

var compLogger = grpclog.Component("grpcgcp")

type gcpLogger struct {
	logger grpclog.LoggerV2
	prefix string
}

// Make sure gcpLogger implements grpclog.LoggerV2.
var _ grpclog.LoggerV2 = (*gcpLogger)(nil)

func NewGCPLogger(logger grpclog.LoggerV2, prefix string) *gcpLogger {
	p := prefix
	if !strings.HasSuffix(p, " ") {
		p = p + " "
	}
	return &gcpLogger{
		logger: logger,
		prefix: p,
	}
}

// Error implements grpclog.LoggerV2.
func (l *gcpLogger) Error(args ...interface{}) {
	l.logger.Error(append([]interface{}{l.prefix}, args)...)
}

// Errorf implements grpclog.LoggerV2.
func (l *gcpLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(l.prefix+format, args...)
}

// Errorln implements grpclog.LoggerV2.
func (l *gcpLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(append([]interface{}{l.prefix}, args)...)
}

// Fatal implements grpclog.LoggerV2.
func (l *gcpLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(append([]interface{}{l.prefix}, args)...)
}

// Fatalf implements grpclog.LoggerV2.
func (l *gcpLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(l.prefix+format, args...)
}

// Fatalln implements grpclog.LoggerV2.
func (l *gcpLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(append([]interface{}{l.prefix}, args)...)
}

// Info implements grpclog.LoggerV2.
func (l *gcpLogger) Info(args ...interface{}) {
	l.logger.Info(append([]interface{}{l.prefix}, args)...)
}

// Infof implements grpclog.LoggerV2.
func (l *gcpLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(l.prefix+format, args...)
}

// Infoln implements grpclog.LoggerV2.
func (l *gcpLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(append([]interface{}{l.prefix}, args)...)
}

// V implements grpclog.LoggerV2.
func (l *gcpLogger) V(level int) bool {
	return l.logger.V(level)
}

// Warning implements grpclog.LoggerV2.
func (l *gcpLogger) Warning(args ...interface{}) {
	l.logger.Warning(append([]interface{}{l.prefix}, args)...)
}

// Warningf implements grpclog.LoggerV2.
func (l *gcpLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(l.prefix+format, args...)
}

// Warningln implements grpclog.LoggerV2.
func (l *gcpLogger) Warningln(args ...interface{}) {
	l.logger.Warningln(append([]interface{}{l.prefix}, args)...)
}
