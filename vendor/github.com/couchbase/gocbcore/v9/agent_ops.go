package gocbcore

import "github.com/couchbase/gocbcore/v9/memd"

// GetCallback is invoked upon completion of a Get operation.
type GetCallback func(*GetResult, error)

// Get retrieves a document.
func (agent *Agent) Get(opts GetOptions, cb GetCallback) (PendingOp, error) {
	return agent.crud.Get(opts, cb)
}

// GetAndTouchCallback is invoked upon completion of a GetAndTouch operation.
type GetAndTouchCallback func(*GetAndTouchResult, error)

// GetAndTouch retrieves a document and updates its expiry.
func (agent *Agent) GetAndTouch(opts GetAndTouchOptions, cb GetAndTouchCallback) (PendingOp, error) {
	return agent.crud.GetAndTouch(opts, cb)
}

// GetAndLockCallback is invoked upon completion of a GetAndLock operation.
type GetAndLockCallback func(*GetAndLockResult, error)

// GetAndLock retrieves a document and locks it.
func (agent *Agent) GetAndLock(opts GetAndLockOptions, cb GetAndLockCallback) (PendingOp, error) {
	return agent.crud.GetAndLock(opts, cb)
}

// GetReplicaCallback is invoked upon completion of a GetReplica operation.
type GetReplicaCallback func(*GetReplicaResult, error)

// GetOneReplica retrieves a document from a replica server.
func (agent *Agent) GetOneReplica(opts GetOneReplicaOptions, cb GetReplicaCallback) (PendingOp, error) {
	return agent.crud.GetOneReplica(opts, cb)
}

// TouchCallback is invoked upon completion of a Touch operation.
type TouchCallback func(*TouchResult, error)

// Touch updates the expiry for a document.
func (agent *Agent) Touch(opts TouchOptions, cb TouchCallback) (PendingOp, error) {
	return agent.crud.Touch(opts, cb)
}

// UnlockCallback is invoked upon completion of a Unlock operation.
type UnlockCallback func(*UnlockResult, error)

// Unlock unlocks a locked document.
func (agent *Agent) Unlock(opts UnlockOptions, cb UnlockCallback) (PendingOp, error) {
	return agent.crud.Unlock(opts, cb)
}

// DeleteCallback is invoked upon completion of a Delete operation.
type DeleteCallback func(*DeleteResult, error)

// Delete removes a document.
func (agent *Agent) Delete(opts DeleteOptions, cb DeleteCallback) (PendingOp, error) {
	return agent.crud.Delete(opts, cb)
}

// StoreCallback is invoked upon completion of a Add, Set or Replace operation.
type StoreCallback func(*StoreResult, error)

// Add stores a document as long as it does not already exist.
func (agent *Agent) Add(opts AddOptions, cb StoreCallback) (PendingOp, error) {
	return agent.crud.Add(opts, cb)
}

// Set stores a document.
func (agent *Agent) Set(opts SetOptions, cb StoreCallback) (PendingOp, error) {
	return agent.crud.Set(opts, cb)
}

// Replace replaces the value of a Couchbase document with another value.
func (agent *Agent) Replace(opts ReplaceOptions, cb StoreCallback) (PendingOp, error) {
	return agent.crud.Replace(opts, cb)
}

// AdjoinCallback is invoked upon completion of a Append or Prepend operation.
type AdjoinCallback func(*AdjoinResult, error)

// Append appends some bytes to a document.
func (agent *Agent) Append(opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	return agent.crud.Append(opts, cb)
}

// Prepend prepends some bytes to a document.
func (agent *Agent) Prepend(opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	return agent.crud.Prepend(opts, cb)
}

// CounterCallback is invoked upon completion of a Increment or Decrement operation.
type CounterCallback func(*CounterResult, error)

// Increment increments the unsigned integer value in a document.
func (agent *Agent) Increment(opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	return agent.crud.Increment(opts, cb)
}

