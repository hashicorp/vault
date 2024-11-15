package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemUnsetVerifiedPublisherRequestBuilder provides operations to call the unsetVerifiedPublisher method.
type ItemUnsetVerifiedPublisherRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemUnsetVerifiedPublisherRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemUnsetVerifiedPublisherRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemUnsetVerifiedPublisherRequestBuilderInternal instantiates a new ItemUnsetVerifiedPublisherRequestBuilder and sets the default values.
func NewItemUnsetVerifiedPublisherRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemUnsetVerifiedPublisherRequestBuilder) {
    m := &ItemUnsetVerifiedPublisherRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/unsetVerifiedPublisher", pathParameters),
    }
    return m
}
// NewItemUnsetVerifiedPublisherRequestBuilder instantiates a new ItemUnsetVerifiedPublisherRequestBuilder and sets the default values.
func NewItemUnsetVerifiedPublisherRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemUnsetVerifiedPublisherRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemUnsetVerifiedPublisherRequestBuilderInternal(urlParams, requestAdapter)
}
// Post unset the verifiedPublisher previously set on an application, removing all verified publisher properties. For more information, see Publisher verification.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/application-unsetverifiedpublisher?view=graph-rest-1.0
func (m *ItemUnsetVerifiedPublisherRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemUnsetVerifiedPublisherRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation unset the verifiedPublisher previously set on an application, removing all verified publisher properties. For more information, see Publisher verification.
// returns a *RequestInformation when successful
func (m *ItemUnsetVerifiedPublisherRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemUnsetVerifiedPublisherRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemUnsetVerifiedPublisherRequestBuilder when successful
func (m *ItemUnsetVerifiedPublisherRequestBuilder) WithUrl(rawUrl string)(*ItemUnsetVerifiedPublisherRequestBuilder) {
    return NewItemUnsetVerifiedPublisherRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
