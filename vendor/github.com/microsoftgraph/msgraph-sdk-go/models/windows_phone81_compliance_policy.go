package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsPhone81CompliancePolicy this class contains compliance settings for Windows 8.1 Mobile.
type WindowsPhone81CompliancePolicy struct {
    DeviceCompliancePolicy
}
// NewWindowsPhone81CompliancePolicy instantiates a new WindowsPhone81CompliancePolicy and sets the default values.
func NewWindowsPhone81CompliancePolicy()(*WindowsPhone81CompliancePolicy) {
    m := &WindowsPhone81CompliancePolicy{
        DeviceCompliancePolicy: *NewDeviceCompliancePolicy(),
    }
    odataTypeValue := "#microsoft.graph.windowsPhone81CompliancePolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsPhone81CompliancePolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsPhone81CompliancePolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsPhone81CompliancePolicy(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsPhone81CompliancePolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceCompliancePolicy.GetFieldDeserializers()
    res["osMaximumVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsMaximumVersion(val)
        }
        return nil
    }
    res["osMinimumVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsMinimumVersion(val)
        }
        return nil
    }
    res["passwordBlockSimple"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordBlockSimple(val)
        }
        return nil
    }
    res["passwordExpirationDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordExpirationDays(val)
        }
        return nil
    }
    res["passwordMinimumCharacterSetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordMinimumCharacterSetCount(val)
        }
        return nil
    }
    res["passwordMinimumLength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordMinimumLength(val)
        }
        return nil
    }
    res["passwordMinutesOfInactivityBeforeLock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordMinutesOfInactivityBeforeLock(val)
        }
        return nil
    }
    res["passwordPreviousPasswordBlockCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordPreviousPasswordBlockCount(val)
        }
        return nil
    }
    res["passwordRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordRequired(val)
        }
        return nil
    }
    res["passwordRequiredType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRequiredPasswordType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordRequiredType(val.(*RequiredPasswordType))
        }
        return nil
    }
    res["storageRequireEncryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRequireEncryption(val)
        }
        return nil
    }
    return res
}
// GetOsMaximumVersion gets the osMaximumVersion property value. Maximum Windows Phone version.
// returns a *string when successful
func (m *WindowsPhone81CompliancePolicy) GetOsMaximumVersion()(*string) {
    val, err := m.GetBackingStore().Get("osMaximumVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOsMinimumVersion gets the osMinimumVersion property value. Minimum Windows Phone version.
// returns a *string when successful
func (m *WindowsPhone81CompliancePolicy) GetOsMinimumVersion()(*string) {
    val, err := m.GetBackingStore().Get("osMinimumVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPasswordBlockSimple gets the passwordBlockSimple property value. Whether or not to block syncing the calendar.
// returns a *bool when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordBlockSimple()(*bool) {
    val, err := m.GetBackingStore().Get("passwordBlockSimple")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. Number of days before the password expires.
// returns a *int32 when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordExpirationDays()(*int32) {
    val, err := m.GetBackingStore().Get("passwordExpirationDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumCharacterSetCount gets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
// returns a *int32 when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordMinimumCharacterSetCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumCharacterSetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. Minimum length of passwords.
// returns a *int32 when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordMinimumLength()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinutesOfInactivityBeforeLock gets the passwordMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a password is required.
// returns a *int32 when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordMinutesOfInactivityBeforeLock()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinutesOfInactivityBeforeLock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordPreviousPasswordBlockCount gets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 0 to 24
// returns a *int32 when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordPreviousPasswordBlockCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordPreviousPasswordBlockCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordRequired gets the passwordRequired property value. Whether or not to require a password.
// returns a *bool when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordRequired()(*bool) {
    val, err := m.GetBackingStore().Get("passwordRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordRequiredType gets the passwordRequiredType property value. Possible values of required passwords.
// returns a *RequiredPasswordType when successful
func (m *WindowsPhone81CompliancePolicy) GetPasswordRequiredType()(*RequiredPasswordType) {
    val, err := m.GetBackingStore().Get("passwordRequiredType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RequiredPasswordType)
    }
    return nil
}
// GetStorageRequireEncryption gets the storageRequireEncryption property value. Require encryption on windows phone devices.
// returns a *bool when successful
func (m *WindowsPhone81CompliancePolicy) GetStorageRequireEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("storageRequireEncryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsPhone81CompliancePolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceCompliancePolicy.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("osMaximumVersion", m.GetOsMaximumVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("osMinimumVersion", m.GetOsMinimumVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordBlockSimple", m.GetPasswordBlockSimple())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordExpirationDays", m.GetPasswordExpirationDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinimumCharacterSetCount", m.GetPasswordMinimumCharacterSetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinimumLength", m.GetPasswordMinimumLength())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinutesOfInactivityBeforeLock", m.GetPasswordMinutesOfInactivityBeforeLock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordPreviousPasswordBlockCount", m.GetPasswordPreviousPasswordBlockCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordRequired", m.GetPasswordRequired())
        if err != nil {
            return err
        }
    }
    if m.GetPasswordRequiredType() != nil {
        cast := (*m.GetPasswordRequiredType()).String()
        err = writer.WriteStringValue("passwordRequiredType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRequireEncryption", m.GetStorageRequireEncryption())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOsMaximumVersion sets the osMaximumVersion property value. Maximum Windows Phone version.
func (m *WindowsPhone81CompliancePolicy) SetOsMaximumVersion(value *string)() {
    err := m.GetBackingStore().Set("osMaximumVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetOsMinimumVersion sets the osMinimumVersion property value. Minimum Windows Phone version.
func (m *WindowsPhone81CompliancePolicy) SetOsMinimumVersion(value *string)() {
    err := m.GetBackingStore().Set("osMinimumVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBlockSimple sets the passwordBlockSimple property value. Whether or not to block syncing the calendar.
func (m *WindowsPhone81CompliancePolicy) SetPasswordBlockSimple(value *bool)() {
    err := m.GetBackingStore().Set("passwordBlockSimple", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. Number of days before the password expires.
func (m *WindowsPhone81CompliancePolicy) SetPasswordExpirationDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordExpirationDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumCharacterSetCount sets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *WindowsPhone81CompliancePolicy) SetPasswordMinimumCharacterSetCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumCharacterSetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. Minimum length of passwords.
func (m *WindowsPhone81CompliancePolicy) SetPasswordMinimumLength(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumLength", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinutesOfInactivityBeforeLock sets the passwordMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a password is required.
func (m *WindowsPhone81CompliancePolicy) SetPasswordMinutesOfInactivityBeforeLock(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinutesOfInactivityBeforeLock", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 0 to 24
func (m *WindowsPhone81CompliancePolicy) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordPreviousPasswordBlockCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequired sets the passwordRequired property value. Whether or not to require a password.
func (m *WindowsPhone81CompliancePolicy) SetPasswordRequired(value *bool)() {
    err := m.GetBackingStore().Set("passwordRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Possible values of required passwords.
func (m *WindowsPhone81CompliancePolicy) SetPasswordRequiredType(value *RequiredPasswordType)() {
    err := m.GetBackingStore().Set("passwordRequiredType", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRequireEncryption sets the storageRequireEncryption property value. Require encryption on windows phone devices.
func (m *WindowsPhone81CompliancePolicy) SetStorageRequireEncryption(value *bool)() {
    err := m.GetBackingStore().Set("storageRequireEncryption", value)
    if err != nil {
        panic(err)
    }
}
type WindowsPhone81CompliancePolicyable interface {
    DeviceCompliancePolicyable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOsMaximumVersion()(*string)
    GetOsMinimumVersion()(*string)
    GetPasswordBlockSimple()(*bool)
    GetPasswordExpirationDays()(*int32)
    GetPasswordMinimumCharacterSetCount()(*int32)
    GetPasswordMinimumLength()(*int32)
    GetPasswordMinutesOfInactivityBeforeLock()(*int32)
    GetPasswordPreviousPasswordBlockCount()(*int32)
    GetPasswordRequired()(*bool)
    GetPasswordRequiredType()(*RequiredPasswordType)
    GetStorageRequireEncryption()(*bool)
    SetOsMaximumVersion(value *string)()
    SetOsMinimumVersion(value *string)()
    SetPasswordBlockSimple(value *bool)()
    SetPasswordExpirationDays(value *int32)()
    SetPasswordMinimumCharacterSetCount(value *int32)()
    SetPasswordMinimumLength(value *int32)()
    SetPasswordMinutesOfInactivityBeforeLock(value *int32)()
    SetPasswordPreviousPasswordBlockCount(value *int32)()
    SetPasswordRequired(value *bool)()
    SetPasswordRequiredType(value *RequiredPasswordType)()
    SetStorageRequireEncryption(value *bool)()
}
