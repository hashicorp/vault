package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageTextInputQuestion struct {
    AccessPackageQuestion
}
// NewAccessPackageTextInputQuestion instantiates a new AccessPackageTextInputQuestion and sets the default values.
func NewAccessPackageTextInputQuestion()(*AccessPackageTextInputQuestion) {
    m := &AccessPackageTextInputQuestion{
        AccessPackageQuestion: *NewAccessPackageQuestion(),
    }
    odataTypeValue := "#microsoft.graph.accessPackageTextInputQuestion"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessPackageTextInputQuestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageTextInputQuestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageTextInputQuestion(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageTextInputQuestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessPackageQuestion.GetFieldDeserializers()
    res["isSingleLineQuestion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSingleLineQuestion(val)
        }
        return nil
    }
    res["regexPattern"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegexPattern(val)
        }
        return nil
    }
    return res
}
// GetIsSingleLineQuestion gets the isSingleLineQuestion property value. Indicates whether the answer is in single or multiple line format.
// returns a *bool when successful
func (m *AccessPackageTextInputQuestion) GetIsSingleLineQuestion()(*bool) {
    val, err := m.GetBackingStore().Get("isSingleLineQuestion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRegexPattern gets the regexPattern property value. The regular expression pattern that any answer to this question must match.
// returns a *string when successful
func (m *AccessPackageTextInputQuestion) GetRegexPattern()(*string) {
    val, err := m.GetBackingStore().Get("regexPattern")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageTextInputQuestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessPackageQuestion.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isSingleLineQuestion", m.GetIsSingleLineQuestion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("regexPattern", m.GetRegexPattern())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsSingleLineQuestion sets the isSingleLineQuestion property value. Indicates whether the answer is in single or multiple line format.
func (m *AccessPackageTextInputQuestion) SetIsSingleLineQuestion(value *bool)() {
    err := m.GetBackingStore().Set("isSingleLineQuestion", value)
    if err != nil {
        panic(err)
    }
}
// SetRegexPattern sets the regexPattern property value. The regular expression pattern that any answer to this question must match.
func (m *AccessPackageTextInputQuestion) SetRegexPattern(value *string)() {
    err := m.GetBackingStore().Set("regexPattern", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageTextInputQuestionable interface {
    AccessPackageQuestionable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsSingleLineQuestion()(*bool)
    GetRegexPattern()(*string)
    SetIsSingleLineQuestion(value *bool)()
    SetRegexPattern(value *string)()
}
