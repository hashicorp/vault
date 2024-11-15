package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventRegistrationPredefinedQuestion struct {
    VirtualEventRegistrationQuestionBase
}
// NewVirtualEventRegistrationPredefinedQuestion instantiates a new VirtualEventRegistrationPredefinedQuestion and sets the default values.
func NewVirtualEventRegistrationPredefinedQuestion()(*VirtualEventRegistrationPredefinedQuestion) {
    m := &VirtualEventRegistrationPredefinedQuestion{
        VirtualEventRegistrationQuestionBase: *NewVirtualEventRegistrationQuestionBase(),
    }
    odataTypeValue := "#microsoft.graph.virtualEventRegistrationPredefinedQuestion"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateVirtualEventRegistrationPredefinedQuestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventRegistrationPredefinedQuestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventRegistrationPredefinedQuestion(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventRegistrationPredefinedQuestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.VirtualEventRegistrationQuestionBase.GetFieldDeserializers()
    res["label"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVirtualEventRegistrationPredefinedQuestionLabel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabel(val.(*VirtualEventRegistrationPredefinedQuestionLabel))
        }
        return nil
    }
    return res
}
// GetLabel gets the label property value. Label of the predefined registration question. It accepts a single line of text: street, city, state, postalCode, countryOrRegion, industry, jobTitle, organization, and unknownFutureValue.
// returns a *VirtualEventRegistrationPredefinedQuestionLabel when successful
func (m *VirtualEventRegistrationPredefinedQuestion) GetLabel()(*VirtualEventRegistrationPredefinedQuestionLabel) {
    val, err := m.GetBackingStore().Get("label")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VirtualEventRegistrationPredefinedQuestionLabel)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventRegistrationPredefinedQuestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.VirtualEventRegistrationQuestionBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetLabel() != nil {
        cast := (*m.GetLabel()).String()
        err = writer.WriteStringValue("label", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLabel sets the label property value. Label of the predefined registration question. It accepts a single line of text: street, city, state, postalCode, countryOrRegion, industry, jobTitle, organization, and unknownFutureValue.
func (m *VirtualEventRegistrationPredefinedQuestion) SetLabel(value *VirtualEventRegistrationPredefinedQuestionLabel)() {
    err := m.GetBackingStore().Set("label", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventRegistrationPredefinedQuestionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    VirtualEventRegistrationQuestionBaseable
    GetLabel()(*VirtualEventRegistrationPredefinedQuestionLabel)
    SetLabel(value *VirtualEventRegistrationPredefinedQuestionLabel)()
}
