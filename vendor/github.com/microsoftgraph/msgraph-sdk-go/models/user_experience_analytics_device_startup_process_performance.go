package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsDeviceStartupProcessPerformance the user experience analytics device startup process performance.
type UserExperienceAnalyticsDeviceStartupProcessPerformance struct {
    Entity
}
// NewUserExperienceAnalyticsDeviceStartupProcessPerformance instantiates a new UserExperienceAnalyticsDeviceStartupProcessPerformance and sets the default values.
func NewUserExperienceAnalyticsDeviceStartupProcessPerformance()(*UserExperienceAnalyticsDeviceStartupProcessPerformance) {
    m := &UserExperienceAnalyticsDeviceStartupProcessPerformance{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsDeviceStartupProcessPerformanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsDeviceStartupProcessPerformanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsDeviceStartupProcessPerformance(), nil
}
// GetDeviceCount gets the deviceCount property value. The count of devices which initiated this process on startup. Supports: $filter, $select, $OrderBy. Read-only.
// returns a *int64 when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetDeviceCount()(*int64) {
    val, err := m.GetBackingStore().Get("deviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCount(val)
        }
        return nil
    }
    res["medianImpactInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMedianImpactInMs(val)
        }
        return nil
    }
    res["processName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessName(val)
        }
        return nil
    }
    res["productName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductName(val)
        }
        return nil
    }
    res["publisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisher(val)
        }
        return nil
    }
    res["totalImpactInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalImpactInMs(val)
        }
        return nil
    }
    return res
}
// GetMedianImpactInMs gets the medianImpactInMs property value. The median impact of startup process on device boot time in milliseconds. Supports: $filter, $select, $OrderBy. Read-only.
// returns a *int64 when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetMedianImpactInMs()(*int64) {
    val, err := m.GetBackingStore().Get("medianImpactInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetProcessName gets the processName property value. The name of the startup process. Examples: outlook, excel. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetProcessName()(*string) {
    val, err := m.GetBackingStore().Get("processName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductName gets the productName property value. The product name of the startup process. Examples: Microsoft Outlook, Microsoft Excel. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetProductName()(*string) {
    val, err := m.GetBackingStore().Get("productName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublisher gets the publisher property value. The publisher of the startup process. Examples: Microsoft Corporation, Contoso Corp. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetPublisher()(*string) {
    val, err := m.GetBackingStore().Get("publisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalImpactInMs gets the totalImpactInMs property value. The total impact of startup process on device boot time in milliseconds. Supports: $filter, $select, $OrderBy. Read-only.
// returns a *int64 when successful
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) GetTotalImpactInMs()(*int64) {
    val, err := m.GetBackingStore().Get("totalImpactInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt64Value("deviceCount", m.GetDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("medianImpactInMs", m.GetMedianImpactInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("processName", m.GetProcessName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productName", m.GetProductName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("publisher", m.GetPublisher())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("totalImpactInMs", m.GetTotalImpactInMs())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceCount sets the deviceCount property value. The count of devices which initiated this process on startup. Supports: $filter, $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) SetDeviceCount(value *int64)() {
    err := m.GetBackingStore().Set("deviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetMedianImpactInMs sets the medianImpactInMs property value. The median impact of startup process on device boot time in milliseconds. Supports: $filter, $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) SetMedianImpactInMs(value *int64)() {
    err := m.GetBackingStore().Set("medianImpactInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessName sets the processName property value. The name of the startup process. Examples: outlook, excel. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) SetProcessName(value *string)() {
    err := m.GetBackingStore().Set("processName", value)
    if err != nil {
        panic(err)
    }
}
// SetProductName sets the productName property value. The product name of the startup process. Examples: Microsoft Outlook, Microsoft Excel. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) SetProductName(value *string)() {
    err := m.GetBackingStore().Set("productName", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisher sets the publisher property value. The publisher of the startup process. Examples: Microsoft Corporation, Contoso Corp. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) SetPublisher(value *string)() {
    err := m.GetBackingStore().Set("publisher", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalImpactInMs sets the totalImpactInMs property value. The total impact of startup process on device boot time in milliseconds. Supports: $filter, $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupProcessPerformance) SetTotalImpactInMs(value *int64)() {
    err := m.GetBackingStore().Set("totalImpactInMs", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsDeviceStartupProcessPerformanceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceCount()(*int64)
    GetMedianImpactInMs()(*int64)
    GetProcessName()(*string)
    GetProductName()(*string)
    GetPublisher()(*string)
    GetTotalImpactInMs()(*int64)
    SetDeviceCount(value *int64)()
    SetMedianImpactInMs(value *int64)()
    SetProcessName(value *string)()
    SetProductName(value *string)()
    SetPublisher(value *string)()
    SetTotalImpactInMs(value *int64)()
}
