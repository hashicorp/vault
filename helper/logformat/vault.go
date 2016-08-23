package logformat

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	log "github.com/mgutz/logxi/v1"
)

const (
	styledefault = iota
	stylejson
)

func NewVaultLogger(level int) log.Logger {
	logger := log.New("vault")
	logger.(*log.DefaultLogger).SetLevel(level)
	logger.(*log.DefaultLogger).SetFormatter(createVaultFormatter())
	return logger
}

func NewVaultLoggerWithWriter(w io.Writer, level int) log.Logger {
	logger := log.NewLogger(w, "vault")
	logger.(*log.DefaultLogger).SetLevel(level)
	logger.(*log.DefaultLogger).SetFormatter(createVaultFormatter())
	return logger
}

func createVaultFormatter() log.Formatter {
	ret := &vaultFormatter{}
	switch os.Getenv("LOGXI_FORMAT") {
	case "vault_json", "vault-json", "vaultjson":
		ret.style = stylejson
	default:
		ret.style = styledefault
	}
	return ret
}

type vaultFormatter struct {
	sync.Mutex
	style int
}

func (v *vaultFormatter) Format(writer io.Writer, level int, msg string, args []interface{}) {
	v.Lock()
	defer v.Unlock()
	switch v.style {
	case stylejson:
		v.formatJSON(writer, level, msg, args)
	default:
		v.formatDefault(writer, level, msg, args)
	}
}

func (v *vaultFormatter) formatDefault(writer io.Writer, level int, msg string, args []interface{}) {
	defer writer.Write([]byte("\n"))

	writer.Write([]byte(time.Now().Local().Format("2006/01/02 15:04:05.000000")))

	switch level {
	case log.LevelError:
		writer.Write([]byte(" [ERR] "))
	case log.LevelWarn:
		writer.Write([]byte(" [WRN] "))
	case log.LevelInfo:
		writer.Write([]byte(" [INF] "))
	case log.LevelNotice:
		writer.Write([]byte(" [NOT] "))
	case log.LevelDebug:
		writer.Write([]byte(" [DBG] "))
	case log.LevelTrace:
		writer.Write([]byte(" [TRC] "))
	default:
		writer.Write([]byte(" [ALL] "))
	}

	writer.Write([]byte(msg))

	if args != nil && len(args) > 0 {
		if len(args)%2 != 0 {
			args = append(args, "[unknown!]")
		}

		for i := 0; i < len(args); i = i + 2 {
			writer.Write([]byte(fmt.Sprintf(" %s=%v", args[i], args[i+1])))
		}
	}
}

type logMsg struct {
	Stamp   string                 `json:"t"`
	Level   string                 `json:"l"`
	Message string                 `json:"m"`
	Args    map[string]interface{} `json:"a,omitempty"`
}

func (v *vaultFormatter) formatJSON(writer io.Writer, level int, msg string, args []interface{}) {
	l := &logMsg{
		Message: msg,
		Stamp:   time.Now().Format("2006-01-02T15:04:05.000000Z07:00"),
	}

	switch level {
	case log.LevelError:
		l.Level = "E"
	case log.LevelWarn:
		l.Level = "W"
	case log.LevelInfo:
		l.Level = "I"
	case log.LevelNotice:
		l.Level = "N"
	case log.LevelDebug:
		l.Level = "D"
	case log.LevelTrace:
		l.Level = "T"
	default:
		l.Level = "A"
	}

	if args != nil && len(args) > 0 {
		l.Args = make(map[string]interface{}, len(args)/2)

		if len(args)%2 != 0 {
			args = append(args, "[unknown!]")
		}

		for i := 0; i < len(args); i = i + 2 {
			if _, ok := args[i].(string); !ok {
				// As this is the logging function not much we can do here
				// without injecting into logs...
				continue
			}
			l.Args[args[i].(string)] = args[i+1]
		}
	}

	enc := json.NewEncoder(writer)
	enc.Encode(l)
}
