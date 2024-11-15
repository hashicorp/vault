package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type NoTrainingNotificationSetting struct {
    EndUserNotificationSetting
}
// NewNoTrainingNotificationSetting instantiates a new NoTrainingNotificationSetting and sets the default values.
func NewNoTrainingNotificationSetting()(*NoTrainingNotificationSetting) {
    m := &NoTrainingNotificationSetting{
        EndUserNotificationSetting: *NewEndUserNotificationSetting(),
    }
    odataTypeValue := "#microsoft.graph.noTrainingNotificationSetting"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateNoTrainingNotificationSettingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateNoTrainingNotificationSettingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewNoTrainingNotificationSetting(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *NoTrainingNotificationSetting) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EndUserNotificationSetting.GetFieldDeserializers()
    res["simulationNotification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSimulationNotificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSimulationNotification(val.(SimulationNotificationable))
        }
        return nil
    }
    return res
}
// GetSimulationNotification gets the simulationNotification property value. The notification for the user who is part of the simulation.
// returns a SimulationNotificationable when successful
func (m *NoTrainingNotificationSetting) GetSimulationNotification()(SimulationNotificationable) {
    val, err := m.GetBackingStore().Get("simulationNotification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SimulationNotificationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *NoTrainingNotificationSetting) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EndUserNotificationSetting.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("simulationNotification", m.GetSimulationNotification())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSimulationNotification sets the simulationNotification property value. The notification for the user who is part of the simulation.
func (m *NoTrainingNotificationSetting) SetSimulationNotification(value SimulationNotificationable)() {
    err := m.GetBackingStore().Set("simulationNotification", value)
    if err != nil {
        panic(err)
    }
}
type NoTrainingNotificationSettingable interface {
    EndUserNotificationSettingable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetSimulationNotification()(SimulationNotificationable)
    SetSimulationNotification(value SimulationNotificationable)()
}
