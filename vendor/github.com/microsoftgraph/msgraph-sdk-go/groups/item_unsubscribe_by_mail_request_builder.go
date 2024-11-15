package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemUnsubscribeByMailRequestBuilder provides operations to call the unsubscribeByMail method.
type ItemUnsubscribeByMailRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemUnsubscribeByMailRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemUnsubscribeByMailRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemUnsubscribeByMailRequestBuilderInternal instantiates a new ItemUnsubscribeByMailRequestBuilder and sets the default values.
func NewItemUnsubscribeByMailRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemUnsubscribeByMailRequestBuilder) {
    m := &ItemUnsubscribeByMailRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/unsubscribeByMail", pathParameters),
    }
    return m
}
// NewItemUnsubscribeByMailRequestBuilder instantiates a new ItemUnsubscribeByMailRequestBuilder and sets the default values.
func NewItemUnsubscribeByMailRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemUnsubscribeByMailRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemUnsubscribeByMailRequestBuilderInternal(urlParams, requestAdapter)
}
// Post calling this method prevents the current user from receiving email notifications for this group about new posts, events, and files in that group. Supported for Microsoft 365 groups only.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/group-unsubscribebymail?view=graph-rest-1.0
func (m *ItemUnsubscribeByMailRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemUnsubscribeByMailRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation calling this method prevents the current user from receiving email notifications for this group about new posts, events, and files in that group. Supported for Microsoft 365 groups only.
// returns a *RequestInformation when successful
func (m *ItemUnsubscribeByMailRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemUnsubscribeByMailRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemUnsubscribeByMailRequestBuilder when successful
func (m *ItemUnsubscribeByMailRequestBuilder) WithUrl(rawUrl string)(*ItemUnsubscribeByMailRequestBuilder) {
    return NewItemUnsubscribeByMailRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
