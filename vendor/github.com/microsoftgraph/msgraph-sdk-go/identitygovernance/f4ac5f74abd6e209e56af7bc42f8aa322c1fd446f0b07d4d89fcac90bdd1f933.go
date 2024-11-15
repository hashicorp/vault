package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder provides operations to call the resume method.
type LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewLifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderInternal instantiates a new LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder and sets the default values.
func NewLifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder) {
    m := &LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/lifecycleWorkflows/workflows/{workflow%2Did}/userProcessingResults/{userProcessingResult%2Did}/taskProcessingResults/{taskProcessingResult%2Did}/microsoft.graph.identityGovernance.resume", pathParameters),
    }
    return m
}
// NewLifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder instantiates a new LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder and sets the default values.
func NewLifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderInternal(urlParams, requestAdapter)
}
// Post resume a task processing result that's inProgress. In the default case an Azure Logic Apps system-assigned managed identity calls this API. For more information, see: Lifecycle Workflows extensibility approach.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/identitygovernance-taskprocessingresult-resume?view=graph-rest-1.0
func (m *LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder) Post(ctx context.Context, body LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeResumePostRequestBodyable, requestConfiguration *LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
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
// ToPostRequestInformation resume a task processing result that's inProgress. In the default case an Azure Logic Apps system-assigned managed identity calls this API. For more information, see: Lifecycle Workflows extensibility approach.
// returns a *RequestInformation when successful
func (m *LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder) ToPostRequestInformation(ctx context.Context, body LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeResumePostRequestBodyable, requestConfiguration *LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder when successful
func (m *LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder) WithUrl(rawUrl string)(*LifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder) {
    return NewLifecycleWorkflowsWorkflowsItemUserProcessingResultsItemTaskProcessingResultsItemMicrosoftGraphIdentityGovernanceResumeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
