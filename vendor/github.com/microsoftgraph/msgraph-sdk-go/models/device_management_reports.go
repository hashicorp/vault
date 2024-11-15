package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementReports singleton entity that acts as a container for all reports functionality.
type DeviceManagementReports struct {
    Entity
}
// NewDeviceManagementReports instantiates a new DeviceManagementReports and sets the default values.
func NewDeviceManagementReports()(*DeviceManagementReports) {
    m := &DeviceManagementReports{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementReportsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceManagementReportsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagementReports(), nil
}
// GetExportJobs gets the exportJobs property value. Entity representing a job to export a report
// returns a []DeviceManagementExportJobable when successful
func (m *DeviceManagementReports) GetExportJobs()([]DeviceManagementExportJobable) {
    val, err := m.GetBackingStore().Get("exportJobs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceManagementExportJobable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceManagementReports) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["exportJobs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceManagementExportJobFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceManagementExportJobable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceManagementExportJobable)
                }
            }
            m.SetExportJobs(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *DeviceManagementReports) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetExportJobs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExportJobs()))
        for i, v := range m.GetExportJobs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("exportJobs", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExportJobs sets the exportJobs property value. Entity representing a job to export a report
func (m *DeviceManagementReports) SetExportJobs(value []DeviceManagementExportJobable)() {
    err := m.GetBackingStore().Set("exportJobs", value)
    if err != nil {
        panic(err)
    }
}
type DeviceManagementReportsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetExportJobs()([]DeviceManagementExportJobable)
    SetExportJobs(value []DeviceManagementExportJobable)()
}
