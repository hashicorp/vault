package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder provides operations to manage the collection of identityContainer entities.
type AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteQueryParameters delete ref of navigation property attributes for identity
type AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteQueryParameters struct {
    // The delete Uri
    Id *string `uriparametername:"%40id"`
}
// AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteQueryParameters
}
// AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetQueryParameters get an identityUserFlowAttribute collection associated with an external identities self-service user flow represented by an externalUsersSelfServiceSignupEventsFlow object. These attributes are collected from the user during the authentication experience defined by the user flow.
type AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetQueryParameters
}
// AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderInternal instantiates a new AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder and sets the default values.
func NewAuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) {
    m := &AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/authenticationEventsFlows/{authenticationEventsFlow%2Did}/graph.externalUsersSelfServiceSignUpEventsFlow/onAttributeCollection/graph.onAttributeCollectionExternalUsersSelfServiceSignUp/attributes/$ref?@id={%40id}{&%24count,%24filter,%24orderby,%24search,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewAuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder instantiates a new AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder and sets the default values.
func NewAuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete ref of navigation property attributes for identity
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) Delete(ctx context.Context, requestConfiguration *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get get an identityUserFlowAttribute collection associated with an external identities self-service user flow represented by an externalUsersSelfServiceSignupEventsFlow object. These attributes are collected from the user during the authentication experience defined by the user flow.
// returns a StringCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/onattributecollectionexternalusersselfservicesignup-list-attributes?view=graph-rest-1.0
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) Get(ctx context.Context, requestConfiguration *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.StringCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateStringCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.StringCollectionResponseable), nil
}
// Post add an attribute to an external identities self-service sign up user flow that's represented by an externalUsersSelfServiceSignupEventsFlow object. You can add both custom and built-in attributes to a user flow. The attribute is added to both the attributeCollection> attributes and attributeCollection> attributeCollectionPage > views collections on the user flow. In the views collection, the attribute is assigned the default settings. You can PATCH the user flow to customize the settings of the attribute on the views object, for example, marking it as required or updating the allowed input types.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/onattributecollectionexternalusersselfservicesignup-post-attributes?view=graph-rest-1.0
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReferenceCreateable, requestConfiguration *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// ToDeleteRequestInformation delete ref of navigation property attributes for identity
// returns a *RequestInformation when successful
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, "{+baseurl}/identity/authenticationEventsFlows/{authenticationEventsFlow%2Did}/graph.externalUsersSelfServiceSignUpEventsFlow/onAttributeCollection/graph.onAttributeCollectionExternalUsersSelfServiceSignUp/attributes/$ref?@id={%40id}", m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get an identityUserFlowAttribute collection associated with an external identities self-service user flow represented by an externalUsersSelfServiceSignupEventsFlow object. These attributes are collected from the user during the authentication experience defined by the user flow.
// returns a *RequestInformation when successful
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, "{+baseurl}/identity/authenticationEventsFlows/{authenticationEventsFlow%2Did}/graph.externalUsersSelfServiceSignUpEventsFlow/onAttributeCollection/graph.onAttributeCollectionExternalUsersSelfServiceSignUp/attributes/$ref{?%24count,%24filter,%24orderby,%24search,%24skip,%24top}", m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPostRequestInformation add an attribute to an external identities self-service sign up user flow that's represented by an externalUsersSelfServiceSignupEventsFlow object. You can add both custom and built-in attributes to a user flow. The attribute is added to both the attributeCollection> attributes and attributeCollection> attributeCollectionPage > views collections on the user flow. In the views collection, the attribute is assigned the default settings. You can PATCH the user flow to customize the settings of the attribute on the views object, for example, marking it as required or updating the allowed input types.
// returns a *RequestInformation when successful
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReferenceCreateable, requestConfiguration *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, "{+baseurl}/identity/authenticationEventsFlows/{authenticationEventsFlow%2Did}/graph.externalUsersSelfServiceSignUpEventsFlow/onAttributeCollection/graph.onAttributeCollectionExternalUsersSelfServiceSignUp/attributes/$ref", m.BaseRequestBuilder.PathParameters)
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
// returns a *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder when successful
func (m *AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) WithUrl(rawUrl string)(*AuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder) {
    return NewAuthenticationEventsFlowsItemGraphExternalUsersSelfServiceSignUpEventsFlowOnAttributeCollectionGraphOnAttributeCollectionExternalUsersSelfServiceSignUpAttributesRefRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
