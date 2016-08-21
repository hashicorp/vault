package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

// scream so user fixes it
const warnImbalancedKey = "FIX_IMBALANCED_PAIRS"
const warnImbalancedPairs = warnImbalancedKey + " => "
const singleArgKey = "_"

func badKeyAtIndex(i int) string {
	return "BAD_KEY_AT_INDEX_" + strconv.Itoa(i)
}

// DefaultLogLog is the default log for this package.
var DefaultLog Logger

// Suppress supresses logging and is useful to supress output in
// in unit tests.
//
// Example
// log.Suppress(true)
// defer log.suppress(false)
func Suppress(quiet bool) {
	silent = quiet
}

var silent bool

// internalLog is the logger used by logxi itself
var InternalLog Logger

type loggerMap struct {
	sync.Mutex
	loggers map[string]Logger
}

var loggers = &loggerMap{
	loggers: map[string]Logger{},
}

func (lm *loggerMap) set(name string, logger Logger) {
	lm.loggers[name] = logger
}

// The assignment character between key-value pairs
var AssignmentChar = ": "

// Separator is the separator to use between key value pairs
//var Separator = "{~}"
var Separator = " "

const ltsvAssignmentChar = ":"
const ltsvSeparator = "\t"

// logxiEnabledMap maps log name patterns to levels
var logxiNameLevelMap map[string]int

// logxiFormat is the formatter kind to create
var logxiFormat string

var colorableStdout io.Writer
var defaultContextLines = 2
var defaultFormat string
var defaultLevel int
var defaultLogxiEnv string
var defaultLogxiFormatEnv string
var defaultMaxCol = 80
var defaultPretty = false
var defaultLogxiColorsEnv string
var defaultTimeFormat string
var disableCallstack bool
var disableCheckKeys bool
var disableColors bool
var home string
var isPretty bool
var isTerminal bool
var isWindows = runtime.GOOS == "windows"
var pkgMutex sync.Mutex
var pool = NewBufferPool()
var timeFormat string
var wd string
var pid = os.Getpid()
var pidStr = strconv.Itoa(os.Getpid())

// KeyMapping is the key map used to print built-in log entry fields.
type KeyMapping struct {
	Level     string
	Message   string
	Name      string
	PID       string
	Time      string
	CallStack string
}

// KeyMap is the key map to use when printing log statements.
var KeyMap = &KeyMapping{
	Level:     "_l",
	Message:   "_m",
	Name:      "_n",
	PID:       "_p",
	Time:      "_t",
	CallStack: "_c",
}

var logxiKeys []string

func setDefaults(isTerminal bool) {
	var err error
	contextLines = defaultContextLines
	wd, err = os.Getwd()
	if err != nil {
		InternalLog.Error("Could not get working directory")
	}

	logxiKeys = []string{KeyMap.Level, KeyMap.Message, KeyMap.Name, KeyMap.Time, KeyMap.CallStack, KeyMap.PID}

	if isTerminal {
		defaultLogxiEnv = "*=WRN"
		defaultLogxiFormatEnv = "happy,fit,maxcol=80,t=15:04:05.000000,context=-1"
		defaultFormat = FormatHappy
		defaultLevel = LevelWarn
		defaultTimeFormat = "15:04:05.000000"
	} else {
		defaultLogxiEnv = "*=ERR"
		defaultLogxiFormatEnv = "JSON,t=2006-01-02T15:04:05-0700"
		defaultFormat = FormatJSON
		defaultLevel = LevelError
		defaultTimeFormat = "2006-01-02T15:04:05-0700"
		disableColors = true
	}

	if isWindows {
		home = os.Getenv("HOMEPATH")
		if os.Getenv("ConEmuANSI") == "ON" {
			defaultLogxiColorsEnv = "key=cyan+h,value,misc=blue+h,source=yellow,TRC,DBG,WRN=yellow+h,INF=green+h,ERR=red+h"
		} else {
			colorableStdout = NewConcurrentWriter(colorable.NewColorableStdout())
			defaultLogxiColorsEnv = "ERR=red,misc=cyan,key=cyan"
		}
		// DefaultScheme is a color scheme optimized for dark background
		// but works well with light backgrounds
	} else {
		home = os.Getenv("HOME")
		term := os.Getenv("TERM")
		if term == "xterm-256color" {
			defaultLogxiColorsEnv = "key=cyan+h,value,misc=blue,source=88,TRC,DBG,WRN=yellow,INF=green+h,ERR=red+h,message=magenta+h"
		} else {
			defaultLogxiColorsEnv = "key=cyan+h,value,misc=blue,source=magenta,TRC,DBG,WRN=yellow,INF=green,ERR=red+h"
		}
	}
}

func isReservedKey(k interface{}) (bool, error) {
	key, ok := k.(string)
	if !ok {
		return false, fmt.Errorf("Key is not a string")
	}

	// check if reserved
	for _, key2 := range logxiKeys {
		if key == key2 {
			return true, nil
		}
	}
	return false, nil
}

func init() {
	colorableStdout = NewConcurrentWriter(os.Stdout)

	isTerminal = isatty.IsTerminal(os.Stdout.Fd())

	// the internal logger to report errors
	if isTerminal {
		InternalLog = NewLogger3(NewConcurrentWriter(os.Stdout), "__logxi", NewTextFormatter("__logxi"))
	} else {
		InternalLog = NewLogger3(NewConcurrentWriter(os.Stdout), "__logxi", NewJSONFormatter("__logxi"))
	}
	InternalLog.SetLevel(LevelError)

	setDefaults(isTerminal)

	RegisterFormatFactory(FormatHappy, formatFactory)
	RegisterFormatFactory(FormatText, formatFactory)
	RegisterFormatFactory(FormatJSON, formatFactory)
	ProcessEnv(readFromEnviron())

	// package logger for users
	DefaultLog = New("~")
}
