package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkforceIntegration struct {
    ChangeTrackedEntity
}
// NewWorkforceIntegration instantiates a new WorkforceIntegration and sets the default values.
func NewWorkforceIntegration()(*WorkforceIntegration) {
    m := &WorkforceIntegration{
        ChangeTrackedEntity: *NewChangeTrackedEntity(),
    }
    odataTypeValue := "#microsoft.graph.workforceIntegration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWorkforceIntegrationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkforceIntegrationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkforceIntegration(), nil
}
// GetApiVersion gets the apiVersion property value. API version for the call back URL. Start with 1.
// returns a *int32 when successful
func (m *WorkforceIntegration) GetApiVersion()(*int32) {
    val, err := m.GetBackingStore().Get("apiVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of the workforce integration.
// returns a *string when successful
func (m *WorkforceIntegration) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEncryption gets the encryption property value. The workforce integration encryption resource.
// returns a WorkforceIntegrationEncryptionable when successful
func (m *WorkforceIntegration) GetEncryption()(WorkforceIntegrationEncryptionable) {
    val, err := m.GetBackingStore().Get("encryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkforceIntegrationEncryptionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkforceIntegration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ChangeTrackedEntity.GetFieldDeserializers()
    res["apiVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApiVersion(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["encryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkforceIntegrationEncryptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEncryption(val.(WorkforceIntegrationEncryptionable))
        }
        return nil
    }
    res["isActive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsActive(val)
        }
        return nil
    }
    res["supportedEntities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWorkforceIntegrationSupportedEntities)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSupportedEntities(val.(*WorkforceIntegrationSupportedEntities))
        }
        return nil
    }
    res["url"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUrl(val)
        }
        return nil
    }
    return res
}
// GetIsActive gets the isActive property value. Indicates whether this workforce integration is currently active and available.
// returns a *bool when successful
func (m *WorkforceIntegration) GetIsActive()(*bool) {
    val, err := m.GetBackingStore().Get("isActive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSupportedEntities gets the supportedEntities property value. The Shifts entities supported for synchronous change notifications. Shifts will make a call back to the url provided on client changes on those entities added here. By default, no entities are supported for change notifications. Possible values are: none, shift, swapRequest, userShiftPreferences, openshift, openShiftRequest, offerShiftRequest, unknownFutureValue.
// returns a *WorkforceIntegrationSupportedEntities when successful
func (m *WorkforceIntegration) GetSupportedEntities()(*WorkforceIntegrationSupportedEntities) {
    val, err := m.GetBackingStore().Get("supportedEntities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WorkforceIntegrationSupportedEntities)
    }
    return nil
}
// GetUrl gets the url property value. Workforce Integration URL for callbacks from the Shifts service.
// returns a *string when successful
func (m *WorkforceIntegration) GetUrl()(*string) {
    val, err := m.GetBackingStore().Get("url")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkforceIntegration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ChangeTrackedEntity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("apiVersion", m.GetApiVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("encryption", m.GetEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isActive", m.GetIsActive())
        if err != nil {
            return err
        }
    }
    if m.GetSupportedEntities() != nil {
        cast := (*m.GetSupportedEntities()).String()
        err = writer.WriteStringValue("supportedEntities", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("url", m.GetUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApiVersion sets the apiVersion property value. API version for the call back URL. Start with 1.
func (m *WorkforceIntegration) SetApiVersion(value *int32)() {
    err := m.GetBackingStore().Set("apiVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of the workforce integration.
func (m *WorkforceIntegration) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEncryption sets the encryption property value. The workforce integration encryption resource.
func (m *WorkforceIntegration) SetEncryption(value WorkforceIntegrationEncryptionable)() {
    err := m.GetBackingStore().Set("encryption", value)
    if err != nil {
        panic(err)
    }
}
// SetIsActive sets the isActive property value. Indicates whether this workforce integration is currently active and available.
func (m *WorkforceIntegration) SetIsActive(value *bool)() {
    err := m.GetBackingStore().Set("isActive", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportedEntities sets the supportedEntities property value. The Shifts entities supported for synchronous change notifications. Shifts will make a call back to the url provided on client changes on those entities added here. By default, no entities are supported for change notifications. Possible values are: none, shift, swapRequest, userShiftPreferences, openshift, openShiftRequest, offerShiftRequest, unknownFutureValue.
func (m *WorkforceIntegration) SetSupportedEntities(value *WorkforceIntegrationSupportedEntities)() {
    err := m.GetBackingStore().Set("supportedEntities", value)
    if err != nil {
        panic(err)
    }
}
// SetUrl sets the url property value. Workforce Integration URL for callbacks from the Shifts service.
func (m *WorkforceIntegration) SetUrl(value *string)() {
    err := m.GetBackingStore().Set("url", value)
    if err != nil {
        panic(err)
    }
}
type WorkforceIntegrationable interface {
    ChangeTrackedEntityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApiVersion()(*int32)
    GetDisplayName()(*string)
    GetEncryption()(WorkforceIntegrationEncryptionable)
    GetIsActive()(*bool)
    GetSupportedEntities()(*WorkforceIntegrationSupportedEntities)
    GetUrl()(*string)
    SetApiVersion(value *int32)()
    SetDisplayName(value *string)()
    SetEncryption(value WorkforceIntegrationEncryptionable)()
    SetIsActive(value *bool)()
    SetSupportedEntities(value *WorkforceIntegrationSupportedEntities)()
    SetUrl(value *string)()
}
