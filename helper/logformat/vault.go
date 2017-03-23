package logformat

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/mgutz/logxi/v1"
)

const (
	styledefault = iota
	stylejson
)

// NewVaultLogger creates a new logger with the specified level and a Vault
// formatter
func NewVaultLogger(level int) log.Logger {
	logger := log.New("vault")
	return setLevelFormatter(logger, level, createVaultFormatter())
}

// NewVaultLoggerWithWriter creates a new logger with the specified level and
// writer and a Vault formatter
func NewVaultLoggerWithWriter(w io.Writer, level int) log.Logger {
	logger := log.NewLogger(w, "vault")
	return setLevelFormatter(logger, level, createVaultFormatter())
}

// Sets the level and formatter on the log, which must be a DefaultLogger
func setLevelFormatter(logger log.Logger, level int, formatter log.Formatter) log.Logger {
	logger.(*log.DefaultLogger).SetLevel(level)
	logger.(*log.DefaultLogger).SetFormatter(formatter)
	return logger
}

// Creates a formatter, checking env vars for the style
func createVaultFormatter() log.Formatter {
	ret := &vaultFormatter{
		Mutex: &sync.Mutex{},
	}
	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	if logFormat == "" {
		logFormat = os.Getenv("LOGXI_FORMAT")
	}
	switch strings.ToLower(logFormat) {
	case "json", "vault_json", "vault-json", "vaultjson":
		ret.style = stylejson
	default:
		ret.style = styledefault
	}
	return ret
}

// Thread safe formatter
type vaultFormatter struct {
	*sync.Mutex
	style  int
	module string
}

func (v *vaultFormatter) Format(writer io.Writer, level int, msg string, args []interface{}) {
	currTime := time.Now()
	v.Lock()
	defer v.Unlock()
	switch v.style {
	case stylejson:
		v.formatJSON(writer, currTime, level, msg, args)
	default:
		v.formatDefault(writer, currTime, level, msg, args)
	}
}

func (v *vaultFormatter) formatDefault(writer io.Writer, currTime time.Time, level int, msg string, args []interface{}) {
	// Write a trailing newline
	defer writer.Write([]byte("\n"))

	writer.Write([]byte(currTime.Local().Format("2006/01/02 15:04:05.000000")))

	switch level {
	case log.LevelCritical:
		writer.Write([]byte(" [CRIT ] "))
	case log.LevelError:
		writer.Write([]byte(" [ERROR] "))
	case log.LevelWarn:
		writer.Write([]byte(" [WARN ] "))
	case log.LevelInfo:
		writer.Write([]byte(" [INFO ] "))
	case log.LevelDebug:
		writer.Write([]byte(" [DEBUG] "))
	case log.LevelTrace:
		writer.Write([]byte(" [TRACE] "))
	default:
		writer.Write([]byte(" [ALL  ] "))
	}

	if v.module != "" {
		writer.Write([]byte(fmt.Sprintf("(%s) ", v.module)))
	}

	writer.Write([]byte(msg))

	if args != nil && len(args) > 0 {
		if len(args)%2 != 0 {
			args = append(args, "[unknown!]")
		}

		writer.Write([]byte(":"))

		for i := 0; i < len(args); i = i + 2 {
			var quote string
			switch args[i+1].(type) {
			case string:
				if strings.ContainsRune(args[i+1].(string), ' ') {
					quote = `"`
				}
			}
			writer.Write([]byte(fmt.Sprintf(" %s=%s%v%s", args[i], quote, args[i+1], quote)))
		}
	}
}

func (v *vaultFormatter) formatJSON(writer io.Writer, currTime time.Time, level int, msg string, args []interface{}) {
	vals := map[string]interface{}{
		"@message":   msg,
		"@timestamp": currTime.Format("2006-01-02T15:04:05.000000Z07:00"),
	}

	var levelStr string
	switch level {
	case log.LevelCritical:
		levelStr = "critical"
	case log.LevelError:
		levelStr = "error"
	case log.LevelWarn:
		levelStr = "warn"
	case log.LevelInfo:
		levelStr = "info"
	case log.LevelDebug:
		levelStr = "debug"
	case log.LevelTrace:
		levelStr = "trace"
	default:
		levelStr = "all"
	}

	vals["@level"] = levelStr

	if v.module != "" {
		vals["@module"] = v.module
	}

	if args != nil && len(args) > 0 {

		if len(args)%2 != 0 {
			args = append(args, "[unknown!]")
		}

		for i := 0; i < len(args); i = i + 2 {
			if _, ok := args[i].(string); !ok {
				// As this is the logging function not much we can do here
				// without injecting into logs...
				continue
			}
			vals[args[i].(string)] = args[i+1]
		}
	}

	enc := json.NewEncoder(writer)
	enc.Encode(vals)
}
