package tenantrelationships

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// TenantRelationshipsRequestBuilder provides operations to manage the tenantRelationship singleton.
type TenantRelationshipsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// TenantRelationshipsRequestBuilderGetQueryParameters get tenantRelationships
type TenantRelationshipsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TenantRelationshipsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TenantRelationshipsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TenantRelationshipsRequestBuilderGetQueryParameters
}
// TenantRelationshipsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TenantRelationshipsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewTenantRelationshipsRequestBuilderInternal instantiates a new TenantRelationshipsRequestBuilder and sets the default values.
func NewTenantRelationshipsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TenantRelationshipsRequestBuilder) {
    m := &TenantRelationshipsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/tenantRelationships{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewTenantRelationshipsRequestBuilder instantiates a new TenantRelationshipsRequestBuilder and sets the default values.
func NewTenantRelationshipsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TenantRelationshipsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTenantRelationshipsRequestBuilderInternal(urlParams, requestAdapter)
}
// DelegatedAdminCustomers provides operations to manage the delegatedAdminCustomers property of the microsoft.graph.tenantRelationship entity.
// returns a *DelegatedAdminCustomersRequestBuilder when successful
func (m *TenantRelationshipsRequestBuilder) DelegatedAdminCustomers()(*DelegatedAdminCustomersRequestBuilder) {
    return NewDelegatedAdminCustomersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DelegatedAdminRelationships provides operations to manage the delegatedAdminRelationships property of the microsoft.graph.tenantRelationship entity.
// returns a *DelegatedAdminRelationshipsRequestBuilder when successful
func (m *TenantRelationshipsRequestBuilder) DelegatedAdminRelationships()(*DelegatedAdminRelationshipsRequestBuilder) {
    return NewDelegatedAdminRelationshipsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FindTenantInformationByDomainNameWithDomainName provides operations to call the findTenantInformationByDomainName method.
// returns a *FindTenantInformationByDomainNameWithDomainNameRequestBuilder when successful
func (m *TenantRelationshipsRequestBuilder) FindTenantInformationByDomainNameWithDomainName(domainName *string)(*FindTenantInformationByDomainNameWithDomainNameRequestBuilder) {
    return NewFindTenantInformationByDomainNameWithDomainNameRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, domainName)
}
// FindTenantInformationByTenantIdWithTenantId provides operations to call the findTenantInformationByTenantId method.
// returns a *FindTenantInformationByTenantIdWithTenantIdRequestBuilder when successful
func (m *TenantRelationshipsRequestBuilder) FindTenantInformationByTenantIdWithTenantId(tenantId *string)(*FindTenantInformationByTenantIdWithTenantIdRequestBuilder) {
    return NewFindTenantInformationByTenantIdWithTenantIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, tenantId)
}
// Get get tenantRelationships
// returns a TenantRelationshipable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TenantRelationshipsRequestBuilder) Get(ctx context.Context, requestConfiguration *TenantRelationshipsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantRelationshipable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTenantRelationshipFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantRelationshipable), nil
}
// MultiTenantOrganization provides operations to manage the multiTenantOrganization property of the microsoft.graph.tenantRelationship entity.
// returns a *MultiTenantOrganizationRequestBuilder when successful
func (m *TenantRelationshipsRequestBuilder) MultiTenantOrganization()(*MultiTenantOrganizationRequestBuilder) {
    return NewMultiTenantOrganizationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update tenantRelationships
// returns a TenantRelationshipable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TenantRelationshipsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantRelationshipable, requestConfiguration *TenantRelationshipsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantRelationshipable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTenantRelationshipFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantRelationshipable), nil
}
// ToGetRequestInformation get tenantRelationships
// returns a *RequestInformation when successful
func (m *TenantRelationshipsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *TenantRelationshipsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update tenantRelationships
// returns a *RequestInformation when successful
func (m *TenantRelationshipsRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantRelationshipable, requestConfiguration *TenantRelationshipsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *TenantRelationshipsRequestBuilder when successful
func (m *TenantRelationshipsRequestBuilder) WithUrl(rawUrl string)(*TenantRelationshipsRequestBuilder) {
    return NewTenantRelationshipsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
