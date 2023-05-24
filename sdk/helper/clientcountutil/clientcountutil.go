package clientcountutil

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
)

type dataGenerationWrapper struct {
	data            *generation.ActivityLogMockInput
	addingToMonth   *generation.Data
	addingToSegment *generation.Segment
	client          *api.Client
	json            []byte
}

func NewActivityLogData(client *api.Client) *dataGenerationWrapper {
	return &dataGenerationWrapper{
		client: client,
		data:   new(generation.ActivityLogMockInput),
	}
}

type ClientOption func(client *generation.Client)

func WithClientNamespace(n string) ClientOption {
	return func(client *generation.Client) {
		client.Namespace = n
	}
}
func WithClientMount(m string) ClientOption {
	return func(client *generation.Client) {
		client.Mount = m
	}
}
func WithClientIsNonEntity() ClientOption {
	return func(client *generation.Client) {
		client.NonEntity = true
	}
}
func WithClientID(id string) ClientOption {
	return func(client *generation.Client) {
		client.Id = id
	}
}

func (d *dataGenerationWrapper) NewCurrentMonthData() *dataGenerationWrapper {
	newMonth := &generation.Data{Month: new(generation.Data_CurrentMonth)}
	d.data.Data = append(d.data.Data, newMonth)
	d.addingToMonth = newMonth
	return d
}
func (d *dataGenerationWrapper) NewMonthDataMonthsAgo(n int) *dataGenerationWrapper {
	newMonth := &generation.Data{Month: &generation.Data_MonthsAgo{MonthsAgo: int32(n)}}
	d.data.Data = append(d.data.Data, newMonth)
	d.addingToMonth = newMonth
	return d
}

func (d *dataGenerationWrapper) ClientsSeen(clients ...*generation.Client) *dataGenerationWrapper {
	if d.addingToSegment == nil {
		d.addingToMonth.GetAll().Clients = append(d.addingToMonth.GetAll().Clients, clients...)
		return d
	}
	d.addingToSegment.Clients.Clients = append(d.addingToSegment.Clients.Clients, clients...)
	return d
}

func (d *dataGenerationWrapper) NewClientSeen(opts ...ClientOption) *dataGenerationWrapper {
	return d.NewClientsSeen(1, opts...)
}
func (d *dataGenerationWrapper) NewClientsSeen(n int, opts ...ClientOption) *dataGenerationWrapper {
	c := new(generation.Client)
	for _, opt := range opts {
		opt(c)
	}
	c.Count = int32(n)
	return d.ClientsSeen(c)
}
func (d *dataGenerationWrapper) RepeatedClientSeen(opts ...ClientOption) *dataGenerationWrapper {
	return d.RepeatedClientsSeen(1, opts...)
}
func (d *dataGenerationWrapper) RepeatedClientsSeen(n int, opts ...ClientOption) *dataGenerationWrapper {
	c := new(generation.Client)
	for _, opt := range opts {
		opt(c)
	}
	c.Repeated = true
	c.Count = int32(n)
	return d.ClientsSeen(c)
}
func (d *dataGenerationWrapper) RepeatedClientSeenFromMonthsAgo(monthsAgo int, opts ...ClientOption) *dataGenerationWrapper {
	return d.RepeatedClientsSeenFromMonthsAgo(1, monthsAgo, opts...)
}
func (d *dataGenerationWrapper) RepeatedClientsSeenFromMonthsAgo(n, monthsAgo int, opts ...ClientOption) *dataGenerationWrapper {
	c := new(generation.Client)
	for _, opt := range opts {
		opt(c)
	}
	c.RepeatedFromMonth = int32(monthsAgo)
	c.Count = int32(n)
	return d.ClientsSeen(c)
}

type SegmentOption func(segment *generation.Segment)

func WithSegmentIndex(n int) SegmentOption {
	return func(segment *generation.Segment) {
		index := int32(n)
		segment.SegmentIndex = &index
	}
}

func (d *dataGenerationWrapper) Segment(opts ...SegmentOption) *dataGenerationWrapper {
	s := &generation.Segment{}
	for _, opt := range opts {
		opt(s)
	}
	d.addingToMonth.GetSegments().Segments = append(d.addingToMonth.GetSegments().Segments, s)
	d.addingToSegment = s
	return d
}

func (d *dataGenerationWrapper) ToJSON() ([]byte, error) {
	if d.json != nil {
		return d.json, nil
	}
	var err error
	return d.json, err
}

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
