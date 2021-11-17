// Package ldif contains an  LDIF parser and marshaller (RFC 2849).
package ldif

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// Entry is one entry in the LDIF
type Entry struct {
	Entry  *ldap.Entry
	Add    *ldap.AddRequest
	Del    *ldap.DelRequest
	Modify *ldap.ModifyRequest
}

// The LDIF struct is used for parsing an LDIF. The Controls
// is used to tell the parser to ignore any controls found
// when parsing (default: false to ignore the controls).
// FoldWidth is used for the line lenght when marshalling.
type LDIF struct {
	Entries    []*Entry
	Version    int
	changeType string
	FoldWidth  int
	Controls   bool
	firstEntry bool
}

// The ParseError holds the error message and the line in the ldif
// where the error occurred.
type ParseError struct {
	Line    int
	Message string
}

// Error implements the error interface
func (e *ParseError) Error() string {
	return fmt.Sprintf("Error in line %d: %s", e.Line, e.Message)
}

var cr byte = '\x0D'
var lf byte = '\x0A'
var sep = string([]byte{cr, lf})
var comment byte = '#'
var space byte = ' '
var spaces = string(space)

// Parse wraps Unmarshal to parse an LDIF from a string
func Parse(str string) (l *LDIF, err error) {
	buf := bytes.NewBuffer([]byte(str))
	l = &LDIF{}
	err = Unmarshal(buf, l)
	return
}

// ParseWithControls wraps Unmarshal to parse an LDIF from
// a string, controls are added to change records
func ParseWithControls(str string) (l *LDIF, err error) {
	buf := bytes.NewBuffer([]byte(str))
	l = &LDIF{Controls: true}
	err = Unmarshal(buf, l)
	return
}

// Unmarshal parses the LDIF from the given io.Reader into the LDIF struct.
// The caller is responsible for closing the io.Reader if that is
// needed.
func Unmarshal(r io.Reader, l *LDIF) (err error) {
	if r == nil {
		return &ParseError{Line: 0, Message: "No reader present"}
	}
	curLine := 0
	l.Version = 0
	l.changeType = ""
	isComment := false

	reader := bufio.NewReader(r)

	var lines []string
	var line, nextLine string
	l.firstEntry = true

	for {
		curLine++
		nextLine, err = reader.ReadString(lf)
		nextLine = strings.TrimRight(nextLine, sep)

		switch err {
		case nil, io.EOF:
			switch len(nextLine) {
			case 0:
				if len(line) == 0 && err == io.EOF {
					return nil
				}
				if len(line) == 0 && len(lines) == 0 {
					continue
				}
				lines = append(lines, line)
				entry, perr := l.parseEntry(lines)
				if perr != nil {
					return &ParseError{Line: curLine, Message: perr.Error()}
				}
				l.Entries = append(l.Entries, entry)
				line = ""
				lines = []string{}
				if err == io.EOF {
					return nil
				}
			default:
				switch nextLine[0] {
				case comment:
					isComment = true
					continue

				case space:
					if isComment {
						continue
					}
					line += nextLine[1:]
					continue

				default:
					isComment = false
					if len(line) != 0 {
						lines = append(lines, line)
					}
					line = nextLine
					continue
				}
			}
		default:
			return &ParseError{Line: curLine, Message: err.Error()}
		}
	}
}

