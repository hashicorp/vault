package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder provides operations to call the usersRegisteredByMethod method.
type AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderInternal instantiates a new AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder and sets the default values.
func NewAuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, includedUserRoles *string, includedUserTypes *string)(*AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder) {
    m := &AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/authenticationMethods/usersRegisteredByMethod(includedUserTypes='{includedUserTypes}',includedUserRoles='{includedUserRoles}')", pathParameters),
    }
    if includedUserRoles != nil {
        m.BaseRequestBuilder.PathParameters["includedUserRoles"] = *includedUserRoles
    }
    if includedUserTypes != nil {
        m.BaseRequestBuilder.PathParameters["includedUserTypes"] = *includedUserTypes
    }
    return m
}
// NewAuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder instantiates a new AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder and sets the default values.
func NewAuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderInternal(urlParams, requestAdapter, nil, nil)
}
// Get invoke function usersRegisteredByMethod
// returns a UserRegistrationMethodSummaryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder) Get(ctx context.Context, requestConfiguration *AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserRegistrationMethodSummaryable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserRegistrationMethodSummaryFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserRegistrationMethodSummaryable), nil
}
// ToGetRequestInformation invoke function usersRegisteredByMethod
// returns a *RequestInformation when successful
func (m *AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder when successful
func (m *AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder) WithUrl(rawUrl string)(*AuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder) {
    return NewAuthenticationMethodsUsersRegisteredByMethodWithIncludedUserTypesWithIncludedUserRolesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
