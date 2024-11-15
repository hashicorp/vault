package directory

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder builds and executes requests for operations under \directory\administrativeUnits\{administrativeUnit-id}\members\{directoryObject-id}
type AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewAdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilderInternal instantiates a new AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder and sets the default values.
func NewAdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) {
    m := &AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/administrativeUnits/{administrativeUnit%2Did}/members/{directoryObject%2Did}", pathParameters),
    }
    return m
}
// NewAdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder instantiates a new AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder and sets the default values.
func NewAdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// GraphApplication casts the previous resource to application.
// returns a *AdministrativeUnitsItemMembersItemGraphApplicationRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) GraphApplication()(*AdministrativeUnitsItemMembersItemGraphApplicationRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemGraphApplicationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphDevice casts the previous resource to device.
// returns a *AdministrativeUnitsItemMembersItemGraphDeviceRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) GraphDevice()(*AdministrativeUnitsItemMembersItemGraphDeviceRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemGraphDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphGroup casts the previous resource to group.
// returns a *AdministrativeUnitsItemMembersItemGraphGroupRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) GraphGroup()(*AdministrativeUnitsItemMembersItemGraphGroupRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemGraphGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphOrgContact casts the previous resource to orgContact.
// returns a *AdministrativeUnitsItemMembersItemGraphOrgContactRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) GraphOrgContact()(*AdministrativeUnitsItemMembersItemGraphOrgContactRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemGraphOrgContactRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphServicePrincipal casts the previous resource to servicePrincipal.
// returns a *AdministrativeUnitsItemMembersItemGraphServicePrincipalRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) GraphServicePrincipal()(*AdministrativeUnitsItemMembersItemGraphServicePrincipalRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemGraphServicePrincipalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphUser casts the previous resource to user.
// returns a *AdministrativeUnitsItemMembersItemGraphUserRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) GraphUser()(*AdministrativeUnitsItemMembersItemGraphUserRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemGraphUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ref provides operations to manage the collection of directory entities.
// returns a *AdministrativeUnitsItemMembersItemRefRequestBuilder when successful
func (m *AdministrativeUnitsItemMembersDirectoryObjectItemRequestBuilder) Ref()(*AdministrativeUnitsItemMembersItemRefRequestBuilder) {
    return NewAdministrativeUnitsItemMembersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
