package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder provides operations to manage the granularMailboxRestoreArtifacts property of the microsoft.graph.exchangeRestoreSession entity.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetQueryParameters get granularMailboxRestoreArtifacts from solutions
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetQueryParameters struct {
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
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetQueryParameters
}
// BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByGranularMailboxRestoreArtifactId provides operations to manage the granularMailboxRestoreArtifacts property of the microsoft.graph.exchangeRestoreSession entity.
// returns a *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) ByGranularMailboxRestoreArtifactId(granularMailboxRestoreArtifactId string)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if granularMailboxRestoreArtifactId != "" {
        urlTplParams["granularMailboxRestoreArtifact%2Did"] = granularMailboxRestoreArtifactId
    }
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsGranularMailboxRestoreArtifactItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderInternal instantiates a new BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder and sets the default values.
func NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) {
    m := &BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/exchangeRestoreSessions/{exchangeRestoreSession%2Did}/granularMailboxRestoreArtifacts{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder instantiates a new BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder and sets the default values.
func NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsCountRequestBuilder when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) Count()(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsCountRequestBuilder) {
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get granularMailboxRestoreArtifacts from solutions
// returns a GranularMailboxRestoreArtifactCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGranularMailboxRestoreArtifactCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactCollectionResponseable), nil
}
// Post create new navigation property to granularMailboxRestoreArtifacts for solutions
// returns a GranularMailboxRestoreArtifactable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
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
// ToGetRequestInformation get granularMailboxRestoreArtifacts from solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to granularMailboxRestoreArtifacts for solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GranularMailboxRestoreArtifactable, requestConfiguration *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder when successful
func (m *BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder) {
    return NewBackupRestoreExchangeRestoreSessionsItemGranularMailboxRestoreArtifactsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
