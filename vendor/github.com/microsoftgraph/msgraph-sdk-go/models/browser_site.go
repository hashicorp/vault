package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BrowserSite singleton entity which is used to specify IE mode site metadata
type BrowserSite struct {
    Entity
}
// NewBrowserSite instantiates a new BrowserSite and sets the default values.
func NewBrowserSite()(*BrowserSite) {
    m := &BrowserSite{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBrowserSiteFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBrowserSiteFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBrowserSite(), nil
}
// GetAllowRedirect gets the allowRedirect property value. Controls the behavior of redirected sites. If true, indicates that the site will open in Internet Explorer 11 or Microsoft Edge even if the site is navigated to as part of a HTTP or meta refresh redirection chain.
// returns a *bool when successful
func (m *BrowserSite) GetAllowRedirect()(*bool) {
    val, err := m.GetBackingStore().Get("allowRedirect")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetComment gets the comment property value. The comment for the site.
// returns a *string when successful
func (m *BrowserSite) GetComment()(*string) {
    val, err := m.GetBackingStore().Get("comment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCompatibilityMode gets the compatibilityMode property value. The compatibilityMode property
// returns a *BrowserSiteCompatibilityMode when successful
func (m *BrowserSite) GetCompatibilityMode()(*BrowserSiteCompatibilityMode) {
    val, err := m.GetBackingStore().Get("compatibilityMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSiteCompatibilityMode)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the site was created.
// returns a *Time when successful
func (m *BrowserSite) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDeletedDateTime gets the deletedDateTime property value. The date and time when the site was deleted.
// returns a *Time when successful
func (m *BrowserSite) GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("deletedDateTime")
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
func (m *BrowserSite) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowRedirect"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowRedirect(val)
        }
        return nil
    }
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
    res["compatibilityMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSiteCompatibilityMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompatibilityMode(val.(*BrowserSiteCompatibilityMode))
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
    res["history"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBrowserSiteHistoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BrowserSiteHistoryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BrowserSiteHistoryable)
                }
            }
            m.SetHistory(res)
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
    res["mergeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSiteMergeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMergeType(val.(*BrowserSiteMergeType))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSiteStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*BrowserSiteStatus))
        }
        return nil
    }
    res["targetEnvironment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBrowserSiteTargetEnvironment)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetEnvironment(val.(*BrowserSiteTargetEnvironment))
        }
        return nil
    }
    res["webUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebUrl(val)
        }
        return nil
    }
    return res
}
// GetHistory gets the history property value. The history of modifications applied to the site.
// returns a []BrowserSiteHistoryable when successful
func (m *BrowserSite) GetHistory()([]BrowserSiteHistoryable) {
    val, err := m.GetBackingStore().Get("history")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BrowserSiteHistoryable)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. The user who last modified the site.
// returns a IdentitySetable when successful
func (m *BrowserSite) GetLastModifiedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the site was last modified.
// returns a *Time when successful
func (m *BrowserSite) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMergeType gets the mergeType property value. The mergeType property
// returns a *BrowserSiteMergeType when successful
func (m *BrowserSite) GetMergeType()(*BrowserSiteMergeType) {
    val, err := m.GetBackingStore().Get("mergeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSiteMergeType)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *BrowserSiteStatus when successful
func (m *BrowserSite) GetStatus()(*BrowserSiteStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSiteStatus)
    }
    return nil
}
// GetTargetEnvironment gets the targetEnvironment property value. The targetEnvironment property
// returns a *BrowserSiteTargetEnvironment when successful
func (m *BrowserSite) GetTargetEnvironment()(*BrowserSiteTargetEnvironment) {
    val, err := m.GetBackingStore().Get("targetEnvironment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BrowserSiteTargetEnvironment)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. The URL of the site.
// returns a *string when successful
func (m *BrowserSite) GetWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("webUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BrowserSite) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowRedirect", m.GetAllowRedirect())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("comment", m.GetComment())
        if err != nil {
            return err
        }
    }
    if m.GetCompatibilityMode() != nil {
        cast := (*m.GetCompatibilityMode()).String()
        err = writer.WriteStringValue("compatibilityMode", &cast)
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
    if m.GetMergeType() != nil {
        cast := (*m.GetMergeType()).String()
        err = writer.WriteStringValue("mergeType", &cast)
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
    if m.GetTargetEnvironment() != nil {
        cast := (*m.GetTargetEnvironment()).String()
        err = writer.WriteStringValue("targetEnvironment", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webUrl", m.GetWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowRedirect sets the allowRedirect property value. Controls the behavior of redirected sites. If true, indicates that the site will open in Internet Explorer 11 or Microsoft Edge even if the site is navigated to as part of a HTTP or meta refresh redirection chain.
func (m *BrowserSite) SetAllowRedirect(value *bool)() {
    err := m.GetBackingStore().Set("allowRedirect", value)
    if err != nil {
        panic(err)
    }
}
// SetComment sets the comment property value. The comment for the site.
func (m *BrowserSite) SetComment(value *string)() {
    err := m.GetBackingStore().Set("comment", value)
    if err != nil {
        panic(err)
    }
}
// SetCompatibilityMode sets the compatibilityMode property value. The compatibilityMode property
func (m *BrowserSite) SetCompatibilityMode(value *BrowserSiteCompatibilityMode)() {
    err := m.GetBackingStore().Set("compatibilityMode", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the site was created.
func (m *BrowserSite) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDeletedDateTime sets the deletedDateTime property value. The date and time when the site was deleted.
func (m *BrowserSite) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("deletedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetHistory sets the history property value. The history of modifications applied to the site.
func (m *BrowserSite) SetHistory(value []BrowserSiteHistoryable)() {
    err := m.GetBackingStore().Set("history", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The user who last modified the site.
func (m *BrowserSite) SetLastModifiedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the site was last modified.
func (m *BrowserSite) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMergeType sets the mergeType property value. The mergeType property
func (m *BrowserSite) SetMergeType(value *BrowserSiteMergeType)() {
    err := m.GetBackingStore().Set("mergeType", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *BrowserSite) SetStatus(value *BrowserSiteStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetEnvironment sets the targetEnvironment property value. The targetEnvironment property
func (m *BrowserSite) SetTargetEnvironment(value *BrowserSiteTargetEnvironment)() {
    err := m.GetBackingStore().Set("targetEnvironment", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. The URL of the site.
func (m *BrowserSite) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type BrowserSiteable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowRedirect()(*bool)
    GetComment()(*string)
    GetCompatibilityMode()(*BrowserSiteCompatibilityMode)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetHistory()([]BrowserSiteHistoryable)
    GetLastModifiedBy()(IdentitySetable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMergeType()(*BrowserSiteMergeType)
    GetStatus()(*BrowserSiteStatus)
    GetTargetEnvironment()(*BrowserSiteTargetEnvironment)
    GetWebUrl()(*string)
    SetAllowRedirect(value *bool)()
    SetComment(value *string)()
    SetCompatibilityMode(value *BrowserSiteCompatibilityMode)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetHistory(value []BrowserSiteHistoryable)()
    SetLastModifiedBy(value IdentitySetable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMergeType(value *BrowserSiteMergeType)()
    SetStatus(value *BrowserSiteStatus)()
    SetTargetEnvironment(value *BrowserSiteTargetEnvironment)()
    SetWebUrl(value *string)()
}
