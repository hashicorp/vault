package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder provides operations to manage the directories property of the microsoft.graph.synchronizationSchema entity.
type ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetQueryParameters contains the collection of directories and all of their objects.
type ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetQueryParameters
}
// ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderInternal instantiates a new ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) {
    m := &ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/synchronization/templates/{synchronizationTemplate%2Did}/schema/directories/{directoryDefinition%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder instantiates a new ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property directories for applications
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Discover provides operations to call the discover method.
// returns a *ItemSynchronizationTemplatesItemSchemaDirectoriesItemDiscoverRequestBuilder when successful
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) Discover()(*ItemSynchronizationTemplatesItemSchemaDirectoriesItemDiscoverRequestBuilder) {
    return NewItemSynchronizationTemplatesItemSchemaDirectoriesItemDiscoverRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get contains the collection of directories and all of their objects.
// returns a DirectoryDefinitionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryDefinitionable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryDefinitionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryDefinitionable), nil
}
// Patch update the navigation property directories in applications
// returns a DirectoryDefinitionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryDefinitionable, requestConfiguration *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryDefinitionable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryDefinitionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryDefinitionable), nil
}
// ToDeleteRequestInformation delete navigation property directories for applications
// returns a *RequestInformation when successful
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation contains the collection of directories and all of their objects.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property directories in applications
// returns a *RequestInformation when successful
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryDefinitionable, requestConfiguration *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder when successful
func (m *ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder) {
    return NewItemSynchronizationTemplatesItemSchemaDirectoriesDirectoryDefinitionItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
