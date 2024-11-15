package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Video struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewVideo instantiates a new Video and sets the default values.
func NewVideo()(*Video) {
    m := &Video{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateVideoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVideoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVideo(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Video) GetAdditionalData()(map[string]any) {
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
// GetAudioBitsPerSample gets the audioBitsPerSample property value. Number of audio bits per sample.
// returns a *int32 when successful
func (m *Video) GetAudioBitsPerSample()(*int32) {
    val, err := m.GetBackingStore().Get("audioBitsPerSample")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAudioChannels gets the audioChannels property value. Number of audio channels.
// returns a *int32 when successful
func (m *Video) GetAudioChannels()(*int32) {
    val, err := m.GetBackingStore().Get("audioChannels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAudioFormat gets the audioFormat property value. Name of the audio format (AAC, MP3, etc.).
// returns a *string when successful
func (m *Video) GetAudioFormat()(*string) {
    val, err := m.GetBackingStore().Get("audioFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAudioSamplesPerSecond gets the audioSamplesPerSecond property value. Number of audio samples per second.
// returns a *int32 when successful
func (m *Video) GetAudioSamplesPerSecond()(*int32) {
    val, err := m.GetBackingStore().Get("audioSamplesPerSecond")
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
func (m *Video) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBitrate gets the bitrate property value. Bit rate of the video in bits per second.
// returns a *int32 when successful
func (m *Video) GetBitrate()(*int32) {
    val, err := m.GetBackingStore().Get("bitrate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDuration gets the duration property value. Duration of the file in milliseconds.
// returns a *int64 when successful
func (m *Video) GetDuration()(*int64) {
    val, err := m.GetBackingStore().Get("duration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Video) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["audioBitsPerSample"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudioBitsPerSample(val)
        }
        return nil
    }
    res["audioChannels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudioChannels(val)
        }
        return nil
    }
    res["audioFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudioFormat(val)
        }
        return nil
    }
    res["audioSamplesPerSecond"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudioSamplesPerSecond(val)
        }
        return nil
    }
    res["bitrate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitrate(val)
        }
        return nil
    }
    res["duration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDuration(val)
        }
        return nil
    }
    res["fourCC"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFourCC(val)
        }
        return nil
    }
    res["frameRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFrameRate(val)
        }
        return nil
    }
    res["height"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHeight(val)
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
    res["width"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWidth(val)
        }
        return nil
    }
    return res
}
// GetFourCC gets the fourCC property value. 'Four character code' name of the video format.
// returns a *string when successful
func (m *Video) GetFourCC()(*string) {
    val, err := m.GetBackingStore().Get("fourCC")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFrameRate gets the frameRate property value. Frame rate of the video.
// returns a *float64 when successful
func (m *Video) GetFrameRate()(*float64) {
    val, err := m.GetBackingStore().Get("frameRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetHeight gets the height property value. Height of the video, in pixels.
// returns a *int32 when successful
func (m *Video) GetHeight()(*int32) {
    val, err := m.GetBackingStore().Get("height")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Video) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWidth gets the width property value. Width of the video, in pixels.
// returns a *int32 when successful
func (m *Video) GetWidth()(*int32) {
    val, err := m.GetBackingStore().Get("width")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Video) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("audioBitsPerSample", m.GetAudioBitsPerSample())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("audioChannels", m.GetAudioChannels())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("audioFormat", m.GetAudioFormat())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("audioSamplesPerSecond", m.GetAudioSamplesPerSecond())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("bitrate", m.GetBitrate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("duration", m.GetDuration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("fourCC", m.GetFourCC())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("frameRate", m.GetFrameRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("height", m.GetHeight())
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
        err := writer.WriteInt32Value("width", m.GetWidth())
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
func (m *Video) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAudioBitsPerSample sets the audioBitsPerSample property value. Number of audio bits per sample.
func (m *Video) SetAudioBitsPerSample(value *int32)() {
    err := m.GetBackingStore().Set("audioBitsPerSample", value)
    if err != nil {
        panic(err)
    }
}
// SetAudioChannels sets the audioChannels property value. Number of audio channels.
func (m *Video) SetAudioChannels(value *int32)() {
    err := m.GetBackingStore().Set("audioChannels", value)
    if err != nil {
        panic(err)
    }
}
// SetAudioFormat sets the audioFormat property value. Name of the audio format (AAC, MP3, etc.).
func (m *Video) SetAudioFormat(value *string)() {
    err := m.GetBackingStore().Set("audioFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetAudioSamplesPerSecond sets the audioSamplesPerSecond property value. Number of audio samples per second.
func (m *Video) SetAudioSamplesPerSecond(value *int32)() {
    err := m.GetBackingStore().Set("audioSamplesPerSecond", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Video) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBitrate sets the bitrate property value. Bit rate of the video in bits per second.
func (m *Video) SetBitrate(value *int32)() {
    err := m.GetBackingStore().Set("bitrate", value)
    if err != nil {
        panic(err)
    }
}
// SetDuration sets the duration property value. Duration of the file in milliseconds.
func (m *Video) SetDuration(value *int64)() {
    err := m.GetBackingStore().Set("duration", value)
    if err != nil {
        panic(err)
    }
}
// SetFourCC sets the fourCC property value. 'Four character code' name of the video format.
func (m *Video) SetFourCC(value *string)() {
    err := m.GetBackingStore().Set("fourCC", value)
    if err != nil {
        panic(err)
    }
}
// SetFrameRate sets the frameRate property value. Frame rate of the video.
func (m *Video) SetFrameRate(value *float64)() {
    err := m.GetBackingStore().Set("frameRate", value)
    if err != nil {
        panic(err)
    }
}
// SetHeight sets the height property value. Height of the video, in pixels.
func (m *Video) SetHeight(value *int32)() {
    err := m.GetBackingStore().Set("height", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Video) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetWidth sets the width property value. Width of the video, in pixels.
func (m *Video) SetWidth(value *int32)() {
    err := m.GetBackingStore().Set("width", value)
    if err != nil {
        panic(err)
    }
}
type Videoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAudioBitsPerSample()(*int32)
    GetAudioChannels()(*int32)
    GetAudioFormat()(*string)
    GetAudioSamplesPerSecond()(*int32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBitrate()(*int32)
    GetDuration()(*int64)
    GetFourCC()(*string)
    GetFrameRate()(*float64)
    GetHeight()(*int32)
    GetOdataType()(*string)
    GetWidth()(*int32)
    SetAudioBitsPerSample(value *int32)()
    SetAudioChannels(value *int32)()
    SetAudioFormat(value *string)()
    SetAudioSamplesPerSecond(value *int32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBitrate(value *int32)()
    SetDuration(value *int64)()
    SetFourCC(value *string)()
    SetFrameRate(value *float64)()
    SetHeight(value *int32)()
    SetOdataType(value *string)()
    SetWidth(value *int32)()
}
