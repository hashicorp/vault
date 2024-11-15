package print

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder provides operations to call the createUploadSession method.
type SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderInternal instantiates a new SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder and sets the default values.
func NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) {
    m := &SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/print/shares/{printerShare%2Did}/jobs/{printJob%2Did}/documents/{printDocument%2Did}/createUploadSession", pathParameters),
    }
    return m
}
// NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder instantiates a new SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder and sets the default values.
func NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderInternal(urlParams, requestAdapter)
}
// Post create an upload session that allows an app to iteratively upload ranges of a binary file linked to the print document. As part of the response, this action returns an upload URL that can be used in subsequent sequential PUT queries. Request headers for each PUT operation can be used to specify the exact range of bytes to be uploaded. This allows transfer to be resumed, in case the network connection is dropped during upload. 
// returns a UploadSessionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/printdocument-createuploadsession?view=graph-rest-1.0
func (m *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) Post(ctx context.Context, body SharesItemJobsItemDocumentsItemCreateUploadSessionPostRequestBodyable, requestConfiguration *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UploadSessionable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUploadSessionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UploadSessionable), nil
}
// ToPostRequestInformation create an upload session that allows an app to iteratively upload ranges of a binary file linked to the print document. As part of the response, this action returns an upload URL that can be used in subsequent sequential PUT queries. Request headers for each PUT operation can be used to specify the exact range of bytes to be uploaded. This allows transfer to be resumed, in case the network connection is dropped during upload. 
// returns a *RequestInformation when successful
func (m *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) ToPostRequestInformation(ctx context.Context, body SharesItemJobsItemDocumentsItemCreateUploadSessionPostRequestBodyable, requestConfiguration *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder when successful
func (m *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) WithUrl(rawUrl string)(*SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) {
    return NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
