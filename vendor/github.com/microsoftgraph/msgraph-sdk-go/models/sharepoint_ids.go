package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SharepointIds struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSharepointIds instantiates a new SharepointIds and sets the default values.
func NewSharepointIds()(*SharepointIds) {
    m := &SharepointIds{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSharepointIdsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharepointIdsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharepointIds(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SharepointIds) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *SharepointIds) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharepointIds) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["listId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetListId(val)
        }
        return nil
    }
    res["listItemId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetListItemId(val)
        }
        return nil
    }
    res["listItemUniqueId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetListItemUniqueId(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
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
    res["siteUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteUrl(val)
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    res["webId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebId(val)
        }
        return nil
    }
    return res
}
// GetListId gets the listId property value. The unique identifier (guid) for the item's list in SharePoint.
// returns a *string when successful
func (m *SharepointIds) GetListId()(*string) {
    val, err := m.GetBackingStore().Get("listId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetListItemId gets the listItemId property value. An integer identifier for the item within the containing list.
// returns a *string when successful
func (m *SharepointIds) GetListItemId()(*string) {
    val, err := m.GetBackingStore().Get("listItemId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetListItemUniqueId gets the listItemUniqueId property value. The unique identifier (guid) for the item within OneDrive for Business or a SharePoint site.
// returns a *string when successful
func (m *SharepointIds) GetListItemUniqueId()(*string) {
    val, err := m.GetBackingStore().Get("listItemUniqueId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SharepointIds) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSiteId gets the siteId property value. The unique identifier (guid) for the item's site collection (SPSite).
// returns a *string when successful
func (m *SharepointIds) GetSiteId()(*string) {
    val, err := m.GetBackingStore().Get("siteId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSiteUrl gets the siteUrl property value. The SharePoint URL for the site that contains the item.
// returns a *string when successful
func (m *SharepointIds) GetSiteUrl()(*string) {
    val, err := m.GetBackingStore().Get("siteUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The unique identifier (guid) for the tenancy.
// returns a *string when successful
func (m *SharepointIds) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWebId gets the webId property value. The unique identifier (guid) for the item's site (SPWeb).
// returns a *string when successful
func (m *SharepointIds) GetWebId()(*string) {
    val, err := m.GetBackingStore().Get("webId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharepointIds) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("listId", m.GetListId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("listItemId", m.GetListItemId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("listItemUniqueId", m.GetListItemUniqueId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("siteId", m.GetSiteId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("siteUrl", m.GetSiteUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("webId", m.GetWebId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SharepointIds) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SharepointIds) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetListId sets the listId property value. The unique identifier (guid) for the item's list in SharePoint.
func (m *SharepointIds) SetListId(value *string)() {
    err := m.GetBackingStore().Set("listId", value)
    if err != nil {
        panic(err)
    }
}
// SetListItemId sets the listItemId property value. An integer identifier for the item within the containing list.
func (m *SharepointIds) SetListItemId(value *string)() {
    err := m.GetBackingStore().Set("listItemId", value)
    if err != nil {
        panic(err)
    }
}
// SetListItemUniqueId sets the listItemUniqueId property value. The unique identifier (guid) for the item within OneDrive for Business or a SharePoint site.
func (m *SharepointIds) SetListItemUniqueId(value *string)() {
    err := m.GetBackingStore().Set("listItemUniqueId", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SharepointIds) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteId sets the siteId property value. The unique identifier (guid) for the item's site collection (SPSite).
func (m *SharepointIds) SetSiteId(value *string)() {
    err := m.GetBackingStore().Set("siteId", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteUrl sets the siteUrl property value. The SharePoint URL for the site that contains the item.
func (m *SharepointIds) SetSiteUrl(value *string)() {
    err := m.GetBackingStore().Set("siteUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The unique identifier (guid) for the tenancy.
func (m *SharepointIds) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetWebId sets the webId property value. The unique identifier (guid) for the item's site (SPWeb).
func (m *SharepointIds) SetWebId(value *string)() {
    err := m.GetBackingStore().Set("webId", value)
    if err != nil {
        panic(err)
    }
}
type SharepointIdsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetListId()(*string)
    GetListItemId()(*string)
    GetListItemUniqueId()(*string)
    GetOdataType()(*string)
    GetSiteId()(*string)
    GetSiteUrl()(*string)
    GetTenantId()(*string)
    GetWebId()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetListId(value *string)()
    SetListItemId(value *string)()
    SetListItemUniqueId(value *string)()
    SetOdataType(value *string)()
    SetSiteId(value *string)()
    SetSiteUrl(value *string)()
    SetTenantId(value *string)()
    SetWebId(value *string)()
}
