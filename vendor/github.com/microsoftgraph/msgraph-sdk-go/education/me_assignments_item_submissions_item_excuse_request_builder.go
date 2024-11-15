package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MeAssignmentsItemSubmissionsItemExcuseRequestBuilder provides operations to call the excuse method.
type MeAssignmentsItemSubmissionsItemExcuseRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MeAssignmentsItemSubmissionsItemExcuseRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeAssignmentsItemSubmissionsItemExcuseRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewMeAssignmentsItemSubmissionsItemExcuseRequestBuilderInternal instantiates a new MeAssignmentsItemSubmissionsItemExcuseRequestBuilder and sets the default values.
func NewMeAssignmentsItemSubmissionsItemExcuseRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemSubmissionsItemExcuseRequestBuilder) {
    m := &MeAssignmentsItemSubmissionsItemExcuseRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/me/assignments/{educationAssignment%2Did}/submissions/{educationSubmission%2Did}/excuse", pathParameters),
    }
    return m
}
// NewMeAssignmentsItemSubmissionsItemExcuseRequestBuilder instantiates a new MeAssignmentsItemSubmissionsItemExcuseRequestBuilder and sets the default values.
func NewMeAssignmentsItemSubmissionsItemExcuseRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemSubmissionsItemExcuseRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeAssignmentsItemSubmissionsItemExcuseRequestBuilderInternal(urlParams, requestAdapter)
}
// Post excuse a submission. Excused submissions aren't included in average grade calculations. Grading rubrics and feedback are deleted. Only teachers can perform this action.  If the Prefer: include-unknown-enum-members request header is provided, the excused submission shows the excused status. Otherwise, the submission status retains the returned status. For more information about how to use this header, see the Examples section.
// returns a EducationSubmissionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/educationsubmission-excuse?view=graph-rest-1.0
func (m *MeAssignmentsItemSubmissionsItemExcuseRequestBuilder) Post(ctx context.Context, requestConfiguration *MeAssignmentsItemSubmissionsItemExcuseRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationSubmissionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable), nil
}
// ToPostRequestInformation excuse a submission. Excused submissions aren't included in average grade calculations. Grading rubrics and feedback are deleted. Only teachers can perform this action.  If the Prefer: include-unknown-enum-members request header is provided, the excused submission shows the excused status. Otherwise, the submission status retains the returned status. For more information about how to use this header, see the Examples section.
// returns a *RequestInformation when successful
func (m *MeAssignmentsItemSubmissionsItemExcuseRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *MeAssignmentsItemSubmissionsItemExcuseRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *MeAssignmentsItemSubmissionsItemExcuseRequestBuilder when successful
func (m *MeAssignmentsItemSubmissionsItemExcuseRequestBuilder) WithUrl(rawUrl string)(*MeAssignmentsItemSubmissionsItemExcuseRequestBuilder) {
    return NewMeAssignmentsItemSubmissionsItemExcuseRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
