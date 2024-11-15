package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ProfileCardProperty struct {
    Entity
}
// NewProfileCardProperty instantiates a new ProfileCardProperty and sets the default values.
func NewProfileCardProperty()(*ProfileCardProperty) {
    m := &ProfileCardProperty{
        Entity: *NewEntity(),
    }
    return m
}
// CreateProfileCardPropertyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateProfileCardPropertyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewProfileCardProperty(), nil
}
// GetAnnotations gets the annotations property value. Allows an administrator to set a custom display label for the directory property and localize it for the users in their tenant.
// returns a []ProfileCardAnnotationable when successful
func (m *ProfileCardProperty) GetAnnotations()([]ProfileCardAnnotationable) {
    val, err := m.GetBackingStore().Get("annotations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProfileCardAnnotationable)
    }
    return nil
}
// GetDirectoryPropertyName gets the directoryPropertyName property value. Identifies a profileCardProperty resource in Get, Update, or Delete operations. Allows an administrator to surface hidden Microsoft Entra ID properties on the Microsoft 365 profile card within their tenant. When present, the Microsoft Entra ID field referenced in this property is visible to all users in your tenant on the contact pane of the profile card. Allowed values for this field are: UserPrincipalName, Fax, StreetAddress, PostalCode, StateOrProvince, Alias, CustomAttribute1,  CustomAttribute2, CustomAttribute3, CustomAttribute4, CustomAttribute5, CustomAttribute6, CustomAttribute7, CustomAttribute8, CustomAttribute9, CustomAttribute10, CustomAttribute11, CustomAttribute12, CustomAttribute13, CustomAttribute14, CustomAttribute15.
// returns a *string when successful
func (m *ProfileCardProperty) GetDirectoryPropertyName()(*string) {
    val, err := m.GetBackingStore().Get("directoryPropertyName")
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
func (m *ProfileCardProperty) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["annotations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateProfileCardAnnotationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ProfileCardAnnotationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ProfileCardAnnotationable)
                }
            }
            m.SetAnnotations(res)
        }
        return nil
    }
    res["directoryPropertyName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirectoryPropertyName(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ProfileCardProperty) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAnnotations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAnnotations()))
        for i, v := range m.GetAnnotations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("annotations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("directoryPropertyName", m.GetDirectoryPropertyName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnnotations sets the annotations property value. Allows an administrator to set a custom display label for the directory property and localize it for the users in their tenant.
func (m *ProfileCardProperty) SetAnnotations(value []ProfileCardAnnotationable)() {
    err := m.GetBackingStore().Set("annotations", value)
    if err != nil {
        panic(err)
    }
}
// SetDirectoryPropertyName sets the directoryPropertyName property value. Identifies a profileCardProperty resource in Get, Update, or Delete operations. Allows an administrator to surface hidden Microsoft Entra ID properties on the Microsoft 365 profile card within their tenant. When present, the Microsoft Entra ID field referenced in this property is visible to all users in your tenant on the contact pane of the profile card. Allowed values for this field are: UserPrincipalName, Fax, StreetAddress, PostalCode, StateOrProvince, Alias, CustomAttribute1,  CustomAttribute2, CustomAttribute3, CustomAttribute4, CustomAttribute5, CustomAttribute6, CustomAttribute7, CustomAttribute8, CustomAttribute9, CustomAttribute10, CustomAttribute11, CustomAttribute12, CustomAttribute13, CustomAttribute14, CustomAttribute15.
func (m *ProfileCardProperty) SetDirectoryPropertyName(value *string)() {
    err := m.GetBackingStore().Set("directoryPropertyName", value)
    if err != nil {
        panic(err)
    }
}
type ProfileCardPropertyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAnnotations()([]ProfileCardAnnotationable)
    GetDirectoryPropertyName()(*string)
    SetAnnotations(value []ProfileCardAnnotationable)()
    SetDirectoryPropertyName(value *string)()
}
