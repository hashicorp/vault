package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IdentityApiConnector struct {
    Entity
}
// NewIdentityApiConnector instantiates a new IdentityApiConnector and sets the default values.
func NewIdentityApiConnector()(*IdentityApiConnector) {
    m := &IdentityApiConnector{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIdentityApiConnectorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIdentityApiConnectorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityApiConnector(), nil
}
// GetAuthenticationConfiguration gets the authenticationConfiguration property value. The object which describes the authentication configuration details for calling the API. Basic and PKCS 12 client certificate are supported.
// returns a ApiAuthenticationConfigurationBaseable when successful
func (m *IdentityApiConnector) GetAuthenticationConfiguration()(ApiAuthenticationConfigurationBaseable) {
    val, err := m.GetBackingStore().Get("authenticationConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ApiAuthenticationConfigurationBaseable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the API connector.
// returns a *string when successful
func (m *IdentityApiConnector) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *IdentityApiConnector) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateApiAuthenticationConfigurationBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationConfiguration(val.(ApiAuthenticationConfigurationBaseable))
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["targetUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetUrl(val)
        }
        return nil
    }
    return res
}
// GetTargetUrl gets the targetUrl property value. The URL of the API endpoint to call.
// returns a *string when successful
func (m *IdentityApiConnector) GetTargetUrl()(*string) {
    val, err := m.GetBackingStore().Get("targetUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IdentityApiConnector) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("authenticationConfiguration", m.GetAuthenticationConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("targetUrl", m.GetTargetUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationConfiguration sets the authenticationConfiguration property value. The object which describes the authentication configuration details for calling the API. Basic and PKCS 12 client certificate are supported.
func (m *IdentityApiConnector) SetAuthenticationConfiguration(value ApiAuthenticationConfigurationBaseable)() {
    err := m.GetBackingStore().Set("authenticationConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the API connector.
func (m *IdentityApiConnector) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetUrl sets the targetUrl property value. The URL of the API endpoint to call.
func (m *IdentityApiConnector) SetTargetUrl(value *string)() {
    err := m.GetBackingStore().Set("targetUrl", value)
    if err != nil {
        panic(err)
    }
}
type IdentityApiConnectorable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationConfiguration()(ApiAuthenticationConfigurationBaseable)
    GetDisplayName()(*string)
    GetTargetUrl()(*string)
    SetAuthenticationConfiguration(value ApiAuthenticationConfigurationBaseable)()
    SetDisplayName(value *string)()
    SetTargetUrl(value *string)()
}
