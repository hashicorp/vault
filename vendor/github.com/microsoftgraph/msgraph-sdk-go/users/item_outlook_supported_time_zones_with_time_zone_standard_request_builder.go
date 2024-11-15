package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder provides operations to call the supportedTimeZones method.
type ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetQueryParameters invoke function supportedTimeZones
type ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetQueryParameters struct {
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
// ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetQueryParameters
}
// NewItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderInternal instantiates a new ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder and sets the default values.
func NewItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, timeZoneStandard *string)(*ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) {
    m := &ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/outlook/supportedTimeZones(TimeZoneStandard='{TimeZoneStandard}'){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    if timeZoneStandard != nil {
        m.BaseRequestBuilder.PathParameters["TimeZoneStandard"] = *timeZoneStandard
    }
    return m
}
// NewItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder instantiates a new ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder and sets the default values.
func NewItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get invoke function supportedTimeZones
// Deprecated: This method is obsolete. Use GetAsSupportedTimeZonesWithTimeZoneStandardGetResponse instead.
// returns a ItemOutlookSupportedTimeZonesWithTimeZoneStandardResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetRequestConfiguration)(ItemOutlookSupportedTimeZonesWithTimeZoneStandardResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemOutlookSupportedTimeZonesWithTimeZoneStandardResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemOutlookSupportedTimeZonesWithTimeZoneStandardResponseable), nil
}
// GetAsSupportedTimeZonesWithTimeZoneStandardGetResponse invoke function supportedTimeZones
// returns a ItemOutlookSupportedTimeZonesWithTimeZoneStandardGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) GetAsSupportedTimeZonesWithTimeZoneStandardGetResponse(ctx context.Context, requestConfiguration *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetRequestConfiguration)(ItemOutlookSupportedTimeZonesWithTimeZoneStandardGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemOutlookSupportedTimeZonesWithTimeZoneStandardGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemOutlookSupportedTimeZonesWithTimeZoneStandardGetResponseable), nil
}
// ToGetRequestInformation invoke function supportedTimeZones
// returns a *RequestInformation when successful
func (m *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder when successful
func (m *ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) WithUrl(rawUrl string)(*ItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder) {
    return NewItemOutlookSupportedTimeZonesWithTimeZoneStandardRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
