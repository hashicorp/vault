package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder provides operations to manage the partners property of the microsoft.graph.crossTenantAccessPolicy entity.
type CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetQueryParameters read the properties and relationships of a partner-specific configuration.
type CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetQueryParameters
}
// CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderInternal instantiates a new CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) {
    m := &CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/crossTenantAccessPolicy/partners/{crossTenantAccessPolicyConfigurationPartner%2DtenantId}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewCrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder instantiates a new CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete a partner-specific configuration in a cross-tenant access policy. If a configuration includes a user synchronization policy, you must first delete the user synchronization policy before you can delete the partner-specific configuration.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/crosstenantaccesspolicyconfigurationpartner-delete?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get read the properties and relationships of a partner-specific configuration.
// returns a CrossTenantAccessPolicyConfigurationPartnerable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/crosstenantaccesspolicyconfigurationpartner-get?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) Get(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CrossTenantAccessPolicyConfigurationPartnerable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCrossTenantAccessPolicyConfigurationPartnerFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CrossTenantAccessPolicyConfigurationPartnerable), nil
}
// IdentitySynchronization provides operations to manage the identitySynchronization property of the microsoft.graph.crossTenantAccessPolicyConfigurationPartner entity.
// returns a *CrossTenantAccessPolicyPartnersItemIdentitySynchronizationRequestBuilder when successful
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) IdentitySynchronization()(*CrossTenantAccessPolicyPartnersItemIdentitySynchronizationRequestBuilder) {
    return NewCrossTenantAccessPolicyPartnersItemIdentitySynchronizationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the properties of a partner-specific configuration.
// returns a CrossTenantAccessPolicyConfigurationPartnerable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/crosstenantaccesspolicyconfigurationpartner-update?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CrossTenantAccessPolicyConfigurationPartnerable, requestConfiguration *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CrossTenantAccessPolicyConfigurationPartnerable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCrossTenantAccessPolicyConfigurationPartnerFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CrossTenantAccessPolicyConfigurationPartnerable), nil
}
// ToDeleteRequestInformation delete a partner-specific configuration in a cross-tenant access policy. If a configuration includes a user synchronization policy, you must first delete the user synchronization policy before you can delete the partner-specific configuration.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of a partner-specific configuration.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a partner-specific configuration.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CrossTenantAccessPolicyConfigurationPartnerable, requestConfiguration *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder when successful
func (m *CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) WithUrl(rawUrl string)(*CrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder) {
    return NewCrossTenantAccessPolicyPartnersCrossTenantAccessPolicyConfigurationPartnerTenantItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
