package applicationtemplates

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemInstantiateRequestBuilder provides operations to call the instantiate method.
type ItemInstantiateRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemInstantiateRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemInstantiateRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemInstantiateRequestBuilderInternal instantiates a new ItemInstantiateRequestBuilder and sets the default values.
func NewItemInstantiateRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemInstantiateRequestBuilder) {
    m := &ItemInstantiateRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applicationTemplates/{applicationTemplate%2Did}/instantiate", pathParameters),
    }
    return m
}
// NewItemInstantiateRequestBuilder instantiates a new ItemInstantiateRequestBuilder and sets the default values.
func NewItemInstantiateRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemInstantiateRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemInstantiateRequestBuilderInternal(urlParams, requestAdapter)
}
// Post add an instance of an application from the Microsoft Entra application gallery into your directory. The application template with ID 8adf8e6e-67b2-4cf2-a259-e3dc5476c621 can be used to add a non-gallery app that you can configure different single-sign on (SSO) modes like SAML SSO and password-based SSO.
// returns a ApplicationServicePrincipalable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/applicationtemplate-instantiate?view=graph-rest-1.0
func (m *ItemInstantiateRequestBuilder) Post(ctx context.Context, body ItemInstantiatePostRequestBodyable, requestConfiguration *ItemInstantiateRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApplicationServicePrincipalable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApplicationServicePrincipalFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApplicationServicePrincipalable), nil
}
// ToPostRequestInformation add an instance of an application from the Microsoft Entra application gallery into your directory. The application template with ID 8adf8e6e-67b2-4cf2-a259-e3dc5476c621 can be used to add a non-gallery app that you can configure different single-sign on (SSO) modes like SAML SSO and password-based SSO.
// returns a *RequestInformation when successful
func (m *ItemInstantiateRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemInstantiatePostRequestBodyable, requestConfiguration *ItemInstantiateRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemInstantiateRequestBuilder when successful
func (m *ItemInstantiateRequestBuilder) WithUrl(rawUrl string)(*ItemInstantiateRequestBuilder) {
    return NewItemInstantiateRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
