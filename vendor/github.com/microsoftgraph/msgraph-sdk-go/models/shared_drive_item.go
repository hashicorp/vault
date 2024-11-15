package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SharedDriveItem struct {
    BaseItem
}
// NewSharedDriveItem instantiates a new SharedDriveItem and sets the default values.
func NewSharedDriveItem()(*SharedDriveItem) {
    m := &SharedDriveItem{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.sharedDriveItem"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSharedDriveItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharedDriveItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedDriveItem(), nil
}
// GetDriveItem gets the driveItem property value. Used to access the underlying driveItem
// returns a DriveItemable when successful
func (m *SharedDriveItem) GetDriveItem()(DriveItemable) {
    val, err := m.GetBackingStore().Get("driveItem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DriveItemable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharedDriveItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["driveItem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDriveItem(val.(DriveItemable))
        }
        return nil
    }
    res["items"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DriveItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DriveItemable)
                }
            }
            m.SetItems(res)
        }
        return nil
    }
    res["list"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateListFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetList(val.(Listable))
        }
        return nil
    }
    res["listItem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetListItem(val.(ListItemable))
        }
        return nil
    }
    res["owner"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOwner(val.(IdentitySetable))
        }
        return nil
    }
    res["permission"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePermissionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermission(val.(Permissionable))
        }
        return nil
    }
    res["root"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoot(val.(DriveItemable))
        }
        return nil
    }
    res["site"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSiteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSite(val.(Siteable))
        }
        return nil
    }
    return res
}
// GetItems gets the items property value. All driveItems contained in the sharing root. This collection cannot be enumerated.
// returns a []DriveItemable when successful
func (m *SharedDriveItem) GetItems()([]DriveItemable) {
    val, err := m.GetBackingStore().Get("items")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemable)
    }
    return nil
}
// GetList gets the list property value. Used to access the underlying list
// returns a Listable when successful
func (m *SharedDriveItem) GetList()(Listable) {
    val, err := m.GetBackingStore().Get("list")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Listable)
    }
    return nil
}
// GetListItem gets the listItem property value. Used to access the underlying listItem
// returns a ListItemable when successful
func (m *SharedDriveItem) GetListItem()(ListItemable) {
    val, err := m.GetBackingStore().Get("listItem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ListItemable)
    }
    return nil
}
// GetOwner gets the owner property value. Information about the owner of the shared item being referenced.
// returns a IdentitySetable when successful
func (m *SharedDriveItem) GetOwner()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("owner")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetPermission gets the permission property value. Used to access the permission representing the underlying sharing link
// returns a Permissionable when successful
func (m *SharedDriveItem) GetPermission()(Permissionable) {
    val, err := m.GetBackingStore().Get("permission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Permissionable)
    }
    return nil
}
// GetRoot gets the root property value. Used to access the underlying driveItem. Deprecated -- use driveItem instead.
// returns a DriveItemable when successful
func (m *SharedDriveItem) GetRoot()(DriveItemable) {
    val, err := m.GetBackingStore().Get("root")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DriveItemable)
    }
    return nil
}
// GetSite gets the site property value. Used to access the underlying site
// returns a Siteable when successful
func (m *SharedDriveItem) GetSite()(Siteable) {
    val, err := m.GetBackingStore().Get("site")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Siteable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharedDriveItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("driveItem", m.GetDriveItem())
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
    {
        err = writer.WriteObjectValue("list", m.GetList())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("listItem", m.GetListItem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("owner", m.GetOwner())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("permission", m.GetPermission())
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
        err = writer.WriteObjectValue("site", m.GetSite())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDriveItem sets the driveItem property value. Used to access the underlying driveItem
func (m *SharedDriveItem) SetDriveItem(value DriveItemable)() {
    err := m.GetBackingStore().Set("driveItem", value)
    if err != nil {
        panic(err)
    }
}
// SetItems sets the items property value. All driveItems contained in the sharing root. This collection cannot be enumerated.
func (m *SharedDriveItem) SetItems(value []DriveItemable)() {
    err := m.GetBackingStore().Set("items", value)
    if err != nil {
        panic(err)
    }
}
// SetList sets the list property value. Used to access the underlying list
func (m *SharedDriveItem) SetList(value Listable)() {
    err := m.GetBackingStore().Set("list", value)
    if err != nil {
        panic(err)
    }
}
// SetListItem sets the listItem property value. Used to access the underlying listItem
func (m *SharedDriveItem) SetListItem(value ListItemable)() {
    err := m.GetBackingStore().Set("listItem", value)
    if err != nil {
        panic(err)
    }
}
// SetOwner sets the owner property value. Information about the owner of the shared item being referenced.
func (m *SharedDriveItem) SetOwner(value IdentitySetable)() {
    err := m.GetBackingStore().Set("owner", value)
    if err != nil {
        panic(err)
    }
}
// SetPermission sets the permission property value. Used to access the permission representing the underlying sharing link
func (m *SharedDriveItem) SetPermission(value Permissionable)() {
    err := m.GetBackingStore().Set("permission", value)
    if err != nil {
        panic(err)
    }
}
// SetRoot sets the root property value. Used to access the underlying driveItem. Deprecated -- use driveItem instead.
func (m *SharedDriveItem) SetRoot(value DriveItemable)() {
    err := m.GetBackingStore().Set("root", value)
    if err != nil {
        panic(err)
    }
}
// SetSite sets the site property value. Used to access the underlying site
func (m *SharedDriveItem) SetSite(value Siteable)() {
    err := m.GetBackingStore().Set("site", value)
    if err != nil {
        panic(err)
    }
}
type SharedDriveItemable interface {
    BaseItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDriveItem()(DriveItemable)
    GetItems()([]DriveItemable)
    GetList()(Listable)
    GetListItem()(ListItemable)
    GetOwner()(IdentitySetable)
    GetPermission()(Permissionable)
    GetRoot()(DriveItemable)
    GetSite()(Siteable)
    SetDriveItem(value DriveItemable)()
    SetItems(value []DriveItemable)()
    SetList(value Listable)()
    SetListItem(value ListItemable)()
    SetOwner(value IdentitySetable)()
    SetPermission(value Permissionable)()
    SetRoot(value DriveItemable)()
    SetSite(value Siteable)()
}
