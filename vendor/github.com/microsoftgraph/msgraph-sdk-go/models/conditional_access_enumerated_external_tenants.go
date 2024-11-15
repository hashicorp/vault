package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ConditionalAccessEnumeratedExternalTenants struct {
    ConditionalAccessExternalTenants
}
// NewConditionalAccessEnumeratedExternalTenants instantiates a new ConditionalAccessEnumeratedExternalTenants and sets the default values.
func NewConditionalAccessEnumeratedExternalTenants()(*ConditionalAccessEnumeratedExternalTenants) {
    m := &ConditionalAccessEnumeratedExternalTenants{
        ConditionalAccessExternalTenants: *NewConditionalAccessExternalTenants(),
    }
    odataTypeValue := "#microsoft.graph.conditionalAccessEnumeratedExternalTenants"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateConditionalAccessEnumeratedExternalTenantsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessEnumeratedExternalTenantsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessEnumeratedExternalTenants(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessEnumeratedExternalTenants) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConditionalAccessExternalTenants.GetFieldDeserializers()
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMembers(res)
        }
        return nil
    }
    return res
}
// GetMembers gets the members property value. A collection of tenant IDs that define the scope of a policy targeting conditional access for guests and external users.
// returns a []string when successful
func (m *ConditionalAccessEnumeratedExternalTenants) GetMembers()([]string) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessEnumeratedExternalTenants) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConditionalAccessExternalTenants.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetMembers() != nil {
        err = writer.WriteCollectionOfStringValues("members", m.GetMembers())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMembers sets the members property value. A collection of tenant IDs that define the scope of a policy targeting conditional access for guests and external users.
func (m *ConditionalAccessEnumeratedExternalTenants) SetMembers(value []string)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessEnumeratedExternalTenantsable interface {
    ConditionalAccessExternalTenantsable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMembers()([]string)
    SetMembers(value []string)()
}
