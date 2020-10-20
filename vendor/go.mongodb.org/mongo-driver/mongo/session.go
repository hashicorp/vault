// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// ErrWrongClient is returned when a user attempts to pass in a session created by a different client than
// the method call is using.
var ErrWrongClient = errors.New("session was not created by this client")

var withTransactionTimeout = 120 * time.Second

// SessionContext combines the context.Context and mongo.Session interfaces. It should be used as the Context arguments
// to operations that should be executed in a session. This type is not goroutine safe and must not be used concurrently
// by multiple goroutines.
//
// There are two ways to create a SessionContext and use it in a session/transaction. The first is to use one of the
// callback-based functions such as WithSession and UseSession. These functions create a SessionContext and pass it to
// the provided callback. The other is to use NewSessionContext to explicitly create a SessionContext.
type SessionContext interface {
	context.Context
	Session
}

type sessionContext struct {
	context.Context
	Session
}

type sessionKey struct {
}

// NewSessionContext creates a new SessionContext associated with the given Context and Session parameters.
func NewSessionContext(ctx context.Context, sess Session) SessionContext {
	return &sessionContext{
		Context: context.WithValue(ctx, sessionKey{}, sess),
		Session: sess,
	}
}

// SessionFromContext extracts the mongo.Session object stored in a Context. This can be used on a SessionContext that
// was created implicitly through one of the callback-based session APIs or explicitly by calling NewSessionContext. If
// there is no Session stored in the provided Context, nil is returned.
func SessionFromContext(ctx context.Context) Session {
	val := ctx.Value(sessionKey{})
	if val == nil {
		return nil
	}

	sess, ok := val.(Session)
	if !ok {
		return nil
	}

	return sess
}

// Session is an interface that represents a MongoDB logical session. Sessions can be used to enable causal consistency
// for a group of operations or to execute operations in an ACID transaction. A new Session can be created from a Client
// instance. A Session created from a Client must only be used to execute operations using that Client or a Database or
// Collection created from that Client. Custom implementations of this interface should not be used in production. For
// more information about sessions, and their use cases, see https://docs.mongodb.com/manual/reference/server-sessions/,
// https://docs.mongodb.com/manual/core/read-isolation-consistency-recency/#causal-consistency, and
// https://docs.mongodb.com/manual/core/transactions/.
//
// StartTransaction starts a new transaction, configured with the given options, on this session. This method will
// return an error if there is already a transaction in-progress for this session.
//
// CommitTransaction commits the active transaction for this session. This method will return an error if there is no
// active transaction for this session or the transaction has been aborted.
//
// AbortTransaction aborts the active transaction for this session. This method will return an error if there is no
// active transaction for this session or the transaction has been committed or aborted.
//
// WithTransaction starts a transaction on this session and runs the fn callback. Errors with the
// TransientTransactionError and UnknownTransactionCommitResult labels are retried for up to 120 seconds. Inside the
// callback, sessCtx must be used as the Context parameter for any operations that should be part of the transaction. If
// the ctx parameter already has a Session attached to it, it will be replaced by this session. The fn callback may be
// run multiple times during WithTransaction due to retry attempts, so it must be idempotent. Non-retryable operation
// errors or any operation errors that occur after the timeout expires will be returned without retrying. For a usage
// example, see the Client.StartSession method documentation.
//
// ClusterTime, OperationTime, Client, and ID return the session's current operation time, the session's current cluster
// time, the Client associated with the session, and the ID document associated with the session, respectively. The ID
// document for a session is in the form {"id": <BSON binary value>}.
//
// EndSession method should abort any existing transactions and close the session.
//
// AdvanceClusterTime and AdvanceOperationTime are for internal use only and must not be called.
type Session interface {
	// Functions to modify session state.
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
	WithTransaction(ctx context.Context, fn func(sessCtx SessionContext) (interface{}, error),
		opts ...*options.TransactionOptions) (interface{}, error)
	EndSession(context.Context)

	// Functions to retrieve session properties.
	ClusterTime() bson.Raw
	OperationTime() *primitive.Timestamp
	Client() *Client
	ID() bson.Raw

	// Functions to modify mutable session properties.
	AdvanceClusterTime(bson.Raw) error
	AdvanceOperationTime(*primitive.Timestamp) error

	session()
}

// XSession is an unstable interface for internal use only.
//
// Deprecated: This interface is unstable because it provides access to a session.Client object, which exists in the
// "x" package. It should not be used by applications and may be changed or removed in any release.
type XSession interface {
	ClientSession() *session.Client
}

// sessionImpl represents a set of sequential operations executed by an application that are related in some way.
type sessionImpl struct {
	clientSession       *session.Client
	client              *Client
	deployment          driver.Deployment
	didCommitAfterStart bool // true if commit was called after start with no other operations
}

var _ Session = &sessionImpl{}
var _ XSession = &sessionImpl{}

// ClientSession implements the XSession interface.
func (s *sessionImpl) ClientSession() *session.Client {
	return s.clientSession
}

// ID implements the Session interface.
func (s *sessionImpl) ID() bson.Raw {
	return bson.Raw(s.clientSession.SessionID)
}

// EndSession implements the Session interface.
func (s *sessionImpl) EndSession(ctx context.Context) {
	if s.clientSession.TransactionInProgress() {
		// ignore all errors aborting during an end session
		_ = s.AbortTransaction(ctx)
	}
	s.clientSession.EndSession()
}

