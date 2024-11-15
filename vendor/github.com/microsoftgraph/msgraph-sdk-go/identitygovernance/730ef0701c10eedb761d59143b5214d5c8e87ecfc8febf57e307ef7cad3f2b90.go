package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder provides operations to manage the scopes property of the microsoft.graph.accessPackageResource entity.
type EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetQueryParameters read-only. Nullable. Supports $expand.
type EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetQueryParameters
}
// EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderInternal instantiates a new EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder and sets the default values.
func NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) {
    m := &EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/entitlementManagement/resourceRequests/{accessPackageResourceRequest%2Did}/catalog/resources/{accessPackageResource%2Did}/roles/{accessPackageResourceRole%2Did}/resource/scopes/{accessPackageResourceScope%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder instantiates a new EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder and sets the default values.
func NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property scopes for identityGovernance
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read-only. Nullable. Supports $expand.
// returns a AccessPackageResourceScopeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceScopeable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageResourceScopeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceScopeable), nil
}
// Patch update the navigation property scopes in identityGovernance
// returns a AccessPackageResourceScopeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceScopeable, requestConfiguration *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceScopeable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageResourceScopeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceScopeable), nil
}
// Resource provides operations to manage the resource property of the microsoft.graph.accessPackageResourceScope entity.
// returns a *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesItemResourceRequestBuilder when successful
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) Resource()(*EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesItemResourceRequestBuilder) {
    return NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesItemResourceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property scopes for identityGovernance
// returns a *RequestInformation when successful
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read-only. Nullable. Supports $expand.
// returns a *RequestInformation when successful
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property scopes in identityGovernance
// returns a *RequestInformation when successful
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageResourceScopeable, requestConfiguration *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder when successful
func (m *EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) WithUrl(rawUrl string)(*EntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder) {
    return NewEntitlementManagementResourceRequestsItemCatalogResourcesItemRolesItemResourceScopesAccessPackageResourceScopeItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
