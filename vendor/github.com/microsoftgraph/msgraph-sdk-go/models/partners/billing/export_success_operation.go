package billing

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ExportSuccessOperation struct {
    Operation
}
// NewExportSuccessOperation instantiates a new ExportSuccessOperation and sets the default values.
func NewExportSuccessOperation()(*ExportSuccessOperation) {
    m := &ExportSuccessOperation{
        Operation: *NewOperation(),
    }
    return m
}
// CreateExportSuccessOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateExportSuccessOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewExportSuccessOperation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ExportSuccessOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Operation.GetFieldDeserializers()
    res["resourceLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateManifestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceLocation(val.(Manifestable))
        }
        return nil
    }
    return res
}
// GetResourceLocation gets the resourceLocation property value. The resourceLocation property
// returns a Manifestable when successful
func (m *ExportSuccessOperation) GetResourceLocation()(Manifestable) {
    val, err := m.GetBackingStore().Get("resourceLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Manifestable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ExportSuccessOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Operation.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("resourceLocation", m.GetResourceLocation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetResourceLocation sets the resourceLocation property value. The resourceLocation property
func (m *ExportSuccessOperation) SetResourceLocation(value Manifestable)() {
    err := m.GetBackingStore().Set("resourceLocation", value)
    if err != nil {
        panic(err)
    }
}
type ExportSuccessOperationable interface {
    Operationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetResourceLocation()(Manifestable)
    SetResourceLocation(value Manifestable)()
}
