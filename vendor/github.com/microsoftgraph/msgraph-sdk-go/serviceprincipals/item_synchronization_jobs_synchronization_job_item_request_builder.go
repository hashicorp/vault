package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationJobsSynchronizationJobItemRequestBuilder provides operations to manage the jobs property of the microsoft.graph.synchronization entity.
type ItemSynchronizationJobsSynchronizationJobItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationJobsSynchronizationJobItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationJobsSynchronizationJobItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetQueryParameters retrieve the existing synchronization job and its properties.
type ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetQueryParameters
}
// ItemSynchronizationJobsSynchronizationJobItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationJobsSynchronizationJobItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BulkUpload provides operations to manage the bulkUpload property of the microsoft.graph.synchronizationJob entity.
// returns a *ItemSynchronizationJobsItemBulkUploadRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) BulkUpload()(*ItemSynchronizationJobsItemBulkUploadRequestBuilder) {
    return NewItemSynchronizationJobsItemBulkUploadRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemSynchronizationJobsSynchronizationJobItemRequestBuilderInternal instantiates a new ItemSynchronizationJobsSynchronizationJobItemRequestBuilder and sets the default values.
func NewItemSynchronizationJobsSynchronizationJobItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) {
    m := &ItemSynchronizationJobsSynchronizationJobItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/synchronization/jobs/{synchronizationJob%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemSynchronizationJobsSynchronizationJobItemRequestBuilder instantiates a new ItemSynchronizationJobsSynchronizationJobItemRequestBuilder and sets the default values.
func NewItemSynchronizationJobsSynchronizationJobItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationJobsSynchronizationJobItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete stop the synchronization job, and permanently delete all the state associated with it. Synchronized accounts are left as-is.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationjob-delete?view=graph-rest-1.0
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve the existing synchronization job and its properties.
// returns a SynchronizationJobable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationjob-get?view=graph-rest-1.0
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SynchronizationJobable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSynchronizationJobFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SynchronizationJobable), nil
}
// Patch update the navigation property jobs in servicePrincipals
// returns a SynchronizationJobable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SynchronizationJobable, requestConfiguration *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SynchronizationJobable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSynchronizationJobFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SynchronizationJobable), nil
}
// Pause provides operations to call the pause method.
// returns a *ItemSynchronizationJobsItemPauseRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Pause()(*ItemSynchronizationJobsItemPauseRequestBuilder) {
    return NewItemSynchronizationJobsItemPauseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ProvisionOnDemand provides operations to call the provisionOnDemand method.
// returns a *ItemSynchronizationJobsItemProvisionOnDemandRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) ProvisionOnDemand()(*ItemSynchronizationJobsItemProvisionOnDemandRequestBuilder) {
    return NewItemSynchronizationJobsItemProvisionOnDemandRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Restart provides operations to call the restart method.
// returns a *ItemSynchronizationJobsItemRestartRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Restart()(*ItemSynchronizationJobsItemRestartRequestBuilder) {
    return NewItemSynchronizationJobsItemRestartRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Schema provides operations to manage the schema property of the microsoft.graph.synchronizationJob entity.
// returns a *ItemSynchronizationJobsItemSchemaRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Schema()(*ItemSynchronizationJobsItemSchemaRequestBuilder) {
    return NewItemSynchronizationJobsItemSchemaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Start provides operations to call the start method.
// returns a *ItemSynchronizationJobsItemStartRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) Start()(*ItemSynchronizationJobsItemStartRequestBuilder) {
    return NewItemSynchronizationJobsItemStartRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation stop the synchronization job, and permanently delete all the state associated with it. Synchronized accounts are left as-is.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation retrieve the existing synchronization job and its properties.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property jobs in servicePrincipals
// returns a *RequestInformation when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SynchronizationJobable, requestConfiguration *ItemSynchronizationJobsSynchronizationJobItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ValidateCredentials provides operations to call the validateCredentials method.
// returns a *ItemSynchronizationJobsItemValidateCredentialsRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) ValidateCredentials()(*ItemSynchronizationJobsItemValidateCredentialsRequestBuilder) {
    return NewItemSynchronizationJobsItemValidateCredentialsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder when successful
func (m *ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationJobsSynchronizationJobItemRequestBuilder) {
    return NewItemSynchronizationJobsSynchronizationJobItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
