package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// ThreatIntelligenceRequestBuilder provides operations to manage the threatIntelligence property of the microsoft.graph.security entity.
type ThreatIntelligenceRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ThreatIntelligenceRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ThreatIntelligenceRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ThreatIntelligenceRequestBuilderGetQueryParameters get threatIntelligence from security
type ThreatIntelligenceRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ThreatIntelligenceRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ThreatIntelligenceRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ThreatIntelligenceRequestBuilderGetQueryParameters
}
// ThreatIntelligenceRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ThreatIntelligenceRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ArticleIndicators provides operations to manage the articleIndicators property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceArticleIndicatorsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) ArticleIndicators()(*ThreatIntelligenceArticleIndicatorsRequestBuilder) {
    return NewThreatIntelligenceArticleIndicatorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Articles provides operations to manage the articles property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceArticlesRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) Articles()(*ThreatIntelligenceArticlesRequestBuilder) {
    return NewThreatIntelligenceArticlesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewThreatIntelligenceRequestBuilderInternal instantiates a new ThreatIntelligenceRequestBuilder and sets the default values.
func NewThreatIntelligenceRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceRequestBuilder) {
    m := &ThreatIntelligenceRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/threatIntelligence{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewThreatIntelligenceRequestBuilder instantiates a new ThreatIntelligenceRequestBuilder and sets the default values.
func NewThreatIntelligenceRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewThreatIntelligenceRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property threatIntelligence for security
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ThreatIntelligenceRequestBuilder) Delete(ctx context.Context, requestConfiguration *ThreatIntelligenceRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get threatIntelligence from security
// returns a ThreatIntelligenceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ThreatIntelligenceRequestBuilder) Get(ctx context.Context, requestConfiguration *ThreatIntelligenceRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ThreatIntelligenceable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateThreatIntelligenceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ThreatIntelligenceable), nil
}
// HostComponents provides operations to manage the hostComponents property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostComponentsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) HostComponents()(*ThreatIntelligenceHostComponentsRequestBuilder) {
    return NewThreatIntelligenceHostComponentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// HostCookies provides operations to manage the hostCookies property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostCookiesRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) HostCookies()(*ThreatIntelligenceHostCookiesRequestBuilder) {
    return NewThreatIntelligenceHostCookiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// HostPairs provides operations to manage the hostPairs property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostPairsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) HostPairs()(*ThreatIntelligenceHostPairsRequestBuilder) {
    return NewThreatIntelligenceHostPairsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// HostPorts provides operations to manage the hostPorts property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostPortsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) HostPorts()(*ThreatIntelligenceHostPortsRequestBuilder) {
    return NewThreatIntelligenceHostPortsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hosts provides operations to manage the hosts property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) Hosts()(*ThreatIntelligenceHostsRequestBuilder) {
    return NewThreatIntelligenceHostsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// HostSslCertificates provides operations to manage the hostSslCertificates property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostSslCertificatesRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) HostSslCertificates()(*ThreatIntelligenceHostSslCertificatesRequestBuilder) {
    return NewThreatIntelligenceHostSslCertificatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// HostTrackers provides operations to manage the hostTrackers property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceHostTrackersRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) HostTrackers()(*ThreatIntelligenceHostTrackersRequestBuilder) {
    return NewThreatIntelligenceHostTrackersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntelligenceProfileIndicators provides operations to manage the intelligenceProfileIndicators property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceIntelligenceProfileIndicatorsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) IntelligenceProfileIndicators()(*ThreatIntelligenceIntelligenceProfileIndicatorsRequestBuilder) {
    return NewThreatIntelligenceIntelligenceProfileIndicatorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntelProfiles provides operations to manage the intelProfiles property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceIntelProfilesRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) IntelProfiles()(*ThreatIntelligenceIntelProfilesRequestBuilder) {
    return NewThreatIntelligenceIntelProfilesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PassiveDnsRecords provides operations to manage the passiveDnsRecords property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligencePassiveDnsRecordsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) PassiveDnsRecords()(*ThreatIntelligencePassiveDnsRecordsRequestBuilder) {
    return NewThreatIntelligencePassiveDnsRecordsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property threatIntelligence in security
// returns a ThreatIntelligenceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ThreatIntelligenceRequestBuilder) Patch(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ThreatIntelligenceable, requestConfiguration *ThreatIntelligenceRequestBuilderPatchRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ThreatIntelligenceable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateThreatIntelligenceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ThreatIntelligenceable), nil
}
// SslCertificates provides operations to manage the sslCertificates property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceSslCertificatesRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) SslCertificates()(*ThreatIntelligenceSslCertificatesRequestBuilder) {
    return NewThreatIntelligenceSslCertificatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Subdomains provides operations to manage the subdomains property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceSubdomainsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) Subdomains()(*ThreatIntelligenceSubdomainsRequestBuilder) {
    return NewThreatIntelligenceSubdomainsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property threatIntelligence for security
// returns a *RequestInformation when successful
func (m *ThreatIntelligenceRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ThreatIntelligenceRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get threatIntelligence from security
// returns a *RequestInformation when successful
func (m *ThreatIntelligenceRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ThreatIntelligenceRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property threatIntelligence in security
// returns a *RequestInformation when successful
func (m *ThreatIntelligenceRequestBuilder) ToPatchRequestInformation(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ThreatIntelligenceable, requestConfiguration *ThreatIntelligenceRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Vulnerabilities provides operations to manage the vulnerabilities property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceVulnerabilitiesRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) Vulnerabilities()(*ThreatIntelligenceVulnerabilitiesRequestBuilder) {
    return NewThreatIntelligenceVulnerabilitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WhoisHistoryRecords provides operations to manage the whoisHistoryRecords property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceWhoisHistoryRecordsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) WhoisHistoryRecords()(*ThreatIntelligenceWhoisHistoryRecordsRequestBuilder) {
    return NewThreatIntelligenceWhoisHistoryRecordsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WhoisRecords provides operations to manage the whoisRecords property of the microsoft.graph.security.threatIntelligence entity.
// returns a *ThreatIntelligenceWhoisRecordsRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) WhoisRecords()(*ThreatIntelligenceWhoisRecordsRequestBuilder) {
    return NewThreatIntelligenceWhoisRecordsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ThreatIntelligenceRequestBuilder when successful
func (m *ThreatIntelligenceRequestBuilder) WithUrl(rawUrl string)(*ThreatIntelligenceRequestBuilder) {
    return NewThreatIntelligenceRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
