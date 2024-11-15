// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"sort"
	"time"
)

// Evaluations is used to query the evaluation endpoints.
type Evaluations struct {
	client *Client
}

// Evaluations returns a new handle on the evaluations.
func (c *Client) Evaluations() *Evaluations {
	return &Evaluations{client: c}
}

// List is used to dump all of the evaluations.
func (e *Evaluations) List(q *QueryOptions) ([]*Evaluation, *QueryMeta, error) {
	var resp []*Evaluation
	qm, err := e.client.query("/v1/evaluations", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(EvalIndexSort(resp))
	return resp, qm, nil
}

func (e *Evaluations) PrefixList(prefix string) ([]*Evaluation, *QueryMeta, error) {
	return e.List(&QueryOptions{Prefix: prefix})
}

// Count is used to get a count of evaluations.
func (e *Evaluations) Count(q *QueryOptions) (*EvalCountResponse, *QueryMeta, error) {
	var resp *EvalCountResponse
	qm, err := e.client.query("/v1/evaluations/count", &resp, q)
	if err != nil {
		return resp, nil, err
	}
	return resp, qm, nil
}

// Info is used to query a single evaluation by its ID.
func (e *Evaluations) Info(evalID string, q *QueryOptions) (*Evaluation, *QueryMeta, error) {
	var resp Evaluation
	qm, err := e.client.query("/v1/evaluation/"+evalID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Delete is used to batch delete evaluations using their IDs.
func (e *Evaluations) Delete(evalIDs []string, w *WriteOptions) (*WriteMeta, error) {
	req := EvalDeleteRequest{
		EvalIDs: evalIDs,
	}
	wm, err := e.client.delete("/v1/evaluations", &req, nil, w)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// DeleteOpts is used to batch delete evaluations using a filter.
func (e *Evaluations) DeleteOpts(req *EvalDeleteRequest, w *WriteOptions) (*EvalDeleteResponse, *WriteMeta, error) {
	resp := &EvalDeleteResponse{}
	wm, err := e.client.delete("/v1/evaluations", &req, resp, w)
	if err != nil {
		return nil, nil, err
	}
	return resp, wm, nil
}

// Allocations is used to retrieve a set of allocations given
// an evaluation ID.
func (e *Evaluations) Allocations(evalID string, q *QueryOptions) ([]*AllocationListStub, *QueryMeta, error) {
	var resp []*AllocationListStub
	qm, err := e.client.query("/v1/evaluation/"+evalID+"/allocations", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(AllocIndexSort(resp))
	return resp, qm, nil
}

const (
	EvalStatusBlocked   = "blocked"
	EvalStatusPending   = "pending"
	EvalStatusComplete  = "complete"
	EvalStatusFailed    = "failed"
	EvalStatusCancelled = "canceled"
)

// Evaluation is used to serialize an evaluation.
type Evaluation struct {
	ID                   string
	Priority             int
	Type                 string
	TriggeredBy          string
	Namespace            string
	JobID                string
	JobModifyIndex       uint64
	NodeID               string
	NodeModifyIndex      uint64
	DeploymentID         string
	Status               string
	StatusDescription    string
	Wait                 time.Duration
	WaitUntil            time.Time
	NextEval             string
	PreviousEval         string
	BlockedEval          string
	RelatedEvals         []*EvaluationStub
	FailedTGAllocs       map[string]*AllocationMetric
	ClassEligibility     map[string]bool
	EscapedComputedClass bool
	QuotaLimitReached    string
	AnnotatePlan         bool
	QueuedAllocations    map[string]int
	SnapshotIndex        uint64
	CreateIndex          uint64
	ModifyIndex          uint64
	CreateTime           int64
	ModifyTime           int64
}

// EvaluationStub is used to serialize parts of an evaluation returned in the
// RelatedEvals field of an Evaluation.
type EvaluationStub struct {
	ID                string
	Priority          int
	Type              string
	TriggeredBy       string
	Namespace         string
	JobID             string
	NodeID            string
	DeploymentID      string
	Status            string
	StatusDescription string
	WaitUntil         time.Time
	NextEval          string
	PreviousEval      string
	BlockedEval       string
	CreateIndex       uint64
	ModifyIndex       uint64
	CreateTime        int64
	ModifyTime        int64
}

type EvalDeleteRequest struct {
	EvalIDs []string
	Filter  string
	WriteRequest
}

type EvalDeleteResponse struct {
	Count int
}

type EvalCountResponse struct {
	Count int
	QueryMeta
}

// EvalIndexSort is a wrapper to sort evaluations by CreateIndex.
// We reverse the test so that we get the highest index first.
type EvalIndexSort []*Evaluation

func (e EvalIndexSort) Len() int {
	return len(e)
}

func (e EvalIndexSort) Less(i, j int) bool {
	return e[i].CreateIndex > e[j].CreateIndex
}

func (e EvalIndexSort) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
