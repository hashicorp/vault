package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SynchronizationTemplate struct {
    Entity
}
// NewSynchronizationTemplate instantiates a new SynchronizationTemplate and sets the default values.
func NewSynchronizationTemplate()(*SynchronizationTemplate) {
    m := &SynchronizationTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSynchronizationTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSynchronizationTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSynchronizationTemplate(), nil
}
// GetApplicationId gets the applicationId property value. Identifier of the application this template belongs to.
// returns a *UUID when successful
func (m *SynchronizationTemplate) GetApplicationId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("applicationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetDefaultEscaped gets the default property value. true if this template is recommended to be the default for the application.
// returns a *bool when successful
func (m *SynchronizationTemplate) GetDefaultEscaped()(*bool) {
    val, err := m.GetBackingStore().Get("defaultEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDescription gets the description property value. Description of the template.
// returns a *string when successful
func (m *SynchronizationTemplate) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDiscoverable gets the discoverable property value. true if this template should appear in the collection of templates available for the application instance (service principal).
// returns a *bool when successful
func (m *SynchronizationTemplate) GetDiscoverable()(*bool) {
    val, err := m.GetBackingStore().Get("discoverable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFactoryTag gets the factoryTag property value. One of the well-known factory tags supported by the synchronization engine. The factoryTag tells the synchronization engine which implementation to use when processing jobs based on this template.
// returns a *string when successful
func (m *SynchronizationTemplate) GetFactoryTag()(*string) {
    val, err := m.GetBackingStore().Get("factoryTag")
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
func (m *SynchronizationTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["applicationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationId(val)
        }
        return nil
    }
    res["default"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultEscaped(val)
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
    res["discoverable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiscoverable(val)
        }
        return nil
    }
    res["factoryTag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFactoryTag(val)
        }
        return nil
    }
    res["metadata"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSynchronizationMetadataEntryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SynchronizationMetadataEntryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SynchronizationMetadataEntryable)
                }
            }
            m.SetMetadata(res)
        }
        return nil
    }
    res["schema"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationSchemaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchema(val.(SynchronizationSchemaable))
        }
        return nil
    }
    return res
}
// GetMetadata gets the metadata property value. Additional extension properties. Unless mentioned explicitly, metadata values should not be changed.
// returns a []SynchronizationMetadataEntryable when successful
func (m *SynchronizationTemplate) GetMetadata()([]SynchronizationMetadataEntryable) {
    val, err := m.GetBackingStore().Get("metadata")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SynchronizationMetadataEntryable)
    }
    return nil
}
// GetSchema gets the schema property value. Default synchronization schema for the jobs based on this template.
// returns a SynchronizationSchemaable when successful
func (m *SynchronizationTemplate) GetSchema()(SynchronizationSchemaable) {
    val, err := m.GetBackingStore().Get("schema")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationSchemaable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SynchronizationTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteUUIDValue("applicationId", m.GetApplicationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("default", m.GetDefaultEscaped())
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
        err = writer.WriteBoolValue("discoverable", m.GetDiscoverable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("factoryTag", m.GetFactoryTag())
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
        err = writer.WriteCollectionOfObjectValues("metadata", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("schema", m.GetSchema())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationId sets the applicationId property value. Identifier of the application this template belongs to.
func (m *SynchronizationTemplate) SetApplicationId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("applicationId", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultEscaped sets the default property value. true if this template is recommended to be the default for the application.
func (m *SynchronizationTemplate) SetDefaultEscaped(value *bool)() {
    err := m.GetBackingStore().Set("defaultEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description of the template.
func (m *SynchronizationTemplate) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDiscoverable sets the discoverable property value. true if this template should appear in the collection of templates available for the application instance (service principal).
func (m *SynchronizationTemplate) SetDiscoverable(value *bool)() {
    err := m.GetBackingStore().Set("discoverable", value)
    if err != nil {
        panic(err)
    }
}
// SetFactoryTag sets the factoryTag property value. One of the well-known factory tags supported by the synchronization engine. The factoryTag tells the synchronization engine which implementation to use when processing jobs based on this template.
func (m *SynchronizationTemplate) SetFactoryTag(value *string)() {
    err := m.GetBackingStore().Set("factoryTag", value)
    if err != nil {
        panic(err)
    }
}
// SetMetadata sets the metadata property value. Additional extension properties. Unless mentioned explicitly, metadata values should not be changed.
func (m *SynchronizationTemplate) SetMetadata(value []SynchronizationMetadataEntryable)() {
    err := m.GetBackingStore().Set("metadata", value)
    if err != nil {
        panic(err)
    }
}
// SetSchema sets the schema property value. Default synchronization schema for the jobs based on this template.
func (m *SynchronizationTemplate) SetSchema(value SynchronizationSchemaable)() {
    err := m.GetBackingStore().Set("schema", value)
    if err != nil {
        panic(err)
    }
}
type SynchronizationTemplateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetDefaultEscaped()(*bool)
    GetDescription()(*string)
    GetDiscoverable()(*bool)
    GetFactoryTag()(*string)
    GetMetadata()([]SynchronizationMetadataEntryable)
    GetSchema()(SynchronizationSchemaable)
    SetApplicationId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetDefaultEscaped(value *bool)()
    SetDescription(value *string)()
    SetDiscoverable(value *bool)()
    SetFactoryTag(value *string)()
    SetMetadata(value []SynchronizationMetadataEntryable)()
    SetSchema(value SynchronizationSchemaable)()
}
