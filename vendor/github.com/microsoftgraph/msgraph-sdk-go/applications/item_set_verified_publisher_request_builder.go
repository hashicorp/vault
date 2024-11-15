package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSetVerifiedPublisherRequestBuilder provides operations to call the setVerifiedPublisher method.
type ItemSetVerifiedPublisherRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSetVerifiedPublisherRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSetVerifiedPublisherRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSetVerifiedPublisherRequestBuilderInternal instantiates a new ItemSetVerifiedPublisherRequestBuilder and sets the default values.
func NewItemSetVerifiedPublisherRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSetVerifiedPublisherRequestBuilder) {
    m := &ItemSetVerifiedPublisherRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/setVerifiedPublisher", pathParameters),
    }
    return m
}
// NewItemSetVerifiedPublisherRequestBuilder instantiates a new ItemSetVerifiedPublisherRequestBuilder and sets the default values.
func NewItemSetVerifiedPublisherRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSetVerifiedPublisherRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSetVerifiedPublisherRequestBuilderInternal(urlParams, requestAdapter)
}
// Post set the verifiedPublisher on an application. For more information, including prerequisites to setting a verified publisher, see Publisher verification.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/application-setverifiedpublisher?view=graph-rest-1.0
func (m *ItemSetVerifiedPublisherRequestBuilder) Post(ctx context.Context, body ItemSetVerifiedPublisherPostRequestBodyable, requestConfiguration *ItemSetVerifiedPublisherRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation set the verifiedPublisher on an application. For more information, including prerequisites to setting a verified publisher, see Publisher verification.
// returns a *RequestInformation when successful
func (m *ItemSetVerifiedPublisherRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemSetVerifiedPublisherPostRequestBodyable, requestConfiguration *ItemSetVerifiedPublisherRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSetVerifiedPublisherRequestBuilder when successful
func (m *ItemSetVerifiedPublisherRequestBuilder) WithUrl(rawUrl string)(*ItemSetVerifiedPublisherRequestBuilder) {
    return NewItemSetVerifiedPublisherRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
