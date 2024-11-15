package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AzureCommunicationServicesUserIdentity struct {
    Identity
}
// NewAzureCommunicationServicesUserIdentity instantiates a new AzureCommunicationServicesUserIdentity and sets the default values.
func NewAzureCommunicationServicesUserIdentity()(*AzureCommunicationServicesUserIdentity) {
    m := &AzureCommunicationServicesUserIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.azureCommunicationServicesUserIdentity"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAzureCommunicationServicesUserIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAzureCommunicationServicesUserIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAzureCommunicationServicesUserIdentity(), nil
}
// GetAzureCommunicationServicesResourceId gets the azureCommunicationServicesResourceId property value. The Azure Communication Services resource ID associated with the user.
// returns a *string when successful
func (m *AzureCommunicationServicesUserIdentity) GetAzureCommunicationServicesResourceId()(*string) {
    val, err := m.GetBackingStore().Get("azureCommunicationServicesResourceId")
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
func (m *AzureCommunicationServicesUserIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["azureCommunicationServicesResourceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureCommunicationServicesResourceId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AzureCommunicationServicesUserIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("azureCommunicationServicesResourceId", m.GetAzureCommunicationServicesResourceId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAzureCommunicationServicesResourceId sets the azureCommunicationServicesResourceId property value. The Azure Communication Services resource ID associated with the user.
func (m *AzureCommunicationServicesUserIdentity) SetAzureCommunicationServicesResourceId(value *string)() {
    err := m.GetBackingStore().Set("azureCommunicationServicesResourceId", value)
    if err != nil {
        panic(err)
    }
}
type AzureCommunicationServicesUserIdentityable interface {
    Identityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAzureCommunicationServicesResourceId()(*string)
    SetAzureCommunicationServicesResourceId(value *string)()
}
