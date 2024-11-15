package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

type eventingManagementProviderCore struct {
	mgmtProvider mgmtProvider

	tracer RequestTracer
	meter  *meterWrapper
}

type eventingRequestOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan
	Context       context.Context
}

type jsonEventingFunction struct {
	Name               string                               `json:"appname"`
	Code               string                               `json:"appcode"`
	Version            string                               `json:"version"`
	EnforceSchema      bool                                 `json:"enforce_schema,omitempty"`
	HandlerUUID        int                                  `json:"handleruuid,omitempty"`
	FunctionInstanceID string                               `json:"function_instance_id,omitempty"`
	Settings           jsonEventingFunctionSettings         `json:"settings"`
	DeploymentConfig   jsonEventingFunctionDeploymentConfig `json:"depcfg"`
	FunctionScope      *jsonEventingFunctionScope           `json:"function_scope,omitempty"`
}

type jsonEventingFunctionScope struct {
	BucketName string `json:"bucket"`
	ScopeName  string `json:"scope"`
}

type jsonEventingFunctionSettings struct {
	CPPWorkerThreadCount   int                                   `json:"cpp_worker_thread_count,omitempty"`
	DCPStreamBoundary      EventingFunctionDCPBoundary           `json:"dcp_stream_boundary,omitempty"`
	Description            string                                `json:"description,omitempty"`
	DeploymentStatus       EventingFunctionDeploymentStatus      `json:"deployment_status"`
	ProcessingStatus       EventingFunctionProcessingStatus      `json:"processing_status"`
	LanguageCompatibility  EventingFunctionLanguageCompatibility `json:"language_compatibility,omitempty"`
	LogLevel               EventingFunctionLogLevel              `json:"log_level,omitempty"`
	ExecutionTimeout       int                                   `json:"execution_timeout,omitempty"`
	LCBInstCapacity        int                                   `json:"lcb_inst_capacity,omitempty"`
	LCBRetryCount          int                                   `json:"lcb_retry_count,omitempty"`
	LCBTimeout             int                                   `json:"lcb_timeout,omitempty"`
	QueryConsistency       QueryScanConsistency                  `json:"n1ql_consistency,omitempty"`
	NumTimerPartitions     int                                   `json:"num_timer_partitions,omitempty"`
	SockBatchSize          int                                   `json:"sock_batch_size,omitempty"`
	TickDuration           int                                   `json:"tick_duration,omitempty"`
	TimerContextSize       int                                   `json:"timer_context_size,omitempty"`
	UserPrefix             string                                `json:"user_prefix,omitempty"`
	BucketCacheSize        int                                   `json:"bucket_cache_size,omitempty"`
	BucketCacheAge         int                                   `json:"bucket_cache_age,omitempty"`
	CurlMaxAllowedRespSize int                                   `json:"curl_max_allowed_resp_size,omitempty"`
	QueryPrepareAll        bool                                  `json:"n1ql_prepare_all,omitempty"`
	WorkerCount            int                                   `json:"worker_count,omitempty"`
	HandlerHeaders         []string                              `json:"handler_headers,omitempty"`
	HandlerFooters         []string                              `json:"handler_footers,omitempty"`
	EnableAppLogRotation   bool                                  `json:"enable_applog_rotation,omitempty"`
	AppLogDir              string                                `json:"app_log_dir,omitempty"`
	AppLogMaxSize          int                                   `json:"app_log_max_size,omitempty"`
	AppLogMaxFiles         int                                   `json:"app_log_max_files,omitempty"`
	CheckpointInterval     int                                   `json:"checkpoint_interval,omitempty"`
}

type jsonEventingFunctionDeploymentConfig struct {
	MetadataBucket     string                                `json:"metadata_bucket"`
	MetadataScope      string                                `json:"metadata_scope,omitempty"`
	MetadataCollection string                                `json:"metadata_collection,omitempty"`
	SourceBucket       string                                `json:"source_bucket"`
	SourceScope        string                                `json:"source_scope,omitempty"`
	SourceCollection   string                                `json:"source_collection,omitempty"`
	BucketBindings     []jsonEventingFunctionBucketBinding   `json:"buckets,omitempty"`
	UrlBindings        []jsonEventingFunctionUrlBinding      `json:"curl,omitempty"`
	ConstantBindings   []jsonEventingFunctionConstantBinding `json:"constants,omitempty"`
}

