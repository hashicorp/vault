// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/axiomhq/hyperloglog"
	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"go.uber.org/atomic"
)

const (
	// activitySubPath is the directory under the system view where
	// the log will be stored.
	activitySubPath        = "counters/activity/"
	activityEntityBasePath = "log/entity/"
	activityTokenBasePath  = "log/directtokens/"
	activityQueryBasePath  = "queries/"
	activityConfigKey      = "config"
	activityIntentLogKey   = "endofmonth"

	// sketch for each month that stores hash of client ids
	distinctClientsBasePath = "log/distinctclients/"

	// for testing purposes (public as needed)
	ActivityLogPrefix = "sys/counters/activity/log/"
	ActivityPrefix    = "sys/counters/activity/"

	// Time to wait on perf standby before sending fragment
	activityFragmentStandbyTime = 10 * time.Minute

	// Time between writes of segment to storage
	activitySegmentInterval = 10 * time.Minute

	// Timeout on RPC calls.
	activityFragmentSendTimeout = 1 * time.Minute

	// Timeout on storage calls.
	activitySegmentWriteTimeout = 1 * time.Minute

	// Number of client records to store per segment. Each ClientRecord may
	// consume upto 99 bytes; rounding it to 100bytes. Considering the storage
	// limit of 512KB per storage entry, we can roughly store 512KB/100bytes =
	// 5241 clients; rounding down to 5000 clients.
	activitySegmentClientCapacity = 5000

	// Maximum number of segments per month. This allows for 700K entities per
	// month; 700K/5K. These limits are geared towards controlling the storage
	// implications of persisting activity logs. If we hit a scenario where the
	// storage consequences are less important in comparison to the accuracy of
	// the client activity, these limits can be further relaxed or even be
	// removed.
	activityLogMaxSegmentPerMonth = 140

	// trackedTWESegmentPeriod is a time period of a little over a month, and represents
	// the amount of time that needs to pass after a 1.9 or later upgrade to result in
	// all fragments and segments no longer storing token counts in the directtokens
	// storage path.
	trackedTWESegmentPeriod = 35 * 24

	// Known types of activity events; there's presently two internal event
	// types (tokens/clients with and without entities), but we're beginning
	// to support additional buckets for e.g., ACME requests.
	nonEntityTokenActivityType = "non-entity-token"
	entityActivityType         = "entity"
)

type segmentInfo struct {
	startTimestamp       int64
	currentClients       *activity.EntityActivityLog
	clientSequenceNumber uint64
	// DEPRECATED
	// This field is needed for backward compatibility with fragments
	// and segments created with vault versions before 1.9.
	tokenCount *activity.TokenCount
}

// ActivityLog tracks unique entity counts and non-entity token counts.
// It handles assembling log fragments (and sending them to the active
// node), writing log segments, and precomputing queries.
type ActivityLog struct {
	core            *Core
	configOverrides *ActivityLogCoreConfig

	// ActivityLog.l protects the configuration settings, except enable, and any modifications
	// to the current segment.
	// Acquire "l" before fragmentLock if both must be held.
	l sync.RWMutex

	// fragmentLock protects enable, partialMonthClientTracker, fragment,
	// standbyFragmentsReceived.
	fragmentLock sync.RWMutex

	// enabled indicates if the activity log is enabled for this cluster.
	// This is protected by fragmentLock so we can check with only
	// a single synchronization call.
	enabled bool

	// log destination
	logger log.Logger

	// metrics sink
	metrics metricsutil.Metrics

	// view is the storage location used by ActivityLog,
	// defaults to sys/activity.
	view *BarrierView

	// nodeID is the ID to use for all fragments that
	// are generated.
	// TODO: use secondary ID when available?
	nodeID string

	// current log fragment (may be nil)
	fragment         *activity.LogFragment
	fragmentCreation time.Time

	// Channel to signal a new fragment has been created
	// so it's appropriate to start the timer.
	newFragmentCh chan struct{}

	// Channel for sending fragment immediately
	sendCh chan struct{}

	// Channel to stop background processing
	doneCh chan struct{}

	// track metadata and contents of the most recent log segment
	currentSegment segmentInfo

	// Fragments received from performance standbys
	standbyFragmentsReceived []*activity.LogFragment

	// precomputed queries
	queryStore          *activity.PrecomputedQueryStore
	defaultReportMonths int
	retentionMonths     int

	// channel closed by delete worker when done
	deleteDone chan struct{}

	// channel closed when deletion at startup is done
	// (for unit test robustness)
	retentionDone         chan struct{}
	computationWorkerDone chan struct{}

	// for testing: is config currently being invalidated. protected by l
	configInvalidationInProgress bool

	// partialMonthClientTracker tracks active clients this month.  Protected by fragmentLock.
	partialMonthClientTracker map[string]*activity.EntityRecord

	inprocessExport *atomic.Bool

	// CensusReportDone is a channel used to signal tests upon successful calls
	// to (CensusReporter).Write() in CensusReport.
	CensusReportDone chan bool

	// CensusReportInterval is the testing configuration for time between
	// Write() calls initiated in CensusReport.
	CensusReportInterval time.Duration

	// clock is used to support manipulating time in unit and integration tests
	clock timeutil.Clock
	// precomputedQueryWritten receives an element whenever a precomputed query
	// is written. It's used for unit testing
	precomputedQueryWritten chan struct{}
}

// These non-persistent configuration options allow us to disable
// parts of the implementation for integration testing.
// The default values should turn everything on.
type ActivityLogCoreConfig struct {
	// Enable activity log even if the feature flag not set
	ForceEnable bool

	DisableFragmentWorker bool

	// Do not start timers to send or persist fragments.
	DisableTimers bool

	// CensusReportInterval is the testing configuration for time
	CensusReportInterval time.Duration

	// MinimumRetentionMonths defines the minimum value for retention
	MinimumRetentionMonths int

	// Clock holds a custom clock to modify time.Now, time.Ticker, time.Timer.
	// If nil, the default functions from the time package are used
	Clock timeutil.Clock

	DisableInvalidation bool
}

// NewActivityLog creates an activity log.
func NewActivityLog(core *Core, logger log.Logger, view *BarrierView, metrics metricsutil.Metrics) (*ActivityLog, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	clock := core.activityLogConfig.Clock
	if clock == nil {
		clock = timeutil.DefaultClock{}
	}
	a := &ActivityLog{
		core:                      core,
		configOverrides:           &core.activityLogConfig,
		logger:                    logger,
		view:                      view,
		metrics:                   metrics,
		nodeID:                    hostname,
		newFragmentCh:             make(chan struct{}, 1),
		sendCh:                    make(chan struct{}, 1), // buffered so it can be triggered by fragment size
		doneCh:                    make(chan struct{}, 1),
		partialMonthClientTracker: make(map[string]*activity.EntityRecord),
		CensusReportInterval:      time.Hour * 1,
		clock:                     clock,
		currentSegment: segmentInfo{
			startTimestamp: 0,
			currentClients: &activity.EntityActivityLog{
				Clients: make([]*activity.EntityRecord, 0),
			},
			// tokenCount is deprecated, but must still exist for the current segment
			// so the fragment that was using TWEs before the 1.9 changes
			// can be flushed to the current segment.
			tokenCount: &activity.TokenCount{
				CountByNamespaceID: make(map[string]uint64),
			},
			clientSequenceNumber: 0,
		},
		standbyFragmentsReceived: make([]*activity.LogFragment, 0),
		inprocessExport:          atomic.NewBool(false),
		precomputedQueryWritten:  make(chan struct{}),
	}

	config, err := a.loadConfigOrDefault(core.activeContext)
	if err != nil {
		return nil, err
	}

	a.SetConfigInit(config)

	a.queryStore = activity.NewPrecomputedQueryStore(
		logger,
		view.SubView(activityQueryBasePath),
		config.RetentionMonths)

	return a, nil
}

// saveCurrentSegmentToStorage updates the record of Entities or
// Non Entity Tokens in persistent storage
// :force: forces a save of tokens/entities even if the in-memory log is empty
func (a *ActivityLog) saveCurrentSegmentToStorage(ctx context.Context, force bool) error {
	// Prevent simultaneous changes to segment
	a.l.Lock()
	defer a.l.Unlock()
	return a.saveCurrentSegmentToStorageLocked(ctx, force)
}

// Must be called with l held.
// :force: forces a save of tokens/entities even if the in-memory log is empty
func (a *ActivityLog) saveCurrentSegmentToStorageLocked(ctx context.Context, force bool) error {
	defer a.metrics.MeasureSinceWithLabels([]string{"core", "activity", "segment_write"},
		a.clock.Now(), []metricsutil.Label{})

	// Swap out the pending fragments
	a.fragmentLock.Lock()
	localFragment := a.fragment
	a.fragment = nil
	standbys := a.standbyFragmentsReceived
	a.standbyFragmentsReceived = make([]*activity.LogFragment, 0)
	a.fragmentLock.Unlock()

	// If segment start time is zero, do not update or write
	// (even if force is true).  This can happen if activityLog is
	// disabled after a save as been triggered.
	if a.currentSegment.startTimestamp == 0 {
		return nil
	}

	// Measure the current fragment
	if localFragment != nil {
		a.metrics.IncrCounterWithLabels([]string{"core", "activity", "fragment_size"},
			float32(len(localFragment.Clients)),
			[]metricsutil.Label{
				{"type", "entity"},
			})
		a.metrics.IncrCounterWithLabels([]string{"core", "activity", "fragment_size"},
			float32(len(localFragment.NonEntityTokens)),
			[]metricsutil.Label{
				{"type", "direct_token"},
			})
	}

	// Collect new entities and new tokens.
	saveChanges := false
	newEntities := make(map[string]*activity.EntityRecord)
	for _, f := range append(standbys, localFragment) {
		if f == nil {
			continue
		}
		if len(f.Clients) != 0 || len(f.NonEntityTokens) != 0 {
			saveChanges = true
		}
		for _, e := range f.Clients {
			// We could sort by timestamp to see which is first.
			// We'll ignore that; the order of the append above means
			// that we choose entries in localFragment over those
			// from standby nodes.
			newEntities[e.ClientID] = e
		}
		// As of 1.9, a fragment should no longer have any NonEntityTokens. However
		// in order to not lose any information about the current segment during the
		// month when the client upgrades to 1.9, we must retain this functionality.
		for ns, val := range f.NonEntityTokens {
			// We track these pre-1.9 values in the old location, which is
			// a.currentSegment.tokenCount, as opposed to the counter that stores tokens
			// without entities that have client IDs, namely
			// a.partialMonthClientTracker.nonEntityCountByNamespaceID. This preserves backward
			// compatibility for the precomputedQueryWorkers and the segment storing
			// logic.
			a.currentSegment.tokenCount.CountByNamespaceID[ns] += val
		}
	}

	if !saveChanges {
		return nil
	}

	// Will all new entities fit?  If not, roll over to a new segment.
	available := activitySegmentClientCapacity - len(a.currentSegment.currentClients.Clients)
	remaining := available - len(newEntities)
	excess := 0
	if remaining < 0 {
		excess = -remaining
	}

	segmentClients := a.currentSegment.currentClients.Clients
	excessClients := make([]*activity.EntityRecord, 0, excess)
	for _, record := range newEntities {
		if available > 0 {
			segmentClients = append(segmentClients, record)
			available -= 1
		} else {
			excessClients = append(excessClients, record)
		}
	}
	a.currentSegment.currentClients.Clients = segmentClients

	err := a.saveCurrentSegmentInternal(ctx, force)
	if err != nil {
		// The current fragment(s) have already been placed into the in-memory
		// segment, but we may lose any excess (in excessClients).
		// There isn't a good way to unwind the transaction on failure,
		// so we may just lose some records.
		return err
	}

	if available <= 0 {
		if a.currentSegment.clientSequenceNumber >= activityLogMaxSegmentPerMonth {
			// Cannot send as Warn because it will repeat too often,
			// and disabling/renabling would be complicated.
			a.logger.Trace("too many segments in current month", "dropped", len(excessClients))
			return nil
		}

		// Rotate to next segment
		a.currentSegment.clientSequenceNumber += 1
		if len(excessClients) > activitySegmentClientCapacity {
			a.logger.Warn("too many new active clients, dropping tail", "clients", len(excessClients))
			excessClients = excessClients[:activitySegmentClientCapacity]
		}
		a.currentSegment.currentClients.Clients = excessClients
		err := a.saveCurrentSegmentInternal(ctx, force)
		if err != nil {
			return err
		}
	}
	return nil
}

