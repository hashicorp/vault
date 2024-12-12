// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package vault

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"google.golang.org/protobuf/encoding/protojson"
)

const helpText = "Create activity log data for testing purposes"

func (b *SystemBackend) activityWritePath() *framework.Path {
	return &framework.Path{
		Pattern:         "internal/counters/activity/write$",
		HelpDescription: helpText,
		HelpSynopsis:    helpText,
		Fields: map[string]*framework.FieldSchema{
			"input": {
				Type:        framework.TypeString,
				Description: "JSON input for generating mock data",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handleActivityWriteData,
				Summary:  "Write activity log data",
			},
		},
	}
}

func (b *SystemBackend) handleActivityWriteData(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	now := time.Now().UTC()

	json := data.Get("input")
	input := &generation.ActivityLogMockInput{}
	err := protojson.Unmarshal([]byte(json.(string)), input)
	if err != nil {
		return logical.ErrorResponse("Invalid input data: %s", err), logical.ErrInvalidRequest
	}
	if len(input.Write) == 0 {
		return logical.ErrorResponse("Missing required \"write\" values"), logical.ErrInvalidRequest
	}
	if len(input.Data) == 0 {
		return logical.ErrorResponse("Missing required \"data\" values"), logical.ErrInvalidRequest
	}

	err = clientcountutil.VerifyInput(input)
	if err != nil {
		return logical.ErrorResponse("Invalid input data: %s", err), logical.ErrInvalidRequest
	}

	numMonths := 0
	for _, month := range input.Data {
		if int(month.GetMonthsAgo()) > numMonths {
			numMonths = int(month.GetMonthsAgo())
		}
	}
	generated := newMultipleMonthsActivityClients(numMonths + 1)
	for _, month := range input.Data {
		err := generated.processMonth(ctx, b.Core, month, now)
		if err != nil {
			return logical.ErrorResponse("failed to process data for month %d", month.GetMonthsAgo()), err
		}
	}

	opts := make(map[generation.WriteOptions]struct{}, len(input.Write))
	for _, opt := range input.Write {
		opts[opt] = struct{}{}
	}
	localPaths, globalPaths, err := generated.write(ctx, opts, b.Core.activityLog, now)
	if err != nil {
		b.logger.Debug("failed to write activity log data", "error", err.Error())
		return logical.ErrorResponse("failed to write data"), err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"local_paths":  localPaths,
			"global_paths": globalPaths,
		},
	}, nil
}

// singleMonthActivityClients holds a single month's client IDs, in the order they were seen
type singleMonthActivityClients struct {
	// globalClients are indexed by ID
	globalClients []*activity.EntityRecord
	// localClients are indexed by ID
	localClients []*activity.EntityRecord
	// predefinedGlobalSegments map from the segment number to the client's index in
	// the clients slice
	predefinedGlobalSegments map[int][]int
	// predefinedLocalSegments map from the segment number to the client's index in
	// the clients slice
	predefinedLocalSegments map[int][]int
	// generationParameters holds the generation request
	generationParameters *generation.Data
}

// multipleMonthsActivityClients holds multiple month's data
type multipleMonthsActivityClients struct {
	// months are in order, with month 0 being the current month and index 1 being 1 month ago
	months []*singleMonthActivityClients
}

func (s *singleMonthActivityClients) addEntityRecord(core *Core, record *activity.EntityRecord, segmentIndex *int, local bool) {
	if !local {
		s.globalClients = append(s.globalClients, record)
	} else {
		s.localClients = append(s.localClients, record)
	}
	if segmentIndex != nil {
		if !local {
			globalIndex := len(s.globalClients) - 1
			s.predefinedGlobalSegments[*segmentIndex] = append(s.predefinedGlobalSegments[*segmentIndex], globalIndex)
		} else {
			localIndex := len(s.localClients) - 1
			s.predefinedLocalSegments[*segmentIndex] = append(s.predefinedLocalSegments[*segmentIndex], localIndex)
		}
	}
}

