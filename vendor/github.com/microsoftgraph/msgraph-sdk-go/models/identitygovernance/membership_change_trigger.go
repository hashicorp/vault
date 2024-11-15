package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MembershipChangeTrigger struct {
    WorkflowExecutionTrigger
}
// NewMembershipChangeTrigger instantiates a new MembershipChangeTrigger and sets the default values.
func NewMembershipChangeTrigger()(*MembershipChangeTrigger) {
    m := &MembershipChangeTrigger{
        WorkflowExecutionTrigger: *NewWorkflowExecutionTrigger(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.membershipChangeTrigger"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMembershipChangeTriggerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMembershipChangeTriggerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMembershipChangeTrigger(), nil
}
// GetChangeType gets the changeType property value. The changeType property
// returns a *MembershipChangeType when successful
func (m *MembershipChangeTrigger) GetChangeType()(*MembershipChangeType) {
    val, err := m.GetBackingStore().Get("changeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MembershipChangeType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MembershipChangeTrigger) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WorkflowExecutionTrigger.GetFieldDeserializers()
    res["changeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMembershipChangeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChangeType(val.(*MembershipChangeType))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *MembershipChangeTrigger) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WorkflowExecutionTrigger.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChangeType() != nil {
        cast := (*m.GetChangeType()).String()
        err = writer.WriteStringValue("changeType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChangeType sets the changeType property value. The changeType property
func (m *MembershipChangeTrigger) SetChangeType(value *MembershipChangeType)() {
    err := m.GetBackingStore().Set("changeType", value)
    if err != nil {
        panic(err)
    }
}
type MembershipChangeTriggerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WorkflowExecutionTriggerable
    GetChangeType()(*MembershipChangeType)
    SetChangeType(value *MembershipChangeType)()
}
