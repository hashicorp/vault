package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AttackSimulationRoot struct {
    Entity
}
// NewAttackSimulationRoot instantiates a new AttackSimulationRoot and sets the default values.
func NewAttackSimulationRoot()(*AttackSimulationRoot) {
    m := &AttackSimulationRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAttackSimulationRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttackSimulationRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttackSimulationRoot(), nil
}
// GetEndUserNotifications gets the endUserNotifications property value. Represents an end user's notification for an attack simulation training.
// returns a []EndUserNotificationable when successful
func (m *AttackSimulationRoot) GetEndUserNotifications()([]EndUserNotificationable) {
    val, err := m.GetBackingStore().Get("endUserNotifications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EndUserNotificationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AttackSimulationRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["endUserNotifications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEndUserNotificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EndUserNotificationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EndUserNotificationable)
                }
            }
            m.SetEndUserNotifications(res)
        }
        return nil
    }
    res["landingPages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLandingPageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LandingPageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LandingPageable)
                }
            }
            m.SetLandingPages(res)
        }
        return nil
    }
    res["loginPages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLoginPageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LoginPageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LoginPageable)
                }
            }
            m.SetLoginPages(res)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttackSimulationOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AttackSimulationOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AttackSimulationOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["payloads"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePayloadFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Payloadable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Payloadable)
                }
            }
            m.SetPayloads(res)
        }
        return nil
    }
    res["simulationAutomations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSimulationAutomationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SimulationAutomationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SimulationAutomationable)
                }
            }
            m.SetSimulationAutomations(res)
        }
        return nil
    }
    res["simulations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSimulationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Simulationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Simulationable)
                }
            }
            m.SetSimulations(res)
        }
        return nil
    }
    res["trainings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTrainingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Trainingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Trainingable)
                }
            }
            m.SetTrainings(res)
        }
        return nil
    }
    return res
}
// GetLandingPages gets the landingPages property value. Represents an attack simulation training landing page.
// returns a []LandingPageable when successful
func (m *AttackSimulationRoot) GetLandingPages()([]LandingPageable) {
    val, err := m.GetBackingStore().Get("landingPages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LandingPageable)
    }
    return nil
}
// GetLoginPages gets the loginPages property value. Represents an attack simulation training login page.
// returns a []LoginPageable when successful
func (m *AttackSimulationRoot) GetLoginPages()([]LoginPageable) {
    val, err := m.GetBackingStore().Get("loginPages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LoginPageable)
    }
    return nil
}
// GetOperations gets the operations property value. Represents an attack simulation training operation.
// returns a []AttackSimulationOperationable when successful
func (m *AttackSimulationRoot) GetOperations()([]AttackSimulationOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AttackSimulationOperationable)
    }
    return nil
}
// GetPayloads gets the payloads property value. Represents an attack simulation training campaign payload in a tenant.
// returns a []Payloadable when successful
func (m *AttackSimulationRoot) GetPayloads()([]Payloadable) {
    val, err := m.GetBackingStore().Get("payloads")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Payloadable)
    }
    return nil
}
// GetSimulationAutomations gets the simulationAutomations property value. Represents simulation automation created to run on a tenant.
// returns a []SimulationAutomationable when successful
func (m *AttackSimulationRoot) GetSimulationAutomations()([]SimulationAutomationable) {
    val, err := m.GetBackingStore().Get("simulationAutomations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SimulationAutomationable)
    }
    return nil
}
// GetSimulations gets the simulations property value. Represents an attack simulation training campaign in a tenant.
// returns a []Simulationable when successful
func (m *AttackSimulationRoot) GetSimulations()([]Simulationable) {
    val, err := m.GetBackingStore().Get("simulations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Simulationable)
    }
    return nil
}
// GetTrainings gets the trainings property value. Represents details about attack simulation trainings.
// returns a []Trainingable when successful
func (m *AttackSimulationRoot) GetTrainings()([]Trainingable) {
    val, err := m.GetBackingStore().Get("trainings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Trainingable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AttackSimulationRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetEndUserNotifications() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEndUserNotifications()))
        for i, v := range m.GetEndUserNotifications() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("endUserNotifications", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLandingPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLandingPages()))
        for i, v := range m.GetLandingPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("landingPages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLoginPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLoginPages()))
        for i, v := range m.GetLoginPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("loginPages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPayloads() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPayloads()))
        for i, v := range m.GetPayloads() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("payloads", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSimulationAutomations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSimulationAutomations()))
        for i, v := range m.GetSimulationAutomations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("simulationAutomations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSimulations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSimulations()))
        for i, v := range m.GetSimulations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("simulations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTrainings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTrainings()))
        for i, v := range m.GetTrainings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("trainings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEndUserNotifications sets the endUserNotifications property value. Represents an end user's notification for an attack simulation training.
func (m *AttackSimulationRoot) SetEndUserNotifications(value []EndUserNotificationable)() {
    err := m.GetBackingStore().Set("endUserNotifications", value)
    if err != nil {
        panic(err)
    }
}
// SetLandingPages sets the landingPages property value. Represents an attack simulation training landing page.
func (m *AttackSimulationRoot) SetLandingPages(value []LandingPageable)() {
    err := m.GetBackingStore().Set("landingPages", value)
    if err != nil {
        panic(err)
    }
}
// SetLoginPages sets the loginPages property value. Represents an attack simulation training login page.
func (m *AttackSimulationRoot) SetLoginPages(value []LoginPageable)() {
    err := m.GetBackingStore().Set("loginPages", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. Represents an attack simulation training operation.
func (m *AttackSimulationRoot) SetOperations(value []AttackSimulationOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPayloads sets the payloads property value. Represents an attack simulation training campaign payload in a tenant.
func (m *AttackSimulationRoot) SetPayloads(value []Payloadable)() {
    err := m.GetBackingStore().Set("payloads", value)
    if err != nil {
        panic(err)
    }
}
// SetSimulationAutomations sets the simulationAutomations property value. Represents simulation automation created to run on a tenant.
func (m *AttackSimulationRoot) SetSimulationAutomations(value []SimulationAutomationable)() {
    err := m.GetBackingStore().Set("simulationAutomations", value)
    if err != nil {
        panic(err)
    }
}
// SetSimulations sets the simulations property value. Represents an attack simulation training campaign in a tenant.
func (m *AttackSimulationRoot) SetSimulations(value []Simulationable)() {
    err := m.GetBackingStore().Set("simulations", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainings sets the trainings property value. Represents details about attack simulation trainings.
func (m *AttackSimulationRoot) SetTrainings(value []Trainingable)() {
    err := m.GetBackingStore().Set("trainings", value)
    if err != nil {
        panic(err)
    }
}
type AttackSimulationRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEndUserNotifications()([]EndUserNotificationable)
    GetLandingPages()([]LandingPageable)
    GetLoginPages()([]LoginPageable)
    GetOperations()([]AttackSimulationOperationable)
    GetPayloads()([]Payloadable)
    GetSimulationAutomations()([]SimulationAutomationable)
    GetSimulations()([]Simulationable)
    GetTrainings()([]Trainingable)
    SetEndUserNotifications(value []EndUserNotificationable)()
    SetLandingPages(value []LandingPageable)()
    SetLoginPages(value []LoginPageable)()
    SetOperations(value []AttackSimulationOperationable)()
    SetPayloads(value []Payloadable)()
    SetSimulationAutomations(value []SimulationAutomationable)()
    SetSimulations(value []Simulationable)()
    SetTrainings(value []Trainingable)()
}
