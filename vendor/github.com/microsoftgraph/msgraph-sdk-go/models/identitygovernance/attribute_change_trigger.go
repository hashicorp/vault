package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AttributeChangeTrigger struct {
    WorkflowExecutionTrigger
}
// NewAttributeChangeTrigger instantiates a new AttributeChangeTrigger and sets the default values.
func NewAttributeChangeTrigger()(*AttributeChangeTrigger) {
    m := &AttributeChangeTrigger{
        WorkflowExecutionTrigger: *NewWorkflowExecutionTrigger(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.attributeChangeTrigger"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAttributeChangeTriggerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttributeChangeTriggerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttributeChangeTrigger(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AttributeChangeTrigger) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WorkflowExecutionTrigger.GetFieldDeserializers()
    res["triggerAttributes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTriggerAttributeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TriggerAttributeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TriggerAttributeable)
                }
            }
            m.SetTriggerAttributes(res)
        }
        return nil
    }
    return res
}
// GetTriggerAttributes gets the triggerAttributes property value. The trigger attribute being changed that triggers the workflowexecutiontrigger of a workflow.)
// returns a []TriggerAttributeable when successful
func (m *AttributeChangeTrigger) GetTriggerAttributes()([]TriggerAttributeable) {
    val, err := m.GetBackingStore().Get("triggerAttributes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TriggerAttributeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AttributeChangeTrigger) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WorkflowExecutionTrigger.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetTriggerAttributes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTriggerAttributes()))
        for i, v := range m.GetTriggerAttributes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("triggerAttributes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTriggerAttributes sets the triggerAttributes property value. The trigger attribute being changed that triggers the workflowexecutiontrigger of a workflow.)
func (m *AttributeChangeTrigger) SetTriggerAttributes(value []TriggerAttributeable)() {
    err := m.GetBackingStore().Set("triggerAttributes", value)
    if err != nil {
        panic(err)
    }
}
type AttributeChangeTriggerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WorkflowExecutionTriggerable
    GetTriggerAttributes()([]TriggerAttributeable)
    SetTriggerAttributes(value []TriggerAttributeable)()
}
