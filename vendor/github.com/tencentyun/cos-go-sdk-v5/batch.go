package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

type BatchService service

type BatchRequestHeaders struct {
	XCosAppid     int          `header:"x-cos-appid" xml:"-" url:"-"`
	ContentLength string       `header:"Content-Length,omitempty" xml:"-" url:"-"`
	ContentType   string       `header:"Content-Type,omitempty" xml:"-" url:"-"`
	Headers       *http.Header `header:"-" xml:"-", url:"-"`
}

// BatchProgressSummary
type BatchProgressSummary struct {
	NumberOfTasksFailed    int `xml:"NumberOfTasksFailed" header:"-" url:"-"`
	NumberOfTasksSucceeded int `xml:"NumberOfTasksSucceeded" header:"-" url:"-"`
	TotalNumberOfTasks     int `xml:"TotalNumberOfTasks" header:"-" url:"-"`
}

// BatchJobReport
type BatchJobReport struct {
	Bucket      string `xml:"Bucket" header:"-" url:"-"`
	Enabled     string `xml:"Enabled" header:"-" url:"-"`
	Format      string `xml:"Format" header:"-" url:"-"`
	Prefix      string `xml:"Prefix,omitempty" header:"-" url:"-"`
	ReportScope string `xml:"ReportScope" header:"-" url:"-"`
}

// BatchJobOperationCopy
type BatchMetadata struct {
	Key   string `xml:"Key" header:"-" url:"-"`
	Value string `xml:"Value" header:"-" url:"-"`
}
type BatchNewObjectMetadata struct {
	CacheControl       string          `xml:"CacheControl,omitempty" header:"-" url:"-"`
	ContentDisposition string          `xml:"ContentDisposition,omitempty" header:"-" url:"-"`
	ContentEncoding    string          `xml:"ContentEncoding,omitempty" header:"-" url:"-"`
	ContentType        string          `xml:"ContentType,omitempty" header:"-" url:"-"`
	HttpExpiresDate    string          `xml:"HttpExpiresDate,omitempty" header:"-" url:"-"`
	SSEAlgorithm       string          `xml:"SSEAlgorithm,omitempty" header:"-" url:"-"`
	UserMetadata       []BatchMetadata `xml:"UserMetadata>member,omitempty" header:"-" url:"-"`
}
type BatchGrantee struct {
	DisplayName    string `xml:"DisplayName,omitempty" header:"-" url:"-"`
	Identifier     string `xml:"Identifier" header:"-" url:"-"`
	TypeIdentifier string `xml:"TypeIdentifier" header:"-" url:"-"`
}
type BatchCOSGrant struct {
	Grantee    *BatchGrantee `xml:"Grantee" header:"-" url:"-"`
	Permission string        `xml:"Permission" header:"-" url:"-"`
}
type BatchAccessControlGrants struct {
	COSGrants *BatchCOSGrant `xml:"COSGrant,omitempty" header:"-" url:"-"`
}
type BatchJobOperationCopy struct {
	AccessControlGrants       *BatchAccessControlGrants `xml:"AccessControlGrants,omitempty" header:"-" url:"-"`
	CannedAccessControlList   string                    `xml:"CannedAccessControlList,omitempty" header:"-" url:"-"`
	MetadataDirective         string                    `xml:"MetadataDirective,omitempty" header:"-" url:"-"`
	ModifiedSinceConstraint   int64                     `xml:"ModifiedSinceConstraint,omitempty" header:"-" url:"-"`
	UnModifiedSinceConstraint int64                     `xml:"UnModifiedSinceConstraint,omitempty" header:"-" url:"-"`
	NewObjectMetadata         *BatchNewObjectMetadata   `xml:"NewObjectMetadata,omitempty" header:"-" url:"-"`
	StorageClass              string                    `xml:"StorageClass,omitempty" header:"-" url:"-"`
	TargetResource            string                    `xml:"TargetResource" header:"-" url:"-"`
}

// BatchInitiateRestoreObject
type BatchInitiateRestoreObject struct {
	ExpirationInDays int    `xml:"ExpirationInDays"`
	JobTier          string `xml:"JobTier"`
}

