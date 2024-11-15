package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AttributeDefinition struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAttributeDefinition instantiates a new AttributeDefinition and sets the default values.
func NewAttributeDefinition()(*AttributeDefinition) {
    m := &AttributeDefinition{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAttributeDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttributeDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttributeDefinition(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AttributeDefinition) GetAdditionalData()(map[string]any) {
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
// GetAnchor gets the anchor property value. true if the attribute should be used as the anchor for the object. Anchor attributes must have a unique value identifying an object, and must be immutable. Default is false. One, and only one, of the object's attributes must be designated as the anchor to support synchronization.
// returns a *bool when successful
func (m *AttributeDefinition) GetAnchor()(*bool) {
    val, err := m.GetBackingStore().Get("anchor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApiExpressions gets the apiExpressions property value. The apiExpressions property
// returns a []StringKeyStringValuePairable when successful
func (m *AttributeDefinition) GetApiExpressions()([]StringKeyStringValuePairable) {
    val, err := m.GetBackingStore().Get("apiExpressions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]StringKeyStringValuePairable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *AttributeDefinition) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCaseExact gets the caseExact property value. true if value of this attribute should be treated as case-sensitive. This setting affects how the synchronization engine detects changes for the attribute.
// returns a *bool when successful
func (m *AttributeDefinition) GetCaseExact()(*bool) {
    val, err := m.GetBackingStore().Get("caseExact")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefaultValue gets the defaultValue property value. The default value of the attribute.
// returns a *string when successful
func (m *AttributeDefinition) GetDefaultValue()(*string) {
    val, err := m.GetBackingStore().Get("defaultValue")
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
func (m *AttributeDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["anchor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnchor(val)
        }
        return nil
    }
    res["apiExpressions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateStringKeyStringValuePairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]StringKeyStringValuePairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(StringKeyStringValuePairable)
                }
            }
            m.SetApiExpressions(res)
        }
        return nil
    }
    res["caseExact"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaseExact(val)
        }
        return nil
    }
    res["defaultValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultValue(val)
        }
        return nil
    }
    res["flowNullValues"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFlowNullValues(val)
        }
        return nil
    }
    res["metadata"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttributeDefinitionMetadataEntryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AttributeDefinitionMetadataEntryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AttributeDefinitionMetadataEntryable)
                }
            }
            m.SetMetadata(res)
        }
        return nil
    }
    res["multivalued"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMultivalued(val)
        }
        return nil
    }
    res["mutability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMutability)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMutability(val.(*Mutability))
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
    res["referencedObjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateReferencedObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ReferencedObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ReferencedObjectable)
                }
            }
            m.SetReferencedObjects(res)
        }
        return nil
    }
    res["required"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequired(val)
        }
        return nil
    }
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAttributeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val.(*AttributeType))
        }
        return nil
    }
    return res
}
// GetFlowNullValues gets the flowNullValues property value. 'true' to allow null values for attributes.
// returns a *bool when successful
func (m *AttributeDefinition) GetFlowNullValues()(*bool) {
    val, err := m.GetBackingStore().Get("flowNullValues")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMetadata gets the metadata property value. Metadata for the given object.
// returns a []AttributeDefinitionMetadataEntryable when successful
func (m *AttributeDefinition) GetMetadata()([]AttributeDefinitionMetadataEntryable) {
    val, err := m.GetBackingStore().Get("metadata")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AttributeDefinitionMetadataEntryable)
    }
    return nil
}
// GetMultivalued gets the multivalued property value. true if an attribute can have multiple values. Default is false.
// returns a *bool when successful
func (m *AttributeDefinition) GetMultivalued()(*bool) {
    val, err := m.GetBackingStore().Get("multivalued")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMutability gets the mutability property value. The mutability property
// returns a *Mutability when successful
func (m *AttributeDefinition) GetMutability()(*Mutability) {
    val, err := m.GetBackingStore().Get("mutability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Mutability)
    }
    return nil
}
// GetName gets the name property value. Name of the attribute. Must be unique within the object definition. Not nullable.
// returns a *string when successful
func (m *AttributeDefinition) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
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
func (m *AttributeDefinition) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReferencedObjects gets the referencedObjects property value. For attributes with reference type, lists referenced objects (for example, the manager attribute would list User as the referenced object).
// returns a []ReferencedObjectable when successful
func (m *AttributeDefinition) GetReferencedObjects()([]ReferencedObjectable) {
    val, err := m.GetBackingStore().Get("referencedObjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ReferencedObjectable)
    }
    return nil
}
// GetRequired gets the required property value. true if attribute is required. Object can not be created if any of the required attributes are missing. If during synchronization, the required attribute has no value, the default value will be used. If default the value was not set, synchronization will record an error.
// returns a *bool when successful
func (m *AttributeDefinition) GetRequired()(*bool) {
    val, err := m.GetBackingStore().Get("required")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTypeEscaped gets the type property value. The type property
// returns a *AttributeType when successful
func (m *AttributeDefinition) GetTypeEscaped()(*AttributeType) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AttributeType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AttributeDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("anchor", m.GetAnchor())
        if err != nil {
            return err
        }
    }
    if m.GetApiExpressions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetApiExpressions()))
        for i, v := range m.GetApiExpressions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("apiExpressions", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("caseExact", m.GetCaseExact())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("defaultValue", m.GetDefaultValue())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("flowNullValues", m.GetFlowNullValues())
        if err != nil {
            return err
        }
    }
    if m.GetMetadata() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMetadata()))
        for i, v := range m.GetMetadata() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("metadata", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("multivalued", m.GetMultivalued())
        if err != nil {
            return err
        }
    }
    if m.GetMutability() != nil {
        cast := (*m.GetMutability()).String()
        err := writer.WriteStringValue("mutability", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
    if m.GetReferencedObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReferencedObjects()))
        for i, v := range m.GetReferencedObjects() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("referencedObjects", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("required", m.GetRequired())
        if err != nil {
            return err
        }
    }
    if m.GetTypeEscaped() != nil {
        cast := (*m.GetTypeEscaped()).String()
        err := writer.WriteStringValue("type", &cast)
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
func (m *AttributeDefinition) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAnchor sets the anchor property value. true if the attribute should be used as the anchor for the object. Anchor attributes must have a unique value identifying an object, and must be immutable. Default is false. One, and only one, of the object's attributes must be designated as the anchor to support synchronization.
func (m *AttributeDefinition) SetAnchor(value *bool)() {
    err := m.GetBackingStore().Set("anchor", value)
    if err != nil {
        panic(err)
    }
}
// SetApiExpressions sets the apiExpressions property value. The apiExpressions property
func (m *AttributeDefinition) SetApiExpressions(value []StringKeyStringValuePairable)() {
    err := m.GetBackingStore().Set("apiExpressions", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AttributeDefinition) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCaseExact sets the caseExact property value. true if value of this attribute should be treated as case-sensitive. This setting affects how the synchronization engine detects changes for the attribute.
func (m *AttributeDefinition) SetCaseExact(value *bool)() {
    err := m.GetBackingStore().Set("caseExact", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultValue sets the defaultValue property value. The default value of the attribute.
func (m *AttributeDefinition) SetDefaultValue(value *string)() {
    err := m.GetBackingStore().Set("defaultValue", value)
    if err != nil {
        panic(err)
    }
}
// SetFlowNullValues sets the flowNullValues property value. 'true' to allow null values for attributes.
func (m *AttributeDefinition) SetFlowNullValues(value *bool)() {
    err := m.GetBackingStore().Set("flowNullValues", value)
    if err != nil {
        panic(err)
    }
}
// SetMetadata sets the metadata property value. Metadata for the given object.
func (m *AttributeDefinition) SetMetadata(value []AttributeDefinitionMetadataEntryable)() {
    err := m.GetBackingStore().Set("metadata", value)
    if err != nil {
        panic(err)
    }
}
// SetMultivalued sets the multivalued property value. true if an attribute can have multiple values. Default is false.
func (m *AttributeDefinition) SetMultivalued(value *bool)() {
    err := m.GetBackingStore().Set("multivalued", value)
    if err != nil {
        panic(err)
    }
}
// SetMutability sets the mutability property value. The mutability property
func (m *AttributeDefinition) SetMutability(value *Mutability)() {
    err := m.GetBackingStore().Set("mutability", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Name of the attribute. Must be unique within the object definition. Not nullable.
func (m *AttributeDefinition) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AttributeDefinition) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetReferencedObjects sets the referencedObjects property value. For attributes with reference type, lists referenced objects (for example, the manager attribute would list User as the referenced object).
func (m *AttributeDefinition) SetReferencedObjects(value []ReferencedObjectable)() {
    err := m.GetBackingStore().Set("referencedObjects", value)
    if err != nil {
        panic(err)
    }
}
// SetRequired sets the required property value. true if attribute is required. Object can not be created if any of the required attributes are missing. If during synchronization, the required attribute has no value, the default value will be used. If default the value was not set, synchronization will record an error.
func (m *AttributeDefinition) SetRequired(value *bool)() {
    err := m.GetBackingStore().Set("required", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. The type property
func (m *AttributeDefinition) SetTypeEscaped(value *AttributeType)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
type AttributeDefinitionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAnchor()(*bool)
    GetApiExpressions()([]StringKeyStringValuePairable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCaseExact()(*bool)
    GetDefaultValue()(*string)
    GetFlowNullValues()(*bool)
    GetMetadata()([]AttributeDefinitionMetadataEntryable)
    GetMultivalued()(*bool)
    GetMutability()(*Mutability)
    GetName()(*string)
    GetOdataType()(*string)
    GetReferencedObjects()([]ReferencedObjectable)
    GetRequired()(*bool)
    GetTypeEscaped()(*AttributeType)
    SetAnchor(value *bool)()
    SetApiExpressions(value []StringKeyStringValuePairable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCaseExact(value *bool)()
    SetDefaultValue(value *string)()
    SetFlowNullValues(value *bool)()
    SetMetadata(value []AttributeDefinitionMetadataEntryable)()
    SetMultivalued(value *bool)()
    SetMutability(value *Mutability)()
    SetName(value *string)()
    SetOdataType(value *string)()
    SetReferencedObjects(value []ReferencedObjectable)()
    SetRequired(value *bool)()
    SetTypeEscaped(value *AttributeType)()
}
