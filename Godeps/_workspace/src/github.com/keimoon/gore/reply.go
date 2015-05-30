package gore

import (
	"io"
	"strconv"
	"time"
)

// Reply type, similar to Hiredis
const (
	ReplyString  = 1
	ReplyArray   = 2
	ReplyInteger = 3
	ReplyNil     = 4
	ReplyStatus  = 5
	ReplyError   = 6
)

// Pair holds a pair of value, such as reply from HGETALL or ZRANGE [WITHSCORES]
type Pair struct {
	First  []byte
	Second []byte
}

// Reply holds redis reply
type Reply struct {
	replyType    int
	integerValue int64
	stringValue  []byte
	arrayValue   []*Reply
}

// Type returns reply type
func (r *Reply) Type() int {
	return r.replyType
}

// String returns string value of a reply
func (r *Reply) String() (string, error) {
	if r.Type() == ReplyNil {
		return "", ErrNil
	}
	if r.Type() != ReplyString && r.Type() != ReplyStatus {
		return "", ErrType
	}
	return string(r.stringValue), nil
}

// Bytes returns string value of a reply as byte array
func (r *Reply) Bytes() ([]byte, error) {
	if r.Type() == ReplyNil {
		return nil, ErrNil
	}
	if r.Type() != ReplyString && r.Type() != ReplyStatus {
		return nil, ErrType
	}
	return r.stringValue, nil
}

// Int returns integer value (for example from INCR) or convert string value to integer if possible
func (r *Reply) Int() (int64, error) {
	if r.Type() == ReplyInteger {
		return r.integerValue, nil
	}
	s, err := r.String()
	if err != nil {
		return 0, err
	}
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, ErrConvert
	}
	return x, nil
}

// Float parses string value to float64
func (r *Reply) Float() (float64, error) {
	s, err := r.String()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

// Bool parses string value to boolean.
// Integer 0 and string "0" returns false while integer 1 and string "1" returns true.
// Other values will cause ErrConvert
func (r *Reply) Bool() (bool, error) {
	switch r.Type() {
	case ReplyString:
		s, _ := r.String()
		if s == "1" || s == "true" {
			return true, nil
		} else if s == "0" || s == "false" {
			return false, nil
		} else {
			return false, ErrConvert
		}
	case ReplyInteger:
		i, _ := r.Int()
		if i == 1 {
			return true, nil
		} else if i == 0 {
			return false, nil
		} else {
			return false, ErrConvert
		}
	default:
		return false, ErrConvert
	}
}

// Integer returns integer value (for example from INCR) or convert string value to integer if possible
func (r *Reply) Integer() (int64, error) {
	return r.Int()
}

// Array returns array value of a reply
func (r *Reply) Array() ([]*Reply, error) {
	if r.Type() == ReplyNil {
		return nil, ErrNil
	}
	if r.Type() != ReplyArray {
		return nil, ErrType
	}
	return r.arrayValue, nil
}

// FixInt returns a fixed size int64
func (r *Reply) FixInt() (int64, error) {
	if r.Type() == ReplyNil {
		return 0, ErrNil
	}
	if r.Type() != ReplyString {
		return 0, ErrType
	}
	return ToFixInt(r.stringValue)
}

// VarInt returns a base-128 encoded int64
func (r *Reply) VarInt() (int64, error) {
	if r.Type() == ReplyNil {
		return 0, ErrNil
	}
	if r.Type() != ReplyString {
		return 0, ErrType
	}
	return ToVarInt(r.stringValue)
}

// Slice parses the reply to a slice. The element of the destination slice
// must be integer, float, boolean, string, []byte, FixInt, VarInt, or a Pair
func (r *Reply) Slice(s interface{}) error {
	if r.Type() == ReplyNil {
		return ErrNil
	}
	if r.Type() != ReplyArray {
		return ErrType
	}
	if len(r.arrayValue) == 0 {
		return nil
	}
	switch s := s.(type) {
	case *[]int:
		*s = make([]int, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].Int()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = int(x)
		}
		return nil
	case *[]int64:
		*s = make([]int64, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].Int()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = x
		}
		return nil
	case *[]float64:
		*s = make([]float64, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].Float()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = x
		}
		return nil
	case *[]bool:
		*s = make([]bool, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].Bool()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = x
		}
		return nil
	case *[]string:
		*s = make([]string, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].String()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = x
		}
		return nil
	case *[][]byte:
		*s = make([][]byte, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].Bytes()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = x
		}
		return nil
	case *[]FixInt:
		*s = make([]FixInt, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].FixInt()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = FixInt(x)
		}
		return nil
	case *[]VarInt:
		*s = make([]VarInt, len(r.arrayValue))
		for i := range r.arrayValue {
			x, err := r.arrayValue[i].VarInt()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = VarInt(x)
		}
		return nil
	case *[]*Pair:
		if len(r.arrayValue)%2 != 0 {
			return ErrType
		}
		*s = make([]*Pair, len(r.arrayValue)/2)
		for i := 0; i < len(r.arrayValue)/2; i++ {
			first, err := r.arrayValue[2*i].Bytes()
			if err != nil && err != ErrNil {
				return err
			}
			second, err := r.arrayValue[2*i+1].Bytes()
			if err != nil && err != ErrNil {
				return err
			}
			(*s)[i] = &Pair{first, second}
		}
		return nil
	default:
		return ErrConvert
	}
	return nil
}

