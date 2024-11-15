package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder provides operations to call the updateAllowedCombinations method.
type AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderInternal instantiates a new AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder and sets the default values.
func NewAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder) {
    m := &AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/authenticationStrengthPolicies/{authenticationStrengthPolicy%2Did}/updateAllowedCombinations", pathParameters),
    }
    return m
}
// NewAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder instantiates a new AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder and sets the default values.
func NewAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderInternal(urlParams, requestAdapter)
}
// Post update the allowedCombinations property of an authenticationStrengthPolicy object. To update other properties of an authenticationStrengthPolicy object, use the Update authenticationStrengthPolicy method.
// returns a UpdateAllowedCombinationsResultable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/authenticationstrengthpolicy-updateallowedcombinations?view=graph-rest-1.0
func (m *AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder) Post(ctx context.Context, body AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsPostRequestBodyable, requestConfiguration *AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UpdateAllowedCombinationsResultable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUpdateAllowedCombinationsResultFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UpdateAllowedCombinationsResultable), nil
}
// ToPostRequestInformation update the allowedCombinations property of an authenticationStrengthPolicy object. To update other properties of an authenticationStrengthPolicy object, use the Update authenticationStrengthPolicy method.
// returns a *RequestInformation when successful
func (m *AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder) ToPostRequestInformation(ctx context.Context, body AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsPostRequestBodyable, requestConfiguration *AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder when successful
func (m *AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder) WithUrl(rawUrl string)(*AuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder) {
    return NewAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
