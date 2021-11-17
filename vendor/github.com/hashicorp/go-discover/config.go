package discover

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Config stores key/value pairs for the discovery
// functions to use.
type Config map[string]string

// Parse parses a "key=val key=val ..." string into a config map. Keys
// and values which contain spaces, backslashes or double-quotes must be
// quoted with double quotes. Use the backslash to escape special
// characters within quoted strings, e.g. "some key"="some \"value\"".
func Parse(s string) (Config, error) {
	return parse(s)
}

// String formats a config map into the "key=val key=val ..."
// understood by Parse. The order of the keys is stable.
func (c Config) String() string {
	// sort 'provider' to the front and keep the keys stable.
	var keys []string
	for k := range c {
		if k != "provider" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	keys = append([]string{"provider"}, keys...)

	quote := func(s string) string {
		if strings.ContainsAny(s, ` "\`) {
			return strconv.Quote(s)
		}
		return s
	}

	var vals []string
	for _, k := range keys {
		v := c[k]
		if v == "" {
			continue
		}
		vals = append(vals, quote(k)+"="+quote(v))
	}
	return strings.Join(vals, " ")
}

func parse(in string) (Config, error) {
	m := Config{}
	s := []rune(strings.TrimSpace(in))
	state := stateKey
	key := ""
	for {
		// exit condition
		if len(s) == 0 {
			break
		}

		// get the next token
		item, val, n := lex(s)
		s = s[n:]
		// fmt.Printf("parse: state: %q item: %q val: '%s' n: %d rest: '%s'\n", state, item, val, n, string(s))

		switch state {

		case stateKey:
			switch item {
			case itemText:
				key = val
				if _, exists := m[key]; exists {
					return nil, fmt.Errorf("%s: duplicate key", key)
				}
				state = stateEqual
			default:
				if val == "" {
					return nil, fmt.Errorf("%s: - equals in key's value, enclosing double-quote needed %s=\"value-with-=-symbol\"", key, key)
				}
				return nil, fmt.Errorf("%s: error with key=value pair %s", key, val)
			}

		case stateEqual:
			switch item {
			case itemEqual:
				state = stateVal
			default:
				return nil, fmt.Errorf("%s: missing '='", key)
			}

		case stateVal:
			switch item {
			case itemText:
				m[key] = val
				state = stateKey
			case itemError:
				return nil, fmt.Errorf("%s: %s", key, val)
			default:
				return nil, fmt.Errorf("%s: missing value", key)
			}
		}
	}

	//fmt.Printf("parse: state: %q rest: '%s'\n", state, string(s))
	switch state {
	case stateEqual:
		return nil, fmt.Errorf("%s: missing '='", key)
	case stateVal:
		return nil, fmt.Errorf("%s: missing value", key)
	}
	if len(m) == 0 {
		return nil, nil
	}
	return m, nil
}

type itemType string

const (
	itemText  itemType = "TEXT"
	itemEqual          = "EQUAL"
	itemError          = "ERROR"
)

func (t itemType) String() string {
	return string(t)
}

type state string

const (

	// lexer states
	stateStart    state = "start"
	stateEqual          = "equal"
	stateText           = "text"
	stateQText          = "qtext"
	stateQTextEnd       = "qtextend"
	stateQTextEsc       = "qtextesc"

	// parser states
	stateKey = "key"
	stateVal = "val"
)

func lex(s []rune) (itemType, string, int) {
	isEqual := func(r rune) bool { return r == '=' }
	isEscape := func(r rune) bool { return r == '\\' }
	isQuote := func(r rune) bool { return r == '"' }
	isSpace := func(r rune) bool { return r == ' ' }

	unquote := func(r []rune) (string, error) {
		v := strings.TrimSpace(string(r))
		return strconv.Unquote(v)
	}

	var quote rune
	state := stateStart
	for i, r := range s {
		// fmt.Println("lex:", "i:", i, "r:", string(r), "state:", string(state), "head:", string(s[:i]), "tail:", string(s[i:]))
		switch state {
		case stateStart:
			switch {
			case isSpace(r):
				// state = stateStart
			case isEqual(r):
				state = stateEqual
			case isQuote(r):
				quote = r
				state = stateQText
			default:
				state = stateText
			}

		case stateEqual:
			return itemEqual, "", i

		case stateText:
			switch {
			case isEqual(r) || isSpace(r):
				v := strings.TrimSpace(string(s[:i]))
				return itemText, v, i
			default:
				// state = stateText
			}

		case stateQText:
			switch {
			case r == quote:
				state = stateQTextEnd
			case isEscape(r):
				state = stateQTextEsc
			default:
				// state = stateQText
			}

		case stateQTextEsc:
			state = stateQText

		case stateQTextEnd:
			v, err := unquote(s[:i])
			if err != nil {
				return itemError, err.Error(), i
			}
			return itemText, v, i
		}
	}

	// fmt.Println("lex:", "state:", string(state))
	switch state {
	case stateEqual:
		return itemEqual, "", len(s)
	case stateQText:
		return itemError, "unbalanced quotes", len(s)
	case stateQTextEsc:
		return itemError, "unterminated escape sequence", len(s)
	case stateQTextEnd:
		v, err := unquote(s)
		if err != nil {
			return itemError, err.Error(), len(s)
		}
		return itemText, v, len(s)
	default:
		return itemText, strings.TrimSpace(string(s)), len(s)
	}
}
