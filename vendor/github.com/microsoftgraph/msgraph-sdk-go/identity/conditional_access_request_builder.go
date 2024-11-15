package identity

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ConditionalAccessRequestBuilder builds and executes requests for operations under \identity\conditionalAccess
type ConditionalAccessRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AuthenticationContextClassReferences provides operations to manage the authenticationContextClassReferences property of the microsoft.graph.conditionalAccessRoot entity.
// returns a *ConditionalAccessAuthenticationContextClassReferencesRequestBuilder when successful
func (m *ConditionalAccessRequestBuilder) AuthenticationContextClassReferences()(*ConditionalAccessAuthenticationContextClassReferencesRequestBuilder) {
    return NewConditionalAccessAuthenticationContextClassReferencesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthenticationStrength provides operations to manage the authenticationStrength property of the microsoft.graph.conditionalAccessRoot entity.
// returns a *ConditionalAccessAuthenticationStrengthRequestBuilder when successful
func (m *ConditionalAccessRequestBuilder) AuthenticationStrength()(*ConditionalAccessAuthenticationStrengthRequestBuilder) {
    return NewConditionalAccessAuthenticationStrengthRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewConditionalAccessRequestBuilderInternal instantiates a new ConditionalAccessRequestBuilder and sets the default values.
func NewConditionalAccessRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConditionalAccessRequestBuilder) {
    m := &ConditionalAccessRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/conditionalAccess", pathParameters),
    }
    return m
}
// NewConditionalAccessRequestBuilder instantiates a new ConditionalAccessRequestBuilder and sets the default values.
func NewConditionalAccessRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConditionalAccessRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewConditionalAccessRequestBuilderInternal(urlParams, requestAdapter)
}
// NamedLocations provides operations to manage the namedLocations property of the microsoft.graph.conditionalAccessRoot entity.
// returns a *ConditionalAccessNamedLocationsRequestBuilder when successful
func (m *ConditionalAccessRequestBuilder) NamedLocations()(*ConditionalAccessNamedLocationsRequestBuilder) {
    return NewConditionalAccessNamedLocationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Policies provides operations to manage the policies property of the microsoft.graph.conditionalAccessRoot entity.
// returns a *ConditionalAccessPoliciesRequestBuilder when successful
func (m *ConditionalAccessRequestBuilder) Policies()(*ConditionalAccessPoliciesRequestBuilder) {
    return NewConditionalAccessPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Templates provides operations to manage the templates property of the microsoft.graph.conditionalAccessRoot entity.
// returns a *ConditionalAccessTemplatesRequestBuilder when successful
func (m *ConditionalAccessRequestBuilder) Templates()(*ConditionalAccessTemplatesRequestBuilder) {
    return NewConditionalAccessTemplatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
