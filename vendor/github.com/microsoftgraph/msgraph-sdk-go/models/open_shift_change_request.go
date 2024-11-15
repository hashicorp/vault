package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OpenShiftChangeRequest struct {
    ScheduleChangeRequest
}
// NewOpenShiftChangeRequest instantiates a new OpenShiftChangeRequest and sets the default values.
func NewOpenShiftChangeRequest()(*OpenShiftChangeRequest) {
    m := &OpenShiftChangeRequest{
        ScheduleChangeRequest: *NewScheduleChangeRequest(),
    }
    odataTypeValue := "#microsoft.graph.openShiftChangeRequest"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOpenShiftChangeRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOpenShiftChangeRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOpenShiftChangeRequest(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OpenShiftChangeRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ScheduleChangeRequest.GetFieldDeserializers()
    res["openShiftId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOpenShiftId(val)
        }
        return nil
    }
    return res
}
// GetOpenShiftId gets the openShiftId property value. ID for the open shift.
// returns a *string when successful
func (m *OpenShiftChangeRequest) GetOpenShiftId()(*string) {
    val, err := m.GetBackingStore().Get("openShiftId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OpenShiftChangeRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ScheduleChangeRequest.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("openShiftId", m.GetOpenShiftId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOpenShiftId sets the openShiftId property value. ID for the open shift.
func (m *OpenShiftChangeRequest) SetOpenShiftId(value *string)() {
    err := m.GetBackingStore().Set("openShiftId", value)
    if err != nil {
        panic(err)
    }
}
type OpenShiftChangeRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ScheduleChangeRequestable
    GetOpenShiftId()(*string)
    SetOpenShiftId(value *string)()
}
