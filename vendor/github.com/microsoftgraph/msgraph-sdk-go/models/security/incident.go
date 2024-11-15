package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Incident struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewIncident instantiates a new Incident and sets the default values.
func NewIncident()(*Incident) {
    m := &Incident{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateIncidentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIncidentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIncident(), nil
}
// GetAlerts gets the alerts property value. The list of related alerts. Supports $expand.
// returns a []Alertable when successful
func (m *Incident) GetAlerts()([]Alertable) {
    val, err := m.GetBackingStore().Get("alerts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Alertable)
    }
    return nil
}
// GetAssignedTo gets the assignedTo property value. Owner of the incident, or null if no owner is assigned. Free editable text.
// returns a *string when successful
func (m *Incident) GetAssignedTo()(*string) {
    val, err := m.GetBackingStore().Get("assignedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetClassification gets the classification property value. The specification for the incident. Possible values are: unknown, falsePositive, truePositive, informationalExpectedActivity, unknownFutureValue.
// returns a *AlertClassification when successful
func (m *Incident) GetClassification()(*AlertClassification) {
    val, err := m.GetBackingStore().Get("classification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AlertClassification)
    }
    return nil
}
// GetComments gets the comments property value. Array of comments created by the Security Operations (SecOps) team when the incident is managed.
// returns a []AlertCommentable when successful
func (m *Incident) GetComments()([]AlertCommentable) {
    val, err := m.GetBackingStore().Get("comments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AlertCommentable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Time when the incident was first created.
// returns a *Time when successful
func (m *Incident) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCustomTags gets the customTags property value. Array of custom tags associated with an incident.
// returns a []string when successful
func (m *Incident) GetCustomTags()([]string) {
    val, err := m.GetBackingStore().Get("customTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDescription gets the description property value. Description of the incident.
// returns a *string when successful
func (m *Incident) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDetermination gets the determination property value. Specifies the determination of the incident. Possible values are: unknown, apt, malware, securityPersonnel, securityTesting, unwantedSoftware, other, multiStagedAttack, compromisedUser, phishing, maliciousUserActivity, clean, insufficientData, confirmedUserActivity, lineOfBusinessApplication, unknownFutureValue.
// returns a *AlertDetermination when successful
func (m *Incident) GetDetermination()(*AlertDetermination) {
    val, err := m.GetBackingStore().Get("determination")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AlertDetermination)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The incident name.
// returns a *string when successful
func (m *Incident) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Incident) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["alerts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAlertFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Alertable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Alertable)
                }
            }
            m.SetAlerts(res)
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
    res["customTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetCustomTags(res)
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
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val)
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
    res["redirectIncidentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRedirectIncidentId(val)
        }
        return nil
    }
    res["resolvingComment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResolvingComment(val)
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
        val, err := n.GetEnumValue(ParseIncidentStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*IncidentStatus))
        }
        return nil
    }
    res["summary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSummary(val)
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
    return res
}
// GetIncidentWebUrl gets the incidentWebUrl property value. The URL for the incident page in the Microsoft 365 Defender portal.
// returns a *string when successful
func (m *Incident) GetIncidentWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("incidentWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. The identity that last modified the incident.
// returns a *string when successful
func (m *Incident) GetLastModifiedBy()(*string) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastUpdateDateTime gets the lastUpdateDateTime property value. Time when the incident was last updated.
// returns a *Time when successful
func (m *Incident) GetLastUpdateDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastUpdateDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRedirectIncidentId gets the redirectIncidentId property value. Only populated in case an incident is grouped with another incident, as part of the logic that processes incidents. In such a case, the status property is redirected.
// returns a *string when successful
func (m *Incident) GetRedirectIncidentId()(*string) {
    val, err := m.GetBackingStore().Get("redirectIncidentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResolvingComment gets the resolvingComment property value. User input that explains the resolution of the incident and the classification choice. This property contains free editable text.
// returns a *string when successful
func (m *Incident) GetResolvingComment()(*string) {
    val, err := m.GetBackingStore().Get("resolvingComment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSeverity gets the severity property value. The severity property
// returns a *AlertSeverity when successful
func (m *Incident) GetSeverity()(*AlertSeverity) {
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
// returns a *IncidentStatus when successful
func (m *Incident) GetStatus()(*IncidentStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*IncidentStatus)
    }
    return nil
}
// GetSummary gets the summary property value. The overview of an attack. When applicable, the summary contains details of what occurred, impacted assets, and the type of attack.
// returns a *string when successful
func (m *Incident) GetSummary()(*string) {
    val, err := m.GetBackingStore().Get("summary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSystemTags gets the systemTags property value. The system tags associated with the incident.
// returns a []string when successful
func (m *Incident) GetSystemTags()([]string) {
    val, err := m.GetBackingStore().Get("systemTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The Microsoft Entra tenant in which the alert was created.
// returns a *string when successful
func (m *Incident) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Incident) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAlerts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAlerts()))
        for i, v := range m.GetAlerts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("alerts", cast)
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
    if m.GetCustomTags() != nil {
        err = writer.WriteCollectionOfStringValues("customTags", m.GetCustomTags())
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
    if m.GetDetermination() != nil {
        cast := (*m.GetDetermination()).String()
        err = writer.WriteStringValue("determination", &cast)
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
    {
        err = writer.WriteStringValue("incidentWebUrl", m.GetIncidentWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lastModifiedBy", m.GetLastModifiedBy())
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
    {
        err = writer.WriteStringValue("redirectIncidentId", m.GetRedirectIncidentId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resolvingComment", m.GetResolvingComment())
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
    {
        err = writer.WriteStringValue("summary", m.GetSummary())
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
    return nil
}
// SetAlerts sets the alerts property value. The list of related alerts. Supports $expand.
func (m *Incident) SetAlerts(value []Alertable)() {
    err := m.GetBackingStore().Set("alerts", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedTo sets the assignedTo property value. Owner of the incident, or null if no owner is assigned. Free editable text.
func (m *Incident) SetAssignedTo(value *string)() {
    err := m.GetBackingStore().Set("assignedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetClassification sets the classification property value. The specification for the incident. Possible values are: unknown, falsePositive, truePositive, informationalExpectedActivity, unknownFutureValue.
func (m *Incident) SetClassification(value *AlertClassification)() {
    err := m.GetBackingStore().Set("classification", value)
    if err != nil {
        panic(err)
    }
}
// SetComments sets the comments property value. Array of comments created by the Security Operations (SecOps) team when the incident is managed.
func (m *Incident) SetComments(value []AlertCommentable)() {
    err := m.GetBackingStore().Set("comments", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Time when the incident was first created.
func (m *Incident) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomTags sets the customTags property value. Array of custom tags associated with an incident.
func (m *Incident) SetCustomTags(value []string)() {
    err := m.GetBackingStore().Set("customTags", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description of the incident.
func (m *Incident) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDetermination sets the determination property value. Specifies the determination of the incident. Possible values are: unknown, apt, malware, securityPersonnel, securityTesting, unwantedSoftware, other, multiStagedAttack, compromisedUser, phishing, maliciousUserActivity, clean, insufficientData, confirmedUserActivity, lineOfBusinessApplication, unknownFutureValue.
func (m *Incident) SetDetermination(value *AlertDetermination)() {
    err := m.GetBackingStore().Set("determination", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The incident name.
func (m *Incident) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIncidentWebUrl sets the incidentWebUrl property value. The URL for the incident page in the Microsoft 365 Defender portal.
func (m *Incident) SetIncidentWebUrl(value *string)() {
    err := m.GetBackingStore().Set("incidentWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The identity that last modified the incident.
func (m *Incident) SetLastModifiedBy(value *string)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdateDateTime sets the lastUpdateDateTime property value. Time when the incident was last updated.
func (m *Incident) SetLastUpdateDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdateDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRedirectIncidentId sets the redirectIncidentId property value. Only populated in case an incident is grouped with another incident, as part of the logic that processes incidents. In such a case, the status property is redirected.
func (m *Incident) SetRedirectIncidentId(value *string)() {
    err := m.GetBackingStore().Set("redirectIncidentId", value)
    if err != nil {
        panic(err)
    }
}
// SetResolvingComment sets the resolvingComment property value. User input that explains the resolution of the incident and the classification choice. This property contains free editable text.
func (m *Incident) SetResolvingComment(value *string)() {
    err := m.GetBackingStore().Set("resolvingComment", value)
    if err != nil {
        panic(err)
    }
}
// SetSeverity sets the severity property value. The severity property
func (m *Incident) SetSeverity(value *AlertSeverity)() {
    err := m.GetBackingStore().Set("severity", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *Incident) SetStatus(value *IncidentStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetSummary sets the summary property value. The overview of an attack. When applicable, the summary contains details of what occurred, impacted assets, and the type of attack.
func (m *Incident) SetSummary(value *string)() {
    err := m.GetBackingStore().Set("summary", value)
    if err != nil {
        panic(err)
    }
}
// SetSystemTags sets the systemTags property value. The system tags associated with the incident.
func (m *Incident) SetSystemTags(value []string)() {
    err := m.GetBackingStore().Set("systemTags", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The Microsoft Entra tenant in which the alert was created.
func (m *Incident) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
type Incidentable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAlerts()([]Alertable)
    GetAssignedTo()(*string)
    GetClassification()(*AlertClassification)
    GetComments()([]AlertCommentable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCustomTags()([]string)
    GetDescription()(*string)
    GetDetermination()(*AlertDetermination)
    GetDisplayName()(*string)
    GetIncidentWebUrl()(*string)
    GetLastModifiedBy()(*string)
    GetLastUpdateDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRedirectIncidentId()(*string)
    GetResolvingComment()(*string)
    GetSeverity()(*AlertSeverity)
    GetStatus()(*IncidentStatus)
    GetSummary()(*string)
    GetSystemTags()([]string)
    GetTenantId()(*string)
    SetAlerts(value []Alertable)()
    SetAssignedTo(value *string)()
    SetClassification(value *AlertClassification)()
    SetComments(value []AlertCommentable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCustomTags(value []string)()
    SetDescription(value *string)()
    SetDetermination(value *AlertDetermination)()
    SetDisplayName(value *string)()
    SetIncidentWebUrl(value *string)()
    SetLastModifiedBy(value *string)()
    SetLastUpdateDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRedirectIncidentId(value *string)()
    SetResolvingComment(value *string)()
    SetSeverity(value *AlertSeverity)()
    SetStatus(value *IncidentStatus)()
    SetSummary(value *string)()
    SetSystemTags(value []string)()
    SetTenantId(value *string)()
}
