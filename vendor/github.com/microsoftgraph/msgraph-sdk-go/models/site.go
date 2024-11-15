package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Site struct {
    BaseItem
}
// NewSite instantiates a new Site and sets the default values.
func NewSite()(*Site) {
    m := &Site{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.site"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSiteFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSiteFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSite(), nil
}
// GetAnalytics gets the analytics property value. Analytics about the view activities that took place on this site.
// returns a ItemAnalyticsable when successful
func (m *Site) GetAnalytics()(ItemAnalyticsable) {
    val, err := m.GetBackingStore().Get("analytics")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemAnalyticsable)
    }
    return nil
}
// GetColumns gets the columns property value. The collection of column definitions reusable across lists under this site.
// returns a []ColumnDefinitionable when successful
func (m *Site) GetColumns()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("columns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetContentTypes gets the contentTypes property value. The collection of content types defined for this site.
// returns a []ContentTypeable when successful
func (m *Site) GetContentTypes()([]ContentTypeable) {
    val, err := m.GetBackingStore().Get("contentTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContentTypeable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The full title for the site. Read-only.
// returns a *string when successful
func (m *Site) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDrive gets the drive property value. The default drive (document library) for this site.
// returns a Driveable when successful
func (m *Site) GetDrive()(Driveable) {
    val, err := m.GetBackingStore().Get("drive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Driveable)
    }
    return nil
}
// GetDrives gets the drives property value. The collection of drives (document libraries) under this site.
// returns a []Driveable when successful
func (m *Site) GetDrives()([]Driveable) {
    val, err := m.GetBackingStore().Get("drives")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Driveable)
    }
    return nil
}
// GetError gets the error property value. The error property
// returns a PublicErrorable when successful
func (m *Site) GetError()(PublicErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PublicErrorable)
    }
    return nil
}
// GetExternalColumns gets the externalColumns property value. The externalColumns property
// returns a []ColumnDefinitionable when successful
func (m *Site) GetExternalColumns()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("externalColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Site) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["analytics"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemAnalyticsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnalytics(val.(ItemAnalyticsable))
        }
        return nil
    }
    res["columns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateColumnDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ColumnDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ColumnDefinitionable)
                }
            }
            m.SetColumns(res)
        }
        return nil
    }
    res["contentTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContentTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContentTypeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContentTypeable)
                }
            }
            m.SetContentTypes(res)
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
    res["drive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDrive(val.(Driveable))
        }
        return nil
    }
    res["drives"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Driveable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Driveable)
                }
            }
            m.SetDrives(res)
        }
        return nil
    }
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePublicErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(PublicErrorable))
        }
        return nil
    }
    res["externalColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateColumnDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ColumnDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ColumnDefinitionable)
                }
            }
            m.SetExternalColumns(res)
        }
        return nil
    }
    res["isPersonalSite"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsPersonalSite(val)
        }
        return nil
    }
    res["items"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBaseItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BaseItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BaseItemable)
                }
            }
            m.SetItems(res)
        }
        return nil
    }
    res["lists"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateListFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Listable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Listable)
                }
            }
            m.SetLists(res)
        }
        return nil
    }
    res["onenote"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnenoteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnenote(val.(Onenoteable))
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRichLongRunningOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RichLongRunningOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RichLongRunningOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["pages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBaseSitePageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BaseSitePageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BaseSitePageable)
                }
            }
            m.SetPages(res)
        }
        return nil
    }
    res["permissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePermissionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Permissionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Permissionable)
                }
            }
            m.SetPermissions(res)
        }
        return nil
    }
    res["root"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoot(val.(Rootable))
        }
        return nil
    }
    res["sharepointIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharepointIdsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharepointIds(val.(SharepointIdsable))
        }
        return nil
    }
    res["siteCollection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSiteCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteCollection(val.(SiteCollectionable))
        }
        return nil
    }
    res["sites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSiteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Siteable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Siteable)
                }
            }
            m.SetSites(res)
        }
        return nil
    }
    return res
}
// GetIsPersonalSite gets the isPersonalSite property value. Identifies whether the site is personal or not. Read-only.
// returns a *bool when successful
func (m *Site) GetIsPersonalSite()(*bool) {
    val, err := m.GetBackingStore().Get("isPersonalSite")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetItems gets the items property value. Used to address any item contained in this site. This collection can't be enumerated.
// returns a []BaseItemable when successful
func (m *Site) GetItems()([]BaseItemable) {
    val, err := m.GetBackingStore().Get("items")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BaseItemable)
    }
    return nil
}
// GetLists gets the lists property value. The collection of lists under this site.
// returns a []Listable when successful
func (m *Site) GetLists()([]Listable) {
    val, err := m.GetBackingStore().Get("lists")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Listable)
    }
    return nil
}
// GetOnenote gets the onenote property value. Calls the OneNote service for notebook related operations.
// returns a Onenoteable when successful
func (m *Site) GetOnenote()(Onenoteable) {
    val, err := m.GetBackingStore().Get("onenote")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Onenoteable)
    }
    return nil
}
// GetOperations gets the operations property value. The collection of long-running operations on the site.
// returns a []RichLongRunningOperationable when successful
func (m *Site) GetOperations()([]RichLongRunningOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RichLongRunningOperationable)
    }
    return nil
}
// GetPages gets the pages property value. The collection of pages in the baseSitePages list in this site.
// returns a []BaseSitePageable when successful
func (m *Site) GetPages()([]BaseSitePageable) {
    val, err := m.GetBackingStore().Get("pages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BaseSitePageable)
    }
    return nil
}
// GetPermissions gets the permissions property value. The permissions associated with the site. Nullable.
// returns a []Permissionable when successful
func (m *Site) GetPermissions()([]Permissionable) {
    val, err := m.GetBackingStore().Get("permissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Permissionable)
    }
    return nil
}
// GetRoot gets the root property value. If present, provides the root site in the site collection. Read-only.
// returns a Rootable when successful
func (m *Site) GetRoot()(Rootable) {
    val, err := m.GetBackingStore().Get("root")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Rootable)
    }
    return nil
}
// GetSharepointIds gets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
// returns a SharepointIdsable when successful
func (m *Site) GetSharepointIds()(SharepointIdsable) {
    val, err := m.GetBackingStore().Get("sharepointIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharepointIdsable)
    }
    return nil
}
// GetSiteCollection gets the siteCollection property value. Provides details about the site's site collection. Available only on the root site. Read-only.
// returns a SiteCollectionable when successful
func (m *Site) GetSiteCollection()(SiteCollectionable) {
    val, err := m.GetBackingStore().Get("siteCollection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SiteCollectionable)
    }
    return nil
}
// GetSites gets the sites property value. The collection of the sub-sites under this site.
// returns a []Siteable when successful
func (m *Site) GetSites()([]Siteable) {
    val, err := m.GetBackingStore().Get("sites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Siteable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Site) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("analytics", m.GetAnalytics())
        if err != nil {
            return err
        }
    }
    if m.GetColumns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetColumns()))
        for i, v := range m.GetColumns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("columns", cast)
        if err != nil {
            return err
        }
    }
    if m.GetContentTypes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetContentTypes()))
        for i, v := range m.GetContentTypes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("contentTypes", cast)
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
        err = writer.WriteObjectValue("drive", m.GetDrive())
        if err != nil {
            return err
        }
    }
    if m.GetDrives() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDrives()))
        for i, v := range m.GetDrives() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("drives", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("error", m.GetError())
        if err != nil {
            return err
        }
    }
    if m.GetExternalColumns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExternalColumns()))
        for i, v := range m.GetExternalColumns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("externalColumns", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isPersonalSite", m.GetIsPersonalSite())
        if err != nil {
            return err
        }
    }
    if m.GetItems() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetItems()))
        for i, v := range m.GetItems() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("items", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLists() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLists()))
        for i, v := range m.GetLists() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("lists", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("onenote", m.GetOnenote())
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPages()))
        for i, v := range m.GetPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("pages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPermissions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPermissions()))
        for i, v := range m.GetPermissions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("permissions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("root", m.GetRoot())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharepointIds", m.GetSharepointIds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("siteCollection", m.GetSiteCollection())
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
    return nil
}
// SetAnalytics sets the analytics property value. Analytics about the view activities that took place on this site.
func (m *Site) SetAnalytics(value ItemAnalyticsable)() {
    err := m.GetBackingStore().Set("analytics", value)
    if err != nil {
        panic(err)
    }
}
// SetColumns sets the columns property value. The collection of column definitions reusable across lists under this site.
func (m *Site) SetColumns(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("columns", value)
    if err != nil {
        panic(err)
    }
}
// SetContentTypes sets the contentTypes property value. The collection of content types defined for this site.
func (m *Site) SetContentTypes(value []ContentTypeable)() {
    err := m.GetBackingStore().Set("contentTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The full title for the site. Read-only.
func (m *Site) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDrive sets the drive property value. The default drive (document library) for this site.
func (m *Site) SetDrive(value Driveable)() {
    err := m.GetBackingStore().Set("drive", value)
    if err != nil {
        panic(err)
    }
}
// SetDrives sets the drives property value. The collection of drives (document libraries) under this site.
func (m *Site) SetDrives(value []Driveable)() {
    err := m.GetBackingStore().Set("drives", value)
    if err != nil {
        panic(err)
    }
}
// SetError sets the error property value. The error property
func (m *Site) SetError(value PublicErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalColumns sets the externalColumns property value. The externalColumns property
func (m *Site) SetExternalColumns(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("externalColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetIsPersonalSite sets the isPersonalSite property value. Identifies whether the site is personal or not. Read-only.
func (m *Site) SetIsPersonalSite(value *bool)() {
    err := m.GetBackingStore().Set("isPersonalSite", value)
    if err != nil {
        panic(err)
    }
}
// SetItems sets the items property value. Used to address any item contained in this site. This collection can't be enumerated.
func (m *Site) SetItems(value []BaseItemable)() {
    err := m.GetBackingStore().Set("items", value)
    if err != nil {
        panic(err)
    }
}
// SetLists sets the lists property value. The collection of lists under this site.
func (m *Site) SetLists(value []Listable)() {
    err := m.GetBackingStore().Set("lists", value)
    if err != nil {
        panic(err)
    }
}
// SetOnenote sets the onenote property value. Calls the OneNote service for notebook related operations.
func (m *Site) SetOnenote(value Onenoteable)() {
    err := m.GetBackingStore().Set("onenote", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The collection of long-running operations on the site.
func (m *Site) SetOperations(value []RichLongRunningOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPages sets the pages property value. The collection of pages in the baseSitePages list in this site.
func (m *Site) SetPages(value []BaseSitePageable)() {
    err := m.GetBackingStore().Set("pages", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissions sets the permissions property value. The permissions associated with the site. Nullable.
func (m *Site) SetPermissions(value []Permissionable)() {
    err := m.GetBackingStore().Set("permissions", value)
    if err != nil {
        panic(err)
    }
}
// SetRoot sets the root property value. If present, provides the root site in the site collection. Read-only.
func (m *Site) SetRoot(value Rootable)() {
    err := m.GetBackingStore().Set("root", value)
    if err != nil {
        panic(err)
    }
}
// SetSharepointIds sets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *Site) SetSharepointIds(value SharepointIdsable)() {
    err := m.GetBackingStore().Set("sharepointIds", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteCollection sets the siteCollection property value. Provides details about the site's site collection. Available only on the root site. Read-only.
func (m *Site) SetSiteCollection(value SiteCollectionable)() {
    err := m.GetBackingStore().Set("siteCollection", value)
    if err != nil {
        panic(err)
    }
}
// SetSites sets the sites property value. The collection of the sub-sites under this site.
func (m *Site) SetSites(value []Siteable)() {
    err := m.GetBackingStore().Set("sites", value)
    if err != nil {
        panic(err)
    }
}
type Siteable interface {
    BaseItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAnalytics()(ItemAnalyticsable)
    GetColumns()([]ColumnDefinitionable)
    GetContentTypes()([]ContentTypeable)
    GetDisplayName()(*string)
    GetDrive()(Driveable)
    GetDrives()([]Driveable)
    GetError()(PublicErrorable)
    GetExternalColumns()([]ColumnDefinitionable)
    GetIsPersonalSite()(*bool)
    GetItems()([]BaseItemable)
    GetLists()([]Listable)
    GetOnenote()(Onenoteable)
    GetOperations()([]RichLongRunningOperationable)
    GetPages()([]BaseSitePageable)
    GetPermissions()([]Permissionable)
    GetRoot()(Rootable)
    GetSharepointIds()(SharepointIdsable)
    GetSiteCollection()(SiteCollectionable)
    GetSites()([]Siteable)
    SetAnalytics(value ItemAnalyticsable)()
    SetColumns(value []ColumnDefinitionable)()
    SetContentTypes(value []ContentTypeable)()
    SetDisplayName(value *string)()
    SetDrive(value Driveable)()
    SetDrives(value []Driveable)()
    SetError(value PublicErrorable)()
    SetExternalColumns(value []ColumnDefinitionable)()
    SetIsPersonalSite(value *bool)()
    SetItems(value []BaseItemable)()
    SetLists(value []Listable)()
    SetOnenote(value Onenoteable)()
    SetOperations(value []RichLongRunningOperationable)()
    SetPages(value []BaseSitePageable)()
    SetPermissions(value []Permissionable)()
    SetRoot(value Rootable)()
    SetSharepointIds(value SharepointIdsable)()
    SetSiteCollection(value SiteCollectionable)()
    SetSites(value []Siteable)()
}
