package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AzureAdTokenAuthentication struct {
    CustomExtensionAuthenticationConfiguration
}
// NewAzureAdTokenAuthentication instantiates a new AzureAdTokenAuthentication and sets the default values.
func NewAzureAdTokenAuthentication()(*AzureAdTokenAuthentication) {
    m := &AzureAdTokenAuthentication{
        CustomExtensionAuthenticationConfiguration: *NewCustomExtensionAuthenticationConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.azureAdTokenAuthentication"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAzureAdTokenAuthenticationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAzureAdTokenAuthenticationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAzureAdTokenAuthentication(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AzureAdTokenAuthentication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CustomExtensionAuthenticationConfiguration.GetFieldDeserializers()
    res["resourceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceId(val)
        }
        return nil
    }
    return res
}
// GetResourceId gets the resourceId property value. The appID of the Microsoft Entra application to use to authenticate an app with a custom extension.
// returns a *string when successful
func (m *AzureAdTokenAuthentication) GetResourceId()(*string) {
    val, err := m.GetBackingStore().Get("resourceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AzureAdTokenAuthentication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CustomExtensionAuthenticationConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("resourceId", m.GetResourceId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetResourceId sets the resourceId property value. The appID of the Microsoft Entra application to use to authenticate an app with a custom extension.
func (m *AzureAdTokenAuthentication) SetResourceId(value *string)() {
    err := m.GetBackingStore().Set("resourceId", value)
    if err != nil {
        panic(err)
    }
}
type AzureAdTokenAuthenticationable interface {
    CustomExtensionAuthenticationConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetResourceId()(*string)
    SetResourceId(value *string)()
}
