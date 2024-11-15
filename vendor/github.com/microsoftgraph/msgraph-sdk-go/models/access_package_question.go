package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageQuestion struct {
    Entity
}
// NewAccessPackageQuestion instantiates a new AccessPackageQuestion and sets the default values.
func NewAccessPackageQuestion()(*AccessPackageQuestion) {
    m := &AccessPackageQuestion{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessPackageQuestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageQuestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.accessPackageMultipleChoiceQuestion":
                        return NewAccessPackageMultipleChoiceQuestion(), nil
                    case "#microsoft.graph.accessPackageTextInputQuestion":
                        return NewAccessPackageTextInputQuestion(), nil
                }
            }
        }
    }
    return NewAccessPackageQuestion(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageQuestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["isAnswerEditable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAnswerEditable(val)
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
    res["localizations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageLocalizedTextFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageLocalizedTextable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageLocalizedTextable)
                }
            }
            m.SetLocalizations(res)
        }
        return nil
    }
    res["sequence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSequence(val)
        }
        return nil
    }
    res["text"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetText(val)
        }
        return nil
    }
    return res
}
// GetIsAnswerEditable gets the isAnswerEditable property value. Specifies whether the requestor is allowed to edit answers to questions for an assignment by posting an update to accessPackageAssignmentRequest.
// returns a *bool when successful
func (m *AccessPackageQuestion) GetIsAnswerEditable()(*bool) {
    val, err := m.GetBackingStore().Get("isAnswerEditable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRequired gets the isRequired property value. Whether the requestor is required to supply an answer or not.
// returns a *bool when successful
func (m *AccessPackageQuestion) GetIsRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocalizations gets the localizations property value. The text of the question represented in a format for a specific locale.
// returns a []AccessPackageLocalizedTextable when successful
func (m *AccessPackageQuestion) GetLocalizations()([]AccessPackageLocalizedTextable) {
    val, err := m.GetBackingStore().Get("localizations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageLocalizedTextable)
    }
    return nil
}
// GetSequence gets the sequence property value. Relative position of this question when displaying a list of questions to the requestor.
// returns a *int32 when successful
func (m *AccessPackageQuestion) GetSequence()(*int32) {
    val, err := m.GetBackingStore().Get("sequence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetText gets the text property value. The text of the question to show to the requestor.
// returns a *string when successful
func (m *AccessPackageQuestion) GetText()(*string) {
    val, err := m.GetBackingStore().Get("text")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageQuestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isAnswerEditable", m.GetIsAnswerEditable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isRequired", m.GetIsRequired())
        if err != nil {
            return err
        }
    }
    if m.GetLocalizations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLocalizations()))
        for i, v := range m.GetLocalizations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("localizations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("sequence", m.GetSequence())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("text", m.GetText())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsAnswerEditable sets the isAnswerEditable property value. Specifies whether the requestor is allowed to edit answers to questions for an assignment by posting an update to accessPackageAssignmentRequest.
func (m *AccessPackageQuestion) SetIsAnswerEditable(value *bool)() {
    err := m.GetBackingStore().Set("isAnswerEditable", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRequired sets the isRequired property value. Whether the requestor is required to supply an answer or not.
func (m *AccessPackageQuestion) SetIsRequired(value *bool)() {
    err := m.GetBackingStore().Set("isRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetLocalizations sets the localizations property value. The text of the question represented in a format for a specific locale.
func (m *AccessPackageQuestion) SetLocalizations(value []AccessPackageLocalizedTextable)() {
    err := m.GetBackingStore().Set("localizations", value)
    if err != nil {
        panic(err)
    }
}
// SetSequence sets the sequence property value. Relative position of this question when displaying a list of questions to the requestor.
func (m *AccessPackageQuestion) SetSequence(value *int32)() {
    err := m.GetBackingStore().Set("sequence", value)
    if err != nil {
        panic(err)
    }
}
// SetText sets the text property value. The text of the question to show to the requestor.
func (m *AccessPackageQuestion) SetText(value *string)() {
    err := m.GetBackingStore().Set("text", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageQuestionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsAnswerEditable()(*bool)
    GetIsRequired()(*bool)
    GetLocalizations()([]AccessPackageLocalizedTextable)
    GetSequence()(*int32)
    GetText()(*string)
    SetIsAnswerEditable(value *bool)()
    SetIsRequired(value *bool)()
    SetLocalizations(value []AccessPackageLocalizedTextable)()
    SetSequence(value *int32)()
    SetText(value *string)()
}
