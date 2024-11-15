package devices

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemRegisteredOwnersDirectoryObjectItemRequestBuilder builds and executes requests for operations under \devices\{device-id}\registeredOwners\{directoryObject-id}
type ItemRegisteredOwnersDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemRegisteredOwnersDirectoryObjectItemRequestBuilderInternal instantiates a new ItemRegisteredOwnersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemRegisteredOwnersDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) {
    m := &ItemRegisteredOwnersDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/devices/{device%2Did}/registeredOwners/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewItemRegisteredOwnersDirectoryObjectItemRequestBuilder instantiates a new ItemRegisteredOwnersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemRegisteredOwnersDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemRegisteredOwnersDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// GraphAppRoleAssignment casts the previous resource to appRoleAssignment.
// returns a *ItemRegisteredOwnersItemGraphAppRoleAssignmentRequestBuilder when successful
func (m *ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) GraphAppRoleAssignment()(*ItemRegisteredOwnersItemGraphAppRoleAssignmentRequestBuilder) {
    return NewItemRegisteredOwnersItemGraphAppRoleAssignmentRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphEndpoint casts the previous resource to endpoint.
// returns a *ItemRegisteredOwnersItemGraphEndpointRequestBuilder when successful
func (m *ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) GraphEndpoint()(*ItemRegisteredOwnersItemGraphEndpointRequestBuilder) {
    return NewItemRegisteredOwnersItemGraphEndpointRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphServicePrincipal casts the previous resource to servicePrincipal.
// returns a *ItemRegisteredOwnersItemGraphServicePrincipalRequestBuilder when successful
func (m *ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) GraphServicePrincipal()(*ItemRegisteredOwnersItemGraphServicePrincipalRequestBuilder) {
    return NewItemRegisteredOwnersItemGraphServicePrincipalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphUser casts the previous resource to user.
// returns a *ItemRegisteredOwnersItemGraphUserRequestBuilder when successful
func (m *ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) GraphUser()(*ItemRegisteredOwnersItemGraphUserRequestBuilder) {
    return NewItemRegisteredOwnersItemGraphUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ref provides operations to manage the collection of device entities.
// returns a *ItemRegisteredOwnersItemRefRequestBuilder when successful
func (m *ItemRegisteredOwnersDirectoryObjectItemRequestBuilder) Ref()(*ItemRegisteredOwnersItemRefRequestBuilder) {
    return NewItemRegisteredOwnersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
