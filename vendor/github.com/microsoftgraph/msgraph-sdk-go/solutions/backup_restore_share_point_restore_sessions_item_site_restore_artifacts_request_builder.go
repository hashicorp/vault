package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder provides operations to manage the siteRestoreArtifacts property of the microsoft.graph.sharePointRestoreSession entity.
type BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetQueryParameters list all the siteRestoreArtifact objects for a sharePointRestoreSession for the tenant.
type BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetQueryParameters struct {
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
// BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetQueryParameters
}
// BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BySiteRestoreArtifactId provides operations to manage the siteRestoreArtifacts property of the microsoft.graph.sharePointRestoreSession entity.
// returns a *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsSiteRestoreArtifactItemRequestBuilder when successful
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) BySiteRestoreArtifactId(siteRestoreArtifactId string)(*BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsSiteRestoreArtifactItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if siteRestoreArtifactId != "" {
        urlTplParams["siteRestoreArtifact%2Did"] = siteRestoreArtifactId
    }
    return NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsSiteRestoreArtifactItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderInternal instantiates a new BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder and sets the default values.
func NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) {
    m := &BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/sharePointRestoreSessions/{sharePointRestoreSession%2Did}/siteRestoreArtifacts{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder instantiates a new BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder and sets the default values.
func NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsCountRequestBuilder when successful
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) Count()(*BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsCountRequestBuilder) {
    return NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get list all the siteRestoreArtifact objects for a sharePointRestoreSession for the tenant.
// returns a SiteRestoreArtifactCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/sharepointrestoresession-list-siterestoreartifacts?view=graph-rest-1.0
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteRestoreArtifactCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSiteRestoreArtifactCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteRestoreArtifactCollectionResponseable), nil
}
// Post create new navigation property to siteRestoreArtifacts for solutions
// returns a SiteRestoreArtifactable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteRestoreArtifactable, requestConfiguration *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteRestoreArtifactable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSiteRestoreArtifactFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteRestoreArtifactable), nil
}
// ToGetRequestInformation list all the siteRestoreArtifact objects for a sharePointRestoreSession for the tenant.
// returns a *RequestInformation when successful
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to siteRestoreArtifacts for solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteRestoreArtifactable, requestConfiguration *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder when successful
func (m *BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder) {
    return NewBackupRestoreSharePointRestoreSessionsItemSiteRestoreArtifactsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
