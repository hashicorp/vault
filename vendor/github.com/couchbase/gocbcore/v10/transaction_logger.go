// nolint: unused
package gocbcore

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type loggableDocKey struct {
	bucket     string
	scope      string
	collection string
	key        []byte
}

func newLoggableDocKey(bucket, scope, collection string, key []byte) loggableDocKey {
	return loggableDocKey{
		bucket:     bucket,
		scope:      scope,
		collection: collection,
		key:        key,
	}
}

func (rdi loggableDocKey) String() string {
	if isLogRedactionLevelFull() || isLogRedactionLevelPartial() {
		return redactUserData(rdi.build())
	}

	return rdi.build()
}

func (rdi loggableDocKey) build() string {
	scope := rdi.scope
	if scope == "" {
		scope = "_default"
	}
	collection := rdi.collection
	if collection == "" {
		collection = "_default"
	}
	return rdi.bucket + "." + scope + "." + collection + "." + string(rdi.key)
}

func (rdi loggableDocKey) redacted() interface{} {
	return redactUserData(rdi.build())
}

type loggableATRKey struct {
	bucket     string
	scope      string
	collection string
	key        []byte
}

func newLoggableATRKey(bucket, scope, collection string, key []byte) loggableATRKey {
	return loggableATRKey{
		bucket:     bucket,
		scope:      scope,
		collection: collection,
		key:        key,
	}
}

func (rdi loggableATRKey) String() string {
	if isLogRedactionLevelFull() {
		return redactMetaData(rdi.build())
	}

	return rdi.build()
}

func (rdi loggableATRKey) build() string {
	scope := rdi.scope
	if scope == "" {
		scope = "_default"
	}
	collection := rdi.collection
	if collection == "" {
		collection = "_default"
	}
	data := rdi.bucket + "." + scope + "." + collection
	if len(rdi.key) > 0 {
		data = data + "." + string(rdi.key)
	}

	return data
}

func (rdi loggableATRKey) redacted() interface{} {
	return redactMetaData(rdi.build())
}

// TransactionLogger is the logger used for logging in transactions.
// Uncommitted: This API may change in the future.
type TransactionLogger interface {
	Log(level LogLevel, offset int, txnID, attemptID, format string, v ...interface{}) error
}

// TransactionLogItem represents an entry in the transaction in memory logging.
type TransactionLogItem struct {
	Level LogLevel

	args      []interface{}
	txnID     string
	attemptID string
	timestamp time.Time
	fmt       string
}

func (item TransactionLogItem) String() string {
	return fmt.Sprintf("%s %s/%s %s", item.timestamp.AppendFormat([]byte{}, "15:04:05.000"), item.txnID, item.attemptID, fmt.Sprintf(item.fmt, item.args...))
}

// InMemoryTransactionLogger logs to memory, also logging WARN and ERROR logs to the SDK logger.
// Uncommitted: This API may change in the future.
type InMemoryTransactionLogger struct {
	lock  sync.Mutex
	items []TransactionLogItem
}

// NewInMemoryTransactionLogger returns a new in memory transaction logger.
// Uncommitted: This API may change in the future.
func NewInMemoryTransactionLogger() *InMemoryTransactionLogger {
	return &InMemoryTransactionLogger{
		items: make([]TransactionLogItem, 0, 256),
	}
}

// Logs returns the set of log items created during the transaction.
func (tl *InMemoryTransactionLogger) Logs() []TransactionLogItem {
	tl.lock.Lock()
	logs := make([]TransactionLogItem, len(tl.items))
	copy(logs, tl.items)
	tl.lock.Unlock()

	return logs
}

// Log logs a new log entry to memory and logs to the SDK logs when the level is WARN or ERROR.
func (tl *InMemoryTransactionLogger) Log(level LogLevel, offset int, txnID, attemptID, fmt string, args ...interface{}) error {
	item := TransactionLogItem{
		Level:     level,
		args:      args,
		txnID:     txnID,
		attemptID: attemptID,
		timestamp: time.Now(),
		fmt:       fmt,
	}
	tl.lock.Lock()
	tl.items = append(tl.items, item)
	tl.lock.Unlock()

	if level <= LogWarn {
		logExf(level, offset, txnID+"/"+attemptID+" "+fmt, args...)
	}

	return nil
}

// NoopTransactionLogger logs to the SDK logs when the level is WARN or ERROR.
// Uncommitted: This API may change in the future.
type NoopTransactionLogger struct {
	logDirectlyBelowLevel LogLevel
}

// NewNoopTransactionLogger returns a new noop transaction logger.
// Uncommitted: This API may change in the future.
func NewNoopTransactionLogger() *NoopTransactionLogger {
	return &NoopTransactionLogger{
		logDirectlyBelowLevel: LogInfo,
	}
}

// Logs returns an empty slice.
func (n *NoopTransactionLogger) Logs() []TransactionLogItem {
	return nil
}

// Log logs to the SDK logs when the level is WARN or ERROR.
func (n *NoopTransactionLogger) Log(level LogLevel, offset int, txnID, attemptID, fmt string, args ...interface{}) error {
	if level < n.logDirectlyBelowLevel {
		logExf(level, offset, txnID+"/"+attemptID+" "+fmt, args...)
	}

	return nil
}

type internalTransactionLogWrapper struct {
	wrapped               TransactionLogger
	logDirectlyBelowLevel LogLevel
	txnID                 string
}

func newInternalTransactionLogger(txnID string, wrapped TransactionLogger) *internalTransactionLogWrapper {
	return &internalTransactionLogWrapper{
		wrapped:               wrapped,
		logDirectlyBelowLevel: LogInfo,
		txnID:                 txnID[:5],
	}
}

func (tl *internalTransactionLogWrapper) logExf(attemptID string, level LogLevel, fmt string, args ...interface{}) {
	if attemptID != "" {
		attemptID = attemptID[:5]
	}

	err := tl.wrapped.Log(level, 1, tl.txnID, attemptID, fmt, args...)
	if err != nil {
		log.Printf("Transaction logger error occurred (%s)\n", err)
	}
}

func (tl *internalTransactionLogWrapper) logDebugf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, LogDebug, format, v...)
}

func (tl *internalTransactionLogWrapper) logSchedf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, LogSched, format, v...)
}

func (tl *internalTransactionLogWrapper) logWarnf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, LogWarn, format, v...)
}

func (tl *internalTransactionLogWrapper) logErrorf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, LogError, format, v...)
}

func (tl *internalTransactionLogWrapper) logInfof(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, LogInfo, format, v...)
}
