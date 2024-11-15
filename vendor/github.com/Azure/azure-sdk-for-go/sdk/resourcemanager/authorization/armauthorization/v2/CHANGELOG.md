# Release History

## 2.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 3.0.0-beta.1 (2023-07-28)
### Breaking Changes

- Field `EffectiveRules` of struct `RoleManagementPolicyAssignmentProperties` has been removed

### Features Added

- New enum type `AccessRecommendationType` with values `AccessRecommendationTypeApprove`, `AccessRecommendationTypeDeny`, `AccessRecommendationTypeNoInfoAvailable`
- New enum type `AccessReviewActorIdentityType` with values `AccessReviewActorIdentityTypeServicePrincipal`, `AccessReviewActorIdentityTypeUser`
- New enum type `AccessReviewApplyResult` with values `AccessReviewApplyResultAppliedSuccessfully`, `AccessReviewApplyResultAppliedSuccessfullyButObjectNotFound`, `AccessReviewApplyResultAppliedWithUnknownFailure`, `AccessReviewApplyResultApplyNotSupported`, `AccessReviewApplyResultApplying`, `AccessReviewApplyResultNew`
- New enum type `AccessReviewDecisionInsightType` with values `AccessReviewDecisionInsightTypeUserSignInInsight`
- New enum type `AccessReviewDecisionPrincipalResourceMembershipType` with values `AccessReviewDecisionPrincipalResourceMembershipTypeDirect`, `AccessReviewDecisionPrincipalResourceMembershipTypeIndirect`
- New enum type `AccessReviewHistoryDefinitionStatus` with values `AccessReviewHistoryDefinitionStatusDone`, `AccessReviewHistoryDefinitionStatusError`, `AccessReviewHistoryDefinitionStatusInProgress`, `AccessReviewHistoryDefinitionStatusRequested`
- New enum type `AccessReviewInstanceReviewersType` with values `AccessReviewInstanceReviewersTypeAssigned`, `AccessReviewInstanceReviewersTypeManagers`, `AccessReviewInstanceReviewersTypeSelf`
- New enum type `AccessReviewInstanceStatus` with values `AccessReviewInstanceStatusApplied`, `AccessReviewInstanceStatusApplying`, `AccessReviewInstanceStatusAutoReviewed`, `AccessReviewInstanceStatusAutoReviewing`, `AccessReviewInstanceStatusCompleted`, `AccessReviewInstanceStatusCompleting`, `AccessReviewInstanceStatusInProgress`, `AccessReviewInstanceStatusInitializing`, `AccessReviewInstanceStatusNotStarted`, `AccessReviewInstanceStatusScheduled`, `AccessReviewInstanceStatusStarting`
- New enum type `AccessReviewRecurrencePatternType` with values `AccessReviewRecurrencePatternTypeAbsoluteMonthly`, `AccessReviewRecurrencePatternTypeWeekly`
- New enum type `AccessReviewRecurrenceRangeType` with values `AccessReviewRecurrenceRangeTypeEndDate`, `AccessReviewRecurrenceRangeTypeNoEnd`, `AccessReviewRecurrenceRangeTypeNumbered`
- New enum type `AccessReviewResult` with values `AccessReviewResultApprove`, `AccessReviewResultDeny`, `AccessReviewResultDontKnow`, `AccessReviewResultNotNotified`, `AccessReviewResultNotReviewed`
- New enum type `AccessReviewReviewerType` with values `AccessReviewReviewerTypeServicePrincipal`, `AccessReviewReviewerTypeUser`
- New enum type `AccessReviewScheduleDefinitionReviewersType` with values `AccessReviewScheduleDefinitionReviewersTypeAssigned`, `AccessReviewScheduleDefinitionReviewersTypeManagers`, `AccessReviewScheduleDefinitionReviewersTypeSelf`
- New enum type `AccessReviewScheduleDefinitionStatus` with values `AccessReviewScheduleDefinitionStatusApplied`, `AccessReviewScheduleDefinitionStatusApplying`, `AccessReviewScheduleDefinitionStatusAutoReviewed`, `AccessReviewScheduleDefinitionStatusAutoReviewing`, `AccessReviewScheduleDefinitionStatusCompleted`, `AccessReviewScheduleDefinitionStatusCompleting`, `AccessReviewScheduleDefinitionStatusInProgress`, `AccessReviewScheduleDefinitionStatusInitializing`, `AccessReviewScheduleDefinitionStatusNotStarted`, `AccessReviewScheduleDefinitionStatusScheduled`, `AccessReviewScheduleDefinitionStatusStarting`
- New enum type `AccessReviewScopeAssignmentState` with values `AccessReviewScopeAssignmentStateActive`, `AccessReviewScopeAssignmentStateEligible`
- New enum type `AccessReviewScopePrincipalType` with values `AccessReviewScopePrincipalTypeGuestUser`, `AccessReviewScopePrincipalTypeRedeemedGuestUser`, `AccessReviewScopePrincipalTypeServicePrincipal`, `AccessReviewScopePrincipalTypeUser`, `AccessReviewScopePrincipalTypeUserGroup`
- New enum type `DecisionResourceType` with values `DecisionResourceTypeAzureRole`
- New enum type `DecisionTargetType` with values `DecisionTargetTypeServicePrincipal`, `DecisionTargetTypeUser`
- New enum type `DefaultDecisionType` with values `DefaultDecisionTypeApprove`, `DefaultDecisionTypeDeny`, `DefaultDecisionTypeRecommendation`
- New enum type `RecordAllDecisionsResult` with values `RecordAllDecisionsResultApprove`, `RecordAllDecisionsResultDeny`
- New enum type `SeverityLevel` with values `SeverityLevelHigh`, `SeverityLevelLow`, `SeverityLevelMedium`
- New function `*AccessReviewDecisionIdentity.GetAccessReviewDecisionIdentity() *AccessReviewDecisionIdentity`
- New function `*AccessReviewDecisionInsightProperties.GetAccessReviewDecisionInsightProperties() *AccessReviewDecisionInsightProperties`
- New function `*AccessReviewDecisionServicePrincipalIdentity.GetAccessReviewDecisionIdentity() *AccessReviewDecisionIdentity`
- New function `*AccessReviewDecisionUserIdentity.GetAccessReviewDecisionIdentity() *AccessReviewDecisionIdentity`
- New function `*AccessReviewDecisionUserSignInInsightProperties.GetAccessReviewDecisionInsightProperties() *AccessReviewDecisionInsightProperties`
- New function `NewAccessReviewDefaultSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewDefaultSettingsClient, error)`
- New function `*AccessReviewDefaultSettingsClient.Get(context.Context, *AccessReviewDefaultSettingsClientGetOptions) (AccessReviewDefaultSettingsClientGetResponse, error)`
- New function `*AccessReviewDefaultSettingsClient.Put(context.Context, AccessReviewScheduleSettings, *AccessReviewDefaultSettingsClientPutOptions) (AccessReviewDefaultSettingsClientPutResponse, error)`
- New function `NewAccessReviewHistoryDefinitionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewHistoryDefinitionClient, error)`
- New function `*AccessReviewHistoryDefinitionClient.Create(context.Context, string, AccessReviewHistoryDefinitionProperties, *AccessReviewHistoryDefinitionClientCreateOptions) (AccessReviewHistoryDefinitionClientCreateResponse, error)`
- New function `*AccessReviewHistoryDefinitionClient.DeleteByID(context.Context, string, *AccessReviewHistoryDefinitionClientDeleteByIDOptions) (AccessReviewHistoryDefinitionClientDeleteByIDResponse, error)`
- New function `NewAccessReviewHistoryDefinitionInstanceClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewHistoryDefinitionInstanceClient, error)`
- New function `*AccessReviewHistoryDefinitionInstanceClient.GenerateDownloadURI(context.Context, string, string, *AccessReviewHistoryDefinitionInstanceClientGenerateDownloadURIOptions) (AccessReviewHistoryDefinitionInstanceClientGenerateDownloadURIResponse, error)`
- New function `NewAccessReviewHistoryDefinitionInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewHistoryDefinitionInstancesClient, error)`
- New function `*AccessReviewHistoryDefinitionInstancesClient.NewListPager(string, *AccessReviewHistoryDefinitionInstancesClientListOptions) *runtime.Pager[AccessReviewHistoryDefinitionInstancesClientListResponse]`
- New function `NewAccessReviewHistoryDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewHistoryDefinitionsClient, error)`
- New function `*AccessReviewHistoryDefinitionsClient.GetByID(context.Context, string, *AccessReviewHistoryDefinitionsClientGetByIDOptions) (AccessReviewHistoryDefinitionsClientGetByIDResponse, error)`
- New function `*AccessReviewHistoryDefinitionsClient.NewListPager(*AccessReviewHistoryDefinitionsClientListOptions) *runtime.Pager[AccessReviewHistoryDefinitionsClientListResponse]`
- New function `NewAccessReviewInstanceClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewInstanceClient, error)`
- New function `*AccessReviewInstanceClient.AcceptRecommendations(context.Context, string, string, *AccessReviewInstanceClientAcceptRecommendationsOptions) (AccessReviewInstanceClientAcceptRecommendationsResponse, error)`
- New function `*AccessReviewInstanceClient.ApplyDecisions(context.Context, string, string, *AccessReviewInstanceClientApplyDecisionsOptions) (AccessReviewInstanceClientApplyDecisionsResponse, error)`
- New function `*AccessReviewInstanceClient.ResetDecisions(context.Context, string, string, *AccessReviewInstanceClientResetDecisionsOptions) (AccessReviewInstanceClientResetDecisionsResponse, error)`
- New function `*AccessReviewInstanceClient.SendReminders(context.Context, string, string, *AccessReviewInstanceClientSendRemindersOptions) (AccessReviewInstanceClientSendRemindersResponse, error)`
- New function `*AccessReviewInstanceClient.Stop(context.Context, string, string, *AccessReviewInstanceClientStopOptions) (AccessReviewInstanceClientStopResponse, error)`
- New function `NewAccessReviewInstanceContactedReviewersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewInstanceContactedReviewersClient, error)`
- New function `*AccessReviewInstanceContactedReviewersClient.NewListPager(string, string, *AccessReviewInstanceContactedReviewersClientListOptions) *runtime.Pager[AccessReviewInstanceContactedReviewersClientListResponse]`
- New function `NewAccessReviewInstanceDecisionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewInstanceDecisionsClient, error)`
- New function `*AccessReviewInstanceDecisionsClient.NewListPager(string, string, *AccessReviewInstanceDecisionsClientListOptions) *runtime.Pager[AccessReviewInstanceDecisionsClientListResponse]`
- New function `NewAccessReviewInstanceMyDecisionsClient(azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewInstanceMyDecisionsClient, error)`
- New function `*AccessReviewInstanceMyDecisionsClient.GetByID(context.Context, string, string, string, *AccessReviewInstanceMyDecisionsClientGetByIDOptions) (AccessReviewInstanceMyDecisionsClientGetByIDResponse, error)`
- New function `*AccessReviewInstanceMyDecisionsClient.NewListPager(string, string, *AccessReviewInstanceMyDecisionsClientListOptions) *runtime.Pager[AccessReviewInstanceMyDecisionsClientListResponse]`
- New function `*AccessReviewInstanceMyDecisionsClient.Patch(context.Context, string, string, string, AccessReviewDecisionProperties, *AccessReviewInstanceMyDecisionsClientPatchOptions) (AccessReviewInstanceMyDecisionsClientPatchResponse, error)`
- New function `NewAccessReviewInstancesAssignedForMyApprovalClient(azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewInstancesAssignedForMyApprovalClient, error)`
- New function `*AccessReviewInstancesAssignedForMyApprovalClient.GetByID(context.Context, string, string, *AccessReviewInstancesAssignedForMyApprovalClientGetByIDOptions) (AccessReviewInstancesAssignedForMyApprovalClientGetByIDResponse, error)`
- New function `*AccessReviewInstancesAssignedForMyApprovalClient.NewListPager(string, *AccessReviewInstancesAssignedForMyApprovalClientListOptions) *runtime.Pager[AccessReviewInstancesAssignedForMyApprovalClientListResponse]`
- New function `NewAccessReviewInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewInstancesClient, error)`
- New function `*AccessReviewInstancesClient.Create(context.Context, string, string, AccessReviewInstanceProperties, *AccessReviewInstancesClientCreateOptions) (AccessReviewInstancesClientCreateResponse, error)`
- New function `*AccessReviewInstancesClient.GetByID(context.Context, string, string, *AccessReviewInstancesClientGetByIDOptions) (AccessReviewInstancesClientGetByIDResponse, error)`
- New function `*AccessReviewInstancesClient.NewListPager(string, *AccessReviewInstancesClientListOptions) *runtime.Pager[AccessReviewInstancesClientListResponse]`
- New function `NewAccessReviewScheduleDefinitionsAssignedForMyApprovalClient(azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewScheduleDefinitionsAssignedForMyApprovalClient, error)`
- New function `*AccessReviewScheduleDefinitionsAssignedForMyApprovalClient.NewListPager(*AccessReviewScheduleDefinitionsAssignedForMyApprovalClientListOptions) *runtime.Pager[AccessReviewScheduleDefinitionsAssignedForMyApprovalClientListResponse]`
- New function `NewAccessReviewScheduleDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessReviewScheduleDefinitionsClient, error)`
- New function `*AccessReviewScheduleDefinitionsClient.CreateOrUpdateByID(context.Context, string, AccessReviewScheduleDefinitionProperties, *AccessReviewScheduleDefinitionsClientCreateOrUpdateByIDOptions) (AccessReviewScheduleDefinitionsClientCreateOrUpdateByIDResponse, error)`
- New function `*AccessReviewScheduleDefinitionsClient.DeleteByID(context.Context, string, *AccessReviewScheduleDefinitionsClientDeleteByIDOptions) (AccessReviewScheduleDefinitionsClientDeleteByIDResponse, error)`
- New function `*AccessReviewScheduleDefinitionsClient.GetByID(context.Context, string, *AccessReviewScheduleDefinitionsClientGetByIDOptions) (AccessReviewScheduleDefinitionsClientGetByIDResponse, error)`
- New function `*AccessReviewScheduleDefinitionsClient.NewListPager(*AccessReviewScheduleDefinitionsClientListOptions) *runtime.Pager[AccessReviewScheduleDefinitionsClientListResponse]`
- New function `*AccessReviewScheduleDefinitionsClient.Stop(context.Context, string, *AccessReviewScheduleDefinitionsClientStopOptions) (AccessReviewScheduleDefinitionsClientStopResponse, error)`
- New function `*AlertConfigurationProperties.GetAlertConfigurationProperties() *AlertConfigurationProperties`
- New function `NewAlertConfigurationsClient(azcore.TokenCredential, *arm.ClientOptions) (*AlertConfigurationsClient, error)`
- New function `*AlertConfigurationsClient.Get(context.Context, string, string, *AlertConfigurationsClientGetOptions) (AlertConfigurationsClientGetResponse, error)`
- New function `*AlertConfigurationsClient.NewListForScopePager(string, *AlertConfigurationsClientListForScopeOptions) *runtime.Pager[AlertConfigurationsClientListForScopeResponse]`
- New function `*AlertConfigurationsClient.Update(context.Context, string, string, AlertConfiguration, *AlertConfigurationsClientUpdateOptions) (AlertConfigurationsClientUpdateResponse, error)`
- New function `NewAlertDefinitionsClient(azcore.TokenCredential, *arm.ClientOptions) (*AlertDefinitionsClient, error)`
- New function `*AlertDefinitionsClient.Get(context.Context, string, string, *AlertDefinitionsClientGetOptions) (AlertDefinitionsClientGetResponse, error)`
- New function `*AlertDefinitionsClient.NewListForScopePager(string, *AlertDefinitionsClientListForScopeOptions) *runtime.Pager[AlertDefinitionsClientListForScopeResponse]`
- New function `*AlertIncidentProperties.GetAlertIncidentProperties() *AlertIncidentProperties`
- New function `NewAlertIncidentsClient(azcore.TokenCredential, *arm.ClientOptions) (*AlertIncidentsClient, error)`
- New function `*AlertIncidentsClient.Get(context.Context, string, string, string, *AlertIncidentsClientGetOptions) (AlertIncidentsClientGetResponse, error)`
- New function `*AlertIncidentsClient.NewListForScopePager(string, string, *AlertIncidentsClientListForScopeOptions) *runtime.Pager[AlertIncidentsClientListForScopeResponse]`
- New function `*AlertIncidentsClient.Remediate(context.Context, string, string, string, *AlertIncidentsClientRemediateOptions) (AlertIncidentsClientRemediateResponse, error)`
- New function `NewAlertOperationClient(azcore.TokenCredential, *arm.ClientOptions) (*AlertOperationClient, error)`
- New function `*AlertOperationClient.Get(context.Context, string, string, *AlertOperationClientGetOptions) (AlertOperationClientGetResponse, error)`
- New function `NewAlertsClient(azcore.TokenCredential, *arm.ClientOptions) (*AlertsClient, error)`
- New function `*AlertsClient.Get(context.Context, string, string, *AlertsClientGetOptions) (AlertsClientGetResponse, error)`
- New function `*AlertsClient.NewListForScopePager(string, *AlertsClientListForScopeOptions) *runtime.Pager[AlertsClientListForScopeResponse]`
- New function `*AlertsClient.BeginRefresh(context.Context, string, string, *AlertsClientBeginRefreshOptions) (*runtime.Poller[AlertsClientRefreshResponse], error)`
- New function `*AlertsClient.BeginRefreshAll(context.Context, string, *AlertsClientBeginRefreshAllOptions) (*runtime.Poller[AlertsClientRefreshAllResponse], error)`
- New function `*AlertsClient.Update(context.Context, string, string, Alert, *AlertsClientUpdateOptions) (AlertsClientUpdateResponse, error)`
- New function `*AzureRolesAssignedOutsidePimAlertConfigurationProperties.GetAlertConfigurationProperties() *AlertConfigurationProperties`
- New function `*AzureRolesAssignedOutsidePimAlertIncidentProperties.GetAlertIncidentProperties() *AlertIncidentProperties`
- New function `*ClientFactory.NewAccessReviewDefaultSettingsClient() *AccessReviewDefaultSettingsClient`
- New function `*ClientFactory.NewAccessReviewHistoryDefinitionClient() *AccessReviewHistoryDefinitionClient`
- New function `*ClientFactory.NewAccessReviewHistoryDefinitionInstanceClient() *AccessReviewHistoryDefinitionInstanceClient`
- New function `*ClientFactory.NewAccessReviewHistoryDefinitionInstancesClient() *AccessReviewHistoryDefinitionInstancesClient`
- New function `*ClientFactory.NewAccessReviewHistoryDefinitionsClient() *AccessReviewHistoryDefinitionsClient`
- New function `*ClientFactory.NewAccessReviewInstanceClient() *AccessReviewInstanceClient`
- New function `*ClientFactory.NewAccessReviewInstanceContactedReviewersClient() *AccessReviewInstanceContactedReviewersClient`
- New function `*ClientFactory.NewAccessReviewInstanceDecisionsClient() *AccessReviewInstanceDecisionsClient`
- New function `*ClientFactory.NewAccessReviewInstanceMyDecisionsClient() *AccessReviewInstanceMyDecisionsClient`
- New function `*ClientFactory.NewAccessReviewInstancesAssignedForMyApprovalClient() *AccessReviewInstancesAssignedForMyApprovalClient`
- New function `*ClientFactory.NewAccessReviewInstancesClient() *AccessReviewInstancesClient`
- New function `*ClientFactory.NewAccessReviewScheduleDefinitionsAssignedForMyApprovalClient() *AccessReviewScheduleDefinitionsAssignedForMyApprovalClient`
- New function `*ClientFactory.NewAccessReviewScheduleDefinitionsClient() *AccessReviewScheduleDefinitionsClient`
- New function `*ClientFactory.NewAlertConfigurationsClient() *AlertConfigurationsClient`
- New function `*ClientFactory.NewAlertDefinitionsClient() *AlertDefinitionsClient`
- New function `*ClientFactory.NewAlertIncidentsClient() *AlertIncidentsClient`
- New function `*ClientFactory.NewAlertOperationClient() *AlertOperationClient`
- New function `*ClientFactory.NewAlertsClient() *AlertsClient`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `*ClientFactory.NewScopeAccessReviewDefaultSettingsClient() *ScopeAccessReviewDefaultSettingsClient`
- New function `*ClientFactory.NewScopeAccessReviewHistoryDefinitionClient() *ScopeAccessReviewHistoryDefinitionClient`
- New function `*ClientFactory.NewScopeAccessReviewHistoryDefinitionInstanceClient() *ScopeAccessReviewHistoryDefinitionInstanceClient`
- New function `*ClientFactory.NewScopeAccessReviewHistoryDefinitionInstancesClient() *ScopeAccessReviewHistoryDefinitionInstancesClient`
- New function `*ClientFactory.NewScopeAccessReviewHistoryDefinitionsClient() *ScopeAccessReviewHistoryDefinitionsClient`
- New function `*ClientFactory.NewScopeAccessReviewInstanceClient() *ScopeAccessReviewInstanceClient`
- New function `*ClientFactory.NewScopeAccessReviewInstanceContactedReviewersClient() *ScopeAccessReviewInstanceContactedReviewersClient`
- New function `*ClientFactory.NewScopeAccessReviewInstanceDecisionsClient() *ScopeAccessReviewInstanceDecisionsClient`
- New function `*ClientFactory.NewScopeAccessReviewInstancesClient() *ScopeAccessReviewInstancesClient`
- New function `*ClientFactory.NewScopeAccessReviewScheduleDefinitionsClient() *ScopeAccessReviewScheduleDefinitionsClient`
- New function `*ClientFactory.NewTenantLevelAccessReviewInstanceContactedReviewersClient() *TenantLevelAccessReviewInstanceContactedReviewersClient`
- New function `*DuplicateRoleCreatedAlertConfigurationProperties.GetAlertConfigurationProperties() *AlertConfigurationProperties`
- New function `*DuplicateRoleCreatedAlertIncidentProperties.GetAlertIncidentProperties() *AlertIncidentProperties`
- New function `NewOperationsClient(azcore.TokenCredential, *arm.ClientOptions) (*OperationsClient, error)`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `NewScopeAccessReviewDefaultSettingsClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewDefaultSettingsClient, error)`
- New function `*ScopeAccessReviewDefaultSettingsClient.Get(context.Context, string, *ScopeAccessReviewDefaultSettingsClientGetOptions) (ScopeAccessReviewDefaultSettingsClientGetResponse, error)`
- New function `*ScopeAccessReviewDefaultSettingsClient.Put(context.Context, string, AccessReviewScheduleSettings, *ScopeAccessReviewDefaultSettingsClientPutOptions) (ScopeAccessReviewDefaultSettingsClientPutResponse, error)`
- New function `NewScopeAccessReviewHistoryDefinitionClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewHistoryDefinitionClient, error)`
- New function `*ScopeAccessReviewHistoryDefinitionClient.Create(context.Context, string, string, AccessReviewHistoryDefinitionProperties, *ScopeAccessReviewHistoryDefinitionClientCreateOptions) (ScopeAccessReviewHistoryDefinitionClientCreateResponse, error)`
- New function `*ScopeAccessReviewHistoryDefinitionClient.DeleteByID(context.Context, string, string, *ScopeAccessReviewHistoryDefinitionClientDeleteByIDOptions) (ScopeAccessReviewHistoryDefinitionClientDeleteByIDResponse, error)`
- New function `NewScopeAccessReviewHistoryDefinitionInstanceClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewHistoryDefinitionInstanceClient, error)`
- New function `*ScopeAccessReviewHistoryDefinitionInstanceClient.GenerateDownloadURI(context.Context, string, string, string, *ScopeAccessReviewHistoryDefinitionInstanceClientGenerateDownloadURIOptions) (ScopeAccessReviewHistoryDefinitionInstanceClientGenerateDownloadURIResponse, error)`
- New function `NewScopeAccessReviewHistoryDefinitionInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewHistoryDefinitionInstancesClient, error)`
- New function `*ScopeAccessReviewHistoryDefinitionInstancesClient.NewListPager(string, string, *ScopeAccessReviewHistoryDefinitionInstancesClientListOptions) *runtime.Pager[ScopeAccessReviewHistoryDefinitionInstancesClientListResponse]`
- New function `NewScopeAccessReviewHistoryDefinitionsClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewHistoryDefinitionsClient, error)`
- New function `*ScopeAccessReviewHistoryDefinitionsClient.GetByID(context.Context, string, string, *ScopeAccessReviewHistoryDefinitionsClientGetByIDOptions) (ScopeAccessReviewHistoryDefinitionsClientGetByIDResponse, error)`
- New function `*ScopeAccessReviewHistoryDefinitionsClient.NewListPager(string, *ScopeAccessReviewHistoryDefinitionsClientListOptions) *runtime.Pager[ScopeAccessReviewHistoryDefinitionsClientListResponse]`
- New function `NewScopeAccessReviewInstanceClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewInstanceClient, error)`
- New function `*ScopeAccessReviewInstanceClient.ApplyDecisions(context.Context, string, string, string, *ScopeAccessReviewInstanceClientApplyDecisionsOptions) (ScopeAccessReviewInstanceClientApplyDecisionsResponse, error)`
- New function `*ScopeAccessReviewInstanceClient.RecordAllDecisions(context.Context, string, string, string, RecordAllDecisionsProperties, *ScopeAccessReviewInstanceClientRecordAllDecisionsOptions) (ScopeAccessReviewInstanceClientRecordAllDecisionsResponse, error)`
- New function `*ScopeAccessReviewInstanceClient.ResetDecisions(context.Context, string, string, string, *ScopeAccessReviewInstanceClientResetDecisionsOptions) (ScopeAccessReviewInstanceClientResetDecisionsResponse, error)`
- New function `*ScopeAccessReviewInstanceClient.SendReminders(context.Context, string, string, string, *ScopeAccessReviewInstanceClientSendRemindersOptions) (ScopeAccessReviewInstanceClientSendRemindersResponse, error)`
- New function `*ScopeAccessReviewInstanceClient.Stop(context.Context, string, string, string, *ScopeAccessReviewInstanceClientStopOptions) (ScopeAccessReviewInstanceClientStopResponse, error)`
- New function `NewScopeAccessReviewInstanceContactedReviewersClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewInstanceContactedReviewersClient, error)`
- New function `*ScopeAccessReviewInstanceContactedReviewersClient.NewListPager(string, string, string, *ScopeAccessReviewInstanceContactedReviewersClientListOptions) *runtime.Pager[ScopeAccessReviewInstanceContactedReviewersClientListResponse]`
- New function `NewScopeAccessReviewInstanceDecisionsClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewInstanceDecisionsClient, error)`
- New function `*ScopeAccessReviewInstanceDecisionsClient.NewListPager(string, string, string, *ScopeAccessReviewInstanceDecisionsClientListOptions) *runtime.Pager[ScopeAccessReviewInstanceDecisionsClientListResponse]`
- New function `NewScopeAccessReviewInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewInstancesClient, error)`
- New function `*ScopeAccessReviewInstancesClient.Create(context.Context, string, string, string, AccessReviewInstanceProperties, *ScopeAccessReviewInstancesClientCreateOptions) (ScopeAccessReviewInstancesClientCreateResponse, error)`
- New function `*ScopeAccessReviewInstancesClient.GetByID(context.Context, string, string, string, *ScopeAccessReviewInstancesClientGetByIDOptions) (ScopeAccessReviewInstancesClientGetByIDResponse, error)`
- New function `*ScopeAccessReviewInstancesClient.NewListPager(string, string, *ScopeAccessReviewInstancesClientListOptions) *runtime.Pager[ScopeAccessReviewInstancesClientListResponse]`
- New function `NewScopeAccessReviewScheduleDefinitionsClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopeAccessReviewScheduleDefinitionsClient, error)`
- New function `*ScopeAccessReviewScheduleDefinitionsClient.CreateOrUpdateByID(context.Context, string, string, AccessReviewScheduleDefinitionProperties, *ScopeAccessReviewScheduleDefinitionsClientCreateOrUpdateByIDOptions) (ScopeAccessReviewScheduleDefinitionsClientCreateOrUpdateByIDResponse, error)`
- New function `*ScopeAccessReviewScheduleDefinitionsClient.DeleteByID(context.Context, string, string, *ScopeAccessReviewScheduleDefinitionsClientDeleteByIDOptions) (ScopeAccessReviewScheduleDefinitionsClientDeleteByIDResponse, error)`
- New function `*ScopeAccessReviewScheduleDefinitionsClient.GetByID(context.Context, string, string, *ScopeAccessReviewScheduleDefinitionsClientGetByIDOptions) (ScopeAccessReviewScheduleDefinitionsClientGetByIDResponse, error)`
- New function `*ScopeAccessReviewScheduleDefinitionsClient.NewListPager(string, *ScopeAccessReviewScheduleDefinitionsClientListOptions) *runtime.Pager[ScopeAccessReviewScheduleDefinitionsClientListResponse]`
- New function `*ScopeAccessReviewScheduleDefinitionsClient.Stop(context.Context, string, string, *ScopeAccessReviewScheduleDefinitionsClientStopOptions) (ScopeAccessReviewScheduleDefinitionsClientStopResponse, error)`
- New function `NewTenantLevelAccessReviewInstanceContactedReviewersClient(azcore.TokenCredential, *arm.ClientOptions) (*TenantLevelAccessReviewInstanceContactedReviewersClient, error)`
- New function `*TenantLevelAccessReviewInstanceContactedReviewersClient.NewListPager(string, string, *TenantLevelAccessReviewInstanceContactedReviewersClientListOptions) *runtime.Pager[TenantLevelAccessReviewInstanceContactedReviewersClientListResponse]`
- New function `*TooManyOwnersAssignedToResourceAlertConfigurationProperties.GetAlertConfigurationProperties() *AlertConfigurationProperties`
- New function `*TooManyOwnersAssignedToResourceAlertIncidentProperties.GetAlertIncidentProperties() *AlertIncidentProperties`
- New function `*TooManyPermanentOwnersAssignedToResourceAlertConfigurationProperties.GetAlertConfigurationProperties() *AlertConfigurationProperties`
- New function `*TooManyPermanentOwnersAssignedToResourceAlertIncidentProperties.GetAlertIncidentProperties() *AlertIncidentProperties`
- New struct `AccessReviewActorIdentity`
- New struct `AccessReviewContactedReviewer`
- New struct `AccessReviewContactedReviewerListResult`
- New struct `AccessReviewContactedReviewerProperties`
- New struct `AccessReviewDecision`
- New struct `AccessReviewDecisionInsight`
- New struct `AccessReviewDecisionListResult`
- New struct `AccessReviewDecisionPrincipalResourceMembership`
- New struct `AccessReviewDecisionProperties`
- New struct `AccessReviewDecisionResource`
- New struct `AccessReviewDecisionServicePrincipalIdentity`
- New struct `AccessReviewDecisionUserIdentity`
- New struct `AccessReviewDecisionUserSignInInsightProperties`
- New struct `AccessReviewDefaultSettings`
- New struct `AccessReviewHistoryDefinition`
- New struct `AccessReviewHistoryDefinitionInstanceListResult`
- New struct `AccessReviewHistoryDefinitionListResult`
- New struct `AccessReviewHistoryDefinitionProperties`
- New struct `AccessReviewHistoryInstance`
- New struct `AccessReviewHistoryInstanceProperties`
- New struct `AccessReviewHistoryScheduleSettings`
- New struct `AccessReviewInstance`
- New struct `AccessReviewInstanceListResult`
- New struct `AccessReviewInstanceProperties`
- New struct `AccessReviewRecurrencePattern`
- New struct `AccessReviewRecurrenceRange`
- New struct `AccessReviewRecurrenceSettings`
- New struct `AccessReviewReviewer`
- New struct `AccessReviewScheduleDefinition`
- New struct `AccessReviewScheduleDefinitionListResult`
- New struct `AccessReviewScheduleDefinitionProperties`
- New struct `AccessReviewScheduleSettings`
- New struct `AccessReviewScope`
- New struct `Alert`
- New struct `AlertConfiguration`
- New struct `AlertConfigurationListResult`
- New struct `AlertDefinition`
- New struct `AlertDefinitionListResult`
- New struct `AlertDefinitionProperties`
- New struct `AlertIncident`
- New struct `AlertIncidentListResult`
- New struct `AlertListResult`
- New struct `AlertOperationResult`
- New struct `AlertProperties`
- New struct `AzureRolesAssignedOutsidePimAlertConfigurationProperties`
- New struct `AzureRolesAssignedOutsidePimAlertIncidentProperties`
- New struct `DuplicateRoleCreatedAlertConfigurationProperties`
- New struct `DuplicateRoleCreatedAlertIncidentProperties`
- New struct `ErrorDefinition`
- New struct `ErrorDefinitionProperties`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `RecordAllDecisionsProperties`
- New struct `TooManyOwnersAssignedToResourceAlertConfigurationProperties`
- New struct `TooManyOwnersAssignedToResourceAlertIncidentProperties`
- New struct `TooManyPermanentOwnersAssignedToResourceAlertConfigurationProperties`
- New struct `TooManyPermanentOwnersAssignedToResourceAlertIncidentProperties`
- New field `Condition`, `ConditionVersion`, `CreatedBy`, `CreatedOn`, `UpdatedBy`, `UpdatedOn` in struct `DenyAssignmentProperties`
- New field `Condition`, `ConditionVersion` in struct `Permission`
- New field `CreatedBy`, `CreatedOn`, `UpdatedBy`, `UpdatedOn` in struct `RoleDefinitionProperties`