// BatchJobOperation
type BatchJobOperation struct {
	PutObjectCopy *BatchJobOperationCopy      `xml:"COSPutObjectCopy,omitempty" header:"-" url:"-"`
	RestoreObject *BatchInitiateRestoreObject `xml:"COSInitiateRestoreObject,omitempty" header:"-" url:"-"`
}

// BatchJobManifest
type BatchJobManifestLocation struct {
	ETag            string `xml:"ETag" header:"-" url:"-"`
	ObjectArn       string `xml:"ObjectArn" header:"-" url:"-"`
	ObjectVersionId string `xml:"ObjectVersionId,omitempty" header:"-" url:"-"`
}
type BatchJobManifestSpec struct {
	Fields []string `xml:"Fields>member,omitempty" header:"-" url:"-"`
	Format string   `xml:"Format" header:"-" url:"-"`
}
type BatchJobManifest struct {
	Location *BatchJobManifestLocation `xml:"Location" header:"-" url:"-"`
	Spec     *BatchJobManifestSpec     `xml:"Spec" header:"-" url:"-"`
}

type BatchCreateJobOptions struct {
	XMLName              xml.Name           `xml:"CreateJobRequest" header:"-" url:"-"`
	ClientRequestToken   string             `xml:"ClientRequestToken" header:"-" url:"-"`
	ConfirmationRequired string             `xml:"ConfirmationRequired,omitempty" header:"-" url:"-"`
	Description          string             `xml:"Description,omitempty" header:"-" url:"-"`
	Manifest             *BatchJobManifest  `xml:"Manifest" header:"-" url:"-"`
	Operation            *BatchJobOperation `xml:"Operation" header:"-" url:"-"`
	Priority             int                `xml:"Priority" header:"-" url:"-"`
	Report               *BatchJobReport    `xml:"Report" header:"-" url:"-"`
	RoleArn              string             `xml:"RoleArn" header:"-" url:"-"`
}

type BatchCreateJobResult struct {
	XMLName xml.Name `xml:"CreateJobResult"`
	JobId   string   `xml:"JobId,omitempty"`
}

func processETag(opt *BatchCreateJobOptions) *BatchCreateJobOptions {
	if opt != nil && opt.Manifest != nil && opt.Manifest.Location != nil {
		opt.Manifest.Location.ETag = "<ETag>" + opt.Manifest.Location.ETag + "</ETag>"
	}
	return opt
}

func (s *BatchService) CreateJob(ctx context.Context, opt *BatchCreateJobOptions, headers *BatchRequestHeaders) (*BatchCreateJobResult, *Response, error) {
	var res BatchCreateJobResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       "/jobs",
		method:    http.MethodPost,
		optHeader: headers,
		body:      opt,
		result:    &res,
	}

	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchJobFailureReasons struct {
	FailureCode   string `xml:"FailureCode" header:"-" url:"-"`
	FailureReason string `xml:"FailureReason" header:"-" url:"-"`
}

type BatchDescribeJob struct {
	ConfirmationRequired string                  `xml:"ConfirmationRequired,omitempty" header:"-" url:"-"`
	CreationTime         string                  `xml:"CreationTime,omitempty" header:"-" url:"-"`
	Description          string                  `xml:"Description,omitempty" header:"-" url:"-"`
	FailureReasons       *BatchJobFailureReasons `xml:"FailureReasons>JobFailure,omitempty" header:"-" url:"-"`
	JobId                string                  `xml:"JobId" header:"-" url:"-"`
	Manifest             *BatchJobManifest       `xml:"Manifest" header:"-" url:"-"`
	Operation            *BatchJobOperation      `xml:"Operation" header:"-" url:"-"`
	Priority             int                     `xml:"Priority" header:"-" url:"-"`
	ProgressSummary      *BatchProgressSummary   `xml:"ProgressSummary" header:"-" url:"-"`
	Report               *BatchJobReport         `xml:"Report,omitempty" header:"-" url:"-"`
	RoleArn              string                  `xml:"RoleArn,omitempty" header:"-" url:"-"`
	Status               string                  `xml:"Status,omitempty" header:"-" url:"-"`
	StatusUpdateReason   string                  `xml:"StatusUpdateReason,omitempty" header:"-" url:"-"`
	SuspendedCause       string                  `xml:"SuspendedCause,omitempty" header:"-" url:"-"`
	SuspendedDate        string                  `xml:"SuspendedDate,omitempty" header:"-" url:"-"`
	TerminationDate      string                  `xml:"TerminationDate,omitempty" header:"-" url:"-"`
}
type BatchDescribeJobResult struct {
	XMLName xml.Name          `xml:"DescribeJobResult"`
	Job     *BatchDescribeJob `xml:"Job,omitempty"`
}

