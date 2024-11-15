package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder provides operations to call the removeGroup method.
type ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderInternal instantiates a new ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder and sets the default values.
func NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) {
    m := &ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/groupLifecyclePolicies/{groupLifecyclePolicy%2Did}/removeGroup", pathParameters),
    }
    return m
}
// NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder instantiates a new ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder and sets the default values.
func NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderInternal(urlParams, requestAdapter)
}
// Post removes a group from a lifecycle policy.
// Deprecated: This method is obsolete. Use PostAsRemoveGroupPostResponse instead.
// returns a ItemGroupLifecyclePoliciesItemRemoveGroupResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/grouplifecyclepolicy-removegroup?view=graph-rest-1.0
func (m *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) Post(ctx context.Context, body ItemGroupLifecyclePoliciesItemRemoveGroupPostRequestBodyable, requestConfiguration *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderPostRequestConfiguration)(ItemGroupLifecyclePoliciesItemRemoveGroupResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemGroupLifecyclePoliciesItemRemoveGroupResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemGroupLifecyclePoliciesItemRemoveGroupResponseable), nil
}
// PostAsRemoveGroupPostResponse removes a group from a lifecycle policy.
// returns a ItemGroupLifecyclePoliciesItemRemoveGroupPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/grouplifecyclepolicy-removegroup?view=graph-rest-1.0
func (m *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) PostAsRemoveGroupPostResponse(ctx context.Context, body ItemGroupLifecyclePoliciesItemRemoveGroupPostRequestBodyable, requestConfiguration *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderPostRequestConfiguration)(ItemGroupLifecyclePoliciesItemRemoveGroupPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemGroupLifecyclePoliciesItemRemoveGroupPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemGroupLifecyclePoliciesItemRemoveGroupPostResponseable), nil
}
// ToPostRequestInformation removes a group from a lifecycle policy.
// returns a *RequestInformation when successful
func (m *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemGroupLifecyclePoliciesItemRemoveGroupPostRequestBodyable, requestConfiguration *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder when successful
func (m *ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) WithUrl(rawUrl string)(*ItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder) {
    return NewItemGroupLifecyclePoliciesItemRemoveGroupRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
