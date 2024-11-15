package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EnumeratedDeviceRegistrationMembership struct {
    DeviceRegistrationMembership
}
// NewEnumeratedDeviceRegistrationMembership instantiates a new EnumeratedDeviceRegistrationMembership and sets the default values.
func NewEnumeratedDeviceRegistrationMembership()(*EnumeratedDeviceRegistrationMembership) {
    m := &EnumeratedDeviceRegistrationMembership{
        DeviceRegistrationMembership: *NewDeviceRegistrationMembership(),
    }
    odataTypeValue := "#microsoft.graph.enumeratedDeviceRegistrationMembership"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEnumeratedDeviceRegistrationMembershipFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEnumeratedDeviceRegistrationMembershipFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEnumeratedDeviceRegistrationMembership(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EnumeratedDeviceRegistrationMembership) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceRegistrationMembership.GetFieldDeserializers()
    res["groups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetGroups(res)
        }
        return nil
    }
    res["users"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetUsers(res)
        }
        return nil
    }
    return res
}
// GetGroups gets the groups property value. The groups property
// returns a []string when successful
func (m *EnumeratedDeviceRegistrationMembership) GetGroups()([]string) {
    val, err := m.GetBackingStore().Get("groups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetUsers gets the users property value. The users property
// returns a []string when successful
func (m *EnumeratedDeviceRegistrationMembership) GetUsers()([]string) {
    val, err := m.GetBackingStore().Get("users")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EnumeratedDeviceRegistrationMembership) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceRegistrationMembership.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetGroups() != nil {
        err = writer.WriteCollectionOfStringValues("groups", m.GetGroups())
        if err != nil {
            return err
        }
    }
    if m.GetUsers() != nil {
        err = writer.WriteCollectionOfStringValues("users", m.GetUsers())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetGroups sets the groups property value. The groups property
func (m *EnumeratedDeviceRegistrationMembership) SetGroups(value []string)() {
    err := m.GetBackingStore().Set("groups", value)
    if err != nil {
        panic(err)
    }
}
// SetUsers sets the users property value. The users property
func (m *EnumeratedDeviceRegistrationMembership) SetUsers(value []string)() {
    err := m.GetBackingStore().Set("users", value)
    if err != nil {
        panic(err)
    }
}
type EnumeratedDeviceRegistrationMembershipable interface {
    DeviceRegistrationMembershipable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetGroups()([]string)
    GetUsers()([]string)
    SetGroups(value []string)()
    SetUsers(value []string)()
}
