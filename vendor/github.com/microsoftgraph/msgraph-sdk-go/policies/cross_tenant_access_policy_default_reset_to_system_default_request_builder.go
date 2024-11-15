package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder provides operations to call the resetToSystemDefault method.
type CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderInternal instantiates a new CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder) {
    m := &CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/crossTenantAccessPolicy/default/resetToSystemDefault", pathParameters),
    }
    return m
}
// NewCrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder instantiates a new CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder and sets the default values.
func NewCrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderInternal(urlParams, requestAdapter)
}
// Post reset any changes made to the default configuration in a cross-tenant access policy back to the system default.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/crosstenantaccesspolicyconfigurationdefault-resettosystemdefault?view=graph-rest-1.0
func (m *CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder) Post(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation reset any changes made to the default configuration in a cross-tenant access policy back to the system default.
// returns a *RequestInformation when successful
func (m *CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder when successful
func (m *CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder) WithUrl(rawUrl string)(*CrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder) {
    return NewCrossTenantAccessPolicyDefaultResetToSystemDefaultRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
