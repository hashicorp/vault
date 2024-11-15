package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430 "github.com/microsoftgraph/msgraph-sdk-go/models/identitygovernance"
)

// LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder provides operations to manage the taskReports property of the microsoft.graph.identityGovernance.workflow entity.
type LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetQueryParameters represents the aggregation of task execution data for tasks within a workflow object.
type LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetQueryParameters
}
// NewLifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderInternal instantiates a new LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder and sets the default values.
func NewLifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) {
    m := &LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/lifecycleWorkflows/workflows/{workflow%2Did}/taskReports/{taskReport%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewLifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder instantiates a new LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder and sets the default values.
func NewLifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get represents the aggregation of task execution data for tasks within a workflow object.
// returns a TaskReportable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) Get(ctx context.Context, requestConfiguration *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetRequestConfiguration)(ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.TaskReportable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.CreateTaskReportFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.TaskReportable), nil
}
// Task provides operations to manage the task property of the microsoft.graph.identityGovernance.taskReport entity.
// returns a *LifecycleWorkflowsWorkflowsItemTaskReportsItemTaskRequestBuilder when successful
func (m *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) Task()(*LifecycleWorkflowsWorkflowsItemTaskReportsItemTaskRequestBuilder) {
    return NewLifecycleWorkflowsWorkflowsItemTaskReportsItemTaskRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TaskDefinition provides operations to manage the taskDefinition property of the microsoft.graph.identityGovernance.taskReport entity.
// returns a *LifecycleWorkflowsWorkflowsItemTaskReportsItemTaskDefinitionRequestBuilder when successful
func (m *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) TaskDefinition()(*LifecycleWorkflowsWorkflowsItemTaskReportsItemTaskDefinitionRequestBuilder) {
    return NewLifecycleWorkflowsWorkflowsItemTaskReportsItemTaskDefinitionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TaskProcessingResults provides operations to manage the taskProcessingResults property of the microsoft.graph.identityGovernance.taskReport entity.
// returns a *LifecycleWorkflowsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder when successful
func (m *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) TaskProcessingResults()(*LifecycleWorkflowsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) {
    return NewLifecycleWorkflowsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation represents the aggregation of task execution data for tasks within a workflow object.
// returns a *RequestInformation when successful
func (m *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder when successful
func (m *LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) WithUrl(rawUrl string)(*LifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder) {
    return NewLifecycleWorkflowsWorkflowsItemTaskReportsTaskReportItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
