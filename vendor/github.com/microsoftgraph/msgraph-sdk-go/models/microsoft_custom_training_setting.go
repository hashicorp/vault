package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MicrosoftCustomTrainingSetting struct {
    TrainingSetting
}
// NewMicrosoftCustomTrainingSetting instantiates a new MicrosoftCustomTrainingSetting and sets the default values.
func NewMicrosoftCustomTrainingSetting()(*MicrosoftCustomTrainingSetting) {
    m := &MicrosoftCustomTrainingSetting{
        TrainingSetting: *NewTrainingSetting(),
    }
    odataTypeValue := "#microsoft.graph.microsoftCustomTrainingSetting"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMicrosoftCustomTrainingSettingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMicrosoftCustomTrainingSettingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftCustomTrainingSetting(), nil
}
// GetCompletionDateTime gets the completionDateTime property value. The completion date and time of the training. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *MicrosoftCustomTrainingSetting) GetCompletionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completionDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MicrosoftCustomTrainingSetting) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TrainingSetting.GetFieldDeserializers()
    res["completionDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletionDateTime(val)
        }
        return nil
    }
    res["trainingAssignmentMappings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMicrosoftTrainingAssignmentMappingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MicrosoftTrainingAssignmentMappingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MicrosoftTrainingAssignmentMappingable)
                }
            }
            m.SetTrainingAssignmentMappings(res)
        }
        return nil
    }
    res["trainingCompletionDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTrainingCompletionDuration)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrainingCompletionDuration(val.(*TrainingCompletionDuration))
        }
        return nil
    }
    return res
}
// GetTrainingAssignmentMappings gets the trainingAssignmentMappings property value. The mapping details of the associated training.
// returns a []MicrosoftTrainingAssignmentMappingable when successful
func (m *MicrosoftCustomTrainingSetting) GetTrainingAssignmentMappings()([]MicrosoftTrainingAssignmentMappingable) {
    val, err := m.GetBackingStore().Get("trainingAssignmentMappings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MicrosoftTrainingAssignmentMappingable)
    }
    return nil
}
// GetTrainingCompletionDuration gets the trainingCompletionDuration property value. The training completion duration that needs to be provided before scheduling the training. Possible values are: week, fortnite, month, unknownFutureValue.
// returns a *TrainingCompletionDuration when successful
func (m *MicrosoftCustomTrainingSetting) GetTrainingCompletionDuration()(*TrainingCompletionDuration) {
    val, err := m.GetBackingStore().Get("trainingCompletionDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TrainingCompletionDuration)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MicrosoftCustomTrainingSetting) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TrainingSetting.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("completionDateTime", m.GetCompletionDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetTrainingAssignmentMappings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTrainingAssignmentMappings()))
        for i, v := range m.GetTrainingAssignmentMappings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("trainingAssignmentMappings", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTrainingCompletionDuration() != nil {
        cast := (*m.GetTrainingCompletionDuration()).String()
        err = writer.WriteStringValue("trainingCompletionDuration", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompletionDateTime sets the completionDateTime property value. The completion date and time of the training. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *MicrosoftCustomTrainingSetting) SetCompletionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completionDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainingAssignmentMappings sets the trainingAssignmentMappings property value. The mapping details of the associated training.
func (m *MicrosoftCustomTrainingSetting) SetTrainingAssignmentMappings(value []MicrosoftTrainingAssignmentMappingable)() {
    err := m.GetBackingStore().Set("trainingAssignmentMappings", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainingCompletionDuration sets the trainingCompletionDuration property value. The training completion duration that needs to be provided before scheduling the training. Possible values are: week, fortnite, month, unknownFutureValue.
func (m *MicrosoftCustomTrainingSetting) SetTrainingCompletionDuration(value *TrainingCompletionDuration)() {
    err := m.GetBackingStore().Set("trainingCompletionDuration", value)
    if err != nil {
        panic(err)
    }
}
type MicrosoftCustomTrainingSettingable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    TrainingSettingable
    GetCompletionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTrainingAssignmentMappings()([]MicrosoftTrainingAssignmentMappingable)
    GetTrainingCompletionDuration()(*TrainingCompletionDuration)
    SetCompletionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTrainingAssignmentMappings(value []MicrosoftTrainingAssignmentMappingable)()
    SetTrainingCompletionDuration(value *TrainingCompletionDuration)()
}