// :force: forces a save of tokens/entities even if the in-memory log is empty
func (a *ActivityLog) saveCurrentSegmentInternal(ctx context.Context, force bool) error {
	_, err := a.saveSegmentEntitiesInternal(ctx, a.currentSegment, force)
	if err != nil {
		return err
	}
	_, err = a.saveSegmentTokensInternal(ctx, a.currentSegment, force)
	return err
}

func (a *ActivityLog) saveSegmentTokensInternal(ctx context.Context, currentSegment segmentInfo, force bool) (string, error) {
	if len(currentSegment.tokenCount.CountByNamespaceID) == 0 && !force {
		return "", nil
	}
	// RFC (VLT-120) defines this as 1-indexed, but it should be 0-indexed
	tokenPath := fmt.Sprintf("%s%d/0", activityTokenBasePath, currentSegment.startTimestamp)
	// We must still allow for the tokenCount of the current segment to
	// be written to storage, since if we remove this code we will incur
	// data loss for one segment's worth of TWEs.
	// We can get away with simply using the oldest version stored because
	// the storing of versions was introduced at the same time as this code.
	oldestVersion, oldestUpgradeTime, err := a.core.FindOldestVersionTimestamp()
	switch {
	case err != nil:
		a.logger.Error(fmt.Sprintf("unable to retrieve oldest version timestamp: %s", err.Error()))
	case len(a.currentSegment.tokenCount.CountByNamespaceID) > 0 &&
		(oldestUpgradeTime.Add(time.Duration(trackedTWESegmentPeriod * time.Hour)).Before(time.Now())):
		a.logger.Error(fmt.Sprintf("storing nonzero token count over a month after vault was upgraded to %s", oldestVersion))
	default:
		if len(a.currentSegment.tokenCount.CountByNamespaceID) > 0 {
			a.logger.Info("storing nonzero token count")
		}
	}
	tokenCount, err := proto.Marshal(a.currentSegment.tokenCount)
	if err != nil {
		return "", err
	}

	a.logger.Trace("writing segment", "path", tokenPath)
	err = a.view.Put(ctx, &logical.StorageEntry{
		Key:   tokenPath,
		Value: tokenCount,
	})
	if err != nil {
		return "", err
	}

	return tokenPath, nil
}

func (a *ActivityLog) saveSegmentEntitiesInternal(ctx context.Context, currentSegment segmentInfo, force bool) (string, error) {
	entityPath := fmt.Sprintf("%s%d/%d", activityEntityBasePath, currentSegment.startTimestamp, currentSegment.clientSequenceNumber)

	for _, client := range a.currentSegment.currentClients.Clients {
		// Explicitly catch and throw clear error message if client ID creation and storage
		// results in a []byte that doesn't assert into a valid string.
		if !utf8.ValidString(client.ClientID) {
			return "", fmt.Errorf("client ID %q is not a valid string:", client.ClientID)
		}
	}

	if len(currentSegment.currentClients.Clients) == 0 && !force {
		return "", nil
	}
	clients, err := proto.Marshal(currentSegment.currentClients)
	if err != nil {
		return entityPath, err
	}

	a.logger.Trace("writing segment", "path", entityPath)
	err = a.view.Put(ctx, &logical.StorageEntry{
		Key:   entityPath,
		Value: clients,
	})
	if err != nil {
		return "", err
	}
	return entityPath, err
}

// parseSegmentNumberFromPath returns the segment number from a path
// (and if it exists - it is the last element in the path)
func parseSegmentNumberFromPath(path string) (int, bool) {
	// as long as both s and sep are not "", len(elems) >= 1
	elems := strings.Split(path, "/")
	segmentNum, err := strconv.Atoi(elems[len(elems)-1])
	if err != nil {
		return 0, false
	}

	return segmentNum, true
}

// availableLogs returns the start_time(s) (in UTC) associated with months for which logs exist,
// sorted last to first
func (a *ActivityLog) availableLogs(ctx context.Context) ([]time.Time, error) {
	paths := make([]string, 0)
	for _, basePath := range []string{activityEntityBasePath, activityTokenBasePath} {
		p, err := a.view.List(ctx, basePath)
		if err != nil {
			return nil, err
		}

		paths = append(paths, p...)
	}

	pathSet := make(map[time.Time]struct{})
	out := make([]time.Time, 0)
	for _, path := range paths {
		// generate a set of unique start times
		time, err := timeutil.ParseTimeFromPath(path)
		if err != nil {
			return nil, err
		}

		if _, present := pathSet[time]; !present {
			pathSet[time] = struct{}{}
			out = append(out, time)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		// sort in reverse order to make processing most recent segment easier
		return out[i].After(out[j])
	})

	a.logger.Trace("scanned existing logs", "out", out)

	return out, nil
}

// getMostRecentActivityLogSegment gets the times (in UTC) associated with the most recent
// contiguous set of activity logs, sorted in decreasing order (latest to earliest)
func (a *ActivityLog) getMostRecentActivityLogSegment(ctx context.Context) ([]time.Time, error) {
	logTimes, err := a.availableLogs(ctx)
	if err != nil {
		return nil, err
	}

	return timeutil.GetMostRecentContiguousMonths(logTimes), nil
}

// getLastEntitySegmentNumber returns the (non-negative) last segment number for the :startTime:, if it exists
func (a *ActivityLog) getLastEntitySegmentNumber(ctx context.Context, startTime time.Time) (uint64, bool, error) {
	p, err := a.view.List(ctx, activityEntityBasePath+fmt.Sprint(startTime.Unix())+"/")
	if err != nil {
		return 0, false, err
	}

	highestNum := -1
	for _, path := range p {
		if num, ok := parseSegmentNumberFromPath(path); ok {
			if num > highestNum {
				highestNum = num
			}
		}
	}

	if highestNum < 0 {
		// numbers less than 0 are invalid. if a negative number is the highest value, there isn't a segment
		return 0, false, nil
	}

	return uint64(highestNum), true, nil
}

// WalkEntitySegments loads each of the entity segments for a particular start time
func (a *ActivityLog) WalkEntitySegments(ctx context.Context, startTime time.Time, hll *hyperloglog.Sketch, walkFn func(*activity.EntityActivityLog, time.Time, *hyperloglog.Sketch) error) error {
	basePath := activityEntityBasePath + fmt.Sprint(startTime.Unix()) + "/"
	pathList, err := a.view.List(ctx, basePath)
	if err != nil {
		return err
	}

	for _, path := range pathList {
		raw, err := a.view.Get(ctx, basePath+path)
		if err != nil {
			return err
		}
		if raw == nil {
			a.logger.Warn("expected log segment not found", "startTime", startTime, "segment", path)
			continue
		}

		out := &activity.EntityActivityLog{}
		err = proto.Unmarshal(raw.Value, out)
		if err != nil {
			return fmt.Errorf("unable to parse segment %v%v: %w", basePath, path, err)
		}
		err = walkFn(out, startTime, hll)
		if err != nil {
			return fmt.Errorf("unable to walk entities: %w", err)
		}
	}
	return nil
}

// WalkTokenSegments loads each of the token segments (expected 1) for a particular start time
func (a *ActivityLog) WalkTokenSegments(ctx context.Context,
	startTime time.Time,
	walkFn func(*activity.TokenCount),
) error {
	basePath := activityTokenBasePath + fmt.Sprint(startTime.Unix()) + "/"
	pathList, err := a.view.List(ctx, basePath)
	if err != nil {
		return err
	}

	for _, path := range pathList {
		raw, err := a.view.Get(ctx, basePath+path)
		if err != nil {
			return err
		}
		if raw == nil {
			a.logger.Trace("no tokens without entities stored without tracking", "startTime", startTime, "segment", path)
			continue
		}
		out := &activity.TokenCount{}
		err = proto.Unmarshal(raw.Value, out)
		if err != nil {
			return fmt.Errorf("unable to parse token segment %v%v: %w", basePath, path, err)
		}
		walkFn(out)
	}
	return nil
}

// loadPriorEntitySegment populates the in-memory tracker for entity IDs that have
// been active "this month"
func (a *ActivityLog) loadPriorEntitySegment(ctx context.Context, startTime time.Time, sequenceNum uint64) error {
	path := activityEntityBasePath + fmt.Sprint(startTime.Unix()) + "/" + strconv.FormatUint(sequenceNum, 10)
	data, err := a.view.Get(ctx, path)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	out := &activity.EntityActivityLog{}
	err = proto.Unmarshal(data.Value, out)
	if err != nil {
		return err
	}

	a.l.RLock()
	a.fragmentLock.Lock()
	// Handle the (unlikely) case where the end of the month has been reached while background loading.
	// Or the feature has been disabled.
	if a.enabled && startTime.Unix() == a.currentSegment.startTimestamp {
		for _, ent := range out.Clients {
			a.partialMonthClientTracker[ent.ClientID] = ent
		}
	}
	a.fragmentLock.Unlock()
	a.l.RUnlock()

	return nil
}

// loadCurrentClientSegment loads the most recent segment (for "this month")
// into memory (to append new entries), and to the partialMonthClientTracker to
// avoid duplication call with fragmentLock and l held.
func (a *ActivityLog) loadCurrentClientSegment(ctx context.Context, startTime time.Time, sequenceNum uint64) error {
	path := activityEntityBasePath + fmt.Sprint(startTime.Unix()) + "/" + strconv.FormatUint(sequenceNum, 10)
	data, err := a.view.Get(ctx, path)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	out := &activity.EntityActivityLog{}
	err = proto.Unmarshal(data.Value, out)
	if err != nil {
		return err
	}

	if !a.core.perfStandby {
		a.currentSegment = segmentInfo{
			startTimestamp: startTime.Unix(),
			currentClients: &activity.EntityActivityLog{
				Clients: out.Clients,
			},
			tokenCount:           a.currentSegment.tokenCount,
			clientSequenceNumber: sequenceNum,
		}
	} else {
		// populate this for edge case checking (if end of month passes while background loading on standby)
		a.currentSegment.startTimestamp = startTime.Unix()
	}

	for _, client := range out.Clients {
		a.partialMonthClientTracker[client.ClientID] = client
	}

	return nil
}

// tokenCountExists checks if there's a token log for :startTime:
// this function should be called with the lock held
func (a *ActivityLog) tokenCountExists(ctx context.Context, startTime time.Time) (bool, error) {
	p, err := a.view.List(ctx, activityTokenBasePath+fmt.Sprint(startTime.Unix())+"/")
	if err != nil {
		return false, err
	}

	for _, path := range p {
		if num, ok := parseSegmentNumberFromPath(path); ok && num == 0 {
			return true, nil
		}
	}

	return false, nil
}

// loadTokenCount populates the in-memory representation of activity token count
// this function should be called with the lock held
func (a *ActivityLog) loadTokenCount(ctx context.Context, startTime time.Time) error {
	tokenCountExists, err := a.tokenCountExists(ctx, startTime)
	if err != nil {
		return err
	}
	if !tokenCountExists {
		return nil
	}

	path := activityTokenBasePath + fmt.Sprint(startTime.Unix()) + "/0"
	data, err := a.view.Get(ctx, path)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	out := &activity.TokenCount{}
	err = proto.Unmarshal(data.Value, out)
	if err != nil {
		return err
	}

	// An empty map is unmarshaled as nil
	if out.CountByNamespaceID == nil {
		out.CountByNamespaceID = make(map[string]uint64)
	}

	// We must load the tokenCount of the current segment into the activity log
	// so that TWEs counted before the introduction of a client ID for TWEs are
	// still reported in the partial client counts.
	a.currentSegment.tokenCount = out

	return nil
}

// entityBackgroundLoader loads entity activity log records for start_date `t`
func (a *ActivityLog) entityBackgroundLoader(ctx context.Context, wg *sync.WaitGroup, t time.Time, seqNums <-chan uint64) {
	defer wg.Done()
	for seqNum := range seqNums {
		select {
		case <-a.doneCh:
			a.logger.Info("background processing told to halt while loading entities", "time", t, "sequence", seqNum)
			return
		default:
		}

		err := a.loadPriorEntitySegment(ctx, t, seqNum)
		if err != nil {
			a.logger.Error("error loading entity activity log", "time", t, "sequence", seqNum, "err", err)
		}
	}
}

// Initialize a new current segment, based on the current time.
// Call with fragmentLock and l held.
func (a *ActivityLog) startNewCurrentLogLocked(now time.Time) {
	a.logger.Trace("initializing new log")
	a.resetCurrentLog()
	a.currentSegment.startTimestamp = now.Unix()
}

// Should be called with fragmentLock and l held.
func (a *ActivityLog) newMonthCurrentLogLocked(currentTime time.Time) {
	a.logger.Trace("continuing log to new month")
	a.resetCurrentLog()
	monthStart := timeutil.StartOfMonth(currentTime.UTC())
	a.currentSegment.startTimestamp = monthStart.Unix()
}

