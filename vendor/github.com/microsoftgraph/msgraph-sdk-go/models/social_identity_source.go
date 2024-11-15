package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SocialIdentitySource struct {
    IdentitySource
}
// NewSocialIdentitySource instantiates a new SocialIdentitySource and sets the default values.
func NewSocialIdentitySource()(*SocialIdentitySource) {
    m := &SocialIdentitySource{
        IdentitySource: *NewIdentitySource(),
    }
    odataTypeValue := "#microsoft.graph.socialIdentitySource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSocialIdentitySourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSocialIdentitySourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSocialIdentitySource(), nil
}
// GetDisplayName gets the displayName property value. The displayName property
// returns a *string when successful
func (m *SocialIdentitySource) GetDisplayName()(*string) {
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
func (m *SocialIdentitySource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IdentitySource.GetFieldDeserializers()
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
    res["socialIdentitySourceType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSocialIdentitySourceType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSocialIdentitySourceType(val.(*SocialIdentitySourceType))
        }
        return nil
    }
    return res
}
// GetSocialIdentitySourceType gets the socialIdentitySourceType property value. The socialIdentitySourceType property
// returns a *SocialIdentitySourceType when successful
func (m *SocialIdentitySource) GetSocialIdentitySourceType()(*SocialIdentitySourceType) {
    val, err := m.GetBackingStore().Get("socialIdentitySourceType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SocialIdentitySourceType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SocialIdentitySource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IdentitySource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetSocialIdentitySourceType() != nil {
        cast := (*m.GetSocialIdentitySourceType()).String()
        err = writer.WriteStringValue("socialIdentitySourceType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The displayName property
func (m *SocialIdentitySource) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetSocialIdentitySourceType sets the socialIdentitySourceType property value. The socialIdentitySourceType property
func (m *SocialIdentitySource) SetSocialIdentitySourceType(value *SocialIdentitySourceType)() {
    err := m.GetBackingStore().Set("socialIdentitySourceType", value)
    if err != nil {
        panic(err)
    }
}
type SocialIdentitySourceable interface {
    IdentitySourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetSocialIdentitySourceType()(*SocialIdentitySourceType)
    SetDisplayName(value *string)()
    SetSocialIdentitySourceType(value *SocialIdentitySourceType)()
}
