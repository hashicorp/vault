package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CloudPcAuditEvent struct {
    Entity
}
// NewCloudPcAuditEvent instantiates a new CloudPcAuditEvent and sets the default values.
func NewCloudPcAuditEvent()(*CloudPcAuditEvent) {
    m := &CloudPcAuditEvent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCloudPcAuditEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCloudPcAuditEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudPcAuditEvent(), nil
}
// GetActivity gets the activity property value. The friendly name of the audit activity.
// returns a *string when successful
func (m *CloudPcAuditEvent) GetActivity()(*string) {
    val, err := m.GetBackingStore().Get("activity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActivityDateTime gets the activityDateTime property value. The date time in UTC when the activity was performed. Read-only.
// returns a *Time when successful
func (m *CloudPcAuditEvent) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("activityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetActivityOperationType gets the activityOperationType property value. The activityOperationType property
// returns a *CloudPcAuditActivityOperationType when successful
func (m *CloudPcAuditEvent) GetActivityOperationType()(*CloudPcAuditActivityOperationType) {
    val, err := m.GetBackingStore().Get("activityOperationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CloudPcAuditActivityOperationType)
    }
    return nil
}
// GetActivityResult gets the activityResult property value. The activityResult property
// returns a *CloudPcAuditActivityResult when successful
func (m *CloudPcAuditEvent) GetActivityResult()(*CloudPcAuditActivityResult) {
    val, err := m.GetBackingStore().Get("activityResult")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CloudPcAuditActivityResult)
    }
    return nil
}
// GetActivityType gets the activityType property value. The type of activity that was performed. Read-only.
// returns a *string when successful
func (m *CloudPcAuditEvent) GetActivityType()(*string) {
    val, err := m.GetBackingStore().Get("activityType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActor gets the actor property value. The actor property
// returns a CloudPcAuditActorable when successful
func (m *CloudPcAuditEvent) GetActor()(CloudPcAuditActorable) {
    val, err := m.GetBackingStore().Get("actor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CloudPcAuditActorable)
    }
    return nil
}
// GetCategory gets the category property value. The category property
// returns a *CloudPcAuditCategory when successful
func (m *CloudPcAuditEvent) GetCategory()(*CloudPcAuditCategory) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CloudPcAuditCategory)
    }
    return nil
}
// GetComponentName gets the componentName property value. The component name for the audit event. Read-only.
// returns a *string when successful
func (m *CloudPcAuditEvent) GetComponentName()(*string) {
    val, err := m.GetBackingStore().Get("componentName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCorrelationId gets the correlationId property value. The client request ID that is used to correlate activity within the system. Read-only.
// returns a *string when successful
func (m *CloudPcAuditEvent) GetCorrelationId()(*string) {
    val, err := m.GetBackingStore().Get("correlationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for the audit event. Read-only.
// returns a *string when successful
func (m *CloudPcAuditEvent) GetDisplayName()(*string) {
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
func (m *CloudPcAuditEvent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
        val, err := n.GetEnumValue(ParseCloudPcAuditActivityOperationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityOperationType(val.(*CloudPcAuditActivityOperationType))
        }
        return nil
    }
    res["activityResult"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCloudPcAuditActivityResult)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityResult(val.(*CloudPcAuditActivityResult))
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
        val, err := n.GetObjectValue(CreateCloudPcAuditActorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActor(val.(CloudPcAuditActorable))
        }
        return nil
    }
    res["category"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCloudPcAuditCategory)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory(val.(*CloudPcAuditCategory))
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
        val, err := n.GetStringValue()
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
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcAuditResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcAuditResourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcAuditResourceable)
                }
            }
            m.SetResources(res)
        }
        return nil
    }
    return res
}
// GetResources gets the resources property value. The list of cloudPcAuditResource objects. Read-only.
// returns a []CloudPcAuditResourceable when successful
func (m *CloudPcAuditEvent) GetResources()([]CloudPcAuditResourceable) {
    val, err := m.GetBackingStore().Get("resources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcAuditResourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CloudPcAuditEvent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    if m.GetActivityOperationType() != nil {
        cast := (*m.GetActivityOperationType()).String()
        err = writer.WriteStringValue("activityOperationType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetActivityResult() != nil {
        cast := (*m.GetActivityResult()).String()
        err = writer.WriteStringValue("activityResult", &cast)
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
    if m.GetCategory() != nil {
        cast := (*m.GetCategory()).String()
        err = writer.WriteStringValue("category", &cast)
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
        err = writer.WriteStringValue("correlationId", m.GetCorrelationId())
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
// SetActivity sets the activity property value. The friendly name of the audit activity.
func (m *CloudPcAuditEvent) SetActivity(value *string)() {
    err := m.GetBackingStore().Set("activity", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityDateTime sets the activityDateTime property value. The date time in UTC when the activity was performed. Read-only.
func (m *CloudPcAuditEvent) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("activityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityOperationType sets the activityOperationType property value. The activityOperationType property
func (m *CloudPcAuditEvent) SetActivityOperationType(value *CloudPcAuditActivityOperationType)() {
    err := m.GetBackingStore().Set("activityOperationType", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityResult sets the activityResult property value. The activityResult property
func (m *CloudPcAuditEvent) SetActivityResult(value *CloudPcAuditActivityResult)() {
    err := m.GetBackingStore().Set("activityResult", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityType sets the activityType property value. The type of activity that was performed. Read-only.
func (m *CloudPcAuditEvent) SetActivityType(value *string)() {
    err := m.GetBackingStore().Set("activityType", value)
    if err != nil {
        panic(err)
    }
}
// SetActor sets the actor property value. The actor property
func (m *CloudPcAuditEvent) SetActor(value CloudPcAuditActorable)() {
    err := m.GetBackingStore().Set("actor", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory sets the category property value. The category property
func (m *CloudPcAuditEvent) SetCategory(value *CloudPcAuditCategory)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetComponentName sets the componentName property value. The component name for the audit event. Read-only.
func (m *CloudPcAuditEvent) SetComponentName(value *string)() {
    err := m.GetBackingStore().Set("componentName", value)
    if err != nil {
        panic(err)
    }
}
// SetCorrelationId sets the correlationId property value. The client request ID that is used to correlate activity within the system. Read-only.
func (m *CloudPcAuditEvent) SetCorrelationId(value *string)() {
    err := m.GetBackingStore().Set("correlationId", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for the audit event. Read-only.
func (m *CloudPcAuditEvent) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetResources sets the resources property value. The list of cloudPcAuditResource objects. Read-only.
func (m *CloudPcAuditEvent) SetResources(value []CloudPcAuditResourceable)() {
    err := m.GetBackingStore().Set("resources", value)
    if err != nil {
        panic(err)
    }
}
type CloudPcAuditEventable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivity()(*string)
    GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetActivityOperationType()(*CloudPcAuditActivityOperationType)
    GetActivityResult()(*CloudPcAuditActivityResult)
    GetActivityType()(*string)
    GetActor()(CloudPcAuditActorable)
    GetCategory()(*CloudPcAuditCategory)
    GetComponentName()(*string)
    GetCorrelationId()(*string)
    GetDisplayName()(*string)
    GetResources()([]CloudPcAuditResourceable)
    SetActivity(value *string)()
    SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetActivityOperationType(value *CloudPcAuditActivityOperationType)()
    SetActivityResult(value *CloudPcAuditActivityResult)()
    SetActivityType(value *string)()
    SetActor(value CloudPcAuditActorable)()
    SetCategory(value *CloudPcAuditCategory)()
    SetComponentName(value *string)()
    SetCorrelationId(value *string)()
    SetDisplayName(value *string)()
    SetResources(value []CloudPcAuditResourceable)()
}
