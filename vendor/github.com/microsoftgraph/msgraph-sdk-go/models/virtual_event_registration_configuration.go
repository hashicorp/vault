package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventRegistrationConfiguration struct {
    Entity
}
// NewVirtualEventRegistrationConfiguration instantiates a new VirtualEventRegistrationConfiguration and sets the default values.
func NewVirtualEventRegistrationConfiguration()(*VirtualEventRegistrationConfiguration) {
    m := &VirtualEventRegistrationConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateVirtualEventRegistrationConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventRegistrationConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.virtualEventWebinarRegistrationConfiguration":
                        return NewVirtualEventWebinarRegistrationConfiguration(), nil
                }
            }
        }
    }
    return NewVirtualEventRegistrationConfiguration(), nil
}
// GetCapacity gets the capacity property value. Total capacity of the virtual event.
// returns a *int32 when successful
func (m *VirtualEventRegistrationConfiguration) GetCapacity()(*int32) {
    val, err := m.GetBackingStore().Get("capacity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventRegistrationConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["capacity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCapacity(val)
        }
        return nil
    }
    res["questions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateVirtualEventRegistrationQuestionBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]VirtualEventRegistrationQuestionBaseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(VirtualEventRegistrationQuestionBaseable)
                }
            }
            m.SetQuestions(res)
        }
        return nil
    }
    res["registrationWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistrationWebUrl(val)
        }
        return nil
    }
    return res
}
// GetQuestions gets the questions property value. Registration questions.
// returns a []VirtualEventRegistrationQuestionBaseable when successful
func (m *VirtualEventRegistrationConfiguration) GetQuestions()([]VirtualEventRegistrationQuestionBaseable) {
    val, err := m.GetBackingStore().Get("questions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]VirtualEventRegistrationQuestionBaseable)
    }
    return nil
}
// GetRegistrationWebUrl gets the registrationWebUrl property value. Registration URL of the virtual event.
// returns a *string when successful
func (m *VirtualEventRegistrationConfiguration) GetRegistrationWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("registrationWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventRegistrationConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("capacity", m.GetCapacity())
        if err != nil {
            return err
        }
    }
    if m.GetQuestions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetQuestions()))
        for i, v := range m.GetQuestions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("questions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registrationWebUrl", m.GetRegistrationWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCapacity sets the capacity property value. Total capacity of the virtual event.
func (m *VirtualEventRegistrationConfiguration) SetCapacity(value *int32)() {
    err := m.GetBackingStore().Set("capacity", value)
    if err != nil {
        panic(err)
    }
}
// SetQuestions sets the questions property value. Registration questions.
func (m *VirtualEventRegistrationConfiguration) SetQuestions(value []VirtualEventRegistrationQuestionBaseable)() {
    err := m.GetBackingStore().Set("questions", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistrationWebUrl sets the registrationWebUrl property value. Registration URL of the virtual event.
func (m *VirtualEventRegistrationConfiguration) SetRegistrationWebUrl(value *string)() {
    err := m.GetBackingStore().Set("registrationWebUrl", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventRegistrationConfigurationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCapacity()(*int32)
    GetQuestions()([]VirtualEventRegistrationQuestionBaseable)
    GetRegistrationWebUrl()(*string)
    SetCapacity(value *int32)()
    SetQuestions(value []VirtualEventRegistrationQuestionBaseable)()
    SetRegistrationWebUrl(value *string)()
}
