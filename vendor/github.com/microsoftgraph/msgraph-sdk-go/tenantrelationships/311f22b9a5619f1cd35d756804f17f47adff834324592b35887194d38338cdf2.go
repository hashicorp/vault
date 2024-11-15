package tenantrelationships

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder provides operations to manage the operations property of the microsoft.graph.delegatedAdminRelationship entity.
type DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetQueryParameters read the properties of a delegatedAdminRelationshipOperation object.
type DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetQueryParameters
}
// DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderInternal instantiates a new DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder and sets the default values.
func NewDelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) {
    m := &DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/tenantRelationships/delegatedAdminRelationships/{delegatedAdminRelationship%2Did}/operations/{delegatedAdminRelationshipOperation%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder instantiates a new DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder and sets the default values.
func NewDelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property operations for tenantRelationships
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read the properties of a delegatedAdminRelationshipOperation object.
// returns a DelegatedAdminRelationshipOperationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/delegatedadminrelationshipoperation-get?view=graph-rest-1.0
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminRelationshipOperationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDelegatedAdminRelationshipOperationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminRelationshipOperationable), nil
}
// Patch update the navigation property operations in tenantRelationships
// returns a DelegatedAdminRelationshipOperationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminRelationshipOperationable, requestConfiguration *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminRelationshipOperationable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDelegatedAdminRelationshipOperationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminRelationshipOperationable), nil
}
// ToDeleteRequestInformation delete navigation property operations for tenantRelationships
// returns a *RequestInformation when successful
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties of a delegatedAdminRelationshipOperation object.
// returns a *RequestInformation when successful
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property operations in tenantRelationships
// returns a *RequestInformation when successful
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedAdminRelationshipOperationable, requestConfiguration *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder when successful
func (m *DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) WithUrl(rawUrl string)(*DelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder) {
    return NewDelegatedAdminRelationshipsItemOperationsDelegatedAdminRelationshipOperationItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
