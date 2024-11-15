package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MeAssignmentsItemCategoriesDeltaRequestBuilder provides operations to call the delta method.
type MeAssignmentsItemCategoriesDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MeAssignmentsItemCategoriesDeltaRequestBuilderGetQueryParameters get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
type MeAssignmentsItemCategoriesDeltaRequestBuilderGetQueryParameters struct {
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
// MeAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MeAssignmentsItemCategoriesDeltaRequestBuilderGetQueryParameters
}
// NewMeAssignmentsItemCategoriesDeltaRequestBuilderInternal instantiates a new MeAssignmentsItemCategoriesDeltaRequestBuilder and sets the default values.
func NewMeAssignmentsItemCategoriesDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemCategoriesDeltaRequestBuilder) {
    m := &MeAssignmentsItemCategoriesDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/me/assignments/{educationAssignment%2Did}/categories/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMeAssignmentsItemCategoriesDeltaRequestBuilder instantiates a new MeAssignmentsItemCategoriesDeltaRequestBuilder and sets the default values.
func NewMeAssignmentsItemCategoriesDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemCategoriesDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeAssignmentsItemCategoriesDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a MeAssignmentsItemCategoriesDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationcategory-delta?view=graph-rest-1.0
func (m *MeAssignmentsItemCategoriesDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *MeAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration)(MeAssignmentsItemCategoriesDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateMeAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(MeAssignmentsItemCategoriesDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
// returns a MeAssignmentsItemCategoriesDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationcategory-delta?view=graph-rest-1.0
func (m *MeAssignmentsItemCategoriesDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *MeAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration)(MeAssignmentsItemCategoriesDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateMeAssignmentsItemCategoriesDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(MeAssignmentsItemCategoriesDeltaGetResponseable), nil
}
// ToGetRequestInformation get a list of newly created or updated educationCategory objects without having to perform a full read of the collection.
// returns a *RequestInformation when successful
func (m *MeAssignmentsItemCategoriesDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MeAssignmentsItemCategoriesDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MeAssignmentsItemCategoriesDeltaRequestBuilder when successful
func (m *MeAssignmentsItemCategoriesDeltaRequestBuilder) WithUrl(rawUrl string)(*MeAssignmentsItemCategoriesDeltaRequestBuilder) {
    return NewMeAssignmentsItemCategoriesDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
