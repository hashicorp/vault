package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder provides operations to call the setOrder method.
type B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewB2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderInternal instantiates a new B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder and sets the default values.
func NewB2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder) {
    m := &B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/b2xUserFlows/{b2xIdentityUserFlow%2Did}/userAttributeAssignments/setOrder", pathParameters),
    }
    return m
}
// NewB2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder instantiates a new B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder and sets the default values.
func NewB2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewB2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderInternal(urlParams, requestAdapter)
}
// Post set the order of identityUserFlowAttributeAssignments being collected within a user flow.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/identityuserflowattributeassignment-setorder?view=graph-rest-1.0
func (m *B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder) Post(ctx context.Context, body B2xUserFlowsItemUserAttributeAssignmentsSetOrderPostRequestBodyable, requestConfiguration *B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation set the order of identityUserFlowAttributeAssignments being collected within a user flow.
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder) ToPostRequestInformation(ctx context.Context, body B2xUserFlowsItemUserAttributeAssignmentsSetOrderPostRequestBodyable, requestConfiguration *B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder when successful
func (m *B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder) WithUrl(rawUrl string)(*B2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder) {
    return NewB2xUserFlowsItemUserAttributeAssignmentsSetOrderRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
