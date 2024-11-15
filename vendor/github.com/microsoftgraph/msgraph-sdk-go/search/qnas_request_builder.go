package search

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e "github.com/microsoftgraph/msgraph-sdk-go/models/search"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// QnasRequestBuilder provides operations to manage the qnas property of the microsoft.graph.searchEntity entity.
type QnasRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// QnasRequestBuilderGetQueryParameters get a list of the qna objects and their properties.
type QnasRequestBuilderGetQueryParameters struct {
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
// QnasRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type QnasRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *QnasRequestBuilderGetQueryParameters
}
// QnasRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type QnasRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByQnaId provides operations to manage the qnas property of the microsoft.graph.searchEntity entity.
// returns a *QnasQnaItemRequestBuilder when successful
func (m *QnasRequestBuilder) ByQnaId(qnaId string)(*QnasQnaItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if qnaId != "" {
        urlTplParams["qna%2Did"] = qnaId
    }
    return NewQnasQnaItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewQnasRequestBuilderInternal instantiates a new QnasRequestBuilder and sets the default values.
func NewQnasRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*QnasRequestBuilder) {
    m := &QnasRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/search/qnas{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewQnasRequestBuilder instantiates a new QnasRequestBuilder and sets the default values.
func NewQnasRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*QnasRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewQnasRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *QnasCountRequestBuilder when successful
func (m *QnasRequestBuilder) Count()(*QnasCountRequestBuilder) {
    return NewQnasCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get a list of the qna objects and their properties.
// returns a QnaCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-searchentity-list-qnas?view=graph-rest-1.0
func (m *QnasRequestBuilder) Get(ctx context.Context, requestConfiguration *QnasRequestBuilderGetRequestConfiguration)(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.QnaCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.CreateQnaCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.QnaCollectionResponseable), nil
}
// Post create a new qna object.
// returns a Qnaable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-searchentity-post-qnas?view=graph-rest-1.0
func (m *QnasRequestBuilder) Post(ctx context.Context, body i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, requestConfiguration *QnasRequestBuilderPostRequestConfiguration)(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.CreateQnaFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable), nil
}
// ToGetRequestInformation get a list of the qna objects and their properties.
// returns a *RequestInformation when successful
func (m *QnasRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *QnasRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create a new qna object.
// returns a *RequestInformation when successful
func (m *QnasRequestBuilder) ToPostRequestInformation(ctx context.Context, body i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, requestConfiguration *QnasRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *QnasRequestBuilder when successful
func (m *QnasRequestBuilder) WithUrl(rawUrl string)(*QnasRequestBuilder) {
    return NewQnasRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
