package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemHomeRealmDiscoveryPoliciesRequestBuilder provides operations to manage the homeRealmDiscoveryPolicies property of the microsoft.graph.servicePrincipal entity.
type ItemHomeRealmDiscoveryPoliciesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemHomeRealmDiscoveryPoliciesRequestBuilderGetQueryParameters list the homeRealmDiscoveryPolicy objects that are assigned to a servicePrincipal.
type ItemHomeRealmDiscoveryPoliciesRequestBuilderGetQueryParameters struct {
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
// ItemHomeRealmDiscoveryPoliciesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemHomeRealmDiscoveryPoliciesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemHomeRealmDiscoveryPoliciesRequestBuilderGetQueryParameters
}
// ByHomeRealmDiscoveryPolicyId gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.servicePrincipals.item.homeRealmDiscoveryPolicies.item collection
// returns a *ItemHomeRealmDiscoveryPoliciesHomeRealmDiscoveryPolicyItemRequestBuilder when successful
func (m *ItemHomeRealmDiscoveryPoliciesRequestBuilder) ByHomeRealmDiscoveryPolicyId(homeRealmDiscoveryPolicyId string)(*ItemHomeRealmDiscoveryPoliciesHomeRealmDiscoveryPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if homeRealmDiscoveryPolicyId != "" {
        urlTplParams["homeRealmDiscoveryPolicy%2Did"] = homeRealmDiscoveryPolicyId
    }
    return NewItemHomeRealmDiscoveryPoliciesHomeRealmDiscoveryPolicyItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemHomeRealmDiscoveryPoliciesRequestBuilderInternal instantiates a new ItemHomeRealmDiscoveryPoliciesRequestBuilder and sets the default values.
func NewItemHomeRealmDiscoveryPoliciesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemHomeRealmDiscoveryPoliciesRequestBuilder) {
    m := &ItemHomeRealmDiscoveryPoliciesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/homeRealmDiscoveryPolicies{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemHomeRealmDiscoveryPoliciesRequestBuilder instantiates a new ItemHomeRealmDiscoveryPoliciesRequestBuilder and sets the default values.
func NewItemHomeRealmDiscoveryPoliciesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemHomeRealmDiscoveryPoliciesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemHomeRealmDiscoveryPoliciesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *ItemHomeRealmDiscoveryPoliciesCountRequestBuilder when successful
func (m *ItemHomeRealmDiscoveryPoliciesRequestBuilder) Count()(*ItemHomeRealmDiscoveryPoliciesCountRequestBuilder) {
    return NewItemHomeRealmDiscoveryPoliciesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get list the homeRealmDiscoveryPolicy objects that are assigned to a servicePrincipal.
// returns a HomeRealmDiscoveryPolicyCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceprincipal-list-homerealmdiscoverypolicies?view=graph-rest-1.0
func (m *ItemHomeRealmDiscoveryPoliciesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemHomeRealmDiscoveryPoliciesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.HomeRealmDiscoveryPolicyCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateHomeRealmDiscoveryPolicyCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.HomeRealmDiscoveryPolicyCollectionResponseable), nil
}
// Ref provides operations to manage the collection of servicePrincipal entities.
// returns a *ItemHomeRealmDiscoveryPoliciesRefRequestBuilder when successful
func (m *ItemHomeRealmDiscoveryPoliciesRequestBuilder) Ref()(*ItemHomeRealmDiscoveryPoliciesRefRequestBuilder) {
    return NewItemHomeRealmDiscoveryPoliciesRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation list the homeRealmDiscoveryPolicy objects that are assigned to a servicePrincipal.
// returns a *RequestInformation when successful
func (m *ItemHomeRealmDiscoveryPoliciesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemHomeRealmDiscoveryPoliciesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemHomeRealmDiscoveryPoliciesRequestBuilder when successful
func (m *ItemHomeRealmDiscoveryPoliciesRequestBuilder) WithUrl(rawUrl string)(*ItemHomeRealmDiscoveryPoliciesRequestBuilder) {
    return NewItemHomeRealmDiscoveryPoliciesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
