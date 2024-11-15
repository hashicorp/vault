package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder provides operations to call the usersRegisteredByFeature method.
type AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AuthenticationMethodsUsersRegisteredByFeatureRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuthenticationMethodsUsersRegisteredByFeatureRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAuthenticationMethodsUsersRegisteredByFeatureRequestBuilderInternal instantiates a new AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder and sets the default values.
func NewAuthenticationMethodsUsersRegisteredByFeatureRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder) {
    m := &AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/authenticationMethods/usersRegisteredByFeature()", pathParameters),
    }
    return m
}
// NewAuthenticationMethodsUsersRegisteredByFeatureRequestBuilder instantiates a new AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder and sets the default values.
func NewAuthenticationMethodsUsersRegisteredByFeatureRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAuthenticationMethodsUsersRegisteredByFeatureRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get the number of users capable of multi-factor authentication, self-service password reset, and passwordless authentication.
// returns a UserRegistrationFeatureSummaryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/authenticationmethodsroot-usersregisteredbyfeature?view=graph-rest-1.0
func (m *AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder) Get(ctx context.Context, requestConfiguration *AuthenticationMethodsUsersRegisteredByFeatureRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserRegistrationFeatureSummaryable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserRegistrationFeatureSummaryFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserRegistrationFeatureSummaryable), nil
}
// ToGetRequestInformation get the number of users capable of multi-factor authentication, self-service password reset, and passwordless authentication.
// returns a *RequestInformation when successful
func (m *AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AuthenticationMethodsUsersRegisteredByFeatureRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder when successful
func (m *AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder) WithUrl(rawUrl string)(*AuthenticationMethodsUsersRegisteredByFeatureRequestBuilder) {
    return NewAuthenticationMethodsUsersRegisteredByFeatureRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
