package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430 "github.com/microsoftgraph/msgraph-sdk-go/models/identitygovernance"
)

// LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder provides operations to call the restore method.
type LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderInternal instantiates a new LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder and sets the default values.
func NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder) {
    m := &LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/lifecycleWorkflows/workflows/{workflow%2Did}/microsoft.graph.identityGovernance.restore", pathParameters),
    }
    return m
}
// NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder instantiates a new LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder and sets the default values.
func NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderInternal(urlParams, requestAdapter)
}
// Post restore a workflow that has been deleted. You can only restore a workflow that was deleted within the last 30 days before Microsoft Entra ID automatically permanently deletes it.
// returns a Workflowable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/identitygovernance-workflow-restore?view=graph-rest-1.0
func (m *LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder) Post(ctx context.Context, requestConfiguration *LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderPostRequestConfiguration)(ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.Workflowable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.CreateWorkflowFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ibf6ed4fc8e373ed2600905053a507c004671ad1749cb4b6b77078a908490c430.Workflowable), nil
}
// ToPostRequestInformation restore a workflow that has been deleted. You can only restore a workflow that was deleted within the last 30 days before Microsoft Entra ID automatically permanently deletes it.
// returns a *RequestInformation when successful
func (m *LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder when successful
func (m *LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder) WithUrl(rawUrl string)(*LifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder) {
    return NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceRestoreRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
