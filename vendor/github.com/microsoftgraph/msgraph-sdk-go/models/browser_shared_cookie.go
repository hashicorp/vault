package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type BrowserSharedCookie struct {
    Entity
}
// NewBrowserSharedCookie instantiates a new BrowserSharedCookie and sets the default values.
func NewBrowserSharedCookie()(*BrowserSharedCookie) {
    m := &BrowserSharedCookie{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBrowserSharedCookieFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBrowserSharedCookieFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBrowserSharedCookie(), nil
}
// GetComment gets the comment property value. The comment for the shared cookie.
// returns a *string when successful
func (m *BrowserSharedCookie) GetComment()(*string) {
    val, err := m.GetBackingStore().Get("comment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the shared cookie was created.
// returns a *Time when successful
func (m *BrowserSharedCookie) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDeletedDateTime gets the deletedDateTime property value. The date and time when the shared cookie was deleted.
// returns a *Time when successful
func (m *BrowserSharedCookie) GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("deletedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the cookie.
// returns a *string when successful
func (m *BrowserSharedCookie) GetDisplayName()(*string) {
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
func (m *BrowserSharedCookie) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["comment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComment(val)
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
    res["deletedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeletedDateTime(val)
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
    res["history"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBrowserSharedCookieHistoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BrowserSharedCookieHistoryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BrowserSharedCookieHistoryable)
                }
            }
            m.SetHistory(res)
        }
        return nil
    }
    res["hostOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHostOnly(val)
        }
        return nil
    }
    res["hostOrDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHostOrDomain(val)
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(IdentitySetable))
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
    res["path"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPath(val)
        }
        return nil
    }
    res["sourceEnvironment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSharedCookieSourceEnvironment)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceEnvironment(val.(*BrowserSharedCookieSourceEnvironment))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSharedCookieStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*BrowserSharedCookieStatus))
        }
        return nil
    }
    return res
}
// GetHistory gets the history property value. The history of modifications applied to the cookie.
// returns a []BrowserSharedCookieHistoryable when successful
func (m *BrowserSharedCookie) GetHistory()([]BrowserSharedCookieHistoryable) {
    val, err := m.GetBackingStore().Get("history")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BrowserSharedCookieHistoryable)
    }
    return nil
}
// GetHostOnly gets the hostOnly property value. Controls whether a cookie is a host-only or domain cookie.
// returns a *bool when successful
func (m *BrowserSharedCookie) GetHostOnly()(*bool) {
    val, err := m.GetBackingStore().Get("hostOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHostOrDomain gets the hostOrDomain property value. The URL of the cookie.
// returns a *string when successful
func (m *BrowserSharedCookie) GetHostOrDomain()(*string) {
    val, err := m.GetBackingStore().Get("hostOrDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. The user who last modified the cookie.
// returns a IdentitySetable when successful
func (m *BrowserSharedCookie) GetLastModifiedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the cookie was last modified.
// returns a *Time when successful
func (m *BrowserSharedCookie) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPath gets the path property value. The path of the cookie.
// returns a *string when successful
func (m *BrowserSharedCookie) GetPath()(*string) {
    val, err := m.GetBackingStore().Get("path")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSourceEnvironment gets the sourceEnvironment property value. The sourceEnvironment property
// returns a *BrowserSharedCookieSourceEnvironment when successful
func (m *BrowserSharedCookie) GetSourceEnvironment()(*BrowserSharedCookieSourceEnvironment) {
    val, err := m.GetBackingStore().Get("sourceEnvironment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSharedCookieSourceEnvironment)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *BrowserSharedCookieStatus when successful
func (m *BrowserSharedCookie) GetStatus()(*BrowserSharedCookieStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSharedCookieStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BrowserSharedCookie) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("comment", m.GetComment())
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
        err = writer.WriteTimeValue("deletedDateTime", m.GetDeletedDateTime())
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
    if m.GetHistory() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHistory()))
        for i, v := range m.GetHistory() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("history", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hostOnly", m.GetHostOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("hostOrDomain", m.GetHostOrDomain())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastModifiedBy", m.GetLastModifiedBy())
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
        err = writer.WriteStringValue("path", m.GetPath())
        if err != nil {
            return err
        }
    }
    if m.GetSourceEnvironment() != nil {
        cast := (*m.GetSourceEnvironment()).String()
        err = writer.WriteStringValue("sourceEnvironment", &cast)
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
    return nil
}
// SetComment sets the comment property value. The comment for the shared cookie.
func (m *BrowserSharedCookie) SetComment(value *string)() {
    err := m.GetBackingStore().Set("comment", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the shared cookie was created.
func (m *BrowserSharedCookie) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDeletedDateTime sets the deletedDateTime property value. The date and time when the shared cookie was deleted.
func (m *BrowserSharedCookie) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("deletedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the cookie.
func (m *BrowserSharedCookie) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetHistory sets the history property value. The history of modifications applied to the cookie.
func (m *BrowserSharedCookie) SetHistory(value []BrowserSharedCookieHistoryable)() {
    err := m.GetBackingStore().Set("history", value)
    if err != nil {
        panic(err)
    }
}
// SetHostOnly sets the hostOnly property value. Controls whether a cookie is a host-only or domain cookie.
func (m *BrowserSharedCookie) SetHostOnly(value *bool)() {
    err := m.GetBackingStore().Set("hostOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetHostOrDomain sets the hostOrDomain property value. The URL of the cookie.
func (m *BrowserSharedCookie) SetHostOrDomain(value *string)() {
    err := m.GetBackingStore().Set("hostOrDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The user who last modified the cookie.
func (m *BrowserSharedCookie) SetLastModifiedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the cookie was last modified.
func (m *BrowserSharedCookie) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPath sets the path property value. The path of the cookie.
func (m *BrowserSharedCookie) SetPath(value *string)() {
    err := m.GetBackingStore().Set("path", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceEnvironment sets the sourceEnvironment property value. The sourceEnvironment property
func (m *BrowserSharedCookie) SetSourceEnvironment(value *BrowserSharedCookieSourceEnvironment)() {
    err := m.GetBackingStore().Set("sourceEnvironment", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *BrowserSharedCookie) SetStatus(value *BrowserSharedCookieStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type BrowserSharedCookieable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetComment()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDisplayName()(*string)
    GetHistory()([]BrowserSharedCookieHistoryable)
    GetHostOnly()(*bool)
    GetHostOrDomain()(*string)
    GetLastModifiedBy()(IdentitySetable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPath()(*string)
    GetSourceEnvironment()(*BrowserSharedCookieSourceEnvironment)
    GetStatus()(*BrowserSharedCookieStatus)
    SetComment(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDisplayName(value *string)()
    SetHistory(value []BrowserSharedCookieHistoryable)()
    SetHostOnly(value *bool)()
    SetHostOrDomain(value *string)()
    SetLastModifiedBy(value IdentitySetable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPath(value *string)()
    SetSourceEnvironment(value *BrowserSharedCookieSourceEnvironment)()
    SetStatus(value *BrowserSharedCookieStatus)()
}
