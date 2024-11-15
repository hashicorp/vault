package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OpenShiftItem struct {
    ShiftItem
}
// NewOpenShiftItem instantiates a new OpenShiftItem and sets the default values.
func NewOpenShiftItem()(*OpenShiftItem) {
    m := &OpenShiftItem{
        ShiftItem: *NewShiftItem(),
    }
    odataTypeValue := "#microsoft.graph.openShiftItem"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOpenShiftItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOpenShiftItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOpenShiftItem(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OpenShiftItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ShiftItem.GetFieldDeserializers()
    res["openSlotCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOpenSlotCount(val)
        }
        return nil
    }
    return res
}
// GetOpenSlotCount gets the openSlotCount property value. Count of the number of slots for the given open shift.
// returns a *int32 when successful
func (m *OpenShiftItem) GetOpenSlotCount()(*int32) {
    val, err := m.GetBackingStore().Get("openSlotCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OpenShiftItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ShiftItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("openSlotCount", m.GetOpenSlotCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOpenSlotCount sets the openSlotCount property value. Count of the number of slots for the given open shift.
func (m *OpenShiftItem) SetOpenSlotCount(value *int32)() {
    err := m.GetBackingStore().Set("openSlotCount", value)
    if err != nil {
        panic(err)
    }
}
type OpenShiftItemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ShiftItemable
    GetOpenSlotCount()(*int32)
    SetOpenSlotCount(value *int32)()
}
