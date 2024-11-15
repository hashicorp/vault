package search

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e "github.com/microsoftgraph/msgraph-sdk-go/models/search"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AcronymsRequestBuilder provides operations to manage the acronyms property of the microsoft.graph.searchEntity entity.
type AcronymsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AcronymsRequestBuilderGetQueryParameters get a list of the acronym objects and their properties.
type AcronymsRequestBuilderGetQueryParameters struct {
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
// AcronymsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AcronymsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AcronymsRequestBuilderGetQueryParameters
}
// AcronymsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AcronymsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByAcronymId provides operations to manage the acronyms property of the microsoft.graph.searchEntity entity.
// returns a *AcronymsAcronymItemRequestBuilder when successful
func (m *AcronymsRequestBuilder) ByAcronymId(acronymId string)(*AcronymsAcronymItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if acronymId != "" {
        urlTplParams["acronym%2Did"] = acronymId
    }
    return NewAcronymsAcronymItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewAcronymsRequestBuilderInternal instantiates a new AcronymsRequestBuilder and sets the default values.
func NewAcronymsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AcronymsRequestBuilder) {
    m := &AcronymsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/search/acronyms{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewAcronymsRequestBuilder instantiates a new AcronymsRequestBuilder and sets the default values.
func NewAcronymsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AcronymsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAcronymsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *AcronymsCountRequestBuilder when successful
func (m *AcronymsRequestBuilder) Count()(*AcronymsCountRequestBuilder) {
    return NewAcronymsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get a list of the acronym objects and their properties.
// returns a AcronymCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-searchentity-list-acronyms?view=graph-rest-1.0
func (m *AcronymsRequestBuilder) Get(ctx context.Context, requestConfiguration *AcronymsRequestBuilderGetRequestConfiguration)(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.AcronymCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.CreateAcronymCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.AcronymCollectionResponseable), nil
}
// Post create a new acronym object.
// returns a Acronymable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-searchentity-post-acronyms?view=graph-rest-1.0
func (m *AcronymsRequestBuilder) Post(ctx context.Context, body i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Acronymable, requestConfiguration *AcronymsRequestBuilderPostRequestConfiguration)(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Acronymable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.CreateAcronymFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Acronymable), nil
}
// ToGetRequestInformation get a list of the acronym objects and their properties.
// returns a *RequestInformation when successful
func (m *AcronymsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AcronymsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create a new acronym object.
// returns a *RequestInformation when successful
func (m *AcronymsRequestBuilder) ToPostRequestInformation(ctx context.Context, body i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Acronymable, requestConfiguration *AcronymsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *AcronymsRequestBuilder when successful
func (m *AcronymsRequestBuilder) WithUrl(rawUrl string)(*AcronymsRequestBuilder) {
    return NewAcronymsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
