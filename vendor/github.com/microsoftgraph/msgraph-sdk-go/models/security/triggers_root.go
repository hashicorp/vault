package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type TriggersRoot struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewTriggersRoot instantiates a new TriggersRoot and sets the default values.
func NewTriggersRoot()(*TriggersRoot) {
    m := &TriggersRoot{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateTriggersRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTriggersRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTriggersRoot(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TriggersRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["retentionEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRetentionEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RetentionEventable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RetentionEventable)
                }
            }
            m.SetRetentionEvents(res)
        }
        return nil
    }
    return res
}
// GetRetentionEvents gets the retentionEvents property value. The retentionEvents property
// returns a []RetentionEventable when successful
func (m *TriggersRoot) GetRetentionEvents()([]RetentionEventable) {
    val, err := m.GetBackingStore().Get("retentionEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RetentionEventable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TriggersRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetRetentionEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRetentionEvents()))
        for i, v := range m.GetRetentionEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("retentionEvents", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRetentionEvents sets the retentionEvents property value. The retentionEvents property
func (m *TriggersRoot) SetRetentionEvents(value []RetentionEventable)() {
    err := m.GetBackingStore().Set("retentionEvents", value)
    if err != nil {
        panic(err)
    }
}
type TriggersRootable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRetentionEvents()([]RetentionEventable)
    SetRetentionEvents(value []RetentionEventable)()
}
