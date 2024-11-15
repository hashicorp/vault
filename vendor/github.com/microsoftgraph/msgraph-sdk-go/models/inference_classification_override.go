package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type InferenceClassificationOverride struct {
    Entity
}
// NewInferenceClassificationOverride instantiates a new InferenceClassificationOverride and sets the default values.
func NewInferenceClassificationOverride()(*InferenceClassificationOverride) {
    m := &InferenceClassificationOverride{
        Entity: *NewEntity(),
    }
    return m
}
// CreateInferenceClassificationOverrideFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateInferenceClassificationOverrideFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInferenceClassificationOverride(), nil
}
// GetClassifyAs gets the classifyAs property value. Specifies how incoming messages from a specific sender should always be classified as. The possible values are: focused, other.
// returns a *InferenceClassificationType when successful
func (m *InferenceClassificationOverride) GetClassifyAs()(*InferenceClassificationType) {
    val, err := m.GetBackingStore().Get("classifyAs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*InferenceClassificationType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *InferenceClassificationOverride) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["classifyAs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseInferenceClassificationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassifyAs(val.(*InferenceClassificationType))
        }
        return nil
    }
    res["senderEmailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSenderEmailAddress(val.(EmailAddressable))
        }
        return nil
    }
    return res
}
// GetSenderEmailAddress gets the senderEmailAddress property value. The email address information of the sender for whom the override is created.
// returns a EmailAddressable when successful
func (m *InferenceClassificationOverride) GetSenderEmailAddress()(EmailAddressable) {
    val, err := m.GetBackingStore().Get("senderEmailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailAddressable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *InferenceClassificationOverride) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetClassifyAs() != nil {
        cast := (*m.GetClassifyAs()).String()
        err = writer.WriteStringValue("classifyAs", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("senderEmailAddress", m.GetSenderEmailAddress())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClassifyAs sets the classifyAs property value. Specifies how incoming messages from a specific sender should always be classified as. The possible values are: focused, other.
func (m *InferenceClassificationOverride) SetClassifyAs(value *InferenceClassificationType)() {
    err := m.GetBackingStore().Set("classifyAs", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderEmailAddress sets the senderEmailAddress property value. The email address information of the sender for whom the override is created.
func (m *InferenceClassificationOverride) SetSenderEmailAddress(value EmailAddressable)() {
    err := m.GetBackingStore().Set("senderEmailAddress", value)
    if err != nil {
        panic(err)
    }
}
type InferenceClassificationOverrideable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClassifyAs()(*InferenceClassificationType)
    GetSenderEmailAddress()(EmailAddressable)
    SetClassifyAs(value *InferenceClassificationType)()
    SetSenderEmailAddress(value EmailAddressable)()
}
