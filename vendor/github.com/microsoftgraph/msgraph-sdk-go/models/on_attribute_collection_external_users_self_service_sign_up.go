package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnAttributeCollectionExternalUsersSelfServiceSignUp struct {
    OnAttributeCollectionHandler
}
// NewOnAttributeCollectionExternalUsersSelfServiceSignUp instantiates a new OnAttributeCollectionExternalUsersSelfServiceSignUp and sets the default values.
func NewOnAttributeCollectionExternalUsersSelfServiceSignUp()(*OnAttributeCollectionExternalUsersSelfServiceSignUp) {
    m := &OnAttributeCollectionExternalUsersSelfServiceSignUp{
        OnAttributeCollectionHandler: *NewOnAttributeCollectionHandler(),
    }
    odataTypeValue := "#microsoft.graph.onAttributeCollectionExternalUsersSelfServiceSignUp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnAttributeCollectionExternalUsersSelfServiceSignUpFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnAttributeCollectionExternalUsersSelfServiceSignUpFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnAttributeCollectionExternalUsersSelfServiceSignUp(), nil
}
// GetAttributeCollectionPage gets the attributeCollectionPage property value. Required. The configuration for how attributes are displayed in the sign-up experience defined by a user flow, like the externalUsersSelfServiceSignupEventsFlow, specifically on the attribute collection page.
// returns a AuthenticationAttributeCollectionPageable when successful
func (m *OnAttributeCollectionExternalUsersSelfServiceSignUp) GetAttributeCollectionPage()(AuthenticationAttributeCollectionPageable) {
    val, err := m.GetBackingStore().Get("attributeCollectionPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthenticationAttributeCollectionPageable)
    }
    return nil
}
// GetAttributes gets the attributes property value. The attributes property
// returns a []IdentityUserFlowAttributeable when successful
func (m *OnAttributeCollectionExternalUsersSelfServiceSignUp) GetAttributes()([]IdentityUserFlowAttributeable) {
    val, err := m.GetBackingStore().Get("attributes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityUserFlowAttributeable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnAttributeCollectionExternalUsersSelfServiceSignUp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OnAttributeCollectionHandler.GetFieldDeserializers()
    res["attributeCollectionPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationAttributeCollectionPageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttributeCollectionPage(val.(AuthenticationAttributeCollectionPageable))
        }
        return nil
    }
    res["attributes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityUserFlowAttributeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentityUserFlowAttributeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentityUserFlowAttributeable)
                }
            }
            m.SetAttributes(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *OnAttributeCollectionExternalUsersSelfServiceSignUp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OnAttributeCollectionHandler.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("attributeCollectionPage", m.GetAttributeCollectionPage())
        if err != nil {
            return err
        }
    }
    if m.GetAttributes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttributes()))
        for i, v := range m.GetAttributes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("attributes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttributeCollectionPage sets the attributeCollectionPage property value. Required. The configuration for how attributes are displayed in the sign-up experience defined by a user flow, like the externalUsersSelfServiceSignupEventsFlow, specifically on the attribute collection page.
func (m *OnAttributeCollectionExternalUsersSelfServiceSignUp) SetAttributeCollectionPage(value AuthenticationAttributeCollectionPageable)() {
    err := m.GetBackingStore().Set("attributeCollectionPage", value)
    if err != nil {
        panic(err)
    }
}
// SetAttributes sets the attributes property value. The attributes property
func (m *OnAttributeCollectionExternalUsersSelfServiceSignUp) SetAttributes(value []IdentityUserFlowAttributeable)() {
    err := m.GetBackingStore().Set("attributes", value)
    if err != nil {
        panic(err)
    }
}
type OnAttributeCollectionExternalUsersSelfServiceSignUpable interface {
    OnAttributeCollectionHandlerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttributeCollectionPage()(AuthenticationAttributeCollectionPageable)
    GetAttributes()([]IdentityUserFlowAttributeable)
    SetAttributeCollectionPage(value AuthenticationAttributeCollectionPageable)()
    SetAttributes(value []IdentityUserFlowAttributeable)()
}
