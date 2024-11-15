package groups

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemRejectedSendersDirectoryObjectItemRequestBuilder builds and executes requests for operations under \groups\{group-id}\rejectedSenders\{directoryObject-id}
type ItemRejectedSendersDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemRejectedSendersDirectoryObjectItemRequestBuilderInternal instantiates a new ItemRejectedSendersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemRejectedSendersDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemRejectedSendersDirectoryObjectItemRequestBuilder) {
    m := &ItemRejectedSendersDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/rejectedSenders/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewItemRejectedSendersDirectoryObjectItemRequestBuilder instantiates a new ItemRejectedSendersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemRejectedSendersDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemRejectedSendersDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemRejectedSendersDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of group entities.
// returns a *ItemRejectedSendersItemRefRequestBuilder when successful
func (m *ItemRejectedSendersDirectoryObjectItemRequestBuilder) Ref()(*ItemRejectedSendersItemRefRequestBuilder) {
    return NewItemRejectedSendersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