type jsonEventingFunctionBucketBinding struct {
	Alias      string                       `json:"alias"`
	Bucket     string                       `json:"bucket_name"`
	Scope      string                       `json:"scope_name,omitempty"`
	Collection string                       `json:"collection_name,omitempty"`
	Access     EventingFunctionBucketAccess `json:"access"`
}

type jsonEventingFunctionUrlBinding struct {
	Hostname               string `json:"hostname"`
	Alias                  string `json:"value"`
	AuthType               string `json:"auth_type"`
	AllowCookies           bool   `json:"allow_cookies"`
	ValidateSSLCertificate bool   `json:"validate_ssl_certificate"`
	Username               string `json:"username,omitempty"`
	Password               string `json:"password,omitempty"`
	BearerKey              string `json:"bearer_key,omitempty"`
}

type jsonEventingFunctionConstantBinding struct {
	Alias   string `json:"value"`
	Literal string `json:"literal"`
}

type jsonEventingFunctionsStatusApp struct {
	CompositeStatus       EventingFunctionStatus           `json:"composite_status"`
	Name                  string                           `json:"name"`
	NumBootstrappingNodes int                              `json:"num_bootstrapping_nodes"`
	NumDeployedNodes      int                              `json:"num_deployed_nodes"`
	DeploymentStatus      EventingFunctionDeploymentStatus `json:"deployment_status"`
	ProcessingStatus      EventingFunctionProcessingStatus `json:"processing_status"`
	FunctionScope         *jsonEventingFunctionScope       `json:"function_scope,omitempty"`
}

type jsonEventingFunctionsStatus struct {
	Apps             []jsonEventingFunctionsStatusApp `json:"apps"`
	NumEventingNodes int                              `json:"num_eventing_nodes"`
}

func (jf *jsonEventingFunction) MatchesScope(scope *Scope) bool {
	if scope == nil {
		return jf.FunctionScope == nil || (jf.FunctionScope.ScopeName == "*" && jf.FunctionScope.BucketName == "*")
	}

	return jf.FunctionScope != nil && jf.FunctionScope.ScopeName == scope.Name() && jf.FunctionScope.BucketName == scope.BucketName()
}

func (jfs *jsonEventingFunctionsStatusApp) MatchesScope(scope *Scope) bool {
	if scope == nil {
		return jfs.FunctionScope == nil || (jfs.FunctionScope.ScopeName == "*" && jfs.FunctionScope.BucketName == "*")
	}

	return jfs.FunctionScope != nil && jfs.FunctionScope.ScopeName == scope.Name() && jfs.FunctionScope.BucketName == scope.BucketName()
}

type eventingResult interface {
	decodeAndFilter(decoder *json.Decoder, scope *Scope) error
}

type eventingFunctions struct {
	functions []EventingFunction
}

func (efs *eventingFunctions) decodeAndFilter(decoder *json.Decoder, scope *Scope) error {
	var jsonFunctions []jsonEventingFunction
	err := decoder.Decode(&jsonFunctions)
	if err != nil {
		return err
	}

	for _, jsonFunc := range jsonFunctions {
		if jsonFunc.MatchesScope(scope) {
			var function EventingFunction
			function.fromJSONEventingFunction(jsonFunc)
			efs.functions = append(efs.functions, function)
		}
	}
	return nil
}

func (ef *EventingFunction) decodeAndFilter(decoder *json.Decoder, scope *Scope) error {
	err := decoder.Decode(&ef)
	return err
}

