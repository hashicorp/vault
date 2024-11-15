package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder provides operations to manage the collection of identityContainer entities.
type B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderInternal instantiates a new B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder and sets the default values.
func NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) {
    m := &B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/b2xUserFlows/{b2xIdentityUserFlow%2Did}/userFlowIdentityProviders/{identityProviderBase%2Did}/$ref", pathParameters),
    }
    return m
}
// NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder instantiates a new B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder and sets the default values.
func NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete ref of navigation property userFlowIdentityProviders for identity
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) Delete(ctx context.Context, requestConfiguration *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderDeleteRequestConfiguration)(error) {
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
// ToDeleteRequestInformation delete ref of navigation property userFlowIdentityProviders for identity
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder when successful
func (m *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) WithUrl(rawUrl string)(*B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) {
    return NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
