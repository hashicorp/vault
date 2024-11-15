package employeeexperience

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder provides operations to manage the engagementAsyncOperations property of the microsoft.graph.employeeExperience entity.
type EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetQueryParameters get an engagementAsyncOperation to track a long-running operation request.
type EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetQueryParameters
}
// EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewEngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderInternal instantiates a new EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder and sets the default values.
func NewEngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) {
    m := &EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/employeeExperience/engagementAsyncOperations/{engagementAsyncOperation%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewEngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder instantiates a new EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder and sets the default values.
func NewEngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property engagementAsyncOperations for employeeExperience
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get get an engagementAsyncOperation to track a long-running operation request.
// returns a EngagementAsyncOperationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/engagementasyncoperation-get?view=graph-rest-1.0
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EngagementAsyncOperationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEngagementAsyncOperationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EngagementAsyncOperationable), nil
}
// Patch update the navigation property engagementAsyncOperations in employeeExperience
// returns a EngagementAsyncOperationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EngagementAsyncOperationable, requestConfiguration *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EngagementAsyncOperationable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEngagementAsyncOperationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EngagementAsyncOperationable), nil
}
// ToDeleteRequestInformation delete navigation property engagementAsyncOperations for employeeExperience
// returns a *RequestInformation when successful
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get an engagementAsyncOperation to track a long-running operation request.
// returns a *RequestInformation when successful
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property engagementAsyncOperations in employeeExperience
// returns a *RequestInformation when successful
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EngagementAsyncOperationable, requestConfiguration *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder when successful
func (m *EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) WithUrl(rawUrl string)(*EngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder) {
    return NewEngagementAsyncOperationsEngagementAsyncOperationItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
