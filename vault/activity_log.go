package vault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/mitchellh/copystructure"
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

	// Number of entity records to store per segment
	// Estimated as 512KiB / 64 bytes = 8192, rounded down
	activitySegmentEntityCapacity = 8000

	// Maximum number of segments per month
	activityLogMaxSegmentPerMonth = 81

	// Number of records (entity or token) to store in a
	// standby fragment before sending it to the active node.
	// Estimates as 8KiB / 64 bytes = 128
	activityFragmentStandbyCapacity = 128

	// Delimiter between the string fields used to generate a client
	// ID for tokens without entities. This is the 0 character, which
	// is a non-printable string. Please see unicode.IsPrint for details.
	clientIDTWEDelimiter = rune('\x00')

	// Delimiter between each policy in the sorted policies used to
	// generate a client ID for tokens without entities. This is the 127
	// character, which is a non-printable string. Please see unicode.IsPrint
	// for details.
	sortedPoliciesTWEDelimiter = rune('\x7F')

	// trackedTWESegmentPeriod is a time period of a little over a month, and represents
	// the amount of time that needs to pass after a 1.9 or later upgrade to result in
	// all fragments and segments no longer storing token counts in the directtokens
	// storage path.
	trackedTWESegmentPeriod = 35 * 24
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

type clients struct {
	distinctEntities    uint64
	distinctNonEntities uint64
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

	// fragmentLock protects enable, activeClients, fragment, standbyFragmentsReceived
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

	// Channel for writing fragment immediately
	writeCh chan struct{}

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
	retentionDone chan struct{}

	// for testing: is config currently being invalidated. protected by l
	configInvalidationInProgress bool

	// clientTracker tracks active clients this month.  Protected by fragmentLock.
	clientTracker *ClientTracker
}

type ClientTracker struct {
	// All known active clients this month; use fragmentLock read-locked
	// to check whether it already exists.
	activeClients               map[string]struct{}
	entityCountByNamespaceID    map[string]uint64
	nonEntityCountByNamespaceID map[string]uint64
}

// These non-persistent configuration options allow us to disable
// parts of the implementation for integration testing.
// The default values should turn everything on.
type ActivityLogCoreConfig struct {
	// Enable activity log even if the feature flag not set
	ForceEnable bool

	// Do not start timers to send or persist fragments.
	DisableTimers bool
}

// NewActivityLog creates an activity log.
func NewActivityLog(core *Core, logger log.Logger, view *BarrierView, metrics metricsutil.Metrics) (*ActivityLog, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	a := &ActivityLog{
		core:            core,
		configOverrides: &core.activityLogConfig,
		logger:          logger,
		view:            view,
		metrics:         metrics,
		nodeID:          hostname,
		newFragmentCh:   make(chan struct{}, 1),
		sendCh:          make(chan struct{}, 1), // buffered so it can be triggered by fragment size
		writeCh:         make(chan struct{}, 1), // same for full segment
		doneCh:          make(chan struct{}, 1),
		clientTracker: &ClientTracker{
			activeClients:               make(map[string]struct{}),
			entityCountByNamespaceID:    make(map[string]uint64),
			nonEntityCountByNamespaceID: make(map[string]uint64),
		},
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
		time.Now(), []metricsutil.Label{})

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
		for _, e := range f.Clients {
			// We could sort by timestamp to see which is first.
			// We'll ignore that; the order of the append above means
			// that we choose entries in localFragment over those
			// from standby nodes.
			newEntities[e.ClientID] = e
			saveChanges = true
		}
		// As of 1.9, a fragment should no longer have any NonEntityTokens. However
		// in order to not lose any information about the current segment during the
		// month when the client upgrades to 1.9, we must retain this functionality.
		for ns, val := range f.NonEntityTokens {
			// We track these pre-1.9 values in the old location, which is
			// a.currentSegment.tokenCount, as opposed to the counter that stores tokens
			// without entities that have client IDs, namely
			// a.clientTracker.nonEntityCountByNamespaceID. This preserves backward
			// compatibility for the precomputedQueryWorkers and the segment storing
			// logic.
			a.currentSegment.tokenCount.CountByNamespaceID[ns] += val
			saveChanges = true
		}
	}

	if !saveChanges {
		return nil
	}

	// Will all new entities fit?  If not, roll over to a new segment.
	available := activitySegmentEntityCapacity - len(a.currentSegment.currentClients.Clients)
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
		if len(excessClients) > activitySegmentEntityCapacity {
			a.logger.Warn("too many new active entities, dropping tail", "entities", len(excessClients))
			excessClients = excessClients[:activitySegmentEntityCapacity]
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
	entityPath := fmt.Sprintf("log/entity/%d/%d", a.currentSegment.startTimestamp, a.currentSegment.clientSequenceNumber)
	// RFC (VLT-120) defines this as 1-indexed, but it should be 0-indexed
	tokenPath := fmt.Sprintf("log/directtokens/%d/0", a.currentSegment.startTimestamp)

	for _, client := range a.currentSegment.currentClients.Clients {
		// Explicitly catch and throw clear error message if client ID creation and storage
		// results in a []byte that doesn't assert into a valid string.
		if !utf8.ValidString(client.ClientID) {
			return fmt.Errorf("client ID %q is not a valid string:", client.ClientID)
		}
	}

	if len(a.currentSegment.currentClients.Clients) > 0 || force {
		entities, err := proto.Marshal(a.currentSegment.currentClients)
		if err != nil {
			return err
		}

		a.logger.Trace("writing segment", "path", entityPath)
		err = a.view.Put(ctx, &logical.StorageEntry{
			Key:   entityPath,
			Value: entities,
		})
		if err != nil {
			return err
		}
	}

	// We must still allow for the tokenCount of the current segment to
	// be written to storage, since if we remove this code we will incur
	// data loss for one segment's worth of TWEs.
	if len(a.currentSegment.tokenCount.CountByNamespaceID) > 0 || force {
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
			return err
		}

		a.logger.Trace("writing segment", "path", tokenPath)
		err = a.view.Put(ctx, &logical.StorageEntry{
			Key:   tokenPath,
			Value: tokenCount,
		})
		if err != nil {
			return err
		}
	}
	return nil
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
func (a *ActivityLog) WalkEntitySegments(ctx context.Context,
	startTime time.Time,
	walkFn func(*activity.EntityActivityLog)) error {
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
		walkFn(out)
	}
	return nil
}

