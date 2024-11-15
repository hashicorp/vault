package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreServiceAppsItemDeactivateRequestBuilder provides operations to call the deactivate method.
type BackupRestoreServiceAppsItemDeactivateRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreServiceAppsItemDeactivateRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreServiceAppsItemDeactivateRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewBackupRestoreServiceAppsItemDeactivateRequestBuilderInternal instantiates a new BackupRestoreServiceAppsItemDeactivateRequestBuilder and sets the default values.
func NewBackupRestoreServiceAppsItemDeactivateRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreServiceAppsItemDeactivateRequestBuilder) {
    m := &BackupRestoreServiceAppsItemDeactivateRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/serviceApps/{serviceApp%2Did}/deactivate", pathParameters),
    }
    return m
}
// NewBackupRestoreServiceAppsItemDeactivateRequestBuilder instantiates a new BackupRestoreServiceAppsItemDeactivateRequestBuilder and sets the default values.
func NewBackupRestoreServiceAppsItemDeactivateRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreServiceAppsItemDeactivateRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreServiceAppsItemDeactivateRequestBuilderInternal(urlParams, requestAdapter)
}
// Post deactivate a serviceApp.
// returns a ServiceAppable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceapp-deactivate?view=graph-rest-1.0
func (m *BackupRestoreServiceAppsItemDeactivateRequestBuilder) Post(ctx context.Context, requestConfiguration *BackupRestoreServiceAppsItemDeactivateRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceAppable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateServiceAppFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceAppable), nil
}
// ToPostRequestInformation deactivate a serviceApp.
// returns a *RequestInformation when successful
func (m *BackupRestoreServiceAppsItemDeactivateRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreServiceAppsItemDeactivateRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *BackupRestoreServiceAppsItemDeactivateRequestBuilder when successful
func (m *BackupRestoreServiceAppsItemDeactivateRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreServiceAppsItemDeactivateRequestBuilder) {
    return NewBackupRestoreServiceAppsItemDeactivateRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
