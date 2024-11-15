package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AttackSimulationSimulationUserCoverage struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAttackSimulationSimulationUserCoverage instantiates a new AttackSimulationSimulationUserCoverage and sets the default values.
func NewAttackSimulationSimulationUserCoverage()(*AttackSimulationSimulationUserCoverage) {
    m := &AttackSimulationSimulationUserCoverage{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAttackSimulationSimulationUserCoverageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttackSimulationSimulationUserCoverageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttackSimulationSimulationUserCoverage(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AttackSimulationSimulationUserCoverage) GetAdditionalData()(map[string]any) {
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
// GetAttackSimulationUser gets the attackSimulationUser property value. User in an attack simulation and training campaign.
// returns a AttackSimulationUserable when successful
func (m *AttackSimulationSimulationUserCoverage) GetAttackSimulationUser()(AttackSimulationUserable) {
    val, err := m.GetBackingStore().Get("attackSimulationUser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AttackSimulationUserable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *AttackSimulationSimulationUserCoverage) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetClickCount gets the clickCount property value. Number of link clicks in the received payloads by the user in attack simulation and training campaigns.
// returns a *int32 when successful
func (m *AttackSimulationSimulationUserCoverage) GetClickCount()(*int32) {
    val, err := m.GetBackingStore().Get("clickCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCompromisedCount gets the compromisedCount property value. Number of compromising actions by the user in attack simulation and training campaigns.
// returns a *int32 when successful
func (m *AttackSimulationSimulationUserCoverage) GetCompromisedCount()(*int32) {
    val, err := m.GetBackingStore().Get("compromisedCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AttackSimulationSimulationUserCoverage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attackSimulationUser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAttackSimulationUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttackSimulationUser(val.(AttackSimulationUserable))
        }
        return nil
    }
    res["clickCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClickCount(val)
        }
        return nil
    }
    res["compromisedCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompromisedCount(val)
        }
        return nil
    }
    res["latestSimulationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLatestSimulationDateTime(val)
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
    res["simulationCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSimulationCount(val)
        }
        return nil
    }
    return res
}
// GetLatestSimulationDateTime gets the latestSimulationDateTime property value. Date and time of the latest attack simulation and training campaign that the user was included in.
// returns a *Time when successful
func (m *AttackSimulationSimulationUserCoverage) GetLatestSimulationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("latestSimulationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AttackSimulationSimulationUserCoverage) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSimulationCount gets the simulationCount property value. Number of attack simulation and training campaigns that the user was included in.
// returns a *int32 when successful
func (m *AttackSimulationSimulationUserCoverage) GetSimulationCount()(*int32) {
    val, err := m.GetBackingStore().Get("simulationCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AttackSimulationSimulationUserCoverage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("attackSimulationUser", m.GetAttackSimulationUser())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("clickCount", m.GetClickCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("compromisedCount", m.GetCompromisedCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("latestSimulationDateTime", m.GetLatestSimulationDateTime())
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
        err := writer.WriteInt32Value("simulationCount", m.GetSimulationCount())
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
func (m *AttackSimulationSimulationUserCoverage) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAttackSimulationUser sets the attackSimulationUser property value. User in an attack simulation and training campaign.
func (m *AttackSimulationSimulationUserCoverage) SetAttackSimulationUser(value AttackSimulationUserable)() {
    err := m.GetBackingStore().Set("attackSimulationUser", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AttackSimulationSimulationUserCoverage) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetClickCount sets the clickCount property value. Number of link clicks in the received payloads by the user in attack simulation and training campaigns.
func (m *AttackSimulationSimulationUserCoverage) SetClickCount(value *int32)() {
    err := m.GetBackingStore().Set("clickCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompromisedCount sets the compromisedCount property value. Number of compromising actions by the user in attack simulation and training campaigns.
func (m *AttackSimulationSimulationUserCoverage) SetCompromisedCount(value *int32)() {
    err := m.GetBackingStore().Set("compromisedCount", value)
    if err != nil {
        panic(err)
    }
}
// SetLatestSimulationDateTime sets the latestSimulationDateTime property value. Date and time of the latest attack simulation and training campaign that the user was included in.
func (m *AttackSimulationSimulationUserCoverage) SetLatestSimulationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("latestSimulationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AttackSimulationSimulationUserCoverage) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSimulationCount sets the simulationCount property value. Number of attack simulation and training campaigns that the user was included in.
func (m *AttackSimulationSimulationUserCoverage) SetSimulationCount(value *int32)() {
    err := m.GetBackingStore().Set("simulationCount", value)
    if err != nil {
        panic(err)
    }
}
type AttackSimulationSimulationUserCoverageable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttackSimulationUser()(AttackSimulationUserable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetClickCount()(*int32)
    GetCompromisedCount()(*int32)
    GetLatestSimulationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetOdataType()(*string)
    GetSimulationCount()(*int32)
    SetAttackSimulationUser(value AttackSimulationUserable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetClickCount(value *int32)()
    SetCompromisedCount(value *int32)()
    SetLatestSimulationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetOdataType(value *string)()
    SetSimulationCount(value *int32)()
}
