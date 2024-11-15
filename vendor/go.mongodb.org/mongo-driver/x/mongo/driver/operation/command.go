// Copyright (C) MongoDB, Inc. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Command is used to run a generic operation.
type Command struct {
	authenticator  driver.Authenticator
	command        bsoncore.Document
	database       string
	deployment     driver.Deployment
	selector       description.ServerSelector
	readPreference *readpref.ReadPref
	clock          *session.ClusterClock
	session        *session.Client
	monitor        *event.CommandMonitor
	resultResponse bsoncore.Document
	resultCursor   *driver.BatchCursor
	crypt          driver.Crypt
	serverAPI      *driver.ServerAPIOptions
	createCursor   bool
	cursorOpts     driver.CursorOptions
	timeout        *time.Duration
	logger         *logger.Logger
}

// NewCommand constructs and returns a new Command. Once the operation is executed, the result may only be accessed via
// the Result() function.
func NewCommand(command bsoncore.Document) *Command {
	return &Command{
		command: command,
	}
}

// NewCursorCommand constructs a new Command. Once the operation is executed, the server response will be used to
// construct a cursor, which can be accessed via the ResultCursor() function.
func NewCursorCommand(command bsoncore.Document, cursorOpts driver.CursorOptions) *Command {
	return &Command{
		command:      command,
		cursorOpts:   cursorOpts,
		createCursor: true,
	}
}

// Result returns the result of executing this operation.
func (c *Command) Result() bsoncore.Document { return c.resultResponse }

// ResultCursor returns the BatchCursor that was constructed using the command response. If the operation was not
// configured to create a cursor (i.e. it was created using NewCommand rather than NewCursorCommand), this function
// will return nil and an error.
func (c *Command) ResultCursor() (*driver.BatchCursor, error) {
	if !c.createCursor {
		return nil, errors.New("command operation was not configured to create a cursor, but a result cursor was requested")
	}
	return c.resultCursor, nil
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (c *Command) Execute(ctx context.Context) error {
	if c.deployment == nil {
		return errors.New("the Command operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn: func(dst []byte, _ description.SelectedServer) ([]byte, error) {
			return append(dst, c.command[4:len(c.command)-1]...), nil
		},
		ProcessResponseFn: func(info driver.ResponseInfo) error {
			c.resultResponse = info.ServerResponse

			if c.createCursor {
				cursorRes, err := driver.NewCursorResponse(info)
				if err != nil {
					return err
				}

				c.resultCursor, err = driver.NewBatchCursor(cursorRes, c.session, c.clock, c.cursorOpts)
				return err
			}

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
		ServerAPI:      c.serverAPI,
		Timeout:        c.timeout,
		Logger:         c.logger,
		Authenticator:  c.authenticator,
	}.Execute(ctx)
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

// ReadPreference set the read preference used with this operation.
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
func (c *Command) Crypt(crypt driver.Crypt) *Command {
	if c == nil {
		c = new(Command)
	}

	c.crypt = crypt
	return c
}

// ServerAPI sets the server API version for this operation.
func (c *Command) ServerAPI(serverAPI *driver.ServerAPIOptions) *Command {
	if c == nil {
		c = new(Command)
	}

	c.serverAPI = serverAPI
	return c
}

// Timeout sets the timeout for this operation.
func (c *Command) Timeout(timeout *time.Duration) *Command {
	if c == nil {
		c = new(Command)
	}

	c.timeout = timeout
	return c
}

// Logger sets the logger for this operation.
func (c *Command) Logger(logger *logger.Logger) *Command {
	if c == nil {
		c = new(Command)
	}

	c.logger = logger
	return c
}

// Authenticator sets the authenticator to use for this operation.
func (c *Command) Authenticator(authenticator driver.Authenticator) *Command {
	if c == nil {
		c = new(Command)
	}

	c.authenticator = authenticator
	return c
}
