package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// UserExperienceAnalyticsCloudManagementDevicesSummary the user experience work from anywhere Cloud management devices summary.
type UserExperienceAnalyticsCloudManagementDevicesSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUserExperienceAnalyticsCloudManagementDevicesSummary instantiates a new UserExperienceAnalyticsCloudManagementDevicesSummary and sets the default values.
func NewUserExperienceAnalyticsCloudManagementDevicesSummary()(*UserExperienceAnalyticsCloudManagementDevicesSummary) {
    m := &UserExperienceAnalyticsCloudManagementDevicesSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUserExperienceAnalyticsCloudManagementDevicesSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsCloudManagementDevicesSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsCloudManagementDevicesSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetAdditionalData()(map[string]any) {
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
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCoManagedDeviceCount gets the coManagedDeviceCount property value. Total number of  co-managed devices. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetCoManagedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("coManagedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["coManagedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCoManagedDeviceCount(val)
        }
        return nil
    }
    res["intuneDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIntuneDeviceCount(val)
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
    res["tenantAttachDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantAttachDeviceCount(val)
        }
        return nil
    }
    return res
}
// GetIntuneDeviceCount gets the intuneDeviceCount property value. The count of intune devices that are not autopilot registerd. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetIntuneDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("intuneDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTenantAttachDeviceCount gets the tenantAttachDeviceCount property value. Total count of tenant attach devices. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) GetTenantAttachDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("tenantAttachDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("coManagedDeviceCount", m.GetCoManagedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("intuneDeviceCount", m.GetIntuneDeviceCount())
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
        err := writer.WriteInt32Value("tenantAttachDeviceCount", m.GetTenantAttachDeviceCount())
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
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCoManagedDeviceCount sets the coManagedDeviceCount property value. Total number of  co-managed devices. Read-only.
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) SetCoManagedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("coManagedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetIntuneDeviceCount sets the intuneDeviceCount property value. The count of intune devices that are not autopilot registerd. Read-only.
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) SetIntuneDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("intuneDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantAttachDeviceCount sets the tenantAttachDeviceCount property value. Total count of tenant attach devices. Read-only.
func (m *UserExperienceAnalyticsCloudManagementDevicesSummary) SetTenantAttachDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("tenantAttachDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsCloudManagementDevicesSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCoManagedDeviceCount()(*int32)
    GetIntuneDeviceCount()(*int32)
    GetOdataType()(*string)
    GetTenantAttachDeviceCount()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCoManagedDeviceCount(value *int32)()
    SetIntuneDeviceCount(value *int32)()
    SetOdataType(value *string)()
    SetTenantAttachDeviceCount(value *int32)()
}
