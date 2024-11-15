package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PermissionGrantPolicy struct {
    PolicyBase
}
// NewPermissionGrantPolicy instantiates a new PermissionGrantPolicy and sets the default values.
func NewPermissionGrantPolicy()(*PermissionGrantPolicy) {
    m := &PermissionGrantPolicy{
        PolicyBase: *NewPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.permissionGrantPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePermissionGrantPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePermissionGrantPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPermissionGrantPolicy(), nil
}
// GetExcludes gets the excludes property value. Condition sets that are excluded in this permission grant policy. Automatically expanded on GET.
// returns a []PermissionGrantConditionSetable when successful
func (m *PermissionGrantPolicy) GetExcludes()([]PermissionGrantConditionSetable) {
    val, err := m.GetBackingStore().Get("excludes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PermissionGrantConditionSetable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PermissionGrantPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PolicyBase.GetFieldDeserializers()
    res["excludes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePermissionGrantConditionSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PermissionGrantConditionSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PermissionGrantConditionSetable)
                }
            }
            m.SetExcludes(res)
        }
        return nil
    }
    res["includes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePermissionGrantConditionSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PermissionGrantConditionSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PermissionGrantConditionSetable)
                }
            }
            m.SetIncludes(res)
        }
        return nil
    }
    return res
}
// GetIncludes gets the includes property value. Condition sets that are included in this permission grant policy. Automatically expanded on GET.
// returns a []PermissionGrantConditionSetable when successful
func (m *PermissionGrantPolicy) GetIncludes()([]PermissionGrantConditionSetable) {
    val, err := m.GetBackingStore().Get("includes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PermissionGrantConditionSetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PermissionGrantPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetExcludes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExcludes()))
        for i, v := range m.GetExcludes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("excludes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetIncludes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIncludes()))
        for i, v := range m.GetIncludes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("includes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExcludes sets the excludes property value. Condition sets that are excluded in this permission grant policy. Automatically expanded on GET.
func (m *PermissionGrantPolicy) SetExcludes(value []PermissionGrantConditionSetable)() {
    err := m.GetBackingStore().Set("excludes", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludes sets the includes property value. Condition sets that are included in this permission grant policy. Automatically expanded on GET.
func (m *PermissionGrantPolicy) SetIncludes(value []PermissionGrantConditionSetable)() {
    err := m.GetBackingStore().Set("includes", value)
    if err != nil {
        panic(err)
    }
}
type PermissionGrantPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PolicyBaseable
    GetExcludes()([]PermissionGrantConditionSetable)
    GetIncludes()([]PermissionGrantConditionSetable)
    SetExcludes(value []PermissionGrantConditionSetable)()
    SetIncludes(value []PermissionGrantConditionSetable)()
}
