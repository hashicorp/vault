package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EnrollmentTroubleshootingEvent event representing an enrollment failure.
type EnrollmentTroubleshootingEvent struct {
    DeviceManagementTroubleshootingEvent
}
// NewEnrollmentTroubleshootingEvent instantiates a new EnrollmentTroubleshootingEvent and sets the default values.
func NewEnrollmentTroubleshootingEvent()(*EnrollmentTroubleshootingEvent) {
    m := &EnrollmentTroubleshootingEvent{
        DeviceManagementTroubleshootingEvent: *NewDeviceManagementTroubleshootingEvent(),
    }
    return m
}
// CreateEnrollmentTroubleshootingEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEnrollmentTroubleshootingEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEnrollmentTroubleshootingEvent(), nil
}
// GetDeviceId gets the deviceId property value. Azure AD device identifier.
// returns a *string when successful
func (m *EnrollmentTroubleshootingEvent) GetDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("deviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnrollmentType gets the enrollmentType property value. Possible ways of adding a mobile device to management.
// returns a *DeviceEnrollmentType when successful
func (m *EnrollmentTroubleshootingEvent) GetEnrollmentType()(*DeviceEnrollmentType) {
    val, err := m.GetBackingStore().Get("enrollmentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceEnrollmentType)
    }
    return nil
}
// GetFailureCategory gets the failureCategory property value. Top level failure categories for enrollment.
// returns a *DeviceEnrollmentFailureReason when successful
func (m *EnrollmentTroubleshootingEvent) GetFailureCategory()(*DeviceEnrollmentFailureReason) {
    val, err := m.GetBackingStore().Get("failureCategory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceEnrollmentFailureReason)
    }
    return nil
}
// GetFailureReason gets the failureReason property value. Detailed failure reason.
// returns a *string when successful
func (m *EnrollmentTroubleshootingEvent) GetFailureReason()(*string) {
    val, err := m.GetBackingStore().Get("failureReason")
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
func (m *EnrollmentTroubleshootingEvent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceManagementTroubleshootingEvent.GetFieldDeserializers()
    res["deviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceId(val)
        }
        return nil
    }
    res["enrollmentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceEnrollmentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnrollmentType(val.(*DeviceEnrollmentType))
        }
        return nil
    }
    res["failureCategory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceEnrollmentFailureReason)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailureCategory(val.(*DeviceEnrollmentFailureReason))
        }
        return nil
    }
    res["failureReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailureReason(val)
        }
        return nil
    }
    res["managedDeviceIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedDeviceIdentifier(val)
        }
        return nil
    }
    res["operatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperatingSystem(val)
        }
        return nil
    }
    res["osVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsVersion(val)
        }
        return nil
    }
    res["userId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserId(val)
        }
        return nil
    }
    return res
}
// GetManagedDeviceIdentifier gets the managedDeviceIdentifier property value. Device identifier created or collected by Intune.
// returns a *string when successful
func (m *EnrollmentTroubleshootingEvent) GetManagedDeviceIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("managedDeviceIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperatingSystem gets the operatingSystem property value. Operating System.
// returns a *string when successful
func (m *EnrollmentTroubleshootingEvent) GetOperatingSystem()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOsVersion gets the osVersion property value. OS Version.
// returns a *string when successful
func (m *EnrollmentTroubleshootingEvent) GetOsVersion()(*string) {
    val, err := m.GetBackingStore().Get("osVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. Identifier for the user that tried to enroll the device.
// returns a *string when successful
func (m *EnrollmentTroubleshootingEvent) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EnrollmentTroubleshootingEvent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceManagementTroubleshootingEvent.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    if m.GetEnrollmentType() != nil {
        cast := (*m.GetEnrollmentType()).String()
        err = writer.WriteStringValue("enrollmentType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFailureCategory() != nil {
        cast := (*m.GetFailureCategory()).String()
        err = writer.WriteStringValue("failureCategory", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("failureReason", m.GetFailureReason())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("managedDeviceIdentifier", m.GetManagedDeviceIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("operatingSystem", m.GetOperatingSystem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("osVersion", m.GetOsVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceId sets the deviceId property value. Azure AD device identifier.
func (m *EnrollmentTroubleshootingEvent) SetDeviceId(value *string)() {
    err := m.GetBackingStore().Set("deviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetEnrollmentType sets the enrollmentType property value. Possible ways of adding a mobile device to management.
func (m *EnrollmentTroubleshootingEvent) SetEnrollmentType(value *DeviceEnrollmentType)() {
    err := m.GetBackingStore().Set("enrollmentType", value)
    if err != nil {
        panic(err)
    }
}
// SetFailureCategory sets the failureCategory property value. Top level failure categories for enrollment.
func (m *EnrollmentTroubleshootingEvent) SetFailureCategory(value *DeviceEnrollmentFailureReason)() {
    err := m.GetBackingStore().Set("failureCategory", value)
    if err != nil {
        panic(err)
    }
}
// SetFailureReason sets the failureReason property value. Detailed failure reason.
func (m *EnrollmentTroubleshootingEvent) SetFailureReason(value *string)() {
    err := m.GetBackingStore().Set("failureReason", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedDeviceIdentifier sets the managedDeviceIdentifier property value. Device identifier created or collected by Intune.
func (m *EnrollmentTroubleshootingEvent) SetManagedDeviceIdentifier(value *string)() {
    err := m.GetBackingStore().Set("managedDeviceIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystem sets the operatingSystem property value. Operating System.
func (m *EnrollmentTroubleshootingEvent) SetOperatingSystem(value *string)() {
    err := m.GetBackingStore().Set("operatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetOsVersion sets the osVersion property value. OS Version.
func (m *EnrollmentTroubleshootingEvent) SetOsVersion(value *string)() {
    err := m.GetBackingStore().Set("osVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. Identifier for the user that tried to enroll the device.
func (m *EnrollmentTroubleshootingEvent) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
type EnrollmentTroubleshootingEventable interface {
    DeviceManagementTroubleshootingEventable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceId()(*string)
    GetEnrollmentType()(*DeviceEnrollmentType)
    GetFailureCategory()(*DeviceEnrollmentFailureReason)
    GetFailureReason()(*string)
    GetManagedDeviceIdentifier()(*string)
    GetOperatingSystem()(*string)
    GetOsVersion()(*string)
    GetUserId()(*string)
    SetDeviceId(value *string)()
    SetEnrollmentType(value *DeviceEnrollmentType)()
    SetFailureCategory(value *DeviceEnrollmentFailureReason)()
    SetFailureReason(value *string)()
    SetManagedDeviceIdentifier(value *string)()
    SetOperatingSystem(value *string)()
    SetOsVersion(value *string)()
    SetUserId(value *string)()
}
