package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemListItemsDeltaRequestBuilder provides operations to call the delta method.
type ItemListItemsDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemListItemsDeltaRequestBuilderGetQueryParameters get newly created, updated, or deleted list items without having to perform a full read of the entire items collection. Your app begins by calling delta without any parameters.The service starts enumerating the hierarchy of the list, returning pages of items, and either an @odata.nextLink or an @odata.deltaLink.Your app should continue calling with the @odata.nextLink until you see an @odata.deltaLink returned. After you received all the changes, you can apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. The delta feed shows the latest state for each item, not each change. If an item was renamed twice, it only shows up once, with its latest name.The same item might appear more than once in a delta feed, for various reasons. You should use the last occurrence you see. Items with this property should be removed from your local state.
type ItemListItemsDeltaRequestBuilderGetQueryParameters struct {
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
// ItemListItemsDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemListItemsDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemListItemsDeltaRequestBuilderGetQueryParameters
}
// NewItemListItemsDeltaRequestBuilderInternal instantiates a new ItemListItemsDeltaRequestBuilder and sets the default values.
func NewItemListItemsDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemListItemsDeltaRequestBuilder) {
    m := &ItemListItemsDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/list/items/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemListItemsDeltaRequestBuilder instantiates a new ItemListItemsDeltaRequestBuilder and sets the default values.
func NewItemListItemsDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemListItemsDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemListItemsDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get newly created, updated, or deleted list items without having to perform a full read of the entire items collection. Your app begins by calling delta without any parameters.The service starts enumerating the hierarchy of the list, returning pages of items, and either an @odata.nextLink or an @odata.deltaLink.Your app should continue calling with the @odata.nextLink until you see an @odata.deltaLink returned. After you received all the changes, you can apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. The delta feed shows the latest state for each item, not each change. If an item was renamed twice, it only shows up once, with its latest name.The same item might appear more than once in a delta feed, for various reasons. You should use the last occurrence you see. Items with this property should be removed from your local state.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ItemListItemsDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/listitem-delta?view=graph-rest-1.0
func (m *ItemListItemsDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemListItemsDeltaRequestBuilderGetRequestConfiguration)(ItemListItemsDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemListItemsDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemListItemsDeltaResponseable), nil
}
// GetAsDeltaGetResponse get newly created, updated, or deleted list items without having to perform a full read of the entire items collection. Your app begins by calling delta without any parameters.The service starts enumerating the hierarchy of the list, returning pages of items, and either an @odata.nextLink or an @odata.deltaLink.Your app should continue calling with the @odata.nextLink until you see an @odata.deltaLink returned. After you received all the changes, you can apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. The delta feed shows the latest state for each item, not each change. If an item was renamed twice, it only shows up once, with its latest name.The same item might appear more than once in a delta feed, for various reasons. You should use the last occurrence you see. Items with this property should be removed from your local state.
// returns a ItemListItemsDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/listitem-delta?view=graph-rest-1.0
func (m *ItemListItemsDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ItemListItemsDeltaRequestBuilderGetRequestConfiguration)(ItemListItemsDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemListItemsDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemListItemsDeltaGetResponseable), nil
}
// ToGetRequestInformation get newly created, updated, or deleted list items without having to perform a full read of the entire items collection. Your app begins by calling delta without any parameters.The service starts enumerating the hierarchy of the list, returning pages of items, and either an @odata.nextLink or an @odata.deltaLink.Your app should continue calling with the @odata.nextLink until you see an @odata.deltaLink returned. After you received all the changes, you can apply them to your local state.To check for changes in the future, call delta again with the @odata.deltaLink from the previous response. The delta feed shows the latest state for each item, not each change. If an item was renamed twice, it only shows up once, with its latest name.The same item might appear more than once in a delta feed, for various reasons. You should use the last occurrence you see. Items with this property should be removed from your local state.
// returns a *RequestInformation when successful
func (m *ItemListItemsDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemListItemsDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemListItemsDeltaRequestBuilder when successful
func (m *ItemListItemsDeltaRequestBuilder) WithUrl(rawUrl string)(*ItemListItemsDeltaRequestBuilder) {
    return NewItemListItemsDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
