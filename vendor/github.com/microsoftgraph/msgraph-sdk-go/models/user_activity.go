package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserActivity struct {
    Entity
}
// NewUserActivity instantiates a new UserActivity and sets the default values.
func NewUserActivity()(*UserActivity) {
    m := &UserActivity{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserActivityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserActivityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserActivity(), nil
}
// GetActivationUrl gets the activationUrl property value. Required. URL used to launch the activity in the best native experience represented by the appId. Might launch a web-based app if no native app exists.
// returns a *string when successful
func (m *UserActivity) GetActivationUrl()(*string) {
    val, err := m.GetBackingStore().Get("activationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetActivitySourceHost gets the activitySourceHost property value. Required. URL for the domain representing the cross-platform identity mapping for the app. Mapping is stored either as a JSON file hosted on the domain or configurable via Windows Dev Center. The JSON file is named cross-platform-app-identifiers and is hosted at root of your HTTPS domain, either at the top level domain or include a sub domain. For example: https://contoso.com or https://myapp.contoso.com but NOT https://myapp.contoso.com/somepath. You must have a unique file and domain (or sub domain) per cross-platform app identity. For example, a separate file and domain is needed for Word vs. PowerPoint.
// returns a *string when successful
func (m *UserActivity) GetActivitySourceHost()(*string) {
    val, err := m.GetBackingStore().Get("activitySourceHost")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppActivityId gets the appActivityId property value. Required. The unique activity ID in the context of the app - supplied by caller and immutable thereafter.
// returns a *string when successful
func (m *UserActivity) GetAppActivityId()(*string) {
    val, err := m.GetBackingStore().Get("appActivityId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppDisplayName gets the appDisplayName property value. Optional. Short text description of the app used to generate the activity for use in cases when the app is not installed on the user’s local device.
// returns a *string when successful
func (m *UserActivity) GetAppDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("appDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentInfo gets the contentInfo property value. Optional. A custom piece of data - JSON-LD extensible description of content according to schema.org syntax.
// returns a UntypedNodeable when successful
func (m *UserActivity) GetContentInfo()(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.UntypedNodeable) {
    val, err := m.GetBackingStore().Get("contentInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.UntypedNodeable)
    }
    return nil
}
// GetContentUrl gets the contentUrl property value. Optional. Used in the event the content can be rendered outside of a native or web-based app experience (for example, a pointer to an item in an RSS feed).
// returns a *string when successful
func (m *UserActivity) GetContentUrl()(*string) {
    val, err := m.GetBackingStore().Get("contentUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Set by the server. DateTime in UTC when the object was created on the server.
// returns a *Time when successful
func (m *UserActivity) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetExpirationDateTime gets the expirationDateTime property value. Set by the server. DateTime in UTC when the object expired on the server.
// returns a *Time when successful
func (m *UserActivity) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("expirationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFallbackUrl gets the fallbackUrl property value. Optional. URL used to launch the activity in a web-based app, if available.
// returns a *string when successful
func (m *UserActivity) GetFallbackUrl()(*string) {
    val, err := m.GetBackingStore().Get("fallbackUrl")
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
func (m *UserActivity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivationUrl(val)
        }
        return nil
    }
    res["activitySourceHost"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivitySourceHost(val)
        }
        return nil
    }
    res["appActivityId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppActivityId(val)
        }
        return nil
    }
    res["appDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppDisplayName(val)
        }
        return nil
    }
    res["contentInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.CreateUntypedNodeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentInfo(val.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.UntypedNodeable))
        }
        return nil
    }
    res["contentUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentUrl(val)
        }
        return nil
    }
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
    res["fallbackUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFallbackUrl(val)
        }
        return nil
    }
    res["historyItems"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateActivityHistoryItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ActivityHistoryItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ActivityHistoryItemable)
                }
            }
            m.SetHistoryItems(res)
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
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*Status))
        }
        return nil
    }
    res["userTimezone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserTimezone(val)
        }
        return nil
    }
    res["visualElements"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVisualInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVisualElements(val.(VisualInfoable))
        }
        return nil
    }
    return res
}
// GetHistoryItems gets the historyItems property value. Optional. NavigationProperty/Containment; navigation property to the activity's historyItems.
// returns a []ActivityHistoryItemable when successful
func (m *UserActivity) GetHistoryItems()([]ActivityHistoryItemable) {
    val, err := m.GetBackingStore().Get("historyItems")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ActivityHistoryItemable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Set by the server. DateTime in UTC when the object was modified on the server.
// returns a *Time when successful
func (m *UserActivity) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetStatus gets the status property value. Set by the server. A status code used to identify valid objects. Values: active, updated, deleted, ignored.
// returns a *Status when successful
func (m *UserActivity) GetStatus()(*Status) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Status)
    }
    return nil
}
// GetUserTimezone gets the userTimezone property value. Optional. The timezone in which the user's device used to generate the activity was located at activity creation time; values supplied as Olson IDs in order to support cross-platform representation.
// returns a *string when successful
func (m *UserActivity) GetUserTimezone()(*string) {
    val, err := m.GetBackingStore().Get("userTimezone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVisualElements gets the visualElements property value. The visualElements property
// returns a VisualInfoable when successful
func (m *UserActivity) GetVisualElements()(VisualInfoable) {
    val, err := m.GetBackingStore().Get("visualElements")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VisualInfoable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserActivity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("activationUrl", m.GetActivationUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activitySourceHost", m.GetActivitySourceHost())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appActivityId", m.GetAppActivityId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appDisplayName", m.GetAppDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("contentInfo", m.GetContentInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentUrl", m.GetContentUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
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
        err = writer.WriteStringValue("fallbackUrl", m.GetFallbackUrl())
        if err != nil {
            return err
        }
    }
    if m.GetHistoryItems() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHistoryItems()))
        for i, v := range m.GetHistoryItems() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("historyItems", cast)
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
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userTimezone", m.GetUserTimezone())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("visualElements", m.GetVisualElements())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivationUrl sets the activationUrl property value. Required. URL used to launch the activity in the best native experience represented by the appId. Might launch a web-based app if no native app exists.
func (m *UserActivity) SetActivationUrl(value *string)() {
    err := m.GetBackingStore().Set("activationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetActivitySourceHost sets the activitySourceHost property value. Required. URL for the domain representing the cross-platform identity mapping for the app. Mapping is stored either as a JSON file hosted on the domain or configurable via Windows Dev Center. The JSON file is named cross-platform-app-identifiers and is hosted at root of your HTTPS domain, either at the top level domain or include a sub domain. For example: https://contoso.com or https://myapp.contoso.com but NOT https://myapp.contoso.com/somepath. You must have a unique file and domain (or sub domain) per cross-platform app identity. For example, a separate file and domain is needed for Word vs. PowerPoint.
func (m *UserActivity) SetActivitySourceHost(value *string)() {
    err := m.GetBackingStore().Set("activitySourceHost", value)
    if err != nil {
        panic(err)
    }
}
// SetAppActivityId sets the appActivityId property value. Required. The unique activity ID in the context of the app - supplied by caller and immutable thereafter.
func (m *UserActivity) SetAppActivityId(value *string)() {
    err := m.GetBackingStore().Set("appActivityId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppDisplayName sets the appDisplayName property value. Optional. Short text description of the app used to generate the activity for use in cases when the app is not installed on the user’s local device.
func (m *UserActivity) SetAppDisplayName(value *string)() {
    err := m.GetBackingStore().Set("appDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetContentInfo sets the contentInfo property value. Optional. A custom piece of data - JSON-LD extensible description of content according to schema.org syntax.
func (m *UserActivity) SetContentInfo(value i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.UntypedNodeable)() {
    err := m.GetBackingStore().Set("contentInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetContentUrl sets the contentUrl property value. Optional. Used in the event the content can be rendered outside of a native or web-based app experience (for example, a pointer to an item in an RSS feed).
func (m *UserActivity) SetContentUrl(value *string)() {
    err := m.GetBackingStore().Set("contentUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Set by the server. DateTime in UTC when the object was created on the server.
func (m *UserActivity) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetExpirationDateTime sets the expirationDateTime property value. Set by the server. DateTime in UTC when the object expired on the server.
func (m *UserActivity) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFallbackUrl sets the fallbackUrl property value. Optional. URL used to launch the activity in a web-based app, if available.
func (m *UserActivity) SetFallbackUrl(value *string)() {
    err := m.GetBackingStore().Set("fallbackUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetHistoryItems sets the historyItems property value. Optional. NavigationProperty/Containment; navigation property to the activity's historyItems.
func (m *UserActivity) SetHistoryItems(value []ActivityHistoryItemable)() {
    err := m.GetBackingStore().Set("historyItems", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Set by the server. DateTime in UTC when the object was modified on the server.
func (m *UserActivity) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Set by the server. A status code used to identify valid objects. Values: active, updated, deleted, ignored.
func (m *UserActivity) SetStatus(value *Status)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetUserTimezone sets the userTimezone property value. Optional. The timezone in which the user's device used to generate the activity was located at activity creation time; values supplied as Olson IDs in order to support cross-platform representation.
func (m *UserActivity) SetUserTimezone(value *string)() {
    err := m.GetBackingStore().Set("userTimezone", value)
    if err != nil {
        panic(err)
    }
}
// SetVisualElements sets the visualElements property value. The visualElements property
func (m *UserActivity) SetVisualElements(value VisualInfoable)() {
    err := m.GetBackingStore().Set("visualElements", value)
    if err != nil {
        panic(err)
    }
}
type UserActivityable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivationUrl()(*string)
    GetActivitySourceHost()(*string)
    GetAppActivityId()(*string)
    GetAppDisplayName()(*string)
    GetContentInfo()(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.UntypedNodeable)
    GetContentUrl()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFallbackUrl()(*string)
    GetHistoryItems()([]ActivityHistoryItemable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetStatus()(*Status)
    GetUserTimezone()(*string)
    GetVisualElements()(VisualInfoable)
    SetActivationUrl(value *string)()
    SetActivitySourceHost(value *string)()
    SetAppActivityId(value *string)()
    SetAppDisplayName(value *string)()
    SetContentInfo(value i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.UntypedNodeable)()
    SetContentUrl(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFallbackUrl(value *string)()
    SetHistoryItems(value []ActivityHistoryItemable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetStatus(value *Status)()
    SetUserTimezone(value *string)()
    SetVisualElements(value VisualInfoable)()
}
