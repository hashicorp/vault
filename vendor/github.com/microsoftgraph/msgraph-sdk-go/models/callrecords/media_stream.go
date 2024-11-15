package callrecords

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MediaStream struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMediaStream instantiates a new MediaStream and sets the default values.
func NewMediaStream()(*MediaStream) {
    m := &MediaStream{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMediaStreamFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMediaStreamFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMediaStream(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MediaStream) GetAdditionalData()(map[string]any) {
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
// GetAudioCodec gets the audioCodec property value. Codec name used to encode audio for transmission on the network. Possible values are: unknown, invalid, cn, pcma, pcmu, amrWide, g722, g7221, g7221c, g729, multiChannelAudio, muchv2, opus, satin, satinFullband, rtAudio8, rtAudio16, silk, silkNarrow, silkWide, siren, xmsRta, unknownFutureValue.
// returns a *AudioCodec when successful
func (m *MediaStream) GetAudioCodec()(*AudioCodec) {
    val, err := m.GetBackingStore().Get("audioCodec")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AudioCodec)
    }
    return nil
}
// GetAverageAudioDegradation gets the averageAudioDegradation property value. Average Network Mean Opinion Score degradation for stream. Represents how much the network loss and jitter has impacted the quality of received audio.
// returns a *float32 when successful
func (m *MediaStream) GetAverageAudioDegradation()(*float32) {
    val, err := m.GetBackingStore().Get("averageAudioDegradation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetAverageAudioNetworkJitter gets the averageAudioNetworkJitter property value. Average jitter for the stream computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
// returns a *ISODuration when successful
func (m *MediaStream) GetAverageAudioNetworkJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("averageAudioNetworkJitter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetAverageBandwidthEstimate gets the averageBandwidthEstimate property value. Average estimated bandwidth available between two endpoints in bits per second.
// returns a *int64 when successful
func (m *MediaStream) GetAverageBandwidthEstimate()(*int64) {
    val, err := m.GetBackingStore().Get("averageBandwidthEstimate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetAverageFreezeDuration gets the averageFreezeDuration property value. Average duration of the received freezing time in the video stream.
// returns a *ISODuration when successful
func (m *MediaStream) GetAverageFreezeDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("averageFreezeDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetAverageJitter gets the averageJitter property value. Average jitter for the stream computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
// returns a *ISODuration when successful
func (m *MediaStream) GetAverageJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("averageJitter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetAveragePacketLossRate gets the averagePacketLossRate property value. Average packet loss rate for stream.
// returns a *float32 when successful
func (m *MediaStream) GetAveragePacketLossRate()(*float32) {
    val, err := m.GetBackingStore().Get("averagePacketLossRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetAverageRatioOfConcealedSamples gets the averageRatioOfConcealedSamples property value. Ratio of the number of audio frames with samples generated by packet loss concealment to the total number of audio frames.
// returns a *float32 when successful
func (m *MediaStream) GetAverageRatioOfConcealedSamples()(*float32) {
    val, err := m.GetBackingStore().Get("averageRatioOfConcealedSamples")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetAverageReceivedFrameRate gets the averageReceivedFrameRate property value. Average frames per second received for all video streams computed over the duration of the session.
// returns a *float32 when successful
func (m *MediaStream) GetAverageReceivedFrameRate()(*float32) {
    val, err := m.GetBackingStore().Get("averageReceivedFrameRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetAverageRoundTripTime gets the averageRoundTripTime property value. Average network propagation round-trip time computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
// returns a *ISODuration when successful
func (m *MediaStream) GetAverageRoundTripTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("averageRoundTripTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetAverageVideoFrameLossPercentage gets the averageVideoFrameLossPercentage property value. Average percentage of video frames lost as displayed to the user.
// returns a *float32 when successful
func (m *MediaStream) GetAverageVideoFrameLossPercentage()(*float32) {
    val, err := m.GetBackingStore().Get("averageVideoFrameLossPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetAverageVideoFrameRate gets the averageVideoFrameRate property value. Average frames per second received for a video stream, computed over the duration of the session.
// returns a *float32 when successful
func (m *MediaStream) GetAverageVideoFrameRate()(*float32) {
    val, err := m.GetBackingStore().Get("averageVideoFrameRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetAverageVideoPacketLossRate gets the averageVideoPacketLossRate property value. Average fraction of packets lost, as specified in RFC 3550, computed over the duration of the session.
// returns a *float32 when successful
func (m *MediaStream) GetAverageVideoPacketLossRate()(*float32) {
    val, err := m.GetBackingStore().Get("averageVideoPacketLossRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *MediaStream) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEndDateTime gets the endDateTime property value. UTC time when the stream ended. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. This field is only available for streams that use the SIP protocol.
// returns a *Time when successful
func (m *MediaStream) GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("endDateTime")
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
func (m *MediaStream) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["audioCodec"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAudioCodec)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudioCodec(val.(*AudioCodec))
        }
        return nil
    }
    res["averageAudioDegradation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageAudioDegradation(val)
        }
        return nil
    }
    res["averageAudioNetworkJitter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageAudioNetworkJitter(val)
        }
        return nil
    }
    res["averageBandwidthEstimate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageBandwidthEstimate(val)
        }
        return nil
    }
    res["averageFreezeDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageFreezeDuration(val)
        }
        return nil
    }
    res["averageJitter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageJitter(val)
        }
        return nil
    }
    res["averagePacketLossRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAveragePacketLossRate(val)
        }
        return nil
    }
    res["averageRatioOfConcealedSamples"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageRatioOfConcealedSamples(val)
        }
        return nil
    }
    res["averageReceivedFrameRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageReceivedFrameRate(val)
        }
        return nil
    }
    res["averageRoundTripTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageRoundTripTime(val)
        }
        return nil
    }
    res["averageVideoFrameLossPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageVideoFrameLossPercentage(val)
        }
        return nil
    }
    res["averageVideoFrameRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageVideoFrameRate(val)
        }
        return nil
    }
    res["averageVideoPacketLossRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageVideoPacketLossRate(val)
        }
        return nil
    }
    res["endDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndDateTime(val)
        }
        return nil
    }
    res["isAudioForwardErrorCorrectionUsed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAudioForwardErrorCorrectionUsed(val)
        }
        return nil
    }
    res["lowFrameRateRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLowFrameRateRatio(val)
        }
        return nil
    }
    res["lowVideoProcessingCapabilityRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLowVideoProcessingCapabilityRatio(val)
        }
        return nil
    }
    res["maxAudioNetworkJitter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxAudioNetworkJitter(val)
        }
        return nil
    }
    res["maxJitter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxJitter(val)
        }
        return nil
    }
    res["maxPacketLossRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxPacketLossRate(val)
        }
        return nil
    }
    res["maxRatioOfConcealedSamples"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxRatioOfConcealedSamples(val)
        }
        return nil
    }
    res["maxRoundTripTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxRoundTripTime(val)
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
    res["packetUtilization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPacketUtilization(val)
        }
        return nil
    }
    res["postForwardErrorCorrectionPacketLossRate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPostForwardErrorCorrectionPacketLossRate(val)
        }
        return nil
    }
    res["rmsFreezeDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRmsFreezeDuration(val)
        }
        return nil
    }
    res["startDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDateTime(val)
        }
        return nil
    }
    res["streamDirection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMediaStreamDirection)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStreamDirection(val.(*MediaStreamDirection))
        }
        return nil
    }
    res["streamId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStreamId(val)
        }
        return nil
    }
    res["videoCodec"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVideoCodec)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVideoCodec(val.(*VideoCodec))
        }
        return nil
    }
    res["wasMediaBypassed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWasMediaBypassed(val)
        }
        return nil
    }
    return res
}
// GetIsAudioForwardErrorCorrectionUsed gets the isAudioForwardErrorCorrectionUsed property value. Indicates whether the forward error correction (FEC) was used at some point during the session. The default value is null.
// returns a *bool when successful
func (m *MediaStream) GetIsAudioForwardErrorCorrectionUsed()(*bool) {
    val, err := m.GetBackingStore().Get("isAudioForwardErrorCorrectionUsed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLowFrameRateRatio gets the lowFrameRateRatio property value. Fraction of the call where frame rate is less than 7.5 frames per second.
// returns a *float32 when successful
func (m *MediaStream) GetLowFrameRateRatio()(*float32) {
    val, err := m.GetBackingStore().Get("lowFrameRateRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetLowVideoProcessingCapabilityRatio gets the lowVideoProcessingCapabilityRatio property value. Fraction of the call that the client is running less than 70% expected video processing capability.
// returns a *float32 when successful
func (m *MediaStream) GetLowVideoProcessingCapabilityRatio()(*float32) {
    val, err := m.GetBackingStore().Get("lowVideoProcessingCapabilityRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetMaxAudioNetworkJitter gets the maxAudioNetworkJitter property value. Maximum of audio network jitter computed over each of the 20 second windows during the session, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
// returns a *ISODuration when successful
func (m *MediaStream) GetMaxAudioNetworkJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("maxAudioNetworkJitter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetMaxJitter gets the maxJitter property value. Maximum jitter for the stream computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
// returns a *ISODuration when successful
func (m *MediaStream) GetMaxJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("maxJitter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetMaxPacketLossRate gets the maxPacketLossRate property value. Maximum packet loss rate for the stream.
// returns a *float32 when successful
func (m *MediaStream) GetMaxPacketLossRate()(*float32) {
    val, err := m.GetBackingStore().Get("maxPacketLossRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetMaxRatioOfConcealedSamples gets the maxRatioOfConcealedSamples property value. Maximum ratio of packets concealed by the healer.
// returns a *float32 when successful
func (m *MediaStream) GetMaxRatioOfConcealedSamples()(*float32) {
    val, err := m.GetBackingStore().Get("maxRatioOfConcealedSamples")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetMaxRoundTripTime gets the maxRoundTripTime property value. Maximum network propagation round-trip time computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
// returns a *ISODuration when successful
func (m *MediaStream) GetMaxRoundTripTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("maxRoundTripTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MediaStream) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPacketUtilization gets the packetUtilization property value. Packet count for the stream.
// returns a *int64 when successful
func (m *MediaStream) GetPacketUtilization()(*int64) {
    val, err := m.GetBackingStore().Get("packetUtilization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetPostForwardErrorCorrectionPacketLossRate gets the postForwardErrorCorrectionPacketLossRate property value. Packet loss rate after FEC has been applied aggregated across all video streams and codecs.
// returns a *float32 when successful
func (m *MediaStream) GetPostForwardErrorCorrectionPacketLossRate()(*float32) {
    val, err := m.GetBackingStore().Get("postForwardErrorCorrectionPacketLossRate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetRmsFreezeDuration gets the rmsFreezeDuration property value. Average duration of the received freezing time in the video stream represented in root mean square.
// returns a *ISODuration when successful
func (m *MediaStream) GetRmsFreezeDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("rmsFreezeDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. UTC time when the stream started. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. This field is only available for streams that use the SIP protocol.
// returns a *Time when successful
func (m *MediaStream) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetStreamDirection gets the streamDirection property value. The streamDirection property
// returns a *MediaStreamDirection when successful
func (m *MediaStream) GetStreamDirection()(*MediaStreamDirection) {
    val, err := m.GetBackingStore().Get("streamDirection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MediaStreamDirection)
    }
    return nil
}
// GetStreamId gets the streamId property value. Unique identifier for the stream.
// returns a *string when successful
func (m *MediaStream) GetStreamId()(*string) {
    val, err := m.GetBackingStore().Get("streamId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVideoCodec gets the videoCodec property value. Codec name used to encode video for transmission on the network. Possible values are: unknown, invalid, av1, h263, h264, h264s, h264uc, h265, rtvc1, rtVideo, xrtvc1, unknownFutureValue.
// returns a *VideoCodec when successful
func (m *MediaStream) GetVideoCodec()(*VideoCodec) {
    val, err := m.GetBackingStore().Get("videoCodec")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VideoCodec)
    }
    return nil
}
// GetWasMediaBypassed gets the wasMediaBypassed property value. True if the media stream bypassed the Mediation Server and went straight between client and PSTN Gateway/PBX, false otherwise.
// returns a *bool when successful
func (m *MediaStream) GetWasMediaBypassed()(*bool) {
    val, err := m.GetBackingStore().Get("wasMediaBypassed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MediaStream) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAudioCodec() != nil {
        cast := (*m.GetAudioCodec()).String()
        err := writer.WriteStringValue("audioCodec", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averageAudioDegradation", m.GetAverageAudioDegradation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("averageAudioNetworkJitter", m.GetAverageAudioNetworkJitter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("averageBandwidthEstimate", m.GetAverageBandwidthEstimate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("averageFreezeDuration", m.GetAverageFreezeDuration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("averageJitter", m.GetAverageJitter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averagePacketLossRate", m.GetAveragePacketLossRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averageRatioOfConcealedSamples", m.GetAverageRatioOfConcealedSamples())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averageReceivedFrameRate", m.GetAverageReceivedFrameRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("averageRoundTripTime", m.GetAverageRoundTripTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averageVideoFrameLossPercentage", m.GetAverageVideoFrameLossPercentage())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averageVideoFrameRate", m.GetAverageVideoFrameRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("averageVideoPacketLossRate", m.GetAverageVideoPacketLossRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isAudioForwardErrorCorrectionUsed", m.GetIsAudioForwardErrorCorrectionUsed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("lowFrameRateRatio", m.GetLowFrameRateRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("lowVideoProcessingCapabilityRatio", m.GetLowVideoProcessingCapabilityRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("maxAudioNetworkJitter", m.GetMaxAudioNetworkJitter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("maxJitter", m.GetMaxJitter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("maxPacketLossRate", m.GetMaxPacketLossRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("maxRatioOfConcealedSamples", m.GetMaxRatioOfConcealedSamples())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("maxRoundTripTime", m.GetMaxRoundTripTime())
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
        err := writer.WriteInt64Value("packetUtilization", m.GetPacketUtilization())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("postForwardErrorCorrectionPacketLossRate", m.GetPostForwardErrorCorrectionPacketLossRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("rmsFreezeDuration", m.GetRmsFreezeDuration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetStreamDirection() != nil {
        cast := (*m.GetStreamDirection()).String()
        err := writer.WriteStringValue("streamDirection", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("streamId", m.GetStreamId())
        if err != nil {
            return err
        }
    }
    if m.GetVideoCodec() != nil {
        cast := (*m.GetVideoCodec()).String()
        err := writer.WriteStringValue("videoCodec", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("wasMediaBypassed", m.GetWasMediaBypassed())
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
func (m *MediaStream) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAudioCodec sets the audioCodec property value. Codec name used to encode audio for transmission on the network. Possible values are: unknown, invalid, cn, pcma, pcmu, amrWide, g722, g7221, g7221c, g729, multiChannelAudio, muchv2, opus, satin, satinFullband, rtAudio8, rtAudio16, silk, silkNarrow, silkWide, siren, xmsRta, unknownFutureValue.
func (m *MediaStream) SetAudioCodec(value *AudioCodec)() {
    err := m.GetBackingStore().Set("audioCodec", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageAudioDegradation sets the averageAudioDegradation property value. Average Network Mean Opinion Score degradation for stream. Represents how much the network loss and jitter has impacted the quality of received audio.
func (m *MediaStream) SetAverageAudioDegradation(value *float32)() {
    err := m.GetBackingStore().Set("averageAudioDegradation", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageAudioNetworkJitter sets the averageAudioNetworkJitter property value. Average jitter for the stream computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
func (m *MediaStream) SetAverageAudioNetworkJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("averageAudioNetworkJitter", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageBandwidthEstimate sets the averageBandwidthEstimate property value. Average estimated bandwidth available between two endpoints in bits per second.
func (m *MediaStream) SetAverageBandwidthEstimate(value *int64)() {
    err := m.GetBackingStore().Set("averageBandwidthEstimate", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageFreezeDuration sets the averageFreezeDuration property value. Average duration of the received freezing time in the video stream.
func (m *MediaStream) SetAverageFreezeDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("averageFreezeDuration", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageJitter sets the averageJitter property value. Average jitter for the stream computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
func (m *MediaStream) SetAverageJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("averageJitter", value)
    if err != nil {
        panic(err)
    }
}
// SetAveragePacketLossRate sets the averagePacketLossRate property value. Average packet loss rate for stream.
func (m *MediaStream) SetAveragePacketLossRate(value *float32)() {
    err := m.GetBackingStore().Set("averagePacketLossRate", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageRatioOfConcealedSamples sets the averageRatioOfConcealedSamples property value. Ratio of the number of audio frames with samples generated by packet loss concealment to the total number of audio frames.
func (m *MediaStream) SetAverageRatioOfConcealedSamples(value *float32)() {
    err := m.GetBackingStore().Set("averageRatioOfConcealedSamples", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageReceivedFrameRate sets the averageReceivedFrameRate property value. Average frames per second received for all video streams computed over the duration of the session.
func (m *MediaStream) SetAverageReceivedFrameRate(value *float32)() {
    err := m.GetBackingStore().Set("averageReceivedFrameRate", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageRoundTripTime sets the averageRoundTripTime property value. Average network propagation round-trip time computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
func (m *MediaStream) SetAverageRoundTripTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("averageRoundTripTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageVideoFrameLossPercentage sets the averageVideoFrameLossPercentage property value. Average percentage of video frames lost as displayed to the user.
func (m *MediaStream) SetAverageVideoFrameLossPercentage(value *float32)() {
    err := m.GetBackingStore().Set("averageVideoFrameLossPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageVideoFrameRate sets the averageVideoFrameRate property value. Average frames per second received for a video stream, computed over the duration of the session.
func (m *MediaStream) SetAverageVideoFrameRate(value *float32)() {
    err := m.GetBackingStore().Set("averageVideoFrameRate", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageVideoPacketLossRate sets the averageVideoPacketLossRate property value. Average fraction of packets lost, as specified in RFC 3550, computed over the duration of the session.
func (m *MediaStream) SetAverageVideoPacketLossRate(value *float32)() {
    err := m.GetBackingStore().Set("averageVideoPacketLossRate", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MediaStream) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEndDateTime sets the endDateTime property value. UTC time when the stream ended. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. This field is only available for streams that use the SIP protocol.
func (m *MediaStream) SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("endDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAudioForwardErrorCorrectionUsed sets the isAudioForwardErrorCorrectionUsed property value. Indicates whether the forward error correction (FEC) was used at some point during the session. The default value is null.
func (m *MediaStream) SetIsAudioForwardErrorCorrectionUsed(value *bool)() {
    err := m.GetBackingStore().Set("isAudioForwardErrorCorrectionUsed", value)
    if err != nil {
        panic(err)
    }
}
// SetLowFrameRateRatio sets the lowFrameRateRatio property value. Fraction of the call where frame rate is less than 7.5 frames per second.
func (m *MediaStream) SetLowFrameRateRatio(value *float32)() {
    err := m.GetBackingStore().Set("lowFrameRateRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetLowVideoProcessingCapabilityRatio sets the lowVideoProcessingCapabilityRatio property value. Fraction of the call that the client is running less than 70% expected video processing capability.
func (m *MediaStream) SetLowVideoProcessingCapabilityRatio(value *float32)() {
    err := m.GetBackingStore().Set("lowVideoProcessingCapabilityRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxAudioNetworkJitter sets the maxAudioNetworkJitter property value. Maximum of audio network jitter computed over each of the 20 second windows during the session, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
func (m *MediaStream) SetMaxAudioNetworkJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("maxAudioNetworkJitter", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxJitter sets the maxJitter property value. Maximum jitter for the stream computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
func (m *MediaStream) SetMaxJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("maxJitter", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxPacketLossRate sets the maxPacketLossRate property value. Maximum packet loss rate for the stream.
func (m *MediaStream) SetMaxPacketLossRate(value *float32)() {
    err := m.GetBackingStore().Set("maxPacketLossRate", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxRatioOfConcealedSamples sets the maxRatioOfConcealedSamples property value. Maximum ratio of packets concealed by the healer.
func (m *MediaStream) SetMaxRatioOfConcealedSamples(value *float32)() {
    err := m.GetBackingStore().Set("maxRatioOfConcealedSamples", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxRoundTripTime sets the maxRoundTripTime property value. Maximum network propagation round-trip time computed as specified in RFC 3550, denoted in ISO 8601 format. For example, 1 second is denoted as 'PT1S', where 'P' is the duration designator, 'T' is the time designator, and 'S' is the second designator.
func (m *MediaStream) SetMaxRoundTripTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("maxRoundTripTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MediaStream) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPacketUtilization sets the packetUtilization property value. Packet count for the stream.
func (m *MediaStream) SetPacketUtilization(value *int64)() {
    err := m.GetBackingStore().Set("packetUtilization", value)
    if err != nil {
        panic(err)
    }
}
// SetPostForwardErrorCorrectionPacketLossRate sets the postForwardErrorCorrectionPacketLossRate property value. Packet loss rate after FEC has been applied aggregated across all video streams and codecs.
func (m *MediaStream) SetPostForwardErrorCorrectionPacketLossRate(value *float32)() {
    err := m.GetBackingStore().Set("postForwardErrorCorrectionPacketLossRate", value)
    if err != nil {
        panic(err)
    }
}
// SetRmsFreezeDuration sets the rmsFreezeDuration property value. Average duration of the received freezing time in the video stream represented in root mean square.
func (m *MediaStream) SetRmsFreezeDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("rmsFreezeDuration", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. UTC time when the stream started. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. This field is only available for streams that use the SIP protocol.
func (m *MediaStream) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetStreamDirection sets the streamDirection property value. The streamDirection property
func (m *MediaStream) SetStreamDirection(value *MediaStreamDirection)() {
    err := m.GetBackingStore().Set("streamDirection", value)
    if err != nil {
        panic(err)
    }
}
// SetStreamId sets the streamId property value. Unique identifier for the stream.
func (m *MediaStream) SetStreamId(value *string)() {
    err := m.GetBackingStore().Set("streamId", value)
    if err != nil {
        panic(err)
    }
}
// SetVideoCodec sets the videoCodec property value. Codec name used to encode video for transmission on the network. Possible values are: unknown, invalid, av1, h263, h264, h264s, h264uc, h265, rtvc1, rtVideo, xrtvc1, unknownFutureValue.
func (m *MediaStream) SetVideoCodec(value *VideoCodec)() {
    err := m.GetBackingStore().Set("videoCodec", value)
    if err != nil {
        panic(err)
    }
}
// SetWasMediaBypassed sets the wasMediaBypassed property value. True if the media stream bypassed the Mediation Server and went straight between client and PSTN Gateway/PBX, false otherwise.
func (m *MediaStream) SetWasMediaBypassed(value *bool)() {
    err := m.GetBackingStore().Set("wasMediaBypassed", value)
    if err != nil {
        panic(err)
    }
}
type MediaStreamable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAudioCodec()(*AudioCodec)
    GetAverageAudioDegradation()(*float32)
    GetAverageAudioNetworkJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetAverageBandwidthEstimate()(*int64)
    GetAverageFreezeDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetAverageJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetAveragePacketLossRate()(*float32)
    GetAverageRatioOfConcealedSamples()(*float32)
    GetAverageReceivedFrameRate()(*float32)
    GetAverageRoundTripTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetAverageVideoFrameLossPercentage()(*float32)
    GetAverageVideoFrameRate()(*float32)
    GetAverageVideoPacketLossRate()(*float32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIsAudioForwardErrorCorrectionUsed()(*bool)
    GetLowFrameRateRatio()(*float32)
    GetLowVideoProcessingCapabilityRatio()(*float32)
    GetMaxAudioNetworkJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetMaxJitter()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetMaxPacketLossRate()(*float32)
    GetMaxRatioOfConcealedSamples()(*float32)
    GetMaxRoundTripTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetOdataType()(*string)
    GetPacketUtilization()(*int64)
    GetPostForwardErrorCorrectionPacketLossRate()(*float32)
    GetRmsFreezeDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetStreamDirection()(*MediaStreamDirection)
    GetStreamId()(*string)
    GetVideoCodec()(*VideoCodec)
    GetWasMediaBypassed()(*bool)
    SetAudioCodec(value *AudioCodec)()
    SetAverageAudioDegradation(value *float32)()
    SetAverageAudioNetworkJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetAverageBandwidthEstimate(value *int64)()
    SetAverageFreezeDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetAverageJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetAveragePacketLossRate(value *float32)()
    SetAverageRatioOfConcealedSamples(value *float32)()
    SetAverageReceivedFrameRate(value *float32)()
    SetAverageRoundTripTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetAverageVideoFrameLossPercentage(value *float32)()
    SetAverageVideoFrameRate(value *float32)()
    SetAverageVideoPacketLossRate(value *float32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIsAudioForwardErrorCorrectionUsed(value *bool)()
    SetLowFrameRateRatio(value *float32)()
    SetLowVideoProcessingCapabilityRatio(value *float32)()
    SetMaxAudioNetworkJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetMaxJitter(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetMaxPacketLossRate(value *float32)()
    SetMaxRatioOfConcealedSamples(value *float32)()
    SetMaxRoundTripTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetOdataType(value *string)()
    SetPacketUtilization(value *int64)()
    SetPostForwardErrorCorrectionPacketLossRate(value *float32)()
    SetRmsFreezeDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetStreamDirection(value *MediaStreamDirection)()
    SetStreamId(value *string)()
    SetVideoCodec(value *VideoCodec)()
    SetWasMediaBypassed(value *bool)()
}
