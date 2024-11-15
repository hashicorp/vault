package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// RoleManagementPolicyAssignmentsItemPolicyRequestBuilder provides operations to manage the policy property of the microsoft.graph.unifiedRoleManagementPolicyAssignment entity.
type RoleManagementPolicyAssignmentsItemPolicyRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetQueryParameters the policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
type RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetQueryParameters
}
// NewRoleManagementPolicyAssignmentsItemPolicyRequestBuilderInternal instantiates a new RoleManagementPolicyAssignmentsItemPolicyRequestBuilder and sets the default values.
func NewRoleManagementPolicyAssignmentsItemPolicyRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RoleManagementPolicyAssignmentsItemPolicyRequestBuilder) {
    m := &RoleManagementPolicyAssignmentsItemPolicyRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies/roleManagementPolicyAssignments/{unifiedRoleManagementPolicyAssignment%2Did}/policy{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewRoleManagementPolicyAssignmentsItemPolicyRequestBuilder instantiates a new RoleManagementPolicyAssignmentsItemPolicyRequestBuilder and sets the default values.
func NewRoleManagementPolicyAssignmentsItemPolicyRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RoleManagementPolicyAssignmentsItemPolicyRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewRoleManagementPolicyAssignmentsItemPolicyRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
// returns a UnifiedRoleManagementPolicyable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *RoleManagementPolicyAssignmentsItemPolicyRequestBuilder) Get(ctx context.Context, requestConfiguration *RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRoleManagementPolicyable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUnifiedRoleManagementPolicyFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UnifiedRoleManagementPolicyable), nil
}
// ToGetRequestInformation the policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
// returns a *RequestInformation when successful
func (m *RoleManagementPolicyAssignmentsItemPolicyRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *RoleManagementPolicyAssignmentsItemPolicyRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *RoleManagementPolicyAssignmentsItemPolicyRequestBuilder when successful
func (m *RoleManagementPolicyAssignmentsItemPolicyRequestBuilder) WithUrl(rawUrl string)(*RoleManagementPolicyAssignmentsItemPolicyRequestBuilder) {
    return NewRoleManagementPolicyAssignmentsItemPolicyRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
