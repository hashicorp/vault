package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder provides operations to manage the onPremisesSynchronization property of the microsoft.graph.directory entity.
type OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetQueryParameters read the properties and relationships of an onPremisesDirectorySynchronization object.
type OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetQueryParameters
}
// OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewOnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderInternal instantiates a new OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder and sets the default values.
func NewOnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) {
    m := &OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/onPremisesSynchronization/{onPremisesDirectorySynchronization%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewOnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder instantiates a new OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder and sets the default values.
func NewOnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property onPremisesSynchronization for directory
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read the properties and relationships of an onPremisesDirectorySynchronization object.
// returns a OnPremisesDirectorySynchronizationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/onpremisesdirectorysynchronization-get?view=graph-rest-1.0
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnPremisesDirectorySynchronizationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOnPremisesDirectorySynchronizationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnPremisesDirectorySynchronizationable), nil
}
// Patch update the properties of an onPremisesDirectorySynchronization object.
// returns a OnPremisesDirectorySynchronizationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/onpremisesdirectorysynchronization-update?view=graph-rest-1.0
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnPremisesDirectorySynchronizationable, requestConfiguration *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnPremisesDirectorySynchronizationable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOnPremisesDirectorySynchronizationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnPremisesDirectorySynchronizationable), nil
}
// ToDeleteRequestInformation delete navigation property onPremisesSynchronization for directory
// returns a *RequestInformation when successful
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of an onPremisesDirectorySynchronization object.
// returns a *RequestInformation when successful
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of an onPremisesDirectorySynchronization object.
// returns a *RequestInformation when successful
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnPremisesDirectorySynchronizationable, requestConfiguration *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder when successful
func (m *OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) WithUrl(rawUrl string)(*OnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder) {
    return NewOnPremisesSynchronizationOnPremisesDirectorySynchronizationItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
