//nolint:unused
package gocb

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type loggableDocKey struct {
	bucket     string
	scope      string
	collection string
	id         string
}

func newLoggableDocKey(bucket, scope, collection string, id string) loggableDocKey {
	return loggableDocKey{
		bucket:     bucket,
		scope:      scope,
		collection: collection,
		id:         id,
	}
}

func (rdi loggableDocKey) String() string {
	scope := rdi.scope
	if scope == "" {
		scope = "_default"
	}
	collection := rdi.collection
	if collection == "" {
		collection = "_default"
	}
	return redactUserDataString(rdi.bucket + "." + scope + "." + collection + "." + rdi.id)
}

// TransactionLogger is the logger used for logging in transactions.
type TransactionLogger interface {
	Logs() []TransactionLogItem
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

// transactionLogger log to memory, also logging WARN and ERROR logs to the SDK logger.
type transactionLogger struct {
	lock                  sync.Mutex
	items                 []TransactionLogItem
	logDirectlyBelowLevel gocbcore.LogLevel
	txnID                 string
}

func newTransactionLogger() *transactionLogger {
	return &transactionLogger{
		logDirectlyBelowLevel: gocbcore.LogInfo,
		items:                 make([]TransactionLogItem, 0, 256),
	}
}

func (tl *transactionLogger) setTxnID(txnID string) {
	tl.txnID = txnID[:5]
}

func (tl *transactionLogger) Logs() []TransactionLogItem {
	tl.lock.Lock()
	logs := make([]TransactionLogItem, len(tl.items))
	copy(logs, tl.items)
	tl.lock.Unlock()

	return logs
}

func (tl *transactionLogger) Log(level gocbcore.LogLevel, offset int, txnID, attemptID, fmt string, args ...interface{}) error {
	item := TransactionLogItem{
		Level:     LogLevel(level),
		args:      args,
		txnID:     txnID,
		attemptID: attemptID,
		timestamp: time.Now(),
		fmt:       fmt,
	}
	tl.lock.Lock()
	tl.items = append(tl.items, item)
	tl.lock.Unlock()

	if level <= gocbcore.LogWarn {
		logExf(LogLevel(level), offset, txnID+"/"+attemptID+" "+fmt, args...)
	}

	return nil
}

func (tl *transactionLogger) logExf(attemptID string, level gocbcore.LogLevel, fmt string, args ...interface{}) {
	if attemptID != "" {
		attemptID = attemptID[:5]
	}

	err := tl.Log(level, 1, tl.txnID, attemptID, fmt, args...)
	if err != nil {
		log.Printf("Transaction logger error occurred (%s)\n", err)
	}
}

func (tl *transactionLogger) logDebugf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, gocbcore.LogDebug, format, v...)
}

func (tl *transactionLogger) logSchedf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, gocbcore.LogSched, format, v...)
}

func (tl *transactionLogger) logWarnf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, gocbcore.LogWarn, format, v...)
}

func (tl *transactionLogger) logErrorf(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, gocbcore.LogError, format, v...)
}

func (tl *transactionLogger) logInfof(attemptID, format string, v ...interface{}) {
	tl.logExf(attemptID, gocbcore.LogInfo, format, v...)
}
