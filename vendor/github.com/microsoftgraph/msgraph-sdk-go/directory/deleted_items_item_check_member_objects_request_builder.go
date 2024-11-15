package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeletedItemsItemCheckMemberObjectsRequestBuilder provides operations to call the checkMemberObjects method.
type DeletedItemsItemCheckMemberObjectsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeletedItemsItemCheckMemberObjectsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedItemsItemCheckMemberObjectsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDeletedItemsItemCheckMemberObjectsRequestBuilderInternal instantiates a new DeletedItemsItemCheckMemberObjectsRequestBuilder and sets the default values.
func NewDeletedItemsItemCheckMemberObjectsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsItemCheckMemberObjectsRequestBuilder) {
    m := &DeletedItemsItemCheckMemberObjectsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/deletedItems/{directoryObject%2Did}/checkMemberObjects", pathParameters),
    }
    return m
}
// NewDeletedItemsItemCheckMemberObjectsRequestBuilder instantiates a new DeletedItemsItemCheckMemberObjectsRequestBuilder and sets the default values.
func NewDeletedItemsItemCheckMemberObjectsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsItemCheckMemberObjectsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeletedItemsItemCheckMemberObjectsRequestBuilderInternal(urlParams, requestAdapter)
}
// Post invoke action checkMemberObjects
// Deprecated: This method is obsolete. Use PostAsCheckMemberObjectsPostResponse instead.
// returns a DeletedItemsItemCheckMemberObjectsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeletedItemsItemCheckMemberObjectsRequestBuilder) Post(ctx context.Context, body DeletedItemsItemCheckMemberObjectsPostRequestBodyable, requestConfiguration *DeletedItemsItemCheckMemberObjectsRequestBuilderPostRequestConfiguration)(DeletedItemsItemCheckMemberObjectsResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateDeletedItemsItemCheckMemberObjectsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(DeletedItemsItemCheckMemberObjectsResponseable), nil
}
// PostAsCheckMemberObjectsPostResponse invoke action checkMemberObjects
// returns a DeletedItemsItemCheckMemberObjectsPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeletedItemsItemCheckMemberObjectsRequestBuilder) PostAsCheckMemberObjectsPostResponse(ctx context.Context, body DeletedItemsItemCheckMemberObjectsPostRequestBodyable, requestConfiguration *DeletedItemsItemCheckMemberObjectsRequestBuilderPostRequestConfiguration)(DeletedItemsItemCheckMemberObjectsPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateDeletedItemsItemCheckMemberObjectsPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(DeletedItemsItemCheckMemberObjectsPostResponseable), nil
}
// ToPostRequestInformation invoke action checkMemberObjects
// returns a *RequestInformation when successful
func (m *DeletedItemsItemCheckMemberObjectsRequestBuilder) ToPostRequestInformation(ctx context.Context, body DeletedItemsItemCheckMemberObjectsPostRequestBodyable, requestConfiguration *DeletedItemsItemCheckMemberObjectsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *DeletedItemsItemCheckMemberObjectsRequestBuilder when successful
func (m *DeletedItemsItemCheckMemberObjectsRequestBuilder) WithUrl(rawUrl string)(*DeletedItemsItemCheckMemberObjectsRequestBuilder) {
    return NewDeletedItemsItemCheckMemberObjectsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
