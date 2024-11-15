package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Subscription struct {
    Entity
}
// NewSubscription instantiates a new Subscription and sets the default values.
func NewSubscription()(*Subscription) {
    m := &Subscription{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSubscriptionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSubscriptionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSubscription(), nil
}
// GetApplicationId gets the applicationId property value. Optional. Identifier of the application used to create the subscription. Read-only.
// returns a *string when successful
func (m *Subscription) GetApplicationId()(*string) {
    val, err := m.GetBackingStore().Get("applicationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetChangeType gets the changeType property value. Required. Indicates the type of change in the subscribed resource that raises a change notification. The supported values are: created, updated, deleted. Multiple values can be combined using a comma-separated list. Note:  Drive root item and list change notifications support only the updated changeType. User and group change notifications support updated and deleted changeType. Use updated to receive notifications when user or group is created, updated, or soft deleted. Use deleted to receive notifications when user or group is permanently deleted.
// returns a *string when successful
func (m *Subscription) GetChangeType()(*string) {
    val, err := m.GetBackingStore().Get("changeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetClientState gets the clientState property value. Optional. Specifies the value of the clientState property sent by the service in each change notification. The maximum length is 128 characters. The client can check that the change notification came from the service by comparing the value of the clientState property sent with the subscription with the value of the clientState property received with each change notification.
// returns a *string when successful
func (m *Subscription) GetClientState()(*string) {
    val, err := m.GetBackingStore().Get("clientState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatorId gets the creatorId property value. Optional. Identifier of the user or service principal that created the subscription. If the app used delegated permissions to create the subscription, this field contains the ID of the signed-in user the app called on behalf of. If the app used application permissions, this field contains the ID of the service principal corresponding to the app. Read-only.
// returns a *string when successful
func (m *Subscription) GetCreatorId()(*string) {
    val, err := m.GetBackingStore().Get("creatorId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEncryptionCertificate gets the encryptionCertificate property value. Optional. A base64-encoded representation of a certificate with a public key used to encrypt resource data in change notifications. Optional but required when includeResourceData is true.
// returns a *string when successful
func (m *Subscription) GetEncryptionCertificate()(*string) {
    val, err := m.GetBackingStore().Get("encryptionCertificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEncryptionCertificateId gets the encryptionCertificateId property value. Optional. A custom app-provided identifier to help identify the certificate needed to decrypt resource data.
// returns a *string when successful
func (m *Subscription) GetEncryptionCertificateId()(*string) {
    val, err := m.GetBackingStore().Get("encryptionCertificateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExpirationDateTime gets the expirationDateTime property value. Required. Specifies the date and time when the webhook subscription expires. The time is in UTC, and can be an amount of time from subscription creation that varies for the resource subscribed to. For the maximum supported subscription length of time, see Subscription lifetime.
// returns a *Time when successful
func (m *Subscription) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
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
func (m *Subscription) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["applicationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationId(val)
        }
        return nil
    }
    res["changeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChangeType(val)
        }
        return nil
    }
    res["clientState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientState(val)
        }
        return nil
    }
    res["creatorId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatorId(val)
        }
        return nil
    }
    res["encryptionCertificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEncryptionCertificate(val)
        }
        return nil
    }
    res["encryptionCertificateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEncryptionCertificateId(val)
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
    res["includeResourceData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncludeResourceData(val)
        }
        return nil
    }
    res["latestSupportedTlsVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLatestSupportedTlsVersion(val)
        }
        return nil
    }
    res["lifecycleNotificationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLifecycleNotificationUrl(val)
        }
        return nil
    }
    res["notificationQueryOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationQueryOptions(val)
        }
        return nil
    }
    res["notificationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationUrl(val)
        }
        return nil
    }
    res["notificationUrlAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationUrlAppId(val)
        }
        return nil
    }
    res["resource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResource(val)
        }
        return nil
    }
    return res
}
// GetIncludeResourceData gets the includeResourceData property value. Optional. When set to true, change notifications include resource data (such as content of a chat message).
// returns a *bool when successful
func (m *Subscription) GetIncludeResourceData()(*bool) {
    val, err := m.GetBackingStore().Get("includeResourceData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLatestSupportedTlsVersion gets the latestSupportedTlsVersion property value. Optional. Specifies the latest version of Transport Layer Security (TLS) that the notification endpoint, specified by notificationUrl, supports. The possible values are: v10, v11, v12, v13. For subscribers whose notification endpoint supports a version lower than the currently recommended version (TLS 1.2), specifying this property by a set timeline allows them to temporarily use their deprecated version of TLS before completing their upgrade to TLS 1.2. For these subscribers, not setting this property per the timeline would result in subscription operations failing. For subscribers whose notification endpoint already supports TLS 1.2, setting this property is optional. In such cases, Microsoft Graph defaults the property to v1_2.
// returns a *string when successful
func (m *Subscription) GetLatestSupportedTlsVersion()(*string) {
    val, err := m.GetBackingStore().Get("latestSupportedTlsVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLifecycleNotificationUrl gets the lifecycleNotificationUrl property value. Required for Teams resources if  the expirationDateTime value is more than 1 hour from now; optional otherwise. The URL of the endpoint that receives lifecycle notifications, including subscriptionRemoved, reauthorizationRequired, and missed notifications. This URL must make use of the HTTPS protocol. For more information, see Reduce missing subscriptions and change notifications.
// returns a *string when successful
func (m *Subscription) GetLifecycleNotificationUrl()(*string) {
    val, err := m.GetBackingStore().Get("lifecycleNotificationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotificationQueryOptions gets the notificationQueryOptions property value. Optional. OData query options for specifying value for the targeting resource. Clients receive notifications when resource reaches the state matching the query options provided here. With this new property in the subscription creation payload along with all existing properties, Webhooks deliver notifications whenever a resource reaches the desired state mentioned in the notificationQueryOptions property. For example, when the print job is completed or when a print job resource isFetchable property value becomes true etc.  Supported only for Universal Print Service. For more information, see Subscribe to change notifications from cloud printing APIs using Microsoft Graph.
// returns a *string when successful
func (m *Subscription) GetNotificationQueryOptions()(*string) {
    val, err := m.GetBackingStore().Get("notificationQueryOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotificationUrl gets the notificationUrl property value. Required. The URL of the endpoint that receives the change notifications. This URL must make use of the HTTPS protocol. Any query string parameter included in the notificationUrl property is included in the HTTP POST request when Microsoft Graph sends the change notifications.
// returns a *string when successful
func (m *Subscription) GetNotificationUrl()(*string) {
    val, err := m.GetBackingStore().Get("notificationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotificationUrlAppId gets the notificationUrlAppId property value. Optional. The app ID that the subscription service can use to generate the validation token. The value allows the client to validate the authenticity of the notification received.
// returns a *string when successful
func (m *Subscription) GetNotificationUrlAppId()(*string) {
    val, err := m.GetBackingStore().Get("notificationUrlAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResource gets the resource property value. Required. Specifies the resource that is monitored for changes. Don't include the base URL (https://graph.microsoft.com/v1.0/). See the possible resource path values for each supported resource.
// returns a *string when successful
func (m *Subscription) GetResource()(*string) {
    val, err := m.GetBackingStore().Get("resource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Subscription) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("applicationId", m.GetApplicationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("changeType", m.GetChangeType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("clientState", m.GetClientState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("creatorId", m.GetCreatorId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("encryptionCertificate", m.GetEncryptionCertificate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("encryptionCertificateId", m.GetEncryptionCertificateId())
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
        err = writer.WriteBoolValue("includeResourceData", m.GetIncludeResourceData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("latestSupportedTlsVersion", m.GetLatestSupportedTlsVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lifecycleNotificationUrl", m.GetLifecycleNotificationUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationQueryOptions", m.GetNotificationQueryOptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationUrl", m.GetNotificationUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationUrlAppId", m.GetNotificationUrlAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resource", m.GetResource())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationId sets the applicationId property value. Optional. Identifier of the application used to create the subscription. Read-only.
func (m *Subscription) SetApplicationId(value *string)() {
    err := m.GetBackingStore().Set("applicationId", value)
    if err != nil {
        panic(err)
    }
}
// SetChangeType sets the changeType property value. Required. Indicates the type of change in the subscribed resource that raises a change notification. The supported values are: created, updated, deleted. Multiple values can be combined using a comma-separated list. Note:  Drive root item and list change notifications support only the updated changeType. User and group change notifications support updated and deleted changeType. Use updated to receive notifications when user or group is created, updated, or soft deleted. Use deleted to receive notifications when user or group is permanently deleted.
func (m *Subscription) SetChangeType(value *string)() {
    err := m.GetBackingStore().Set("changeType", value)
    if err != nil {
        panic(err)
    }
}
// SetClientState sets the clientState property value. Optional. Specifies the value of the clientState property sent by the service in each change notification. The maximum length is 128 characters. The client can check that the change notification came from the service by comparing the value of the clientState property sent with the subscription with the value of the clientState property received with each change notification.
func (m *Subscription) SetClientState(value *string)() {
    err := m.GetBackingStore().Set("clientState", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatorId sets the creatorId property value. Optional. Identifier of the user or service principal that created the subscription. If the app used delegated permissions to create the subscription, this field contains the ID of the signed-in user the app called on behalf of. If the app used application permissions, this field contains the ID of the service principal corresponding to the app. Read-only.
func (m *Subscription) SetCreatorId(value *string)() {
    err := m.GetBackingStore().Set("creatorId", value)
    if err != nil {
        panic(err)
    }
}
// SetEncryptionCertificate sets the encryptionCertificate property value. Optional. A base64-encoded representation of a certificate with a public key used to encrypt resource data in change notifications. Optional but required when includeResourceData is true.
func (m *Subscription) SetEncryptionCertificate(value *string)() {
    err := m.GetBackingStore().Set("encryptionCertificate", value)
    if err != nil {
        panic(err)
    }
}
// SetEncryptionCertificateId sets the encryptionCertificateId property value. Optional. A custom app-provided identifier to help identify the certificate needed to decrypt resource data.
func (m *Subscription) SetEncryptionCertificateId(value *string)() {
    err := m.GetBackingStore().Set("encryptionCertificateId", value)
    if err != nil {
        panic(err)
    }
}
// SetExpirationDateTime sets the expirationDateTime property value. Required. Specifies the date and time when the webhook subscription expires. The time is in UTC, and can be an amount of time from subscription creation that varies for the resource subscribed to. For the maximum supported subscription length of time, see Subscription lifetime.
func (m *Subscription) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeResourceData sets the includeResourceData property value. Optional. When set to true, change notifications include resource data (such as content of a chat message).
func (m *Subscription) SetIncludeResourceData(value *bool)() {
    err := m.GetBackingStore().Set("includeResourceData", value)
    if err != nil {
        panic(err)
    }
}
// SetLatestSupportedTlsVersion sets the latestSupportedTlsVersion property value. Optional. Specifies the latest version of Transport Layer Security (TLS) that the notification endpoint, specified by notificationUrl, supports. The possible values are: v10, v11, v12, v13. For subscribers whose notification endpoint supports a version lower than the currently recommended version (TLS 1.2), specifying this property by a set timeline allows them to temporarily use their deprecated version of TLS before completing their upgrade to TLS 1.2. For these subscribers, not setting this property per the timeline would result in subscription operations failing. For subscribers whose notification endpoint already supports TLS 1.2, setting this property is optional. In such cases, Microsoft Graph defaults the property to v1_2.
func (m *Subscription) SetLatestSupportedTlsVersion(value *string)() {
    err := m.GetBackingStore().Set("latestSupportedTlsVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetLifecycleNotificationUrl sets the lifecycleNotificationUrl property value. Required for Teams resources if  the expirationDateTime value is more than 1 hour from now; optional otherwise. The URL of the endpoint that receives lifecycle notifications, including subscriptionRemoved, reauthorizationRequired, and missed notifications. This URL must make use of the HTTPS protocol. For more information, see Reduce missing subscriptions and change notifications.
func (m *Subscription) SetLifecycleNotificationUrl(value *string)() {
    err := m.GetBackingStore().Set("lifecycleNotificationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationQueryOptions sets the notificationQueryOptions property value. Optional. OData query options for specifying value for the targeting resource. Clients receive notifications when resource reaches the state matching the query options provided here. With this new property in the subscription creation payload along with all existing properties, Webhooks deliver notifications whenever a resource reaches the desired state mentioned in the notificationQueryOptions property. For example, when the print job is completed or when a print job resource isFetchable property value becomes true etc.  Supported only for Universal Print Service. For more information, see Subscribe to change notifications from cloud printing APIs using Microsoft Graph.
func (m *Subscription) SetNotificationQueryOptions(value *string)() {
    err := m.GetBackingStore().Set("notificationQueryOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationUrl sets the notificationUrl property value. Required. The URL of the endpoint that receives the change notifications. This URL must make use of the HTTPS protocol. Any query string parameter included in the notificationUrl property is included in the HTTP POST request when Microsoft Graph sends the change notifications.
func (m *Subscription) SetNotificationUrl(value *string)() {
    err := m.GetBackingStore().Set("notificationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationUrlAppId sets the notificationUrlAppId property value. Optional. The app ID that the subscription service can use to generate the validation token. The value allows the client to validate the authenticity of the notification received.
func (m *Subscription) SetNotificationUrlAppId(value *string)() {
    err := m.GetBackingStore().Set("notificationUrlAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetResource sets the resource property value. Required. Specifies the resource that is monitored for changes. Don't include the base URL (https://graph.microsoft.com/v1.0/). See the possible resource path values for each supported resource.
func (m *Subscription) SetResource(value *string)() {
    err := m.GetBackingStore().Set("resource", value)
    if err != nil {
        panic(err)
    }
}
type Subscriptionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationId()(*string)
    GetChangeType()(*string)
    GetClientState()(*string)
    GetCreatorId()(*string)
    GetEncryptionCertificate()(*string)
    GetEncryptionCertificateId()(*string)
    GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIncludeResourceData()(*bool)
    GetLatestSupportedTlsVersion()(*string)
    GetLifecycleNotificationUrl()(*string)
    GetNotificationQueryOptions()(*string)
    GetNotificationUrl()(*string)
    GetNotificationUrlAppId()(*string)
    GetResource()(*string)
    SetApplicationId(value *string)()
    SetChangeType(value *string)()
    SetClientState(value *string)()
    SetCreatorId(value *string)()
    SetEncryptionCertificate(value *string)()
    SetEncryptionCertificateId(value *string)()
    SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIncludeResourceData(value *bool)()
    SetLatestSupportedTlsVersion(value *string)()
    SetLifecycleNotificationUrl(value *string)()
    SetNotificationQueryOptions(value *string)()
    SetNotificationUrl(value *string)()
    SetNotificationUrlAppId(value *string)()
    SetResource(value *string)()
}
