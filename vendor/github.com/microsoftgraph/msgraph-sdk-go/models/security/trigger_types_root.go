package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type TriggerTypesRoot struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewTriggerTypesRoot instantiates a new TriggerTypesRoot and sets the default values.
func NewTriggerTypesRoot()(*TriggerTypesRoot) {
    m := &TriggerTypesRoot{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateTriggerTypesRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTriggerTypesRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTriggerTypesRoot(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TriggerTypesRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["retentionEventTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRetentionEventTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RetentionEventTypeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RetentionEventTypeable)
                }
            }
            m.SetRetentionEventTypes(res)
        }
        return nil
    }
    return res
}
// GetRetentionEventTypes gets the retentionEventTypes property value. The retentionEventTypes property
// returns a []RetentionEventTypeable when successful
func (m *TriggerTypesRoot) GetRetentionEventTypes()([]RetentionEventTypeable) {
    val, err := m.GetBackingStore().Get("retentionEventTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RetentionEventTypeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TriggerTypesRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetRetentionEventTypes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRetentionEventTypes()))
        for i, v := range m.GetRetentionEventTypes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("retentionEventTypes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRetentionEventTypes sets the retentionEventTypes property value. The retentionEventTypes property
func (m *TriggerTypesRoot) SetRetentionEventTypes(value []RetentionEventTypeable)() {
    err := m.GetBackingStore().Set("retentionEventTypes", value)
    if err != nil {
        panic(err)
    }
}
type TriggerTypesRootable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRetentionEventTypes()([]RetentionEventTypeable)
    SetRetentionEventTypes(value []RetentionEventTypeable)()
}
