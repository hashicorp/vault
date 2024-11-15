package identity

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder builds and executes requests for operations under \identity\authenticationEventsFlows\{authenticationEventsFlow-id}\conditions\applications
type AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewAuthenticationEventsFlowsItemConditionsApplicationsRequestBuilderInternal instantiates a new AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder and sets the default values.
func NewAuthenticationEventsFlowsItemConditionsApplicationsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder) {
    m := &AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/authenticationEventsFlows/{authenticationEventsFlow%2Did}/conditions/applications", pathParameters),
    }
    return m
}
// NewAuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder instantiates a new AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder and sets the default values.
func NewAuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAuthenticationEventsFlowsItemConditionsApplicationsRequestBuilderInternal(urlParams, requestAdapter)
}
// IncludeApplications provides operations to manage the includeApplications property of the microsoft.graph.authenticationConditionsApplications entity.
// returns a *AuthenticationEventsFlowsItemConditionsApplicationsIncludeApplicationsRequestBuilder when successful
func (m *AuthenticationEventsFlowsItemConditionsApplicationsRequestBuilder) IncludeApplications()(*AuthenticationEventsFlowsItemConditionsApplicationsIncludeApplicationsRequestBuilder) {
    return NewAuthenticationEventsFlowsItemConditionsApplicationsIncludeApplicationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
