package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ApplePushNotificationCertificate apple push notification certificate.
type ApplePushNotificationCertificate struct {
    Entity
}
// NewApplePushNotificationCertificate instantiates a new ApplePushNotificationCertificate and sets the default values.
func NewApplePushNotificationCertificate()(*ApplePushNotificationCertificate) {
    m := &ApplePushNotificationCertificate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateApplePushNotificationCertificateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateApplePushNotificationCertificateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewApplePushNotificationCertificate(), nil
}
// GetAppleIdentifier gets the appleIdentifier property value. Apple Id of the account used to create the MDM push certificate.
// returns a *string when successful
func (m *ApplePushNotificationCertificate) GetAppleIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("appleIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCertificate gets the certificate property value. Not yet documented
// returns a *string when successful
func (m *ApplePushNotificationCertificate) GetCertificate()(*string) {
    val, err := m.GetBackingStore().Get("certificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCertificateSerialNumber gets the certificateSerialNumber property value. Certificate serial number. This property is read-only.
// returns a *string when successful
func (m *ApplePushNotificationCertificate) GetCertificateSerialNumber()(*string) {
    val, err := m.GetBackingStore().Get("certificateSerialNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCertificateUploadFailureReason gets the certificateUploadFailureReason property value. The reason the certificate upload failed.
// returns a *string when successful
func (m *ApplePushNotificationCertificate) GetCertificateUploadFailureReason()(*string) {
    val, err := m.GetBackingStore().Get("certificateUploadFailureReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCertificateUploadStatus gets the certificateUploadStatus property value. The certificate upload status.
// returns a *string when successful
func (m *ApplePushNotificationCertificate) GetCertificateUploadStatus()(*string) {
    val, err := m.GetBackingStore().Get("certificateUploadStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExpirationDateTime gets the expirationDateTime property value. The expiration date and time for Apple push notification certificate.
// returns a *Time when successful
func (m *ApplePushNotificationCertificate) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
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
func (m *ApplePushNotificationCertificate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appleIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppleIdentifier(val)
        }
        return nil
    }
    res["certificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificate(val)
        }
        return nil
    }
    res["certificateSerialNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificateSerialNumber(val)
        }
        return nil
    }
    res["certificateUploadFailureReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificateUploadFailureReason(val)
        }
        return nil
    }
    res["certificateUploadStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificateUploadStatus(val)
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
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["topicIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTopicIdentifier(val)
        }
        return nil
    }
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Last modified date and time for Apple push notification certificate.
// returns a *Time when successful
func (m *ApplePushNotificationCertificate) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTopicIdentifier gets the topicIdentifier property value. Topic Id.
// returns a *string when successful
func (m *ApplePushNotificationCertificate) GetTopicIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("topicIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ApplePushNotificationCertificate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appleIdentifier", m.GetAppleIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("certificate", m.GetCertificate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("certificateUploadFailureReason", m.GetCertificateUploadFailureReason())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("certificateUploadStatus", m.GetCertificateUploadStatus())
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
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("topicIdentifier", m.GetTopicIdentifier())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppleIdentifier sets the appleIdentifier property value. Apple Id of the account used to create the MDM push certificate.
func (m *ApplePushNotificationCertificate) SetAppleIdentifier(value *string)() {
    err := m.GetBackingStore().Set("appleIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificate sets the certificate property value. Not yet documented
func (m *ApplePushNotificationCertificate) SetCertificate(value *string)() {
    err := m.GetBackingStore().Set("certificate", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificateSerialNumber sets the certificateSerialNumber property value. Certificate serial number. This property is read-only.
func (m *ApplePushNotificationCertificate) SetCertificateSerialNumber(value *string)() {
    err := m.GetBackingStore().Set("certificateSerialNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificateUploadFailureReason sets the certificateUploadFailureReason property value. The reason the certificate upload failed.
func (m *ApplePushNotificationCertificate) SetCertificateUploadFailureReason(value *string)() {
    err := m.GetBackingStore().Set("certificateUploadFailureReason", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificateUploadStatus sets the certificateUploadStatus property value. The certificate upload status.
func (m *ApplePushNotificationCertificate) SetCertificateUploadStatus(value *string)() {
    err := m.GetBackingStore().Set("certificateUploadStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetExpirationDateTime sets the expirationDateTime property value. The expiration date and time for Apple push notification certificate.
func (m *ApplePushNotificationCertificate) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Last modified date and time for Apple push notification certificate.
func (m *ApplePushNotificationCertificate) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTopicIdentifier sets the topicIdentifier property value. Topic Id.
func (m *ApplePushNotificationCertificate) SetTopicIdentifier(value *string)() {
    err := m.GetBackingStore().Set("topicIdentifier", value)
    if err != nil {
        panic(err)
    }
}
type ApplePushNotificationCertificateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppleIdentifier()(*string)
    GetCertificate()(*string)
    GetCertificateSerialNumber()(*string)
    GetCertificateUploadFailureReason()(*string)
    GetCertificateUploadStatus()(*string)
    GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTopicIdentifier()(*string)
    SetAppleIdentifier(value *string)()
    SetCertificate(value *string)()
    SetCertificateSerialNumber(value *string)()
    SetCertificateUploadFailureReason(value *string)()
    SetCertificateUploadStatus(value *string)()
    SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTopicIdentifier(value *string)()
}