// Map converts the reply into a map[string]string.
// It will return error unless the reply is an array reply from HGETALL,
// or SENTINEL master
func (r *Reply) Map() (map[string]string, error) {
	if r.IsNil() {
		return nil, ErrNil
	}
	if !r.IsArray() {
		return nil, ErrType
	}
	if len(r.arrayValue)%2 != 0 {
		return nil, ErrType
	}
	m := make(map[string]string)
	for i := 0; i < len(r.arrayValue)/2; i++ {
		first, err := r.arrayValue[2*i].String()
		if err != nil {
			continue
		}
		second, err := r.arrayValue[2*i+1].String()
		if err != nil && err != ErrNil {
			continue
		}
		m[first] = second
	}
	return m, nil
}

// Error returns error message
func (r *Reply) Error() (string, error) {
	if r.Type() != ReplyError {
		return "", ErrType
	}
	return string(r.stringValue), nil
}

// IsNil checks if reply is nil or not
func (r *Reply) IsNil() bool {
	return r.Type() == ReplyNil
}

// IsStatus checks if reply is a status or not
func (r *Reply) IsStatus() bool {
	return r.Type() == ReplyStatus
}

// IsOk checks if reply is a status reply with "OK" response
func (r *Reply) IsOk() bool {
	return r == okReply
}

// IsString checks if reply is a string reply or not
func (r *Reply) IsString() bool {
	return r.Type() == ReplyString
}

// IsInteger checks if reply is a integer reply or not
func (r *Reply) IsInteger() bool {
	return r.Type() == ReplyInteger
}

// IsArray checks if reply is an array reply or not
func (r *Reply) IsArray() bool {
	return r.Type() == ReplyArray
}

// IsError checks if reply is an error reply or not
func (r *Reply) IsError() bool {
	return r.Type() == ReplyError
}

// Receive safely read a reply from conn
func Receive(conn *Conn) (r *Reply, err error) {
	conn.Lock()
	defer func() {
		if err != nil {
			conn.state = connStateNotConnected
			conn.Unlock()
			conn.fail()
		} else {
			conn.Unlock()
		}
	}()
	if conn.RequestTimeout != 0 {
		conn.tcpConn.SetReadDeadline(time.Now().Add(conn.RequestTimeout))
	}
	return readReply(conn)
}

// Motivated by redigo. Good job, man
func readReply(conn *Conn) (*Reply, error) {
	line, err := readLine(conn)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, ErrRead
	}
	switch line[0] {
	case '+':
		switch {
		case len(line) == 3 && line[1] == 'O' && line[2] == 'K':
			return okReply, nil
		case len(line) == 5 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G':
			return pongReply, nil
		case len(line) == 7 && line[1] == 'Q' && line[2] == 'U' && line[3] == 'E' &&
			line[4] == 'U' && line[5] == 'E' && line[6] == 'D':
			return queuedReply, nil
		default:
			return &Reply{
				replyType:   ReplyStatus,
				stringValue: line[1:],
			}, nil

		}
	case '-':
		return &Reply{
			replyType:   ReplyError,
			stringValue: line[1:],
		}, nil
	case ':':
		intValue, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return nil, ErrRead
		}
		return &Reply{
			replyType:    ReplyInteger,
			integerValue: intValue,
		}, nil
	case '$':
		l, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return nil, ErrRead
		}
		if l < 0 {
			return &Reply{
				replyType: ReplyNil,
			}, nil
		}
		b := make([]byte, l)
		_, err = io.ReadFull(conn.rb, b)
		if err != nil {
			return nil, ErrRead
		}
		line, err = readLine(conn)
		if err != nil || len(line) != 0 {
			return nil, ErrRead
		}
		return &Reply{
			replyType:   ReplyString,
			stringValue: b,
		}, nil
	case '*':
		l, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return nil, ErrRead
		}
		if l < 0 {
			return &Reply{
				replyType: ReplyNil,
			}, nil
		}
		replyArray := make([]*Reply, l)
		for i := range replyArray {
			replyArray[i], err = readReply(conn)
			if err != nil {
				return nil, err
			}
		}
		return &Reply{
			replyType:  ReplyArray,
			arrayValue: replyArray,
		}, nil
	default:
		return nil, ErrRead
	}
}

func readLine(conn *Conn) ([]byte, error) {
	b, err := conn.rb.ReadSlice('\n')
	if err != nil {
		return nil, ErrRead
	}
	i := len(b) - 2
	if i < 0 || b[i] != '\r' {
		return nil, ErrRead
	}
	return b[:i], nil
}

var (
	okReply = &Reply{
		replyType:   ReplyStatus,
		stringValue: []byte{'O', 'K'},
	}
	pongReply = &Reply{
		replyType:   ReplyStatus,
		stringValue: []byte{'P', 'O', 'N', 'G'},
	}
	queuedReply = &Reply{
		replyType:   ReplyStatus,
		stringValue: []byte{'Q', 'U', 'E', 'U', 'E', 'D'},
	}
)
