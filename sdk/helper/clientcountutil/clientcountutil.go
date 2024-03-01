// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package clientcountutil provides a library to generate activity log data for
// testing.
package clientcountutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"google.golang.org/protobuf/encoding/protojson"
)

// ActivityLogDataGenerator holds an ActivityLogMockInput. Users can create the
// generator with NewActivityLogData(), add content to the generator using
// the fluent API methods, and generate and write the JSON representation of the
// input to the Vault API.
type ActivityLogDataGenerator struct {
	data            *generation.ActivityLogMockInput
	addingToMonth   *generation.Data
	addingToSegment *generation.Segment
	client          *api.Client
}

// NewActivityLogData creates a new instance of an activity log data generator
// The type returned by this function cannot be called concurrently
func NewActivityLogData(client *api.Client) *ActivityLogDataGenerator {
	return &ActivityLogDataGenerator{
		client: client,
		data:   new(generation.ActivityLogMockInput),
	}
}

// NewCurrentMonthData opens a new month of data for the current month. All
// clients will continue to be added to this month until a new month is created
// with NewPreviousMonthData.
func (d *ActivityLogDataGenerator) NewCurrentMonthData() *ActivityLogDataGenerator {
	return d.newMonth(&generation.Data{Month: &generation.Data_CurrentMonth{CurrentMonth: true}})
}

// NewPreviousMonthData opens a new month of data, where the clients will be
// recorded as having been seen monthsAgo months ago. All clients will continue
// to be added to this month until a new month is created with
// NewPreviousMonthData or NewCurrentMonthData.
func (d *ActivityLogDataGenerator) NewPreviousMonthData(monthsAgo int) *ActivityLogDataGenerator {
	return d.newMonth(&generation.Data{Month: &generation.Data_MonthsAgo{MonthsAgo: int32(monthsAgo)}})
}

func (d *ActivityLogDataGenerator) newMonth(newMonth *generation.Data) *ActivityLogDataGenerator {
	d.data.Data = append(d.data.Data, newMonth)
	d.addingToMonth = newMonth
	d.addingToSegment = nil
	return d
}

// MonthOption holds an option that can be set for the entire month
type MonthOption func(m *generation.Data)

// WithMaximumSegmentIndex sets the maximum segment index for the segments in
// the open month. Set this value in order to set how many indexes the data
// should be split across. This must include any empty or skipped indexes. For
// example, say that you would like all of your data split across indexes 0 and
// 3, with the following empty and skipped indexes:
//
//	empty indexes: [2]
//	skipped indexes: [1]
//
// To accomplish that, you will need to call WithMaximumSegmentIndex(3).
// This value will be ignored if you have called Segment() for the open month
// If not set, all data will be in 1 segment.
func WithMaximumSegmentIndex(n int) MonthOption {
	return func(m *generation.Data) {
		m.NumSegments = int32(n)
	}
}

// WithEmptySegmentIndexes sets which segment indexes should be empty for the
// segments in the open month. If you use this option, you must either:
//  1. ensure that you've called Segment() for the open month
//  2. use WithMaximumSegmentIndex() to set the total number of segments
//
// If you haven't set either of those values then this option will be ignored,
// unless you included 0 as an empty segment index in which case only an empty
// segment will be created.
func WithEmptySegmentIndexes(i ...int) MonthOption {
	return func(m *generation.Data) {
		indexes := make([]int32, 0, len(i))
		for _, index := range i {
			indexes = append(indexes, int32(index))
		}
		m.EmptySegmentIndexes = indexes
	}
}

// WithSkipSegmentIndexes sets which segment indexes should be skipped for the
// segments in the open month. If you use this option, you must either:
//  1. ensure that you've called Segment() for the open month
//  2. use WithMaximumSegmentIndex() to set the total number of segments
//
// If you haven't set either of those values then this option will be ignored,
// unless you included 0 as a skipped segment index in which case no segments
// will be created.
func WithSkipSegmentIndexes(i ...int) MonthOption {
	return func(m *generation.Data) {
		indexes := make([]int32, 0, len(i))
		for _, index := range i {
			indexes = append(indexes, int32(index))
		}
		m.SkipSegmentIndexes = indexes
	}
}

// SetMonthOptions can be called at any time to set options for the open month
func (d *ActivityLogDataGenerator) SetMonthOptions(opts ...MonthOption) *ActivityLogDataGenerator {
	for _, opt := range opts {
		opt(d.addingToMonth)
	}
	return d
}

// ClientOption defines additional options for the client
// This type and the functions that return it are here for ease of use. A user
// could also choose to create the *generation.Client themselves, without using
// a ClientOption
type ClientOption func(client *generation.Client)

// WithClientNamespace sets the namespace for the client
func WithClientNamespace(n string) ClientOption {
	return func(client *generation.Client) {
		client.Namespace = n
	}
}

// WithClientMount sets the mount path for the client
func WithClientMount(m string) ClientOption {
	return func(client *generation.Client) {
		client.Mount = m
	}
}

