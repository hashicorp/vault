package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder provides operations to manage the multiTenantOrganizationIdentitySynchronization property of the microsoft.graph.policyTemplate entity.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetQueryParameters get the cross-tenant access policy template with user synchronization settings for a multitenant organization.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetQueryParameters
}
// CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderInternal instantiates a new CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) {
    m := &CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/crossTenantAccessPolicy/templates/multiTenantOrganizationIdentitySynchronization{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder instantiates a new CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property multiTenantOrganizationIdentitySynchronization for policies
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) Delete(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get the cross-tenant access policy template with user synchronization settings for a multitenant organization.
// returns a MultiTenantOrganizationIdentitySyncPolicyTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/multitenantorganizationidentitysyncpolicytemplate-get?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) Get(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationIdentitySyncPolicyTemplateable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMultiTenantOrganizationIdentitySyncPolicyTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationIdentitySyncPolicyTemplateable), nil
}
// Patch update the cross-tenant access policy template with user synchronization settings for a multitenant organization.
// returns a MultiTenantOrganizationIdentitySyncPolicyTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/multitenantorganizationidentitysyncpolicytemplate-update?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationIdentitySyncPolicyTemplateable, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationIdentitySyncPolicyTemplateable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMultiTenantOrganizationIdentitySyncPolicyTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationIdentitySyncPolicyTemplateable), nil
}
// ToDeleteRequestInformation delete navigation property multiTenantOrganizationIdentitySynchronization for policies
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get the cross-tenant access policy template with user synchronization settings for a multitenant organization.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the cross-tenant access policy template with user synchronization settings for a multitenant organization.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MultiTenantOrganizationIdentitySyncPolicyTemplateable, requestConfiguration *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder when successful
func (m *CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) WithUrl(rawUrl string)(*CrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder) {
    return NewCrossTenantAccessPolicyTemplatesMultiTenantOrganizationIdentitySynchronizationRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
