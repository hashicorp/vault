package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SharePointIdentity struct {
    Identity
}
// NewSharePointIdentity instantiates a new SharePointIdentity and sets the default values.
func NewSharePointIdentity()(*SharePointIdentity) {
    m := &SharePointIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.sharePointIdentity"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSharePointIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharePointIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharePointIdentity(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharePointIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["loginName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLoginName(val)
        }
        return nil
    }
    return res
}
// GetLoginName gets the loginName property value. The sign in name of the SharePoint identity.
// returns a *string when successful
func (m *SharePointIdentity) GetLoginName()(*string) {
    val, err := m.GetBackingStore().Get("loginName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharePointIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("loginName", m.GetLoginName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLoginName sets the loginName property value. The sign in name of the SharePoint identity.
func (m *SharePointIdentity) SetLoginName(value *string)() {
    err := m.GetBackingStore().Set("loginName", value)
    if err != nil {
        panic(err)
    }
}
type SharePointIdentityable interface {
    Identityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLoginName()(*string)
    SetLoginName(value *string)()
}
