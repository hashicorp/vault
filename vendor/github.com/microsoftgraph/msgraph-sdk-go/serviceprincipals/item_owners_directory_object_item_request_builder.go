package serviceprincipals

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemOwnersDirectoryObjectItemRequestBuilder builds and executes requests for operations under \servicePrincipals\{servicePrincipal-id}\owners\{directoryObject-id}
type ItemOwnersDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewItemOwnersDirectoryObjectItemRequestBuilderInternal instantiates a new ItemOwnersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemOwnersDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOwnersDirectoryObjectItemRequestBuilder) {
    m := &ItemOwnersDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/owners/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewItemOwnersDirectoryObjectItemRequestBuilder instantiates a new ItemOwnersDirectoryObjectItemRequestBuilder and sets the default values.
func NewItemOwnersDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOwnersDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemOwnersDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// GraphAppRoleAssignment casts the previous resource to appRoleAssignment.
// returns a *ItemOwnersItemGraphAppRoleAssignmentRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphAppRoleAssignment()(*ItemOwnersItemGraphAppRoleAssignmentRequestBuilder) {
    return NewItemOwnersItemGraphAppRoleAssignmentRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphEndpoint casts the previous resource to endpoint.
// returns a *ItemOwnersItemGraphEndpointRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) GraphEndpoint()(*ItemOwnersItemGraphEndpointRequestBuilder) {
    return NewItemOwnersItemGraphEndpointRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
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
// Ref provides operations to manage the collection of servicePrincipal entities.
// returns a *ItemOwnersItemRefRequestBuilder when successful
func (m *ItemOwnersDirectoryObjectItemRequestBuilder) Ref()(*ItemOwnersItemRefRequestBuilder) {
    return NewItemOwnersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
