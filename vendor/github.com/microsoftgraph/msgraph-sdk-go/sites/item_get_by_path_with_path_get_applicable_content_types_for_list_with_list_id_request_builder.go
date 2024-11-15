package sites

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder provides operations to call the getApplicableContentTypesForList method.
type ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetQueryParameters get site contentTypes that can be added to a list.
type ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetQueryParameters struct {
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
// ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetQueryParameters
}
// NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderInternal instantiates a new ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder and sets the default values.
func NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, listId *string)(*ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) {
    m := &ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/sites/{site%2Did}/getByPath(path='{path}')/getApplicableContentTypesForList(listId='{listId}'){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    if listId != nil {
        m.BaseRequestBuilder.PathParameters["listId"] = *listId
    }
    return m
}
// NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder instantiates a new ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder and sets the default values.
func NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get get site contentTypes that can be added to a list.
// Deprecated: This method is obsolete. Use GetAsGetApplicableContentTypesForListWithListIdGetResponse instead.
// returns a ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/site-getapplicablecontenttypesforlist?view=graph-rest-1.0
func (m *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetRequestConfiguration)(ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemGetByPathWithPathGetApplicableContentTypesForListWithListIdResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdResponseable), nil
}
// GetAsGetApplicableContentTypesForListWithListIdGetResponse get site contentTypes that can be added to a list.
// returns a ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/site-getapplicablecontenttypesforlist?view=graph-rest-1.0
func (m *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) GetAsGetApplicableContentTypesForListWithListIdGetResponse(ctx context.Context, requestConfiguration *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetRequestConfiguration)(ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemGetByPathWithPathGetApplicableContentTypesForListWithListIdGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdGetResponseable), nil
}
// ToGetRequestInformation get site contentTypes that can be added to a list.
// returns a *RequestInformation when successful
func (m *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder when successful
func (m *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) WithUrl(rawUrl string)(*ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) {
    return NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
