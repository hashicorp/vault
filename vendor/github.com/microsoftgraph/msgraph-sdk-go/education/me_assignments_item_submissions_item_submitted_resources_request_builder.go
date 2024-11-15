package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder provides operations to manage the submittedResources property of the microsoft.graph.educationSubmission entity.
type MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetQueryParameters get submittedResources from education
type MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetQueryParameters struct {
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
// MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetQueryParameters
}
// MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByEducationSubmissionResourceId provides operations to manage the submittedResources property of the microsoft.graph.educationSubmission entity.
// returns a *MeAssignmentsItemSubmissionsItemSubmittedResourcesEducationSubmissionResourceItemRequestBuilder when successful
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) ByEducationSubmissionResourceId(educationSubmissionResourceId string)(*MeAssignmentsItemSubmissionsItemSubmittedResourcesEducationSubmissionResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if educationSubmissionResourceId != "" {
        urlTplParams["educationSubmissionResource%2Did"] = educationSubmissionResourceId
    }
    return NewMeAssignmentsItemSubmissionsItemSubmittedResourcesEducationSubmissionResourceItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewMeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderInternal instantiates a new MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder and sets the default values.
func NewMeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) {
    m := &MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/me/assignments/{educationAssignment%2Did}/submissions/{educationSubmission%2Did}/submittedResources{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder instantiates a new MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder and sets the default values.
func NewMeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *MeAssignmentsItemSubmissionsItemSubmittedResourcesCountRequestBuilder when successful
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) Count()(*MeAssignmentsItemSubmissionsItemSubmittedResourcesCountRequestBuilder) {
    return NewMeAssignmentsItemSubmissionsItemSubmittedResourcesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get submittedResources from education
// returns a EducationSubmissionResourceCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) Get(ctx context.Context, requestConfiguration *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionResourceCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationSubmissionResourceCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionResourceCollectionResponseable), nil
}
// Post create new navigation property to submittedResources for education
// returns a EducationSubmissionResourceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionResourceable, requestConfiguration *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionResourceable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationSubmissionResourceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionResourceable), nil
}
// ToGetRequestInformation get submittedResources from education
// returns a *RequestInformation when successful
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to submittedResources for education
// returns a *RequestInformation when successful
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionResourceable, requestConfiguration *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder when successful
func (m *MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) WithUrl(rawUrl string)(*MeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder) {
    return NewMeAssignmentsItemSubmissionsItemSubmittedResourcesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
