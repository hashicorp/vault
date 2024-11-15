package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CommunicationsApplicationIdentity struct {
    Identity
}
// NewCommunicationsApplicationIdentity instantiates a new CommunicationsApplicationIdentity and sets the default values.
func NewCommunicationsApplicationIdentity()(*CommunicationsApplicationIdentity) {
    m := &CommunicationsApplicationIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.communicationsApplicationIdentity"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCommunicationsApplicationIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCommunicationsApplicationIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCommunicationsApplicationIdentity(), nil
}
// GetApplicationType gets the applicationType property value. First-party Microsoft application that presents this identity.
// returns a *string when successful
func (m *CommunicationsApplicationIdentity) GetApplicationType()(*string) {
    val, err := m.GetBackingStore().Get("applicationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CommunicationsApplicationIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["applicationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationType(val)
        }
        return nil
    }
    res["hidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHidden(val)
        }
        return nil
    }
    return res
}
// GetHidden gets the hidden property value. True if the participant shouldn't be shown in other participants' rosters.
// returns a *bool when successful
func (m *CommunicationsApplicationIdentity) GetHidden()(*bool) {
    val, err := m.GetBackingStore().Get("hidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CommunicationsApplicationIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("applicationType", m.GetApplicationType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hidden", m.GetHidden())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationType sets the applicationType property value. First-party Microsoft application that presents this identity.
func (m *CommunicationsApplicationIdentity) SetApplicationType(value *string)() {
    err := m.GetBackingStore().Set("applicationType", value)
    if err != nil {
        panic(err)
    }
}
// SetHidden sets the hidden property value. True if the participant shouldn't be shown in other participants' rosters.
func (m *CommunicationsApplicationIdentity) SetHidden(value *bool)() {
    err := m.GetBackingStore().Set("hidden", value)
    if err != nil {
        panic(err)
    }
}
type CommunicationsApplicationIdentityable interface {
    Identityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationType()(*string)
    GetHidden()(*bool)
    SetApplicationType(value *string)()
    SetHidden(value *bool)()
}
