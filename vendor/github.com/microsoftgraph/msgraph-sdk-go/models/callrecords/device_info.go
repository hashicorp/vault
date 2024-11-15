package callrecords

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type DeviceInfo struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDeviceInfo instantiates a new DeviceInfo and sets the default values.
func NewDeviceInfo()(*DeviceInfo) {
    m := &DeviceInfo{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDeviceInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceInfo(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DeviceInfo) GetAdditionalData()(map[string]any) {
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
func (m *DeviceInfo) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCaptureDeviceDriver gets the captureDeviceDriver property value. Name of the capture device driver used by the media endpoint.
// returns a *string when successful
func (m *DeviceInfo) GetCaptureDeviceDriver()(*string) {
    val, err := m.GetBackingStore().Get("captureDeviceDriver")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCaptureDeviceName gets the captureDeviceName property value. Name of the capture device used by the media endpoint.
// returns a *string when successful
func (m *DeviceInfo) GetCaptureDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("captureDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCaptureNotFunctioningEventRatio gets the captureNotFunctioningEventRatio property value. Fraction of the call that the media endpoint detected the capture device was not working properly.
// returns a *float32 when successful
func (m *DeviceInfo) GetCaptureNotFunctioningEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("captureNotFunctioningEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetCpuInsufficentEventRatio gets the cpuInsufficentEventRatio property value. Fraction of the call that the media endpoint detected the CPU resources available were insufficient and caused poor quality of the audio sent and received.
// returns a *float32 when successful
func (m *DeviceInfo) GetCpuInsufficentEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("cpuInsufficentEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetDeviceClippingEventRatio gets the deviceClippingEventRatio property value. Fraction of the call that the media endpoint detected clipping in the captured audio that caused poor quality of the audio being sent.
// returns a *float32 when successful
func (m *DeviceInfo) GetDeviceClippingEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("deviceClippingEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetDeviceGlitchEventRatio gets the deviceGlitchEventRatio property value. Fraction of the call that the media endpoint detected glitches or gaps in the audio played or captured that caused poor quality of the audio being sent or received.
// returns a *float32 when successful
func (m *DeviceInfo) GetDeviceGlitchEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("deviceGlitchEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["captureDeviceDriver"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaptureDeviceDriver(val)
        }
        return nil
    }
    res["captureDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaptureDeviceName(val)
        }
        return nil
    }
    res["captureNotFunctioningEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaptureNotFunctioningEventRatio(val)
        }
        return nil
    }
    res["cpuInsufficentEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCpuInsufficentEventRatio(val)
        }
        return nil
    }
    res["deviceClippingEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceClippingEventRatio(val)
        }
        return nil
    }
    res["deviceGlitchEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceGlitchEventRatio(val)
        }
        return nil
    }
    res["howlingEventCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHowlingEventCount(val)
        }
        return nil
    }
    res["initialSignalLevelRootMeanSquare"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitialSignalLevelRootMeanSquare(val)
        }
        return nil
    }
    res["lowSpeechLevelEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLowSpeechLevelEventRatio(val)
        }
        return nil
    }
    res["lowSpeechToNoiseEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLowSpeechToNoiseEventRatio(val)
        }
        return nil
    }
    res["micGlitchRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicGlitchRate(val)
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
    res["receivedNoiseLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReceivedNoiseLevel(val)
        }
        return nil
    }
    res["receivedSignalLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReceivedSignalLevel(val)
        }
        return nil
    }
    res["renderDeviceDriver"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRenderDeviceDriver(val)
        }
        return nil
    }
    res["renderDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRenderDeviceName(val)
        }
        return nil
    }
    res["renderMuteEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRenderMuteEventRatio(val)
        }
        return nil
    }
    res["renderNotFunctioningEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRenderNotFunctioningEventRatio(val)
        }
        return nil
    }
    res["renderZeroVolumeEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRenderZeroVolumeEventRatio(val)
        }
        return nil
    }
    res["sentNoiseLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentNoiseLevel(val)
        }
        return nil
    }
    res["sentSignalLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentSignalLevel(val)
        }
        return nil
    }
    res["speakerGlitchRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpeakerGlitchRate(val)
        }
        return nil
    }
    return res
}
// GetHowlingEventCount gets the howlingEventCount property value. Number of times during the call that the media endpoint detected howling or screeching audio.
// returns a *int32 when successful
func (m *DeviceInfo) GetHowlingEventCount()(*int32) {
    val, err := m.GetBackingStore().Get("howlingEventCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetInitialSignalLevelRootMeanSquare gets the initialSignalLevelRootMeanSquare property value. The root mean square (RMS) of the incoming signal of up to the first 30 seconds of the call.
// returns a *float32 when successful
func (m *DeviceInfo) GetInitialSignalLevelRootMeanSquare()(*float32) {
    val, err := m.GetBackingStore().Get("initialSignalLevelRootMeanSquare")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetLowSpeechLevelEventRatio gets the lowSpeechLevelEventRatio property value. Fraction of the call that the media endpoint detected low speech level that caused poor quality of the audio being sent.
// returns a *float32 when successful
func (m *DeviceInfo) GetLowSpeechLevelEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("lowSpeechLevelEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetLowSpeechToNoiseEventRatio gets the lowSpeechToNoiseEventRatio property value. Fraction of the call that the media endpoint detected low speech to noise level that caused poor quality of the audio being sent.
// returns a *float32 when successful
func (m *DeviceInfo) GetLowSpeechToNoiseEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("lowSpeechToNoiseEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetMicGlitchRate gets the micGlitchRate property value. Glitches per 5 minute interval for the media endpoint's microphone.
// returns a *float32 when successful
func (m *DeviceInfo) GetMicGlitchRate()(*float32) {
    val, err := m.GetBackingStore().Get("micGlitchRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DeviceInfo) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReceivedNoiseLevel gets the receivedNoiseLevel property value. Average energy level of received audio for audio classified as mono noise or left channel of stereo noise by the media endpoint.
// returns a *int32 when successful
func (m *DeviceInfo) GetReceivedNoiseLevel()(*int32) {
    val, err := m.GetBackingStore().Get("receivedNoiseLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetReceivedSignalLevel gets the receivedSignalLevel property value. Average energy level of received audio for audio classified as mono speech, or left channel of stereo speech by the media endpoint.
// returns a *int32 when successful
func (m *DeviceInfo) GetReceivedSignalLevel()(*int32) {
    val, err := m.GetBackingStore().Get("receivedSignalLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRenderDeviceDriver gets the renderDeviceDriver property value. Name of the render device driver used by the media endpoint.
// returns a *string when successful
func (m *DeviceInfo) GetRenderDeviceDriver()(*string) {
    val, err := m.GetBackingStore().Get("renderDeviceDriver")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRenderDeviceName gets the renderDeviceName property value. Name of the render device used by the media endpoint.
// returns a *string when successful
func (m *DeviceInfo) GetRenderDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("renderDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRenderMuteEventRatio gets the renderMuteEventRatio property value. Fraction of the call that media endpoint detected device render is muted.
// returns a *float32 when successful
func (m *DeviceInfo) GetRenderMuteEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("renderMuteEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetRenderNotFunctioningEventRatio gets the renderNotFunctioningEventRatio property value. Fraction of the call that the media endpoint detected the render device was not working properly.
// returns a *float32 when successful
func (m *DeviceInfo) GetRenderNotFunctioningEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("renderNotFunctioningEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetRenderZeroVolumeEventRatio gets the renderZeroVolumeEventRatio property value. Fraction of the call that media endpoint detected device render volume is set to 0.
// returns a *float32 when successful
func (m *DeviceInfo) GetRenderZeroVolumeEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("renderZeroVolumeEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetSentNoiseLevel gets the sentNoiseLevel property value. Average energy level of sent audio for audio classified as mono noise or left channel of stereo noise by the media endpoint.
// returns a *int32 when successful
func (m *DeviceInfo) GetSentNoiseLevel()(*int32) {
    val, err := m.GetBackingStore().Get("sentNoiseLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSentSignalLevel gets the sentSignalLevel property value. Average energy level of sent audio for audio classified as mono speech, or left channel of stereo speech by the media endpoint.
// returns a *int32 when successful
func (m *DeviceInfo) GetSentSignalLevel()(*int32) {
    val, err := m.GetBackingStore().Get("sentSignalLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSpeakerGlitchRate gets the speakerGlitchRate property value. Glitches per 5 minute internal for the media endpoint's loudspeaker.
// returns a *float32 when successful
func (m *DeviceInfo) GetSpeakerGlitchRate()(*float32) {
    val, err := m.GetBackingStore().Get("speakerGlitchRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("captureDeviceDriver", m.GetCaptureDeviceDriver())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("captureDeviceName", m.GetCaptureDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("captureNotFunctioningEventRatio", m.GetCaptureNotFunctioningEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("cpuInsufficentEventRatio", m.GetCpuInsufficentEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("deviceClippingEventRatio", m.GetDeviceClippingEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("deviceGlitchEventRatio", m.GetDeviceGlitchEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("howlingEventCount", m.GetHowlingEventCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("initialSignalLevelRootMeanSquare", m.GetInitialSignalLevelRootMeanSquare())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("lowSpeechLevelEventRatio", m.GetLowSpeechLevelEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("lowSpeechToNoiseEventRatio", m.GetLowSpeechToNoiseEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("micGlitchRate", m.GetMicGlitchRate())
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
        err := writer.WriteInt32Value("receivedNoiseLevel", m.GetReceivedNoiseLevel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("receivedSignalLevel", m.GetReceivedSignalLevel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("renderDeviceDriver", m.GetRenderDeviceDriver())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("renderDeviceName", m.GetRenderDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("renderMuteEventRatio", m.GetRenderMuteEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("renderNotFunctioningEventRatio", m.GetRenderNotFunctioningEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("renderZeroVolumeEventRatio", m.GetRenderZeroVolumeEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("sentNoiseLevel", m.GetSentNoiseLevel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("sentSignalLevel", m.GetSentSignalLevel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("speakerGlitchRate", m.GetSpeakerGlitchRate())
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
func (m *DeviceInfo) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DeviceInfo) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCaptureDeviceDriver sets the captureDeviceDriver property value. Name of the capture device driver used by the media endpoint.
func (m *DeviceInfo) SetCaptureDeviceDriver(value *string)() {
    err := m.GetBackingStore().Set("captureDeviceDriver", value)
    if err != nil {
        panic(err)
    }
}
// SetCaptureDeviceName sets the captureDeviceName property value. Name of the capture device used by the media endpoint.
func (m *DeviceInfo) SetCaptureDeviceName(value *string)() {
    err := m.GetBackingStore().Set("captureDeviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetCaptureNotFunctioningEventRatio sets the captureNotFunctioningEventRatio property value. Fraction of the call that the media endpoint detected the capture device was not working properly.
func (m *DeviceInfo) SetCaptureNotFunctioningEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("captureNotFunctioningEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetCpuInsufficentEventRatio sets the cpuInsufficentEventRatio property value. Fraction of the call that the media endpoint detected the CPU resources available were insufficient and caused poor quality of the audio sent and received.
func (m *DeviceInfo) SetCpuInsufficentEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("cpuInsufficentEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceClippingEventRatio sets the deviceClippingEventRatio property value. Fraction of the call that the media endpoint detected clipping in the captured audio that caused poor quality of the audio being sent.
func (m *DeviceInfo) SetDeviceClippingEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("deviceClippingEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceGlitchEventRatio sets the deviceGlitchEventRatio property value. Fraction of the call that the media endpoint detected glitches or gaps in the audio played or captured that caused poor quality of the audio being sent or received.
func (m *DeviceInfo) SetDeviceGlitchEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("deviceGlitchEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetHowlingEventCount sets the howlingEventCount property value. Number of times during the call that the media endpoint detected howling or screeching audio.
func (m *DeviceInfo) SetHowlingEventCount(value *int32)() {
    err := m.GetBackingStore().Set("howlingEventCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInitialSignalLevelRootMeanSquare sets the initialSignalLevelRootMeanSquare property value. The root mean square (RMS) of the incoming signal of up to the first 30 seconds of the call.
func (m *DeviceInfo) SetInitialSignalLevelRootMeanSquare(value *float32)() {
    err := m.GetBackingStore().Set("initialSignalLevelRootMeanSquare", value)
    if err != nil {
        panic(err)
    }
}
// SetLowSpeechLevelEventRatio sets the lowSpeechLevelEventRatio property value. Fraction of the call that the media endpoint detected low speech level that caused poor quality of the audio being sent.
func (m *DeviceInfo) SetLowSpeechLevelEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("lowSpeechLevelEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetLowSpeechToNoiseEventRatio sets the lowSpeechToNoiseEventRatio property value. Fraction of the call that the media endpoint detected low speech to noise level that caused poor quality of the audio being sent.
func (m *DeviceInfo) SetLowSpeechToNoiseEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("lowSpeechToNoiseEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetMicGlitchRate sets the micGlitchRate property value. Glitches per 5 minute interval for the media endpoint's microphone.
func (m *DeviceInfo) SetMicGlitchRate(value *float32)() {
    err := m.GetBackingStore().Set("micGlitchRate", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DeviceInfo) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetReceivedNoiseLevel sets the receivedNoiseLevel property value. Average energy level of received audio for audio classified as mono noise or left channel of stereo noise by the media endpoint.
func (m *DeviceInfo) SetReceivedNoiseLevel(value *int32)() {
    err := m.GetBackingStore().Set("receivedNoiseLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetReceivedSignalLevel sets the receivedSignalLevel property value. Average energy level of received audio for audio classified as mono speech, or left channel of stereo speech by the media endpoint.
func (m *DeviceInfo) SetReceivedSignalLevel(value *int32)() {
    err := m.GetBackingStore().Set("receivedSignalLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetRenderDeviceDriver sets the renderDeviceDriver property value. Name of the render device driver used by the media endpoint.
func (m *DeviceInfo) SetRenderDeviceDriver(value *string)() {
    err := m.GetBackingStore().Set("renderDeviceDriver", value)
    if err != nil {
        panic(err)
    }
}
// SetRenderDeviceName sets the renderDeviceName property value. Name of the render device used by the media endpoint.
func (m *DeviceInfo) SetRenderDeviceName(value *string)() {
    err := m.GetBackingStore().Set("renderDeviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetRenderMuteEventRatio sets the renderMuteEventRatio property value. Fraction of the call that media endpoint detected device render is muted.
func (m *DeviceInfo) SetRenderMuteEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("renderMuteEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetRenderNotFunctioningEventRatio sets the renderNotFunctioningEventRatio property value. Fraction of the call that the media endpoint detected the render device was not working properly.
func (m *DeviceInfo) SetRenderNotFunctioningEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("renderNotFunctioningEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetRenderZeroVolumeEventRatio sets the renderZeroVolumeEventRatio property value. Fraction of the call that media endpoint detected device render volume is set to 0.
func (m *DeviceInfo) SetRenderZeroVolumeEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("renderZeroVolumeEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetSentNoiseLevel sets the sentNoiseLevel property value. Average energy level of sent audio for audio classified as mono noise or left channel of stereo noise by the media endpoint.
func (m *DeviceInfo) SetSentNoiseLevel(value *int32)() {
    err := m.GetBackingStore().Set("sentNoiseLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetSentSignalLevel sets the sentSignalLevel property value. Average energy level of sent audio for audio classified as mono speech, or left channel of stereo speech by the media endpoint.
func (m *DeviceInfo) SetSentSignalLevel(value *int32)() {
    err := m.GetBackingStore().Set("sentSignalLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetSpeakerGlitchRate sets the speakerGlitchRate property value. Glitches per 5 minute internal for the media endpoint's loudspeaker.
func (m *DeviceInfo) SetSpeakerGlitchRate(value *float32)() {
    err := m.GetBackingStore().Set("speakerGlitchRate", value)
    if err != nil {
        panic(err)
    }
}
type DeviceInfoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCaptureDeviceDriver()(*string)
    GetCaptureDeviceName()(*string)
    GetCaptureNotFunctioningEventRatio()(*float32)
    GetCpuInsufficentEventRatio()(*float32)
    GetDeviceClippingEventRatio()(*float32)
    GetDeviceGlitchEventRatio()(*float32)
    GetHowlingEventCount()(*int32)
    GetInitialSignalLevelRootMeanSquare()(*float32)
    GetLowSpeechLevelEventRatio()(*float32)
    GetLowSpeechToNoiseEventRatio()(*float32)
    GetMicGlitchRate()(*float32)
    GetOdataType()(*string)
    GetReceivedNoiseLevel()(*int32)
    GetReceivedSignalLevel()(*int32)
    GetRenderDeviceDriver()(*string)
    GetRenderDeviceName()(*string)
    GetRenderMuteEventRatio()(*float32)
    GetRenderNotFunctioningEventRatio()(*float32)
    GetRenderZeroVolumeEventRatio()(*float32)
    GetSentNoiseLevel()(*int32)
    GetSentSignalLevel()(*int32)
    GetSpeakerGlitchRate()(*float32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCaptureDeviceDriver(value *string)()
    SetCaptureDeviceName(value *string)()
    SetCaptureNotFunctioningEventRatio(value *float32)()
    SetCpuInsufficentEventRatio(value *float32)()
    SetDeviceClippingEventRatio(value *float32)()
    SetDeviceGlitchEventRatio(value *float32)()
    SetHowlingEventCount(value *int32)()
    SetInitialSignalLevelRootMeanSquare(value *float32)()
    SetLowSpeechLevelEventRatio(value *float32)()
    SetLowSpeechToNoiseEventRatio(value *float32)()
    SetMicGlitchRate(value *float32)()
    SetOdataType(value *string)()
    SetReceivedNoiseLevel(value *int32)()
    SetReceivedSignalLevel(value *int32)()
    SetRenderDeviceDriver(value *string)()
    SetRenderDeviceName(value *string)()
    SetRenderMuteEventRatio(value *float32)()
    SetRenderNotFunctioningEventRatio(value *float32)()
    SetRenderZeroVolumeEventRatio(value *float32)()
    SetSentNoiseLevel(value *int32)()
    SetSentSignalLevel(value *int32)()
    SetSpeakerGlitchRate(value *float32)()
}