// WithClientIsNonEntity sets whether the client is an entity client or a non-
// entity token client
func WithClientIsNonEntity() ClientOption {
	return WithClientType("non-entity")
}

// WithClientType sets the client type to the given string. If this client type
// is not "entity", then the client will be counted in the activity log as a
// non-entity client
func WithClientType(typ string) ClientOption {
	return func(client *generation.Client) {
		client.ClientType = typ
	}
}

// WithClientID sets the ID for the client
func WithClientID(id string) ClientOption {
	return func(client *generation.Client) {
		client.Id = id
	}
}

// ClientsSeen adds clients to the month that was most recently opened with
// NewPreviousMonthData or NewCurrentMonthData.
func (d *ActivityLogDataGenerator) ClientsSeen(clients ...*generation.Client) *ActivityLogDataGenerator {
	if d.addingToSegment == nil {
		if d.addingToMonth.Clients == nil {
			d.addingToMonth.Clients = &generation.Data_All{All: &generation.Clients{}}
		}
		d.addingToMonth.GetAll().Clients = append(d.addingToMonth.GetAll().Clients, clients...)
		return d
	}
	d.addingToSegment.Clients.Clients = append(d.addingToSegment.Clients.Clients, clients...)
	return d
}

// NewClientSeen adds 1 new client with the given options to the most recently
// opened month.
func (d *ActivityLogDataGenerator) NewClientSeen(opts ...ClientOption) *ActivityLogDataGenerator {
	return d.NewClientsSeen(1, opts...)
}

// NewClientsSeen adds n new clients with the given options to the most recently
// opened month.
func (d *ActivityLogDataGenerator) NewClientsSeen(n int, opts ...ClientOption) *ActivityLogDataGenerator {
	c := new(generation.Client)
	for _, opt := range opts {
		opt(c)
	}
	c.Count = int32(n)
	return d.ClientsSeen(c)
}

// RepeatedClientSeen adds 1 client that was seen in the previous month to
// the month that was most recently opened. This client will have the attributes
// described by the provided options.
func (d *ActivityLogDataGenerator) RepeatedClientSeen(opts ...ClientOption) *ActivityLogDataGenerator {
	return d.RepeatedClientsSeen(1, opts...)
}

// RepeatedClientsSeen adds n clients that were seen in the previous month to
// the month that was most recently opened. These clients will have the
// attributes described by provided options.
func (d *ActivityLogDataGenerator) RepeatedClientsSeen(n int, opts ...ClientOption) *ActivityLogDataGenerator {
	c := new(generation.Client)
	for _, opt := range opts {
		opt(c)
	}
	c.Repeated = true
	c.Count = int32(n)
	return d.ClientsSeen(c)
}

// RepeatedClientSeenFromMonthsAgo adds 1 client that was seen in monthsAgo
// month to the month that was most recently opened. This client will have the
// attributes described by provided options.
func (d *ActivityLogDataGenerator) RepeatedClientSeenFromMonthsAgo(monthsAgo int, opts ...ClientOption) *ActivityLogDataGenerator {
	return d.RepeatedClientsSeenFromMonthsAgo(1, monthsAgo, opts...)
}

// RepeatedClientsSeenFromMonthsAgo adds n clients that were seen in monthsAgo
// month to the month that was most recently opened. These clients will have the
// attributes described by provided options.
func (d *ActivityLogDataGenerator) RepeatedClientsSeenFromMonthsAgo(n, monthsAgo int, opts ...ClientOption) *ActivityLogDataGenerator {
	c := new(generation.Client)
	for _, opt := range opts {
		opt(c)
	}
	c.RepeatedFromMonth = int32(monthsAgo)
	c.Count = int32(n)
	return d.ClientsSeen(c)
}

// SegmentOption defines additional options for the segment
type SegmentOption func(segment *generation.Segment)

// WithSegmentIndex sets the index for the segment to n. If this option is not
// provided, the segment will be given the next consecutive index
func WithSegmentIndex(n int) SegmentOption {
	return func(segment *generation.Segment) {
		index := int32(n)
		segment.SegmentIndex = &index
	}
}

// Segment starts a segment within the current month. All clients will be added
// to this segment, until either Segment is called again to create a new open
// segment, or NewPreviousMonthData or NewCurrentMonthData is called to open a
// new month.
func (d *ActivityLogDataGenerator) Segment(opts ...SegmentOption) *ActivityLogDataGenerator {
	s := &generation.Segment{
		Clients: &generation.Clients{},
	}
	for _, opt := range opts {
		opt(s)
	}
	if d.addingToMonth.GetSegments() == nil {
		d.addingToMonth.Clients = &generation.Data_Segments{Segments: &generation.Segments{}}
	}
	d.addingToMonth.GetSegments().Segments = append(d.addingToMonth.GetSegments().Segments, s)
	d.addingToSegment = s
	return d
}

