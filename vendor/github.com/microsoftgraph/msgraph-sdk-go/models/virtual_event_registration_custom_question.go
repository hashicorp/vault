package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventRegistrationCustomQuestion struct {
    VirtualEventRegistrationQuestionBase
}
// NewVirtualEventRegistrationCustomQuestion instantiates a new VirtualEventRegistrationCustomQuestion and sets the default values.
func NewVirtualEventRegistrationCustomQuestion()(*VirtualEventRegistrationCustomQuestion) {
    m := &VirtualEventRegistrationCustomQuestion{
        VirtualEventRegistrationQuestionBase: *NewVirtualEventRegistrationQuestionBase(),
    }
    odataTypeValue := "#microsoft.graph.virtualEventRegistrationCustomQuestion"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateVirtualEventRegistrationCustomQuestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventRegistrationCustomQuestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventRegistrationCustomQuestion(), nil
}
// GetAnswerChoices gets the answerChoices property value. Answer choices when answerInputType is singleChoice or multiChoice.
// returns a []string when successful
func (m *VirtualEventRegistrationCustomQuestion) GetAnswerChoices()([]string) {
    val, err := m.GetBackingStore().Get("answerChoices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAnswerInputType gets the answerInputType property value. Input type of the registration question answer. Possible values are text, multilineText, singleChoice, multiChoice, boolean, and unknownFutureValue.
// returns a *VirtualEventRegistrationQuestionAnswerInputType when successful
func (m *VirtualEventRegistrationCustomQuestion) GetAnswerInputType()(*VirtualEventRegistrationQuestionAnswerInputType) {
    val, err := m.GetBackingStore().Get("answerInputType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VirtualEventRegistrationQuestionAnswerInputType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventRegistrationCustomQuestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.VirtualEventRegistrationQuestionBase.GetFieldDeserializers()
    res["answerChoices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAnswerChoices(res)
        }
        return nil
    }
    res["answerInputType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVirtualEventRegistrationQuestionAnswerInputType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnswerInputType(val.(*VirtualEventRegistrationQuestionAnswerInputType))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *VirtualEventRegistrationCustomQuestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.VirtualEventRegistrationQuestionBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAnswerChoices() != nil {
        err = writer.WriteCollectionOfStringValues("answerChoices", m.GetAnswerChoices())
        if err != nil {
            return err
        }
    }
    if m.GetAnswerInputType() != nil {
        cast := (*m.GetAnswerInputType()).String()
        err = writer.WriteStringValue("answerInputType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnswerChoices sets the answerChoices property value. Answer choices when answerInputType is singleChoice or multiChoice.
func (m *VirtualEventRegistrationCustomQuestion) SetAnswerChoices(value []string)() {
    err := m.GetBackingStore().Set("answerChoices", value)
    if err != nil {
        panic(err)
    }
}
// SetAnswerInputType sets the answerInputType property value. Input type of the registration question answer. Possible values are text, multilineText, singleChoice, multiChoice, boolean, and unknownFutureValue.
func (m *VirtualEventRegistrationCustomQuestion) SetAnswerInputType(value *VirtualEventRegistrationQuestionAnswerInputType)() {
    err := m.GetBackingStore().Set("answerInputType", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventRegistrationCustomQuestionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    VirtualEventRegistrationQuestionBaseable
    GetAnswerChoices()([]string)
    GetAnswerInputType()(*VirtualEventRegistrationQuestionAnswerInputType)
    SetAnswerChoices(value []string)()
    SetAnswerInputType(value *VirtualEventRegistrationQuestionAnswerInputType)()
}
