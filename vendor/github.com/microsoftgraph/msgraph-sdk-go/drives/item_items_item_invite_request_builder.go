package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemInviteRequestBuilder provides operations to call the invite method.
type ItemItemsItemInviteRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemInviteRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemInviteRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemInviteRequestBuilderInternal instantiates a new ItemItemsItemInviteRequestBuilder and sets the default values.
func NewItemItemsItemInviteRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemInviteRequestBuilder) {
    m := &ItemItemsItemInviteRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/invite", pathParameters),
    }
    return m
}
// NewItemItemsItemInviteRequestBuilder instantiates a new ItemItemsItemInviteRequestBuilder and sets the default values.
func NewItemItemsItemInviteRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemInviteRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemInviteRequestBuilderInternal(urlParams, requestAdapter)
}
// Post sends a sharing invitation for a driveItem.A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.
// Deprecated: This method is obsolete. Use PostAsInvitePostResponse instead.
// returns a ItemItemsItemInviteResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/driveitem-invite?view=graph-rest-1.0
func (m *ItemItemsItemInviteRequestBuilder) Post(ctx context.Context, body ItemItemsItemInvitePostRequestBodyable, requestConfiguration *ItemItemsItemInviteRequestBuilderPostRequestConfiguration)(ItemItemsItemInviteResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemItemsItemInviteResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemItemsItemInviteResponseable), nil
}
// PostAsInvitePostResponse sends a sharing invitation for a driveItem.A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.
// returns a ItemItemsItemInvitePostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/driveitem-invite?view=graph-rest-1.0
func (m *ItemItemsItemInviteRequestBuilder) PostAsInvitePostResponse(ctx context.Context, body ItemItemsItemInvitePostRequestBodyable, requestConfiguration *ItemItemsItemInviteRequestBuilderPostRequestConfiguration)(ItemItemsItemInvitePostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemItemsItemInvitePostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemItemsItemInvitePostResponseable), nil
}
// ToPostRequestInformation sends a sharing invitation for a driveItem.A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.
// returns a *RequestInformation when successful
func (m *ItemItemsItemInviteRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemItemsItemInvitePostRequestBodyable, requestConfiguration *ItemItemsItemInviteRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemItemsItemInviteRequestBuilder when successful
func (m *ItemItemsItemInviteRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemInviteRequestBuilder) {
    return NewItemItemsItemInviteRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
