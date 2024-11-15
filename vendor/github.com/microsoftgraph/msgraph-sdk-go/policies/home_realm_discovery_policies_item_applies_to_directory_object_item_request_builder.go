package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder provides operations to manage the appliesTo property of the microsoft.graph.stsPolicy entity.
type HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetQueryParameters get appliesTo from policies
type HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetQueryParameters
}
// NewHomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderInternal instantiates a new HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder and sets the default values.
func NewHomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) {
    m := &HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/homeRealmDiscoveryPolicies/{homeRealmDiscoveryPolicy%2Did}/appliesTo/{directoryObject%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewHomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder instantiates a new HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder and sets the default values.
func NewHomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewHomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get appliesTo from policies
// returns a DirectoryObjectable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) Get(ctx context.Context, requestConfiguration *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable), nil
}
// ToGetRequestInformation get appliesTo from policies
// returns a *RequestInformation when successful
func (m *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder when successful
func (m *HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) WithUrl(rawUrl string)(*HomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) {
    return NewHomeRealmDiscoveryPoliciesItemAppliesToDirectoryObjectItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
