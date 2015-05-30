package gore

import (
	"time"
)

// Pipeline keeps a list of command for sending to redis once, saving network roundtrip
type Pipeline struct {
	commands []*Command
}

// NewPipeline returns new Pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		commands: []*Command{},
	}
}

// Add appends new commands to the pipeline
func (p *Pipeline) Add(cmd ...*Command) {
	p.commands = append(p.commands, cmd...)
}

// Reset clears all command in the pipeline
func (p *Pipeline) Reset() {
	p.commands = []*Command{}
}

// Run sends the pipeline and returns a slice of Reply
func (p *Pipeline) Run(conn *Conn) (r []*Reply, err error) {
	if len(p.commands) == 0 {
		return nil, nil
	}
	if conn.state != connStateConnected {
		return nil, ErrNotConnected
	}
	conn.Lock()
	defer func() {
                conn.Unlock()
                if err != nil {
                        conn.fail()
                }
	}()
	if conn.RequestTimeout != 0 {
		conn.tcpConn.SetWriteDeadline(time.Now().Add(conn.RequestTimeout * time.Duration(len(p.commands)/10+1)))
	}
	for _, cmd := range p.commands {
		err = cmd.writeCommand(conn)
		if err != nil {
			return nil, ErrWrite
		}
	}
	err = conn.wb.Flush()
	if err != nil {
		return nil, ErrWrite
	}
	if conn.RequestTimeout != 0 {
		conn.tcpConn.SetReadDeadline(time.Now().Add(conn.RequestTimeout * time.Duration(len(p.commands)/10+1)))
	}
	replies := make([]*Reply, len(p.commands))
	for i := range replies {
		rep, err := readReply(conn)
		if err != nil {
			return nil, err
		}
		replies[i] = rep
	}
	return replies, nil
}