// populateSegments converts a month of clients into a segmented map. The map's
// keys are the segment index, and the value are the clients that were seen in
// that index. If the value is an empty slice, then it's an empty index. If the
// value is nil, then it's a skipped index
func (s *singleMonthActivityClients) populateSegments(predefinedSegments map[int][]int, clients []*activity.EntityRecord) (map[int][]*activity.EntityRecord, error) {
	segments := make(map[int][]*activity.EntityRecord)
	ignoreIndexes := make(map[int]struct{})
	skipIndexes := s.generationParameters.SkipSegmentIndexes
	emptyIndexes := s.generationParameters.EmptySegmentIndexes

	for _, i := range skipIndexes {
		segments[int(i)] = nil
		ignoreIndexes[int(i)] = struct{}{}
	}
	for _, i := range emptyIndexes {
		segments[int(i)] = make([]*activity.EntityRecord, 0, 0)
		ignoreIndexes[int(i)] = struct{}{}
	}

	// if we have predefined segments, then we can construct the map using those
	if len(predefinedSegments) > 0 {
		for segment, clientIndexes := range predefinedSegments {
			clientsInSegment := make([]*activity.EntityRecord, 0, len(clientIndexes))
			for _, idx := range clientIndexes {
				clientsInSegment = append(clientsInSegment, clients[idx])
			}
			segments[segment] = clientsInSegment
		}
		return segments, nil
	}

	// determine how many segments are necessary to store the clients for this month
	// using the default storage limits
	numNecessarySegments := len(clients) / ActivitySegmentClientCapacity
	if len(clients)%ActivitySegmentClientCapacity != 0 {
		numNecessarySegments++
	}
	totalSegmentCount := numNecessarySegments

	// override the segment count if set by client
	if s.generationParameters.GetNumSegments() > 0 {
		totalSegmentCount = int(s.generationParameters.GetNumSegments())
	}

	numNonUsable := len(skipIndexes) + len(emptyIndexes)
	usableSegmentCount := totalSegmentCount - numNonUsable
	if usableSegmentCount <= 0 {
		return nil, fmt.Errorf("num segments %d is too low, it must be greater than %d (%d skipped indexes + %d empty indexes)", totalSegmentCount, numNonUsable, len(skipIndexes), len(emptyIndexes))
	}

	// determine how many clients should be in each segment
	segmentSizes := len(clients) / usableSegmentCount
	if len(clients)%usableSegmentCount != 0 {
		segmentSizes++
	}

	if segmentSizes > ActivitySegmentClientCapacity {
		return nil, fmt.Errorf("the number of segments is too low, it must be greater than %d in order to meet storage limits", numNecessarySegments)
	}

	clientIndex := 0
	for i := 0; i < totalSegmentCount; i++ {
		if clientIndex >= len(clients) {
			break
		}
		if _, ok := ignoreIndexes[i]; ok {
			continue
		}
		for len(segments[i]) < segmentSizes && clientIndex < len(clients) {
			segments[i] = append(segments[i], clients[clientIndex])
			clientIndex++
		}
	}
	return segments, nil
}

// addNewClients generates clients according to the given parameters, and adds them to the month
// the client will always have the mountAccessor as its mount accessor
func (s *singleMonthActivityClients) addNewClients(c *generation.Client, mountAccessor string, segmentIndex *int, monthsAgo int32, now time.Time, core *Core) error {
	count := 1
	if c.Count > 1 {
		count = int(c.Count)
	}
	ts := timeutil.MonthsPreviousTo(int(monthsAgo), now)

	// identify is client is local or global
	isLocal, err := isClientLocal(core, c.ClientType, mountAccessor)
	if err != nil {
		return err
	}

	isNonEntity := c.ClientType != entityActivityType
	for i := 0; i < count; i++ {
		record := &activity.EntityRecord{
			ClientID:      c.Id,
			NamespaceID:   c.Namespace,
			MountAccessor: mountAccessor,
			NonEntity:     isNonEntity,
			ClientType:    c.ClientType,
			Timestamp:     ts.Unix(),
		}
		if record.ClientID == "" {
			var err error
			record.ClientID, err = uuid.GenerateUUID()
			if err != nil {
				return err
			}
		}

		s.addEntityRecord(core, record, segmentIndex, isLocal)
	}
	return nil
}

