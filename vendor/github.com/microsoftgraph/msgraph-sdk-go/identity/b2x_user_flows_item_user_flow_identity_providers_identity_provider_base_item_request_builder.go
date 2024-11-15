package identity

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder builds and executes requests for operations under \identity\b2xUserFlows\{b2xIdentityUserFlow-id}\userFlowIdentityProviders\{identityProviderBase-id}
type B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewB2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilderInternal instantiates a new B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder and sets the default values.
func NewB2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder) {
    m := &B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/b2xUserFlows/{b2xIdentityUserFlow%2Did}/userFlowIdentityProviders/{identityProviderBase%2Did}", pathParameters),
    }
    return m
}
// NewB2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder instantiates a new B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder and sets the default values.
func NewB2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewB2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of identityContainer entities.
// returns a *B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder when successful
func (m *B2xUserFlowsItemUserFlowIdentityProvidersIdentityProviderBaseItemRequestBuilder) Ref()(*B2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilder) {
    return NewB2xUserFlowsItemUserFlowIdentityProvidersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
