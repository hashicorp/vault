package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SslCertificateEntity struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSslCertificateEntity instantiates a new SslCertificateEntity and sets the default values.
func NewSslCertificateEntity()(*SslCertificateEntity) {
    m := &SslCertificateEntity{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSslCertificateEntityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSslCertificateEntityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSslCertificateEntity(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SslCertificateEntity) GetAdditionalData()(map[string]any) {
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
// GetAddress gets the address property value. A physical address of the entity.
// returns a PhysicalAddressable when successful
func (m *SslCertificateEntity) GetAddress()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("address")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PhysicalAddressable)
    }
    return nil
}
// GetAlternateNames gets the alternateNames property value. Alternate names for this entity that are part of the certificate.
// returns a []string when successful
func (m *SslCertificateEntity) GetAlternateNames()([]string) {
    val, err := m.GetBackingStore().Get("alternateNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *SslCertificateEntity) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCommonName gets the commonName property value. A common name for this entity.
// returns a *string when successful
func (m *SslCertificateEntity) GetCommonName()(*string) {
    val, err := m.GetBackingStore().Get("commonName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmail gets the email property value. An email for this entity.
// returns a *string when successful
func (m *SslCertificateEntity) GetEmail()(*string) {
    val, err := m.GetBackingStore().Get("email")
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
func (m *SslCertificateEntity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["address"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAddress(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PhysicalAddressable))
        }
        return nil
    }
    res["alternateNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAlternateNames(res)
        }
        return nil
    }
    res["commonName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCommonName(val)
        }
        return nil
    }
    res["email"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmail(val)
        }
        return nil
    }
    res["givenName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGivenName(val)
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
    res["organizationName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrganizationName(val)
        }
        return nil
    }
    res["organizationUnitName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrganizationUnitName(val)
        }
        return nil
    }
    res["serialNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSerialNumber(val)
        }
        return nil
    }
    res["surname"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSurname(val)
        }
        return nil
    }
    return res
}
// GetGivenName gets the givenName property value. If the entity is a person, this is the person's given name (first name).
// returns a *string when successful
func (m *SslCertificateEntity) GetGivenName()(*string) {
    val, err := m.GetBackingStore().Get("givenName")
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
func (m *SslCertificateEntity) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrganizationName gets the organizationName property value. If the entity is an organization, this is the name of the organization.
// returns a *string when successful
func (m *SslCertificateEntity) GetOrganizationName()(*string) {
    val, err := m.GetBackingStore().Get("organizationName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrganizationUnitName gets the organizationUnitName property value. If the entity is an organization, this communicates if a unit in the organization is named on the entity.
// returns a *string when successful
func (m *SslCertificateEntity) GetOrganizationUnitName()(*string) {
    val, err := m.GetBackingStore().Get("organizationUnitName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSerialNumber gets the serialNumber property value. A serial number assigned to the entity; usually only available if the entity is the issuer.
// returns a *string when successful
func (m *SslCertificateEntity) GetSerialNumber()(*string) {
    val, err := m.GetBackingStore().Get("serialNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSurname gets the surname property value. If the entity is a person, this is the person's surname (last name).
// returns a *string when successful
func (m *SslCertificateEntity) GetSurname()(*string) {
    val, err := m.GetBackingStore().Get("surname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SslCertificateEntity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("address", m.GetAddress())
        if err != nil {
            return err
        }
    }
    if m.GetAlternateNames() != nil {
        err := writer.WriteCollectionOfStringValues("alternateNames", m.GetAlternateNames())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("commonName", m.GetCommonName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("givenName", m.GetGivenName())
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
        err := writer.WriteStringValue("organizationName", m.GetOrganizationName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("organizationUnitName", m.GetOrganizationUnitName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("serialNumber", m.GetSerialNumber())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("surname", m.GetSurname())
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
func (m *SslCertificateEntity) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAddress sets the address property value. A physical address of the entity.
func (m *SslCertificateEntity) SetAddress(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PhysicalAddressable)() {
    err := m.GetBackingStore().Set("address", value)
    if err != nil {
        panic(err)
    }
}
// SetAlternateNames sets the alternateNames property value. Alternate names for this entity that are part of the certificate.
func (m *SslCertificateEntity) SetAlternateNames(value []string)() {
    err := m.GetBackingStore().Set("alternateNames", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SslCertificateEntity) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCommonName sets the commonName property value. A common name for this entity.
func (m *SslCertificateEntity) SetCommonName(value *string)() {
    err := m.GetBackingStore().Set("commonName", value)
    if err != nil {
        panic(err)
    }
}
// SetEmail sets the email property value. An email for this entity.
func (m *SslCertificateEntity) SetEmail(value *string)() {
    err := m.GetBackingStore().Set("email", value)
    if err != nil {
        panic(err)
    }
}
// SetGivenName sets the givenName property value. If the entity is a person, this is the person's given name (first name).
func (m *SslCertificateEntity) SetGivenName(value *string)() {
    err := m.GetBackingStore().Set("givenName", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SslCertificateEntity) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOrganizationName sets the organizationName property value. If the entity is an organization, this is the name of the organization.
func (m *SslCertificateEntity) SetOrganizationName(value *string)() {
    err := m.GetBackingStore().Set("organizationName", value)
    if err != nil {
        panic(err)
    }
}
// SetOrganizationUnitName sets the organizationUnitName property value. If the entity is an organization, this communicates if a unit in the organization is named on the entity.
func (m *SslCertificateEntity) SetOrganizationUnitName(value *string)() {
    err := m.GetBackingStore().Set("organizationUnitName", value)
    if err != nil {
        panic(err)
    }
}
// SetSerialNumber sets the serialNumber property value. A serial number assigned to the entity; usually only available if the entity is the issuer.
func (m *SslCertificateEntity) SetSerialNumber(value *string)() {
    err := m.GetBackingStore().Set("serialNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetSurname sets the surname property value. If the entity is a person, this is the person's surname (last name).
func (m *SslCertificateEntity) SetSurname(value *string)() {
    err := m.GetBackingStore().Set("surname", value)
    if err != nil {
        panic(err)
    }
}
type SslCertificateEntityable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAddress()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PhysicalAddressable)
    GetAlternateNames()([]string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCommonName()(*string)
    GetEmail()(*string)
    GetGivenName()(*string)
    GetOdataType()(*string)
    GetOrganizationName()(*string)
    GetOrganizationUnitName()(*string)
    GetSerialNumber()(*string)
    GetSurname()(*string)
    SetAddress(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PhysicalAddressable)()
    SetAlternateNames(value []string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCommonName(value *string)()
    SetEmail(value *string)()
    SetGivenName(value *string)()
    SetOdataType(value *string)()
    SetOrganizationName(value *string)()
    SetOrganizationUnitName(value *string)()
    SetSerialNumber(value *string)()
    SetSurname(value *string)()
}
