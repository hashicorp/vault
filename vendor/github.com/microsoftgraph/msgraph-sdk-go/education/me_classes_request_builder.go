package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MeClassesRequestBuilder provides operations to manage the classes property of the microsoft.graph.educationUser entity.
type MeClassesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MeClassesRequestBuilderGetQueryParameters get the educationClass resources an educationUser is a member of.
type MeClassesRequestBuilderGetQueryParameters struct {
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
// MeClassesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeClassesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MeClassesRequestBuilderGetQueryParameters
}
// ByEducationClassId provides operations to manage the classes property of the microsoft.graph.educationUser entity.
// returns a *MeClassesEducationClassItemRequestBuilder when successful
func (m *MeClassesRequestBuilder) ByEducationClassId(educationClassId string)(*MeClassesEducationClassItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if educationClassId != "" {
        urlTplParams["educationClass%2Did"] = educationClassId
    }
    return NewMeClassesEducationClassItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewMeClassesRequestBuilderInternal instantiates a new MeClassesRequestBuilder and sets the default values.
func NewMeClassesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeClassesRequestBuilder) {
    m := &MeClassesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/me/classes{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMeClassesRequestBuilder instantiates a new MeClassesRequestBuilder and sets the default values.
func NewMeClassesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeClassesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeClassesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *MeClassesCountRequestBuilder when successful
func (m *MeClassesRequestBuilder) Count()(*MeClassesCountRequestBuilder) {
    return NewMeClassesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get the educationClass resources an educationUser is a member of.
// returns a EducationClassCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationuser-list-classes?view=graph-rest-1.0
func (m *MeClassesRequestBuilder) Get(ctx context.Context, requestConfiguration *MeClassesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationClassCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationClassCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationClassCollectionResponseable), nil
}
// ToGetRequestInformation get the educationClass resources an educationUser is a member of.
// returns a *RequestInformation when successful
func (m *MeClassesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MeClassesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MeClassesRequestBuilder when successful
func (m *MeClassesRequestBuilder) WithUrl(rawUrl string)(*MeClassesRequestBuilder) {
    return NewMeClassesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
