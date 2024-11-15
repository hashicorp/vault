package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type HealthIssue struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewHealthIssue instantiates a new HealthIssue and sets the default values.
func NewHealthIssue()(*HealthIssue) {
    m := &HealthIssue{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateHealthIssueFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHealthIssueFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHealthIssue(), nil
}
// GetAdditionalInformation gets the additionalInformation property value. Contains additional information about the issue, such as a list of items to fix.
// returns a []string when successful
func (m *HealthIssue) GetAdditionalInformation()([]string) {
    val, err := m.GetBackingStore().Get("additionalInformation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the health issue was generated. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *HealthIssue) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. Contains more detailed information about the health issue.
// returns a *string when successful
func (m *HealthIssue) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name of the health issue.
// returns a *string when successful
func (m *HealthIssue) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDomainNames gets the domainNames property value. A list of the fully qualified domain names of the domains or the sensors the health issue is related to.
// returns a []string when successful
func (m *HealthIssue) GetDomainNames()([]string) {
    val, err := m.GetBackingStore().Get("domainNames")
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
func (m *HealthIssue) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["additionalInformation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAdditionalInformation(res)
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
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["domainNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetDomainNames(res)
        }
        return nil
    }
    res["healthIssueType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseHealthIssueType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHealthIssueType(val.(*HealthIssueType))
        }
        return nil
    }
    res["issueTypeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIssueTypeId(val)
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["recommendations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetRecommendations(res)
        }
        return nil
    }
    res["recommendedActionCommands"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetRecommendedActionCommands(res)
        }
        return nil
    }
    res["sensorDNSNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSensorDNSNames(res)
        }
        return nil
    }
    res["severity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseHealthIssueSeverity)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSeverity(val.(*HealthIssueSeverity))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseHealthIssueStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*HealthIssueStatus))
        }
        return nil
    }
    return res
}
// GetHealthIssueType gets the healthIssueType property value. The type of the health issue. The possible values are: sensor, global, unknownFutureValue. For a list of all health issues and their identifiers, see Microsoft Defender for Identity health issues.
// returns a *HealthIssueType when successful
func (m *HealthIssue) GetHealthIssueType()(*HealthIssueType) {
    val, err := m.GetBackingStore().Get("healthIssueType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*HealthIssueType)
    }
    return nil
}
// GetIssueTypeId gets the issueTypeId property value. The type identifier of the health issue. For a list of all health issues and their identifiers, see Microsoft Defender for Identity health issues.
// returns a *string when successful
func (m *HealthIssue) GetIssueTypeId()(*string) {
    val, err := m.GetBackingStore().Get("issueTypeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the health issue was last updated. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *HealthIssue) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRecommendations gets the recommendations property value. A list of recommended actions that can be taken to resolve the issue effectively and efficiently. These actions might include instructions for further investigation and aren't limited to prewritten responses.
// returns a []string when successful
func (m *HealthIssue) GetRecommendations()([]string) {
    val, err := m.GetBackingStore().Get("recommendations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetRecommendedActionCommands gets the recommendedActionCommands property value. A list of commands from the PowerShell module for the product that can be used to resolve the issue, if available. If no commands can be used to solve the issue, this property is empty. The commands, if present, provide a quick and efficient way to address the issue. These commands run in sequence for the single recommended fix.
// returns a []string when successful
func (m *HealthIssue) GetRecommendedActionCommands()([]string) {
    val, err := m.GetBackingStore().Get("recommendedActionCommands")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSensorDNSNames gets the sensorDNSNames property value. A list of the DNS names of the sensors the health issue is related to.
// returns a []string when successful
func (m *HealthIssue) GetSensorDNSNames()([]string) {
    val, err := m.GetBackingStore().Get("sensorDNSNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSeverity gets the severity property value. The severity of the health issue. The possible values are: low, medium, high, unknownFutureValue.
// returns a *HealthIssueSeverity when successful
func (m *HealthIssue) GetSeverity()(*HealthIssueSeverity) {
    val, err := m.GetBackingStore().Get("severity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*HealthIssueSeverity)
    }
    return nil
}
// GetStatus gets the status property value. The status of the health issue. The possible values are: open, closed, suppressed, unknownFutureValue.
// returns a *HealthIssueStatus when successful
func (m *HealthIssue) GetStatus()(*HealthIssueStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*HealthIssueStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *HealthIssue) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAdditionalInformation() != nil {
        err = writer.WriteCollectionOfStringValues("additionalInformation", m.GetAdditionalInformation())
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
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetDomainNames() != nil {
        err = writer.WriteCollectionOfStringValues("domainNames", m.GetDomainNames())
        if err != nil {
            return err
        }
    }
    if m.GetHealthIssueType() != nil {
        cast := (*m.GetHealthIssueType()).String()
        err = writer.WriteStringValue("healthIssueType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("issueTypeId", m.GetIssueTypeId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetRecommendations() != nil {
        err = writer.WriteCollectionOfStringValues("recommendations", m.GetRecommendations())
        if err != nil {
            return err
        }
    }
    if m.GetRecommendedActionCommands() != nil {
        err = writer.WriteCollectionOfStringValues("recommendedActionCommands", m.GetRecommendedActionCommands())
        if err != nil {
            return err
        }
    }
    if m.GetSensorDNSNames() != nil {
        err = writer.WriteCollectionOfStringValues("sensorDNSNames", m.GetSensorDNSNames())
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
    return nil
}
// SetAdditionalInformation sets the additionalInformation property value. Contains additional information about the issue, such as a list of items to fix.
func (m *HealthIssue) SetAdditionalInformation(value []string)() {
    err := m.GetBackingStore().Set("additionalInformation", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the health issue was generated. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *HealthIssue) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Contains more detailed information about the health issue.
func (m *HealthIssue) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name of the health issue.
func (m *HealthIssue) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDomainNames sets the domainNames property value. A list of the fully qualified domain names of the domains or the sensors the health issue is related to.
func (m *HealthIssue) SetDomainNames(value []string)() {
    err := m.GetBackingStore().Set("domainNames", value)
    if err != nil {
        panic(err)
    }
}
// SetHealthIssueType sets the healthIssueType property value. The type of the health issue. The possible values are: sensor, global, unknownFutureValue. For a list of all health issues and their identifiers, see Microsoft Defender for Identity health issues.
func (m *HealthIssue) SetHealthIssueType(value *HealthIssueType)() {
    err := m.GetBackingStore().Set("healthIssueType", value)
    if err != nil {
        panic(err)
    }
}
// SetIssueTypeId sets the issueTypeId property value. The type identifier of the health issue. For a list of all health issues and their identifiers, see Microsoft Defender for Identity health issues.
func (m *HealthIssue) SetIssueTypeId(value *string)() {
    err := m.GetBackingStore().Set("issueTypeId", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the health issue was last updated. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *HealthIssue) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendations sets the recommendations property value. A list of recommended actions that can be taken to resolve the issue effectively and efficiently. These actions might include instructions for further investigation and aren't limited to prewritten responses.
func (m *HealthIssue) SetRecommendations(value []string)() {
    err := m.GetBackingStore().Set("recommendations", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendedActionCommands sets the recommendedActionCommands property value. A list of commands from the PowerShell module for the product that can be used to resolve the issue, if available. If no commands can be used to solve the issue, this property is empty. The commands, if present, provide a quick and efficient way to address the issue. These commands run in sequence for the single recommended fix.
func (m *HealthIssue) SetRecommendedActionCommands(value []string)() {
    err := m.GetBackingStore().Set("recommendedActionCommands", value)
    if err != nil {
        panic(err)
    }
}
// SetSensorDNSNames sets the sensorDNSNames property value. A list of the DNS names of the sensors the health issue is related to.
func (m *HealthIssue) SetSensorDNSNames(value []string)() {
    err := m.GetBackingStore().Set("sensorDNSNames", value)
    if err != nil {
        panic(err)
    }
}
// SetSeverity sets the severity property value. The severity of the health issue. The possible values are: low, medium, high, unknownFutureValue.
func (m *HealthIssue) SetSeverity(value *HealthIssueSeverity)() {
    err := m.GetBackingStore().Set("severity", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status of the health issue. The possible values are: open, closed, suppressed, unknownFutureValue.
func (m *HealthIssue) SetStatus(value *HealthIssueStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type HealthIssueable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAdditionalInformation()([]string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetDomainNames()([]string)
    GetHealthIssueType()(*HealthIssueType)
    GetIssueTypeId()(*string)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRecommendations()([]string)
    GetRecommendedActionCommands()([]string)
    GetSensorDNSNames()([]string)
    GetSeverity()(*HealthIssueSeverity)
    GetStatus()(*HealthIssueStatus)
    SetAdditionalInformation(value []string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetDomainNames(value []string)()
    SetHealthIssueType(value *HealthIssueType)()
    SetIssueTypeId(value *string)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRecommendations(value []string)()
    SetRecommendedActionCommands(value []string)()
    SetSensorDNSNames(value []string)()
    SetSeverity(value *HealthIssueSeverity)()
    SetStatus(value *HealthIssueStatus)()
}
