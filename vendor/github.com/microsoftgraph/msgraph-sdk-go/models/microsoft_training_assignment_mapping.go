package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MicrosoftTrainingAssignmentMapping struct {
    TrainingSetting
}
// NewMicrosoftTrainingAssignmentMapping instantiates a new MicrosoftTrainingAssignmentMapping and sets the default values.
func NewMicrosoftTrainingAssignmentMapping()(*MicrosoftTrainingAssignmentMapping) {
    m := &MicrosoftTrainingAssignmentMapping{
        TrainingSetting: *NewTrainingSetting(),
    }
    odataTypeValue := "#microsoft.graph.microsoftTrainingAssignmentMapping"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMicrosoftTrainingAssignmentMappingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMicrosoftTrainingAssignmentMappingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftTrainingAssignmentMapping(), nil
}
// GetAssignedTo gets the assignedTo property value. A user collection that specifies to whom the training should be assigned. Possible values are: none, allUsers, clickedPayload, compromised, reportedPhish, readButNotClicked, didNothing, unknownFutureValue.
// returns a []TrainingAssignedTo when successful
func (m *MicrosoftTrainingAssignmentMapping) GetAssignedTo()([]TrainingAssignedTo) {
    val, err := m.GetBackingStore().Get("assignedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TrainingAssignedTo)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MicrosoftTrainingAssignmentMapping) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TrainingSetting.GetFieldDeserializers()
    res["assignedTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseTrainingAssignedTo)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TrainingAssignedTo, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*TrainingAssignedTo))
                }
            }
            m.SetAssignedTo(res)
        }
        return nil
    }
    res["training"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTrainingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTraining(val.(Trainingable))
        }
        return nil
    }
    return res
}
// GetTraining gets the training property value. The training property
// returns a Trainingable when successful
func (m *MicrosoftTrainingAssignmentMapping) GetTraining()(Trainingable) {
    val, err := m.GetBackingStore().Get("training")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Trainingable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MicrosoftTrainingAssignmentMapping) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TrainingSetting.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignedTo() != nil {
        err = writer.WriteCollectionOfStringValues("assignedTo", SerializeTrainingAssignedTo(m.GetAssignedTo()))
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("training", m.GetTraining())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignedTo sets the assignedTo property value. A user collection that specifies to whom the training should be assigned. Possible values are: none, allUsers, clickedPayload, compromised, reportedPhish, readButNotClicked, didNothing, unknownFutureValue.
func (m *MicrosoftTrainingAssignmentMapping) SetAssignedTo(value []TrainingAssignedTo)() {
    err := m.GetBackingStore().Set("assignedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetTraining sets the training property value. The training property
func (m *MicrosoftTrainingAssignmentMapping) SetTraining(value Trainingable)() {
    err := m.GetBackingStore().Set("training", value)
    if err != nil {
        panic(err)
    }
}
type MicrosoftTrainingAssignmentMappingable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    TrainingSettingable
    GetAssignedTo()([]TrainingAssignedTo)
    GetTraining()(Trainingable)
    SetAssignedTo(value []TrainingAssignedTo)()
    SetTraining(value Trainingable)()
}
