package log

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"
)

type bufferWriter interface {
	Write(p []byte) (nn int, err error)
	WriteRune(r rune) (n int, err error)
	WriteString(s string) (n int, err error)
}

// JSONFormatter is a fast, efficient JSON formatter optimized for logging.
//
// * log entry keys are not escaped
//   Who uses complex keys when coding? Checked by HappyDevFormatter in case user does.
//   Nested object keys are escaped by json.Marshal().
// * Primitive types uses strconv
// * Logger reserved key values (time, log name, level) require no conversion
// * sync.Pool buffer for bytes.Buffer
type JSONFormatter struct {
	name string
}

// NewJSONFormatter creates a new instance of JSONFormatter.
func NewJSONFormatter(name string) *JSONFormatter {
	return &JSONFormatter{name: name}
}

func (jf *JSONFormatter) writeString(buf bufferWriter, s string) {
	b, err := json.Marshal(s)
	if err != nil {
		InternalLog.Error("Could not json.Marshal string.", "str", s)
		buf.WriteString(`"Could not marshal this key's string"`)
		return
	}
	buf.Write(b)
}

func (jf *JSONFormatter) writeError(buf bufferWriter, err error) {
	jf.writeString(buf, err.Error())
	jf.set(buf, KeyMap.CallStack, string(debug.Stack()))
	return
}

func (jf *JSONFormatter) appendValue(buf bufferWriter, val interface{}) {
	if val == nil {
		buf.WriteString("null")
		return
	}

	// always show error stack even at cost of some performance. there's
	// nothing worse than looking at production logs without a clue
	if err, ok := val.(error); ok {
		jf.writeError(buf, err)
		return
	}

	value := reflect.ValueOf(val)
	kind := value.Kind()
	if kind == reflect.Ptr {
		if value.IsNil() {
			buf.WriteString("null")
			return
		}
		value = value.Elem()
		kind = value.Kind()
	}
	switch kind {
	case reflect.Bool:
		if value.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(value.Int(), 10))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteString(strconv.FormatUint(value.Uint(), 10))

	case reflect.Float32:
		buf.WriteString(strconv.FormatFloat(value.Float(), 'g', -1, 32))

	case reflect.Float64:
		buf.WriteString(strconv.FormatFloat(value.Float(), 'g', -1, 64))

	default:
		var err error
		var b []byte
		if stringer, ok := val.(fmt.Stringer); ok {
			b, err = json.Marshal(stringer.String())
		} else {
			b, err = json.Marshal(val)
		}

		if err != nil {
			InternalLog.Error("Could not json.Marshal value: ", "formatter", "JSONFormatter", "err", err.Error())
			if s, ok := val.(string); ok {
				b, err = json.Marshal(s)
			} else if s, ok := val.(fmt.Stringer); ok {
				b, err = json.Marshal(s.String())
			} else {
				b, err = json.Marshal(fmt.Sprintf("%#v", val))
			}

			if err != nil {
				// should never get here, but JSONFormatter should never panic
				msg := "Could not Sprintf value"
				InternalLog.Error(msg)
				buf.WriteString(`"` + msg + `"`)
				return
			}
		}
		buf.Write(b)
	}
}

func (jf *JSONFormatter) set(buf bufferWriter, key string, val interface{}) {
	// WARNING: assumes this is not first key
	buf.WriteString(`, "`)
	buf.WriteString(key)
	buf.WriteString(`":`)
	jf.appendValue(buf, val)
}

// Format formats log entry as JSON.
func (jf *JSONFormatter) Format(writer io.Writer, level int, msg string, args []interface{}) {
	buf := pool.Get()
	defer pool.Put(buf)

	const lead = `", "`
	const colon = `":"`

	buf.WriteString(`{"`)
	buf.WriteString(KeyMap.Time)
	buf.WriteString(`":"`)
	buf.WriteString(time.Now().Format(timeFormat))

	buf.WriteString(`", "`)
	buf.WriteString(KeyMap.PID)
	buf.WriteString(`":"`)
	buf.WriteString(pidStr)

	buf.WriteString(`", "`)
	buf.WriteString(KeyMap.Level)
	buf.WriteString(`":"`)
	buf.WriteString(LevelMap[level])

	buf.WriteString(`", "`)
	buf.WriteString(KeyMap.Name)
	buf.WriteString(`":"`)
	buf.WriteString(jf.name)

	buf.WriteString(`", "`)
	buf.WriteString(KeyMap.Message)
	buf.WriteString(`":`)
	jf.appendValue(buf, msg)

	var lenArgs = len(args)
	if lenArgs > 0 {
		if lenArgs == 1 {
			jf.set(buf, singleArgKey, args[0])
		} else if lenArgs%2 == 0 {
			for i := 0; i < lenArgs; i += 2 {
				if key, ok := args[i].(string); ok {
					if key == "" {
						// show key is invalid
						jf.set(buf, badKeyAtIndex(i), args[i+1])
					} else {
						jf.set(buf, key, args[i+1])
					}
				} else {
					// show key is invalid
					jf.set(buf, badKeyAtIndex(i), args[i+1])
				}
			}
		} else {
			jf.set(buf, warnImbalancedKey, args)
		}
	}
	buf.WriteString("}\n")
	buf.WriteTo(writer)
}

// LogEntry returns the JSON log entry object built by Format(). Used by
// HappyDevFormatter to ensure any data logged while developing properly
// logs in production.
func (jf *JSONFormatter) LogEntry(level int, msg string, args []interface{}) map[string]interface{} {
	buf := pool.Get()
	defer pool.Put(buf)
	jf.Format(buf, level, msg, args)
	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	if err != nil {
		panic("Unable to unmarhsal entry from JSONFormatter: " + err.Error() + " \"" + string(buf.Bytes()) + "\"")
	}
	return entry
}