func (es *EventingStatus) decodeAndFilter(decoder *json.Decoder, scope *Scope) error {
	var js jsonEventingFunctionsStatus
	err := decoder.Decode(&js)
	if err != nil {
		return err
	}
	var functions []EventingFunctionState
	for _, f := range js.Apps {
		if f.MatchesScope(scope) {
			functions = append(functions, EventingFunctionState{
				Name:                  f.Name,
				Status:                f.CompositeStatus,
				NumBootstrappingNodes: f.NumBootstrappingNodes,
				NumDeployedNodes:      f.NumDeployedNodes,
				DeploymentStatus:      f.DeploymentStatus,
				ProcessingStatus:      f.ProcessingStatus,
			})
		}
	}

	es.NumEventingNodes = js.NumEventingNodes
	es.Functions = functions

	return nil
}

func (ef *EventingFunction) toJSONEventingFunction() jsonEventingFunction {
	var bucketBindings []jsonEventingFunctionBucketBinding
	for _, b := range ef.BucketBindings {
		bucketBindings = append(bucketBindings, jsonEventingFunctionBucketBinding{
			Alias:      b.Alias,
			Bucket:     b.Name.Bucket,
			Scope:      b.Name.Scope,
			Collection: b.Name.Collection,
			Access:     b.Access,
		})
	}
	var urlBindings []jsonEventingFunctionUrlBinding
	for _, b := range ef.UrlBindings {
		urlBindings = append(urlBindings, jsonEventingFunctionUrlBinding{
			Hostname:               b.Hostname,
			Alias:                  b.Alias,
			AuthType:               b.Auth.Method(),
			AllowCookies:           b.AllowCookies,
			ValidateSSLCertificate: b.ValidateSSLCertificate,
			Username:               b.Auth.Username(),
			Password:               b.Auth.Password(),
			BearerKey:              b.Auth.Key(),
		})
	}
	var constantBindings []jsonEventingFunctionConstantBinding
	for _, b := range ef.ConstantBindings {
		constantBindings = append(constantBindings, jsonEventingFunctionConstantBinding(b))
	}

	return jsonEventingFunction{
		Name:               ef.Name,
		Code:               ef.Code,
		Version:            ef.Version,
		EnforceSchema:      ef.EnforceSchema,
		HandlerUUID:        ef.HandlerUUID,
		FunctionInstanceID: ef.FunctionInstanceID,
		Settings: jsonEventingFunctionSettings{
			CPPWorkerThreadCount:   ef.Settings.CPPWorkerThreadCount,
			DCPStreamBoundary:      ef.Settings.DCPStreamBoundary,
			Description:            ef.Settings.Description,
			DeploymentStatus:       ef.Settings.DeploymentStatus,
			ProcessingStatus:       ef.Settings.ProcessingStatus,
			LanguageCompatibility:  ef.Settings.LanguageCompatibility,
			LogLevel:               ef.Settings.LogLevel,
			ExecutionTimeout:       int(ef.Settings.ExecutionTimeout.Seconds()),
			LCBInstCapacity:        ef.Settings.LCBInstCapacity,
			LCBRetryCount:          ef.Settings.LCBRetryCount,
			LCBTimeout:             int(ef.Settings.LCBTimeout.Seconds()),
			QueryConsistency:       ef.Settings.QueryConsistency,
			NumTimerPartitions:     ef.Settings.NumTimerPartitions,
			SockBatchSize:          ef.Settings.SockBatchSize,
			TickDuration:           int(ef.Settings.TickDuration.Milliseconds()),
			TimerContextSize:       ef.Settings.TimerContextSize,
			UserPrefix:             ef.Settings.UserPrefix,
			BucketCacheSize:        ef.Settings.BucketCacheSize,
			BucketCacheAge:         ef.Settings.BucketCacheAge,
			CurlMaxAllowedRespSize: ef.Settings.CurlMaxAllowedRespSize,
			QueryPrepareAll:        ef.Settings.QueryPrepareAll,
			WorkerCount:            ef.Settings.WorkerCount,
			HandlerHeaders:         ef.Settings.HandlerHeaders,
			HandlerFooters:         ef.Settings.HandlerFooters,
			EnableAppLogRotation:   ef.Settings.EnableAppLogRotation,
			AppLogDir:              ef.Settings.AppLogDir,
			AppLogMaxSize:          ef.Settings.AppLogMaxSize,
			AppLogMaxFiles:         ef.Settings.AppLogMaxFiles,
			CheckpointInterval:     int(ef.Settings.CheckpointInterval.Seconds()),
		},
		DeploymentConfig: jsonEventingFunctionDeploymentConfig{
			MetadataBucket:     ef.MetadataKeyspace.Bucket,
			MetadataScope:      ef.MetadataKeyspace.Scope,
			MetadataCollection: ef.MetadataKeyspace.Collection,
			SourceBucket:       ef.SourceKeyspace.Bucket,
			SourceScope:        ef.SourceKeyspace.Scope,
			SourceCollection:   ef.SourceKeyspace.Collection,
			BucketBindings:     bucketBindings,
			UrlBindings:        urlBindings,
			ConstantBindings:   constantBindings,
		},
	}
}

