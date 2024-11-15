package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder provides operations to manage the instances property of the microsoft.graph.accessReviewHistoryDefinition entity.
type AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetQueryParameters if the accessReviewHistoryDefinition is a recurring definition, instances represent each recurrence. A definition that doesn't recur will have exactly one instance.
type AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetQueryParameters
}
// AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderInternal instantiates a new AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder and sets the default values.
func NewAccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) {
    m := &AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/accessReviews/historyDefinitions/{accessReviewHistoryDefinition%2Did}/instances/{accessReviewHistoryInstance%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewAccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder instantiates a new AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder and sets the default values.
func NewAccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property instances for identityGovernance
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// GenerateDownloadUri provides operations to call the generateDownloadUri method.
// returns a *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) GenerateDownloadUri()(*AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) {
    return NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get if the accessReviewHistoryDefinition is a recurring definition, instances represent each recurrence. A definition that doesn't recur will have exactly one instance.
// returns a AccessReviewHistoryInstanceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) Get(ctx context.Context, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessReviewHistoryInstanceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable), nil
}
// Patch update the navigation property instances in identityGovernance
// returns a AccessReviewHistoryInstanceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessReviewHistoryInstanceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable), nil
}
// ToDeleteRequestInformation delete navigation property instances for identityGovernance
// returns a *RequestInformation when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation if the accessReviewHistoryDefinition is a recurring definition, instances represent each recurrence. A definition that doesn't recur will have exactly one instance.
// returns a *RequestInformation when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property instances in identityGovernance
// returns a *RequestInformation when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) WithUrl(rawUrl string)(*AccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder) {
    return NewAccessReviewsHistoryDefinitionsItemInstancesAccessReviewHistoryInstanceItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
