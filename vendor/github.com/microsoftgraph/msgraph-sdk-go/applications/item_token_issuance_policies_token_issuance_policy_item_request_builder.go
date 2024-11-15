package applications

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder builds and executes requests for operations under \applications\{application-id}\tokenIssuancePolicies\{tokenIssuancePolicy-id}
type ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilderInternal instantiates a new ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder and sets the default values.
func NewItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder) {
    m := &ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/tokenIssuancePolicies/{tokenIssuancePolicy%2Did}", pathParameters),
    }
    return m
}
// NewItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder instantiates a new ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder and sets the default values.
func NewItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of application entities.
// returns a *ItemTokenIssuancePoliciesItemRefRequestBuilder when successful
func (m *ItemTokenIssuancePoliciesTokenIssuancePolicyItemRequestBuilder) Ref()(*ItemTokenIssuancePoliciesItemRefRequestBuilder) {
    return NewItemTokenIssuancePoliciesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
