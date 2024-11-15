package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSitesItemGetByPathWithPathRequestBuilder provides operations to call the getByPath method.
type ItemSitesItemGetByPathWithPathRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSitesItemGetByPathWithPathRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesItemGetByPathWithPathRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Analytics provides operations to manage the analytics property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathAnalyticsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Analytics()(*ItemSitesItemGetByPathWithPathAnalyticsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathAnalyticsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Columns provides operations to manage the columns property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathColumnsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Columns()(*ItemSitesItemGetByPathWithPathColumnsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemSitesItemGetByPathWithPathRequestBuilderInternal instantiates a new ItemSitesItemGetByPathWithPathRequestBuilder and sets the default values.
func NewItemSitesItemGetByPathWithPathRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, path *string)(*ItemSitesItemGetByPathWithPathRequestBuilder) {
    m := &ItemSitesItemGetByPathWithPathRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/sites/{site%2Did}/getByPath(path='{path}')", pathParameters),
    }
    if path != nil {
        m.BaseRequestBuilder.PathParameters["path"] = *path
    }
    return m
}
// NewItemSitesItemGetByPathWithPathRequestBuilder instantiates a new ItemSitesItemGetByPathWithPathRequestBuilder and sets the default values.
func NewItemSitesItemGetByPathWithPathRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesItemGetByPathWithPathRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSitesItemGetByPathWithPathRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// ContentTypes provides operations to manage the contentTypes property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathContentTypesRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) ContentTypes()(*ItemSitesItemGetByPathWithPathContentTypesRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathContentTypesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreatedByUser provides operations to manage the createdByUser property of the microsoft.graph.baseItem entity.
// returns a *ItemSitesItemGetByPathWithPathCreatedByUserRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) CreatedByUser()(*ItemSitesItemGetByPathWithPathCreatedByUserRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathCreatedByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drive provides operations to manage the drive property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathDriveRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Drive()(*ItemSitesItemGetByPathWithPathDriveRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathDriveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drives provides operations to manage the drives property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathDrivesRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Drives()(*ItemSitesItemGetByPathWithPathDrivesRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathDrivesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExternalColumns provides operations to manage the externalColumns property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathExternalColumnsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) ExternalColumns()(*ItemSitesItemGetByPathWithPathExternalColumnsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathExternalColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function getByPath
// returns a Siteable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSitesItemGetByPathWithPathRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSiteFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable), nil
}
// GetActivitiesByInterval provides operations to call the getActivitiesByInterval method.
// returns a *ItemSitesItemGetByPathWithPathGetActivitiesByIntervalRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) GetActivitiesByInterval()(*ItemSitesItemGetByPathWithPathGetActivitiesByIntervalRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathGetActivitiesByIntervalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
// returns a *ItemSitesItemGetByPathWithPathGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*ItemSitesItemGetByPathWithPathGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, interval, startDateTime)
}
// GetApplicableContentTypesForListWithListId provides operations to call the getApplicableContentTypesForList method.
// returns a *ItemSitesItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) GetApplicableContentTypesForListWithListId(listId *string)(*ItemSitesItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, listId)
}
// Items provides operations to manage the items property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathItemsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Items()(*ItemSitesItemGetByPathWithPathItemsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathItemsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastModifiedByUser provides operations to manage the lastModifiedByUser property of the microsoft.graph.baseItem entity.
// returns a *ItemSitesItemGetByPathWithPathLastModifiedByUserRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) LastModifiedByUser()(*ItemSitesItemGetByPathWithPathLastModifiedByUserRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathLastModifiedByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Lists provides operations to manage the lists property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathListsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Lists()(*ItemSitesItemGetByPathWithPathListsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathListsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Onenote provides operations to manage the onenote property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathOnenoteRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Onenote()(*ItemSitesItemGetByPathWithPathOnenoteRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathOnenoteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Operations provides operations to manage the operations property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathOperationsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Operations()(*ItemSitesItemGetByPathWithPathOperationsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pages provides operations to manage the pages property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathPagesRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Pages()(*ItemSitesItemGetByPathWithPathPagesRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathPagesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Permissions provides operations to manage the permissions property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathPermissionsRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Permissions()(*ItemSitesItemGetByPathWithPathPermissionsRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathPermissionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sites provides operations to manage the sites property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathSitesRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) Sites()(*ItemSitesItemGetByPathWithPathSitesRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathSitesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TermStore provides operations to manage the termStore property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathTermStoreRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) TermStore()(*ItemSitesItemGetByPathWithPathTermStoreRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathTermStoreRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TermStores provides operations to manage the termStores property of the microsoft.graph.site entity.
// returns a *ItemSitesItemGetByPathWithPathTermStoresRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) TermStores()(*ItemSitesItemGetByPathWithPathTermStoresRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathTermStoresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function getByPath
// returns a *RequestInformation when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSitesItemGetByPathWithPathRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemSitesItemGetByPathWithPathRequestBuilder when successful
func (m *ItemSitesItemGetByPathWithPathRequestBuilder) WithUrl(rawUrl string)(*ItemSitesItemGetByPathWithPathRequestBuilder) {
    return NewItemSitesItemGetByPathWithPathRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
