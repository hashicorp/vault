package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder provides operations to call the startWorkingTime method.
type ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderInternal instantiates a new ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder and sets the default values.
func NewItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder) {
    m := &ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/solutions/workingTimeSchedule/startWorkingTime", pathParameters),
    }
    return m
}
// NewItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder instantiates a new ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder and sets the default values.
func NewItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderInternal(urlParams, requestAdapter)
}
// Post trigger the policies associated with the start of working hours for a specific user.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/workingtimeschedule-startworkingtime?view=graph-rest-1.0
func (m *ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation trigger the policies associated with the start of working hours for a specific user.
// returns a *RequestInformation when successful
func (m *ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder when successful
func (m *ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder) WithUrl(rawUrl string)(*ItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder) {
    return NewItemSolutionsWorkingTimeScheduleStartWorkingTimeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
