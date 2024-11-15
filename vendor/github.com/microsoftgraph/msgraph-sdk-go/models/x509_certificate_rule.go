package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type X509CertificateRule struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewX509CertificateRule instantiates a new X509CertificateRule and sets the default values.
func NewX509CertificateRule()(*X509CertificateRule) {
    m := &X509CertificateRule{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateX509CertificateRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateX509CertificateRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewX509CertificateRule(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *X509CertificateRule) GetAdditionalData()(map[string]any) {
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
func (m *X509CertificateRule) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *X509CertificateRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["identifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentifier(val)
        }
        return nil
    }
    res["issuerSubjectIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIssuerSubjectIdentifier(val)
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
    res["policyOidIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyOidIdentifier(val)
        }
        return nil
    }
    res["x509CertificateAuthenticationMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseX509CertificateAuthenticationMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetX509CertificateAuthenticationMode(val.(*X509CertificateAuthenticationMode))
        }
        return nil
    }
    res["x509CertificateRequiredAffinityLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseX509CertificateAffinityLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetX509CertificateRequiredAffinityLevel(val.(*X509CertificateAffinityLevel))
        }
        return nil
    }
    res["x509CertificateRuleType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseX509CertificateRuleType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetX509CertificateRuleType(val.(*X509CertificateRuleType))
        }
        return nil
    }
    return res
}
// GetIdentifier gets the identifier property value. The identifier of the X.509 certificate. Required.
// returns a *string when successful
func (m *X509CertificateRule) GetIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("identifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIssuerSubjectIdentifier gets the issuerSubjectIdentifier property value. The issuerSubjectIdentifier property
// returns a *string when successful
func (m *X509CertificateRule) GetIssuerSubjectIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("issuerSubjectIdentifier")
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
func (m *X509CertificateRule) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPolicyOidIdentifier gets the policyOidIdentifier property value. The policyOidIdentifier property
// returns a *string when successful
func (m *X509CertificateRule) GetPolicyOidIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("policyOidIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetX509CertificateAuthenticationMode gets the x509CertificateAuthenticationMode property value. The type of strong authentication mode. The possible values are: x509CertificateSingleFactor, x509CertificateMultiFactor, unknownFutureValue. Required.
// returns a *X509CertificateAuthenticationMode when successful
func (m *X509CertificateRule) GetX509CertificateAuthenticationMode()(*X509CertificateAuthenticationMode) {
    val, err := m.GetBackingStore().Get("x509CertificateAuthenticationMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*X509CertificateAuthenticationMode)
    }
    return nil
}
// GetX509CertificateRequiredAffinityLevel gets the x509CertificateRequiredAffinityLevel property value. The x509CertificateRequiredAffinityLevel property
// returns a *X509CertificateAffinityLevel when successful
func (m *X509CertificateRule) GetX509CertificateRequiredAffinityLevel()(*X509CertificateAffinityLevel) {
    val, err := m.GetBackingStore().Get("x509CertificateRequiredAffinityLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*X509CertificateAffinityLevel)
    }
    return nil
}
// GetX509CertificateRuleType gets the x509CertificateRuleType property value. The type of the X.509 certificate mode configuration rule. The possible values are: issuerSubject, policyOID, unknownFutureValue. Required.
// returns a *X509CertificateRuleType when successful
func (m *X509CertificateRule) GetX509CertificateRuleType()(*X509CertificateRuleType) {
    val, err := m.GetBackingStore().Get("x509CertificateRuleType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*X509CertificateRuleType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *X509CertificateRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("identifier", m.GetIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("issuerSubjectIdentifier", m.GetIssuerSubjectIdentifier())
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
        err := writer.WriteStringValue("policyOidIdentifier", m.GetPolicyOidIdentifier())
        if err != nil {
            return err
        }
    }
    if m.GetX509CertificateAuthenticationMode() != nil {
        cast := (*m.GetX509CertificateAuthenticationMode()).String()
        err := writer.WriteStringValue("x509CertificateAuthenticationMode", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetX509CertificateRequiredAffinityLevel() != nil {
        cast := (*m.GetX509CertificateRequiredAffinityLevel()).String()
        err := writer.WriteStringValue("x509CertificateRequiredAffinityLevel", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetX509CertificateRuleType() != nil {
        cast := (*m.GetX509CertificateRuleType()).String()
        err := writer.WriteStringValue("x509CertificateRuleType", &cast)
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
func (m *X509CertificateRule) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *X509CertificateRule) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIdentifier sets the identifier property value. The identifier of the X.509 certificate. Required.
func (m *X509CertificateRule) SetIdentifier(value *string)() {
    err := m.GetBackingStore().Set("identifier", value)
    if err != nil {
        panic(err)
    }
}
// SetIssuerSubjectIdentifier sets the issuerSubjectIdentifier property value. The issuerSubjectIdentifier property
func (m *X509CertificateRule) SetIssuerSubjectIdentifier(value *string)() {
    err := m.GetBackingStore().Set("issuerSubjectIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *X509CertificateRule) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyOidIdentifier sets the policyOidIdentifier property value. The policyOidIdentifier property
func (m *X509CertificateRule) SetPolicyOidIdentifier(value *string)() {
    err := m.GetBackingStore().Set("policyOidIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetX509CertificateAuthenticationMode sets the x509CertificateAuthenticationMode property value. The type of strong authentication mode. The possible values are: x509CertificateSingleFactor, x509CertificateMultiFactor, unknownFutureValue. Required.
func (m *X509CertificateRule) SetX509CertificateAuthenticationMode(value *X509CertificateAuthenticationMode)() {
    err := m.GetBackingStore().Set("x509CertificateAuthenticationMode", value)
    if err != nil {
        panic(err)
    }
}
// SetX509CertificateRequiredAffinityLevel sets the x509CertificateRequiredAffinityLevel property value. The x509CertificateRequiredAffinityLevel property
func (m *X509CertificateRule) SetX509CertificateRequiredAffinityLevel(value *X509CertificateAffinityLevel)() {
    err := m.GetBackingStore().Set("x509CertificateRequiredAffinityLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetX509CertificateRuleType sets the x509CertificateRuleType property value. The type of the X.509 certificate mode configuration rule. The possible values are: issuerSubject, policyOID, unknownFutureValue. Required.
func (m *X509CertificateRule) SetX509CertificateRuleType(value *X509CertificateRuleType)() {
    err := m.GetBackingStore().Set("x509CertificateRuleType", value)
    if err != nil {
        panic(err)
    }
}
type X509CertificateRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIdentifier()(*string)
    GetIssuerSubjectIdentifier()(*string)
    GetOdataType()(*string)
    GetPolicyOidIdentifier()(*string)
    GetX509CertificateAuthenticationMode()(*X509CertificateAuthenticationMode)
    GetX509CertificateRequiredAffinityLevel()(*X509CertificateAffinityLevel)
    GetX509CertificateRuleType()(*X509CertificateRuleType)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIdentifier(value *string)()
    SetIssuerSubjectIdentifier(value *string)()
    SetOdataType(value *string)()
    SetPolicyOidIdentifier(value *string)()
    SetX509CertificateAuthenticationMode(value *X509CertificateAuthenticationMode)()
    SetX509CertificateRequiredAffinityLevel(value *X509CertificateAffinityLevel)()
    SetX509CertificateRuleType(value *X509CertificateRuleType)()
}
