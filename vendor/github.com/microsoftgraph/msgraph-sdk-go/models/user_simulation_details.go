package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type UserSimulationDetails struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUserSimulationDetails instantiates a new UserSimulationDetails and sets the default values.
func NewUserSimulationDetails()(*UserSimulationDetails) {
    m := &UserSimulationDetails{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUserSimulationDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserSimulationDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserSimulationDetails(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UserSimulationDetails) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAssignedTrainingsCount gets the assignedTrainingsCount property value. Number of trainings assigned to a user in an attack simulation and training campaign.
// returns a *int32 when successful
func (m *UserSimulationDetails) GetAssignedTrainingsCount()(*int32) {
    val, err := m.GetBackingStore().Get("assignedTrainingsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *UserSimulationDetails) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCompletedTrainingsCount gets the completedTrainingsCount property value. Number of trainings completed by a user in an attack simulation and training campaign.
// returns a *int32 when successful
func (m *UserSimulationDetails) GetCompletedTrainingsCount()(*int32) {
    val, err := m.GetBackingStore().Get("completedTrainingsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCompromisedDateTime gets the compromisedDateTime property value. Date and time of the compromising online action by a user in an attack simulation and training campaign.
// returns a *Time when successful
func (m *UserSimulationDetails) GetCompromisedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("compromisedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserSimulationDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["assignedTrainingsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignedTrainingsCount(val)
        }
        return nil
    }
    res["completedTrainingsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedTrainingsCount(val)
        }
        return nil
    }
    res["compromisedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompromisedDateTime(val)
        }
        return nil
    }
    res["inProgressTrainingsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInProgressTrainingsCount(val)
        }
        return nil
    }
    res["isCompromised"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsCompromised(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["reportedPhishDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReportedPhishDateTime(val)
        }
        return nil
    }
    res["simulationEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserSimulationEventInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserSimulationEventInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserSimulationEventInfoable)
                }
            }
            m.SetSimulationEvents(res)
        }
        return nil
    }
    res["simulationUser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAttackSimulationUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSimulationUser(val.(AttackSimulationUserable))
        }
        return nil
    }
    res["trainingEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserTrainingEventInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserTrainingEventInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserTrainingEventInfoable)
                }
            }
            m.SetTrainingEvents(res)
        }
        return nil
    }
    return res
}
// GetInProgressTrainingsCount gets the inProgressTrainingsCount property value. Number of trainings in progress by a user in an attack simulation and training campaign.
// returns a *int32 when successful
func (m *UserSimulationDetails) GetInProgressTrainingsCount()(*int32) {
    val, err := m.GetBackingStore().Get("inProgressTrainingsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetIsCompromised gets the isCompromised property value. Indicates whether a user was compromised in an attack simulation and training campaign.
// returns a *bool when successful
func (m *UserSimulationDetails) GetIsCompromised()(*bool) {
    val, err := m.GetBackingStore().Get("isCompromised")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UserSimulationDetails) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReportedPhishDateTime gets the reportedPhishDateTime property value. Date and time when a user reported the delivered payload as phishing in the attack simulation and training campaign.
// returns a *Time when successful
func (m *UserSimulationDetails) GetReportedPhishDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("reportedPhishDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSimulationEvents gets the simulationEvents property value. List of simulation events of a user in the attack simulation and training campaign.
// returns a []UserSimulationEventInfoable when successful
func (m *UserSimulationDetails) GetSimulationEvents()([]UserSimulationEventInfoable) {
    val, err := m.GetBackingStore().Get("simulationEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserSimulationEventInfoable)
    }
    return nil
}
// GetSimulationUser gets the simulationUser property value. User in an attack simulation and training campaign.
// returns a AttackSimulationUserable when successful
func (m *UserSimulationDetails) GetSimulationUser()(AttackSimulationUserable) {
    val, err := m.GetBackingStore().Get("simulationUser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AttackSimulationUserable)
    }
    return nil
}
// GetTrainingEvents gets the trainingEvents property value. List of training events of a user in the attack simulation and training campaign.
// returns a []UserTrainingEventInfoable when successful
func (m *UserSimulationDetails) GetTrainingEvents()([]UserTrainingEventInfoable) {
    val, err := m.GetBackingStore().Get("trainingEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserTrainingEventInfoable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserSimulationDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("assignedTrainingsCount", m.GetAssignedTrainingsCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("completedTrainingsCount", m.GetCompletedTrainingsCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("compromisedDateTime", m.GetCompromisedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("inProgressTrainingsCount", m.GetInProgressTrainingsCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isCompromised", m.GetIsCompromised())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("reportedPhishDateTime", m.GetReportedPhishDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetSimulationEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSimulationEvents()))
        for i, v := range m.GetSimulationEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("simulationEvents", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("simulationUser", m.GetSimulationUser())
        if err != nil {
            return err
        }
    }
    if m.GetTrainingEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTrainingEvents()))
        for i, v := range m.GetTrainingEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("trainingEvents", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserSimulationDetails) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedTrainingsCount sets the assignedTrainingsCount property value. Number of trainings assigned to a user in an attack simulation and training campaign.
func (m *UserSimulationDetails) SetAssignedTrainingsCount(value *int32)() {
    err := m.GetBackingStore().Set("assignedTrainingsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UserSimulationDetails) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCompletedTrainingsCount sets the completedTrainingsCount property value. Number of trainings completed by a user in an attack simulation and training campaign.
func (m *UserSimulationDetails) SetCompletedTrainingsCount(value *int32)() {
    err := m.GetBackingStore().Set("completedTrainingsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompromisedDateTime sets the compromisedDateTime property value. Date and time of the compromising online action by a user in an attack simulation and training campaign.
func (m *UserSimulationDetails) SetCompromisedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("compromisedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetInProgressTrainingsCount sets the inProgressTrainingsCount property value. Number of trainings in progress by a user in an attack simulation and training campaign.
func (m *UserSimulationDetails) SetInProgressTrainingsCount(value *int32)() {
    err := m.GetBackingStore().Set("inProgressTrainingsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetIsCompromised sets the isCompromised property value. Indicates whether a user was compromised in an attack simulation and training campaign.
func (m *UserSimulationDetails) SetIsCompromised(value *bool)() {
    err := m.GetBackingStore().Set("isCompromised", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserSimulationDetails) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetReportedPhishDateTime sets the reportedPhishDateTime property value. Date and time when a user reported the delivered payload as phishing in the attack simulation and training campaign.
func (m *UserSimulationDetails) SetReportedPhishDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("reportedPhishDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSimulationEvents sets the simulationEvents property value. List of simulation events of a user in the attack simulation and training campaign.
func (m *UserSimulationDetails) SetSimulationEvents(value []UserSimulationEventInfoable)() {
    err := m.GetBackingStore().Set("simulationEvents", value)
    if err != nil {
        panic(err)
    }
}
// SetSimulationUser sets the simulationUser property value. User in an attack simulation and training campaign.
func (m *UserSimulationDetails) SetSimulationUser(value AttackSimulationUserable)() {
    err := m.GetBackingStore().Set("simulationUser", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainingEvents sets the trainingEvents property value. List of training events of a user in the attack simulation and training campaign.
func (m *UserSimulationDetails) SetTrainingEvents(value []UserTrainingEventInfoable)() {
    err := m.GetBackingStore().Set("trainingEvents", value)
    if err != nil {
        panic(err)
    }
}
type UserSimulationDetailsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignedTrainingsCount()(*int32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCompletedTrainingsCount()(*int32)
    GetCompromisedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetInProgressTrainingsCount()(*int32)
    GetIsCompromised()(*bool)
    GetOdataType()(*string)
    GetReportedPhishDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSimulationEvents()([]UserSimulationEventInfoable)
    GetSimulationUser()(AttackSimulationUserable)
    GetTrainingEvents()([]UserTrainingEventInfoable)
    SetAssignedTrainingsCount(value *int32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCompletedTrainingsCount(value *int32)()
    SetCompromisedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetInProgressTrainingsCount(value *int32)()
    SetIsCompromised(value *bool)()
    SetOdataType(value *string)()
    SetReportedPhishDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSimulationEvents(value []UserSimulationEventInfoable)()
    SetSimulationUser(value AttackSimulationUserable)()
    SetTrainingEvents(value []UserTrainingEventInfoable)()
}
