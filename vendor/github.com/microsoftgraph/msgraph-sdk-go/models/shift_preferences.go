package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ShiftPreferences struct {
    ChangeTrackedEntity
}
// NewShiftPreferences instantiates a new ShiftPreferences and sets the default values.
func NewShiftPreferences()(*ShiftPreferences) {
    m := &ShiftPreferences{
        ChangeTrackedEntity: *NewChangeTrackedEntity(),
    }
    odataTypeValue := "#microsoft.graph.shiftPreferences"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateShiftPreferencesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateShiftPreferencesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewShiftPreferences(), nil
}
// GetAvailability gets the availability property value. Availability of the user to be scheduled for work and its recurrence pattern.
// returns a []ShiftAvailabilityable when successful
func (m *ShiftPreferences) GetAvailability()([]ShiftAvailabilityable) {
    val, err := m.GetBackingStore().Get("availability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ShiftAvailabilityable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ShiftPreferences) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ChangeTrackedEntity.GetFieldDeserializers()
    res["availability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateShiftAvailabilityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ShiftAvailabilityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ShiftAvailabilityable)
                }
            }
            m.SetAvailability(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ShiftPreferences) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ChangeTrackedEntity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAvailability() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAvailability()))
        for i, v := range m.GetAvailability() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("availability", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAvailability sets the availability property value. Availability of the user to be scheduled for work and its recurrence pattern.
func (m *ShiftPreferences) SetAvailability(value []ShiftAvailabilityable)() {
    err := m.GetBackingStore().Set("availability", value)
    if err != nil {
        panic(err)
    }
}
type ShiftPreferencesable interface {
    ChangeTrackedEntityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAvailability()([]ShiftAvailabilityable)
    SetAvailability(value []ShiftAvailabilityable)()
}
