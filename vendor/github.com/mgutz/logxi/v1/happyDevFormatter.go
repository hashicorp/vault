package log

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/mgutz/ansi"
)

// colorScheme defines a color theme for HappyDevFormatter
type colorScheme struct {
	Key     string
	Message string
	Value   string
	Misc    string
	Source  string

	Trace string
	Debug string
	Info  string
	Warn  string
	Error string
}

var indent = "  "
var maxCol = defaultMaxCol
var theme *colorScheme

func parseKVList(s, separator string) map[string]string {
	pairs := strings.Split(s, separator)
	if len(pairs) == 0 {
		return nil
	}
	m := map[string]string{}
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		parts := strings.Split(pair, "=")
		switch len(parts) {
		case 1:
			m[parts[0]] = ""
		case 2:
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func parseTheme(theme string) *colorScheme {
	m := parseKVList(theme, ",")
	cs := &colorScheme{}
	var wildcard string

	var color = func(key string) string {
		if disableColors {
			return ""
		}
		style := m[key]
		c := ansi.ColorCode(style)
		if c == "" {
			c = wildcard
		}
		//fmt.Printf("plain=%b [%s] %s=%q\n", ansi.DefaultFG, key, style, c)

		return c
	}
	wildcard = color("*")

	if wildcard != ansi.Reset {
		cs.Key = wildcard
		cs.Value = wildcard
		cs.Misc = wildcard
		cs.Source = wildcard
		cs.Message = wildcard

		cs.Trace = wildcard
		cs.Debug = wildcard
		cs.Warn = wildcard
		cs.Info = wildcard
		cs.Error = wildcard
	}

	cs.Key = color("key")
	cs.Value = color("value")
	cs.Misc = color("misc")
	cs.Source = color("source")
	cs.Message = color("message")

	cs.Trace = color("TRC")
	cs.Debug = color("DBG")
	cs.Warn = color("WRN")
	cs.Info = color("INF")
	cs.Error = color("ERR")
	return cs
}

// HappyDevFormatter is the formatter used for terminals. It is
// colorful, dev friendly and provides meaningful logs when
// warnings and errors occur.
//
// HappyDevFormatter does not worry about performance. It's at least 3-4X
// slower than JSONFormatter since it delegates to JSONFormatter to marshal
// then unmarshal JSON. Then it does other stuff like read source files, sort
// keys all to give a developer more information.
//
// SHOULD NOT be used in production for extended period of time. However, it
// works fine in SSH terminals and binary deployments.
type HappyDevFormatter struct {
	name string
	col  int
	// always use the production formatter
	jsonFormatter *JSONFormatter
}

// NewHappyDevFormatter returns a new instance of HappyDevFormatter.
func NewHappyDevFormatter(name string) *HappyDevFormatter {
	jf := NewJSONFormatter(name)
	return &HappyDevFormatter{
		name:          name,
		jsonFormatter: jf,
	}
}

func (hd *HappyDevFormatter) writeKey(buf bufferWriter, key string) {
	// assumes this is not the first key
	hd.writeString(buf, Separator)
	if key == "" {
		return
	}
	buf.WriteString(theme.Key)
	hd.writeString(buf, key)
	hd.writeString(buf, AssignmentChar)
	if !disableColors {
		buf.WriteString(ansi.Reset)
	}
}

func (hd *HappyDevFormatter) set(buf bufferWriter, key string, value interface{}, color string) {
	var str string
	if s, ok := value.(string); ok {
		str = s
	} else if s, ok := value.(fmt.Stringer); ok {
		str = s.String()
	} else {
		str = fmt.Sprintf("%v", value)
	}
	val := strings.Trim(str, "\n ")
	if (isPretty && key != "") || hd.col+len(key)+2+len(val) >= maxCol {
		buf.WriteString("\n")
		hd.col = 0
		hd.writeString(buf, indent)
	}
	hd.writeKey(buf, key)
	if color != "" {
		buf.WriteString(color)
	}
	hd.writeString(buf, val)
	if color != "" && !disableColors {
		buf.WriteString(ansi.Reset)
	}
}

// Write a string and tracks the position of the string so we can break lines
// cleanly. Do not send ANSI escape sequences, just raw strings
func (hd *HappyDevFormatter) writeString(buf bufferWriter, s string) {
	buf.WriteString(s)
	hd.col += len(s)
}

func (hd *HappyDevFormatter) getContext(color string) string {
	if disableCallstack {
		return ""
	}
	frames := parseDebugStack(string(debug.Stack()), 5, true)
	if len(frames) == 0 {
		return ""
	}
	for _, frame := range frames {
		context := frame.String(color, theme.Source)
		if context != "" {
			return context
		}
	}
	return ""
}

func (hd *HappyDevFormatter) getLevelContext(level int, entry map[string]interface{}) (message string, context string, color string) {

	switch level {
	case LevelTrace:
		color = theme.Trace
		context = hd.getContext(color)
		context += "\n"
	case LevelDebug:
		color = theme.Debug
	case LevelInfo:
		color = theme.Info
	// case LevelWarn:
	// 	color = theme.Warn
	// 	context = hd.getContext(color)
	// 	context += "\n"
	case LevelWarn, LevelError, LevelFatal:

		// warnings return an error but if it does not have an error
		// then print line info only
		if level == LevelWarn {
			color = theme.Warn
			kv := entry[KeyMap.CallStack]
			if kv == nil {
				context = hd.getContext(color)
				context += "\n"
				break
			}
		} else {
			color = theme.Error
		}

		if disableCallstack || contextLines == -1 {
			context = trimDebugStack(string(debug.Stack()))
			break
		}
		frames := parseLogxiStack(entry, 4, true)
		if frames == nil {
			frames = parseDebugStack(string(debug.Stack()), 4, true)
		}

		if len(frames) == 0 {
			break
		}
		errbuf := pool.Get()
		defer pool.Put(errbuf)
		lines := 0
		for _, frame := range frames {
			err := frame.readSource(contextLines)
			if err != nil {
				// by setting to empty, the original stack is used
				errbuf.Reset()
				break
			}
			ctx := frame.String(color, theme.Source)
			if ctx == "" {
				continue
			}
			errbuf.WriteString(ctx)
			errbuf.WriteRune('\n')
			lines++
		}
		context = errbuf.String()
	default:
		panic("should never get here")
	}
	return message, context, color
}

// Format a log entry.
func (hd *HappyDevFormatter) Format(writer io.Writer, level int, msg string, args []interface{}) {
	buf := pool.Get()
	defer pool.Put(buf)

	if len(args) == 1 {
		args = append(args, 0)
		copy(args[1:], args[0:])
		args[0] = singleArgKey
	}

	// warn about reserved, bad and complex keys
	for i := 0; i < len(args); i += 2 {
		isReserved, err := isReservedKey(args[i])
		if err != nil {
			InternalLog.Error("Key is not a string.", "err", fmt.Errorf("args[%d]=%v", i, args[i]))
		} else if isReserved {
			InternalLog.Fatal("Key conflicts with reserved key. Avoiding using single rune keys.", "key", args[i].(string))
		} else {
			// Ensure keys are simple strings. The JSONFormatter doesn't escape
			// keys as a performance tradeoff. This panics if the JSON key
			// value has a different value than a simple quoted string.
			key := args[i].(string)
			b, err := json.Marshal(key)
			if err != nil {
				panic("Key is invalid. " + err.Error())
			}
			if string(b) != `"`+key+`"` {
				panic("Key is complex. Use simpler key for: " + fmt.Sprintf("%q", key))
			}
		}
	}

	// use the production JSON formatter to format the log first. This
	// ensures JSON will marshal/unmarshal correctly in production.
	entry := hd.jsonFormatter.LogEntry(level, msg, args)

	// reset the column tracker used for fancy formatting
	hd.col = 0

	// timestamp
	buf.WriteString(theme.Misc)
	hd.writeString(buf, entry[KeyMap.Time].(string))
	if !disableColors {
		buf.WriteString(ansi.Reset)
	}

	// emphasize warnings and errors
	message, context, color := hd.getLevelContext(level, entry)
	if message == "" {
		message = entry[KeyMap.Message].(string)
	}

	// DBG, INF ...
	hd.set(buf, "", entry[KeyMap.Level].(string), color)
	// logger name
	hd.set(buf, "", entry[KeyMap.Name], theme.Misc)
	// message from user
	hd.set(buf, "", message, theme.Message)

	// Preserve key order in the sequencethey were added by developer.This
	// makes it easier for developers to follow the log.
	order := []string{}
	lenArgs := len(args)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= lenArgs {
			continue
		}
		if key, ok := args[i].(string); ok {
			order = append(order, key)
		} else {
			order = append(order, badKeyAtIndex(i))
		}
	}

	for _, key := range order {
		// skip reserved keys which were already added to buffer above
		isReserved, err := isReservedKey(key)
		if err != nil {
			panic("key is invalid. Should never get here. " + err.Error())
		} else if isReserved {
			continue
		}
		hd.set(buf, key, entry[key], theme.Value)
	}

	addLF := true
	hasCallStack := entry[KeyMap.CallStack] != nil
	// WRN,ERR file, line number context

	if context != "" {
		// warnings and traces are single line, space can be optimized
		if level == LevelTrace || (level == LevelWarn && !hasCallStack) {
			// gets rid of "in "
			idx := strings.IndexRune(context, 'n')
			hd.set(buf, "in", context[idx+2:], color)
		} else {
			buf.WriteRune('\n')
			if !disableColors {
				buf.WriteString(color)
			}
			addLF = context[len(context)-1:len(context)] != "\n"
			buf.WriteString(context)
			if !disableColors {
				buf.WriteString(ansi.Reset)
			}
		}
	} else if hasCallStack {
		hd.set(buf, "", entry[KeyMap.CallStack], color)
	}
	if addLF {
		buf.WriteRune('\n')
	}
	buf.WriteTo(writer)
}
