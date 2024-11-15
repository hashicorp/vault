package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SocialIdentityProvider struct {
    IdentityProviderBase
}
// NewSocialIdentityProvider instantiates a new SocialIdentityProvider and sets the default values.
func NewSocialIdentityProvider()(*SocialIdentityProvider) {
    m := &SocialIdentityProvider{
        IdentityProviderBase: *NewIdentityProviderBase(),
    }
    odataTypeValue := "#microsoft.graph.socialIdentityProvider"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSocialIdentityProviderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSocialIdentityProviderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSocialIdentityProvider(), nil
}
// GetClientId gets the clientId property value. The identifier for the client application obtained when registering the application with the identity provider. Required.
// returns a *string when successful
func (m *SocialIdentityProvider) GetClientId()(*string) {
    val, err := m.GetBackingStore().Get("clientId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetClientSecret gets the clientSecret property value. The client secret for the application that is obtained when the application is registered with the identity provider. This is write-only. A read operation returns . Required.
// returns a *string when successful
func (m *SocialIdentityProvider) GetClientSecret()(*string) {
    val, err := m.GetBackingStore().Get("clientSecret")
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
func (m *SocialIdentityProvider) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IdentityProviderBase.GetFieldDeserializers()
    res["clientId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientId(val)
        }
        return nil
    }
    res["clientSecret"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientSecret(val)
        }
        return nil
    }
    res["identityProviderType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentityProviderType(val)
        }
        return nil
    }
    return res
}
// GetIdentityProviderType gets the identityProviderType property value. For a B2B scenario, possible values: Google, Facebook. For a B2C scenario, possible values: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat. Required.
// returns a *string when successful
func (m *SocialIdentityProvider) GetIdentityProviderType()(*string) {
    val, err := m.GetBackingStore().Get("identityProviderType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SocialIdentityProvider) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IdentityProviderBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("clientId", m.GetClientId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("clientSecret", m.GetClientSecret())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityProviderType", m.GetIdentityProviderType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClientId sets the clientId property value. The identifier for the client application obtained when registering the application with the identity provider. Required.
func (m *SocialIdentityProvider) SetClientId(value *string)() {
    err := m.GetBackingStore().Set("clientId", value)
    if err != nil {
        panic(err)
    }
}
// SetClientSecret sets the clientSecret property value. The client secret for the application that is obtained when the application is registered with the identity provider. This is write-only. A read operation returns . Required.
func (m *SocialIdentityProvider) SetClientSecret(value *string)() {
    err := m.GetBackingStore().Set("clientSecret", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityProviderType sets the identityProviderType property value. For a B2B scenario, possible values: Google, Facebook. For a B2C scenario, possible values: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat. Required.
func (m *SocialIdentityProvider) SetIdentityProviderType(value *string)() {
    err := m.GetBackingStore().Set("identityProviderType", value)
    if err != nil {
        panic(err)
    }
}
type SocialIdentityProviderable interface {
    IdentityProviderBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClientId()(*string)
    GetClientSecret()(*string)
    GetIdentityProviderType()(*string)
    SetClientId(value *string)()
    SetClientSecret(value *string)()
    SetIdentityProviderType(value *string)()
}
