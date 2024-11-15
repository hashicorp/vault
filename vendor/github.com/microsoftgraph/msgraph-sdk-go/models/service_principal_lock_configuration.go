package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ServicePrincipalLockConfiguration struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewServicePrincipalLockConfiguration instantiates a new ServicePrincipalLockConfiguration and sets the default values.
func NewServicePrincipalLockConfiguration()(*ServicePrincipalLockConfiguration) {
    m := &ServicePrincipalLockConfiguration{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateServicePrincipalLockConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServicePrincipalLockConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServicePrincipalLockConfiguration(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ServicePrincipalLockConfiguration) GetAdditionalData()(map[string]any) {
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
// GetAllProperties gets the allProperties property value. Enables locking all sensitive properties. The sensitive properties are keyCredentials, passwordCredentials, and tokenEncryptionKeyId.
// returns a *bool when successful
func (m *ServicePrincipalLockConfiguration) GetAllProperties()(*bool) {
    val, err := m.GetBackingStore().Get("allProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ServicePrincipalLockConfiguration) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCredentialsWithUsageSign gets the credentialsWithUsageSign property value. Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Sign.
// returns a *bool when successful
func (m *ServicePrincipalLockConfiguration) GetCredentialsWithUsageSign()(*bool) {
    val, err := m.GetBackingStore().Get("credentialsWithUsageSign")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCredentialsWithUsageVerify gets the credentialsWithUsageVerify property value. Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Verify. This locks OAuth service principals.
// returns a *bool when successful
func (m *ServicePrincipalLockConfiguration) GetCredentialsWithUsageVerify()(*bool) {
    val, err := m.GetBackingStore().Get("credentialsWithUsageVerify")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServicePrincipalLockConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllProperties(val)
        }
        return nil
    }
    res["credentialsWithUsageSign"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCredentialsWithUsageSign(val)
        }
        return nil
    }
    res["credentialsWithUsageVerify"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCredentialsWithUsageVerify(val)
        }
        return nil
    }
    res["isEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabled(val)
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
    res["tokenEncryptionKeyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTokenEncryptionKeyId(val)
        }
        return nil
    }
    return res
}
// GetIsEnabled gets the isEnabled property value. Enables or disables service principal lock configuration. To allow the sensitive properties to be updated, update this property to false to disable the lock on the service principal.
// returns a *bool when successful
func (m *ServicePrincipalLockConfiguration) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ServicePrincipalLockConfiguration) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTokenEncryptionKeyId gets the tokenEncryptionKeyId property value. Locks the tokenEncryptionKeyId property for modification on the service principal.
// returns a *bool when successful
func (m *ServicePrincipalLockConfiguration) GetTokenEncryptionKeyId()(*bool) {
    val, err := m.GetBackingStore().Get("tokenEncryptionKeyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServicePrincipalLockConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allProperties", m.GetAllProperties())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("credentialsWithUsageSign", m.GetCredentialsWithUsageSign())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("credentialsWithUsageVerify", m.GetCredentialsWithUsageVerify())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
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
        err := writer.WriteBoolValue("tokenEncryptionKeyId", m.GetTokenEncryptionKeyId())
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
func (m *ServicePrincipalLockConfiguration) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllProperties sets the allProperties property value. Enables locking all sensitive properties. The sensitive properties are keyCredentials, passwordCredentials, and tokenEncryptionKeyId.
func (m *ServicePrincipalLockConfiguration) SetAllProperties(value *bool)() {
    err := m.GetBackingStore().Set("allProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ServicePrincipalLockConfiguration) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCredentialsWithUsageSign sets the credentialsWithUsageSign property value. Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Sign.
func (m *ServicePrincipalLockConfiguration) SetCredentialsWithUsageSign(value *bool)() {
    err := m.GetBackingStore().Set("credentialsWithUsageSign", value)
    if err != nil {
        panic(err)
    }
}
// SetCredentialsWithUsageVerify sets the credentialsWithUsageVerify property value. Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Verify. This locks OAuth service principals.
func (m *ServicePrincipalLockConfiguration) SetCredentialsWithUsageVerify(value *bool)() {
    err := m.GetBackingStore().Set("credentialsWithUsageVerify", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. Enables or disables service principal lock configuration. To allow the sensitive properties to be updated, update this property to false to disable the lock on the service principal.
func (m *ServicePrincipalLockConfiguration) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ServicePrincipalLockConfiguration) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenEncryptionKeyId sets the tokenEncryptionKeyId property value. Locks the tokenEncryptionKeyId property for modification on the service principal.
func (m *ServicePrincipalLockConfiguration) SetTokenEncryptionKeyId(value *bool)() {
    err := m.GetBackingStore().Set("tokenEncryptionKeyId", value)
    if err != nil {
        panic(err)
    }
}
type ServicePrincipalLockConfigurationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllProperties()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCredentialsWithUsageSign()(*bool)
    GetCredentialsWithUsageVerify()(*bool)
    GetIsEnabled()(*bool)
    GetOdataType()(*string)
    GetTokenEncryptionKeyId()(*bool)
    SetAllProperties(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCredentialsWithUsageSign(value *bool)()
    SetCredentialsWithUsageVerify(value *bool)()
    SetIsEnabled(value *bool)()
    SetOdataType(value *string)()
    SetTokenEncryptionKeyId(value *bool)()
}
