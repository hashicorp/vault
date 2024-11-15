package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Drive struct {
    BaseItem
}
// NewDrive instantiates a new Drive and sets the default values.
func NewDrive()(*Drive) {
    m := &Drive{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.drive"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDriveFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDriveFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDrive(), nil
}
// GetBundles gets the bundles property value. Collection of bundles (albums and multi-select-shared sets of items). Only in personal OneDrive.
// returns a []DriveItemable when successful
func (m *Drive) GetBundles()([]DriveItemable) {
    val, err := m.GetBackingStore().Get("bundles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemable)
    }
    return nil
}
// GetDriveType gets the driveType property value. Describes the type of drive represented by this resource. OneDrive personal drives return personal. OneDrive for Business returns business. SharePoint document libraries return documentLibrary. Read-only.
// returns a *string when successful
func (m *Drive) GetDriveType()(*string) {
    val, err := m.GetBackingStore().Get("driveType")
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
func (m *Drive) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["bundles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetBundles(res)
        }
        return nil
    }
    res["driveType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDriveType(val)
        }
        return nil
    }
    res["following"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFollowing(res)
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
    res["quota"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateQuotaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuota(val.(Quotaable))
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
    res["sharePointIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharepointIdsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharePointIds(val.(SharepointIdsable))
        }
        return nil
    }
    res["special"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSpecial(res)
        }
        return nil
    }
    res["system"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSystemFacetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSystem(val.(SystemFacetable))
        }
        return nil
    }
    return res
}
// GetFollowing gets the following property value. The list of items the user is following. Only in OneDrive for Business.
// returns a []DriveItemable when successful
func (m *Drive) GetFollowing()([]DriveItemable) {
    val, err := m.GetBackingStore().Get("following")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemable)
    }
    return nil
}
// GetItems gets the items property value. All items contained in the drive. Read-only. Nullable.
// returns a []DriveItemable when successful
func (m *Drive) GetItems()([]DriveItemable) {
    val, err := m.GetBackingStore().Get("items")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemable)
    }
    return nil
}
// GetList gets the list property value. For drives in SharePoint, the underlying document library list. Read-only. Nullable.
// returns a Listable when successful
func (m *Drive) GetList()(Listable) {
    val, err := m.GetBackingStore().Get("list")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Listable)
    }
    return nil
}
// GetOwner gets the owner property value. Optional. The user account that owns the drive. Read-only.
// returns a IdentitySetable when successful
func (m *Drive) GetOwner()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("owner")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetQuota gets the quota property value. Optional. Information about the drive's storage space quota. Read-only.
// returns a Quotaable when successful
func (m *Drive) GetQuota()(Quotaable) {
    val, err := m.GetBackingStore().Get("quota")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Quotaable)
    }
    return nil
}
// GetRoot gets the root property value. The root folder of the drive. Read-only.
// returns a DriveItemable when successful
func (m *Drive) GetRoot()(DriveItemable) {
    val, err := m.GetBackingStore().Get("root")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DriveItemable)
    }
    return nil
}
// GetSharePointIds gets the sharePointIds property value. The sharePointIds property
// returns a SharepointIdsable when successful
func (m *Drive) GetSharePointIds()(SharepointIdsable) {
    val, err := m.GetBackingStore().Get("sharePointIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharepointIdsable)
    }
    return nil
}
// GetSpecial gets the special property value. Collection of common folders available in OneDrive. Read-only. Nullable.
// returns a []DriveItemable when successful
func (m *Drive) GetSpecial()([]DriveItemable) {
    val, err := m.GetBackingStore().Get("special")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemable)
    }
    return nil
}
// GetSystem gets the system property value. If present, indicates that it's a system-managed drive. Read-only.
// returns a SystemFacetable when successful
func (m *Drive) GetSystem()(SystemFacetable) {
    val, err := m.GetBackingStore().Get("system")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SystemFacetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Drive) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBundles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBundles()))
        for i, v := range m.GetBundles() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("bundles", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("driveType", m.GetDriveType())
        if err != nil {
            return err
        }
    }
    if m.GetFollowing() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFollowing()))
        for i, v := range m.GetFollowing() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("following", cast)
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
        err = writer.WriteObjectValue("owner", m.GetOwner())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("quota", m.GetQuota())
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
        err = writer.WriteObjectValue("sharePointIds", m.GetSharePointIds())
        if err != nil {
            return err
        }
    }
    if m.GetSpecial() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSpecial()))
        for i, v := range m.GetSpecial() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("special", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("system", m.GetSystem())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBundles sets the bundles property value. Collection of bundles (albums and multi-select-shared sets of items). Only in personal OneDrive.
func (m *Drive) SetBundles(value []DriveItemable)() {
    err := m.GetBackingStore().Set("bundles", value)
    if err != nil {
        panic(err)
    }
}
// SetDriveType sets the driveType property value. Describes the type of drive represented by this resource. OneDrive personal drives return personal. OneDrive for Business returns business. SharePoint document libraries return documentLibrary. Read-only.
func (m *Drive) SetDriveType(value *string)() {
    err := m.GetBackingStore().Set("driveType", value)
    if err != nil {
        panic(err)
    }
}
// SetFollowing sets the following property value. The list of items the user is following. Only in OneDrive for Business.
func (m *Drive) SetFollowing(value []DriveItemable)() {
    err := m.GetBackingStore().Set("following", value)
    if err != nil {
        panic(err)
    }
}
// SetItems sets the items property value. All items contained in the drive. Read-only. Nullable.
func (m *Drive) SetItems(value []DriveItemable)() {
    err := m.GetBackingStore().Set("items", value)
    if err != nil {
        panic(err)
    }
}
// SetList sets the list property value. For drives in SharePoint, the underlying document library list. Read-only. Nullable.
func (m *Drive) SetList(value Listable)() {
    err := m.GetBackingStore().Set("list", value)
    if err != nil {
        panic(err)
    }
}
// SetOwner sets the owner property value. Optional. The user account that owns the drive. Read-only.
func (m *Drive) SetOwner(value IdentitySetable)() {
    err := m.GetBackingStore().Set("owner", value)
    if err != nil {
        panic(err)
    }
}
// SetQuota sets the quota property value. Optional. Information about the drive's storage space quota. Read-only.
func (m *Drive) SetQuota(value Quotaable)() {
    err := m.GetBackingStore().Set("quota", value)
    if err != nil {
        panic(err)
    }
}
// SetRoot sets the root property value. The root folder of the drive. Read-only.
func (m *Drive) SetRoot(value DriveItemable)() {
    err := m.GetBackingStore().Set("root", value)
    if err != nil {
        panic(err)
    }
}
// SetSharePointIds sets the sharePointIds property value. The sharePointIds property
func (m *Drive) SetSharePointIds(value SharepointIdsable)() {
    err := m.GetBackingStore().Set("sharePointIds", value)
    if err != nil {
        panic(err)
    }
}
// SetSpecial sets the special property value. Collection of common folders available in OneDrive. Read-only. Nullable.
func (m *Drive) SetSpecial(value []DriveItemable)() {
    err := m.GetBackingStore().Set("special", value)
    if err != nil {
        panic(err)
    }
}
// SetSystem sets the system property value. If present, indicates that it's a system-managed drive. Read-only.
func (m *Drive) SetSystem(value SystemFacetable)() {
    err := m.GetBackingStore().Set("system", value)
    if err != nil {
        panic(err)
    }
}
type Driveable interface {
    BaseItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBundles()([]DriveItemable)
    GetDriveType()(*string)
    GetFollowing()([]DriveItemable)
    GetItems()([]DriveItemable)
    GetList()(Listable)
    GetOwner()(IdentitySetable)
    GetQuota()(Quotaable)
    GetRoot()(DriveItemable)
    GetSharePointIds()(SharepointIdsable)
    GetSpecial()([]DriveItemable)
    GetSystem()(SystemFacetable)
    SetBundles(value []DriveItemable)()
    SetDriveType(value *string)()
    SetFollowing(value []DriveItemable)()
    SetItems(value []DriveItemable)()
    SetList(value Listable)()
    SetOwner(value IdentitySetable)()
    SetQuota(value Quotaable)()
    SetRoot(value DriveItemable)()
    SetSharePointIds(value SharepointIdsable)()
    SetSpecial(value []DriveItemable)()
    SetSystem(value SystemFacetable)()
}
