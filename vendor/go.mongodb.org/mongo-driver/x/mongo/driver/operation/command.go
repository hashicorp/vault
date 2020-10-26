// NOTE: This file is maintained by hand because operationgen cannot generate it.

package operation

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Command is used to run a generic operation.
type Command struct {
	command        bsoncore.Document
	readConcern    *readconcern.ReadConcern
	database       string
	deployment     driver.Deployment
	selector       description.ServerSelector
	readPreference *readpref.ReadPref
	clock          *session.ClusterClock
	session        *session.Client
	monitor        *event.CommandMonitor
	result         bsoncore.Document
	srvr           driver.Server
	desc           description.Server
	crypt          *driver.Crypt
}

// NewCommand constructs and returns a new Command.
func NewCommand(command bsoncore.Document) *Command { return &Command{command: command} }

// Result returns the result of executing this operation.
func (c *Command) Result() bsoncore.Document { return c.result }

// ResultCursor parses the command response as a cursor and returns the resulting BatchCursor.
func (c *Command) ResultCursor(opts driver.CursorOptions) (*driver.BatchCursor, error) {
	cursorRes, err := driver.NewCursorResponse(c.result, c.srvr, c.desc)
	if err != nil {
		return nil, err
	}

	return driver.NewBatchCursor(cursorRes, c.session, c.clock, opts)
}

// Execute runs this operations and returns an error if the operaiton did not execute successfully.
func (c *Command) Execute(ctx context.Context) error {
	if c.deployment == nil {
		return errors.New("the Command operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn: func(dst []byte, desc description.SelectedServer) ([]byte, error) {
			return append(dst, c.command[4:len(c.command)-1]...), nil
		},
		ProcessResponseFn: func(resp bsoncore.Document, srvr driver.Server, desc description.Server, currIndex int) error {
			c.result = resp
			c.srvr = srvr
			c.desc = desc
			return nil
		},
		Client:         c.session,
		Clock:          c.clock,
		CommandMonitor: c.monitor,
		Database:       c.database,
		Deployment:     c.deployment,
		ReadPreference: c.readPreference,
		Selector:       c.selector,
		Crypt:          c.crypt,
	}.Execute(ctx, nil)
}

// Command sets the command to be run.
func (c *Command) Command(command bsoncore.Document) *Command {
	if c == nil {
		c = new(Command)
	}

	c.command = command
	return c
}

// Session sets the session for this operation.
func (c *Command) Session(session *session.Client) *Command {
	if c == nil {
		c = new(Command)
	}

	c.session = session
	return c
}

// ClusterClock sets the cluster clock for this operation.
func (c *Command) ClusterClock(clock *session.ClusterClock) *Command {
	if c == nil {
		c = new(Command)
	}

	c.clock = clock
	return c
}

// CommandMonitor sets the monitor to use for APM events.
func (c *Command) CommandMonitor(monitor *event.CommandMonitor) *Command {
	if c == nil {
		c = new(Command)
	}

	c.monitor = monitor
	return c
}

// Database sets the database to run this operation against.
func (c *Command) Database(database string) *Command {
	if c == nil {
		c = new(Command)
	}

	c.database = database
	return c
}

// Deployment sets the deployment to use for this operation.
func (c *Command) Deployment(deployment driver.Deployment) *Command {
	if c == nil {
		c = new(Command)
	}

	c.deployment = deployment
	return c
}

// ReadConcern specifies the read concern for this operation.
func (c *Command) ReadConcern(readConcern *readconcern.ReadConcern) *Command {
	if c == nil {
		c = new(Command)
	}

	c.readConcern = readConcern
	return c
}

// ReadPreference set the read prefernce used with this operation.
func (c *Command) ReadPreference(readPreference *readpref.ReadPref) *Command {
	if c == nil {
		c = new(Command)
	}

	c.readPreference = readPreference
	return c
}

// ServerSelector sets the selector used to retrieve a server.
func (c *Command) ServerSelector(selector description.ServerSelector) *Command {
	if c == nil {
		c = new(Command)
	}

	c.selector = selector
	return c
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (c *Command) Crypt(crypt *driver.Crypt) *Command {
	if c == nil {
		c = new(Command)
	}

	c.crypt = crypt
	return c
}
