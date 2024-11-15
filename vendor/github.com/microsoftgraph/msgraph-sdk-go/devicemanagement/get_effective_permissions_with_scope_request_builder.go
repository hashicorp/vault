package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// GetEffectivePermissionsWithScopeRequestBuilder provides operations to call the getEffectivePermissions method.
type GetEffectivePermissionsWithScopeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// GetEffectivePermissionsWithScopeRequestBuilderGetQueryParameters retrieves the effective permissions of the currently authenticated user
type GetEffectivePermissionsWithScopeRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// GetEffectivePermissionsWithScopeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type GetEffectivePermissionsWithScopeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *GetEffectivePermissionsWithScopeRequestBuilderGetQueryParameters
}
// NewGetEffectivePermissionsWithScopeRequestBuilderInternal instantiates a new GetEffectivePermissionsWithScopeRequestBuilder and sets the default values.
func NewGetEffectivePermissionsWithScopeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, scope *string)(*GetEffectivePermissionsWithScopeRequestBuilder) {
    m := &GetEffectivePermissionsWithScopeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/getEffectivePermissions(scope='{scope}'){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    if scope != nil {
        m.BaseRequestBuilder.PathParameters["scope"] = *scope
    }
    return m
}
// NewGetEffectivePermissionsWithScopeRequestBuilder instantiates a new GetEffectivePermissionsWithScopeRequestBuilder and sets the default values.
func NewGetEffectivePermissionsWithScopeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*GetEffectivePermissionsWithScopeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewGetEffectivePermissionsWithScopeRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get retrieves the effective permissions of the currently authenticated user
// Deprecated: This method is obsolete. Use GetAsGetEffectivePermissionsWithScopeGetResponse instead.
// returns a GetEffectivePermissionsWithScopeResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-rbac-devicemanagement-geteffectivepermissions?view=graph-rest-1.0
func (m *GetEffectivePermissionsWithScopeRequestBuilder) Get(ctx context.Context, requestConfiguration *GetEffectivePermissionsWithScopeRequestBuilderGetRequestConfiguration)(GetEffectivePermissionsWithScopeResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateGetEffectivePermissionsWithScopeResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(GetEffectivePermissionsWithScopeResponseable), nil
}
// GetAsGetEffectivePermissionsWithScopeGetResponse retrieves the effective permissions of the currently authenticated user
// returns a GetEffectivePermissionsWithScopeGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-rbac-devicemanagement-geteffectivepermissions?view=graph-rest-1.0
func (m *GetEffectivePermissionsWithScopeRequestBuilder) GetAsGetEffectivePermissionsWithScopeGetResponse(ctx context.Context, requestConfiguration *GetEffectivePermissionsWithScopeRequestBuilderGetRequestConfiguration)(GetEffectivePermissionsWithScopeGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateGetEffectivePermissionsWithScopeGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(GetEffectivePermissionsWithScopeGetResponseable), nil
}
// ToGetRequestInformation retrieves the effective permissions of the currently authenticated user
// returns a *RequestInformation when successful
func (m *GetEffectivePermissionsWithScopeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *GetEffectivePermissionsWithScopeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *GetEffectivePermissionsWithScopeRequestBuilder when successful
func (m *GetEffectivePermissionsWithScopeRequestBuilder) WithUrl(rawUrl string)(*GetEffectivePermissionsWithScopeRequestBuilder) {
    return NewGetEffectivePermissionsWithScopeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
