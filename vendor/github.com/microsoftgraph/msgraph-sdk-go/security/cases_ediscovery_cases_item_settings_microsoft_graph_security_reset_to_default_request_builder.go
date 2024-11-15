package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder provides operations to call the resetToDefault method.
type CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderInternal instantiates a new CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder and sets the default values.
func NewCasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder) {
    m := &CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/cases/ediscoveryCases/{ediscoveryCase%2Did}/settings/microsoft.graph.security.resetToDefault", pathParameters),
    }
    return m
}
// NewCasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder instantiates a new CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder and sets the default values.
func NewCasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderInternal(urlParams, requestAdapter)
}
// Post reset a caseSettings object to the default values.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/security-ediscoverycasesettings-resettodefault?view=graph-rest-1.0
func (m *CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder) Post(ctx context.Context, requestConfiguration *CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation reset a caseSettings object to the default values.
// returns a *RequestInformation when successful
func (m *CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder when successful
func (m *CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder) WithUrl(rawUrl string)(*CasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder) {
    return NewCasesEdiscoveryCasesItemSettingsMicrosoftGraphSecurityResetToDefaultRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
