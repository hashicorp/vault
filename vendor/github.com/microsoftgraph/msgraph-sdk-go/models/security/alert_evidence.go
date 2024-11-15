package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AlertEvidence struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAlertEvidence instantiates a new AlertEvidence and sets the default values.
func NewAlertEvidence()(*AlertEvidence) {
    m := &AlertEvidence{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAlertEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAlertEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.security.amazonResourceEvidence":
                        return NewAmazonResourceEvidence(), nil
                    case "#microsoft.graph.security.analyzedMessageEvidence":
                        return NewAnalyzedMessageEvidence(), nil
                    case "#microsoft.graph.security.azureResourceEvidence":
                        return NewAzureResourceEvidence(), nil
                    case "#microsoft.graph.security.blobContainerEvidence":
                        return NewBlobContainerEvidence(), nil
                    case "#microsoft.graph.security.blobEvidence":
                        return NewBlobEvidence(), nil
                    case "#microsoft.graph.security.cloudApplicationEvidence":
                        return NewCloudApplicationEvidence(), nil
                    case "#microsoft.graph.security.cloudLogonRequestEvidence":
                        return NewCloudLogonRequestEvidence(), nil
                    case "#microsoft.graph.security.cloudLogonSessionEvidence":
                        return NewCloudLogonSessionEvidence(), nil
                    case "#microsoft.graph.security.containerEvidence":
                        return NewContainerEvidence(), nil
                    case "#microsoft.graph.security.containerImageEvidence":
                        return NewContainerImageEvidence(), nil
                    case "#microsoft.graph.security.containerRegistryEvidence":
                        return NewContainerRegistryEvidence(), nil
                    case "#microsoft.graph.security.deviceEvidence":
                        return NewDeviceEvidence(), nil
                    case "#microsoft.graph.security.dnsEvidence":
                        return NewDnsEvidence(), nil
                    case "#microsoft.graph.security.fileEvidence":
                        return NewFileEvidence(), nil
                    case "#microsoft.graph.security.fileHashEvidence":
                        return NewFileHashEvidence(), nil
                    case "#microsoft.graph.security.gitHubOrganizationEvidence":
                        return NewGitHubOrganizationEvidence(), nil
                    case "#microsoft.graph.security.gitHubRepoEvidence":
                        return NewGitHubRepoEvidence(), nil
                    case "#microsoft.graph.security.gitHubUserEvidence":
                        return NewGitHubUserEvidence(), nil
                    case "#microsoft.graph.security.googleCloudResourceEvidence":
                        return NewGoogleCloudResourceEvidence(), nil
                    case "#microsoft.graph.security.hostLogonSessionEvidence":
                        return NewHostLogonSessionEvidence(), nil
                    case "#microsoft.graph.security.ioTDeviceEvidence":
                        return NewIoTDeviceEvidence(), nil
                    case "#microsoft.graph.security.ipEvidence":
                        return NewIpEvidence(), nil
                    case "#microsoft.graph.security.kubernetesClusterEvidence":
                        return NewKubernetesClusterEvidence(), nil
                    case "#microsoft.graph.security.kubernetesControllerEvidence":
                        return NewKubernetesControllerEvidence(), nil
                    case "#microsoft.graph.security.kubernetesNamespaceEvidence":
                        return NewKubernetesNamespaceEvidence(), nil
                    case "#microsoft.graph.security.kubernetesPodEvidence":
                        return NewKubernetesPodEvidence(), nil
                    case "#microsoft.graph.security.kubernetesSecretEvidence":
                        return NewKubernetesSecretEvidence(), nil
                    case "#microsoft.graph.security.kubernetesServiceAccountEvidence":
                        return NewKubernetesServiceAccountEvidence(), nil
                    case "#microsoft.graph.security.kubernetesServiceEvidence":
                        return NewKubernetesServiceEvidence(), nil
                    case "#microsoft.graph.security.mailboxConfigurationEvidence":
                        return NewMailboxConfigurationEvidence(), nil
                    case "#microsoft.graph.security.mailboxEvidence":
                        return NewMailboxEvidence(), nil
                    case "#microsoft.graph.security.mailClusterEvidence":
                        return NewMailClusterEvidence(), nil
                    case "#microsoft.graph.security.malwareEvidence":
                        return NewMalwareEvidence(), nil
                    case "#microsoft.graph.security.networkConnectionEvidence":
                        return NewNetworkConnectionEvidence(), nil
                    case "#microsoft.graph.security.nicEvidence":
                        return NewNicEvidence(), nil
                    case "#microsoft.graph.security.oauthApplicationEvidence":
                        return NewOauthApplicationEvidence(), nil
                    case "#microsoft.graph.security.processEvidence":
                        return NewProcessEvidence(), nil
                    case "#microsoft.graph.security.registryKeyEvidence":
                        return NewRegistryKeyEvidence(), nil
                    case "#microsoft.graph.security.registryValueEvidence":
                        return NewRegistryValueEvidence(), nil
                    case "#microsoft.graph.security.sasTokenEvidence":
                        return NewSasTokenEvidence(), nil
                    case "#microsoft.graph.security.securityGroupEvidence":
                        return NewSecurityGroupEvidence(), nil
                    case "#microsoft.graph.security.servicePrincipalEvidence":
                        return NewServicePrincipalEvidence(), nil
                    case "#microsoft.graph.security.submissionMailEvidence":
                        return NewSubmissionMailEvidence(), nil
                    case "#microsoft.graph.security.urlEvidence":
                        return NewUrlEvidence(), nil
                    case "#microsoft.graph.security.userEvidence":
                        return NewUserEvidence(), nil
                }
            }
        }
    }
    return NewAlertEvidence(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AlertEvidence) GetAdditionalData()(map[string]any) {
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
func (m *AlertEvidence) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the evidence was created and added to the alert. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *AlertEvidence) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDetailedRoles gets the detailedRoles property value. Detailed description of the entity role/s in an alert. Values are free-form.
// returns a []string when successful
func (m *AlertEvidence) GetDetailedRoles()([]string) {
    val, err := m.GetBackingStore().Get("detailedRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AlertEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["detailedRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetDetailedRoles(res)
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
    res["remediationStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEvidenceRemediationStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemediationStatus(val.(*EvidenceRemediationStatus))
        }
        return nil
    }
    res["remediationStatusDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemediationStatusDetails(val)
        }
        return nil
    }
    res["roles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseEvidenceRole)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EvidenceRole, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*EvidenceRole))
                }
            }
            m.SetRoles(res)
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetTags(res)
        }
        return nil
    }
    res["verdict"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEvidenceVerdict)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVerdict(val.(*EvidenceVerdict))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AlertEvidence) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRemediationStatus gets the remediationStatus property value. The remediationStatus property
