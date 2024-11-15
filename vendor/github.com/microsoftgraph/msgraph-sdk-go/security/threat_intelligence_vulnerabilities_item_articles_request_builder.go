package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder provides operations to manage the articles property of the microsoft.graph.security.vulnerability entity.
type ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetQueryParameters articles related to this vulnerability.
type ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetQueryParameters struct {
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
// ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetQueryParameters
}
// ByArticleId provides operations to manage the articles property of the microsoft.graph.security.vulnerability entity.
// returns a *ThreatIntelligenceVulnerabilitiesItemArticlesArticleItemRequestBuilder when successful
func (m *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) ByArticleId(articleId string)(*ThreatIntelligenceVulnerabilitiesItemArticlesArticleItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if articleId != "" {
        urlTplParams["article%2Did"] = articleId
    }
    return NewThreatIntelligenceVulnerabilitiesItemArticlesArticleItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderInternal instantiates a new ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder and sets the default values.
func NewThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) {
    m := &ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/threatIntelligence/vulnerabilities/{vulnerability%2Did}/articles{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder instantiates a new ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder and sets the default values.
func NewThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *ThreatIntelligenceVulnerabilitiesItemArticlesCountRequestBuilder when successful
func (m *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) Count()(*ThreatIntelligenceVulnerabilitiesItemArticlesCountRequestBuilder) {
    return NewThreatIntelligenceVulnerabilitiesItemArticlesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get articles related to this vulnerability.
// returns a ArticleCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) Get(ctx context.Context, requestConfiguration *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ArticleCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateArticleCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ArticleCollectionResponseable), nil
}
// ToGetRequestInformation articles related to this vulnerability.
// returns a *RequestInformation when successful
func (m *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder when successful
func (m *ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) WithUrl(rawUrl string)(*ThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder) {
    return NewThreatIntelligenceVulnerabilitiesItemArticlesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
