package gore

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Command sent to redis
type Command struct {
	name string
	args []interface{}
}

// NewCommand returns a new Command
func NewCommand(name string, args ...interface{}) *Command {
	return &Command{
		name: strings.TrimSpace(name),
		args: args,
	}
}

// Run sends command to redis
func (cmd *Command) Run(conn *Conn) (r *Reply, err error) {
	conn.Lock()
	if conn.state != connStateConnected {
		conn.Unlock()
		return nil, ErrNotConnected
	}
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
		conn.tcpConn.SetWriteDeadline(time.Now().Add(conn.RequestTimeout))
	}
	err = cmd.writeCommand(conn)
	if err != nil {
		return nil, ErrWrite
	}
	err = conn.wb.Flush()
	if err != nil {
		return nil, ErrWrite
	}
	if conn.RequestTimeout != 0 {
		conn.tcpConn.SetReadDeadline(time.Now().Add(conn.RequestTimeout))
	}
	return readReply(conn)
}

// Send safely sends a command over conn
func (cmd *Command) Send(conn *Conn) (err error) {
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
		conn.tcpConn.SetWriteDeadline(time.Now().Add(conn.RequestTimeout))
	}
	err = cmd.writeCommand(conn)
	if err != nil {
		return ErrWrite
	}
	err = conn.wb.Flush()
	if err != nil {
		return ErrWrite
	}
	return nil
}

func (cmd *Command) writeCommand(conn *Conn) error {
	cmdLen := strconv.FormatInt(int64(len(cmd.args))+1, 10)
	_, err := conn.wb.WriteString("*" + cmdLen + "\r\n")
	if err != nil {
		return err
	}
	err = writeString(cmd.name, conn)
	if err != nil {
		return err
	}
	for _, arg := range cmd.args {
		err = writeBytes(convertString(arg), conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertString(arg interface{}) []byte {
	switch arg := arg.(type) {
	case string:
		return []byte(arg)
	case []byte:
		return arg
	case int:
		return []byte(strconv.FormatInt(int64(arg), 10))
	case int64:
		return []byte(strconv.FormatInt(arg, 10))
	case float64:
		return []byte(strconv.FormatFloat(arg, 'g', -1, 64))
	case FixInt:
		return arg.Bytes()
	case VarInt:
		return arg.Bytes()
	case bool:
		if arg {
			return []byte("1")
		}
		return []byte("0")
	case nil:
		return []byte("")
	default:
		return []byte(fmt.Sprint(arg))
	}
}

func writeString(s string, conn *Conn) error {
	l := strconv.FormatInt(int64(len(s)), 10)
	_, err := conn.wb.WriteString("$" + l + "\r\n" + s + "\r\n")
	return err
}

func writeBytes(b []byte, conn *Conn) error {
	l := strconv.FormatInt(int64(len(b)), 10)
	conn.wb.WriteString("$" + l + "\r\n")
	conn.wb.Write(b)
	_, err := conn.wb.WriteString("\r\n")
	return err
}
