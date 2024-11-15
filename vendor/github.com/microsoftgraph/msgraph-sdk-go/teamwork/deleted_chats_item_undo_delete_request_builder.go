package teamwork

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeletedChatsItemUndoDeleteRequestBuilder provides operations to call the undoDelete method.
type DeletedChatsItemUndoDeleteRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeletedChatsItemUndoDeleteRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedChatsItemUndoDeleteRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDeletedChatsItemUndoDeleteRequestBuilderInternal instantiates a new DeletedChatsItemUndoDeleteRequestBuilder and sets the default values.
func NewDeletedChatsItemUndoDeleteRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedChatsItemUndoDeleteRequestBuilder) {
    m := &DeletedChatsItemUndoDeleteRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/teamwork/deletedChats/{deletedChat%2Did}/undoDelete", pathParameters),
    }
    return m
}
// NewDeletedChatsItemUndoDeleteRequestBuilder instantiates a new DeletedChatsItemUndoDeleteRequestBuilder and sets the default values.
func NewDeletedChatsItemUndoDeleteRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedChatsItemUndoDeleteRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeletedChatsItemUndoDeleteRequestBuilderInternal(urlParams, requestAdapter)
}
// Post restore a  deletedChat to an active chat.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/deletedchat-undodelete?view=graph-rest-1.0
func (m *DeletedChatsItemUndoDeleteRequestBuilder) Post(ctx context.Context, requestConfiguration *DeletedChatsItemUndoDeleteRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation restore a  deletedChat to an active chat.
// returns a *RequestInformation when successful
func (m *DeletedChatsItemUndoDeleteRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *DeletedChatsItemUndoDeleteRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *DeletedChatsItemUndoDeleteRequestBuilder when successful
func (m *DeletedChatsItemUndoDeleteRequestBuilder) WithUrl(rawUrl string)(*DeletedChatsItemUndoDeleteRequestBuilder) {
    return NewDeletedChatsItemUndoDeleteRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
