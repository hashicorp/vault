package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type X509CertificateCombinationConfiguration struct {
    AuthenticationCombinationConfiguration
}
// NewX509CertificateCombinationConfiguration instantiates a new X509CertificateCombinationConfiguration and sets the default values.
func NewX509CertificateCombinationConfiguration()(*X509CertificateCombinationConfiguration) {
    m := &X509CertificateCombinationConfiguration{
        AuthenticationCombinationConfiguration: *NewAuthenticationCombinationConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.x509CertificateCombinationConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateX509CertificateCombinationConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateX509CertificateCombinationConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewX509CertificateCombinationConfiguration(), nil
}
// GetAllowedIssuerSkis gets the allowedIssuerSkis property value. A list of allowed subject key identifier values.
// returns a []string when successful
func (m *X509CertificateCombinationConfiguration) GetAllowedIssuerSkis()([]string) {
    val, err := m.GetBackingStore().Get("allowedIssuerSkis")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAllowedPolicyOIDs gets the allowedPolicyOIDs property value. A list of allowed policy OIDs.
// returns a []string when successful
func (m *X509CertificateCombinationConfiguration) GetAllowedPolicyOIDs()([]string) {
    val, err := m.GetBackingStore().Get("allowedPolicyOIDs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *X509CertificateCombinationConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationCombinationConfiguration.GetFieldDeserializers()
    res["allowedIssuerSkis"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetAllowedIssuerSkis(res)
        }
        return nil
    }
    res["allowedPolicyOIDs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetAllowedPolicyOIDs(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *X509CertificateCombinationConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationCombinationConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedIssuerSkis() != nil {
        err = writer.WriteCollectionOfStringValues("allowedIssuerSkis", m.GetAllowedIssuerSkis())
        if err != nil {
            return err
        }
    }
    if m.GetAllowedPolicyOIDs() != nil {
        err = writer.WriteCollectionOfStringValues("allowedPolicyOIDs", m.GetAllowedPolicyOIDs())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedIssuerSkis sets the allowedIssuerSkis property value. A list of allowed subject key identifier values.
func (m *X509CertificateCombinationConfiguration) SetAllowedIssuerSkis(value []string)() {
    err := m.GetBackingStore().Set("allowedIssuerSkis", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedPolicyOIDs sets the allowedPolicyOIDs property value. A list of allowed policy OIDs.
func (m *X509CertificateCombinationConfiguration) SetAllowedPolicyOIDs(value []string)() {
    err := m.GetBackingStore().Set("allowedPolicyOIDs", value)
    if err != nil {
        panic(err)
    }
}
type X509CertificateCombinationConfigurationable interface {
    AuthenticationCombinationConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedIssuerSkis()([]string)
    GetAllowedPolicyOIDs()([]string)
    SetAllowedIssuerSkis(value []string)()
    SetAllowedPolicyOIDs(value []string)()
}
