package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SiteProtectionUnit struct {
    ProtectionUnitBase
}
// NewSiteProtectionUnit instantiates a new SiteProtectionUnit and sets the default values.
func NewSiteProtectionUnit()(*SiteProtectionUnit) {
    m := &SiteProtectionUnit{
        ProtectionUnitBase: *NewProtectionUnitBase(),
    }
    odataTypeValue := "#microsoft.graph.siteProtectionUnit"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSiteProtectionUnitFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSiteProtectionUnitFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSiteProtectionUnit(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SiteProtectionUnit) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ProtectionUnitBase.GetFieldDeserializers()
    res["siteId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteId(val)
        }
        return nil
    }
    res["siteName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteName(val)
        }
        return nil
    }
    res["siteWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteWebUrl(val)
        }
        return nil
    }
    return res
}
// GetSiteId gets the siteId property value. Unique identifier of the SharePoint site.
// returns a *string when successful
func (m *SiteProtectionUnit) GetSiteId()(*string) {
    val, err := m.GetBackingStore().Get("siteId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSiteName gets the siteName property value. Name of the SharePoint site.
// returns a *string when successful
func (m *SiteProtectionUnit) GetSiteName()(*string) {
    val, err := m.GetBackingStore().Get("siteName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSiteWebUrl gets the siteWebUrl property value. The web URL of the SharePoint site.
// returns a *string when successful
func (m *SiteProtectionUnit) GetSiteWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("siteWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SiteProtectionUnit) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ProtectionUnitBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("siteId", m.GetSiteId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSiteId sets the siteId property value. Unique identifier of the SharePoint site.
func (m *SiteProtectionUnit) SetSiteId(value *string)() {
    err := m.GetBackingStore().Set("siteId", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteName sets the siteName property value. Name of the SharePoint site.
func (m *SiteProtectionUnit) SetSiteName(value *string)() {
    err := m.GetBackingStore().Set("siteName", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteWebUrl sets the siteWebUrl property value. The web URL of the SharePoint site.
func (m *SiteProtectionUnit) SetSiteWebUrl(value *string)() {
    err := m.GetBackingStore().Set("siteWebUrl", value)
    if err != nil {
        panic(err)
    }
}
type SiteProtectionUnitable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ProtectionUnitBaseable
    GetSiteId()(*string)
    GetSiteName()(*string)
    GetSiteWebUrl()(*string)
    SetSiteId(value *string)()
    SetSiteName(value *string)()
    SetSiteWebUrl(value *string)()
}
