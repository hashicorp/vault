package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EntitlementManagementAssignmentsAdditionalAccessRequestBuilder provides operations to call the additionalAccess method.
type EntitlementManagementAssignmentsAdditionalAccessRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetQueryParameters in Microsoft Entra Entitlement Management, retrieve a collection of accessPackageAssignment objects that indicate a target user has an assignment to a specified access package and also an assignment to another, potentially incompatible, access package.  This can be used to prepare to configure the incompatible access packages for a specific access package.
type EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetQueryParameters
}
// NewEntitlementManagementAssignmentsAdditionalAccessRequestBuilderInternal instantiates a new EntitlementManagementAssignmentsAdditionalAccessRequestBuilder and sets the default values.
func NewEntitlementManagementAssignmentsAdditionalAccessRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) {
    m := &EntitlementManagementAssignmentsAdditionalAccessRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/entitlementManagement/assignments/additionalAccess(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewEntitlementManagementAssignmentsAdditionalAccessRequestBuilder instantiates a new EntitlementManagementAssignmentsAdditionalAccessRequestBuilder and sets the default values.
func NewEntitlementManagementAssignmentsAdditionalAccessRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementAssignmentsAdditionalAccessRequestBuilderInternal(urlParams, requestAdapter)
}
// Get in Microsoft Entra Entitlement Management, retrieve a collection of accessPackageAssignment objects that indicate a target user has an assignment to a specified access package and also an assignment to another, potentially incompatible, access package.  This can be used to prepare to configure the incompatible access packages for a specific access package.
// Deprecated: This method is obsolete. Use GetAsAdditionalAccessGetResponse instead.
// returns a EntitlementManagementAssignmentsAdditionalAccessResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accesspackageassignment-additionalaccess?view=graph-rest-1.0
func (m *EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetRequestConfiguration)(EntitlementManagementAssignmentsAdditionalAccessResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateEntitlementManagementAssignmentsAdditionalAccessResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(EntitlementManagementAssignmentsAdditionalAccessResponseable), nil
}
// GetAsAdditionalAccessGetResponse in Microsoft Entra Entitlement Management, retrieve a collection of accessPackageAssignment objects that indicate a target user has an assignment to a specified access package and also an assignment to another, potentially incompatible, access package.  This can be used to prepare to configure the incompatible access packages for a specific access package.
// returns a EntitlementManagementAssignmentsAdditionalAccessGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accesspackageassignment-additionalaccess?view=graph-rest-1.0
func (m *EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) GetAsAdditionalAccessGetResponse(ctx context.Context, requestConfiguration *EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetRequestConfiguration)(EntitlementManagementAssignmentsAdditionalAccessGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateEntitlementManagementAssignmentsAdditionalAccessGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(EntitlementManagementAssignmentsAdditionalAccessGetResponseable), nil
}
// ToGetRequestInformation in Microsoft Entra Entitlement Management, retrieve a collection of accessPackageAssignment objects that indicate a target user has an assignment to a specified access package and also an assignment to another, potentially incompatible, access package.  This can be used to prepare to configure the incompatible access packages for a specific access package.
// returns a *RequestInformation when successful
func (m *EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementAssignmentsAdditionalAccessRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EntitlementManagementAssignmentsAdditionalAccessRequestBuilder when successful
func (m *EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) WithUrl(rawUrl string)(*EntitlementManagementAssignmentsAdditionalAccessRequestBuilder) {
    return NewEntitlementManagementAssignmentsAdditionalAccessRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
