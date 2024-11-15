package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookRequestBuilder provides operations to manage the workbook property of the microsoft.graph.driveItem entity.
type ItemItemsItemWorkbookRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemItemsItemWorkbookRequestBuilderGetQueryParameters for files that are Excel spreadsheets, access to the workbook API to work with the spreadsheet's contents. Nullable.
type ItemItemsItemWorkbookRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemItemsItemWorkbookRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemItemsItemWorkbookRequestBuilderGetQueryParameters
}
// ItemItemsItemWorkbookRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Application provides operations to manage the application property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookApplicationRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Application()(*ItemItemsItemWorkbookApplicationRequestBuilder) {
    return NewItemItemsItemWorkbookApplicationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CloseSession provides operations to call the closeSession method.
// returns a *ItemItemsItemWorkbookCloseSessionRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) CloseSession()(*ItemItemsItemWorkbookCloseSessionRequestBuilder) {
    return NewItemItemsItemWorkbookCloseSessionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Comments provides operations to manage the comments property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookCommentsRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Comments()(*ItemItemsItemWorkbookCommentsRequestBuilder) {
    return NewItemItemsItemWorkbookCommentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemItemsItemWorkbookRequestBuilderInternal instantiates a new ItemItemsItemWorkbookRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookRequestBuilder) {
    m := &ItemItemsItemWorkbookRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookRequestBuilder instantiates a new ItemItemsItemWorkbookRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateSession provides operations to call the createSession method.
// returns a *ItemItemsItemWorkbookCreateSessionRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) CreateSession()(*ItemItemsItemWorkbookCreateSessionRequestBuilder) {
    return NewItemItemsItemWorkbookCreateSessionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delete delete navigation property workbook for drives
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookRequestBuilderDeleteRequestConfiguration)(error) {
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
// Functions provides operations to manage the functions property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookFunctionsRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Functions()(*ItemItemsItemWorkbookFunctionsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get for files that are Excel spreadsheets, access to the workbook API to work with the spreadsheet's contents. Nullable.
// returns a Workbookable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Workbookable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Workbookable), nil
}
// Names provides operations to manage the names property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookNamesRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Names()(*ItemItemsItemWorkbookNamesRequestBuilder) {
    return NewItemItemsItemWorkbookNamesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Operations provides operations to manage the operations property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookOperationsRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Operations()(*ItemItemsItemWorkbookOperationsRequestBuilder) {
    return NewItemItemsItemWorkbookOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property workbook in drives
// returns a Workbookable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Workbookable, requestConfiguration *ItemItemsItemWorkbookRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Workbookable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Workbookable), nil
}
// RefreshSession provides operations to call the refreshSession method.
// returns a *ItemItemsItemWorkbookRefreshSessionRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) RefreshSession()(*ItemItemsItemWorkbookRefreshSessionRequestBuilder) {
    return NewItemItemsItemWorkbookRefreshSessionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SessionInfoResourceWithKey provides operations to call the sessionInfoResource method.
// returns a *ItemItemsItemWorkbookSessionInfoResourceWithKeyRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) SessionInfoResourceWithKey(key *string)(*ItemItemsItemWorkbookSessionInfoResourceWithKeyRequestBuilder) {
    return NewItemItemsItemWorkbookSessionInfoResourceWithKeyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, key)
}
// TableRowOperationResultWithKey provides operations to call the tableRowOperationResult method.
// returns a *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) TableRowOperationResultWithKey(key *string)(*ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) {
    return NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, key)
}
// Tables provides operations to manage the tables property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookTablesRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Tables()(*ItemItemsItemWorkbookTablesRequestBuilder) {
    return NewItemItemsItemWorkbookTablesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property workbook for drives
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation for files that are Excel spreadsheets, access to the workbook API to work with the spreadsheet's contents. Nullable.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property workbook in drives
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Workbookable, requestConfiguration *ItemItemsItemWorkbookRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemItemsItemWorkbookRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookRequestBuilder) {
    return NewItemItemsItemWorkbookRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheets provides operations to manage the worksheets property of the microsoft.graph.workbook entity.
// returns a *ItemItemsItemWorkbookWorksheetsRequestBuilder when successful
func (m *ItemItemsItemWorkbookRequestBuilder) Worksheets()(*ItemItemsItemWorkbookWorksheetsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
