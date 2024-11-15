package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrintTask struct {
    Entity
}
// NewPrintTask instantiates a new PrintTask and sets the default values.
func NewPrintTask()(*PrintTask) {
    m := &PrintTask{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrintTaskFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintTaskFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintTask(), nil
}
// GetDefinition gets the definition property value. The definition property
// returns a PrintTaskDefinitionable when successful
func (m *PrintTask) GetDefinition()(PrintTaskDefinitionable) {
    val, err := m.GetBackingStore().Get("definition")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintTaskDefinitionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrintTask) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["definition"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintTaskDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefinition(val.(PrintTaskDefinitionable))
        }
        return nil
    }
    res["parentUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentUrl(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintTaskStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(PrintTaskStatusable))
        }
        return nil
    }
    res["trigger"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintTaskTriggerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrigger(val.(PrintTaskTriggerable))
        }
        return nil
    }
    return res
}
// GetParentUrl gets the parentUrl property value. The URL for the print entity that triggered this task. For example, https://graph.microsoft.com/v1.0/print/printers/{printerId}/jobs/{jobId}. Read-only.
// returns a *string when successful
func (m *PrintTask) GetParentUrl()(*string) {
    val, err := m.GetBackingStore().Get("parentUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a PrintTaskStatusable when successful
func (m *PrintTask) GetStatus()(PrintTaskStatusable) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintTaskStatusable)
    }
    return nil
}
// GetTrigger gets the trigger property value. The trigger property
// returns a PrintTaskTriggerable when successful
func (m *PrintTask) GetTrigger()(PrintTaskTriggerable) {
    val, err := m.GetBackingStore().Get("trigger")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintTaskTriggerable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrintTask) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("definition", m.GetDefinition())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentUrl", m.GetParentUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("trigger", m.GetTrigger())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDefinition sets the definition property value. The definition property
func (m *PrintTask) SetDefinition(value PrintTaskDefinitionable)() {
    err := m.GetBackingStore().Set("definition", value)
    if err != nil {
        panic(err)
    }
}
// SetParentUrl sets the parentUrl property value. The URL for the print entity that triggered this task. For example, https://graph.microsoft.com/v1.0/print/printers/{printerId}/jobs/{jobId}. Read-only.
func (m *PrintTask) SetParentUrl(value *string)() {
    err := m.GetBackingStore().Set("parentUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *PrintTask) SetStatus(value PrintTaskStatusable)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTrigger sets the trigger property value. The trigger property
func (m *PrintTask) SetTrigger(value PrintTaskTriggerable)() {
    err := m.GetBackingStore().Set("trigger", value)
    if err != nil {
        panic(err)
    }
}
type PrintTaskable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDefinition()(PrintTaskDefinitionable)
    GetParentUrl()(*string)
    GetStatus()(PrintTaskStatusable)
    GetTrigger()(PrintTaskTriggerable)
    SetDefinition(value PrintTaskDefinitionable)()
    SetParentUrl(value *string)()
    SetStatus(value PrintTaskStatusable)()
    SetTrigger(value PrintTaskTriggerable)()
}
