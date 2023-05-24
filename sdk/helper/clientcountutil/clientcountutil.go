// Package clientcountutil provides a library to generate activity log data for
// testing.
package clientcountutil

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"google.golang.org/protobuf/encoding/protojson"
)

type dataGenerationWrapper struct {
	data            *generation.ActivityLogMockInput
	addingToMonth   *generation.Data
	addingToSegment *generation.Segment
	client          *api.Client
}

// NewActivityLogData creates a new instance of an activity log data generator
// The type returned by this function cannot be called concurrently
func NewActivityLogData(client *api.Client) *dataGenerationWrapper {
	return &dataGenerationWrapper{
		client: client,
		data:   new(generation.ActivityLogMockInput),
	}
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
	return func(client *generation.Client) {
		client.NonEntity = true
	}
}

// WithClientID sets the ID for the client
func WithClientID(id string) ClientOption {
	return func(client *generation.Client) {
		client.Id = id
	}
}

// NewCurrentMonthData opens a new month of data for the current month. All
// clients will continue to be added to this month until a new month is created
// with NewMonthDataMonthsAgo.
func (d *dataGenerationWrapper) NewCurrentMonthData() *dataGenerationWrapper {
	d.newMonth(&generation.Data{Month: new(generation.Data_CurrentMonth)})
	return d
}

// NewMonthDataMonthsAgo opens a new month of data, where the clients will be
// recorded as having been seen monthsAgo months ago. All clients will continue
// to be added to this month until a new month is created with
// NewMonthDataMonthsAgo or NewCurrentMonthData.
func (d *dataGenerationWrapper) NewMonthDataMonthsAgo(monthsAgo int) *dataGenerationWrapper {
	d.newMonth(&generation.Data{Month: &generation.Data_MonthsAgo{MonthsAgo: int32(monthsAgo)}})
	return d
}

func (d *dataGenerationWrapper) newMonth(newMonth *generation.Data) {
	d.data.Data = append(d.data.Data, newMonth)
	d.addingToMonth = newMonth
	d.addingToSegment = nil
}

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
func (d *dataGenerationWrapper) SetMonthOptions(opts ...MonthOption) *dataGenerationWrapper {
	for _, opt := range opts {
		opt(d.addingToMonth)
	}
	return d
}

// ClientsSeen adds clients to the month that was most recently opened with
// NewMonthDataMonthsAgo or NewCurrentMonthData.
func (d *dataGenerationWrapper) ClientsSeen(clients ...*generation.Client) *dataGenerationWrapper {
	if d.addingToSegment == nil {
		d.addingToMonth.GetAll().Clients = append(d.addingToMonth.GetAll().Clients, clients...)
		return d
	}
	d.addingToSegment.Clients.Clients = append(d.addingToSegment.Clients.Clients, clients...)
	return d
}

// NewClientSeen adds 1 new client with the given options to the most recently
// opened month.
func (d *dataGenerationWrapper) NewClientSeen(opts ...ClientOption) *dataGenerationWrapper {
	return d.NewClientsSeen(1, opts...)
}

// NewClientsSeen adds n new clients with the given options to the most recently
// opened month.
func (d *dataGenerationWrapper) NewClientsSeen(n int, opts ...ClientOption) *dataGenerationWrapper {
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
func (d *dataGenerationWrapper) RepeatedClientSeen(opts ...ClientOption) *dataGenerationWrapper {
	return d.RepeatedClientsSeen(1, opts...)
}

// RepeatedClientsSeen adds n clients that were seen in the previous month to
// the month that was most recently opened. These clients will have the
// attributes described by provided options.
func (d *dataGenerationWrapper) RepeatedClientsSeen(n int, opts ...ClientOption) *dataGenerationWrapper {
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
func (d *dataGenerationWrapper) RepeatedClientSeenFromMonthsAgo(monthsAgo int, opts ...ClientOption) *dataGenerationWrapper {
	return d.RepeatedClientsSeenFromMonthsAgo(1, monthsAgo, opts...)
}

// RepeatedClientsSeenFromMonthsAgo adds n clients that were seen in monthsAgo
// month to the month that was most recently opened. These clients will have the
// attributes described by provided options.
func (d *dataGenerationWrapper) RepeatedClientsSeenFromMonthsAgo(n, monthsAgo int, opts ...ClientOption) *dataGenerationWrapper {
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
// segment, or NewMonthDataMonthsAgo or NewCurrentMonthData is called to open a
// new month.
func (d *dataGenerationWrapper) Segment(opts ...SegmentOption) *dataGenerationWrapper {
	s := &generation.Segment{}
	for _, opt := range opts {
		opt(s)
	}
	d.addingToMonth.GetSegments().Segments = append(d.addingToMonth.GetSegments().Segments, s)
	d.addingToSegment = s
	return d
}

// ToJSON returns the JSON representation of the data
func (d *dataGenerationWrapper) ToJSON() ([]byte, error) {
	return protojson.Marshal(d.data)
}

// ToProto returns the ActivityLogMockInput protobuf
func (d *dataGenerationWrapper) ToProto() *generation.ActivityLogMockInput {
	return d.data
}

// Write writes the data to the API with the given write options. The method
// returns the new paths that have been written. Note that the API endpoint will
// only be present when Vault has been compiled with the "testonly" flag.
func (d *dataGenerationWrapper) Write(ctx context.Context, writeOptions ...generation.WriteOptions) ([]string, error) {
	d.data.Write = writeOptions
	data, err := d.ToJSON()
	if err != nil {
		return nil, err
	}
	resp, err := d.client.Logical().WriteBytesWithContext(ctx, "sys/internal/counters/activity/write", data)
	if err != nil {
		return nil, err
	}
	paths := resp.Data["paths"]
	castedPaths, ok := paths.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid paths data: %v", paths)
	}
	return castedPaths, nil
}
