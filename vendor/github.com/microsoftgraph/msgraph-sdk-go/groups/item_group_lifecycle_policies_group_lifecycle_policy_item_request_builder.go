package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder provides operations to manage the groupLifecyclePolicies property of the microsoft.graph.group entity.
type ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetQueryParameters the collection of lifecycle policies for this group. Read-only. Nullable.
type ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetQueryParameters
}
// ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AddGroup provides operations to call the addGroup method.
// returns a *ItemGroupLifecyclePoliciesItemAddGroupRequestBuilder when successful
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) AddGroup()(*ItemGroupLifecyclePoliciesItemAddGroupRequestBuilder) {
    return NewItemGroupLifecyclePoliciesItemAddGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderInternal instantiates a new ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder and sets the default values.
func NewItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) {
    m := &ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/groupLifecyclePolicies/{groupLifecyclePolicy%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder instantiates a new ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder and sets the default values.
func NewItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property groupLifecyclePolicies for groups
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the collection of lifecycle policies for this group. Read-only. Nullable.
// returns a GroupLifecyclePolicyable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GroupLifecyclePolicyable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGroupLifecyclePolicyFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GroupLifecyclePolicyable), nil
}
// Patch update the navigation property groupLifecyclePolicies in groups
// returns a GroupLifecyclePolicyable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GroupLifecyclePolicyable, requestConfiguration *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GroupLifecyclePolicyable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGroupLifecyclePolicyFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GroupLifecyclePolicyable), nil
}
// RemoveGroup provides operations to call the removeGroup method.
// returns a *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder when successful
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) RemoveGroup()(*ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) {
    return NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property groupLifecyclePolicies for groups
// returns a *RequestInformation when successful
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the collection of lifecycle policies for this group. Read-only. Nullable.
// returns a *RequestInformation when successful
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property groupLifecyclePolicies in groups
// returns a *RequestInformation when successful
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.GroupLifecyclePolicyable, requestConfiguration *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder when successful
func (m *ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) WithUrl(rawUrl string)(*ItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder) {
    return NewItemGroupLifecyclePoliciesGroupLifecyclePolicyItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