// ToJSON returns the JSON representation of the data
func (d *ActivityLogDataGenerator) ToJSON() ([]byte, error) {
	return protojson.Marshal(d.data)
}

// ToProto returns the ActivityLogMockInput protobuf
func (d *ActivityLogDataGenerator) ToProto() *generation.ActivityLogMockInput {
	return d.data
}

// Write writes the data to the API with the given write options. The method
// returns the new paths that have been written. Note that the API endpoint will
// only be present when Vault has been compiled with the "testonly" flag.
func (d *ActivityLogDataGenerator) Write(ctx context.Context, writeOptions ...generation.WriteOptions) ([]string, error) {
	d.data.Write = writeOptions
	err := VerifyInput(d.data)
	if err != nil {
		return nil, err
	}
	data, err := d.ToJSON()
	if err != nil {
		return nil, err
	}
	resp, err := d.client.Logical().WriteWithContext(ctx, "sys/internal/counters/activity/write", map[string]interface{}{"input": string(data)})
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("received no data")
	}
	paths := resp.Data["paths"]
	castedPaths, ok := paths.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid paths data: %v", paths)
	}
	returnPaths := make([]string, 0, len(castedPaths))
	for _, path := range castedPaths {
		returnPaths = append(returnPaths, path.(string))
	}
	return returnPaths, nil
}

// VerifyInput checks that the input data is valid
func VerifyInput(input *generation.ActivityLogMockInput) error {
	// mapping from monthsAgo to the month's data
	months := make(map[int32]*generation.Data)

	// this keeps track of the index of the earliest month. We need to verify
	// that this month doesn't have any repeated clients
	earliestMonthsAgo := int32(0)

	// this map holds a set of the month indexes for any RepeatedFromMonth
	// values. Each element will be checked to ensure month that should be
	// repeated from exists in the input data
	repeatedFromMonths := make(map[int32]struct{})

	for _, month := range input.Data {
		monthsAgo := month.GetMonthsAgo()
		if monthsAgo > earliestMonthsAgo {
			earliestMonthsAgo = monthsAgo
		}

		// verify that no monthsAgo value is repeated
		if _, seen := months[monthsAgo]; seen {
			return fmt.Errorf("multiple months with monthsAgo %d", monthsAgo)
		}
		months[monthsAgo] = month

		// the number of segments should be correct
		if month.NumSegments > 0 && int(month.NumSegments)-len(month.GetSkipSegmentIndexes())-len(month.GetEmptySegmentIndexes()) <= 0 {
			return fmt.Errorf("number of segments %d is too small. It must be large enough to include the empty (%v) and skipped (%v) segments", month.NumSegments, month.GetSkipSegmentIndexes(), month.GetEmptySegmentIndexes())
		}

		if segments := month.GetSegments(); segments != nil {
			if month.NumSegments > 0 {
				return errors.New("cannot specify both number of segments and create segmented data")
			}

			segmentIndexes := make(map[int32]struct{})
			for _, segment := range segments.Segments {

				// collect any RepeatedFromMonth values
				for _, client := range segment.GetClients().GetClients() {
					if repeatFrom := client.RepeatedFromMonth; repeatFrom > 0 {
						repeatedFromMonths[repeatFrom] = struct{}{}
					}
				}

				// verify that no segment indexes are repeated
				segmentIndex := segment.SegmentIndex
				if segmentIndex == nil {
					continue
				}
				if _, seen := segmentIndexes[*segmentIndex]; seen {
					return fmt.Errorf("cannot have repeated segment index %d", *segmentIndex)
				}
				segmentIndexes[*segmentIndex] = struct{}{}
			}
		} else {
			for _, client := range month.GetAll().GetClients() {
				// collect any RepeatedFromMonth values
				if repeatFrom := client.RepeatedFromMonth; repeatFrom > 0 {
					repeatedFromMonths[repeatFrom] = struct{}{}
				}
			}
		}
	}

	// check that the corresponding month exists for all the RepeatedFromMonth
	// values
	for repeated := range repeatedFromMonths {
		if _, ok := months[repeated]; !ok {
			return fmt.Errorf("cannot repeat from %d months ago", repeated)
		}
	}
	// the earliest month can't have any repeated clients, because there are no
	// earlier months to repeat from
	earliestMonth := months[earliestMonthsAgo]
	repeatedClients := false
	if all := earliestMonth.GetAll(); all != nil {
		for _, client := range all.GetClients() {
			repeatedClients = repeatedClients || client.Repeated || client.RepeatedFromMonth != 0
		}
	} else {
		for _, segment := range earliestMonth.GetSegments().GetSegments() {
			for _, client := range segment.GetClients().GetClients() {
				repeatedClients = repeatedClients || client.Repeated || client.RepeatedFromMonth != 0
			}
		}
	}

	if repeatedClients {
		return fmt.Errorf("%d months ago cannot have repeated clients, because it is the earliest month", earliestMonthsAgo)
	}

	return nil
}