// WalkTokenSegments loads each of the token segments (expected 1) for a particular start time
func (a *ActivityLog) WalkTokenSegments(ctx context.Context,
	startTime time.Time,
	walkFn func(*activity.TokenCount)) error {
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
			a.clientTracker.addClient(ent)
		}
	}
	a.fragmentLock.Unlock()
	a.l.RUnlock()

	return nil
}

// loadCurrentClientSegment loads the most recent segment (for "this month") into memory
// (to append new entries), and to the activeClients to avoid duplication
// call with fragmentLock and l held
func (a *ActivityLog) loadCurrentClientSegment(ctx context.Context, startTime time.Time, sequenceNum uint64) error {
	path := activityEntityBasePath + fmt.Sprint(startTime.Unix()) + "/" + strconv.FormatUint(sequenceNum, 10)
	data, err := a.view.Get(ctx, path)
	if err != nil {
		return err
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

	for _, ent := range out.Clients {
		a.clientTracker.addClient(ent)
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
	a.clientTracker.activeClients = make(map[string]struct{})
	a.clientTracker.entityCountByNamespaceID = make(map[string]uint64)
	a.clientTracker.nonEntityCountByNamespaceID = make(map[string]uint64)
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
		a.startNewCurrentLogLocked(time.Now().UTC())
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

	// check for segments out of retention period, if it has changed
	go a.retentionWorker(ctx, time.Now(), a.retentionMonths)
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
	logger := c.baseLogger.Named("activity")
	c.AddLogger(logger)

	if os.Getenv("VAULT_DISABLE_ACTIVITY_LOG") != "" {
		logger.Info("activity log disabled via environment variable")
		return nil
	}

	view := c.systemBarrierView.SubView(activitySubPath)

	manager, err := NewActivityLog(c, logger, view, c.metricSink)
	if err != nil {
		return err
	}
	c.activityLog = manager

	// load activity log for "this month" into memory
	err = manager.refreshFromStoredLog(manager.core.activeContext, wg, time.Now().UTC())
	if err != nil {
		return err
	}

	// Start the background worker, depending on type
	// Lock already held here, can't use .PerfStandby()
	// The workers need to know the current segment time.
	if c.perfStandby {
		go manager.perfStandbyFragmentWorker(ctx)
	} else {
		go manager.activeFragmentWorker(ctx)

		// Check for any intent log, in the background
		go manager.precomputedQueryWorker(ctx)

		// Catch up on garbage collection
		// Signal when this is done so that unit tests can proceed.
		manager.retentionDone = make(chan struct{})
		go func() {
			manager.retentionWorker(ctx, time.Now(), manager.retentionMonths)
			close(manager.retentionDone)
		}()
	}

	return nil
}

// stopActivityLog removes the ActivityLog from Core
// and frees any resources.
func (c *Core) stopActivityLog() {
	// preSeal may run before startActivityLog got a chance to complete.
	if c.activityLog != nil {
		// Shut down background worker
		close(c.activityLog.doneCh)
	}

	c.activityLog = nil
}

func (a *ActivityLog) StartOfNextMonth() time.Time {
	a.l.RLock()
	defer a.l.RUnlock()
	var segmentStart time.Time
	if a.currentSegment.startTimestamp == 0 {
		segmentStart = time.Now().UTC()
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
	timer := time.NewTimer(time.Duration(0))
	fragmentWaiting := false
	// Eat first event, so timer is stopped
	<-timer.C

	endOfMonth := time.NewTimer(a.StartOfNextMonth().Sub(time.Now()))
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
			a.clientTracker.activeClients = make(map[string]struct{})
			a.clientTracker.entityCountByNamespaceID = make(map[string]uint64)
			a.clientTracker.nonEntityCountByNamespaceID = make(map[string]uint64)
			a.fragmentLock.Unlock()

			// Set timer for next month.
			// The current segment *probably* hasn't been set yet (via invalidation),
			// so don't rely on it.
			target := timeutil.StartOfNextMonth(time.Now().UTC())
			endOfMonth.Reset(target.Sub(time.Now()))
		}
	}
}

// activeFragmentWorker handles scheduling the write of the next
// segment.  It runs on active nodes only.
func (a *ActivityLog) activeFragmentWorker(ctx context.Context) {
	ticker := time.NewTicker(activitySegmentInterval)

	endOfMonth := time.NewTimer(a.StartOfNextMonth().Sub(time.Now()))
	if a.configOverrides.DisableTimers {
		endOfMonth.Stop()
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
		case <-a.writeCh:
			a.logger.Trace("writing segment on request")
			writeFunc()

			// Reset the schedule to wait 10 minutes from this forced write.
			ticker.Stop()
			ticker = time.NewTicker(activitySegmentInterval)

			// Simpler, but ticker.Reset was introduced in go 1.15:
			// ticker.Reset(activitySegmentInterval)
		case currentTime := <-endOfMonth.C:
			err := a.HandleEndOfMonth(ctx, currentTime.UTC())
			if err != nil {
				a.logger.Error("failed to perform end of month rotation", "error", err)
			}

			// Garbage collect any segments or queries based on the immediate
			// value of retentionMonths.
			a.l.RLock()
			go a.retentionWorker(ctx, currentTime.UTC(), a.retentionMonths)
			a.l.RUnlock()

			delta := a.StartOfNextMonth().Sub(time.Now())
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

	prevSegmentTimestamp := a.currentSegment.startTimestamp
	nextSegmentTimestamp := timeutil.StartOfMonth(currentTime.UTC()).Unix()

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
	a.AddClientToFragment(entityID, namespaceID, timestamp, false)
}

// AddClientToFragment checks a client ID for uniqueness and
// if not already present, adds it to the current fragment.
// The timestamp is a Unix timestamp *without* nanoseconds, as that
// is what token.CreationTime uses.
func (a *ActivityLog) AddClientToFragment(clientID string, namespaceID string, timestamp int64, isTWE bool) {
	// Check whether entity ID already recorded
	var present bool

	a.fragmentLock.RLock()
	if a.enabled {
		_, present = a.clientTracker.activeClients[clientID]
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
	_, present = a.clientTracker.activeClients[clientID]
	if present {
		return
	}

	a.createCurrentFragment()

	clientRecord := &activity.EntityRecord{
		ClientID:    clientID,
		NamespaceID: namespaceID,
		Timestamp:   timestamp,
	}

	// Track whether the clientID corresponds to a token without an entity or not.
	// This field is backward compatible, as the default is 0, so records created
	// from pre-1.9 activityLog code will automatically be marked as having an entity.
	if isTWE {
		clientRecord.NonEntity = true
	}

	a.fragment.Clients = append(a.fragment.Clients, clientRecord)
	a.clientTracker.addClient(clientRecord)
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
		a.fragmentCreation = time.Now().UTC()

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
		a.clientTracker.addClient(e)
	}

	a.standbyFragmentsReceived = append(a.standbyFragmentsReceived, fragment)

	// TODO: check if current segment is full and should be written
}

type ClientCountResponse struct {
	DistinctEntities int `json:"distinct_entities"`
	NonEntityTokens  int `json:"non_entity_tokens"`
	Clients          int `json:"clients"`
}

type ClientCountInNamespace struct {
	NamespaceID   string              `json:"namespace_id"`
	NamespacePath string              `json:"namespace_path"`
	Counts        ClientCountResponse `json:"counts"`
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

func (a *ActivityLog) handleQuery(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	queryNS, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	pq, err := a.queryStore.Get(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}
	if pq == nil {
		return nil, nil
	}

	responseData := make(map[string]interface{})
	responseData["start_time"] = pq.StartTime.Format(time.RFC3339)
	responseData["end_time"] = pq.EndTime.Format(time.RFC3339)
	byNamespace := make([]*ClientCountInNamespace, 0)

	totalEntities := 0
	totalTokens := 0

	for _, nsRecord := range pq.Namespaces {
		ns, err := NamespaceByID(ctx, nsRecord.NamespaceID, a.core)
		if err != nil {
			return nil, err
		}
		if a.includeInResponse(queryNS, ns) {
			var displayPath string
			if ns == nil {
				displayPath = fmt.Sprintf("deleted namespace %q", nsRecord.NamespaceID)
			} else {
				displayPath = ns.Path
			}
			byNamespace = append(byNamespace, &ClientCountInNamespace{
				NamespaceID:   nsRecord.NamespaceID,
				NamespacePath: displayPath,
				Counts: ClientCountResponse{
					DistinctEntities: int(nsRecord.Entities),
					NonEntityTokens:  int(nsRecord.NonEntityTokens),
					Clients:          int(nsRecord.Entities + nsRecord.NonEntityTokens),
				},
			})
			totalEntities += int(nsRecord.Entities)
			totalTokens += int(nsRecord.NonEntityTokens)
		}
	}

	sort.Slice(byNamespace, func(i, j int) bool {
		return byNamespace[i].Counts.Clients > byNamespace[j].Counts.Clients
	})

	responseData["by_namespace"] = byNamespace
	responseData["total"] = &ClientCountResponse{
		DistinctEntities: totalEntities,
		NonEntityTokens:  totalTokens,
		Clients:          totalEntities + totalTokens,
	}
	return responseData, nil
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
// and returns the corresponding Client ID.
// This currently occurs on token usage only.
func (a *ActivityLog) HandleTokenUsage(entry *logical.TokenEntry) string {
	// First, check if a is enabled, so as to avoid the cost of creating an ID for
	// tokens without entities in the case where it not.
	a.fragmentLock.RLock()
	if !a.enabled {
		a.fragmentLock.RUnlock()
		return ""
	}
	a.fragmentLock.RUnlock()

	// Do not count wrapping tokens in client count
	if IsWrappingToken(entry) {
		return ""
	}

	// Do not count root tokens in client count.
	if entry.IsRoot() {
		return ""
	}

	// Parse an entry's client ID and add it to the activity log
	clientID, isTWE := a.CreateClientID(entry)
	a.AddClientToFragment(clientID, entry.NamespaceID, entry.CreationTime, isTWE)
	return clientID
}

// CreateClientID returns the client ID, and a boolean which is false if the clientID
// has an entity, and true otherwise
func (a *ActivityLog) CreateClientID(entry *logical.TokenEntry) (string, bool) {
	var clientIDInputBuilder strings.Builder

	// if entry has an associated entity ID, return it
	if entry.EntityID != "" {
		return entry.EntityID, false
	}

	// The entry is associated with a TWE (token without entity). In this case
	// we must create a client ID by calculating the following formula:
	// clientID = SHA256(sorted policies + namespace)

	// Step 1: Copy entry policies to a new struct
	sortedPolicies := make([]string, len(entry.Policies))
	copy(sortedPolicies, entry.Policies)

	// Step 2: Sort and join copied policies
	sort.Strings(sortedPolicies)
	for _, pol := range sortedPolicies {
		clientIDInputBuilder.WriteRune(sortedPoliciesTWEDelimiter)
		clientIDInputBuilder.WriteString(pol)
	}

	// Step 3: Add namespace ID
	clientIDInputBuilder.WriteRune(clientIDTWEDelimiter)
	clientIDInputBuilder.WriteString(entry.NamespaceID)

	if clientIDInputBuilder.Len() == 0 {
		a.logger.Error("vault token with no entity ID, policies, or namespace was recorded " +
			"in the activity log")
		return "", true
	}
	// Step 4: Remove the first character in the string, as it's an unnecessary delimiter
	clientIDInput := clientIDInputBuilder.String()[1:]

	// Step 5: Hash the sum
	hashed := sha256.Sum256([]byte(clientIDInput))
	return base64.StdEncoding.EncodeToString(hashed[:]), true
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
	// Base retention period on the month we are generating (even in the past)--- time.Now()
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

	times, err := a.getMostRecentActivityLogSegment(ctx)
	if err != nil {
		a.logger.Warn("could not list recent segments", "error", err)
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

	// "times" is already in reverse order, start building the per-namespace maps
	// from the last month backward

	type NamespaceCounts struct {
		// entityID -> present
		Entities map[string]struct{}
		// count. This exists for backward compatibility
		Tokens uint64
		// clientID -> present
		NonEntities map[string]struct{}
	}
	byNamespace := make(map[string]*NamespaceCounts)

	createNs := func(namespaceID string) {
		if _, namespacePresent := byNamespace[namespaceID]; !namespacePresent {
			byNamespace[namespaceID] = &NamespaceCounts{
				Entities:    make(map[string]struct{}),
				Tokens:      0,
				NonEntities: make(map[string]struct{}),
			}
		}
	}

	walkEntities := func(l *activity.EntityActivityLog) {
		for _, e := range l.Clients {
			createNs(e.NamespaceID)
			if e.NonEntity == true {
				byNamespace[e.NamespaceID].NonEntities[e.ClientID] = struct{}{}
			} else {
				byNamespace[e.NamespaceID].Entities[e.ClientID] = struct{}{}
			}
		}
	}
	walkTokens := func(l *activity.TokenCount) {
		for nsID, v := range l.CountByNamespaceID {
			createNs(nsID)
			byNamespace[nsID].Tokens += v
		}
	}

	endTime := timeutil.EndOfMonth(time.Unix(lastMonth, 0).UTC())
	activePeriodStart := timeutil.MonthsPreviousTo(a.defaultReportMonths, endTime)
	// If not enough data, report as much as we have in the window
	if activePeriodStart.Before(times[len(times)-1]) {
		activePeriodStart = times[len(times)-1]
	}

	for _, startTime := range times {
		// Do not work back further than the current retention window,
		// which will just get deleted anyway.
		if startTime.Before(retentionWindow) {
			break
		}

		err = a.WalkEntitySegments(ctx, startTime, walkEntities)
		if err != nil {
			a.logger.Warn("failed to load previous segments", "error", err)
			return err
		}
		err = a.WalkTokenSegments(ctx, startTime, walkTokens)
		if err != nil {
			a.logger.Warn("failed to load previous token counts", "error", err)
			return err
		}

		// Save the work to date in a record
		pq := &activity.PrecomputedQuery{
			StartTime:  startTime,
			EndTime:    endTime,
			Namespaces: make([]*activity.NamespaceRecord, 0, len(byNamespace)),
		}

		for nsID, counts := range byNamespace {
			pq.Namespaces = append(pq.Namespaces, &activity.NamespaceRecord{
				NamespaceID:     nsID,
				Entities:        uint64(len(counts.Entities)),
				NonEntityTokens: counts.Tokens + uint64(len(counts.NonEntities)),
			})

			// If this is the most recent month, or the start of the reporting period, output
			// a metric for each namespace.
			if startTime == times[0] {
				a.metrics.SetGaugeWithLabels(
					[]string{"identity", "entity", "active", "monthly"},
					float32(len(counts.Entities)),
					[]metricsutil.Label{
						{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
					},
				)
				a.metrics.SetGaugeWithLabels(
					[]string{"identity", "nonentity", "active", "monthly"},
					float32(len(counts.NonEntities))+float32(counts.Tokens),
					[]metricsutil.Label{
						{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
					},
				)
			} else if startTime == activePeriodStart {
				a.metrics.SetGaugeWithLabels(
					[]string{"identity", "entity", "active", "reporting_period"},
					float32(len(counts.Entities)),
					[]metricsutil.Label{
						{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
					},
				)
				a.metrics.SetGaugeWithLabels(
					[]string{"identity", "nonentity", "active", "reporting_period"},
					float32(len(counts.NonEntities))+float32(counts.Tokens),
					[]metricsutil.Label{
						{Name: "namespace", Value: a.namespaceToLabel(ctx, nsID)},
					},
				)
			}
		}

		err = a.queryStore.Put(ctx, pq)
		if err != nil {
			a.logger.Warn("failed to store precomputed query", "error", err)
		}
	}

	// delete the intent log
	a.view.Delete(ctx, activityIntentLogKey)

	a.logger.Info("finished computing queries", "month", endTime)

	return nil
}

// goroutine to delete any segments or precomputed queries older than
// the retention period.
// We expect the return value won't be checked, so log errors as they occur
// (but for unit testing having the error return should help.)
func (a *ActivityLog) retentionWorker(ctx context.Context, currentTime time.Time, retentionMonths int) error {
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
	count := len(a.clientTracker.activeClients)

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

// partialMonthClientCount returns the number of clients used so far this month.
// If activity log is not enabled, the response will be nil
func (a *ActivityLog) partialMonthClientCount(ctx context.Context) (map[string]interface{}, error) {
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()

	if !a.enabled {
		// nothing to count
		return nil, nil
	}
	byNamespace := make([]*ClientCountInNamespace, 0)
	responseData := make(map[string]interface{})
	totalEntities := 0
	totalTokens := 0
	nonEntityTokensMapInterface, err := copystructure.Copy(a.clientTracker.nonEntityCountByNamespaceID)
	if err != nil {
		return nil, fmt.Errorf("error making deep copy of nonEntityCounts: %+w", err)
	}
	nonEntityTokensMap := nonEntityTokensMapInterface.(map[string]uint64)
	// Merge the tokenCounts created pre-clientID with the newly counted
	// clientID tokens, if tokenCounts exist.
	for nsID, count := range a.currentSegment.tokenCount.CountByNamespaceID {
		nonEntityTokensMap[nsID] += count
	}
	clientCountTable := createClientCountTable(a.clientTracker.entityCountByNamespaceID, nonEntityTokensMap)
	queryNS, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	for nsID, clients := range clientCountTable {
		ns, err := NamespaceByID(ctx, nsID, a.core)
		if err != nil {
			return nil, err
		}

		// Only include namespaces that are the queryNS or within it.  If queryNS is the
		// root namespace, include all namespaces, even those which have been deleted.
		if a.includeInResponse(queryNS, ns) {
			var displayPath string
			if ns == nil {
				displayPath = fmt.Sprintf("deleted namespace %q", nsID)
			} else {
				displayPath = ns.Path
			}

			byNamespace = append(byNamespace, &ClientCountInNamespace{
				NamespaceID:   nsID,
				NamespacePath: displayPath,
				Counts: ClientCountResponse{
					DistinctEntities: int(clients.distinctEntities),
					NonEntityTokens:  int(clients.distinctNonEntities),
					Clients:          int(clients.distinctEntities + clients.distinctNonEntities),
				},
			})

			totalEntities += int(clients.distinctEntities)
			totalTokens += int(clients.distinctNonEntities)

		}
	}

	sort.Slice(byNamespace, func(i, j int) bool {
		return byNamespace[i].NamespaceID < byNamespace[j].NamespaceID
	})

	responseData["by_namespace"] = byNamespace
	responseData["distinct_entities"] = totalEntities
	responseData["non_entity_tokens"] = totalTokens
	responseData["clients"] = totalEntities + totalTokens

	return responseData, nil
}

// createClientCountTable maps the entitycount and token count to the namespace id.
func createClientCountTable(entityMap map[string]uint64, nonEntityMap map[string]uint64) map[string]*clients {
	clientCountTable := make(map[string]*clients)
	for nsID, count := range entityMap {
		if _, ok := clientCountTable[nsID]; !ok {
			clientCountTable[nsID] = &clients{distinctEntities: 0, distinctNonEntities: 0}
		}
		clientCountTable[nsID].distinctEntities += count

	}

	for nsID, count := range nonEntityMap {
		if _, ok := clientCountTable[nsID]; !ok {
			clientCountTable[nsID] = &clients{distinctEntities: 0, distinctNonEntities: 0}
		}
		clientCountTable[nsID].distinctNonEntities += count
	}
	return clientCountTable
}

func (ct *ClientTracker) addClient(e *activity.EntityRecord) {
	if _, ok := ct.activeClients[e.ClientID]; !ok {
		ct.activeClients[e.ClientID] = struct{}{}
		if e.NonEntity == true {
			ct.nonEntityCountByNamespaceID[e.NamespaceID] += 1
		} else {
			ct.entityCountByNamespaceID[e.NamespaceID] += 1
		}
	}
}
