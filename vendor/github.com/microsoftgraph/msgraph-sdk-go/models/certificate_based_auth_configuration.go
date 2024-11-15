package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CertificateBasedAuthConfiguration struct {
    Entity
}
// NewCertificateBasedAuthConfiguration instantiates a new CertificateBasedAuthConfiguration and sets the default values.
func NewCertificateBasedAuthConfiguration()(*CertificateBasedAuthConfiguration) {
    m := &CertificateBasedAuthConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCertificateBasedAuthConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCertificateBasedAuthConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCertificateBasedAuthConfiguration(), nil
}
// GetCertificateAuthorities gets the certificateAuthorities property value. Collection of certificate authorities which creates a trusted certificate chain.
// returns a []CertificateAuthorityable when successful
func (m *CertificateBasedAuthConfiguration) GetCertificateAuthorities()([]CertificateAuthorityable) {
    val, err := m.GetBackingStore().Get("certificateAuthorities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CertificateAuthorityable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CertificateBasedAuthConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["certificateAuthorities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCertificateAuthorityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CertificateAuthorityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CertificateAuthorityable)
                }
            }
            m.SetCertificateAuthorities(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CertificateBasedAuthConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCertificateAuthorities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCertificateAuthorities()))
        for i, v := range m.GetCertificateAuthorities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("certificateAuthorities", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCertificateAuthorities sets the certificateAuthorities property value. Collection of certificate authorities which creates a trusted certificate chain.
func (m *CertificateBasedAuthConfiguration) SetCertificateAuthorities(value []CertificateAuthorityable)() {
    err := m.GetBackingStore().Set("certificateAuthorities", value)
    if err != nil {
        panic(err)
    }
}
type CertificateBasedAuthConfigurationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCertificateAuthorities()([]CertificateAuthorityable)
    SetCertificateAuthorities(value []CertificateAuthorityable)()
}
