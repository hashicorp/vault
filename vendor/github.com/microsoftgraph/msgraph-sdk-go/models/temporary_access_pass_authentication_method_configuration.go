package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TemporaryAccessPassAuthenticationMethodConfiguration struct {
    AuthenticationMethodConfiguration
}
// NewTemporaryAccessPassAuthenticationMethodConfiguration instantiates a new TemporaryAccessPassAuthenticationMethodConfiguration and sets the default values.
func NewTemporaryAccessPassAuthenticationMethodConfiguration()(*TemporaryAccessPassAuthenticationMethodConfiguration) {
    m := &TemporaryAccessPassAuthenticationMethodConfiguration{
        AuthenticationMethodConfiguration: *NewAuthenticationMethodConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.temporaryAccessPassAuthenticationMethodConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTemporaryAccessPassAuthenticationMethodConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTemporaryAccessPassAuthenticationMethodConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTemporaryAccessPassAuthenticationMethodConfiguration(), nil
}
// GetDefaultLength gets the defaultLength property value. Default length in characters of a Temporary Access Pass object. Must be between 8 and 48 characters.
// returns a *int32 when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetDefaultLength()(*int32) {
    val, err := m.GetBackingStore().Get("defaultLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDefaultLifetimeInMinutes gets the defaultLifetimeInMinutes property value. Default lifetime in minutes for a Temporary Access Pass. Value can be any integer between the minimumLifetimeInMinutes and maximumLifetimeInMinutes.
// returns a *int32 when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetDefaultLifetimeInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("defaultLifetimeInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethodConfiguration.GetFieldDeserializers()
    res["defaultLength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultLength(val)
        }
        return nil
    }
    res["defaultLifetimeInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultLifetimeInMinutes(val)
        }
        return nil
    }
    res["includeTargets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationMethodTargetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodTargetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationMethodTargetable)
                }
            }
            m.SetIncludeTargets(res)
        }
        return nil
    }
    res["isUsableOnce"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsUsableOnce(val)
        }
        return nil
    }
    res["maximumLifetimeInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaximumLifetimeInMinutes(val)
        }
        return nil
    }
    res["minimumLifetimeInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumLifetimeInMinutes(val)
        }
        return nil
    }
    return res
}
// GetIncludeTargets gets the includeTargets property value. A collection of groups that are enabled to use the authentication method.
// returns a []AuthenticationMethodTargetable when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetIncludeTargets()([]AuthenticationMethodTargetable) {
    val, err := m.GetBackingStore().Get("includeTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodTargetable)
    }
    return nil
}
// GetIsUsableOnce gets the isUsableOnce property value. If true, all the passes in the tenant will be restricted to one-time use. If false, passes in the tenant can be created to be either one-time use or reusable.
// returns a *bool when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetIsUsableOnce()(*bool) {
    val, err := m.GetBackingStore().Get("isUsableOnce")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMaximumLifetimeInMinutes gets the maximumLifetimeInMinutes property value. Maximum lifetime in minutes for any Temporary Access Pass created in the tenant. Value can be between 10 and 43200 minutes (equivalent to 30 days).
// returns a *int32 when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetMaximumLifetimeInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("maximumLifetimeInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumLifetimeInMinutes gets the minimumLifetimeInMinutes property value. Minimum lifetime in minutes for any Temporary Access Pass created in the tenant. Value can be between 10 and 43200 minutes (equivalent to 30 days).
// returns a *int32 when successful
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) GetMinimumLifetimeInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("minimumLifetimeInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethodConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("defaultLength", m.GetDefaultLength())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("defaultLifetimeInMinutes", m.GetDefaultLifetimeInMinutes())
        if err != nil {
            return err
        }
    }
    if m.GetIncludeTargets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIncludeTargets()))
        for i, v := range m.GetIncludeTargets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("includeTargets", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isUsableOnce", m.GetIsUsableOnce())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("maximumLifetimeInMinutes", m.GetMaximumLifetimeInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minimumLifetimeInMinutes", m.GetMinimumLifetimeInMinutes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDefaultLength sets the defaultLength property value. Default length in characters of a Temporary Access Pass object. Must be between 8 and 48 characters.
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) SetDefaultLength(value *int32)() {
    err := m.GetBackingStore().Set("defaultLength", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultLifetimeInMinutes sets the defaultLifetimeInMinutes property value. Default lifetime in minutes for a Temporary Access Pass. Value can be any integer between the minimumLifetimeInMinutes and maximumLifetimeInMinutes.
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) SetDefaultLifetimeInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("defaultLifetimeInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeTargets sets the includeTargets property value. A collection of groups that are enabled to use the authentication method.
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) SetIncludeTargets(value []AuthenticationMethodTargetable)() {
    err := m.GetBackingStore().Set("includeTargets", value)
    if err != nil {
        panic(err)
    }
}
// SetIsUsableOnce sets the isUsableOnce property value. If true, all the passes in the tenant will be restricted to one-time use. If false, passes in the tenant can be created to be either one-time use or reusable.
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) SetIsUsableOnce(value *bool)() {
    err := m.GetBackingStore().Set("isUsableOnce", value)
    if err != nil {
        panic(err)
    }
}
// SetMaximumLifetimeInMinutes sets the maximumLifetimeInMinutes property value. Maximum lifetime in minutes for any Temporary Access Pass created in the tenant. Value can be between 10 and 43200 minutes (equivalent to 30 days).
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) SetMaximumLifetimeInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("maximumLifetimeInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumLifetimeInMinutes sets the minimumLifetimeInMinutes property value. Minimum lifetime in minutes for any Temporary Access Pass created in the tenant. Value can be between 10 and 43200 minutes (equivalent to 30 days).
func (m *TemporaryAccessPassAuthenticationMethodConfiguration) SetMinimumLifetimeInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("minimumLifetimeInMinutes", value)
    if err != nil {
        panic(err)
    }
}
type TemporaryAccessPassAuthenticationMethodConfigurationable interface {
    AuthenticationMethodConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDefaultLength()(*int32)
    GetDefaultLifetimeInMinutes()(*int32)
    GetIncludeTargets()([]AuthenticationMethodTargetable)
    GetIsUsableOnce()(*bool)
    GetMaximumLifetimeInMinutes()(*int32)
    GetMinimumLifetimeInMinutes()(*int32)
    SetDefaultLength(value *int32)()
    SetDefaultLifetimeInMinutes(value *int32)()
    SetIncludeTargets(value []AuthenticationMethodTargetable)()
    SetIsUsableOnce(value *bool)()
    SetMaximumLifetimeInMinutes(value *int32)()
    SetMinimumLifetimeInMinutes(value *int32)()
}
