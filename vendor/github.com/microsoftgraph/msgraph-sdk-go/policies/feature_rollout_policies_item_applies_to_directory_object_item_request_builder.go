package policies

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder builds and executes requests for operations under \policies\featureRolloutPolicies\{featureRolloutPolicy-id}\appliesTo\{directoryObject-id}
type FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewFeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilderInternal instantiates a new FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder and sets the default values.
func NewFeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) {
    m := &FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/featureRolloutPolicies/{featureRolloutPolicy%2Did}/appliesTo/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewFeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder instantiates a new FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder and sets the default values.
func NewFeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewFeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of policyRoot entities.
// returns a *FeatureRolloutPoliciesItemAppliesToItemRefRequestBuilder when successful
func (m *FeatureRolloutPoliciesItemAppliesToDirectoryObjectItemRequestBuilder) Ref()(*FeatureRolloutPoliciesItemAppliesToItemRefRequestBuilder) {
    return NewFeatureRolloutPoliciesItemAppliesToItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
