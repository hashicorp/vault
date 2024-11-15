package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder provides operations to manage the delegatedPermissionClassifications property of the microsoft.graph.servicePrincipal entity.
type ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetQueryParameters get delegatedPermissionClassifications from servicePrincipals
type ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetQueryParameters
}
// ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderInternal instantiates a new ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder and sets the default values.
func NewItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) {
    m := &ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/delegatedPermissionClassifications/{delegatedPermissionClassification%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder instantiates a new ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder and sets the default values.
func NewItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete deletes a delegatedPermissionClassification which had previously been set for a delegated permission.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceprincipal-delete-delegatedpermissionclassifications?view=graph-rest-1.0
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get delegatedPermissionClassifications from servicePrincipals
// returns a DelegatedPermissionClassificationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedPermissionClassificationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDelegatedPermissionClassificationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedPermissionClassificationable), nil
}
// Patch update the navigation property delegatedPermissionClassifications in servicePrincipals
// returns a DelegatedPermissionClassificationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedPermissionClassificationable, requestConfiguration *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedPermissionClassificationable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDelegatedPermissionClassificationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedPermissionClassificationable), nil
}
// ToDeleteRequestInformation deletes a delegatedPermissionClassification which had previously been set for a delegated permission.
// returns a *RequestInformation when successful
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get delegatedPermissionClassifications from servicePrincipals
// returns a *RequestInformation when successful
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property delegatedPermissionClassifications in servicePrincipals
// returns a *RequestInformation when successful
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DelegatedPermissionClassificationable, requestConfiguration *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder when successful
func (m *ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) WithUrl(rawUrl string)(*ItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder) {
    return NewItemDelegatedPermissionClassificationsDelegatedPermissionClassificationItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
