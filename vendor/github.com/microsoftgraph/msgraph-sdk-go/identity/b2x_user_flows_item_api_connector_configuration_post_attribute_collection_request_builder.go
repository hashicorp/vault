package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder provides operations to manage the postAttributeCollection property of the microsoft.graph.userFlowApiConnectorConfiguration entity.
type B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetQueryParameters get postAttributeCollection from identity
type B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetQueryParameters
}
// B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderInternal instantiates a new B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder and sets the default values.
func NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) {
    m := &B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/b2xUserFlows/{b2xIdentityUserFlow%2Did}/apiConnectorConfiguration/postAttributeCollection{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder instantiates a new B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder and sets the default values.
func NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property postAttributeCollection for identity
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) Delete(ctx context.Context, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get postAttributeCollection from identity
// returns a IdentityApiConnectorable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) Get(ctx context.Context, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityApiConnectorable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
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
// Patch update the navigation property postAttributeCollection in identity
// returns a IdentityApiConnectorable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityApiConnectorable, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityApiConnectorable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
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
// Ref provides operations to manage the collection of identityContainer entities.
// returns a *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRefRequestBuilder when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) Ref()(*B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRefRequestBuilder) {
    return NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property postAttributeCollection for identity
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get postAttributeCollection from identity
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property postAttributeCollection in identity
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityApiConnectorable, requestConfiguration *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// UploadClientCertificate provides operations to call the uploadClientCertificate method.
// returns a *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionUploadClientCertificateRequestBuilder when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) UploadClientCertificate()(*B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionUploadClientCertificateRequestBuilder) {
    return NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionUploadClientCertificateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder when successful
func (m *B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) WithUrl(rawUrl string)(*B2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder) {
    return NewB2xUserFlowsItemApiConnectorConfigurationPostAttributeCollectionRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
