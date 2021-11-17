package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

// AnalyticsIndexManager provides methods for performing Couchbase Analytics index management.
type AnalyticsIndexManager struct {
	aProvider    analyticsIndexQueryProvider
	mgmtProvider mgmtProvider

	globalTimeout time.Duration
	tracer        RequestTracer
	meter         *meterWrapper
}

type analyticsIndexQueryProvider interface {
	AnalyticsQuery(statement string, opts *AnalyticsOptions) (*AnalyticsResult, error)
}

func (am *AnalyticsIndexManager) doAnalyticsQuery(q string, opts *AnalyticsOptions) ([][]byte, error) {
	if opts.Timeout == 0 {
		opts.Timeout = am.globalTimeout
	}

	result, err := am.aProvider.AnalyticsQuery(q, opts)
	if err != nil {
		return nil, err
	}

	var rows [][]byte
	for result.Next() {
		var row json.RawMessage
		err := result.Row(&row)
		if err != nil {
			logWarnf("management operation failed to read row: %s", err)
		} else {
			rows = append(rows, row)
		}
	}
	err = result.Err()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (am *AnalyticsIndexManager) doMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error) {
	resp, err := am.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type jsonAnalyticsDataset struct {
	DatasetName   string `json:"DatasetName"`
	DataverseName string `json:"DataverseName"`
	LinkName      string `json:"LinkName"`
	BucketName    string `json:"BucketName"`
}

type jsonAnalyticsIndex struct {
	IndexName     string `json:"IndexName"`
	DatasetName   string `json:"DatasetName"`
	DataverseName string `json:"DataverseName"`
	IsPrimary     bool   `json:"IsPrimary"`
}

// AnalyticsDataset contains information about an analytics dataset.
type AnalyticsDataset struct {
	Name          string
	DataverseName string
	LinkName      string
	BucketName    string
}

func (ad *AnalyticsDataset) fromData(data jsonAnalyticsDataset) error {
	ad.Name = data.DatasetName
	ad.DataverseName = data.DataverseName
	ad.LinkName = data.LinkName
	ad.BucketName = data.BucketName

	return nil
}

// AnalyticsIndex contains information about an analytics index.
type AnalyticsIndex struct {
	Name          string
	DatasetName   string
	DataverseName string
	IsPrimary     bool
}

func (ai *AnalyticsIndex) fromData(data jsonAnalyticsIndex) error {
	ai.Name = data.IndexName
	ai.DatasetName = data.DatasetName
	ai.DataverseName = data.DataverseName
	ai.IsPrimary = data.IsPrimary

	return nil
}

// CreateAnalyticsDataverseOptions is the set of options available to the AnalyticsManager CreateDataverse operation.
type CreateAnalyticsDataverseOptions struct {
	IgnoreIfExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

func (am *AnalyticsIndexManager) uncompoundName(dataverse string) string {
	dvPieces := strings.Split(dataverse, "/")
	return "`" + strings.Join(dvPieces, "`.`") + "`"
}

// CreateDataverse creates a new analytics dataset.
func (am *AnalyticsIndexManager) CreateDataverse(dataverseName string, opts *CreateAnalyticsDataverseOptions) error {
	if opts == nil {
		opts = &CreateAnalyticsDataverseOptions{}
	}

	if dataverseName == "" {
		return invalidArgumentsError{
			message: "dataset name cannot be empty",
		}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_create_dataverse", start)

	var ignoreStr string
	if opts.IgnoreIfExists {
		ignoreStr = "IF NOT EXISTS"
	}

	q := fmt.Sprintf("CREATE DATAVERSE %s %s", am.uncompoundName(dataverseName), ignoreStr)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_create_dataverse", "management")
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:         opts.Timeout,
		RetryStrategy:   opts.RetryStrategy,
		ParentSpan:      span,
		ClientContextID: uuid.New().String(),
		Context:         opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// DropAnalyticsDataverseOptions is the set of options available to the AnalyticsManager DropDataverse operation.
type DropAnalyticsDataverseOptions struct {
	IgnoreIfNotExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropDataverse drops an analytics dataset.
func (am *AnalyticsIndexManager) DropDataverse(dataverseName string, opts *DropAnalyticsDataverseOptions) error {
	if opts == nil {
		opts = &DropAnalyticsDataverseOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_drop_dataverse", start)

	var ignoreStr string
	if opts.IgnoreIfNotExists {
		ignoreStr = "IF EXISTS"
	}

	q := fmt.Sprintf("DROP DATAVERSE %s %s", am.uncompoundName(dataverseName), ignoreStr)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_drop_dataverse", "management")
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return err
}

// CreateAnalyticsDatasetOptions is the set of options available to the AnalyticsManager CreateDataset operation.
type CreateAnalyticsDatasetOptions struct {
	IgnoreIfExists bool
	Condition      string
	DataverseName  string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateDataset creates a new analytics dataset.
func (am *AnalyticsIndexManager) CreateDataset(datasetName, bucketName string, opts *CreateAnalyticsDatasetOptions) error {
	if opts == nil {
		opts = &CreateAnalyticsDatasetOptions{}
	}

	if datasetName == "" {
		return invalidArgumentsError{
			message: "dataset name cannot be empty",
		}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_create_dataset", start)

	var ignoreStr string
	if opts.IgnoreIfExists {
		ignoreStr = "IF NOT EXISTS"
	}

	var where string
	if opts.Condition != "" {
		if !strings.HasPrefix(strings.ToUpper(opts.Condition), "WHERE") {
			where = "WHERE "
		}
		where += opts.Condition
	}

	if opts.DataverseName == "" {
		datasetName = fmt.Sprintf("`%s`", datasetName)
	} else {
		datasetName = fmt.Sprintf("%s.`%s`", am.uncompoundName(opts.DataverseName), datasetName)
	}

	q := fmt.Sprintf("CREATE DATASET %s %s ON `%s` %s", ignoreStr, datasetName, bucketName, where)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_create_dataset", "management")
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// DropAnalyticsDatasetOptions is the set of options available to the AnalyticsManager DropDataset operation.
type DropAnalyticsDatasetOptions struct {
	IgnoreIfNotExists bool
	DataverseName     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropDataset drops an analytics dataset.
func (am *AnalyticsIndexManager) DropDataset(datasetName string, opts *DropAnalyticsDatasetOptions) error {
	if opts == nil {
		opts = &DropAnalyticsDatasetOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_drop_dataset", start)

	var ignoreStr string
	if opts.IgnoreIfNotExists {
		ignoreStr = "IF EXISTS"
	}

	if opts.DataverseName == "" {
		datasetName = fmt.Sprintf("`%s`", datasetName)
	} else {
		datasetName = fmt.Sprintf("%s.`%s`", am.uncompoundName(opts.DataverseName), datasetName)
	}

	q := fmt.Sprintf("DROP DATASET %s %s", datasetName, ignoreStr)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_drop_dataset", "management")
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAllAnalyticsDatasetsOptions is the set of options available to the AnalyticsManager GetAllDatasets operation.
type GetAllAnalyticsDatasetsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllDatasets gets all analytics datasets.
func (am *AnalyticsIndexManager) GetAllDatasets(opts *GetAllAnalyticsDatasetsOptions) ([]AnalyticsDataset, error) {
	if opts == nil {
		opts = &GetAllAnalyticsDatasetsOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_get_all_datasets", start)

	q := "SELECT d.* FROM Metadata.`Dataset` d WHERE d.DataverseName <> \"Metadata\""
	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_get_all_datasets", "management")
	span.SetAttribute("db.statement", q)
	defer span.End()

	rows, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return nil, err
	}

	datasets := make([]AnalyticsDataset, len(rows))
	for rowIdx, row := range rows {
		var datasetData jsonAnalyticsDataset
		err := json.Unmarshal(row, &datasetData)
		if err != nil {
			return nil, err
		}

		err = datasets[rowIdx].fromData(datasetData)
		if err != nil {
			return nil, err
		}
	}

	return datasets, nil
}

// CreateAnalyticsIndexOptions is the set of options available to the AnalyticsManager CreateIndex operation.
type CreateAnalyticsIndexOptions struct {
	IgnoreIfExists bool
	DataverseName  string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateIndex creates a new analytics dataset.
func (am *AnalyticsIndexManager) CreateIndex(datasetName, indexName string, fields map[string]string, opts *CreateAnalyticsIndexOptions) error {
	if opts == nil {
		opts = &CreateAnalyticsIndexOptions{}
	}

	if indexName == "" {
		return invalidArgumentsError{
			message: "index name cannot be empty",
		}
	}
	if len(fields) <= 0 {
		return invalidArgumentsError{
			message: "you must specify at least one field to index",
		}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_create_index", start)

	var ignoreStr string
	if opts.IgnoreIfExists {
		ignoreStr = "IF NOT EXISTS"
	}

	var indexFields []string
	for name, typ := range fields {
		indexFields = append(indexFields, name+":"+typ)
	}

	if opts.DataverseName == "" {
		datasetName = fmt.Sprintf("`%s`", datasetName)
	} else {
		datasetName = fmt.Sprintf("%s.`%s`", am.uncompoundName(opts.DataverseName), datasetName)
	}

	q := fmt.Sprintf("CREATE INDEX `%s` %s ON %s (%s)", indexName, ignoreStr, datasetName, strings.Join(indexFields, ","))

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_create_index", "management")
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// DropAnalyticsIndexOptions is the set of options available to the AnalyticsManager DropIndex operation.
type DropAnalyticsIndexOptions struct {
	IgnoreIfNotExists bool
	DataverseName     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropIndex drops an analytics index.
func (am *AnalyticsIndexManager) DropIndex(datasetName, indexName string, opts *DropAnalyticsIndexOptions) error {
	if opts == nil {
		opts = &DropAnalyticsIndexOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_drop_index", start)

	var ignoreStr string
	if opts.IgnoreIfNotExists {
		ignoreStr = "IF EXISTS"
	}

	if opts.DataverseName == "" {
		datasetName = fmt.Sprintf("`%s`", datasetName)
	} else {
		datasetName = fmt.Sprintf("%s.`%s`", am.uncompoundName(opts.DataverseName), datasetName)
	}

	q := fmt.Sprintf("DROP INDEX %s.%s %s", datasetName, indexName, ignoreStr)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_drop_index", "management")
	span.SetAttribute("db.statement", q)
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAllAnalyticsIndexesOptions is the set of options available to the AnalyticsManager GetAllIndexes operation.
type GetAllAnalyticsIndexesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllIndexes gets all analytics indexes.
func (am *AnalyticsIndexManager) GetAllIndexes(opts *GetAllAnalyticsIndexesOptions) ([]AnalyticsIndex, error) {
	if opts == nil {
		opts = &GetAllAnalyticsIndexesOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_get_all_indexes", start)

	q := "SELECT d.* FROM Metadata.`Index` d WHERE d.DataverseName <> \"Metadata\""
	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_get_all_indexes", "management")
	defer span.End()

	rows, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return nil, err
	}

	indexes := make([]AnalyticsIndex, len(rows))
	for rowIdx, row := range rows {
		var indexData jsonAnalyticsIndex
		err := json.Unmarshal(row, &indexData)
		if err != nil {
			return nil, err
		}

		err = indexes[rowIdx].fromData(indexData)
		if err != nil {
			return nil, err
		}
	}

	return indexes, nil
}

// ConnectAnalyticsLinkOptions is the set of options available to the AnalyticsManager ConnectLink operation.
type ConnectAnalyticsLinkOptions struct {
	LinkName      string
	DataverseName string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ConnectLink connects an analytics link.
func (am *AnalyticsIndexManager) ConnectLink(opts *ConnectAnalyticsLinkOptions) error {
	if opts == nil {
		opts = &ConnectAnalyticsLinkOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_connect_link", start)

	linkName := opts.LinkName
	if linkName == "" {
		linkName = "Local"
	}
	if opts.DataverseName != "" {
		linkName = fmt.Sprintf("%s.`%s`", am.uncompoundName(opts.DataverseName), linkName)
	}

	q := fmt.Sprintf("CONNECT LINK %s", linkName)
	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_connect_link", "management")
	span.SetAttribute("db.statement", q)
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// DisconnectAnalyticsLinkOptions is the set of options available to the AnalyticsManager DisconnectLink operation.
type DisconnectAnalyticsLinkOptions struct {
	LinkName      string
	DataverseName string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DisconnectLink disconnects an analytics link.
func (am *AnalyticsIndexManager) DisconnectLink(opts *DisconnectAnalyticsLinkOptions) error {
	if opts == nil {
		opts = &DisconnectAnalyticsLinkOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_disconnect_link", start)

	linkName := opts.LinkName
	if linkName == "" {
		linkName = "Local"
	}
	if opts.DataverseName != "" {
		linkName = fmt.Sprintf("%s.`%s`", am.uncompoundName(opts.DataverseName), linkName)
	}

	q := fmt.Sprintf("DISCONNECT LINK %s", linkName)
	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_disconnect_link", "management")
	defer span.End()

	_, err := am.doAnalyticsQuery(q, &AnalyticsOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetPendingMutationsAnalyticsOptions is the set of options available to the user manager GetPendingMutations operation.
type GetPendingMutationsAnalyticsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetPendingMutations returns the number of pending mutations for all indexes in the form of dataverse.dataset:mutations.
func (am *AnalyticsIndexManager) GetPendingMutations(opts *GetPendingMutationsAnalyticsOptions) (map[string]map[string]int, error) {
	if opts == nil {
		opts = &GetPendingMutationsAnalyticsOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_get_pending_mutations", start)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_get_pending_mutations", "management")
	span.SetAttribute("db.operation", "GET /analytics/node/agg/stats/remaining")
	defer span.End()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = am.globalTimeout
	}

	req := mgmtRequest{
		Service:       ServiceTypeAnalytics,
		Method:        "GET",
		Path:          "/analytics/node/agg/stats/remaining",
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       timeout,
		parentSpanCtx: span.Context(),
	}
	resp, err := am.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, makeMgmtBadStatusError("failed to get pending mutations", &req, resp)
	}

	pending := make(map[string]map[string]int)
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&pending)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		logDebugf("Failed to close socket (%s)", err)
	}

	return pending, nil
}

// AnalyticsLink describes an external or remote analytics link, used to access data external to the cluster.
type AnalyticsLink interface {
	// Name returns the name of this link.
	Name() string
	// DataverseName returns the name of the dataverse that this link belongs to.
	DataverseName() string
	// FormEncode encodes the link into a form data representation, to be sent as the body of a CreateLink or ReplaceLink
	// request.
	FormEncode() ([]byte, error)
	// Validate is used by CreateLink and ReplaceLink to ensure that the link is valid.
	Validate() error
	// LinkType returns the type of analytics type this link is.
	LinkType() AnalyticsLinkType
}

// NewCouchbaseRemoteAnalyticsLinkOptions are the options available when creating a new CouchbaseRemoteAnalyticsLink.
type NewCouchbaseRemoteAnalyticsLinkOptions struct {
	Encryption CouchbaseRemoteAnalyticsEncryptionSettings
	Username   string
	Password   string
}

// NewCouchbaseRemoteAnalyticsLink creates a new CouchbaseRemoteAnalyticsLink.
// Scope is the analytics scope in the form of "bucket/scope".
func NewCouchbaseRemoteAnalyticsLink(linkName, hostname, dataverseName string,
	opts *NewCouchbaseRemoteAnalyticsLinkOptions) *CouchbaseRemoteAnalyticsLink {
	if opts == nil {
		opts = &NewCouchbaseRemoteAnalyticsLinkOptions{}
	}
	return &CouchbaseRemoteAnalyticsLink{
		Dataverse:  dataverseName,
		LinkName:   linkName,
		Hostname:   hostname,
		Encryption: opts.Encryption,
		Username:   opts.Username,
		Password:   opts.Password,
	}
}

// NewS3ExternalAnalyticsLinkOptions are the options available when creating a new S3ExternalAnalyticsLink.
type NewS3ExternalAnalyticsLinkOptions struct {
	SessionToken    string
	ServiceEndpoint string
}

// NewS3ExternalAnalyticsLink creates a new S3ExternalAnalyticsLink with the scope field populated.
// Scope is the analytics scope in the form of "bucket/scope".
func NewS3ExternalAnalyticsLink(linkName, dataverseName, accessKeyID, secretAccessKey, region string,
	opts *NewS3ExternalAnalyticsLinkOptions) *S3ExternalAnalyticsLink {
	if opts == nil {
		opts = &NewS3ExternalAnalyticsLinkOptions{}
	}
	return &S3ExternalAnalyticsLink{
		Dataverse:       dataverseName,
		LinkName:        linkName,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    opts.SessionToken,
		Region:          region,
		ServiceEndpoint: opts.ServiceEndpoint,
	}
}

// NewAzureBlobExternalAnalyticsLinkOptions are the options available when creating a new AzureBlobExternalAnalyticsLink.
// VOLATILE: This API is subject to change at any time.
type NewAzureBlobExternalAnalyticsLinkOptions struct {
	ConnectionString      string
	AccountName           string
	AccountKey            string
	SharedAccessSignature string
	BlobEndpoint          string
	EndpointSuffix        string
}

// NewAzureBlobExternalAnalyticsLink creates a new AzureBlobExternalAnalyticsLink.
// VOLATILE: This API is subject to change at any time.
func NewAzureBlobExternalAnalyticsLink(linkName, dataverseName string,
	opts *NewAzureBlobExternalAnalyticsLinkOptions) *AzureBlobExternalAnalyticsLink {
	if opts == nil {
		opts = &NewAzureBlobExternalAnalyticsLinkOptions{}
	}
	return &AzureBlobExternalAnalyticsLink{
		Dataverse:             dataverseName,
		LinkName:              linkName,
		ConnectionString:      opts.ConnectionString,
		AccountName:           opts.AccountName,
		AccountKey:            opts.AccountKey,
		SharedAccessSignature: opts.SharedAccessSignature,
		BlobEndpoint:          opts.BlobEndpoint,
		EndpointSuffix:        opts.EndpointSuffix,
	}
}

// CreateAnalyticsLinkOptions is the set of options available to the analytics manager CreateLink
// function.
type CreateAnalyticsLinkOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// CreateLink creates an analytics link.
func (am *AnalyticsIndexManager) CreateLink(link AnalyticsLink, opts *CreateAnalyticsLinkOptions) error {
	if opts == nil {
		opts = &CreateAnalyticsLinkOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_create_link", start)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_create_link", "management")
	defer span.End()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = am.globalTimeout
	}

	if err := link.Validate(); err != nil {
		return err
	}

	endpoint := am.endpointFromLink(link)
	span.SetAttribute("db.operation", "POST "+endpoint)

	eSpan := createSpan(am.tracer, span, "request_encoding", "")
	data, err := link.FormEncode()
	eSpan.End()
	if err != nil {
		return err
	}

	req := mgmtRequest{
		Service:       ServiceTypeAnalytics,
		Method:        "POST",
		Path:          endpoint,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       timeout,
		parentSpanCtx: span.Context(),
		Body:          data,
		ContentType:   "application/x-www-form-urlencoded",
	}

	resp, err := am.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return am.tryParseLinkErrorMessage(&req, resp)
	}

	return nil
}

// ReplaceAnalyticsLinkOptions is the set of options available to the analytics manager ReplaceLink
// function.
type ReplaceAnalyticsLinkOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ReplaceLink modifies an existing analytics link.
func (am *AnalyticsIndexManager) ReplaceLink(link AnalyticsLink, opts *ReplaceAnalyticsLinkOptions) error {
	if opts == nil {
		opts = &ReplaceAnalyticsLinkOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_replace_link", start)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_replace_link", "management")
	defer span.End()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = am.globalTimeout
	}

	if err := link.Validate(); err != nil {
		return err
	}

	endpoint := am.endpointFromLink(link)
	span.SetAttribute("db.operation", "PUT "+endpoint)

	eSpan := createSpan(am.tracer, span, "request_encoding", "")
	data, err := link.FormEncode()
	eSpan.End()
	if err != nil {
		return err
	}

	req := mgmtRequest{
		Service:       ServiceTypeAnalytics,
		Method:        "PUT",
		Path:          endpoint,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       timeout,
		parentSpanCtx: span.Context(),
		Body:          data,
		ContentType:   "application/x-www-form-urlencoded",
	}

	resp, err := am.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return am.tryParseLinkErrorMessage(&req, resp)
	}

	return nil
}

// DropAnalyticsLinkOptions is the set of options available to the analytics manager DropLink
// function.
type DropAnalyticsLinkOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropLink removes an existing external analytics link from specified scope.
// dataverseName can be given in the form of "namepart" or "namepart1/namepart2".
// Only available against Couchbase Server 7.0+.
func (am *AnalyticsIndexManager) DropLink(linkName, dataverseName string, opts *DropAnalyticsLinkOptions) error {
	if opts == nil {
		opts = &DropAnalyticsLinkOptions{}
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_drop_link", start)

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_drop_link", "management")
	defer span.End()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = am.globalTimeout
	}

	var payload []byte
	var endpoint string
	if strings.Contains(dataverseName, "/") {
		endpoint = fmt.Sprintf("/analytics/link/%s/%s", url.PathEscape(dataverseName), linkName)
	} else {
		endpoint = "/analytics/link"
		values := url.Values{}
		values.Add("dataverse", dataverseName)
		values.Add("name", linkName)

		eSpan := createSpan(am.tracer, span, spanNameRequestEncoding, "management")
		payload = []byte(values.Encode())
		eSpan.End()
	}
	span.SetAttribute("db.operation", "DELETE "+endpoint)

	req := mgmtRequest{
		Service:       ServiceTypeAnalytics,
		Method:        "DELETE",
		Path:          endpoint,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       timeout,
		parentSpanCtx: span.Context(),
		ContentType:   "application/x-www-form-urlencoded",
		Body:          payload,
	}

	resp, err := am.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return am.tryParseLinkErrorMessage(&req, resp)
	}

	return nil
}

// GetAnalyticsLinksOptions are the options available to the AnalyticsManager GetLinks function.
type GetAnalyticsLinksOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Dataverse restricts the results to a given dataverse, can be given in the form of "namepart" or "namepart1/namepart2".
	Dataverse string
	// LinkType restricts the results to the given link type.
	LinkType AnalyticsLinkType
	// Name restricts the results to the link with the specified name.
	// If set then `Scope` must also be set.
	Name string

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetLinks retrieves all external or remote analytics links.
func (am *AnalyticsIndexManager) GetLinks(opts *GetAnalyticsLinksOptions) ([]AnalyticsLink, error) {
	if opts == nil {
		opts = &GetAnalyticsLinksOptions{}
	}

	if opts.Name != "" && opts.Dataverse == "" {
		return nil, makeInvalidArgumentsError("when name is set then dataverse must also be set")
	}

	start := time.Now()
	defer am.meter.ValueRecord(meterValueServiceManagement, "manager_analytics_get_all_links", start)

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = am.globalTimeout
	}

	var querystring []string
	var endpoint string
	if strings.Contains(opts.Dataverse, "/") {
		endpoint = fmt.Sprintf("/analytics/link/%s", url.PathEscape(opts.Dataverse))

		if opts.Name != "" {
			endpoint = fmt.Sprintf("%s/%s", endpoint, opts.Name)
		}
		if opts.LinkType != "" {
			querystring = append(querystring, fmt.Sprintf("type=%s", opts.LinkType))
		}
	} else {
		endpoint = "/analytics/link"

		if opts.Dataverse != "" {
			querystring = append(querystring, "dataverse="+opts.Dataverse)
			if opts.Name != "" {
				querystring = append(querystring, "name="+opts.Name)
			}
		}
		if opts.LinkType != "" {
			querystring = append(querystring, fmt.Sprintf("type=%s", opts.LinkType))
		}
	}

	if len(querystring) > 0 {
		endpoint = endpoint + "?" + strings.Join(querystring, "&")
	}

	span := createSpan(am.tracer, opts.ParentSpan, "manager_analytics_get_all_links", "management")
	span.SetAttribute("db.operation", "GET "+endpoint)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeAnalytics,
		Method:        "GET",
		Path:          endpoint,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       timeout,
		parentSpanCtx: span.Context(),
		IsIdempotent:  true,
	}

	resp, err := am.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, am.tryParseLinkErrorMessage(&req, resp)
	}

	var jsonLinks []map[string]interface{}
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&jsonLinks)
	if err != nil {
		return nil, err
	}

	var links []AnalyticsLink
	for _, jsonLink := range jsonLinks {
		linkType, ok := jsonLink["type"]
		if !ok {
			logWarnf("External analytics link missing type field, skipping")
			continue
		}

		linkTypeStr, ok := linkType.(string)
		if !ok {
			logWarnf("External analytics link type field not a string, skipping")
			continue
		}

		link := am.linkFromJSON(AnalyticsLinkType(linkTypeStr), jsonLink)
		if link == nil {
			logWarnf("External analytics link type %s unknown, skipping", linkTypeStr)
			continue
		}

		links = append(links, link)
	}

	return links, nil
}

func (am *AnalyticsIndexManager) fieldFromJSONMapAsString(name string, json map[string]interface{}) string {
	field, ok := json[name]
	if !ok {
		return ""
	}

	strField, ok := field.(string)
	if !ok {
		return ""
	}

	return strField
}

func (am *AnalyticsIndexManager) endpointFromLink(link AnalyticsLink) string {
	var endpoint string
	switch l := link.(type) {
	case *CouchbaseRemoteAnalyticsLink:
		if strings.Contains(l.Dataverse, "/") {
			endpoint = fmt.Sprintf("/analytics/link/%s/%s", url.PathEscape(l.Dataverse), l.LinkName)
		} else {
			endpoint = "/analytics/link"
		}
	case *S3ExternalAnalyticsLink:
		if strings.Contains(l.Dataverse, "/") {
			endpoint = fmt.Sprintf("/analytics/link/%s/%s", url.PathEscape(l.Dataverse), l.LinkName)
		} else {
			endpoint = "/analytics/link"
		}
	case *AzureBlobExternalAnalyticsLink:
		endpoint = fmt.Sprintf("/analytics/link/%s/%s", url.PathEscape(l.Dataverse), l.LinkName)
	default:
		endpoint = "/analytics/link"
	}
	return endpoint
}

func (am *AnalyticsIndexManager) tryParseLinkErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read bucket manager response body: %s", err)
		return nil
	}

	if strings.Contains(strings.ToLower(string(b)), "24055") {
		return makeGenericMgmtError(ErrAnalyticsLinkExists, req, resp)
	}
	if strings.Contains(strings.ToLower(string(b)), "24034") {
		return makeGenericMgmtError(ErrDataverseNotFound, req, resp)
	}

	return makeGenericMgmtError(errors.New(string(b)), req, resp)
}

func (am *AnalyticsIndexManager) linkFromJSON(linkType AnalyticsLinkType, jsonLink map[string]interface{}) AnalyticsLink {
	dataverse := am.fieldFromJSONMapAsString("dataverse", jsonLink)
	if dataverse == "" {
		dataverse = am.fieldFromJSONMapAsString("scope", jsonLink)
	}
	switch linkType {
	case AnalyticsLinkTypeCouchbaseRemote:
		encryptionLevel := am.fieldFromJSONMapAsString("encryption", jsonLink)
		return &CouchbaseRemoteAnalyticsLink{
			Dataverse: dataverse,
			LinkName:  am.fieldFromJSONMapAsString("name", jsonLink),
			Hostname:  am.fieldFromJSONMapAsString("activeHostname", jsonLink),
			Encryption: CouchbaseRemoteAnalyticsEncryptionSettings{
				EncryptionLevel:   analyticsEncryptionLevelFromString(encryptionLevel),
				Certificate:       []byte(am.fieldFromJSONMapAsString("certificate", jsonLink)),
				ClientCertificate: []byte(am.fieldFromJSONMapAsString("clientCertificate", jsonLink)),
			},
			Username: am.fieldFromJSONMapAsString("username", jsonLink),
		}
	case AnalyticsLinkTypeS3External:
		return &S3ExternalAnalyticsLink{
			Dataverse:       dataverse,
			LinkName:        am.fieldFromJSONMapAsString("name", jsonLink),
			AccessKeyID:     am.fieldFromJSONMapAsString("accessKeyId", jsonLink),
			Region:          am.fieldFromJSONMapAsString("region", jsonLink),
			ServiceEndpoint: am.fieldFromJSONMapAsString("serviceEndpoint", jsonLink),
		}
	case AnalyticsLinkTypeAzureExternal:
		return &AzureBlobExternalAnalyticsLink{
			Dataverse:      dataverse,
			LinkName:       am.fieldFromJSONMapAsString("name", jsonLink),
			AccountName:    am.fieldFromJSONMapAsString("accountName", jsonLink),
			BlobEndpoint:   am.fieldFromJSONMapAsString("blobEndpoint", jsonLink),
			EndpointSuffix: am.fieldFromJSONMapAsString("endpointSuffix", jsonLink),
		}
	default:
		return nil
	}
}
