package gore

import (
	"time"
)

// Transaction implements MULTI/EXEC/WATCH protocol of redis.
// Transaction must be used with connection pool, or undefined behavior.
// may orcur when used from multiple goroutine.
type Transaction struct {
	conn     *Conn
	commands []*Command
}

// NewTransaction returns new transaction
func NewTransaction(conn *Conn) *Transaction {
	return &Transaction{
		conn:     conn,
		commands: []*Command{NewCommand("MULTI")},
	}
}

// Watch watches some keys. If the key has been changed before Exec,
// the transaction will be aborted
func (t *Transaction) Watch(key ...string) error {
	if len(key) == 0 {
		return nil
	}
	args := make([]interface{}, len(key))
	for i := range key {
		args[i] = key[i]
	}
	_, err := NewCommand("WATCH", args...).Run(t.conn)
	return err
}

// Add appends commands to the transaction
func (t *Transaction) Add(cmd ...*Command) {
	t.commands = append(t.commands, cmd...)
}

// Commit commits the whole transaction.
// If transaction fail, ErrTransactionAborted is returned.
// If watched key has been modified, ErrKeyChanged is returned.
func (t *Transaction) Commit() ([]*Reply, error) {
	if t.conn.state != connStateConnected {
		return nil, ErrNotConnected
	}
	t.commands = append(t.commands, NewCommand("EXEC"))
	if t.conn.RequestTimeout != 0 {
		t.conn.tcpConn.SetWriteDeadline(time.Now().Add(t.conn.RequestTimeout * time.Duration(len(t.commands)/10+1)))
	}
	for _, cmd := range t.commands {
		err := cmd.writeCommand(t.conn)
		if err != nil {
			t.conn.fail()
			return nil, ErrWrite
		}
	}
	err := t.conn.wb.Flush()
	if err != nil {
		t.conn.fail()
		return nil, ErrWrite
	}
	if t.conn.RequestTimeout != 0 {
		t.conn.tcpConn.SetReadDeadline(time.Now().Add(t.conn.RequestTimeout * time.Duration(len(t.commands)/10+1)))
	}
	replies := make([]*Reply, len(t.commands))
	for i := range replies {
		rep, err := readReply(t.conn)
		if err != nil {
			t.conn.fail()
			return nil, err
		}
		replies[i] = rep
	}
	execReply := replies[len(replies)-1]
	if execReply.IsError() {
		return replies, ErrTransactionAborted
	}
	if execReply.IsNil() {
		return replies, ErrKeyChanged
	}
	if !execReply.IsArray() {
		t.conn.fail()
		return replies, ErrType
	}
	return execReply.Array()
}

// Discard discards the transaction
func (t *Transaction) Discard() error {
	_, err := NewCommand("MULTI").Run(t.conn)
	if err == nil {
		_, err = NewCommand("DISCARD").Run(t.conn)
	}
	return err
}
