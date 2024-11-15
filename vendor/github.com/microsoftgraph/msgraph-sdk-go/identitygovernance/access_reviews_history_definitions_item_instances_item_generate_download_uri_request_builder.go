package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder provides operations to call the generateDownloadUri method.
type AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderInternal instantiates a new AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder and sets the default values.
func NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) {
    m := &AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/accessReviews/historyDefinitions/{accessReviewHistoryDefinition%2Did}/instances/{accessReviewHistoryInstance%2Did}/generateDownloadUri", pathParameters),
    }
    return m
}
// NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder instantiates a new AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder and sets the default values.
func NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderInternal(urlParams, requestAdapter)
}
// Post generates a URI for an accessReviewHistoryInstance object the status for which is done. Each URI can be used to retrieve the instance's review history data. Each URI is valid for 24 hours and can be retrieved by fetching the downloadUri property from the accessReviewHistoryInstance object.
// returns a AccessReviewHistoryInstanceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/accessreviewhistoryinstance-generatedownloaduri?view=graph-rest-1.0
func (m *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) Post(ctx context.Context, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewHistoryInstanceable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation generates a URI for an accessReviewHistoryInstance object the status for which is done. Each URI can be used to retrieve the instance's review history data. Each URI is valid for 24 hours and can be retrieved by fetching the downloadUri property from the accessReviewHistoryInstance object.
// returns a *RequestInformation when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder when successful
func (m *AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) WithUrl(rawUrl string)(*AccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder) {
    return NewAccessReviewsHistoryDefinitionsItemInstancesItemGenerateDownloadUriRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
