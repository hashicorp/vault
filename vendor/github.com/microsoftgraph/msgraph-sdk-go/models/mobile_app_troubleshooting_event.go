package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MobileAppTroubleshootingEvent struct {
    Entity
}
// NewMobileAppTroubleshootingEvent instantiates a new MobileAppTroubleshootingEvent and sets the default values.
func NewMobileAppTroubleshootingEvent()(*MobileAppTroubleshootingEvent) {
    m := &MobileAppTroubleshootingEvent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMobileAppTroubleshootingEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMobileAppTroubleshootingEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMobileAppTroubleshootingEvent(), nil
}
// GetAppLogCollectionRequests gets the appLogCollectionRequests property value. Indicates collection of App Log Upload Request.
// returns a []AppLogCollectionRequestable when successful
func (m *MobileAppTroubleshootingEvent) GetAppLogCollectionRequests()([]AppLogCollectionRequestable) {
    val, err := m.GetBackingStore().Get("appLogCollectionRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppLogCollectionRequestable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MobileAppTroubleshootingEvent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appLogCollectionRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppLogCollectionRequestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppLogCollectionRequestable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppLogCollectionRequestable)
                }
            }
            m.SetAppLogCollectionRequests(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *MobileAppTroubleshootingEvent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppLogCollectionRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppLogCollectionRequests()))
        for i, v := range m.GetAppLogCollectionRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appLogCollectionRequests", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppLogCollectionRequests sets the appLogCollectionRequests property value. Indicates collection of App Log Upload Request.
func (m *MobileAppTroubleshootingEvent) SetAppLogCollectionRequests(value []AppLogCollectionRequestable)() {
    err := m.GetBackingStore().Set("appLogCollectionRequests", value)
    if err != nil {
        panic(err)
    }
}
type MobileAppTroubleshootingEventable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppLogCollectionRequests()([]AppLogCollectionRequestable)
    SetAppLogCollectionRequests(value []AppLogCollectionRequestable)()
}
