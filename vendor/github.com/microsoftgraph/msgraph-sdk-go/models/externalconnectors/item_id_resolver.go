package externalconnectors

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ItemIdResolver struct {
    UrlToItemResolverBase
}
// NewItemIdResolver instantiates a new ItemIdResolver and sets the default values.
func NewItemIdResolver()(*ItemIdResolver) {
    m := &ItemIdResolver{
        UrlToItemResolverBase: *NewUrlToItemResolverBase(),
    }
    odataTypeValue := "#microsoft.graph.externalConnectors.itemIdResolver"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateItemIdResolverFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemIdResolverFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemIdResolver(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemIdResolver) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UrlToItemResolverBase.GetFieldDeserializers()
    res["itemId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetItemId(val)
        }
        return nil
    }
    res["urlMatchInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUrlMatchInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUrlMatchInfo(val.(UrlMatchInfoable))
        }
        return nil
    }
    return res
}
// GetItemId gets the itemId property value. Pattern that specifies how to form the ID of the external item that the URL represents. The named groups from the regular expression in urlPattern within the urlMatchInfo can be referenced by inserting the group name inside curly brackets.
// returns a *string when successful
func (m *ItemIdResolver) GetItemId()(*string) {
    val, err := m.GetBackingStore().Get("itemId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUrlMatchInfo gets the urlMatchInfo property value. Configurations to match and resolve URL.
// returns a UrlMatchInfoable when successful
func (m *ItemIdResolver) GetUrlMatchInfo()(UrlMatchInfoable) {
    val, err := m.GetBackingStore().Get("urlMatchInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UrlMatchInfoable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemIdResolver) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UrlToItemResolverBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("itemId", m.GetItemId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("urlMatchInfo", m.GetUrlMatchInfo())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetItemId sets the itemId property value. Pattern that specifies how to form the ID of the external item that the URL represents. The named groups from the regular expression in urlPattern within the urlMatchInfo can be referenced by inserting the group name inside curly brackets.
func (m *ItemIdResolver) SetItemId(value *string)() {
    err := m.GetBackingStore().Set("itemId", value)
    if err != nil {
        panic(err)
    }
}
// SetUrlMatchInfo sets the urlMatchInfo property value. Configurations to match and resolve URL.
func (m *ItemIdResolver) SetUrlMatchInfo(value UrlMatchInfoable)() {
    err := m.GetBackingStore().Set("urlMatchInfo", value)
    if err != nil {
        panic(err)
    }
}
type ItemIdResolverable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UrlToItemResolverBaseable
    GetItemId()(*string)
    GetUrlMatchInfo()(UrlMatchInfoable)
    SetItemId(value *string)()
    SetUrlMatchInfo(value UrlMatchInfoable)()
}