func (l *LDIF) parseEntry(lines []string) (entry *Entry, err error) {
	if len(lines) == 0 {
		return nil, errors.New("empty entry?")
	}

	if l.firstEntry && strings.HasPrefix(lines[0], "version:") {
		l.firstEntry = false
		line := strings.TrimLeft(lines[0][8:], spaces)
		if l.Version, err = strconv.Atoi(line); err != nil {
			return nil, err
		}

		if l.Version != 1 {
			return nil, errors.New("Invalid version spec " + string(line))
		}

		l.Version = 1
		if len(lines) == 1 {
			return nil, nil
		}
		lines = lines[1:]
	}
	l.firstEntry = false

	if len(lines) == 0 {
		return nil, nil
	}

	if !strings.HasPrefix(lines[0], "dn:") {
		return nil, errors.New("missing 'dn:'")
	}
	_, val, err := l.parseLine(lines[0])
	if err != nil {
		return nil, err
	}
	dn := val

	if len(lines) == 1 {
		return nil, errors.New("only a dn: line")
	}
	lines = lines[1:]

	var controls []ldap.Control
	controls, lines, err = l.parseControls(lines)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(lines[0], "changetype:") {
		_, val, err := l.parseLine(lines[0])
		if err != nil {
			return nil, err
		}
		l.changeType = val
		if len(lines) > 1 {
			lines = lines[1:]
		}
	}
	switch l.changeType {
	case "":
		if len(controls) != 0 {
			return nil, errors.New("controls found without changetype")
		}
		attrs, err := l.parseAttrs(lines)
		if err != nil {
			return nil, err
		}
		return &Entry{Entry: ldap.NewEntry(dn, attrs)}, nil

	case "add":
		attrs, err := l.parseAttrs(lines)
		if err != nil {
			return nil, err
		}
		// FIXME: controls for add - see https://github.com/go-ldap/ldap/issues/81
		add := ldap.NewAddRequest(dn, controls)
		for attr, vals := range attrs {
			add.Attribute(attr, vals)
		}
		return &Entry{Add: add}, nil

	case "delete":
		if len(lines) > 1 {
			return nil, errors.New("no attributes allowed for changetype delete")
		}
		return &Entry{Del: ldap.NewDelRequest(dn, controls)}, nil

	case "modify":
		// FIXME: controls for modify - see https://github.com/go-ldap/ldap/issues/81
		mod := ldap.NewModifyRequest(dn, controls)
		var op, attribute string
		var values []string
		if lines[len(lines)-1] != "-" {
			return nil, errors.New("modify request does not close with a single dash")
		}

		for i := 0; i < len(lines); i++ {
			if lines[i] == "-" {
				switch op {
				case "":
					return nil, fmt.Errorf("empty operation")
				case "add":
					mod.Add(attribute, values)
					op = ""
					attribute = ""
					values = nil
				case "replace":
					mod.Replace(attribute, values)
					op = ""
					attribute = ""
					values = nil
				case "delete":
					mod.Delete(attribute, values)
					op = ""
					attribute = ""
					values = nil
				default:
					return nil, fmt.Errorf("invalid operation %s in modify request", op)
				}
				continue
			}
			attr, val, err := l.parseLine(lines[i])
			if err != nil {
				return nil, err
			}
			if op == "" {
				op = attr
				attribute = val
			} else {
				if attr != attribute {
					return nil, fmt.Errorf("invalid attribute %s in %s request for %s", attr, op, attribute)
				}
				values = append(values, val)
			}
		}
		return &Entry{Modify: mod}, nil

	case "moddn", "modrdn":
		return nil, fmt.Errorf("unsupported changetype %s", l.changeType)

	default:
		return nil, fmt.Errorf("invalid changetype %s", l.changeType)
	}
}

func (l *LDIF) parseAttrs(lines []string) (map[string][]string, error) {
	attrs := make(map[string][]string)
	for i := 0; i < len(lines); i++ {
		attr, val, err := l.parseLine(lines[i])
		if err != nil {
			return nil, err
		}
		attrs[attr] = append(attrs[attr], val)
	}
	return attrs, nil
}

func (l *LDIF) parseLine(line string) (attr, val string, err error) {
	off := 0
	for len(line) > off && line[off] != ':' {
		off++
		if off >= len(line) {
			err = fmt.Errorf("Missing : in line `%s`", line)
			return
		}
	}
	if off == len(line) {
		err = fmt.Errorf("Missing : in the line `%s`", line)
		return
	}

	if off > len(line)-2 {
		err = errors.New("empty value")
		// FIXME: this is allowed for some attributes, e.g. seeAlso
		return
	}

	attr = line[0:off]
	if err = validAttr(attr); err != nil {
		attr = ""
		val = ""
		return
	}

	switch line[off+1] {
	case ':':
		val, err = decodeBase64(strings.TrimLeft(line[off+2:], spaces))
		if err != nil {
			return
		}

	case '<':
		val, err = readURLValue(strings.TrimLeft(line[off+2:], spaces))
		if err != nil {
			return
		}

	default:
		val = strings.TrimLeft(line[off+1:], spaces)
	}

	return
}