// Initialize a new current segment, based on the given time
// should be called with fragmentLock and l held.
func (a *ActivityLog) newSegmentAtGivenTime(t time.Time) {
	timestamp := t.Unix()

	a.logger.Trace("starting a segment", "timestamp", timestamp)
	a.resetCurrentLog()
	a.currentSegment.startTimestamp = timestamp
}

// Reset all the current segment state.
// Should be called with fragmentLock and l held.
func (a *ActivityLog) resetCurrentLog() {
	a.currentSegment.startTimestamp = 0
	a.currentSegment.currentClients = &activity.EntityActivityLog{
		Clients: make([]*activity.EntityRecord, 0),
	}

	// We must still initialize the tokenCount to recieve tokenCounts from fragments
	// during the month where customers upgrade to 1.9
	a.currentSegment.tokenCount = &activity.TokenCount{
		CountByNamespaceID: make(map[string]uint64),
	}

	a.currentSegment.clientSequenceNumber = 0

	a.fragment = nil
	a.partialMonthClientTracker = make(map[string]*activity.EntityRecord)

	a.standbyFragmentsReceived = make([]*activity.LogFragment, 0)
}

func (a *ActivityLog) deleteLogWorker(ctx context.Context, startTimestamp int64, whenDone chan struct{}) {
	entityPath := fmt.Sprintf("%v%v/", activityEntityBasePath, startTimestamp)
	tokenPath := fmt.Sprintf("%v%v/", activityTokenBasePath, startTimestamp)

	entitySegments, err := a.view.List(ctx, entityPath)
	if err != nil {
		a.logger.Error("could not list entity paths", "error", err)
		return
	}
	for _, p := range entitySegments {
		err = a.view.Delete(ctx, entityPath+p)
		if err != nil {
			a.logger.Error("could not delete entity log", "error", err)
		}
	}

	tokenSegments, err := a.view.List(ctx, tokenPath)
	if err != nil {
		a.logger.Error("could not list token paths", "error", err)
		return
	}
	for _, p := range tokenSegments {
		err = a.view.Delete(ctx, tokenPath+p)
		if err != nil {
			a.logger.Error("could not delete token log", "error", err)
		}
	}

	// Allow whoever started this as a goroutine to wait for it to finish.
	close(whenDone)
}

func (a *ActivityLog) WaitForDeletion() {
	a.l.Lock()
	// May be nil, if never set
	doneCh := a.deleteDone
	a.l.Unlock()
	if doneCh != nil {
		select {
		case <-doneCh:
			break
		}
	}
}

// refreshFromStoredLog loads the appropriate entities/tokencounts for active and performance standbys
// the most recent segment is loaded synchronously, and older segments are loaded in the background
// this function expects stateLock to be held
func (a *ActivityLog) refreshFromStoredLog(ctx context.Context, wg *sync.WaitGroup, now time.Time) error {
	a.l.Lock()
	defer a.l.Unlock()
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()

	decreasingLogTimes, err := a.getMostRecentActivityLogSegment(ctx)
	if err != nil {
		return err
	}

	if len(decreasingLogTimes) == 0 {
		if a.enabled {
			if a.core.perfStandby {
				// reset the log without updating the timestamp
				a.resetCurrentLog()
			} else {
				a.startNewCurrentLogLocked(now)
			}
		}

		return nil
	}

	mostRecent := decreasingLogTimes[0]

	if !a.enabled {
		a.logger.Debug("activity log not enabled, skipping refresh from storage")
		if !a.core.perfStandby && timeutil.IsCurrentMonth(mostRecent, now) {
			a.logger.Debug("activity log is disabled, cleaning up logs for the current month")
			go a.deleteLogWorker(ctx, mostRecent.Unix(), make(chan struct{}))
		}

		return nil
	}

	if timeutil.IsPreviousMonth(mostRecent, now) {
		// no activity logs to load for this month. if we are enabled, interpret
		// it as having missed the rotation, so let it fall through and load
		// if we missed generating the precomputed query, activeFragmentWorker()
		// will clean things up when it runs next

		a.logger.Trace("no log segments for current month", "mostRecent", mostRecent)
		a.logger.Info("rotating activity log to new month")
	} else if mostRecent.After(now) {
		// we can't do anything if the most recent log is in the future
		a.logger.Warn("timestamp from log to load is in the future", "timestamp", mostRecent)
		return nil
	} else if !timeutil.IsCurrentMonth(mostRecent, now) {
		// the most recent log in storage is 2+ months in the past

		a.logger.Warn("most recent log in storage is 2 or more months in the past.", "timestamp", mostRecent)
		if a.core.perfStandby {
			// reset the log without updating the timestamp
			a.resetCurrentLog()
		} else {
			a.startNewCurrentLogLocked(now)
		}

		return nil
	}

	// load token counts from storage into memory. As of 1.9, this functionality
	// is still required since without it, we would lose replicated TWE counts for the
	// current segment.
	if !a.core.perfStandby {
		err = a.loadTokenCount(ctx, mostRecent)
		if err != nil {
			return err
		}
	}

	// load entity logs from storage into memory
	lastSegment, segmentsExist, err := a.getLastEntitySegmentNumber(ctx, mostRecent)
	if err != nil {
		return err
	}
	if !segmentsExist {
		a.logger.Trace("no entity segments for current month")
		return nil
	}

	err = a.loadCurrentClientSegment(ctx, mostRecent, lastSegment)
	if err != nil || lastSegment == 0 {
		return err
	}
	lastSegment--

	seqNums := make(chan uint64, lastSegment+1)
	wg.Add(1)
	go a.entityBackgroundLoader(ctx, wg, mostRecent, seqNums)

	for n := int(lastSegment); n >= 0; n-- {
		seqNums <- uint64(n)
	}
	close(seqNums)

	return nil
}

// This version is used during construction
func (a *ActivityLog) SetConfigInit(config activityConfig) {
	switch config.Enabled {
	case "enable":
		a.enabled = true
	case "default":
		a.enabled = activityLogEnabledDefault
	case "disable":
		a.enabled = false
	}

	if a.configOverrides.ForceEnable {
		a.enabled = true
	}

	a.defaultReportMonths = config.DefaultReportMonths
	a.retentionMonths = config.RetentionMonths

	if a.retentionMonths < a.configOverrides.MinimumRetentionMonths {
		a.retentionMonths = a.configOverrides.MinimumRetentionMonths
	}

	if a.configOverrides.CensusReportInterval > 0 {
		a.CensusReportInterval = a.configOverrides.CensusReportInterval
	}
}

// This version reacts to user changes
func (a *ActivityLog) SetConfig(ctx context.Context, config activityConfig) {
	a.l.Lock()
	defer a.l.Unlock()

	// enabled is protected by fragmentLock
	a.fragmentLock.Lock()
	originalEnabled := a.enabled
	switch config.Enabled {
	case "enable":
		a.enabled = true
	case "default":
		a.enabled = activityLogEnabledDefault
	case "disable":
		a.enabled = false
	}

	if a.enabled != originalEnabled {
		a.logger.Info("activity log enable changed", "original", originalEnabled, "current", a.enabled)
	}

	if !a.enabled && a.currentSegment.startTimestamp != 0 {
		a.logger.Trace("deleting current segment")
		a.deleteDone = make(chan struct{})
		// this is called from a request under stateLock, so use activeContext
		go a.deleteLogWorker(a.core.activeContext, a.currentSegment.startTimestamp, a.deleteDone)
		a.resetCurrentLog()
	}

	forceSave := false
	if a.enabled && a.currentSegment.startTimestamp == 0 {
		a.startNewCurrentLogLocked(a.clock.Now().UTC())
		// Force a save so we can distinguish between
		//
		// Month N-1: present
		// Month N: <blank because we missed the month end>
		//
		// and
		//
		// Month N-1: present
		// Month N: <blank because disabled and re-enabled>
		forceSave = true
	}
	a.fragmentLock.Unlock()

	if forceSave {
		// l is still held here
		a.saveCurrentSegmentInternal(ctx, true)
	}

	a.defaultReportMonths = config.DefaultReportMonths
	a.retentionMonths = config.RetentionMonths
	if a.retentionMonths < a.configOverrides.MinimumRetentionMonths {
		a.retentionMonths = a.configOverrides.MinimumRetentionMonths
	}

	// check for segments out of retention period, if it has changed
	go a.retentionWorker(ctx, a.clock.Now(), a.retentionMonths)
}

// update the enable flag and reset the current log
func (a *ActivityLog) SetConfigStandby(ctx context.Context, config activityConfig) {
	a.l.Lock()
	defer a.l.Unlock()

	// enable is protected by fragmentLock
	a.fragmentLock.Lock()
	originalEnabled := a.enabled
	switch config.Enabled {
	case "enable":
		a.enabled = true
	case "default":
		a.enabled = activityLogEnabledDefault
	case "disable":
		a.enabled = false
	}

	if a.enabled != originalEnabled {
		a.logger.Info("activity log enable changed", "original", originalEnabled, "current", a.enabled)
		a.resetCurrentLog()
	}
	a.fragmentLock.Unlock()
}

func (a *ActivityLog) queriesAvailable(ctx context.Context) (bool, error) {
	if a.queryStore == nil {
		return false, nil
	}
	return a.queryStore.QueriesAvailable(ctx)
}

// setupActivityLog hooks up the singleton ActivityLog into Core.
func (c *Core) setupActivityLog(ctx context.Context, wg *sync.WaitGroup) error {
	c.activityLogLock.Lock()
	defer c.activityLogLock.Unlock()
	return c.setupActivityLogLocked(ctx, wg)
}

// setupActivityLogLocked hooks up the singleton ActivityLog into Core.
// this function should be called with activityLogLock.
func (c *Core) setupActivityLogLocked(ctx context.Context, wg *sync.WaitGroup) error {
	logger := c.baseLogger.Named("activity")

	if os.Getenv("VAULT_DISABLE_ACTIVITY_LOG") != "" {
		if c.CensusLicensingEnabled() {
			logger.Warn("activity log disabled via environment variable while reporting is enabled. " +
				"Reporting will override, and the activity log will be enabled")
		} else {
			logger.Info("activity log disabled via environment variable")
			return nil
		}
	}

	view := c.systemBarrierView.SubView(activitySubPath)

	manager, err := NewActivityLog(c, logger, view, c.metricSink)
	if err != nil {
		return err
	}
	c.activityLog = manager

	// load activity log for "this month" into memory
	err = manager.refreshFromStoredLog(manager.core.activeContext, wg, manager.clock.Now().UTC())
	if err != nil {
		return err
	}

	// Start the background worker, depending on type
	// Lock already held here, can't use .PerfStandby()
	// The workers need to know the current segment time.
	if c.perfStandby {
		if !c.activityLogConfig.DisableFragmentWorker {
			go manager.perfStandbyFragmentWorker(ctx)
		}
	} else {
		if !c.activityLogConfig.DisableFragmentWorker {
			go manager.activeFragmentWorker(ctx)
		}

		// Check for any intent log, in the background
		manager.computationWorkerDone = make(chan struct{})
		go func() {
			manager.precomputedQueryWorker(ctx)
			close(manager.computationWorkerDone)
		}()

		// Catch up on garbage collection
		// Signal when this is done so that unit tests can proceed.
		manager.retentionDone = make(chan struct{})
		go func(months int) {
			manager.retentionWorker(ctx, manager.clock.Now(), months)
			close(manager.retentionDone)
		}(manager.retentionMonths)

		manager.CensusReportDone = make(chan bool)
		go c.activityLog.CensusReport(ctx, c.CensusAgent(), c.BillingStart())
	}

	return nil
}

// stopActivityLogLocked removes the ActivityLog from Core
// and frees any resources.
// this function should be called with activityLogLock
func (c *Core) stopActivityLogLocked() {
	// preSeal may run before startActivityLog got a chance to complete.
	if c.activityLog != nil {
		// Shut down background worker
		close(c.activityLog.doneCh)
	}

	c.activityLog = nil
}

// stopActivityLog removes the ActivityLog from Core
// and frees any resources.
func (c *Core) stopActivityLog() {
	c.activityLogLock.Lock()
	defer c.activityLogLock.Unlock()
	c.stopActivityLogLocked()
}

func (a *ActivityLog) StartOfNextMonth() time.Time {
	a.l.RLock()
	defer a.l.RUnlock()
	var segmentStart time.Time
	if a.currentSegment.startTimestamp == 0 {
		segmentStart = a.clock.Now().UTC()
	} else {
		segmentStart = time.Unix(a.currentSegment.startTimestamp, 0).UTC()
	}
	// Basing this on the segment start will mean we trigger EOM rollover when
	// necessary because we were down.
	return timeutil.StartOfNextMonth(segmentStart)
}

