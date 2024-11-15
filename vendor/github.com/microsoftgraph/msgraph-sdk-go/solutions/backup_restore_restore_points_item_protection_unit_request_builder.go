package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreRestorePointsItemProtectionUnitRequestBuilder provides operations to manage the protectionUnit property of the microsoft.graph.restorePoint entity.
type BackupRestoreRestorePointsItemProtectionUnitRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetQueryParameters the site, drive, or mailbox units that are protected under a protection policy.
type BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetQueryParameters
}
// NewBackupRestoreRestorePointsItemProtectionUnitRequestBuilderInternal instantiates a new BackupRestoreRestorePointsItemProtectionUnitRequestBuilder and sets the default values.
func NewBackupRestoreRestorePointsItemProtectionUnitRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreRestorePointsItemProtectionUnitRequestBuilder) {
    m := &BackupRestoreRestorePointsItemProtectionUnitRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/restorePoints/{restorePoint%2Did}/protectionUnit{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewBackupRestoreRestorePointsItemProtectionUnitRequestBuilder instantiates a new BackupRestoreRestorePointsItemProtectionUnitRequestBuilder and sets the default values.
func NewBackupRestoreRestorePointsItemProtectionUnitRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreRestorePointsItemProtectionUnitRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreRestorePointsItemProtectionUnitRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the site, drive, or mailbox units that are protected under a protection policy.
// returns a ProtectionUnitBaseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreRestorePointsItemProtectionUnitRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ProtectionUnitBaseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateProtectionUnitBaseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ProtectionUnitBaseable), nil
}
// ToGetRequestInformation the site, drive, or mailbox units that are protected under a protection policy.
// returns a *RequestInformation when successful
func (m *BackupRestoreRestorePointsItemProtectionUnitRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreRestorePointsItemProtectionUnitRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreRestorePointsItemProtectionUnitRequestBuilder when successful
func (m *BackupRestoreRestorePointsItemProtectionUnitRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreRestorePointsItemProtectionUnitRequestBuilder) {
    return NewBackupRestoreRestorePointsItemProtectionUnitRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
