package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder provides operations to manage the requestor property of the microsoft.graph.accessPackageAssignmentRequest entity.
type EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetQueryParameters the subject who requested or, if a direct assignment, was assigned. Read-only. Nullable. Supports $expand.
type EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetQueryParameters
}
// NewEntitlementManagementAssignmentRequestsItemRequestorRequestBuilderInternal instantiates a new EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder and sets the default values.
func NewEntitlementManagementAssignmentRequestsItemRequestorRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder) {
    m := &EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/entitlementManagement/assignmentRequests/{accessPackageAssignmentRequest%2Did}/requestor{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewEntitlementManagementAssignmentRequestsItemRequestorRequestBuilder instantiates a new EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder and sets the default values.
func NewEntitlementManagementAssignmentRequestsItemRequestorRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementAssignmentRequestsItemRequestorRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the subject who requested or, if a direct assignment, was assigned. Read-only. Nullable. Supports $expand.
// returns a AccessPackageSubjectable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageSubjectable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageSubjectFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageSubjectable), nil
}
// ToGetRequestInformation the subject who requested or, if a direct assignment, was assigned. Read-only. Nullable. Supports $expand.
// returns a *RequestInformation when successful
func (m *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder when successful
func (m *EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder) WithUrl(rawUrl string)(*EntitlementManagementAssignmentRequestsItemRequestorRequestBuilder) {
    return NewEntitlementManagementAssignmentRequestsItemRequestorRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