func (ef *EventingFunction) fromJSONEventingFunction(jf jsonEventingFunction) {
	var bucketBindings []EventingFunctionBucketBinding
	for _, b := range jf.DeploymentConfig.BucketBindings {
		bucketBindings = append(bucketBindings, EventingFunctionBucketBinding{
			Alias: b.Alias,
			Name: EventingFunctionKeyspace{
				Bucket:     b.Bucket,
				Scope:      b.Scope,
				Collection: b.Collection,
			},
			Access: b.Access,
		})
	}
	var urlBindings []EventingFunctionUrlBinding
	for _, b := range jf.DeploymentConfig.UrlBindings {
		var auth EventingFunctionUrlAuth
		switch b.AuthType {
		case "no-auth":
			auth = EventingFunctionUrlNoAuth{}
		case "basic":
			auth = EventingFunctionUrlAuthBasic{
				User: b.Username,
			}
		case "digest":
			auth = EventingFunctionUrlAuthDigest{
				User: b.Username,
			}
		case "bearer":
			auth = EventingFunctionUrlAuthBearer{}
		}

		urlBindings = append(urlBindings, EventingFunctionUrlBinding{
			Hostname:               b.Hostname,
			Alias:                  b.Alias,
			Auth:                   auth,
			AllowCookies:           b.AllowCookies,
			ValidateSSLCertificate: b.ValidateSSLCertificate,
		})
	}
	var constantBindings []EventingFunctionConstantBinding
	for _, b := range jf.DeploymentConfig.ConstantBindings {
		constantBindings = append(constantBindings, EventingFunctionConstantBinding(b))
	}

	ef.Name = jf.Name
	ef.Code = jf.Code
	ef.Version = jf.Version
	ef.EnforceSchema = jf.EnforceSchema
	ef.HandlerUUID = jf.HandlerUUID
	ef.FunctionInstanceID = jf.FunctionInstanceID
	ef.Settings = EventingFunctionSettings{
		CPPWorkerThreadCount:   jf.Settings.CPPWorkerThreadCount,
		DCPStreamBoundary:      jf.Settings.DCPStreamBoundary,
		Description:            jf.Settings.Description,
		DeploymentStatus:       jf.Settings.DeploymentStatus,
		ProcessingStatus:       jf.Settings.ProcessingStatus,
		LanguageCompatibility:  jf.Settings.LanguageCompatibility,
		LogLevel:               jf.Settings.LogLevel,
		ExecutionTimeout:       time.Duration(jf.Settings.ExecutionTimeout) * time.Second,
		LCBInstCapacity:        jf.Settings.LCBInstCapacity,
		LCBRetryCount:          jf.Settings.LCBRetryCount,
		LCBTimeout:             time.Duration(jf.Settings.LCBTimeout) * time.Second,
		QueryConsistency:       jf.Settings.QueryConsistency,
		NumTimerPartitions:     jf.Settings.NumTimerPartitions,
		SockBatchSize:          jf.Settings.SockBatchSize,
		TickDuration:           time.Duration(jf.Settings.TickDuration) * time.Millisecond,
		TimerContextSize:       jf.Settings.TimerContextSize,
		UserPrefix:             jf.Settings.UserPrefix,
		BucketCacheSize:        jf.Settings.BucketCacheSize,
		BucketCacheAge:         jf.Settings.BucketCacheAge,
		CurlMaxAllowedRespSize: jf.Settings.CurlMaxAllowedRespSize,
		QueryPrepareAll:        jf.Settings.QueryPrepareAll,
		WorkerCount:            jf.Settings.WorkerCount,
		HandlerHeaders:         jf.Settings.HandlerHeaders,
		HandlerFooters:         jf.Settings.HandlerFooters,
		EnableAppLogRotation:   jf.Settings.EnableAppLogRotation,
		AppLogDir:              jf.Settings.AppLogDir,
		AppLogMaxSize:          jf.Settings.AppLogMaxSize,
		AppLogMaxFiles:         jf.Settings.AppLogMaxFiles,
		CheckpointInterval:     time.Duration(jf.Settings.CheckpointInterval) * time.Second,
	}
	ef.MetadataKeyspace = EventingFunctionKeyspace{
		Bucket:     jf.DeploymentConfig.MetadataBucket,
		Scope:      jf.DeploymentConfig.MetadataScope,
		Collection: jf.DeploymentConfig.MetadataCollection,
	}
	ef.SourceKeyspace = EventingFunctionKeyspace{
		Bucket:     jf.DeploymentConfig.SourceBucket,
		Scope:      jf.DeploymentConfig.SourceScope,
		Collection: jf.DeploymentConfig.SourceCollection,
	}
	ef.BucketBindings = bucketBindings
	ef.UrlBindings = urlBindings
	ef.ConstantBindings = constantBindings
}