// Decrement decrements the unsigned integer value in a document.
func (agent *Agent) Decrement(opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	return agent.crud.Decrement(opts, cb)
}

// GetRandomCallback is invoked upon completion of a GetRandom operation.
type GetRandomCallback func(*GetRandomResult, error)

// GetRandom retrieves the key and value of a random document stored within Couchbase Server.
func (agent *Agent) GetRandom(opts GetRandomOptions, cb GetRandomCallback) (PendingOp, error) {
	return agent.crud.GetRandom(opts, cb)
}

// GetMetaCallback is invoked upon completion of a GetMeta operation.
type GetMetaCallback func(*GetMetaResult, error)

// GetMeta retrieves a document along with some internal Couchbase meta-data.
func (agent *Agent) GetMeta(opts GetMetaOptions, cb GetMetaCallback) (PendingOp, error) {
	return agent.crud.GetMeta(opts, cb)
}

// SetMetaCallback is invoked upon completion of a SetMeta operation.
type SetMetaCallback func(*SetMetaResult, error)

// SetMeta stores a document along with setting some internal Couchbase meta-data.
func (agent *Agent) SetMeta(opts SetMetaOptions, cb SetMetaCallback) (PendingOp, error) {
	return agent.crud.SetMeta(opts, cb)
}

// DeleteMetaCallback is invoked upon completion of a DeleteMeta operation.
type DeleteMetaCallback func(*DeleteMetaResult, error)

// DeleteMeta deletes a document along with setting some internal Couchbase meta-data.
func (agent *Agent) DeleteMeta(opts DeleteMetaOptions, cb DeleteMetaCallback) (PendingOp, error) {
	return agent.crud.DeleteMeta(opts, cb)
}

// StatsCallback is invoked upon completion of a Stats operation.
type StatsCallback func(*StatsResult, error)

// Stats retrieves statistics information from the server.  Note that as this
// function is an aggregator across numerous servers, there are no guarantees
// about the consistency of the results.  Occasionally, some nodes may not be
// represented in the results, or there may be conflicting information between
// multiple nodes (a vbucket active on two separate nodes at once).
func (agent *Agent) Stats(opts StatsOptions, cb StatsCallback) (PendingOp, error) {
	return agent.stats.Stats(opts, cb)
}

// ObserveCallback is invoked upon completion of a Observe operation.
type ObserveCallback func(*ObserveResult, error)

// Observe retrieves the current CAS and persistence state for a document.
func (agent *Agent) Observe(opts ObserveOptions, cb ObserveCallback) (PendingOp, error) {
	return agent.observe.Observe(opts, cb)
}

// ObserveVbCallback is invoked upon completion of a ObserveVb operation.
type ObserveVbCallback func(*ObserveVbResult, error)

// ObserveVb retrieves the persistence state sequence numbers for a particular VBucket
// and includes additional details not included by the basic version.
func (agent *Agent) ObserveVb(opts ObserveVbOptions, cb ObserveVbCallback) (PendingOp, error) {
	return agent.observe.ObserveVb(opts, cb)
}

// SubDocOp defines a per-operation structure to be passed to MutateIn
// or LookupIn for performing many sub-document operations.
type SubDocOp struct {
	Op    memd.SubDocOpType
	Flags memd.SubdocFlag
	Path  string
	Value []byte
}

// LookupInCallback is invoked upon completion of a LookupIn operation.
type LookupInCallback func(*LookupInResult, error)

// LookupIn performs a multiple-lookup sub-document operation on a document.
func (agent *Agent) LookupIn(opts LookupInOptions, cb LookupInCallback) (PendingOp, error) {
	return agent.crud.LookupIn(opts, cb)
}

// MutateInCallback is invoked upon completion of a MutateIn operation.
type MutateInCallback func(*MutateInResult, error)

// MutateIn performs a multiple-mutation sub-document operation on a document.
func (agent *Agent) MutateIn(opts MutateInOptions, cb MutateInCallback) (PendingOp, error) {
	return agent.crud.MutateIn(opts, cb)
}

