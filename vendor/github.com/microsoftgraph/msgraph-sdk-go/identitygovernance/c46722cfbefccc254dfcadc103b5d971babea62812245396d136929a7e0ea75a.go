package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430 "github.com/microsoftgraph/msgraph-sdk-go/models/identitygovernance"
)

// LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder provides operations to manage the taskProcessingResults property of the microsoft.graph.identityGovernance.taskReport entity.
type LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetQueryParameters the related lifecycle workflow taskProcessingResults.
type LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetQueryParameters
}
// ByTaskProcessingResultId provides operations to manage the taskProcessingResults property of the microsoft.graph.identityGovernance.taskReport entity.
// returns a *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsTaskProcessingResultItemRequestBuilder when successful
func (m *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) ByTaskProcessingResultId(taskProcessingResultId string)(*LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsTaskProcessingResultItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if taskProcessingResultId != "" {
        urlTplParams["taskProcessingResult%2Did"] = taskProcessingResultId
    }
    return NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsTaskProcessingResultItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderInternal instantiates a new LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder and sets the default values.
func NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) {
    m := &LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/lifecycleWorkflows/deletedItems/workflows/{workflow%2Did}/taskReports/{taskReport%2Did}/taskProcessingResults{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder instantiates a new LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder and sets the default values.
func NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsCountRequestBuilder when successful
func (m *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) Count()(*LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsCountRequestBuilder) {
    return NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the related lifecycle workflow taskProcessingResults.
// returns a TaskProcessingResultCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) Get(ctx context.Context, requestConfiguration *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetRequestConfiguration)(ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.TaskProcessingResultCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.CreateTaskProcessingResultCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.TaskProcessingResultCollectionResponseable), nil
}
// ToGetRequestInformation the related lifecycle workflow taskProcessingResults.
// returns a *RequestInformation when successful
func (m *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder when successful
func (m *LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) WithUrl(rawUrl string)(*LifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder) {
    return NewLifecycleWorkflowsDeletedItemsWorkflowsItemTaskReportsItemTaskProcessingResultsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
