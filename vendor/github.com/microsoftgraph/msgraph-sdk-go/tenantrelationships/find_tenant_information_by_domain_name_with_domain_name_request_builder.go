package tenantrelationships

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// FindTenantInformationByDomainNameWithDomainNameRequestBuilder provides operations to call the findTenantInformationByDomainName method.
type FindTenantInformationByDomainNameWithDomainNameRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// FindTenantInformationByDomainNameWithDomainNameRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type FindTenantInformationByDomainNameWithDomainNameRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewFindTenantInformationByDomainNameWithDomainNameRequestBuilderInternal instantiates a new FindTenantInformationByDomainNameWithDomainNameRequestBuilder and sets the default values.
func NewFindTenantInformationByDomainNameWithDomainNameRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, domainName *string)(*FindTenantInformationByDomainNameWithDomainNameRequestBuilder) {
    m := &FindTenantInformationByDomainNameWithDomainNameRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/tenantRelationships/findTenantInformationByDomainName(domainName='{domainName}')", pathParameters),
    }
    if domainName != nil {
        m.BaseRequestBuilder.PathParameters["domainName"] = *domainName
    }
    return m
}
// NewFindTenantInformationByDomainNameWithDomainNameRequestBuilder instantiates a new FindTenantInformationByDomainNameWithDomainNameRequestBuilder and sets the default values.
func NewFindTenantInformationByDomainNameWithDomainNameRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*FindTenantInformationByDomainNameWithDomainNameRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewFindTenantInformationByDomainNameWithDomainNameRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get given a domain name, search for a tenant and read its tenantInformation. You can use this API to validate tenant information and use the tenantId to configure cross-tenant access settings between you and the tenant.
// returns a TenantInformationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/tenantrelationship-findtenantinformationbydomainname?view=graph-rest-1.0
func (m *FindTenantInformationByDomainNameWithDomainNameRequestBuilder) Get(ctx context.Context, requestConfiguration *FindTenantInformationByDomainNameWithDomainNameRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantInformationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTenantInformationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TenantInformationable), nil
}
// ToGetRequestInformation given a domain name, search for a tenant and read its tenantInformation. You can use this API to validate tenant information and use the tenantId to configure cross-tenant access settings between you and the tenant.
// returns a *RequestInformation when successful
func (m *FindTenantInformationByDomainNameWithDomainNameRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *FindTenantInformationByDomainNameWithDomainNameRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *FindTenantInformationByDomainNameWithDomainNameRequestBuilder when successful
func (m *FindTenantInformationByDomainNameWithDomainNameRequestBuilder) WithUrl(rawUrl string)(*FindTenantInformationByDomainNameWithDomainNameRequestBuilder) {
    return NewFindTenantInformationByDomainNameWithDomainNameRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
