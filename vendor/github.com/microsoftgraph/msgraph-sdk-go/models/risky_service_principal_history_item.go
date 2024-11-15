package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RiskyServicePrincipalHistoryItem struct {
    RiskyServicePrincipal
}
// NewRiskyServicePrincipalHistoryItem instantiates a new RiskyServicePrincipalHistoryItem and sets the default values.
func NewRiskyServicePrincipalHistoryItem()(*RiskyServicePrincipalHistoryItem) {
    m := &RiskyServicePrincipalHistoryItem{
        RiskyServicePrincipal: *NewRiskyServicePrincipal(),
    }
    return m
}
// CreateRiskyServicePrincipalHistoryItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRiskyServicePrincipalHistoryItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRiskyServicePrincipalHistoryItem(), nil
}
// GetActivity gets the activity property value. The activity related to service principal risk level change.
// returns a RiskServicePrincipalActivityable when successful
func (m *RiskyServicePrincipalHistoryItem) GetActivity()(RiskServicePrincipalActivityable) {
    val, err := m.GetBackingStore().Get("activity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RiskServicePrincipalActivityable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RiskyServicePrincipalHistoryItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RiskyServicePrincipal.GetFieldDeserializers()
    res["activity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRiskServicePrincipalActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivity(val.(RiskServicePrincipalActivityable))
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
    return res
}
// GetInitiatedBy gets the initiatedBy property value. The identifier of the actor of the operation.
// returns a *string when successful
func (m *RiskyServicePrincipalHistoryItem) GetInitiatedBy()(*string) {
    val, err := m.GetBackingStore().Get("initiatedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RiskyServicePrincipalHistoryItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RiskyServicePrincipal.Serialize(writer)
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
    return nil
}
// SetActivity sets the activity property value. The activity related to service principal risk level change.
func (m *RiskyServicePrincipalHistoryItem) SetActivity(value RiskServicePrincipalActivityable)() {
    err := m.GetBackingStore().Set("activity", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiatedBy sets the initiatedBy property value. The identifier of the actor of the operation.
func (m *RiskyServicePrincipalHistoryItem) SetInitiatedBy(value *string)() {
    err := m.GetBackingStore().Set("initiatedBy", value)
    if err != nil {
        panic(err)
    }
}
type RiskyServicePrincipalHistoryItemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RiskyServicePrincipalable
    GetActivity()(RiskServicePrincipalActivityable)
    GetInitiatedBy()(*string)
    SetActivity(value RiskServicePrincipalActivityable)()
    SetInitiatedBy(value *string)()
}
