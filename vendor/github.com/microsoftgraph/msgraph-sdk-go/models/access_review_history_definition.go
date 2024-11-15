package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewHistoryDefinition struct {
    Entity
}
// NewAccessReviewHistoryDefinition instantiates a new AccessReviewHistoryDefinition and sets the default values.
func NewAccessReviewHistoryDefinition()(*AccessReviewHistoryDefinition) {
    m := &AccessReviewHistoryDefinition{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessReviewHistoryDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewHistoryDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewHistoryDefinition(), nil
}
// GetCreatedBy gets the createdBy property value. The createdBy property
// returns a UserIdentityable when successful
func (m *AccessReviewHistoryDefinition) GetCreatedBy()(UserIdentityable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserIdentityable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Timestamp when the access review definition was created.
// returns a *Time when successful
func (m *AccessReviewHistoryDefinition) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDecisions gets the decisions property value. Determines which review decisions will be included in the fetched review history data if specified. Optional on create. All decisions are included by default if no decisions are provided on create. Possible values are: approve, deny, dontKnow, notReviewed, and notNotified.
// returns a []AccessReviewHistoryDecisionFilter when successful
func (m *AccessReviewHistoryDefinition) GetDecisions()([]AccessReviewHistoryDecisionFilter) {
    val, err := m.GetBackingStore().Get("decisions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewHistoryDecisionFilter)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name for the access review history data collection. Required.
// returns a *string when successful
func (m *AccessReviewHistoryDefinition) GetDisplayName()(*string) {
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
func (m *AccessReviewHistoryDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(UserIdentityable))
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
    res["decisions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseAccessReviewHistoryDecisionFilter)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewHistoryDecisionFilter, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*AccessReviewHistoryDecisionFilter))
                }
            }
            m.SetDecisions(res)
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
    res["instances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewHistoryInstanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewHistoryInstanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewHistoryInstanceable)
                }
            }
            m.SetInstances(res)
        }
        return nil
    }
    res["reviewHistoryPeriodEndDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReviewHistoryPeriodEndDateTime(val)
        }
        return nil
    }
    res["reviewHistoryPeriodStartDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReviewHistoryPeriodStartDateTime(val)
        }
        return nil
    }
    res["scheduleSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewHistoryScheduleSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduleSettings(val.(AccessReviewHistoryScheduleSettingsable))
        }
        return nil
    }
    res["scopes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewScopeable)
                }
            }
            m.SetScopes(res)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessReviewHistoryStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*AccessReviewHistoryStatus))
        }
        return nil
    }
    return res
}
// GetInstances gets the instances property value. If the accessReviewHistoryDefinition is a recurring definition, instances represent each recurrence. A definition that doesn't recur will have exactly one instance.
// returns a []AccessReviewHistoryInstanceable when successful
func (m *AccessReviewHistoryDefinition) GetInstances()([]AccessReviewHistoryInstanceable) {
    val, err := m.GetBackingStore().Get("instances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewHistoryInstanceable)
    }
    return nil
}
// GetReviewHistoryPeriodEndDateTime gets the reviewHistoryPeriodEndDateTime property value. A timestamp. Reviews ending on or before this date will be included in the fetched history data. Only required if scheduleSettings isn't defined.
// returns a *Time when successful
func (m *AccessReviewHistoryDefinition) GetReviewHistoryPeriodEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("reviewHistoryPeriodEndDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetReviewHistoryPeriodStartDateTime gets the reviewHistoryPeriodStartDateTime property value. A timestamp. Reviews starting on or before this date will be included in the fetched history data. Only required if scheduleSettings isn't defined.
// returns a *Time when successful
func (m *AccessReviewHistoryDefinition) GetReviewHistoryPeriodStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("reviewHistoryPeriodStartDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetScheduleSettings gets the scheduleSettings property value. The settings for a recurring access review history definition series. Only required if reviewHistoryPeriodStartDateTime or reviewHistoryPeriodEndDateTime aren't defined. Not supported yet.
// returns a AccessReviewHistoryScheduleSettingsable when successful
func (m *AccessReviewHistoryDefinition) GetScheduleSettings()(AccessReviewHistoryScheduleSettingsable) {
    val, err := m.GetBackingStore().Get("scheduleSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewHistoryScheduleSettingsable)
    }
    return nil
}
// GetScopes gets the scopes property value. Used to scope what reviews are included in the fetched history data. Fetches reviews whose scope matches with this provided scope. Required.
// returns a []AccessReviewScopeable when successful
func (m *AccessReviewHistoryDefinition) GetScopes()([]AccessReviewScopeable) {
    val, err := m.GetBackingStore().Get("scopes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewScopeable)
    }
    return nil
}
// GetStatus gets the status property value. Represents the status of the review history data collection. The possible values are: done, inProgress, error, requested, unknownFutureValue.
// returns a *AccessReviewHistoryStatus when successful
func (m *AccessReviewHistoryDefinition) GetStatus()(*AccessReviewHistoryStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessReviewHistoryStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewHistoryDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
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
    if m.GetDecisions() != nil {
        err = writer.WriteCollectionOfStringValues("decisions", SerializeAccessReviewHistoryDecisionFilter(m.GetDecisions()))
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
    if m.GetInstances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInstances()))
        for i, v := range m.GetInstances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("instances", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("reviewHistoryPeriodEndDateTime", m.GetReviewHistoryPeriodEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("reviewHistoryPeriodStartDateTime", m.GetReviewHistoryPeriodStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("scheduleSettings", m.GetScheduleSettings())
        if err != nil {
            return err
        }
    }
    if m.GetScopes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScopes()))
        for i, v := range m.GetScopes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scopes", cast)
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
// SetCreatedBy sets the createdBy property value. The createdBy property
func (m *AccessReviewHistoryDefinition) SetCreatedBy(value UserIdentityable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Timestamp when the access review definition was created.
func (m *AccessReviewHistoryDefinition) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDecisions sets the decisions property value. Determines which review decisions will be included in the fetched review history data if specified. Optional on create. All decisions are included by default if no decisions are provided on create. Possible values are: approve, deny, dontKnow, notReviewed, and notNotified.
func (m *AccessReviewHistoryDefinition) SetDecisions(value []AccessReviewHistoryDecisionFilter)() {
    err := m.GetBackingStore().Set("decisions", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name for the access review history data collection. Required.
func (m *AccessReviewHistoryDefinition) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetInstances sets the instances property value. If the accessReviewHistoryDefinition is a recurring definition, instances represent each recurrence. A definition that doesn't recur will have exactly one instance.
func (m *AccessReviewHistoryDefinition) SetInstances(value []AccessReviewHistoryInstanceable)() {
    err := m.GetBackingStore().Set("instances", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewHistoryPeriodEndDateTime sets the reviewHistoryPeriodEndDateTime property value. A timestamp. Reviews ending on or before this date will be included in the fetched history data. Only required if scheduleSettings isn't defined.
func (m *AccessReviewHistoryDefinition) SetReviewHistoryPeriodEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("reviewHistoryPeriodEndDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewHistoryPeriodStartDateTime sets the reviewHistoryPeriodStartDateTime property value. A timestamp. Reviews starting on or before this date will be included in the fetched history data. Only required if scheduleSettings isn't defined.
func (m *AccessReviewHistoryDefinition) SetReviewHistoryPeriodStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("reviewHistoryPeriodStartDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduleSettings sets the scheduleSettings property value. The settings for a recurring access review history definition series. Only required if reviewHistoryPeriodStartDateTime or reviewHistoryPeriodEndDateTime aren't defined. Not supported yet.
func (m *AccessReviewHistoryDefinition) SetScheduleSettings(value AccessReviewHistoryScheduleSettingsable)() {
    err := m.GetBackingStore().Set("scheduleSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetScopes sets the scopes property value. Used to scope what reviews are included in the fetched history data. Fetches reviews whose scope matches with this provided scope. Required.
func (m *AccessReviewHistoryDefinition) SetScopes(value []AccessReviewScopeable)() {
    err := m.GetBackingStore().Set("scopes", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Represents the status of the review history data collection. The possible values are: done, inProgress, error, requested, unknownFutureValue.
func (m *AccessReviewHistoryDefinition) SetStatus(value *AccessReviewHistoryStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewHistoryDefinitionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCreatedBy()(UserIdentityable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDecisions()([]AccessReviewHistoryDecisionFilter)
    GetDisplayName()(*string)
    GetInstances()([]AccessReviewHistoryInstanceable)
    GetReviewHistoryPeriodEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetReviewHistoryPeriodStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetScheduleSettings()(AccessReviewHistoryScheduleSettingsable)
    GetScopes()([]AccessReviewScopeable)
    GetStatus()(*AccessReviewHistoryStatus)
    SetCreatedBy(value UserIdentityable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDecisions(value []AccessReviewHistoryDecisionFilter)()
    SetDisplayName(value *string)()
    SetInstances(value []AccessReviewHistoryInstanceable)()
    SetReviewHistoryPeriodEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetReviewHistoryPeriodStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetScheduleSettings(value AccessReviewHistoryScheduleSettingsable)()
    SetScopes(value []AccessReviewScopeable)()
    SetStatus(value *AccessReviewHistoryStatus)()
}
