package identityprotection

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder provides operations to manage the history property of the microsoft.graph.riskyServicePrincipal entity.
type RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetQueryParameters represents the risk history of Microsoft Entra service principals.
type RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetQueryParameters
}
// RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewRiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderInternal instantiates a new RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder and sets the default values.
func NewRiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) {
    m := &RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityProtection/riskyServicePrincipals/{riskyServicePrincipal%2Did}/history/{riskyServicePrincipalHistoryItem%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewRiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder instantiates a new RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder and sets the default values.
func NewRiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewRiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property history for identityProtection
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get represents the risk history of Microsoft Entra service principals.
// returns a RiskyServicePrincipalHistoryItemable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) Get(ctx context.Context, requestConfiguration *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RiskyServicePrincipalHistoryItemable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateRiskyServicePrincipalHistoryItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RiskyServicePrincipalHistoryItemable), nil
}
// Patch update the navigation property history in identityProtection
// returns a RiskyServicePrincipalHistoryItemable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RiskyServicePrincipalHistoryItemable, requestConfiguration *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RiskyServicePrincipalHistoryItemable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateRiskyServicePrincipalHistoryItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RiskyServicePrincipalHistoryItemable), nil
}
// ToDeleteRequestInformation delete navigation property history for identityProtection
// returns a *RequestInformation when successful
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation represents the risk history of Microsoft Entra service principals.
// returns a *RequestInformation when successful
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property history in identityProtection
// returns a *RequestInformation when successful
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RiskyServicePrincipalHistoryItemable, requestConfiguration *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder when successful
func (m *RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) WithUrl(rawUrl string)(*RiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder) {
    return NewRiskyServicePrincipalsItemHistoryRiskyServicePrincipalHistoryItemItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
