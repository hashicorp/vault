package rolemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder provides operations to manage the resourceActions property of the microsoft.graph.unifiedRbacResourceNamespace entity.
type DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetQueryParameters get resourceActions from roleManagement
type DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetQueryParameters
}
// DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderInternal instantiates a new DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder and sets the default values.
func NewDirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) {
    m := &DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/roleManagement/directory/resourceNamespaces/{unifiedRbacResourceNamespace%2Did}/resourceActions/{unifiedRbacResourceAction%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder instantiates a new DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder and sets the default values.
func NewDirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property resourceActions for roleManagement
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get resourceActions from roleManagement
// returns a UnifiedRbacResourceActionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRbacResourceActionable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUnifiedRbacResourceActionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRbacResourceActionable), nil
}
// Patch update the navigation property resourceActions in roleManagement
// returns a UnifiedRbacResourceActionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRbacResourceActionable, requestConfiguration *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRbacResourceActionable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUnifiedRbacResourceActionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRbacResourceActionable), nil
}
// ToDeleteRequestInformation delete navigation property resourceActions for roleManagement
// returns a *RequestInformation when successful
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get resourceActions from roleManagement
// returns a *RequestInformation when successful
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property resourceActions in roleManagement
// returns a *RequestInformation when successful
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRbacResourceActionable, requestConfiguration *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder when successful
func (m *DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) WithUrl(rawUrl string)(*DirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder) {
    return NewDirectoryResourceNamespacesItemResourceActionsUnifiedRbacResourceActionItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
