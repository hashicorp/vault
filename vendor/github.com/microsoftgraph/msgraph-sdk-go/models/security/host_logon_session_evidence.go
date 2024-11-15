package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type HostLogonSessionEvidence struct {
    AlertEvidence
}
// NewHostLogonSessionEvidence instantiates a new HostLogonSessionEvidence and sets the default values.
func NewHostLogonSessionEvidence()(*HostLogonSessionEvidence) {
    m := &HostLogonSessionEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.hostLogonSessionEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateHostLogonSessionEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHostLogonSessionEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHostLogonSessionEvidence(), nil
}
// GetAccount gets the account property value. The account property
// returns a UserEvidenceable when successful
func (m *HostLogonSessionEvidence) GetAccount()(UserEvidenceable) {
    val, err := m.GetBackingStore().Get("account")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserEvidenceable)
    }
    return nil
}
// GetEndUtcDateTime gets the endUtcDateTime property value. The endUtcDateTime property
// returns a *Time when successful
func (m *HostLogonSessionEvidence) GetEndUtcDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("endUtcDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *HostLogonSessionEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["endUtcDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndUtcDateTime(val)
        }
        return nil
    }
    res["host"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHost(val.(DeviceEvidenceable))
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
    return res
}
// GetHost gets the host property value. The host property
// returns a DeviceEvidenceable when successful
func (m *HostLogonSessionEvidence) GetHost()(DeviceEvidenceable) {
    val, err := m.GetBackingStore().Get("host")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceEvidenceable)
    }
    return nil
}
// GetSessionId gets the sessionId property value. The sessionId property
// returns a *string when successful
func (m *HostLogonSessionEvidence) GetSessionId()(*string) {
    val, err := m.GetBackingStore().Get("sessionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartUtcDateTime gets the startUtcDateTime property value. The startUtcDateTime property
// returns a *Time when successful
func (m *HostLogonSessionEvidence) GetStartUtcDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startUtcDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *HostLogonSessionEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteTimeValue("endUtcDateTime", m.GetEndUtcDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("host", m.GetHost())
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
    return nil
}
// SetAccount sets the account property value. The account property
func (m *HostLogonSessionEvidence) SetAccount(value UserEvidenceable)() {
    err := m.GetBackingStore().Set("account", value)
    if err != nil {
        panic(err)
    }
}
// SetEndUtcDateTime sets the endUtcDateTime property value. The endUtcDateTime property
func (m *HostLogonSessionEvidence) SetEndUtcDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("endUtcDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetHost sets the host property value. The host property
func (m *HostLogonSessionEvidence) SetHost(value DeviceEvidenceable)() {
    err := m.GetBackingStore().Set("host", value)
    if err != nil {
        panic(err)
    }
}
// SetSessionId sets the sessionId property value. The sessionId property
func (m *HostLogonSessionEvidence) SetSessionId(value *string)() {
    err := m.GetBackingStore().Set("sessionId", value)
    if err != nil {
        panic(err)
    }
}
// SetStartUtcDateTime sets the startUtcDateTime property value. The startUtcDateTime property
func (m *HostLogonSessionEvidence) SetStartUtcDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startUtcDateTime", value)
    if err != nil {
        panic(err)
    }
}
type HostLogonSessionEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccount()(UserEvidenceable)
    GetEndUtcDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetHost()(DeviceEvidenceable)
    GetSessionId()(*string)
    GetStartUtcDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetAccount(value UserEvidenceable)()
    SetEndUtcDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetHost(value DeviceEvidenceable)()
    SetSessionId(value *string)()
    SetStartUtcDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
