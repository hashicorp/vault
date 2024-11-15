package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsDriveItemItemRequestBuilder provides operations to manage the items property of the microsoft.graph.drive entity.
type ItemItemsDriveItemItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsDriveItemItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsDriveItemItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemItemsDriveItemItemRequestBuilderGetQueryParameters all items contained in the drive. Read-only. Nullable.
type ItemItemsDriveItemItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemItemsDriveItemItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsDriveItemItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemItemsDriveItemItemRequestBuilderGetQueryParameters
}
// ItemItemsDriveItemItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsDriveItemItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Analytics provides operations to manage the analytics property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemAnalyticsRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Analytics()(*ItemItemsItemAnalyticsRequestBuilder) {
    return NewItemItemsItemAnalyticsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AssignSensitivityLabel provides operations to call the assignSensitivityLabel method.
// returns a *ItemItemsItemAssignSensitivityLabelRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) AssignSensitivityLabel()(*ItemItemsItemAssignSensitivityLabelRequestBuilder) {
    return NewItemItemsItemAssignSensitivityLabelRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Checkin provides operations to call the checkin method.
// returns a *ItemItemsItemCheckinRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Checkin()(*ItemItemsItemCheckinRequestBuilder) {
    return NewItemItemsItemCheckinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Checkout provides operations to call the checkout method.
// returns a *ItemItemsItemCheckoutRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Checkout()(*ItemItemsItemCheckoutRequestBuilder) {
    return NewItemItemsItemCheckoutRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemChildrenRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Children()(*ItemItemsItemChildrenRequestBuilder) {
    return NewItemItemsItemChildrenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemItemsDriveItemItemRequestBuilderInternal instantiates a new ItemItemsDriveItemItemRequestBuilder and sets the default values.
func NewItemItemsDriveItemItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsDriveItemItemRequestBuilder) {
    m := &ItemItemsDriveItemItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemItemsDriveItemItemRequestBuilder instantiates a new ItemItemsDriveItemItemRequestBuilder and sets the default values.
func NewItemItemsDriveItemItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsDriveItemItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsDriveItemItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Content provides operations to manage the media for the drive entity.
// returns a *ItemItemsItemContentRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Content()(*ItemItemsItemContentRequestBuilder) {
    return NewItemItemsItemContentRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Copy provides operations to call the copy method.
// returns a *ItemItemsItemCopyRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Copy()(*ItemItemsItemCopyRequestBuilder) {
    return NewItemItemsItemCopyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreatedByUser provides operations to manage the createdByUser property of the microsoft.graph.baseItem entity.
// returns a *ItemItemsItemCreatedByUserRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) CreatedByUser()(*ItemItemsItemCreatedByUserRequestBuilder) {
    return NewItemItemsItemCreatedByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreateLink provides operations to call the createLink method.
// returns a *ItemItemsItemCreateLinkRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) CreateLink()(*ItemItemsItemCreateLinkRequestBuilder) {
    return NewItemItemsItemCreateLinkRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreateUploadSession provides operations to call the createUploadSession method.
// returns a *ItemItemsItemCreateUploadSessionRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) CreateUploadSession()(*ItemItemsItemCreateUploadSessionRequestBuilder) {
    return NewItemItemsItemCreateUploadSessionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delete delete navigation property items for drives
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsDriveItemItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemItemsDriveItemItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Delta provides operations to call the delta method.
// returns a *ItemItemsItemDeltaRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Delta()(*ItemItemsItemDeltaRequestBuilder) {
    return NewItemItemsItemDeltaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeltaWithToken provides operations to call the delta method.
// returns a *ItemItemsItemDeltaWithTokenRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) DeltaWithToken(token *string)(*ItemItemsItemDeltaWithTokenRequestBuilder) {
    return NewItemItemsItemDeltaWithTokenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, token)
}
// ExtractSensitivityLabels provides operations to call the extractSensitivityLabels method.
// returns a *ItemItemsItemExtractSensitivityLabelsRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) ExtractSensitivityLabels()(*ItemItemsItemExtractSensitivityLabelsRequestBuilder) {
    return NewItemItemsItemExtractSensitivityLabelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Follow provides operations to call the follow method.
// returns a *ItemItemsItemFollowRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Follow()(*ItemItemsItemFollowRequestBuilder) {
    return NewItemItemsItemFollowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get all items contained in the drive. Read-only. Nullable.
