package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SimulationReportOverview struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSimulationReportOverview instantiates a new SimulationReportOverview and sets the default values.
func NewSimulationReportOverview()(*SimulationReportOverview) {
    m := &SimulationReportOverview{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSimulationReportOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSimulationReportOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSimulationReportOverview(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SimulationReportOverview) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *SimulationReportOverview) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SimulationReportOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["recommendedActions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRecommendedActionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RecommendedActionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RecommendedActionable)
                }
            }
            m.SetRecommendedActions(res)
        }
        return nil
    }
    res["resolvedTargetsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResolvedTargetsCount(val)
        }
        return nil
    }
    res["simulationEventsContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSimulationEventsContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSimulationEventsContent(val.(SimulationEventsContentable))
        }
        return nil
    }
    res["trainingEventsContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTrainingEventsContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrainingEventsContent(val.(TrainingEventsContentable))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SimulationReportOverview) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecommendedActions gets the recommendedActions property value. List of recommended actions for a tenant to improve its security posture based on the attack simulation and training campaign attack type.
// returns a []RecommendedActionable when successful
func (m *SimulationReportOverview) GetRecommendedActions()([]RecommendedActionable) {
    val, err := m.GetBackingStore().Get("recommendedActions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RecommendedActionable)
    }
    return nil
}
// GetResolvedTargetsCount gets the resolvedTargetsCount property value. Number of valid users in the attack simulation and training campaign.
// returns a *int32 when successful
func (m *SimulationReportOverview) GetResolvedTargetsCount()(*int32) {
    val, err := m.GetBackingStore().Get("resolvedTargetsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSimulationEventsContent gets the simulationEventsContent property value. Summary of simulation events in the attack simulation and training campaign.
// returns a SimulationEventsContentable when successful
func (m *SimulationReportOverview) GetSimulationEventsContent()(SimulationEventsContentable) {
    val, err := m.GetBackingStore().Get("simulationEventsContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SimulationEventsContentable)
    }
    return nil
}
// GetTrainingEventsContent gets the trainingEventsContent property value. Summary of assigned trainings in the attack simulation and training campaign.
// returns a TrainingEventsContentable when successful
func (m *SimulationReportOverview) GetTrainingEventsContent()(TrainingEventsContentable) {
    val, err := m.GetBackingStore().Get("trainingEventsContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TrainingEventsContentable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SimulationReportOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRecommendedActions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRecommendedActions()))
        for i, v := range m.GetRecommendedActions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("recommendedActions", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("resolvedTargetsCount", m.GetResolvedTargetsCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("simulationEventsContent", m.GetSimulationEventsContent())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("trainingEventsContent", m.GetTrainingEventsContent())
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
func (m *SimulationReportOverview) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SimulationReportOverview) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SimulationReportOverview) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendedActions sets the recommendedActions property value. List of recommended actions for a tenant to improve its security posture based on the attack simulation and training campaign attack type.
func (m *SimulationReportOverview) SetRecommendedActions(value []RecommendedActionable)() {
    err := m.GetBackingStore().Set("recommendedActions", value)
    if err != nil {
        panic(err)
    }
}
// SetResolvedTargetsCount sets the resolvedTargetsCount property value. Number of valid users in the attack simulation and training campaign.
func (m *SimulationReportOverview) SetResolvedTargetsCount(value *int32)() {
    err := m.GetBackingStore().Set("resolvedTargetsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSimulationEventsContent sets the simulationEventsContent property value. Summary of simulation events in the attack simulation and training campaign.
func (m *SimulationReportOverview) SetSimulationEventsContent(value SimulationEventsContentable)() {
    err := m.GetBackingStore().Set("simulationEventsContent", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainingEventsContent sets the trainingEventsContent property value. Summary of assigned trainings in the attack simulation and training campaign.
func (m *SimulationReportOverview) SetTrainingEventsContent(value TrainingEventsContentable)() {
    err := m.GetBackingStore().Set("trainingEventsContent", value)
    if err != nil {
        panic(err)
    }
}
type SimulationReportOverviewable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetRecommendedActions()([]RecommendedActionable)
    GetResolvedTargetsCount()(*int32)
    GetSimulationEventsContent()(SimulationEventsContentable)
    GetTrainingEventsContent()(TrainingEventsContentable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetRecommendedActions(value []RecommendedActionable)()
    SetResolvedTargetsCount(value *int32)()
    SetSimulationEventsContent(value SimulationEventsContentable)()
    SetTrainingEventsContent(value TrainingEventsContentable)()
}
