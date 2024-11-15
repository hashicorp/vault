package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder provides operations to manage the restorePoint property of the microsoft.graph.restoreArtifactBase entity.
type BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetQueryParameters represents the date and time when an artifact is protected by a protectionPolicy and can be restored.
type BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetQueryParameters
}
// NewBackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderInternal instantiates a new BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder and sets the default values.
func NewBackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder) {
    m := &BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/exchangeRestoreSessions/{exchangeRestoreSession%2Did}/mailboxRestoreArtifacts/{mailboxRestoreArtifact%2Did}/restorePoint{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewBackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder instantiates a new BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder and sets the default values.
func NewBackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderInternal(urlParams, requestAdapter)
}
// Get represents the date and time when an artifact is protected by a protectionPolicy and can be restored.
// returns a RestorePointable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RestorePointable, error) {
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
func (m *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder when successful
func (m *BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder) {
    return NewBackupRestoreExchangeRestoreSessionsItemMailboxRestoreArtifactsItemRestorePointRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
