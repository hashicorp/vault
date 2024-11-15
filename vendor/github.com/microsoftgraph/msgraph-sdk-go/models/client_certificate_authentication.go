package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ClientCertificateAuthentication struct {
    ApiAuthenticationConfigurationBase
}
// NewClientCertificateAuthentication instantiates a new ClientCertificateAuthentication and sets the default values.
func NewClientCertificateAuthentication()(*ClientCertificateAuthentication) {
    m := &ClientCertificateAuthentication{
        ApiAuthenticationConfigurationBase: *NewApiAuthenticationConfigurationBase(),
    }
    odataTypeValue := "#microsoft.graph.clientCertificateAuthentication"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateClientCertificateAuthenticationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateClientCertificateAuthenticationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewClientCertificateAuthentication(), nil
}
// GetCertificateList gets the certificateList property value. The list of certificates uploaded for this API connector.
// returns a []Pkcs12CertificateInformationable when successful
func (m *ClientCertificateAuthentication) GetCertificateList()([]Pkcs12CertificateInformationable) {
    val, err := m.GetBackingStore().Get("certificateList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Pkcs12CertificateInformationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ClientCertificateAuthentication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ApiAuthenticationConfigurationBase.GetFieldDeserializers()
    res["certificateList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePkcs12CertificateInformationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Pkcs12CertificateInformationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Pkcs12CertificateInformationable)
                }
            }
            m.SetCertificateList(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ClientCertificateAuthentication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ApiAuthenticationConfigurationBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCertificateList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCertificateList()))
        for i, v := range m.GetCertificateList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("certificateList", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCertificateList sets the certificateList property value. The list of certificates uploaded for this API connector.
func (m *ClientCertificateAuthentication) SetCertificateList(value []Pkcs12CertificateInformationable)() {
    err := m.GetBackingStore().Set("certificateList", value)
    if err != nil {
        panic(err)
    }
}
type ClientCertificateAuthenticationable interface {
    ApiAuthenticationConfigurationBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCertificateList()([]Pkcs12CertificateInformationable)
    SetCertificateList(value []Pkcs12CertificateInformationable)()
}
