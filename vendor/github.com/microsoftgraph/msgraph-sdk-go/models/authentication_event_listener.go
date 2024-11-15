package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationEventListener struct {
    Entity
}
// NewAuthenticationEventListener instantiates a new AuthenticationEventListener and sets the default values.
func NewAuthenticationEventListener()(*AuthenticationEventListener) {
    m := &AuthenticationEventListener{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationEventListenerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationEventListenerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.onAttributeCollectionListener":
                        return NewOnAttributeCollectionListener(), nil
                    case "#microsoft.graph.onAuthenticationMethodLoadStartListener":
                        return NewOnAuthenticationMethodLoadStartListener(), nil
                    case "#microsoft.graph.onInteractiveAuthFlowStartListener":
                        return NewOnInteractiveAuthFlowStartListener(), nil
                    case "#microsoft.graph.onTokenIssuanceStartListener":
                        return NewOnTokenIssuanceStartListener(), nil
                    case "#microsoft.graph.onUserCreateStartListener":
                        return NewOnUserCreateStartListener(), nil
                }
            }
        }
    }
    return NewAuthenticationEventListener(), nil
}
// GetAuthenticationEventsFlowId gets the authenticationEventsFlowId property value. Indicates the authenticationEventListener is associated with an authenticationEventsFlow. Read-only.
// returns a *string when successful
func (m *AuthenticationEventListener) GetAuthenticationEventsFlowId()(*string) {
    val, err := m.GetBackingStore().Get("authenticationEventsFlowId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConditions gets the conditions property value. The conditions on which this authenticationEventListener should trigger.
// returns a AuthenticationConditionsable when successful
func (m *AuthenticationEventListener) GetConditions()(AuthenticationConditionsable) {
    val, err := m.GetBackingStore().Get("conditions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthenticationConditionsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthenticationEventListener) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationEventsFlowId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationEventsFlowId(val)
        }
        return nil
    }
    res["conditions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationConditionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditions(val.(AuthenticationConditionsable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AuthenticationEventListener) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("authenticationEventsFlowId", m.GetAuthenticationEventsFlowId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conditions", m.GetConditions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationEventsFlowId sets the authenticationEventsFlowId property value. Indicates the authenticationEventListener is associated with an authenticationEventsFlow. Read-only.
func (m *AuthenticationEventListener) SetAuthenticationEventsFlowId(value *string)() {
    err := m.GetBackingStore().Set("authenticationEventsFlowId", value)
    if err != nil {
        panic(err)
    }
}
// SetConditions sets the conditions property value. The conditions on which this authenticationEventListener should trigger.
func (m *AuthenticationEventListener) SetConditions(value AuthenticationConditionsable)() {
    err := m.GetBackingStore().Set("conditions", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationEventListenerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationEventsFlowId()(*string)
    GetConditions()(AuthenticationConditionsable)
    SetAuthenticationEventsFlowId(value *string)()
    SetConditions(value AuthenticationConditionsable)()
}
