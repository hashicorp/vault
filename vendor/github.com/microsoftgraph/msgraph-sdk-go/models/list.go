package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type List struct {
    BaseItem
}
// NewList instantiates a new List and sets the default values.
func NewList()(*List) {
    m := &List{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.list"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateListFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateListFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewList(), nil
}
// GetColumns gets the columns property value. The collection of field definitions for this list.
// returns a []ColumnDefinitionable when successful
func (m *List) GetColumns()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("columns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetContentTypes gets the contentTypes property value. The collection of content types present in this list.
// returns a []ContentTypeable when successful
func (m *List) GetContentTypes()([]ContentTypeable) {
    val, err := m.GetBackingStore().Get("contentTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContentTypeable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The displayable title of the list.
// returns a *string when successful
func (m *List) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDrive gets the drive property value. Allows access to the list as a drive resource with driveItems. Only present on document libraries.
// returns a Driveable when successful
func (m *List) GetDrive()(Driveable) {
    val, err := m.GetBackingStore().Get("drive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Driveable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *List) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
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
    res["items"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ListItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ListItemable)
                }
            }
            m.SetItems(res)
        }
        return nil
    }
    res["list"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateListInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetList(val.(ListInfoable))
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
    res["subscriptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubscriptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Subscriptionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Subscriptionable)
                }
            }
            m.SetSubscriptions(res)
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
// GetItems gets the items property value. All items contained in the list.
// returns a []ListItemable when successful
func (m *List) GetItems()([]ListItemable) {
    val, err := m.GetBackingStore().Get("items")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ListItemable)
    }
    return nil
}
// GetList gets the list property value. Contains more details about the list.
// returns a ListInfoable when successful
func (m *List) GetList()(ListInfoable) {
    val, err := m.GetBackingStore().Get("list")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ListInfoable)
    }
    return nil
}
// GetOperations gets the operations property value. The collection of long-running operations on the list.
// returns a []RichLongRunningOperationable when successful
func (m *List) GetOperations()([]RichLongRunningOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RichLongRunningOperationable)
    }
    return nil
}
// GetSharepointIds gets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
// returns a SharepointIdsable when successful
func (m *List) GetSharepointIds()(SharepointIdsable) {
    val, err := m.GetBackingStore().Get("sharepointIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharepointIdsable)
    }
    return nil
}
// GetSubscriptions gets the subscriptions property value. The set of subscriptions on the list.
// returns a []Subscriptionable when successful
func (m *List) GetSubscriptions()([]Subscriptionable) {
    val, err := m.GetBackingStore().Get("subscriptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Subscriptionable)
    }
    return nil
}
// GetSystem gets the system property value. If present, indicates that the list is system-managed. Read-only.
// returns a SystemFacetable when successful
func (m *List) GetSystem()(SystemFacetable) {
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
func (m *List) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
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
    {
        err = writer.WriteObjectValue("sharepointIds", m.GetSharepointIds())
        if err != nil {
            return err
        }
    }
    if m.GetSubscriptions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSubscriptions()))
        for i, v := range m.GetSubscriptions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("subscriptions", cast)
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
// SetColumns sets the columns property value. The collection of field definitions for this list.
func (m *List) SetColumns(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("columns", value)
    if err != nil {
        panic(err)
    }
}
// SetContentTypes sets the contentTypes property value. The collection of content types present in this list.
func (m *List) SetContentTypes(value []ContentTypeable)() {
    err := m.GetBackingStore().Set("contentTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The displayable title of the list.
func (m *List) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDrive sets the drive property value. Allows access to the list as a drive resource with driveItems. Only present on document libraries.
func (m *List) SetDrive(value Driveable)() {
    err := m.GetBackingStore().Set("drive", value)
    if err != nil {
        panic(err)
    }
}
// SetItems sets the items property value. All items contained in the list.
func (m *List) SetItems(value []ListItemable)() {
    err := m.GetBackingStore().Set("items", value)
    if err != nil {
        panic(err)
    }
}
// SetList sets the list property value. Contains more details about the list.
func (m *List) SetList(value ListInfoable)() {
    err := m.GetBackingStore().Set("list", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The collection of long-running operations on the list.
func (m *List) SetOperations(value []RichLongRunningOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetSharepointIds sets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *List) SetSharepointIds(value SharepointIdsable)() {
    err := m.GetBackingStore().Set("sharepointIds", value)
    if err != nil {
        panic(err)
    }
}
// SetSubscriptions sets the subscriptions property value. The set of subscriptions on the list.
func (m *List) SetSubscriptions(value []Subscriptionable)() {
    err := m.GetBackingStore().Set("subscriptions", value)
    if err != nil {
        panic(err)
    }
}
// SetSystem sets the system property value. If present, indicates that the list is system-managed. Read-only.
func (m *List) SetSystem(value SystemFacetable)() {
    err := m.GetBackingStore().Set("system", value)
    if err != nil {
        panic(err)
    }
}
type Listable interface {
    BaseItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetColumns()([]ColumnDefinitionable)
    GetContentTypes()([]ContentTypeable)
    GetDisplayName()(*string)
    GetDrive()(Driveable)
    GetItems()([]ListItemable)
    GetList()(ListInfoable)
    GetOperations()([]RichLongRunningOperationable)
    GetSharepointIds()(SharepointIdsable)
    GetSubscriptions()([]Subscriptionable)
    GetSystem()(SystemFacetable)
    SetColumns(value []ColumnDefinitionable)()
    SetContentTypes(value []ContentTypeable)()
    SetDisplayName(value *string)()
    SetDrive(value Driveable)()
    SetItems(value []ListItemable)()
    SetList(value ListInfoable)()
    SetOperations(value []RichLongRunningOperationable)()
    SetSharepointIds(value SharepointIdsable)()
    SetSubscriptions(value []Subscriptionable)()
    SetSystem(value SystemFacetable)()
}