// perfStandbyFragmentWorker handles scheduling fragments
// to send via RPC; it runs on perf standby nodes only.
func (a *ActivityLog) perfStandbyFragmentWorker(ctx context.Context) {
	timer := a.clock.NewTimer(time.Duration(0))
	fragmentWaiting := false
	// Eat first event, so timer is stopped
	<-timer.C

	endOfMonth := a.clock.NewTimer(a.StartOfNextMonth().Sub(a.clock.Now()))
	if a.configOverrides.DisableTimers {
		endOfMonth.Stop()
	}

	sendFunc := func() {
		ctx, cancel := context.WithTimeout(ctx, activityFragmentSendTimeout)
		defer cancel()
		err := a.sendCurrentFragment(ctx)
		if err != nil {
			a.logger.Warn("activity log fragment lost", "error", err)
		}
	}

	for {
		select {
		case <-a.doneCh:
			// Shutting down activity log.
			if fragmentWaiting && !timer.Stop() {
				<-timer.C
			}
			if !endOfMonth.Stop() {
				<-endOfMonth.C
			}
			return
		case <-a.newFragmentCh:
			// New fragment created, start the timer if not
			// already running
			if !fragmentWaiting {
				fragmentWaiting = true
				if !a.configOverrides.DisableTimers {
					a.logger.Trace("reset fragment timer")
					timer.Reset(activityFragmentStandbyTime)
				}
			}
		case <-timer.C:
			a.logger.Trace("sending fragment on timer expiration")
			fragmentWaiting = false
			sendFunc()
		case <-a.sendCh:
			a.logger.Trace("sending fragment on request")
			// It might be that we get sendCh before fragmentCh
			// if a fragment is created and then immediately fills
			// up to its limit. So we attempt to send even if the timer's
			// not running.
			if fragmentWaiting {
				fragmentWaiting = false
				if !timer.Stop() {
					<-timer.C
				}
			}
			sendFunc()
		case <-endOfMonth.C:
			a.logger.Trace("sending fragment on end of month")
			// Flush the current fragment, if any
			if fragmentWaiting {
				fragmentWaiting = false
				if !timer.Stop() {
					<-timer.C
				}
			}
			sendFunc()

			// clear active entity set
			a.fragmentLock.Lock()
			a.partialMonthClientTracker = make(map[string]*activity.EntityRecord)

			a.fragmentLock.Unlock()

			// Set timer for next month.
			// The current segment *probably* hasn't been set yet (via invalidation),
			// so don't rely on it.
			target := timeutil.StartOfNextMonth(a.clock.Now().UTC())
			endOfMonth.Reset(target.Sub(a.clock.Now()))
		}
	}
}

// activeFragmentWorker handles scheduling the write of the next
// segment.  It runs on active nodes only.
func (a *ActivityLog) activeFragmentWorker(ctx context.Context) {
	ticker := a.clock.NewTicker(activitySegmentInterval)

	endOfMonth := a.clock.NewTimer(a.StartOfNextMonth().Sub(a.clock.Now()))
	if a.configOverrides.DisableTimers {
		endOfMonth.Stop()
	}

	endOfMonthChannel := endOfMonth.C
	if a.core.activityLogConfig.DisableTimers {
		endOfMonthChannel = nil
	}

	writeFunc := func() {
		ctx, cancel := context.WithTimeout(ctx, activitySegmentWriteTimeout)
		defer cancel()
		err := a.saveCurrentSegmentToStorage(ctx, false)
		if err != nil {
			a.logger.Warn("activity log segment not saved, current fragment lost", "error", err)
		}
	}

	// we modify the doneCh in some tests, so let's make sure we don't trip
	// the race detector
	a.l.RLock()
	doneCh := a.doneCh
	a.l.RUnlock()

	for {
		select {
		case <-doneCh:
			// Shutting down activity log.
			ticker.Stop()
			return
		case <-a.newFragmentCh:
			// Just eat the message; the ticker is
			// already running so we don't need to start it.
			// (But we might change the behavior in the future.)
			a.logger.Trace("new local fragment created")
			continue
		case <-ticker.C:
			// It's harder to disable a Ticker so we'll just ignore it.
			if a.configOverrides.DisableTimers {
				continue
			}
			a.logger.Trace("writing segment on timer expiration")
			writeFunc()
		case currentTime := <-endOfMonthChannel:
			err := a.HandleEndOfMonth(ctx, currentTime.UTC())
			if err != nil {
				a.logger.Error("failed to perform end of month rotation", "error", err)
			}

			// Garbage collect any segments or queries based on the immediate
			// value of retentionMonths.
			a.l.RLock()
			go a.retentionWorker(ctx, currentTime.UTC(), a.retentionMonths)
			a.l.RUnlock()

			delta := a.StartOfNextMonth().Sub(a.clock.Now())
			if delta < 20*time.Minute {
				delta = 20 * time.Minute
			}
			a.logger.Trace("scheduling next month", "delta", delta)
			endOfMonth.Reset(delta)
		}
	}
}

type ActivityIntentLog struct {
	PreviousMonth int64 `json:"previous_month"`
	NextMonth     int64 `json:"next_month"`
}

// Handle rotation to end-of-month
// currentTime is an argument for unit-testing purposes
func (a *ActivityLog) HandleEndOfMonth(ctx context.Context, currentTime time.Time) error {
	// Hold lock to prevent segment or enable changing,
	// disable will apply to *next* month.
	a.l.Lock()
	defer a.l.Unlock()

	a.fragmentLock.RLock()
	// Don't bother if disabled
	// since l is locked earlier (and SetConfig() is the only way enabled can change)
	// we don't need to worry about enabled changing during this work
	enabled := a.enabled
	a.fragmentLock.RUnlock()
	if !enabled {
		return nil
	}

	a.logger.Trace("starting end of month processing", "rolloverTime", currentTime)

	err := a.writeIntentLog(ctx, a.currentSegment.startTimestamp, currentTime)
	if err != nil {
		return err
	}
	// Save the current segment; this does not guarantee that fragment will be
	// empty when it returns, but dropping some measurements is acceptable.
	// We use force=true here in case an entry didn't appear this month
	err = a.saveCurrentSegmentToStorageLocked(ctx, true)

	// Don't return this error, just log it, we are done with that segment anyway.
	if err != nil {
		a.logger.Warn("last save of segment failed", "error", err)
	}

	// Advance the log; no need to force a save here because we have
	// the intent log written already.
	//
	// On recovery refreshFromStoredLog() will see we're no longer
	// in the previous month, and recover by calling newMonthCurrentLog
	// again and triggering the precomputed query.
	a.fragmentLock.Lock()
	a.newMonthCurrentLogLocked(currentTime)
	a.fragmentLock.Unlock()

	// Work on precomputed queries in background
	go a.precomputedQueryWorker(ctx)

	return nil
}

// writeIntentLog writes out an intent log for the month
// prevSegmentTimestamp is the timestamp of the segment that we would like to
// transform into a precomputed query.
// nextSegment is the timestamp for the next segment. When invoked by end of
// month processing, this will be the current time and should be in a different
// month than the prevSegmentTimestamp
func (a *ActivityLog) writeIntentLog(ctx context.Context, prevSegmentTimestamp int64, nextSegment time.Time) error {
	nextSegmentTimestamp := timeutil.StartOfMonth(nextSegment.UTC()).Unix()

	// Write out an intent log for the rotation with the current and new segment times.
	intentLog := &ActivityIntentLog{
		PreviousMonth: prevSegmentTimestamp,
		NextMonth:     nextSegmentTimestamp,
	}
	entry, err := logical.StorageEntryJSON(activityIntentLogKey, intentLog)
	if err != nil {
		return err
	}
	err = a.view.Put(ctx, entry)
	if err != nil {
		return err
	}
	return nil
}

// ResetActivityLog is used to extract the current fragment(s) during
// integration testing, so that it can be checked in a race-free way.
func (c *Core) ResetActivityLog() []*activity.LogFragment {
	c.stateLock.RLock()
	a := c.activityLog
	c.stateLock.RUnlock()
	if a == nil {
		return nil
	}

	allFragments := make([]*activity.LogFragment, 1)
	a.fragmentLock.Lock()
	allFragments[0] = a.fragment
	a.fragment = nil

	allFragments = append(allFragments, a.standbyFragmentsReceived...)
	a.standbyFragmentsReceived = make([]*activity.LogFragment, 0)
	a.fragmentLock.Unlock()
	return allFragments
}

func (a *ActivityLog) AddEntityToFragment(entityID string, namespaceID string, timestamp int64) {
	a.AddClientToFragment(entityID, namespaceID, timestamp, false, "")
}

// AddClientToFragment checks a client ID for uniqueness and
// if not already present, adds it to the current fragment.
//
// See note below about AddActivityToFragment.
func (a *ActivityLog) AddClientToFragment(clientID string, namespaceID string, timestamp int64, isTWE bool, mountAccessor string) {
	// TWE == token without entity
	if isTWE {
		a.AddActivityToFragment(clientID, namespaceID, timestamp, nonEntityTokenActivityType, mountAccessor)
		return
	}

	a.AddActivityToFragment(clientID, namespaceID, timestamp, entityActivityType, mountAccessor)
}

// AddActivityToFragment adds a client count event of any type to
// add to the current fragment. ClientIDs must be unique across
// all types; if not already present, we will add it to the current
// fragment. The timestamp is a Unix timestamp *without* nanoseconds,
// as that is what token.CreationTime uses.
func (a *ActivityLog) AddActivityToFragment(clientID string, namespaceID string, timestamp int64, activityType string, mountAccessor string) {
	// Check whether entity ID already recorded
	var present bool

	// TODO: This hack enables separate tracking of events without handling
	// separate storage buckets for counting these event types. Consider
	// removing if the event type is otherwise clear; notably though, this
	// does help ensure clientID uniqueness across different types of tokens,
	// assuming it does not break any other downstream systems.
	if activityType != nonEntityTokenActivityType && activityType != entityActivityType {
		clientID = activityType + "." + clientID
	}

	a.fragmentLock.RLock()
	if a.enabled {
		_, present = a.partialMonthClientTracker[clientID]
	} else {
		present = true
	}
	a.fragmentLock.RUnlock()
	if present {
		return
	}

	// Update current fragment with new active entity
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()

	// Re-check entity ID after re-acquiring lock
	_, present = a.partialMonthClientTracker[clientID]
	if present {
		return
	}

	a.createCurrentFragment()

	clientRecord := &activity.EntityRecord{
		ClientID:      clientID,
		NamespaceID:   namespaceID,
		Timestamp:     timestamp,
		MountAccessor: mountAccessor,
		ClientType:    activityType,
	}

	// Track whether the clientID corresponds to a token without an entity or not.
	// This field is backward compatible, as the default is 0, so records created
	// from pre-1.9 activityLog code will automatically be marked as having an entity.
	if activityType != entityActivityType {
		// TODO: This part needs to be modified potentially for separate
		// storage buckets of custom event types. Consider setting the above
		// condition to activityType == nonEntityTokenEventType in the future.
		clientRecord.NonEntity = true
	}

	a.fragment.Clients = append(a.fragment.Clients, clientRecord)
	a.partialMonthClientTracker[clientRecord.ClientID] = clientRecord
}

// Create the current fragment if it doesn't already exist.
// Must be called with the lock held.
func (a *ActivityLog) createCurrentFragment() {
	if a.fragment == nil {
		a.fragment = &activity.LogFragment{
			OriginatingNode: a.nodeID,
			Clients:         make([]*activity.EntityRecord, 0, 120),
			NonEntityTokens: make(map[string]uint64),
		}
		a.fragmentCreation = a.clock.Now().UTC()

		// Signal that a new segment is available, start
		// the timer to send it.
		a.newFragmentCh <- struct{}{}
	}
}

func (a *ActivityLog) receivedFragment(fragment *activity.LogFragment) {
	a.logger.Trace("received fragment from standby", "node", fragment.OriginatingNode)

	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()

	if !a.enabled {
		return
	}

	for _, e := range fragment.Clients {
		a.partialMonthClientTracker[e.ClientID] = e
	}

	a.standbyFragmentsReceived = append(a.standbyFragmentsReceived, fragment)

	// TODO: check if current segment is full and should be written
}

type ResponseCounts struct {
	DistinctEntities int `json:"distinct_entities"`
	EntityClients    int `json:"entity_clients"`
	NonEntityTokens  int `json:"non_entity_tokens"`
	NonEntityClients int `json:"non_entity_clients"`
	Clients          int `json:"clients"`
}

