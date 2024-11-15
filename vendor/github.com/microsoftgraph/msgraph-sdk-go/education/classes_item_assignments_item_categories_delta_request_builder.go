package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder provides operations to call the delta method.
type ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetQueryParameters get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
type ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetQueryParameters struct {
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
// ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetQueryParameters
}
// NewClassesItemAssignmentsItemCategoriesDeltaRequestBuilderInternal instantiates a new ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder and sets the default values.
func NewClassesItemAssignmentsItemCategoriesDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) {
    m := &ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/classes/{educationClass%2Did}/assignments/{educationAssignment%2Did}/categories/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewClassesItemAssignmentsItemCategoriesDeltaRequestBuilder instantiates a new ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder and sets the default values.
func NewClassesItemAssignmentsItemCategoriesDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewClassesItemAssignmentsItemCategoriesDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ClassesItemAssignmentsItemCategoriesDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationcategory-delta?view=graph-rest-1.0
func (m *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration)(ClassesItemAssignmentsItemCategoriesDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateClassesItemAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ClassesItemAssignmentsItemCategoriesDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
// returns a ClassesItemAssignmentsItemCategoriesDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationcategory-delta?view=graph-rest-1.0
func (m *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration)(ClassesItemAssignmentsItemCategoriesDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateClassesItemAssignmentsItemCategoriesDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ClassesItemAssignmentsItemCategoriesDeltaGetResponseable), nil
}
// ToGetRequestInformation get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
// returns a *RequestInformation when successful
func (m *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder when successful
func (m *ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) WithUrl(rawUrl string)(*ClassesItemAssignmentsItemCategoriesDeltaRequestBuilder) {
    return NewClassesItemAssignmentsItemCategoriesDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
