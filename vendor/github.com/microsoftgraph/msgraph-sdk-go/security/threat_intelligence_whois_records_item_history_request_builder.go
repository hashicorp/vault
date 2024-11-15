package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder provides operations to manage the history property of the microsoft.graph.security.whoisRecord entity.
type ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetQueryParameters the collection of historical records associated to this WHOIS object.
type ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetQueryParameters
}
// ByWhoisHistoryRecordId provides operations to manage the history property of the microsoft.graph.security.whoisRecord entity.
// returns a *ThreatIntelligenceWhoisRecordsItemHistoryWhoisHistoryRecordItemRequestBuilder when successful
func (m *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) ByWhoisHistoryRecordId(whoisHistoryRecordId string)(*ThreatIntelligenceWhoisRecordsItemHistoryWhoisHistoryRecordItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if whoisHistoryRecordId != "" {
        urlTplParams["whoisHistoryRecord%2Did"] = whoisHistoryRecordId
    }
    return NewThreatIntelligenceWhoisRecordsItemHistoryWhoisHistoryRecordItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderInternal instantiates a new ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder and sets the default values.
func NewThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) {
    m := &ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/threatIntelligence/whoisRecords/{whoisRecord%2Did}/history{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder instantiates a new ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder and sets the default values.
func NewThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *ThreatIntelligenceWhoisRecordsItemHistoryCountRequestBuilder when successful
func (m *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) Count()(*ThreatIntelligenceWhoisRecordsItemHistoryCountRequestBuilder) {
    return NewThreatIntelligenceWhoisRecordsItemHistoryCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the collection of historical records associated to this WHOIS object.
// returns a WhoisHistoryRecordCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) Get(ctx context.Context, requestConfiguration *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.WhoisHistoryRecordCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateWhoisHistoryRecordCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.WhoisHistoryRecordCollectionResponseable), nil
}
// ToGetRequestInformation the collection of historical records associated to this WHOIS object.
// returns a *RequestInformation when successful
func (m *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder when successful
func (m *ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) WithUrl(rawUrl string)(*ThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder) {
    return NewThreatIntelligenceWhoisRecordsItemHistoryRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
