package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AuditEventsGetAuditCategoriesRequestBuilder provides operations to call the getAuditCategories method.
type AuditEventsGetAuditCategoriesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AuditEventsGetAuditCategoriesRequestBuilderGetQueryParameters not yet documented
type AuditEventsGetAuditCategoriesRequestBuilderGetQueryParameters struct {
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
// AuditEventsGetAuditCategoriesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuditEventsGetAuditCategoriesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AuditEventsGetAuditCategoriesRequestBuilderGetQueryParameters
}
// NewAuditEventsGetAuditCategoriesRequestBuilderInternal instantiates a new AuditEventsGetAuditCategoriesRequestBuilder and sets the default values.
func NewAuditEventsGetAuditCategoriesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuditEventsGetAuditCategoriesRequestBuilder) {
    m := &AuditEventsGetAuditCategoriesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/auditEvents/getAuditCategories(){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewAuditEventsGetAuditCategoriesRequestBuilder instantiates a new AuditEventsGetAuditCategoriesRequestBuilder and sets the default values.
func NewAuditEventsGetAuditCategoriesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuditEventsGetAuditCategoriesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAuditEventsGetAuditCategoriesRequestBuilderInternal(urlParams, requestAdapter)
}
// Get not yet documented
// Deprecated: This method is obsolete. Use GetAsGetAuditCategoriesGetResponse instead.
// returns a AuditEventsGetAuditCategoriesResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-auditing-auditevent-getauditcategories?view=graph-rest-1.0
func (m *AuditEventsGetAuditCategoriesRequestBuilder) Get(ctx context.Context, requestConfiguration *AuditEventsGetAuditCategoriesRequestBuilderGetRequestConfiguration)(AuditEventsGetAuditCategoriesResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateAuditEventsGetAuditCategoriesResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(AuditEventsGetAuditCategoriesResponseable), nil
}
// GetAsGetAuditCategoriesGetResponse not yet documented
// returns a AuditEventsGetAuditCategoriesGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-auditing-auditevent-getauditcategories?view=graph-rest-1.0
func (m *AuditEventsGetAuditCategoriesRequestBuilder) GetAsGetAuditCategoriesGetResponse(ctx context.Context, requestConfiguration *AuditEventsGetAuditCategoriesRequestBuilderGetRequestConfiguration)(AuditEventsGetAuditCategoriesGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateAuditEventsGetAuditCategoriesGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(AuditEventsGetAuditCategoriesGetResponseable), nil
}
// ToGetRequestInformation not yet documented
// returns a *RequestInformation when successful
func (m *AuditEventsGetAuditCategoriesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AuditEventsGetAuditCategoriesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *AuditEventsGetAuditCategoriesRequestBuilder when successful
func (m *AuditEventsGetAuditCategoriesRequestBuilder) WithUrl(rawUrl string)(*AuditEventsGetAuditCategoriesRequestBuilder) {
    return NewAuditEventsGetAuditCategoriesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
