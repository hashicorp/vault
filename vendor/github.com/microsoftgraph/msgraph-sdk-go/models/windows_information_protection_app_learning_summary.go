package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtectionAppLearningSummary windows Information Protection AppLearning Summary entity.
type WindowsInformationProtectionAppLearningSummary struct {
    Entity
}
// NewWindowsInformationProtectionAppLearningSummary instantiates a new WindowsInformationProtectionAppLearningSummary and sets the default values.
func NewWindowsInformationProtectionAppLearningSummary()(*WindowsInformationProtectionAppLearningSummary) {
    m := &WindowsInformationProtectionAppLearningSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWindowsInformationProtectionAppLearningSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsInformationProtectionAppLearningSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsInformationProtectionAppLearningSummary(), nil
}
// GetApplicationName gets the applicationName property value. Application Name
// returns a *string when successful
func (m *WindowsInformationProtectionAppLearningSummary) GetApplicationName()(*string) {
    val, err := m.GetBackingStore().Get("applicationName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetApplicationType gets the applicationType property value. Possible types of Application
// returns a *ApplicationType when successful
func (m *WindowsInformationProtectionAppLearningSummary) GetApplicationType()(*ApplicationType) {
    val, err := m.GetBackingStore().Get("applicationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ApplicationType)
    }
    return nil
}
// GetDeviceCount gets the deviceCount property value. Device Count
// returns a *int32 when successful
func (m *WindowsInformationProtectionAppLearningSummary) GetDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("deviceCount")
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
func (m *WindowsInformationProtectionAppLearningSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["applicationName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationName(val)
        }
        return nil
    }
    res["applicationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseApplicationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationType(val.(*ApplicationType))
        }
        return nil
    }
    res["deviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCount(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WindowsInformationProtectionAppLearningSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("applicationName", m.GetApplicationName())
        if err != nil {
            return err
        }
    }
    if m.GetApplicationType() != nil {
        cast := (*m.GetApplicationType()).String()
        err = writer.WriteStringValue("applicationType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deviceCount", m.GetDeviceCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationName sets the applicationName property value. Application Name
func (m *WindowsInformationProtectionAppLearningSummary) SetApplicationName(value *string)() {
    err := m.GetBackingStore().Set("applicationName", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationType sets the applicationType property value. Possible types of Application
func (m *WindowsInformationProtectionAppLearningSummary) SetApplicationType(value *ApplicationType)() {
    err := m.GetBackingStore().Set("applicationType", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCount sets the deviceCount property value. Device Count
func (m *WindowsInformationProtectionAppLearningSummary) SetDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("deviceCount", value)
    if err != nil {
        panic(err)
    }
}
type WindowsInformationProtectionAppLearningSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationName()(*string)
    GetApplicationType()(*ApplicationType)
    GetDeviceCount()(*int32)
    SetApplicationName(value *string)()
    SetApplicationType(value *ApplicationType)()
    SetDeviceCount(value *int32)()
}
