package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder provides operations to call the replyAll method.
type ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderInternal instantiates a new ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder and sets the default values.
func NewItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder) {
    m := &ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/mailFolders/{mailFolder%2Did}/childFolders/{mailFolder%2Did1}/messages/{message%2Did}/replyAll", pathParameters),
    }
    return m
}
// NewItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder instantiates a new ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder and sets the default values.
func NewItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderInternal(urlParams, requestAdapter)
}
// Post reply to all recipients of a message using either JSON or MIME format. When using JSON format:- Specify either a comment or the body property of the message parameter. Specifying both will return an HTTP 400 Bad Request error.- If the original message specifies a recipient in the replyTo property, per Internet Message Format (RFC 2822), send the reply to the recipients in replyTo and not the recipient in the from property. When using MIME format:- Provide the applicable Internet message headers and the MIME content, all encoded in base64 format in the request body.- Add any attachments and S/MIME properties to the MIME content. This method saves the message in the Sent Items folder. Alternatively, create a draft to reply-all to a message and send it later.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/message-replyall?view=graph-rest-1.0
func (m *ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder) Post(ctx context.Context, body ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllPostRequestBodyable, requestConfiguration *ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation reply to all recipients of a message using either JSON or MIME format. When using JSON format:- Specify either a comment or the body property of the message parameter. Specifying both will return an HTTP 400 Bad Request error.- If the original message specifies a recipient in the replyTo property, per Internet Message Format (RFC 2822), send the reply to the recipients in replyTo and not the recipient in the from property. When using MIME format:- Provide the applicable Internet message headers and the MIME content, all encoded in base64 format in the request body.- Add any attachments and S/MIME properties to the MIME content. This method saves the message in the Sent Items folder. Alternatively, create a draft to reply-all to a message and send it later.
// returns a *RequestInformation when successful
func (m *ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllPostRequestBodyable, requestConfiguration *ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder when successful
func (m *ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder) WithUrl(rawUrl string)(*ItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder) {
    return NewItemMailFoldersItemChildFoldersItemMessagesItemReplyAllRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
