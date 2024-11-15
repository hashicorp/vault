package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TargetManager struct {
    SubjectSet
}
// NewTargetManager instantiates a new TargetManager and sets the default values.
func NewTargetManager()(*TargetManager) {
    m := &TargetManager{
        SubjectSet: *NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.targetManager"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTargetManagerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTargetManagerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTargetManager(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TargetManager) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
// GetManagerLevel gets the managerLevel property value. Manager level, between 1 and 4. The direct manager is 1.
// returns a *int32 when successful
func (m *TargetManager) GetManagerLevel()(*int32) {
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
func (m *TargetManager) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetManagerLevel sets the managerLevel property value. Manager level, between 1 and 4. The direct manager is 1.
func (m *TargetManager) SetManagerLevel(value *int32)() {
    err := m.GetBackingStore().Set("managerLevel", value)
    if err != nil {
        panic(err)
    }
}
type TargetManagerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SubjectSetable
    GetManagerLevel()(*int32)
    SetManagerLevel(value *int32)()
}
