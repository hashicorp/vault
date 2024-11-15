package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SynchronizationRule struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSynchronizationRule instantiates a new SynchronizationRule and sets the default values.
func NewSynchronizationRule()(*SynchronizationRule) {
    m := &SynchronizationRule{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSynchronizationRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSynchronizationRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSynchronizationRule(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SynchronizationRule) GetAdditionalData()(map[string]any) {
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
func (m *SynchronizationRule) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetContainerFilter gets the containerFilter property value. The containerFilter property
// returns a ContainerFilterable when successful
func (m *SynchronizationRule) GetContainerFilter()(ContainerFilterable) {
    val, err := m.GetBackingStore().Get("containerFilter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContainerFilterable)
    }
    return nil
}
// GetEditable gets the editable property value. true if the synchronization rule can be customized; false if this rule is read-only and shouldn't be changed.
// returns a *bool when successful
func (m *SynchronizationRule) GetEditable()(*bool) {
    val, err := m.GetBackingStore().Get("editable")
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
func (m *SynchronizationRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["containerFilter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContainerFilterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContainerFilter(val.(ContainerFilterable))
        }
        return nil
    }
    res["editable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEditable(val)
        }
        return nil
    }
    res["groupFilter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGroupFilterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupFilter(val.(GroupFilterable))
        }
        return nil
    }
    res["id"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetId(val)
        }
        return nil
    }
    res["metadata"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetMetadata(res)
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
    res["objectMappings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateObjectMappingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ObjectMappingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ObjectMappingable)
                }
            }
            m.SetObjectMappings(res)
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
    res["priority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPriority(val)
        }
        return nil
    }
    res["sourceDirectoryName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceDirectoryName(val)
        }
        return nil
    }
    res["targetDirectoryName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetDirectoryName(val)
        }
        return nil
    }
    return res
}
// GetGroupFilter gets the groupFilter property value. The groupFilter property
// returns a GroupFilterable when successful
func (m *SynchronizationRule) GetGroupFilter()(GroupFilterable) {
    val, err := m.GetBackingStore().Get("groupFilter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(GroupFilterable)
    }
    return nil
}
// GetId gets the id property value. Synchronization rule identifier. Must be one of the identifiers recognized by the synchronization engine. Supported rule identifiers can be found in the synchronization template returned by the API.
// returns a *string when successful
func (m *SynchronizationRule) GetId()(*string) {
    val, err := m.GetBackingStore().Get("id")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMetadata gets the metadata property value. Additional extension properties. Unless instructed explicitly by the support team, metadata values shouldn't be changed.
// returns a []StringKeyStringValuePairable when successful
func (m *SynchronizationRule) GetMetadata()([]StringKeyStringValuePairable) {
    val, err := m.GetBackingStore().Get("metadata")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]StringKeyStringValuePairable)
    }
    return nil
}
// GetName gets the name property value. Human-readable name of the synchronization rule. Not nullable.
// returns a *string when successful
func (m *SynchronizationRule) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetObjectMappings gets the objectMappings property value. Collection of object mappings supported by the rule. Tells the synchronization engine which objects should be synchronized.
// returns a []ObjectMappingable when successful
func (m *SynchronizationRule) GetObjectMappings()([]ObjectMappingable) {
    val, err := m.GetBackingStore().Get("objectMappings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ObjectMappingable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SynchronizationRule) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPriority gets the priority property value. Priority relative to other rules in the synchronizationSchema. Rules with the lowest priority number will be processed first.
// returns a *int32 when successful
func (m *SynchronizationRule) GetPriority()(*int32) {
    val, err := m.GetBackingStore().Get("priority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSourceDirectoryName gets the sourceDirectoryName property value. Name of the source directory. Must match one of the directory definitions in synchronizationSchema.
// returns a *string when successful
func (m *SynchronizationRule) GetSourceDirectoryName()(*string) {
    val, err := m.GetBackingStore().Get("sourceDirectoryName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTargetDirectoryName gets the targetDirectoryName property value. Name of the target directory. Must match one of the directory definitions in synchronizationSchema.
// returns a *string when successful
func (m *SynchronizationRule) GetTargetDirectoryName()(*string) {
    val, err := m.GetBackingStore().Get("targetDirectoryName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SynchronizationRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("containerFilter", m.GetContainerFilter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("editable", m.GetEditable())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("groupFilter", m.GetGroupFilter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("id", m.GetId())
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
        err := writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetObjectMappings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetObjectMappings()))
        for i, v := range m.GetObjectMappings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("objectMappings", cast)
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
        err := writer.WriteInt32Value("priority", m.GetPriority())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("sourceDirectoryName", m.GetSourceDirectoryName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("targetDirectoryName", m.GetTargetDirectoryName())
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
func (m *SynchronizationRule) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SynchronizationRule) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetContainerFilter sets the containerFilter property value. The containerFilter property
func (m *SynchronizationRule) SetContainerFilter(value ContainerFilterable)() {
    err := m.GetBackingStore().Set("containerFilter", value)
    if err != nil {
        panic(err)
    }
}
// SetEditable sets the editable property value. true if the synchronization rule can be customized; false if this rule is read-only and shouldn't be changed.
func (m *SynchronizationRule) SetEditable(value *bool)() {
    err := m.GetBackingStore().Set("editable", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupFilter sets the groupFilter property value. The groupFilter property
func (m *SynchronizationRule) SetGroupFilter(value GroupFilterable)() {
    err := m.GetBackingStore().Set("groupFilter", value)
    if err != nil {
        panic(err)
    }
}
// SetId sets the id property value. Synchronization rule identifier. Must be one of the identifiers recognized by the synchronization engine. Supported rule identifiers can be found in the synchronization template returned by the API.
func (m *SynchronizationRule) SetId(value *string)() {
    err := m.GetBackingStore().Set("id", value)
    if err != nil {
        panic(err)
    }
}
// SetMetadata sets the metadata property value. Additional extension properties. Unless instructed explicitly by the support team, metadata values shouldn't be changed.
func (m *SynchronizationRule) SetMetadata(value []StringKeyStringValuePairable)() {
    err := m.GetBackingStore().Set("metadata", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Human-readable name of the synchronization rule. Not nullable.
func (m *SynchronizationRule) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetObjectMappings sets the objectMappings property value. Collection of object mappings supported by the rule. Tells the synchronization engine which objects should be synchronized.
func (m *SynchronizationRule) SetObjectMappings(value []ObjectMappingable)() {
    err := m.GetBackingStore().Set("objectMappings", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SynchronizationRule) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPriority sets the priority property value. Priority relative to other rules in the synchronizationSchema. Rules with the lowest priority number will be processed first.
func (m *SynchronizationRule) SetPriority(value *int32)() {
    err := m.GetBackingStore().Set("priority", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceDirectoryName sets the sourceDirectoryName property value. Name of the source directory. Must match one of the directory definitions in synchronizationSchema.
func (m *SynchronizationRule) SetSourceDirectoryName(value *string)() {
    err := m.GetBackingStore().Set("sourceDirectoryName", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetDirectoryName sets the targetDirectoryName property value. Name of the target directory. Must match one of the directory definitions in synchronizationSchema.
func (m *SynchronizationRule) SetTargetDirectoryName(value *string)() {
    err := m.GetBackingStore().Set("targetDirectoryName", value)
    if err != nil {
        panic(err)
    }
}
type SynchronizationRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetContainerFilter()(ContainerFilterable)
    GetEditable()(*bool)
    GetGroupFilter()(GroupFilterable)
    GetId()(*string)
    GetMetadata()([]StringKeyStringValuePairable)
    GetName()(*string)
    GetObjectMappings()([]ObjectMappingable)
    GetOdataType()(*string)
    GetPriority()(*int32)
    GetSourceDirectoryName()(*string)
    GetTargetDirectoryName()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetContainerFilter(value ContainerFilterable)()
    SetEditable(value *bool)()
    SetGroupFilter(value GroupFilterable)()
    SetId(value *string)()
    SetMetadata(value []StringKeyStringValuePairable)()
    SetName(value *string)()
    SetObjectMappings(value []ObjectMappingable)()
    SetOdataType(value *string)()
    SetPriority(value *int32)()
    SetSourceDirectoryName(value *string)()
    SetTargetDirectoryName(value *string)()
}
