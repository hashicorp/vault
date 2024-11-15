package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type DocumentSet struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDocumentSet instantiates a new DocumentSet and sets the default values.
func NewDocumentSet()(*DocumentSet) {
    m := &DocumentSet{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDocumentSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDocumentSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDocumentSet(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DocumentSet) GetAdditionalData()(map[string]any) {
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
// GetAllowedContentTypes gets the allowedContentTypes property value. Content types allowed in document set.
// returns a []ContentTypeInfoable when successful
func (m *DocumentSet) GetAllowedContentTypes()([]ContentTypeInfoable) {
    val, err := m.GetBackingStore().Get("allowedContentTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContentTypeInfoable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *DocumentSet) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDefaultContents gets the defaultContents property value. Default contents of document set.
// returns a []DocumentSetContentable when successful
func (m *DocumentSet) GetDefaultContents()([]DocumentSetContentable) {
    val, err := m.GetBackingStore().Get("defaultContents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DocumentSetContentable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DocumentSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowedContentTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContentTypeInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContentTypeInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContentTypeInfoable)
                }
            }
            m.SetAllowedContentTypes(res)
        }
        return nil
    }
    res["defaultContents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDocumentSetContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DocumentSetContentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DocumentSetContentable)
                }
            }
            m.SetDefaultContents(res)
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
    res["propagateWelcomePageChanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPropagateWelcomePageChanges(val)
        }
        return nil
    }
    res["sharedColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSharedColumns(res)
        }
        return nil
    }
    res["shouldPrefixNameToFile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShouldPrefixNameToFile(val)
        }
        return nil
    }
    res["welcomePageColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetWelcomePageColumns(res)
        }
        return nil
    }
    res["welcomePageUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWelcomePageUrl(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DocumentSet) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPropagateWelcomePageChanges gets the propagateWelcomePageChanges property value. Specifies whether to push welcome page changes to inherited content types.
// returns a *bool when successful
func (m *DocumentSet) GetPropagateWelcomePageChanges()(*bool) {
    val, err := m.GetBackingStore().Get("propagateWelcomePageChanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSharedColumns gets the sharedColumns property value. The sharedColumns property
// returns a []ColumnDefinitionable when successful
func (m *DocumentSet) GetSharedColumns()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("sharedColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetShouldPrefixNameToFile gets the shouldPrefixNameToFile property value. Indicates whether to add the name of the document set to each file name.
// returns a *bool when successful
func (m *DocumentSet) GetShouldPrefixNameToFile()(*bool) {
    val, err := m.GetBackingStore().Get("shouldPrefixNameToFile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWelcomePageColumns gets the welcomePageColumns property value. The welcomePageColumns property
// returns a []ColumnDefinitionable when successful
func (m *DocumentSet) GetWelcomePageColumns()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("welcomePageColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetWelcomePageUrl gets the welcomePageUrl property value. Welcome page absolute URL.
// returns a *string when successful
func (m *DocumentSet) GetWelcomePageUrl()(*string) {
    val, err := m.GetBackingStore().Get("welcomePageUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DocumentSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAllowedContentTypes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAllowedContentTypes()))
        for i, v := range m.GetAllowedContentTypes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("allowedContentTypes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDefaultContents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDefaultContents()))
        for i, v := range m.GetDefaultContents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("defaultContents", cast)
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
        err := writer.WriteBoolValue("propagateWelcomePageChanges", m.GetPropagateWelcomePageChanges())
        if err != nil {
            return err
        }
    }
    if m.GetSharedColumns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharedColumns()))
        for i, v := range m.GetSharedColumns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("sharedColumns", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("shouldPrefixNameToFile", m.GetShouldPrefixNameToFile())
        if err != nil {
            return err
        }
    }
    if m.GetWelcomePageColumns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWelcomePageColumns()))
        for i, v := range m.GetWelcomePageColumns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("welcomePageColumns", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("welcomePageUrl", m.GetWelcomePageUrl())
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
func (m *DocumentSet) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedContentTypes sets the allowedContentTypes property value. Content types allowed in document set.
func (m *DocumentSet) SetAllowedContentTypes(value []ContentTypeInfoable)() {
    err := m.GetBackingStore().Set("allowedContentTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DocumentSet) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDefaultContents sets the defaultContents property value. Default contents of document set.
func (m *DocumentSet) SetDefaultContents(value []DocumentSetContentable)() {
    err := m.GetBackingStore().Set("defaultContents", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DocumentSet) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPropagateWelcomePageChanges sets the propagateWelcomePageChanges property value. Specifies whether to push welcome page changes to inherited content types.
func (m *DocumentSet) SetPropagateWelcomePageChanges(value *bool)() {
    err := m.GetBackingStore().Set("propagateWelcomePageChanges", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedColumns sets the sharedColumns property value. The sharedColumns property
func (m *DocumentSet) SetSharedColumns(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("sharedColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetShouldPrefixNameToFile sets the shouldPrefixNameToFile property value. Indicates whether to add the name of the document set to each file name.
func (m *DocumentSet) SetShouldPrefixNameToFile(value *bool)() {
    err := m.GetBackingStore().Set("shouldPrefixNameToFile", value)
    if err != nil {
        panic(err)
    }
}
// SetWelcomePageColumns sets the welcomePageColumns property value. The welcomePageColumns property
func (m *DocumentSet) SetWelcomePageColumns(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("welcomePageColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetWelcomePageUrl sets the welcomePageUrl property value. Welcome page absolute URL.
func (m *DocumentSet) SetWelcomePageUrl(value *string)() {
    err := m.GetBackingStore().Set("welcomePageUrl", value)
    if err != nil {
        panic(err)
    }
}
type DocumentSetable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedContentTypes()([]ContentTypeInfoable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDefaultContents()([]DocumentSetContentable)
    GetOdataType()(*string)
    GetPropagateWelcomePageChanges()(*bool)
    GetSharedColumns()([]ColumnDefinitionable)
    GetShouldPrefixNameToFile()(*bool)
    GetWelcomePageColumns()([]ColumnDefinitionable)
    GetWelcomePageUrl()(*string)
    SetAllowedContentTypes(value []ContentTypeInfoable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDefaultContents(value []DocumentSetContentable)()
    SetOdataType(value *string)()
    SetPropagateWelcomePageChanges(value *bool)()
    SetSharedColumns(value []ColumnDefinitionable)()
    SetShouldPrefixNameToFile(value *bool)()
    SetWelcomePageColumns(value []ColumnDefinitionable)()
    SetWelcomePageUrl(value *string)()
}