// returns a *EvidenceRemediationStatus when successful
func (m *AlertEvidence) GetRemediationStatus()(*EvidenceRemediationStatus) {
    val, err := m.GetBackingStore().Get("remediationStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EvidenceRemediationStatus)
    }
    return nil
}
// GetRemediationStatusDetails gets the remediationStatusDetails property value. Details about the remediation status.
// returns a *string when successful
func (m *AlertEvidence) GetRemediationStatusDetails()(*string) {
    val, err := m.GetBackingStore().Get("remediationStatusDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoles gets the roles property value. The role/s that an evidence entity represents in an alert, for example, an IP address that is associated with an attacker has the evidence role Attacker.
// returns a []EvidenceRole when successful
func (m *AlertEvidence) GetRoles()([]EvidenceRole) {
    val, err := m.GetBackingStore().Get("roles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EvidenceRole)
    }
    return nil
}
// GetTags gets the tags property value. Array of custom tags associated with an evidence instance, for example, to denote a group of devices, high-value assets, etc.
// returns a []string when successful
func (m *AlertEvidence) GetTags()([]string) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetVerdict gets the verdict property value. The verdict property
// returns a *EvidenceVerdict when successful
func (m *AlertEvidence) GetVerdict()(*EvidenceVerdict) {
    val, err := m.GetBackingStore().Get("verdict")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EvidenceVerdict)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AlertEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetDetailedRoles() != nil {
        err := writer.WriteCollectionOfStringValues("detailedRoles", m.GetDetailedRoles())
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
    if m.GetRemediationStatus() != nil {
        cast := (*m.GetRemediationStatus()).String()
        err := writer.WriteStringValue("remediationStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("remediationStatusDetails", m.GetRemediationStatusDetails())
        if err != nil {
            return err
        }
    }
    if m.GetRoles() != nil {
        err := writer.WriteCollectionOfStringValues("roles", SerializeEvidenceRole(m.GetRoles()))
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        err := writer.WriteCollectionOfStringValues("tags", m.GetTags())
        if err != nil {
            return err
        }
    }
    if m.GetVerdict() != nil {
        cast := (*m.GetVerdict()).String()
        err := writer.WriteStringValue("verdict", &cast)
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
func (m *AlertEvidence) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AlertEvidence) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the evidence was created and added to the alert. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *AlertEvidence) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDetailedRoles sets the detailedRoles property value. Detailed description of the entity role/s in an alert. Values are free-form.
func (m *AlertEvidence) SetDetailedRoles(value []string)() {
    err := m.GetBackingStore().Set("detailedRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AlertEvidence) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRemediationStatus sets the remediationStatus property value. The remediationStatus property
func (m *AlertEvidence) SetRemediationStatus(value *EvidenceRemediationStatus)() {
    err := m.GetBackingStore().Set("remediationStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetRemediationStatusDetails sets the remediationStatusDetails property value. Details about the remediation status.
func (m *AlertEvidence) SetRemediationStatusDetails(value *string)() {
    err := m.GetBackingStore().Set("remediationStatusDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetRoles sets the roles property value. The role/s that an evidence entity represents in an alert, for example, an IP address that is associated with an attacker has the evidence role Attacker.
func (m *AlertEvidence) SetRoles(value []EvidenceRole)() {
    err := m.GetBackingStore().Set("roles", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. Array of custom tags associated with an evidence instance, for example, to denote a group of devices, high-value assets, etc.
func (m *AlertEvidence) SetTags(value []string)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetVerdict sets the verdict property value. The verdict property
func (m *AlertEvidence) SetVerdict(value *EvidenceVerdict)() {
    err := m.GetBackingStore().Set("verdict", value)
    if err != nil {
        panic(err)
    }
}
type AlertEvidenceable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDetailedRoles()([]string)
    GetOdataType()(*string)
    GetRemediationStatus()(*EvidenceRemediationStatus)
    GetRemediationStatusDetails()(*string)
    GetRoles()([]EvidenceRole)
    GetTags()([]string)
    GetVerdict()(*EvidenceVerdict)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDetailedRoles(value []string)()
    SetOdataType(value *string)()
    SetRemediationStatus(value *EvidenceRemediationStatus)()
    SetRemediationStatusDetails(value *string)()
    SetRoles(value []EvidenceRole)()
    SetTags(value []string)()
    SetVerdict(value *EvidenceVerdict)()
}
