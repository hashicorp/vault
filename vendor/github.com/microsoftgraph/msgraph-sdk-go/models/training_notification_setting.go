package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TrainingNotificationSetting struct {
    EndUserNotificationSetting
}
// NewTrainingNotificationSetting instantiates a new TrainingNotificationSetting and sets the default values.
func NewTrainingNotificationSetting()(*TrainingNotificationSetting) {
    m := &TrainingNotificationSetting{
        EndUserNotificationSetting: *NewEndUserNotificationSetting(),
    }
    odataTypeValue := "#microsoft.graph.trainingNotificationSetting"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTrainingNotificationSettingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTrainingNotificationSettingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTrainingNotificationSetting(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TrainingNotificationSetting) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EndUserNotificationSetting.GetFieldDeserializers()
    res["trainingAssignment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBaseEndUserNotificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrainingAssignment(val.(BaseEndUserNotificationable))
        }
        return nil
    }
    res["trainingReminder"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTrainingReminderNotificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrainingReminder(val.(TrainingReminderNotificationable))
        }
        return nil
    }
    return res
}
// GetTrainingAssignment gets the trainingAssignment property value. Training assignment details.
// returns a BaseEndUserNotificationable when successful
func (m *TrainingNotificationSetting) GetTrainingAssignment()(BaseEndUserNotificationable) {
    val, err := m.GetBackingStore().Get("trainingAssignment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BaseEndUserNotificationable)
    }
    return nil
}
// GetTrainingReminder gets the trainingReminder property value. Training reminder details.
// returns a TrainingReminderNotificationable when successful
func (m *TrainingNotificationSetting) GetTrainingReminder()(TrainingReminderNotificationable) {
    val, err := m.GetBackingStore().Get("trainingReminder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TrainingReminderNotificationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TrainingNotificationSetting) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EndUserNotificationSetting.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("trainingAssignment", m.GetTrainingAssignment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("trainingReminder", m.GetTrainingReminder())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTrainingAssignment sets the trainingAssignment property value. Training assignment details.
func (m *TrainingNotificationSetting) SetTrainingAssignment(value BaseEndUserNotificationable)() {
    err := m.GetBackingStore().Set("trainingAssignment", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainingReminder sets the trainingReminder property value. Training reminder details.
func (m *TrainingNotificationSetting) SetTrainingReminder(value TrainingReminderNotificationable)() {
    err := m.GetBackingStore().Set("trainingReminder", value)
    if err != nil {
        panic(err)
    }
}
type TrainingNotificationSettingable interface {
    EndUserNotificationSettingable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetTrainingAssignment()(BaseEndUserNotificationable)
    GetTrainingReminder()(TrainingReminderNotificationable)
    SetTrainingAssignment(value BaseEndUserNotificationable)()
    SetTrainingReminder(value TrainingReminderNotificationable)()
}
