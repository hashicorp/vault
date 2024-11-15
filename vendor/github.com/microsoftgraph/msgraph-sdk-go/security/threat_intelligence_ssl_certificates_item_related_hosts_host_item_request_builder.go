package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder provides operations to manage the relatedHosts property of the microsoft.graph.security.sslCertificate entity.
type ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetQueryParameters the host resources related with this sslCertificate.
type ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetQueryParameters
}
// NewThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderInternal instantiates a new ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder and sets the default values.
func NewThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder) {
    m := &ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/threatIntelligence/sslCertificates/{sslCertificate%2Did}/relatedHosts/{host%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder instantiates a new ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder and sets the default values.
func NewThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the host resources related with this sslCertificate.
// returns a Hostable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.Hostable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateHostFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.Hostable), nil
}
// ToGetRequestInformation the host resources related with this sslCertificate.
// returns a *RequestInformation when successful
func (m *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder when successful
func (m *ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder) WithUrl(rawUrl string)(*ThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder) {
    return NewThreatIntelligenceSslCertificatesItemRelatedHostsHostItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
