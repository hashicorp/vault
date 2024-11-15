package pgtype

import (
	"bytes"
	"fmt"
)

type UntypedTextMultirange struct {
	Elements []string
}

func ParseUntypedTextMultirange(src string) (*UntypedTextMultirange, error) {
	utmr := &UntypedTextMultirange{}
	utmr.Elements = make([]string, 0)

	buf := bytes.NewBufferString(src)

	skipWhitespace(buf)

	r, _, err := buf.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("invalid array: %v", err)
	}

	if r != '{' {
		return nil, fmt.Errorf("invalid multirange, expected '{': %v", err)
	}

parseValueLoop:
	for {
		r, _, err = buf.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("invalid multirange: %v", err)
		}

		switch r {
		case ',': // skip range separator
		case '}':
			break parseValueLoop
		default:
			buf.UnreadRune()
			value, err := parseRange(buf)
			if err != nil {
				return nil, fmt.Errorf("invalid multirange value: %v", err)
			}
			utmr.Elements = append(utmr.Elements, value)
		}
	}

	skipWhitespace(buf)

	if buf.Len() > 0 {
		return nil, fmt.Errorf("unexpected trailing data: %v", buf.String())
	}

	return utmr, nil

}

func parseRange(buf *bytes.Buffer) (string, error) {

	s := &bytes.Buffer{}

	boundSepRead := false
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case ',', '}':
			if r == ',' && !boundSepRead {
				boundSepRead = true
				break
			}
			buf.UnreadRune()
			return s.String(), nil
		}

		s.WriteRune(r)
	}
}
