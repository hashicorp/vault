package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AzureCommunicationServicesUserConversationMember struct {
    ConversationMember
}
// NewAzureCommunicationServicesUserConversationMember instantiates a new AzureCommunicationServicesUserConversationMember and sets the default values.
func NewAzureCommunicationServicesUserConversationMember()(*AzureCommunicationServicesUserConversationMember) {
    m := &AzureCommunicationServicesUserConversationMember{
        ConversationMember: *NewConversationMember(),
    }
    odataTypeValue := "#microsoft.graph.azureCommunicationServicesUserConversationMember"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAzureCommunicationServicesUserConversationMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAzureCommunicationServicesUserConversationMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAzureCommunicationServicesUserConversationMember(), nil
}
// GetAzureCommunicationServicesId gets the azureCommunicationServicesId property value. Azure Communication Services ID of the user.
// returns a *string when successful
func (m *AzureCommunicationServicesUserConversationMember) GetAzureCommunicationServicesId()(*string) {
    val, err := m.GetBackingStore().Get("azureCommunicationServicesId")
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
func (m *AzureCommunicationServicesUserConversationMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConversationMember.GetFieldDeserializers()
    res["azureCommunicationServicesId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureCommunicationServicesId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AzureCommunicationServicesUserConversationMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConversationMember.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("azureCommunicationServicesId", m.GetAzureCommunicationServicesId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAzureCommunicationServicesId sets the azureCommunicationServicesId property value. Azure Communication Services ID of the user.
func (m *AzureCommunicationServicesUserConversationMember) SetAzureCommunicationServicesId(value *string)() {
    err := m.GetBackingStore().Set("azureCommunicationServicesId", value)
    if err != nil {
        panic(err)
    }
}
type AzureCommunicationServicesUserConversationMemberable interface {
    ConversationMemberable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAzureCommunicationServicesId()(*string)
    SetAzureCommunicationServicesId(value *string)()
}
