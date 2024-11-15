package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RequestorManager struct {
    SubjectSet
}
// NewRequestorManager instantiates a new RequestorManager and sets the default values.
func NewRequestorManager()(*RequestorManager) {
    m := &RequestorManager{
        SubjectSet: *NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.requestorManager"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRequestorManagerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRequestorManagerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRequestorManager(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RequestorManager) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectSet.GetFieldDeserializers()
    res["managerLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagerLevel(val)
        }
        return nil
    }
    return res
}
// GetManagerLevel gets the managerLevel property value. The hierarchical level of the manager with respect to the requestor. For example, the direct manager of a requestor would have a managerLevel of 1, while the manager of the requestor's manager would have a managerLevel of 2. Default value for managerLevel is 1. Possible values for this property range from 1 to 2.
// returns a *int32 when successful
func (m *RequestorManager) GetManagerLevel()(*int32) {
    val, err := m.GetBackingStore().Get("managerLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RequestorManager) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectSet.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("managerLevel", m.GetManagerLevel())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetManagerLevel sets the managerLevel property value. The hierarchical level of the manager with respect to the requestor. For example, the direct manager of a requestor would have a managerLevel of 1, while the manager of the requestor's manager would have a managerLevel of 2. Default value for managerLevel is 1. Possible values for this property range from 1 to 2.
func (m *RequestorManager) SetManagerLevel(value *int32)() {
    err := m.GetBackingStore().Set("managerLevel", value)
    if err != nil {
        panic(err)
    }
}
type RequestorManagerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SubjectSetable
    GetManagerLevel()(*int32)
    SetManagerLevel(value *int32)()
}
