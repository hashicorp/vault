package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnAuthenticationMethodLoadStartListener struct {
    AuthenticationEventListener
}
// NewOnAuthenticationMethodLoadStartListener instantiates a new OnAuthenticationMethodLoadStartListener and sets the default values.
func NewOnAuthenticationMethodLoadStartListener()(*OnAuthenticationMethodLoadStartListener) {
    m := &OnAuthenticationMethodLoadStartListener{
        AuthenticationEventListener: *NewAuthenticationEventListener(),
    }
    odataTypeValue := "#microsoft.graph.onAuthenticationMethodLoadStartListener"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnAuthenticationMethodLoadStartListenerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnAuthenticationMethodLoadStartListenerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnAuthenticationMethodLoadStartListener(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnAuthenticationMethodLoadStartListener) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationEventListener.GetFieldDeserializers()
    res["handler"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnAuthenticationMethodLoadStartHandlerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHandler(val.(OnAuthenticationMethodLoadStartHandlerable))
        }
        return nil
    }
    return res
}
// GetHandler gets the handler property value. Required. Configuration for what to invoke if the event resolves to this listener. This property lets us define potential handler configurations per-event.
// returns a OnAuthenticationMethodLoadStartHandlerable when successful
func (m *OnAuthenticationMethodLoadStartListener) GetHandler()(OnAuthenticationMethodLoadStartHandlerable) {
    val, err := m.GetBackingStore().Get("handler")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnAuthenticationMethodLoadStartHandlerable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnAuthenticationMethodLoadStartListener) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetHandler sets the handler property value. Required. Configuration for what to invoke if the event resolves to this listener. This property lets us define potential handler configurations per-event.
func (m *OnAuthenticationMethodLoadStartListener) SetHandler(value OnAuthenticationMethodLoadStartHandlerable)() {
    err := m.GetBackingStore().Set("handler", value)
    if err != nil {
        panic(err)
    }
}
type OnAuthenticationMethodLoadStartListenerable interface {
    AuthenticationEventListenerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetHandler()(OnAuthenticationMethodLoadStartHandlerable)
    SetHandler(value OnAuthenticationMethodLoadStartHandlerable)()
}
