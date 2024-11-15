package groups

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemOwnersDirectoryObjectItemRequestBuilder builds and executes requests for operations under \groups\{group-id}\owners\{directoryObject-id}
type ItemOwnersDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemOwnersDirectoryObjectItemRequestBuilderInternal instantiates a new ItemOwnersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemOwnersDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOwnersDirectoryObjectItemRequestBuilder) {
    m := &ItemOwnersDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/owners/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewItemOwnersDirectoryObjectItemRequestBuilder instantiates a new ItemOwnersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemOwnersDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOwnersDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemOwnersDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// GraphApplication casts the previous resource to application.
// returns a *ItemOwnersItemGraphApplicationRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphApplication()(*ItemOwnersItemGraphApplicationRequestBuilder) {
    return NewItemOwnersItemGraphApplicationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphDevice casts the previous resource to device.
// returns a *ItemOwnersItemGraphDeviceRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphDevice()(*ItemOwnersItemGraphDeviceRequestBuilder) {
    return NewItemOwnersItemGraphDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphGroup casts the previous resource to group.
// returns a *ItemOwnersItemGraphGroupRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphGroup()(*ItemOwnersItemGraphGroupRequestBuilder) {
    return NewItemOwnersItemGraphGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphOrgContact casts the previous resource to orgContact.
// returns a *ItemOwnersItemGraphOrgContactRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphOrgContact()(*ItemOwnersItemGraphOrgContactRequestBuilder) {
    return NewItemOwnersItemGraphOrgContactRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphServicePrincipal casts the previous resource to servicePrincipal.
// returns a *ItemOwnersItemGraphServicePrincipalRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphServicePrincipal()(*ItemOwnersItemGraphServicePrincipalRequestBuilder) {
    return NewItemOwnersItemGraphServicePrincipalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphUser casts the previous resource to user.
// returns a *ItemOwnersItemGraphUserRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphUser()(*ItemOwnersItemGraphUserRequestBuilder) {
    return NewItemOwnersItemGraphUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ref provides operations to manage the collection of group entities.
// returns a *ItemOwnersItemRefRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) Ref()(*ItemOwnersItemRefRequestBuilder) {
    return NewItemOwnersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
