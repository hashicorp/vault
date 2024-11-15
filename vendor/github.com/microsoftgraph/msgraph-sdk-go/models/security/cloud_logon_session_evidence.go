package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CloudLogonSessionEvidence struct {
    AlertEvidence
}
// NewCloudLogonSessionEvidence instantiates a new CloudLogonSessionEvidence and sets the default values.
func NewCloudLogonSessionEvidence()(*CloudLogonSessionEvidence) {
    m := &CloudLogonSessionEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.cloudLogonSessionEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCloudLogonSessionEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCloudLogonSessionEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudLogonSessionEvidence(), nil
}
// GetAccount gets the account property value. The account associated with the sign-in session.
// returns a UserEvidenceable when successful
func (m *CloudLogonSessionEvidence) GetAccount()(UserEvidenceable) {
    val, err := m.GetBackingStore().Get("account")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserEvidenceable)
    }
    return nil
}
// GetBrowser gets the browser property value. The browser that is used for the sign-in, if known.
// returns a *string when successful
func (m *CloudLogonSessionEvidence) GetBrowser()(*string) {
    val, err := m.GetBackingStore().Get("browser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. The friendly name of the device, if known.
// returns a *string when successful
func (m *CloudLogonSessionEvidence) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
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
func (m *CloudLogonSessionEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["account"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccount(val.(UserEvidenceable))
        }
        return nil
    }
    res["browser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowser(val)
        }
        return nil
    }
    res["deviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceName(val)
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
    res["previousLogonDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreviousLogonDateTime(val)
        }
        return nil
    }
    res["protocol"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtocol(val)
        }
        return nil
    }
    res["sessionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSessionId(val)
        }
        return nil
    }
    res["startUtcDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartUtcDateTime(val)
        }
        return nil
    }
    res["userAgent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserAgent(val)
        }
        return nil
    }
    return res
}
// GetOperatingSystem gets the operatingSystem property value. The operating system that the device is running, if known.
// returns a *string when successful
func (m *CloudLogonSessionEvidence) GetOperatingSystem()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreviousLogonDateTime gets the previousLogonDateTime property value. The previous sign-in time for this account, if known.
// returns a *Time when successful
func (m *CloudLogonSessionEvidence) GetPreviousLogonDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("previousLogonDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetProtocol gets the protocol property value. The authentication protocol that is used in this session, if known.
// returns a *string when successful
func (m *CloudLogonSessionEvidence) GetProtocol()(*string) {
    val, err := m.GetBackingStore().Get("protocol")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSessionId gets the sessionId property value. The session ID for the account reported in the alert.
// returns a *string when successful
func (m *CloudLogonSessionEvidence) GetSessionId()(*string) {
    val, err := m.GetBackingStore().Get("sessionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartUtcDateTime gets the startUtcDateTime property value. The session start time, if known.
// returns a *Time when successful
func (m *CloudLogonSessionEvidence) GetStartUtcDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startUtcDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetUserAgent gets the userAgent property value. The user agent that is used for the sign-in, if known.
// returns a *string when successful
func (m *CloudLogonSessionEvidence) GetUserAgent()(*string) {
    val, err := m.GetBackingStore().Get("userAgent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CloudLogonSessionEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("account", m.GetAccount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("browser", m.GetBrowser())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceName", m.GetDeviceName())
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
        err = writer.WriteTimeValue("previousLogonDateTime", m.GetPreviousLogonDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("protocol", m.GetProtocol())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sessionId", m.GetSessionId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startUtcDateTime", m.GetStartUtcDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userAgent", m.GetUserAgent())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccount sets the account property value. The account associated with the sign-in session.
func (m *CloudLogonSessionEvidence) SetAccount(value UserEvidenceable)() {
    err := m.GetBackingStore().Set("account", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowser sets the browser property value. The browser that is used for the sign-in, if known.
func (m *CloudLogonSessionEvidence) SetBrowser(value *string)() {
    err := m.GetBackingStore().Set("browser", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. The friendly name of the device, if known.
func (m *CloudLogonSessionEvidence) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystem sets the operatingSystem property value. The operating system that the device is running, if known.
func (m *CloudLogonSessionEvidence) SetOperatingSystem(value *string)() {
    err := m.GetBackingStore().Set("operatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviousLogonDateTime sets the previousLogonDateTime property value. The previous sign-in time for this account, if known.
func (m *CloudLogonSessionEvidence) SetPreviousLogonDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("previousLogonDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetProtocol sets the protocol property value. The authentication protocol that is used in this session, if known.
func (m *CloudLogonSessionEvidence) SetProtocol(value *string)() {
    err := m.GetBackingStore().Set("protocol", value)
    if err != nil {
        panic(err)
    }
}
// SetSessionId sets the sessionId property value. The session ID for the account reported in the alert.
func (m *CloudLogonSessionEvidence) SetSessionId(value *string)() {
    err := m.GetBackingStore().Set("sessionId", value)
    if err != nil {
        panic(err)
    }
}
// SetStartUtcDateTime sets the startUtcDateTime property value. The session start time, if known.
func (m *CloudLogonSessionEvidence) SetStartUtcDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startUtcDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAgent sets the userAgent property value. The user agent that is used for the sign-in, if known.
func (m *CloudLogonSessionEvidence) SetUserAgent(value *string)() {
    err := m.GetBackingStore().Set("userAgent", value)
    if err != nil {
        panic(err)
    }
}
type CloudLogonSessionEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccount()(UserEvidenceable)
    GetBrowser()(*string)
    GetDeviceName()(*string)
    GetOperatingSystem()(*string)
    GetPreviousLogonDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetProtocol()(*string)
    GetSessionId()(*string)
    GetStartUtcDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetUserAgent()(*string)
    SetAccount(value UserEvidenceable)()
    SetBrowser(value *string)()
    SetDeviceName(value *string)()
    SetOperatingSystem(value *string)()
    SetPreviousLogonDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetProtocol(value *string)()
    SetSessionId(value *string)()
    SetStartUtcDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetUserAgent(value *string)()
}
