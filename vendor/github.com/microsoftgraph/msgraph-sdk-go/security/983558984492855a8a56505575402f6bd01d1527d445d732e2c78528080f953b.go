package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder provides operations to call the release method.
type CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderInternal instantiates a new CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder and sets the default values.
func NewCasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder) {
    m := &CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/cases/ediscoveryCases/{ediscoveryCase%2Did}/noncustodialDataSources/{ediscoveryNoncustodialDataSource%2Did}/microsoft.graph.security.release", pathParameters),
    }
    return m
}
// NewCasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder instantiates a new CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder and sets the default values.
func NewCasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderInternal(urlParams, requestAdapter)
}
// Post release the non-custodial data source from the case.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/security-ediscoverynoncustodialdatasource-release?view=graph-rest-1.0
func (m *CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder) Post(ctx context.Context, requestConfiguration *CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation release the non-custodial data source from the case.
// returns a *RequestInformation when successful
func (m *CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder when successful
func (m *CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder) WithUrl(rawUrl string)(*CasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder) {
    return NewCasesEdiscoveryCasesItemNoncustodialDataSourcesItemMicrosoftGraphSecurityReleaseRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
