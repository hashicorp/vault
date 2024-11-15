package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnAttributeCollectionListener struct {
    AuthenticationEventListener
}
// NewOnAttributeCollectionListener instantiates a new OnAttributeCollectionListener and sets the default values.
func NewOnAttributeCollectionListener()(*OnAttributeCollectionListener) {
    m := &OnAttributeCollectionListener{
        AuthenticationEventListener: *NewAuthenticationEventListener(),
    }
    odataTypeValue := "#microsoft.graph.onAttributeCollectionListener"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnAttributeCollectionListenerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnAttributeCollectionListenerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnAttributeCollectionListener(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnAttributeCollectionListener) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationEventListener.GetFieldDeserializers()
    res["handler"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnAttributeCollectionHandlerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHandler(val.(OnAttributeCollectionHandlerable))
        }
        return nil
    }
    return res
}
// GetHandler gets the handler property value. Required. Configuration for what to invoke if the event resolves to this listener.
// returns a OnAttributeCollectionHandlerable when successful
func (m *OnAttributeCollectionListener) GetHandler()(OnAttributeCollectionHandlerable) {
    val, err := m.GetBackingStore().Get("handler")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnAttributeCollectionHandlerable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnAttributeCollectionListener) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetHandler sets the handler property value. Required. Configuration for what to invoke if the event resolves to this listener.
func (m *OnAttributeCollectionListener) SetHandler(value OnAttributeCollectionHandlerable)() {
    err := m.GetBackingStore().Set("handler", value)
    if err != nil {
        panic(err)
    }
}
type OnAttributeCollectionListenerable interface {
    AuthenticationEventListenerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetHandler()(OnAttributeCollectionHandlerable)
    SetHandler(value OnAttributeCollectionHandlerable)()
}
