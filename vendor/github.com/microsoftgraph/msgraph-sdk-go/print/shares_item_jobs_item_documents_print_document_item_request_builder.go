package print

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder provides operations to manage the documents property of the microsoft.graph.printJob entity.
type SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetQueryParameters get documents from print
type SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetQueryParameters
}
// SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewSharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderInternal instantiates a new SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder and sets the default values.
func NewSharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) {
    m := &SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/print/shares/{printerShare%2Did}/jobs/{printJob%2Did}/documents/{printDocument%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewSharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder instantiates a new SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder and sets the default values.
func NewSharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Content provides operations to manage the media for the print entity.
// returns a *SharesItemJobsItemDocumentsItemValueContentRequestBuilder when successful
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) Content()(*SharesItemJobsItemDocumentsItemValueContentRequestBuilder) {
    return NewSharesItemJobsItemDocumentsItemValueContentRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreateUploadSession provides operations to call the createUploadSession method.
// returns a *SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder when successful
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) CreateUploadSession()(*SharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilder) {
    return NewSharesItemJobsItemDocumentsItemCreateUploadSessionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delete delete navigation property documents for print
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get documents from print
// returns a PrintDocumentable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) Get(ctx context.Context, requestConfiguration *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintDocumentable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrintDocumentFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintDocumentable), nil
}
// Patch update the navigation property documents in print
// returns a PrintDocumentable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintDocumentable, requestConfiguration *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintDocumentable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrintDocumentFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintDocumentable), nil
}
// ToDeleteRequestInformation delete navigation property documents for print
// returns a *RequestInformation when successful
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get documents from print
// returns a *RequestInformation when successful
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property documents in print
// returns a *RequestInformation when successful
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintDocumentable, requestConfiguration *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder when successful
func (m *SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) WithUrl(rawUrl string)(*SharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder) {
    return NewSharesItemJobsItemDocumentsPrintDocumentItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