// processMonth populates a month of client data
func (m *multipleMonthsActivityClients) processMonth(ctx context.Context, core *Core, month *generation.Data, now time.Time) error {
	// default to using the root namespace and the first mount on the root namespace
	mounts, err := core.ListMounts()
	if err != nil {
		return err
	}
	defaultMountAccessorRootNS := ""
	for _, mount := range mounts {
		if mount.NamespaceID == namespace.RootNamespaceID {
			defaultMountAccessorRootNS = mount.Accessor
			break
		}
	}
	m.months[month.GetMonthsAgo()].generationParameters = month
	add := func(c []*generation.Client, segmentIndex *int) error {
		for _, clients := range c {
			mountAccessor := defaultMountAccessorRootNS

			if clients.Namespace == "" {
				clients.Namespace = namespace.RootNamespaceID
			}
			if clients.ClientType == "" {
				clients.ClientType = entityActivityType
			}

			if clients.Namespace != namespace.RootNamespaceID && !strings.HasSuffix(clients.Namespace, "/") {
				clients.Namespace += "/"
			}
			// verify that the namespace exists
			ns := core.namespaceByPath(clients.Namespace)
			if ns.ID == namespace.RootNamespaceID && clients.Namespace != namespace.RootNamespaceID {
				return fmt.Errorf("unable to find namespace %s", clients.Namespace)
			}
			clients.Namespace = ns.ID

			// verify that the mount exists
			if clients.Mount != "" {
				if !strings.HasSuffix(clients.Mount, "/") {
					clients.Mount += "/"
				}
				nctx := namespace.ContextWithNamespace(ctx, ns)
				mountEntry := core.router.MatchingMountEntry(nctx, clients.Mount)
				if mountEntry == nil {
					return fmt.Errorf("unable to find matching mount in namespace %s", ns.Path)
				}
				mountAccessor = mountEntry.Accessor
			}

			if clients.Namespace != namespace.RootNamespaceID && clients.Mount == "" {
				// if we're not using the root namespace, find a mount on the namespace that we are using
				found := false
				for _, mount := range mounts {
					if mount.NamespaceID == ns.ID {
						mountAccessor = mount.Accessor
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("unable to find matching mount in namespace %s", ns.Path)
				}
			}

			err = m.addClientToMonth(month.GetMonthsAgo(), clients, mountAccessor, segmentIndex, now, core)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if month.GetAll() != nil {
		return add(month.GetAll().GetClients(), nil)
	}
	predefinedSegments := month.GetSegments()
	for i, segment := range predefinedSegments.GetSegments() {
		index := i
		if segment.SegmentIndex != nil {
			index = int(*segment.SegmentIndex)
		}
		err = add(segment.GetClients().GetClients(), &index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *multipleMonthsActivityClients) addClientToMonth(monthsAgo int32, c *generation.Client, mountAccessor string, segmentIndex *int, now time.Time, core *Core) error {
	if c.Repeated || c.RepeatedFromMonth > 0 {
		return m.addRepeatedClients(monthsAgo, c, mountAccessor, segmentIndex, core)
	}
	return m.months[monthsAgo].addNewClients(c, mountAccessor, segmentIndex, monthsAgo, now, core)
}

func (m *multipleMonthsActivityClients) addRepeatedClients(monthsAgo int32, c *generation.Client, mountAccessor string, segmentIndex *int, core *Core) error {
	addingTo := m.months[monthsAgo]
	repeatedFromMonth := monthsAgo + 1
	if c.RepeatedFromMonth > 0 {
		repeatedFromMonth = c.RepeatedFromMonth
	}
	repeatedFrom := m.months[repeatedFromMonth]

	// identify is client is local or global
	isLocal, err := isClientLocal(core, c.ClientType, mountAccessor)
	if err != nil {
		return err
	}

	numClients := 1
	if c.Count > 0 {
		numClients = int(c.Count)
	}

	repeatedClients := repeatedFrom.globalClients
	if isLocal {
		repeatedClients = repeatedFrom.localClients
	}
	for _, client := range repeatedClients {
		if c.ClientType == client.ClientType && mountAccessor == client.MountAccessor && c.Namespace == client.NamespaceID {
			addingTo.addEntityRecord(core, client, segmentIndex, isLocal)
			numClients--
			if numClients == 0 {
				break
			}
		}
	}
	if numClients > 0 {
		return fmt.Errorf("missing repeated %d clients matching given parameters", numClients)
	}
	return nil
}

// isClientLocal checks whether the given client is on a local mount.
// In all other cases, we will assume it is a global client.
func isClientLocal(core *Core, clientType string, mountAccessor string) (bool, error) {
	// Tokens are not replicated to performance secondary clusters
	if clientType == nonEntityTokenActivityType {
		return true, nil
	}
	mountEntry := core.router.MatchingMountByAccessor(mountAccessor)
	// If the mount entry is nil, this means the mount has been deleted. We will assume it was replicated because we do not want to
	// over count clients
	if mountEntry != nil && mountEntry.Local {
		return true, nil
	}

	return false, nil
}

func (m *multipleMonthsActivityClients) addMissingCurrentMonth() {
	missing := m.months[0].generationParameters == nil &&
		len(m.months) > 1 &&
		m.months[1].generationParameters != nil
	if !missing {
		return
	}
	m.months[0].generationParameters = &generation.Data{EmptySegmentIndexes: []int32{0}, NumSegments: 2}
}

func (m *multipleMonthsActivityClients) timestampForMonth(i int, now time.Time) time.Time {
	if i > 0 {
		return timeutil.StartOfMonth(timeutil.MonthsPreviousTo(i, now))
	}
	return now
}

func (m *multipleMonthsActivityClients) write(ctx context.Context, opts map[generation.WriteOptions]struct{}, activityLog *ActivityLog, now time.Time) ([]string, []string, error) {
	globalPaths := []string{}
	localPaths := []string{}

	_, writePQ := opts[generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES]
	_, writeDistinctClients := opts[generation.WriteOptions_WRITE_DISTINCT_CLIENTS]
	_, writeIntentLog := opts[generation.WriteOptions_WRITE_INTENT_LOGS]

	m.addMissingCurrentMonth()

	for i, month := range m.months {
		if month.generationParameters == nil {
			continue
		}
		timestamp := m.timestampForMonth(i, now)
		if len(month.globalClients) > 0 {
			globalSegments, err := month.populateSegments(month.predefinedGlobalSegments, month.globalClients)
			if err != nil {
				return nil, nil, err
			}
			for segmentIndex, segment := range globalSegments {
				if segment == nil {
					// skip the index
					continue
				}
				entityPath, err := activityLog.saveSegmentEntitiesInternal(ctx, segmentInfo{
					startTimestamp:       timestamp.Unix(),
					currentClients:       &activity.EntityActivityLog{Clients: segment},
					clientSequenceNumber: uint64(segmentIndex),
					tokenCount:           &activity.TokenCount{},
				}, true, activityGlobalPathPrefix)
				if err != nil {
					return nil, nil, err
				}
				globalPaths = append(globalPaths, entityPath)
			}
		}
		if len(month.localClients) > 0 {
			localSegments, err := month.populateSegments(month.predefinedLocalSegments, month.localClients)
			if err != nil {
				return nil, nil, err
			}
			for segmentIndex, segment := range localSegments {
				if segment == nil {
					// skip the index
					continue
				}
				entityPath, err := activityLog.saveSegmentEntitiesInternal(ctx, segmentInfo{
					startTimestamp:       timestamp.Unix(),
					currentClients:       &activity.EntityActivityLog{Clients: segment},
					clientSequenceNumber: uint64(segmentIndex),
					tokenCount:           &activity.TokenCount{},
				}, true, activityLocalPathPrefix)
				if err != nil {
					return nil, nil, err
				}
				localPaths = append(localPaths, entityPath)
			}
		}
	}
	if writePQ || writeDistinctClients {
		// start with the oldest month of data, and create precomputed queries
		// up to that month
		pqWg := sync.WaitGroup{}
		for i := len(m.months) - 1; i > 0; i-- {
			pqWg.Add(1)
			go func(i int) {
				defer pqWg.Done()
				activityLog.precomputedQueryWorker(ctx, &ActivityIntentLog{
					PreviousMonth: m.timestampForMonth(i, now).Unix(),
					NextMonth:     now.Unix(),
				})
			}(i)
		}
		pqWg.Wait()
	}
	if writeIntentLog {
		err := activityLog.writeIntentLog(ctx, m.latestTimestamp(now, false).Unix(), m.latestTimestamp(now, true).UTC())
		if err != nil {
			return nil, nil, err
		}
	}
	wg := sync.WaitGroup{}
	err := activityLog.refreshFromStoredLog(ctx, &wg, now)
	if err != nil {
		return nil, nil, err
	}
	wg.Wait()
	return localPaths, globalPaths, nil
}

func (m *multipleMonthsActivityClients) latestTimestamp(now time.Time, includeCurrentMonth bool) time.Time {
	for i, month := range m.months {
		if month.generationParameters != nil && (i != 0 || includeCurrentMonth) {
			return timeutil.StartOfMonth(timeutil.MonthsPreviousTo(i, now))
		}
	}
	return time.Time{}
}

func (m *multipleMonthsActivityClients) earliestTimestamp(now time.Time) time.Time {
	for i := len(m.months) - 1; i >= 0; i-- {
		month := m.months[i]
		if month.generationParameters != nil {
			return timeutil.StartOfMonth(timeutil.MonthsPreviousTo(i, now))
		}
	}
	return time.Time{}
}

func newMultipleMonthsActivityClients(numberOfMonths int) *multipleMonthsActivityClients {
	m := &multipleMonthsActivityClients{
		months: make([]*singleMonthActivityClients, numberOfMonths),
	}
	for i := 0; i < numberOfMonths; i++ {
		m.months[i] = &singleMonthActivityClients{
			predefinedGlobalSegments: make(map[int][]int),
			predefinedLocalSegments:  make(map[int][]int),
		}
	}
	return m
}

func newProtoSegmentReader(segments map[int][]*activity.EntityRecord) SegmentReader {
	allRecords := make([][]*activity.EntityRecord, 0, len(segments))
	for _, records := range segments {
		if segments == nil {
			continue
		}
		allRecords = append(allRecords, records)
	}
	return &sliceSegmentReader{
		records: allRecords,
	}
}

type sliceSegmentReader struct {
	records [][]*activity.EntityRecord
	i       int
}

// ReadGlobalEntity here is a dummy implementation.
// Segment reader is never used when writing using the ClientCountUtil library
func (p *sliceSegmentReader) ReadGlobalEntity(ctx context.Context) (*activity.EntityActivityLog, error) {
	if p.i == len(p.records) {
		return nil, io.EOF
	}
	record := p.records[p.i]
	p.i++
	return &activity.EntityActivityLog{Clients: record}, nil
}

// ReadLocalEntity here is a dummy implementation.
// Segment reader is never used when writing using the ClientCountUtil library
func (p *sliceSegmentReader) ReadLocalEntity(ctx context.Context) (*activity.EntityActivityLog, error) {
	if p.i == len(p.records) {
		return nil, io.EOF
	}
	record := p.records[p.i]
	p.i++
	return &activity.EntityActivityLog{Clients: record}, nil
}

func (p *sliceSegmentReader) ReadToken(ctx context.Context) (*activity.TokenCount, error) {
	return nil, io.EOF
}
