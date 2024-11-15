package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder provides operations to manage the oneDriveForBusinessProtectionPolicies property of the microsoft.graph.backupRestoreRoot entity.
type BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetQueryParameters the list of OneDrive for Business protection policies in the tenant.
type BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetQueryParameters
}
// BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewBackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderInternal instantiates a new BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder and sets the default values.
func NewBackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) {
    m := &BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/backupRestore/oneDriveForBusinessProtectionPolicies/{oneDriveForBusinessProtectionPolicy%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewBackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder instantiates a new BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder and sets the default values.
func NewBackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property oneDriveForBusinessProtectionPolicies for solutions
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DriveInclusionRules provides operations to manage the driveInclusionRules property of the microsoft.graph.oneDriveForBusinessProtectionPolicy entity.
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) DriveInclusionRules()(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveInclusionRulesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DriveProtectionUnits provides operations to manage the driveProtectionUnits property of the microsoft.graph.oneDriveForBusinessProtectionPolicy entity.
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveProtectionUnitsRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) DriveProtectionUnits()(*BackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveProtectionUnitsRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesItemDriveProtectionUnitsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the list of OneDrive for Business protection policies in the tenant.
// returns a OneDriveForBusinessProtectionPolicyable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) Get(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OneDriveForBusinessProtectionPolicyable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOneDriveForBusinessProtectionPolicyFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OneDriveForBusinessProtectionPolicyable), nil
}
// Patch update the protection policy for the OneDrive service in Microsoft 365. This method adds a driveProtectionUnit to or removes it from a oneDriveForBusinessProtectionPolicy object.
// returns a OneDriveForBusinessProtectionPolicyable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/onedriveforbusinessprotectionpolicy-update?view=graph-rest-1.0
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OneDriveForBusinessProtectionPolicyable, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OneDriveForBusinessProtectionPolicyable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOneDriveForBusinessProtectionPolicyFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OneDriveForBusinessProtectionPolicyable), nil
}
// ToDeleteRequestInformation delete navigation property oneDriveForBusinessProtectionPolicies for solutions
// returns a *RequestInformation when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the list of OneDrive for Business protection policies in the tenant.
// returns a *RequestInformation when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the protection policy for the OneDrive service in Microsoft 365. This method adds a driveProtectionUnit to or removes it from a oneDriveForBusinessProtectionPolicy object.
// returns a *RequestInformation when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OneDriveForBusinessProtectionPolicyable, requestConfiguration *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder when successful
func (m *BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) WithUrl(rawUrl string)(*BackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder) {
    return NewBackupRestoreOneDriveForBusinessProtectionPoliciesOneDriveForBusinessProtectionPolicyItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
