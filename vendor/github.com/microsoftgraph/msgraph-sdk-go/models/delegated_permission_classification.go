package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DelegatedPermissionClassification struct {
    Entity
}
// NewDelegatedPermissionClassification instantiates a new DelegatedPermissionClassification and sets the default values.
func NewDelegatedPermissionClassification()(*DelegatedPermissionClassification) {
    m := &DelegatedPermissionClassification{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDelegatedPermissionClassificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDelegatedPermissionClassificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDelegatedPermissionClassification(), nil
}
// GetClassification gets the classification property value. The classification value. Possible values: low, medium (preview), high (preview). Doesn't support $filter.
// returns a *PermissionClassificationType when successful
func (m *DelegatedPermissionClassification) GetClassification()(*PermissionClassificationType) {
    val, err := m.GetBackingStore().Get("classification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PermissionClassificationType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DelegatedPermissionClassification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["classification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePermissionClassificationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassification(val.(*PermissionClassificationType))
        }
        return nil
    }
    res["permissionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermissionId(val)
        }
        return nil
    }
    res["permissionName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermissionName(val)
        }
        return nil
    }
    return res
}
// GetPermissionId gets the permissionId property value. The unique identifier (id) for the delegated permission listed in the oauth2PermissionScopes collection of the servicePrincipal. Required on create. Doesn't support $filter.
// returns a *string when successful
func (m *DelegatedPermissionClassification) GetPermissionId()(*string) {
    val, err := m.GetBackingStore().Get("permissionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPermissionName gets the permissionName property value. The claim value (value) for the delegated permission listed in the oauth2PermissionScopes collection of the servicePrincipal. Doesn't support $filter.
// returns a *string when successful
func (m *DelegatedPermissionClassification) GetPermissionName()(*string) {
    val, err := m.GetBackingStore().Get("permissionName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DelegatedPermissionClassification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetClassification() != nil {
        cast := (*m.GetClassification()).String()
        err = writer.WriteStringValue("classification", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("permissionId", m.GetPermissionId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("permissionName", m.GetPermissionName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClassification sets the classification property value. The classification value. Possible values: low, medium (preview), high (preview). Doesn't support $filter.
func (m *DelegatedPermissionClassification) SetClassification(value *PermissionClassificationType)() {
    err := m.GetBackingStore().Set("classification", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionId sets the permissionId property value. The unique identifier (id) for the delegated permission listed in the oauth2PermissionScopes collection of the servicePrincipal. Required on create. Doesn't support $filter.
func (m *DelegatedPermissionClassification) SetPermissionId(value *string)() {
    err := m.GetBackingStore().Set("permissionId", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionName sets the permissionName property value. The claim value (value) for the delegated permission listed in the oauth2PermissionScopes collection of the servicePrincipal. Doesn't support $filter.
func (m *DelegatedPermissionClassification) SetPermissionName(value *string)() {
    err := m.GetBackingStore().Set("permissionName", value)
    if err != nil {
        panic(err)
    }
}
type DelegatedPermissionClassificationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClassification()(*PermissionClassificationType)
    GetPermissionId()(*string)
    GetPermissionName()(*string)
    SetClassification(value *PermissionClassificationType)()
    SetPermissionId(value *string)()
    SetPermissionName(value *string)()
}
