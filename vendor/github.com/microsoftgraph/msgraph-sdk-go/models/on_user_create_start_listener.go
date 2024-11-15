package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnUserCreateStartListener struct {
    AuthenticationEventListener
}
// NewOnUserCreateStartListener instantiates a new OnUserCreateStartListener and sets the default values.
func NewOnUserCreateStartListener()(*OnUserCreateStartListener) {
    m := &OnUserCreateStartListener{
        AuthenticationEventListener: *NewAuthenticationEventListener(),
    }
    odataTypeValue := "#microsoft.graph.onUserCreateStartListener"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnUserCreateStartListenerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnUserCreateStartListenerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnUserCreateStartListener(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnUserCreateStartListener) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationEventListener.GetFieldDeserializers()
    res["handler"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnUserCreateStartHandlerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHandler(val.(OnUserCreateStartHandlerable))
        }
        return nil
    }
    return res
}
// GetHandler gets the handler property value. Required. Configuration for what to invoke if the event resolves to this listener. This lets us define potential handler configurations per-event.
// returns a OnUserCreateStartHandlerable when successful
func (m *OnUserCreateStartListener) GetHandler()(OnUserCreateStartHandlerable) {
    val, err := m.GetBackingStore().Get("handler")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnUserCreateStartHandlerable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnUserCreateStartListener) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationEventListener.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("handler", m.GetHandler())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetHandler sets the handler property value. Required. Configuration for what to invoke if the event resolves to this listener. This lets us define potential handler configurations per-event.
func (m *OnUserCreateStartListener) SetHandler(value OnUserCreateStartHandlerable)() {
    err := m.GetBackingStore().Set("handler", value)
    if err != nil {
        panic(err)
    }
}
type OnUserCreateStartListenerable interface {
    AuthenticationEventListenerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetHandler()(OnUserCreateStartHandlerable)
    SetHandler(value OnUserCreateStartHandlerable)()
}
