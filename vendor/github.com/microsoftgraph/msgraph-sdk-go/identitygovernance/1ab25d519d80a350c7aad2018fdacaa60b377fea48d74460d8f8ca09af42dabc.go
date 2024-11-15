package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder provides operations to manage the resourceRoleScopes property of the microsoft.graph.accessPackage entity.
type EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetQueryParameters the resource roles and scopes in this access package.
type EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetQueryParameters
}
// EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewEntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderInternal instantiates a new EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder and sets the default values.
func NewEntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) {
    m := &EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/entitlementManagement/accessPackages/{accessPackage%2Did}/resourceRoleScopes/{accessPackageResourceRoleScope%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewEntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder instantiates a new EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder and sets the default values.
func NewEntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete remove a accessPackageResourceRoleScope from an accessPackage list of resource role scopes.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accesspackage-delete-resourcerolescopes?view=graph-rest-1.0
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the resource roles and scopes in this access package.
// returns a AccessPackageResourceRoleScopeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceRoleScopeable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageResourceRoleScopeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceRoleScopeable), nil
}
// Patch update the navigation property resourceRoleScopes in identityGovernance
// returns a AccessPackageResourceRoleScopeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceRoleScopeable, requestConfiguration *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceRoleScopeable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageResourceRoleScopeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceRoleScopeable), nil
}
// Role provides operations to manage the role property of the microsoft.graph.accessPackageResourceRoleScope entity.
// returns a *EntitlementManagementAccessPackagesItemResourceRoleScopesItemRoleRequestBuilder when successful
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) Role()(*EntitlementManagementAccessPackagesItemResourceRoleScopesItemRoleRequestBuilder) {
    return NewEntitlementManagementAccessPackagesItemResourceRoleScopesItemRoleRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Scope provides operations to manage the scope property of the microsoft.graph.accessPackageResourceRoleScope entity.
// returns a *EntitlementManagementAccessPackagesItemResourceRoleScopesItemScopeRequestBuilder when successful
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) Scope()(*EntitlementManagementAccessPackagesItemResourceRoleScopesItemScopeRequestBuilder) {
    return NewEntitlementManagementAccessPackagesItemResourceRoleScopesItemScopeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation remove a accessPackageResourceRoleScope from an accessPackage list of resource role scopes.
// returns a *RequestInformation when successful
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the resource roles and scopes in this access package.
// returns a *RequestInformation when successful
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property resourceRoleScopes in identityGovernance
// returns a *RequestInformation when successful
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceRoleScopeable, requestConfiguration *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder when successful
func (m *EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) WithUrl(rawUrl string)(*EntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder) {
    return NewEntitlementManagementAccessPackagesItemResourceRoleScopesAccessPackageResourceRoleScopeItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
