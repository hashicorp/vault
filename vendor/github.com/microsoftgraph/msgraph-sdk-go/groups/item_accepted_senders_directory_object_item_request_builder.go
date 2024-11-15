package groups

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemAcceptedSendersDirectoryObjectItemRequestBuilder builds and executes requests for operations under \groups\{group-id}\acceptedSenders\{directoryObject-id}
type ItemAcceptedSendersDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemAcceptedSendersDirectoryObjectItemRequestBuilderInternal instantiates a new ItemAcceptedSendersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemAcceptedSendersDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAcceptedSendersDirectoryObjectItemRequestBuilder) {
    m := &ItemAcceptedSendersDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/acceptedSenders/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewItemAcceptedSendersDirectoryObjectItemRequestBuilder instantiates a new ItemAcceptedSendersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemAcceptedSendersDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAcceptedSendersDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemAcceptedSendersDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of group entities.
// returns a *ItemAcceptedSendersItemRefRequestBuilder when successful
func (m *ItemAcceptedSendersDirectoryObjectItemRequestBuilder) Ref()(*ItemAcceptedSendersItemRefRequestBuilder) {
    return NewItemAcceptedSendersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
