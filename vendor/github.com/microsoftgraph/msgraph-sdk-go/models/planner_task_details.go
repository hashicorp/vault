package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PlannerTaskDetails struct {
    Entity
}
// NewPlannerTaskDetails instantiates a new PlannerTaskDetails and sets the default values.
func NewPlannerTaskDetails()(*PlannerTaskDetails) {
    m := &PlannerTaskDetails{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerTaskDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePlannerTaskDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerTaskDetails(), nil
}
// GetChecklist gets the checklist property value. The collection of checklist items on the task.
// returns a PlannerChecklistItemsable when successful
func (m *PlannerTaskDetails) GetChecklist()(PlannerChecklistItemsable) {
    val, err := m.GetBackingStore().Get("checklist")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerChecklistItemsable)
    }
    return nil
}
// GetDescription gets the description property value. Description of the task.
// returns a *string when successful
func (m *PlannerTaskDetails) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PlannerTaskDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["checklist"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerChecklistItemsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChecklist(val.(PlannerChecklistItemsable))
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["previewType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePlannerPreviewType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreviewType(val.(*PlannerPreviewType))
        }
        return nil
    }
    res["references"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerExternalReferencesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReferences(val.(PlannerExternalReferencesable))
        }
        return nil
    }
    return res
}
// GetPreviewType gets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference. When set to automatic the displayed preview is chosen by the app viewing the task.
// returns a *PlannerPreviewType when successful
func (m *PlannerTaskDetails) GetPreviewType()(*PlannerPreviewType) {
    val, err := m.GetBackingStore().Get("previewType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PlannerPreviewType)
    }
    return nil
}
// GetReferences gets the references property value. The collection of references on the task.
// returns a PlannerExternalReferencesable when successful
func (m *PlannerTaskDetails) GetReferences()(PlannerExternalReferencesable) {
    val, err := m.GetBackingStore().Get("references")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerExternalReferencesable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PlannerTaskDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("checklist", m.GetChecklist())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    if m.GetPreviewType() != nil {
        cast := (*m.GetPreviewType()).String()
        err = writer.WriteStringValue("previewType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("references", m.GetReferences())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChecklist sets the checklist property value. The collection of checklist items on the task.
func (m *PlannerTaskDetails) SetChecklist(value PlannerChecklistItemsable)() {
    err := m.GetBackingStore().Set("checklist", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description of the task.
func (m *PlannerTaskDetails) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviewType sets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference. When set to automatic the displayed preview is chosen by the app viewing the task.
func (m *PlannerTaskDetails) SetPreviewType(value *PlannerPreviewType)() {
    err := m.GetBackingStore().Set("previewType", value)
    if err != nil {
        panic(err)
    }
}
// SetReferences sets the references property value. The collection of references on the task.
func (m *PlannerTaskDetails) SetReferences(value PlannerExternalReferencesable)() {
    err := m.GetBackingStore().Set("references", value)
    if err != nil {
        panic(err)
    }
}
type PlannerTaskDetailsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChecklist()(PlannerChecklistItemsable)
    GetDescription()(*string)
    GetPreviewType()(*PlannerPreviewType)
    GetReferences()(PlannerExternalReferencesable)
    SetChecklist(value PlannerChecklistItemsable)()
    SetDescription(value *string)()
    SetPreviewType(value *PlannerPreviewType)()
    SetReferences(value PlannerExternalReferencesable)()
}
