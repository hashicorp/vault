package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder provides operations to manage the appManagementPolicies property of the microsoft.graph.servicePrincipal entity.
type ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetQueryParameters the appManagementPolicy applied to this application.
type ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetQueryParameters
}
// NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal instantiates a new ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder and sets the default values.
func NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    m := &ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/appManagementPolicies/{appManagementPolicy%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder instantiates a new ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder and sets the default values.
func NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the appManagementPolicy applied to this application.
// returns a AppManagementPolicyable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AppManagementPolicyable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAppManagementPolicyFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AppManagementPolicyable), nil
}
// ToGetRequestInformation the appManagementPolicy applied to this application.
// returns a *RequestInformation when successful
func (m *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder when successful
func (m *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) WithUrl(rawUrl string)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    return NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
