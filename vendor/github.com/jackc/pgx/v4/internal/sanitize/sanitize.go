package sanitize

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Part is either a string or an int. A string is raw SQL. An int is a
// argument placeholder.
type Part interface{}

type Query struct {
	Parts []Part
}

// utf.DecodeRune returns the utf8.RuneError for errors. But that is actually rune U+FFFD -- the unicode replacement
// character. utf8.RuneError is not an error if it is also width 3.
//
// https://github.com/jackc/pgx/issues/1380
const replacementcharacterwidth = 3

func (q *Query) Sanitize(args ...interface{}) (string, error) {
	argUse := make([]bool, len(args))
	buf := &bytes.Buffer{}

	for _, part := range q.Parts {
		var str string
		switch part := part.(type) {
		case string:
			str = part
		case int:
			argIdx := part - 1
			if argIdx >= len(args) {
				return "", fmt.Errorf("insufficient arguments")
			}
			arg := args[argIdx]
			switch arg := arg.(type) {
			case nil:
				str = "null"
			case int64:
				str = strconv.FormatInt(arg, 10)
			case float64:
				str = strconv.FormatFloat(arg, 'f', -1, 64)
			case bool:
				str = strconv.FormatBool(arg)
			case []byte:
				str = QuoteBytes(arg)
			case string:
				str = QuoteString(arg)
			case time.Time:
				str = arg.Truncate(time.Microsecond).Format("'2006-01-02 15:04:05.999999999Z07:00:00'")
			default:
				return "", fmt.Errorf("invalid arg type: %T", arg)
			}
			argUse[argIdx] = true

			// Prevent SQL injection via Line Comment Creation
			// https://github.com/jackc/pgx/security/advisories/GHSA-m7wr-2xf7-cm9p
			str = " " + str + " "
		default:
			return "", fmt.Errorf("invalid Part type: %T", part)
		}
		buf.WriteString(str)
	}

	for i, used := range argUse {
		if !used {
			return "", fmt.Errorf("unused argument: %d", i)
		}
	}
	return buf.String(), nil
}

func NewQuery(sql string) (*Query, error) {
	l := &sqlLexer{
		src:     sql,
		stateFn: rawState,
	}

	for l.stateFn != nil {
		l.stateFn = l.stateFn(l)
	}

	query := &Query{Parts: l.parts}

	return query, nil
}

func QuoteString(str string) string {
	return "'" + strings.ReplaceAll(str, "'", "''") + "'"
}

func QuoteBytes(buf []byte) string {
	return `'\x` + hex.EncodeToString(buf) + "'"
}

type sqlLexer struct {
	src     string
	start   int
	pos     int
	nested  int // multiline comment nesting level.
	stateFn stateFn
	parts   []Part
}

type stateFn func(*sqlLexer) stateFn

func rawState(l *sqlLexer) stateFn {
	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		switch r {
		case 'e', 'E':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune == '\'' {
				l.pos += width
				return escapeStringState
			}
		case '\'':
			return singleQuoteState
		case '"':
			return doubleQuoteState
		case '$':
			nextRune, _ := utf8.DecodeRuneInString(l.src[l.pos:])
			if '0' <= nextRune && nextRune <= '9' {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos-width])
				}
				l.start = l.pos
				return placeholderState
			}
		case '-':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune == '-' {
				l.pos += width
				return oneLineCommentState
			}
		case '/':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune == '*' {
				l.pos += width
				return multilineCommentState
			}
		case utf8.RuneError:
			if width != replacementcharacterwidth {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos])
					l.start = l.pos
				}
				return nil
			}
		}
	}
}

func singleQuoteState(l *sqlLexer) stateFn {
	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		switch r {
		case '\'':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune != '\'' {
				return rawState
			}
			l.pos += width
		case utf8.RuneError:
			if width != replacementcharacterwidth {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos])
					l.start = l.pos
				}
				return nil
			}
		}
	}
}

func doubleQuoteState(l *sqlLexer) stateFn {
	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		switch r {
		case '"':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune != '"' {
				return rawState
			}
			l.pos += width
		case utf8.RuneError:
			if width != replacementcharacterwidth {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos])
					l.start = l.pos
				}
				return nil
			}
		}
	}
}

// placeholderState consumes a placeholder value. The $ must have already has
// already been consumed. The first rune must be a digit.
func placeholderState(l *sqlLexer) stateFn {
	num := 0

	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		if '0' <= r && r <= '9' {
			num *= 10
			num += int(r - '0')
		} else {
			l.parts = append(l.parts, num)
			l.pos -= width
			l.start = l.pos
			return rawState
		}
	}
}

func escapeStringState(l *sqlLexer) stateFn {
	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		switch r {
		case '\\':
			_, width = utf8.DecodeRuneInString(l.src[l.pos:])
			l.pos += width
		case '\'':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune != '\'' {
				return rawState
			}
			l.pos += width
		case utf8.RuneError:
			if width != replacementcharacterwidth {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos])
					l.start = l.pos
				}
				return nil
			}
		}
	}
}

func oneLineCommentState(l *sqlLexer) stateFn {
	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		switch r {
		case '\\':
			_, width = utf8.DecodeRuneInString(l.src[l.pos:])
			l.pos += width
		case '\n', '\r':
			return rawState
		case utf8.RuneError:
			if width != replacementcharacterwidth {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos])
					l.start = l.pos
				}
				return nil
			}
		}
	}
}

func multilineCommentState(l *sqlLexer) stateFn {
	for {
		r, width := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += width

		switch r {
		case '/':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune == '*' {
				l.pos += width
				l.nested++
			}
		case '*':
			nextRune, width := utf8.DecodeRuneInString(l.src[l.pos:])
			if nextRune != '/' {
				continue
			}

			l.pos += width
			if l.nested == 0 {
				return rawState
			}
			l.nested--

		case utf8.RuneError:
			if width != replacementcharacterwidth {
				if l.pos-l.start > 0 {
					l.parts = append(l.parts, l.src[l.start:l.pos])
					l.start = l.pos
				}
				return nil
			}
		}
	}
}

// SanitizeSQL replaces placeholder values with args. It quotes and escapes args
// as necessary. This function is only safe when standard_conforming_strings is
// on.
func SanitizeSQL(sql string, args ...interface{}) (string, error) {
	query, err := NewQuery(sql)
	if err != nil {
		return "", err
	}
	return query.Sanitize(args...)
}
