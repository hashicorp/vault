package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CertificateAuthority struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCertificateAuthority instantiates a new CertificateAuthority and sets the default values.
func NewCertificateAuthority()(*CertificateAuthority) {
    m := &CertificateAuthority{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCertificateAuthorityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCertificateAuthorityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCertificateAuthority(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CertificateAuthority) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *CertificateAuthority) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCertificate gets the certificate property value. Required. The base64 encoded string representing the public certificate.
// returns a []byte when successful
func (m *CertificateAuthority) GetCertificate()([]byte) {
    val, err := m.GetBackingStore().Get("certificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetCertificateRevocationListUrl gets the certificateRevocationListUrl property value. The URL of the certificate revocation list.
// returns a *string when successful
func (m *CertificateAuthority) GetCertificateRevocationListUrl()(*string) {
    val, err := m.GetBackingStore().Get("certificateRevocationListUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeltaCertificateRevocationListUrl gets the deltaCertificateRevocationListUrl property value. The URL contains the list of all revoked certificates since the last time a full certificate revocaton list was created.
// returns a *string when successful
func (m *CertificateAuthority) GetDeltaCertificateRevocationListUrl()(*string) {
    val, err := m.GetBackingStore().Get("deltaCertificateRevocationListUrl")
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
func (m *CertificateAuthority) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["certificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificate(val)
        }
        return nil
    }
    res["certificateRevocationListUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificateRevocationListUrl(val)
        }
        return nil
    }
    res["deltaCertificateRevocationListUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeltaCertificateRevocationListUrl(val)
        }
        return nil
    }
    res["isRootAuthority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRootAuthority(val)
        }
        return nil
    }
    res["issuer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIssuer(val)
        }
        return nil
    }
    res["issuerSki"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIssuerSki(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    return res
}
// GetIsRootAuthority gets the isRootAuthority property value. Required. true if the trusted certificate is a root authority, false if the trusted certificate is an intermediate authority.
// returns a *bool when successful
func (m *CertificateAuthority) GetIsRootAuthority()(*bool) {
    val, err := m.GetBackingStore().Get("isRootAuthority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIssuer gets the issuer property value. The issuer of the certificate, calculated from the certificate value. Read-only.
// returns a *string when successful
func (m *CertificateAuthority) GetIssuer()(*string) {
    val, err := m.GetBackingStore().Get("issuer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIssuerSki gets the issuerSki property value. The subject key identifier of the certificate, calculated from the certificate value. Read-only.
// returns a *string when successful
func (m *CertificateAuthority) GetIssuerSki()(*string) {
    val, err := m.GetBackingStore().Get("issuerSki")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *CertificateAuthority) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CertificateAuthority) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteByteArrayValue("certificate", m.GetCertificate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("certificateRevocationListUrl", m.GetCertificateRevocationListUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deltaCertificateRevocationListUrl", m.GetDeltaCertificateRevocationListUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRootAuthority", m.GetIsRootAuthority())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("issuer", m.GetIssuer())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("issuerSki", m.GetIssuerSki())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CertificateAuthority) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CertificateAuthority) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCertificate sets the certificate property value. Required. The base64 encoded string representing the public certificate.
func (m *CertificateAuthority) SetCertificate(value []byte)() {
    err := m.GetBackingStore().Set("certificate", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificateRevocationListUrl sets the certificateRevocationListUrl property value. The URL of the certificate revocation list.
func (m *CertificateAuthority) SetCertificateRevocationListUrl(value *string)() {
    err := m.GetBackingStore().Set("certificateRevocationListUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetDeltaCertificateRevocationListUrl sets the deltaCertificateRevocationListUrl property value. The URL contains the list of all revoked certificates since the last time a full certificate revocaton list was created.
func (m *CertificateAuthority) SetDeltaCertificateRevocationListUrl(value *string)() {
    err := m.GetBackingStore().Set("deltaCertificateRevocationListUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRootAuthority sets the isRootAuthority property value. Required. true if the trusted certificate is a root authority, false if the trusted certificate is an intermediate authority.
func (m *CertificateAuthority) SetIsRootAuthority(value *bool)() {
    err := m.GetBackingStore().Set("isRootAuthority", value)
    if err != nil {
        panic(err)
    }
}
// SetIssuer sets the issuer property value. The issuer of the certificate, calculated from the certificate value. Read-only.
func (m *CertificateAuthority) SetIssuer(value *string)() {
    err := m.GetBackingStore().Set("issuer", value)
    if err != nil {
        panic(err)
    }
}
// SetIssuerSki sets the issuerSki property value. The subject key identifier of the certificate, calculated from the certificate value. Read-only.
func (m *CertificateAuthority) SetIssuerSki(value *string)() {
    err := m.GetBackingStore().Set("issuerSki", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *CertificateAuthority) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type CertificateAuthorityable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCertificate()([]byte)
    GetCertificateRevocationListUrl()(*string)
    GetDeltaCertificateRevocationListUrl()(*string)
    GetIsRootAuthority()(*bool)
    GetIssuer()(*string)
    GetIssuerSki()(*string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCertificate(value []byte)()
    SetCertificateRevocationListUrl(value *string)()
    SetDeltaCertificateRevocationListUrl(value *string)()
    SetIsRootAuthority(value *bool)()
    SetIssuer(value *string)()
    SetIssuerSki(value *string)()
    SetOdataType(value *string)()
}