func (s *BatchService) DescribeJob(ctx context.Context, id string, headers *BatchRequestHeaders) (*BatchDescribeJobResult, *Response, error) {
	var res BatchDescribeJobResult
	u := fmt.Sprintf("/jobs/%s", id)
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       u,
		method:    http.MethodGet,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchListJobsOptions struct {
	JobStatuses string `url:"jobStatuses,omitempty" header:"-" xml:"-"`
	MaxResults  int    `url:"maxResults,omitempty" header:"-" xml:"-"`
	NextToken   string `url:"nextToken,omitempty" header:"-" xml:"-"`
}

type BatchListJobsMember struct {
	CreationTime    string                `xml:"CreationTime,omitempty" header:"-" url:"-"`
	Description     string                `xml:"Description,omitempty" header:"-" url:"-"`
	JobId           string                `xml:"JobId,omitempty" header:"-" url:"-"`
	Operation       string                `xml:"Operation,omitempty" header:"-" url:"-"`
	Priority        int                   `xml:"Priority,omitempty" header:"-" url:"-"`
	ProgressSummary *BatchProgressSummary `xml:"ProgressSummary,omitempty" header:"-" url:"-"`
	Status          string                `xml:"Status,omitempty" header:"-" url:"-"`
	TerminationDate string                `xml:"TerminationDate,omitempty" header:"-" url:"-"`
}
type BatchListJobs struct {
	Members []BatchListJobsMember `xml:"member,omitempty" header:"-" url:"-"`
}
type BatchListJobsResult struct {
	XMLName   xml.Name       `xml:"ListJobsResult"`
	Jobs      *BatchListJobs `xml:"Jobs,omitempty"`
	NextToken string         `xml:"NextToken,omitempty"`
}

func (s *BatchService) ListJobs(ctx context.Context, opt *BatchListJobsOptions, headers *BatchRequestHeaders) (*BatchListJobsResult, *Response, error) {
	var res BatchListJobsResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       "/jobs",
		method:    http.MethodGet,
		optQuery:  opt,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchUpdatePriorityOptions struct {
	JobId    string `url:"-" header:"-" xml:"-"`
	Priority int    `url:"priority" header:"-" xml:"-"`
}
type BatchUpdatePriorityResult struct {
	XMLName  xml.Name `xml:"UpdateJobPriorityResult"`
	JobId    string   `xml:"JobId,omitempty"`
	Priority int      `xml:"Priority,omitempty"`
}

func (s *BatchService) UpdateJobPriority(ctx context.Context, opt *BatchUpdatePriorityOptions, headers *BatchRequestHeaders) (*BatchUpdatePriorityResult, *Response, error) {
	u := fmt.Sprintf("/jobs/%s/priority", opt.JobId)
	var res BatchUpdatePriorityResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       u,
		method:    http.MethodPost,
		optQuery:  opt,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchUpdateStatusOptions struct {
	JobId              string `header:"-" url:"-" xml:"-"`
	RequestedJobStatus string `url:"requestedJobStatus" header:"-" xml:"-"`
	StatusUpdateReason string `url:"statusUpdateReason,omitempty" header:"-", xml:"-"`
}
type BatchUpdateStatusResult struct {
	XMLName            xml.Name `xml:"UpdateJobStatusResult"`
	JobId              string   `xml:"JobId,omitempty"`
	Status             string   `xml:"Status,omitempty"`
	StatusUpdateReason string   `xml:"StatusUpdateReason,omitempty"`
}

func (s *BatchService) UpdateJobStatus(ctx context.Context, opt *BatchUpdateStatusOptions, headers *BatchRequestHeaders) (*BatchUpdateStatusResult, *Response, error) {
	u := fmt.Sprintf("/jobs/%s/status", opt.JobId)
	var res BatchUpdateStatusResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       u,
		method:    http.MethodPost,
		optQuery:  opt,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}
