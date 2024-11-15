package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder provides operations to manage the termsAndConditions property of the microsoft.graph.termsAndConditionsAcceptanceStatus entity.
type TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetQueryParameters navigation link to the terms and conditions that are assigned.
type TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetQueryParameters
}
// NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderInternal instantiates a new TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder and sets the default values.
func NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) {
    m := &TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/termsAndConditions/{termsAndConditions%2Did}/acceptanceStatuses/{termsAndConditionsAcceptanceStatus%2Did}/termsAndConditions{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder instantiates a new TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder and sets the default values.
func NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get navigation link to the terms and conditions that are assigned.
// returns a TermsAndConditionsable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) Get(ctx context.Context, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTermsAndConditionsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsable), nil
}
// ToGetRequestInformation navigation link to the terms and conditions that are assigned.
// returns a *RequestInformation when successful
func (m *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder when successful
func (m *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) WithUrl(rawUrl string)(*TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) {
    return NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
