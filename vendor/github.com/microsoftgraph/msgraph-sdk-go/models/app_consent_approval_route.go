package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AppConsentApprovalRoute struct {
    Entity
}
// NewAppConsentApprovalRoute instantiates a new AppConsentApprovalRoute and sets the default values.
func NewAppConsentApprovalRoute()(*AppConsentApprovalRoute) {
    m := &AppConsentApprovalRoute{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAppConsentApprovalRouteFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppConsentApprovalRouteFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppConsentApprovalRoute(), nil
}
// GetAppConsentRequests gets the appConsentRequests property value. A collection of appConsentRequest objects representing apps for which admin consent has been requested by one or more users.
// returns a []AppConsentRequestable when successful
func (m *AppConsentApprovalRoute) GetAppConsentRequests()([]AppConsentRequestable) {
    val, err := m.GetBackingStore().Get("appConsentRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppConsentRequestable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AppConsentApprovalRoute) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appConsentRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppConsentRequestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppConsentRequestable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppConsentRequestable)
                }
            }
            m.SetAppConsentRequests(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AppConsentApprovalRoute) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppConsentRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppConsentRequests()))
        for i, v := range m.GetAppConsentRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appConsentRequests", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppConsentRequests sets the appConsentRequests property value. A collection of appConsentRequest objects representing apps for which admin consent has been requested by one or more users.
func (m *AppConsentApprovalRoute) SetAppConsentRequests(value []AppConsentRequestable)() {
    err := m.GetBackingStore().Set("appConsentRequests", value)
    if err != nil {
        panic(err)
    }
}
type AppConsentApprovalRouteable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppConsentRequests()([]AppConsentRequestable)
    SetAppConsentRequests(value []AppConsentRequestable)()
}
