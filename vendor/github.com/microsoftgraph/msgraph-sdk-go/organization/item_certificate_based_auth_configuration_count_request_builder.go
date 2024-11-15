package organization

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemCertificateBasedAuthConfigurationCountRequestBuilder provides operations to count the resources in the collection.
type ItemCertificateBasedAuthConfigurationCountRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemCertificateBasedAuthConfigurationCountRequestBuilderGetQueryParameters get the number of the resource
type ItemCertificateBasedAuthConfigurationCountRequestBuilderGetQueryParameters struct {
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
}
// ItemCertificateBasedAuthConfigurationCountRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemCertificateBasedAuthConfigurationCountRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemCertificateBasedAuthConfigurationCountRequestBuilderGetQueryParameters
}
// NewItemCertificateBasedAuthConfigurationCountRequestBuilderInternal instantiates a new ItemCertificateBasedAuthConfigurationCountRequestBuilder and sets the default values.
func NewItemCertificateBasedAuthConfigurationCountRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemCertificateBasedAuthConfigurationCountRequestBuilder) {
    m := &ItemCertificateBasedAuthConfigurationCountRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/organization/{organization%2Did}/certificateBasedAuthConfiguration/$count{?%24filter,%24search}", pathParameters),
    }
    return m
}
// NewItemCertificateBasedAuthConfigurationCountRequestBuilder instantiates a new ItemCertificateBasedAuthConfigurationCountRequestBuilder and sets the default values.
func NewItemCertificateBasedAuthConfigurationCountRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemCertificateBasedAuthConfigurationCountRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemCertificateBasedAuthConfigurationCountRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get the number of the resource
// returns a *int32 when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemCertificateBasedAuthConfigurationCountRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemCertificateBasedAuthConfigurationCountRequestBuilderGetRequestConfiguration)(*int32, error) {
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
func (m *ItemCertificateBasedAuthConfigurationCountRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemCertificateBasedAuthConfigurationCountRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemCertificateBasedAuthConfigurationCountRequestBuilder when successful
func (m *ItemCertificateBasedAuthConfigurationCountRequestBuilder) WithUrl(rawUrl string)(*ItemCertificateBasedAuthConfigurationCountRequestBuilder) {
    return NewItemCertificateBasedAuthConfigurationCountRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
