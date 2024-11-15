package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder provides operations to count the resources in the collection.
type BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetQueryParameters get the number of the resource
type BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetQueryParameters struct {
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
}
// BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetQueryParameters
}
// NewBackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderInternal instantiates a new BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder and sets the default values.
func NewBackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder) {
    m := &BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/sharePointProtectionPolicies/{sharePointProtectionPolicy%2Did}/siteInclusionRules/$count{?%24filter,%24search}", pathParameters),
    }
    return m
}
// NewBackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder instantiates a new BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder and sets the default values.
func NewBackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get the number of the resource
// returns a *int32 when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetRequestConfiguration)(*int32, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.SendPrimitive(ctx, requestInfo, "int32", errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(*int32), nil
}
// ToGetRequestInformation get the number of the resource
// returns a *RequestInformation when successful
func (m *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "text/plain;q=0.9")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder when successful
func (m *BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder) {
    return NewBackupRestoreSharePointProtectionPoliciesItemSiteInclusionRulesCountRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
