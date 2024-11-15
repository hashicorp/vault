package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRbacResourceNamespace struct {
    Entity
}
// NewUnifiedRbacResourceNamespace instantiates a new UnifiedRbacResourceNamespace and sets the default values.
func NewUnifiedRbacResourceNamespace()(*UnifiedRbacResourceNamespace) {
    m := &UnifiedRbacResourceNamespace{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUnifiedRbacResourceNamespaceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRbacResourceNamespaceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRbacResourceNamespace(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRbacResourceNamespace) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["resourceActions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRbacResourceActionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRbacResourceActionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRbacResourceActionable)
                }
            }
            m.SetResourceActions(res)
        }
        return nil
    }
    return res
}
// GetName gets the name property value. The name property
// returns a *string when successful
func (m *UnifiedRbacResourceNamespace) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResourceActions gets the resourceActions property value. The resourceActions property
// returns a []UnifiedRbacResourceActionable when successful
func (m *UnifiedRbacResourceNamespace) GetResourceActions()([]UnifiedRbacResourceActionable) {
    val, err := m.GetBackingStore().Get("resourceActions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRbacResourceActionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRbacResourceNamespace) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetResourceActions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceActions()))
        for i, v := range m.GetResourceActions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resourceActions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetName sets the name property value. The name property
func (m *UnifiedRbacResourceNamespace) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceActions sets the resourceActions property value. The resourceActions property
func (m *UnifiedRbacResourceNamespace) SetResourceActions(value []UnifiedRbacResourceActionable)() {
    err := m.GetBackingStore().Set("resourceActions", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRbacResourceNamespaceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetName()(*string)
    GetResourceActions()([]UnifiedRbacResourceActionable)
    SetName(value *string)()
    SetResourceActions(value []UnifiedRbacResourceActionable)()
}
