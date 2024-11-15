package print

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// SharesItemAllowedGroupsGroupItemRequestBuilder builds and executes requests for operations under \print\shares\{printerShare-id}\allowedGroups\{group-id}
type SharesItemAllowedGroupsGroupItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewSharesItemAllowedGroupsGroupItemRequestBuilderInternal instantiates a new SharesItemAllowedGroupsGroupItemRequestBuilder and sets the default values.
func NewSharesItemAllowedGroupsGroupItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemAllowedGroupsGroupItemRequestBuilder) {
    m := &SharesItemAllowedGroupsGroupItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/print/shares/{printerShare%2Did}/allowedGroups/{group%2Did}", pathParameters),
    }
    return m
}
// NewSharesItemAllowedGroupsGroupItemRequestBuilder instantiates a new SharesItemAllowedGroupsGroupItemRequestBuilder and sets the default values.
func NewSharesItemAllowedGroupsGroupItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemAllowedGroupsGroupItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSharesItemAllowedGroupsGroupItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of print entities.
// returns a *SharesItemAllowedGroupsItemRefRequestBuilder when successful
func (m *SharesItemAllowedGroupsGroupItemRequestBuilder) Ref()(*SharesItemAllowedGroupsItemRefRequestBuilder) {
    return NewSharesItemAllowedGroupsItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServiceProvisioningErrors the serviceProvisioningErrors property
// returns a *SharesItemAllowedGroupsItemServiceProvisioningErrorsRequestBuilder when successful
func (m *SharesItemAllowedGroupsGroupItemRequestBuilder) ServiceProvisioningErrors()(*SharesItemAllowedGroupsItemServiceProvisioningErrorsRequestBuilder) {
    return NewSharesItemAllowedGroupsItemServiceProvisioningErrorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
