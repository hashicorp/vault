package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreRequestBuilder provides operations to manage the backupRestore property of the microsoft.graph.solutionsRoot entity.
type BackupRestoreRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BackupRestoreRequestBuilderGetQueryParameters get the serviceStatus of the Microsoft 365 Backup Storage service in a tenant.
type BackupRestoreRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BackupRestoreRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreRequestBuilderGetQueryParameters
}
// BackupRestoreRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewBackupRestoreRequestBuilderInternal instantiates a new BackupRestoreRequestBuilder and sets the default values.
func NewBackupRestoreRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreRequestBuilder) {
    m := &BackupRestoreRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewBackupRestoreRequestBuilder instantiates a new BackupRestoreRequestBuilder and sets the default values.
func NewBackupRestoreRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property backupRestore for solutions
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreRequestBuilder) Delete(ctx context.Context, requestConfiguration *BackupRestoreRequestBuilderDeleteRequestConfiguration)(error) {
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
// DriveInclusionRules provides operations to manage the driveInclusionRules property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreDriveInclusionRulesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) DriveInclusionRules()(*BackupRestoreDriveInclusionRulesRequestBuilder) {
    return NewBackupRestoreDriveInclusionRulesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DriveProtectionUnits provides operations to manage the driveProtectionUnits property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreDriveProtectionUnitsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) DriveProtectionUnits()(*BackupRestoreDriveProtectionUnitsRequestBuilder) {
    return NewBackupRestoreDriveProtectionUnitsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Enable provides operations to call the enable method.
// returns a *BackupRestoreEnableRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) Enable()(*BackupRestoreEnableRequestBuilder) {
    return NewBackupRestoreEnableRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExchangeProtectionPolicies provides operations to manage the exchangeProtectionPolicies property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreExchangeProtectionPoliciesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) ExchangeProtectionPolicies()(*BackupRestoreExchangeProtectionPoliciesRequestBuilder) {
    return NewBackupRestoreExchangeProtectionPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExchangeRestoreSessions provides operations to manage the exchangeRestoreSessions property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreExchangeRestoreSessionsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) ExchangeRestoreSessions()(*BackupRestoreExchangeRestoreSessionsRequestBuilder) {
    return NewBackupRestoreExchangeRestoreSessionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get the serviceStatus of the Microsoft 365 Backup Storage service in a tenant.
// returns a BackupRestoreRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/backuprestoreroot-get?view=graph-rest-1.0
func (m *BackupRestoreRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BackupRestoreRootable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateBackupRestoreRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BackupRestoreRootable), nil
}
// MailboxInclusionRules provides operations to manage the mailboxInclusionRules property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreMailboxInclusionRulesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) MailboxInclusionRules()(*BackupRestoreMailboxInclusionRulesRequestBuilder) {
    return NewBackupRestoreMailboxInclusionRulesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MailboxProtectionUnits provides operations to manage the mailboxProtectionUnits property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreMailboxProtectionUnitsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) MailboxProtectionUnits()(*BackupRestoreMailboxProtectionUnitsRequestBuilder) {
    return NewBackupRestoreMailboxProtectionUnitsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OneDriveForBusinessProtectionPolicies provides operations to manage the oneDriveForBusinessProtectionPolicies property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) OneDriveForBusinessProtectionPolicies()(*BackupRestoreOneDriveForBusinessProtectionPoliciesRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OneDriveForBusinessRestoreSessions provides operations to manage the oneDriveForBusinessRestoreSessions property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreOneDriveForBusinessRestoreSessionsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) OneDriveForBusinessRestoreSessions()(*BackupRestoreOneDriveForBusinessRestoreSessionsRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessRestoreSessionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property backupRestore in solutions
// returns a BackupRestoreRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BackupRestoreRootable, requestConfiguration *BackupRestoreRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BackupRestoreRootable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateBackupRestoreRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BackupRestoreRootable), nil
}
// ProtectionPolicies provides operations to manage the protectionPolicies property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreProtectionPoliciesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) ProtectionPolicies()(*BackupRestoreProtectionPoliciesRequestBuilder) {
    return NewBackupRestoreProtectionPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ProtectionUnits provides operations to manage the protectionUnits property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreProtectionUnitsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) ProtectionUnits()(*BackupRestoreProtectionUnitsRequestBuilder) {
    return NewBackupRestoreProtectionUnitsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RestorePoints provides operations to manage the restorePoints property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreRestorePointsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) RestorePoints()(*BackupRestoreRestorePointsRequestBuilder) {
    return NewBackupRestoreRestorePointsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RestoreSessions provides operations to manage the restoreSessions property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreRestoreSessionsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) RestoreSessions()(*BackupRestoreRestoreSessionsRequestBuilder) {
    return NewBackupRestoreRestoreSessionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServiceApps provides operations to manage the serviceApps property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreServiceAppsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) ServiceApps()(*BackupRestoreServiceAppsRequestBuilder) {
    return NewBackupRestoreServiceAppsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SharePointProtectionPolicies provides operations to manage the sharePointProtectionPolicies property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreSharePointProtectionPoliciesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) SharePointProtectionPolicies()(*BackupRestoreSharePointProtectionPoliciesRequestBuilder) {
    return NewBackupRestoreSharePointProtectionPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SharePointRestoreSessions provides operations to manage the sharePointRestoreSessions property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreSharePointRestoreSessionsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) SharePointRestoreSessions()(*BackupRestoreSharePointRestoreSessionsRequestBuilder) {
    return NewBackupRestoreSharePointRestoreSessionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SiteInclusionRules provides operations to manage the siteInclusionRules property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreSiteInclusionRulesRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) SiteInclusionRules()(*BackupRestoreSiteInclusionRulesRequestBuilder) {
    return NewBackupRestoreSiteInclusionRulesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SiteProtectionUnits provides operations to manage the siteProtectionUnits property of the microsoft.graph.backupRestoreRoot entity.
// returns a *BackupRestoreSiteProtectionUnitsRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) SiteProtectionUnits()(*BackupRestoreSiteProtectionUnitsRequestBuilder) {
    return NewBackupRestoreSiteProtectionUnitsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property backupRestore for solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get the serviceStatus of the Microsoft 365 Backup Storage service in a tenant.
// returns a *RequestInformation when successful
func (m *BackupRestoreRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property backupRestore in solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BackupRestoreRootable, requestConfiguration *BackupRestoreRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreRequestBuilder when successful
func (m *BackupRestoreRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreRequestBuilder) {
    return NewBackupRestoreRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
