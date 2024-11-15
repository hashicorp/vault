package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WindowsHelloForBusinessAuthenticationMethod struct {
    AuthenticationMethod
}
// NewWindowsHelloForBusinessAuthenticationMethod instantiates a new WindowsHelloForBusinessAuthenticationMethod and sets the default values.
func NewWindowsHelloForBusinessAuthenticationMethod()(*WindowsHelloForBusinessAuthenticationMethod) {
    m := &WindowsHelloForBusinessAuthenticationMethod{
        AuthenticationMethod: *NewAuthenticationMethod(),
    }
    odataTypeValue := "#microsoft.graph.windowsHelloForBusinessAuthenticationMethod"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsHelloForBusinessAuthenticationMethodFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsHelloForBusinessAuthenticationMethodFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsHelloForBusinessAuthenticationMethod(), nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time that this Windows Hello for Business key was registered.
// returns a *Time when successful
func (m *WindowsHelloForBusinessAuthenticationMethod) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDevice gets the device property value. The registered device on which this Windows Hello for Business key resides. Supports $expand. When you get a user's Windows Hello for Business registration information, this property is returned only on a single GET and when you specify ?$expand. For example, GET /users/admin@contoso.com/authentication/windowsHelloForBusinessMethods/_jpuR-TGZtk6aQCLF3BQjA2?$expand=device.
// returns a Deviceable when successful
func (m *WindowsHelloForBusinessAuthenticationMethod) GetDevice()(Deviceable) {
    val, err := m.GetBackingStore().Get("device")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Deviceable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the device on which Windows Hello for Business is registered
// returns a *string when successful
func (m *WindowsHelloForBusinessAuthenticationMethod) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *WindowsHelloForBusinessAuthenticationMethod) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethod.GetFieldDeserializers()
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["device"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevice(val.(Deviceable))
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
    res["keyStrength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationMethodKeyStrength)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyStrength(val.(*AuthenticationMethodKeyStrength))
        }
        return nil
    }
    return res
}
// GetKeyStrength gets the keyStrength property value. Key strength of this Windows Hello for Business key. Possible values are: normal, weak, unknown.
// returns a *AuthenticationMethodKeyStrength when successful
func (m *WindowsHelloForBusinessAuthenticationMethod) GetKeyStrength()(*AuthenticationMethodKeyStrength) {
    val, err := m.GetBackingStore().Get("keyStrength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationMethodKeyStrength)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsHelloForBusinessAuthenticationMethod) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethod.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("device", m.GetDevice())
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
    if m.GetKeyStrength() != nil {
        cast := (*m.GetKeyStrength()).String()
        err = writer.WriteStringValue("keyStrength", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time that this Windows Hello for Business key was registered.
func (m *WindowsHelloForBusinessAuthenticationMethod) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDevice sets the device property value. The registered device on which this Windows Hello for Business key resides. Supports $expand. When you get a user's Windows Hello for Business registration information, this property is returned only on a single GET and when you specify ?$expand. For example, GET /users/admin@contoso.com/authentication/windowsHelloForBusinessMethods/_jpuR-TGZtk6aQCLF3BQjA2?$expand=device.
func (m *WindowsHelloForBusinessAuthenticationMethod) SetDevice(value Deviceable)() {
    err := m.GetBackingStore().Set("device", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the device on which Windows Hello for Business is registered
func (m *WindowsHelloForBusinessAuthenticationMethod) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyStrength sets the keyStrength property value. Key strength of this Windows Hello for Business key. Possible values are: normal, weak, unknown.
func (m *WindowsHelloForBusinessAuthenticationMethod) SetKeyStrength(value *AuthenticationMethodKeyStrength)() {
    err := m.GetBackingStore().Set("keyStrength", value)
    if err != nil {
        panic(err)
    }
}
type WindowsHelloForBusinessAuthenticationMethodable interface {
    AuthenticationMethodable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDevice()(Deviceable)
    GetDisplayName()(*string)
    GetKeyStrength()(*AuthenticationMethodKeyStrength)
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDevice(value Deviceable)()
    SetDisplayName(value *string)()
    SetKeyStrength(value *AuthenticationMethodKeyStrength)()
}
