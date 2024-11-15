package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder provides operations to manage the authenticationMethodModes property of the microsoft.graph.authenticationStrengthRoot entity.
type ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetQueryParameters names and descriptions of all valid authentication method modes in the system.
type ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetQueryParameters
}
// ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderInternal instantiates a new ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder and sets the default values.
func NewConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) {
    m := &ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/conditionalAccess/authenticationStrength/authenticationMethodModes/{authenticationMethodModeDetail%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder instantiates a new ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder and sets the default values.
func NewConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property authenticationMethodModes for identity
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get names and descriptions of all valid authentication method modes in the system.
// returns a AuthenticationMethodModeDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationMethodModeDetailable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAuthenticationMethodModeDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationMethodModeDetailable), nil
}
// Patch update the navigation property authenticationMethodModes in identity
// returns a AuthenticationMethodModeDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationMethodModeDetailable, requestConfiguration *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationMethodModeDetailable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAuthenticationMethodModeDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationMethodModeDetailable), nil
}
// ToDeleteRequestInformation delete navigation property authenticationMethodModes for identity
// returns a *RequestInformation when successful
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation names and descriptions of all valid authentication method modes in the system.
// returns a *RequestInformation when successful
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// ToPatchRequestInformation update the navigation property authenticationMethodModes in identity
// returns a *RequestInformation when successful
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationMethodModeDetailable, requestConfiguration *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder when successful
func (m *ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) WithUrl(rawUrl string)(*ConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder) {
    return NewConditionalAccessAuthenticationStrengthAuthenticationMethodModesAuthenticationMethodModeDetailItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
