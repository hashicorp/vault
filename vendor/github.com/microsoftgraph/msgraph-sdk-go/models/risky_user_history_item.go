package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RiskyUserHistoryItem struct {
    RiskyUser
}
// NewRiskyUserHistoryItem instantiates a new RiskyUserHistoryItem and sets the default values.
func NewRiskyUserHistoryItem()(*RiskyUserHistoryItem) {
    m := &RiskyUserHistoryItem{
        RiskyUser: *NewRiskyUser(),
    }
    return m
}
// CreateRiskyUserHistoryItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRiskyUserHistoryItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRiskyUserHistoryItem(), nil
}
// GetActivity gets the activity property value. The activity related to user risk level change.
// returns a RiskUserActivityable when successful
func (m *RiskyUserHistoryItem) GetActivity()(RiskUserActivityable) {
    val, err := m.GetBackingStore().Get("activity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RiskUserActivityable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RiskyUserHistoryItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RiskyUser.GetFieldDeserializers()
    res["activity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRiskUserActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivity(val.(RiskUserActivityable))
        }
        return nil
    }
    res["initiatedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitiatedBy(val)
        }
        return nil
    }
    res["userId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserId(val)
        }
        return nil
    }
    return res
}
// GetInitiatedBy gets the initiatedBy property value. The ID of actor that does the operation.
// returns a *string when successful
func (m *RiskyUserHistoryItem) GetInitiatedBy()(*string) {
    val, err := m.GetBackingStore().Get("initiatedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. The ID of the user.
// returns a *string when successful
func (m *RiskyUserHistoryItem) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RiskyUserHistoryItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RiskyUser.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("activity", m.GetActivity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("initiatedBy", m.GetInitiatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivity sets the activity property value. The activity related to user risk level change.
func (m *RiskyUserHistoryItem) SetActivity(value RiskUserActivityable)() {
    err := m.GetBackingStore().Set("activity", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiatedBy sets the initiatedBy property value. The ID of actor that does the operation.
func (m *RiskyUserHistoryItem) SetInitiatedBy(value *string)() {
    err := m.GetBackingStore().Set("initiatedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. The ID of the user.
func (m *RiskyUserHistoryItem) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
type RiskyUserHistoryItemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RiskyUserable
    GetActivity()(RiskUserActivityable)
    GetInitiatedBy()(*string)
    GetUserId()(*string)
    SetActivity(value RiskUserActivityable)()
    SetInitiatedBy(value *string)()
    SetUserId(value *string)()
}
