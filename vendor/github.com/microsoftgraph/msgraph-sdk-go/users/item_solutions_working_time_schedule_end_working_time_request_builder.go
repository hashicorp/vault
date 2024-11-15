package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder provides operations to call the endWorkingTime method.
type ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderInternal instantiates a new ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder and sets the default values.
func NewItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder) {
    m := &ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/solutions/workingTimeSchedule/endWorkingTime", pathParameters),
    }
    return m
}
// NewItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder instantiates a new ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder and sets the default values.
func NewItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderInternal(urlParams, requestAdapter)
}
// Post trigger the policies associated with the end of working hours for a specific user.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/workingtimeschedule-endworkingtime?view=graph-rest-1.0
func (m *ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation trigger the policies associated with the end of working hours for a specific user.
// returns a *RequestInformation when successful
func (m *ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder when successful
func (m *ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder) WithUrl(rawUrl string)(*ItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder) {
    return NewItemSolutionsWorkingTimeScheduleEndWorkingTimeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
