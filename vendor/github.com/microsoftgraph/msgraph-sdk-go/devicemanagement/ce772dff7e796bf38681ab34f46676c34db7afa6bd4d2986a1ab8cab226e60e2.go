package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder provides operations to manage the acceptanceStatuses property of the microsoft.graph.termsAndConditions entity.
type TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetQueryParameters read properties and relationships of the termsAndConditionsAcceptanceStatus object.
type TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetQueryParameters
}
// TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewTermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderInternal instantiates a new TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder and sets the default values.
func NewTermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) {
    m := &TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/termsAndConditions/{termsAndConditions%2Did}/acceptanceStatuses/{termsAndConditionsAcceptanceStatus%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewTermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder instantiates a new TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder and sets the default values.
func NewTermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete deletes a termsAndConditionsAcceptanceStatus.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-companyterms-termsandconditionsacceptancestatus-delete?view=graph-rest-1.0
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read properties and relationships of the termsAndConditionsAcceptanceStatus object.
// returns a TermsAndConditionsAcceptanceStatusable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-companyterms-termsandconditionsacceptancestatus-get?view=graph-rest-1.0
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) Get(ctx context.Context, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsAcceptanceStatusable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTermsAndConditionsAcceptanceStatusFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsAcceptanceStatusable), nil
}
// Patch update the properties of a termsAndConditionsAcceptanceStatus object.
// returns a TermsAndConditionsAcceptanceStatusable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-companyterms-termsandconditionsacceptancestatus-update?view=graph-rest-1.0
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsAcceptanceStatusable, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsAcceptanceStatusable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTermsAndConditionsAcceptanceStatusFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsAcceptanceStatusable), nil
}
// TermsAndConditions provides operations to manage the termsAndConditions property of the microsoft.graph.termsAndConditionsAcceptanceStatus entity.
// returns a *TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder when successful
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) TermsAndConditions()(*TermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilder) {
    return NewTermsAndConditionsItemAcceptanceStatusesItemTermsAndConditionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation deletes a termsAndConditionsAcceptanceStatus.
// returns a *RequestInformation when successful
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read properties and relationships of the termsAndConditionsAcceptanceStatus object.
// returns a *RequestInformation when successful
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a termsAndConditionsAcceptanceStatus object.
// returns a *RequestInformation when successful
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TermsAndConditionsAcceptanceStatusable, requestConfiguration *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder when successful
func (m *TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) WithUrl(rawUrl string)(*TermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder) {
    return NewTermsAndConditionsItemAcceptanceStatusesTermsAndConditionsAcceptanceStatusItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
