package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Fido2AuthenticationMethodConfiguration struct {
    AuthenticationMethodConfiguration
}
// NewFido2AuthenticationMethodConfiguration instantiates a new Fido2AuthenticationMethodConfiguration and sets the default values.
func NewFido2AuthenticationMethodConfiguration()(*Fido2AuthenticationMethodConfiguration) {
    m := &Fido2AuthenticationMethodConfiguration{
        AuthenticationMethodConfiguration: *NewAuthenticationMethodConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.fido2AuthenticationMethodConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateFido2AuthenticationMethodConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFido2AuthenticationMethodConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFido2AuthenticationMethodConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Fido2AuthenticationMethodConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethodConfiguration.GetFieldDeserializers()
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
    res["isAttestationEnforced"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAttestationEnforced(val)
        }
        return nil
    }
    res["isSelfServiceRegistrationAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSelfServiceRegistrationAllowed(val)
        }
        return nil
    }
    res["keyRestrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFido2KeyRestrictionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyRestrictions(val.(Fido2KeyRestrictionsable))
        }
        return nil
    }
    return res
}
// GetIncludeTargets gets the includeTargets property value. A collection of groups that are enabled to use the authentication method.
// returns a []AuthenticationMethodTargetable when successful
func (m *Fido2AuthenticationMethodConfiguration) GetIncludeTargets()([]AuthenticationMethodTargetable) {
    val, err := m.GetBackingStore().Get("includeTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodTargetable)
    }
    return nil
}
// GetIsAttestationEnforced gets the isAttestationEnforced property value. Determines whether attestation must be enforced for FIDO2 security key registration.
// returns a *bool when successful
func (m *Fido2AuthenticationMethodConfiguration) GetIsAttestationEnforced()(*bool) {
    val, err := m.GetBackingStore().Get("isAttestationEnforced")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSelfServiceRegistrationAllowed gets the isSelfServiceRegistrationAllowed property value. Determines if users can register new FIDO2 security keys.
// returns a *bool when successful
func (m *Fido2AuthenticationMethodConfiguration) GetIsSelfServiceRegistrationAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("isSelfServiceRegistrationAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyRestrictions gets the keyRestrictions property value. Controls whether key restrictions are enforced on FIDO2 security keys, either allowing or disallowing certain key types as defined by Authenticator Attestation GUID (AAGUID), an identifier that indicates the type (e.g. make and model) of the authenticator.
// returns a Fido2KeyRestrictionsable when successful
func (m *Fido2AuthenticationMethodConfiguration) GetKeyRestrictions()(Fido2KeyRestrictionsable) {
    val, err := m.GetBackingStore().Get("keyRestrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Fido2KeyRestrictionsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Fido2AuthenticationMethodConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethodConfiguration.Serialize(writer)
    if err != nil {
        return err
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
        err = writer.WriteBoolValue("isAttestationEnforced", m.GetIsAttestationEnforced())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSelfServiceRegistrationAllowed", m.GetIsSelfServiceRegistrationAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("keyRestrictions", m.GetKeyRestrictions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIncludeTargets sets the includeTargets property value. A collection of groups that are enabled to use the authentication method.
func (m *Fido2AuthenticationMethodConfiguration) SetIncludeTargets(value []AuthenticationMethodTargetable)() {
    err := m.GetBackingStore().Set("includeTargets", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAttestationEnforced sets the isAttestationEnforced property value. Determines whether attestation must be enforced for FIDO2 security key registration.
func (m *Fido2AuthenticationMethodConfiguration) SetIsAttestationEnforced(value *bool)() {
    err := m.GetBackingStore().Set("isAttestationEnforced", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSelfServiceRegistrationAllowed sets the isSelfServiceRegistrationAllowed property value. Determines if users can register new FIDO2 security keys.
func (m *Fido2AuthenticationMethodConfiguration) SetIsSelfServiceRegistrationAllowed(value *bool)() {
    err := m.GetBackingStore().Set("isSelfServiceRegistrationAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyRestrictions sets the keyRestrictions property value. Controls whether key restrictions are enforced on FIDO2 security keys, either allowing or disallowing certain key types as defined by Authenticator Attestation GUID (AAGUID), an identifier that indicates the type (e.g. make and model) of the authenticator.
func (m *Fido2AuthenticationMethodConfiguration) SetKeyRestrictions(value Fido2KeyRestrictionsable)() {
    err := m.GetBackingStore().Set("keyRestrictions", value)
    if err != nil {
        panic(err)
    }
}
type Fido2AuthenticationMethodConfigurationable interface {
    AuthenticationMethodConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIncludeTargets()([]AuthenticationMethodTargetable)
    GetIsAttestationEnforced()(*bool)
    GetIsSelfServiceRegistrationAllowed()(*bool)
    GetKeyRestrictions()(Fido2KeyRestrictionsable)
    SetIncludeTargets(value []AuthenticationMethodTargetable)()
    SetIsAttestationEnforced(value *bool)()
    SetIsSelfServiceRegistrationAllowed(value *bool)()
    SetKeyRestrictions(value Fido2KeyRestrictionsable)()
}
