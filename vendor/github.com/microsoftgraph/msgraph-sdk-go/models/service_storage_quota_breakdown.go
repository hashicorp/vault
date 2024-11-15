package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServiceStorageQuotaBreakdown struct {
    StorageQuotaBreakdown
}
// NewServiceStorageQuotaBreakdown instantiates a new ServiceStorageQuotaBreakdown and sets the default values.
func NewServiceStorageQuotaBreakdown()(*ServiceStorageQuotaBreakdown) {
    m := &ServiceStorageQuotaBreakdown{
        StorageQuotaBreakdown: *NewStorageQuotaBreakdown(),
    }
    return m
}
// CreateServiceStorageQuotaBreakdownFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceStorageQuotaBreakdownFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceStorageQuotaBreakdown(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServiceStorageQuotaBreakdown) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.StorageQuotaBreakdown.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *ServiceStorageQuotaBreakdown) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.StorageQuotaBreakdown.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type ServiceStorageQuotaBreakdownable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    StorageQuotaBreakdownable
}
