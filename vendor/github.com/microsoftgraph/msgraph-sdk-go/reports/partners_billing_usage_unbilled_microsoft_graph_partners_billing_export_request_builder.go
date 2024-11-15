package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4 "github.com/microsoftgraph/msgraph-sdk-go/models/partners/billing"
)

// PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder provides operations to call the export method.
type PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderInternal instantiates a new PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder and sets the default values.
func NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder) {
    m := &PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/partners/billing/usage/unbilled/microsoft.graph.partners.billing.export", pathParameters),
    }
    return m
}
// NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder instantiates a new PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder and sets the default values.
func NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderInternal(urlParams, requestAdapter)
}
// Post export unbilled Azure usage data for a specific billing period and currency.
// returns a Operationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/partners-billing-unbilledusage-export?view=graph-rest-1.0
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder) Post(ctx context.Context, body PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBodyable, requestConfiguration *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderPostRequestConfiguration)(ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.Operationable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.CreateOperationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.Operationable), nil
}
// ToPostRequestInformation export unbilled Azure usage data for a specific billing period and currency.
// returns a *RequestInformation when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder) ToPostRequestInformation(ctx context.Context, body PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBodyable, requestConfiguration *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder) WithUrl(rawUrl string)(*PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder) {
    return NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
