package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationMethodConfiguration struct {
    Entity
}
// NewAuthenticationMethodConfiguration instantiates a new AuthenticationMethodConfiguration and sets the default values.
func NewAuthenticationMethodConfiguration()(*AuthenticationMethodConfiguration) {
    m := &AuthenticationMethodConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationMethodConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationMethodConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.emailAuthenticationMethodConfiguration":
                        return NewEmailAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.fido2AuthenticationMethodConfiguration":
                        return NewFido2AuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.microsoftAuthenticatorAuthenticationMethodConfiguration":
                        return NewMicrosoftAuthenticatorAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.smsAuthenticationMethodConfiguration":
                        return NewSmsAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.softwareOathAuthenticationMethodConfiguration":
                        return NewSoftwareOathAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.temporaryAccessPassAuthenticationMethodConfiguration":
                        return NewTemporaryAccessPassAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.voiceAuthenticationMethodConfiguration":
                        return NewVoiceAuthenticationMethodConfiguration(), nil
                    case "#microsoft.graph.x509CertificateAuthenticationMethodConfiguration":
                        return NewX509CertificateAuthenticationMethodConfiguration(), nil
                }
            }
        }
    }
    return NewAuthenticationMethodConfiguration(), nil
}
// GetExcludeTargets gets the excludeTargets property value. Groups of users that are excluded from a policy.
// returns a []ExcludeTargetable when successful
func (m *AuthenticationMethodConfiguration) GetExcludeTargets()([]ExcludeTargetable) {
    val, err := m.GetBackingStore().Get("excludeTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ExcludeTargetable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthenticationMethodConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["excludeTargets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExcludeTargetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ExcludeTargetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ExcludeTargetable)
                }
            }
            m.SetExcludeTargets(res)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationMethodState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*AuthenticationMethodState))
        }
        return nil
    }
    return res
}
// GetState gets the state property value. The state of the policy. Possible values are: enabled, disabled.
// returns a *AuthenticationMethodState when successful
func (m *AuthenticationMethodConfiguration) GetState()(*AuthenticationMethodState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationMethodState)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationMethodConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetExcludeTargets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExcludeTargets()))
        for i, v := range m.GetExcludeTargets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("excludeTargets", cast)
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExcludeTargets sets the excludeTargets property value. Groups of users that are excluded from a policy.
func (m *AuthenticationMethodConfiguration) SetExcludeTargets(value []ExcludeTargetable)() {
    err := m.GetBackingStore().Set("excludeTargets", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state of the policy. Possible values are: enabled, disabled.
func (m *AuthenticationMethodConfiguration) SetState(value *AuthenticationMethodState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationMethodConfigurationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetExcludeTargets()([]ExcludeTargetable)
    GetState()(*AuthenticationMethodState)
    SetExcludeTargets(value []ExcludeTargetable)()
    SetState(value *AuthenticationMethodState)()
}
