package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreRestorePointsSearchRequestBuilder provides operations to call the search method.
type BackupRestoreRestorePointsSearchRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreRestorePointsSearchRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreRestorePointsSearchRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewBackupRestoreRestorePointsSearchRequestBuilderInternal instantiates a new BackupRestoreRestorePointsSearchRequestBuilder and sets the default values.
func NewBackupRestoreRestorePointsSearchRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreRestorePointsSearchRequestBuilder) {
    m := &BackupRestoreRestorePointsSearchRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/restorePoints/search", pathParameters),
    }
    return m
}
// NewBackupRestoreRestorePointsSearchRequestBuilder instantiates a new BackupRestoreRestorePointsSearchRequestBuilder and sets the default values.
func NewBackupRestoreRestorePointsSearchRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreRestorePointsSearchRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreRestorePointsSearchRequestBuilderInternal(urlParams, requestAdapter)
}
// Post search for the restorePoint objects associated with a protectionUnit.
// returns a RestorePointSearchResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/restorepoint-search?view=graph-rest-1.0
func (m *BackupRestoreRestorePointsSearchRequestBuilder) Post(ctx context.Context, body BackupRestoreRestorePointsSearchPostRequestBodyable, requestConfiguration *BackupRestoreRestorePointsSearchRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RestorePointSearchResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateRestorePointSearchResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RestorePointSearchResponseable), nil
}
// ToPostRequestInformation search for the restorePoint objects associated with a protectionUnit.
// returns a *RequestInformation when successful
func (m *BackupRestoreRestorePointsSearchRequestBuilder) ToPostRequestInformation(ctx context.Context, body BackupRestoreRestorePointsSearchPostRequestBodyable, requestConfiguration *BackupRestoreRestorePointsSearchRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *BackupRestoreRestorePointsSearchRequestBuilder when successful
func (m *BackupRestoreRestorePointsSearchRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreRestorePointsSearchRequestBuilder) {
    return NewBackupRestoreRestorePointsSearchRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
