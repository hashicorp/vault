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
	paths, err := generated.write(ctx, opts, b.Core.activityLog, now)
	if err != nil {
		b.logger.Debug("failed to write activity log data", "error", err.Error())
		return logical.ErrorResponse("failed to write data"), err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"paths": paths,
		},
	}, nil
}

// singleMonthActivityClients holds a single month's client IDs, in the order they were seen
type singleMonthActivityClients struct {
	// clients are indexed by ID
	clients []*activity.EntityRecord
	// predefinedSegments map from the segment number to the client's index in
	// the clients slice
	predefinedSegments map[int][]int
	// generationParameters holds the generation request
	generationParameters *generation.Data
}

// multipleMonthsActivityClients holds multiple month's data
type multipleMonthsActivityClients struct {
	// months are in order, with month 0 being the current month and index 1 being 1 month ago
	months []*singleMonthActivityClients
}

func (s *singleMonthActivityClients) addEntityRecord(record *activity.EntityRecord, segmentIndex *int) {
	s.clients = append(s.clients, record)
	if segmentIndex != nil {
		index := len(s.clients) - 1
		s.predefinedSegments[*segmentIndex] = append(s.predefinedSegments[*segmentIndex], index)
	}
}

// populateSegments converts a month of clients into a segmented map. The map's
// keys are the segment index, and the value are the clients that were seen in
// that index. If the value is an empty slice, then it's an empty index. If the
// value is nil, then it's a skipped index
func (s *singleMonthActivityClients) populateSegments() (map[int][]*activity.EntityRecord, error) {
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
	if len(s.predefinedSegments) > 0 {
		for segment, clientIndexes := range s.predefinedSegments {
			clientsInSegment := make([]*activity.EntityRecord, 0, len(clientIndexes))
			for _, idx := range clientIndexes {
				clientsInSegment = append(clientsInSegment, s.clients[idx])
			}
			segments[segment] = clientsInSegment
		}
		return segments, nil
	}

	// determine how many segments are necessary to store the clients for this month
	// using the default storage limits
	numNecessarySegments := len(s.clients) / ActivitySegmentClientCapacity
	if len(s.clients)%ActivitySegmentClientCapacity != 0 {
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
	segmentSizes := len(s.clients) / usableSegmentCount
	if len(s.clients)%usableSegmentCount != 0 {
		segmentSizes++
	}

	if segmentSizes > ActivitySegmentClientCapacity {
		return nil, fmt.Errorf("the number of segments is too low, it must be greater than %d in order to meet storage limits", numNecessarySegments)
	}

	clientIndex := 0
	for i := 0; i < totalSegmentCount; i++ {
		if clientIndex >= len(s.clients) {
			break
		}
		if _, ok := ignoreIndexes[i]; ok {
			continue
		}
		for len(segments[i]) < segmentSizes && clientIndex < len(s.clients) {
			segments[i] = append(segments[i], s.clients[clientIndex])
			clientIndex++
		}
	}
	return segments, nil
}

// addNewClients generates clients according to the given parameters, and adds them to the month
// the client will always have the mountAccessor as its mount accessor
func (s *singleMonthActivityClients) addNewClients(c *generation.Client, mountAccessor string, segmentIndex *int, monthsAgo int32, now time.Time) error {
	count := 1
	if c.Count > 1 {
		count = int(c.Count)
	}
	isNonEntity := c.ClientType != entityActivityType
	ts := timeutil.MonthsPreviousTo(int(monthsAgo), now)

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
		s.addEntityRecord(record, segmentIndex)
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

			err = m.addClientToMonth(month.GetMonthsAgo(), clients, mountAccessor, segmentIndex, now)
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

func (m *multipleMonthsActivityClients) addClientToMonth(monthsAgo int32, c *generation.Client, mountAccessor string, segmentIndex *int, now time.Time) error {
	if c.Repeated || c.RepeatedFromMonth > 0 {
		return m.addRepeatedClients(monthsAgo, c, mountAccessor, segmentIndex)
	}
	return m.months[monthsAgo].addNewClients(c, mountAccessor, segmentIndex, monthsAgo, now)
}

func (m *multipleMonthsActivityClients) addRepeatedClients(monthsAgo int32, c *generation.Client, mountAccessor string, segmentIndex *int) error {
	addingTo := m.months[monthsAgo]
	repeatedFromMonth := monthsAgo + 1
	if c.RepeatedFromMonth > 0 {
		repeatedFromMonth = c.RepeatedFromMonth
	}
	repeatedFrom := m.months[repeatedFromMonth]
	numClients := 1
	if c.Count > 0 {
		numClients = int(c.Count)
	}
	for _, client := range repeatedFrom.clients {
		if c.ClientType == client.ClientType && mountAccessor == client.MountAccessor && c.Namespace == client.NamespaceID {
			addingTo.addEntityRecord(client, segmentIndex)
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

func (m *multipleMonthsActivityClients) write(ctx context.Context, opts map[generation.WriteOptions]struct{}, activityLog *ActivityLog, now time.Time) ([]string, error) {
	paths := []string{}

	_, writePQ := opts[generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES]
	_, writeDistinctClients := opts[generation.WriteOptions_WRITE_DISTINCT_CLIENTS]
	_, writeIntentLog := opts[generation.WriteOptions_WRITE_INTENT_LOGS]

	m.addMissingCurrentMonth()

	for i, month := range m.months {
		if month.generationParameters == nil {
			continue
		}
		timestamp := m.timestampForMonth(i, now)
		segments, err := month.populateSegments()
		if err != nil {
			return nil, err
		}
		for segmentIndex, segment := range segments {
			if segment == nil {
				// skip the index
				continue
			}
			entityPath, err := activityLog.saveSegmentEntitiesInternal(ctx, segmentInfo{
				startTimestamp:       timestamp.Unix(),
				currentClients:       &activity.EntityActivityLog{Clients: segment},
				clientSequenceNumber: uint64(segmentIndex),
				tokenCount:           &activity.TokenCount{},
			}, true)
			if err != nil {
				return nil, err
			}
			paths = append(paths, entityPath)
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
			return nil, err
		}
	}
	wg := sync.WaitGroup{}
	err := activityLog.refreshFromStoredLog(ctx, &wg, now)
	if err != nil {
		return nil, err
	}
	wg.Wait()
	return paths, nil
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
			predefinedSegments: make(map[int][]int),
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

func (p *sliceSegmentReader) ReadToken(ctx context.Context) (*activity.TokenCount, error) {
	return nil, io.EOF
}

func (p *sliceSegmentReader) ReadEntity(ctx context.Context) (*activity.EntityActivityLog, error) {
	if p.i == len(p.records) {
		return nil, io.EOF
	}
	record := p.records[p.i]
	p.i++
	return &activity.EntityActivityLog{Clients: record}, nil
}
