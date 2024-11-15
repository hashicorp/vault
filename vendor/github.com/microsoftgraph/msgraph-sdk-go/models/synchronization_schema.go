package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SynchronizationSchema struct {
    Entity
}
// NewSynchronizationSchema instantiates a new SynchronizationSchema and sets the default values.
func NewSynchronizationSchema()(*SynchronizationSchema) {
    m := &SynchronizationSchema{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSynchronizationSchemaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSynchronizationSchemaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSynchronizationSchema(), nil
}
// GetDirectories gets the directories property value. Contains the collection of directories and all of their objects.
// returns a []DirectoryDefinitionable when successful
func (m *SynchronizationSchema) GetDirectories()([]DirectoryDefinitionable) {
    val, err := m.GetBackingStore().Get("directories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryDefinitionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SynchronizationSchema) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["directories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryDefinitionable)
                }
            }
            m.SetDirectories(res)
        }
        return nil
    }
    res["synchronizationRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSynchronizationRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SynchronizationRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SynchronizationRuleable)
                }
            }
            m.SetSynchronizationRules(res)
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetSynchronizationRules gets the synchronizationRules property value. A collection of synchronization rules configured for the synchronizationJob or synchronizationTemplate.
// returns a []SynchronizationRuleable when successful
func (m *SynchronizationSchema) GetSynchronizationRules()([]SynchronizationRuleable) {
    val, err := m.GetBackingStore().Get("synchronizationRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SynchronizationRuleable)
    }
    return nil
}
// GetVersion gets the version property value. The version of the schema, updated automatically with every schema change.
// returns a *string when successful
func (m *SynchronizationSchema) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SynchronizationSchema) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDirectories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDirectories()))
        for i, v := range m.GetDirectories() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("directories", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSynchronizationRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSynchronizationRules()))
        for i, v := range m.GetSynchronizationRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("synchronizationRules", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDirectories sets the directories property value. Contains the collection of directories and all of their objects.
func (m *SynchronizationSchema) SetDirectories(value []DirectoryDefinitionable)() {
    err := m.GetBackingStore().Set("directories", value)
    if err != nil {
        panic(err)
    }
}
// SetSynchronizationRules sets the synchronizationRules property value. A collection of synchronization rules configured for the synchronizationJob or synchronizationTemplate.
func (m *SynchronizationSchema) SetSynchronizationRules(value []SynchronizationRuleable)() {
    err := m.GetBackingStore().Set("synchronizationRules", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The version of the schema, updated automatically with every schema change.
func (m *SynchronizationSchema) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type SynchronizationSchemaable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDirectories()([]DirectoryDefinitionable)
    GetSynchronizationRules()([]SynchronizationRuleable)
    GetVersion()(*string)
    SetDirectories(value []DirectoryDefinitionable)()
    SetSynchronizationRules(value []SynchronizationRuleable)()
    SetVersion(value *string)()
}
