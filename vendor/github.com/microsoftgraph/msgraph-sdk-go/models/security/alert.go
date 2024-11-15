package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Alert struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewAlert instantiates a new Alert and sets the default values.
func NewAlert()(*Alert) {
    m := &Alert{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateAlertFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAlertFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAlert(), nil
}
// GetActorDisplayName gets the actorDisplayName property value. The adversary or activity group that is associated with this alert.
// returns a *string when successful
func (m *Alert) GetActorDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("actorDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAdditionalDataProperty gets the additionalData property value. A collection of other alert properties, including user-defined properties. Any custom details defined in the alert, and any dynamic content in the alert details, are stored here.
// returns a Dictionaryable when successful
func (m *Alert) GetAdditionalDataProperty()(Dictionaryable) {
    val, err := m.GetBackingStore().Get("additionalDataProperty")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Dictionaryable)
    }
    return nil
}
// GetAlertPolicyId gets the alertPolicyId property value. The ID of the policy that generated the alert, and populated when there is a specific policy that generated the alert, whether configured by a customer or a built-in policy.
// returns a *string when successful
func (m *Alert) GetAlertPolicyId()(*string) {
    val, err := m.GetBackingStore().Get("alertPolicyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAlertWebUrl gets the alertWebUrl property value. URL for the Microsoft 365 Defender portal alert page.
// returns a *string when successful
func (m *Alert) GetAlertWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("alertWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAssignedTo gets the assignedTo property value. Owner of the alert, or null if no owner is assigned.
// returns a *string when successful
func (m *Alert) GetAssignedTo()(*string) {
    val, err := m.GetBackingStore().Get("assignedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory gets the category property value. The attack kill-chain category that the alert belongs to. Aligned with the MITRE ATT&CK framework.
// returns a *string when successful
func (m *Alert) GetCategory()(*string) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetClassification gets the classification property value. Specifies whether the alert represents a true threat. Possible values are: unknown, falsePositive, truePositive, informationalExpectedActivity, unknownFutureValue.
// returns a *AlertClassification when successful
func (m *Alert) GetClassification()(*AlertClassification) {
    val, err := m.GetBackingStore().Get("classification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AlertClassification)
    }
    return nil
}
// GetComments gets the comments property value. Array of comments created by the Security Operations (SecOps) team during the alert management process.
// returns a []AlertCommentable when successful
func (m *Alert) GetComments()([]AlertCommentable) {
    val, err := m.GetBackingStore().Get("comments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AlertCommentable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Time when Microsoft 365 Defender created the alert.
// returns a *Time when successful
func (m *Alert) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. String value describing each alert.
// returns a *string when successful
func (m *Alert) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDetectionSource gets the detectionSource property value. Detection technology or sensor that identified the notable component or activity. Possible values are: unknown, microsoftDefenderForEndpoint, antivirus, smartScreen, customTi, microsoftDefenderForOffice365, automatedInvestigation, microsoftThreatExperts, customDetection, microsoftDefenderForIdentity, cloudAppSecurity, microsoft365Defender, azureAdIdentityProtection, manual, microsoftDataLossPrevention, appGovernancePolicy, appGovernanceDetection, unknownFutureValue, microsoftDefenderForCloud, microsoftDefenderForIoT, microsoftDefenderForServers, microsoftDefenderForStorage, microsoftDefenderForDNS, microsoftDefenderForDatabases, microsoftDefenderForContainers, microsoftDefenderForNetwork, microsoftDefenderForAppService, microsoftDefenderForKeyVault, microsoftDefenderForResourceManager, microsoftDefenderForApiManagement, microsoftSentinel, nrtAlerts, scheduledAlerts, microsoftDefenderThreatIntelligenceAnalytics, builtInMl. You must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: microsoftDefenderForCloud, microsoftDefenderForIoT, microsoftDefenderForServers, microsoftDefenderForStorage, microsoftDefenderForDNS, microsoftDefenderForDatabases, microsoftDefenderForContainers, microsoftDefenderForNetwork, microsoftDefenderForAppService, microsoftDefenderForKeyVault, microsoftDefenderForResourceManager, microsoftDefenderForApiManagement, microsoftSentinel, nrtAlerts, scheduledAlerts, microsoftDefenderThreatIntelligenceAnalytics, builtInMl.
// returns a *DetectionSource when successful
func (m *Alert) GetDetectionSource()(*DetectionSource) {
    val, err := m.GetBackingStore().Get("detectionSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DetectionSource)
    }
    return nil
}
// GetDetectorId gets the detectorId property value. The ID of the detector that triggered the alert.
// returns a *string when successful
func (m *Alert) GetDetectorId()(*string) {
    val, err := m.GetBackingStore().Get("detectorId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDetermination gets the determination property value. Specifies the result of the investigation, whether the alert represents a true attack and if so, the nature of the attack. Possible values are: unknown, apt, malware, securityPersonnel, securityTesting, unwantedSoftware, other, multiStagedAttack, compromisedAccount, phishing, maliciousUserActivity, notMalicious, notEnoughDataToValidate, confirmedUserActivity, lineOfBusinessApplication, unknownFutureValue.
// returns a *AlertDetermination when successful
func (m *Alert) GetDetermination()(*AlertDetermination) {
    val, err := m.GetBackingStore().Get("determination")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AlertDetermination)
    }
    return nil
}
// GetEvidence gets the evidence property value. Collection of evidence related to the alert.
// returns a []AlertEvidenceable when successful
func (m *Alert) GetEvidence()([]AlertEvidenceable) {
    val, err := m.GetBackingStore().Get("evidence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AlertEvidenceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Alert) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["actorDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActorDisplayName(val)
        }
        return nil
    }
    res["additionalData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDictionaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAdditionalDataProperty(val.(Dictionaryable))
        }
        return nil
    }
    res["alertPolicyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlertPolicyId(val)
        }
        return nil
    }
    res["alertWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlertWebUrl(val)
        }
        return nil
    }
    res["assignedTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignedTo(val)
        }
        return nil
    }
    res["category"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory(val)
        }
        return nil
    }
    res["classification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAlertClassification)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassification(val.(*AlertClassification))
        }
        return nil
    }
    res["comments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAlertCommentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AlertCommentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AlertCommentable)
                }
            }
            m.SetComments(res)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["detectionSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDetectionSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetectionSource(val.(*DetectionSource))
        }
        return nil
    }
    res["detectorId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetectorId(val)
        }
        return nil
    }
    res["determination"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAlertDetermination)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetermination(val.(*AlertDetermination))
        }
        return nil
    }
    res["evidence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAlertEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AlertEvidenceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AlertEvidenceable)
                }
            }
            m.SetEvidence(res)
        }
        return nil
    }
    res["firstActivityDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirstActivityDateTime(val)
        }
        return nil
    }
    res["incidentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncidentId(val)
        }
        return nil
    }
    res["incidentWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncidentWebUrl(val)
        }
        return nil
    }
    res["lastActivityDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastActivityDateTime(val)
        }
        return nil
    }
    res["lastUpdateDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastUpdateDateTime(val)
        }
        return nil
    }
    res["mitreTechniques"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMitreTechniques(res)
        }
        return nil
    }
    res["productName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductName(val)
        }
        return nil
    }
    res["providerAlertId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProviderAlertId(val)
        }
        return nil
    }
    res["recommendedActions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecommendedActions(val)
        }
        return nil
    }
    res["resolvedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResolvedDateTime(val)
        }
        return nil
    }
    res["serviceSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServiceSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceSource(val.(*ServiceSource))
        }
        return nil
    }
    res["severity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAlertSeverity)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSeverity(val.(*AlertSeverity))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAlertStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*AlertStatus))
        }
        return nil
    }
    res["systemTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetSystemTags(res)
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    res["threatDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetThreatDisplayName(val)
        }
        return nil
    }
    res["threatFamilyName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetThreatFamilyName(val)
        }
        return nil
    }
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val)
        }
        return nil
    }
    return res
}
// GetFirstActivityDateTime gets the firstActivityDateTime property value. The earliest activity associated with the alert.
// returns a *Time when successful
func (m *Alert) GetFirstActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("firstActivityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetIncidentId gets the incidentId property value. Unique identifier to represent the incident this alert resource is associated with.
// returns a *string when successful
func (m *Alert) GetIncidentId()(*string) {
    val, err := m.GetBackingStore().Get("incidentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIncidentWebUrl gets the incidentWebUrl property value. URL for the incident page in the Microsoft 365 Defender portal.
// returns a *string when successful
func (m *Alert) GetIncidentWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("incidentWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastActivityDateTime gets the lastActivityDateTime property value. The oldest activity associated with the alert.
// returns a *Time when successful
func (m *Alert) GetLastActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastActivityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastUpdateDateTime gets the lastUpdateDateTime property value. Time when the alert was last updated at Microsoft 365 Defender.
// returns a *Time when successful
func (m *Alert) GetLastUpdateDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastUpdateDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMitreTechniques gets the mitreTechniques property value. The attack techniques, as aligned with the MITRE ATT&CK framework.
// returns a []string when successful
func (m *Alert) GetMitreTechniques()([]string) {
    val, err := m.GetBackingStore().Get("mitreTechniques")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetProductName gets the productName property value. The name of the product which published this alert.
// returns a *string when successful
func (m *Alert) GetProductName()(*string) {
    val, err := m.GetBackingStore().Get("productName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProviderAlertId gets the providerAlertId property value. The ID of the alert as it appears in the security provider product that generated the alert.
// returns a *string when successful
func (m *Alert) GetProviderAlertId()(*string) {
    val, err := m.GetBackingStore().Get("providerAlertId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecommendedActions gets the recommendedActions property value. Recommended response and remediation actions to take in the event this alert was generated.
// returns a *string when successful
func (m *Alert) GetRecommendedActions()(*string) {
    val, err := m.GetBackingStore().Get("recommendedActions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResolvedDateTime gets the resolvedDateTime property value. Time when the alert was resolved.
// returns a *Time when successful
func (m *Alert) GetResolvedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("resolvedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetServiceSource gets the serviceSource property value. The serviceSource property
// returns a *ServiceSource when successful
func (m *Alert) GetServiceSource()(*ServiceSource) {
    val, err := m.GetBackingStore().Get("serviceSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServiceSource)
    }
    return nil
}
// GetSeverity gets the severity property value. The severity property
// returns a *AlertSeverity when successful
func (m *Alert) GetSeverity()(*AlertSeverity) {
    val, err := m.GetBackingStore().Get("severity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AlertSeverity)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *AlertStatus when successful
func (m *Alert) GetStatus()(*AlertStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AlertStatus)
    }
    return nil
}
// GetSystemTags gets the systemTags property value. The system tags associated with the alert.
// returns a []string when successful
func (m *Alert) GetSystemTags()([]string) {
    val, err := m.GetBackingStore().Get("systemTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The Microsoft Entra tenant the alert was created in.
// returns a *string when successful
func (m *Alert) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetThreatDisplayName gets the threatDisplayName property value. The threat associated with this alert.
// returns a *string when successful
func (m *Alert) GetThreatDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("threatDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetThreatFamilyName gets the threatFamilyName property value. Threat family associated with this alert.
// returns a *string when successful
func (m *Alert) GetThreatFamilyName()(*string) {
    val, err := m.GetBackingStore().Get("threatFamilyName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTitle gets the title property value. Brief identifying string value describing the alert.
// returns a *string when successful
func (m *Alert) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Alert) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("actorDisplayName", m.GetActorDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("additionalData", m.GetAdditionalDataProperty())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("alertPolicyId", m.GetAlertPolicyId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("alertWebUrl", m.GetAlertWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("assignedTo", m.GetAssignedTo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("category", m.GetCategory())
        if err != nil {
            return err
        }
    }
    if m.GetClassification() != nil {
        cast := (*m.GetClassification()).String()
        err = writer.WriteStringValue("classification", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetComments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetComments()))
        for i, v := range m.GetComments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("comments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    if m.GetDetectionSource() != nil {
        cast := (*m.GetDetectionSource()).String()
        err = writer.WriteStringValue("detectionSource", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("detectorId", m.GetDetectorId())
        if err != nil {
            return err
        }
    }
    if m.GetDetermination() != nil {
        cast := (*m.GetDetermination()).String()
        err = writer.WriteStringValue("determination", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetEvidence() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEvidence()))
        for i, v := range m.GetEvidence() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("evidence", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("firstActivityDateTime", m.GetFirstActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("incidentId", m.GetIncidentId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("incidentWebUrl", m.GetIncidentWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastActivityDateTime", m.GetLastActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastUpdateDateTime", m.GetLastUpdateDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetMitreTechniques() != nil {
        err = writer.WriteCollectionOfStringValues("mitreTechniques", m.GetMitreTechniques())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productName", m.GetProductName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("providerAlertId", m.GetProviderAlertId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recommendedActions", m.GetRecommendedActions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("resolvedDateTime", m.GetResolvedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetServiceSource() != nil {
        cast := (*m.GetServiceSource()).String()
        err = writer.WriteStringValue("serviceSource", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSeverity() != nil {
        cast := (*m.GetSeverity()).String()
        err = writer.WriteStringValue("severity", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSystemTags() != nil {
        err = writer.WriteCollectionOfStringValues("systemTags", m.GetSystemTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("threatDisplayName", m.GetThreatDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("threatFamilyName", m.GetThreatFamilyName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActorDisplayName sets the actorDisplayName property value. The adversary or activity group that is associated with this alert.
func (m *Alert) SetActorDisplayName(value *string)() {
    err := m.GetBackingStore().Set("actorDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalDataProperty sets the additionalData property value. A collection of other alert properties, including user-defined properties. Any custom details defined in the alert, and any dynamic content in the alert details, are stored here.
func (m *Alert) SetAdditionalDataProperty(value Dictionaryable)() {
    err := m.GetBackingStore().Set("additionalDataProperty", value)
    if err != nil {
        panic(err)
    }
}
// SetAlertPolicyId sets the alertPolicyId property value. The ID of the policy that generated the alert, and populated when there is a specific policy that generated the alert, whether configured by a customer or a built-in policy.
func (m *Alert) SetAlertPolicyId(value *string)() {
    err := m.GetBackingStore().Set("alertPolicyId", value)
    if err != nil {
        panic(err)
    }
}
// SetAlertWebUrl sets the alertWebUrl property value. URL for the Microsoft 365 Defender portal alert page.
func (m *Alert) SetAlertWebUrl(value *string)() {
    err := m.GetBackingStore().Set("alertWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedTo sets the assignedTo property value. Owner of the alert, or null if no owner is assigned.
func (m *Alert) SetAssignedTo(value *string)() {
    err := m.GetBackingStore().Set("assignedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory sets the category property value. The attack kill-chain category that the alert belongs to. Aligned with the MITRE ATT&CK framework.
func (m *Alert) SetCategory(value *string)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetClassification sets the classification property value. Specifies whether the alert represents a true threat. Possible values are: unknown, falsePositive, truePositive, informationalExpectedActivity, unknownFutureValue.
func (m *Alert) SetClassification(value *AlertClassification)() {
    err := m.GetBackingStore().Set("classification", value)
    if err != nil {
        panic(err)
    }
}
// SetComments sets the comments property value. Array of comments created by the Security Operations (SecOps) team during the alert management process.
func (m *Alert) SetComments(value []AlertCommentable)() {
    err := m.GetBackingStore().Set("comments", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Time when Microsoft 365 Defender created the alert.
func (m *Alert) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. String value describing each alert.
func (m *Alert) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDetectionSource sets the detectionSource property value. Detection technology or sensor that identified the notable component or activity. Possible values are: unknown, microsoftDefenderForEndpoint, antivirus, smartScreen, customTi, microsoftDefenderForOffice365, automatedInvestigation, microsoftThreatExperts, customDetection, microsoftDefenderForIdentity, cloudAppSecurity, microsoft365Defender, azureAdIdentityProtection, manual, microsoftDataLossPrevention, appGovernancePolicy, appGovernanceDetection, unknownFutureValue, microsoftDefenderForCloud, microsoftDefenderForIoT, microsoftDefenderForServers, microsoftDefenderForStorage, microsoftDefenderForDNS, microsoftDefenderForDatabases, microsoftDefenderForContainers, microsoftDefenderForNetwork, microsoftDefenderForAppService, microsoftDefenderForKeyVault, microsoftDefenderForResourceManager, microsoftDefenderForApiManagement, microsoftSentinel, nrtAlerts, scheduledAlerts, microsoftDefenderThreatIntelligenceAnalytics, builtInMl. You must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: microsoftDefenderForCloud, microsoftDefenderForIoT, microsoftDefenderForServers, microsoftDefenderForStorage, microsoftDefenderForDNS, microsoftDefenderForDatabases, microsoftDefenderForContainers, microsoftDefenderForNetwork, microsoftDefenderForAppService, microsoftDefenderForKeyVault, microsoftDefenderForResourceManager, microsoftDefenderForApiManagement, microsoftSentinel, nrtAlerts, scheduledAlerts, microsoftDefenderThreatIntelligenceAnalytics, builtInMl.
func (m *Alert) SetDetectionSource(value *DetectionSource)() {
    err := m.GetBackingStore().Set("detectionSource", value)
    if err != nil {
        panic(err)
    }
}
// SetDetectorId sets the detectorId property value. The ID of the detector that triggered the alert.
func (m *Alert) SetDetectorId(value *string)() {
    err := m.GetBackingStore().Set("detectorId", value)
    if err != nil {
        panic(err)
    }
}
// SetDetermination sets the determination property value. Specifies the result of the investigation, whether the alert represents a true attack and if so, the nature of the attack. Possible values are: unknown, apt, malware, securityPersonnel, securityTesting, unwantedSoftware, other, multiStagedAttack, compromisedAccount, phishing, maliciousUserActivity, notMalicious, notEnoughDataToValidate, confirmedUserActivity, lineOfBusinessApplication, unknownFutureValue.
func (m *Alert) SetDetermination(value *AlertDetermination)() {
    err := m.GetBackingStore().Set("determination", value)
    if err != nil {
        panic(err)
    }
}
// SetEvidence sets the evidence property value. Collection of evidence related to the alert.
func (m *Alert) SetEvidence(value []AlertEvidenceable)() {
    err := m.GetBackingStore().Set("evidence", value)
    if err != nil {
        panic(err)
    }
}
// SetFirstActivityDateTime sets the firstActivityDateTime property value. The earliest activity associated with the alert.
func (m *Alert) SetFirstActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("firstActivityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetIncidentId sets the incidentId property value. Unique identifier to represent the incident this alert resource is associated with.
func (m *Alert) SetIncidentId(value *string)() {
    err := m.GetBackingStore().Set("incidentId", value)
    if err != nil {
        panic(err)
    }
}
// SetIncidentWebUrl sets the incidentWebUrl property value. URL for the incident page in the Microsoft 365 Defender portal.
func (m *Alert) SetIncidentWebUrl(value *string)() {
    err := m.GetBackingStore().Set("incidentWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetLastActivityDateTime sets the lastActivityDateTime property value. The oldest activity associated with the alert.
func (m *Alert) SetLastActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastActivityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdateDateTime sets the lastUpdateDateTime property value. Time when the alert was last updated at Microsoft 365 Defender.
func (m *Alert) SetLastUpdateDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdateDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMitreTechniques sets the mitreTechniques property value. The attack techniques, as aligned with the MITRE ATT&CK framework.
func (m *Alert) SetMitreTechniques(value []string)() {
    err := m.GetBackingStore().Set("mitreTechniques", value)
    if err != nil {
        panic(err)
    }
}
// SetProductName sets the productName property value. The name of the product which published this alert.
func (m *Alert) SetProductName(value *string)() {
    err := m.GetBackingStore().Set("productName", value)
    if err != nil {
        panic(err)
    }
}
// SetProviderAlertId sets the providerAlertId property value. The ID of the alert as it appears in the security provider product that generated the alert.
func (m *Alert) SetProviderAlertId(value *string)() {
    err := m.GetBackingStore().Set("providerAlertId", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendedActions sets the recommendedActions property value. Recommended response and remediation actions to take in the event this alert was generated.
func (m *Alert) SetRecommendedActions(value *string)() {
    err := m.GetBackingStore().Set("recommendedActions", value)
    if err != nil {
        panic(err)
    }
}
// SetResolvedDateTime sets the resolvedDateTime property value. Time when the alert was resolved.
func (m *Alert) SetResolvedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("resolvedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceSource sets the serviceSource property value. The serviceSource property
func (m *Alert) SetServiceSource(value *ServiceSource)() {
    err := m.GetBackingStore().Set("serviceSource", value)
    if err != nil {
        panic(err)
    }
}
// SetSeverity sets the severity property value. The severity property
func (m *Alert) SetSeverity(value *AlertSeverity)() {
    err := m.GetBackingStore().Set("severity", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *Alert) SetStatus(value *AlertStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetSystemTags sets the systemTags property value. The system tags associated with the alert.
func (m *Alert) SetSystemTags(value []string)() {
    err := m.GetBackingStore().Set("systemTags", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The Microsoft Entra tenant the alert was created in.
func (m *Alert) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetThreatDisplayName sets the threatDisplayName property value. The threat associated with this alert.
func (m *Alert) SetThreatDisplayName(value *string)() {
    err := m.GetBackingStore().Set("threatDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetThreatFamilyName sets the threatFamilyName property value. Threat family associated with this alert.
func (m *Alert) SetThreatFamilyName(value *string)() {
    err := m.GetBackingStore().Set("threatFamilyName", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. Brief identifying string value describing the alert.
func (m *Alert) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
type Alertable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActorDisplayName()(*string)
    GetAdditionalDataProperty()(Dictionaryable)
    GetAlertPolicyId()(*string)
    GetAlertWebUrl()(*string)
    GetAssignedTo()(*string)
    GetCategory()(*string)
    GetClassification()(*AlertClassification)
    GetComments()([]AlertCommentable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDetectionSource()(*DetectionSource)
    GetDetectorId()(*string)
    GetDetermination()(*AlertDetermination)
    GetEvidence()([]AlertEvidenceable)
    GetFirstActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIncidentId()(*string)
    GetIncidentWebUrl()(*string)
    GetLastActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastUpdateDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMitreTechniques()([]string)
    GetProductName()(*string)
    GetProviderAlertId()(*string)
    GetRecommendedActions()(*string)
    GetResolvedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetServiceSource()(*ServiceSource)
    GetSeverity()(*AlertSeverity)
    GetStatus()(*AlertStatus)
    GetSystemTags()([]string)
    GetTenantId()(*string)
    GetThreatDisplayName()(*string)
    GetThreatFamilyName()(*string)
    GetTitle()(*string)
    SetActorDisplayName(value *string)()
    SetAdditionalDataProperty(value Dictionaryable)()
    SetAlertPolicyId(value *string)()
    SetAlertWebUrl(value *string)()
    SetAssignedTo(value *string)()
    SetCategory(value *string)()
    SetClassification(value *AlertClassification)()
    SetComments(value []AlertCommentable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDetectionSource(value *DetectionSource)()
    SetDetectorId(value *string)()
    SetDetermination(value *AlertDetermination)()
    SetEvidence(value []AlertEvidenceable)()
    SetFirstActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIncidentId(value *string)()
    SetIncidentWebUrl(value *string)()
    SetLastActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastUpdateDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMitreTechniques(value []string)()
    SetProductName(value *string)()
    SetProviderAlertId(value *string)()
    SetRecommendedActions(value *string)()
    SetResolvedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetServiceSource(value *ServiceSource)()
    SetSeverity(value *AlertSeverity)()
    SetStatus(value *AlertStatus)()
    SetSystemTags(value []string)()
    SetTenantId(value *string)()
    SetThreatDisplayName(value *string)()
    SetThreatFamilyName(value *string)()
    SetTitle(value *string)()
}
