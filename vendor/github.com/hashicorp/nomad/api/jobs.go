// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/cronexpr"
)

const (
	// JobTypeService indicates a long-running processes
	JobTypeService = "service"

	// JobTypeBatch indicates a short-lived process
	JobTypeBatch = "batch"

	// JobTypeSystem indicates a system process that should run on all clients
	JobTypeSystem = "system"

	// JobTypeSysbatch indicates a short-lived system process that should run
	// on all clients.
	JobTypeSysbatch = "sysbatch"

	// JobDefaultPriority is the default priority if not specified.
	JobDefaultPriority = 50

	// PeriodicSpecCron is used for a cron spec.
	PeriodicSpecCron = "cron"

	// DefaultNamespace is the default namespace.
	DefaultNamespace = "default"

	// For Job configuration, GlobalRegion is a sentinel region value
	// that users may specify to indicate the job should be run on
	// the region of the node that the job was submitted to.
	// For Client configuration, if no region information is given,
	// the client node will default to be part of the GlobalRegion.
	GlobalRegion = "global"
)

const (
	// RegisterEnforceIndexErrPrefix is the prefix to use in errors caused by
	// enforcing the job modify index during registers.
	RegisterEnforceIndexErrPrefix = "Enforcing job modify index"
)

const (
	// JobPeriodicLaunchSuffix is the string appended to the periodic jobs ID
	// when launching derived instances of it.
	JobPeriodicLaunchSuffix = "/periodic-"

	// JobDispatchLaunchSuffix is the string appended to the parameterized job's ID
	// when dispatching instances of it.
	JobDispatchLaunchSuffix = "/dispatch-"
)

// Jobs is used to access the job-specific endpoints.
type Jobs struct {
	client *Client
}

// JobsParseRequest is used for arguments of the /v1/jobs/parse endpoint
type JobsParseRequest struct {
	// JobHCL is an hcl jobspec
	JobHCL string

	// Variables are HCL2 variables associated with the job. Only works with hcl2.
	//
	// Interpreted as if it were the content of a variables file.
	Variables string

	// Canonicalize is a flag as to if the server should return default values
	// for unset fields
	Canonicalize bool
}

// Jobs returns a handle on the jobs endpoints.
func (c *Client) Jobs() *Jobs {
	return &Jobs{client: c}
}

// ParseHCL is used to convert the HCL representation of a Job to JSON server side.
// To parse the HCL client side see package github.com/hashicorp/nomad/jobspec
// Use ParseHCLOpts if you need to customize JobsParseRequest.
func (j *Jobs) ParseHCL(jobHCL string, canonicalize bool) (*Job, error) {
	req := &JobsParseRequest{
		JobHCL:       jobHCL,
		Canonicalize: canonicalize,
	}
	return j.ParseHCLOpts(req)
}

// ParseHCLOpts is used to request the server convert the HCL representation of a
// Job to JSON on our behalf. Only accepts HCL2 jobs as input.
func (j *Jobs) ParseHCLOpts(req *JobsParseRequest) (*Job, error) {
	var job Job
	_, err := j.client.put("/v1/jobs/parse", req, &job, nil)
	return &job, err
}

func (j *Jobs) Validate(job *Job, q *WriteOptions) (*JobValidateResponse, *WriteMeta, error) {
	var resp JobValidateResponse
	req := &JobValidateRequest{Job: job}
	if q != nil {
		req.WriteRequest = WriteRequest{Region: q.Region}
	}
	wm, err := j.client.put("/v1/validate/job", req, &resp, q)
	return &resp, wm, err
}

// RegisterOptions is used to pass through job registration parameters
type RegisterOptions struct {
	EnforceIndex   bool
	ModifyIndex    uint64
	PolicyOverride bool
	PreserveCounts bool
	EvalPriority   int
	Submission     *JobSubmission
}

// Register is used to register a new job. It returns the ID
// of the evaluation, along with any errors encountered.
func (j *Jobs) Register(job *Job, q *WriteOptions) (*JobRegisterResponse, *WriteMeta, error) {
	return j.RegisterOpts(job, nil, q)
}

// EnforceRegister is used to register a job enforcing its job modify index.
func (j *Jobs) EnforceRegister(job *Job, modifyIndex uint64, q *WriteOptions) (*JobRegisterResponse, *WriteMeta, error) {
	opts := RegisterOptions{EnforceIndex: true, ModifyIndex: modifyIndex}
	return j.RegisterOpts(job, &opts, q)
}

