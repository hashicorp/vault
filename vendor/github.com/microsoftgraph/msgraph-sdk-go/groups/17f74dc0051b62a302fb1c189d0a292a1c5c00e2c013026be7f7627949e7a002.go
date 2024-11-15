package groups

import (
    "context"
    i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274 "strconv"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder provides operations to call the getRecentNotebooks method.
type ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetQueryParameters get a list of recentNotebook instances that have been accessed by the signed-in user.
type ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetQueryParameters
}
// NewItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderInternal instantiates a new ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder and sets the default values.
func NewItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, includePersonalNotebooks *bool)(*ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) {
    m := &ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/sites/{site%2Did}/onenote/notebooks/getRecentNotebooks(includePersonalNotebooks={includePersonalNotebooks}){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    if includePersonalNotebooks != nil {
        m.BaseRequestBuilder.PathParameters["includePersonalNotebooks"] = i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274.FormatBool(*includePersonalNotebooks)
    }
    return m
}
// NewItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder instantiates a new ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder and sets the default values.
func NewItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get get a list of recentNotebook instances that have been accessed by the signed-in user.
// Deprecated: This method is obsolete. Use GetAsGetRecentNotebooksWithIncludePersonalNotebooksGetResponse instead.
// returns a ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/notebook-getrecentnotebooks?view=graph-rest-1.0
func (m *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetRequestConfiguration)(ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksResponseable), nil
}
// GetAsGetRecentNotebooksWithIncludePersonalNotebooksGetResponse get a list of recentNotebook instances that have been accessed by the signed-in user.
// returns a ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/notebook-getrecentnotebooks?view=graph-rest-1.0
func (m *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) GetAsGetRecentNotebooksWithIncludePersonalNotebooksGetResponse(ctx context.Context, requestConfiguration *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetRequestConfiguration)(ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksGetResponseable), nil
}
// ToGetRequestInformation get a list of recentNotebook instances that have been accessed by the signed-in user.
// returns a *RequestInformation when successful
func (m *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder when successful
func (m *ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) WithUrl(rawUrl string)(*ItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder) {
    return NewItemSitesItemOnenoteNotebooksGetRecentNotebooksWithIncludePersonalNotebooksRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
