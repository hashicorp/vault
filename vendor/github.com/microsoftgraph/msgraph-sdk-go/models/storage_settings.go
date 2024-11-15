package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type StorageSettings struct {
    Entity
}
// NewStorageSettings instantiates a new StorageSettings and sets the default values.
func NewStorageSettings()(*StorageSettings) {
    m := &StorageSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateStorageSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateStorageSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewStorageSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *StorageSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["quota"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUnifiedStorageQuotaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuota(val.(UnifiedStorageQuotaable))
        }
        return nil
    }
    return res
}
// GetQuota gets the quota property value. The quota property
// returns a UnifiedStorageQuotaable when successful
func (m *StorageSettings) GetQuota()(UnifiedStorageQuotaable) {
    val, err := m.GetBackingStore().Get("quota")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UnifiedStorageQuotaable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *StorageSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("quota", m.GetQuota())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetQuota sets the quota property value. The quota property
func (m *StorageSettings) SetQuota(value UnifiedStorageQuotaable)() {
    err := m.GetBackingStore().Set("quota", value)
    if err != nil {
        panic(err)
    }
}
type StorageSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetQuota()(UnifiedStorageQuotaable)
    SetQuota(value UnifiedStorageQuotaable)()
}
