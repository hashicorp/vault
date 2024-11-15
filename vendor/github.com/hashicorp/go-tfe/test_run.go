// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ TestRuns = (*testRuns)(nil)

// TestRuns describes all the test run related methods that the Terraform
// Enterprise API supports.
//
// **Note: These methods are still in BETA and subject to change.**
type TestRuns interface {
	// List all the test runs for a given private registry module.
	List(ctx context.Context, moduleID RegistryModuleID, options *TestRunListOptions) (*TestRunList, error)

	// Read a test run by its ID.
	Read(ctx context.Context, moduleID RegistryModuleID, testRunID string) (*TestRun, error)

	// Create a new test run with the given options.
	Create(ctx context.Context, options TestRunCreateOptions) (*TestRun, error)

	// Logs retrieves the logs for a test run by its ID.
	Logs(ctx context.Context, moduleID RegistryModuleID, testRunID string) (io.Reader, error)

	// Cancel a test run by its ID.
	Cancel(ctx context.Context, moduleID RegistryModuleID, testRunID string) error

	// ForceCancel a test run by its ID.
	ForceCancel(ctx context.Context, moduleID RegistryModuleID, testRunID string) error
}

// testRuns implements TestRuns.
type testRuns struct {
	client *Client
}

// TestRunStatus represents the status of a test run.
type TestRunStatus string

// List all available test run statuses.
const (
	TestRunPending  TestRunStatus = "pending"
	TestRunQueued   TestRunStatus = "queued"
	TestRunRunning  TestRunStatus = "running"
	TestRunErrored  TestRunStatus = "errored"
	TestRunCanceled TestRunStatus = "canceled"
	TestRunFinished TestRunStatus = "finished"
)

// TestStatus represents the status of an individual test within an overall test
// run.
type TestStatus string

// List all available test statuses.
const (
	TestPending TestStatus = "pending"
	TestSkip    TestStatus = "skip"
	TestPass    TestStatus = "pass"
	TestFail    TestStatus = "fail"
	TestError   TestStatus = "error"
)

