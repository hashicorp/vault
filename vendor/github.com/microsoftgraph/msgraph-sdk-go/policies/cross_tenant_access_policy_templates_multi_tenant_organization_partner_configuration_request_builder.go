package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder provides operations to manage the multiTenantOrganizationPartnerConfiguration property of the microsoft.graph.policyTemplate entity.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetQueryParameters get the cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetQueryParameters
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderInternal instantiates a new CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) {
    m := &CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/crossTenantAccessPolicy/templates/multiTenantOrganizationPartnerConfiguration{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder instantiates a new CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property multiTenantOrganizationPartnerConfiguration for policies
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) Delete(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get the cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
// returns a MultiTenantOrganizationPartnerConfigurationTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/multitenantorganizationpartnerconfigurationtemplate-get?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) Get(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationPartnerConfigurationTemplateable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMultiTenantOrganizationPartnerConfigurationTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationPartnerConfigurationTemplateable), nil
}
// Patch update the cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
// returns a MultiTenantOrganizationPartnerConfigurationTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/multitenantorganizationpartnerconfigurationtemplate-update?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationPartnerConfigurationTemplateable, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationPartnerConfigurationTemplateable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMultiTenantOrganizationPartnerConfigurationTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationPartnerConfigurationTemplateable), nil
}
// ToDeleteRequestInformation delete navigation property multiTenantOrganizationPartnerConfiguration for policies
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get the cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationPartnerConfigurationTemplateable, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) WithUrl(rawUrl string)(*CrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder) {
    return NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationPartnerConfigurationRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
