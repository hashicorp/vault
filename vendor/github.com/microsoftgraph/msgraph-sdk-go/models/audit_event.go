package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuditEvent a class containing the properties for Audit Event.
type AuditEvent struct {
    Entity
}
// NewAuditEvent instantiates a new AuditEvent and sets the default values.
func NewAuditEvent()(*AuditEvent) {
    m := &AuditEvent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuditEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuditEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuditEvent(), nil
}
// GetActivity gets the activity property value. Friendly name of the activity.
// returns a *string when successful
func (m *AuditEvent) GetActivity()(*string) {
    val, err := m.GetBackingStore().Get("activity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActivityDateTime gets the activityDateTime property value. The date time in UTC when the activity was performed.
// returns a *Time when successful
func (m *AuditEvent) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("activityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetActivityOperationType gets the activityOperationType property value. The HTTP operation type of the activity.
// returns a *string when successful
func (m *AuditEvent) GetActivityOperationType()(*string) {
    val, err := m.GetBackingStore().Get("activityOperationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActivityResult gets the activityResult property value. The result of the activity.
// returns a *string when successful
func (m *AuditEvent) GetActivityResult()(*string) {
    val, err := m.GetBackingStore().Get("activityResult")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActivityType gets the activityType property value. The type of activity that was being performed.
// returns a *string when successful
func (m *AuditEvent) GetActivityType()(*string) {
    val, err := m.GetBackingStore().Get("activityType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActor gets the actor property value. AAD user and application that are associated with the audit event.
// returns a AuditActorable when successful
func (m *AuditEvent) GetActor()(AuditActorable) {
    val, err := m.GetBackingStore().Get("actor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuditActorable)
    }
    return nil
}
// GetCategory gets the category property value. Audit category.
// returns a *string when successful
func (m *AuditEvent) GetCategory()(*string) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetComponentName gets the componentName property value. Component name.
// returns a *string when successful
func (m *AuditEvent) GetComponentName()(*string) {
    val, err := m.GetBackingStore().Get("componentName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCorrelationId gets the correlationId property value. The client request Id that is used to correlate activity within the system.
// returns a *UUID when successful
func (m *AuditEvent) GetCorrelationId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("correlationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Event display name.
// returns a *string when successful
func (m *AuditEvent) GetDisplayName()(*string) {
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
func (m *AuditEvent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivity(val)
        }
        return nil
    }
    res["activityDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityDateTime(val)
        }
        return nil
    }
    res["activityOperationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityOperationType(val)
        }
        return nil
    }
    res["activityResult"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityResult(val)
        }
        return nil
    }
    res["activityType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityType(val)
        }
        return nil
    }
    res["actor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuditActorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActor(val.(AuditActorable))
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
    res["componentName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComponentName(val)
        }
        return nil
    }
    res["correlationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCorrelationId(val)
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
    res["resources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuditResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuditResourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuditResourceable)
                }
            }
            m.SetResources(res)
        }
        return nil
    }
    return res
}
// GetResources gets the resources property value. Resources being modified.
// returns a []AuditResourceable when successful
func (m *AuditEvent) GetResources()([]AuditResourceable) {
    val, err := m.GetBackingStore().Get("resources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuditResourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuditEvent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("activity", m.GetActivity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("activityDateTime", m.GetActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityOperationType", m.GetActivityOperationType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityResult", m.GetActivityResult())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityType", m.GetActivityType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("actor", m.GetActor())
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
    {
        err = writer.WriteStringValue("componentName", m.GetComponentName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("correlationId", m.GetCorrelationId())
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
    if m.GetResources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResources()))
        for i, v := range m.GetResources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivity sets the activity property value. Friendly name of the activity.
func (m *AuditEvent) SetActivity(value *string)() {
    err := m.GetBackingStore().Set("activity", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityDateTime sets the activityDateTime property value. The date time in UTC when the activity was performed.
func (m *AuditEvent) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("activityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityOperationType sets the activityOperationType property value. The HTTP operation type of the activity.
func (m *AuditEvent) SetActivityOperationType(value *string)() {
    err := m.GetBackingStore().Set("activityOperationType", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityResult sets the activityResult property value. The result of the activity.
func (m *AuditEvent) SetActivityResult(value *string)() {
    err := m.GetBackingStore().Set("activityResult", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityType sets the activityType property value. The type of activity that was being performed.
func (m *AuditEvent) SetActivityType(value *string)() {
    err := m.GetBackingStore().Set("activityType", value)
    if err != nil {
        panic(err)
    }
}
// SetActor sets the actor property value. AAD user and application that are associated with the audit event.
func (m *AuditEvent) SetActor(value AuditActorable)() {
    err := m.GetBackingStore().Set("actor", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory sets the category property value. Audit category.
func (m *AuditEvent) SetCategory(value *string)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetComponentName sets the componentName property value. Component name.
func (m *AuditEvent) SetComponentName(value *string)() {
    err := m.GetBackingStore().Set("componentName", value)
    if err != nil {
        panic(err)
    }
}
// SetCorrelationId sets the correlationId property value. The client request Id that is used to correlate activity within the system.
func (m *AuditEvent) SetCorrelationId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("correlationId", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Event display name.
func (m *AuditEvent) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetResources sets the resources property value. Resources being modified.
func (m *AuditEvent) SetResources(value []AuditResourceable)() {
    err := m.GetBackingStore().Set("resources", value)
    if err != nil {
        panic(err)
    }
}
type AuditEventable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivity()(*string)
    GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetActivityOperationType()(*string)
    GetActivityResult()(*string)
    GetActivityType()(*string)
    GetActor()(AuditActorable)
    GetCategory()(*string)
    GetComponentName()(*string)
    GetCorrelationId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetDisplayName()(*string)
    GetResources()([]AuditResourceable)
    SetActivity(value *string)()
    SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetActivityOperationType(value *string)()
    SetActivityResult(value *string)()
    SetActivityType(value *string)()
    SetActor(value AuditActorable)()
    SetCategory(value *string)()
    SetComponentName(value *string)()
    SetCorrelationId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetDisplayName(value *string)()
    SetResources(value []AuditResourceable)()
}
