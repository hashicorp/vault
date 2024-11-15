package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnPremisesDirectorySynchronization struct {
    Entity
}
// NewOnPremisesDirectorySynchronization instantiates a new OnPremisesDirectorySynchronization and sets the default values.
func NewOnPremisesDirectorySynchronization()(*OnPremisesDirectorySynchronization) {
    m := &OnPremisesDirectorySynchronization{
        Entity: *NewEntity(),
    }
    return m
}
// CreateOnPremisesDirectorySynchronizationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnPremisesDirectorySynchronizationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnPremisesDirectorySynchronization(), nil
}
// GetConfiguration gets the configuration property value. Consists of configurations that can be fine-tuned and impact the on-premises directory synchronization process for a tenant. Nullable.
// returns a OnPremisesDirectorySynchronizationConfigurationable when successful
func (m *OnPremisesDirectorySynchronization) GetConfiguration()(OnPremisesDirectorySynchronizationConfigurationable) {
    val, err := m.GetBackingStore().Get("configuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnPremisesDirectorySynchronizationConfigurationable)
    }
    return nil
}
// GetFeatures gets the features property value. The features property
// returns a OnPremisesDirectorySynchronizationFeatureable when successful
func (m *OnPremisesDirectorySynchronization) GetFeatures()(OnPremisesDirectorySynchronizationFeatureable) {
    val, err := m.GetBackingStore().Get("features")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnPremisesDirectorySynchronizationFeatureable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnPremisesDirectorySynchronization) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["configuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnPremisesDirectorySynchronizationConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfiguration(val.(OnPremisesDirectorySynchronizationConfigurationable))
        }
        return nil
    }
    res["features"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnPremisesDirectorySynchronizationFeatureFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatures(val.(OnPremisesDirectorySynchronizationFeatureable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *OnPremisesDirectorySynchronization) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("configuration", m.GetConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("features", m.GetFeatures())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConfiguration sets the configuration property value. Consists of configurations that can be fine-tuned and impact the on-premises directory synchronization process for a tenant. Nullable.
func (m *OnPremisesDirectorySynchronization) SetConfiguration(value OnPremisesDirectorySynchronizationConfigurationable)() {
    err := m.GetBackingStore().Set("configuration", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatures sets the features property value. The features property
func (m *OnPremisesDirectorySynchronization) SetFeatures(value OnPremisesDirectorySynchronizationFeatureable)() {
    err := m.GetBackingStore().Set("features", value)
    if err != nil {
        panic(err)
    }
}
type OnPremisesDirectorySynchronizationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConfiguration()(OnPremisesDirectorySynchronizationConfigurationable)
    GetFeatures()(OnPremisesDirectorySynchronizationFeatureable)
    SetConfiguration(value OnPremisesDirectorySynchronizationConfigurationable)()
    SetFeatures(value OnPremisesDirectorySynchronizationFeatureable)()
}
