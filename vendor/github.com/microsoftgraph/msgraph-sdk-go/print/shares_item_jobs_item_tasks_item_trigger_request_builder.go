package print

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SharesItemJobsItemTasksItemTriggerRequestBuilder provides operations to manage the trigger property of the microsoft.graph.printTask entity.
type SharesItemJobsItemTasksItemTriggerRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SharesItemJobsItemTasksItemTriggerRequestBuilderGetQueryParameters the printTaskTrigger that triggered this task's execution. Read-only.
type SharesItemJobsItemTasksItemTriggerRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SharesItemJobsItemTasksItemTriggerRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SharesItemJobsItemTasksItemTriggerRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SharesItemJobsItemTasksItemTriggerRequestBuilderGetQueryParameters
}
// NewSharesItemJobsItemTasksItemTriggerRequestBuilderInternal instantiates a new SharesItemJobsItemTasksItemTriggerRequestBuilder and sets the default values.
func NewSharesItemJobsItemTasksItemTriggerRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemJobsItemTasksItemTriggerRequestBuilder) {
    m := &SharesItemJobsItemTasksItemTriggerRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/print/shares/{printerShare%2Did}/jobs/{printJob%2Did}/tasks/{printTask%2Did}/trigger{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewSharesItemJobsItemTasksItemTriggerRequestBuilder instantiates a new SharesItemJobsItemTasksItemTriggerRequestBuilder and sets the default values.
func NewSharesItemJobsItemTasksItemTriggerRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemJobsItemTasksItemTriggerRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSharesItemJobsItemTasksItemTriggerRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the printTaskTrigger that triggered this task's execution. Read-only.
// returns a PrintTaskTriggerable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SharesItemJobsItemTasksItemTriggerRequestBuilder) Get(ctx context.Context, requestConfiguration *SharesItemJobsItemTasksItemTriggerRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintTaskTriggerable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrintTaskTriggerFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintTaskTriggerable), nil
}
// ToGetRequestInformation the printTaskTrigger that triggered this task's execution. Read-only.
// returns a *RequestInformation when successful
func (m *SharesItemJobsItemTasksItemTriggerRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SharesItemJobsItemTasksItemTriggerRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SharesItemJobsItemTasksItemTriggerRequestBuilder when successful
func (m *SharesItemJobsItemTasksItemTriggerRequestBuilder) WithUrl(rawUrl string)(*SharesItemJobsItemTasksItemTriggerRequestBuilder) {
    return NewSharesItemJobsItemTasksItemTriggerRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