// RegisterOpts is used to register a new job with the passed RegisterOpts. It
// returns the ID of the evaluation, along with any errors encountered.
func (j *Jobs) RegisterOpts(job *Job, opts *RegisterOptions, q *WriteOptions) (*JobRegisterResponse, *WriteMeta, error) {
	// Format the request
	req := &JobRegisterRequest{Job: job}
	if opts != nil {
		if opts.EnforceIndex {
			req.EnforceIndex = true
			req.JobModifyIndex = opts.ModifyIndex
		}
		req.PolicyOverride = opts.PolicyOverride
		req.PreserveCounts = opts.PreserveCounts
		req.EvalPriority = opts.EvalPriority
		req.Submission = opts.Submission
	}

	var resp JobRegisterResponse
	wm, err := j.client.put("/v1/jobs", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

type JobListFields struct {
	Meta bool
}
type JobListOptions struct {
	Fields *JobListFields
}

// List is used to list all of the existing jobs.
func (j *Jobs) List(q *QueryOptions) ([]*JobListStub, *QueryMeta, error) {
	return j.ListOptions(nil, q)
}

// List is used to list all of the existing jobs.
func (j *Jobs) ListOptions(opts *JobListOptions, q *QueryOptions) ([]*JobListStub, *QueryMeta, error) {
	var resp []*JobListStub

	destinationURL := "/v1/jobs"

	if opts != nil && opts.Fields != nil {
		qp := url.Values{}
		qp.Add("meta", fmt.Sprint(opts.Fields.Meta))
		destinationURL = destinationURL + "?" + qp.Encode()
	}

	qm, err := j.client.query(destinationURL, &resp, q)
	if err != nil {
		return nil, qm, err
	}
	sort.Sort(JobIDSort(resp))
	return resp, qm, nil
}

// PrefixList is used to list all existing jobs that match the prefix.
func (j *Jobs) PrefixList(prefix string) ([]*JobListStub, *QueryMeta, error) {
	return j.List(&QueryOptions{Prefix: prefix})
}

// Info is used to retrieve information about a particular
// job given its unique ID.
func (j *Jobs) Info(jobID string, q *QueryOptions) (*Job, *QueryMeta, error) {
	var resp Job
	qm, err := j.client.query("/v1/job/"+url.PathEscape(jobID), &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Scale is used to scale a job.
func (j *Jobs) Scale(jobID, group string, count *int, message string, error bool, meta map[string]interface{},
	q *WriteOptions) (*JobRegisterResponse, *WriteMeta, error) {

	var count64 *int64
	if count != nil {
		count64 = pointerOf(int64(*count))
	}
	req := &ScalingRequest{
		Count: count64,
		Target: map[string]string{
			"Job":   jobID,
			"Group": group,
		},
		Error:   error,
		Message: message,
		Meta:    meta,
	}
	var resp JobRegisterResponse
	qm, err := j.client.put(fmt.Sprintf("/v1/job/%s/scale", url.PathEscape(jobID)), req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// ScaleWithRequest is used to scale a job, giving the caller complete control
// over the ScalingRequest
func (j *Jobs) ScaleWithRequest(jobID string, req *ScalingRequest, q *WriteOptions) (*JobRegisterResponse, *WriteMeta, error) {
	var resp JobRegisterResponse
	qm, err := j.client.put(fmt.Sprintf("/v1/job/%s/scale", url.PathEscape(jobID)), req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// ScaleStatus is used to retrieve information about a particular
// job given its unique ID.
func (j *Jobs) ScaleStatus(jobID string, q *QueryOptions) (*JobScaleStatusResponse, *QueryMeta, error) {
	var resp JobScaleStatusResponse
	qm, err := j.client.query(fmt.Sprintf("/v1/job/%s/scale", url.PathEscape(jobID)), &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Versions is used to retrieve all versions of a particular job given its
// unique ID.
func (j *Jobs) Versions(jobID string, diffs bool, q *QueryOptions) ([]*Job, []*JobDiff, *QueryMeta, error) {
	opts := &VersionsOptions{
		Diffs: diffs,
	}
	return j.VersionsOpts(jobID, opts, q)
}

// VersionByTag is used to retrieve a job version by its VersionTag name.
func (j *Jobs) VersionByTag(jobID, tag string, q *QueryOptions) (*Job, *QueryMeta, error) {
	versions, _, qm, err := j.Versions(jobID, false, q)
	if err != nil {
		return nil, nil, err
	}

	// Find the version with the matching tag
	for _, version := range versions {
		if version.VersionTag != nil && version.VersionTag.Name == tag {
			return version, qm, nil
		}
	}

	return nil, nil, fmt.Errorf("version tag %s not found for job %s", tag, jobID)
}

type VersionsOptions struct {
	Diffs       bool
	DiffTag     string
	DiffVersion *uint64
}

func (j *Jobs) VersionsOpts(jobID string, opts *VersionsOptions, q *QueryOptions) ([]*Job, []*JobDiff, *QueryMeta, error) {
	var resp JobVersionsResponse

	qp := url.Values{}
	if opts != nil {
		qp.Add("diffs", strconv.FormatBool(opts.Diffs))
		if opts.DiffTag != "" {
			qp.Add("diff_tag", opts.DiffTag)
		}
		if opts.DiffVersion != nil {
			qp.Add("diff_version", strconv.FormatUint(*opts.DiffVersion, 10))
		}
	}

	qm, err := j.client.query(fmt.Sprintf("/v1/job/%s/versions?%s", url.PathEscape(jobID), qp.Encode()), &resp, q)
	if err != nil {
		return nil, nil, nil, err
	}
	return resp.Versions, resp.Diffs, qm, nil
}

// Submission is used to retrieve the original submitted source of a job given its
// namespace, jobID, and version number. The original source might not be available,
// which case nil is returned with no error.
func (j *Jobs) Submission(jobID string, version int, q *QueryOptions) (*JobSubmission, *QueryMeta, error) {
	var sub JobSubmission
	s := fmt.Sprintf("/v1/job/%s/submission?version=%d", url.PathEscape(jobID), version)
	qm, err := j.client.query(s, &sub, q)
	if err != nil {
		return nil, nil, err
	}
	return &sub, qm, nil
}

// Allocations is used to return the allocs for a given job ID.
func (j *Jobs) Allocations(jobID string, allAllocs bool, q *QueryOptions) ([]*AllocationListStub, *QueryMeta, error) {
	var resp []*AllocationListStub
	u, err := url.Parse("/v1/job/" + url.PathEscape(jobID) + "/allocations")
	if err != nil {
		return nil, nil, err
	}

	v := u.Query()
	v.Add("all", strconv.FormatBool(allAllocs))
	u.RawQuery = v.Encode()

	qm, err := j.client.query(u.String(), &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(AllocIndexSort(resp))
	return resp, qm, nil
}

// Deployments is used to query the deployments associated with the given job
// ID.
func (j *Jobs) Deployments(jobID string, all bool, q *QueryOptions) ([]*Deployment, *QueryMeta, error) {
	var resp []*Deployment
	u, err := url.Parse("/v1/job/" + url.PathEscape(jobID) + "/deployments")
	if err != nil {
		return nil, nil, err
	}

	v := u.Query()
	v.Add("all", strconv.FormatBool(all))
	u.RawQuery = v.Encode()
	qm, err := j.client.query(u.String(), &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(DeploymentIndexSort(resp))
	return resp, qm, nil
}

// LatestDeployment is used to query for the latest deployment associated with
// the given job ID.
func (j *Jobs) LatestDeployment(jobID string, q *QueryOptions) (*Deployment, *QueryMeta, error) {
	var resp *Deployment
	qm, err := j.client.query("/v1/job/"+url.PathEscape(jobID)+"/deployment", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Evaluations is used to query the evaluations associated with the given job
// ID.
func (j *Jobs) Evaluations(jobID string, q *QueryOptions) ([]*Evaluation, *QueryMeta, error) {
	var resp []*Evaluation
	qm, err := j.client.query("/v1/job/"+url.PathEscape(jobID)+"/evaluations", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(EvalIndexSort(resp))
	return resp, qm, nil
}

// Deregister is used to remove an existing job. If purge is set to true, the job
// is deregistered and purged from the system versus still being queryable and
// eventually GC'ed from the system. Most callers should not specify purge.
func (j *Jobs) Deregister(jobID string, purge bool, q *WriteOptions) (string, *WriteMeta, error) {
	var resp JobDeregisterResponse
	wm, err := j.client.delete(fmt.Sprintf("/v1/job/%v?purge=%t", url.PathEscape(jobID), purge), nil, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

// DeregisterOptions is used to pass through job deregistration parameters
type DeregisterOptions struct {
	// If Purge is set to true, the job is deregistered and purged from the
	// system versus still being queryable and eventually GC'ed from the
	// system. Most callers should not specify purge.
	Purge bool

	// If Global is set to true, all regions of a multiregion job will be
	// stopped.
	Global bool

	// EvalPriority is an optional priority to use on any evaluation created as
	// a result on this job deregistration. This value must be between 1-100
	// inclusively, where a larger value corresponds to a higher priority. This
	// is useful when an operator wishes to push through a job deregistration
	// in busy clusters with a large evaluation backlog.
	EvalPriority int

	// NoShutdownDelay, if set to true, will override the group and
	// task shutdown_delay configuration and ignore the delay for any
	// allocations stopped as a result of this Deregister call.
	NoShutdownDelay bool
}

// DeregisterOpts is used to remove an existing job. See DeregisterOptions
// for parameters.
func (j *Jobs) DeregisterOpts(jobID string, opts *DeregisterOptions, q *WriteOptions) (string, *WriteMeta, error) {
	var resp JobDeregisterResponse

	// The base endpoint to add query params to.
	endpoint := "/v1/job/" + url.PathEscape(jobID)

	// Protect against nil opts. url.Values expects a string, and so using
	// fmt.Sprintf is the best way to do this.
	if opts != nil {
		endpoint += fmt.Sprintf("?purge=%t&global=%t&eval_priority=%v&no_shutdown_delay=%t",
			opts.Purge, opts.Global, opts.EvalPriority, opts.NoShutdownDelay)
	}

	wm, err := j.client.delete(endpoint, nil, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

// ForceEvaluate is used to force-evaluate an existing job.
func (j *Jobs) ForceEvaluate(jobID string, q *WriteOptions) (string, *WriteMeta, error) {
	var resp JobRegisterResponse
	wm, err := j.client.put("/v1/job/"+url.PathEscape(jobID)+"/evaluate", nil, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

// EvaluateWithOpts is used to force-evaluate an existing job and takes additional options
// for whether to force reschedule failed allocations
func (j *Jobs) EvaluateWithOpts(jobID string, opts EvalOptions, q *WriteOptions) (string, *WriteMeta, error) {
	req := &JobEvaluateRequest{
		JobID:       jobID,
		EvalOptions: opts,
	}

	var resp JobRegisterResponse
	wm, err := j.client.put("/v1/job/"+url.PathEscape(jobID)+"/evaluate", req, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

// PeriodicForce spawns a new instance of the periodic job and returns the eval ID
func (j *Jobs) PeriodicForce(jobID string, q *WriteOptions) (string, *WriteMeta, error) {
	var resp periodicForceResponse
	wm, err := j.client.put("/v1/job/"+url.PathEscape(jobID)+"/periodic/force", nil, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

// PlanOptions is used to pass through job planning parameters
type PlanOptions struct {
	Diff           bool
	PolicyOverride bool
}

func (j *Jobs) Plan(job *Job, diff bool, q *WriteOptions) (*JobPlanResponse, *WriteMeta, error) {
	opts := PlanOptions{Diff: diff}
	return j.PlanOpts(job, &opts, q)
}

func (j *Jobs) PlanOpts(job *Job, opts *PlanOptions, q *WriteOptions) (*JobPlanResponse, *WriteMeta, error) {
	if job == nil {
		return nil, nil, errors.New("must pass non-nil job")
	}
	if job.ID == nil {
		return nil, nil, errors.New("job is missing ID")
	}

	// Setup the request
	req := &JobPlanRequest{
		Job: job,
	}
	if opts != nil {
		req.Diff = opts.Diff
		req.PolicyOverride = opts.PolicyOverride
	}

	var resp JobPlanResponse
	wm, err := j.client.put("/v1/job/"+url.PathEscape(*job.ID)+"/plan", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

func (j *Jobs) Summary(jobID string, q *QueryOptions) (*JobSummary, *QueryMeta, error) {
	var resp JobSummary
	qm, err := j.client.query("/v1/job/"+url.PathEscape(jobID)+"/summary", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

func (j *Jobs) Dispatch(jobID string, meta map[string]string,
	payload []byte, idPrefixTemplate string, q *WriteOptions) (*JobDispatchResponse, *WriteMeta, error) {
	var resp JobDispatchResponse
	req := &JobDispatchRequest{
		JobID:            jobID,
		Meta:             meta,
		Payload:          payload,
		IdPrefixTemplate: idPrefixTemplate,
	}
	wm, err := j.client.put("/v1/job/"+url.PathEscape(jobID)+"/dispatch", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Revert is used to revert the given job to the passed version. If
// enforceVersion is set, the job is only reverted if the current version is at
// the passed version.
func (j *Jobs) Revert(jobID string, version uint64, enforcePriorVersion *uint64,
	q *WriteOptions, consulToken, vaultToken string) (*JobRegisterResponse, *WriteMeta, error) {

	var resp JobRegisterResponse
	req := &JobRevertRequest{
		JobID:               jobID,
		JobVersion:          version,
		EnforcePriorVersion: enforcePriorVersion,
		ConsulToken:         consulToken,
		VaultToken:          vaultToken,
	}
	wm, err := j.client.put("/v1/job/"+url.PathEscape(jobID)+"/revert", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Stable is used to mark a job version's stability.
func (j *Jobs) Stable(jobID string, version uint64, stable bool,
	q *WriteOptions) (*JobStabilityResponse, *WriteMeta, error) {

	var resp JobStabilityResponse
	req := &JobStabilityRequest{
		JobID:      jobID,
		JobVersion: version,
		Stable:     stable,
	}
	wm, err := j.client.put("/v1/job/"+url.PathEscape(jobID)+"/stable", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Services is used to return a list of service registrations associated to the
// specified jobID.
func (j *Jobs) Services(jobID string, q *QueryOptions) ([]*ServiceRegistration, *QueryMeta, error) {
	var resp []*ServiceRegistration
	qm, err := j.client.query("/v1/job/"+jobID+"/services", &resp, q)
	return resp, qm, err
}

// periodicForceResponse is used to deserialize a force response
type periodicForceResponse struct {
	EvalID string
}

// UpdateStrategy defines a task groups update strategy.
type UpdateStrategy struct {
	Stagger          *time.Duration `mapstructure:"stagger" hcl:"stagger,optional"`
	MaxParallel      *int           `mapstructure:"max_parallel" hcl:"max_parallel,optional"`
	HealthCheck      *string        `mapstructure:"health_check" hcl:"health_check,optional"`
	MinHealthyTime   *time.Duration `mapstructure:"min_healthy_time" hcl:"min_healthy_time,optional"`
	HealthyDeadline  *time.Duration `mapstructure:"healthy_deadline" hcl:"healthy_deadline,optional"`
	ProgressDeadline *time.Duration `mapstructure:"progress_deadline" hcl:"progress_deadline,optional"`
	Canary           *int           `mapstructure:"canary" hcl:"canary,optional"`
	AutoRevert       *bool          `mapstructure:"auto_revert" hcl:"auto_revert,optional"`
	AutoPromote      *bool          `mapstructure:"auto_promote" hcl:"auto_promote,optional"`
}

// DefaultUpdateStrategy provides a baseline that can be used to upgrade
// jobs with the old policy or for populating field defaults.
func DefaultUpdateStrategy() *UpdateStrategy {
	return &UpdateStrategy{
		Stagger:          pointerOf(30 * time.Second),
		MaxParallel:      pointerOf(1),
		HealthCheck:      pointerOf("checks"),
		MinHealthyTime:   pointerOf(10 * time.Second),
		HealthyDeadline:  pointerOf(5 * time.Minute),
		ProgressDeadline: pointerOf(10 * time.Minute),
		AutoRevert:       pointerOf(false),
		Canary:           pointerOf(0),
		AutoPromote:      pointerOf(false),
	}
}

func (u *UpdateStrategy) Copy() *UpdateStrategy {
	if u == nil {
		return nil
	}

	copy := new(UpdateStrategy)

	if u.Stagger != nil {
		copy.Stagger = pointerOf(*u.Stagger)
	}

	if u.MaxParallel != nil {
		copy.MaxParallel = pointerOf(*u.MaxParallel)
	}

	if u.HealthCheck != nil {
		copy.HealthCheck = pointerOf(*u.HealthCheck)
	}

	if u.MinHealthyTime != nil {
		copy.MinHealthyTime = pointerOf(*u.MinHealthyTime)
	}

	if u.HealthyDeadline != nil {
		copy.HealthyDeadline = pointerOf(*u.HealthyDeadline)
	}

	if u.ProgressDeadline != nil {
		copy.ProgressDeadline = pointerOf(*u.ProgressDeadline)
	}

	if u.AutoRevert != nil {
		copy.AutoRevert = pointerOf(*u.AutoRevert)
	}

	if u.Canary != nil {
		copy.Canary = pointerOf(*u.Canary)
	}

	if u.AutoPromote != nil {
		copy.AutoPromote = pointerOf(*u.AutoPromote)
	}

	return copy
}

func (u *UpdateStrategy) Merge(o *UpdateStrategy) {
	if o == nil {
		return
	}

	if o.Stagger != nil {
		u.Stagger = pointerOf(*o.Stagger)
	}

	if o.MaxParallel != nil {
		u.MaxParallel = pointerOf(*o.MaxParallel)
	}

	if o.HealthCheck != nil {
		u.HealthCheck = pointerOf(*o.HealthCheck)
	}

	if o.MinHealthyTime != nil {
		u.MinHealthyTime = pointerOf(*o.MinHealthyTime)
	}

	if o.HealthyDeadline != nil {
		u.HealthyDeadline = pointerOf(*o.HealthyDeadline)
	}

	if o.ProgressDeadline != nil {
		u.ProgressDeadline = pointerOf(*o.ProgressDeadline)
	}

	if o.AutoRevert != nil {
		u.AutoRevert = pointerOf(*o.AutoRevert)
	}

	if o.Canary != nil {
		u.Canary = pointerOf(*o.Canary)
	}

	if o.AutoPromote != nil {
		u.AutoPromote = pointerOf(*o.AutoPromote)
	}
}

func (u *UpdateStrategy) Canonicalize() {
	d := DefaultUpdateStrategy()

	if u.MaxParallel == nil {
		u.MaxParallel = d.MaxParallel
	}

	if u.Stagger == nil {
		u.Stagger = d.Stagger
	}

	if u.HealthCheck == nil {
		u.HealthCheck = d.HealthCheck
	}

	if u.HealthyDeadline == nil {
		u.HealthyDeadline = d.HealthyDeadline
	}

	if u.ProgressDeadline == nil {
		u.ProgressDeadline = d.ProgressDeadline
	}

	if u.MinHealthyTime == nil {
		u.MinHealthyTime = d.MinHealthyTime
	}

	if u.AutoRevert == nil {
		u.AutoRevert = d.AutoRevert
	}

	if u.Canary == nil {
		u.Canary = d.Canary
	}

	if u.AutoPromote == nil {
		u.AutoPromote = d.AutoPromote
	}
}

// Empty returns whether the UpdateStrategy is empty or has user defined values.
func (u *UpdateStrategy) Empty() bool {
	if u == nil {
		return true
	}

	if u.Stagger != nil && *u.Stagger != 0 {
		return false
	}

	if u.MaxParallel != nil && *u.MaxParallel != 0 {
		return false
	}

	if u.HealthCheck != nil && *u.HealthCheck != "" {
		return false
	}

	if u.MinHealthyTime != nil && *u.MinHealthyTime != 0 {
		return false
	}

	if u.HealthyDeadline != nil && *u.HealthyDeadline != 0 {
		return false
	}

	if u.ProgressDeadline != nil && *u.ProgressDeadline != 0 {
		return false
	}

	if u.AutoRevert != nil && *u.AutoRevert {
		return false
	}

	if u.AutoPromote != nil && *u.AutoPromote {
		return false
	}

	if u.Canary != nil && *u.Canary != 0 {
		return false
	}

	return true
}

type Multiregion struct {
	Strategy *MultiregionStrategy `hcl:"strategy,block"`
	Regions  []*MultiregionRegion `hcl:"region,block"`
}

func (m *Multiregion) Canonicalize() {
	if m.Strategy == nil {
		m.Strategy = &MultiregionStrategy{
			MaxParallel: pointerOf(0),
			OnFailure:   pointerOf(""),
		}
	} else {
		if m.Strategy.MaxParallel == nil {
			m.Strategy.MaxParallel = pointerOf(0)
		}
		if m.Strategy.OnFailure == nil {
			m.Strategy.OnFailure = pointerOf("")
		}
	}
	if m.Regions == nil {
		m.Regions = []*MultiregionRegion{}
	}
	for _, region := range m.Regions {
		if region.Count == nil {
			region.Count = pointerOf(1)
		}
		if region.Datacenters == nil {
			region.Datacenters = []string{}
		}
		if region.Meta == nil {
			region.Meta = map[string]string{}
		}
	}
}

func (m *Multiregion) Copy() *Multiregion {
	if m == nil {
		return nil
	}
	copy := new(Multiregion)
	if m.Strategy != nil {
		copy.Strategy = new(MultiregionStrategy)
		copy.Strategy.MaxParallel = pointerOf(*m.Strategy.MaxParallel)
		copy.Strategy.OnFailure = pointerOf(*m.Strategy.OnFailure)
	}
	for _, region := range m.Regions {
		copyRegion := new(MultiregionRegion)
		copyRegion.Name = region.Name
		copyRegion.Count = pointerOf(*region.Count)
		copyRegion.Datacenters = append(copyRegion.Datacenters, region.Datacenters...)
		copyRegion.NodePool = region.NodePool
		for k, v := range region.Meta {
			copyRegion.Meta[k] = v
		}

		copy.Regions = append(copy.Regions, copyRegion)
	}
	return copy
}

type MultiregionStrategy struct {
	MaxParallel *int    `mapstructure:"max_parallel" hcl:"max_parallel,optional"`
	OnFailure   *string `mapstructure:"on_failure" hcl:"on_failure,optional"`
}

type MultiregionRegion struct {
	Name        string            `hcl:",label"`
	Count       *int              `hcl:"count,optional"`
	Datacenters []string          `hcl:"datacenters,optional"`
	NodePool    string            `hcl:"node_pool,optional"`
	Meta        map[string]string `hcl:"meta,block"`
}

// PeriodicConfig is for serializing periodic config for a job.
type PeriodicConfig struct {
	Enabled         *bool    `hcl:"enabled,optional"`
	Spec            *string  `hcl:"cron,optional"`
	Specs           []string `hcl:"crons,optional"`
	SpecType        *string
	ProhibitOverlap *bool   `mapstructure:"prohibit_overlap" hcl:"prohibit_overlap,optional"`
	TimeZone        *string `mapstructure:"time_zone" hcl:"time_zone,optional"`
}

func (p *PeriodicConfig) Canonicalize() {
	if p.Enabled == nil {
		p.Enabled = pointerOf(true)
	}
	if p.Spec == nil {
		p.Spec = pointerOf("")
	}
	if p.Specs == nil {
		p.Specs = []string{}
	}
	if p.SpecType == nil {
		p.SpecType = pointerOf(PeriodicSpecCron)
	}
	if p.ProhibitOverlap == nil {
		p.ProhibitOverlap = pointerOf(false)
	}
	if p.TimeZone == nil || *p.TimeZone == "" {
		p.TimeZone = pointerOf("UTC")
	}
}

// Next returns the closest time instant matching the spec that is after the
// passed time. If no matching instance exists, the zero value of time.Time is
// returned. The `time.Location` of the returned value matches that of the
// passed time.
func (p *PeriodicConfig) Next(fromTime time.Time) (time.Time, error) {
	// Single spec parsing
	if p != nil && *p.SpecType == PeriodicSpecCron {
		if p.Spec != nil && *p.Spec != "" {
			return cronParseNext(fromTime, *p.Spec)
		}
	}

	// multiple specs parsing
	var nextTime time.Time
	for _, spec := range p.Specs {
		t, err := cronParseNext(fromTime, spec)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed parsing cron expression %s: %v", spec, err)
		}
		if nextTime.IsZero() || t.Before(nextTime) {
			nextTime = t
		}
	}
	return nextTime, nil
}

// cronParseNext is a helper that parses the next time for the given expression
// but captures any panic that may occur in the underlying library.
// ---  THIS FUNCTION IS REPLICATED IN nomad/structs/structs.go
// and should be kept in sync.
func cronParseNext(fromTime time.Time, spec string) (t time.Time, err error) {
	defer func() {
		if recover() != nil {
			t = time.Time{}
			err = fmt.Errorf("failed parsing cron expression: %q", spec)
		}
	}()
	exp, err := cronexpr.Parse(spec)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed parsing cron expression: %s: %v", spec, err)
	}
	return exp.Next(fromTime), nil
}

func (p *PeriodicConfig) GetLocation() (*time.Location, error) {
	if p.TimeZone == nil || *p.TimeZone == "" {
		return time.UTC, nil
	}

	return time.LoadLocation(*p.TimeZone)
}

// ParameterizedJobConfig is used to configure the parameterized job.
type ParameterizedJobConfig struct {
	Payload      string   `hcl:"payload,optional"`
	MetaRequired []string `mapstructure:"meta_required" hcl:"meta_required,optional"`
	MetaOptional []string `mapstructure:"meta_optional" hcl:"meta_optional,optional"`
}

// JobSubmission is used to hold information about the original content of a job
// specification being submitted to Nomad.
//
// At any time a JobSubmission may be nil, indicating no information is known about
// the job submission.
type JobSubmission struct {
	// Source contains the original job definition (may be in the format of
	// hcl1, hcl2, or json). HCL1 jobs can no longer be parsed.
	Source string

	// Format indicates what the Source content was (hcl1, hcl2, or json). HCL1
	// jobs can no longer be parsed.
	Format string

	// VariableFlags contains the CLI "-var" flag arguments as submitted with the
	// job (hcl2 only).
	VariableFlags map[string]string

	// Variables contains the opaque variables configuration as coming from
	// a var-file or the WebUI variables input (hcl2 only).
	Variables string
}

type JobUIConfig struct {
	Description string       `hcl:"description,optional"`
	Links       []*JobUILink `hcl:"link,block"`
}

type JobUILink struct {
	Label string `hcl:"label,optional"`
	URL   string `hcl:"url,optional"`
}

func (j *JobUIConfig) Canonicalize() {
	if j == nil {
		return
	}

	if len(j.Links) == 0 {
		j.Links = nil
	}
}

func (j *JobUIConfig) Copy() *JobUIConfig {
	if j == nil {
		return nil
	}

	copy := new(JobUIConfig)
	copy.Description = j.Description

	for _, link := range j.Links {
		copy.Links = append(copy.Links, link.Copy())
	}

	return copy
}

func (j *JobUILink) Copy() *JobUILink {
	if j == nil {
		return nil
	}

	return &JobUILink{
		Label: j.Label,
		URL:   j.URL,
	}
}

type JobVersionTag struct {
	Name        string
	Description string
	TaggedTime  int64
}

func (j *JobVersionTag) Copy() *JobVersionTag {
	if j == nil {
		return nil
	}

	return &JobVersionTag{
		Name:        j.Name,
		Description: j.Description,
		TaggedTime:  j.TaggedTime,
	}
}

func (js *JobSubmission) Canonicalize() {
	if js == nil {
		return
	}

	if len(js.VariableFlags) == 0 {
		js.VariableFlags = nil
	}

	// if there are multiline variables, make sure we escape the newline
	// characters to preserve them. This way, when the job gets stopped and
	// restarted in the UI, variable values will be parsed correctly.
	for k, v := range js.VariableFlags {
		if strings.Contains(v, "\n") {
			js.VariableFlags[k] = strings.ReplaceAll(v, "\n", "\\n")
		}
	}
}

func (js *JobSubmission) Copy() *JobSubmission {
	if js == nil {
		return nil
	}

	return &JobSubmission{
		Source:        js.Source,
		Format:        js.Format,
		VariableFlags: maps.Clone(js.VariableFlags),
		Variables:     js.Variables,
	}
}

// Job is used to serialize a job.
type Job struct {
	/* Fields parsed from HCL config */

	Region           *string                 `hcl:"region,optional"`
	Namespace        *string                 `hcl:"namespace,optional"`
	ID               *string                 `hcl:"id,optional"`
	Name             *string                 `hcl:"name,optional"`
	Type             *string                 `hcl:"type,optional"`
	Priority         *int                    `hcl:"priority,optional"`
	AllAtOnce        *bool                   `mapstructure:"all_at_once" hcl:"all_at_once,optional"`
	Datacenters      []string                `hcl:"datacenters,optional"`
	NodePool         *string                 `mapstructure:"node_pool" hcl:"node_pool,optional"`
	Constraints      []*Constraint           `hcl:"constraint,block"`
	Affinities       []*Affinity             `hcl:"affinity,block"`
	TaskGroups       []*TaskGroup            `hcl:"group,block"`
	Update           *UpdateStrategy         `hcl:"update,block"`
	Multiregion      *Multiregion            `hcl:"multiregion,block"`
	Spreads          []*Spread               `hcl:"spread,block"`
	Periodic         *PeriodicConfig         `hcl:"periodic,block"`
	ParameterizedJob *ParameterizedJobConfig `hcl:"parameterized,block"`
	Reschedule       *ReschedulePolicy       `hcl:"reschedule,block"`
	Migrate          *MigrateStrategy        `hcl:"migrate,block"`
	Meta             map[string]string       `hcl:"meta,block"`
	ConsulToken      *string                 `mapstructure:"consul_token" hcl:"consul_token,optional"`
	VaultToken       *string                 `mapstructure:"vault_token" hcl:"vault_token,optional"`
	UI               *JobUIConfig            `hcl:"ui,block"`

	/* Fields set by server, not sourced from job config file */

	Stop                     *bool
	ParentID                 *string
	Dispatched               bool
	DispatchIdempotencyToken *string
	Payload                  []byte
	ConsulNamespace          *string `mapstructure:"consul_namespace"`
	VaultNamespace           *string `mapstructure:"vault_namespace"`
	NomadTokenID             *string `mapstructure:"nomad_token_id"`
	Status                   *string
	StatusDescription        *string
	Stable                   *bool
	Version                  *uint64
	SubmitTime               *int64
	CreateIndex              *uint64
	ModifyIndex              *uint64
	JobModifyIndex           *uint64
	VersionTag               *JobVersionTag
}

// IsPeriodic returns whether a job is periodic.
func (j *Job) IsPeriodic() bool {
	return j.Periodic != nil
}

// IsParameterized returns whether a job is parameterized job.
func (j *Job) IsParameterized() bool {
	return j.ParameterizedJob != nil && !j.Dispatched
}

// IsMultiregion returns whether a job is a multiregion job
func (j *Job) IsMultiregion() bool {
	return j.Multiregion != nil && j.Multiregion.Regions != nil && len(j.Multiregion.Regions) > 0
}

func (j *Job) Canonicalize() {
	if j.ID == nil {
		j.ID = pointerOf("")
	}
	if j.Name == nil {
		j.Name = pointerOf(*j.ID)
	}
	if j.ParentID == nil {
		j.ParentID = pointerOf("")
	}
	if j.Namespace == nil {
		j.Namespace = pointerOf(DefaultNamespace)
	}
	if j.Priority == nil {
		j.Priority = pointerOf(JobDefaultPriority)
	}
	if j.Stop == nil {
		j.Stop = pointerOf(false)
	}
	if j.Region == nil {
		j.Region = pointerOf(GlobalRegion)
	}
	if j.NodePool == nil {
		j.NodePool = pointerOf("")
	}
	if j.Type == nil {
		j.Type = pointerOf("service")
	}
	if j.AllAtOnce == nil {
		j.AllAtOnce = pointerOf(false)
	}
	if j.ConsulToken == nil {
		j.ConsulToken = pointerOf("")
	}
	if j.ConsulNamespace == nil {
		j.ConsulNamespace = pointerOf("")
	}
	if j.VaultToken == nil {
		j.VaultToken = pointerOf("")
	}
	if j.VaultNamespace == nil {
		j.VaultNamespace = pointerOf("")
	}
	if j.NomadTokenID == nil {
		j.NomadTokenID = pointerOf("")
	}
	if j.Status == nil {
		j.Status = pointerOf("")
	}
	if j.StatusDescription == nil {
		j.StatusDescription = pointerOf("")
	}
	if j.Stable == nil {
		j.Stable = pointerOf(false)
	}
	if j.Version == nil {
		j.Version = pointerOf(uint64(0))
	}
	if j.CreateIndex == nil {
		j.CreateIndex = pointerOf(uint64(0))
	}
	if j.ModifyIndex == nil {
		j.ModifyIndex = pointerOf(uint64(0))
	}
	if j.JobModifyIndex == nil {
		j.JobModifyIndex = pointerOf(uint64(0))
	}
	if j.Periodic != nil {
		j.Periodic.Canonicalize()
	}
	if j.Update != nil {
		j.Update.Canonicalize()
	} else if *j.Type == JobTypeService {
		j.Update = DefaultUpdateStrategy()
	}
	if j.Multiregion != nil {
		j.Multiregion.Canonicalize()
	}

	for _, tg := range j.TaskGroups {
		tg.Canonicalize(j)
	}

	for _, spread := range j.Spreads {
		spread.Canonicalize()
	}
	for _, a := range j.Affinities {
		a.Canonicalize()
	}

	if j.UI != nil {
		j.UI.Canonicalize()
	}
}

// LookupTaskGroup finds a task group by name
func (j *Job) LookupTaskGroup(name string) *TaskGroup {
	for _, tg := range j.TaskGroups {
		if *tg.Name == name {
			return tg
		}
	}
	return nil
}

// JobSummary summarizes the state of the allocations of a job
type JobSummary struct {
	JobID     string
	Namespace string
	Summary   map[string]TaskGroupSummary
	Children  *JobChildrenSummary

	// Raft Indexes
	CreateIndex uint64
	ModifyIndex uint64
}

// JobChildrenSummary contains the summary of children job status
type JobChildrenSummary struct {
	Pending int64
	Running int64
	Dead    int64
}

func (jc *JobChildrenSummary) Sum() int {
	if jc == nil {
		return 0
	}

	return int(jc.Pending + jc.Running + jc.Dead)
}

// TaskGroup summarizes the state of all the allocations of a particular
// TaskGroup
type TaskGroupSummary struct {
	Queued   int
	Complete int
	Failed   int
	Running  int
	Starting int
	Lost     int
	Unknown  int
}

// JobListStub is used to return a subset of information about
// jobs during list operations.
type JobListStub struct {
	ID                string
	ParentID          string
	Name              string
	Namespace         string `json:",omitempty"`
	Datacenters       []string
	Type              string
	Priority          int
	Periodic          bool
	ParameterizedJob  bool
	Stop              bool
	Status            string
	StatusDescription string
	JobSummary        *JobSummary
	CreateIndex       uint64
	ModifyIndex       uint64
	JobModifyIndex    uint64
	SubmitTime        int64
	Meta              map[string]string `json:",omitempty"`
}

// JobIDSort is used to sort jobs by their job ID's.
type JobIDSort []*JobListStub

func (j JobIDSort) Len() int {
	return len(j)
}

func (j JobIDSort) Less(a, b int) bool {
	return j[a].ID < j[b].ID
}

func (j JobIDSort) Swap(a, b int) {
	j[a], j[b] = j[b], j[a]
}

// NewServiceJob creates and returns a new service-style job
// for long-lived processes using the provided name, ID, and
// relative job priority.
func NewServiceJob(id, name, region string, pri int) *Job {
	return newJob(id, name, region, JobTypeService, pri)
}

// NewBatchJob creates and returns a new batch-style job for
// short-lived processes using the provided name and ID along
// with the relative job priority.
func NewBatchJob(id, name, region string, pri int) *Job {
	return newJob(id, name, region, JobTypeBatch, pri)
}

// NewSystemJob creates and returns a new system-style job for processes
// designed to run on all clients, using the provided name and ID along with
// the relative job priority.
func NewSystemJob(id, name, region string, pri int) *Job {
	return newJob(id, name, region, JobTypeSystem, pri)
}

// NewSysbatchJob creates and returns a new sysbatch-style job for short-lived
// processes designed to run on all clients, using the provided name and ID
// along with the relative job priority.
func NewSysbatchJob(id, name, region string, pri int) *Job {
	return newJob(id, name, region, JobTypeSysbatch, pri)
}

// newJob is used to create a new Job struct.
func newJob(id, name, region, typ string, pri int) *Job {
	return &Job{
		Region:   &region,
		ID:       &id,
		Name:     &name,
		Type:     &typ,
		Priority: &pri,
	}
}

// SetMeta is used to set arbitrary k/v pairs of metadata on a job.
func (j *Job) SetMeta(key, val string) *Job {
	if j.Meta == nil {
		j.Meta = make(map[string]string)
	}
	j.Meta[key] = val
	return j
}

// AddDatacenter is used to add a datacenter to a job.
func (j *Job) AddDatacenter(dc string) *Job {
	j.Datacenters = append(j.Datacenters, dc)
	return j
}

// Constrain is used to add a constraint to a job.
func (j *Job) Constrain(c *Constraint) *Job {
	j.Constraints = append(j.Constraints, c)
	return j
}

// AddAffinity is used to add an affinity to a job.
func (j *Job) AddAffinity(a *Affinity) *Job {
	j.Affinities = append(j.Affinities, a)
	return j
}

// AddTaskGroup adds a task group to an existing job.
func (j *Job) AddTaskGroup(grp *TaskGroup) *Job {
	j.TaskGroups = append(j.TaskGroups, grp)
	return j
}

// AddPeriodicConfig adds a periodic config to an existing job.
func (j *Job) AddPeriodicConfig(cfg *PeriodicConfig) *Job {
	j.Periodic = cfg
	return j
}

func (j *Job) AddSpread(s *Spread) *Job {
	j.Spreads = append(j.Spreads, s)
	return j
}

type WriteRequest struct {
	// The target region for this write
	Region string

	// Namespace is the target namespace for this write
	Namespace string

	// SecretID is the secret ID of an ACL token
	SecretID string
}

// JobValidateRequest is used to validate a job
type JobValidateRequest struct {
	Job *Job
	WriteRequest
}

// JobValidateResponse is the response from validate request
type JobValidateResponse struct {
	// DriverConfigValidated indicates whether the agent validated the driver
	// config
	DriverConfigValidated bool

	// ValidationErrors is a list of validation errors
	ValidationErrors []string

	// Error is a string version of any error that may have occurred
	Error string

	// Warnings contains any warnings about the given job. These may include
	// deprecation warnings.
	Warnings string
}

// JobRevertRequest is used to revert a job to a prior version.
type JobRevertRequest struct {
	// JobID is the ID of the job  being reverted
	JobID string

	// JobVersion the version to revert to.
	JobVersion uint64

	// EnforcePriorVersion if set will enforce that the job is at the given
	// version before reverting.
	EnforcePriorVersion *uint64

	// ConsulToken is the Consul token that proves the submitter of the job revert
	// has access to the Service Identity policies associated with the job's
	// Consul Connect enabled services. This field is only used to transfer the
	// token and is not stored after the Job revert.
	ConsulToken string `json:",omitempty"`

	// VaultToken is the Vault token that proves the submitter of the job revert
	// has access to any Vault policies specified in the targeted job version. This
	// field is only used to authorize the revert and is not stored after the Job
	// revert.
	VaultToken string `json:",omitempty"`

	WriteRequest
}

// JobRegisterRequest is used to update a job
type JobRegisterRequest struct {
	Submission *JobSubmission
	Job        *Job

	// If EnforceIndex is set then the job will only be registered if the passed
	// JobModifyIndex matches the current Jobs index. If the index is zero, the
	// register only occurs if the job is new.
	EnforceIndex   bool   `json:",omitempty"`
	JobModifyIndex uint64 `json:",omitempty"`
	PolicyOverride bool   `json:",omitempty"`
	PreserveCounts bool   `json:",omitempty"`

	// EvalPriority is an optional priority to use on any evaluation created as
	// a result on this job registration. This value must be between 1-100
	// inclusively, where a larger value corresponds to a higher priority. This
	// is useful when an operator wishes to push through a job registration in
	// busy clusters with a large evaluation backlog. This avoids needing to
	// change the job priority which also impacts preemption.
	EvalPriority int `json:",omitempty"`

	WriteRequest
}

// JobRegisterResponse is used to respond to a job registration
type JobRegisterResponse struct {
	EvalID          string
	EvalCreateIndex uint64
	JobModifyIndex  uint64

	// Warnings contains any warnings about the given job. These may include
	// deprecation warnings.
	Warnings string

	QueryMeta
}

// JobDeregisterResponse is used to respond to a job deregistration
type JobDeregisterResponse struct {
	EvalID          string
	EvalCreateIndex uint64
	JobModifyIndex  uint64
	QueryMeta
}

type JobPlanRequest struct {
	Job            *Job
	Diff           bool
	PolicyOverride bool
	WriteRequest
}

type JobPlanResponse struct {
	JobModifyIndex     uint64
	CreatedEvals       []*Evaluation
	Diff               *JobDiff
	Annotations        *PlanAnnotations
	FailedTGAllocs     map[string]*AllocationMetric
	NextPeriodicLaunch time.Time

	// Warnings contains any warnings about the given job. These may include
	// deprecation warnings.
	Warnings string
}

type JobDiff struct {
	Type       string
	ID         string
	Fields     []*FieldDiff
	Objects    []*ObjectDiff
	TaskGroups []*TaskGroupDiff
}

type TaskGroupDiff struct {
	Type    string
	Name    string
	Fields  []*FieldDiff
	Objects []*ObjectDiff
	Tasks   []*TaskDiff
	Updates map[string]uint64
}

type TaskDiff struct {
	Type        string
	Name        string
	Fields      []*FieldDiff
	Objects     []*ObjectDiff
	Annotations []string
}

type FieldDiff struct {
	Type        string
	Name        string
	Old, New    string
	Annotations []string
}

type ObjectDiff struct {
	Type    string
	Name    string
	Fields  []*FieldDiff
	Objects []*ObjectDiff
}

type PlanAnnotations struct {
	DesiredTGUpdates map[string]*DesiredUpdates
	PreemptedAllocs  []*AllocationListStub
}

type DesiredUpdates struct {
	Ignore            uint64
	Place             uint64
	Migrate           uint64
	Stop              uint64
	InPlaceUpdate     uint64
	DestructiveUpdate uint64
	Canary            uint64
	Preemptions       uint64
}

type JobDispatchRequest struct {
	JobID            string
	Payload          []byte
	Meta             map[string]string
	IdPrefixTemplate string
}

type JobDispatchResponse struct {
	DispatchedJobID string
	EvalID          string
	EvalCreateIndex uint64
	JobCreateIndex  uint64
	WriteMeta
}

// JobVersionsResponse is used for a job get versions request
type JobVersionsResponse struct {
	Versions []*Job
	Diffs    []*JobDiff
	QueryMeta
}

// JobSubmissionResponse is used for a job get submission request
type JobSubmissionResponse struct {
	Submission *JobSubmission
	QueryMeta
}

// JobStabilityRequest is used to marked a job as stable.
type JobStabilityRequest struct {
	// Job to set the stability on
	JobID      string
	JobVersion uint64

	// Set the stability
	Stable bool
	WriteRequest
}

// JobStabilityResponse is the response when marking a job as stable.
type JobStabilityResponse struct {
	JobModifyIndex uint64
	WriteMeta
}

// JobEvaluateRequest is used when we just need to re-evaluate a target job
type JobEvaluateRequest struct {
	JobID       string
	EvalOptions EvalOptions
	WriteRequest
}

// EvalOptions is used to encapsulate options when forcing a job evaluation
type EvalOptions struct {
	ForceReschedule bool
}

// ActionExec is used to run a pre-defined command inside a running task.
// The call blocks until command terminates (or an error occurs), and returns the exit code.
func (j *Jobs) ActionExec(ctx context.Context,
	alloc *Allocation, job string, task string, tty bool, command []string,
	action string,
	stdin io.Reader, stdout, stderr io.Writer,
	terminalSizeCh <-chan TerminalSize, q *QueryOptions) (exitCode int, err error) {

	s := &execSession{
		client:  j.client,
		alloc:   alloc,
		job:     job,
		task:    task,
		tty:     tty,
		command: command,
		action:  action,

		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,

		terminalSizeCh: terminalSizeCh,
		q:              q,
	}

	return s.run(ctx)
}

// JobStatusesRequest is used to get statuses for jobs,
// their allocations and deployments.
type JobStatusesRequest struct {
	// Jobs may be optionally provided to request a subset of specific jobs.
	Jobs []NamespacedID
	// IncludeChildren will include child (batch) jobs in the response.
	IncludeChildren bool
}

type TagVersionRequest struct {
	Version     uint64
	Description string
	WriteRequest
}

func (j *Jobs) TagVersion(jobID string, version uint64, name string, description string, q *WriteOptions) (*WriteMeta, error) {
	var tagRequest = &TagVersionRequest{
		Version:     version,
		Description: description,
	}

	return j.client.put("/v1/job/"+url.PathEscape(jobID)+"/versions/"+name+"/tag", tagRequest, nil, q)
}

func (j *Jobs) UntagVersion(jobID string, name string, q *WriteOptions) (*WriteMeta, error) {
	return j.client.delete("/v1/job/"+url.PathEscape(jobID)+"/versions/"+name+"/tag", nil, nil, q)
}