## 2.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 2.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.0.0 (2022-09-26)
### Breaking Changes

- Function `*RoleAssignmentsClient.NewListForResourcePager` parameter(s) have been changed from `(string, string, string, string, string, *RoleAssignmentsClientListForResourceOptions)` to `(string, string, string, string, *RoleAssignmentsClientListForResourceOptions)`
- Type of `RoleAssignment.Properties` has been changed from `*RoleAssignmentPropertiesWithScope` to `*RoleAssignmentProperties`
- Function `*RoleAssignmentsClient.NewListPager` has been renamed to `*RoleAssignmentsClient.NewListForSubscriptionPager`

### Features Added

- New function `*DenyAssignmentsClient.Get(context.Context, string, string, *DenyAssignmentsClientGetOptions) (DenyAssignmentsClientGetResponse, error)`
- New function `*DenyAssignmentsClient.NewListForScopePager(string, *DenyAssignmentsClientListForScopeOptions) *runtime.Pager[DenyAssignmentsClientListForScopeResponse]`
- New function `*DenyAssignmentsClient.NewListForResourcePager(string, string, string, string, string, *DenyAssignmentsClientListForResourceOptions) *runtime.Pager[DenyAssignmentsClientListForResourceResponse]`
- New function `*DenyAssignmentsClient.NewListForResourceGroupPager(string, *DenyAssignmentsClientListForResourceGroupOptions) *runtime.Pager[DenyAssignmentsClientListForResourceGroupResponse]`
- New function `*DenyAssignmentsClient.GetByID(context.Context, string, *DenyAssignmentsClientGetByIDOptions) (DenyAssignmentsClientGetByIDResponse, error)`
- New function `*DenyAssignmentsClient.NewListPager(*DenyAssignmentsClientListOptions) *runtime.Pager[DenyAssignmentsClientListResponse]`
- New function `NewDenyAssignmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DenyAssignmentsClient, error)`
- New struct `DenyAssignment`
- New struct `DenyAssignmentFilter`
- New struct `DenyAssignmentListResult`
- New struct `DenyAssignmentPermission`
- New struct `DenyAssignmentProperties`
- New struct `DenyAssignmentsClient`
- New struct `DenyAssignmentsClientGetByIDOptions`
- New struct `DenyAssignmentsClientGetByIDResponse`
- New struct `DenyAssignmentsClientGetOptions`
- New struct `DenyAssignmentsClientGetResponse`
- New struct `DenyAssignmentsClientListForResourceGroupOptions`
- New struct `DenyAssignmentsClientListForResourceGroupResponse`
- New struct `DenyAssignmentsClientListForResourceOptions`
- New struct `DenyAssignmentsClientListForResourceResponse`
- New struct `DenyAssignmentsClientListForScopeOptions`
- New struct `DenyAssignmentsClientListForScopeResponse`
- New struct `DenyAssignmentsClientListOptions`
- New struct `DenyAssignmentsClientListResponse`
- New struct `ValidationResponse`
- New struct `ValidationResponseErrorInfo`
- New field `TenantID` in struct `RoleAssignmentsClientGetByIDOptions`
- New field `DataActions` in struct `Permission`
- New field `NotDataActions` in struct `Permission`
- New field `TenantID` in struct `RoleAssignmentsClientListForResourceOptions`
- New field `UpdatedBy` in struct `RoleAssignmentProperties`
- New field `Condition` in struct `RoleAssignmentProperties`
- New field `CreatedOn` in struct `RoleAssignmentProperties`
- New field `UpdatedOn` in struct `RoleAssignmentProperties`
- New field `CreatedBy` in struct `RoleAssignmentProperties`
- New field `ConditionVersion` in struct `RoleAssignmentProperties`
- New field `DelegatedManagedIdentityResourceID` in struct `RoleAssignmentProperties`
- New field `Description` in struct `RoleAssignmentProperties`
- New field `PrincipalType` in struct `RoleAssignmentProperties`
- New field `Scope` in struct `RoleAssignmentProperties`
- New field `TenantID` in struct `RoleAssignmentsClientDeleteByIDOptions`
- New field `IsDataAction` in struct `ProviderOperation`
- New field `TenantID` in struct `RoleAssignmentsClientDeleteOptions`
- New field `Type` in struct `RoleDefinitionFilter`
- New field `TenantID` in struct `RoleAssignmentsClientListForResourceGroupOptions`
- New field `TenantID` in struct `RoleAssignmentsClientGetOptions`
- New field `SkipToken` in struct `RoleAssignmentsClientListForScopeOptions`
- New field `TenantID` in struct `RoleAssignmentsClientListForScopeOptions`


## 1.0.0 (2022-06-02)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
