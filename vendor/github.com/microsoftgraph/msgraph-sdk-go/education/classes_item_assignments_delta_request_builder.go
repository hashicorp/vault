package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ClassesItemAssignmentsDeltaRequestBuilder provides operations to call the delta method.
type ClassesItemAssignmentsDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ClassesItemAssignmentsDeltaRequestBuilderGetQueryParameters get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
type ClassesItemAssignmentsDeltaRequestBuilderGetQueryParameters struct {
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
// ClassesItemAssignmentsDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ClassesItemAssignmentsDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ClassesItemAssignmentsDeltaRequestBuilderGetQueryParameters
}
// NewClassesItemAssignmentsDeltaRequestBuilderInternal instantiates a new ClassesItemAssignmentsDeltaRequestBuilder and sets the default values.
func NewClassesItemAssignmentsDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemAssignmentsDeltaRequestBuilder) {
    m := &ClassesItemAssignmentsDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/classes/{educationClass%2Did}/assignments/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewClassesItemAssignmentsDeltaRequestBuilder instantiates a new ClassesItemAssignmentsDeltaRequestBuilder and sets the default values.
func NewClassesItemAssignmentsDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemAssignmentsDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewClassesItemAssignmentsDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ClassesItemAssignmentsDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationassignment-delta?view=graph-rest-1.0
func (m *ClassesItemAssignmentsDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ClassesItemAssignmentsDeltaRequestBuilderGetRequestConfiguration)(ClassesItemAssignmentsDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateClassesItemAssignmentsDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ClassesItemAssignmentsDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
// returns a ClassesItemAssignmentsDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationassignment-delta?view=graph-rest-1.0
func (m *ClassesItemAssignmentsDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ClassesItemAssignmentsDeltaRequestBuilderGetRequestConfiguration)(ClassesItemAssignmentsDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateClassesItemAssignmentsDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ClassesItemAssignmentsDeltaGetResponseable), nil
}
// ToGetRequestInformation get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
// returns a *RequestInformation when successful
func (m *ClassesItemAssignmentsDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ClassesItemAssignmentsDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ClassesItemAssignmentsDeltaRequestBuilder when successful
func (m *ClassesItemAssignmentsDeltaRequestBuilder) WithUrl(rawUrl string)(*ClassesItemAssignmentsDeltaRequestBuilder) {
    return NewClassesItemAssignmentsDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
