package pgx

import (
	"io"

	"github.com/jackc/pgx/pgproto3"
)

func (c *Conn) readUntilCopyOutResponse() error {
	for {
		msg, err := c.rxMsg()
		if err != nil {
			return err
		}

		switch msg := msg.(type) {
		case *pgproto3.CopyOutResponse:
			return nil
		default:
			err = c.processContextFreeMsg(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Conn) CopyToWriter(w io.Writer, sql string, args ...interface{}) (CommandTag, error) {
	if err := c.sendSimpleQuery(sql, args...); err != nil {
		return "", err
	}

	if err := c.readUntilCopyOutResponse(); err != nil {
		return "", err
	}

	for {
		msg, err := c.rxMsg()
		if err != nil {
			return "", err
		}

		switch msg := msg.(type) {
		case *pgproto3.CopyDone:
			break
		case *pgproto3.CopyData:
			_, err := w.Write(msg.Data)
			if err != nil {
				c.die(err)
				return "", err
			}
		case *pgproto3.ReadyForQuery:
			c.rxReadyForQuery(msg)
			return "", nil
		case *pgproto3.CommandComplete:
			return CommandTag(msg.CommandTag), nil
		case *pgproto3.ErrorResponse:
			return "", c.rxErrorResponse(msg)
		default:
			return "", c.processContextFreeMsg(msg)
		}
	}
}
