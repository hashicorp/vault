// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

// Recommendations is used to query the recommendations endpoints.
type Recommendations struct {
	client *Client
}

// Recommendations returns a new handle on the recommendations endpoints.
func (c *Client) Recommendations() *Recommendations {
	return &Recommendations{client: c}
}

// List is used to dump all of the recommendations in the cluster
func (r *Recommendations) List(q *QueryOptions) ([]*Recommendation, *QueryMeta, error) {
	var resp []*Recommendation
	qm, err := r.client.query("/v1/recommendations", &resp, q)
	if err != nil {
		return nil, qm, err
	}
	return resp, qm, nil
}

// Info is used to return information on a single recommendation
func (r *Recommendations) Info(id string, q *QueryOptions) (*Recommendation, *QueryMeta, error) {
	var resp Recommendation
	qm, err := r.client.query("/v1/recommendation/"+id, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Upsert is used to create or update a recommendation
func (r *Recommendations) Upsert(rec *Recommendation, q *WriteOptions) (*Recommendation, *WriteMeta, error) {
	var resp Recommendation
	wm, err := r.client.put("/v1/recommendation", rec, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Delete is used to delete a list of recommendations
func (r *Recommendations) Delete(ids []string, q *WriteOptions) (*WriteMeta, error) {
	req := &RecommendationApplyRequest{
		Apply:   []string{},
		Dismiss: ids,
	}
	wm, err := r.client.put("/v1/recommendations/apply", req, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Apply is used to apply a set of recommendations
func (r *Recommendations) Apply(ids []string, policyOverride bool) (
	*RecommendationApplyResponse, *WriteMeta, error) {
	req := &RecommendationApplyRequest{
		Apply:          ids,
		PolicyOverride: policyOverride,
	}
	var resp RecommendationApplyResponse
	wm, err := r.client.put("/v1/recommendations/apply", req, &resp, nil)
	if err != nil {
		return nil, nil, err
	}
	resp.WriteMeta = *wm
	return &resp, wm, nil
}

// Recommendation is used to serialize a recommendation.
type Recommendation struct {
	ID             string
	Region         string
	Namespace      string
	JobID          string
	JobVersion     uint64
	Group          string
	Task           string
	Resource       string
	Value          int
	Current        int
	Meta           map[string]interface{}
	Stats          map[string]float64
	EnforceVersion bool

	SubmitTime int64

	CreateIndex uint64
	ModifyIndex uint64
}

// RecommendationApplyRequest is used to apply and/or dismiss a set of recommendations
type RecommendationApplyRequest struct {
	Apply          []string
	Dismiss        []string
	PolicyOverride bool
}

// RecommendationApplyResponse is used to apply a set of recommendations
type RecommendationApplyResponse struct {
	UpdatedJobs []*SingleRecommendationApplyResult
	Errors      []*SingleRecommendationApplyError
	WriteMeta
}

type SingleRecommendationApplyResult struct {
	Namespace       string
	JobID           string
	JobModifyIndex  uint64
	EvalID          string
	EvalCreateIndex uint64
	Warnings        string
	Recommendations []string
}

type SingleRecommendationApplyError struct {
	Namespace       string
	JobID           string
	Recommendations []string
	Error           string
}