type ResponseNamespace struct {
	NamespaceID   string           `json:"namespace_id"`
	NamespacePath string           `json:"namespace_path"`
	Counts        ResponseCounts   `json:"counts"`
	Mounts        []*ResponseMount `json:"mounts"`
}

type ResponseMonth struct {
	Timestamp  string               `json:"timestamp"`
	Counts     *ResponseCounts      `json:"counts"`
	Namespaces []*ResponseNamespace `json:"namespaces"`
	NewClients *ResponseNewClients  `json:"new_clients" mapstructure:"new_clients"`
}

type ResponseNewClients struct {
	Counts     *ResponseCounts      `json:"counts"`
	Namespaces []*ResponseNamespace `json:"namespaces"`
}

type ResponseMount struct {
	MountPath string          `json:"mount_path"`
	Counts    *ResponseCounts `json:"counts"`
}

// ActivityLogInjectResponse injects a precomputed query into storage for testing.
func (c *Core) ActivityLogInjectResponse(ctx context.Context, pq *activity.PrecomputedQuery) error {
	c.stateLock.RLock()
	a := c.activityLog
	c.stateLock.RUnlock()
	if a == nil {
		return errors.New("nil activity log")
	}
	return a.queryStore.Put(ctx, pq)
}

func (a *ActivityLog) includeInResponse(query *namespace.Namespace, record *namespace.Namespace) bool {
	if record == nil {
		// Deleted namespace, only include in root queries
		return query.ID == namespace.RootNamespaceID
	}
	return record.HasParent(query)
}

func (a *ActivityLog) DefaultStartTime(endTime time.Time) time.Time {
	// If end time is September 30, then start time should be
	// October 1st to get 12 months of data.
	a.l.RLock()
	defer a.l.RUnlock()

	monthStart := timeutil.StartOfMonth(endTime)
	return monthStart.AddDate(0, -a.defaultReportMonths+1, 0)
}

func (a *ActivityLog) handleQuery(ctx context.Context, startTime, endTime time.Time, limitNamespaces int) (map[string]interface{}, error) {
	var computePartial bool

	// Change the start time to the beginning of the month, and the end time to be the end
	// of the month.
	startTime = timeutil.StartOfMonth(startTime)
	endTime = timeutil.EndOfMonth(endTime)

	// If the endTime of the query is the current month, request data from the queryStore
	// with the endTime equal to the end of the last month, and add in the current month
	// data.
	precomputedQueryEndTime := endTime
	if timeutil.IsCurrentMonth(endTime, a.clock.Now().UTC()) {
		precomputedQueryEndTime = timeutil.EndOfMonth(timeutil.MonthsPreviousTo(1, timeutil.StartOfMonth(endTime)))
		computePartial = true
	}

	pq := &activity.PrecomputedQuery{}
	if startTime.After(precomputedQueryEndTime) && timeutil.IsCurrentMonth(startTime, a.clock.Now().UTC()) {
		// We're only calculating the partial month client count. Skip the precomputation
		// get call.
		pq = &activity.PrecomputedQuery{
			StartTime:  startTime,
			EndTime:    endTime,
			Namespaces: make([]*activity.NamespaceRecord, 0),
			Months:     make([]*activity.MonthRecord, 0),
		}
	} else {
		storedQuery, err := a.queryStore.Get(ctx, startTime, precomputedQueryEndTime)
		if err != nil {
			return nil, err
		}
		if storedQuery == nil {
			// If the storedQuery is nil, that means there's no historical data to process. But, it's possible there's
			// still current month data to process, so rather than returning a 204, let's proceed along like we're
			// just querying the current month.
			storedQuery = &activity.PrecomputedQuery{
				StartTime:  startTime,
				EndTime:    endTime,
				Namespaces: make([]*activity.NamespaceRecord, 0),
				Months:     make([]*activity.MonthRecord, 0),
			}
		}
		pq = storedQuery
	}

	// Calculate the namespace response breakdowns and totals for entities and tokens from the initial
	// namespace data.
	totalEntities, totalTokens, byNamespaceResponse, err := a.calculateByNamespaceResponseForQuery(ctx, pq.Namespaces)
	if err != nil {
		return nil, err
	}

	// If we need to add the current month's client counts into the total, compute the namespace
	// breakdown for the current month as well.
	var partialByMonth map[int64]*processMonth
	var partialByNamespace map[string]*processByNamespace
	var totalEntitiesCurrent int
	var totalTokensCurrent int
	var byNamespaceResponseCurrent []*ResponseNamespace
	if computePartial {
		// Traverse through current month's activitylog data and group clients
		// into months and namespaces
		a.fragmentLock.RLock()
		partialByMonth, partialByNamespace = a.populateNamespaceAndMonthlyBreakdowns()
		a.fragmentLock.RUnlock()

		// Convert the byNamespace breakdowns into structs that are
		// consumable by the /activity endpoint, so as to reuse code between these two
		// endpoints.
		byNamespaceComputation := a.transformALNamespaceBreakdowns(partialByNamespace)

		// Calculate the namespace response breakdowns and totals for entities and tokens from the initial
		// namespace data.
		totalEntitiesCurrent, totalTokensCurrent, byNamespaceResponseCurrent, err = a.calculateByNamespaceResponseForQuery(ctx, byNamespaceComputation)
		if err != nil {
			return nil, err
		}

		// Create a mapping of namespace id to slice index, so that we can efficiently update our results without
		// having to traverse the entire namespace response slice every time.
		nsrMap := make(map[string]int)
		for i, nr := range byNamespaceResponse {
			nsrMap[nr.NamespaceID] = i
		}

		// Rather than blindly appending, which will create duplicates, check our existing counts against the current
		// month counts, and append or update as necessary. We also want to account for mounts and their counts.
		for _, nrc := range byNamespaceResponseCurrent {
			if ndx, ok := nsrMap[nrc.NamespaceID]; ok {
				existingRecord := byNamespaceResponse[ndx]

				// Create a map of the existing mounts, so we don't duplicate them
				mountMap := make(map[string]*ResponseCounts)
				for _, erm := range existingRecord.Mounts {
					mountMap[erm.MountPath] = erm.Counts
				}

				existingRecord.Counts.EntityClients += nrc.Counts.EntityClients
				existingRecord.Counts.Clients += nrc.Counts.Clients
				existingRecord.Counts.DistinctEntities += nrc.Counts.DistinctEntities
				existingRecord.Counts.NonEntityClients += nrc.Counts.NonEntityClients
				existingRecord.Counts.NonEntityTokens += nrc.Counts.NonEntityTokens

				// Check the current month mounts against the existing mounts and if there are matches, update counts
				// accordingly. If there is no match, append the new mount to the existing mounts, so it will be counted
				// later.
				for _, nrcMount := range nrc.Mounts {
					if existingRecordMountCounts, ook := mountMap[nrcMount.MountPath]; ook {
						existingRecordMountCounts.EntityClients += nrcMount.Counts.EntityClients
						existingRecordMountCounts.Clients += nrcMount.Counts.Clients
						existingRecordMountCounts.DistinctEntities += nrcMount.Counts.DistinctEntities
						existingRecordMountCounts.NonEntityClients += nrcMount.Counts.NonEntityClients
						existingRecordMountCounts.NonEntityTokens += nrcMount.Counts.NonEntityTokens
					} else {
						existingRecord.Mounts = append(existingRecord.Mounts, nrcMount)
					}
				}
			} else {
				byNamespaceResponse = append(byNamespaceResponse, nrc)
			}
		}
	}

	// Sort clients within each namespace
	a.sortALResponseNamespaces(byNamespaceResponse)

	if limitNamespaces > 0 {
		totalEntities, totalTokens, byNamespaceResponse = a.limitNamespacesInALResponse(byNamespaceResponse, limitNamespaces)
	}

	distinctEntitiesResponse := totalEntities
	if computePartial {
		currentMonth, err := a.computeCurrentMonthForBillingPeriod(ctx, partialByMonth, startTime, endTime)
		if err != nil {
			return nil, err
		}

		// Add the namespace attribution for the current month to the newly computed current month value. Note
		// that transformMonthBreakdowns calculates a superstruct of the required namespace struct due to its
		// primary use-case being for precomputedQueryWorker, but we will reuse this code for brevity and extract
		// the namespaces from it.
		currentMonthNamespaceAttribution := a.transformMonthBreakdowns(partialByMonth)

		// Ensure that there is only one element in this list -- if not, warn.
		if len(currentMonthNamespaceAttribution) > 1 {
			a.logger.Warn("more than one month worth of namespace and mount attribution calculated for "+
				"current month values", "number of months", len(currentMonthNamespaceAttribution))
		}
		if len(currentMonthNamespaceAttribution) == 0 {
			a.logger.Warn("no month data found, returning query with no namespace attribution for current month")
		} else {
			currentMonth.Namespaces = currentMonthNamespaceAttribution[0].Namespaces
			currentMonth.NewClients.Namespaces = currentMonthNamespaceAttribution[0].NewClients.Namespaces
		}
		pq.Months = append(pq.Months, currentMonth)
		distinctEntitiesResponse += pq.Months[len(pq.Months)-1].NewClients.Counts.EntityClients
	}

	// Now populate the response based on breakdowns.
	responseData := make(map[string]interface{})
	responseData["start_time"] = pq.StartTime.Format(time.RFC3339)

	// If we computed partial counts, we should return the actual end time we computed counts for, not the pre-computed
	// query end time. If we don't do this, the end_time in the response doesn't match the actual data in the response,
	// which is confusing. Note that regardless of what end time is given, if it falls within the current month, it will
	// be set to the end of the current month. This is definitely suboptimal, and possibly confusing, but still an
	// improvement over using the pre-computed query end time.
	if computePartial {
		responseData["end_time"] = endTime.Format(time.RFC3339)
	} else {
		responseData["end_time"] = pq.EndTime.Format(time.RFC3339)
	}

	responseData["by_namespace"] = byNamespaceResponse
	responseData["total"] = &ResponseCounts{
		DistinctEntities: distinctEntitiesResponse,
		EntityClients:    totalEntities + totalEntitiesCurrent,
		NonEntityTokens:  totalTokens + totalTokensCurrent,
		NonEntityClients: totalTokens + totalTokensCurrent,
		Clients:          totalEntities + totalEntitiesCurrent + totalTokens + totalTokensCurrent,
	}

	// Create and populate the month response structs based on the monthly breakdown.
	months, err := a.prepareMonthsResponseForQuery(ctx, pq.Months)
	if err != nil {
		return nil, err
	}

	// Sort the months and clients within each month before adding the months to the response
	a.sortActivityLogMonthsResponse(months)

	// Modify the final month output to make response more consumable based on API request
	months = a.modifyResponseMonths(months, startTime, endTime)
	responseData["months"] = months

	return responseData, nil
}

// modifyResponseMonths fills out various parts of the query structure to help
// activity log clients parse the returned query.
func (a *ActivityLog) modifyResponseMonths(months []*ResponseMonth, start time.Time, end time.Time) []*ResponseMonth {
	if len(months) == 0 {
		return months
	}
	start = timeutil.StartOfMonth(start)
	end = timeutil.EndOfMonth(end)
	if timeutil.IsCurrentMonth(end, a.clock.Now().UTC()) {
		end = timeutil.EndOfMonth(timeutil.StartOfMonth(end).AddDate(0, -1, 0))
	}
	modifiedResponseMonths := make([]*ResponseMonth, 0)
	firstMonth, err := time.Parse(time.RFC3339, months[0].Timestamp)
	if err != nil {
		return months
	}
	for start.Before(firstMonth) && !timeutil.IsCurrentMonth(start, firstMonth) {
		monthPlaceholder := &ResponseMonth{Timestamp: start.UTC().Format(time.RFC3339)}
		modifiedResponseMonths = append(modifiedResponseMonths, monthPlaceholder)
		start = timeutil.StartOfMonth(start.AddDate(0, 1, 0))
	}
	modifiedResponseMonths = append(modifiedResponseMonths, months...)
	lastMonthStart, err := time.Parse(time.RFC3339, modifiedResponseMonths[len(modifiedResponseMonths)-1].Timestamp)
	if err != nil {
		return modifiedResponseMonths
	}
	lastMonth := timeutil.EndOfMonth(lastMonthStart)
	for lastMonth.Before(end) && !timeutil.IsCurrentMonth(end, lastMonth) {
		lastMonth = timeutil.StartOfMonth(lastMonth).AddDate(0, 1, 0)
		monthPlaceholder := &ResponseMonth{Timestamp: lastMonth.UTC().Format(time.RFC3339)}
		modifiedResponseMonths = append(modifiedResponseMonths, monthPlaceholder)

		// reset lastMonth to be the end of the month so we can make an apt comparison
		// in the next loop iteration
		lastMonth = timeutil.EndOfMonth(lastMonth)
	}
	return modifiedResponseMonths
}