// N1QLQueryCallback is invoked upon completion of a N1QLQuery operation.
type N1QLQueryCallback func(*N1QLRowReader, error)

// N1QLQuery executes a N1QL query
func (agent *Agent) N1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	return agent.n1ql.N1QLQuery(opts, cb)
}

// PreparedN1QLQuery executes a prepared N1QL query
func (agent *Agent) PreparedN1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	return agent.n1ql.PreparedN1QLQuery(opts, cb)
}

// AnalyticsQueryCallback is invoked upon completion of a AnalyticsQuery operation.
type AnalyticsQueryCallback func(*AnalyticsRowReader, error)

// AnalyticsQuery executes an analytics query
func (agent *Agent) AnalyticsQuery(opts AnalyticsQueryOptions, cb AnalyticsQueryCallback) (PendingOp, error) {
	return agent.analytics.AnalyticsQuery(opts, cb)
}

// SearchQueryCallback is invoked upon completion of a SearchQuery operation.
type SearchQueryCallback func(*SearchRowReader, error)

// SearchQuery executes a Search query
func (agent *Agent) SearchQuery(opts SearchQueryOptions, cb SearchQueryCallback) (PendingOp, error) {
	return agent.search.SearchQuery(opts, cb)
}

// ViewQueryCallback is invoked upon completion of a ViewQuery operation.
type ViewQueryCallback func(*ViewQueryRowReader, error)

// ViewQuery executes a view query
func (agent *Agent) ViewQuery(opts ViewQueryOptions, cb ViewQueryCallback) (PendingOp, error) {
	return agent.views.ViewQuery(opts, cb)
}

// DoHTTPRequestCallback is invoked upon completion of a DoHTTPRequest operation.
type DoHTTPRequestCallback func(*HTTPResponse, error)

// DoHTTPRequest will perform an HTTP request against one of the HTTP
// services which are available within the SDK.
func (agent *Agent) DoHTTPRequest(req *HTTPRequest, cb DoHTTPRequestCallback) (PendingOp, error) {
	return agent.http.DoHTTPRequest(req, cb)
}

// GetCollectionManifestCallback is invoked upon completion of a GetCollectionManifest operation.
type GetCollectionManifestCallback func(*GetCollectionManifestResult, error)

// GetCollectionManifest fetches the current server manifest. This function will not update the client's collection
// id cache.
func (agent *Agent) GetCollectionManifest(opts GetCollectionManifestOptions, cb GetCollectionManifestCallback) (PendingOp, error) {
	return agent.collections.GetCollectionManifest(opts, cb)
}

// GetCollectionIDCallback is invoked upon completion of a GetCollectionID operation.
type GetCollectionIDCallback func(*GetCollectionIDResult, error)

// GetCollectionID fetches the collection id and manifest id that the collection belongs to, given a scope name
// and collection name. This function will also prime the client's collection id cache.
func (agent *Agent) GetCollectionID(scopeName string, collectionName string, opts GetCollectionIDOptions, cb GetCollectionIDCallback) (PendingOp, error) {
	return agent.collections.GetCollectionID(scopeName, collectionName, opts, cb)
}

// PingCallback is invoked upon completion of a PingKv operation.
type PingCallback func(*PingResult, error)

// Ping pings all of the servers we are connected to and returns
// a report regarding the pings that were performed.
func (agent *Agent) Ping(opts PingOptions, cb PingCallback) (PendingOp, error) {
	return agent.diagnostics.Ping(opts, cb)
}

// Diagnostics returns diagnostics information about the client.
// Mainly containing a list of open connections and their current
// states.
func (agent *Agent) Diagnostics(opts DiagnosticsOptions) (*DiagnosticInfo, error) {
	return agent.diagnostics.Diagnostics(opts)
}

// WaitUntilReadyCallback is invoked upon completion of a WaitUntilReady operation.
type WaitUntilReadyCallback func(*WaitUntilReadyResult, error)
