package tenantrelationships

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder provides operations to manage the serviceManagementDetails property of the microsoft.graph.delegatedAdminCustomer entity.
type DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetQueryParameters contains the management details of a service in the customer tenant that's managed by delegated administration.
type DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetQueryParameters
}
// DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderInternal instantiates a new DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder and sets the default values.
func NewDelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) {
    m := &DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/tenantRelationships/delegatedAdminCustomers/{delegatedAdminCustomer%2Did}/serviceManagementDetails/{delegatedAdminServiceManagementDetail%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder instantiates a new DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder and sets the default values.
func NewDelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property serviceManagementDetails for tenantRelationships
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
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
// Get contains the management details of a service in the customer tenant that's managed by delegated administration.
// returns a DelegatedAdminServiceManagementDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminServiceManagementDetailable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDelegatedAdminServiceManagementDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminServiceManagementDetailable), nil
}
// Patch update the navigation property serviceManagementDetails in tenantRelationships
// returns a DelegatedAdminServiceManagementDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminServiceManagementDetailable, requestConfiguration *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminServiceManagementDetailable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDelegatedAdminServiceManagementDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminServiceManagementDetailable), nil
}
// ToDeleteRequestInformation delete navigation property serviceManagementDetails for tenantRelationships
// returns a *RequestInformation when successful
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation contains the management details of a service in the customer tenant that's managed by delegated administration.
// returns a *RequestInformation when successful
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property serviceManagementDetails in tenantRelationships
// returns a *RequestInformation when successful
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminServiceManagementDetailable, requestConfiguration *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder when successful
func (m *DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) WithUrl(rawUrl string)(*DelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder) {
    return NewDelegatedAdminCustomersItemServiceManagementDetailsDelegatedAdminServiceManagementDetailItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