type activityConfig struct {
	// DefaultReportMonths are the default number of months that are returned on
	// a report. The zero value uses the system default of 12.
	DefaultReportMonths int `json:"default_report_months"`

	// RetentionMonths defines the number of months we want to retain data. The
	// zero value uses the system default of 24 months.
	RetentionMonths int `json:"retention_months"`

	// Enabled is one of enable, disable, default.
	Enabled string `json:"enabled"`

	CensusReportInterval time.Duration `json:"census_report_interval"`
}

func defaultActivityConfig() activityConfig {
	return activityConfig{
		DefaultReportMonths: 12,
		RetentionMonths:     24,
		Enabled:             "default",
	}
}

func (a *ActivityLog) loadConfigOrDefault(ctx context.Context) (activityConfig, error) {
	// Load from storage
	var config activityConfig
	configRaw, err := a.view.Get(ctx, activityConfigKey)
	if err != nil {
		return config, err
	}
	if configRaw == nil {
		return defaultActivityConfig(), nil
	}

	if err := configRaw.DecodeJSON(&config); err != nil {
		return config, err
	}

	return config, nil
}

// HandleTokenUsage adds the TokenEntry to the current fragment of the activity log
// This currently occurs on token usage only.
func (a *ActivityLog) HandleTokenUsage(ctx context.Context, entry *logical.TokenEntry, clientID string, isTWE bool) error {
	// First, check if a is enabled, so as to avoid the cost of creating an ID for
	// tokens without entities in the case where it not.
	a.fragmentLock.RLock()
	if !a.enabled {
		a.fragmentLock.RUnlock()
		return nil
	}
	a.fragmentLock.RUnlock()

	// Do not count wrapping tokens in client count
	if IsWrappingToken(entry) {
		return nil
	}

	// Do not count root tokens in client count.
	if entry.IsRoot() {
		return nil
	}

	// Tokens created for the purpose of Link should bypass counting for billing purposes
	if entry.InternalMeta != nil && entry.InternalMeta[IgnoreForBilling] == "true" {
		return nil
	}

	// Look up the mount accessor of the auth method that issued the token, taking care to resolve the token path
	// against the token namespace, which may not be the same as the request namespace!
	tokenNS, err := NamespaceByID(ctx, entry.NamespaceID, a.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}
	tokenCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	mountAccessor := ""
	mountEntry := a.core.router.MatchingMountEntry(tokenCtx, entry.Path)
	if mountEntry != nil {
		mountAccessor = mountEntry.Accessor
	}

	// Parse an entry's client ID and add it to the activity log
	a.AddClientToFragment(clientID, entry.NamespaceID, entry.CreationTime, isTWE, mountAccessor)
	return nil
}

func (a *ActivityLog) namespaceToLabel(ctx context.Context, nsID string) string {
	ns, err := NamespaceByID(ctx, nsID, a.core)
	if err != nil || ns == nil {
		return fmt.Sprintf("deleted-%v", nsID)
	}
	if ns.Path == "" {
		return "root"
	}
	return ns.Path
}

type (
	summaryByNamespace map[string]*processByNamespace
	summaryByMount     map[string]*processMount
	summaryByMonth     map[int64]*processMonth
)

type processCounts struct {
	// entityID -> present
	Entities map[string]struct{}
	// count. This exists for backward compatibility
	Tokens uint64
	// clientID -> present
	NonEntities map[string]struct{}
}

func newProcessCounts() *processCounts {
	return &processCounts{
		Entities:    make(map[string]struct{}),
		Tokens:      0,
		NonEntities: make(map[string]struct{}),
	}
}

func (p *processCounts) delete(client *activity.EntityRecord) {
	if !p.contains(client) {
		return
	}
	if client.NonEntity {
		delete(p.NonEntities, client.ClientID)
	} else {
		delete(p.Entities, client.ClientID)
	}
}

func (p *processCounts) add(client *activity.EntityRecord) {
	if client.NonEntity {
		p.NonEntities[client.ClientID] = struct{}{}
	} else {
		p.Entities[client.ClientID] = struct{}{}
	}
}

func (p *processCounts) contains(client *activity.EntityRecord) bool {
	if client.NonEntity {
		_, ok := p.NonEntities[client.ClientID]
		return ok
	}
	_, ok := p.Entities[client.ClientID]
	return ok
}

type processMount struct {
	Counts *processCounts
}

func newProcessMount() *processMount {
	return &processMount{
		Counts: newProcessCounts(),
	}
}

func (p *processMount) add(client *activity.EntityRecord) {
	p.Counts.add(client)
}

func (p *processMount) delete(client *activity.EntityRecord) {
	p.Counts.delete(client)
}

func (s summaryByMount) add(client *activity.EntityRecord) {
	if _, present := s[client.MountAccessor]; !present {
		s[client.MountAccessor] = newProcessMount()
	}
	s[client.MountAccessor].add(client)
}

func (s summaryByMount) delete(client *activity.EntityRecord) {
	if m, present := s[client.MountAccessor]; present {
		m.delete(client)
	}
}

type processByNamespace struct {
	Counts *processCounts
	Mounts summaryByMount
}

func newByNamespace() *processByNamespace {
	return &processByNamespace{
		Counts: newProcessCounts(),
		Mounts: make(summaryByMount),
	}
}

func (p *processByNamespace) add(client *activity.EntityRecord) {
	p.Counts.add(client)
	p.Mounts.add(client)
}

func (p *processByNamespace) delete(client *activity.EntityRecord) {
	p.Counts.delete(client)
	p.Mounts.delete(client)
}

func (s summaryByNamespace) add(client *activity.EntityRecord) {
	if _, present := s[client.NamespaceID]; !present {
		s[client.NamespaceID] = newByNamespace()
	}
	s[client.NamespaceID].add(client)
}

func (s summaryByNamespace) delete(client *activity.EntityRecord) {
	if n, present := s[client.NamespaceID]; present {
		n.delete(client)
	}
}

type processNewClients struct {
	Counts     *processCounts
	Namespaces summaryByNamespace
}

func newProcessNewClients() *processNewClients {
	return &processNewClients{
		Counts:     newProcessCounts(),
		Namespaces: make(summaryByNamespace),
	}
}

func (p *processNewClients) add(client *activity.EntityRecord) {
	p.Counts.add(client)
	p.Namespaces.add(client)
}

func (p *processNewClients) delete(client *activity.EntityRecord) {
	p.Counts.delete(client)
	p.Namespaces.delete(client)
}

type processMonth struct {
	Counts     *processCounts
	Namespaces summaryByNamespace
	NewClients *processNewClients
}

func newProcessMonth() *processMonth {
	return &processMonth{
		Counts:     newProcessCounts(),
		Namespaces: make(summaryByNamespace),
		NewClients: newProcessNewClients(),
	}
}

func (p *processMonth) add(client *activity.EntityRecord) {
	p.Counts.add(client)
	p.NewClients.add(client)
	p.Namespaces.add(client)
}

func (s summaryByMonth) add(client *activity.EntityRecord, startTime time.Time) {
	monthTimestamp := timeutil.StartOfMonth(startTime).UTC().Unix()
	if _, present := s[monthTimestamp]; !present {
		s[monthTimestamp] = newProcessMonth()
	}
	s[monthTimestamp].add(client)
}

// processClientRecord parses the client record e and stores the breakdowns in
// the maps provided.
func processClientRecord(e *activity.EntityRecord, byNamespace summaryByNamespace, byMonth summaryByMonth, startTime time.Time) {
	byNamespace.add(e)
	byMonth.add(e, startTime)
}

// handleEntitySegment processes the record and adds it to the correct month/
// namespace breakdown maps, as well as to the hyperloglog for the month. New
// clients are deduplicated in opts.byMonth so that clients will only appear in
// the first month in which they are seen.
// This method must be called in reverse chronological order of the months (with
// the most recent month being called before previous months)
func (a *ActivityLog) handleEntitySegment(l *activity.EntityActivityLog, segmentTime time.Time, hll *hyperloglog.Sketch, opts pqOptions) error {
	for _, e := range l.Clients {

		processClientRecord(e, opts.byNamespace, opts.byMonth, segmentTime)
		hll.Insert([]byte(e.ClientID))

		// step forward in time through the months to check if the client is
		// present. If it is, delete it. This is because the client should only
		// be reported as new in the earliest month that it was seen
		finalMonth := timeutil.StartOfMonth(opts.activePeriodEnd).UTC()
		for currMonth := timeutil.StartOfMonth(segmentTime).UTC(); currMonth.Before(finalMonth); currMonth = timeutil.StartOfNextMonth(currMonth).UTC() {
			// Invalidate the client from being a new client in the next month
			next := timeutil.StartOfNextMonth(currMonth).UTC().Unix()
			if _, present := opts.byMonth[next]; present {
				// delete from the new clients map for the next month
				// this will handle deleting from the per-namespace and per-mount maps of NewClients
				opts.byMonth[next].NewClients.delete(e)
			}
		}
	}

	return nil
}

// breakdownTokenSegment handles a TokenCount record, adding it to the namespace breakdown
func (a *ActivityLog) breakdownTokenSegment(l *activity.TokenCount, byNamespace map[string]*processByNamespace) {
	for nsID, v := range l.CountByNamespaceID {
		if _, present := byNamespace[nsID]; !present {
			byNamespace[nsID] = newByNamespace()
		}
		byNamespace[nsID].Counts.Tokens += v
	}
}

func (a *ActivityLog) writePrecomputedQuery(ctx context.Context, segmentTime time.Time, opts pqOptions) error {
	pq := &activity.PrecomputedQuery{
		StartTime:  segmentTime,
		EndTime:    opts.endTime,
		Namespaces: make([]*activity.NamespaceRecord, 0, len(opts.byNamespace)),
		Months:     make([]*activity.MonthRecord, 0, len(opts.byMonth)),
	}
	// this will transform the byMonth map into the correctly formatted protobuf
	pq.Months = a.transformMonthBreakdowns(opts.byMonth)

	// the byNamespace map also needs to be transformed into a protobuf
	for nsID, entry := range opts.byNamespace {
		mountRecord := make([]*activity.MountRecord, 0, len(entry.Mounts))
		for mountAccessor, mountData := range entry.Mounts {
			mountRecord = append(mountRecord, &activity.MountRecord{
				MountPath: a.mountAccessorToMountPath(mountAccessor),
				Counts: &activity.CountsRecord{
					EntityClients:    len(mountData.Counts.Entities),
					NonEntityClients: int(mountData.Counts.Tokens) + len(mountData.Counts.NonEntities),
				},
			})
		}

		pq.Namespaces = append(pq.Namespaces, &activity.NamespaceRecord{
			NamespaceID:     nsID,
			Entities:        uint64(len(entry.Counts.Entities)),
			NonEntityTokens: entry.Counts.Tokens + uint64(len(entry.Counts.NonEntities)),
			Mounts:          mountRecord,
		})
	}
	err := a.queryStore.Put(ctx, pq)
	if err != nil {
		a.logger.Warn("failed to store precomputed query", "error", err)
	}
	return nil
}

// pqOptions holds fields that will be used when creating precomputed queries
// These fields will remain the same for every segment that a precomputed query worker is handling
type pqOptions struct {
	byNamespace map[string]*processByNamespace
	byMonth     map[int64]*processMonth
	// endTime sets the end time of the precomputed query.
	// When invoked on schedule by the precomputedQueryWorker, this is the end of the month that just finished.
	endTime time.Time
	// activePeriodStart is the earliest date in our retention window
	activePeriodStart time.Time
	// activePeriodEnd is the latest date in our retention window.
	// When invoked on schedule by the precomputedQueryWorker, this will be the timestamp of the most recent segment
	// that's present in storage
	activePeriodEnd time.Time
}

