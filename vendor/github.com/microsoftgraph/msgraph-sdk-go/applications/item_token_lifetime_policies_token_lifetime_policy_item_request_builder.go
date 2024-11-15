package applications

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder builds and executes requests for operations under \applications\{application-id}\tokenLifetimePolicies\{tokenLifetimePolicy-id}
type ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilderInternal instantiates a new ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder and sets the default values.
func NewItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder) {
    m := &ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/tokenLifetimePolicies/{tokenLifetimePolicy%2Did}", pathParameters),
    }
    return m
}
// NewItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder instantiates a new ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder and sets the default values.
func NewItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of application entities.
// returns a *ItemTokenLifetimePoliciesItemRefRequestBuilder when successful
func (m *ItemTokenLifetimePoliciesTokenLifetimePolicyItemRequestBuilder) Ref()(*ItemTokenLifetimePoliciesItemRefRequestBuilder) {
    return NewItemTokenLifetimePoliciesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