func (emp *eventingManagementProviderCore) doMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error) {
	resp, err := emp.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (emp *eventingManagementProviderCore) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read eventing function response body: %s", err)
		return nil
	}

	var baseErr error
	strBody := string(b)
	if strings.Contains(strBody, "ERR_APP_NOT_FOUND_TS") {
		baseErr = ErrEventingFunctionNotFound
	} else if strings.Contains(strBody, "ERR_APP_NOT_DEPLOYED") {
		baseErr = ErrEventingFunctionNotDeployed
	} else if strings.Contains(strBody, "ERR_HANDLER_COMPILATION") {
		baseErr = ErrEventingFunctionCompilationFailure
	} else if strings.Contains(strBody, "ERR_SRC_MB_SAME") {
		baseErr = ErrEventingFunctionIdenticalKeyspace
	} else if strings.Contains(strBody, "ERR_APP_NOT_BOOTSTRAPPED") {
		baseErr = ErrEventingFunctionNotBootstrapped
	} else if strings.Contains(strBody, "ERR_APP_NOT_UNDEPLOYED") {
		baseErr = ErrEventingFunctionDeployed
	} else if strings.Contains(strBody, "ERR_COLLECTION_MISSING") {
		baseErr = ErrCollectionNotFound
	} else if strings.Contains(strBody, "ERR_BUCKET_MISSING") {
		baseErr = ErrBucketNotFound
	} else {
		baseErr = errors.New(string(b))
	}

	return makeGenericMgmtError(baseErr, req, resp, strBody)
}

func (emp *eventingManagementProviderCore) scopedPath(path string, scope *Scope) string {
	if scope == nil {
		return path
	}
	return fmt.Sprintf("%s?bucket=%s&scope=%s", path, url.PathEscape(scope.BucketName()), url.PathEscape(scope.Name()))
}

