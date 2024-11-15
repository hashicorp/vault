package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamworkHostedContent struct {
    Entity
}
// NewTeamworkHostedContent instantiates a new TeamworkHostedContent and sets the default values.
func NewTeamworkHostedContent()(*TeamworkHostedContent) {
    m := &TeamworkHostedContent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamworkHostedContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamworkHostedContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.chatMessageHostedContent":
                        return NewChatMessageHostedContent(), nil
                }
            }
        }
    }
    return NewTeamworkHostedContent(), nil
}
// GetContentBytes gets the contentBytes property value. Write only. Bytes for the hosted content (such as images).
// returns a []byte when successful
func (m *TeamworkHostedContent) GetContentBytes()([]byte) {
    val, err := m.GetBackingStore().Get("contentBytes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetContentType gets the contentType property value. Write only. Content type. such as image/png, image/jpg.
// returns a *string when successful
func (m *TeamworkHostedContent) GetContentType()(*string) {
    val, err := m.GetBackingStore().Get("contentType")
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
func (m *TeamworkHostedContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["contentBytes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentBytes(val)
        }
        return nil
    }
    res["contentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentType(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *TeamworkHostedContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteByteArrayValue("contentBytes", m.GetContentBytes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentType", m.GetContentType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContentBytes sets the contentBytes property value. Write only. Bytes for the hosted content (such as images).
func (m *TeamworkHostedContent) SetContentBytes(value []byte)() {
    err := m.GetBackingStore().Set("contentBytes", value)
    if err != nil {
        panic(err)
    }
}
// SetContentType sets the contentType property value. Write only. Content type. such as image/png, image/jpg.
func (m *TeamworkHostedContent) SetContentType(value *string)() {
    err := m.GetBackingStore().Set("contentType", value)
    if err != nil {
        panic(err)
    }
}
type TeamworkHostedContentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContentBytes()([]byte)
    GetContentType()(*string)
    SetContentBytes(value []byte)()
    SetContentType(value *string)()
}
