package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder provides operations to count the resources in the collection.
type MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetQueryParameters get the number of the resource
type MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetQueryParameters struct {
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
}
// MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetQueryParameters
}
// NewMobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderInternal instantiates a new MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder and sets the default values.
func NewMobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder) {
    m := &MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/graph.microsoftStoreForBusinessApp/$count{?%24filter,%24search}", pathParameters),
    }
    return m
}
// NewMobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder instantiates a new MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder and sets the default values.
func NewMobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get the number of the resource
// returns a *int32 when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder) Get(ctx context.Context, requestConfiguration *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetRequestConfiguration)(*int32, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.SendPrimitive(ctx, requestInfo, "int32", errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(*int32), nil
}
// ToGetRequestInformation get the number of the resource
// returns a *RequestInformation when successful
func (m *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "text/plain;q=0.9")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder when successful
func (m *MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder) WithUrl(rawUrl string)(*MobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder) {
    return NewMobileAppsGraphMicrosoftStoreForBusinessAppCountRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
