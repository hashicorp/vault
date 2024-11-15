package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SimulationNotification struct {
    BaseEndUserNotification
}
// NewSimulationNotification instantiates a new SimulationNotification and sets the default values.
func NewSimulationNotification()(*SimulationNotification) {
    m := &SimulationNotification{
        BaseEndUserNotification: *NewBaseEndUserNotification(),
    }
    odataTypeValue := "#microsoft.graph.simulationNotification"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSimulationNotificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSimulationNotificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSimulationNotification(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SimulationNotification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseEndUserNotification.GetFieldDeserializers()
    res["targettedUserType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTargettedUserType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargettedUserType(val.(*TargettedUserType))
        }
        return nil
    }
    return res
}
// GetTargettedUserType gets the targettedUserType property value. Target user type. Possible values are: unknown, clicked, compromised, allUsers, unknownFutureValue.
// returns a *TargettedUserType when successful
func (m *SimulationNotification) GetTargettedUserType()(*TargettedUserType) {
    val, err := m.GetBackingStore().Get("targettedUserType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TargettedUserType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SimulationNotification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseEndUserNotification.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetTargettedUserType() != nil {
        cast := (*m.GetTargettedUserType()).String()
        err = writer.WriteStringValue("targettedUserType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTargettedUserType sets the targettedUserType property value. Target user type. Possible values are: unknown, clicked, compromised, allUsers, unknownFutureValue.
func (m *SimulationNotification) SetTargettedUserType(value *TargettedUserType)() {
    err := m.GetBackingStore().Set("targettedUserType", value)
    if err != nil {
        panic(err)
    }
}
type SimulationNotificationable interface {
    BaseEndUserNotificationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetTargettedUserType()(*TargettedUserType)
    SetTargettedUserType(value *TargettedUserType)()
}
