package search

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e "github.com/microsoftgraph/msgraph-sdk-go/models/search"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// QnasQnaItemRequestBuilder provides operations to manage the qnas property of the microsoft.graph.searchEntity entity.
type QnasQnaItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// QnasQnaItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type QnasQnaItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// QnasQnaItemRequestBuilderGetQueryParameters read the properties and relationships of a qna object.
type QnasQnaItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// QnasQnaItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type QnasQnaItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *QnasQnaItemRequestBuilderGetQueryParameters
}
// QnasQnaItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type QnasQnaItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewQnasQnaItemRequestBuilderInternal instantiates a new QnasQnaItemRequestBuilder and sets the default values.
func NewQnasQnaItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*QnasQnaItemRequestBuilder) {
    m := &QnasQnaItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/search/qnas/{qna%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewQnasQnaItemRequestBuilder instantiates a new QnasQnaItemRequestBuilder and sets the default values.
func NewQnasQnaItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*QnasQnaItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewQnasQnaItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete a qna object.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-qna-delete?view=graph-rest-1.0
func (m *QnasQnaItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *QnasQnaItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read the properties and relationships of a qna object.
// returns a Qnaable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-qna-get?view=graph-rest-1.0
func (m *QnasQnaItemRequestBuilder) Get(ctx context.Context, requestConfiguration *QnasQnaItemRequestBuilderGetRequestConfiguration)(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
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
// Patch update the properties of a qna object.
// returns a Qnaable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/search-qna-update?view=graph-rest-1.0
func (m *QnasQnaItemRequestBuilder) Patch(ctx context.Context, body i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, requestConfiguration *QnasQnaItemRequestBuilderPatchRequestConfiguration)(i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
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
// ToDeleteRequestInformation delete a qna object.
// returns a *RequestInformation when successful
func (m *QnasQnaItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *QnasQnaItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of a qna object.
// returns a *RequestInformation when successful
func (m *QnasQnaItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *QnasQnaItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a qna object.
// returns a *RequestInformation when successful
func (m *QnasQnaItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body i517b35a40b7cc3c50a0c7990c48f2ec92f4c4d36a97445a2aebfdc3c0071c22e.Qnaable, requestConfiguration *QnasQnaItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *QnasQnaItemRequestBuilder when successful
func (m *QnasQnaItemRequestBuilder) WithUrl(rawUrl string)(*QnasQnaItemRequestBuilder) {
    return NewQnasQnaItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
