package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Bitlocker struct {
    Entity
}
// NewBitlocker instantiates a new Bitlocker and sets the default values.
func NewBitlocker()(*Bitlocker) {
    m := &Bitlocker{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBitlockerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBitlockerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBitlocker(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Bitlocker) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["recoveryKeys"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBitlockerRecoveryKeyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BitlockerRecoveryKeyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BitlockerRecoveryKeyable)
                }
            }
            m.SetRecoveryKeys(res)
        }
        return nil
    }
    return res
}
// GetRecoveryKeys gets the recoveryKeys property value. The recovery keys associated with the bitlocker entity.
// returns a []BitlockerRecoveryKeyable when successful
func (m *Bitlocker) GetRecoveryKeys()([]BitlockerRecoveryKeyable) {
    val, err := m.GetBackingStore().Get("recoveryKeys")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BitlockerRecoveryKeyable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Bitlocker) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetRecoveryKeys() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRecoveryKeys()))
        for i, v := range m.GetRecoveryKeys() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("recoveryKeys", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRecoveryKeys sets the recoveryKeys property value. The recovery keys associated with the bitlocker entity.
func (m *Bitlocker) SetRecoveryKeys(value []BitlockerRecoveryKeyable)() {
    err := m.GetBackingStore().Set("recoveryKeys", value)
    if err != nil {
        panic(err)
    }
}
type Bitlockerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRecoveryKeys()([]BitlockerRecoveryKeyable)
    SetRecoveryKeys(value []BitlockerRecoveryKeyable)()
}
