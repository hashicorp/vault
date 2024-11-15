package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BrowserSiteList a singleton entity which is used to specify IE mode site list metadata
type BrowserSiteList struct {
    Entity
}
// NewBrowserSiteList instantiates a new BrowserSiteList and sets the default values.
func NewBrowserSiteList()(*BrowserSiteList) {
    m := &BrowserSiteList{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBrowserSiteListFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBrowserSiteListFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBrowserSiteList(), nil
}
// GetDescription gets the description property value. The description of the site list.
// returns a *string when successful
func (m *BrowserSiteList) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the site list.
// returns a *string when successful
func (m *BrowserSiteList) GetDisplayName()(*string) {
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
func (m *BrowserSiteList) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["publishedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublishedBy(val.(IdentitySetable))
        }
        return nil
    }
    res["publishedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublishedDateTime(val)
        }
        return nil
    }
    res["revision"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRevision(val)
        }
        return nil
    }
    res["sharedCookies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBrowserSharedCookieFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BrowserSharedCookieable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BrowserSharedCookieable)
                }
            }
            m.SetSharedCookies(res)
        }
        return nil
    }
    res["sites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBrowserSiteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BrowserSiteable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BrowserSiteable)
                }
            }
            m.SetSites(res)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSiteListStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*BrowserSiteListStatus))
        }
        return nil
    }
    return res
}
// GetLastModifiedBy gets the lastModifiedBy property value. The user who last modified the site list.
// returns a IdentitySetable when successful
func (m *BrowserSiteList) GetLastModifiedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the site list was last modified.
// returns a *Time when successful
func (m *BrowserSiteList) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPublishedBy gets the publishedBy property value. The user who published the site list.
// returns a IdentitySetable when successful
func (m *BrowserSiteList) GetPublishedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("publishedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetPublishedDateTime gets the publishedDateTime property value. The date and time when the site list was published.
// returns a *Time when successful
func (m *BrowserSiteList) GetPublishedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("publishedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRevision gets the revision property value. The current revision of the site list.
// returns a *string when successful
func (m *BrowserSiteList) GetRevision()(*string) {
    val, err := m.GetBackingStore().Get("revision")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSharedCookies gets the sharedCookies property value. A collection of shared cookies defined for the site list.
// returns a []BrowserSharedCookieable when successful
func (m *BrowserSiteList) GetSharedCookies()([]BrowserSharedCookieable) {
    val, err := m.GetBackingStore().Get("sharedCookies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BrowserSharedCookieable)
    }
    return nil
}
// GetSites gets the sites property value. A collection of sites defined for the site list.
// returns a []BrowserSiteable when successful
func (m *BrowserSiteList) GetSites()([]BrowserSiteable) {
    val, err := m.GetBackingStore().Get("sites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BrowserSiteable)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *BrowserSiteListStatus when successful
func (m *BrowserSiteList) GetStatus()(*BrowserSiteListStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSiteListStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BrowserSiteList) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteObjectValue("publishedBy", m.GetPublishedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("publishedDateTime", m.GetPublishedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("revision", m.GetRevision())
        if err != nil {
            return err
        }
    }
    if m.GetSharedCookies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharedCookies()))
        for i, v := range m.GetSharedCookies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sharedCookies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSites() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSites()))
        for i, v := range m.GetSites() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sites", cast)
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
// SetDescription sets the description property value. The description of the site list.
func (m *BrowserSiteList) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the site list.
func (m *BrowserSiteList) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The user who last modified the site list.
func (m *BrowserSiteList) SetLastModifiedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the site list was last modified.
func (m *BrowserSiteList) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPublishedBy sets the publishedBy property value. The user who published the site list.
func (m *BrowserSiteList) SetPublishedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("publishedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetPublishedDateTime sets the publishedDateTime property value. The date and time when the site list was published.
func (m *BrowserSiteList) SetPublishedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("publishedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRevision sets the revision property value. The current revision of the site list.
func (m *BrowserSiteList) SetRevision(value *string)() {
    err := m.GetBackingStore().Set("revision", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedCookies sets the sharedCookies property value. A collection of shared cookies defined for the site list.
func (m *BrowserSiteList) SetSharedCookies(value []BrowserSharedCookieable)() {
    err := m.GetBackingStore().Set("sharedCookies", value)
    if err != nil {
        panic(err)
    }
}
// SetSites sets the sites property value. A collection of sites defined for the site list.
func (m *BrowserSiteList) SetSites(value []BrowserSiteable)() {
    err := m.GetBackingStore().Set("sites", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *BrowserSiteList) SetStatus(value *BrowserSiteListStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type BrowserSiteListable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetLastModifiedBy()(IdentitySetable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPublishedBy()(IdentitySetable)
    GetPublishedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRevision()(*string)
    GetSharedCookies()([]BrowserSharedCookieable)
    GetSites()([]BrowserSiteable)
    GetStatus()(*BrowserSiteListStatus)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetLastModifiedBy(value IdentitySetable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPublishedBy(value IdentitySetable)()
    SetPublishedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRevision(value *string)()
    SetSharedCookies(value []BrowserSharedCookieable)()
    SetSites(value []BrowserSiteable)()
    SetStatus(value *BrowserSiteListStatus)()
}
