package communications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// PresencesItemSetStatusMessageRequestBuilder provides operations to call the setStatusMessage method.
type PresencesItemSetStatusMessageRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// PresencesItemSetStatusMessageRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PresencesItemSetStatusMessageRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewPresencesItemSetStatusMessageRequestBuilderInternal instantiates a new PresencesItemSetStatusMessageRequestBuilder and sets the default values.
func NewPresencesItemSetStatusMessageRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PresencesItemSetStatusMessageRequestBuilder) {
    m := &PresencesItemSetStatusMessageRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/communications/presences/{presence%2Did}/setStatusMessage", pathParameters),
    }
    return m
}
// NewPresencesItemSetStatusMessageRequestBuilder instantiates a new PresencesItemSetStatusMessageRequestBuilder and sets the default values.
func NewPresencesItemSetStatusMessageRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PresencesItemSetStatusMessageRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPresencesItemSetStatusMessageRequestBuilderInternal(urlParams, requestAdapter)
}
// Post set a presence status message for a user. An optional expiration date and time can be supplied.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/presence-setstatusmessage?view=graph-rest-1.0
func (m *PresencesItemSetStatusMessageRequestBuilder) Post(ctx context.Context, body PresencesItemSetStatusMessagePostRequestBodyable, requestConfiguration *PresencesItemSetStatusMessageRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
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
// ToPostRequestInformation set a presence status message for a user. An optional expiration date and time can be supplied.
// returns a *RequestInformation when successful
func (m *PresencesItemSetStatusMessageRequestBuilder) ToPostRequestInformation(ctx context.Context, body PresencesItemSetStatusMessagePostRequestBodyable, requestConfiguration *PresencesItemSetStatusMessageRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *PresencesItemSetStatusMessageRequestBuilder when successful
func (m *PresencesItemSetStatusMessageRequestBuilder) WithUrl(rawUrl string)(*PresencesItemSetStatusMessageRequestBuilder) {
    return NewPresencesItemSetStatusMessageRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
