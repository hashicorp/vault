package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemClaimsMappingPoliciesRequestBuilder provides operations to manage the claimsMappingPolicies property of the microsoft.graph.servicePrincipal entity.
type ItemClaimsMappingPoliciesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemClaimsMappingPoliciesRequestBuilderGetQueryParameters list the claimsMappingPolicy objects that are assigned to a servicePrincipal.
type ItemClaimsMappingPoliciesRequestBuilderGetQueryParameters struct {
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
// ItemClaimsMappingPoliciesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemClaimsMappingPoliciesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemClaimsMappingPoliciesRequestBuilderGetQueryParameters
}
// ByClaimsMappingPolicyId gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.servicePrincipals.item.claimsMappingPolicies.item collection
// returns a *ItemClaimsMappingPoliciesClaimsMappingPolicyItemRequestBuilder when successful
func (m *ItemClaimsMappingPoliciesRequestBuilder) ByClaimsMappingPolicyId(claimsMappingPolicyId string)(*ItemClaimsMappingPoliciesClaimsMappingPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if claimsMappingPolicyId != "" {
        urlTplParams["claimsMappingPolicy%2Did"] = claimsMappingPolicyId
    }
    return NewItemClaimsMappingPoliciesClaimsMappingPolicyItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemClaimsMappingPoliciesRequestBuilderInternal instantiates a new ItemClaimsMappingPoliciesRequestBuilder and sets the default values.
func NewItemClaimsMappingPoliciesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemClaimsMappingPoliciesRequestBuilder) {
    m := &ItemClaimsMappingPoliciesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/claimsMappingPolicies{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemClaimsMappingPoliciesRequestBuilder instantiates a new ItemClaimsMappingPoliciesRequestBuilder and sets the default values.
func NewItemClaimsMappingPoliciesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemClaimsMappingPoliciesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemClaimsMappingPoliciesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *ItemClaimsMappingPoliciesCountRequestBuilder when successful
func (m *ItemClaimsMappingPoliciesRequestBuilder) Count()(*ItemClaimsMappingPoliciesCountRequestBuilder) {
    return NewItemClaimsMappingPoliciesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get list the claimsMappingPolicy objects that are assigned to a servicePrincipal.
// returns a ClaimsMappingPolicyCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceprincipal-list-claimsmappingpolicies?view=graph-rest-1.0
func (m *ItemClaimsMappingPoliciesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemClaimsMappingPoliciesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ClaimsMappingPolicyCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateClaimsMappingPolicyCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ClaimsMappingPolicyCollectionResponseable), nil
}
// Ref provides operations to manage the collection of servicePrincipal entities.
// returns a *ItemClaimsMappingPoliciesRefRequestBuilder when successful
func (m *ItemClaimsMappingPoliciesRequestBuilder) Ref()(*ItemClaimsMappingPoliciesRefRequestBuilder) {
    return NewItemClaimsMappingPoliciesRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation list the claimsMappingPolicy objects that are assigned to a servicePrincipal.
// returns a *RequestInformation when successful
func (m *ItemClaimsMappingPoliciesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemClaimsMappingPoliciesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemClaimsMappingPoliciesRequestBuilder when successful
func (m *ItemClaimsMappingPoliciesRequestBuilder) WithUrl(rawUrl string)(*ItemClaimsMappingPoliciesRequestBuilder) {
    return NewItemClaimsMappingPoliciesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
