package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder provides operations to manage the restorePoint property of the microsoft.graph.restoreArtifactBase entity.
type BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetQueryParameters represents the date and time when an artifact is protected by a protectionPolicy and can be restored.
type BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetQueryParameters
}
// NewBackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderInternal instantiates a new BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder and sets the default values.
func NewBackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder) {
    m := &BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/oneDriveForBusinessRestoreSessions/{oneDriveForBusinessRestoreSession%2Did}/driveRestoreArtifacts/{driveRestoreArtifact%2Did}/restorePoint{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewBackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder instantiates a new BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder and sets the default values.
func NewBackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderInternal(urlParams, requestAdapter)
}
// Get represents the date and time when an artifact is protected by a protectionPolicy and can be restored.
// returns a RestorePointable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RestorePointable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateRestorePointFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RestorePointable), nil
}
// ToGetRequestInformation represents the date and time when an artifact is protected by a protectionPolicy and can be restored.
// returns a *RequestInformation when successful
func (m *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessRestoreSessionsItemDriveRestoreArtifactsItemRestorePointRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