func (emp *eventingManagementProviderCore) doRequest(scope *Scope, path string, method string, opName string, function *EventingFunction,
	target eventingResult, opts eventingRequestOptions) error {
	start := time.Now()
	defer emp.meter.ValueRecord(meterValueServiceManagement, opName, start)

	if opName != "get_all_functions" && opName != "functions_status" {
		path = emp.scopedPath(path, scope)
	}

	op := "manager_eventing_" + opName
	span := createSpan(emp.tracer, opts.ParentSpan, op, "management")
	span.SetAttribute("db.operation", method+" "+path)
	if scope == nil {
		span.SetAttribute("db.name", "*")
		span.SetAttribute("db.couchbase.scope", "*")
	} else {
		span.SetAttribute("db.name", scope.BucketName())
		span.SetAttribute("db.couchbase.scope", scope.Name())
	}
	defer span.End()

	var b []byte
	if function != nil {
		jsonFunction := function.toJSONEventingFunction()

		// Injecting the function scope for the scope-level UpsertFunction
		if scope != nil {
			jsonFunction.FunctionScope = &jsonEventingFunctionScope{
				ScopeName:  scope.Name(),
				BucketName: scope.BucketName(),
			}
		}

		var err error
		b, err = json.Marshal(jsonFunction)
		if err != nil {
			return err
		}
	}

	req := mgmtRequest{
		Service:       ServiceTypeEventing,
		Method:        method,
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
		Body:          b,
	}
	resp, err := emp.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := emp.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return idxErr
		}

		return makeMgmtBadStatusError("failed eventing "+opName, &req, resp)
	}

	if target != nil {
		jsonDec := json.NewDecoder(resp.Body)
		err = target.decodeAndFilter(jsonDec, scope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (emp *eventingManagementProviderCore) UpsertFunction(scope *Scope, function EventingFunction, opts *UpsertEventingFunctionOptions) error {
	if opts == nil {
		opts = &UpsertEventingFunctionOptions{}
	}

	return emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s", url.PathEscape(function.Name)), "POST",
		"upsert_function", &function, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

func (emp *eventingManagementProviderCore) DropFunction(scope *Scope, name string, opts *DropEventingFunctionOptions) error {
	if opts == nil {
		opts = &DropEventingFunctionOptions{}
	}

	return emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s", url.PathEscape(name)), "DELETE",
		"drop_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

func (emp *eventingManagementProviderCore) DeployFunction(scope *Scope, name string, opts *DeployEventingFunctionOptions) error {
	if opts == nil {
		opts = &DeployEventingFunctionOptions{}
	}

	return emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s/deploy", url.PathEscape(name)), "POST",
		"deploy_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

func (emp *eventingManagementProviderCore) UndeployFunction(scope *Scope, name string, opts *UndeployEventingFunctionOptions) error {
	if opts == nil {
		opts = &UndeployEventingFunctionOptions{}
	}

	return emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s/undeploy", url.PathEscape(name)), "POST",
		"undeploy_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

func (emp *eventingManagementProviderCore) GetAllFunctions(scope *Scope, opts *GetAllEventingFunctionsOptions) ([]EventingFunction, error) {
	if opts == nil {
		opts = &GetAllEventingFunctionsOptions{}
	}

	var functions eventingFunctions
	err := emp.doRequest(scope, "/api/v1/functions", "GET",
		"get_all_functions", nil, &functions, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	if err != nil {
		return nil, err
	}

	return functions.functions, nil
}

func (emp *eventingManagementProviderCore) GetFunction(scope *Scope, name string, opts *GetEventingFunctionOptions) (*EventingFunction, error) {
	if opts == nil {
		opts = &GetEventingFunctionOptions{}
	}

	var function EventingFunction
	err := emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s", url.PathEscape(name)), "GET",
		"get_function", nil, &function, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	if err != nil {
		return nil, err
	}

	return &function, nil
}

func (emp *eventingManagementProviderCore) PauseFunction(scope *Scope, name string, opts *PauseEventingFunctionOptions) error {
	if opts == nil {
		opts = &PauseEventingFunctionOptions{}
	}

	return emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s/pause", url.PathEscape(name)), "POST",
		"pause_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

func (emp *eventingManagementProviderCore) ResumeFunction(scope *Scope, name string, opts *ResumeEventingFunctionOptions) error {
	if opts == nil {
		opts = &ResumeEventingFunctionOptions{}
	}

	return emp.doRequest(scope, fmt.Sprintf("/api/v1/functions/%s/resume", url.PathEscape(name)), "POST",
		"resume_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

func (emp *eventingManagementProviderCore) FunctionsStatus(scope *Scope, opts *EventingFunctionsStatusOptions) (*EventingStatus, error) {
	if opts == nil {
		opts = &EventingFunctionsStatusOptions{}
	}

	var functions EventingStatus
	err := emp.doRequest(scope, "/api/v1/status", "GET",
		"functions_status", nil, &functions, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	if err != nil {
		return nil, err
	}

	return &functions, nil
}
