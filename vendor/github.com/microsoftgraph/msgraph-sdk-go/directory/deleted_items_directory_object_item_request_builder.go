package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeletedItemsDirectoryObjectItemRequestBuilder provides operations to manage the deletedItems property of the microsoft.graph.directory entity.
type DeletedItemsDirectoryObjectItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeletedItemsDirectoryObjectItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedItemsDirectoryObjectItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DeletedItemsDirectoryObjectItemRequestBuilderGetQueryParameters retrieve the properties of a recently deleted application, group, servicePrincipal, administrative unit, or user object from deleted items.
type DeletedItemsDirectoryObjectItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DeletedItemsDirectoryObjectItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedItemsDirectoryObjectItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeletedItemsDirectoryObjectItemRequestBuilderGetQueryParameters
}
// CheckMemberGroups provides operations to call the checkMemberGroups method.
// returns a *DeletedItemsItemCheckMemberGroupsRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) CheckMemberGroups()(*DeletedItemsItemCheckMemberGroupsRequestBuilder) {
    return NewDeletedItemsItemCheckMemberGroupsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CheckMemberObjects provides operations to call the checkMemberObjects method.
// returns a *DeletedItemsItemCheckMemberObjectsRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) CheckMemberObjects()(*DeletedItemsItemCheckMemberObjectsRequestBuilder) {
    return NewDeletedItemsItemCheckMemberObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewDeletedItemsDirectoryObjectItemRequestBuilderInternal instantiates a new DeletedItemsDirectoryObjectItemRequestBuilder and sets the default values.
func NewDeletedItemsDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsDirectoryObjectItemRequestBuilder) {
    m := &DeletedItemsDirectoryObjectItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/deletedItems/{directoryObject%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDeletedItemsDirectoryObjectItemRequestBuilder instantiates a new DeletedItemsDirectoryObjectItemRequestBuilder and sets the default values.
func NewDeletedItemsDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsDirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeletedItemsDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete permanently delete a recently deleted application, group, servicePrincipal, or user object from deleted items. After an item is permanently deleted, it cannot be restored. Administrative units cannot be permanently deleted by using the deletedItems API. Soft-deleted administrative units will be permanently deleted 30 days after initial deletion unless they are restored.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/directory-deleteditems-delete?view=graph-rest-1.0
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DeletedItemsDirectoryObjectItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
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
// Get retrieve the properties of a recently deleted application, group, servicePrincipal, administrative unit, or user object from deleted items.
// returns a DirectoryObjectable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/directory-deleteditems-get?view=graph-rest-1.0
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DeletedItemsDirectoryObjectItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable), nil
}
// GetMemberGroups provides operations to call the getMemberGroups method.
// returns a *DeletedItemsItemGetMemberGroupsRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GetMemberGroups()(*DeletedItemsItemGetMemberGroupsRequestBuilder) {
    return NewDeletedItemsItemGetMemberGroupsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetMemberObjects provides operations to call the getMemberObjects method.
// returns a *DeletedItemsItemGetMemberObjectsRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GetMemberObjects()(*DeletedItemsItemGetMemberObjectsRequestBuilder) {
    return NewDeletedItemsItemGetMemberObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphAdministrativeUnit casts the previous resource to administrativeUnit.
// returns a *DeletedItemsItemGraphAdministrativeUnitRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GraphAdministrativeUnit()(*DeletedItemsItemGraphAdministrativeUnitRequestBuilder) {
    return NewDeletedItemsItemGraphAdministrativeUnitRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphApplication casts the previous resource to application.
// returns a *DeletedItemsItemGraphApplicationRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GraphApplication()(*DeletedItemsItemGraphApplicationRequestBuilder) {
    return NewDeletedItemsItemGraphApplicationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphDevice casts the previous resource to device.
// returns a *DeletedItemsItemGraphDeviceRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GraphDevice()(*DeletedItemsItemGraphDeviceRequestBuilder) {
    return NewDeletedItemsItemGraphDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphGroup casts the previous resource to group.
// returns a *DeletedItemsItemGraphGroupRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GraphGroup()(*DeletedItemsItemGraphGroupRequestBuilder) {
    return NewDeletedItemsItemGraphGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphServicePrincipal casts the previous resource to servicePrincipal.
// returns a *DeletedItemsItemGraphServicePrincipalRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GraphServicePrincipal()(*DeletedItemsItemGraphServicePrincipalRequestBuilder) {
    return NewDeletedItemsItemGraphServicePrincipalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphUser casts the previous resource to user.
// returns a *DeletedItemsItemGraphUserRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) GraphUser()(*DeletedItemsItemGraphUserRequestBuilder) {
    return NewDeletedItemsItemGraphUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Restore provides operations to call the restore method.
// returns a *DeletedItemsItemRestoreRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) Restore()(*DeletedItemsItemRestoreRequestBuilder) {
    return NewDeletedItemsItemRestoreRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation permanently delete a recently deleted application, group, servicePrincipal, or user object from deleted items. After an item is permanently deleted, it cannot be restored. Administrative units cannot be permanently deleted by using the deletedItems API. Soft-deleted administrative units will be permanently deleted 30 days after initial deletion unless they are restored.
// returns a *RequestInformation when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *DeletedItemsDirectoryObjectItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation retrieve the properties of a recently deleted application, group, servicePrincipal, administrative unit, or user object from deleted items.
// returns a *RequestInformation when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeletedItemsDirectoryObjectItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *DeletedItemsDirectoryObjectItemRequestBuilder when successful
func (m *DeletedItemsDirectoryObjectItemRequestBuilder) WithUrl(rawUrl string)(*DeletedItemsDirectoryObjectItemRequestBuilder) {
    return NewDeletedItemsDirectoryObjectItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
