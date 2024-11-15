package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemDeltaRequestBuilder provides operations to call the delta method.
type ItemItemsItemDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemDeltaRequestBuilderGetQueryParameters track changes in a driveItem and its children over time. Your app begins by calling delta without any parameters.The service starts enumerating the drive's hierarchy, returning pages of items and either an @odata.nextLink or an @odata.deltaLink, as described below.Your app should continue calling with the @odata.nextLink until you no longer see an @odata.nextLink returned, or you see a response with an empty set of changes. After you have finished receiving all the changes, you may apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. Deleted items are returned with the deleted facet.Items with this property set should be removed from your local state.
type ItemItemsItemDeltaRequestBuilderGetQueryParameters struct {
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
// ItemItemsItemDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemItemsItemDeltaRequestBuilderGetQueryParameters
}
// NewItemItemsItemDeltaRequestBuilderInternal instantiates a new ItemItemsItemDeltaRequestBuilder and sets the default values.
func NewItemItemsItemDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemDeltaRequestBuilder) {
    m := &ItemItemsItemDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemItemsItemDeltaRequestBuilder instantiates a new ItemItemsItemDeltaRequestBuilder and sets the default values.
func NewItemItemsItemDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get track changes in a driveItem and its children over time. Your app begins by calling delta without any parameters.The service starts enumerating the drive's hierarchy, returning pages of items and either an @odata.nextLink or an @odata.deltaLink, as described below.Your app should continue calling with the @odata.nextLink until you no longer see an @odata.nextLink returned, or you see a response with an empty set of changes. After you have finished receiving all the changes, you may apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. Deleted items are returned with the deleted facet.Items with this property set should be removed from your local state.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ItemItemsItemDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/driveitem-delta?view=graph-rest-1.0
func (m *ItemItemsItemDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemDeltaRequestBuilderGetRequestConfiguration)(ItemItemsItemDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemItemsItemDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemItemsItemDeltaResponseable), nil
}
// GetAsDeltaGetResponse track changes in a driveItem and its children over time. Your app begins by calling delta without any parameters.The service starts enumerating the drive's hierarchy, returning pages of items and either an @odata.nextLink or an @odata.deltaLink, as described below.Your app should continue calling with the @odata.nextLink until you no longer see an @odata.nextLink returned, or you see a response with an empty set of changes. After you have finished receiving all the changes, you may apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. Deleted items are returned with the deleted facet.Items with this property set should be removed from your local state.
// returns a ItemItemsItemDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/driveitem-delta?view=graph-rest-1.0
func (m *ItemItemsItemDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ItemItemsItemDeltaRequestBuilderGetRequestConfiguration)(ItemItemsItemDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemItemsItemDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemItemsItemDeltaGetResponseable), nil
}
// ToGetRequestInformation track changes in a driveItem and its children over time. Your app begins by calling delta without any parameters.The service starts enumerating the drive's hierarchy, returning pages of items and either an @odata.nextLink or an @odata.deltaLink, as described below.Your app should continue calling with the @odata.nextLink until you no longer see an @odata.nextLink returned, or you see a response with an empty set of changes. After you have finished receiving all the changes, you may apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. Deleted items are returned with the deleted facet.Items with this property set should be removed from your local state.
// returns a *RequestInformation when successful
func (m *ItemItemsItemDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemItemsItemDeltaRequestBuilder when successful
func (m *ItemItemsItemDeltaRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemDeltaRequestBuilder) {
    return NewItemItemsItemDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
