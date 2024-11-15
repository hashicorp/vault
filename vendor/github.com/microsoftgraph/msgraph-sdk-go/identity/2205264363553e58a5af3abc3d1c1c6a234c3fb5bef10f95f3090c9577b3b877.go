package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder provides operations to call the uploadClientCertificate method.
type B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewB2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderInternal instantiates a new B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder and sets the default values.
func NewB2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder) {
    m := &B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/b2xUserFlows/{b2xIdentityUserFlow%2Did}/apiConnectorConfiguration/postFederationSignup/uploadClientCertificate", pathParameters),
    }
    return m
}
// NewB2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder instantiates a new B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder and sets the default values.
func NewB2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewB2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderInternal(urlParams, requestAdapter)
}
// Post upload a PKCS 12 format key (.pfx) to an API connector's authentication configuration. The input is a base-64 encoded value of the PKCS 12 certificate contents. This method returns an apiConnector.
// returns a IdentityApiConnectorable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/identityapiconnector-uploadclientcertificate?view=graph-rest-1.0
func (m *B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder) Post(ctx context.Context, body B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificatePostRequestBodyable, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityApiConnectorable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentityApiConnectorFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityApiConnectorable), nil
}
// ToPostRequestInformation upload a PKCS 12 format key (.pfx) to an API connector's authentication configuration. The input is a base-64 encoded value of the PKCS 12 certificate contents. This method returns an apiConnector.
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder) ToPostRequestInformation(ctx context.Context, body B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificatePostRequestBodyable, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder) WithUrl(rawUrl string)(*B2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder) {
    return NewB2xUserFlowsItemApiConnectorConfigurationPostFederationSignupUploadClientCertificateRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