func (l *LDIF) parseControls(lines []string) ([]ldap.Control, []string, error) {
	var controls []ldap.Control
	for {
		if !strings.HasPrefix(lines[0], "control:") {
			break
		}
		if !l.Controls {
			if len(lines) == 1 {
				return nil, nil, errors.New("only controls found")
			}
			lines = lines[1:]
			continue
		}

		_, val, err := l.parseLine(lines[0])
		if err != nil {
			return nil, nil, err
		}

		var oid, ctrlValue string
		criticality := false

		parts := strings.SplitN(val, " ", 3)
		if err = validOID(parts[0]); err != nil {
			return nil, nil, fmt.Errorf("%s is not a valid oid: %s", oid, err)
		}
		oid = parts[0]

		if len(parts) > 1 {
			switch parts[1] {
			case "true":
				criticality = true
				if len(parts) > 2 {
					parts[1] = parts[2]
					parts = parts[0:2]
				}
			case "false":
				criticality = false
				if len(parts) > 2 {
					parts[1] = parts[2]
					parts = parts[0:2]
				}
			}
		}
		if len(parts) == 2 {
			ctrlValue = parts[1]
		}
		if ctrlValue == "" {
			switch oid {
			case ldap.ControlTypeManageDsaIT:
				controls = append(controls, &ldap.ControlManageDsaIT{Criticality: criticality})
			default:
				return nil, nil, fmt.Errorf("unsupported control found: %s", oid)
			}
		} else {
			switch ctrlValue[0] { // where is this documented?
			case ':':
				if len(ctrlValue) == 1 {
					return nil, nil, errors.New("missing value for base64 encoded control value")
				}
				ctrlValue, err = decodeBase64(strings.TrimLeft(ctrlValue[1:], spaces))
				if err != nil {
					return nil, nil, err
				}
				if ctrlValue == "" {
					return nil, nil, errors.New("base64 decoded to empty value")
				}

			case '<':
				if len(ctrlValue) == 1 {
					return nil, nil, errors.New("missing value for url control value")
				}
				ctrlValue, err = readURLValue(strings.TrimLeft(ctrlValue[1:], spaces))
				if err != nil {
					return nil, nil, err
				}
				if ctrlValue == "" {
					return nil, nil, errors.New("url resolved to an empty value")
				}
			}
			// TODO:
			// convert ctrlValue to *ber.Packet and decode with something like
			//   ctrl := ldap.DecodeControl()
			// ... FIXME: the controls need a Decode() interface
			// so we can just do a
			//   ctrl := ldap.ControlByOID(oid) // returns an empty &ControlSomething{}
			//   ctrl.Decode((*ber.Packet)(ctrlValue))
			//   ctrl.Criticality = criticality
			// that should be usable in github.com/go-ldap/ldap/control.go also
			// to decode the incoming control
			// controls = append(controls, ctrl)
			return nil, nil, fmt.Errorf("controls with values are not supported, oid: %s", oid)
		}

		if len(lines) == 1 {
			return nil, nil, errors.New("only controls found")
		}
		lines = lines[1:]
	}
	return controls, lines, nil
}

func readURLValue(val string) (string, error) {
	u, err := url.Parse(val)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %s", err)
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("unsupported URL scheme %s", u.Scheme)
	}
	data, err := ioutil.ReadFile(toPath(u))
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %s", u.Path, err)
	}
	val = string(data) // FIXME: safe?
	return val, nil
}

func decodeBase64(enc string) (string, error) {
	dec := make([]byte, base64.StdEncoding.DecodedLen(len([]byte(enc))))
	n, err := base64.StdEncoding.Decode(dec, []byte(enc))
	if err != nil {
		return "", err
	}
	return string(dec[:n]), nil
}

func validOID(oid string) error {
	lastDot := true
	for _, c := range oid {
		switch {
		case c == '.' && lastDot:
			return errors.New("OID with at least 2 consecutive dots")
		case c == '.':
			lastDot = true
		case c >= '0' && c <= '9':
			lastDot = false
		default:
			return errors.New("Invalid character in OID")
		}
	}
	return nil
}

func validAttr(attr string) error {
	if len(attr) == 0 {
		return errors.New("empty attribute name")
	}
	switch {
	case attr[0] >= 'A' && attr[0] <= 'Z':
		// A-Z
	case attr[0] >= 'a' && attr[0] <= 'z':
		// a-z
	default:
		if attr[0] >= '0' && attr[0] <= '9' {
			return validOID(attr)
		}
		return errors.New("invalid first character in attribute")
	}
	for i := 1; i < len(attr); i++ {
		c := attr[i]
		switch {
		case c >= '0' && c <= '9':
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		case c == '-':
		case c == ';':
		default:
			return errors.New("invalid character in attribute name")
		}
	}
	return nil
}

// AllEntries returns all *ldap.Entries in the LDIF
func (l *LDIF) AllEntries() (entries []*ldap.Entry) {
	for _, entry := range l.Entries {
		if entry.Entry != nil {
			entries = append(entries, entry.Entry)
		}
	}
	return entries
}