// WithTransaction implements the Session interface.
func (s *sessionImpl) WithTransaction(ctx context.Context, fn func(sessCtx SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (interface{}, error) {
	timeout := time.NewTimer(withTransactionTimeout)
	defer timeout.Stop()
	var err error
	for {
		err = s.StartTransaction(opts...)
		if err != nil {
			return nil, err
		}

		res, err := fn(NewSessionContext(ctx, s))
		if err != nil {
			if s.clientSession.TransactionRunning() {
				_ = s.AbortTransaction(ctx)
			}

			select {
			case <-timeout.C:
				return nil, err
			default:
			}

			if cerr, ok := err.(CommandError); ok {
				if cerr.HasErrorLabel(driver.TransientTransactionError) {
					continue
				}
			}
			return res, err
		}

		err = s.clientSession.CheckAbortTransaction()
		if err != nil {
			return res, nil
		}

	CommitLoop:
		for {
			err = s.CommitTransaction(ctx)
			if err == nil {
				return res, nil
			}

			select {
			case <-timeout.C:
				return res, err
			default:
			}

			if cerr, ok := err.(CommandError); ok {
				if cerr.HasErrorLabel(driver.UnknownTransactionCommitResult) && !cerr.IsMaxTimeMSExpiredError() {
					continue
				}
				if cerr.HasErrorLabel(driver.TransientTransactionError) {
					break CommitLoop
				}
			}
			return res, err
		}
	}
}

// StartTransaction implements the Session interface.
func (s *sessionImpl) StartTransaction(opts ...*options.TransactionOptions) error {
	err := s.clientSession.CheckStartTransaction()
	if err != nil {
		return err
	}

	s.didCommitAfterStart = false

	topts := options.MergeTransactionOptions(opts...)
	coreOpts := &session.TransactionOptions{
		ReadConcern:    topts.ReadConcern,
		ReadPreference: topts.ReadPreference,
		WriteConcern:   topts.WriteConcern,
		MaxCommitTime:  topts.MaxCommitTime,
	}

	return s.clientSession.StartTransaction(coreOpts)
}

// AbortTransaction implements the Session interface.
func (s *sessionImpl) AbortTransaction(ctx context.Context) error {
	err := s.clientSession.CheckAbortTransaction()
	if err != nil {
		return err
	}

	// Do not run the abort command if the transaction is in starting state
	if s.clientSession.TransactionStarting() || s.didCommitAfterStart {
		return s.clientSession.AbortTransaction()
	}

	selector := makePinnedSelector(s.clientSession, description.WriteSelector())

	s.clientSession.Aborting = true
	_ = operation.NewAbortTransaction().Session(s.clientSession).ClusterClock(s.client.clock).Database("admin").
		Deployment(s.deployment).WriteConcern(s.clientSession.CurrentWc).ServerSelector(selector).
		Retry(driver.RetryOncePerCommand).CommandMonitor(s.client.monitor).
		RecoveryToken(bsoncore.Document(s.clientSession.RecoveryToken)).Execute(ctx)

	s.clientSession.Aborting = false
	_ = s.clientSession.AbortTransaction()

	return nil
}

// CommitTransaction implements the Session interface.
func (s *sessionImpl) CommitTransaction(ctx context.Context) error {
	err := s.clientSession.CheckCommitTransaction()
	if err != nil {
		return err
	}

	// Do not run the commit command if the transaction is in started state
	if s.clientSession.TransactionStarting() || s.didCommitAfterStart {
		s.didCommitAfterStart = true
		return s.clientSession.CommitTransaction()
	}

	if s.clientSession.TransactionCommitted() {
		s.clientSession.RetryingCommit = true
	}

	selector := makePinnedSelector(s.clientSession, description.WriteSelector())

	s.clientSession.Committing = true
	op := operation.NewCommitTransaction().
		Session(s.clientSession).ClusterClock(s.client.clock).Database("admin").Deployment(s.deployment).
		WriteConcern(s.clientSession.CurrentWc).ServerSelector(selector).Retry(driver.RetryOncePerCommand).
		CommandMonitor(s.client.monitor).RecoveryToken(bsoncore.Document(s.clientSession.RecoveryToken))
	if s.clientSession.CurrentMct != nil {
		op.MaxTimeMS(int64(*s.clientSession.CurrentMct / time.Millisecond))
	}

	err = op.Execute(ctx)
	s.clientSession.Committing = false
	commitErr := s.clientSession.CommitTransaction()

	// We set the write concern to majority for subsequent calls to CommitTransaction.
	s.clientSession.UpdateCommitTransactionWriteConcern()

	if err != nil {
		return replaceErrors(err)
	}
	return commitErr
}

// ClusterTime implements the Session interface.
func (s *sessionImpl) ClusterTime() bson.Raw {
	return s.clientSession.ClusterTime
}

// AdvanceClusterTime implements the Session interface.
func (s *sessionImpl) AdvanceClusterTime(d bson.Raw) error {
	return s.clientSession.AdvanceClusterTime(d)
}

// OperationTime implements the Session interface.
func (s *sessionImpl) OperationTime() *primitive.Timestamp {
	return s.clientSession.OperationTime
}

// AdvanceOperationTime implements the Session interface.
func (s *sessionImpl) AdvanceOperationTime(ts *primitive.Timestamp) error {
	return s.clientSession.AdvanceOperationTime(ts)
}

// Client implements the Session interface.
func (s *sessionImpl) Client() *Client {
	return s.client
}

// session implements the Session interface.
func (*sessionImpl) session() {
}

// sessionFromContext checks for a sessionImpl in the argued context and returns the session if it
// exists
func sessionFromContext(ctx context.Context) *session.Client {
	s := ctx.Value(sessionKey{})
	if ses, ok := s.(*sessionImpl); ses != nil && ok {
		return ses.clientSession
	}

	return nil
}
