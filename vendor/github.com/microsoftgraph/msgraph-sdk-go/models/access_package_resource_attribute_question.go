package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageResourceAttributeQuestion struct {
    AccessPackageResourceAttributeSource
}
// NewAccessPackageResourceAttributeQuestion instantiates a new AccessPackageResourceAttributeQuestion and sets the default values.
func NewAccessPackageResourceAttributeQuestion()(*AccessPackageResourceAttributeQuestion) {
    m := &AccessPackageResourceAttributeQuestion{
        AccessPackageResourceAttributeSource: *NewAccessPackageResourceAttributeSource(),
    }
    odataTypeValue := "#microsoft.graph.accessPackageResourceAttributeQuestion"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessPackageResourceAttributeQuestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageResourceAttributeQuestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageResourceAttributeQuestion(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageResourceAttributeQuestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessPackageResourceAttributeSource.GetFieldDeserializers()
    res["question"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageQuestionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuestion(val.(AccessPackageQuestionable))
        }
        return nil
    }
    return res
}
// GetQuestion gets the question property value. The question property
// returns a AccessPackageQuestionable when successful
func (m *AccessPackageResourceAttributeQuestion) GetQuestion()(AccessPackageQuestionable) {
    val, err := m.GetBackingStore().Get("question")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageQuestionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageResourceAttributeQuestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessPackageResourceAttributeSource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("question", m.GetQuestion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetQuestion sets the question property value. The question property
func (m *AccessPackageResourceAttributeQuestion) SetQuestion(value AccessPackageQuestionable)() {
    err := m.GetBackingStore().Set("question", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageResourceAttributeQuestionable interface {
    AccessPackageResourceAttributeSourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetQuestion()(AccessPackageQuestionable)
    SetQuestion(value AccessPackageQuestionable)()
}
