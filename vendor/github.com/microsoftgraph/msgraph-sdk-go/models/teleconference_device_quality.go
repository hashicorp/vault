package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type TeleconferenceDeviceQuality struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewTeleconferenceDeviceQuality instantiates a new TeleconferenceDeviceQuality and sets the default values.
func NewTeleconferenceDeviceQuality()(*TeleconferenceDeviceQuality) {
    m := &TeleconferenceDeviceQuality{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateTeleconferenceDeviceQualityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeleconferenceDeviceQualityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeleconferenceDeviceQuality(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *TeleconferenceDeviceQuality) GetAdditionalData()(map[string]any) {
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
func (m *TeleconferenceDeviceQuality) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCallChainId gets the callChainId property value. A unique identifier for all  the participant calls in a conference or a unique identifier for two participant calls in P2P call. This needs to be copied over from Microsoft.Graph.Call.CallChainId.
// returns a *UUID when successful
func (m *TeleconferenceDeviceQuality) GetCallChainId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("callChainId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetCloudServiceDeploymentEnvironment gets the cloudServiceDeploymentEnvironment property value. A geo-region where the service is deployed, such as ProdNoam.
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetCloudServiceDeploymentEnvironment()(*string) {
    val, err := m.GetBackingStore().Get("cloudServiceDeploymentEnvironment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCloudServiceDeploymentId gets the cloudServiceDeploymentId property value. A unique deployment identifier assigned by Azure.
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetCloudServiceDeploymentId()(*string) {
    val, err := m.GetBackingStore().Get("cloudServiceDeploymentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCloudServiceInstanceName gets the cloudServiceInstanceName property value. The Azure deployed cloud service instance name, such as FrontEndIN3.
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetCloudServiceInstanceName()(*string) {
    val, err := m.GetBackingStore().Get("cloudServiceInstanceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCloudServiceName gets the cloudServiceName property value. The Azure deployed cloud service name, such as contoso.cloudapp.net.
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetCloudServiceName()(*string) {
    val, err := m.GetBackingStore().Get("cloudServiceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceDescription gets the deviceDescription property value. Any additional description, such as VTC Bldg 30/21.
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetDeviceDescription()(*string) {
    val, err := m.GetBackingStore().Get("deviceDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. The user media agent name, such as Cisco SX80.
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeleconferenceDeviceQuality) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["callChainId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallChainId(val)
        }
        return nil
    }
    res["cloudServiceDeploymentEnvironment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudServiceDeploymentEnvironment(val)
        }
        return nil
    }
    res["cloudServiceDeploymentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudServiceDeploymentId(val)
        }
        return nil
    }
    res["cloudServiceInstanceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudServiceInstanceName(val)
        }
        return nil
    }
    res["cloudServiceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudServiceName(val)
        }
        return nil
    }
    res["deviceDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceDescription(val)
        }
        return nil
    }
    res["deviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceName(val)
        }
        return nil
    }
    res["mediaLegId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaLegId(val)
        }
        return nil
    }
    res["mediaQualityList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeleconferenceDeviceMediaQualityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeleconferenceDeviceMediaQualityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeleconferenceDeviceMediaQualityable)
                }
            }
            m.SetMediaQualityList(res)
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
    res["participantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParticipantId(val)
        }
        return nil
    }
    return res
}
// GetMediaLegId gets the mediaLegId property value. A unique identifier for a specific media leg of a participant in a conference.  One participant can have multiple media leg identifiers if retargeting happens. CVI partner assigns this value.
// returns a *UUID when successful
func (m *TeleconferenceDeviceQuality) GetMediaLegId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("mediaLegId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetMediaQualityList gets the mediaQualityList property value. The list of media qualities in a media session (call), such as audio quality, video quality, and/or screen sharing quality.
// returns a []TeleconferenceDeviceMediaQualityable when successful
func (m *TeleconferenceDeviceQuality) GetMediaQualityList()([]TeleconferenceDeviceMediaQualityable) {
    val, err := m.GetBackingStore().Get("mediaQualityList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeleconferenceDeviceMediaQualityable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *TeleconferenceDeviceQuality) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetParticipantId gets the participantId property value. A unique identifier for a specific participant in a conference. The CVI partner needs to copy over Call.MyParticipantId to this property.
// returns a *UUID when successful
func (m *TeleconferenceDeviceQuality) GetParticipantId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("participantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeleconferenceDeviceQuality) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteUUIDValue("callChainId", m.GetCallChainId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("cloudServiceDeploymentEnvironment", m.GetCloudServiceDeploymentEnvironment())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("cloudServiceDeploymentId", m.GetCloudServiceDeploymentId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("cloudServiceInstanceName", m.GetCloudServiceInstanceName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("cloudServiceName", m.GetCloudServiceName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deviceDescription", m.GetDeviceDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deviceName", m.GetDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteUUIDValue("mediaLegId", m.GetMediaLegId())
        if err != nil {
            return err
        }
    }
    if m.GetMediaQualityList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMediaQualityList()))
        for i, v := range m.GetMediaQualityList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("mediaQualityList", cast)
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
        err := writer.WriteUUIDValue("participantId", m.GetParticipantId())
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
func (m *TeleconferenceDeviceQuality) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *TeleconferenceDeviceQuality) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCallChainId sets the callChainId property value. A unique identifier for all  the participant calls in a conference or a unique identifier for two participant calls in P2P call. This needs to be copied over from Microsoft.Graph.Call.CallChainId.
func (m *TeleconferenceDeviceQuality) SetCallChainId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("callChainId", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudServiceDeploymentEnvironment sets the cloudServiceDeploymentEnvironment property value. A geo-region where the service is deployed, such as ProdNoam.
func (m *TeleconferenceDeviceQuality) SetCloudServiceDeploymentEnvironment(value *string)() {
    err := m.GetBackingStore().Set("cloudServiceDeploymentEnvironment", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudServiceDeploymentId sets the cloudServiceDeploymentId property value. A unique deployment identifier assigned by Azure.
func (m *TeleconferenceDeviceQuality) SetCloudServiceDeploymentId(value *string)() {
    err := m.GetBackingStore().Set("cloudServiceDeploymentId", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudServiceInstanceName sets the cloudServiceInstanceName property value. The Azure deployed cloud service instance name, such as FrontEndIN3.
func (m *TeleconferenceDeviceQuality) SetCloudServiceInstanceName(value *string)() {
    err := m.GetBackingStore().Set("cloudServiceInstanceName", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudServiceName sets the cloudServiceName property value. The Azure deployed cloud service name, such as contoso.cloudapp.net.
func (m *TeleconferenceDeviceQuality) SetCloudServiceName(value *string)() {
    err := m.GetBackingStore().Set("cloudServiceName", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceDescription sets the deviceDescription property value. Any additional description, such as VTC Bldg 30/21.
func (m *TeleconferenceDeviceQuality) SetDeviceDescription(value *string)() {
    err := m.GetBackingStore().Set("deviceDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. The user media agent name, such as Cisco SX80.
func (m *TeleconferenceDeviceQuality) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaLegId sets the mediaLegId property value. A unique identifier for a specific media leg of a participant in a conference.  One participant can have multiple media leg identifiers if retargeting happens. CVI partner assigns this value.
func (m *TeleconferenceDeviceQuality) SetMediaLegId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("mediaLegId", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaQualityList sets the mediaQualityList property value. The list of media qualities in a media session (call), such as audio quality, video quality, and/or screen sharing quality.
func (m *TeleconferenceDeviceQuality) SetMediaQualityList(value []TeleconferenceDeviceMediaQualityable)() {
    err := m.GetBackingStore().Set("mediaQualityList", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TeleconferenceDeviceQuality) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetParticipantId sets the participantId property value. A unique identifier for a specific participant in a conference. The CVI partner needs to copy over Call.MyParticipantId to this property.
func (m *TeleconferenceDeviceQuality) SetParticipantId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("participantId", value)
    if err != nil {
        panic(err)
    }
}
type TeleconferenceDeviceQualityable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCallChainId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetCloudServiceDeploymentEnvironment()(*string)
    GetCloudServiceDeploymentId()(*string)
    GetCloudServiceInstanceName()(*string)
    GetCloudServiceName()(*string)
    GetDeviceDescription()(*string)
    GetDeviceName()(*string)
    GetMediaLegId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetMediaQualityList()([]TeleconferenceDeviceMediaQualityable)
    GetOdataType()(*string)
    GetParticipantId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCallChainId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetCloudServiceDeploymentEnvironment(value *string)()
    SetCloudServiceDeploymentId(value *string)()
    SetCloudServiceInstanceName(value *string)()
    SetCloudServiceName(value *string)()
    SetDeviceDescription(value *string)()
    SetDeviceName(value *string)()
    SetMediaLegId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetMediaQualityList(value []TeleconferenceDeviceMediaQualityable)()
    SetOdataType(value *string)()
    SetParticipantId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
}
