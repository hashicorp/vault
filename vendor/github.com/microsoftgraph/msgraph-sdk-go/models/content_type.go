package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ContentType struct {
    Entity
}
// NewContentType instantiates a new ContentType and sets the default values.
func NewContentType()(*ContentType) {
    m := &ContentType{
        Entity: *NewEntity(),
    }
    return m
}
// CreateContentTypeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateContentTypeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewContentType(), nil
}
// GetAssociatedHubsUrls gets the associatedHubsUrls property value. List of canonical URLs for hub sites with which this content type is associated to. This will contain all hub sites where this content type is queued to be enforced or is already enforced. Enforcing a content type means that the content type is applied to the lists in the enforced sites.
// returns a []string when successful
func (m *ContentType) GetAssociatedHubsUrls()([]string) {
    val, err := m.GetBackingStore().Get("associatedHubsUrls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetBase gets the base property value. Parent contentType from which this content type is derived.
// returns a ContentTypeable when successful
func (m *ContentType) GetBase()(ContentTypeable) {
    val, err := m.GetBackingStore().Get("base")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContentTypeable)
    }
    return nil
}
// GetBaseTypes gets the baseTypes property value. The collection of content types that are ancestors of this content type.
// returns a []ContentTypeable when successful
func (m *ContentType) GetBaseTypes()([]ContentTypeable) {
    val, err := m.GetBackingStore().Get("baseTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContentTypeable)
    }
    return nil
}
// GetColumnLinks gets the columnLinks property value. The collection of columns that are required by this content type.
// returns a []ColumnLinkable when successful
func (m *ContentType) GetColumnLinks()([]ColumnLinkable) {
    val, err := m.GetBackingStore().Get("columnLinks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnLinkable)
    }
    return nil
}
// GetColumnPositions gets the columnPositions property value. Column order information in a content type.
// returns a []ColumnDefinitionable when successful
func (m *ContentType) GetColumnPositions()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("columnPositions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetColumns gets the columns property value. The collection of column definitions for this content type.
// returns a []ColumnDefinitionable when successful
func (m *ContentType) GetColumns()([]ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("columns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ColumnDefinitionable)
    }
    return nil
}
// GetDescription gets the description property value. The descriptive text for the item.
// returns a *string when successful
func (m *ContentType) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDocumentSet gets the documentSet property value. Document Set metadata.
// returns a DocumentSetable when successful
func (m *ContentType) GetDocumentSet()(DocumentSetable) {
    val, err := m.GetBackingStore().Get("documentSet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DocumentSetable)
    }
    return nil
}
// GetDocumentTemplate gets the documentTemplate property value. Document template metadata. To make sure that documents have consistent content across a site and its subsites, you can associate a Word, Excel, or PowerPoint template with a site content type.
// returns a DocumentSetContentable when successful
func (m *ContentType) GetDocumentTemplate()(DocumentSetContentable) {
    val, err := m.GetBackingStore().Get("documentTemplate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DocumentSetContentable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ContentType) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["associatedHubsUrls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetAssociatedHubsUrls(res)
        }
        return nil
    }
    res["base"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContentTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBase(val.(ContentTypeable))
        }
        return nil
    }
    res["baseTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetBaseTypes(res)
        }
        return nil
    }
    res["columnLinks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateColumnLinkFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ColumnLinkable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ColumnLinkable)
                }
            }
            m.SetColumnLinks(res)
        }
        return nil
    }
    res["columnPositions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetColumnPositions(res)
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
    res["documentSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDocumentSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDocumentSet(val.(DocumentSetable))
        }
        return nil
    }
    res["documentTemplate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDocumentSetContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDocumentTemplate(val.(DocumentSetContentable))
        }
        return nil
    }
    res["group"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroup(val)
        }
        return nil
    }
    res["hidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHidden(val)
        }
        return nil
    }
    res["inheritedFrom"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemReferenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInheritedFrom(val.(ItemReferenceable))
        }
        return nil
    }
    res["isBuiltIn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsBuiltIn(val)
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["order"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContentTypeOrderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrder(val.(ContentTypeOrderable))
        }
        return nil
    }
    res["parentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentId(val)
        }
        return nil
    }
    res["propagateChanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPropagateChanges(val)
        }
        return nil
    }
    res["readOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReadOnly(val)
        }
        return nil
    }
    res["sealed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSealed(val)
        }
        return nil
    }
    return res
}
// GetGroup gets the group property value. The name of the group this content type belongs to. Helps organize related content types.
// returns a *string when successful
func (m *ContentType) GetGroup()(*string) {
    val, err := m.GetBackingStore().Get("group")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHidden gets the hidden property value. Indicates whether the content type is hidden in the list's 'New' menu.
// returns a *bool when successful
func (m *ContentType) GetHidden()(*bool) {
    val, err := m.GetBackingStore().Get("hidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInheritedFrom gets the inheritedFrom property value. If this content type is inherited from another scope (like a site), provides a reference to the item where the content type is defined.
// returns a ItemReferenceable when successful
func (m *ContentType) GetInheritedFrom()(ItemReferenceable) {
    val, err := m.GetBackingStore().Get("inheritedFrom")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemReferenceable)
    }
    return nil
}
// GetIsBuiltIn gets the isBuiltIn property value. Specifies if a content type is a built-in content type.
// returns a *bool when successful
func (m *ContentType) GetIsBuiltIn()(*bool) {
    val, err := m.GetBackingStore().Get("isBuiltIn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetName gets the name property value. The name of the content type.
// returns a *string when successful
func (m *ContentType) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrder gets the order property value. Specifies the order in which the content type appears in the selection UI.
// returns a ContentTypeOrderable when successful
func (m *ContentType) GetOrder()(ContentTypeOrderable) {
    val, err := m.GetBackingStore().Get("order")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContentTypeOrderable)
    }
    return nil
}
// GetParentId gets the parentId property value. The unique identifier of the content type.
// returns a *string when successful
func (m *ContentType) GetParentId()(*string) {
    val, err := m.GetBackingStore().Get("parentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPropagateChanges gets the propagateChanges property value. If true, any changes made to the content type are pushed to inherited content types and lists that implement the content type.
// returns a *bool when successful
func (m *ContentType) GetPropagateChanges()(*bool) {
    val, err := m.GetBackingStore().Get("propagateChanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetReadOnly gets the readOnly property value. If true, the content type can't be modified unless this value is first set to false.
// returns a *bool when successful
func (m *ContentType) GetReadOnly()(*bool) {
    val, err := m.GetBackingStore().Get("readOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSealed gets the sealed property value. If true, the content type can't be modified by users or through push-down operations. Only site collection administrators can seal or unseal content types.
// returns a *bool when successful
func (m *ContentType) GetSealed()(*bool) {
    val, err := m.GetBackingStore().Get("sealed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ContentType) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssociatedHubsUrls() != nil {
        err = writer.WriteCollectionOfStringValues("associatedHubsUrls", m.GetAssociatedHubsUrls())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("base", m.GetBase())
        if err != nil {
            return err
        }
    }
    if m.GetBaseTypes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBaseTypes()))
        for i, v := range m.GetBaseTypes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("baseTypes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetColumnLinks() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetColumnLinks()))
        for i, v := range m.GetColumnLinks() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("columnLinks", cast)
        if err != nil {
            return err
        }
    }
    if m.GetColumnPositions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetColumnPositions()))
        for i, v := range m.GetColumnPositions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("columnPositions", cast)
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
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("documentSet", m.GetDocumentSet())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("documentTemplate", m.GetDocumentTemplate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("group", m.GetGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hidden", m.GetHidden())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("inheritedFrom", m.GetInheritedFrom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isBuiltIn", m.GetIsBuiltIn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("order", m.GetOrder())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentId", m.GetParentId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("propagateChanges", m.GetPropagateChanges())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("readOnly", m.GetReadOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("sealed", m.GetSealed())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssociatedHubsUrls sets the associatedHubsUrls property value. List of canonical URLs for hub sites with which this content type is associated to. This will contain all hub sites where this content type is queued to be enforced or is already enforced. Enforcing a content type means that the content type is applied to the lists in the enforced sites.
func (m *ContentType) SetAssociatedHubsUrls(value []string)() {
    err := m.GetBackingStore().Set("associatedHubsUrls", value)
    if err != nil {
        panic(err)
    }
}
// SetBase sets the base property value. Parent contentType from which this content type is derived.
func (m *ContentType) SetBase(value ContentTypeable)() {
    err := m.GetBackingStore().Set("base", value)
    if err != nil {
        panic(err)
    }
}
// SetBaseTypes sets the baseTypes property value. The collection of content types that are ancestors of this content type.
func (m *ContentType) SetBaseTypes(value []ContentTypeable)() {
    err := m.GetBackingStore().Set("baseTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetColumnLinks sets the columnLinks property value. The collection of columns that are required by this content type.
func (m *ContentType) SetColumnLinks(value []ColumnLinkable)() {
    err := m.GetBackingStore().Set("columnLinks", value)
    if err != nil {
        panic(err)
    }
}
// SetColumnPositions sets the columnPositions property value. Column order information in a content type.
func (m *ContentType) SetColumnPositions(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("columnPositions", value)
    if err != nil {
        panic(err)
    }
}
// SetColumns sets the columns property value. The collection of column definitions for this content type.
func (m *ContentType) SetColumns(value []ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("columns", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The descriptive text for the item.
func (m *ContentType) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDocumentSet sets the documentSet property value. Document Set metadata.
func (m *ContentType) SetDocumentSet(value DocumentSetable)() {
    err := m.GetBackingStore().Set("documentSet", value)
    if err != nil {
        panic(err)
    }
}
// SetDocumentTemplate sets the documentTemplate property value. Document template metadata. To make sure that documents have consistent content across a site and its subsites, you can associate a Word, Excel, or PowerPoint template with a site content type.
func (m *ContentType) SetDocumentTemplate(value DocumentSetContentable)() {
    err := m.GetBackingStore().Set("documentTemplate", value)
    if err != nil {
        panic(err)
    }
}
// SetGroup sets the group property value. The name of the group this content type belongs to. Helps organize related content types.
func (m *ContentType) SetGroup(value *string)() {
    err := m.GetBackingStore().Set("group", value)
    if err != nil {
        panic(err)
    }
}
// SetHidden sets the hidden property value. Indicates whether the content type is hidden in the list's 'New' menu.
func (m *ContentType) SetHidden(value *bool)() {
    err := m.GetBackingStore().Set("hidden", value)
    if err != nil {
        panic(err)
    }
}
// SetInheritedFrom sets the inheritedFrom property value. If this content type is inherited from another scope (like a site), provides a reference to the item where the content type is defined.
func (m *ContentType) SetInheritedFrom(value ItemReferenceable)() {
    err := m.GetBackingStore().Set("inheritedFrom", value)
    if err != nil {
        panic(err)
    }
}
// SetIsBuiltIn sets the isBuiltIn property value. Specifies if a content type is a built-in content type.
func (m *ContentType) SetIsBuiltIn(value *bool)() {
    err := m.GetBackingStore().Set("isBuiltIn", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The name of the content type.
func (m *ContentType) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetOrder sets the order property value. Specifies the order in which the content type appears in the selection UI.
func (m *ContentType) SetOrder(value ContentTypeOrderable)() {
    err := m.GetBackingStore().Set("order", value)
    if err != nil {
        panic(err)
    }
}
// SetParentId sets the parentId property value. The unique identifier of the content type.
func (m *ContentType) SetParentId(value *string)() {
    err := m.GetBackingStore().Set("parentId", value)
    if err != nil {
        panic(err)
    }
}
// SetPropagateChanges sets the propagateChanges property value. If true, any changes made to the content type are pushed to inherited content types and lists that implement the content type.
func (m *ContentType) SetPropagateChanges(value *bool)() {
    err := m.GetBackingStore().Set("propagateChanges", value)
    if err != nil {
        panic(err)
    }
}
// SetReadOnly sets the readOnly property value. If true, the content type can't be modified unless this value is first set to false.
func (m *ContentType) SetReadOnly(value *bool)() {
    err := m.GetBackingStore().Set("readOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetSealed sets the sealed property value. If true, the content type can't be modified by users or through push-down operations. Only site collection administrators can seal or unseal content types.
func (m *ContentType) SetSealed(value *bool)() {
    err := m.GetBackingStore().Set("sealed", value)
    if err != nil {
        panic(err)
    }
}
type ContentTypeable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssociatedHubsUrls()([]string)
    GetBase()(ContentTypeable)
    GetBaseTypes()([]ContentTypeable)
    GetColumnLinks()([]ColumnLinkable)
    GetColumnPositions()([]ColumnDefinitionable)
    GetColumns()([]ColumnDefinitionable)
    GetDescription()(*string)
    GetDocumentSet()(DocumentSetable)
    GetDocumentTemplate()(DocumentSetContentable)
    GetGroup()(*string)
    GetHidden()(*bool)
    GetInheritedFrom()(ItemReferenceable)
    GetIsBuiltIn()(*bool)
    GetName()(*string)
    GetOrder()(ContentTypeOrderable)
    GetParentId()(*string)
    GetPropagateChanges()(*bool)
    GetReadOnly()(*bool)
    GetSealed()(*bool)
    SetAssociatedHubsUrls(value []string)()
    SetBase(value ContentTypeable)()
    SetBaseTypes(value []ContentTypeable)()
    SetColumnLinks(value []ColumnLinkable)()
    SetColumnPositions(value []ColumnDefinitionable)()
    SetColumns(value []ColumnDefinitionable)()
    SetDescription(value *string)()
    SetDocumentSet(value DocumentSetable)()
    SetDocumentTemplate(value DocumentSetContentable)()
    SetGroup(value *string)()
    SetHidden(value *bool)()
    SetInheritedFrom(value ItemReferenceable)()
    SetIsBuiltIn(value *bool)()
    SetName(value *string)()
    SetOrder(value ContentTypeOrderable)()
    SetParentId(value *string)()
    SetPropagateChanges(value *bool)()
    SetReadOnly(value *bool)()
    SetSealed(value *bool)()
}
