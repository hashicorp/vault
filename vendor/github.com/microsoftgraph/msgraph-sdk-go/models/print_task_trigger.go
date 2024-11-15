package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrintTaskTrigger struct {
    Entity
}
// NewPrintTaskTrigger instantiates a new PrintTaskTrigger and sets the default values.
func NewPrintTaskTrigger()(*PrintTaskTrigger) {
    m := &PrintTaskTrigger{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrintTaskTriggerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintTaskTriggerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintTaskTrigger(), nil
}
// GetDefinition gets the definition property value. The definition property
// returns a PrintTaskDefinitionable when successful
func (m *PrintTaskTrigger) GetDefinition()(PrintTaskDefinitionable) {
    val, err := m.GetBackingStore().Get("definition")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintTaskDefinitionable)
    }
    return nil
}
// GetEvent gets the event property value. The event property
// returns a *PrintEvent when successful
func (m *PrintTaskTrigger) GetEvent()(*PrintEvent) {
    val, err := m.GetBackingStore().Get("event")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintEvent)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrintTaskTrigger) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["event"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintEvent)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEvent(val.(*PrintEvent))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *PrintTaskTrigger) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    if m.GetEvent() != nil {
        cast := (*m.GetEvent()).String()
        err = writer.WriteStringValue("event", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDefinition sets the definition property value. The definition property
func (m *PrintTaskTrigger) SetDefinition(value PrintTaskDefinitionable)() {
    err := m.GetBackingStore().Set("definition", value)
    if err != nil {
        panic(err)
    }
}
// SetEvent sets the event property value. The event property
func (m *PrintTaskTrigger) SetEvent(value *PrintEvent)() {
    err := m.GetBackingStore().Set("event", value)
    if err != nil {
        panic(err)
    }
}
type PrintTaskTriggerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDefinition()(PrintTaskDefinitionable)
    GetEvent()(*PrintEvent)
    SetDefinition(value PrintTaskDefinitionable)()
    SetEvent(value *PrintEvent)()
}