// returns a DriveItemable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsDriveItemItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsDriveItemItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable), nil
}
// GetActivitiesByInterval provides operations to call the getActivitiesByInterval method.
// returns a *ItemItemsItemGetActivitiesByIntervalRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) GetActivitiesByInterval()(*ItemItemsItemGetActivitiesByIntervalRequestBuilder) {
    return NewItemItemsItemGetActivitiesByIntervalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
// returns a *ItemItemsItemGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*ItemItemsItemGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return NewItemItemsItemGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, interval, startDateTime)
}
// Invite provides operations to call the invite method.
// returns a *ItemItemsItemInviteRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Invite()(*ItemItemsItemInviteRequestBuilder) {
    return NewItemItemsItemInviteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastModifiedByUser provides operations to manage the lastModifiedByUser property of the microsoft.graph.baseItem entity.
// returns a *ItemItemsItemLastModifiedByUserRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) LastModifiedByUser()(*ItemItemsItemLastModifiedByUserRequestBuilder) {
    return NewItemItemsItemLastModifiedByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemListItemRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) ListItem()(*ItemItemsItemListItemRequestBuilder) {
    return NewItemItemsItemListItemRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property items in drives
// returns a DriveItemable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsDriveItemItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, requestConfiguration *ItemItemsDriveItemItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable), nil
}
// PermanentDelete provides operations to call the permanentDelete method.
// returns a *ItemItemsItemPermanentDeleteRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) PermanentDelete()(*ItemItemsItemPermanentDeleteRequestBuilder) {
    return NewItemItemsItemPermanentDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Permissions provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemPermissionsRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Permissions()(*ItemItemsItemPermissionsRequestBuilder) {
    return NewItemItemsItemPermissionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Preview provides operations to call the preview method.
// returns a *ItemItemsItemPreviewRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Preview()(*ItemItemsItemPreviewRequestBuilder) {
    return NewItemItemsItemPreviewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Restore provides operations to call the restore method.
// returns a *ItemItemsItemRestoreRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Restore()(*ItemItemsItemRestoreRequestBuilder) {
    return NewItemItemsItemRestoreRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RetentionLabel provides operations to manage the retentionLabel property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemRetentionLabelRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) RetentionLabel()(*ItemItemsItemRetentionLabelRequestBuilder) {
    return NewItemItemsItemRetentionLabelRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SearchWithQ provides operations to call the search method.
// returns a *ItemItemsItemSearchWithQRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) SearchWithQ(q *string)(*ItemItemsItemSearchWithQRequestBuilder) {
    return NewItemItemsItemSearchWithQRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, q)
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemSubscriptionsRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Subscriptions()(*ItemItemsItemSubscriptionsRequestBuilder) {
    return NewItemItemsItemSubscriptionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemThumbnailsRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Thumbnails()(*ItemItemsItemThumbnailsRequestBuilder) {
    return NewItemItemsItemThumbnailsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property items for drives
// returns a *RequestInformation when successful
func (m *ItemItemsDriveItemItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemItemsDriveItemItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation all items contained in the drive. Read-only. Nullable.
// returns a *RequestInformation when successful
func (m *ItemItemsDriveItemItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsDriveItemItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property items in drives
// returns a *RequestInformation when successful
func (m *ItemItemsDriveItemItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, requestConfiguration *ItemItemsDriveItemItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// Unfollow provides operations to call the unfollow method.
// returns a *ItemItemsItemUnfollowRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Unfollow()(*ItemItemsItemUnfollowRequestBuilder) {
    return NewItemItemsItemUnfollowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ValidatePermission provides operations to call the validatePermission method.
// returns a *ItemItemsItemValidatePermissionRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) ValidatePermission()(*ItemItemsItemValidatePermissionRequestBuilder) {
    return NewItemItemsItemValidatePermissionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemVersionsRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Versions()(*ItemItemsItemVersionsRequestBuilder) {
    return NewItemItemsItemVersionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsDriveItemItemRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) WithUrl(rawUrl string)(*ItemItemsDriveItemItemRequestBuilder) {
    return NewItemItemsDriveItemItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Workbook provides operations to manage the workbook property of the microsoft.graph.driveItem entity.
// returns a *ItemItemsItemWorkbookRequestBuilder when successful
func (m *ItemItemsDriveItemItemRequestBuilder) Workbook()(*ItemItemsItemWorkbookRequestBuilder) {
    return NewItemItemsItemWorkbookRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
