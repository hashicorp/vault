package sites

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemGetByPathWithPathRequestBuilder provides operations to call the getByPath method.
type ItemGetByPathWithPathRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemGetByPathWithPathRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGetByPathWithPathRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Analytics provides operations to manage the analytics property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathAnalyticsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Analytics()(*ItemGetByPathWithPathAnalyticsRequestBuilder) {
    return NewItemGetByPathWithPathAnalyticsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Columns provides operations to manage the columns property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathColumnsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Columns()(*ItemGetByPathWithPathColumnsRequestBuilder) {
    return NewItemGetByPathWithPathColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemGetByPathWithPathRequestBuilderInternal instantiates a new ItemGetByPathWithPathRequestBuilder and sets the default values.
func NewItemGetByPathWithPathRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, path *string)(*ItemGetByPathWithPathRequestBuilder) {
    m := &ItemGetByPathWithPathRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/sites/{site%2Did}/getByPath(path='{path}')", pathParameters),
    }
    if path != nil {
        m.BaseRequestBuilder.PathParameters["path"] = *path
    }
    return m
}
// NewItemGetByPathWithPathRequestBuilder instantiates a new ItemGetByPathWithPathRequestBuilder and sets the default values.
func NewItemGetByPathWithPathRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGetByPathWithPathRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemGetByPathWithPathRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// ContentTypes provides operations to manage the contentTypes property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathContentTypesRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) ContentTypes()(*ItemGetByPathWithPathContentTypesRequestBuilder) {
    return NewItemGetByPathWithPathContentTypesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreatedByUser provides operations to manage the createdByUser property of the microsoft.graph.baseItem entity.
// returns a *ItemGetByPathWithPathCreatedByUserRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) CreatedByUser()(*ItemGetByPathWithPathCreatedByUserRequestBuilder) {
    return NewItemGetByPathWithPathCreatedByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drive provides operations to manage the drive property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathDriveRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Drive()(*ItemGetByPathWithPathDriveRequestBuilder) {
    return NewItemGetByPathWithPathDriveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drives provides operations to manage the drives property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathDrivesRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Drives()(*ItemGetByPathWithPathDrivesRequestBuilder) {
    return NewItemGetByPathWithPathDrivesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExternalColumns provides operations to manage the externalColumns property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathExternalColumnsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) ExternalColumns()(*ItemGetByPathWithPathExternalColumnsRequestBuilder) {
    return NewItemGetByPathWithPathExternalColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function getByPath
// returns a Siteable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemGetByPathWithPathRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemGetByPathWithPathRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable, error) {
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
// returns a *ItemGetByPathWithPathGetActivitiesByIntervalRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) GetActivitiesByInterval()(*ItemGetByPathWithPathGetActivitiesByIntervalRequestBuilder) {
    return NewItemGetByPathWithPathGetActivitiesByIntervalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
// returns a *ItemGetByPathWithPathGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*ItemGetByPathWithPathGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return NewItemGetByPathWithPathGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, interval, startDateTime)
}
// GetApplicableContentTypesForListWithListId provides operations to call the getApplicableContentTypesForList method.
// returns a *ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) GetApplicableContentTypesForListWithListId(listId *string)(*ItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilder) {
    return NewItemGetByPathWithPathGetApplicableContentTypesForListWithListIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, listId)
}
// Items provides operations to manage the items property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathItemsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Items()(*ItemGetByPathWithPathItemsRequestBuilder) {
    return NewItemGetByPathWithPathItemsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastModifiedByUser provides operations to manage the lastModifiedByUser property of the microsoft.graph.baseItem entity.
// returns a *ItemGetByPathWithPathLastModifiedByUserRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) LastModifiedByUser()(*ItemGetByPathWithPathLastModifiedByUserRequestBuilder) {
    return NewItemGetByPathWithPathLastModifiedByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Lists provides operations to manage the lists property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathListsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Lists()(*ItemGetByPathWithPathListsRequestBuilder) {
    return NewItemGetByPathWithPathListsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Onenote provides operations to manage the onenote property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathOnenoteRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Onenote()(*ItemGetByPathWithPathOnenoteRequestBuilder) {
    return NewItemGetByPathWithPathOnenoteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Operations provides operations to manage the operations property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathOperationsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Operations()(*ItemGetByPathWithPathOperationsRequestBuilder) {
    return NewItemGetByPathWithPathOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pages provides operations to manage the pages property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathPagesRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Pages()(*ItemGetByPathWithPathPagesRequestBuilder) {
    return NewItemGetByPathWithPathPagesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Permissions provides operations to manage the permissions property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathPermissionsRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Permissions()(*ItemGetByPathWithPathPermissionsRequestBuilder) {
    return NewItemGetByPathWithPathPermissionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sites provides operations to manage the sites property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathSitesRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) Sites()(*ItemGetByPathWithPathSitesRequestBuilder) {
    return NewItemGetByPathWithPathSitesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TermStore provides operations to manage the termStore property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathTermStoreRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) TermStore()(*ItemGetByPathWithPathTermStoreRequestBuilder) {
    return NewItemGetByPathWithPathTermStoreRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TermStores provides operations to manage the termStores property of the microsoft.graph.site entity.
// returns a *ItemGetByPathWithPathTermStoresRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) TermStores()(*ItemGetByPathWithPathTermStoresRequestBuilder) {
    return NewItemGetByPathWithPathTermStoresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function getByPath
// returns a *RequestInformation when successful
func (m *ItemGetByPathWithPathRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemGetByPathWithPathRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemGetByPathWithPathRequestBuilder when successful
func (m *ItemGetByPathWithPathRequestBuilder) WithUrl(rawUrl string)(*ItemGetByPathWithPathRequestBuilder) {
    return NewItemGetByPathWithPathRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
