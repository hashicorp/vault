package log

/*
http://en.wikipedia.org/wiki/Syslog

Code	Severity	Keyword
0	Emergency	emerg (panic)	System is unusable.

	A "panic" condition usually affecting multiple apps/servers/sites. At this
	level it would usually notify all tech staff on call.

1	Alert	alert	Action must be taken immediately.

	Should be corrected immediately, therefore notify staff who can fix the
	problem. An example would be the loss of a primary ISP connection.

2	Critical	crit	Critical conditions.

	Should be corrected immediately, but indicates failure in a secondary
	system, an example is a loss of a backup ISP connection.

3	Error	err (error)	Error conditions.

	Non-urgent failures, these should be relayed to developers or admins; each
	item must be resolved within a given time.

4	Warning	warning (warn)	Warning conditions.

	Warning messages, not an error, but indication that an error will occur if
	action is not taken, e.g. file system 85% full - each item must be resolved
	within a given time.

5	Notice	notice	Normal but significant condition.

	Events that are unusual but not error conditions - might be summarized in
	an email to developers or admins to spot potential problems - no immediate
	action required.

6	Informational	info	Informational messages.

	Normal operational messages - may be harvested for reporting, measuring
	throughput, etc. - no action required.

7	Debug	debug	Debug-level messages.

	Info useful to developers for debugging the application, not useful during operations.
*/

const (
	// LevelEnv chooses level from LOGXI environment variable or defaults
	// to LevelInfo
	LevelEnv = -10000

	// LevelOff means logging is disabled for logger. This should always
	// be first
	LevelOff = -1000

	// LevelEmergency is usually 0 but that is also the "zero" value
	// for Go, which means whenever we do any lookup in string -> int
	// map 0 is returned (not good).
	LevelEmergency = -1

	// LevelAlert means action must be taken immediately.
	LevelAlert = 1

	// LevelFatal means it should be corrected immediately, eg cannot connect to database.
	LevelFatal = 2

	// LevelCritical is alias for LevelFatal
	LevelCritical = 2

	// LevelError is a non-urgen failure to notify devlopers or admins
	LevelError = 3

	// LevelWarn indiates an error will occur if action is not taken, eg file system 85% full
	LevelWarn = 4

	// LevelNotice is normal but significant condition.
	LevelNotice = 5

	// LevelInfo is info level
	LevelInfo = 6

	// LevelDebug is debug level
	LevelDebug = 7

	// LevelTrace is trace level and displays file and line in terminal
	LevelTrace = 10

	// LevelAll is all levels
	LevelAll = 1000
)

// FormatHappy uses HappyDevFormatter
const FormatHappy = "happy"

// FormatText uses TextFormatter
const FormatText = "text"

// FormatJSON uses JSONFormatter
const FormatJSON = "JSON"

// FormatEnv selects formatter based on LOGXI_FORMAT environment variable
const FormatEnv = ""

// LevelMap maps int enums to string level.
var LevelMap = map[int]string{
	LevelFatal: "FTL",
	LevelError: "ERR",
	LevelWarn:  "WRN",
	LevelInfo:  "INF",
	LevelDebug: "DBG",
	LevelTrace: "TRC",
}

// LevelMap maps int enums to string level.
var LevelAtoi = map[string]int{
	"OFF": LevelOff,
	"FTL": LevelFatal,
	"ERR": LevelError,
	"WRN": LevelWarn,
	"INF": LevelInfo,
	"DBG": LevelDebug,
	"TRC": LevelTrace,
	"ALL": LevelAll,

	"off":   LevelOff,
	"fatal": LevelFatal,
	"error": LevelError,
	"warn":  LevelWarn,
	"info":  LevelInfo,
	"debug": LevelDebug,
	"trace": LevelTrace,
	"all":   LevelAll,
}

// Logger is the interface for logging.
type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{})
	Log(level int, msg string, args []interface{})

	SetLevel(int)
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	// Error, Fatal not needed, those SHOULD always be logged
}
