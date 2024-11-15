package linodego

import (
	"log"
	"os"
)

//nolint:unused
type httpLogger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

//nolint:unused
type logger struct {
	l *log.Logger
}

//nolint:unused
func createLogger() *logger {
	l := &logger{l: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

//nolint:unused
var _ httpLogger = (*logger)(nil)

//nolint:unused
func (l *logger) Errorf(format string, v ...interface{}) {
	l.output("ERROR RESTY "+format, v...)
}

//nolint:unused
func (l *logger) Warnf(format string, v ...interface{}) {
	l.output("WARN RESTY "+format, v...)
}

//nolint:unused
func (l *logger) Debugf(format string, v ...interface{}) {
	l.output("DEBUG RESTY "+format, v...)
}

//nolint:unused
func (l *logger) output(format string, v ...interface{}) { //nolint:goprintffuncname
	if len(v) == 0 {
		l.l.Print(format)
		return
	}
	l.l.Printf(format, v...)
}
