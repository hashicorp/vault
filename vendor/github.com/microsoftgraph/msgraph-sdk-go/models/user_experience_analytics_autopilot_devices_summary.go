package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// UserExperienceAnalyticsAutopilotDevicesSummary the user experience analytics summary of Devices not windows autopilot ready.
type UserExperienceAnalyticsAutopilotDevicesSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUserExperienceAnalyticsAutopilotDevicesSummary instantiates a new UserExperienceAnalyticsAutopilotDevicesSummary and sets the default values.
func NewUserExperienceAnalyticsAutopilotDevicesSummary()(*UserExperienceAnalyticsAutopilotDevicesSummary) {
    m := &UserExperienceAnalyticsAutopilotDevicesSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUserExperienceAnalyticsAutopilotDevicesSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsAutopilotDevicesSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsAutopilotDevicesSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetAdditionalData()(map[string]any) {
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
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDevicesNotAutopilotRegistered gets the devicesNotAutopilotRegistered property value. The count of intune devices that are not autopilot registerd. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetDevicesNotAutopilotRegistered()(*int32) {
    val, err := m.GetBackingStore().Get("devicesNotAutopilotRegistered")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDevicesWithoutAutopilotProfileAssigned gets the devicesWithoutAutopilotProfileAssigned property value. The count of intune devices not autopilot profile assigned. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetDevicesWithoutAutopilotProfileAssigned()(*int32) {
    val, err := m.GetBackingStore().Get("devicesWithoutAutopilotProfileAssigned")
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
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["devicesNotAutopilotRegistered"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevicesNotAutopilotRegistered(val)
        }
        return nil
    }
    res["devicesWithoutAutopilotProfileAssigned"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevicesWithoutAutopilotProfileAssigned(val)
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
    res["totalWindows10DevicesWithoutTenantAttached"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalWindows10DevicesWithoutTenantAttached(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalWindows10DevicesWithoutTenantAttached gets the totalWindows10DevicesWithoutTenantAttached property value. The count of windows 10 devices that are Intune and co-managed. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) GetTotalWindows10DevicesWithoutTenantAttached()(*int32) {
    val, err := m.GetBackingStore().Get("totalWindows10DevicesWithoutTenantAttached")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("devicesNotAutopilotRegistered", m.GetDevicesNotAutopilotRegistered())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("devicesWithoutAutopilotProfileAssigned", m.GetDevicesWithoutAutopilotProfileAssigned())
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
        err := writer.WriteInt32Value("totalWindows10DevicesWithoutTenantAttached", m.GetTotalWindows10DevicesWithoutTenantAttached())
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
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDevicesNotAutopilotRegistered sets the devicesNotAutopilotRegistered property value. The count of intune devices that are not autopilot registerd. Read-only.
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) SetDevicesNotAutopilotRegistered(value *int32)() {
    err := m.GetBackingStore().Set("devicesNotAutopilotRegistered", value)
    if err != nil {
        panic(err)
    }
}
// SetDevicesWithoutAutopilotProfileAssigned sets the devicesWithoutAutopilotProfileAssigned property value. The count of intune devices not autopilot profile assigned. Read-only.
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) SetDevicesWithoutAutopilotProfileAssigned(value *int32)() {
    err := m.GetBackingStore().Set("devicesWithoutAutopilotProfileAssigned", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalWindows10DevicesWithoutTenantAttached sets the totalWindows10DevicesWithoutTenantAttached property value. The count of windows 10 devices that are Intune and co-managed. Read-only.
func (m *UserExperienceAnalyticsAutopilotDevicesSummary) SetTotalWindows10DevicesWithoutTenantAttached(value *int32)() {
    err := m.GetBackingStore().Set("totalWindows10DevicesWithoutTenantAttached", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsAutopilotDevicesSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDevicesNotAutopilotRegistered()(*int32)
    GetDevicesWithoutAutopilotProfileAssigned()(*int32)
    GetOdataType()(*string)
    GetTotalWindows10DevicesWithoutTenantAttached()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDevicesNotAutopilotRegistered(value *int32)()
    SetDevicesWithoutAutopilotProfileAssigned(value *int32)()
    SetOdataType(value *string)()
    SetTotalWindows10DevicesWithoutTenantAttached(value *int32)()
}
