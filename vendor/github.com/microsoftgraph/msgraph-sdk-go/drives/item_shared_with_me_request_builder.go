package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSharedWithMeRequestBuilder provides operations to call the sharedWithMe method.
type ItemSharedWithMeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSharedWithMeRequestBuilderGetQueryParameters get a list of driveItem objects shared with the owner of a drive. The driveItems returned from the sharedWithMe method always include the remoteItem facet that indicates they are items from a different drive.
type ItemSharedWithMeRequestBuilderGetQueryParameters struct {
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
// ItemSharedWithMeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSharedWithMeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSharedWithMeRequestBuilderGetQueryParameters
}
// NewItemSharedWithMeRequestBuilderInternal instantiates a new ItemSharedWithMeRequestBuilder and sets the default values.
func NewItemSharedWithMeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSharedWithMeRequestBuilder) {
    m := &ItemSharedWithMeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/sharedWithMe(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSharedWithMeRequestBuilder instantiates a new ItemSharedWithMeRequestBuilder and sets the default values.
func NewItemSharedWithMeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSharedWithMeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSharedWithMeRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a list of driveItem objects shared with the owner of a drive. The driveItems returned from the sharedWithMe method always include the remoteItem facet that indicates they are items from a different drive.
// Deprecated: This method is obsolete. Use GetAsSharedWithMeGetResponse instead.
// returns a ItemSharedWithMeResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/drive-sharedwithme?view=graph-rest-1.0
func (m *ItemSharedWithMeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSharedWithMeRequestBuilderGetRequestConfiguration)(ItemSharedWithMeResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSharedWithMeResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSharedWithMeResponseable), nil
}
// GetAsSharedWithMeGetResponse get a list of driveItem objects shared with the owner of a drive. The driveItems returned from the sharedWithMe method always include the remoteItem facet that indicates they are items from a different drive.
// returns a ItemSharedWithMeGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/drive-sharedwithme?view=graph-rest-1.0
func (m *ItemSharedWithMeRequestBuilder) GetAsSharedWithMeGetResponse(ctx context.Context, requestConfiguration *ItemSharedWithMeRequestBuilderGetRequestConfiguration)(ItemSharedWithMeGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSharedWithMeGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSharedWithMeGetResponseable), nil
}
// ToGetRequestInformation get a list of driveItem objects shared with the owner of a drive. The driveItems returned from the sharedWithMe method always include the remoteItem facet that indicates they are items from a different drive.
// returns a *RequestInformation when successful
func (m *ItemSharedWithMeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSharedWithMeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSharedWithMeRequestBuilder when successful
func (m *ItemSharedWithMeRequestBuilder) WithUrl(rawUrl string)(*ItemSharedWithMeRequestBuilder) {
    return NewItemSharedWithMeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
