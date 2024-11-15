package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder provides operations to manage the driveInclusionRules property of the microsoft.graph.oneDriveForBusinessProtectionPolicy entity.
type BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetQueryParameters get a list of the driveProtectionRule objects that are associated with a OneDrive for Business protection policy. An inclusion rule indicates that a protection policy should contain protection units that match the specified rule criteria. The initial status of a protection rule upon creation is active. After the rule is applied, the state is either completed or completedWithErrors.
type BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetQueryParameters
}
// ByDriveProtectionRuleId provides operations to manage the driveInclusionRules property of the microsoft.graph.oneDriveForBusinessProtectionPolicy entity.
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesDriveProtectionRuleItemRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) ByDriveProtectionRuleId(driveProtectionRuleId string)(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesDriveProtectionRuleItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if driveProtectionRuleId != "" {
        urlTplParams["driveProtectionRule%2Did"] = driveProtectionRuleId
    }
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesDriveProtectionRuleItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderInternal instantiates a new BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder and sets the default values.
func NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) {
    m := &BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/oneDriveForBusinessProtectionPolicies/{oneDriveForBusinessProtectionPolicy%2Did}/driveInclusionRules{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder instantiates a new BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder and sets the default values.
func NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesCountRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) Count()(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesCountRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get a list of the driveProtectionRule objects that are associated with a OneDrive for Business protection policy. An inclusion rule indicates that a protection policy should contain protection units that match the specified rule criteria. The initial status of a protection rule upon creation is active. After the rule is applied, the state is either completed or completedWithErrors.
// returns a DriveProtectionRuleCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/onedriveforbusinessprotectionpolicy-list-driveinclusionrules?view=graph-rest-1.0
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveProtectionRuleCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveProtectionRuleCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveProtectionRuleCollectionResponseable), nil
}
// ToGetRequestInformation get a list of the driveProtectionRule objects that are associated with a OneDrive for Business protection policy. An inclusion rule indicates that a protection policy should contain protection units that match the specified rule criteria. The initial status of a protection rule upon creation is active. After the rule is applied, the state is either completed or completedWithErrors.
// returns a *RequestInformation when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