// TestRun represents a Terraform Enterprise test run.
type TestRun struct {
	ID               string                  `jsonapi:"primary,test-runs"`
	Status           TestRunStatus           `jsonapi:"attr,status"`
	StatusTimestamps TestRunStatusTimestamps `jsonapi:"attr,status-timestamps"`
	TestStatus       TestStatus              `jsonapi:"attr,test-status"`
	TestsPassed      int                     `jsonapi:"attr,tests-passed"`
	TestsFailed      int                     `jsonapi:"attr,tests-failed"`
	TestsErrored     int                     `jsonapi:"attr,tests-errored"`
	TestsSkipped     int                     `jsonapi:"attr,tests-skipped"`
	LogReadURL       string                  `jsonapi:"attr,log-read-url"`

	// Relations
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`
	RegistryModule       *RegistryModule       `jsonapi:"relation,registry-module"`
}

// TestRunStatusTimestamps holds the timestamps for individual test run
// statuses.
type TestRunStatusTimestamps struct {
	CanceledAt      time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	ErroredAt       time.Time `jsonapi:"attr,errored-at,rfc3339"`
	FinishedAt      time.Time `jsonapi:"attr,finished-at,rfc3339"`
	ForceCanceledAt time.Time `jsonapi:"attr,force-canceled-at,rfc3339"`
	QueuedAt        time.Time `jsonapi:"attr,queued-at,rfc3339"`
	StartedAt       time.Time `jsonapi:"attr,started-at,rfc3339"`
}

// TestRunCreateOptions represents the options for creating a run.
type TestRunCreateOptions struct {
	// Type is a public field utitilized by JSON:API to set the resource type
	// via the field tag. It is not a user-defined value and does not need to
	// be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,test-runs"`

	// If non-empty, requests that only a subset of testing files within the
	// ConfigurationVersion should be executed.
	Filters []string `jsonapi:"attr,filters,omitempty"`

	// Specifies the directory within the ConfigurationVersion that test files
	// should be loaded from. Defaults to "tests" if empty.
	TestDirectory *string `jsonapi:"attr,test-directory,omitempty"`

	// Verbose prints out the plan and state files for each run block that is
	// executed by this TestRun.
	Verbose *bool `jsonapi:"attr,verbose,omitempty"`

	// Variables allows you to specify terraform input variables for
	// a particular run, prioritized over variables defined on the workspace.
	Variables []*RunVariable `jsonapi:"attr,variables,omitempty"`

	// ConfigurationVersion specifies the configuration version to use for this
	// test run.
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`

	// RegistryModule specifies the registry module this test run should be
	// assigned to.
	RegistryModule *RegistryModule `jsonapi:"relation,registry-module"`
}

// TestRunList represents a list of test runs.
type TestRunList struct {
	*Pagination
	Items []*TestRun
}

// TestRunListOptions represents the options for listing runs.
type TestRunListOptions struct {
	ListOptions
}

// List all the test runs for a given private registry module.
func (s *testRuns) List(ctx context.Context, moduleID RegistryModuleID, options *TestRunListOptions) (*TestRunList, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", testRunsPath(moduleID), options)
	if err != nil {
		return nil, err
	}

	trl := &TestRunList{}
	err = req.Do(ctx, trl)
	if err != nil {
		return nil, err
	}

	return trl, nil
}

// Read a test run by its ID.
func (s *testRuns) Read(ctx context.Context, moduleID RegistryModuleID, testRunID string) (*TestRun, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if !validStringID(&testRunID) {
		return nil, ErrInvalidTestRunID
	}

	u := fmt.Sprintf("%s/%s", testRunsPath(moduleID), url.PathEscape(testRunID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tr := &TestRun{}
	err = req.Do(ctx, tr)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

// Create a new test run with the given options.
func (s *testRuns) Create(ctx context.Context, options TestRunCreateOptions) (*TestRun, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	moduleID := RegistryModuleID{
		Organization: options.RegistryModule.Organization.Name,
		Name:         options.RegistryModule.Name,
		Provider:     options.RegistryModule.Provider,
		Namespace:    options.RegistryModule.Namespace,
		RegistryName: options.RegistryModule.RegistryName,
	}

	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", testRunsPath(moduleID), &options)
	if err != nil {
		return nil, err
	}

	tr := &TestRun{}
	err = req.Do(ctx, tr)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

// Logs retrieves the logs for a test run by its ID.
func (s *testRuns) Logs(ctx context.Context, moduleID RegistryModuleID, testRunID string) (io.Reader, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if !validStringID(&testRunID) {
		return nil, ErrInvalidTestRunID
	}

	tr, err := s.Read(ctx, moduleID, testRunID)
	if err != nil {
		return nil, err
	}

	if tr.LogReadURL == "" {
		return nil, fmt.Errorf("test run %s does not have a log URL", testRunID)
	}

	u, err := url.Parse(tr.LogReadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid log URL: %w", err)
	}

	done := func() (bool, error) {
		tr, err := s.Read(ctx, moduleID, testRunID)
		if err != nil {
			return false, err
		}

		switch tr.Status {
		case TestRunErrored, TestRunCanceled, TestRunFinished:
			return true, nil
		default:
			return false, nil
		}
	}

	return &LogReader{
		client: s.client,
		ctx:    ctx,
		done:   done,
		logURL: u,
	}, nil
}

// Cancel a test run by its ID.
func (s *testRuns) Cancel(ctx context.Context, moduleID RegistryModuleID, testRunID string) error {
	if err := moduleID.valid(); err != nil {
		return err
	}

	if !validStringID(&testRunID) {
		return ErrInvalidTestRunID
	}

	u := fmt.Sprintf("%s/%s/cancel", testRunsPath(moduleID), url.PathEscape(testRunID))
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// ForceCancel a test run by its ID.
func (s *testRuns) ForceCancel(ctx context.Context, moduleID RegistryModuleID, testRunID string) error {
	if err := moduleID.valid(); err != nil {
		return err
	}

	if !validStringID(&testRunID) {
		return ErrInvalidTestRunID
	}

	u := fmt.Sprintf("%s/%s/force-cancel", testRunsPath(moduleID), url.PathEscape(testRunID))
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o TestRunCreateOptions) valid() error {
	if o.ConfigurationVersion == nil {
		return ErrInvalidConfigVersionID
	}

	if o.RegistryModule == nil {
		return ErrRequiredRegistryModule
	}

	if o.RegistryModule.Organization == nil {
		return ErrRequiredOrg
	}

	return nil
}

func testRunsPath(moduleID RegistryModuleID) string {
	return fmt.Sprintf("organizations/%s/tests/registry-modules/%s/%s/%s/%s/test-runs",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(string(moduleID.RegistryName)),
		url.PathEscape(moduleID.Namespace),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider))
}
