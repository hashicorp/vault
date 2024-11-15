package callrecords

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Media struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMedia instantiates a new Media and sets the default values.
func NewMedia()(*Media) {
    m := &Media{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMediaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMediaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMedia(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Media) GetAdditionalData()(map[string]any) {
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
func (m *Media) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCalleeDevice gets the calleeDevice property value. Device information associated with the callee endpoint of this media.
// returns a DeviceInfoable when successful
func (m *Media) GetCalleeDevice()(DeviceInfoable) {
    val, err := m.GetBackingStore().Get("calleeDevice")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceInfoable)
    }
    return nil
}
// GetCalleeNetwork gets the calleeNetwork property value. Network information associated with the callee endpoint of this media.
// returns a NetworkInfoable when successful
func (m *Media) GetCalleeNetwork()(NetworkInfoable) {
    val, err := m.GetBackingStore().Get("calleeNetwork")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(NetworkInfoable)
    }
    return nil
}
// GetCallerDevice gets the callerDevice property value. Device information associated with the caller endpoint of this media.
// returns a DeviceInfoable when successful
func (m *Media) GetCallerDevice()(DeviceInfoable) {
    val, err := m.GetBackingStore().Get("callerDevice")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceInfoable)
    }
    return nil
}
// GetCallerNetwork gets the callerNetwork property value. Network information associated with the caller endpoint of this media.
// returns a NetworkInfoable when successful
func (m *Media) GetCallerNetwork()(NetworkInfoable) {
    val, err := m.GetBackingStore().Get("callerNetwork")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(NetworkInfoable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Media) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["calleeDevice"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCalleeDevice(val.(DeviceInfoable))
        }
        return nil
    }
    res["calleeNetwork"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateNetworkInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCalleeNetwork(val.(NetworkInfoable))
        }
        return nil
    }
    res["callerDevice"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallerDevice(val.(DeviceInfoable))
        }
        return nil
    }
    res["callerNetwork"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateNetworkInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallerNetwork(val.(NetworkInfoable))
        }
        return nil
    }
    res["label"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabel(val)
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
    res["streams"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMediaStreamFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MediaStreamable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MediaStreamable)
                }
            }
            m.SetStreams(res)
        }
        return nil
    }
    return res
}
// GetLabel gets the label property value. How the media was identified during media negotiation stage.
// returns a *string when successful
func (m *Media) GetLabel()(*string) {
    val, err := m.GetBackingStore().Get("label")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Media) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStreams gets the streams property value. Network streams associated with this media.
// returns a []MediaStreamable when successful
func (m *Media) GetStreams()([]MediaStreamable) {
    val, err := m.GetBackingStore().Get("streams")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MediaStreamable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Media) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("calleeDevice", m.GetCalleeDevice())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("calleeNetwork", m.GetCalleeNetwork())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("callerDevice", m.GetCallerDevice())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("callerNetwork", m.GetCallerNetwork())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("label", m.GetLabel())
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
    if m.GetStreams() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetStreams()))
        for i, v := range m.GetStreams() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("streams", cast)
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
func (m *Media) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Media) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCalleeDevice sets the calleeDevice property value. Device information associated with the callee endpoint of this media.
func (m *Media) SetCalleeDevice(value DeviceInfoable)() {
    err := m.GetBackingStore().Set("calleeDevice", value)
    if err != nil {
        panic(err)
    }
}
// SetCalleeNetwork sets the calleeNetwork property value. Network information associated with the callee endpoint of this media.
func (m *Media) SetCalleeNetwork(value NetworkInfoable)() {
    err := m.GetBackingStore().Set("calleeNetwork", value)
    if err != nil {
        panic(err)
    }
}
// SetCallerDevice sets the callerDevice property value. Device information associated with the caller endpoint of this media.
func (m *Media) SetCallerDevice(value DeviceInfoable)() {
    err := m.GetBackingStore().Set("callerDevice", value)
    if err != nil {
        panic(err)
    }
}
// SetCallerNetwork sets the callerNetwork property value. Network information associated with the caller endpoint of this media.
func (m *Media) SetCallerNetwork(value NetworkInfoable)() {
    err := m.GetBackingStore().Set("callerNetwork", value)
    if err != nil {
        panic(err)
    }
}
// SetLabel sets the label property value. How the media was identified during media negotiation stage.
func (m *Media) SetLabel(value *string)() {
    err := m.GetBackingStore().Set("label", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Media) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetStreams sets the streams property value. Network streams associated with this media.
func (m *Media) SetStreams(value []MediaStreamable)() {
    err := m.GetBackingStore().Set("streams", value)
    if err != nil {
        panic(err)
    }
}
type Mediaable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCalleeDevice()(DeviceInfoable)
    GetCalleeNetwork()(NetworkInfoable)
    GetCallerDevice()(DeviceInfoable)
    GetCallerNetwork()(NetworkInfoable)
    GetLabel()(*string)
    GetOdataType()(*string)
    GetStreams()([]MediaStreamable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCalleeDevice(value DeviceInfoable)()
    SetCalleeNetwork(value NetworkInfoable)()
    SetCallerDevice(value DeviceInfoable)()
    SetCallerNetwork(value NetworkInfoable)()
    SetLabel(value *string)()
    SetOdataType(value *string)()
    SetStreams(value []MediaStreamable)()
}
