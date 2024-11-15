package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AttributeMapping struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAttributeMapping instantiates a new AttributeMapping and sets the default values.
func NewAttributeMapping()(*AttributeMapping) {
    m := &AttributeMapping{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAttributeMappingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttributeMappingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttributeMapping(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AttributeMapping) GetAdditionalData()(map[string]any) {
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
func (m *AttributeMapping) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDefaultValue gets the defaultValue property value. Default value to be used in case the source property was evaluated to null. Optional.
// returns a *string when successful
func (m *AttributeMapping) GetDefaultValue()(*string) {
    val, err := m.GetBackingStore().Get("defaultValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExportMissingReferences gets the exportMissingReferences property value. For internal use only.
// returns a *bool when successful
func (m *AttributeMapping) GetExportMissingReferences()(*bool) {
    val, err := m.GetBackingStore().Get("exportMissingReferences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AttributeMapping) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["exportMissingReferences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExportMissingReferences(val)
        }
        return nil
    }
    res["flowBehavior"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAttributeFlowBehavior)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFlowBehavior(val.(*AttributeFlowBehavior))
        }
        return nil
    }
    res["flowType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAttributeFlowType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFlowType(val.(*AttributeFlowType))
        }
        return nil
    }
    res["matchingPriority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMatchingPriority(val)
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
    res["source"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAttributeMappingSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSource(val.(AttributeMappingSourceable))
        }
        return nil
    }
    res["targetAttributeName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetAttributeName(val)
        }
        return nil
    }
    return res
}
// GetFlowBehavior gets the flowBehavior property value. The flowBehavior property
// returns a *AttributeFlowBehavior when successful
func (m *AttributeMapping) GetFlowBehavior()(*AttributeFlowBehavior) {
    val, err := m.GetBackingStore().Get("flowBehavior")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AttributeFlowBehavior)
    }
    return nil
}
// GetFlowType gets the flowType property value. The flowType property
// returns a *AttributeFlowType when successful
func (m *AttributeMapping) GetFlowType()(*AttributeFlowType) {
    val, err := m.GetBackingStore().Get("flowType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AttributeFlowType)
    }
    return nil
}
// GetMatchingPriority gets the matchingPriority property value. If higher than 0, this attribute will be used to perform an initial match of the objects between source and target directories. The synchronization engine will try to find the matching object using attribute with lowest value of matching priority first. If not found, the attribute with the next matching priority will be used, and so on a until match is found or no more matching attributes are left. Only attributes that are expected to have unique values, such as email, should be used as matching attributes.
// returns a *int32 when successful
func (m *AttributeMapping) GetMatchingPriority()(*int32) {
    val, err := m.GetBackingStore().Get("matchingPriority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AttributeMapping) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSource gets the source property value. Defines how a value should be extracted (or transformed) from the source object.
// returns a AttributeMappingSourceable when successful
func (m *AttributeMapping) GetSource()(AttributeMappingSourceable) {
    val, err := m.GetBackingStore().Get("source")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AttributeMappingSourceable)
    }
    return nil
}
// GetTargetAttributeName gets the targetAttributeName property value. Name of the attribute on the target object.
// returns a *string when successful
func (m *AttributeMapping) GetTargetAttributeName()(*string) {
    val, err := m.GetBackingStore().Get("targetAttributeName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AttributeMapping) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("defaultValue", m.GetDefaultValue())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("exportMissingReferences", m.GetExportMissingReferences())
        if err != nil {
            return err
        }
    }
    if m.GetFlowBehavior() != nil {
        cast := (*m.GetFlowBehavior()).String()
        err := writer.WriteStringValue("flowBehavior", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFlowType() != nil {
        cast := (*m.GetFlowType()).String()
        err := writer.WriteStringValue("flowType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("matchingPriority", m.GetMatchingPriority())
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
        err := writer.WriteObjectValue("source", m.GetSource())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("targetAttributeName", m.GetTargetAttributeName())
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
func (m *AttributeMapping) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AttributeMapping) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDefaultValue sets the defaultValue property value. Default value to be used in case the source property was evaluated to null. Optional.
func (m *AttributeMapping) SetDefaultValue(value *string)() {
    err := m.GetBackingStore().Set("defaultValue", value)
    if err != nil {
        panic(err)
    }
}
// SetExportMissingReferences sets the exportMissingReferences property value. For internal use only.
func (m *AttributeMapping) SetExportMissingReferences(value *bool)() {
    err := m.GetBackingStore().Set("exportMissingReferences", value)
    if err != nil {
        panic(err)
    }
}
// SetFlowBehavior sets the flowBehavior property value. The flowBehavior property
func (m *AttributeMapping) SetFlowBehavior(value *AttributeFlowBehavior)() {
    err := m.GetBackingStore().Set("flowBehavior", value)
    if err != nil {
        panic(err)
    }
}
// SetFlowType sets the flowType property value. The flowType property
func (m *AttributeMapping) SetFlowType(value *AttributeFlowType)() {
    err := m.GetBackingStore().Set("flowType", value)
    if err != nil {
        panic(err)
    }
}
// SetMatchingPriority sets the matchingPriority property value. If higher than 0, this attribute will be used to perform an initial match of the objects between source and target directories. The synchronization engine will try to find the matching object using attribute with lowest value of matching priority first. If not found, the attribute with the next matching priority will be used, and so on a until match is found or no more matching attributes are left. Only attributes that are expected to have unique values, such as email, should be used as matching attributes.
func (m *AttributeMapping) SetMatchingPriority(value *int32)() {
    err := m.GetBackingStore().Set("matchingPriority", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AttributeMapping) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSource sets the source property value. Defines how a value should be extracted (or transformed) from the source object.
func (m *AttributeMapping) SetSource(value AttributeMappingSourceable)() {
    err := m.GetBackingStore().Set("source", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetAttributeName sets the targetAttributeName property value. Name of the attribute on the target object.
func (m *AttributeMapping) SetTargetAttributeName(value *string)() {
    err := m.GetBackingStore().Set("targetAttributeName", value)
    if err != nil {
        panic(err)
    }
}
type AttributeMappingable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDefaultValue()(*string)
    GetExportMissingReferences()(*bool)
    GetFlowBehavior()(*AttributeFlowBehavior)
    GetFlowType()(*AttributeFlowType)
    GetMatchingPriority()(*int32)
    GetOdataType()(*string)
    GetSource()(AttributeMappingSourceable)
    GetTargetAttributeName()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDefaultValue(value *string)()
    SetExportMissingReferences(value *bool)()
    SetFlowBehavior(value *AttributeFlowBehavior)()
    SetFlowType(value *AttributeFlowType)()
    SetMatchingPriority(value *int32)()
    SetOdataType(value *string)()
    SetSource(value AttributeMappingSourceable)()
    SetTargetAttributeName(value *string)()
}
