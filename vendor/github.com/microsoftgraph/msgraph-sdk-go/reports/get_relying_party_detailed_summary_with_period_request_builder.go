package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder provides operations to call the getRelyingPartyDetailedSummary method.
type GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetQueryParameters get a summary of AD FS relying parties information.
type GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetQueryParameters struct {
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
// GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetQueryParameters
}
// NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilderInternal instantiates a new GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder and sets the default values.
func NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, period *string)(*GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) {
    m := &GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/getRelyingPartyDetailedSummary(period='{period}'){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    if period != nil {
        m.BaseRequestBuilder.PathParameters["period"] = *period
    }
    return m
}
// NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilder instantiates a new GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder and sets the default values.
func NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get get a summary of AD FS relying parties information.
// Deprecated: This method is obsolete. Use GetAsGetRelyingPartyDetailedSummaryWithPeriodGetResponse instead.
// returns a GetRelyingPartyDetailedSummaryWithPeriodResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/reportroot-getrelyingpartydetailedsummary?view=graph-rest-1.0
func (m *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) Get(ctx context.Context, requestConfiguration *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetRequestConfiguration)(GetRelyingPartyDetailedSummaryWithPeriodResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateGetRelyingPartyDetailedSummaryWithPeriodResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(GetRelyingPartyDetailedSummaryWithPeriodResponseable), nil
}
// GetAsGetRelyingPartyDetailedSummaryWithPeriodGetResponse get a summary of AD FS relying parties information.
// returns a GetRelyingPartyDetailedSummaryWithPeriodGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/reportroot-getrelyingpartydetailedsummary?view=graph-rest-1.0
func (m *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) GetAsGetRelyingPartyDetailedSummaryWithPeriodGetResponse(ctx context.Context, requestConfiguration *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetRequestConfiguration)(GetRelyingPartyDetailedSummaryWithPeriodGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateGetRelyingPartyDetailedSummaryWithPeriodGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(GetRelyingPartyDetailedSummaryWithPeriodGetResponseable), nil
}
// ToGetRequestInformation get a summary of AD FS relying parties information.
// returns a *RequestInformation when successful
func (m *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder when successful
func (m *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) WithUrl(rawUrl string)(*GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) {
    return NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
