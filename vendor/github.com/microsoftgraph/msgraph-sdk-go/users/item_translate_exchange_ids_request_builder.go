package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemTranslateExchangeIdsRequestBuilder provides operations to call the translateExchangeIds method.
type ItemTranslateExchangeIdsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemTranslateExchangeIdsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTranslateExchangeIdsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemTranslateExchangeIdsRequestBuilderInternal instantiates a new ItemTranslateExchangeIdsRequestBuilder and sets the default values.
func NewItemTranslateExchangeIdsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTranslateExchangeIdsRequestBuilder) {
    m := &ItemTranslateExchangeIdsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/translateExchangeIds", pathParameters),
    }
    return m
}
// NewItemTranslateExchangeIdsRequestBuilder instantiates a new ItemTranslateExchangeIdsRequestBuilder and sets the default values.
func NewItemTranslateExchangeIdsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTranslateExchangeIdsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTranslateExchangeIdsRequestBuilderInternal(urlParams, requestAdapter)
}
// Post translate identifiers of Outlook-related resources between formats.
// Deprecated: This method is obsolete. Use PostAsTranslateExchangeIdsPostResponse instead.
// returns a ItemTranslateExchangeIdsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/user-translateexchangeids?view=graph-rest-1.0
func (m *ItemTranslateExchangeIdsRequestBuilder) Post(ctx context.Context, body ItemTranslateExchangeIdsPostRequestBodyable, requestConfiguration *ItemTranslateExchangeIdsRequestBuilderPostRequestConfiguration)(ItemTranslateExchangeIdsResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemTranslateExchangeIdsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemTranslateExchangeIdsResponseable), nil
}
// PostAsTranslateExchangeIdsPostResponse translate identifiers of Outlook-related resources between formats.
// returns a ItemTranslateExchangeIdsPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/user-translateexchangeids?view=graph-rest-1.0
func (m *ItemTranslateExchangeIdsRequestBuilder) PostAsTranslateExchangeIdsPostResponse(ctx context.Context, body ItemTranslateExchangeIdsPostRequestBodyable, requestConfiguration *ItemTranslateExchangeIdsRequestBuilderPostRequestConfiguration)(ItemTranslateExchangeIdsPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemTranslateExchangeIdsPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemTranslateExchangeIdsPostResponseable), nil
}
// ToPostRequestInformation translate identifiers of Outlook-related resources between formats.
// returns a *RequestInformation when successful
func (m *ItemTranslateExchangeIdsRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemTranslateExchangeIdsPostRequestBodyable, requestConfiguration *ItemTranslateExchangeIdsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemTranslateExchangeIdsRequestBuilder when successful
func (m *ItemTranslateExchangeIdsRequestBuilder) WithUrl(rawUrl string)(*ItemTranslateExchangeIdsRequestBuilder) {
    return NewItemTranslateExchangeIdsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
