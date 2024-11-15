package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AgreementAcceptance struct {
    Entity
}
// NewAgreementAcceptance instantiates a new AgreementAcceptance and sets the default values.
func NewAgreementAcceptance()(*AgreementAcceptance) {
    m := &AgreementAcceptance{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAgreementAcceptanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAgreementAcceptanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementAcceptance(), nil
}
// GetAgreementFileId gets the agreementFileId property value. The identifier of the agreement file accepted by the user.
// returns a *string when successful
func (m *AgreementAcceptance) GetAgreementFileId()(*string) {
    val, err := m.GetBackingStore().Get("agreementFileId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAgreementId gets the agreementId property value. The identifier of the agreement.
// returns a *string when successful
func (m *AgreementAcceptance) GetAgreementId()(*string) {
    val, err := m.GetBackingStore().Get("agreementId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceDisplayName gets the deviceDisplayName property value. The display name of the device used for accepting the agreement.
// returns a *string when successful
func (m *AgreementAcceptance) GetDeviceDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("deviceDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceId gets the deviceId property value. The unique identifier of the device used for accepting the agreement. Supports $filter (eq) and eq for null values.
// returns a *string when successful
func (m *AgreementAcceptance) GetDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("deviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceOSType gets the deviceOSType property value. The operating system used to accept the agreement.
// returns a *string when successful
func (m *AgreementAcceptance) GetDeviceOSType()(*string) {
    val, err := m.GetBackingStore().Get("deviceOSType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceOSVersion gets the deviceOSVersion property value. The operating system version of the device used to accept the agreement.
// returns a *string when successful
func (m *AgreementAcceptance) GetDeviceOSVersion()(*string) {
    val, err := m.GetBackingStore().Get("deviceOSVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExpirationDateTime gets the expirationDateTime property value. The expiration date time of the acceptance. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and eq for null values.
// returns a *Time when successful
func (m *AgreementAcceptance) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("expirationDateTime")
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
func (m *AgreementAcceptance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["agreementFileId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAgreementFileId(val)
        }
        return nil
    }
    res["agreementId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAgreementId(val)
        }
        return nil
    }
    res["deviceDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceDisplayName(val)
        }
        return nil
    }
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
    res["deviceOSType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceOSType(val)
        }
        return nil
    }
    res["deviceOSVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceOSVersion(val)
        }
        return nil
    }
    res["expirationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpirationDateTime(val)
        }
        return nil
    }
    res["recordedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecordedDateTime(val)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAgreementAcceptanceState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*AgreementAcceptanceState))
        }
        return nil
    }
    res["userDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserDisplayName(val)
        }
        return nil
    }
    res["userEmail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserEmail(val)
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
    res["userPrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserPrincipalName(val)
        }
        return nil
    }
    return res
}
// GetRecordedDateTime gets the recordedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *AgreementAcceptance) GetRecordedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("recordedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetState gets the state property value. The state of the agreement acceptance. Possible values are: accepted, declined. Supports $filter (eq).
// returns a *AgreementAcceptanceState when successful
func (m *AgreementAcceptance) GetState()(*AgreementAcceptanceState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AgreementAcceptanceState)
    }
    return nil
}
// GetUserDisplayName gets the userDisplayName property value. Display name of the user when the acceptance was recorded.
// returns a *string when successful
func (m *AgreementAcceptance) GetUserDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("userDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserEmail gets the userEmail property value. Email of the user when the acceptance was recorded.
// returns a *string when successful
func (m *AgreementAcceptance) GetUserEmail()(*string) {
    val, err := m.GetBackingStore().Get("userEmail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. The identifier of the user who accepted the agreement. Supports $filter (eq).
// returns a *string when successful
func (m *AgreementAcceptance) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. UPN of the user when the acceptance was recorded.
// returns a *string when successful
func (m *AgreementAcceptance) GetUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("userPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AgreementAcceptance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("agreementFileId", m.GetAgreementFileId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("agreementId", m.GetAgreementId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceDisplayName", m.GetDeviceDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceOSType", m.GetDeviceOSType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceOSVersion", m.GetDeviceOSVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("expirationDateTime", m.GetExpirationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("recordedDateTime", m.GetRecordedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userDisplayName", m.GetUserDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userEmail", m.GetUserEmail())
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
    {
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAgreementFileId sets the agreementFileId property value. The identifier of the agreement file accepted by the user.
func (m *AgreementAcceptance) SetAgreementFileId(value *string)() {
    err := m.GetBackingStore().Set("agreementFileId", value)
    if err != nil {
        panic(err)
    }
}
// SetAgreementId sets the agreementId property value. The identifier of the agreement.
func (m *AgreementAcceptance) SetAgreementId(value *string)() {
    err := m.GetBackingStore().Set("agreementId", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceDisplayName sets the deviceDisplayName property value. The display name of the device used for accepting the agreement.
func (m *AgreementAcceptance) SetDeviceDisplayName(value *string)() {
    err := m.GetBackingStore().Set("deviceDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceId sets the deviceId property value. The unique identifier of the device used for accepting the agreement. Supports $filter (eq) and eq for null values.
func (m *AgreementAcceptance) SetDeviceId(value *string)() {
    err := m.GetBackingStore().Set("deviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceOSType sets the deviceOSType property value. The operating system used to accept the agreement.
func (m *AgreementAcceptance) SetDeviceOSType(value *string)() {
    err := m.GetBackingStore().Set("deviceOSType", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceOSVersion sets the deviceOSVersion property value. The operating system version of the device used to accept the agreement.
func (m *AgreementAcceptance) SetDeviceOSVersion(value *string)() {
    err := m.GetBackingStore().Set("deviceOSVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetExpirationDateTime sets the expirationDateTime property value. The expiration date time of the acceptance. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and eq for null values.
func (m *AgreementAcceptance) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRecordedDateTime sets the recordedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *AgreementAcceptance) SetRecordedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("recordedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state of the agreement acceptance. Possible values are: accepted, declined. Supports $filter (eq).
func (m *AgreementAcceptance) SetState(value *AgreementAcceptanceState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetUserDisplayName sets the userDisplayName property value. Display name of the user when the acceptance was recorded.
func (m *AgreementAcceptance) SetUserDisplayName(value *string)() {
    err := m.GetBackingStore().Set("userDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserEmail sets the userEmail property value. Email of the user when the acceptance was recorded.
func (m *AgreementAcceptance) SetUserEmail(value *string)() {
    err := m.GetBackingStore().Set("userEmail", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. The identifier of the user who accepted the agreement. Supports $filter (eq).
func (m *AgreementAcceptance) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. UPN of the user when the acceptance was recorded.
func (m *AgreementAcceptance) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
type AgreementAcceptanceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAgreementFileId()(*string)
    GetAgreementId()(*string)
    GetDeviceDisplayName()(*string)
    GetDeviceId()(*string)
    GetDeviceOSType()(*string)
    GetDeviceOSVersion()(*string)
    GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRecordedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetState()(*AgreementAcceptanceState)
    GetUserDisplayName()(*string)
    GetUserEmail()(*string)
    GetUserId()(*string)
    GetUserPrincipalName()(*string)
    SetAgreementFileId(value *string)()
    SetAgreementId(value *string)()
    SetDeviceDisplayName(value *string)()
    SetDeviceId(value *string)()
    SetDeviceOSType(value *string)()
    SetDeviceOSVersion(value *string)()
    SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRecordedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetState(value *AgreementAcceptanceState)()
    SetUserDisplayName(value *string)()
    SetUserEmail(value *string)()
    SetUserId(value *string)()
    SetUserPrincipalName(value *string)()
}
