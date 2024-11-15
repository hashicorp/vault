package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EmailAuthenticationMethodConfiguration struct {
    AuthenticationMethodConfiguration
}
// NewEmailAuthenticationMethodConfiguration instantiates a new EmailAuthenticationMethodConfiguration and sets the default values.
func NewEmailAuthenticationMethodConfiguration()(*EmailAuthenticationMethodConfiguration) {
    m := &EmailAuthenticationMethodConfiguration{
        AuthenticationMethodConfiguration: *NewAuthenticationMethodConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.emailAuthenticationMethodConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEmailAuthenticationMethodConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEmailAuthenticationMethodConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmailAuthenticationMethodConfiguration(), nil
}
// GetAllowExternalIdToUseEmailOtp gets the allowExternalIdToUseEmailOtp property value. Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled, unknownFutureValue. Tenants in the default state who didn't use public preview have email OTP enabled beginning in October 2021.
// returns a *ExternalEmailOtpState when successful
func (m *EmailAuthenticationMethodConfiguration) GetAllowExternalIdToUseEmailOtp()(*ExternalEmailOtpState) {
    val, err := m.GetBackingStore().Get("allowExternalIdToUseEmailOtp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ExternalEmailOtpState)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EmailAuthenticationMethodConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethodConfiguration.GetFieldDeserializers()
    res["allowExternalIdToUseEmailOtp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseExternalEmailOtpState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowExternalIdToUseEmailOtp(val.(*ExternalEmailOtpState))
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
    return res
}
// GetIncludeTargets gets the includeTargets property value. A collection of groups that are enabled to use the authentication method.
// returns a []AuthenticationMethodTargetable when successful
func (m *EmailAuthenticationMethodConfiguration) GetIncludeTargets()([]AuthenticationMethodTargetable) {
    val, err := m.GetBackingStore().Get("includeTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodTargetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EmailAuthenticationMethodConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethodConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowExternalIdToUseEmailOtp() != nil {
        cast := (*m.GetAllowExternalIdToUseEmailOtp()).String()
        err = writer.WriteStringValue("allowExternalIdToUseEmailOtp", &cast)
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
    return nil
}
// SetAllowExternalIdToUseEmailOtp sets the allowExternalIdToUseEmailOtp property value. Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled, unknownFutureValue. Tenants in the default state who didn't use public preview have email OTP enabled beginning in October 2021.
func (m *EmailAuthenticationMethodConfiguration) SetAllowExternalIdToUseEmailOtp(value *ExternalEmailOtpState)() {
    err := m.GetBackingStore().Set("allowExternalIdToUseEmailOtp", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeTargets sets the includeTargets property value. A collection of groups that are enabled to use the authentication method.
func (m *EmailAuthenticationMethodConfiguration) SetIncludeTargets(value []AuthenticationMethodTargetable)() {
    err := m.GetBackingStore().Set("includeTargets", value)
    if err != nil {
        panic(err)
    }
}
type EmailAuthenticationMethodConfigurationable interface {
    AuthenticationMethodConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowExternalIdToUseEmailOtp()(*ExternalEmailOtpState)
    GetIncludeTargets()([]AuthenticationMethodTargetable)
    SetAllowExternalIdToUseEmailOtp(value *ExternalEmailOtpState)()
    SetIncludeTargets(value []AuthenticationMethodTargetable)()
}
