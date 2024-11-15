package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TrainingReminderNotification struct {
    BaseEndUserNotification
}
// NewTrainingReminderNotification instantiates a new TrainingReminderNotification and sets the default values.
func NewTrainingReminderNotification()(*TrainingReminderNotification) {
    m := &TrainingReminderNotification{
        BaseEndUserNotification: *NewBaseEndUserNotification(),
    }
    odataTypeValue := "#microsoft.graph.trainingReminderNotification"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTrainingReminderNotificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTrainingReminderNotificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTrainingReminderNotification(), nil
}
// GetDeliveryFrequency gets the deliveryFrequency property value. Configurable frequency for the reminder email introduced during simulation creation. Possible values are: unknown, weekly, biWeekly, unknownFutureValue.
// returns a *NotificationDeliveryFrequency when successful
func (m *TrainingReminderNotification) GetDeliveryFrequency()(*NotificationDeliveryFrequency) {
    val, err := m.GetBackingStore().Get("deliveryFrequency")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*NotificationDeliveryFrequency)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TrainingReminderNotification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseEndUserNotification.GetFieldDeserializers()
    res["deliveryFrequency"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseNotificationDeliveryFrequency)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryFrequency(val.(*NotificationDeliveryFrequency))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *TrainingReminderNotification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseEndUserNotification.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDeliveryFrequency() != nil {
        cast := (*m.GetDeliveryFrequency()).String()
        err = writer.WriteStringValue("deliveryFrequency", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeliveryFrequency sets the deliveryFrequency property value. Configurable frequency for the reminder email introduced during simulation creation. Possible values are: unknown, weekly, biWeekly, unknownFutureValue.
func (m *TrainingReminderNotification) SetDeliveryFrequency(value *NotificationDeliveryFrequency)() {
    err := m.GetBackingStore().Set("deliveryFrequency", value)
    if err != nil {
        panic(err)
    }
}
type TrainingReminderNotificationable interface {
    BaseEndUserNotificationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeliveryFrequency()(*NotificationDeliveryFrequency)
    SetDeliveryFrequency(value *NotificationDeliveryFrequency)()
}
