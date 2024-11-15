package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ProvisionedIdentity struct {
    Identity
}
// NewProvisionedIdentity instantiates a new ProvisionedIdentity and sets the default values.
func NewProvisionedIdentity()(*ProvisionedIdentity) {
    m := &ProvisionedIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.provisionedIdentity"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateProvisionedIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateProvisionedIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewProvisionedIdentity(), nil
}
// GetDetails gets the details property value. Details of the identity.
// returns a DetailsInfoable when successful
func (m *ProvisionedIdentity) GetDetails()(DetailsInfoable) {
    val, err := m.GetBackingStore().Get("details")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DetailsInfoable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ProvisionedIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["details"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDetailsInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetails(val.(DetailsInfoable))
        }
        return nil
    }
    res["identityType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentityType(val)
        }
        return nil
    }
    return res
}
// GetIdentityType gets the identityType property value. Type of identity that has been provisioned, such as 'user' or 'group'. Supports $filter (eq, contains).
// returns a *string when successful
func (m *ProvisionedIdentity) GetIdentityType()(*string) {
    val, err := m.GetBackingStore().Get("identityType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ProvisionedIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("details", m.GetDetails())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityType", m.GetIdentityType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDetails sets the details property value. Details of the identity.
func (m *ProvisionedIdentity) SetDetails(value DetailsInfoable)() {
    err := m.GetBackingStore().Set("details", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityType sets the identityType property value. Type of identity that has been provisioned, such as 'user' or 'group'. Supports $filter (eq, contains).
func (m *ProvisionedIdentity) SetIdentityType(value *string)() {
    err := m.GetBackingStore().Set("identityType", value)
    if err != nil {
        panic(err)
    }
}
type ProvisionedIdentityable interface {
    Identityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDetails()(DetailsInfoable)
    GetIdentityType()(*string)
    SetDetails(value DetailsInfoable)()
    SetIdentityType(value *string)()
}