// segmentToPrecomputedQuery processes a single segment
func (a *ActivityLog) segmentToPrecomputedQuery(ctx context.Context, segmentTime time.Time, reader SegmentReader, opts pqOptions) error {
	hyperloglog, err := a.CreateOrFetchHyperlogLog(ctx, segmentTime)
	if err != nil {
		// We were unable to create or fetch the hll, but we should still
		// continue with our precomputation
		a.logger.Warn("unable to create or fetch hyperloglog", "start time", segmentTime, "error", err)
	}

	// Iterate through entities, adding them to the hyperloglog and the summary maps in opts
	for {
		entity, err := reader.ReadEntity(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			a.logger.Warn("failed to read segment", "error", err)
			return err
		}
		err = a.handleEntitySegment(entity, segmentTime, hyperloglog, opts)
		if err != nil {
			a.logger.Warn("failed to handle entity segment", "error", err)
			return err
		}
	}

	// Store the hyperloglog
	err = a.StoreHyperlogLog(ctx, segmentTime, hyperloglog)
	if err != nil {
		a.logger.Warn("failed to store hyperloglog for month", "start time", segmentTime, "error", err)
	}

	// Iterate through any tokens and add them to per namespace map
	for {
		token, err := reader.ReadToken(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			a.logger.Warn("failed to load token counts", "error", err)
			return err
		}
		a.breakdownTokenSegment(token, opts.byNamespace)
	}

	// write metrics
	for nsID, entry := range opts.byNamespace {
		// If this is the most recent month, or the start of the reporting period, output
		// a metric for each namespace.
		if segmentTime == opts.activePeriodEnd {
			a.metrics.SetGaugeWithLabels(
				[]string{"identity", "entity", "active", "monthly"},
				float32(len(entry.Counts.Entities)),
				[]metricsutil.Label{
					{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
				},
			)
			a.metrics.SetGaugeWithLabels(
				[]string{"identity", "nonentity", "active", "monthly"},
				float32(len(entry.Counts.NonEntities))+float32(entry.Counts.Tokens),
				[]metricsutil.Label{
					{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
				},
			)
		} else if segmentTime == opts.activePeriodStart {
			a.metrics.SetGaugeWithLabels(
				[]string{"identity", "entity", "active", "reporting_period"},
				float32(len(entry.Counts.Entities)),
				[]metricsutil.Label{
					{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
				},
			)
			a.metrics.SetGaugeWithLabels(
				[]string{"identity", "nonentity", "active", "reporting_period"},
				float32(len(entry.Counts.NonEntities))+float32(entry.Counts.Tokens),
				[]metricsutil.Label{
					{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
				},
			)
		}
	}

	// convert the maps to the proper format and write them as precomputed queries
	return a.writePrecomputedQuery(ctx, segmentTime, opts)
}

// goroutine to process the request in the intent log, creating precomputed queries.
// We expect the return value won't be checked, so log errors as they occur
// (but for unit testing having the error return should help.)
func (a *ActivityLog) precomputedQueryWorker(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Cancel the context if activity log is shut down.
	// This will cause the next storage operation to fail.
	a.l.RLock()
	// doneCh is modified in some tests, so we don't want to access that member
	// without a lock, but we don't want to hold the lock for the entire lifetime
	// of this goroutine.  Passing the channel to the goroutine works here because
	// no tests depend on us accessing the new doneCh after modifying the field.
	go func(done chan struct{}) {
		select {
		case <-done:
			cancel()
		case <-ctx.Done():
			break
		}
	}(a.doneCh)
	a.l.RUnlock()

	// Load the intent log
	rawIntentLog, err := a.view.Get(ctx, activityIntentLogKey)
	if err != nil {
		a.logger.Warn("could not load intent log", "error", err)
		return err
	}
	if rawIntentLog == nil {
		a.logger.Trace("no intent log found")
		return err
	}
	var intent ActivityIntentLog
	err = json.Unmarshal(rawIntentLog.Value, &intent)
	if err != nil {
		a.logger.Warn("could not parse intent log", "error", err)
		return err
	}

	// currentMonth could change (from another month end) after we release the lock.
	// But, it's not critical to correct operation; this is a check for intent logs that are
	// too old, and startTimestamp should only go forward (unless it is zero.)
	// If there's an intent log, finish it even if the feature is currently disabled.
	a.l.RLock()
	currentMonth := a.currentSegment.startTimestamp
	// Base retention period on the month we are generating (even in the past)--- a.clock.Now()
	// would work but this will be easier to control in tests.
	retentionWindow := timeutil.MonthsPreviousTo(a.retentionMonths, time.Unix(intent.NextMonth, 0).UTC())
	a.l.RUnlock()
	if currentMonth != 0 && intent.NextMonth != currentMonth {
		a.logger.Warn("intent log does not match current segment",
			"intent", intent.NextMonth, "current", currentMonth)
		return errors.New("intent log is too far in the past")
	}

	lastMonth := intent.PreviousMonth
	a.logger.Info("computing queries", "month", time.Unix(lastMonth, 0).UTC())

	times, err := a.availableLogs(ctx)
	if err != nil {
		a.logger.Warn("could not list available logs", "error", err)
		return err
	}
	if len(times) == 0 {
		a.logger.Warn("no months in storage")
		a.view.Delete(ctx, activityIntentLogKey)
		return errors.New("previous month not found")
	}
	if times[0].Unix() != lastMonth {
		a.logger.Warn("last month not in storage", "latest", times[0].Unix())
		a.view.Delete(ctx, activityIntentLogKey)
		return errors.New("previous month not found")
	}

	byNamespace := make(map[string]*processByNamespace)
	byMonth := make(map[int64]*processMonth)

	endTime := timeutil.EndOfMonth(time.Unix(lastMonth, 0).UTC())
	activePeriodStart := timeutil.MonthsPreviousTo(a.defaultReportMonths, endTime)
	// If not enough data, report as much as we have in the window
	if activePeriodStart.Before(times[len(times)-1]) {
		activePeriodStart = times[len(times)-1]
	}
	opts := pqOptions{
		byNamespace:       byNamespace,
		byMonth:           byMonth,
		endTime:           endTime,
		activePeriodStart: activePeriodStart,
		activePeriodEnd:   times[0],
	}
	// "times" is already in reverse order, start building the per-namespace maps
	// from the last month backward
	for _, startTime := range times {
		// Do not work back further than the current retention window,
		// which will just get deleted anyway.
		if startTime.Before(retentionWindow) {
			break
		}
		reader, err := a.NewSegmentFileReader(ctx, startTime)
		if err != nil {
			return err
		}
		err = a.segmentToPrecomputedQuery(ctx, startTime, reader, opts)
		if err != nil {
			return err
		}
	}

	// delete the intent log
	a.view.Delete(ctx, activityIntentLogKey)

	a.logger.Info("finished computing queries", "month", endTime)

	select {
	case a.precomputedQueryWritten <- struct{}{}:
	default:
	}
	return nil
}

// goroutine to delete any segments or precomputed queries older than
// the retention period.
// We expect the return value won't be checked, so log errors as they occur
// (but for unit testing having the error return should help.)
func (a *ActivityLog) retentionWorker(ctx context.Context, currentTime time.Time, retentionMonths int) error {
	if a.core.activityLogConfig.DisableTimers {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Cancel the context if activity log is shut down.
	// This will cause the next storage operation to fail.
	a.l.RLock()
	doneCh := a.doneCh
	a.l.RUnlock()
	go func() {
		select {
		case <-doneCh:
			cancel()
		case <-ctx.Done():
			break
		}
	}()

	// everything >= the threshold is OK
	retentionThreshold := timeutil.MonthsPreviousTo(retentionMonths, currentTime)

	available, err := a.availableLogs(ctx)
	if err != nil {
		a.logger.Warn("could not list segments", "error", err)
		return err
	}
	for _, t := range available {
		// One at a time seems OK
		if t.Before(retentionThreshold) {
			a.logger.Trace("deleting segments", "startTime", t)
			a.deleteLogWorker(ctx, t.Unix(), make(chan struct{}))
		}
	}

	if a.queryStore != nil {
		err = a.queryStore.DeleteQueriesBefore(ctx, retentionThreshold)
		if err != nil {
			a.logger.Warn("deletion of queries failed", "error", err)
		}
		return err
	}

	return nil
}

// Periodic report of number of active entities, with the current month.
// We don't break this down by namespace because that would require going to storage (that information
// is not currently stored in memory.)
func (a *ActivityLog) PartialMonthMetrics(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()
	if !a.enabled {
		// Empty list
		return []metricsutil.GaugeLabelValues{}, nil
	}
	count := len(a.partialMonthClientTracker)

	return []metricsutil.GaugeLabelValues{
		{
			Labels: []metricsutil.Label{},
			Value:  float32(count),
		},
	}, nil
}

func (c *Core) activeEntityGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	c.stateLock.RLock()
	a := c.activityLog
	c.stateLock.RUnlock()
	if a == nil {
		return []metricsutil.GaugeLabelValues{}, nil
	}
	return a.PartialMonthMetrics(ctx)
}

// populateNamespaceAndMonthlyBreakdowns traverses the partial month data
// stored in memory and groups them by months and namespaces.
func (a *ActivityLog) populateNamespaceAndMonthlyBreakdowns() (map[int64]*processMonth, map[string]*processByNamespace) {
	// Parse the monthly clients and prepare the breakdowns.
	byNamespace := make(map[string]*processByNamespace)
	byMonth := make(map[int64]*processMonth)
	for _, e := range a.partialMonthClientTracker {
		processClientRecord(e, byNamespace, byMonth, a.clock.Now())
	}
	return byMonth, byNamespace
}

// transformMonthBreakdowns converts a map of unix timestamp -> processMonth to
// a slice of MonthRecord
func (a *ActivityLog) transformMonthBreakdowns(byMonth map[int64]*processMonth) []*activity.MonthRecord {
	monthly := make([]*activity.MonthRecord, 0)
	processByNamespaces := func(nsMap map[string]*processByNamespace) []*activity.MonthlyNamespaceRecord {
		nsRecord := make([]*activity.MonthlyNamespaceRecord, 0, len(nsMap))
		for nsID, nsData := range nsMap {
			// Process mount specific data within a namespace within a given month
			mountRecord := make([]*activity.MountRecord, 0, len(nsMap[nsID].Mounts))
			for mountAccessor, mountData := range nsMap[nsID].Mounts {
				mountRecord = append(mountRecord, &activity.MountRecord{
					MountPath: a.mountAccessorToMountPath(mountAccessor),
					Counts: &activity.CountsRecord{
						EntityClients:    len(mountData.Counts.Entities),
						NonEntityClients: int(mountData.Counts.Tokens) + len(mountData.Counts.NonEntities),
					},
				})
			}

			// Process ns specific data within a given month
			nsRecord = append(nsRecord, &activity.MonthlyNamespaceRecord{
				NamespaceID: nsID,
				Counts: &activity.CountsRecord{
					EntityClients:    len(nsData.Counts.Entities),
					NonEntityClients: int(nsData.Counts.Tokens) + len(nsData.Counts.NonEntities),
				},
				Mounts: mountRecord,
			})
		}
		return nsRecord
	}
	for timestamp, monthData := range byMonth {
		newClientsNSRecord := processByNamespaces(monthData.NewClients.Namespaces)
		newClientRecord := &activity.NewClientRecord{
			Counts: &activity.CountsRecord{
				EntityClients:    len(monthData.NewClients.Counts.Entities),
				NonEntityClients: int(monthData.NewClients.Counts.Tokens) + len(monthData.NewClients.Counts.NonEntities),
			},
			Namespaces: newClientsNSRecord,
		}

		// Process all the months
		monthly = append(monthly, &activity.MonthRecord{
			Timestamp: timestamp,
			Counts: &activity.CountsRecord{
				EntityClients:    len(monthData.Counts.Entities),
				NonEntityClients: int(monthData.Counts.Tokens) + len(monthData.Counts.NonEntities),
			},
			Namespaces: processByNamespaces(monthData.Namespaces),
			NewClients: newClientRecord,
		})
	}
	return monthly
}

func (a *ActivityLog) calculateByNamespaceResponseForQuery(ctx context.Context, byNamespace []*activity.NamespaceRecord) (int, int, []*ResponseNamespace, error) {
	queryNS, err := namespace.FromContext(ctx)
	if err != nil {
		return 0, 0, nil, err
	}
	byNamespaceResponse := make([]*ResponseNamespace, 0)
	totalEntities := 0
	totalTokens := 0

	for _, nsRecord := range byNamespace {
		ns, err := NamespaceByID(ctx, nsRecord.NamespaceID, a.core)
		if err != nil {
			return 0, 0, nil, err
		}
		if a.includeInResponse(queryNS, ns) {
			mountResponse := make([]*ResponseMount, 0, len(nsRecord.Mounts))
			for _, mountRecord := range nsRecord.Mounts {
				mountResponse = append(mountResponse, &ResponseMount{
					MountPath: mountRecord.MountPath,
					Counts: &ResponseCounts{
						DistinctEntities: int(mountRecord.Counts.EntityClients),
						EntityClients:    int(mountRecord.Counts.EntityClients),
						NonEntityClients: int(mountRecord.Counts.NonEntityClients),
						NonEntityTokens:  int(mountRecord.Counts.NonEntityClients),
						Clients:          int(mountRecord.Counts.EntityClients + mountRecord.Counts.NonEntityClients),
					},
				})
			}
			// Sort the mounts in descending order of usage
			sort.Slice(mountResponse, func(i, j int) bool {
				return mountResponse[i].Counts.Clients > mountResponse[j].Counts.Clients
			})

			var displayPath string
			if ns == nil {
				displayPath = fmt.Sprintf("deleted namespace %q", nsRecord.NamespaceID)
			} else {
				displayPath = ns.Path
			}
			byNamespaceResponse = append(byNamespaceResponse, &ResponseNamespace{
				NamespaceID:   nsRecord.NamespaceID,
				NamespacePath: displayPath,
				Counts: ResponseCounts{
					DistinctEntities: int(nsRecord.Entities),
					EntityClients:    int(nsRecord.Entities),
					NonEntityTokens:  int(nsRecord.NonEntityTokens),
					NonEntityClients: int(nsRecord.NonEntityTokens),
					Clients:          int(nsRecord.Entities + nsRecord.NonEntityTokens),
				},
				Mounts: mountResponse,
			})
			totalEntities += int(nsRecord.Entities)
			totalTokens += int(nsRecord.NonEntityTokens)
		}
	}
	return totalEntities, totalTokens, byNamespaceResponse, nil
}

func (a *ActivityLog) prepareMonthsResponseForQuery(ctx context.Context, byMonth []*activity.MonthRecord) ([]*ResponseMonth, error) {
	months := make([]*ResponseMonth, 0, len(byMonth))
	for _, monthsRecord := range byMonth {
		newClientsResponse := &ResponseNewClients{}
		if int(monthsRecord.NewClients.Counts.EntityClients+monthsRecord.NewClients.Counts.NonEntityClients) != 0 {
			newClientsNSResponse, err := a.prepareNamespaceResponse(ctx, monthsRecord.NewClients.Namespaces)
			if err != nil {
				return nil, err
			}
			newClientsResponse.Counts = &ResponseCounts{
				EntityClients:    int(monthsRecord.NewClients.Counts.EntityClients),
				NonEntityClients: int(monthsRecord.NewClients.Counts.NonEntityClients),
				Clients:          int(monthsRecord.NewClients.Counts.EntityClients + monthsRecord.NewClients.Counts.NonEntityClients),
			}
			newClientsResponse.Namespaces = newClientsNSResponse
		}

		monthResponse := &ResponseMonth{
			Timestamp: time.Unix(monthsRecord.Timestamp, 0).UTC().Format(time.RFC3339),
		}
		if int(monthsRecord.Counts.EntityClients+monthsRecord.Counts.NonEntityClients) != 0 {
			nsResponse, err := a.prepareNamespaceResponse(ctx, monthsRecord.Namespaces)
			if err != nil {
				return nil, err
			}
			monthResponse.Counts = &ResponseCounts{
				EntityClients:    int(monthsRecord.Counts.EntityClients),
				NonEntityClients: int(monthsRecord.Counts.NonEntityClients),
				Clients:          int(monthsRecord.Counts.EntityClients + monthsRecord.Counts.NonEntityClients),
			}
			monthResponse.Namespaces = nsResponse
			monthResponse.NewClients = newClientsResponse
			months = append(months, monthResponse)
		}
	}
	return months, nil
}

// prepareNamespaceResponse populates the namespace portion of the activity log response struct
// from
func (a *ActivityLog) prepareNamespaceResponse(ctx context.Context, nsRecords []*activity.MonthlyNamespaceRecord) ([]*ResponseNamespace, error) {
	queryNS, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	nsResponse := make([]*ResponseNamespace, 0, len(nsRecords))
	for _, nsRecord := range nsRecords {
		if int(nsRecord.Counts.EntityClients) == 0 && int(nsRecord.Counts.NonEntityClients) == 0 {
			continue
		}

		ns, err := NamespaceByID(ctx, nsRecord.NamespaceID, a.core)
		if err != nil {
			return nil, err
		}
		if a.includeInResponse(queryNS, ns) {
			mountResponse := make([]*ResponseMount, 0, len(nsRecord.Mounts))
			for _, mountRecord := range nsRecord.Mounts {
				if int(mountRecord.Counts.EntityClients) == 0 && int(mountRecord.Counts.NonEntityClients) == 0 {
					continue
				}

				mountResponse = append(mountResponse, &ResponseMount{
					MountPath: mountRecord.MountPath,
					Counts: &ResponseCounts{
						EntityClients:    int(mountRecord.Counts.EntityClients),
						NonEntityClients: int(mountRecord.Counts.NonEntityClients),
						Clients:          int(mountRecord.Counts.EntityClients + mountRecord.Counts.NonEntityClients),
					},
				})
			}

			var displayPath string
			if ns == nil {
				displayPath = fmt.Sprintf("deleted namespace %q", nsRecord.NamespaceID)
			} else {
				displayPath = ns.Path
			}
			nsResponse = append(nsResponse, &ResponseNamespace{
				NamespaceID:   nsRecord.NamespaceID,
				NamespacePath: displayPath,
				Counts: ResponseCounts{
					EntityClients:    int(nsRecord.Counts.EntityClients),
					NonEntityClients: int(nsRecord.Counts.NonEntityClients),
					Clients:          int(nsRecord.Counts.EntityClients + nsRecord.Counts.NonEntityClients),
				},
				Mounts: mountResponse,
			})
		}
	}
	return nsResponse, nil
}

// partialMonthClientCount returns the number of clients used so far this month.
// If activity log is not enabled, the response will be nil
func (a *ActivityLog) partialMonthClientCount(ctx context.Context) (map[string]interface{}, error) {
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()

	if !a.enabled {
		// nothing to count
		return nil, nil
	}

	// Traverse through current month's activitylog data and group clients
	// into months and namespaces
	byMonth, byNamespace := a.populateNamespaceAndMonthlyBreakdowns()

	// Convert the byNamespace breakdowns into structs that are
	// consumable by the /activity endpoint, so as to reuse code between these two
	// endpoints.
	byNamespaceComputation := a.transformALNamespaceBreakdowns(byNamespace)

	// Calculate the namespace response breakdowns and totals for entities and tokens from the initial
	// namespace data.
	totalEntities, totalTokens, byNamespaceResponse, err := a.calculateByNamespaceResponseForQuery(ctx, byNamespaceComputation)
	if err != nil {
		return nil, err
	}

	// Sort clients within each namespace
	a.sortALResponseNamespaces(byNamespaceResponse)

	// Now populate the response based on breakdowns.
	responseData := make(map[string]interface{})
	responseData["by_namespace"] = byNamespaceResponse
	responseData["distinct_entities"] = totalEntities
	responseData["entity_clients"] = totalEntities
	responseData["non_entity_tokens"] = totalTokens
	responseData["non_entity_clients"] = totalTokens
	responseData["clients"] = totalEntities + totalTokens

	// The partialMonthClientCount should not have more than one month worth of data.
	// If it does, something has gone wrong and we should warn that the activity log data
	// might be inaccurate.
	if len(byMonth) != 1 {
		monthTimestamps := make([]string, 0)
		for timestamp := range byMonth {
			dateTimeString := time.Unix(timestamp, 0).UTC().Format(time.RFC3339)
			monthTimestamps = append(monthTimestamps, dateTimeString)
		}
		a.logger.Error("more or less than one month of data recorded in current month's activity log", "timestamps", monthTimestamps)
	}

	// Convert the byMonth breakdowns into structs that are
	// consumable by the /activity endpoint, so as to reuse code between these two
	// endpoints.
	monthlyComputation := a.transformMonthBreakdowns(byMonth)

	// Create and populate the month response structs based on the monthly breakdown.
	months, err := a.prepareMonthsResponseForQuery(ctx, monthlyComputation)
	if err != nil {
		return nil, err
	}

	// Sort the months and clients within each month before adding the months to the response
	a.sortActivityLogMonthsResponse(months)
	responseData["months"] = months

	return responseData, nil
}

func (a *ActivityLog) writeExport(ctx context.Context, rw http.ResponseWriter, format string, startTime, endTime time.Time) error {
	// For capacity reasons only allow a single in-process export at a time.
	// TODO do we really need to do this?
	if !a.inprocessExport.CAS(false, true) {
		return fmt.Errorf("existing export in progress")
	}
	defer a.inprocessExport.Store(false)

	// Find the months with activity log data that are between the start and end
	// months. We want to walk this in cronological order so the oldest instance of a
	// client usage is recorded, not the most recent.
	times, err := a.availableLogs(ctx)
	if err != nil {
		a.logger.Warn("failed to list available log segments", "error", err)
		return fmt.Errorf("failed to list available log segments: %w", err)
	}
	sort.Slice(times, func(i, j int) bool {
		// sort in chronological order to produce the output we want showing what
		// month an entity first had activity.
		return times[i].Before(times[j])
	})

	// Filter over just the months we care about
	filteredList := make([]time.Time, 0, len(times))
	for _, t := range times {
		if timeutil.InRange(t, startTime, endTime) {
			filteredList = append(filteredList, t)
		}
	}
	if len(filteredList) == 0 {
		a.logger.Info("no data to export", "start_time", startTime, "end_time", endTime)
		return fmt.Errorf("no data to export in provided time range")
	}

	actualStartTime := filteredList[len(filteredList)-1]
	a.logger.Trace("choose start time for export", "actualStartTime", actualStartTime, "months_included", filteredList)

	// Add headers here because we start to immediately write in the csv encoder
	// constructor.
	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"activity_export_%d_to_%d.%s\"", actualStartTime.Unix(), endTime.Unix(), format))
	rw.Header().Add("Content-Type", fmt.Sprintf("application/%s", format))

	var encoder encoder
	switch format {
	case "json":
		encoder = newJSONEncoder(rw)
	case "csv":
		var err error
		encoder, err = newCSVEncoder(rw)
		if err != nil {
			return fmt.Errorf("failed to create csv encoder: %w", err)
		}
	default:
		return fmt.Errorf("invalid format: %s", format)
	}

	a.logger.Info("starting activity log export", "start_time", startTime, "end_time", endTime, "format", format)

	dedupedIds := make(map[string]struct{})

	walkEntities := func(l *activity.EntityActivityLog, startTime time.Time, hll *hyperloglog.Sketch) error {
		for _, e := range l.Clients {
			if _, ok := dedupedIds[e.ClientID]; ok {
				continue
			}

			dedupedIds[e.ClientID] = struct{}{}
			err := encoder.Encode(e)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// For each month in the filtered list walk all the log segments

	for _, startTime := range filteredList {
		err := a.WalkEntitySegments(ctx, startTime, nil, walkEntities)
		if err != nil {
			a.logger.Error("failed to load segments for export", "error", err)
			return fmt.Errorf("failed to load segments for export: %w", err)
		}
	}

	// Flush and error check the encoder. This is neccessary for buffered
	// encoders like csv.
	encoder.Flush()
	if err := encoder.Error(); err != nil {
		a.logger.Error("failed to flush export encoding", "error", err)
		return fmt.Errorf("failed to flush export encoding: %w", err)
	}

	return nil
}

type encoder interface {
	Encode(*activity.EntityRecord) error
	Flush()
	Error() error
}

var _ encoder = (*jsonEncoder)(nil)

type jsonEncoder struct {
	e *json.Encoder
}

func newJSONEncoder(w io.Writer) *jsonEncoder {
	return &jsonEncoder{
		e: json.NewEncoder(w),
	}
}

func (j *jsonEncoder) Encode(er *activity.EntityRecord) error {
	return j.e.Encode(er)
}

// Flush is a no-op because json.Encoder doesn't buffer data
func (j *jsonEncoder) Flush() {}

// Error is a no-op because flushing is a no-op.
func (j *jsonEncoder) Error() error { return nil }

var _ encoder = (*csvEncoder)(nil)

type csvEncoder struct {
	*csv.Writer
}

func newCSVEncoder(w io.Writer) (*csvEncoder, error) {
	writer := csv.NewWriter(w)

	err := writer.Write([]string{
		"client_id",
		"namespace_id",
		"timestamp",
		"non_entity",
		"mount_accessor",
	})
	if err != nil {
		return nil, err
	}

	return &csvEncoder{
		Writer: writer,
	}, nil
}

// Encode converts an export bundle into a set of strings and writes them to the
// csv writer.
func (c *csvEncoder) Encode(e *activity.EntityRecord) error {
	return c.Writer.Write([]string{
		e.ClientID,
		e.NamespaceID,
		fmt.Sprintf("%d", e.Timestamp),
		fmt.Sprintf("%t", e.NonEntity),
		e.MountAccessor,
	})
}
