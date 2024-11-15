package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrivilegedAccessRoot struct {
    Entity
}
// NewPrivilegedAccessRoot instantiates a new PrivilegedAccessRoot and sets the default values.
func NewPrivilegedAccessRoot()(*PrivilegedAccessRoot) {
    m := &PrivilegedAccessRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrivilegedAccessRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrivilegedAccessRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrivilegedAccessRoot(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrivilegedAccessRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["group"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrivilegedAccessGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroup(val.(PrivilegedAccessGroupable))
        }
        return nil
    }
    return res
}
// GetGroup gets the group property value. A group that's governed through Privileged Identity Management (PIM).
// returns a PrivilegedAccessGroupable when successful
func (m *PrivilegedAccessRoot) GetGroup()(PrivilegedAccessGroupable) {
    val, err := m.GetBackingStore().Get("group")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrivilegedAccessGroupable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrivilegedAccessRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("group", m.GetGroup())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetGroup sets the group property value. A group that's governed through Privileged Identity Management (PIM).
func (m *PrivilegedAccessRoot) SetGroup(value PrivilegedAccessGroupable)() {
    err := m.GetBackingStore().Set("group", value)
    if err != nil {
        panic(err)
    }
}
type PrivilegedAccessRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetGroup()(PrivilegedAccessGroupable)
    SetGroup(value PrivilegedAccessGroupable)()
}
