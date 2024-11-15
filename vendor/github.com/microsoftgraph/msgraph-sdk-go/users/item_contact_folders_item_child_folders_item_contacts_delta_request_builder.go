package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder provides operations to call the delta method.
type ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetQueryParameters get a set of contacts that have been added, deleted, or updated in a specified folder. A delta function call for contacts in a folder is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the contacts in that folder. This allows you to maintain and synchronize a local store of a user's contacts without having to fetch the entire set of contacts from the server every time.  
type ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetQueryParameters struct {
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
// ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetQueryParameters
}
// NewItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderInternal instantiates a new ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder and sets the default values.
func NewItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) {
    m := &ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/contactFolders/{contactFolder%2Did}/childFolders/{contactFolder%2Did1}/contacts/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder instantiates a new ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder and sets the default values.
func NewItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a set of contacts that have been added, deleted, or updated in a specified folder. A delta function call for contacts in a folder is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the contacts in that folder. This allows you to maintain and synchronize a local store of a user's contacts without having to fetch the entire set of contacts from the server every time.  
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ItemContactFoldersItemChildFoldersItemContactsDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/contact-delta?view=graph-rest-1.0
func (m *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetRequestConfiguration)(ItemContactFoldersItemChildFoldersItemContactsDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemContactFoldersItemChildFoldersItemContactsDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemContactFoldersItemChildFoldersItemContactsDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a set of contacts that have been added, deleted, or updated in a specified folder. A delta function call for contacts in a folder is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the contacts in that folder. This allows you to maintain and synchronize a local store of a user's contacts without having to fetch the entire set of contacts from the server every time.  
// returns a ItemContactFoldersItemChildFoldersItemContactsDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/contact-delta?view=graph-rest-1.0
func (m *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetRequestConfiguration)(ItemContactFoldersItemChildFoldersItemContactsDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemContactFoldersItemChildFoldersItemContactsDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemContactFoldersItemChildFoldersItemContactsDeltaGetResponseable), nil
}
// ToGetRequestInformation get a set of contacts that have been added, deleted, or updated in a specified folder. A delta function call for contacts in a folder is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the contacts in that folder. This allows you to maintain and synchronize a local store of a user's contacts without having to fetch the entire set of contacts from the server every time.  
// returns a *RequestInformation when successful
func (m *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder when successful
func (m *ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) WithUrl(rawUrl string)(*ItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder) {
    return NewItemContactFoldersItemChildFoldersItemContactsDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
