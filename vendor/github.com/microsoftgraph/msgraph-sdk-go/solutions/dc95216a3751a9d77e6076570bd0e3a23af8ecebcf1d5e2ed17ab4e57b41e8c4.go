package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder provides operations to manage the granularMailboxRestoreArtifacts property of the microsoft.graph.exchangeRestoreSession entity.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetQueryParameters get granularMailboxRestoreArtifacts from solutions
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetQueryParameters
}
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderInternal instantiates a new BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder and sets the default values.
func NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) {
    m := &BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/exchangeRestoreSessions/{exchangeRestoreSession%2Did}/granularMailboxRestoreArtifacts/{granularMailboxRestoreArtifact%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder instantiates a new BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder and sets the default values.
func NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property granularMailboxRestoreArtifacts for solutions
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get granularMailboxRestoreArtifacts from solutions
// returns a GranularMailboxRestoreArtifactable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGranularMailboxRestoreArtifactFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable), nil
}
// Patch update the navigation property granularMailboxRestoreArtifacts in solutions
// returns a GranularMailboxRestoreArtifactable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGranularMailboxRestoreArtifactFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable), nil
}
// RestorePoint provides operations to manage the restorePoint property of the microsoft.graph.restoreArtifactBase entity.
// returns a *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsItemRestorePointRequestBuilder when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) RestorePoint()(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsItemRestorePointRequestBuilder) {
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsItemRestorePointRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property granularMailboxRestoreArtifacts for solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get granularMailboxRestoreArtifacts from solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property granularMailboxRestoreArtifacts in solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) {
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
