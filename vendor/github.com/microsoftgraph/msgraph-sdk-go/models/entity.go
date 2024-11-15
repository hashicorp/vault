package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Entity struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewEntity instantiates a new Entity and sets the default values.
func NewEntity()(*Entity) {
    m := &Entity{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateEntityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEntityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.aadUserConversationMember":
                        return NewAadUserConversationMember(), nil
                    case "#microsoft.graph.accessPackage":
                        return NewAccessPackage(), nil
                    case "#microsoft.graph.accessPackageAssignment":
                        return NewAccessPackageAssignment(), nil
                    case "#microsoft.graph.accessPackageAssignmentPolicy":
                        return NewAccessPackageAssignmentPolicy(), nil
                    case "#microsoft.graph.accessPackageAssignmentRequest":
                        return NewAccessPackageAssignmentRequest(), nil
                    case "#microsoft.graph.accessPackageAssignmentRequestWorkflowExtension":
                        return NewAccessPackageAssignmentRequestWorkflowExtension(), nil
                    case "#microsoft.graph.accessPackageAssignmentWorkflowExtension":
                        return NewAccessPackageAssignmentWorkflowExtension(), nil
                    case "#microsoft.graph.accessPackageCatalog":
                        return NewAccessPackageCatalog(), nil
                    case "#microsoft.graph.accessPackageMultipleChoiceQuestion":
                        return NewAccessPackageMultipleChoiceQuestion(), nil
                    case "#microsoft.graph.accessPackageQuestion":
                        return NewAccessPackageQuestion(), nil
                    case "#microsoft.graph.accessPackageResource":
                        return NewAccessPackageResource(), nil
                    case "#microsoft.graph.accessPackageResourceEnvironment":
                        return NewAccessPackageResourceEnvironment(), nil
                    case "#microsoft.graph.accessPackageResourceRequest":
                        return NewAccessPackageResourceRequest(), nil
                    case "#microsoft.graph.accessPackageResourceRole":
                        return NewAccessPackageResourceRole(), nil
                    case "#microsoft.graph.accessPackageResourceRoleScope":
                        return NewAccessPackageResourceRoleScope(), nil
                    case "#microsoft.graph.accessPackageResourceScope":
                        return NewAccessPackageResourceScope(), nil
                    case "#microsoft.graph.accessPackageSubject":
                        return NewAccessPackageSubject(), nil
                    case "#microsoft.graph.accessPackageTextInputQuestion":
                        return NewAccessPackageTextInputQuestion(), nil
                    case "#microsoft.graph.accessReviewHistoryDefinition":
                        return NewAccessReviewHistoryDefinition(), nil
                    case "#microsoft.graph.accessReviewHistoryInstance":
                        return NewAccessReviewHistoryInstance(), nil
                    case "#microsoft.graph.accessReviewInstance":
                        return NewAccessReviewInstance(), nil
                    case "#microsoft.graph.accessReviewInstanceDecisionItem":
                        return NewAccessReviewInstanceDecisionItem(), nil
                    case "#microsoft.graph.accessReviewReviewer":
                        return NewAccessReviewReviewer(), nil
                    case "#microsoft.graph.accessReviewScheduleDefinition":
                        return NewAccessReviewScheduleDefinition(), nil
                    case "#microsoft.graph.accessReviewSet":
                        return NewAccessReviewSet(), nil
                    case "#microsoft.graph.accessReviewStage":
                        return NewAccessReviewStage(), nil
                    case "#microsoft.graph.activityBasedTimeoutPolicy":
                        return NewActivityBasedTimeoutPolicy(), nil
                    case "#microsoft.graph.activityHistoryItem":
                        return NewActivityHistoryItem(), nil
                    case "#microsoft.graph.addLargeGalleryViewOperation":
                        return NewAddLargeGalleryViewOperation(), nil
                    case "#microsoft.graph.adminConsentRequestPolicy":
                        return NewAdminConsentRequestPolicy(), nil
                    case "#microsoft.graph.administrativeUnit":
                        return NewAdministrativeUnit(), nil
                    case "#microsoft.graph.adminMicrosoft365Apps":
                        return NewAdminMicrosoft365Apps(), nil
                    case "#microsoft.graph.adminReportSettings":
                        return NewAdminReportSettings(), nil
                    case "#microsoft.graph.agreement":
                        return NewAgreement(), nil
                    case "#microsoft.graph.agreementAcceptance":
                        return NewAgreementAcceptance(), nil
                    case "#microsoft.graph.agreementFile":
                        return NewAgreementFile(), nil
                    case "#microsoft.graph.agreementFileLocalization":
                        return NewAgreementFileLocalization(), nil
                    case "#microsoft.graph.agreementFileProperties":
                        return NewAgreementFileProperties(), nil
                    case "#microsoft.graph.agreementFileVersion":
                        return NewAgreementFileVersion(), nil
                    case "#microsoft.graph.alert":
                        return NewAlert(), nil
                    case "#microsoft.graph.allowedValue":
                        return NewAllowedValue(), nil
                    case "#microsoft.graph.androidCompliancePolicy":
                        return NewAndroidCompliancePolicy(), nil
                    case "#microsoft.graph.androidCustomConfiguration":
                        return NewAndroidCustomConfiguration(), nil
                    case "#microsoft.graph.androidGeneralDeviceConfiguration":
                        return NewAndroidGeneralDeviceConfiguration(), nil
                    case "#microsoft.graph.androidLobApp":
                        return NewAndroidLobApp(), nil
                    case "#microsoft.graph.androidManagedAppProtection":
                        return NewAndroidManagedAppProtection(), nil
                    case "#microsoft.graph.androidManagedAppRegistration":
                        return NewAndroidManagedAppRegistration(), nil
                    case "#microsoft.graph.androidStoreApp":
                        return NewAndroidStoreApp(), nil
                    case "#microsoft.graph.androidWorkProfileCompliancePolicy":
                        return NewAndroidWorkProfileCompliancePolicy(), nil
                    case "#microsoft.graph.androidWorkProfileCustomConfiguration":
                        return NewAndroidWorkProfileCustomConfiguration(), nil
                    case "#microsoft.graph.androidWorkProfileGeneralDeviceConfiguration":
                        return NewAndroidWorkProfileGeneralDeviceConfiguration(), nil
                    case "#microsoft.graph.anonymousGuestConversationMember":
                        return NewAnonymousGuestConversationMember(), nil
                    case "#microsoft.graph.appCatalogs":
                        return NewAppCatalogs(), nil
                    case "#microsoft.graph.appConsentApprovalRoute":
                        return NewAppConsentApprovalRoute(), nil
                    case "#microsoft.graph.appConsentRequest":
                        return NewAppConsentRequest(), nil
                    case "#microsoft.graph.appleDeviceFeaturesConfigurationBase":
                        return NewAppleDeviceFeaturesConfigurationBase(), nil
                    case "#microsoft.graph.appleManagedIdentityProvider":
                        return NewAppleManagedIdentityProvider(), nil
                    case "#microsoft.graph.applePushNotificationCertificate":
                        return NewApplePushNotificationCertificate(), nil
                    case "#microsoft.graph.application":
                        return NewApplication(), nil
                    case "#microsoft.graph.applicationTemplate":
                        return NewApplicationTemplate(), nil
                    case "#microsoft.graph.appLogCollectionRequest":
                        return NewAppLogCollectionRequest(), nil
                    case "#microsoft.graph.appManagementPolicy":
                        return NewAppManagementPolicy(), nil
                    case "#microsoft.graph.appRoleAssignment":
                        return NewAppRoleAssignment(), nil
                    case "#microsoft.graph.approval":
                        return NewApproval(), nil
                    case "#microsoft.graph.approvalStage":
                        return NewApprovalStage(), nil
                    case "#microsoft.graph.appScope":
                        return NewAppScope(), nil
                    case "#microsoft.graph.associatedTeamInfo":
                        return NewAssociatedTeamInfo(), nil
                    case "#microsoft.graph.attachment":
                        return NewAttachment(), nil
                    case "#microsoft.graph.attachmentBase":
                        return NewAttachmentBase(), nil
                    case "#microsoft.graph.attachmentSession":
                        return NewAttachmentSession(), nil
                    case "#microsoft.graph.attackSimulationOperation":
                        return NewAttackSimulationOperation(), nil
                    case "#microsoft.graph.attackSimulationRoot":
                        return NewAttackSimulationRoot(), nil
                    case "#microsoft.graph.attendanceRecord":
                        return NewAttendanceRecord(), nil
                    case "#microsoft.graph.attributeMappingFunctionSchema":
                        return NewAttributeMappingFunctionSchema(), nil
                    case "#microsoft.graph.attributeSet":
                        return NewAttributeSet(), nil
                    case "#microsoft.graph.audioRoutingGroup":
                        return NewAudioRoutingGroup(), nil
                    case "#microsoft.graph.auditEvent":
                        return NewAuditEvent(), nil
                    case "#microsoft.graph.auditLogRoot":
                        return NewAuditLogRoot(), nil
                    case "#microsoft.graph.authentication":
                        return NewAuthentication(), nil
                    case "#microsoft.graph.authenticationCombinationConfiguration":
                        return NewAuthenticationCombinationConfiguration(), nil
                    case "#microsoft.graph.authenticationContextClassReference":
                        return NewAuthenticationContextClassReference(), nil
                    case "#microsoft.graph.authenticationEventListener":
                        return NewAuthenticationEventListener(), nil
                    case "#microsoft.graph.authenticationEventsFlow":
                        return NewAuthenticationEventsFlow(), nil
                    case "#microsoft.graph.authenticationFlowsPolicy":
                        return NewAuthenticationFlowsPolicy(), nil
                    case "#microsoft.graph.authenticationMethod":
                        return NewAuthenticationMethod(), nil
                    case "#microsoft.graph.authenticationMethodConfiguration":
                        return NewAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.authenticationMethodModeDetail":
                        return NewAuthenticationMethodModeDetail(), nil
                    case "#microsoft.graph.authenticationMethodsPolicy":
                        return NewAuthenticationMethodsPolicy(), nil
                    case "#microsoft.graph.authenticationMethodsRoot":
                        return NewAuthenticationMethodsRoot(), nil
                    case "#microsoft.graph.authenticationMethodTarget":
                        return NewAuthenticationMethodTarget(), nil
                    case "#microsoft.graph.authenticationStrengthPolicy":
                        return NewAuthenticationStrengthPolicy(), nil
                    case "#microsoft.graph.authenticationStrengthRoot":
                        return NewAuthenticationStrengthRoot(), nil
                    case "#microsoft.graph.authoredNote":
                        return NewAuthoredNote(), nil
                    case "#microsoft.graph.authorizationPolicy":
                        return NewAuthorizationPolicy(), nil
                    case "#microsoft.graph.azureCommunicationServicesUserConversationMember":
                        return NewAzureCommunicationServicesUserConversationMember(), nil
                    case "#microsoft.graph.b2xIdentityUserFlow":
                        return NewB2xIdentityUserFlow(), nil
                    case "#microsoft.graph.backupRestoreRoot":
                        return NewBackupRestoreRoot(), nil
                    case "#microsoft.graph.baseItem":
                        return NewBaseItem(), nil
                    case "#microsoft.graph.baseItemVersion":
                        return NewBaseItemVersion(), nil
                    case "#microsoft.graph.baseSitePage":
                        return NewBaseSitePage(), nil
                    case "#microsoft.graph.bitlocker":
                        return NewBitlocker(), nil
                    case "#microsoft.graph.bitlockerRecoveryKey":
                        return NewBitlockerRecoveryKey(), nil
                    case "#microsoft.graph.bookingAppointment":
                        return NewBookingAppointment(), nil
                    case "#microsoft.graph.bookingBusiness":
                        return NewBookingBusiness(), nil
                    case "#microsoft.graph.bookingCurrency":
                        return NewBookingCurrency(), nil
                    case "#microsoft.graph.bookingCustomer":
                        return NewBookingCustomer(), nil
                    case "#microsoft.graph.bookingCustomerBase":
                        return NewBookingCustomerBase(), nil
                    case "#microsoft.graph.bookingCustomQuestion":
                        return NewBookingCustomQuestion(), nil
                    case "#microsoft.graph.bookingService":
                        return NewBookingService(), nil
                    case "#microsoft.graph.bookingStaffMember":
                        return NewBookingStaffMember(), nil
                    case "#microsoft.graph.bookingStaffMemberBase":
                        return NewBookingStaffMemberBase(), nil
                    case "#microsoft.graph.browserSharedCookie":
                        return NewBrowserSharedCookie(), nil
                    case "#microsoft.graph.browserSite":
                        return NewBrowserSite(), nil
                    case "#microsoft.graph.browserSiteList":
                        return NewBrowserSiteList(), nil
                    case "#microsoft.graph.builtInIdentityProvider":
                        return NewBuiltInIdentityProvider(), nil
                    case "#microsoft.graph.bulkUpload":
                        return NewBulkUpload(), nil
                    case "#microsoft.graph.calendar":
                        return NewCalendar(), nil
                    case "#microsoft.graph.calendarGroup":
                        return NewCalendarGroup(), nil
                    case "#microsoft.graph.calendarPermission":
                        return NewCalendarPermission(), nil
                    case "#microsoft.graph.calendarSharingMessage":
                        return NewCalendarSharingMessage(), nil
                    case "#microsoft.graph.call":
                        return NewCall(), nil
                    case "#microsoft.graph.callRecording":
                        return NewCallRecording(), nil
                    case "#microsoft.graph.callTranscript":
                        return NewCallTranscript(), nil
                    case "#microsoft.graph.cancelMediaProcessingOperation":
                        return NewCancelMediaProcessingOperation(), nil
                    case "#microsoft.graph.canvasLayout":
                        return NewCanvasLayout(), nil
                    case "#microsoft.graph.certificateBasedAuthConfiguration":
                        return NewCertificateBasedAuthConfiguration(), nil
                    case "#microsoft.graph.changeTrackedEntity":
                        return NewChangeTrackedEntity(), nil
                    case "#microsoft.graph.channel":
                        return NewChannel(), nil
                    case "#microsoft.graph.chat":
                        return NewChat(), nil
                    case "#microsoft.graph.chatMessage":
                        return NewChatMessage(), nil
                    case "#microsoft.graph.chatMessageHostedContent":
                        return NewChatMessageHostedContent(), nil
                    case "#microsoft.graph.chatMessageInfo":
                        return NewChatMessageInfo(), nil
                    case "#microsoft.graph.checklistItem":
                        return NewChecklistItem(), nil
                    case "#microsoft.graph.claimsMappingPolicy":
                        return NewClaimsMappingPolicy(), nil
                    case "#microsoft.graph.cloudClipboardItem":
                        return NewCloudClipboardItem(), nil
                    case "#microsoft.graph.cloudClipboardRoot":
                        return NewCloudClipboardRoot(), nil
                    case "#microsoft.graph.cloudPC":
                        return NewCloudPC(), nil
                    case "#microsoft.graph.cloudPcAuditEvent":
                        return NewCloudPcAuditEvent(), nil
                    case "#microsoft.graph.cloudPcDeviceImage":
                        return NewCloudPcDeviceImage(), nil
                    case "#microsoft.graph.cloudPcGalleryImage":
                        return NewCloudPcGalleryImage(), nil
                    case "#microsoft.graph.cloudPcOnPremisesConnection":
                        return NewCloudPcOnPremisesConnection(), nil
                    case "#microsoft.graph.cloudPcProvisioningPolicy":
                        return NewCloudPcProvisioningPolicy(), nil
                    case "#microsoft.graph.cloudPcProvisioningPolicyAssignment":
                        return NewCloudPcProvisioningPolicyAssignment(), nil
                    case "#microsoft.graph.cloudPcUserSetting":
                        return NewCloudPcUserSetting(), nil
                    case "#microsoft.graph.cloudPcUserSettingAssignment":
                        return NewCloudPcUserSettingAssignment(), nil
                    case "#microsoft.graph.columnDefinition":
                        return NewColumnDefinition(), nil
                    case "#microsoft.graph.columnLink":
                        return NewColumnLink(), nil
                    case "#microsoft.graph.commsOperation":
                        return NewCommsOperation(), nil
                    case "#microsoft.graph.community":
                        return NewCommunity(), nil
                    case "#microsoft.graph.companySubscription":
                        return NewCompanySubscription(), nil
                    case "#microsoft.graph.complianceManagementPartner":
                        return NewComplianceManagementPartner(), nil
                    case "#microsoft.graph.conditionalAccessPolicy":
                        return NewConditionalAccessPolicy(), nil
                    case "#microsoft.graph.conditionalAccessRoot":
                        return NewConditionalAccessRoot(), nil
                    case "#microsoft.graph.conditionalAccessTemplate":
                        return NewConditionalAccessTemplate(), nil
                    case "#microsoft.graph.connectedOrganization":
                        return NewConnectedOrganization(), nil
                    case "#microsoft.graph.contact":
                        return NewContact(), nil
                    case "#microsoft.graph.contactFolder":
                        return NewContactFolder(), nil
                    case "#microsoft.graph.contentSharingSession":
                        return NewContentSharingSession(), nil
                    case "#microsoft.graph.contentType":
                        return NewContentType(), nil
                    case "#microsoft.graph.contract":
                        return NewContract(), nil
                    case "#microsoft.graph.conversation":
                        return NewConversation(), nil
                    case "#microsoft.graph.conversationMember":
                        return NewConversationMember(), nil
                    case "#microsoft.graph.conversationThread":
                        return NewConversationThread(), nil
                    case "#microsoft.graph.countryNamedLocation":
                        return NewCountryNamedLocation(), nil
                    case "#microsoft.graph.crossTenantAccessPolicy":
                        return NewCrossTenantAccessPolicy(), nil
                    case "#microsoft.graph.crossTenantAccessPolicyConfigurationDefault":
                        return NewCrossTenantAccessPolicyConfigurationDefault(), nil
                    case "#microsoft.graph.customAuthenticationExtension":
                        return NewCustomAuthenticationExtension(), nil
                    case "#microsoft.graph.customCalloutExtension":
                        return NewCustomCalloutExtension(), nil
                    case "#microsoft.graph.customExtensionStageSetting":
                        return NewCustomExtensionStageSetting(), nil
                    case "#microsoft.graph.customSecurityAttributeDefinition":
                        return NewCustomSecurityAttributeDefinition(), nil
                    case "#microsoft.graph.dataPolicyOperation":
                        return NewDataPolicyOperation(), nil
                    case "#microsoft.graph.defaultManagedAppProtection":
                        return NewDefaultManagedAppProtection(), nil
                    case "#microsoft.graph.delegatedAdminAccessAssignment":
                        return NewDelegatedAdminAccessAssignment(), nil
                    case "#microsoft.graph.delegatedAdminCustomer":
                        return NewDelegatedAdminCustomer(), nil
                    case "#microsoft.graph.delegatedAdminRelationship":
                        return NewDelegatedAdminRelationship(), nil
                    case "#microsoft.graph.delegatedAdminRelationshipOperation":
                        return NewDelegatedAdminRelationshipOperation(), nil
                    case "#microsoft.graph.delegatedAdminRelationshipRequest":
                        return NewDelegatedAdminRelationshipRequest(), nil
                    case "#microsoft.graph.delegatedAdminServiceManagementDetail":
                        return NewDelegatedAdminServiceManagementDetail(), nil
                    case "#microsoft.graph.delegatedPermissionClassification":
                        return NewDelegatedPermissionClassification(), nil
                    case "#microsoft.graph.deletedChat":
                        return NewDeletedChat(), nil
                    case "#microsoft.graph.deletedItemContainer":
                        return NewDeletedItemContainer(), nil
                    case "#microsoft.graph.deletedTeam":
                        return NewDeletedTeam(), nil
                    case "#microsoft.graph.deltaParticipants":
                        return NewDeltaParticipants(), nil
                    case "#microsoft.graph.detectedApp":
                        return NewDetectedApp(), nil
                    case "#microsoft.graph.device":
                        return NewDevice(), nil
                    case "#microsoft.graph.deviceAndAppManagementRoleAssignment":
                        return NewDeviceAndAppManagementRoleAssignment(), nil
                    case "#microsoft.graph.deviceAndAppManagementRoleDefinition":
                        return NewDeviceAndAppManagementRoleDefinition(), nil
                    case "#microsoft.graph.deviceAppManagement":
                        return NewDeviceAppManagement(), nil
                    case "#microsoft.graph.deviceCategory":
                        return NewDeviceCategory(), nil
                    case "#microsoft.graph.deviceComplianceActionItem":
                        return NewDeviceComplianceActionItem(), nil
                    case "#microsoft.graph.deviceComplianceDeviceOverview":
                        return NewDeviceComplianceDeviceOverview(), nil
                    case "#microsoft.graph.deviceComplianceDeviceStatus":
                        return NewDeviceComplianceDeviceStatus(), nil
                    case "#microsoft.graph.deviceCompliancePolicy":
                        return NewDeviceCompliancePolicy(), nil
                    case "#microsoft.graph.deviceCompliancePolicyAssignment":
                        return NewDeviceCompliancePolicyAssignment(), nil
                    case "#microsoft.graph.deviceCompliancePolicyDeviceStateSummary":
                        return NewDeviceCompliancePolicyDeviceStateSummary(), nil
                    case "#microsoft.graph.deviceCompliancePolicySettingStateSummary":
                        return NewDeviceCompliancePolicySettingStateSummary(), nil
                    case "#microsoft.graph.deviceCompliancePolicyState":
                        return NewDeviceCompliancePolicyState(), nil
                    case "#microsoft.graph.deviceComplianceScheduledActionForRule":
                        return NewDeviceComplianceScheduledActionForRule(), nil
                    case "#microsoft.graph.deviceComplianceSettingState":
                        return NewDeviceComplianceSettingState(), nil
                    case "#microsoft.graph.deviceComplianceUserOverview":
                        return NewDeviceComplianceUserOverview(), nil
                    case "#microsoft.graph.deviceComplianceUserStatus":
                        return NewDeviceComplianceUserStatus(), nil
                    case "#microsoft.graph.deviceConfiguration":
                        return NewDeviceConfiguration(), nil
                    case "#microsoft.graph.deviceConfigurationAssignment":
                        return NewDeviceConfigurationAssignment(), nil
                    case "#microsoft.graph.deviceConfigurationDeviceOverview":
                        return NewDeviceConfigurationDeviceOverview(), nil
                    case "#microsoft.graph.deviceConfigurationDeviceStateSummary":
                        return NewDeviceConfigurationDeviceStateSummary(), nil
                    case "#microsoft.graph.deviceConfigurationDeviceStatus":
                        return NewDeviceConfigurationDeviceStatus(), nil
                    case "#microsoft.graph.deviceConfigurationState":
                        return NewDeviceConfigurationState(), nil
                    case "#microsoft.graph.deviceConfigurationUserOverview":
                        return NewDeviceConfigurationUserOverview(), nil
                    case "#microsoft.graph.deviceConfigurationUserStatus":
                        return NewDeviceConfigurationUserStatus(), nil
                    case "#microsoft.graph.deviceEnrollmentConfiguration":
                        return NewDeviceEnrollmentConfiguration(), nil
                    case "#microsoft.graph.deviceEnrollmentLimitConfiguration":
                        return NewDeviceEnrollmentLimitConfiguration(), nil
                    case "#microsoft.graph.deviceEnrollmentPlatformRestrictionsConfiguration":
                        return NewDeviceEnrollmentPlatformRestrictionsConfiguration(), nil
                    case "#microsoft.graph.deviceEnrollmentWindowsHelloForBusinessConfiguration":
                        return NewDeviceEnrollmentWindowsHelloForBusinessConfiguration(), nil
                    case "#microsoft.graph.deviceInstallState":
                        return NewDeviceInstallState(), nil
                    case "#microsoft.graph.deviceLocalCredentialInfo":
                        return NewDeviceLocalCredentialInfo(), nil
                    case "#microsoft.graph.deviceLogCollectionResponse":
                        return NewDeviceLogCollectionResponse(), nil
                    case "#microsoft.graph.deviceManagement":
                        return NewDeviceManagement(), nil
                    case "#microsoft.graph.deviceManagementCachedReportConfiguration":
                        return NewDeviceManagementCachedReportConfiguration(), nil
                    case "#microsoft.graph.deviceManagementExchangeConnector":
                        return NewDeviceManagementExchangeConnector(), nil
                    case "#microsoft.graph.deviceManagementExportJob":
                        return NewDeviceManagementExportJob(), nil
                    case "#microsoft.graph.deviceManagementPartner":
                        return NewDeviceManagementPartner(), nil
                    case "#microsoft.graph.deviceManagementReports":
                        return NewDeviceManagementReports(), nil
                    case "#microsoft.graph.deviceManagementTroubleshootingEvent":
                        return NewDeviceManagementTroubleshootingEvent(), nil
                    case "#microsoft.graph.deviceRegistrationPolicy":
                        return NewDeviceRegistrationPolicy(), nil
                    case "#microsoft.graph.directory":
                        return NewDirectory(), nil
                    case "#microsoft.graph.directoryAudit":
                        return NewDirectoryAudit(), nil
                    case "#microsoft.graph.directoryDefinition":
                        return NewDirectoryDefinition(), nil
                    case "#microsoft.graph.directoryObject":
                        return NewDirectoryObject(), nil
                    case "#microsoft.graph.directoryObjectPartnerReference":
                        return NewDirectoryObjectPartnerReference(), nil
                    case "#microsoft.graph.directoryRole":
                        return NewDirectoryRole(), nil
                    case "#microsoft.graph.directoryRoleTemplate":
                        return NewDirectoryRoleTemplate(), nil
                    case "#microsoft.graph.documentSetVersion":
                        return NewDocumentSetVersion(), nil
                    case "#microsoft.graph.domain":
                        return NewDomain(), nil
                    case "#microsoft.graph.domainDnsCnameRecord":
                        return NewDomainDnsCnameRecord(), nil
                    case "#microsoft.graph.domainDnsMxRecord":
                        return NewDomainDnsMxRecord(), nil
                    case "#microsoft.graph.domainDnsRecord":
                        return NewDomainDnsRecord(), nil
                    case "#microsoft.graph.domainDnsSrvRecord":
                        return NewDomainDnsSrvRecord(), nil
                    case "#microsoft.graph.domainDnsTxtRecord":
                        return NewDomainDnsTxtRecord(), nil
                    case "#microsoft.graph.domainDnsUnavailableRecord":
                        return NewDomainDnsUnavailableRecord(), nil
                    case "#microsoft.graph.drive":
                        return NewDrive(), nil
                    case "#microsoft.graph.driveItem":
                        return NewDriveItem(), nil
                    case "#microsoft.graph.driveItemVersion":
                        return NewDriveItemVersion(), nil
                    case "#microsoft.graph.driveProtectionRule":
                        return NewDriveProtectionRule(), nil
                    case "#microsoft.graph.driveProtectionUnit":
                        return NewDriveProtectionUnit(), nil
                    case "#microsoft.graph.driveRestoreArtifact":
                        return NewDriveRestoreArtifact(), nil
                    case "#microsoft.graph.eBookInstallSummary":
                        return NewEBookInstallSummary(), nil
                    case "#microsoft.graph.edge":
                        return NewEdge(), nil
                    case "#microsoft.graph.editionUpgradeConfiguration":
                        return NewEditionUpgradeConfiguration(), nil
                    case "#microsoft.graph.educationAssignment":
                        return NewEducationAssignment(), nil
                    case "#microsoft.graph.educationAssignmentDefaults":
                        return NewEducationAssignmentDefaults(), nil
                    case "#microsoft.graph.educationAssignmentResource":
                        return NewEducationAssignmentResource(), nil
                    case "#microsoft.graph.educationAssignmentSettings":
                        return NewEducationAssignmentSettings(), nil
                    case "#microsoft.graph.educationCategory":
                        return NewEducationCategory(), nil
                    case "#microsoft.graph.educationClass":
                        return NewEducationClass(), nil
                    case "#microsoft.graph.educationFeedbackOutcome":
                        return NewEducationFeedbackOutcome(), nil
                    case "#microsoft.graph.educationFeedbackResourceOutcome":
                        return NewEducationFeedbackResourceOutcome(), nil
                    case "#microsoft.graph.educationGradingCategory":
                        return NewEducationGradingCategory(), nil
                    case "#microsoft.graph.educationModule":
                        return NewEducationModule(), nil
                    case "#microsoft.graph.educationModuleResource":
                        return NewEducationModuleResource(), nil
                    case "#microsoft.graph.educationOrganization":
                        return NewEducationOrganization(), nil
                    case "#microsoft.graph.educationOutcome":
                        return NewEducationOutcome(), nil
                    case "#microsoft.graph.educationPointsOutcome":
                        return NewEducationPointsOutcome(), nil
                    case "#microsoft.graph.educationRubric":
                        return NewEducationRubric(), nil
                    case "#microsoft.graph.educationRubricOutcome":
                        return NewEducationRubricOutcome(), nil
                    case "#microsoft.graph.educationSchool":
                        return NewEducationSchool(), nil
                    case "#microsoft.graph.educationSubmission":
                        return NewEducationSubmission(), nil
                    case "#microsoft.graph.educationSubmissionResource":
                        return NewEducationSubmissionResource(), nil
                    case "#microsoft.graph.educationUser":
                        return NewEducationUser(), nil
                    case "#microsoft.graph.emailAuthenticationMethod":
                        return NewEmailAuthenticationMethod(), nil
                    case "#microsoft.graph.emailAuthenticationMethodConfiguration":
                        return NewEmailAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.emailFileAssessmentRequest":
                        return NewEmailFileAssessmentRequest(), nil
                    case "#microsoft.graph.employeeExperienceUser":
                        return NewEmployeeExperienceUser(), nil
                    case "#microsoft.graph.endpoint":
                        return NewEndpoint(), nil
                    case "#microsoft.graph.endUserNotification":
                        return NewEndUserNotification(), nil
                    case "#microsoft.graph.endUserNotificationDetail":
                        return NewEndUserNotificationDetail(), nil
                    case "#microsoft.graph.engagementAsyncOperation":
                        return NewEngagementAsyncOperation(), nil
                    case "#microsoft.graph.enrollmentConfigurationAssignment":
                        return NewEnrollmentConfigurationAssignment(), nil
                    case "#microsoft.graph.enrollmentTroubleshootingEvent":
                        return NewEnrollmentTroubleshootingEvent(), nil
                    case "#microsoft.graph.enterpriseCodeSigningCertificate":
                        return NewEnterpriseCodeSigningCertificate(), nil
                    case "#microsoft.graph.entitlementManagement":
                        return NewEntitlementManagement(), nil
                    case "#microsoft.graph.entitlementManagementSettings":
                        return NewEntitlementManagementSettings(), nil
                    case "#microsoft.graph.event":
                        return NewEvent(), nil
                    case "#microsoft.graph.eventMessage":
                        return NewEventMessage(), nil
                    case "#microsoft.graph.eventMessageRequest":
                        return NewEventMessageRequest(), nil
                    case "#microsoft.graph.eventMessageResponse":
                        return NewEventMessageResponse(), nil
                    case "#microsoft.graph.exchangeProtectionPolicy":
                        return NewExchangeProtectionPolicy(), nil
                    case "#microsoft.graph.exchangeRestoreSession":
                        return NewExchangeRestoreSession(), nil
                    case "#microsoft.graph.extension":
                        return NewExtension(), nil
                    case "#microsoft.graph.extensionProperty":
                        return NewExtensionProperty(), nil
                    case "#microsoft.graph.externalDomainName":
                        return NewExternalDomainName(), nil
                    case "#microsoft.graph.externalUsersSelfServiceSignUpEventsFlow":
                        return NewExternalUsersSelfServiceSignUpEventsFlow(), nil
                    case "#microsoft.graph.featureRolloutPolicy":
                        return NewFeatureRolloutPolicy(), nil
                    case "#microsoft.graph.federatedIdentityCredential":
                        return NewFederatedIdentityCredential(), nil
                    case "#microsoft.graph.fido2AuthenticationMethod":
                        return NewFido2AuthenticationMethod(), nil
                    case "#microsoft.graph.fido2AuthenticationMethodConfiguration":
                        return NewFido2AuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.fido2CombinationConfiguration":
                        return NewFido2CombinationConfiguration(), nil
                    case "#microsoft.graph.fieldValueSet":
                        return NewFieldValueSet(), nil
                    case "#microsoft.graph.fileAssessmentRequest":
                        return NewFileAssessmentRequest(), nil
                    case "#microsoft.graph.fileAttachment":
                        return NewFileAttachment(), nil
                    case "#microsoft.graph.fileStorage":
                        return NewFileStorage(), nil
                    case "#microsoft.graph.fileStorageContainer":
                        return NewFileStorageContainer(), nil
                    case "#microsoft.graph.filterOperatorSchema":
                        return NewFilterOperatorSchema(), nil
                    case "#microsoft.graph.governanceInsight":
                        return NewGovernanceInsight(), nil
                    case "#microsoft.graph.granularMailboxRestoreArtifact":
                        return NewGranularMailboxRestoreArtifact(), nil
                    case "#microsoft.graph.group":
                        return NewGroup(), nil
                    case "#microsoft.graph.groupLifecyclePolicy":
                        return NewGroupLifecyclePolicy(), nil
                    case "#microsoft.graph.groupSetting":
                        return NewGroupSetting(), nil
                    case "#microsoft.graph.groupSettingTemplate":
                        return NewGroupSettingTemplate(), nil
                    case "#microsoft.graph.homeRealmDiscoveryPolicy":
                        return NewHomeRealmDiscoveryPolicy(), nil
                    case "#microsoft.graph.horizontalSection":
                        return NewHorizontalSection(), nil
                    case "#microsoft.graph.horizontalSectionColumn":
                        return NewHorizontalSectionColumn(), nil
                    case "#microsoft.graph.identityApiConnector":
                        return NewIdentityApiConnector(), nil
                    case "#microsoft.graph.identityBuiltInUserFlowAttribute":
                        return NewIdentityBuiltInUserFlowAttribute(), nil
                    case "#microsoft.graph.identityContainer":
                        return NewIdentityContainer(), nil
                    case "#microsoft.graph.identityCustomUserFlowAttribute":
                        return NewIdentityCustomUserFlowAttribute(), nil
                    case "#microsoft.graph.identityProvider":
                        return NewIdentityProvider(), nil
                    case "#microsoft.graph.identityProviderBase":
                        return NewIdentityProviderBase(), nil
                    case "#microsoft.graph.identitySecurityDefaultsEnforcementPolicy":
                        return NewIdentitySecurityDefaultsEnforcementPolicy(), nil
                    case "#microsoft.graph.identityUserFlow":
                        return NewIdentityUserFlow(), nil
                    case "#microsoft.graph.identityUserFlowAttribute":
                        return NewIdentityUserFlowAttribute(), nil
                    case "#microsoft.graph.identityUserFlowAttributeAssignment":
                        return NewIdentityUserFlowAttributeAssignment(), nil
                    case "#microsoft.graph.importedWindowsAutopilotDeviceIdentity":
                        return NewImportedWindowsAutopilotDeviceIdentity(), nil
                    case "#microsoft.graph.importedWindowsAutopilotDeviceIdentityUpload":
                        return NewImportedWindowsAutopilotDeviceIdentityUpload(), nil
                    case "#microsoft.graph.inferenceClassification":
                        return NewInferenceClassification(), nil
                    case "#microsoft.graph.inferenceClassificationOverride":
                        return NewInferenceClassificationOverride(), nil
                    case "#microsoft.graph.insightsSettings":
                        return NewInsightsSettings(), nil
                    case "#microsoft.graph.internalDomainFederation":
                        return NewInternalDomainFederation(), nil
                    case "#microsoft.graph.internetExplorerMode":
                        return NewInternetExplorerMode(), nil
                    case "#microsoft.graph.invitation":
                        return NewInvitation(), nil
                    case "#microsoft.graph.inviteParticipantsOperation":
                        return NewInviteParticipantsOperation(), nil
                    case "#microsoft.graph.iosCertificateProfile":
                        return NewIosCertificateProfile(), nil
                    case "#microsoft.graph.iosCompliancePolicy":
                        return NewIosCompliancePolicy(), nil
                    case "#microsoft.graph.iosCustomConfiguration":
                        return NewIosCustomConfiguration(), nil
                    case "#microsoft.graph.iosDeviceFeaturesConfiguration":
                        return NewIosDeviceFeaturesConfiguration(), nil
                    case "#microsoft.graph.iosGeneralDeviceConfiguration":
                        return NewIosGeneralDeviceConfiguration(), nil
                    case "#microsoft.graph.iosiPadOSWebClip":
                        return NewIosiPadOSWebClip(), nil
                    case "#microsoft.graph.iosLobApp":
                        return NewIosLobApp(), nil
                    case "#microsoft.graph.iosLobAppProvisioningConfigurationAssignment":
                        return NewIosLobAppProvisioningConfigurationAssignment(), nil
                    case "#microsoft.graph.iosManagedAppProtection":
                        return NewIosManagedAppProtection(), nil
                    case "#microsoft.graph.iosManagedAppRegistration":
                        return NewIosManagedAppRegistration(), nil
                    case "#microsoft.graph.iosMobileAppConfiguration":
                        return NewIosMobileAppConfiguration(), nil
                    case "#microsoft.graph.iosStoreApp":
                        return NewIosStoreApp(), nil
                    case "#microsoft.graph.iosUpdateConfiguration":
                        return NewIosUpdateConfiguration(), nil
                    case "#microsoft.graph.iosUpdateDeviceStatus":
                        return NewIosUpdateDeviceStatus(), nil
                    case "#microsoft.graph.iosVppApp":
                        return NewIosVppApp(), nil
                    case "#microsoft.graph.iosVppEBook":
                        return NewIosVppEBook(), nil
                    case "#microsoft.graph.iosVppEBookAssignment":
                        return NewIosVppEBookAssignment(), nil
                    case "#microsoft.graph.ipNamedLocation":
                        return NewIpNamedLocation(), nil
                    case "#microsoft.graph.itemActivity":
                        return NewItemActivity(), nil
                    case "#microsoft.graph.itemActivityStat":
                        return NewItemActivityStat(), nil
                    case "#microsoft.graph.itemAnalytics":
                        return NewItemAnalytics(), nil
                    case "#microsoft.graph.itemAttachment":
                        return NewItemAttachment(), nil
                    case "#microsoft.graph.itemInsights":
                        return NewItemInsights(), nil
                    case "#microsoft.graph.itemRetentionLabel":
                        return NewItemRetentionLabel(), nil
                    case "#microsoft.graph.landingPage":
                        return NewLandingPage(), nil
                    case "#microsoft.graph.landingPageDetail":
                        return NewLandingPageDetail(), nil
                    case "#microsoft.graph.learningAssignment":
                        return NewLearningAssignment(), nil
                    case "#microsoft.graph.learningContent":
                        return NewLearningContent(), nil
                    case "#microsoft.graph.learningCourseActivity":
                        return NewLearningCourseActivity(), nil
                    case "#microsoft.graph.learningProvider":
                        return NewLearningProvider(), nil
                    case "#microsoft.graph.learningSelfInitiatedCourse":
                        return NewLearningSelfInitiatedCourse(), nil
                    case "#microsoft.graph.licenseDetails":
                        return NewLicenseDetails(), nil
                    case "#microsoft.graph.linkedResource":
                        return NewLinkedResource(), nil
                    case "#microsoft.graph.list":
                        return NewList(), nil
                    case "#microsoft.graph.listItem":
                        return NewListItem(), nil
                    case "#microsoft.graph.listItemVersion":
                        return NewListItemVersion(), nil
                    case "#microsoft.graph.localizedNotificationMessage":
                        return NewLocalizedNotificationMessage(), nil
                    case "#microsoft.graph.loginPage":
                        return NewLoginPage(), nil
                    case "#microsoft.graph.longRunningOperation":
                        return NewLongRunningOperation(), nil
                    case "#microsoft.graph.m365AppsInstallationOptions":
                        return NewM365AppsInstallationOptions(), nil
                    case "#microsoft.graph.macOSCompliancePolicy":
                        return NewMacOSCompliancePolicy(), nil
                    case "#microsoft.graph.macOSCustomConfiguration":
                        return NewMacOSCustomConfiguration(), nil
                    case "#microsoft.graph.macOSDeviceFeaturesConfiguration":
                        return NewMacOSDeviceFeaturesConfiguration(), nil
                    case "#microsoft.graph.macOSDmgApp":
                        return NewMacOSDmgApp(), nil
                    case "#microsoft.graph.macOSGeneralDeviceConfiguration":
                        return NewMacOSGeneralDeviceConfiguration(), nil
                    case "#microsoft.graph.macOSLobApp":
                        return NewMacOSLobApp(), nil
                    case "#microsoft.graph.macOSMicrosoftDefenderApp":
                        return NewMacOSMicrosoftDefenderApp(), nil
                    case "#microsoft.graph.macOSMicrosoftEdgeApp":
                        return NewMacOSMicrosoftEdgeApp(), nil
                    case "#microsoft.graph.macOSOfficeSuiteApp":
                        return NewMacOSOfficeSuiteApp(), nil
                    case "#microsoft.graph.mailAssessmentRequest":
                        return NewMailAssessmentRequest(), nil
                    case "#microsoft.graph.mailboxProtectionRule":
                        return NewMailboxProtectionRule(), nil
                    case "#microsoft.graph.mailboxProtectionUnit":
                        return NewMailboxProtectionUnit(), nil
                    case "#microsoft.graph.mailboxRestoreArtifact":
                        return NewMailboxRestoreArtifact(), nil
                    case "#microsoft.graph.mailFolder":
                        return NewMailFolder(), nil
                    case "#microsoft.graph.mailSearchFolder":
                        return NewMailSearchFolder(), nil
                    case "#microsoft.graph.malwareStateForWindowsDevice":
                        return NewMalwareStateForWindowsDevice(), nil
                    case "#microsoft.graph.managedAndroidLobApp":
                        return NewManagedAndroidLobApp(), nil
                    case "#microsoft.graph.managedAndroidStoreApp":
                        return NewManagedAndroidStoreApp(), nil
                    case "#microsoft.graph.managedApp":
                        return NewManagedApp(), nil
                    case "#microsoft.graph.managedAppConfiguration":
                        return NewManagedAppConfiguration(), nil
                    case "#microsoft.graph.managedAppOperation":
                        return NewManagedAppOperation(), nil
                    case "#microsoft.graph.managedAppPolicy":
                        return NewManagedAppPolicy(), nil
                    case "#microsoft.graph.managedAppPolicyDeploymentSummary":
                        return NewManagedAppPolicyDeploymentSummary(), nil
                    case "#microsoft.graph.managedAppProtection":
                        return NewManagedAppProtection(), nil
                    case "#microsoft.graph.managedAppRegistration":
                        return NewManagedAppRegistration(), nil
                    case "#microsoft.graph.managedAppStatus":
                        return NewManagedAppStatus(), nil
                    case "#microsoft.graph.managedAppStatusRaw":
                        return NewManagedAppStatusRaw(), nil
                    case "#microsoft.graph.managedDevice":
                        return NewManagedDevice(), nil
                    case "#microsoft.graph.managedDeviceMobileAppConfiguration":
                        return NewManagedDeviceMobileAppConfiguration(), nil
                    case "#microsoft.graph.managedDeviceMobileAppConfigurationAssignment":
                        return NewManagedDeviceMobileAppConfigurationAssignment(), nil
                    case "#microsoft.graph.managedDeviceMobileAppConfigurationDeviceStatus":
                        return NewManagedDeviceMobileAppConfigurationDeviceStatus(), nil
                    case "#microsoft.graph.managedDeviceMobileAppConfigurationDeviceSummary":
                        return NewManagedDeviceMobileAppConfigurationDeviceSummary(), nil
                    case "#microsoft.graph.managedDeviceMobileAppConfigurationUserStatus":
                        return NewManagedDeviceMobileAppConfigurationUserStatus(), nil
                    case "#microsoft.graph.managedDeviceMobileAppConfigurationUserSummary":
                        return NewManagedDeviceMobileAppConfigurationUserSummary(), nil
                    case "#microsoft.graph.managedDeviceOverview":
                        return NewManagedDeviceOverview(), nil
                    case "#microsoft.graph.managedEBook":
                        return NewManagedEBook(), nil
                    case "#microsoft.graph.managedEBookAssignment":
                        return NewManagedEBookAssignment(), nil
                    case "#microsoft.graph.managedIOSLobApp":
                        return NewManagedIOSLobApp(), nil
                    case "#microsoft.graph.managedIOSStoreApp":
                        return NewManagedIOSStoreApp(), nil
                    case "#microsoft.graph.managedMobileApp":
                        return NewManagedMobileApp(), nil
                    case "#microsoft.graph.managedMobileLobApp":
                        return NewManagedMobileLobApp(), nil
                    case "#microsoft.graph.mdmWindowsInformationProtectionPolicy":
                        return NewMdmWindowsInformationProtectionPolicy(), nil
                    case "#microsoft.graph.meetingAttendanceReport":
                        return NewMeetingAttendanceReport(), nil
                    case "#microsoft.graph.membershipOutlierInsight":
                        return NewMembershipOutlierInsight(), nil
                    case "#microsoft.graph.message":
                        return NewMessage(), nil
                    case "#microsoft.graph.messageRule":
                        return NewMessageRule(), nil
                    case "#microsoft.graph.microsoftAccountUserConversationMember":
                        return NewMicrosoftAccountUserConversationMember(), nil
                    case "#microsoft.graph.microsoftAuthenticatorAuthenticationMethod":
                        return NewMicrosoftAuthenticatorAuthenticationMethod(), nil
                    case "#microsoft.graph.microsoftAuthenticatorAuthenticationMethodConfiguration":
                        return NewMicrosoftAuthenticatorAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.microsoftAuthenticatorAuthenticationMethodTarget":
                        return NewMicrosoftAuthenticatorAuthenticationMethodTarget(), nil
                    case "#microsoft.graph.microsoftStoreForBusinessApp":
                        return NewMicrosoftStoreForBusinessApp(), nil
                    case "#microsoft.graph.mobileApp":
                        return NewMobileApp(), nil
                    case "#microsoft.graph.mobileAppAssignment":
                        return NewMobileAppAssignment(), nil
                    case "#microsoft.graph.mobileAppCategory":
                        return NewMobileAppCategory(), nil
                    case "#microsoft.graph.mobileAppContent":
                        return NewMobileAppContent(), nil
                    case "#microsoft.graph.mobileAppContentFile":
                        return NewMobileAppContentFile(), nil
                    case "#microsoft.graph.mobileAppTroubleshootingEvent":
                        return NewMobileAppTroubleshootingEvent(), nil
                    case "#microsoft.graph.mobileContainedApp":
                        return NewMobileContainedApp(), nil
                    case "#microsoft.graph.mobileLobApp":
                        return NewMobileLobApp(), nil
                    case "#microsoft.graph.mobileThreatDefenseConnector":
                        return NewMobileThreatDefenseConnector(), nil
                    case "#microsoft.graph.multiTenantOrganization":
                        return NewMultiTenantOrganization(), nil
                    case "#microsoft.graph.multiTenantOrganizationIdentitySyncPolicyTemplate":
                        return NewMultiTenantOrganizationIdentitySyncPolicyTemplate(), nil
                    case "#microsoft.graph.multiTenantOrganizationJoinRequestRecord":
                        return NewMultiTenantOrganizationJoinRequestRecord(), nil
                    case "#microsoft.graph.multiTenantOrganizationMember":
                        return NewMultiTenantOrganizationMember(), nil
                    case "#microsoft.graph.multiTenantOrganizationPartnerConfigurationTemplate":
                        return NewMultiTenantOrganizationPartnerConfigurationTemplate(), nil
                    case "#microsoft.graph.multiValueLegacyExtendedProperty":
                        return NewMultiValueLegacyExtendedProperty(), nil
                    case "#microsoft.graph.muteParticipantOperation":
                        return NewMuteParticipantOperation(), nil
                    case "#microsoft.graph.namedLocation":
                        return NewNamedLocation(), nil
                    case "#microsoft.graph.notebook":
                        return NewNotebook(), nil
                    case "#microsoft.graph.notificationMessageTemplate":
                        return NewNotificationMessageTemplate(), nil
                    case "#microsoft.graph.oAuth2PermissionGrant":
                        return NewOAuth2PermissionGrant(), nil
                    case "#microsoft.graph.offerShiftRequest":
                        return NewOfferShiftRequest(), nil
                    case "#microsoft.graph.officeGraphInsights":
                        return NewOfficeGraphInsights(), nil
                    case "#microsoft.graph.onAttributeCollectionListener":
                        return NewOnAttributeCollectionListener(), nil
                    case "#microsoft.graph.onAuthenticationMethodLoadStartListener":
                        return NewOnAuthenticationMethodLoadStartListener(), nil
                    case "#microsoft.graph.oneDriveForBusinessProtectionPolicy":
                        return NewOneDriveForBusinessProtectionPolicy(), nil
                    case "#microsoft.graph.oneDriveForBusinessRestoreSession":
                        return NewOneDriveForBusinessRestoreSession(), nil
                    case "#microsoft.graph.onenote":
                        return NewOnenote(), nil
                    case "#microsoft.graph.onenoteEntityBaseModel":
                        return NewOnenoteEntityBaseModel(), nil
                    case "#microsoft.graph.onenoteEntityHierarchyModel":
                        return NewOnenoteEntityHierarchyModel(), nil
                    case "#microsoft.graph.onenoteEntitySchemaObjectModel":
                        return NewOnenoteEntitySchemaObjectModel(), nil
                    case "#microsoft.graph.onenoteOperation":
                        return NewOnenoteOperation(), nil
                    case "#microsoft.graph.onenotePage":
                        return NewOnenotePage(), nil
                    case "#microsoft.graph.onenoteResource":
                        return NewOnenoteResource(), nil
                    case "#microsoft.graph.onenoteSection":
                        return NewOnenoteSection(), nil
                    case "#microsoft.graph.onInteractiveAuthFlowStartListener":
                        return NewOnInteractiveAuthFlowStartListener(), nil
                    case "#microsoft.graph.onlineMeeting":
                        return NewOnlineMeeting(), nil
                    case "#microsoft.graph.onlineMeetingBase":
                        return NewOnlineMeetingBase(), nil
                    case "#microsoft.graph.onPremisesConditionalAccessSettings":
                        return NewOnPremisesConditionalAccessSettings(), nil
                    case "#microsoft.graph.onPremisesDirectorySynchronization":
                        return NewOnPremisesDirectorySynchronization(), nil
                    case "#microsoft.graph.onTokenIssuanceStartCustomExtension":
                        return NewOnTokenIssuanceStartCustomExtension(), nil
                    case "#microsoft.graph.onTokenIssuanceStartListener":
                        return NewOnTokenIssuanceStartListener(), nil
                    case "#microsoft.graph.onUserCreateStartListener":
                        return NewOnUserCreateStartListener(), nil
                    case "#microsoft.graph.openShift":
                        return NewOpenShift(), nil
                    case "#microsoft.graph.openShiftChangeRequest":
                        return NewOpenShiftChangeRequest(), nil
                    case "#microsoft.graph.openTypeExtension":
                        return NewOpenTypeExtension(), nil
                    case "#microsoft.graph.operation":
                        return NewOperation(), nil
                    case "#microsoft.graph.organization":
                        return NewOrganization(), nil
                    case "#microsoft.graph.organizationalBranding":
                        return NewOrganizationalBranding(), nil
                    case "#microsoft.graph.organizationalBrandingLocalization":
                        return NewOrganizationalBrandingLocalization(), nil
                    case "#microsoft.graph.organizationalBrandingProperties":
                        return NewOrganizationalBrandingProperties(), nil
                    case "#microsoft.graph.orgContact":
                        return NewOrgContact(), nil
                    case "#microsoft.graph.outlookCategory":
                        return NewOutlookCategory(), nil
                    case "#microsoft.graph.outlookItem":
                        return NewOutlookItem(), nil
                    case "#microsoft.graph.outlookUser":
                        return NewOutlookUser(), nil
                    case "#microsoft.graph.participant":
                        return NewParticipant(), nil
                    case "#microsoft.graph.participantJoiningNotification":
                        return NewParticipantJoiningNotification(), nil
                    case "#microsoft.graph.participantLeftNotification":
                        return NewParticipantLeftNotification(), nil
                    case "#microsoft.graph.partners":
                        return NewPartners(), nil
                    case "#microsoft.graph.passwordAuthenticationMethod":
                        return NewPasswordAuthenticationMethod(), nil
                    case "#microsoft.graph.payload":
                        return NewPayload(), nil
                    case "#microsoft.graph.peopleAdminSettings":
                        return NewPeopleAdminSettings(), nil
                    case "#microsoft.graph.permission":
                        return NewPermission(), nil
                    case "#microsoft.graph.permissionGrantConditionSet":
                        return NewPermissionGrantConditionSet(), nil
                    case "#microsoft.graph.permissionGrantPolicy":
                        return NewPermissionGrantPolicy(), nil
                    case "#microsoft.graph.person":
                        return NewPerson(), nil
                    case "#microsoft.graph.phoneAuthenticationMethod":
                        return NewPhoneAuthenticationMethod(), nil
                    case "#microsoft.graph.pinnedChatMessageInfo":
                        return NewPinnedChatMessageInfo(), nil
                    case "#microsoft.graph.place":
                        return NewPlace(), nil
                    case "#microsoft.graph.planner":
                        return NewPlanner(), nil
                    case "#microsoft.graph.plannerAssignedToTaskBoardTaskFormat":
                        return NewPlannerAssignedToTaskBoardTaskFormat(), nil
                    case "#microsoft.graph.plannerBucket":
                        return NewPlannerBucket(), nil
                    case "#microsoft.graph.plannerBucketTaskBoardTaskFormat":
                        return NewPlannerBucketTaskBoardTaskFormat(), nil
                    case "#microsoft.graph.plannerGroup":
                        return NewPlannerGroup(), nil
                    case "#microsoft.graph.plannerPlan":
                        return NewPlannerPlan(), nil
                    case "#microsoft.graph.plannerPlanDetails":
                        return NewPlannerPlanDetails(), nil
                    case "#microsoft.graph.plannerProgressTaskBoardTaskFormat":
                        return NewPlannerProgressTaskBoardTaskFormat(), nil
                    case "#microsoft.graph.plannerTask":
                        return NewPlannerTask(), nil
                    case "#microsoft.graph.plannerTaskDetails":
                        return NewPlannerTaskDetails(), nil
                    case "#microsoft.graph.plannerUser":
                        return NewPlannerUser(), nil
                    case "#microsoft.graph.playPromptOperation":
                        return NewPlayPromptOperation(), nil
                    case "#microsoft.graph.policyBase":
                        return NewPolicyBase(), nil
                    case "#microsoft.graph.policyRoot":
                        return NewPolicyRoot(), nil
                    case "#microsoft.graph.policyTemplate":
                        return NewPolicyTemplate(), nil
                    case "#microsoft.graph.post":
                        return NewPost(), nil
                    case "#microsoft.graph.presence":
                        return NewPresence(), nil
                    case "#microsoft.graph.printConnector":
                        return NewPrintConnector(), nil
                    case "#microsoft.graph.printDocument":
                        return NewPrintDocument(), nil
                    case "#microsoft.graph.printer":
                        return NewPrinter(), nil
                    case "#microsoft.graph.printerBase":
                        return NewPrinterBase(), nil
                    case "#microsoft.graph.printerCreateOperation":
                        return NewPrinterCreateOperation(), nil
                    case "#microsoft.graph.printerShare":
                        return NewPrinterShare(), nil
                    case "#microsoft.graph.printJob":
                        return NewPrintJob(), nil
                    case "#microsoft.graph.printOperation":
                        return NewPrintOperation(), nil
                    case "#microsoft.graph.printService":
                        return NewPrintService(), nil
                    case "#microsoft.graph.printServiceEndpoint":
                        return NewPrintServiceEndpoint(), nil
                    case "#microsoft.graph.printTask":
                        return NewPrintTask(), nil
                    case "#microsoft.graph.printTaskDefinition":
                        return NewPrintTaskDefinition(), nil
                    case "#microsoft.graph.printTaskTrigger":
                        return NewPrintTaskTrigger(), nil
                    case "#microsoft.graph.printUsage":
                        return NewPrintUsage(), nil
                    case "#microsoft.graph.printUsageByPrinter":
                        return NewPrintUsageByPrinter(), nil
                    case "#microsoft.graph.printUsageByUser":
                        return NewPrintUsageByUser(), nil
                    case "#microsoft.graph.privilegedAccessGroup":
                        return NewPrivilegedAccessGroup(), nil
                    case "#microsoft.graph.privilegedAccessGroupAssignmentSchedule":
                        return NewPrivilegedAccessGroupAssignmentSchedule(), nil
                    case "#microsoft.graph.privilegedAccessGroupAssignmentScheduleInstance":
                        return NewPrivilegedAccessGroupAssignmentScheduleInstance(), nil
                    case "#microsoft.graph.privilegedAccessGroupAssignmentScheduleRequest":
                        return NewPrivilegedAccessGroupAssignmentScheduleRequest(), nil
                    case "#microsoft.graph.privilegedAccessGroupEligibilitySchedule":
                        return NewPrivilegedAccessGroupEligibilitySchedule(), nil
                    case "#microsoft.graph.privilegedAccessGroupEligibilityScheduleInstance":
                        return NewPrivilegedAccessGroupEligibilityScheduleInstance(), nil
                    case "#microsoft.graph.privilegedAccessGroupEligibilityScheduleRequest":
                        return NewPrivilegedAccessGroupEligibilityScheduleRequest(), nil
                    case "#microsoft.graph.privilegedAccessRoot":
                        return NewPrivilegedAccessRoot(), nil
                    case "#microsoft.graph.privilegedAccessSchedule":
                        return NewPrivilegedAccessSchedule(), nil
                    case "#microsoft.graph.privilegedAccessScheduleInstance":
                        return NewPrivilegedAccessScheduleInstance(), nil
                    case "#microsoft.graph.privilegedAccessScheduleRequest":
                        return NewPrivilegedAccessScheduleRequest(), nil
                    case "#microsoft.graph.profileCardProperty":
                        return NewProfileCardProperty(), nil
                    case "#microsoft.graph.profilePhoto":
                        return NewProfilePhoto(), nil
                    case "#microsoft.graph.pronounsSettings":
                        return NewPronounsSettings(), nil
                    case "#microsoft.graph.protectionPolicyBase":
                        return NewProtectionPolicyBase(), nil
                    case "#microsoft.graph.protectionRuleBase":
                        return NewProtectionRuleBase(), nil
                    case "#microsoft.graph.protectionUnitBase":
                        return NewProtectionUnitBase(), nil
                    case "#microsoft.graph.provisioningObjectSummary":
                        return NewProvisioningObjectSummary(), nil
                    case "#microsoft.graph.rbacApplication":
                        return NewRbacApplication(), nil
                    case "#microsoft.graph.recordOperation":
                        return NewRecordOperation(), nil
                    case "#microsoft.graph.referenceAttachment":
                        return NewReferenceAttachment(), nil
                    case "#microsoft.graph.relyingPartyDetailedSummary":
                        return NewRelyingPartyDetailedSummary(), nil
                    case "#microsoft.graph.remoteAssistancePartner":
                        return NewRemoteAssistancePartner(), nil
                    case "#microsoft.graph.remoteDesktopSecurityConfiguration":
                        return NewRemoteDesktopSecurityConfiguration(), nil
                    case "#microsoft.graph.request":
                        return NewRequest(), nil
                    case "#microsoft.graph.resellerDelegatedAdminRelationship":
                        return NewResellerDelegatedAdminRelationship(), nil
                    case "#microsoft.graph.resourceOperation":
                        return NewResourceOperation(), nil
                    case "#microsoft.graph.resourceSpecificPermissionGrant":
                        return NewResourceSpecificPermissionGrant(), nil
                    case "#microsoft.graph.restoreArtifactBase":
                        return NewRestoreArtifactBase(), nil
                    case "#microsoft.graph.restorePoint":
                        return NewRestorePoint(), nil
                    case "#microsoft.graph.restoreSessionBase":
                        return NewRestoreSessionBase(), nil
                    case "#microsoft.graph.richLongRunningOperation":
                        return NewRichLongRunningOperation(), nil
                    case "#microsoft.graph.riskDetection":
                        return NewRiskDetection(), nil
                    case "#microsoft.graph.riskyServicePrincipal":
                        return NewRiskyServicePrincipal(), nil
                    case "#microsoft.graph.riskyServicePrincipalHistoryItem":
                        return NewRiskyServicePrincipalHistoryItem(), nil
                    case "#microsoft.graph.riskyUser":
                        return NewRiskyUser(), nil
                    case "#microsoft.graph.riskyUserHistoryItem":
                        return NewRiskyUserHistoryItem(), nil
                    case "#microsoft.graph.roleAssignment":
                        return NewRoleAssignment(), nil
                    case "#microsoft.graph.roleDefinition":
                        return NewRoleDefinition(), nil
                    case "#microsoft.graph.room":
                        return NewRoom(), nil
                    case "#microsoft.graph.roomList":
                        return NewRoomList(), nil
                    case "#microsoft.graph.samlOrWsFedExternalDomainFederation":
                        return NewSamlOrWsFedExternalDomainFederation(), nil
                    case "#microsoft.graph.samlOrWsFedProvider":
                        return NewSamlOrWsFedProvider(), nil
                    case "#microsoft.graph.schedule":
                        return NewSchedule(), nil
                    case "#microsoft.graph.scheduleChangeRequest":
                        return NewScheduleChangeRequest(), nil
                    case "#microsoft.graph.schedulingGroup":
                        return NewSchedulingGroup(), nil
                    case "#microsoft.graph.schemaExtension":
                        return NewSchemaExtension(), nil
                    case "#microsoft.graph.scopedRoleMembership":
                        return NewScopedRoleMembership(), nil
                    case "#microsoft.graph.searchEntity":
                        return NewSearchEntity(), nil
                    case "#microsoft.graph.sectionGroup":
                        return NewSectionGroup(), nil
                    case "#microsoft.graph.secureScore":
                        return NewSecureScore(), nil
                    case "#microsoft.graph.secureScoreControlProfile":
                        return NewSecureScoreControlProfile(), nil
                    case "#microsoft.graph.security":
                        return NewSecurity(), nil
                    case "#microsoft.graph.securityReportsRoot":
                        return NewSecurityReportsRoot(), nil
                    case "#microsoft.graph.sendDtmfTonesOperation":
                        return NewSendDtmfTonesOperation(), nil
                    case "#microsoft.graph.serviceAnnouncement":
                        return NewServiceAnnouncement(), nil
                    case "#microsoft.graph.serviceAnnouncementAttachment":
                        return NewServiceAnnouncementAttachment(), nil
                    case "#microsoft.graph.serviceAnnouncementBase":
                        return NewServiceAnnouncementBase(), nil
                    case "#microsoft.graph.serviceApp":
                        return NewServiceApp(), nil
                    case "#microsoft.graph.serviceHealth":
                        return NewServiceHealth(), nil
                    case "#microsoft.graph.serviceHealthIssue":
                        return NewServiceHealthIssue(), nil
                    case "#microsoft.graph.servicePrincipal":
                        return NewServicePrincipal(), nil
                    case "#microsoft.graph.servicePrincipalRiskDetection":
                        return NewServicePrincipalRiskDetection(), nil
                    case "#microsoft.graph.serviceStorageQuotaBreakdown":
                        return NewServiceStorageQuotaBreakdown(), nil
                    case "#microsoft.graph.serviceUpdateMessage":
                        return NewServiceUpdateMessage(), nil
                    case "#microsoft.graph.settingStateDeviceSummary":
                        return NewSettingStateDeviceSummary(), nil
                    case "#microsoft.graph.sharedDriveItem":
                        return NewSharedDriveItem(), nil
                    case "#microsoft.graph.sharedInsight":
                        return NewSharedInsight(), nil
                    case "#microsoft.graph.sharedPCConfiguration":
                        return NewSharedPCConfiguration(), nil
                    case "#microsoft.graph.sharedWithChannelTeamInfo":
                        return NewSharedWithChannelTeamInfo(), nil
                    case "#microsoft.graph.sharepoint":
                        return NewSharepoint(), nil
                    case "#microsoft.graph.sharePointProtectionPolicy":
                        return NewSharePointProtectionPolicy(), nil
                    case "#microsoft.graph.sharePointRestoreSession":
                        return NewSharePointRestoreSession(), nil
                    case "#microsoft.graph.sharepointSettings":
                        return NewSharepointSettings(), nil
                    case "#microsoft.graph.shift":
                        return NewShift(), nil
                    case "#microsoft.graph.shiftPreferences":
                        return NewShiftPreferences(), nil
                    case "#microsoft.graph.signIn":
                        return NewSignIn(), nil
                    case "#microsoft.graph.simulation":
                        return NewSimulation(), nil
                    case "#microsoft.graph.simulationAutomation":
                        return NewSimulationAutomation(), nil
                    case "#microsoft.graph.simulationAutomationRun":
                        return NewSimulationAutomationRun(), nil
                    case "#microsoft.graph.singleValueLegacyExtendedProperty":
                        return NewSingleValueLegacyExtendedProperty(), nil
                    case "#microsoft.graph.site":
                        return NewSite(), nil
                    case "#microsoft.graph.sitePage":
                        return NewSitePage(), nil
                    case "#microsoft.graph.siteProtectionRule":
                        return NewSiteProtectionRule(), nil
                    case "#microsoft.graph.siteProtectionUnit":
                        return NewSiteProtectionUnit(), nil
                    case "#microsoft.graph.siteRestoreArtifact":
                        return NewSiteRestoreArtifact(), nil
                    case "#microsoft.graph.skypeForBusinessUserConversationMember":
                        return NewSkypeForBusinessUserConversationMember(), nil
                    case "#microsoft.graph.skypeUserConversationMember":
                        return NewSkypeUserConversationMember(), nil
                    case "#microsoft.graph.smsAuthenticationMethodConfiguration":
                        return NewSmsAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.smsAuthenticationMethodTarget":
                        return NewSmsAuthenticationMethodTarget(), nil
                    case "#microsoft.graph.socialIdentityProvider":
                        return NewSocialIdentityProvider(), nil
                    case "#microsoft.graph.softwareOathAuthenticationMethod":
                        return NewSoftwareOathAuthenticationMethod(), nil
                    case "#microsoft.graph.softwareOathAuthenticationMethodConfiguration":
                        return NewSoftwareOathAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.softwareUpdateStatusSummary":
                        return NewSoftwareUpdateStatusSummary(), nil
                    case "#microsoft.graph.standardWebPart":
                        return NewStandardWebPart(), nil
                    case "#microsoft.graph.startHoldMusicOperation":
                        return NewStartHoldMusicOperation(), nil
                    case "#microsoft.graph.stopHoldMusicOperation":
                        return NewStopHoldMusicOperation(), nil
                    case "#microsoft.graph.storageQuotaBreakdown":
                        return NewStorageQuotaBreakdown(), nil
                    case "#microsoft.graph.storageSettings":
                        return NewStorageSettings(), nil
                    case "#microsoft.graph.stsPolicy":
                        return NewStsPolicy(), nil
                    case "#microsoft.graph.subjectRightsRequest":
                        return NewSubjectRightsRequest(), nil
                    case "#microsoft.graph.subscribedSku":
                        return NewSubscribedSku(), nil
                    case "#microsoft.graph.subscribeToToneOperation":
                        return NewSubscribeToToneOperation(), nil
                    case "#microsoft.graph.subscription":
                        return NewSubscription(), nil
                    case "#microsoft.graph.swapShiftsChangeRequest":
                        return NewSwapShiftsChangeRequest(), nil
                    case "#microsoft.graph.synchronization":
                        return NewSynchronization(), nil
                    case "#microsoft.graph.synchronizationJob":
                        return NewSynchronizationJob(), nil
                    case "#microsoft.graph.synchronizationSchema":
                        return NewSynchronizationSchema(), nil
                    case "#microsoft.graph.synchronizationTemplate":
                        return NewSynchronizationTemplate(), nil
                    case "#microsoft.graph.targetDeviceGroup":
                        return NewTargetDeviceGroup(), nil
                    case "#microsoft.graph.targetedManagedAppConfiguration":
                        return NewTargetedManagedAppConfiguration(), nil
                    case "#microsoft.graph.targetedManagedAppPolicyAssignment":
                        return NewTargetedManagedAppPolicyAssignment(), nil
                    case "#microsoft.graph.targetedManagedAppProtection":
                        return NewTargetedManagedAppProtection(), nil
                    case "#microsoft.graph.taskFileAttachment":
                        return NewTaskFileAttachment(), nil
                    case "#microsoft.graph.team":
                        return NewTeam(), nil
                    case "#microsoft.graph.teamInfo":
                        return NewTeamInfo(), nil
                    case "#microsoft.graph.teamsApp":
                        return NewTeamsApp(), nil
                    case "#microsoft.graph.teamsAppDefinition":
                        return NewTeamsAppDefinition(), nil
                    case "#microsoft.graph.teamsAppInstallation":
                        return NewTeamsAppInstallation(), nil
                    case "#microsoft.graph.teamsAppSettings":
                        return NewTeamsAppSettings(), nil
                    case "#microsoft.graph.teamsAsyncOperation":
                        return NewTeamsAsyncOperation(), nil
                    case "#microsoft.graph.teamsTab":
                        return NewTeamsTab(), nil
                    case "#microsoft.graph.teamsTemplate":
                        return NewTeamsTemplate(), nil
                    case "#microsoft.graph.teamwork":
                        return NewTeamwork(), nil
                    case "#microsoft.graph.teamworkBot":
                        return NewTeamworkBot(), nil
                    case "#microsoft.graph.teamworkHostedContent":
                        return NewTeamworkHostedContent(), nil
                    case "#microsoft.graph.teamworkTag":
                        return NewTeamworkTag(), nil
                    case "#microsoft.graph.teamworkTagMember":
                        return NewTeamworkTagMember(), nil
                    case "#microsoft.graph.telecomExpenseManagementPartner":
                        return NewTelecomExpenseManagementPartner(), nil
                    case "#microsoft.graph.temporaryAccessPassAuthenticationMethod":
                        return NewTemporaryAccessPassAuthenticationMethod(), nil
                    case "#microsoft.graph.temporaryAccessPassAuthenticationMethodConfiguration":
                        return NewTemporaryAccessPassAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.tenantAppManagementPolicy":
                        return NewTenantAppManagementPolicy(), nil
                    case "#microsoft.graph.termsAndConditions":
                        return NewTermsAndConditions(), nil
                    case "#microsoft.graph.termsAndConditionsAcceptanceStatus":
                        return NewTermsAndConditionsAcceptanceStatus(), nil
                    case "#microsoft.graph.termsAndConditionsAssignment":
                        return NewTermsAndConditionsAssignment(), nil
                    case "#microsoft.graph.termsOfUseContainer":
                        return NewTermsOfUseContainer(), nil
                    case "#microsoft.graph.textWebPart":
                        return NewTextWebPart(), nil
                    case "#microsoft.graph.threatAssessmentRequest":
                        return NewThreatAssessmentRequest(), nil
                    case "#microsoft.graph.threatAssessmentResult":
                        return NewThreatAssessmentResult(), nil
                    case "#microsoft.graph.thumbnailSet":
                        return NewThumbnailSet(), nil
                    case "#microsoft.graph.timeOff":
                        return NewTimeOff(), nil
                    case "#microsoft.graph.timeOffReason":
                        return NewTimeOffReason(), nil
                    case "#microsoft.graph.timeOffRequest":
                        return NewTimeOffRequest(), nil
                    case "#microsoft.graph.todo":
                        return NewTodo(), nil
                    case "#microsoft.graph.todoTask":
                        return NewTodoTask(), nil
                    case "#microsoft.graph.todoTaskList":
                        return NewTodoTaskList(), nil
                    case "#microsoft.graph.tokenIssuancePolicy":
                        return NewTokenIssuancePolicy(), nil
                    case "#microsoft.graph.tokenLifetimePolicy":
                        return NewTokenLifetimePolicy(), nil
                    case "#microsoft.graph.training":
                        return NewTraining(), nil
                    case "#microsoft.graph.trainingLanguageDetail":
                        return NewTrainingLanguageDetail(), nil
                    case "#microsoft.graph.trending":
                        return NewTrending(), nil
                    case "#microsoft.graph.unifiedRbacResourceAction":
                        return NewUnifiedRbacResourceAction(), nil
                    case "#microsoft.graph.unifiedRbacResourceNamespace":
                        return NewUnifiedRbacResourceNamespace(), nil
                    case "#microsoft.graph.unifiedRoleAssignment":
                        return NewUnifiedRoleAssignment(), nil
                    case "#microsoft.graph.unifiedRoleAssignmentSchedule":
                        return NewUnifiedRoleAssignmentSchedule(), nil
                    case "#microsoft.graph.unifiedRoleAssignmentScheduleInstance":
                        return NewUnifiedRoleAssignmentScheduleInstance(), nil
                    case "#microsoft.graph.unifiedRoleAssignmentScheduleRequest":
                        return NewUnifiedRoleAssignmentScheduleRequest(), nil
                    case "#microsoft.graph.unifiedRoleDefinition":
                        return NewUnifiedRoleDefinition(), nil
                    case "#microsoft.graph.unifiedRoleEligibilitySchedule":
                        return NewUnifiedRoleEligibilitySchedule(), nil
                    case "#microsoft.graph.unifiedRoleEligibilityScheduleInstance":
                        return NewUnifiedRoleEligibilityScheduleInstance(), nil
                    case "#microsoft.graph.unifiedRoleEligibilityScheduleRequest":
                        return NewUnifiedRoleEligibilityScheduleRequest(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicy":
                        return NewUnifiedRoleManagementPolicy(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyApprovalRule":
                        return NewUnifiedRoleManagementPolicyApprovalRule(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyAssignment":
                        return NewUnifiedRoleManagementPolicyAssignment(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule":
                        return NewUnifiedRoleManagementPolicyAuthenticationContextRule(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyEnablementRule":
                        return NewUnifiedRoleManagementPolicyEnablementRule(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule":
                        return NewUnifiedRoleManagementPolicyExpirationRule(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule":
                        return NewUnifiedRoleManagementPolicyNotificationRule(), nil
                    case "#microsoft.graph.unifiedRoleManagementPolicyRule":
                        return NewUnifiedRoleManagementPolicyRule(), nil
                    case "#microsoft.graph.unifiedRoleScheduleBase":
                        return NewUnifiedRoleScheduleBase(), nil
                    case "#microsoft.graph.unifiedRoleScheduleInstanceBase":
                        return NewUnifiedRoleScheduleInstanceBase(), nil
                    case "#microsoft.graph.unifiedStorageQuota":
                        return NewUnifiedStorageQuota(), nil
                    case "#microsoft.graph.unmuteParticipantOperation":
                        return NewUnmuteParticipantOperation(), nil
                    case "#microsoft.graph.updateRecordingStatusOperation":
                        return NewUpdateRecordingStatusOperation(), nil
                    case "#microsoft.graph.urlAssessmentRequest":
                        return NewUrlAssessmentRequest(), nil
                    case "#microsoft.graph.usedInsight":
                        return NewUsedInsight(), nil
                    case "#microsoft.graph.user":
                        return NewUser(), nil
                    case "#microsoft.graph.userActivity":
                        return NewUserActivity(), nil
                    case "#microsoft.graph.userConsentRequest":
                        return NewUserConsentRequest(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthApplicationPerformance":
                        return NewUserExperienceAnalyticsAppHealthApplicationPerformance(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails":
                        return NewUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthAppPerformanceByAppVersionDeviceId":
                        return NewUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDeviceId(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthAppPerformanceByOSVersion":
                        return NewUserExperienceAnalyticsAppHealthAppPerformanceByOSVersion(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthDeviceModelPerformance":
                        return NewUserExperienceAnalyticsAppHealthDeviceModelPerformance(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthDevicePerformance":
                        return NewUserExperienceAnalyticsAppHealthDevicePerformance(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthDevicePerformanceDetails":
                        return NewUserExperienceAnalyticsAppHealthDevicePerformanceDetails(), nil
                    case "#microsoft.graph.userExperienceAnalyticsAppHealthOSVersionPerformance":
                        return NewUserExperienceAnalyticsAppHealthOSVersionPerformance(), nil
                    case "#microsoft.graph.userExperienceAnalyticsBaseline":
                        return NewUserExperienceAnalyticsBaseline(), nil
                    case "#microsoft.graph.userExperienceAnalyticsCategory":
                        return NewUserExperienceAnalyticsCategory(), nil
                    case "#microsoft.graph.userExperienceAnalyticsDevicePerformance":
                        return NewUserExperienceAnalyticsDevicePerformance(), nil
                    case "#microsoft.graph.userExperienceAnalyticsDeviceScores":
                        return NewUserExperienceAnalyticsDeviceScores(), nil
                    case "#microsoft.graph.userExperienceAnalyticsDeviceStartupHistory":
                        return NewUserExperienceAnalyticsDeviceStartupHistory(), nil
                    case "#microsoft.graph.userExperienceAnalyticsDeviceStartupProcess":
                        return NewUserExperienceAnalyticsDeviceStartupProcess(), nil
                    case "#microsoft.graph.userExperienceAnalyticsDeviceStartupProcessPerformance":
                        return NewUserExperienceAnalyticsDeviceStartupProcessPerformance(), nil
                    case "#microsoft.graph.userExperienceAnalyticsMetric":
                        return NewUserExperienceAnalyticsMetric(), nil
                    case "#microsoft.graph.userExperienceAnalyticsMetricHistory":
                        return NewUserExperienceAnalyticsMetricHistory(), nil
                    case "#microsoft.graph.userExperienceAnalyticsModelScores":
                        return NewUserExperienceAnalyticsModelScores(), nil
                    case "#microsoft.graph.userExperienceAnalyticsOverview":
                        return NewUserExperienceAnalyticsOverview(), nil
                    case "#microsoft.graph.userExperienceAnalyticsScoreHistory":
                        return NewUserExperienceAnalyticsScoreHistory(), nil
                    case "#microsoft.graph.userExperienceAnalyticsWorkFromAnywhereDevice":
                        return NewUserExperienceAnalyticsWorkFromAnywhereDevice(), nil
                    case "#microsoft.graph.userExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric":
                        return NewUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric(), nil
                    case "#microsoft.graph.userExperienceAnalyticsWorkFromAnywhereMetric":
                        return NewUserExperienceAnalyticsWorkFromAnywhereMetric(), nil
                    case "#microsoft.graph.userExperienceAnalyticsWorkFromAnywhereModelPerformance":
                        return NewUserExperienceAnalyticsWorkFromAnywhereModelPerformance(), nil
                    case "#microsoft.graph.userFlowLanguageConfiguration":
                        return NewUserFlowLanguageConfiguration(), nil
                    case "#microsoft.graph.userFlowLanguagePage":
                        return NewUserFlowLanguagePage(), nil
                    case "#microsoft.graph.userInsightsSettings":
                        return NewUserInsightsSettings(), nil
                    case "#microsoft.graph.userInstallStateSummary":
                        return NewUserInstallStateSummary(), nil
                    case "#microsoft.graph.userRegistrationDetails":
                        return NewUserRegistrationDetails(), nil
                    case "#microsoft.graph.userScopeTeamsAppInstallation":
                        return NewUserScopeTeamsAppInstallation(), nil
                    case "#microsoft.graph.userSettings":
                        return NewUserSettings(), nil
                    case "#microsoft.graph.userSignInInsight":
                        return NewUserSignInInsight(), nil
                    case "#microsoft.graph.userSolutionRoot":
                        return NewUserSolutionRoot(), nil
                    case "#microsoft.graph.userStorage":
                        return NewUserStorage(), nil
                    case "#microsoft.graph.userTeamwork":
                        return NewUserTeamwork(), nil
                    case "#microsoft.graph.verticalSection":
                        return NewVerticalSection(), nil
                    case "#microsoft.graph.virtualEndpoint":
                        return NewVirtualEndpoint(), nil
                    case "#microsoft.graph.virtualEvent":
                        return NewVirtualEvent(), nil
                    case "#microsoft.graph.virtualEventPresenter":
                        return NewVirtualEventPresenter(), nil
                    case "#microsoft.graph.virtualEventRegistration":
                        return NewVirtualEventRegistration(), nil
                    case "#microsoft.graph.virtualEventRegistrationConfiguration":
                        return NewVirtualEventRegistrationConfiguration(), nil
                    case "#microsoft.graph.virtualEventRegistrationCustomQuestion":
                        return NewVirtualEventRegistrationCustomQuestion(), nil
                    case "#microsoft.graph.virtualEventRegistrationPredefinedQuestion":
                        return NewVirtualEventRegistrationPredefinedQuestion(), nil
                    case "#microsoft.graph.virtualEventRegistrationQuestionBase":
                        return NewVirtualEventRegistrationQuestionBase(), nil
                    case "#microsoft.graph.virtualEventSession":
                        return NewVirtualEventSession(), nil
                    case "#microsoft.graph.virtualEventsRoot":
                        return NewVirtualEventsRoot(), nil
                    case "#microsoft.graph.virtualEventTownhall":
                        return NewVirtualEventTownhall(), nil
                    case "#microsoft.graph.virtualEventWebinar":
                        return NewVirtualEventWebinar(), nil
                    case "#microsoft.graph.virtualEventWebinarRegistrationConfiguration":
                        return NewVirtualEventWebinarRegistrationConfiguration(), nil
                    case "#microsoft.graph.voiceAuthenticationMethodConfiguration":
                        return NewVoiceAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.vppToken":
                        return NewVppToken(), nil
                    case "#microsoft.graph.webApp":
                        return NewWebApp(), nil
                    case "#microsoft.graph.webPart":
                        return NewWebPart(), nil
                    case "#microsoft.graph.win32LobApp":
                        return NewWin32LobApp(), nil
                    case "#microsoft.graph.windows10CompliancePolicy":
                        return NewWindows10CompliancePolicy(), nil
                    case "#microsoft.graph.windows10CustomConfiguration":
                        return NewWindows10CustomConfiguration(), nil
                    case "#microsoft.graph.windows10EndpointProtectionConfiguration":
                        return NewWindows10EndpointProtectionConfiguration(), nil
                    case "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration":
                        return NewWindows10EnrollmentCompletionPageConfiguration(), nil
                    case "#microsoft.graph.windows10EnterpriseModernAppManagementConfiguration":
                        return NewWindows10EnterpriseModernAppManagementConfiguration(), nil
                    case "#microsoft.graph.windows10GeneralConfiguration":
                        return NewWindows10GeneralConfiguration(), nil
                    case "#microsoft.graph.windows10MobileCompliancePolicy":
                        return NewWindows10MobileCompliancePolicy(), nil
                    case "#microsoft.graph.windows10SecureAssessmentConfiguration":
                        return NewWindows10SecureAssessmentConfiguration(), nil
                    case "#microsoft.graph.windows10TeamGeneralConfiguration":
                        return NewWindows10TeamGeneralConfiguration(), nil
                    case "#microsoft.graph.windows81CompliancePolicy":
                        return NewWindows81CompliancePolicy(), nil
                    case "#microsoft.graph.windows81GeneralConfiguration":
                        return NewWindows81GeneralConfiguration(), nil
                    case "#microsoft.graph.windowsAppX":
                        return NewWindowsAppX(), nil
                    case "#microsoft.graph.windowsAutopilotDeploymentProfile":
                        return NewWindowsAutopilotDeploymentProfile(), nil
                    case "#microsoft.graph.windowsAutopilotDeploymentProfileAssignment":
                        return NewWindowsAutopilotDeploymentProfileAssignment(), nil
                    case "#microsoft.graph.windowsAutopilotDeviceIdentity":
                        return NewWindowsAutopilotDeviceIdentity(), nil
                    case "#microsoft.graph.windowsDefenderAdvancedThreatProtectionConfiguration":
                        return NewWindowsDefenderAdvancedThreatProtectionConfiguration(), nil
                    case "#microsoft.graph.windowsDeviceMalwareState":
                        return NewWindowsDeviceMalwareState(), nil
                    case "#microsoft.graph.windowsHelloForBusinessAuthenticationMethod":
                        return NewWindowsHelloForBusinessAuthenticationMethod(), nil
                    case "#microsoft.graph.windowsInformationProtection":
                        return NewWindowsInformationProtection(), nil
                    case "#microsoft.graph.windowsInformationProtectionAppLearningSummary":
                        return NewWindowsInformationProtectionAppLearningSummary(), nil
                    case "#microsoft.graph.windowsInformationProtectionAppLockerFile":
                        return NewWindowsInformationProtectionAppLockerFile(), nil
                    case "#microsoft.graph.windowsInformationProtectionNetworkLearningSummary":
                        return NewWindowsInformationProtectionNetworkLearningSummary(), nil
                    case "#microsoft.graph.windowsInformationProtectionPolicy":
                        return NewWindowsInformationProtectionPolicy(), nil
                    case "#microsoft.graph.windowsMalwareInformation":
                        return NewWindowsMalwareInformation(), nil
                    case "#microsoft.graph.windowsMicrosoftEdgeApp":
                        return NewWindowsMicrosoftEdgeApp(), nil
                    case "#microsoft.graph.windowsMobileMSI":
                        return NewWindowsMobileMSI(), nil
                    case "#microsoft.graph.windowsPhone81CompliancePolicy":
                        return NewWindowsPhone81CompliancePolicy(), nil
                    case "#microsoft.graph.windowsPhone81CustomConfiguration":
                        return NewWindowsPhone81CustomConfiguration(), nil
                    case "#microsoft.graph.windowsPhone81GeneralConfiguration":
                        return NewWindowsPhone81GeneralConfiguration(), nil
                    case "#microsoft.graph.windowsProtectionState":
                        return NewWindowsProtectionState(), nil
                    case "#microsoft.graph.windowsSetting":
                        return NewWindowsSetting(), nil
                    case "#microsoft.graph.windowsSettingInstance":
                        return NewWindowsSettingInstance(), nil
                    case "#microsoft.graph.windowsUniversalAppX":
                        return NewWindowsUniversalAppX(), nil
                    case "#microsoft.graph.windowsUniversalAppXContainedApp":
                        return NewWindowsUniversalAppXContainedApp(), nil
                    case "#microsoft.graph.windowsUpdateForBusinessConfiguration":
                        return NewWindowsUpdateForBusinessConfiguration(), nil
                    case "#microsoft.graph.windowsWebApp":
                        return NewWindowsWebApp(), nil
                    case "#microsoft.graph.workbook":
                        return NewWorkbook(), nil
                    case "#microsoft.graph.workbookApplication":
                        return NewWorkbookApplication(), nil
                    case "#microsoft.graph.workbookChart":
                        return NewWorkbookChart(), nil
                    case "#microsoft.graph.workbookChartAreaFormat":
                        return NewWorkbookChartAreaFormat(), nil
                    case "#microsoft.graph.workbookChartAxes":
                        return NewWorkbookChartAxes(), nil
                    case "#microsoft.graph.workbookChartAxis":
                        return NewWorkbookChartAxis(), nil
                    case "#microsoft.graph.workbookChartAxisFormat":
                        return NewWorkbookChartAxisFormat(), nil
                    case "#microsoft.graph.workbookChartAxisTitle":
                        return NewWorkbookChartAxisTitle(), nil
                    case "#microsoft.graph.workbookChartAxisTitleFormat":
                        return NewWorkbookChartAxisTitleFormat(), nil
                    case "#microsoft.graph.workbookChartDataLabelFormat":
                        return NewWorkbookChartDataLabelFormat(), nil
                    case "#microsoft.graph.workbookChartDataLabels":
                        return NewWorkbookChartDataLabels(), nil
                    case "#microsoft.graph.workbookChartFill":
                        return NewWorkbookChartFill(), nil
                    case "#microsoft.graph.workbookChartFont":
                        return NewWorkbookChartFont(), nil
                    case "#microsoft.graph.workbookChartGridlines":
                        return NewWorkbookChartGridlines(), nil
                    case "#microsoft.graph.workbookChartGridlinesFormat":
                        return NewWorkbookChartGridlinesFormat(), nil
                    case "#microsoft.graph.workbookChartLegend":
                        return NewWorkbookChartLegend(), nil
                    case "#microsoft.graph.workbookChartLegendFormat":
                        return NewWorkbookChartLegendFormat(), nil
                    case "#microsoft.graph.workbookChartLineFormat":
                        return NewWorkbookChartLineFormat(), nil
                    case "#microsoft.graph.workbookChartPoint":
                        return NewWorkbookChartPoint(), nil
                    case "#microsoft.graph.workbookChartPointFormat":
                        return NewWorkbookChartPointFormat(), nil
                    case "#microsoft.graph.workbookChartSeries":
                        return NewWorkbookChartSeries(), nil
                    case "#microsoft.graph.workbookChartSeriesFormat":
                        return NewWorkbookChartSeriesFormat(), nil
                    case "#microsoft.graph.workbookChartTitle":
                        return NewWorkbookChartTitle(), nil
                    case "#microsoft.graph.workbookChartTitleFormat":
                        return NewWorkbookChartTitleFormat(), nil
                    case "#microsoft.graph.workbookComment":
                        return NewWorkbookComment(), nil
                    case "#microsoft.graph.workbookCommentReply":
                        return NewWorkbookCommentReply(), nil
                    case "#microsoft.graph.workbookFilter":
                        return NewWorkbookFilter(), nil
                    case "#microsoft.graph.workbookFormatProtection":
                        return NewWorkbookFormatProtection(), nil
                    case "#microsoft.graph.workbookFunctionResult":
                        return NewWorkbookFunctionResult(), nil
                    case "#microsoft.graph.workbookFunctions":
                        return NewWorkbookFunctions(), nil
                    case "#microsoft.graph.workbookNamedItem":
                        return NewWorkbookNamedItem(), nil
                    case "#microsoft.graph.workbookOperation":
                        return NewWorkbookOperation(), nil
                    case "#microsoft.graph.workbookPivotTable":
                        return NewWorkbookPivotTable(), nil
                    case "#microsoft.graph.workbookRange":
                        return NewWorkbookRange(), nil
                    case "#microsoft.graph.workbookRangeBorder":
                        return NewWorkbookRangeBorder(), nil
                    case "#microsoft.graph.workbookRangeFill":
                        return NewWorkbookRangeFill(), nil
                    case "#microsoft.graph.workbookRangeFont":
                        return NewWorkbookRangeFont(), nil
                    case "#microsoft.graph.workbookRangeFormat":
                        return NewWorkbookRangeFormat(), nil
                    case "#microsoft.graph.workbookRangeSort":
                        return NewWorkbookRangeSort(), nil
                    case "#microsoft.graph.workbookRangeView":
                        return NewWorkbookRangeView(), nil
                    case "#microsoft.graph.workbookTable":
                        return NewWorkbookTable(), nil
                    case "#microsoft.graph.workbookTableColumn":
                        return NewWorkbookTableColumn(), nil
                    case "#microsoft.graph.workbookTableRow":
                        return NewWorkbookTableRow(), nil
                    case "#microsoft.graph.workbookTableSort":
                        return NewWorkbookTableSort(), nil
                    case "#microsoft.graph.workbookWorksheet":
                        return NewWorkbookWorksheet(), nil
                    case "#microsoft.graph.workbookWorksheetProtection":
                        return NewWorkbookWorksheetProtection(), nil
                    case "#microsoft.graph.workforceIntegration":
                        return NewWorkforceIntegration(), nil
                    case "#microsoft.graph.workingTimeSchedule":
                        return NewWorkingTimeSchedule(), nil
                    case "#microsoft.graph.x509CertificateAuthenticationMethodConfiguration":
                        return NewX509CertificateAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.x509CertificateCombinationConfiguration":
                        return NewX509CertificateCombinationConfiguration(), nil
                }
            }
        }
    }
    return NewEntity(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Entity) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *Entity) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Entity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["id"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetId(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    return res
}
// GetId gets the id property value. The unique identifier for an entity. Read-only.
// returns a *string when successful
func (m *Entity) GetId()(*string) {
    val, err := m.GetBackingStore().Get("id")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Entity) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Entity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("id", m.GetId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Entity) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Entity) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetId sets the id property value. The unique identifier for an entity. Read-only.
func (m *Entity) SetId(value *string)() {
    err := m.GetBackingStore().Set("id", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Entity) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type Entityable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetId()(*string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetId(value *string)()
    SetOdataType(value *string)()
}
