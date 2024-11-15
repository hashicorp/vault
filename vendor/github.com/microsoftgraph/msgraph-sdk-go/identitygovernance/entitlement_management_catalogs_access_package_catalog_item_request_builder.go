package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder provides operations to manage the catalogs property of the microsoft.graph.entitlementManagement entity.
type EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetQueryParameters retrieve the properties and relationships of an accessPackageCatalog object.
type EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetQueryParameters
}
// EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AccessPackages provides operations to manage the accessPackages property of the microsoft.graph.accessPackageCatalog entity.
// returns a *EntitlementManagementCatalogsItemAccessPackagesRequestBuilder when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) AccessPackages()(*EntitlementManagementCatalogsItemAccessPackagesRequestBuilder) {
    return NewEntitlementManagementCatalogsItemAccessPackagesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewEntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderInternal instantiates a new EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder and sets the default values.
func NewEntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) {
    m := &EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/entitlementManagement/catalogs/{accessPackageCatalog%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewEntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder instantiates a new EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder and sets the default values.
func NewEntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CustomWorkflowExtensions provides operations to manage the customWorkflowExtensions property of the microsoft.graph.accessPackageCatalog entity.
// returns a *EntitlementManagementCatalogsItemCustomWorkflowExtensionsRequestBuilder when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) CustomWorkflowExtensions()(*EntitlementManagementCatalogsItemCustomWorkflowExtensionsRequestBuilder) {
    return NewEntitlementManagementCatalogsItemCustomWorkflowExtensionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delete delete an accessPackageCatalog.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accesspackagecatalog-delete?view=graph-rest-1.0
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve the properties and relationships of an accessPackageCatalog object.
// returns a AccessPackageCatalogable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accesspackagecatalog-get?view=graph-rest-1.0
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageCatalogable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageCatalogFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageCatalogable), nil
}
// Patch update an existing accessPackageCatalog object to change one or more of its properties, such as the display name or description.
// returns a AccessPackageCatalogable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accesspackagecatalog-update?view=graph-rest-1.0
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageCatalogable, requestConfiguration *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageCatalogable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageCatalogFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageCatalogable), nil
}
// ResourceRoles provides operations to manage the resourceRoles property of the microsoft.graph.accessPackageCatalog entity.
// returns a *EntitlementManagementCatalogsItemResourceRolesRequestBuilder when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) ResourceRoles()(*EntitlementManagementCatalogsItemResourceRolesRequestBuilder) {
    return NewEntitlementManagementCatalogsItemResourceRolesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Resources provides operations to manage the resources property of the microsoft.graph.accessPackageCatalog entity.
// returns a *EntitlementManagementCatalogsItemResourcesRequestBuilder when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) Resources()(*EntitlementManagementCatalogsItemResourcesRequestBuilder) {
    return NewEntitlementManagementCatalogsItemResourcesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ResourceScopes provides operations to manage the resourceScopes property of the microsoft.graph.accessPackageCatalog entity.
// returns a *EntitlementManagementCatalogsItemResourceScopesRequestBuilder when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) ResourceScopes()(*EntitlementManagementCatalogsItemResourceScopesRequestBuilder) {
    return NewEntitlementManagementCatalogsItemResourceScopesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete an accessPackageCatalog.
// returns a *RequestInformation when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation retrieve the properties and relationships of an accessPackageCatalog object.
// returns a *RequestInformation when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update an existing accessPackageCatalog object to change one or more of its properties, such as the display name or description.
// returns a *RequestInformation when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageCatalogable, requestConfiguration *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder when successful
func (m *EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) WithUrl(rawUrl string)(*EntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder) {
    return NewEntitlementManagementCatalogsAccessPackageCatalogItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
