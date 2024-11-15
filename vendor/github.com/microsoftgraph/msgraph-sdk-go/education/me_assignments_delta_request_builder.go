package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MeAssignmentsDeltaRequestBuilder provides operations to call the delta method.
type MeAssignmentsDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MeAssignmentsDeltaRequestBuilderGetQueryParameters get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
type MeAssignmentsDeltaRequestBuilderGetQueryParameters struct {
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
// MeAssignmentsDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeAssignmentsDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MeAssignmentsDeltaRequestBuilderGetQueryParameters
}
// NewMeAssignmentsDeltaRequestBuilderInternal instantiates a new MeAssignmentsDeltaRequestBuilder and sets the default values.
func NewMeAssignmentsDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsDeltaRequestBuilder) {
    m := &MeAssignmentsDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/me/assignments/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMeAssignmentsDeltaRequestBuilder instantiates a new MeAssignmentsDeltaRequestBuilder and sets the default values.
func NewMeAssignmentsDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeAssignmentsDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a MeAssignmentsDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationassignment-delta?view=graph-rest-1.0
func (m *MeAssignmentsDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *MeAssignmentsDeltaRequestBuilderGetRequestConfiguration)(MeAssignmentsDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateMeAssignmentsDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(MeAssignmentsDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
// returns a MeAssignmentsDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationassignment-delta?view=graph-rest-1.0
func (m *MeAssignmentsDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *MeAssignmentsDeltaRequestBuilderGetRequestConfiguration)(MeAssignmentsDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateMeAssignmentsDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(MeAssignmentsDeltaGetResponseable), nil
}
// ToGetRequestInformation get a list of newly-created or updated assignments without reading the whole collection. A teacher or an application running with application permissions can see all assignment objects for the class. Students can only see assignments that are assigned to them.
// returns a *RequestInformation when successful
func (m *MeAssignmentsDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MeAssignmentsDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MeAssignmentsDeltaRequestBuilder when successful
func (m *MeAssignmentsDeltaRequestBuilder) WithUrl(rawUrl string)(*MeAssignmentsDeltaRequestBuilder) {
    return NewMeAssignmentsDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
