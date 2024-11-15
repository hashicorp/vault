package applications

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder builds and executes requests for operations under \applications\{application-id}\appManagementPolicies\{appManagementPolicy-id}
type ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal instantiates a new ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder and sets the default values.
func NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    m := &ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/appManagementPolicies/{appManagementPolicy%2Did}", pathParameters),
    }
    return m
}
// NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder instantiates a new ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder and sets the default values.
func NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of application entities.
// returns a *ItemAppManagementPoliciesItemRefRequestBuilder when successful
func (m *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) Ref()(*ItemAppManagementPoliciesItemRefRequestBuilder) {
    return NewItemAppManagementPoliciesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
