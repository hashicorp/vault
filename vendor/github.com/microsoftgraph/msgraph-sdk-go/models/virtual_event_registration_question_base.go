package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventRegistrationQuestionBase struct {
    Entity
}
// NewVirtualEventRegistrationQuestionBase instantiates a new VirtualEventRegistrationQuestionBase and sets the default values.
func NewVirtualEventRegistrationQuestionBase()(*VirtualEventRegistrationQuestionBase) {
    m := &VirtualEventRegistrationQuestionBase{
        Entity: *NewEntity(),
    }
    return m
}
// CreateVirtualEventRegistrationQuestionBaseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventRegistrationQuestionBaseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.virtualEventRegistrationCustomQuestion":
                        return NewVirtualEventRegistrationCustomQuestion(), nil
                    case "#microsoft.graph.virtualEventRegistrationPredefinedQuestion":
                        return NewVirtualEventRegistrationPredefinedQuestion(), nil
                }
            }
        }
    }
    return NewVirtualEventRegistrationQuestionBase(), nil
}
// GetDisplayName gets the displayName property value. Display name of the registration question.
// returns a *string when successful
func (m *VirtualEventRegistrationQuestionBase) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventRegistrationQuestionBase) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
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
    return res
}
// GetIsRequired gets the isRequired property value. Indicates whether an answer to the question is required. The default value is false.
// returns a *bool when successful
func (m *VirtualEventRegistrationQuestionBase) GetIsRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventRegistrationQuestionBase) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
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
    return nil
}
// SetDisplayName sets the displayName property value. Display name of the registration question.
func (m *VirtualEventRegistrationQuestionBase) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRequired sets the isRequired property value. Indicates whether an answer to the question is required. The default value is false.
func (m *VirtualEventRegistrationQuestionBase) SetIsRequired(value *bool)() {
    err := m.GetBackingStore().Set("isRequired", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventRegistrationQuestionBaseable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetIsRequired()(*bool)
    SetDisplayName(value *string)()
    SetIsRequired(value *bool)()
}
