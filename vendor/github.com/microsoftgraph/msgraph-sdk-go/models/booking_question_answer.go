package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type BookingQuestionAnswer struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewBookingQuestionAnswer instantiates a new BookingQuestionAnswer and sets the default values.
func NewBookingQuestionAnswer()(*BookingQuestionAnswer) {
    m := &BookingQuestionAnswer{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateBookingQuestionAnswerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBookingQuestionAnswerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingQuestionAnswer(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *BookingQuestionAnswer) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAnswer gets the answer property value. The answer given by the user in case the answerInputType is text.
// returns a *string when successful
func (m *BookingQuestionAnswer) GetAnswer()(*string) {
    val, err := m.GetBackingStore().Get("answer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAnswerInputType gets the answerInputType property value. The expected answer type. The possible values are: text, radioButton, unknownFutureValue.
// returns a *AnswerInputType when successful
func (m *BookingQuestionAnswer) GetAnswerInputType()(*AnswerInputType) {
    val, err := m.GetBackingStore().Get("answerInputType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AnswerInputType)
    }
    return nil
}
// GetAnswerOptions gets the answerOptions property value. In case the answerInputType is radioButton, this will consists of a list of possible answer values.
// returns a []string when successful
func (m *BookingQuestionAnswer) GetAnswerOptions()([]string) {
    val, err := m.GetBackingStore().Get("answerOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *BookingQuestionAnswer) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *BookingQuestionAnswer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["answer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnswer(val)
        }
        return nil
    }
    res["answerInputType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAnswerInputType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnswerInputType(val.(*AnswerInputType))
        }
        return nil
    }
    res["answerOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAnswerOptions(res)
        }
        return nil
    }
    res["isRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRequired(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["question"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuestion(val)
        }
        return nil
    }
    res["questionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuestionId(val)
        }
        return nil
    }
    res["selectedOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSelectedOptions(res)
        }
        return nil
    }
    return res
}
// GetIsRequired gets the isRequired property value. Indicates whether it is mandatory to answer the custom question.
// returns a *bool when successful
func (m *BookingQuestionAnswer) GetIsRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *BookingQuestionAnswer) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQuestion gets the question property value. The question.
// returns a *string when successful
func (m *BookingQuestionAnswer) GetQuestion()(*string) {
    val, err := m.GetBackingStore().Get("question")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQuestionId gets the questionId property value. The ID of the custom question.
// returns a *string when successful
func (m *BookingQuestionAnswer) GetQuestionId()(*string) {
    val, err := m.GetBackingStore().Get("questionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSelectedOptions gets the selectedOptions property value. The answers selected by the user.
// returns a []string when successful
func (m *BookingQuestionAnswer) GetSelectedOptions()([]string) {
    val, err := m.GetBackingStore().Get("selectedOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BookingQuestionAnswer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("answer", m.GetAnswer())
        if err != nil {
            return err
        }
    }
    if m.GetAnswerInputType() != nil {
        cast := (*m.GetAnswerInputType()).String()
        err := writer.WriteStringValue("answerInputType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetAnswerOptions() != nil {
        err := writer.WriteCollectionOfStringValues("answerOptions", m.GetAnswerOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRequired", m.GetIsRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("question", m.GetQuestion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("questionId", m.GetQuestionId())
        if err != nil {
            return err
        }
    }
    if m.GetSelectedOptions() != nil {
        err := writer.WriteCollectionOfStringValues("selectedOptions", m.GetSelectedOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BookingQuestionAnswer) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAnswer sets the answer property value. The answer given by the user in case the answerInputType is text.
func (m *BookingQuestionAnswer) SetAnswer(value *string)() {
    err := m.GetBackingStore().Set("answer", value)
    if err != nil {
        panic(err)
    }
}
// SetAnswerInputType sets the answerInputType property value. The expected answer type. The possible values are: text, radioButton, unknownFutureValue.
func (m *BookingQuestionAnswer) SetAnswerInputType(value *AnswerInputType)() {
    err := m.GetBackingStore().Set("answerInputType", value)
    if err != nil {
        panic(err)
    }
}
// SetAnswerOptions sets the answerOptions property value. In case the answerInputType is radioButton, this will consists of a list of possible answer values.
func (m *BookingQuestionAnswer) SetAnswerOptions(value []string)() {
    err := m.GetBackingStore().Set("answerOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *BookingQuestionAnswer) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIsRequired sets the isRequired property value. Indicates whether it is mandatory to answer the custom question.
func (m *BookingQuestionAnswer) SetIsRequired(value *bool)() {
    err := m.GetBackingStore().Set("isRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BookingQuestionAnswer) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetQuestion sets the question property value. The question.
func (m *BookingQuestionAnswer) SetQuestion(value *string)() {
    err := m.GetBackingStore().Set("question", value)
    if err != nil {
        panic(err)
    }
}
// SetQuestionId sets the questionId property value. The ID of the custom question.
func (m *BookingQuestionAnswer) SetQuestionId(value *string)() {
    err := m.GetBackingStore().Set("questionId", value)
    if err != nil {
        panic(err)
    }
}
// SetSelectedOptions sets the selectedOptions property value. The answers selected by the user.
func (m *BookingQuestionAnswer) SetSelectedOptions(value []string)() {
    err := m.GetBackingStore().Set("selectedOptions", value)
    if err != nil {
        panic(err)
    }
}
type BookingQuestionAnswerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAnswer()(*string)
    GetAnswerInputType()(*AnswerInputType)
    GetAnswerOptions()([]string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIsRequired()(*bool)
    GetOdataType()(*string)
    GetQuestion()(*string)
    GetQuestionId()(*string)
    GetSelectedOptions()([]string)
    SetAnswer(value *string)()
    SetAnswerInputType(value *AnswerInputType)()
    SetAnswerOptions(value []string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIsRequired(value *bool)()
    SetOdataType(value *string)()
    SetQuestion(value *string)()
    SetQuestionId(value *string)()
    SetSelectedOptions(value []string)()
}
