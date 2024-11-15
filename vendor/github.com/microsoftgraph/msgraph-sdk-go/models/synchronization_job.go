package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SynchronizationJob struct {
    Entity
}
// NewSynchronizationJob instantiates a new SynchronizationJob and sets the default values.
func NewSynchronizationJob()(*SynchronizationJob) {
    m := &SynchronizationJob{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSynchronizationJobFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSynchronizationJobFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSynchronizationJob(), nil
}
// GetBulkUpload gets the bulkUpload property value. The bulk upload operation for the job.
// returns a BulkUploadable when successful
func (m *SynchronizationJob) GetBulkUpload()(BulkUploadable) {
    val, err := m.GetBackingStore().Get("bulkUpload")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BulkUploadable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SynchronizationJob) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["bulkUpload"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBulkUploadFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBulkUpload(val.(BulkUploadable))
        }
        return nil
    }
    res["schedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchedule(val.(SynchronizationScheduleable))
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
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(SynchronizationStatusable))
        }
        return nil
    }
    res["synchronizationJobSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValuePairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValuePairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValuePairable)
                }
            }
            m.SetSynchronizationJobSettings(res)
        }
        return nil
    }
    res["templateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTemplateId(val)
        }
        return nil
    }
    return res
}
// GetSchedule gets the schedule property value. Schedule used to run the job. Read-only.
// returns a SynchronizationScheduleable when successful
func (m *SynchronizationJob) GetSchedule()(SynchronizationScheduleable) {
    val, err := m.GetBackingStore().Get("schedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationScheduleable)
    }
    return nil
}
// GetSchema gets the schema property value. The synchronization schema configured for the job.
// returns a SynchronizationSchemaable when successful
func (m *SynchronizationJob) GetSchema()(SynchronizationSchemaable) {
    val, err := m.GetBackingStore().Get("schema")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationSchemaable)
    }
    return nil
}
// GetStatus gets the status property value. Status of the job, which includes when the job was last run, current job state, and errors.
// returns a SynchronizationStatusable when successful
func (m *SynchronizationJob) GetStatus()(SynchronizationStatusable) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationStatusable)
    }
    return nil
}
// GetSynchronizationJobSettings gets the synchronizationJobSettings property value. Settings associated with the job. Some settings are inherited from the template.
// returns a []KeyValuePairable when successful
func (m *SynchronizationJob) GetSynchronizationJobSettings()([]KeyValuePairable) {
    val, err := m.GetBackingStore().Get("synchronizationJobSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValuePairable)
    }
    return nil
}
// GetTemplateId gets the templateId property value. Identifier of the synchronization template this job is based on.
// returns a *string when successful
func (m *SynchronizationJob) GetTemplateId()(*string) {
    val, err := m.GetBackingStore().Get("templateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SynchronizationJob) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("bulkUpload", m.GetBulkUpload())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("schedule", m.GetSchedule())
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
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    if m.GetSynchronizationJobSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSynchronizationJobSettings()))
        for i, v := range m.GetSynchronizationJobSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("synchronizationJobSettings", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("templateId", m.GetTemplateId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBulkUpload sets the bulkUpload property value. The bulk upload operation for the job.
func (m *SynchronizationJob) SetBulkUpload(value BulkUploadable)() {
    err := m.GetBackingStore().Set("bulkUpload", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedule sets the schedule property value. Schedule used to run the job. Read-only.
func (m *SynchronizationJob) SetSchedule(value SynchronizationScheduleable)() {
    err := m.GetBackingStore().Set("schedule", value)
    if err != nil {
        panic(err)
    }
}
// SetSchema sets the schema property value. The synchronization schema configured for the job.
func (m *SynchronizationJob) SetSchema(value SynchronizationSchemaable)() {
    err := m.GetBackingStore().Set("schema", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Status of the job, which includes when the job was last run, current job state, and errors.
func (m *SynchronizationJob) SetStatus(value SynchronizationStatusable)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetSynchronizationJobSettings sets the synchronizationJobSettings property value. Settings associated with the job. Some settings are inherited from the template.
func (m *SynchronizationJob) SetSynchronizationJobSettings(value []KeyValuePairable)() {
    err := m.GetBackingStore().Set("synchronizationJobSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetTemplateId sets the templateId property value. Identifier of the synchronization template this job is based on.
func (m *SynchronizationJob) SetTemplateId(value *string)() {
    err := m.GetBackingStore().Set("templateId", value)
    if err != nil {
        panic(err)
    }
}
type SynchronizationJobable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBulkUpload()(BulkUploadable)
    GetSchedule()(SynchronizationScheduleable)
    GetSchema()(SynchronizationSchemaable)
    GetStatus()(SynchronizationStatusable)
    GetSynchronizationJobSettings()([]KeyValuePairable)
    GetTemplateId()(*string)
    SetBulkUpload(value BulkUploadable)()
    SetSchedule(value SynchronizationScheduleable)()
    SetSchema(value SynchronizationSchemaable)()
    SetStatus(value SynchronizationStatusable)()
    SetSynchronizationJobSettings(value []KeyValuePairable)()
    SetTemplateId(value *string)()
}
