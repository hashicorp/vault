package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationAssignmentIndividualRecipient struct {
    EducationAssignmentRecipient
}
// NewEducationAssignmentIndividualRecipient instantiates a new EducationAssignmentIndividualRecipient and sets the default values.
func NewEducationAssignmentIndividualRecipient()(*EducationAssignmentIndividualRecipient) {
    m := &EducationAssignmentIndividualRecipient{
        EducationAssignmentRecipient: *NewEducationAssignmentRecipient(),
    }
    odataTypeValue := "#microsoft.graph.educationAssignmentIndividualRecipient"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationAssignmentIndividualRecipientFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationAssignmentIndividualRecipientFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationAssignmentIndividualRecipient(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationAssignmentIndividualRecipient) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationAssignmentRecipient.GetFieldDeserializers()
    res["recipients"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetRecipients(res)
        }
        return nil
    }
    return res
}
// GetRecipients gets the recipients property value. A collection of IDs of the recipients.
// returns a []string when successful
func (m *EducationAssignmentIndividualRecipient) GetRecipients()([]string) {
    val, err := m.GetBackingStore().Get("recipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationAssignmentIndividualRecipient) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationAssignmentRecipient.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetRecipients() != nil {
        err = writer.WriteCollectionOfStringValues("recipients", m.GetRecipients())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRecipients sets the recipients property value. A collection of IDs of the recipients.
func (m *EducationAssignmentIndividualRecipient) SetRecipients(value []string)() {
    err := m.GetBackingStore().Set("recipients", value)
    if err != nil {
        panic(err)
    }
}
type EducationAssignmentIndividualRecipientable interface {
    EducationAssignmentRecipientable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRecipients()([]string)
    SetRecipients(value []string)()
}
