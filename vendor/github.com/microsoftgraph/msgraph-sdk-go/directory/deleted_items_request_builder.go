package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeletedItemsRequestBuilder provides operations to manage the deletedItems property of the microsoft.graph.directory entity.
type DeletedItemsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeletedItemsRequestBuilderGetQueryParameters retrieve the properties of a recently deleted application, group, servicePrincipal, administrative unit, or user object from deleted items.
type DeletedItemsRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// DeletedItemsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedItemsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeletedItemsRequestBuilderGetQueryParameters
}
// ByDirectoryObjectId provides operations to manage the deletedItems property of the microsoft.graph.directory entity.
// returns a *DeletedItemsDirectoryObjectItemRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) ByDirectoryObjectId(directoryObjectId string)(*DeletedItemsDirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if directoryObjectId != "" {
        urlTplParams["directoryObject%2Did"] = directoryObjectId
    }
    return NewDeletedItemsDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewDeletedItemsRequestBuilderInternal instantiates a new DeletedItemsRequestBuilder and sets the default values.
func NewDeletedItemsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsRequestBuilder) {
    m := &DeletedItemsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/deletedItems{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewDeletedItemsRequestBuilder instantiates a new DeletedItemsRequestBuilder and sets the default values.
func NewDeletedItemsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeletedItemsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *DeletedItemsCountRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) Count()(*DeletedItemsCountRequestBuilder) {
    return NewDeletedItemsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get retrieve the properties of a recently deleted application, group, servicePrincipal, administrative unit, or user object from deleted items.
// returns a DirectoryObjectCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeletedItemsRequestBuilder) Get(ctx context.Context, requestConfiguration *DeletedItemsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable), nil
}
// GetAvailableExtensionProperties provides operations to call the getAvailableExtensionProperties method.
// returns a *DeletedItemsGetAvailableExtensionPropertiesRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GetAvailableExtensionProperties()(*DeletedItemsGetAvailableExtensionPropertiesRequestBuilder) {
    return NewDeletedItemsGetAvailableExtensionPropertiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetByIds provides operations to call the getByIds method.
// returns a *DeletedItemsGetByIdsRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GetByIds()(*DeletedItemsGetByIdsRequestBuilder) {
    return NewDeletedItemsGetByIdsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphAdministrativeUnit casts the previous resource to administrativeUnit.
// returns a *DeletedItemsGraphAdministrativeUnitRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GraphAdministrativeUnit()(*DeletedItemsGraphAdministrativeUnitRequestBuilder) {
    return NewDeletedItemsGraphAdministrativeUnitRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphApplication casts the previous resource to application.
// returns a *DeletedItemsGraphApplicationRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GraphApplication()(*DeletedItemsGraphApplicationRequestBuilder) {
    return NewDeletedItemsGraphApplicationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphDevice casts the previous resource to device.
// returns a *DeletedItemsGraphDeviceRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GraphDevice()(*DeletedItemsGraphDeviceRequestBuilder) {
    return NewDeletedItemsGraphDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphGroup casts the previous resource to group.
// returns a *DeletedItemsGraphGroupRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GraphGroup()(*DeletedItemsGraphGroupRequestBuilder) {
    return NewDeletedItemsGraphGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphServicePrincipal casts the previous resource to servicePrincipal.
// returns a *DeletedItemsGraphServicePrincipalRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GraphServicePrincipal()(*DeletedItemsGraphServicePrincipalRequestBuilder) {
    return NewDeletedItemsGraphServicePrincipalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphUser casts the previous resource to user.
// returns a *DeletedItemsGraphUserRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) GraphUser()(*DeletedItemsGraphUserRequestBuilder) {
    return NewDeletedItemsGraphUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation retrieve the properties of a recently deleted application, group, servicePrincipal, administrative unit, or user object from deleted items.
// returns a *RequestInformation when successful
func (m *DeletedItemsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeletedItemsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ValidateProperties provides operations to call the validateProperties method.
// returns a *DeletedItemsValidatePropertiesRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) ValidateProperties()(*DeletedItemsValidatePropertiesRequestBuilder) {
    return NewDeletedItemsValidatePropertiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *DeletedItemsRequestBuilder when successful
func (m *DeletedItemsRequestBuilder) WithUrl(rawUrl string)(*DeletedItemsRequestBuilder) {
    return NewDeletedItemsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
